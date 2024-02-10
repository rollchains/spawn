package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/cosmos/btcutil/bech32"
	"github.com/spf13/cobra"
	"github.com/strangelove-ventures/simapp"
)

type SpawnNewConfig struct {
	ProjectName  string
	Bech32Prefix string
	AppName      string
	AppDirName   string
	BinaryName   string
	TokenDenom   string

	IgnoreFiles []string

	Debugging bool

	DisabledFeatures []string
}

func (cfg *SpawnNewConfig) Validate() error {
	if strings.ContainsAny(cfg.ProjectName, `~!@#$%^&*()_+{}|:"<>?/.,;'[]\=-`) {
		return fmt.Errorf("project name cannot contain special characters %s", cfg.ProjectName)
	}

	return nil
}

const (
	FlagWalletPrefix = "bech32"
	FlagBinaryName   = "bin"
	FlagDebugging    = "debug"
	FlagTokenDenom   = "denom"

	FlagDisabled = "disable"
	FlagNoGit    = "no-git"
)

var (
	IgnoredFiles      = []string{"embed.go", "heighliner/"}
	SupportedFeatures = []string{"tokenfactory", "poa", "globalfee", "wasm", "icahost", "icacontroller"}
)

func init() {
	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().String(FlagBinaryName, "appd", "binary name")
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, "disable features: "+strings.Join(SupportedFeatures, ", "))
	newChain.Flags().String(FlagTokenDenom, "stake", "token denom")
	newChain.Flags().Bool(FlagNoGit, false, "git init base")
}

var newChain = &cobra.Command{
	Use:   "new-chain [project-name]",
	Short: "Create a new project",
	Example: fmt.Sprintf(
		`spawn new rollchain --%s=cosmos --%s=appd --%s=token --%s=tokenfactory,poa,globalfee`,
		FlagWalletPrefix, FlagBinaryName, FlagTokenDenom, FlagDisabled,
	),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {
		projName := strings.ToLower(args[0])
		appName := cases.Title(language.AmericanEnglish).String(projName) + "App"

		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinaryName)
		denom, _ := cmd.Flags().GetString(FlagTokenDenom)

		debug, _ := cmd.Flags().GetBool(FlagDebugging)

		disabled, _ := cmd.Flags().GetStringSlice(FlagDisabled)

		ignoreGitInit, _ := cmd.Flags().GetBool(FlagNoGit)

		cfg := &SpawnNewConfig{
			ProjectName:  projName,
			Bech32Prefix: walletPrefix,
			AppName:      appName,
			AppDirName:   "." + projName,
			BinaryName:   binName,
			TokenDenom:   denom,
			Debugging:    debug,

			// by default everything is on, then we remove what the user wants to disable
			DisabledFeatures: disabled,
		}
		if err := cfg.Validate(); err != nil {
			fmt.Println("Error validating config:", err)
			return
		}

		NewChain(cfg)

		// Create the base git repo
		if !ignoreGitInit {
			// if git already exists, don't init
			if err := execCommand("git", "init", projName, "--quiet"); err != nil {
				fmt.Println("Error initializing git:", err)
			}
			if err := os.Chdir(projName); err != nil {
				fmt.Println("Error changing to project directory:", err)
			}
			if err := execCommand("git", "add", "."); err != nil {
				fmt.Println("Error adding files to git:", err)
			}
			if err := execCommand("git", "commit", "-m", "initial commit", "--quiet"); err != nil {
				fmt.Println("Error committing initial files:", err)
			}
		}

		// Announce how to use it
		fmt.Printf("\n\n🎉 New blockchain '%s' generated!\n", projName)
		fmt.Println("🏅Getting started:")
		fmt.Println("  - $ cd " + projName)
		fmt.Println("  - $ make testnet      # build & start a testnet")
		fmt.Println("  - $ make testnet-ibc  # build & start an ibc testnet")
		fmt.Printf("  - $ make install      # build the %s binary\n", binName)
		fmt.Println("  - $ make local-image  # build docker image")
	},
}

func NewChain(cfg *SpawnNewConfig) {
	NewDirName := cfg.ProjectName
	bech32Prefix := cfg.Bech32Prefix
	projName := cfg.ProjectName
	appName := cfg.AppName
	appDirName := cfg.AppDirName
	binaryName := cfg.BinaryName
	Debugging := cfg.Debugging
	disabled := cfg.DisabledFeatures

	fmt.Println("Disabled features:", disabled)

	goModName := fmt.Sprintf("github.com/strangelove-ventures/%s", NewDirName)

	fmt.Println("Spawning new app:", NewDirName)

	if err := os.MkdirAll(NewDirName, 0755); err != nil {
		panic(err)
	}

	err := fs.WalkDir(simapp.SimApp, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(NewDirName, relPath)

		// if Debugging {
		// 	fmt.Printf("relPath: %s, newPath: %s\n", relPath, newPath)
		// }

		if relPath == "." {
			return nil
		}

		if d.IsDir() {
			// if relPath is a dir, continue walking
			return nil
		}

		for _, ignoreFile := range IgnoredFiles {
			if strings.HasSuffix(newPath, ignoreFile) || strings.HasPrefix(newPath, ignoreFile) {
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
		fileContent = removeDisabledFeatures(disabled, relPath, fileContent)

		fc := string(fileContent)

		if fc == "REMOVE" {
			// don't save this file
			return nil
		}

		if relPath == path.Join("scripts", "test_node.sh") {
			fc = strings.ReplaceAll(fc, "export BINARY=${BINARY:-wasmd}", fmt.Sprintf("export BINARY=${BINARY:-%s}", binaryName))
			fc = strings.ReplaceAll(fc, "export DENOM=${DENOM:-token}", fmt.Sprintf("export DENOM=${DENOM:-%s}", cfg.TokenDenom))
		}

		if relPath == "Dockerfile" {
			fc = strings.ReplaceAll(fc, "wasmd", binaryName)
		}

		// TODO: regex would be nicer for replacing incase it changes up stream. may never though. Also limit to specific files?
		fc = strings.ReplaceAll(fc, ".wasmd", appDirName)
		fc = strings.ReplaceAll(fc, `const appName = "WasmApp"`, fmt.Sprintf(`const appName = "%s"`, appName))
		fc = strings.ReplaceAll(fc, `Bech32Prefix = "wasm"`, fmt.Sprintf(`Bech32Prefix = "%s"`, bech32Prefix))
		fc = strings.ReplaceAll(fc, "github.com/strangelove-ventures/simapp", goModName)

		// MakeFileReplace
		fc = strings.ReplaceAll(fc, "https://github.com/CosmWasm/wasmd.git", fmt.Sprintf("https://%s.git", goModName))
		fc = strings.ReplaceAll(fc, "version.Name=wasm", fmt.Sprintf("version.Name=%s", appName)) // ldflags
		fc = strings.ReplaceAll(fc, "version.AppName=wasmd", fmt.Sprintf("version.AppName=%s", binaryName))
		fc = strings.ReplaceAll(fc, "cmd/wasmd", fmt.Sprintf("cmd/%s", binaryName))
		fc = strings.ReplaceAll(fc, "build/wasmd", fmt.Sprintf("build/%s", binaryName))
		fc = strings.ReplaceAll(fc, "wasmd keys", fmt.Sprintf("%s keys", binaryName)) // testnet

		// heighliner (not working atm)
		fc = strings.ReplaceAll(fc, "docker build . -t wasmd:local", fmt.Sprintf(`docker build . -t %s:local`, strings.ToLower(projName)))
		// TODO: remember to make the below path.Join
		// fc = strings.ReplaceAll(fc, "heighliner build -c wasmd --local --dockerfile=cosmos -f chains.yaml", fmt.Sprintf(`heighliner build -c %s --local --dockerfile=cosmos -f chains.yaml`, strings.ToLower(appName)))
		// if strings.HasSuffix(relPath, "chains.yaml") {
		// 	fc = strings.ReplaceAll(fc, "myappname", strings.ToLower(appName))
		// 	fc = strings.ReplaceAll(fc, "/go/bin/wasmd", fmt.Sprintf("/go/bin/%s", binaryName))
		// }

		// local-interchain config
		if strings.HasSuffix(relPath, "testnet.json") {
			fc = strings.ReplaceAll(fc, `"repository": "wasmd"`, fmt.Sprintf(`"repository": "%s"`, strings.ToLower(projName)))
			fc = strings.ReplaceAll(fc, `"bech32_prefix": "wasm"`, fmt.Sprintf(`"bech32_prefix": "%s"`, bech32Prefix))
			fc = strings.ReplaceAll(fc, "appName", projName)
			fc = strings.ReplaceAll(fc, "mydenom", cfg.TokenDenom)
			fc = strings.ReplaceAll(fc, "wasmd", binaryName)

			// making dynamic would be nice
			for _, addr := range []string{"wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g", "wasm1efd63aw40lxf3n4mhf7dzhjkr453axursysrvp"} {
				// bech32 convert to the new prefix
				_, bz, err := bech32.Decode(addr, 100)
				if err != nil {
					panic(err)
				}

				newAddr, err := bech32.Encode(bech32Prefix, bz)
				if err != nil {
					panic(err)
				}

				fc = strings.ReplaceAll(fc, addr, newAddr)
			}
		}

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

// Removes disabled features from the files specified
func removeDisabledFeatures(disabled []string, relativePath string, fileContent []byte) []byte {
	for _, name := range disabled {
		switch strings.ToLower(name) {
		case "tokenfactory", "token-factory", "tf":
			fileContent = removeTokenFactory(relativePath, fileContent)
		case "poa":
			fileContent = removePoa(relativePath, fileContent)
		case "globalfee":
			fileContent = removeGlobalFee(relativePath, fileContent)
		case "wasm", "cosmwasm":
			fileContent = removeWasm(relativePath, fileContent)
			continue
		case "icahost":
			// what about all ICA?
			// fileContent = removeICAHost(relativePath, fileContent)
			continue
		case "icacontroller":
			// fileContent = removeICAController(relativePath, fileContent)
			continue
		}
	}

	// remove any left over `// spawntag:` comments
	fileContent = RemoveTaggedLines("", string(fileContent), false)

	return fileContent
}

// Removes all references from the tokenfactory file
func removeTokenFactory(relativePath string, fileContent []byte) []byte {
	if relativePath == "go.mod" || relativePath == "go.sum" {
		fileContent = RemoveGoModImport("github.com/reecepbcups/tokenfactory", fileContent)
	}

	if relativePath == path.Join("app", "app.go") {
		fileContent = RemoveGeneralModule("tokenfactory", string(fileContent))
	}

	if relativePath == path.Join("scripts", "test_node.sh") {
		fileContent = RemoveGeneralModule("tokenfactory", string(fileContent))
	}

	return fileContent
}

func removePoa(relativePath string, fileContent []byte) []byte {
	if relativePath == "go.mod" || relativePath == "go.sum" {
		fileContent = RemoveGoModImport("github.com/strangelove-ventures/poa", fileContent)
	}

	if relativePath == path.Join("app", "app.go") || relativePath == path.Join("app", "ante.go") {
		fileContent = RemoveGeneralModule("poa", string(fileContent))
	}

	if relativePath == path.Join("scripts", "test_node.sh") {
		fileContent = RemoveGeneralModule("poa", string(fileContent))
	}

	return fileContent
}

func removeGlobalFee(relativePath string, fileContent []byte) []byte {

	fileContent = HandleCommentSwaps("globalfee", string(fileContent))
	fileContent = RemoveTaggedLines("globalfee", string(fileContent), true)

	if relativePath == "go.mod" || relativePath == "go.sum" {
		fileContent = RemoveGoModImport("github.com/reecepbcups/globalfee", fileContent)
	}

	if relativePath == path.Join("app", "app.go") || relativePath == path.Join("app", "ante.go") {
		fileContent = RemoveGeneralModule("globalfee", string(fileContent))
		fileContent = RemoveGeneralModule("GlobalFee", string(fileContent))
	}

	if relativePath == path.Join("scripts", "test_node.sh") {
		fileContent = RemoveGeneralModule("globalfee", string(fileContent))
	}

	return fileContent
}

func removeWasm(relativePath string, fileContent []byte) []byte {

	// remove any line with spawntag:wasm
	// if strings.Contains(string(fileContent), "spawntag:wasm") {}
	fileContent = RemoveTaggedLines("wasm", string(fileContent), true)

	if relativePath == "go.mod" || relativePath == "go.sum" {
		fileContent = RemoveGoModImport("github.com/CosmWasm/wasmd", fileContent)
		fileContent = RemoveGoModImport("github.com/CosmWasm/wasmvm", fileContent)
	}

	if relativePath == path.Join("app", "app.go") || relativePath == path.Join("app", "ante.go") {
		for _, w := range []string{
			"WasmKeeper", "wasmtypes", "wasmStack",
			"wasmOpts", "TXCounterStoreService", "WasmConfig",
			"wasmDir", "tokenfactorybindings", "github.com/CosmWasm/wasmd", "wasmvm",
		} {
			fileContent = RemoveGeneralModule(w, string(fileContent))
		}

	}

	if relativePath == path.Join("app", "ante.go") {
		fileContent = RemoveGeneralModule("wasm", string(fileContent))
	}

	if relativePath == path.Join("app", "encoding.go") {
		fileContent = RemoveGeneralModule("wasmkeeper", string(fileContent))
	}

	if relativePath == path.Join("app", "sim_test.go") {
		fileContent = RemoveGeneralModule("wasm", string(fileContent))
	}

	if relativePath == path.Join("app", "app_test.go") {
		fileContent = RemoveGeneralModule("wasmOpts", string(fileContent))
		fileContent = RemoveGeneralModule("wasmkeeper", string(fileContent))
	}

	if relativePath == path.Join("app", "test_support.go") {
		fileContent = RemoveGeneralModule("wasm", string(fileContent))
	}

	if relativePath == path.Join("app", "test_helpers.go") {
		for _, w := range []string{"emptyWasmOptions", "wasmkeeper", "WasmOpts", "wasmOpts"} {
			fileContent = RemoveGeneralModule(w, string(fileContent))
		}

	}

	if relativePath == path.Join("app", "wasm.go") {
		fileContent = []byte("REMOVE")
	}

	if relativePath == path.Join("cmd", "wasmd", "commands.go") {
		for _, w := range []string{"wasm", "wasmOpts", "wasmcli", "wasmtypes"} {
			fileContent = RemoveGeneralModule(w, string(fileContent))
		}
	}

	if relativePath == path.Join("cmd", "wasmd", "root.go") {
		for _, w := range []string{"wasmtypes", "wasmkeeper"} {
			fileContent = RemoveGeneralModule(w, string(fileContent))
		}
	}

	return fileContent
}

// Sometimes if we remove a module, we want to delete one line and use another.
func HandleCommentSwaps(name string, fileContent string) []byte {
	newContent := make([]string, 0, len(strings.Split(fileContent, "\n")))

	uncomment := fmt.Sprintf("?spawntag:%s", name)

	for idx, line := range strings.Split(fileContent, "\n") {
		hasUncommentTag := strings.Contains(line, uncomment)
		if hasUncommentTag {
			line = strings.Replace(line, "//", "", 1)
			line = strings.TrimRight(strings.Replace(line, fmt.Sprintf("// %s", uncomment), "", 1), " ")
			fmt.Printf("uncomment %s: %d, %s\n", name, idx, line)
		}

		newContent = append(newContent, line)
	}

	return []byte(strings.Join(newContent, "\n"))
}

const expectedFormat = "// spawntag:"

// RemoveTaggedLines deletes tagged lines or just removes the comment if desired.
func RemoveTaggedLines(name string, fileContent string, deleteLine bool) []byte {
	newContent := make([]string, 0, len(strings.Split(fileContent, "\n")))

	startIdx := -1
	for idx, line := range strings.Split(fileContent, "\n") {
		// TODO: regex anything in between // and spawntag such as spaces, symbols, etc.
		line = strings.ReplaceAll(line, "//spawntag:", expectedFormat) // just QOL for us to not tear our hair out

		hasTag := strings.Contains(line, fmt.Sprintf("spawntag:%s", name))
		hasMultiLineTag := strings.Contains(line, fmt.Sprintf("!spawntag:%s", name))

		// if the line has a tag, and the tag starts with a !, then we will continue until we find the end of the tag with another.
		if startIdx != -1 {
			if !hasMultiLineTag {
				continue
			}

			startIdx = -1
			fmt.Println("endIdx:", idx, line)
			continue
		}

		if hasMultiLineTag {
			if !deleteLine {
				continue
			}

			startIdx = idx
			fmt.Printf("startIdx %s: %d, %s\n", name, idx, line)
			continue
		}

		if hasTag {
			if deleteLine {
				continue
			}

			line = strings.Split(line, expectedFormat)[0]
			line = strings.TrimRight(line, " ")
		}

		newContent = append(newContent, line)
	}

	return []byte(strings.Join(newContent, "\n"))
}

// RemoveGeneralModule removes any matching names from the fileContent.
// i.e. if moduleFind is "tokenfactory" any lines with "tokenfactory" will be removed
// including comments.
// If an import or other line depends on a solo module a user wishes to remove, add a comment to the line
// such as `// tag:tokenfactory` to also remove other lines within the simapp template
func RemoveGeneralModule(removeText string, fileContent string) []byte {
	newContent := make([]string, 0, len(strings.Split(fileContent, "\n")))

	startIdx := -1
	for idx, line := range strings.Split(fileContent, "\n") {
		// if we are in a startIdx, then we need to continue until we find the close parenthesis (i.e. NewKeeper)
		if startIdx != -1 {
			fmt.Printf("rm %s startIdx: %d, %s\n", removeText, idx, line)
			if strings.TrimSpace(line) == ")" || strings.TrimSpace(line) == "}" {
				fmt.Println("endIdx:", idx, line)
				startIdx = -1
				continue
			}

			continue
		}

		lineHas := strings.Contains(line, removeText)

		if lineHas && (strings.HasSuffix(strings.TrimSpace(line), "(") || strings.HasSuffix(strings.TrimSpace(line), "{")) {
			startIdx = idx
			fmt.Printf("startIdx %s: %d, %s\n", removeText, idx, line)
			continue
		}

		if lineHas {
			fmt.Printf("rm %s: %d, %s\n", removeText, idx, line)
			continue
		}

		newContent = append(newContent, line)
	}

	return []byte(strings.Join(newContent, "\n"))
}

// given a go mod, remove a line within the file content
func RemoveGoModImport(module string, fileContent []byte) []byte {
	fcs := string(fileContent)
	lines := strings.Split(fcs, "\n")

	newLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.Contains(line, module) {
			newLines = append(newLines, line)
		}
	}

	return []byte(strings.Join(newLines, "\n"))
}
