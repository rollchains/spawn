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
	var err error

	NewDirName := cfg.ProjectName
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

		if myFileContent.IsPath(path.Join("interchaintest", "setup.go")) {
			myFileContent.ReplaceAll( // must be first
				`ibc.NewDockerImage("wasmd", "local", "1025:1025")`,
				fmt.Sprintf(`ibc.NewDockerImage("%s", "local", "1025:1025")`, strings.ToLower(cfg.ProjectName)),
			)
			myFileContent.ReplaceAll("mydenom", cfg.TokenDenom)
			myFileContent.ReplaceAll("appName", cfg.ProjectName)
			myFileContent.ReplaceAll(`Binary  = "wasmd"`, fmt.Sprintf(`Binary  = "%s"`, cfg.BinaryName)) // else it would replace the Cosmwasm/wasmd import path
			myFileContent.ReplaceAll(`Bech32 = "wasm"`, fmt.Sprintf(`Bech32 = "%s"`, cfg.Bech32Prefix))

			myFileContent.FindAndReplaceAddressBech32("wasm", cfg.Bech32Prefix, Debugging)

		}

		// Removes any modules references after we modify interchaintest values
		myFileContent.RemoveDisabledFeatures(cfg)

		return myFileContent.Save()
	})
	if err != nil {
		fmt.Println(err)
	}

	if cfg.GitInitOnCreate {
		cfg.GitInitNewProjectRepo()
	}
}
