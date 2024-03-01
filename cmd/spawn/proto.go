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
			for name, module := range modules {
				fmt.Println("Module: ", name)

				modulePath := path.Join(cwd, "x", name, "keeper") // hardcode for keeper is less than ideal, but will do for now

				for _, services := range module {
					// fmt.Println("FileType: ", fileType)
					for _, service := range services {
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

							// open and read the contents of f
							content, err := os.ReadFile(path.Join(modulePath, f.Name()))
							if err != nil {
								fmt.Println("Error: ", err)
							}

							t := isFileQueryOrMsgServer(content)
							if t == "none" {
								continue
							}

							fmt.Println("File: ", f, t)

							// find any line with `func ` in it
							currentMethods := make([]string, 0)
							lines := strings.Split(string(content), "\n")
							for _, line := range lines {
								// receiver func
								if strings.Contains(line, "func (") {
									// fmt.Println("  func: ", line)
									currentMethods = append(currentMethods, line)
								}
							}

							fmt.Println("Current methods: ", currentMethods)

							// get missing methods
							// missingMethods := make([]string, 0)
							// for _, method := range service {
							// 	found := false
							// 	for _, currentMethod := range currentMethods {
							// 		if strings.Contains(currentMethod, method.Name) {
							// 			found = true
							// 			break
							// 		}
							// 	}

							// 	if !found {
							// 		missingMethods = append(missingMethods, method.Name)
							// 	}
							// }
						}

					}
				}
			}

		},
	}

	return cmd
}

// if a file only has an UpdateParams, but UpdateParams and OtherMethod are in proto, just add the OtherMethod signature
func modifyFileContentsToMatchProtoDefinitions() {

}

func isFileQueryOrMsgServer(bz []byte) string {
	s := strings.ToLower(string(bz))

	if strings.Contains(s, "queryserver") || strings.Contains(s, "querier") {
		return "query"
	}

	if strings.Contains(s, "msgserver") || strings.Contains(s, "msgservice") {
		return "tx"
	}

	return "none"
}
