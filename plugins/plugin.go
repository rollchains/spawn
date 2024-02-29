package plugins

import "github.com/spf13/cobra"

// SpawnPlugin is the interface that we're exposing as a plugin.
type SpawnPlugin interface {
	Cmd() func() *cobra.Command
}

var _ SpawnPlugin = &SpawnPluginBase{}

type SpawnPluginBase struct {
	Command func() *cobra.Command
}

// Cmd implements SpawnPlugin.
func (s *SpawnPluginBase) Cmd() func() *cobra.Command {
	return s.Command
}
