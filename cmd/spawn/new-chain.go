package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"gitub.com/strangelove-ventures/spawn/spawn"
)

var (
	SupportedModules = items{
		{ID: "tokenfactory", IsSelected: true, Details: "Native token minting, sending, and burning on the chain"},
		{ID: "poa", IsSelected: true, Details: "Proof-of-Authority consensus algorithm (permissioned network)"},
		{ID: "globalfee", IsSelected: true, Details: "Static minimum fee(s) for all transactions, controlled by governance"},
		{ID: "cosmwasm", IsSelected: true, Details: "Cosmos smart contracts"},
	}
)

const (
	FlagWalletPrefix = "bech32"
	FlagBinDaemon    = "bin"
	FlagDebugging    = "debug"
	FlagTokenDenom   = "denom"
	FlagGithubOrg    = "org"
	FlagDisabled     = "disable"
	FlagNoGit        = "no-git"
	FlagBypassPrompt = "bypass-prompt"
)

func init() {
	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().StringP(FlagBinDaemon, "b", "simd", "binary name")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
	newChain.Flags().String(FlagTokenDenom, "token", "bank token denomination")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, "disable features: "+SupportedModules.String())
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
	newChain.Flags().Bool(FlagNoGit, false, "git init base")
	newChain.Flags().Bool(FlagBypassPrompt, false, "bypass UI prompt")

	newChain.Flags().SetNormalizeFunc(normalizeWhitelistVarRun)
}

var newChain = &cobra.Command{
	Use:   "new-chain [project-name]",
	Short: "Create a new project",
	Example: fmt.Sprintf(
		`spawn new rollchain --%s=cosmos --%s=simd --%s=token --%s=tokenfactory,poa,globalfee`,
		FlagWalletPrefix, FlagBinDaemon, FlagTokenDenom, FlagDisabled,
	),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"new", "init", "create"},
	Run: func(cmd *cobra.Command, args []string) {
		projName := strings.ToLower(args[0])
		homeDir := "." + projName

		disabled, _ := cmd.Flags().GetStringSlice(FlagDisabled)
		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinDaemon)
		denom, _ := cmd.Flags().GetString(FlagTokenDenom)
		debug, _ := cmd.Flags().GetBool(FlagDebugging)
		ignoreGitInit, _ := cmd.Flags().GetBool(FlagNoGit)
		githubOrg, _ := cmd.Flags().GetString(FlagGithubOrg)

		bypassPrompt, _ := cmd.Flags().GetBool(FlagBypassPrompt)
		if len(disabled) == 0 && !bypassPrompt {
			items, err := selectItems(0, SupportedModules, true)
			if err != nil {
				fmt.Println("Error selecting disabled:", err)
				return
			}
			disabled = items.NOTSlice()
		}

		cfg := &spawn.NewChainConfig{
			ProjectName:     projName,
			Bech32Prefix:    walletPrefix,
			HomeDir:         homeDir,
			BinDaemon:       binName,
			Denom:           denom,
			Debug:           debug,
			GithubOrg:       githubOrg,
			IgnoreGitInit:   ignoreGitInit,
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
	case "bypass", "skip", "force", "prompt-bypass", "bypass-ui", "no-ui":
		name = FlagBypassPrompt
	}

	return pflag.NormalizedName(name)
}
