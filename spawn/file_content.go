package spawn

import (
	"fmt"
	"os"
	"path"
	"strings"
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

func (fc *FileContent) ReplaceAll(old, new string) {
	fc.Contents = strings.ReplaceAll(fc.Contents, old, new)
}

func (fc *FileContent) IsPath(relPath string) bool {
	return strings.HasSuffix(fc.RelativePath, relPath)
}

func (fc *FileContent) IsPathPrefixed(relPath string) bool {
	return strings.HasPrefix(fc.RelativePath, relPath)
}

func (fc *FileContent) HasIgnoreFile() bool {
	for _, ignoreFile := range IgnoredFiles {
		if fc.IsPath(ignoreFile) || fc.IsPathPrefixed(ignoreFile) {
			return true
		}
	}
	return false
}

func (fc *FileContent) DeleteContents(path string) {
	if fc.IsPath(path) {
		fmt.Println("Deleting contents for", path)
		fc.Contents = ""
	}
}

func (fc *FileContent) ReplaceTestNodeScript(cfg *NewChainConfig) {
	if fc.IsPath(path.Join("scripts", "test_node.sh")) {
		fc.ReplaceAll("export BINARY=${BINARY:-wasmd}", fmt.Sprintf("export BINARY=${BINARY:-%s}", cfg.BinaryName))
		fc.ReplaceAll("export DENOM=${DENOM:-token}", fmt.Sprintf("export DENOM=${DENOM:-%s}", cfg.TokenDenom))
	}
}

func (fc *FileContent) ReplaceDockerFile(cfg *NewChainConfig) {
	if fc.IsPath("Dockerfile") {
		fc.ReplaceAll("wasmd", cfg.BinaryName)
	}
}

func (fc *FileContent) ReplaceApp(cfg *NewChainConfig) {
	if fc.IsPath(path.Join("app", "app.go")) {
		fc.ReplaceAll(".wasmd", cfg.AppDirName)
		fc.ReplaceAll(`const appName = "WasmApp"`, fmt.Sprintf(`const appName = "%s"`, cfg.AppName))
		fc.ReplaceAll(`Bech32Prefix = "wasm"`, fmt.Sprintf(`Bech32Prefix = "%s"`, cfg.Bech32Prefix))
	}
}

// ReplaceEverywhereReplaces any file content that matches anywhere in the file regardless of location.
func (fc *FileContent) ReplaceEverywhere(cfg *NewChainConfig) {
	fc.ReplaceAll("github.com/strangelove-ventures/simapp", cfg.GithubPath())

	wasmBin := path.Join("cmd", "wasmd")
	if fc.IsPath(wasmBin) {
		newBinPath := path.Join("cmd", cfg.BinaryName)
		fc.NewPath = strings.ReplaceAll(fc.NewPath, wasmBin, newBinPath)
	}

}

func (fc *FileContent) ReplaceMakeFile(cfg *NewChainConfig) {
	bin := cfg.BinaryName

	fc.ReplaceAll("https://github.com/CosmWasm/wasmd.git", fmt.Sprintf("https://%s.git", cfg.GithubPath()))
	fc.ReplaceAll("version.Name=wasm", fmt.Sprintf("version.Name=%s", cfg.AppName)) // ldflags
	fc.ReplaceAll("version.AppName=wasmd", fmt.Sprintf("version.AppName=%s", bin))
	fc.ReplaceAll("cmd/wasmd", fmt.Sprintf("cmd/%s", bin))
	fc.ReplaceAll("build/wasmd", fmt.Sprintf("build/%s", bin))
	fc.ReplaceAll("wasmd keys", fmt.Sprintf("%s keys", bin)) // for testnet

	fc.ReplaceAll("docker build . -t wasmd:local", fmt.Sprintf(`docker build . -t %s:local`, strings.ToLower(cfg.ProjectName)))

	fc.ReplaceAll("heighliner build -c wasmd", fmt.Sprintf(`heighliner build -c %s`, strings.ToLower(cfg.ProjectName)))
	if fc.IsPath("chains.yaml") {
		fc.ReplaceAll("myappname", strings.ToLower(cfg.ProjectName))
		fc.ReplaceAll("/go/bin/wasmd", fmt.Sprintf("/go/bin/%s", bin))
	}

}

func (fc *FileContent) ReplaceLocalInterchainJSON(cfg *NewChainConfig) {
	if fc.IsPath("testnet.json") { // this matches testnet.json and ibc-testnet.json
		fc.ReplaceAll(`"repository": "wasmd"`, fmt.Sprintf(`"repository": "%s"`, strings.ToLower(cfg.ProjectName)))
		fc.ReplaceAll(`"bech32_prefix": "wasm"`, fmt.Sprintf(`"bech32_prefix": "%s"`, cfg.Bech32Prefix))
		fc.ReplaceAll("appName", cfg.ProjectName)
		fc.ReplaceAll("mydenom", cfg.TokenDenom)
		fc.ReplaceAll("wasmd", cfg.BinaryName)

		fc.FindAndReplaceAddressBech32("wasm", cfg.Bech32Prefix, cfg.Debugging)
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
