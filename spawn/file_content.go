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
	// Content []byte // TODO: maybe save as string? Then we convert to bytes when saving?

	Contents string // add a way to iterate over the \n lines?
}

func NewFileContent(relativePath, newPath string) *FileContent {
	return &FileContent{
		RelativePath: relativePath,
		NewPath:      newPath,
		Contents:     "",
	}
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
	if fc.RelativePath == path.Join("scripts", "test_node.sh") {
		fc.Contents = strings.ReplaceAll(fc.Contents, "export BINARY=${BINARY:-wasmd}", fmt.Sprintf("export BINARY=${BINARY:-%s}", cfg.BinaryName))
		fc.Contents = strings.ReplaceAll(fc.Contents, "export DENOM=${DENOM:-token}", fmt.Sprintf("export DENOM=${DENOM:-%s}", cfg.TokenDenom))
	}
}

func (fc *FileContent) ReplaceDockerFile(cfg *NewChainConfig) {
	if fc.RelativePath == "Dockerfile" {
		fc.Contents = strings.ReplaceAll(fc.Contents, "wasmd", cfg.BinaryName)
	}
}

func (fc *FileContent) ReplaceApp(cfg *NewChainConfig) {
	if fc.RelativePath == path.Join("app", "app.go") { // TODO: Ends with app?
		fc.Contents = strings.ReplaceAll(fc.Contents, ".wasmd", cfg.AppDirName)
		fc.Contents = strings.ReplaceAll(fc.Contents, `const appName = "WasmApp"`, fmt.Sprintf(`const appName = "%s"`, cfg.AppName))
		fc.Contents = strings.ReplaceAll(fc.Contents, `Bech32Prefix = "wasm"`, fmt.Sprintf(`Bech32Prefix = "%s"`, cfg.Bech32Prefix))
	}
}

// ReplaceEverywhereReplaces any file content that matches anywhere in the file regardless of location.
func (fc *FileContent) ReplaceEverywhere(cfg *NewChainConfig) {

	fc.Contents = strings.ReplaceAll(fc.Contents, "github.com/strangelove-ventures/simapp", cfg.GithubPath())

	// if the relPath is cmd/wasmd, replace it to be cmd/binaryName
	wasmBin := path.Join("cmd", "wasmd")
	if strings.HasPrefix(fc.RelativePath, wasmBin) {
		fc.NewPath = strings.ReplaceAll(fc.NewPath, wasmBin, path.Join("cmd", cfg.BinaryName))
	}

}

func (fc *FileContent) ReplaceMakeFile(cfg *NewChainConfig) {

	bin := cfg.BinaryName

	fc.Contents = strings.ReplaceAll(fc.Contents, "https://github.com/CosmWasm/wasmd.git", fmt.Sprintf("https://%s.git", cfg.GithubPath()))
	fc.Contents = strings.ReplaceAll(fc.Contents, "version.Name=wasm", fmt.Sprintf("version.Name=%s", cfg.AppName)) // ldflags
	fc.Contents = strings.ReplaceAll(fc.Contents, "version.AppName=wasmd", fmt.Sprintf("version.AppName=%s", bin))
	fc.Contents = strings.ReplaceAll(fc.Contents, "cmd/wasmd", fmt.Sprintf("cmd/%s", bin))
	fc.Contents = strings.ReplaceAll(fc.Contents, "build/wasmd", fmt.Sprintf("build/%s", bin))
	fc.Contents = strings.ReplaceAll(fc.Contents, "wasmd keys", fmt.Sprintf("%s keys", bin)) // for testnet

	// heighliner (not working atm)
	fc.Contents = strings.ReplaceAll(fc.Contents, "docker build . -t wasmd:local", fmt.Sprintf(`docker build . -t %s:local`, strings.ToLower(cfg.ProjectName)))

	// TODO: remember to make the below path.Join
	// fc = strings.ReplaceAll(fc, "heighliner build -c wasmd --local --dockerfile=cosmos -f chains.yaml", fmt.Sprintf(`heighliner build -c %s --local --dockerfile=cosmos -f chains.yaml`, strings.ToLower(appName)))
	// if strings.HasSuffix(relPath, "chains.yaml") {
	// 	fc = strings.ReplaceAll(fc, "myappname", strings.ToLower(appName))
	// 	fc = strings.ReplaceAll(fc, "/go/bin/wasmd", fmt.Sprintf("/go/bin/%s", binaryName))
	// }

}

func (fc *FileContent) ReplaceLocalInterchainJSON(cfg *NewChainConfig) {

	// local-interchain config
	if strings.HasSuffix(fc.RelativePath, "testnet.json") {
		fc.Contents = strings.ReplaceAll(fc.Contents, `"repository": "wasmd"`, fmt.Sprintf(`"repository": "%s"`, strings.ToLower(cfg.ProjectName)))
		fc.Contents = strings.ReplaceAll(fc.Contents, `"bech32_prefix": "wasm"`, fmt.Sprintf(`"bech32_prefix": "%s"`, cfg.Bech32Prefix))
		fc.Contents = strings.ReplaceAll(fc.Contents, "appName", cfg.ProjectName)
		fc.Contents = strings.ReplaceAll(fc.Contents, "mydenom", cfg.TokenDenom)
		fc.Contents = strings.ReplaceAll(fc.Contents, "wasmd", cfg.BinaryName)

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

			fc.Contents = strings.ReplaceAll(fc.Contents, addr, newAddr)
		}
	}

}

func (fc *FileContent) Save() error {
	if fc.Contents == "" {
		fmt.Printf("Save() No contents for %s. Not saving\n", fc.NewPath)
		return nil
	}

	if err := os.MkdirAll(path.Dir(fc.NewPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fc.NewPath, []byte(fc.Contents), 0644)
}
