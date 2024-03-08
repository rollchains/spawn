package spawn

import (
	"fmt"
	"path"
	"strings"
)

// Given a string, return the reduced name for the module
// e.g. "tf" and "token-factory" both return "tokenfactory"
func AliasName(name string) string {
	switch strings.ToLower(name) {
	case "tokenfactory", "token-factory", "tf":
		return "tokenfactory"
	case "proof-of-authority", "poa", "proofofauthority", "poauthority":
		return "poa"
	case "globalfee", "global-fee":
		return "globalfee"
	case "wasm", "cosmwasm", "cw":
		return "cosmwasm"
	case "wasmlc", "wasm-lc", "cwlc", "cosmwasm-lc", "wasm-light-client", "08wasm", "08-wasm":
		return "wasmlc"
	case "ibc-packetforward", "packetforward", "pfm":
		return "packetforward"
	case "ignite", "ignite-cli":
		return "ignite"
	default:
		panic(fmt.Sprintf("AliasName: unknown feature to remove %s", name))
	}
}

// Removes disabled features from the files specified
func (fc *FileContent) RemoveDisabledFeatures(cfg *NewChainConfig) {

	isWasmLCDisabled := false
	for _, name := range cfg.DisabledModules {
		if AliasName(name) == "wasmlc" {
			isWasmLCDisabled = true
			break
		}
	}

	for _, name := range cfg.DisabledModules {

		base := AliasName(name)

		// must match MainAliasNames return
		switch strings.ToLower(base) {
		case "tokenfactory":
			fc.RemoveTokenFactory()
		case "poa":
			fc.RemovePOA()
		case "globalfee":
			fc.RemoveGlobalFee()
		case "cosmwasm":
			fc.RemoveCosmWasm(isWasmLCDisabled)
		case "wasmlc":
			fc.RemoveWasmLightClient()
		case "packetforward":
			fc.RemovePacketForward()
		case "ignite":
			fc.RemoveIgniteCLI()
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
		path.Join("app", "app.go"),
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
		path.Join("app", "app.go"),
		path.Join("app", "ante.go"),
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
		path.Join("app", "app.go"),
		path.Join("app", "ante.go"),
		path.Join("scripts", "test_node.sh"),
		path.Join("interchaintest", "setup.go"),
	)

	fc.RemoveModuleFromText("GlobalFee", path.Join("app", "app.go"))
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
		path.Join("app", "app.go"),
	)
}

func (fc *FileContent) RemovePacketForward() {
	text := "packetforward"
	fc.RemoveGoModImport("github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward")

	fc.RemoveModuleFromText(text,
		path.Join("app", "app.go"),
		path.Join("workflows", "interchaintest-e2e.yml"),
	)
	fc.RemoveModuleFromText("PacketForward", path.Join("app", "app.go"))

	fc.DeleteContents(path.Join("interchaintest", "packetforward_test.go"))
}
