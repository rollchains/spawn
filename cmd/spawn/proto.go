package main

import (
	"fmt"
	"os"
	"path"
	"strings"

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

			protoPath := path.Join(cwd, "proto")

			modules := spawn.GetModuleMapFromProto(protoPath)

			// module name -> RPC methods
			missing := make(map[string][]*spawn.ProtoRPC, 0)

			filePaths := make(map[string]map[spawn.FileType]string, 0)

			// TODO: currently if using multiple modules, it will run the code 2 times for generating missing methods
			for name, rpcMethods := range modules {
				fmt.Println("\n------------- Module: ", name)

				modulePath := path.Join(cwd, "x", name, "keeper") // hardcode for keeper is less than ideal, but will do for now

				// currentMethods := make(map[spawn.FileType][]string, 0) // tx/query -> methods
				txMethods := make([]string, 0)
				queryMethods := make([]string, 0)

				msgServerFile := ""
				queryServerFile := ""

				for _, rpc := range rpcMethods {
					rpc := rpc
					rpc.Module = name

					fmt.Println("\nService: ", rpc)

					// get files in service.Location
					files, err := os.ReadDir(modulePath)
					if err != nil {
						fmt.Println("Error: ", err)
					}

					for _, f := range files {
						if strings.HasSuffix(f.Name(), "_test.go") {
							continue
						}

						content, err := os.ReadFile(path.Join(modulePath, f.Name()))
						if err != nil {
							fmt.Println("Error: ", err)
						}

						// if the file type is not the expected, continue
						// if the content of this file is not the same as the service we are tying to use, continue
						if rpc.FType != isFileQueryOrMsgServer(content) {
							continue
						}

						fmt.Println(" = File: ", f.Name())

						switch rpc.FType {
						case spawn.Tx:
							msgServerFile = path.Join(modulePath, f.Name())
						case spawn.Query:
							queryServerFile = path.Join(modulePath, f.Name())
						default:
							fmt.Println("Error: ", "Unknown FileType")
							panic("RUT ROE RAGGY")
						}

						// find any line with `func ` in it

						lines := strings.Split(string(content), "\n")
						for _, line := range lines {
							// receiver func
							if strings.Contains(line, "func (") {
								// fmt.Println("  func: ", line)
								// currentMethods = append(currentMethods, line)

								// if the method is already in the currentMethods, skip
								// if not, add to missingMethods

								// if _, ok := currentMethods[t]; !ok {
								// 	currentMethods[t] = make([]string, 0)
								// }

								if !strings.Contains(line, rpc.Name) {
									// TODO: put missing here?
									continue
								} else {
									// if the method is already in the currentMethods, skip
									// fmt.Println("    method: ", line)
									// currentMethods[t] = append(currentMethods[t], parseReceiverMethodName(line))
									switch rpc.FType {
									case spawn.Tx:
										txMethods = append(txMethods, parseReceiverMethodName(line))
									case spawn.Query:
										queryMethods = append(queryMethods, parseReceiverMethodName(line))
									default:
										fmt.Println("Error: ", "Unknown FileType")
										panic("RUT ROE RAGGY")
									}
								}
							}
						}
					}

					filePaths[name] = map[spawn.FileType]string{
						spawn.Tx:    msgServerFile,
						spawn.Query: queryServerFile,
					}
				}

				// map[query:[Params] tx:[UpdateParams]]
				fmt.Println("\n-------- Current Tx Methods: ", txMethods)
				fmt.Println("-------- Current Query Methods: ", queryMethods)

				// print modules
				// fmt.Println("\nModules: ", modules)

				// iterate services again and apply any missing methods to the file

				for name, rpcs := range modules {
					for _, rpc := range rpcs {
						// current := currentMethods[fileType]

						var current []string
						switch rpc.FType {
						case spawn.Tx:
							current = txMethods
						case spawn.Query:
							current = queryMethods
						default:
							fmt.Println("Error: ", "Unknown FileType")
							panic("RUT ROE RAGGY")
						}

						fmt.Println("   Current: ", rpc.FType, current)

						// ft := service.FType // tx or spawn, found in currentMethods
						// if service.Name

						found := false
						for _, method := range current {
							method := method
							if method == rpc.Name {
								found = true
								break
							}
						}

						// if !found {
						// TODO: fix this hack, yuck, *spits*. phew
						// it's the wrong type, fix
						// if fileType != service.FType {
						// 	service.FType = fileType
						// 	// fmt.Println("  - Skipping: ", fileType, service.FType, service)
						// 	// continue
						// }

						if found {
							fmt.Println("  - Found: ", rpc.Name)
							continue
						}

						if _, ok := missing[name]; !ok {
							missing[name] = make([]*spawn.ProtoRPC, 0)
						}

						// MISSING METHOD
						// fmt.Println("Missing method: ", service.Name, service.Req, service.Res, service.Location)
						missing[name] = append(missing[name], rpc)
						fmt.Println("  - Missing: ", rpc.Name)
					}
				}

				fmt.Println("\n\nMissing: ")
				for name, rpcs := range missing {
					fmt.Println("  - Module: ", name)
					for _, rpc := range rpcs {
						fmt.Println("    - ", rpc.Name, rpc.Req, rpc.Res)
					}
				}
			}

			// print filePaths
			fmt.Println("\nFile Paths: ", filePaths)

			// get filePaths[miss.Module]
			fmt.Println("\n File Path Specific: ", filePaths["cnd"])
			fmt.Println("\n File Path Specific2: ", filePaths["cnd"][spawn.Query])

			for _, missed := range missing {
				for _, miss := range missed {
					miss := miss
					// get miss.FType from filePaths
					// print miss.Module and miss.FType
					fmt.Println("Module: ", miss.Module, "FType: ", miss.FType)

					p := filePaths[miss.Module][miss.FType]
					fmt.Println("File: ", p)

					content, err := os.ReadFile(p)
					if err != nil {
						panic(fmt.Sprintf("Error: %s, file: %s", err.Error(), p))
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
	panic("unimplemented")
	return &types.%s{}, nil
}
`, miss.Name, miss.Name, miss.Req, miss.Res, miss.Res)
						// append to the file content after a new line at the end
						content = append(content, []byte("\n"+code)...)
					case spawn.Query:
						fmt.Println("Append to Query")
						code := fmt.Sprintf(`// %s implements types.QueryServer.
func (k Querier) %s(goCtx context.Context, req *types.%s) (*types.%s, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	panic("unimplemented")
	return &types.%s{}, nil
}
`, miss.Name, miss.Name, miss.Req, miss.Res, miss.Res)

						// append to the file content after a new line at the end
						content = append(content, []byte("\n"+code)...)

					}

					if err := os.WriteFile(p, content, 0644); err != nil {
						fmt.Println("Error: ", err)
					}
				}

			}

		},
	}

	return cmd
}

func parseReceiverMethodName(f string) string {
	// given a string of text like `func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {`
	// parse out Params, req and the response

	name := ""

	f = strings.ReplaceAll(f, "func (", "")
	parts := strings.Split(f, ") ")
	name = strings.Split(parts[1], "(")[0]

	return strings.Trim(name, " ")
}

func isFileQueryOrMsgServer(bz []byte) spawn.FileType {
	s := strings.ToLower(string(bz))

	if strings.Contains(s, "queryserver") || strings.Contains(s, "querier") {
		return spawn.Query
	}

	if strings.Contains(s, "msgserver") || strings.Contains(s, "msgservice") {
		return spawn.Tx
	}

	return spawn.None
}
