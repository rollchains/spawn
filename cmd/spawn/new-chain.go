package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/strangelove-ventures/simapp"
)

type SpawnNewConfig struct {
	ProjectName  string
	Bech32Prefix string
	AppName      string
	AppDirName   string
	BinaryName   string

	IgnoreFiles []string

	Debugging bool
}

const (
	FlagWalletPrefix = "bech32"
	FlagBinaryName   = "bin"
	FlagDebugging    = "debug"
)

var IgnoredFiles = []string{"generate.sh", "embed.go"}

func init() {
	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().String(FlagBinaryName, "appd", "binary name")
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
}

// TODO: reduce required inputs here. (or make them flags with defaults?)
var newChain = &cobra.Command{
	Use:     "new-chain [project-name]",
	Short:   "List all current chains or outputs a current config information",
	Example: fmt.Sprintf(`spawn new-chain my-project --%s=cosmos --%s=appd`, FlagWalletPrefix, FlagBinaryName),
	Args:    cobra.ExactArgs(1),
	// ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	return GetFiles(), cobra.ShellCompDirectiveNoFileComp
	// },
	Run: func(cmd *cobra.Command, args []string) {
		projName := strings.ToLower(args[0])
		appName := strings.Title(projName) + "App"

		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinaryName)

		debug, _ := cmd.Flags().GetBool(FlagDebugging)

		cfg := SpawnNewConfig{
			ProjectName:  projName,
			Bech32Prefix: walletPrefix,
			AppName:      appName,
			AppDirName:   "." + projName,
			BinaryName:   binName,

			Debugging: debug,
		}

		NewChain(cfg)

	},
}

func NewChain(cfg SpawnNewConfig) {
	NewDirName := cfg.ProjectName
	bech32Prefix := cfg.Bech32Prefix
	appName := cfg.AppName
	appDirName := cfg.AppDirName
	binaryName := cfg.BinaryName
	Debugging := cfg.Debugging

	goModName := fmt.Sprintf("github.com/strangelove-ventures/%s", NewDirName)

	fmt.Println("Spawning new app:", NewDirName)

	// create NewDirName directory
	if err := os.MkdirAll(NewDirName, 0755); err != nil {
		panic(err)
	}

	err := fs.WalkDir(simapp.SimApp, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(NewDirName, relPath)

		if Debugging {
			fmt.Println("relPath", relPath)
			fmt.Println("newPath", newPath)
		}

		if relPath == "." {
			return nil
		}

		// if relPath is a dir, continue
		if d.IsDir() {
			return nil
		}

		for _, ignoreFile := range IgnoredFiles {
			if strings.HasSuffix(newPath, ignoreFile) {
				if Debugging {
					fmt.Println("ignoring", newPath)
				}
				return nil
			}
		}

		// grab the file contents from path
		fileContent, err := simapp.SimApp.ReadFile(relPath)
		if err != nil {
			return err
		}
		fc := string(fileContent)

		// TODO: regex would be nicer for replacing incase it changes up stream. may never though. Also limit to specific files?
		fc = strings.ReplaceAll(fc, ".wasmd", appDirName)
		fc = strings.ReplaceAll(fc, `const appName = "WasmApp"`, fmt.Sprintf(`const appName = "%s"`, appName))
		fc = strings.ReplaceAll(fc, `Bech32Prefix = "wasm"`, fmt.Sprintf(`Bech32Prefix = "%s"`, bech32Prefix))
		fc = strings.ReplaceAll(fc, "github.com/strangelove-ventures/simapp", goModName)

		// MakeFileReplace
		fc = strings.ReplaceAll(fc, "https://github.com/CosmWasm/wasmd.git", fmt.Sprintf("https://%s.git", goModName))
		fc = strings.ReplaceAll(fc, "version.Name=wasm", fmt.Sprintf("version.Name=%s", appName)) // ldflags
		fc = strings.ReplaceAll(fc, "version.AppName=wasmd", fmt.Sprintf("version.AppName=%s", binaryName))
		fc = strings.ReplaceAll(fc, "github.com/CosmWasm/wasmd/app.Bech32Prefix=wasm", fmt.Sprintf("%s/app.Bech32Prefix=%s", goModName, bech32Prefix))
		fc = strings.ReplaceAll(fc, "cmd/wasmd", fmt.Sprintf("cmd/%s", binaryName))
		fc = strings.ReplaceAll(fc, "build/wasmd", fmt.Sprintf("build/%s", binaryName))

		// heighliner
		if strings.HasSuffix(relPath, "chains.yaml") {
			fc = strings.ReplaceAll(fc, "MyAppName", appName)
			fc = strings.ReplaceAll(fc, "MyAppBinary", binaryName)
		}
		fc = strings.ReplaceAll(fc, "heighliner build -c juno --local -f ./chains.yaml", fmt.Sprintf(`heighliner build -c %s --local -f ./chains.yaml`, strings.ToLower(appName)))

		// if the relPath is cmd/wasmd, replace it to be cmd/binaryName
		if strings.HasPrefix(relPath, "cmd/wasmd") {
			newPath = strings.ReplaceAll(newPath, "cmd/wasmd", fmt.Sprintf("cmd/%s", binaryName))
		}

		if err := os.MkdirAll(path.Dir(newPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(newPath, []byte(fc), 0644); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
