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

		if err := SetupModuleProtoBase(GetLogger(), extName); err != nil {
			logger.Error("Error setting up module", err)
			return
		}

		// sets up the files in x/
		if err := SetupModuleExtensionFiles(GetLogger(), extName); err != nil {
			logger.Error("Error setting up module", err)
			return
		}

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
	extNameTitle := strings.Title(extName)

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

	// import paths
	newImports := []string{
		fmt.Sprintf(`%s "%s/x/%s"`, extName, moduleName, extName),
		fmt.Sprintf(`%skeeper "%s/x/%s/keeper"`, extName, moduleName, extName),
		fmt.Sprintf(`%stypes "%s/x/%s/types"`, extName, moduleName, extName),
	}
	appGoLines = appendNewImportsToSource(appGoPath, newImports, appGoLines)

	// find the ModuleManager  within the ChainApp struct to add the new keeper
	// insert the new keeper type here at appModuleManagerLine -2
	appModuleManagerLine := findLineWithText(appGoLines, "*module.Manager")
	fmt.Println("appModuleManager", appModuleManagerLine)
	appGoLines = append(appGoLines[:appModuleManagerLine-2], append([]string{fmt.Sprintf(`	%sKeeper          %skeeper.Keeper`, extNameTitle, extName)}, appGoLines[appModuleManagerLine-2:]...)...)

	// find line storetypes.NewKVStoreKeys, and get the final line which ends with just )
	start, end := findLinesWithText(appGoLines, "storetypes.NewKVStoreKeys")
	fmt.Println("start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.StoreKey,`, extName)}, appGoLines[end-1:]...)...)
	// fmt.Println("lines", strings.Join(appGoLines[start:end+1], "\n"))

	// find text for app.EvidenceKeeper = *evidenceKeeper
	evidenceTextLine := findLineWithText(appGoLines, "app.EvidenceKeeper = *evidenceKeeper")
	keeperText := fmt.Sprintf(`	// Create the %s Keeper
	app.%sKeeper = %skeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[%stypes.StoreKey]),
		logger,
	)`+"\n", extName, extNameTitle, extName, extName)
	appGoLines = append(appGoLines[:evidenceTextLine+2], append([]string{keeperText}, appGoLines[evidenceTextLine+2:]...)...)

	// find app.ModuleManager = module.NewManager( lines
	start, end = findLinesWithText(appGoLines, "app.ModuleManager = module.NewManager")
	fmt.Println("start", start, "end", end)
	newAppModuleText := fmt.Sprintf(`		%s.NewAppModule(appCodec, app.%sKeeper),`+"\n", extName, extNameTitle)
	appGoLines = append(appGoLines[:end-1], append([]string{newAppModuleText}, appGoLines[end-1:]...)...)

	// begin Blockers
	start, end = findLinesWithText(appGoLines, "SetOrderBeginBlockers(")
	fmt.Println("start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.ModuleName,`, extName)}, appGoLines[end-1:]...)...)

	// end blockers
	start, end = findLinesWithText(appGoLines, "SetOrderEndBlockers(")
	fmt.Println("start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.ModuleName,`, extName)}, appGoLines[end-1:]...)...)

	// genesis module order
	start, end = findLinesWithText(appGoLines, "genesisModuleOrder := []string")
	fmt.Println("start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.ModuleName,`, extName)}, appGoLines[end-1:]...)...)

	// module params (being removed in SDK v51.)
	start, end = findLinesWithText(appGoLines, "initParamsKeeper(appCodec")
	fmt.Println("start", start, "end", end)
	appGoLines = append(appGoLines[:end-3], append([]string{fmt.Sprintf(`	paramsKeeper.Subspace(%stypes.ModuleName)`, extName)}, appGoLines[end-3:]...)...)

	// print appGoLines
	// fmt.Println("appGoLines", strings.Join(appGoLines, "\n"))

	// save the new app.go
	return os.WriteFile(appGoPath, []byte(strings.Join(appGoLines, "\n")), 0644)
}

// --- source ---

// finds the line with text, until the closing line which has a ) or }.
// TODO Very similar to RemoveModuleFromText / RemoveTaggedLines.
func findLinesWithText(source []string, text string) (startIdx, endIdx int) {
	startMultiLineFind := false
	for idx, line := range source {
		if startMultiLineFind {
			if strings.TrimSpace(line) == ")" || strings.TrimSpace(line) == "}" {
				return startIdx, idx + 1
			}
		}

		if strings.Contains(line, text) {
			startMultiLineFind = true
			startIdx = idx
			continue
		}
	}

	return 0, 0
}

func findLineWithText(source []string, text string) (lineNum int) {
	for i, line := range source {
		if strings.Contains(line, text) {
			return i
		}
	}

	return 0
}

// --- import source ---
func appendNewImportsToSource(filePath string, newImports, oldSource []string) []string {
	imports, start, end, err := parseImports(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(imports)
	fmt.Println(start, end)

	for _, newImport := range newImports {
		imports = append(imports, "\t"+newImport)
	}

	return append(oldSource[:start], append(imports, oldSource[end-1:]...)...)
}

func parseImports(filePath string) ([]string, int, int, error) {
	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, 0, 0, err
	}

	// Split the content into lines
	lines := strings.Split(string(content), "\n")

	// Find the import block and its boundaries
	importStartLine := -1
	importEndLine := -1
	for i, line := range lines {
		if strings.Contains(line, "import (") {
			importStartLine = i + 1 // Line numbers start from 1
		} else if importStartLine != -1 && strings.Contains(line, ")") {
			importEndLine = i + 1 // Line numbers start from 1
			break
		}
	}

	// If no import block found, return empty slice and line numbers as 0
	if importStartLine == -1 || importEndLine == -1 {
		return []string{}, 0, 0, nil
	}

	// Extract import strings within the import block
	var imports []string
	for _, line := range lines[importStartLine:importEndLine] {
		if strings.Contains(line, "\"") {
			imports = append(imports, line)
		}
	}

	return imports, importStartLine, importEndLine, nil
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

		// TODO: Just replace all example -> the new x/ name ?
		fc.ReplaceAll("github.com/strangelove-ventures/simapp", moduleName)
		fc.ReplaceAll("x/example", fmt.Sprintf("x/%s", extName))
		fc.ReplaceAll("package example", fmt.Sprintf("package %s", extName))
		fc.ReplaceAll("example", extName) // doing here, will see if it works as expected.

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
