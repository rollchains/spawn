package main

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	plugins "gitub.com/strangelove-ventures/spawn/plugins"
	"gitub.com/strangelove-ventures/spawn/spawn"
)

// Set in the makefile ld_flags on compile
var SpawnVersion = ""

var LogLevelFlag = "log-level"

var appPlugins map[string]spawn.Greeter

func main() {
	appPlugins = loadPlugins()
	fmt.Println("ssss", appPlugins)

	rootCmd.AddCommand(newChain)
	rootCmd.AddCommand(LocalICCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(ModuleCmd())

	for name := range appPlugins {
		fmt.Println("name", name)
		rootCmd.AddCommand(PluginCmd(name))
	}

	rootCmd.PersistentFlags().String(LogLevelFlag, "info", "log level (debug, info, warn, error)")

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

func PluginCmd(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Plugin " + name,
		Run: func(cmd *cobra.Command, args []string) {
			// This should be in the plugin interface for interaction
			fmt.Println("Plugin", name)
		},
	}
}

var rootCmd = &cobra.Command{
	Use:   "spawn",
	Short: "Entry into the Interchain",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: false,
	},
	Run: func(cmd *cobra.Command, args []string) {

		// var handshakeConfig = plugin.HandshakeConfig{
		// 	ProtocolVersion:  1,
		// 	MagicCookieKey:   "BASIC_PLUGIN",
		// 	MagicCookieValue: "hello",
		// }

		// // pluginMap is the map of plugins we can dispense.
		// var pluginMap = map[string]plugin.Plugin{
		// 	"greeter": &plugins.GreeterPlugin{},
		// }

		// logger := hclog.New(&hclog.LoggerOptions{
		// 	Name: "plugin",
		// 	// Output: os.Stdout,
		// 	Level: hclog.Error,
		// })

		// // We're a host! Start by launching the plugin process.
		// client := plugin.NewClient(&plugin.ClientConfig{
		// 	HandshakeConfig: handshakeConfig,
		// 	Plugins:         pluginMap,
		// 	Cmd:             exec.Command("./plugins/greeter"), // go build -o ./plugin/greeter ./plugin/greeter_impl.go
		// 	Logger:          logger,
		// })
		// defer client.Kill()

		// // Connect via RPC
		// rpcClient, err := client.Client()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // Request the plugin
		// raw, err := rpcClient.Dispense("greeter")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// // We should have a Greeter now! This feels like a normal interface
		// // implementation but is in fact over an RPC connection.
		// greeter := raw.(plugins.Greeter)
		// fmt.Println(greeter.Greet())

		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

func loadPlugins() map[string]spawn.Greeter { // or plugin.Plugin ?
	// plugins.Plugins contains them all
	f := plugins.PluginsFS

	pairings := make(map[string]spawn.Greeter)

	fs.WalkDir(f, ".", func(relPath string, d fs.DirEntry, e error) error {
		if d.IsDir() {
			// TODO: iterate internal and have them as sub commands
			return nil
		}

		// removes '.' and any files with extensions
		if strings.Contains(relPath, ".") {
			return nil
		}

		// print relPath
		fmt.Println("relPath", relPath)
		// name, cookie := strings.Split(relPath, "-")[0], strings.Split(relPath, "-")[1]
		// fmt.Println("name", name)
		// fmt.Println("cookie", cookie)

		var handshakeConfig = plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		}

		// // pluginMap is the map of plugins we can dispense.
		var pluginMap = map[string]plugin.Plugin{
			relPath: &spawn.GreeterPlugin{},
		}

		logger := hclog.New(&hclog.LoggerOptions{
			Name: "plugin",
			// Output: os.Stdout,
			Level: hclog.Error,
		})

		// We're a host! Start by launching the plugin process.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Cmd:             exec.Command("./plugins/" + relPath),
			Logger:          logger,
		})
		defer client.Kill()

		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			log.Fatal(err)
		}

		// // Request the plugin
		raw, err := rpcClient.Dispense(relPath)
		if err != nil {
			log.Fatal(err)
		}

		// We should have a Greeter now! This feels like a normal interface
		// implementation but is in fact over an RPC connection.
		sp := raw.(spawn.Greeter)
		fmt.Println("interaction", sp.Greet())

		pairings[relPath] = sp

		return nil
	})

	return pairings
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
		log.Fatal(err)
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
