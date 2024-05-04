package main

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/rollchains/spawn/spawn"
	"github.com/spf13/cobra"
)

// Set in the makefile ld_flags on compile
var SpawnVersion = ""

var LogLevelFlag = "log-level"

func main() {

	rootCmd.AddCommand(newChain)
	rootCmd.AddCommand(LocalICCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(ModuleCmd())
	rootCmd.AddCommand(ProtoServiceGenerate())

	applyPluginCmds()

	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error while executing your CLI. Err: %v\n", err)
		os.Exit(1)
	}
}

func GetLogger() *slog.Logger {
	w := os.Stderr

	logLevel := parseLogLevelFromFlags()

	slog.SetDefault(slog.New(
		// TODO: Windows support colored logs: https://github.com/mattn/go-colorable `tint.NewHandler(colorable.NewColorable(w), nil)`
		tint.NewHandler(w, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.Kitchen,
			// Enables colors only if the terminal supports it
			NoColor: !isatty.IsTerminal(w.Fd()),
		}),
	))

	return slog.Default()
}

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

		info, err := spawn.ParseCobraCLICmd(abspath)
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

var rootCmd = &cobra.Command{
	Use:   "spawn",
	Short: "Entry into the Interchain",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: false,
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of spawn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(SpawnVersion)
	},
}

func parseLogLevelFromFlags() slog.Level {
	logLevel, err := rootCmd.PersistentFlags().GetString(LogLevelFlag)
	if err != nil {
		return slog.LevelInfo
	}

	var lvl slog.Level

	switch strings.ToLower(logLevel) {
	case "debug", "d", "dbg":
		lvl = slog.LevelDebug
	case "info", "i", "inf":
		lvl = slog.LevelInfo
	case "warn", "w", "wrn":
		lvl = slog.LevelWarn
	case "error", "e", "err", "fatal", "f", "ftl":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	return lvl
}
