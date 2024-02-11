package spawn

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cosmos/btcutil/bech32"
)

type FileContent struct {
	// The path from within the embeded FileSystem
	RelativePath string
	// The new location of the file
	NewPath string
	// The contents of the file from the embededFileSystem (initially unmodified)
	Content []byte // TODO: maybe save as string? Then we convert to bytes when saving?
}

func (fc *FileContent) String() string {
	return fmt.Sprintf("RelativePath: %s, NewPath: %s", fc.RelativePath, fc.NewPath)
}

func (fc *FileContent) HasIgnoreFile() bool {
	for _, ignoreFile := range IgnoredFiles {
		// or contains?
		if strings.HasSuffix(fc.NewPath, ignoreFile) || strings.HasPrefix(fc.NewPath, ignoreFile) {
			return true
		}
	}
	return false
}

func (fc *FileContent) ReplaceTestNodeScript(cfg *NewChainConfig) {
	c := string(fc.Content)
	if fc.RelativePath == path.Join("scripts", "test_node.sh") {
		c = strings.ReplaceAll(c, "export BINARY=${BINARY:-wasmd}", fmt.Sprintf("export BINARY=${BINARY:-%s}", cfg.BinaryName))
		c = strings.ReplaceAll(c, "export DENOM=${DENOM:-token}", fmt.Sprintf("export DENOM=${DENOM:-%s}", cfg.TokenDenom))
	}
	fc.Content = []byte(c)
}

func (fc *FileContent) ReplaceDockerFile(cfg *NewChainConfig) {
	c := string(fc.Content)
	if fc.RelativePath == "Dockerfile" {
		c = strings.ReplaceAll(c, "wasmd", cfg.BinaryName)
	}
	fc.Content = []byte(c)
}

func (fc *FileContent) ReplaceApp(cfg *NewChainConfig) {
	c := string(fc.Content)
	if fc.RelativePath == path.Join("app", "app.go") { // TODO: Ends with app?
		c = strings.ReplaceAll(c, ".wasmd", cfg.AppDirName)
		c = strings.ReplaceAll(c, `const appName = "WasmApp"`, fmt.Sprintf(`const appName = "%s"`, cfg.AppName))
		c = strings.ReplaceAll(c, `Bech32Prefix = "wasm"`, fmt.Sprintf(`Bech32Prefix = "%s"`, cfg.Bech32Prefix))
	}
	fc.Content = []byte(c)
}

// ReplaceEverywhereReplaces any file content that matches anywhere in the file regardless of location.
func (fc *FileContent) ReplaceEverywhere(cfg *NewChainConfig) {
	c := string(fc.Content)

	c = strings.ReplaceAll(c, "github.com/strangelove-ventures/simapp", cfg.GithubPath())

	// if the relPath is cmd/wasmd, replace it to be cmd/binaryName
	wasmBin := path.Join("cmd", "wasmd")
	if strings.HasPrefix(fc.RelativePath, wasmBin) {
		fc.NewPath = strings.ReplaceAll(fc.NewPath, wasmBin, path.Join("cmd", cfg.BinaryName))
	}

	fc.Content = []byte(c)
}

func (fc *FileContent) ReplaceMakeFile(cfg *NewChainConfig) {
	c := string(fc.Content)

	bin := cfg.BinaryName

	c = strings.ReplaceAll(c, "https://github.com/CosmWasm/wasmd.git", fmt.Sprintf("https://%s.git", cfg.GithubPath()))
	c = strings.ReplaceAll(c, "version.Name=wasm", fmt.Sprintf("version.Name=%s", cfg.AppName)) // ldflags
	c = strings.ReplaceAll(c, "version.AppName=wasmd", fmt.Sprintf("version.AppName=%s", bin))
	c = strings.ReplaceAll(c, "cmd/wasmd", fmt.Sprintf("cmd/%s", bin))
	c = strings.ReplaceAll(c, "build/wasmd", fmt.Sprintf("build/%s", bin))
	c = strings.ReplaceAll(c, "wasmd keys", fmt.Sprintf("%s keys", bin)) // for testnet

	// heighliner (not working atm)
	c = strings.ReplaceAll(c, "docker build . -t wasmd:local", fmt.Sprintf(`docker build . -t %s:local`, strings.ToLower(cfg.ProjectName)))

	// TODO: remember to make the below path.Join
	// fc = strings.ReplaceAll(fc, "heighliner build -c wasmd --local --dockerfile=cosmos -f chains.yaml", fmt.Sprintf(`heighliner build -c %s --local --dockerfile=cosmos -f chains.yaml`, strings.ToLower(appName)))
	// if strings.HasSuffix(relPath, "chains.yaml") {
	// 	fc = strings.ReplaceAll(fc, "myappname", strings.ToLower(appName))
	// 	fc = strings.ReplaceAll(fc, "/go/bin/wasmd", fmt.Sprintf("/go/bin/%s", binaryName))
	// }

	fc.Content = []byte(c)
}

func (fc *FileContent) ReplaceLocalInterchainJSON(cfg *NewChainConfig) {
	c := string(fc.Content)

	// local-interchain config
	if strings.HasSuffix(fc.RelativePath, "testnet.json") {
		c = strings.ReplaceAll(c, `"repository": "wasmd"`, fmt.Sprintf(`"repository": "%s"`, strings.ToLower(cfg.ProjectName)))
		c = strings.ReplaceAll(c, `"bech32_prefix": "wasm"`, fmt.Sprintf(`"bech32_prefix": "%s"`, cfg.Bech32Prefix))
		c = strings.ReplaceAll(c, "appName", cfg.ProjectName)
		c = strings.ReplaceAll(c, "mydenom", cfg.TokenDenom)
		c = strings.ReplaceAll(c, "wasmd", cfg.BinaryName)

		// TODO: make dynamic so we can perform on any file.
		// if \"(wasm1...)", grab value in group & bech32 replace
		for _, addr := range []string{"wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g", "wasm1efd63aw40lxf3n4mhf7dzhjkr453axursysrvp"} {
			// bech32 convert to the new prefix
			_, bz, err := bech32.Decode(addr, 100)
			if err != nil {
				panic(err)
			}

			newAddr, err := bech32.Encode(cfg.Bech32Prefix, bz)
			if err != nil {
				panic(err)
			}

			c = strings.ReplaceAll(c, addr, newAddr)
		}
	}

	fc.Content = []byte(c)
}

func (fc *FileContent) Save() error {
	if err := os.MkdirAll(path.Dir(fc.NewPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fc.NewPath, fc.Content, 0644)
}
