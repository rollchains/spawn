package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Set in the makefile ld_flags on compile
var SpawnVersion = ""

func main() {
	rootCmd.AddCommand(newChain)
	rootCmd.AddCommand(LocalICCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error while executing your CLI. Err: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "spawn",
	Short: "Entry into the Interchain",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: false,
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of spawn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(SpawnVersion)
	},
}
