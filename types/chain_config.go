package spawntypes

import (
	"fmt"
	"strings"
)

type NewChainConfig struct {
	ProjectName     string
	Bech32Prefix    string
	AppName         string
	AppDirName      string
	BinaryName      string
	TokenDenom      string
	GithubOrg       string
	GitInitOnCreate bool

	Debugging bool

	DisabledFeatures []string
}

func (cfg *NewChainConfig) Validate() error {
	if strings.ContainsAny(cfg.ProjectName, `~!@#$%^&*()_+{}|:"<>?/.,;'[]\=-`) {
		return fmt.Errorf("project name cannot contain special characters %s", cfg.ProjectName)
	}

	return nil
}

func (cfg *NewChainConfig) AnnounceSuccessfulBuild() {
	projName := cfg.ProjectName
	binName := cfg.BinaryName

	fmt.Printf("\n\n🎉 New blockchain '%s' generated!\n", projName)
	fmt.Println("🏅Getting started:")
	fmt.Println("  - $ cd " + projName)
	fmt.Println("  - $ make testnet      # build & start a testnet")
	fmt.Println("  - $ make testnet-ibc  # build & start an ibc testnet")
	fmt.Printf("  - $ make install      # build the %s binary\n", binName)
	fmt.Println("  - $ make local-image  # build docker image")
}
