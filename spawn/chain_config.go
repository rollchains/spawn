package spawn

import (
	"embed"
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
	// ProjectName is the name of the new chain
	ProjectName string
	// Bech32Prefix is the new wallet prefix
	Bech32Prefix string
	// The home directory of the new chain (e.g. .simapp/)
	HomeDir string
	// BinDaemon is the name of the binary. (e.g. appd)
	BinDaemon string
	// Denom is the token denomination (e.g. stake, uatom, etc.)
	Denom string
	// GithubOrg is the github organization name to use for the module
	GithubOrg string
	// IgnoreGitInit is a flag to ignore git init
	IgnoreGitInit bool
	// Debug is a flag to enable debug logging
	Debug bool

	DisabledModules []string
}

func (cfg *NewChainConfig) Validate() error {
	if strings.ContainsAny(cfg.ProjectName, `~!@#$%^&*()_+{}|:"<>?/.,;'[]\=-`) {
		return fmt.Errorf("project name cannot contain special characters %s", cfg.ProjectName)
	}

	return nil
}

func (cfg *NewChainConfig) AnnounceSuccessfulBuild() {
	projName := cfg.ProjectName
	bin := cfg.BinDaemon

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
	disabled := cfg.DisabledModules

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

	if !cfg.IgnoreGitInit {
		cfg.GitInitNewProjectRepo()
	}
}

func (cfg *NewChainConfig) SetupMainChainApp() error {
	newDirName := cfg.ProjectName

	simappFS := simapp.SimAppFS
	return fs.WalkDir(simappFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(newDirName, relPath)
		fc, err := cfg.getFileContent(newPath, simappFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
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

		return fc.Save(cfg.Debug)
	})
}

func (cfg *NewChainConfig) SetupInterchainTest() error {
	newDirName := cfg.ProjectName
	debug := cfg.Debug

	// Interchaintest e2e is a nested submodule. go.mod is renamed to go.mod_ to avoid conflicts
	// It will be unwound during unpacking to properly nest it.
	ictestFS := simapp.ICTestFS
	return fs.WalkDir(ictestFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(newDirName, relPath)

		// work around to make nested embed.FS happy.
		if strings.HasSuffix(newPath, "go.mod_") {
			newPath = strings.ReplaceAll(newPath, "go.mod_", "go.mod")
		}

		fc, err := cfg.getFileContent(newPath, ictestFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
		}

		if fc.IsPath(path.Join("interchaintest", "setup.go")) {
			fc.ReplaceAll( // must be first
				`ibc.NewDockerImage("wasmd", "local", "1025:1025")`,
				fmt.Sprintf(`ibc.NewDockerImage("%s", "local", "1025:1025")`, strings.ToLower(cfg.ProjectName)),
			)
			fc.ReplaceAll("mydenom", cfg.Denom)
			fc.ReplaceAll("appName", cfg.ProjectName)
			fc.ReplaceAll(`Binary  = "wasmd"`, fmt.Sprintf(`Binary  = "%s"`, cfg.BinDaemon)) // else it would replace the Cosmwasm/wasmd import path
			fc.ReplaceAll(`Bech32 = "wasm"`, fmt.Sprintf(`Bech32 = "%s"`, cfg.Bech32Prefix))

			fc.FindAndReplaceAddressBech32("wasm", cfg.Bech32Prefix, debug)

		}

		// Removes any modules references after we modify interchaintest values
		fc.RemoveDisabledFeatures(cfg)

		return fc.Save(cfg.Debug)
	})
}

func (cfg *NewChainConfig) getFileContent(newFilePath string, fs embed.FS, relPath string, d fs.DirEntry) (*FileContent, error) {
	if relPath == "." {
		return nil, nil
	}

	if d.IsDir() {
		return nil, nil
	}

	fc := NewFileContent(relPath, newFilePath)

	if fc.HasIgnoreFile() {
		if cfg.Debug {
			fmt.Println("[!] Ignoring File: ", fc.NewPath)
		}
		return nil, nil
	}

	if cfg.Debug {
		fmt.Println(fc)
	}

	// Read the file contents from the embedded FS
	if fileContent, err := fs.ReadFile(relPath); err != nil {
		return nil, err
	} else {
		fc.Contents = string(fileContent)
	}

	return fc, nil
}
