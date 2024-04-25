package spawn

import (
	"fmt"
	"path"
	"strings"
)

var (
	TokenFactory       = "tokenfactory"
	POA                = "poa"
	GlobalFee          = "globalfee"
	CosmWasm           = "cosmwasm"
	WasmLC             = "wasmlc"
	PacketForward      = "packetforward"
	IBCRateLimit       = "ibc-ratelimit"
	Ignite             = "ignite"
	InterchainSecurity = "ics"
	Staking            = "staking" // if ICS is used, we remove staking

	appGo   = path.Join("app", "app.go")
	appAnte = path.Join("app", "ante.go")
)

// Given a string, return the reduced name for the module
// e.g. "tf" and "token-factory" both return "tokenfactory"
func AliasName(name string) string {
	switch strings.ToLower(name) {
	case TokenFactory, "token-factory", "tf":
		return "tokenfactory"
	case POA, "proof-of-authority", "proofofauthority", "poauthority":
		return POA
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
	case IBCRateLimit, "ibc-rate-limit":
		return IBCRateLimit
	case InterchainSecurity, "interchain-security":
		return InterchainSecurity
	case Staking:
		return Staking
	default:
		panic(fmt.Sprintf("AliasName: unknown feature to remove %s", name))
	}
}

// Removes disabled features from the files specified
func (fc *FileContent) RemoveDisabledFeatures(cfg *NewChainConfig) {

	isWasmLCDisabled := false
	isUsingICS := true
	for _, name := range cfg.DisabledModules {
		if AliasName(name) == "wasmlc" {
			isWasmLCDisabled = true
			break
		}
		if AliasName(name) == InterchainSecurity {
			isUsingICS = false
		}
	}

	if isUsingICS {
		cfg.DisabledModules = append(cfg.DisabledModules, Staking)
	}

	for _, name := range cfg.DisabledModules {
		base := AliasName(name)

		switch strings.ToLower(base) {
		case TokenFactory:
			fc.RemoveTokenFactory()
		case POA:
			fc.RemovePOA()
		case GlobalFee:
			fc.RemoveGlobalFee()
		case CosmWasm:
			fc.RemoveCosmWasm(isWasmLCDisabled)
		case WasmLC:
			fc.RemoveWasmLightClient()
		case PacketForward:
			fc.RemovePacketForward()
		case IBCRateLimit:
			fc.RemoveIBCRateLimit()
		case Ignite:
			fc.RemoveIgniteCLI()
		case InterchainSecurity:
			fc.RemoveInterchainSecurity()
		case Staking:
			fc.RemoveStaking()
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
		path.Join("interchaintest", "setup.go"),
		path.Join("workflows", "interchaintest-e2e.yml"),
	)

	fc.DeleteContents(path.Join("interchaintest", "tokenfactory_test.go"))
}

func (fc *FileContent) RemovePOA() {
	text := "poa"
	fc.RemoveGoModImport("github.com/strangelove-ventures/poa")

	fc.RemoveModuleFromText(text,
		appGo,
		appAnte,
		path.Join("scripts", "test_node.sh"),
		path.Join("interchaintest", "setup.go"),
		path.Join("workflows", "interchaintest-e2e.yml"),
	)

	fc.DeleteContents(path.Join("interchaintest", "poa_test.go"))
	fc.DeleteContents(path.Join("interchaintest", "poa.go")) // helpers
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

	fc.DeleteContents(path.Join("app", "wasm.go"))

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

	fc.DeleteContents(path.Join("interchaintest", "cosmwasm_test.go"))
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

	fc.RemoveModuleFromText(text,
		appGo,
		path.Join("workflows", "interchaintest-e2e.yml"),
	)
	fc.RemoveModuleFromText("PacketForward", appGo)

	fc.DeleteContents(path.Join("interchaintest", "packetforward_test.go"))
}

func (fc *FileContent) RemoveIBCRateLimit() {
	text := "ratelimit"
	fc.RemoveGoModImport("github.com/Stride-Labs/ibc-rate-limiting")

	fc.RemoveModuleFromText(text,
		appGo,
		path.Join("workflows", "interchaintest-e2e.yml"),
	)
	fc.RemoveModuleFromText("RatelimitKeeper", appGo)

	// fc.DeleteContents(path.Join("interchaintest", "rate_limit_test.go"))
}

func (fc *FileContent) RemoveIgniteCLI() {
	fc.RemoveLineWithAnyMatch("starport scaffolding")
}

func (fc *FileContent) RemoveInterchainSecurity() {
	fc.RemoveLineWithAnyMatch("ibcconsumerkeeper")
	fc.RemoveLineWithAnyMatch("ibcconsumertypes")
	fc.RemoveGoModImport("github.com/cosmos/interchain-security")

	fc.RemoveModuleFromText("ConsumerKeeper", appGo)
}

// Remove staking module if using a custom impl like the ICS Consumer
func (fc *FileContent) RemoveStaking() {
	fc.RemovePOA() // if we already removed we should be fine

	fc.RemoveTaggedLines("staking", true)
	fc.RemoveModuleFromText("StakingKeeper", appGo)
	fc.RemoveModuleFromText("stakingtypes", appGo)

	// for removing ICS.
	fc.RemoveMint()
	fc.RemoveDistribution()
}

func (fc *FileContent) RemoveMint() {
	fc.RemoveModuleFromText("minttypes", appGo)
	fc.RemoveModuleFromText("mint", appGo)
	fc.RemoveModuleFromText("MintKeeper", appGo)
}

func (fc *FileContent) RemoveDistribution() {
	fc.RemoveModuleFromText("distrtypes", appGo)
}
