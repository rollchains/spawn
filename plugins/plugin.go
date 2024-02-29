package plugins

import "github.com/spf13/cobra"

// SpawnPlugin is the interface that we're exposing as a plugin.
type SpawnPlugin interface {
	Cmd() *cobra.Command
}

var _ SpawnPlugin = &SpawnPluginBase{}

type SpawnPluginBase struct {
	cmd *cobra.Command
}

func NewSpawnPluginBase(cmd *cobra.Command) *SpawnPluginBase {
	return &SpawnPluginBase{
		cmd: cmd,
	}
}

// Cmd implements SpawnPlugin.
func (s *SpawnPluginBase) Cmd() *cobra.Command {
	return s.cmd
}
