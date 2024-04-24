package plugins

import "github.com/spf13/cobra"

// SpawnPlugin is the interface that we're exposing as a plugin.
type SpawnPlugin interface {
	Cmd() *cobra.Command
	Version() string
}

var _ SpawnPlugin = &SpawnPluginBase{}

type SpawnPluginBase struct {
	cmd     *cobra.Command
	version string
}

func NewSpawnPluginBase(cmd *cobra.Command) *SpawnPluginBase {
	return &SpawnPluginBase{
		cmd:     cmd,
		version: "v0.0.1",
	}
}

func (s *SpawnPluginBase) Cmd() *cobra.Command {
	return s.cmd
}

func (s *SpawnPluginBase) Version() string {
	return s.version
}
