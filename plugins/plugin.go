package plugins

import "github.com/spf13/cobra"

// SpawnPlugin is the interface that we're exposing as a plugin.
type SpawnPlugin interface {
	Name() string
	Cmd() func() *cobra.Command
}

var _ SpawnPlugin = &SpawnPluginBase{}

type SpawnPluginBase struct {
	PluginName string
	Command    func() *cobra.Command
}

// Name implements SpawnPlugin.
func (s *SpawnPluginBase) Name() string {
	return s.PluginName
}

// Cmd implements SpawnPlugin.
func (s *SpawnPluginBase) Cmd() func() *cobra.Command {
	return s.Command
}
