package spawn

import (
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
