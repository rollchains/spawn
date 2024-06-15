package spawn

import (
	"fmt"
	"path"
	"strings"
)

// !NOTE:
// - Always remove the ModuleKeeper before removing other types
// - Handle ComentSwaps before removing lines

var (
	TokenFactory        = "tokenfactory"
	POA                 = "poa"
	POS                 = "staking" // if ICS is used, we remove staking
	GlobalFee           = "globalfee"
	CosmWasm            = "cosmwasm"
	WasmLC              = "wasmlc"
	PacketForward       = "packetforward"
	IBCRateLimit        = "ibc-ratelimit"
	Ignite              = "ignite"
	InterchainSecurity  = "ics"
	OptimisticExecution = "optimistic-execution"

	appGo   = path.Join("app", "app.go")
	appAnte = path.Join("app", "ante.go")
)

// used for fuzz testing
var AllFeatures = []string{
	TokenFactory, POA, GlobalFee, CosmWasm, WasmLC,
	PacketForward, IBCRateLimit, Ignite, InterchainSecurity, POS,
}

// Given a string, return the reduced name for the module
// e.g. "tf" and "token-factory" both return "tokenfactory"
func AliasName(name string) string {
	switch strings.ToLower(name) {
	case TokenFactory, "token-factory", "tf":
		return "tokenfactory"
	case POA, "proof-of-authority", "proofofauthority", "poauthority":
		return POA
	case POS, "proof-of-stake", "staking", "pos":
		return POS
	case GlobalFee, "global-fee":
		return GlobalFee
	case CosmWasm, "wasm", "cw":
		return CosmWasm
	case WasmLC, "wasm-lc", "cwlc", "cosmwasm-lc", "wasm-light-client",
		"08wasm", "08-wasm", "08wasmlc", "08wasm-lc", "08-wasm-lc", "08-wasmlc":
		return WasmLC
	case PacketForward, "ibc-packetforward", "pfm":
		return PacketForward
	case Ignite, "ignite-cli":
		return Ignite
	case OptimisticExecution, "optimisticexecution", "optimistic-exec":
		return OptimisticExecution
	case IBCRateLimit, "ibc-rate-limit", "ratelimit":
		return IBCRateLimit
	case InterchainSecurity, "interchain-security":
		return InterchainSecurity
	default:
		panic(fmt.Sprintf("AliasName: unknown feature to remove: %s", name))
	}
}

// Removes disabled features from the files specified
// NOTE: Ensure you call `SetProperFeaturePairs` before calling this function
func (fc *FileContent) RemoveDisabledFeatures(cfg *NewChainConfig) {
	for _, name := range cfg.DisabledModules {
		base := AliasName(name)

		switch strings.ToLower(base) {
		// consensus
		case POA:
			fc.RemovePOA()
		case POS:
			fc.RemoveStaking()
		case InterchainSecurity:
			fc.RemoveInterchainSecurity()
		// modules
		case TokenFactory:
			fc.RemoveTokenFactory()
		case GlobalFee:
			fc.RemoveGlobalFee()
		case CosmWasm:
			fc.RemoveCosmWasm(cfg.IsFeatureDisabled(WasmLC))
		case WasmLC:
			fc.RemoveWasmLightClient()
		case PacketForward:
			fc.RemovePacketForward()
		case IBCRateLimit:
			fc.RemoveIBCRateLimit()
		// other
		case Ignite:
			fc.RemoveIgniteCLI()
		case OptimisticExecution:
			fc.RemoveOptimisticExecution()
		default:
			panic(fmt.Sprintf("unknown feature to remove %s", name))
		}
	}

	// remove any left over `// spawntag:` comments
	fc.RemoveTaggedLines("", false)
}

func (fc *FileContent) RemoveTokenFactory() {
	text := "tokenfactory"
	fc.RemoveGoModImport("github.com/strangelove-ventures/tokenfactory")

	fc.RemoveModuleFromText(text,
		appGo,
		path.Join("scripts", "test_node.sh"),
		path.Join("scripts", "test_ics_node.sh"),
		path.Join("interchaintest", "setup.go"),
		path.Join("workflows", "interchaintest-e2e.yml"),
	)

	fc.DeleteFile(path.Join("interchaintest", "tokenfactory_test.go"))
}

func (fc *FileContent) RemovePOA() {
	text := "poa"
	fc.RemoveGoModImport("github.com/strangelove-ventures/poa")

	fc.RemoveModuleFromText(text,
		appGo,
		appAnte,
		path.Join("scripts", "test_node.sh"),
		path.Join("scripts", "test_ics_node.sh"),
		path.Join("interchaintest", "setup.go"),
		path.Join("workflows", "interchaintest-e2e.yml"),
	)

	fc.DeleteFile(path.Join("interchaintest", "poa_test.go"))
	fc.DeleteFile(path.Join("interchaintest", "poa.go")) // helpers
}

func (fc *FileContent) RemoveGlobalFee() {
	text := "globalfee"
	fc.RemoveGoModImport("github.com/strangelove-ventures/globalfee")

	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	fc.RemoveModuleFromText(text,
		appGo,
		appAnte,
		path.Join("scripts", "test_node.sh"),
		path.Join("scripts", "test_ics_node.sh"),
		path.Join("interchaintest", "setup.go"),
	)

	fc.RemoveModuleFromText("GlobalFee", appGo)
}

func (fc *FileContent) RemoveCosmWasm(isWasmClientDisabled bool) {
	text := "wasm"
	fc.RemoveGoModImport("github.com/CosmWasm/wasmd")

	if isWasmClientDisabled {
		fc.RemoveGoModImport("github.com/CosmWasm/wasmvm")
	}

	fc.RemoveTaggedLines(text, true)

	fc.DeleteFile(path.Join("app", "wasm.go"))

	for _, word := range []string{
		"WasmKeeper", "wasmtypes", "wasmStack",
		"wasmOpts", "TXCounterStoreService", "WasmConfig",
		"wasmDir", "tokenfactorybindings", "github.com/CosmWasm/wasmd",
	} {
		fc.RemoveModuleFromText(word,
			appGo,
			appAnte,
		)
	}

	fc.RemoveModuleFromText("wasmkeeper",
		path.Join("app", "encoding.go"),
		path.Join("app", "app_test.go"),
		path.Join("app", "test_helpers.go"),
		path.Join("cmd", "wasmd", "root.go"),
	)

	fc.RemoveModuleFromText(text,
		appAnte,
		path.Join("app", "sim_test.go"),
		path.Join("app", "test_helpers.go"),
		path.Join("app", "test_support.go"),
		path.Join("interchaintest", "setup.go"),
		path.Join("cmd", "wasmd", "commands.go"),
		path.Join("app", "app_test.go"),
		path.Join("cmd", "wasmd", "root.go"),
		path.Join("workflows", "interchaintest-e2e.yml"),
	)

	fc.DeleteFile(path.Join("interchaintest", "cosmwasm_test.go"))
	fc.DeleteDirectoryContents(path.Join("interchaintest", "contracts"))
}

func (fc *FileContent) RemoveWasmLightClient() {
	// tag <spawntag:08wasmlc is used instead so it does not match spawntag:wasm
	text := "08wasmlc"
	fc.RemoveGoModImport("github.com/cosmos/ibc-go/modules/light-clients/08-wasm")

	fc.RemoveTaggedLines(text, true)

	fc.RemoveModuleFromText("wasmlc",
		appGo,
	)
}

func (fc *FileContent) RemovePacketForward() {
	text := "packetforward"
	fc.RemoveGoModImport("github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward")

	fc.RemoveModuleFromText(text, appGo)
	fc.RemoveModuleFromText("PacketForward", appGo)

	fc.removePacketForwardTestOnly()
}

func (fc *FileContent) RemoveIBCRateLimit() {
	text := "ratelimit"
	fc.RemoveGoModImport("github.com/Stride-Labs/ibc-rate-limiting")

	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	fc.RemoveModuleFromText("RatelimitKeeper", appGo)
	fc.RemoveModuleFromText(text,
		appGo,
		path.Join("workflows", "interchaintest-e2e.yml"),
	)

	fc.DeleteFile(path.Join("interchaintest", "ibc_rate_limit_test.go"))
}

func (fc *FileContent) RemoveIgniteCLI() {
	fc.RemoveLineWithAnyMatch("starport scaffolding")
}

func (fc *FileContent) RemoveOptimisticExecution() {
	fc.RemoveTaggedLines(OptimisticExecution, true)
}

func (fc *FileContent) RemoveInterchainSecurity() {
	text := "ics"

	fc.RemoveGoModImport("github.com/cosmos/interchain-security")

	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	// remove from e2e
	fc.RemoveModuleFromText(text, path.Join("workflows", "interchaintest-e2e.yml"))

	fc.RemoveModuleFromText("ibcconsumerkeeper.NewNonZeroKeeper", appGo)
	fc.RemoveModuleFromText("ConsumerKeeper", appGo)
	fc.RemoveModuleFromText("ScopedIBCConsumerKeeper", appGo)

	fc.RemoveLineWithAnyMatch("ibcconsumerkeeper")
	fc.RemoveLineWithAnyMatch("ibcconsumertypes")
	fc.RemoveLineWithAnyMatch("consumerante")

	fc.DeleteFile(path.Join("cmd", "wasmd", "ics_consumer.go"))
	fc.DeleteFile(path.Join("scripts", "test_ics_node.sh"))

	// TODO: remove any ictest related

}

// Remove staking module if using a custom impl like the ICS Consumer
func (fc *FileContent) RemoveStaking() {
	fc.RemovePOA() // if we already removed we should be fine

	text := "staking"
	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	fc.RemoveModuleFromText("StakingKeeper", appGo)
	fc.RemoveModuleFromText("stakingtypes", appGo)

	// TODO: depends on staking bond denom. Fix? (idk how ICS does this atm)
	fc.RemoveModuleFromText("globalfeeante", appAnte)

	// delete core modules which depend on staking
	fc.RemoveMint()
	fc.RemoveDistribution()
	fc.RemoveGov()

	// delete test helpers

	fc.DeleteFile(path.Join("app", "sim_test.go"))
	fc.DeleteFile(path.Join("app", "test_helpers.go"))
	fc.DeleteFile(path.Join("app", "test_support.go"))
	fc.DeleteFile(path.Join("app", "app_test.go"))
	fc.DeleteFile(path.Join("cmd", "wasmd", "testnet.go")) // TODO(nit): switch this to be cfg.BinDaemon instead? (check actual path vs relative)

	// Since we will be using ICS (test_ics_node.sh)
	fc.DeleteFile(path.Join("scripts", "test_node.sh"))

	fc.removePacketForwardTestOnly()

	// fix: make sh-testnet
	if fc.ContainsPath("Makefile") {
		fc.ReplaceAll("test_node.sh", "test_ics_node.sh")
	}
}

func (fc *FileContent) RemoveMint() {
	// NOTE: be careful, tenderMINT has 'mint' suffix in it. Which can match
	text := "mint"
	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	// TODO: Fix this so it does not break
	fc.RemoveModuleFromText("MintKeeper", appGo)
	fc.RemoveModuleFromText("mintkeeper", appGo)
	fc.RemoveLineWithAnyMatch("minttypes.")
}

func (fc *FileContent) RemoveGov() {
	text := "gov"
	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	fc.RemoveModuleFromText("GovKeeper", appGo)

	fc.RemoveModuleFromText("govtypes.StoreKey,", appGo)
	fc.RemoveModuleFromText("govtypes.ModuleName,", appGo) // begin blockers, genesis, etc. note the ','
}

func (fc *FileContent) RemoveDistribution() {
	text := "distribution"
	fc.HandleCommentSwaps(text)
	fc.RemoveTaggedLines(text, true)

	fc.RemoveModuleFromText("distrtypes", appGo)
	fc.RemoveModuleFromText("DistrKeeper", appGo)
	fc.RemoveModuleFromText("distrkeeper", appGo)
}

// removePacketForwardTestOnly removes the interchaintest and workflow runner for PFM
func (fc *FileContent) removePacketForwardTestOnly() {
	text := "packetforward"
	fc.RemoveModuleFromText(text,
		path.Join("workflows", "interchaintest-e2e.yml"),
	)
	fc.DeleteFile(path.Join("interchaintest", "packetforward_test.go"))
}
