package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/strangelove-ventures/simapp"
	"gitub.com/strangelove-ventures/spawn/spawn"
)

var moduleCmd = &cobra.Command{
	Use:     "module [name]",
	Short:   "Create a new module scaffolding",
	Example: `spawn module mymodule`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"m", "mod", "proto"},
	Run: func(cmd *cobra.Command, args []string) {

		// TODO: Are special characters allowed?
		// ext name is the x/ 'module' name.
		extName := strings.ToLower(args[0])

		// TODO: `module name` will regen if it does not already exist,
		// else it will add in a base template. (Smart create/edit)

		// if does not exist:
		SetupModule(GetLogger(), extName)
		// else:
		// make proto-gen for the user / refresh?

		// Announce the new module & how to code gen
		fmt.Printf("\n🎉 New Module '%s' generated!\n", extName)
		fmt.Println("🏅Generate Go Code:")
		fmt.Println("  - $ make proto-gen       # convert proto -> code and depinject")
	},
}

func SetupModule(logger *slog.Logger, extName string) error {
	protoFS := simapp.ProtoModuleFS

	if err := os.MkdirAll("proto", 0755); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory", err)
		return err
	}

	moduleName := readCurrentModuleName(path.Join(cwd, "go.mod"))
	moduleNameProto := convertModuleNameToProto(moduleName)

	return fs.WalkDir(protoFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(cwd, relPath)
		fc, err := spawn.GetFileContent(logger, newPath, protoFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
		}

		// rename proto path for the new module
		exampleProtoPath := path.Join("proto", "example")
		if fc.ContainsPath(exampleProtoPath) {
			newBinPath := path.Join("proto", extName)
			fc.NewPath = strings.ReplaceAll(fc.NewPath, exampleProtoPath, newBinPath)
		}

		// any file content that has github.com/strangelove-ventures/simapp replace to moduleName

		fc.ReplaceAll("github.com/strangelove-ventures/simapp", moduleName)
		fc.ReplaceAll("strangelove_ventures.simapp", moduleNameProto)

		// TODO: maybe juts a straight up replace all on 'example' here instead?
		fc.ReplaceAll("example.module.v1", fmt.Sprintf("%s.module.v1", extName))
		fc.ReplaceAll("x/example", fmt.Sprintf("x/%s", extName))
		fc.ReplaceAll("example/Params", fmt.Sprintf("%s/Params", extName))
		fc.ReplaceAll("example/v1/params", fmt.Sprintf("%s/v1/params", extName))
		fc.ReplaceAll("package example.v1", fmt.Sprintf("package %s.v1", extName))
		fc.ReplaceAll(`import "example`, fmt.Sprintf(`import "%s`, extName))

		// TODO: set the values in the keepers / msg server automatically

		return fc.Save()
	})
}

// readCurrentModuleName reads the module name from the go.mod file on the host machine.
func readCurrentModuleName(loc string) string {
	if !strings.HasSuffix(loc, "go.mod") {
		loc = path.Join(loc, "go.mod")
	}

	// read file from path into a []byte
	var fileContent []byte
	fileContent, err := os.ReadFile(loc)
	if err != nil {
		fmt.Println("Error reading file", err)
		return ""
	}

	lines := strings.Split(string(fileContent), "\n")
	for _, line := range lines {
		if strings.Contains(line, "module") {
			return strings.Split(line, " ")[1]
		}
	}

	return ""
}

func convertModuleNameToProto(moduleName string) string {
	// github.com/rollchains/myproject -> rollchains.myproject
	text := strings.Replace(moduleName, "github.com/", "", 1)
	return strings.Replace(text, "/", ".", -1)
}
