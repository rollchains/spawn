package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"gitub.com/strangelove-ventures/spawn/spawn"
)

// Here is a real implementation of Greeter
type ExamplePlugin struct {
	logger hclog.Logger
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
		"greeter": &spawn.GreeterPlugin{Impl: greeter},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
