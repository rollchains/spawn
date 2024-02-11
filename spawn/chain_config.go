package spawn

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/strangelove-ventures/simapp"
)

var (
	IgnoredFiles = []string{"embed.go", "heighliner/"}
)

type NewChainConfig struct {
	ProjectName     string // What is the diff between ProjectName and AppName? Can I merge these together?
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
	bin := cfg.BinaryName

	fmt.Printf("\n\n🎉 New blockchain '%s' generated!\n", projName)
	fmt.Println("🏅Getting started:")
	fmt.Println("  - $ cd " + projName)
	fmt.Println("  - $ make testnet      # build & start a testnet")
	fmt.Println("  - $ make testnet-ibc  # build & start an ibc testnet")
	fmt.Println("  - $ make install      # build the " + bin + " binary")
	fmt.Println("  - $ make local-image  # build docker image")
}

func (cfg *NewChainConfig) GithubPath() string {
	return fmt.Sprintf("github.com/%s/%s", cfg.GithubOrg, cfg.ProjectName)
}

func (cfg *NewChainConfig) NewChain() {
	NewDirName := cfg.ProjectName
	disabled := cfg.DisabledFeatures

	fmt.Println("Spawning new app:", NewDirName)
	fmt.Println("Disabled features:", disabled)

	if err := os.MkdirAll(NewDirName, 0755); err != nil {
		panic(err)
	}

	if err := cfg.SetupMainChainApp(); err != nil {
		fmt.Println(fmt.Errorf("error setting up main chain app: %s", err))
	}

	if err := cfg.SetupInterchainTest(); err != nil {
		fmt.Println(fmt.Errorf("error setting up interchain test: %s", err))
	}

	if cfg.GitInitOnCreate {
		cfg.GitInitNewProjectRepo()
	}
}

func (cfg *NewChainConfig) SetupMainChainApp() error {
	NewDirName := cfg.ProjectName
	Debugging := cfg.Debugging

	return fs.WalkDir(simapp.SimAppFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		if relPath == "." {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		fc := NewFileContent(relPath, path.Join(NewDirName, relPath))

		if fc.HasIgnoreFile() {
			if Debugging {
				fmt.Println("[!] Ignoring File: ", fc.NewPath)
			}
			return nil
		}

		if cfg.Debugging {
			fmt.Println(fc)
		}

		// Read the file contents from the embedded FS
		if fileContent, err := simapp.SimAppFS.ReadFile(relPath); err != nil {
			return err
		} else {
			// Save the file's content to the struct
			fc.Contents = string(fileContent)
		}

		// Removes any modules we care nothing about
		fc.RemoveDisabledFeatures(cfg)

		// scripts/test_node.sh
		fc.ReplaceTestNodeScript(cfg)
		// .github/workflows/interchaintest-e2e.yml
		fc.ReplaceGithubActionWorkflows(cfg)
		// Dockerfile
		fc.ReplaceDockerFile(cfg)
		// app/app.go
		fc.ReplaceApp(cfg)
		// Makefile
		fc.ReplaceMakeFile(cfg)
		// *testnet.json (chains/ directory)
		fc.ReplaceLocalInterchainJSON(cfg)

		// *All Files
		fc.ReplaceEverywhere(cfg)

		return fc.Save()
	})
}

func (cfg *NewChainConfig) SetupInterchainTest() error {
	NewDirName := cfg.ProjectName
	Debugging := cfg.Debugging

	// Interchaintest e2e is a nested submodule. go.mod is renamed to go.mod_ to avoid conflicts
	// It will be unwound during unpacking to properly nest it.
	return fs.WalkDir(simapp.ICTestFS, ".", func(relPath string, d fs.DirEntry, e error) error {
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

		fc := NewFileContent(relPath, newPath)

		if fc.HasIgnoreFile() {
			if Debugging {
				fmt.Println("[!] Ignoring File: ", fc.NewPath)
			}
			return nil
		}

		if cfg.Debugging {
			fmt.Println(fc)
		}

		// Read the file contents from the embedded FS
		if fileContent, err := simapp.ICTestFS.ReadFile(relPath); err != nil {
			return err
		} else {
			// Save the file's content to the struct
			fc.Contents = string(fileContent)
		}

		if fc.IsPath(path.Join("interchaintest", "setup.go")) {
			fc.ReplaceAll( // must be first
				`ibc.NewDockerImage("wasmd", "local", "1025:1025")`,
				fmt.Sprintf(`ibc.NewDockerImage("%s", "local", "1025:1025")`, strings.ToLower(cfg.ProjectName)),
			)
			fc.ReplaceAll("mydenom", cfg.TokenDenom)
			fc.ReplaceAll("appName", cfg.ProjectName)
			fc.ReplaceAll(`Binary  = "wasmd"`, fmt.Sprintf(`Binary  = "%s"`, cfg.BinaryName)) // else it would replace the Cosmwasm/wasmd import path
			fc.ReplaceAll(`Bech32 = "wasm"`, fmt.Sprintf(`Bech32 = "%s"`, cfg.Bech32Prefix))

			fc.FindAndReplaceAddressBech32("wasm", cfg.Bech32Prefix, Debugging)

		}

		// Removes any modules references after we modify interchaintest values
		fc.RemoveDisabledFeatures(cfg)

		return fc.Save()
	})
}
