package main

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	spawn "gitub.com/strangelove-ventures/spawn/spawn"
)

var _ spawn.Greeter = (*BasePlugin)(nil)

type BasePlugin struct{}

// Greet implements spawn.Greeter.
func (b *BasePlugin) Greet() string {
	panic("unimplemented")
}

// func (g *BasePlugin) Interact() string {
// 	fmt.Println("message from GreeterHello.Greet")

// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err)
// 		return "Error getting current working directory"
// 	}

// 	f, err := os.Create(path.Join(cwd, "example.txt"))
// 	if err != nil {
// 		fmt.Println(err)
// 		return "Error creating file"
// 	}
// 	defer f.Close()

// 	return fmt.Sprintf("File created: %s", f.Name())
// }

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
	// logger := hclog.New(&hclog.LoggerOptions{
	// 	Level:      hclog.Trace,
	// 	Output:     os.Stderr,
	// 	JSONFormat: true,
	// })

	ep := &BasePlugin{
		// logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"example": &spawn.GreeterPlugin{Impl: ep},
	}

	// logger.Debug("message from plugin", "foo", "bar")
	fmt.Println("message from plugin")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
