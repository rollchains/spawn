package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"gitub.com/rollchains/spawn/spawn"
)

var (
	SupportedFeatures = items{
		{ID: "proof-of-authority", IsSelected: true, Details: "Proof-of-Authority consensus algorithm (permissioned network)"},
		{ID: "tokenfactory", IsSelected: true, Details: "Native token minting, sending, and burning on the chain"},
		{ID: "globalfee", IsSelected: true, Details: "Static minimum fee(s) for all transactions, controlled by governance"},
		{ID: "cosmwasm", IsSelected: false, Details: "Cosmos smart contracts"},
		{ID: "ignite-cli", IsSelected: false, Details: "Ignite-CLI Support"},
		{ID: "ibc-packetforward", IsSelected: true, Details: "Packet forwarding (for IBC)"},
	}
)

const (
	FlagWalletPrefix   = "bech32"
	FlagBinDaemon      = "bin"
	FlagDebugging      = "debug"
	FlagTokenDenom     = "denom"
	FlagGithubOrg      = "org"
	FlagDisabled       = "disable"
	FlagNoGit          = "skip-git"
	FlagBypassPrompt   = "bypass-prompt"
	FlagIgniteCLIOptIn = "ignite-cli"
)

func init() {
	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().StringP(FlagBinDaemon, "b", "simd", "binary name")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
	newChain.Flags().String(FlagTokenDenom, "token", "bank token denomination")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, "disable features: "+SupportedFeatures.String())
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
	newChain.Flags().Bool(FlagNoGit, false, "ignore git init")
	newChain.Flags().Bool(FlagBypassPrompt, false, "bypass UI prompt")
	newChain.Flags().Bool(FlagIgniteCLIOptIn, false, "opt-in to Ignite CLI support")

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
		logger := GetLogger()

		projName := strings.ToLower(args[0])
		homeDir := "." + projName

		disabled, _ := cmd.Flags().GetStringSlice(FlagDisabled)
		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinDaemon)
		denom, _ := cmd.Flags().GetString(FlagTokenDenom)
		ignoreGitInit, _ := cmd.Flags().GetBool(FlagNoGit)
		githubOrg, _ := cmd.Flags().GetString(FlagGithubOrg)
		useIgnite, _ := cmd.Flags().GetBool(FlagIgniteCLIOptIn)

		if useIgnite {
			for _, feat := range SupportedFeatures {
				if strings.Contains(feat.ID, "ignite") {
					feat.IsSelected = true
				}
			}
		}

		bypassPrompt, _ := cmd.Flags().GetBool(FlagBypassPrompt)
		if len(disabled) == 0 && !bypassPrompt {
			items, err := selectItems(0, SupportedFeatures, true)
			if err != nil {
				logger.Error("Error selecting disabled", "err", err)
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
			GithubOrg:       githubOrg,
			IgnoreGitInit:   ignoreGitInit,
			DisabledModules: disabled,
			Logger:          logger,
		}
		if err := cfg.Validate(); err != nil {
			logger.Error("Error validating config", "err", err)
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
	case "token", "denomination", "coin":
		name = FlagTokenDenom
	case "no-git", "ignore-git":
		name = FlagNoGit
	}

	return pflag.NormalizedName(name)
}
