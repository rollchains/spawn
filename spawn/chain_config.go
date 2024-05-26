package spawn

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"path"
	"path/filepath"
	"strings"

	"github.com/rollchains/spawn/simapp"
	"github.com/spf13/afero"
	localictypes "github.com/strangelove-ventures/interchaintest/local-interchain/interchain/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

var (
	IgnoredFiles = []string{"embed.go", "heighliner/"}

	CosmosHubProvider = localictypes.
				ChainCosmosHub("localcosmos-1").
				SetDockerImage(ibc.NewDockerImage("", "v15.1.0", "1025:1025")).
				SetDefaultSDKv47Genesis(2)
)

type NewChainConfig struct {
	// ProjectName is the name of the new chain
	ProjectName string
	// Bech32Prefix is the new wallet prefix
	Bech32Prefix string
	// The home directory of the new chain (e.g. .simapp) within the binary
	// This should typically be prefixed with a period.
	HomeDir string
	// BinDaemon is the name of the binary. (e.g. appd)
	BinDaemon string
	// Denom is the token denomination (e.g. stake, uatom, etc.)
	Denom string
	// GithubOrg is the github organization name to use for the module
	GithubOrg string
	// IgnoreGitInit is a flag to ignore git init
	IgnoreGitInit bool

	DisabledModules []string

	Logger *slog.Logger

	isUsingICS bool

	FileSystem afero.Fs
}

func (cfg NewChainConfig) Run(doAnnounce bool) {
	if err := cfg.Validate(); err != nil {
		cfg.Logger.Error("Error validating config", "err", err)
		return
	}

	cfg.NewChain()
	if doAnnounce {
		cfg.AnnounceSuccessfulBuild()
	}
}

// SetProperFeaturePairs ensures modules that are meant to be disabled, are.
// ex: if ICS is enabled, disable staking if it is not already disabled
// Normalizes the names, removes any parent dependencies, and removes duplicates
func (cfg *NewChainConfig) SetProperFeaturePairs() {
	d := RemoveDuplicates(cfg.DisabledModules)

	isUsingICS := true
	for _, name := range d {
		if AliasName(name) == InterchainSecurity {
			isUsingICS = false
		}
	}
	cfg.isUsingICS = isUsingICS

	// TODO: maybe we allow so democratic consumers are allowed?
	// Remove staking if ICS is in use
	if isUsingICS && !cfg.IsFeatureDisabled(Staking) {
		d = append(d, Staking)
	}

	cfg.DisabledModules = d
	cfg.Logger.Debug("Disabled features", "features", cfg.DisabledModules)
}

func RemoveDuplicates(disabled []string) []string {
	names := make(map[string]bool)
	for _, d := range disabled {
		names[d] = true
	}

	newDisabled := []string{}
	for d := range names {
		newDisabled = append(newDisabled, d)
	}

	return newDisabled
}

// NormalizeDisabledNames normalizes the names, removes any parent dependencies, and removes duplicates.
// It then returns the cleaned list of disabled modules.
func NormalizeDisabledNames(disabled []string, improperPairs map[string][]string) []string {
	for i, name := range disabled {
		// normalize disabled to standard aliases
		alias := AliasName(name)
		disabled[i] = alias

		// if we disable a feature which has disabled dependency, we need to disable those too
		if deps, ok := improperPairs[alias]; ok {
			// duplicates will arise, removed in the next step
			disabled = append(disabled, deps...)
		}
	}

	return RemoveDuplicates(disabled)
}

func (cfg *NewChainConfig) IsFeatureDisabled(featName string) bool {
	for _, feat := range cfg.DisabledModules {
		if AliasName(feat) == AliasName(featName) {
			return true
		}
	}
	return false
}

func (cfg *NewChainConfig) Validate() error {
	if strings.ContainsAny(cfg.ProjectName, `~!@#$%^&*()_+{}|:"<>?/.,;'[]\=-`) {
		return fmt.Errorf("project name cannot contain special characters %s", cfg.ProjectName)
	}

	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return nil
}

func (cfg *NewChainConfig) AnnounceSuccessfulBuild() {
	projName := cfg.ProjectName
	// bin := cfg.BinDaemon

	// no logger here, straight to stdout
	fmt.Printf("\n🎉 New blockchain '%s' generated!\n", projName)
	fmt.Println("🏅 Getting started:")
	fmt.Println("  - $ cd " + projName)
	fmt.Printf("  - $ gh repo create %s --source=. --remote=upstream --push --private\n", projName)
	fmt.Println("  - $ spawn module new <name>   # generate a new module scaffolding")
	fmt.Println("  - $ make testnet              # build & start a testnet with IBC")
}

func (cfg *NewChainConfig) GithubPath() string {
	return fmt.Sprintf("github.com/%s/%s", cfg.GithubOrg, cfg.ProjectName)
}

func (cfg *NewChainConfig) NewChain() {
	NewDirName := cfg.ProjectName
	logger := cfg.Logger

	// Set proper pairings for modules to be disabled if others are enabled
	cfg.SetProperFeaturePairs()

	logger.Debug("Spawning new app", "app", NewDirName)
	logger.Debug("Disabled features", "features", cfg.DisabledModules)

	if err := cfg.FileSystem.MkdirAll(NewDirName, 0755); err != nil {
		panic(err)
	}

	if err := cfg.SetupMainChainApp(); err != nil {
		logger.Error("Error setting up main chain app", "err", err)
	}

	if err := cfg.SetupInterchainTest(); err != nil {
		logger.Error("Error setting up interchain test", "err", err)
	}

	// setup local-interchain testnets
	// *testnet.json (chains/ directory)
	cfg.SetupLocalInterchainJSON()

	if !cfg.IgnoreGitInit {
		cfg.GitInitNewProjectRepo()
	}
}

func (cfg *NewChainConfig) SetupMainChainApp() error {
	newDirName := cfg.ProjectName

	simappFS := simapp.SimAppFS
	return fs.WalkDir(simappFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(newDirName, relPath)
		fc, err := GetFileContent(cfg.Logger, newPath, simappFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
		}

		// .github/workflows/interchaintest-e2e.yml (required to replace docker image in workflow)
		fc.ReplaceGithubActionWorkflows(cfg)
		// Dockerfile
		fc.ReplaceDockerFile(cfg)
		// scripts/test_node.sh
		fc.ReplaceTestNodeScript(cfg)
		// app/app.go
		fc.ReplaceApp(cfg)
		// Makefile
		fc.ReplaceMakeFile(cfg)
		// *All Files
		fc.ReplaceEverywhere(cfg)
		// Removes any modules we care nothing about
		fc.RemoveDisabledFeatures(cfg)

		// TODO:
		// if err := fc.FormatGoFile(); err != nil {
		// 	return err
		// }

		return fc.Save(cfg.FileSystem)
	})
}

func (cfg *NewChainConfig) SetupInterchainTest() error {
	newDirName := cfg.ProjectName

	// Interchaintest e2e is a nested submodule. go.mod is renamed to go.mod_ to avoid conflicts
	// It will be unwound during unpacking to properly nest it.
	ictestFS := simapp.ICTestFS

	isFormated := make(map[string]bool)

	return fs.WalkDir(ictestFS, ".", func(relPath string, d fs.DirEntry, e error) error {
		newPath := path.Join(newDirName, relPath)

		// work around to make nested embed.FS happy.
		if strings.HasSuffix(newPath, "go.mod_") {
			newPath = strings.ReplaceAll(newPath, "go.mod_", "go.mod")
		}

		fc, err := GetFileContent(cfg.Logger, newPath, ictestFS, relPath, d)
		if err != nil {
			return err
		} else if fc == nil {
			return nil
		}

		if !fc.ContainsPath("interchaintest") {
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

			fc.FindAndReplaceAddressBech32("wasm", cfg.Bech32Prefix)

		}

		// *All Files
		fc.ReplaceEverywhere(cfg)

		// Removes any modules references after we modify interchaintest values
		fc.RemoveDisabledFeatures(cfg)

		// is key in isFormated
		if !isFormated[fc.NewPath] {
			if err := fc.FormatGoFile(); err != nil {
				return err
			}
			isFormated[fc.NewPath] = true
		}

		return fc.Save(cfg.FileSystem)
	})
}

// TODO: allow selecting for other chains to generate from (ethos, saga)
// SetupLocalInterchainJSON sets up the local-interchain testnets configuration files.
func (cfg *NewChainConfig) SetupLocalInterchainJSON() {
	c := localictypes.NewChainBuilder(cfg.ProjectName, "localchain-1", cfg.BinDaemon, cfg.Denom, cfg.Bech32Prefix).
		SetBlockTime("1000ms").
		SetDockerImage(ibc.NewDockerImage(strings.ToLower(cfg.ProjectName), "local", "")).
		SetTrustingPeriod("336h").
		SetHostPortOverride(localictypes.BaseHostPortOverride()).
		SetDefaultSDKv47Genesis(2)

	if cfg.isUsingICS {
		c.SetICSConsumerLink("localcosmos-1")

		// ICS won't have gov, mint, staking, etc.
		c.Genesis.Modify = []cosmos.GenesisKV{}
		c.SetGenesis(c.Genesis)
	} else {
		// make this is an IBC testnet for POA/POS chains
		c.SetAppendedIBCPathLink(CosmosHubProvider)
	}

	cc := localictypes.NewChainsConfig(c, CosmosHubProvider)

	fPath := fmt.Sprintf("%s/chains/testnet.json", cfg.ProjectName)

	// switch case on cfg.FileSystem and the type it is
	switch cfg.FileSystem.(type) {
	case *afero.OsFs:
		if err := cc.SaveJSON(fPath); err != nil {
			panic(err)
		}
	case *afero.MemMapFs:
		if err := cfg.FileSystem.MkdirAll(filepath.Dir(fPath), 0777); err != nil {
			panic(fmt.Errorf("fs testing: failed to create directory: %w", err))
		}
		bz, err := json.MarshalIndent(cfg, "", "    ")
		if err != nil {
			panic(fmt.Errorf("fs testing: failed to marshal chains config: %w", err))
		}

		f, err := cfg.FileSystem.Create(fPath)
		if err != nil {
			panic(fmt.Errorf("fs testing: failed to create file: %w", err))
		}

		_, err = f.Write(bz)
		if err != nil {
			panic(fmt.Errorf("fs testing: failed to write file: %w", err))
		}
	default:
		panic(fmt.Errorf("fs testing: unknown fs type: %T", cfg.FileSystem))
	}
}

func GetFileContent(logger *slog.Logger, newFilePath string, fs embed.FS, relPath string, d fs.DirEntry) (*FileContent, error) {
	if relPath == "." {
		return nil, nil
	}

	if d.IsDir() {
		return nil, nil
	}

	fc := NewFileContent(logger, relPath, newFilePath)

	if fc.HasIgnoreFile() {
		logger.Debug("Ignoring file", "file", fc.NewPath)
		return nil, nil
	}

	// Read the file contents from the embedded FS
	if fileContent, err := fs.ReadFile(relPath); err != nil {
		return nil, err
	} else {
		fc.Contents = string(fileContent)
	}

	return fc, nil
}
