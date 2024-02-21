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

	textcases "golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

var moduleCmd = &cobra.Command{
	Use:     "module [name]",
	Short:   "Create a new module scaffolding",
	Example: `spawn module mymodule`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"m", "mod", "proto", "ext", "extension"},
	Run: func(cmd *cobra.Command, args []string) {
		logger := GetLogger()

		// ext name is the x/ cosmos module name.
		extName := strings.ToLower(args[0])

		specialChars := "!@#$%^&*()_+{}|-:<>?`=[]\\;',./~"
		for _, char := range specialChars {
			if strings.Contains(extName, string(char)) {
				logger.Error("Special characters are not allowed in module names")
				return
			}
		}

		// TODO: don't err here and instead Prompt UI to perform actions?
		// TODO: protoc-gen, generate msg_server from proto files, etc?
		cwd, err := os.Getwd()
		if err != nil {
			logger.Error("Error getting current working directory", err)
			return
		}
		if _, err := os.Stat(path.Join(cwd, "x", extName)); err == nil {
			logger.Error("TODO: Module already exists in x/.", "module", extName)
			return
		}

		// Setup Proto files to match the new x/ cosmos module name & go.mod module namespace (i.e. github org).
		if err := SetupModuleProtoBase(GetLogger(), extName); err != nil {
			logger.Error("Error setting up proto for module", err)
			return
		}

		// sets up the files in x/
		if err := SetupModuleExtensionFiles(GetLogger(), extName); err != nil {
			logger.Error("Error setting up x/ module files", err)
			return
		}

		// Import the files to app.go
		if err := AddModuleToAppGo(GetLogger(), extName); err != nil {
			logger.Error("Error adding new x/ module to app.go", err)
			return
		}

		// Announce the new module & how to code gen the proto files.
		fmt.Printf("\n🎉 New Module '%s' generated!\n", extName)
		fmt.Println("🏅Generate Go Code:")
		fmt.Println("  - $ make proto-gen       # convert proto -> code + generate depinject api")
	},
}

// SetupModuleProtoBase iterates through the proto embedded fs and replaces the paths and goMod names to match
// the new desired module.
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

	goModName := spawn.ReadCurrentGoModuleName(path.Join(cwd, "go.mod"))
	protoNamespace := convertGoModuleNameToProtoNamespace(goModName)

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

		fc.ReplaceAll("github.com/strangelove-ventures/simapp", goModName)
		fc.ReplaceAll("strangelove_ventures.simapp", protoNamespace)

		// replace example -> the new x/ name
		fc.ReplaceAll("example", extName)

		// TODO: set the values in the keepers / msg server automatically

		return fc.Save()
	})
}

// SetupModuleExtensionFiles iterates through the x/example embedded fs and replaces the paths and goMod names to match
// the new desired module.
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

	goModName := spawn.ReadCurrentGoModuleName(path.Join(cwd, "go.mod"))

	// copy x/example to x/extName
	return fs.WalkDir(extFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(cwd, relPath)
		fc, err := spawn.GetFileContent(logger, newPath, extFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
		}

		logger.Debug("file content", "path", fc.NewPath, "content", fc.Contents)

		// rename x/example path for the new module
		examplePath := path.Join("x", "example")
		if fc.ContainsPath(examplePath) {
			newBinPath := path.Join("x", extName)
			fc.NewPath = strings.ReplaceAll(fc.NewPath, examplePath, newBinPath)
		}

		fc.ReplaceAll("github.com/strangelove-ventures/simapp", goModName)
		fc.ReplaceAll("x/example", fmt.Sprintf("x/%s", extName))
		fc.ReplaceAll("package example", fmt.Sprintf("package %s", extName))
		fc.ReplaceAll("example", extName)

		return fc.Save()
	})
}

// AddModuleToAppGo adds the new module to the app.go file.
func AddModuleToAppGo(logger *slog.Logger, extName string) error {
	extNameTitle := textcases.Title(lang.AmericanEnglish).String(extName)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory", err)
		return err
	}

	goModName := spawn.ReadCurrentGoModuleName(path.Join(cwd, "go.mod"))

	appGoPath := path.Join(cwd, "app", "app.go")
	fmt.Println("appGoPath", appGoPath)

	var buffer []byte
	buffer, err = os.ReadFile(appGoPath)
	if err != nil {
		fmt.Println("Error reading file", err)
		return err
	}

	// Gets the source code of the app.go file line by line.
	appGoLines := strings.Split(string(buffer), "\n")

	// generates the new imports for the module
	appGoLines = appendNewImportsToSource(
		appGoPath, // reads file imports from this location
		appGoLines,
		[]string{
			// example "github.com/rollchain/simapp/x/example"
			fmt.Sprintf(`%s "%s/x/%s"`, extName, goModName, extName),
			fmt.Sprintf(`%skeeper "%s/x/%s/keeper"`, extName, goModName, extName),
			fmt.Sprintf(`%stypes "%s/x/%s/types"`, extName, goModName, extName),
		},
	)

	// Add keeper to the ChainApp struct.
	appModuleManagerLine := spawn.FindLineWithText(appGoLines, "*module.Manager")
	logger.Debug("module manager", "extName", extName, "line", appModuleManagerLine)
	appGoLines = append(appGoLines[:appModuleManagerLine-2], append([]string{fmt.Sprintf(`	%sKeeper %skeeper.Keeper`, extNameTitle, extName)}, appGoLines[appModuleManagerLine-2:]...)...)

	// Setup the new module store key.
	start, end := spawn.FindLinesWithText(appGoLines, "NewKVStoreKeys(")
	logger.Debug("store key", "extName", extName, "start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.StoreKey,`, extName)}, appGoLines[end-1:]...)...)

	// Initialize the new module keeper.
	evidenceTextLine := spawn.FindLineWithText(appGoLines, "app.EvidenceKeeper = *evidenceKeeper")
	logger.Debug("evidence keeper", "extName", extName, "line", evidenceTextLine)
	keeperText := fmt.Sprintf(`	// Create the %s Keeper
	app.%sKeeper = %skeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[%stypes.StoreKey]),
		logger,
	)`+"\n", extName, extNameTitle, extName, extName)
	appGoLines = append(appGoLines[:evidenceTextLine+2], append([]string{keeperText}, appGoLines[evidenceTextLine+2:]...)...)

	// Register the app module.
	start, end = spawn.FindLinesWithText(appGoLines, "NewManager(")
	logger.Debug("module manager", "extName", extName, "start", start, "end", end)
	newAppModuleText := fmt.Sprintf(`		%s.NewAppModule(appCodec, app.%sKeeper),`+"\n", extName, extNameTitle)
	appGoLines = append(appGoLines[:end-1], append([]string{newAppModuleText}, appGoLines[end-1:]...)...)

	// Set the begin block order of the new module.
	start, end = spawn.FindLinesWithText(appGoLines, "SetOrderBeginBlockers(")
	logger.Debug("begin block order", "extName", extName, "start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.ModuleName,`, extName)}, appGoLines[end-1:]...)...)

	// Set the end block order of the new module.
	start, end = spawn.FindLinesWithText(appGoLines, "SetOrderEndBlockers(")
	logger.Debug("end block order", "extName", extName, "start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.ModuleName,`, extName)}, appGoLines[end-1:]...)...)

	// Set the genesis module order of the new module.
	start, end = spawn.FindLinesWithText(appGoLines, "genesisModuleOrder := []string")
	logger.Debug("genesis module order", "extName", extName, "start", start, "end", end)
	appGoLines = append(appGoLines[:end-1], append([]string{fmt.Sprintf(`		%stypes.ModuleName,`, extName)}, appGoLines[end-1:]...)...)

	// Register the params to x/params module. (Removed in SDK v51)
	start, end = spawn.FindLinesWithText(appGoLines, "func initParamsKeeper(")
	logger.Debug("initParamsKeeper register", "extName", extName, "start", start, "end", end)
	appGoLines = append(appGoLines[:end-3], append([]string{fmt.Sprintf(`	paramsKeeper.Subspace(%stypes.ModuleName)`, extName)}, appGoLines[end-3:]...)...)

	return os.WriteFile(appGoPath, []byte(strings.Join(appGoLines, "\n")), 0644)
}

// appendNewImportsToSource appends new imports to the source file at the end of the import block (before the closing `)` ).
func appendNewImportsToSource(filePath string, oldSource, newImports []string) []string {
	imports, start, end, err := spawn.ParseFileImports(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, newImport := range newImports {
		imports = append(imports, "\t"+newImport)
	}

	return append(oldSource[:start], append(imports, oldSource[end-1:]...)...)
}

// convertGoModuleNameToProtoNamespace converts the github.com/*/* module name to a proto module compatible name.
// i.e. github.com/rollchains/myproject -> rollchains.myproject
func convertGoModuleNameToProtoNamespace(moduleName string) string {
	text := strings.Replace(moduleName, "github.com/", "", 1)
	return strings.Replace(text, "/", ".", -1)
}
