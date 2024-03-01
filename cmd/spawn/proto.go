package main

import (
	"fmt"
	"io/fs"
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

			dirs, err := os.ReadDir(protoPath)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			absDirs := make([]string, 0)
			for _, dir := range dirs {
				if !dir.IsDir() {
					continue
				}

				if len(args) > 0 && dir.Name() != args[0] {
					continue
				}

				absDirs = append(absDirs, path.Join(protoPath, dir.Name()))
			}

			fmt.Println("Found dirs: ", absDirs)

			// walk it

			modules := make(map[string]map[spawn.FileType][]spawn.ProtoService)

			fs.WalkDir(os.DirFS(protoPath), ".", func(relPath string, d fs.DirEntry, e error) error {
				if !strings.HasSuffix(relPath, ".proto") {
					return nil
				}

				// read file content
				content, err := os.ReadFile(path.Join(protoPath, relPath))
				if err != nil {
					fmt.Println("Error: ", err)
				}

				fileType := spawn.SortContentToFileType(content)

				parent := path.Dir(relPath)
				parent = strings.Split(parent, "/")[0]

				// add/append to modules
				if _, ok := modules[parent]; !ok {
					modules[parent] = make(map[spawn.FileType][]spawn.ProtoService)
				}

				goPkgDir := spawn.GetGoPackageLocationOfFiles(content)

				switch fileType {
				case spawn.Tx:
					fmt.Println("File is a transaction")
					tx := spawn.ProtoServiceParser(content, goPkgDir)
					modules[parent][spawn.Tx] = append(modules[parent][spawn.Tx], tx...)

				case spawn.Query:
					fmt.Println("File is a query")
					query := spawn.ProtoServiceParser(content, goPkgDir)
					modules[parent][spawn.Query] = append(modules[parent][spawn.Query], query...)
				case spawn.None:
					fmt.Println("File is neither a transaction nor a query")
				}

				return nil
			})

			fmt.Printf("Modules: %+v\n", modules)

			// TODO:
			// - Go find types.MsgServer in the module (may also need to parse this data from the proto file & save to the map)
			// - Find the Name of the service, then overwrite the Req and Res types
			// tx:[{Name:UpdateParams Req:MsgUpdateParams Res:MsgUpdateParamsResponse}]] -> `func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {`
			// do for querier too.

			// iterate overall, find files containing the proto messages, and then set new stubs automatically
			// map[module][spawn.FileType]string
			filePaths := make(map[string]map[spawn.FileType]string, 0)

			for name, module := range modules {
				fmt.Println("\n------------- Module: ", name)

				modulePath := path.Join(cwd, "x", name, "keeper") // hardcode for keeper is less than ideal, but will do for now

				currentMethods := make(map[spawn.FileType][]string, 0) // tx/query -> methods

				msgServerFile := ""
				queryServerFile := ""

				for fileType, services := range module {
					for idx, service := range services {
						fmt.Println("\nService: ", service)

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

							// This limits so we only check if the service is a type and also the file is the same type
							t := isFileQueryOrMsgServer(content)
							service.FType = t // set the file type for future iteration to set missing methods
							service.Module = name
							services[idx] = service
							module[fileType] = services
							modules[name] = module

							switch t {
							case "tx":
								msgServerFile = path.Join(modulePath, f.Name())
								if fileType != spawn.Tx {
									continue
								}
							case "query":
								queryServerFile = path.Join(modulePath, f.Name())
								if fileType != spawn.Query {
									continue
								}
							case "none":
								continue
							}

							fmt.Println(" = File: ", f.Name(), t)

							// find any line with `func ` in it

							lines := strings.Split(string(content), "\n")
							for _, line := range lines {
								// receiver func
								if strings.Contains(line, "func (") {
									// fmt.Println("  func: ", line)
									// currentMethods = append(currentMethods, line)

									// if the method is already in the currentMethods, skip
									// if not, add to missingMethods

									if _, ok := currentMethods[t]; !ok {
										currentMethods[t] = make([]string, 0)
									}

									if !strings.Contains(line, service.Name) {
										continue
									}

									// if the method is already in the currentMethods, skip
									if strings.Contains(line, service.Name) {
										// fmt.Println("    method: ", line)
										currentMethods[t] = append(currentMethods[t], parseReceiverMethodName(line))
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
				fmt.Println("\n-------- Current Methods: ", currentMethods)

				// print modules
				// fmt.Println("\nModules: ", modules)

				// iterate services again and apply any missing methods to the file
				missing := make(map[spawn.FileType][]spawn.ProtoService, 0)

				for _, module := range modules {
					for fileType, services := range module {
						fmt.Println("\nFile Type: ", fileType)
						current := currentMethods[fileType]
						fmt.Println("   Current: ", fileType, current)

						for _, service := range services {
							service := service
							fmt.Println(" - Service: ", service)
							// ft := service.FType // tx or spawn, found in currentMethods
							// if service.Name

							found := false
							for _, method := range current {
								method := method
								if method == service.Name {
									found = true
									break
								}
							}

							if !found {

								if fileType != service.FType {
									fmt.Println("  - Skipping: ", fileType, service.FType, service)
									continue
								}

								// MISSING METHOD
								// fmt.Println("Missing method: ", service.Name, service.Req, service.Res, service.Location)
								missing[fileType] = append(missing[fileType], service)
							}
						}

						// print missing
					}

					fmt.Println("\n\nMissing: ", missing)
				}

				// print filePaths
				fmt.Println("\nFile Paths: ", filePaths)

				for _, missed := range missing {
					for _, miss := range missed {
						// get miss.FType from filePaths
						p := filePaths[miss.Module][miss.FType]
						fmt.Println("File: ", p)

						content, err := os.ReadFile(p)
						if err != nil {
							fmt.Println("Error: ", err)
						}
						fmt.Println("Content: ", string(content))

						// append to the file
						fmt.Println("Append to file: ", miss.FType, miss.Name, miss.Req, miss.Res)

						switch miss.FType {
						case spawn.Tx:
							fmt.Println("Append to Tx")
						case spawn.Query:
							fmt.Println("Append to Query")

							// FeeShare implements types.QueryServer.
							// func (k Querier) FeeShare(context.Context, *types.QueryFeeShareRequest) (*types.QueryFeeShareResponse, error) {
							// 	panic("unimplemented")
							// }
							code := `// ` + miss.Name + ` implements types.QueryServer.
func (k Querier) ` + miss.Name + `(c context.Context, req *types.` + miss.Req + `) (*types.` + miss.Res + `, error) {
	panic("unimplemented")
	// ctx := sdk.UnwrapSDKContext(c)
	return &types.` + miss.Res + `{}, nil
}`

							// append to the file content after a new line at the end
							content = append(content, []byte("\n"+code)...)

							fmt.Println("New Content: ", string(content))

							if err := os.WriteFile(p, content, 0644); err != nil {
								fmt.Println("Error: ", err)
							}

						}

					}

				}

			}

		},
	}

	return cmd
}

// type methodReceiver struct {
// 	Name string
// 	Req  string
// 	Res  string
// }

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
