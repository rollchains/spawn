package main

import (
	"log"
	"os"
	"path"

	plugins "github.com/rollchains/spawn/plugins"
	"github.com/spf13/cobra"
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
func (e *SpawnMainExamplePlugin) Cmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   cmdName,
		Short: cmdName + " plugin command",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "touch-file [name]",
		Short: "An example plugin sub command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			fileName := args[0]

			filePath := path.Join(cwd, fileName)
			file, err := os.Create(filePath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			cmd.Printf("Created file: %s\n", filePath)
		},
	})

	return rootCmd
}
