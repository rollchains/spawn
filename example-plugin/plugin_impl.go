package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"gitub.com/strangelove-ventures/spawn/spawn"
)

var _ spawn.SpawnPlugin = (*ExamplePlugin)(nil)

// Here is a real implementation of Greeter
type ExamplePlugin struct {
	logger hclog.Logger
}

// Cmd implements spawn.SpawnPlugin.
func (g *ExamplePlugin) Cmd() *cobra.Command {
	g.logger.Debug("message from ExamplePlugin.Cmd")
	rootCmd := cobra.Command{
		Use:     "my-cmd",
		Aliases: []string{"mc"},
		Short:   "my-cmd short description",
		Args:    cobra.NoArgs,
		Example: "my-cmd",
		Run: func(cmd *cobra.Command, args []string) {
			// This should be in the plugin interface for interaction
			fmt.Println("Plugin", "my-cmd from the plugin !!!")
		},
	}

	return &rootCmd
}

func (g *ExamplePlugin) Greet() string {
	g.logger.Debug("message from GreeterHello.Greet")
	return "Hello!"
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Error,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	greeter := &ExamplePlugin{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &spawn.SpawnPluginBase{Impl: greeter},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
