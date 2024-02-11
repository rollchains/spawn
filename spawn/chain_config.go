package spawn

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/cosmos/btcutil/bech32"
	"github.com/strangelove-ventures/simapp"
)

var (
	IgnoredFiles = []string{"embed.go", "heighliner/"}
)

type NewChainConfig struct {
	ProjectName     string
	Bech32Prefix    string
	AppName         string
	AppDirName      string
	BinaryName      string
	TokenDenom      string
	GithubOrg       string
	GitInitOnCreate bool

	Debugging bool

	DisabledFeatures []string
}

func (cfg *NewChainConfig) Validate() error {
	if strings.ContainsAny(cfg.ProjectName, `~!@#$%^&*()_+{}|:"<>?/.,;'[]\=-`) {
		return fmt.Errorf("project name cannot contain special characters %s", cfg.ProjectName)
	}

	return nil
}

func (cfg *NewChainConfig) AnnounceSuccessfulBuild() {
	projName := cfg.ProjectName
	binName := cfg.BinaryName

	fmt.Printf("\n\n🎉 New blockchain '%s' generated!\n", projName)
	fmt.Println("🏅Getting started:")
	fmt.Println("  - $ cd " + projName)
	fmt.Println("  - $ make testnet      # build & start a testnet")
	fmt.Println("  - $ make testnet-ibc  # build & start an ibc testnet")
	fmt.Printf("  - $ make install      # build the %s binary\n", binName)
	fmt.Println("  - $ make local-image  # build docker image")
}

func (cfg *NewChainConfig) GithubPath() string {
	return fmt.Sprintf("github.com/%s/%s", cfg.GithubOrg, cfg.ProjectName)
}

// --- Logic ---

func (cfg *NewChainConfig) NewChain() {
	var err error

	NewDirName := cfg.ProjectName
	// bech32Prefix := cfg.Bech32Prefix
	// projName := cfg.ProjectName
	// appName := cfg.AppName
	// appDirName := cfg.AppDirName
	// binaryName := cfg.BinaryName
	Debugging := cfg.Debugging
	disabled := cfg.DisabledFeatures

	fmt.Println("Spawning new app:", NewDirName)
	fmt.Println("Disabled features:", disabled)

	// Create the new project directory
	if err := os.MkdirAll(NewDirName, 0755); err != nil {
		panic(err)
	}

	err = fs.WalkDir(simapp.SimAppFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		if relPath == "." {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		myFileContent := NewFileContent(relPath, path.Join(NewDirName, relPath))

		if myFileContent.HasIgnoreFile() {
			if Debugging {
				fmt.Println("[!] Ignoring File: ", myFileContent.NewPath)
			}
			return nil
		}

		if cfg.Debugging {
			fmt.Println(myFileContent)
		}

		// Read the file contents from the embedded FS
		if fileContent, err := simapp.SimAppFS.ReadFile(relPath); err != nil {
			return err
		} else {
			// Save the file's content to the struct
			myFileContent.Contents = string(fileContent)
		}

		// Removes any modules we care nothing about
		myFileContent.RemoveDisabledFeatures(cfg)

		// scripts/test_node.sh
		myFileContent.ReplaceTestNodeScript(cfg)
		// Dockerfile
		myFileContent.ReplaceDockerFile(cfg)
		// app/app.go
		myFileContent.ReplaceApp(cfg)
		// Makefile
		myFileContent.ReplaceMakeFile(cfg)
		// *testnet.json (chains/ directory)
		myFileContent.ReplaceLocalInterchainJSON(cfg)

		// *All Files
		myFileContent.ReplaceEverywhere(cfg)

		return myFileContent.Save()
	})
	if err != nil {
		fmt.Println(err)
	}

	// TODO: - fc.IsPath
	// Interchaintest e2e is a nested submodule. go.mod is renamed to go.mod_ to avoid conflicts
	// It will be unwound during unpacking to properly nest it.
	err = fs.WalkDir(simapp.ICTestFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(NewDirName, relPath)

		// work around to make nested embed.FS happy.
		if strings.HasSuffix(newPath, "go.mod_") {
			newPath = strings.ReplaceAll(newPath, "go.mod_", "go.mod")
		}

		if relPath == "." {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		myFileContent := NewFileContent(relPath, newPath)

		if myFileContent.HasIgnoreFile() {
			if Debugging {
				fmt.Println("[!] Ignoring File: ", myFileContent.NewPath)
			}
			return nil
		}

		if cfg.Debugging {
			fmt.Println(myFileContent)
		}

		// Read the file contents from the embedded FS
		if fileContent, err := simapp.ICTestFS.ReadFile(relPath); err != nil {
			return err
		} else {
			// Save the file's content to the struct
			myFileContent.Contents = string(fileContent)
		}

		// Removes any modules references within interchaintest that we do not care about
		myFileContent.RemoveDisabledFeatures(cfg)

		if myFileContent.IsPath(path.Join("interchaintest", "setup.go")) {
			myFileContent.ReplaceAll( // must be first
				`ibc.NewDockerImage("wasmd", "local", "1025:1025")`,
				fmt.Sprintf(`ibc.NewDockerImage("%s", "local", "1025:1025")`, strings.ToLower(cfg.ProjectName)),
			)
			myFileContent.ReplaceAll("mydenom", cfg.TokenDenom)
			myFileContent.ReplaceAll("appName", cfg.ProjectName)
			myFileContent.ReplaceAll(`Binary  = "wasmd"`, fmt.Sprintf(`Binary  = "%s"`, cfg.BinaryName)) // else it would replace the Cosmwasm/wasmd import path
			myFileContent.ReplaceAll(`Bech32 = "wasm"`, fmt.Sprintf(`Bech32 = "%s"`, cfg.Bech32Prefix))

			// making dynamic would be nice (req: regex. Would always be \"wasm1.*\" or something like that)
			// gov, acc0, acc1
			for _, addr := range []string{"wasm10d07y265gmmuvt4z0w9aw880jnsr700js7zslc", "wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g", "wasm1efd63aw40lxf3n4mhf7dzhjkr453axursysrvp"} {
				_, bz, err := bech32.Decode(addr, 100)
				if err != nil {
					panic(err)
				}

				newAddr, err := bech32.Encode(cfg.Bech32Prefix, bz)
				if err != nil {
					panic(err)
				}

				myFileContent.ReplaceAll(addr, newAddr)
			}

		}

		return myFileContent.Save()
	})
	if err != nil {
		fmt.Println(err)
	}

	if cfg.GitInitOnCreate {
		cfg.GitInitNewProjectRepo()
	}
}

// Removes disabled features from the files specified
func (fc *FileContent) RemoveDisabledFeatures(cfg *NewChainConfig) {
	for _, name := range cfg.DisabledFeatures {
		switch strings.ToLower(name) {
		case "tokenfactory", "token-factory", "tf":
			fc.RemoveTokenFactory()
		case "poa":
			fc.RemovePOA()
		case "globalfee":
			fc.RemoveGlobalFee()
		case "wasm", "cosmwasm", "cw":
			fc.RemoveCosmWasm()
		default:
			// is this acceptable? or should we just print and continue?
			panic("unknown feature to remove " + name)
		}
	}

	// remove any left over `// spawntag:` comments
	fc.RemoveTaggedLines("", false)
}

func (fc *FileContent) RemoveTokenFactory() {
	text := "tokenfactory"
	fc.RemoveGoModImport("github.com/reecepbcups/tokenfactory")

	fc.RemoveModuleFromText(text, path.Join("app", "app.go"))
	fc.RemoveModuleFromText(text, path.Join("scripts", "test_node.sh"))
	fc.RemoveModuleFromText(text, path.Join("interchaintest", "setup.go"))
}

func (fc *FileContent) RemovePOA() {
	text := "poa"
	fc.RemoveGoModImport("github.com/strangelove-ventures/poa")

	fc.RemoveModuleFromText(text,
		path.Join("app", "app.go"),
		path.Join("app", "ante.go"),
		path.Join("scripts", "test_node.sh"),
		path.Join("interchaintest", "setup.go"),
	)
}

func (fc *FileContent) RemoveGlobalFee() {
	text := "globalfee"
	fc.RemoveGoModImport("github.com/reecepbcups/globalfee")

	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	fc.RemoveModuleFromText(text,
		path.Join("app", "app.go"),
		path.Join("app", "ante.go"),
		path.Join("scripts", "test_node.sh"),
		path.Join("interchaintest", "setup.go"),
	)

	fc.RemoveModuleFromText("GlobalFee", path.Join("app", "app.go"))
}

func (fc *FileContent) RemoveCosmWasm() {
	text := "wasm"
	fc.RemoveGoModImport("github.com/CosmWasm/wasmd")
	fc.RemoveGoModImport("github.com/CosmWasm/wasmvm")

	fc.RemoveTaggedLines(text, true)

	fc.DeleteContents(path.Join("app", "wasm.go"))

	for _, word := range []string{
		"WasmKeeper", "wasmtypes", "wasmStack",
		"wasmOpts", "TXCounterStoreService", "WasmConfig",
		"wasmDir", "tokenfactorybindings", "github.com/CosmWasm/wasmd", "wasmvm",
	} {
		fc.RemoveModuleFromText(word,
			path.Join("app", "app.go"),
			path.Join("app", "ante.go"),
		)
	}

	fc.RemoveModuleFromText("wasmkeeper",
		path.Join("app", "encoding.go"),
		path.Join("app", "app_test.go"),
		path.Join("app", "test_helpers.go"),
		path.Join("cmd", "wasmd", "root.go"),
	)

	fc.RemoveModuleFromText(text,
		path.Join("app", "ante.go"),
		path.Join("app", "sim_test.go"),
		path.Join("app", "test_helpers.go"),
		path.Join("app", "test_support.go"),
		path.Join("interchaintest", "setup.go"),
		path.Join("cmd", "wasmd", "commands.go"),
		path.Join("app", "app_test.go"),
		path.Join("cmd", "wasmd", "root.go"),
	)
}

// given a go mod, remove line(s) with the importPath present.
func (fc *FileContent) RemoveGoModImport(importPath string) {
	if !fc.IsPath("go.mod") && !fc.IsPath("go.sum") {
		return
	}

	fmt.Println("removing go.mod import", fc.RelativePath, "for", importPath)

	lines := strings.Split(fc.Contents, "\n")

	newLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.Contains(line, importPath) {
			newLines = append(newLines, line)
		}
	}

	fc.Contents = strings.Join(newLines, "\n")
}
