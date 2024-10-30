package spawn

import (
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/cosmos/btcutil/bech32"
	"golang.org/x/tools/imports"
)

// bech32 prefix found within the files
const baseBech32Prefix = "wasm"

var (
	// AlreadyChecked allows for better debugging to reduce fc spam
	AlreadyCheckedDeletion      = make(map[string]bool)
	AlreadyCheckedGoModDeletion = make(map[string]bool)
)

type FileContent struct {
	// The path from within the embedded FileSystem
	RelativePath string
	// The new location of the file
	NewPath string
	// The contents of the file from the embededFileSystem (initially unmodified)
	Contents string

	Logger *slog.Logger
}

func NewFileContent(logger *slog.Logger, relativePath, newPath string) *FileContent {
	return &FileContent{
		RelativePath: relativePath,
		NewPath:      newPath,
		Contents:     "",
		Logger:       logger,
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

func (fc *FileContent) ContainsPath(relPath string) bool {
	return strings.Contains(fc.RelativePath, relPath)
}

func (fc *FileContent) InPaths(relPaths []string) bool {
	for _, relPath := range relPaths {
		if fc.IsPath(relPath) {
			return true
		}
	}
	return false
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

// DeleteFile sets the content of the file to nothing. On save, the file is ignored if there is no content.
func (fc *FileContent) DeleteFile(path string) {
	if fc.IsPath(path) && !AlreadyCheckedDeletion[path] {
		AlreadyCheckedDeletion[path] = true

		fc.Logger.Debug("Deleting contents for", "path", path)
		fc.Contents = ""
	}
}

func (fc *FileContent) DeleteDirectoryContents(path string) {
	if fc.ContainsPath(path) {
		fc.Logger.Debug("Deleting contents for", "path", path)
		fc.Contents = ""
	}
}

func (fc *FileContent) ReplaceTestNodeScript(cfg *NewChainConfig) {
	if fc.IsPath(path.Join("scripts", "test_node.sh")) || fc.IsPath(path.Join("scripts", "test_ics_node.sh")) {
		fc.ReplaceAll("export BINARY=${BINARY:-wasmd}", fmt.Sprintf("export BINARY=${BINARY:-%s}", cfg.BinDaemon))
		fc.ReplaceAll("export DENOM=${DENOM:-token}", fmt.Sprintf("export DENOM=${DENOM:-%s}", cfg.Denom))

		fc.ReplaceAll(`export HOME_DIR=$(eval echo "${HOME_DIR:-"~/.simapp"}")`, fmt.Sprintf(`export HOME_DIR=$(eval echo "${HOME_DIR:-"~/%s"}")`, cfg.HomeDir))
		fc.ReplaceAll(`HOME_DIR="~/.simapp"`, fmt.Sprintf(`HOME_DIR="~/%s"`, cfg.HomeDir))

		fc.FindAndReplaceAddressBech32(baseBech32Prefix, cfg.Bech32Prefix)
	}
}

func (fc *FileContent) ReplaceGithubActionWorkflows(cfg *NewChainConfig) {
	if fc.IsPath(path.Join(".github", "workflows", "interchaintest-e2e.yml")) {
		fc.ReplaceAll("wasmd:local", fmt.Sprintf("%s:local", strings.ToLower(cfg.ProjectName)))
	}
	if fc.IsPath(path.Join(".github", "workflows", "docker-release.yml")) {
		fc.ReplaceAll("/go/bin/wasmd", fmt.Sprintf("/go/bin/%s", cfg.BinDaemon))
	}
	if fc.ContainsPath(".goreleaser.yaml") {
		fc.ReplaceAll("rollchains", cfg.GithubOrg)
		fc.ReplaceAll("simapp", cfg.ProjectName)
		fc.ReplaceAll("wasmd", cfg.BinDaemon)
	}
}

func (fc *FileContent) ReplaceDockerFile(cfg *NewChainConfig) {
	if fc.IsPath("Dockerfile") {
		fc.ReplaceAll("wasmd", cfg.BinDaemon)
	}
}

func (fc *FileContent) ReplaceApp(cfg *NewChainConfig) {
	if fc.IsPath(path.Join("app", "app.go")) {
		fc.ReplaceAll(".myapplicationd", cfg.HomeDir)
		fc.ReplaceAll(`CosmosSimApp`, cfg.ProjectName)
		fc.ReplaceAll(`mybechprefix`, cfg.Bech32Prefix)
	}
}

// ReplaceEverywhereReplaces any file content that matches anywhere in the file regardless of location.
func (fc *FileContent) ReplaceEverywhere(cfg *NewChainConfig) {
	fc.ReplaceAll("github.com/rollchains/spawn/simapp", cfg.GithubPath())

	wasmBin := path.Join("cmd", "wasmd")
	if fc.ContainsPath(wasmBin) {
		newBinPath := path.Join("cmd", cfg.BinDaemon)
		fc.NewPath = strings.ReplaceAll(fc.NewPath, wasmBin, newBinPath)
	}

	if fc.IsPath(path.Join("interchaintest", "poa_test.go")) {
		fc.FindAndReplaceAddressBech32(baseBech32Prefix, cfg.Bech32Prefix)
	}

}

func (fc *FileContent) ReplaceMakeFile(cfg *NewChainConfig) {
	bin := cfg.BinDaemon

	fc.ReplaceAll("https://github.com/rollchains/spawn/simapp.git", fmt.Sprintf("https://%s.git", cfg.GithubPath()))
	fc.ReplaceAll("version.Name=wasm", fmt.Sprintf("version.Name=%s", cfg.ProjectName)) // ldflags
	fc.ReplaceAll("version.AppName=wasmd", fmt.Sprintf("version.AppName=%s", bin))
	fc.ReplaceAll("cmd/wasmd", fmt.Sprintf("cmd/%s", bin))
	fc.ReplaceAll("build/wasmd", fmt.Sprintf("build/%s", bin))
	fc.ReplaceAll("wasmd keys", fmt.Sprintf("%s keys", bin))     // for testnet
	fc.ReplaceAll("wasmd config", fmt.Sprintf("%s config", bin)) // for local config

	fc.ReplaceAll("docker build . -t wasmd:local", fmt.Sprintf(`docker build . -t %s:local`, strings.ToLower(cfg.ProjectName)))

	fc.ReplaceAll("heighliner build -c wasmd", fmt.Sprintf(`heighliner build -c %s`, strings.ToLower(cfg.ProjectName)))
	if fc.IsPath("chains.yaml") {
		fc.ReplaceAll("myappname", strings.ToLower(cfg.ProjectName))
		fc.ReplaceAll("/go/bin/wasmd", fmt.Sprintf("/go/bin/%s", bin))
	}
}

// FindAndReplaceStandardWalletsBech32 finds a prefix1... address and replaces it with a new prefix1... address
// This works for both standard wallets (38 length after prefix1) and also smart contracts (58)
func (fc *FileContent) FindAndReplaceAddressBech32(oldPrefix, newPrefix string) {
	oldPrefix = strings.TrimSuffix(oldPrefix, "1")
	newPrefix = strings.TrimSuffix(newPrefix, "1")

	// 58 must be first to match smart contracts fully else it would only match the first 38
	// e.g. wasm10d07y265gmmuvt4z0w9aw880jnsr700js7zslc & wasm1qsrercqegvs4ye0yqg93knv73ye5dc3prqwd6jcdcuj8ggp6w0usrfxlpt
	r := regexp.MustCompile(oldPrefix + `1([0-9a-z]{58}|[0-9a-z]{38})`)

	foundAddrs := r.FindAllString(fc.Contents, -1)
	fc.Logger.Debug("Regex: Found Addresses", "addresses", foundAddrs, "path", fc.NewPath)

	for _, addr := range foundAddrs {
		_, bz, err := bech32.Decode(addr, 100)
		if err != nil {
			panic(fmt.Sprintf("error decoding bech32 address: %s. err: %s", addr, err.Error()))
		}

		newAddr, err := bech32.Encode(newPrefix, bz)
		if err != nil {
			panic(fmt.Sprintf("error encoding bech32 address: %s. err: %s", addr, err.Error()))
		}

		fc.ReplaceAll(addr, newAddr)
	}
}

// given a go mod, remove line(s) with the importPath present.
func (fc *FileContent) RemoveGoModImport(importPath string) {
	if !fc.IsPath("go.mod") && !fc.IsPath("go.sum") {
		return
	}

	if AlreadyCheckedGoModDeletion[fc.NewPath] {
		return
	}
	AlreadyCheckedGoModDeletion[fc.NewPath] = true

	fc.Logger.Debug("removing go.mod import", "path", fc.RelativePath, "import", importPath)

	lines := strings.Split(fc.Contents, "\n")

	newLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.Contains(line, importPath) {
			newLines = append(newLines, line)
		}
	}

	fc.Contents = strings.Join(newLines, "\n")

}

func (fc *FileContent) FormatGoFile() error {
	// Removes unused imports & tidies up the files
	if strings.HasSuffix(fc.NewPath, ".go") && len(fc.Contents) > 0 {
		newSrc, err := imports.Process(fc.NewPath, []byte(fc.Contents), nil)
		if err != nil {
			fc.Logger.Error("error processing imports", "err", err, "file", fc.NewPath)
			return fmt.Errorf("error processing imports: %w", err)
		}

		bz, err := format.Source(newSrc)
		if err != nil {
			fc.Logger.Error("error formatting go file", "err", err, "file", fc.NewPath)
			return fmt.Errorf("error formatting go file: %w", err)
		}

		fc.Contents = string(bz)
	}

	return nil
}

func (fc *FileContent) Save() error {
	if fc.Contents == "" {
		fc.Logger.Debug("Save() No contents for", "path", fc.NewPath)
		return nil
	}

	// 777 is used as some users have weird setup / group environments w/ MacOS. 766 is more ideal but .sh files need to be executed.
	if err := os.MkdirAll(path.Dir(fc.NewPath), 0777); err != nil {
		return err
	}

	return os.WriteFile(fc.NewPath, []byte(fc.Contents), 0777)
}
