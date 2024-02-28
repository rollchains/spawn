package main

import (
	"log"

	"github.com/spf13/cobra"
	plugins "gitub.com/strangelove-ventures/spawn/plugins"
)

// Make the plugin public
var Plugin SpawnMainExamplePlugin

var _ plugins.SpawnPlugin = &SpawnMainExamplePlugin{}

type SpawnMainExamplePlugin struct {
	Impl plugins.SpawnPluginBase
}

// Name implements plugins.SpawnPlugin.
func (e *SpawnMainExamplePlugin) Name() string {
	return "example"
}

// Cmd implements plugins.SpawnPlugin.
func (e *SpawnMainExamplePlugin) Cmd() func() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   e.Name(),
		Short: e.Name() + " plugin command",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}

	// add a sub command to the root command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "sub-cmd",
		Short: "An example plugin sub command",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Hello from the example plugin sub command!")
		},
	})

	return func() *cobra.Command {
		return rootCmd
	}
}
