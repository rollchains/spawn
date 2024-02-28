package spawn

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

var _ SpawnPlugin = (*SpawnRPC)(nil)
var _ SpawnPluginRPC = (*SpawnRPCServer)(nil)

// SpawnPlugin is the interface that we're exposing as a plugin.
type SpawnPlugin interface {
	Greet() string
	Cmd() *cobra.Command
}

type SpawnPluginRPC interface {
	Greet(args interface{}, resp *string) error
	Cmd(resp *cobra.Command) error
}

// Here is an implementation that talks over RPC
type SpawnRPC struct{ client *rpc.Client }

// Cmd implements SpawnPlugin.
func (g *SpawnRPC) Cmd() *cobra.Command {
	var resp cobra.Command

	err := g.client.Call("Plugin.Cmd", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return &resp
}

func (g *SpawnRPC) Greet() string {
	var resp string
	err := g.client.Call("Plugin.Greet", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type SpawnRPCServer struct {
	// This is the real implementation
	Impl SpawnPlugin
}

// Cmd implements SpawnPluginRPC.
func (s *SpawnRPCServer) Cmd(resp *cobra.Command) error {
	*resp = *s.Impl.Cmd()
	return nil
}

func (s *SpawnRPCServer) Greet(args interface{}, resp *string) error {
	*resp = s.Impl.Greet()
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.

var _ plugin.Plugin = (*SpawnPluginBase)(nil)

type SpawnPluginBase struct {
	// Impl Injection
	Impl SpawnPlugin
}

func (p *SpawnPluginBase) Server(*plugin.MuxBroker) (interface{}, error) {
	return &SpawnRPCServer{Impl: p.Impl}, nil
}

func (SpawnPluginBase) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &SpawnRPC{client: c}, nil
}
