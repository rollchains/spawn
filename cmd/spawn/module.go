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
		logger := GetLogger()

		// ext name is the x/ 'module' name.
		extName := strings.ToLower(args[0])

		specialChars := "!@#$%^&*()_+{}|-:<>?`=[]\\;',./~"
		for _, char := range specialChars {
			if strings.Contains(extName, string(char)) {
				logger.Error("Special characters are not allowed in module names")
				return
			}
		}

		// cwd, err := os.Getwd()
		// if err != nil {
		// 	logger.Error("Error getting current working directory", err)
		// 	return
		// }

		// TODO:
		// see if cwd/x/extName exists
		// if _, err := os.Stat(path.Join(cwd, "x", extName)); err == nil {
		// 	logger.Error("TODO: Module already exists in x/. (Prompt UI to perform actions? (protoc-gen, generate keeper, setup to app.go, etc?))", "module", extName)
		// 	return
		// }

		// if err := SetupModuleProtoBase(GetLogger(), extName); err != nil {
		// 	logger.Error("Error setting up module", err)
		// 	return
		// }

		// // sets up the files in x/
		// if err := SetupModuleExtensionFiles(GetLogger(), extName); err != nil {
		// 	logger.Error("Error setting up module", err)
		// 	return
		// }

		if err := AddModuleToAppGo(GetLogger(), extName); err != nil {
			logger.Error("Error adding module to app.go", err)
			return
		}

		// TODO: Add the module base to the app.go

		// Announce the new module & how to code gen
		fmt.Printf("\n🎉 New Module '%s' generated!\n", extName)
		fmt.Println("🏅Generate Go Code:")
		fmt.Println("  - $ make proto-gen       # convert proto -> code and depinject")
	},
}

// last thing: add the module to app.go base.
func AddModuleToAppGo(logger *slog.Logger, extName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory", err)
		return err
	}

	moduleName := readCurrentModuleName(path.Join(cwd, "go.mod"))

	appGoPath := path.Join(cwd, "app", "app.go")
	fmt.Println("appGoPath", appGoPath)

	var buffer []byte
	buffer, err = os.ReadFile(appGoPath)
	if err != nil {
		fmt.Println("Error reading file", err)
		return err
	}

	// convert bytes to string
	appGoContent := string(buffer)

	// print the appGoContent
	// fmt.Println("appGoContent", appGoContent)

	appGoLines := strings.Split(appGoContent, "\n")

	// iterate line by line with appGoContent
	// newAppGo := make([]string, len(appGoLines))

	// stopImport := false

	// get the line index of "import ("

	// split appGoContent vby new lines

	imports := make([]string, 0)
	importLinesIndex := [2]int{}
	searchingForImport := false
	for idx, line := range appGoLines {
		// get the line with import (
		if strings.Contains(line, "import (") {
			fmt.Println("found import ( at line", idx)
			importLinesIndex[0] = idx + 1
			searchingForImport = true
			continue
		}

		if searchingForImport {
			if strings.Contains(line, ")") {
				fmt.Println("found ) at line", idx)
				searchingForImport = false
				importLinesIndex[1] = idx + 1
				break
			}
			imports = append(imports, line) // \t"my_import_path"
		}
	}

	// append new imports to this import list in the format: `moduleName/x/extName`, moduleName/x/extName/keeper, moduleName/x/extName/types
	imports = append(imports, fmt.Sprintf("\t\"%s/x/%s\"", moduleName, extName))
	imports = append(imports, fmt.Sprintf("\t\"%s/x/%s/keeper\"", moduleName, extName))
	imports = append(imports, fmt.Sprintf("\t\"%s/x/%s/types\"", moduleName, extName))

	// print importLinesIndex
	fmt.Println("importLinesIndex", importLinesIndex)

	// iterate over the lines in appGoContent
	newAppGo := make([]string, len(appGoLines))
	for idx, line := range appGoLines {
		// if we hit either of importLinesIndex, skip those and then append the new imports
		if idx >= importLinesIndex[0] && idx <= importLinesIndex[1] {
			continue
		}

		// append the line to newAppGo
		newAppGo = append(newAppGo, line)

	// print imports
	fmt.Println("imports", imports)

	// for idx, line := range appGoLines {
	// 	// fmt.Println("line", line)

	// 	if len(importSection) > 0 && !stopImport {
	// 		// wait until we get to the end
	// 		if strings.Contains(line, ")") {
	// 			// append the new module import data here.
	// 			importSection = append(importSection, fmt.Sprintf("\t\"%s/x/%s\"\n", moduleName, extName))
	// 			// newAppGo = append(newAppGo, importSection...)
	// 			stopImport = true
	// 		}
	// 		// else {
	// 		// 	importSection = append(importSection, line)
	// 		// }
	// 	}

	// 	// if line contains import (, start capturing lines
	// 	if strings.Contains(line, "import (") {
	// 		// capture all lines until the end )
	// 		importSection = append(importSection, line)
	// 	}
	// }

	// print importSection
	// fmt.Println("importSection", importSection)

	// TODO: find import (, and capture all lines until the end ). Then prepend the new module import data here.
	// Same for mac perms? (add input to signal if you want it to become a module acc in the maccPerms { bracket).
	// Find '// Custom' or ModuleManager in the next struct section (type ChainApp struct {)
	// Find app.EvidenceKeeper = *evidenceKeeper and append the NewKeeper lines after it
	// find app.ModuleManager = module.NewManager and append the NewAppModule with the basic setup at the end
	// SetOrderBeginBlockers, SetOrderEndBlockers, genesisModuleOrder & paramsKeeper.Subspace

	return nil
}

func SetupModuleExtensionFiles(logger *slog.Logger, extName string) error {
	extFS := simapp.ExtensionFS

	if err := os.MkdirAll(path.Join("x", extName), 0755); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory", err)
		return err
	}

	moduleName := readCurrentModuleName(path.Join(cwd, "go.mod"))
	// moduleNameProto := convertModuleNameToProto(moduleName)

	// copy x/example to x/extName
	return fs.WalkDir(extFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(cwd, relPath)
		fc, err := spawn.GetFileContent(logger, newPath, extFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
		}

		fmt.Println("newPath", newPath)

		// rename x/example path for the new module
		examplePath := path.Join("x", "example")
		if fc.ContainsPath(examplePath) {
			newBinPath := path.Join("x", extName)
			fc.NewPath = strings.ReplaceAll(fc.NewPath, examplePath, newBinPath)
		}

		fc.ReplaceAll("github.com/strangelove-ventures/simapp", moduleName)
		fc.ReplaceAll("x/example", fmt.Sprintf("x/%s", extName))

		return fc.Save()
	})
}

func SetupModuleProtoBase(logger *slog.Logger, extName string) error {
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

		// replace example -> the new x/ name
		fc.ReplaceAll("example", extName)

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
