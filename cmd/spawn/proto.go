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
			logger := GetLogger()

			cwd, err := os.Getwd()
			if err != nil {
				logger.Error("Error", "error", err)
			}

			missingRPCMethods, err := spawn.GetMissingRPCMethodsFromModuleProto(logger, cwd)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			hasChanges := false
			for _, v := range missingRPCMethods {
				if len(v) > 0 {
					hasChanges = true
					break
				}
			}
			if !hasChanges {
				logger.Info("No missing methods to apply")
				return
			}

			if err := spawn.ApplyMissingRPCMethodsToGoSourceFiles(logger, missingRPCMethods); err != nil {
				logger.Error("Error", "error", err)
			}
		},
	}

	return cmd
}
