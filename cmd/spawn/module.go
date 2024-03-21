package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/rollchains/simapp"
	"github.com/rollchains/spawn/spawn"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	textcases "golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

const (
	FlagIsIBCMiddleware = "ibc-middleware"
)

type features struct {
	ibcMiddleware bool
}

func normalizeModuleFlags(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "ibcmiddleware", "middleware":
		name = FlagIsIBCMiddleware
	}

	return pflag.NormalizedName(name)
}

func ModuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "module",
		Short:   "Add, Remove, and Manage modules",
		Aliases: []string{"m", "mod", "proto", "ext", "extension"},
	}

	cmd.AddCommand(
		NewCmd(),
		// TODO: remove, import/add from upstream -> app.go
	)

	return cmd
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "new [name]",
		Short:   "Create a new module scaffolding",
		Example: `spawn module new mymodule [--ibc-middleware]`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"c", "create"},
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

			cwd, err := os.Getwd()
			if err != nil {
				logger.Error("Error getting current working directory", err)
				return
			}
			if _, err := os.Stat(path.Join(cwd, "x", extName)); err == nil {
				logger.Error("Module already exists in x/.", "module", extName)
				return
			}

			isIBCMiddleware, err := cmd.Flags().GetBool(FlagIsIBCMiddleware)
			if err != nil {
				logger.Error("Error getting IBC Middleware flag", err)
				return
			}

			feats := &features{
				ibcMiddleware: isIBCMiddleware,
			}

			// Setup Proto files to match the new x/ cosmos module name & go.mod module namespace (i.e. github org).
			if err := SetupModuleProtoBase(GetLogger(), extName, feats); err != nil {
				logger.Error("Error setting up proto for module", err)
				return
			}

			// sets up the files in x/
			if err := SetupModuleExtensionFiles(GetLogger(), extName, feats); err != nil {
				logger.Error("Error setting up x/ module files", err)
				return
			}

			// Import the files to app.go
			if err := AddModuleToAppGo(GetLogger(), extName, feats); err != nil {
				logger.Error("Error adding new x/ module to app.go", err)
				return
			}

			// Announce the new module & how to code gen the proto files.
			fmt.Printf("\n🎉 New Module '%s' generated!\n", extName)
			fmt.Println("🏅Generate Go Code:")
			fmt.Println("  - $ make proto-gen         # proto -> types + dependency injection API")
			fmt.Println("  - $ spawn stub-interfaces  # proto -> MsgServer (txs) + Querier generation")
		},
	}

	cmd.Flags().Bool(FlagIsIBCMiddleware, false, "Set the module as an IBC Middleware module")
	cmd.Flags().SetNormalizeFunc(normalizeModuleFlags)

	return cmd
}

// SetupModuleProtoBase iterates through the proto embedded fs and replaces the paths and goMod names to match
// the new desired module.
func SetupModuleProtoBase(logger *slog.Logger, extName string, feats *features) error {
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

	moduleName := "example"
	if feats.ibcMiddleware {
		moduleName = "ibcmiddleware"
	}

	logger.Debug("proto namespace", "goModName", goModName, "protoNamespace", protoNamespace, "moduleName", moduleName)

	return fs.WalkDir(protoFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(cwd, relPath)
		fc, err := spawn.GetFileContent(logger, newPath, protoFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
		}

		// ignore emebeded files for modules we are not working with
		switch moduleName {
		case "example":
			if strings.Contains(fc.NewPath, "ibcmiddleware") {
				return nil
			}
		case "ibcmiddleware":
			if strings.Contains(fc.NewPath, "example") {
				return nil
			}
		}

		// rename proto path for the new module
		exampleProtoPath := path.Join("proto", moduleName)
		if fc.ContainsPath(exampleProtoPath) {
			newBinPath := path.Join("proto", extName)
			fc.NewPath = strings.ReplaceAll(fc.NewPath, exampleProtoPath, newBinPath)
		}

		fc.ReplaceAll("github.com/rollchains/simapp", goModName)
		fc.ReplaceAll("strangelove_ventures.simapp", protoNamespace)

		// replace example -> the new x/ name
		fc.ReplaceAll(moduleName, extName)

		return fc.Save()
	})
}

// SetupModuleExtensionFiles iterates through the x/example embedded fs and replaces the paths and goMod names to match
// the new desired module.
func SetupModuleExtensionFiles(logger *slog.Logger, extName string, feats *features) error {
	extFS := simapp.ExtensionFS

	if err := os.MkdirAll(path.Join("x", extName), 0755); err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory", err)
		return err
	}

	moduleName := "example"
	if feats.ibcMiddleware {
		moduleName = "ibcmiddleware"
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

		// ignore emebeded files for modules we are not working with
		switch moduleName {
		case "example":
			if strings.Contains(fc.NewPath, "ibcmiddleware") {
				return nil
			}
		case "ibcmiddleware":
			if strings.Contains(fc.NewPath, "example") {
				return nil
			}
		}

		logger.Debug("file content", "path", fc.NewPath, "content", fc.Contents)

		// rename x/<module> path for the new module
		examplePath := path.Join("x", moduleName)
		if fc.ContainsPath(examplePath) {
			newBinPath := path.Join("x", extName)
			fc.NewPath = strings.ReplaceAll(fc.NewPath, examplePath, newBinPath)
		}

		fc.ReplaceAll("github.com/rollchains/simapp", goModName)
		fc.ReplaceAll(fmt.Sprintf("x/%s", moduleName), fmt.Sprintf("x/%s", extName))
		fc.ReplaceAll(fmt.Sprintf("package %s", moduleName), fmt.Sprintf("package %s", extName))
		fc.ReplaceAll(moduleName, extName)

		return fc.Save()
	})
}

// AddModuleToAppGo adds the new module to the app.go file.
func AddModuleToAppGo(logger *slog.Logger, extName string, feats *features) error {
	extNameTitle := textcases.Title(lang.AmericanEnglish).String(extName)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory", err)
		return err
	}

	goModName := spawn.ReadCurrentGoModuleName(path.Join(cwd, "go.mod"))

	appGoPath := path.Join(cwd, "app", "app.go")
	logger.Debug("app.go path", "path", appGoPath)

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

	var keeperText string
	if feats.ibcMiddleware {
		keeperText = fmt.Sprintf(`	// Create the %s Middleware Keeper
	app.%sKeeper = %skeeper.NewKeeper(
		appCodec,
		app.MsgServiceRouter(),
		app.IBCKeeper.ChannelKeeper,
	)`+"\n", extName, extNameTitle, extName)
	} else {
		keeperText = fmt.Sprintf(`	// Create the %s Keeper
	app.%sKeeper = %skeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[%stypes.StoreKey]),
		logger,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)`+"\n", extName, extNameTitle, extName, extName)
	}

	appGoLines = append(appGoLines[:evidenceTextLine+2], append([]string{keeperText}, appGoLines[evidenceTextLine+2:]...)...)

	// Register the app module.
	start, end = spawn.FindLinesWithText(appGoLines, "NewManager(")
	logger.Debug("module manager", "extName", extName, "start", start, "end", end)

	var newAppModuleText string
	if feats.ibcMiddleware {
		newAppModuleText = fmt.Sprintf(`		%s.NewAppModule(app.%sKeeper),`+"\n", extName, extNameTitle)
	} else {
		newAppModuleText = fmt.Sprintf(`		%s.NewAppModule(appCodec, app.%sKeeper),`+"\n", extName, extNameTitle)
	}

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
