package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitub.com/strangelove-ventures/spawn/spawn"
)

func ProtoServiceGenerate() *cobra.Command {
	cmd := &cobra.Command{
		// TODO: this name is ew
		// TODO: Put this in the make file on proto-gen? (after)
		Use:     "service-generate [module]",
		Short:   "Auto generate the MsgService stubs from proto -> Cosmos-SDK",
		Example: `spawn service-generate mymodule`,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"sg"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("service-generate called")

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error: ", err)
			}

			missingRPCMethods := spawn.GetMissingRPCMethodsFromModuleProto(cwd)
			if len(missingRPCMethods) == 0 {
				fmt.Println("No missing methods")
				return
			}

			spawn.ApplyMissingRPCMethodsToGoSourceFiles(missingRPCMethods)
		},
	}

	return cmd
}
