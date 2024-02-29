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

				fmt.Println("debugging... : ", parent, relPath, fileType)

				switch fileType {
				case spawn.Tx:
					fmt.Println("File is a transaction")
					tx := spawn.ProtoServiceParser(content)
					modules[parent][spawn.Tx] = append(modules[parent][spawn.Tx], tx...)

				case spawn.Query:
					fmt.Println("File is a query")
					query := spawn.ProtoServiceParser(content)
					modules[parent][spawn.Query] = append(modules[parent][spawn.Query], query...)
				case spawn.None:
					fmt.Println("File is neither a transaction nor a query")
				}

				return nil
			})

			fmt.Printf("Modules: %+v\n", modules)

		},
	}

	return cmd
}
