package spawn

import (
	"fmt"
	"path"
	"strings"
)

// Removes disabled features from the files specified
func (fc *FileContent) RemoveDisabledFeatures(cfg *NewChainConfig) {
	for _, name := range cfg.DisabledModules {
		switch strings.ToLower(name) {
		case "tokenfactory", "token-factory", "tf":
			fc.RemoveTokenFactory()
		case "poa":
			fc.RemovePOA()
		case "globalfee":
			fc.RemoveGlobalFee()
		case "wasm", "cosmwasm", "cw":
			fc.RemoveCosmWasm()
		case "packetforward", "pfm":
			fc.RemovePacketForward()
		default:
			panic(fmt.Sprintf("unknown feature to remove %s", name))
		}
	}

	// remove any left over `// spawntag:` comments
	fc.RemoveTaggedLines("", false)
}

func (fc *FileContent) RemoveTokenFactory() {
	text := "tokenfactory"
	fc.RemoveGoModImport("github.com/reecepbcups/tokenfactory")

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
		path.Join("workflows", "interchaintest-e2e.yml"),
	)

	fc.DeleteContents(path.Join("interchaintest", "cosmwasm_test.go"))
	fc.DeleteDirectoryContents(path.Join("interchaintest", "contracts"))
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
