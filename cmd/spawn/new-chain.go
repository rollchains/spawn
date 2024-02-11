package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"gitub.com/strangelove-ventures/spawn/spawn"
)

var (
	SupportedFeatures = []string{"tokenfactory", "poa", "globalfee", "wasm", "icahost", "icacontroller"}
)

const (
	FlagWalletPrefix = "bech32"
	FlagBinDaemon    = "bin"
	FlagDebugging    = "debug"
	FlagTokenDenom   = "denom"
	FlagGithubOrg    = "org"
	FlagDisabled     = "disable"
	FlagNoGit        = "no-git"
)

func init() {
	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().StringP(FlagBinDaemon, "b", "appd", "binary name")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
	newChain.Flags().String(FlagTokenDenom, "stake", "token denom")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, "disable features: "+strings.Join(SupportedFeatures, ", "))
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
	newChain.Flags().Bool(FlagNoGit, false, "git init base")

	newChain.Flags().SetNormalizeFunc(normalizeWhitelistVarRun)
}

var newChain = &cobra.Command{
	Use:   "new-chain [project-name]",
	Short: "Create a new project",
	Example: fmt.Sprintf(
		`spawn new rollchain --%s=cosmos --%s=appd --%s=token --%s=tokenfactory,poa,globalfee`,
		FlagWalletPrefix, FlagBinDaemon, FlagTokenDenom, FlagDisabled,
	),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"new", "init"},
	Run: func(cmd *cobra.Command, args []string) {
		projName := strings.ToLower(args[0])
		homeDir := "." + projName

		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinDaemon)
		denom, _ := cmd.Flags().GetString(FlagTokenDenom)
		debug, _ := cmd.Flags().GetBool(FlagDebugging)
		disabled, _ := cmd.Flags().GetStringSlice(FlagDisabled)
		ignoreGitInit, _ := cmd.Flags().GetBool(FlagNoGit)
		githubOrg, _ := cmd.Flags().GetString(FlagGithubOrg)

		cfg := &spawn.NewChainConfig{
			ProjectName:  projName,
			Bech32Prefix: walletPrefix,
			HomeDir:      homeDir,
			BinDaemon:    binName,
			Denom:        denom,
			Debug:        debug,
			GithubOrg:    githubOrg,

			IgnoreGitInit: ignoreGitInit,

			// by default everything is on, then we remove what the user wants to disable
			DisabledModules: disabled,
		}
		if err := cfg.Validate(); err != nil {
			fmt.Println("Error validating config:", err)
			return
		}

		cfg.NewChain()
		cfg.AnnounceSuccessfulBuild()
	},
}

func normalizeWhitelistVarRun(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "binary":
		name = FlagBinDaemon
	case "disabled":
		name = FlagDisabled
	}

	return pflag.NormalizedName(name)
}
