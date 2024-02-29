package main

import (
	"log"

	"github.com/spf13/cobra"
	plugins "gitub.com/strangelove-ventures/spawn/plugins"
)

// Make the plugin public
var Plugin SpawnMainExamplePlugin

var _ plugins.SpawnPlugin = &SpawnMainExamplePlugin{}

const (
	cmdName = "example"
)

type SpawnMainExamplePlugin struct {
	Impl plugins.SpawnPluginBase
}

// Cmd implements plugins.SpawnPlugin.
func (e *SpawnMainExamplePlugin) Cmd() func() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   cmdName,
		Short: cmdName + " plugin command",
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
