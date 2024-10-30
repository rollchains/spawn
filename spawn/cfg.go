package spawn

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/rollchains/spawn/simapp"
	"github.com/rollchains/spawn/spawn/types"
	localictypes "github.com/strangelove-ventures/interchaintest/local-interchain/interchain/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
)

var (
	// errFileText is used to store the contents of a failed file on save to help with debugging
	errFileText       = ""
	CosmosHubProvider *localictypes.Chain
	IgnoredFiles      = []string{"embed.go", "heighliner/"}
	isAlphaFn         = regexp.MustCompile(`^[A-Za-z]+$`).MatchString
)

func init() {
	CosmosHubProvider = localictypes.
		ChainCosmosHub("localcosmos-1").
		SetDockerImage(ibc.NewDockerImage("", "v15.1.0", "1025:1025")).
		SetBlockTime("2000ms").
		SetDefaultSDKv47Genesis(2)

	// override default genesis
	CosmosHubProvider.Genesis.Modify = []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", CosmosHubProvider.Denom),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.amount", "1"),
	}
}

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
	IgnoreGitInit   bool
	DisabledModules []string
	Logger          *slog.Logger
}

// NodeHome returns the full path to the node home directory
// ex: $HOME/.simapp
func (cfg NewChainConfig) NodeHome() string {
	return path.Join("$HOME", cfg.HomeDir)
}

func (cfg NewChainConfig) ValidateAndRun(doAnnounce bool) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("error validating config: %w", err)
	}

	if err := cfg.CreateNewChain(); err != nil {
		return fmt.Errorf("error creating new chain: %w", err)
	}

	if doAnnounce {
		cfg.AnnounceSuccessfulBuild()
	}

	return nil
}

// SetProperFeaturePairs ensures modules that are meant to be disabled, are.
// ex: if ICS is enabled, disable staking if it is not already disabled
// Normalizes the names, removes any parent dependencies, and removes duplicates
func (cfg *NewChainConfig) SetProperFeaturePairs() {
	d := RemoveDuplicates(cfg.DisabledModules)

	if cfg.IsFeatureEnabled(InterchainSecurity) {
		d = append(d, POA)
	}

	cfg.DisabledModules = d
	cfg.Logger.Debug("SetProperFeaturePairs Disabled features", "features", cfg.DisabledModules)
}

func (cfg *NewChainConfig) Validate() error {
	if cfg.ProjectName == "" {
		return types.ErrCfgEmptyProject
	}

	if strings.ContainsAny(cfg.ProjectName, `~!@#$%^&*()_+{}|:"<>?/.,;'[]\=-`) {
		return types.ErrCfgProjSpecialChars
	}

	// TODO(breaks): new chain named matrixposibcmodule && cd && spawn new nsibc --ibc-module moves in all template modules too
	// Probably since it fines it in the import namespace, so when moving files it just coppies all over.
	// The module.go `switch moduleName` is the culprit.
	if strings.Contains(cfg.ProjectName, "ibcmodule") || strings.Contains(cfg.ProjectName, "ibcmiddleware") {
		return fmt.Errorf("project name cannot contain 'ibcmodule' or 'ibcmiddleware'")
	}

	if cfg.GithubOrg == "" {
		return types.ErrCfgEmptyOrg
	}

	minDenomLen := 3
	if len(cfg.Denom) < minDenomLen {
		return types.ErrExpectedRange(types.ErrCfgDenomTooShort, minDenomLen, len(cfg.Denom))
	}

	minBinLen := 2
	if len(cfg.BinDaemon) < minBinLen {
		return types.ErrExpectedRange(types.ErrCfgBinTooShort, minBinLen, len(cfg.BinDaemon))
	}

	if cfg.Bech32Prefix == "" {
		return types.ErrCfgEmptyBech32
	}

	cfg.Bech32Prefix = strings.ToLower(cfg.Bech32Prefix)
	if !isAlphaFn(cfg.Bech32Prefix) {
		return types.ErrCfgBech32Alpha
	}

	minHomeLen := 2
	if len(cfg.HomeDir) < minHomeLen {
		return types.ErrExpectedRange(types.ErrCfgHomeDirTooShort, minHomeLen, len(cfg.HomeDir))
	}

	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	return nil
}

func (cfg *NewChainConfig) AnnounceSuccessfulBuild() {
	projName := cfg.ProjectName

	// no logger here, straight to stdout
	fmt.Printf("\n🎉 New blockchain '%s' generated!\n", projName)
	fmt.Println("🏅 Getting started:")
	fmt.Println("  - $ cd " + projName)
	fmt.Printf("  - $ gh repo create %s --source=. --remote=upstream --push --private\n", projName)
	fmt.Println("  - $ spawn module new <name>   # generate a new module scaffolding")
	fmt.Println("  - $ make testnet              # build & start a testnet with IBC")
	fmt.Println("  - $ make generate-webapp      # create cosmos webapp template")
}

func (cfg *NewChainConfig) GithubPath() string {
	return fmt.Sprintf("github.com/%s/%s", cfg.GithubOrg, cfg.ProjectName)
}

func (cfg *NewChainConfig) CreateNewChain() error {
	NewDirName := cfg.ProjectName
	logger := cfg.Logger

	// Set proper pairings for modules to be disabled if others are enabled
	cfg.SetProperFeaturePairs()

	logger.Info("Spawning new app", "name", NewDirName)
	logger.Debug("NewChain Disabled features", "features", cfg.DisabledModules)

	if err := os.MkdirAll(NewDirName, 0755); err != nil {
		logger.Error("Error creating directory", "err", err)
		return fmt.Errorf("error creating directory: %w", err)
	}

	logger.Info("Setting up main chain app", "name", NewDirName)
	if err := cfg.SetupMainChainApp(); err != nil {
		logger.Error("Error setting up main chain app", "err", err, "file", debugErrorFile(logger, NewDirName))
		return fmt.Errorf("error setting up main chain app: %w", err)
	}

	logger.Info("Setting up interchain test", "name", NewDirName, "note", "this may take a minute while it downloads dependencies")
	if err := cfg.SetupInterchainTest(); err != nil {
		logger.Error("Error setting up interchain test", "err", err, "file", debugErrorFile(logger, NewDirName))
		return fmt.Errorf("error setting up interchain test: %w", err)
	}

	logger.Info("Setting up chain metadata")
	cfg.MetadataFile().SaveJSON(fmt.Sprintf("%s/chain_metadata.json", NewDirName))
	cfg.ChainRegistryFile().SaveJSON(fmt.Sprintf("%s/chain_registry.json", NewDirName))
	cfg.ChainRegistryAssetsFile().SaveJSON(fmt.Sprintf("%s/chain_registry_assets.json", NewDirName))

	// setup local-interchain testnets
	// *testnet.json (chains/ directory)
	logger.Info("Setting up local interchain JSON")
	cfg.SetupLocalInterchainJSON()

	cfg.MakeModTidy()

	if !cfg.IgnoreGitInit {
		cfg.GitInitNewProjectRepo()
	}

	if cfg.IsFeatureEnabled("block-explorer") {
		cfg.NewPingPubExplorer()
	}

	return nil
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

		errFileText = fc.Contents
		if err := fc.FormatGoFile(); err != nil {
			return err
		}

		if err := fc.Save(); err != nil {
			return err
		}

		return nil
	})
}

func (cfg *NewChainConfig) SetupInterchainTest() error {
	newDirName := cfg.ProjectName

	// Interchaintest e2e is a nested submodule. go.mod is renamed to go.mod_ to avoid conflicts
	// It will be unwound during unpacking to properly nest it.
	ictestFS := simapp.ICTestFS
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

		if fc.IsPath(path.Join("interchaintest", "setup.go")) {
			fc.ReplaceAll( // must be first
				`ibc.NewDockerImage("wasmd", "local", "1025:1025")`,
				fmt.Sprintf(`ibc.NewDockerImage("%s", "local", "1025:1025")`, strings.ToLower(cfg.ProjectName)),
			)
			fc.ReplaceAll("mydenom", cfg.Denom)
			fc.ReplaceAll("appName", cfg.ProjectName)
			fc.ReplaceAll(`Binary  = "wasmd"`, fmt.Sprintf(`Binary  = "%s"`, cfg.BinDaemon)) // else it would replace the Cosmwasm/wasmd import path
			fc.ReplaceAll(`mybechprefix`, cfg.Bech32Prefix)

			fc.FindAndReplaceAddressBech32("wasm", cfg.Bech32Prefix)

		}

		// *All Files
		fc.ReplaceEverywhere(cfg)

		// Removes any modules references after we modify interchaintest values
		fc.RemoveDisabledFeatures(cfg)

		errFileText = fc.Contents
		if err := fc.FormatGoFile(); err != nil {
			return err
		}

		return fc.Save()
	})
}

// TODO: allow selecting for other chains to generate from (ethos, saga)
// SetupLocalInterchainJSON sets up the local-interchain testnets configuration files.
func (cfg *NewChainConfig) SetupLocalInterchainJSON() {
	c := localictypes.NewChainBuilder(cfg.ProjectName, "localchain-1", cfg.BinDaemon, cfg.Denom, cfg.Bech32Prefix).
		SetBlockTime("2000ms").
		SetDockerImage(ibc.NewDockerImage(strings.ToLower(cfg.ProjectName), "local", "")).
		SetTrustingPeriod("336h").
		SetHostPortOverride(localictypes.BaseHostPortOverride()).
		SetDefaultSDKv47Genesis(2)

	c.Genesis.Modify = []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", c.Denom),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.amount", "1"),
	}

	c.ConfigFileOverrides = []localictypes.ConfigFileOverrides{
		{
			File: "config/app.toml",
			Paths: testutil.Toml{
				"api.enabled-unsafe-cors": true,
			},
		},
		{
			File: "config/config.toml",
			Paths: testutil.Toml{
				"rpc.cors_allowed_origins": []string{"*"},
			},
		},
	}

	// Create a chain that is thisnetwork -> cosmoshub
	if cfg.IsFeatureEnabled(InterchainSecurity) {
		c.SetICSConsumerLink("localcosmos-1")
	} else {
		// Standalone testnet with no IBC connections
		if err := localictypes.NewChainsConfig(c).SaveJSON(fmt.Sprintf("%s/chains/standalone.json", cfg.ProjectName)); err != nil {
			panic(err)
		}

		// make this is an IBC testnet for POA/POS chains
		c.SetAppendedIBCPathLink(CosmosHubProvider)
	}

	cc := localictypes.NewChainsConfig(c, CosmosHubProvider)
	if err := cc.SaveJSON(fmt.Sprintf("%s/chains/testnet.json", cfg.ProjectName)); err != nil {
		panic(err)
	}

	// Create a testnet that is thisnetwork -> thisnetwork (great for IBC module testing)
	// To complex for now to support (ICS1+Chain & ICS2+Chain2)
	if !cfg.IsFeatureEnabled(InterchainSecurity) {
		chainB := localictypes.NewChainBuilder(cfg.ProjectName, "localchain-2", cfg.BinDaemon, cfg.Denom, cfg.Bech32Prefix).
			SetBlockTime("2000ms").
			SetDockerImage(ibc.NewDockerImage(strings.ToLower(cfg.ProjectName), "local", "")).
			SetTrustingPeriod("336h").
			SetDefaultSDKv47Genesis(2)

		c.SetIBCPaths([]string{}) // clear IBC paths
		c.SetAppendedIBCPathLink(chainB)

		cc = localictypes.NewChainsConfig(c, chainB)
		if err := cc.SaveJSON(fmt.Sprintf("%s/chains/self-ibc.json", cfg.ProjectName)); err != nil {
			panic(err)
		}
	}

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

func (cfg *NewChainConfig) IsFeatureEnabled(feat string) bool {
	featAlias := AliasName(feat)
	for _, disabledFeat := range cfg.DisabledModules {
		if AliasName(disabledFeat) == featAlias {
			return false
		}
	}
	return true
}

// debugErrorFile saves the errored file to a debug directory for easier debugging.
// Returning the path to the file.
func debugErrorFile(logger *slog.Logger, newDirname string) string {
	debugDir := "debugging"
	fname := fmt.Sprintf("debug-error-%s-%s.go", newDirname, time.Now().Format("2006-01-02-15-04-05"))

	if err := os.MkdirAll(debugDir, 0755); err != nil {
		logger.Error("Error creating debug directory", "err", err)
		return ""
	}

	fullPath := path.Join(debugDir, fname)
	if err := os.WriteFile(fullPath, []byte(errFileText), 0766); err != nil {
		logger.Error("Error saving debug file", "err", err)
	}

	return fullPath
}
