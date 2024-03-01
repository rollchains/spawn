package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitub.com/strangelove-ventures/spawn/spawn"
)

func ProtoServiceGenerate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stub-gen [module (optional)]",
		Short:   "Auto generate the MsgService & Querier from proto -> Cosmos-SDK methods",
		Long:    `Auto generate the interface stubs for the types.QueryServer and types.MsgServer for your module. If no module is provided, it will do for all modules in your proto folder.`,
		Example: `spawn stub-gen [module_name]`,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{
			"stub", "stub-generate", "stub-interface", "stub-interfaces",
			"service-generate", "sg",
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := GetLogger()

			cwd, err := os.Getwd()
			if err != nil {
				logger.Error("Error", "error", err)
				return
			}

			missingRPCMethods, err := spawn.GetMissingRPCMethodsFromModuleProto(logger, cwd)
			if err != nil {
				fmt.Println("Error: ", err)
				return
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
				return
			}
		},
	}

	return cmd
}
