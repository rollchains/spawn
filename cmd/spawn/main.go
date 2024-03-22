package main

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path"
	"plugin"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/rollchains/spawn/plugins"
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

func applyPluginCmds() {
	plugins := &cobra.Command{
		Use:     "plugins",
		Short:   "Manage plugins",
		Aliases: []string{"plugin", "plug", "pl"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}

	for _, plugin := range loadPlugins() {
		plugins.AddCommand(plugin.Cmd())
	}

	rootCmd.AddCommand(plugins)
}

func loadPlugins() map[string]*plugins.SpawnPluginBase {
	p := make(map[string]*plugins.SpawnPluginBase)

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

		if !strings.Contains(relPath, ".so") {
			return nil
		}

		absPath := path.Join(pluginsDir, relPath)

		// read the absolute path
		plug, err := plugin.Open(absPath)
		if err != nil {
			logger.Error(fmt.Sprintf("Error opening plugin: %v", err))
			return nil
		}

		base, err := plug.Lookup("Plugin")
		if err != nil {
			logger.Error(fmt.Sprintf("Error looking up symbol: %v", err))
			return nil
		}

		pluginInstance, ok := base.(plugins.SpawnPlugin)
		if !ok {
			logger.Error(fmt.Sprintf("Plugin %s does not implement the SpawnPlugin interface. Skipping", absPath))
			return nil
		}

		p[relPath] = plugins.NewSpawnPluginBase(pluginInstance.Cmd())

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
