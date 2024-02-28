package main

import (
	"github.com/spf13/cobra"
	plugins "gitub.com/strangelove-ventures/spawn/plugins"
)

var _ plugins.SpawnPlugin = &ExamplePlugin{}

type ExamplePlugin struct {
	Impl plugins.SpawnPluginBase
}

// Cmd implements plugins.SpawnPlugin.
func (e *ExamplePlugin) Cmd() func() *cobra.Command {
	return func() *cobra.Command {
		return &cobra.Command{
			Use:   "example",
			Short: "An example plugin",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println("Hello from the example plugin!")
			},
		}
	}
}

// Cmd implements plugins.SpawnPlugin.
// func (e *ExamplePlugin) Cmd() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "example",
// 		Short: "An example plugin",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			cmd.Println("Hello from the example plugin!")
// 		},
// 	}
// }

var Plugin ExamplePlugin
