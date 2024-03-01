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

			// module name -> RPC methods
			// missing := make(map[string][]*spawn.ProtoRPC, 0)
			// filePaths := make(map[string]map[spawn.FileType]string, 0)

			mm := spawn.GetCurrentRPCMethodsFromModuleProto(cwd)
			missing := mm.Missing

			// print filePaths
			// fmt.Println("\nFile Paths: ", filePaths) // these are saved for each rpc service now instead

			if len(missing) == 0 {
				fmt.Println("No missing methods")
				return
			}

			for _, missed := range missing {
				for _, miss := range missed {
					miss := miss
					// get miss.FType from filePaths
					// print miss.Module and miss.FType
					fmt.Println("Module: ", miss.Module, "FType: ", miss.FType)

					// fileLoc := filePaths[miss.Module][miss.FType]
					fileLoc := miss.FileLoc
					fmt.Println("File: ", fileLoc)

					content, err := os.ReadFile(fileLoc)
					if err != nil {
						panic(fmt.Sprintf("Error: %s, file: %s", err.Error(), fileLoc))
					}
					// fmt.Println("Content: ", string(content))

					// append to the file
					fmt.Println("Append to file: ", miss.FType, miss.Name, miss.Req, miss.Res)

					switch miss.FType {
					case spawn.Tx:
						fmt.Println("Append to Tx")
						code := fmt.Sprintf(`// %s implements types.MsgServer.
func (ms msgServer) %s(ctx context.Context, msg *types.%s) (*types.%s, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	panic("%s is unimplemented")
	return &types.%s{}, nil
}
`, miss.Name, miss.Name, miss.Req, miss.Res, miss.Name, miss.Res)
						// append to the file content after a new line at the end
						content = append(content, []byte("\n"+code)...)
					case spawn.Query:
						fmt.Println("Append to Query")
						code := fmt.Sprintf(`// %s implements types.QueryServer.
func (k Querier) %s(goCtx context.Context, req *types.%s) (*types.%s, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	panic("%s is unimplemented")
	return &types.%s{}, nil
}
`, miss.Name, miss.Name, miss.Req, miss.Res, miss.Name, miss.Res)

						// append to the file content after a new line at the end
						content = append(content, []byte("\n"+code)...)

					}

					if err := os.WriteFile(fileLoc, content, 0644); err != nil {
						fmt.Println("Error: ", err)
					}
				}

			}

		},
	}

	return cmd
}
