package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

var PluginsCmd = &cobra.Command{
	Use:     "plugins",
	Short:   "Spawn Plugins",
	Aliases: []string{"plugin", "plug", "pl"},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

func applyPluginCmds() {
	for name, abspath := range loadPlugins() {
		name := name
		abspath := abspath

		info, err := ParseCobraCLICmd(abspath)
		if err != nil {
			GetLogger().Warn("error parsing the CLI commands from the plugin", "name", name, "error", err)
			continue
		}

		execCmd := &cobra.Command{
			Use:   name,
			Short: info.Description,
			Run: func(cmd *cobra.Command, args []string) {
				output, err := exec.Command(abspath, args...).CombinedOutput()
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println(string(output))
			},
		}
		PluginsCmd.AddCommand(execCmd)
	}

	rootCmd.AddCommand(PluginsCmd)
}

// returns name and path
func loadPlugins() map[string]string {
	p := make(map[string]string)

	logger := GetLogger()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	pluginsDir := path.Join(homeDir, ".spawn", "plugins")

	d := os.DirFS(pluginsDir)
	if _, err := d.Open("."); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(pluginsDir, 0755); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	err = fs.WalkDir(d, ".", func(relPath string, d fs.DirEntry, e error) error {
		if d.IsDir() {
			return nil
		}

		// /home/username/.spawn/plugins/myplugin
		absPath := path.Join(pluginsDir, relPath)

		// ensure path exist
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			logger.Error(fmt.Sprintf("Plugin %s does not exist. Skipping", absPath))
			return nil
		}

		name := path.Base(absPath)
		p[name] = absPath
		return nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error walking the path %s: %v", pluginsDir, err))
		panic(err)
	}

	return p
}
