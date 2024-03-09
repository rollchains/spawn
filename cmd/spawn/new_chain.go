package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"gitub.com/rollchains/spawn/spawn"
)

var (
	// SupportedFeatures is a list of all features that can be toggled.
	// - UI: uses the IsSelected
	// - CLI: all are enabled by default. Must opt out.
	SupportedFeatures = items{
		{ID: "proof-of-authority", IsSelected: true, Details: "Proof-of-Authority consensus algorithm (permissioned network)"},
		{ID: "tokenfactory", IsSelected: true, Details: "Native token minting, sending, and burning on the chain"},
		{ID: "globalfee", IsSelected: true, Details: "Static minimum fee(s) for all transactions, controlled by governance"},
		{ID: "ibc-packetforward", IsSelected: true, Details: "Packet forwarding (for IBC)"},
		{ID: "cosmwasm", IsSelected: false, Details: "Cosmos smart contracts"},
		{ID: "wasm-light-client", IsSelected: false, Details: "08 Wasm Light Client"},
		{ID: "ignite-cli", IsSelected: false, Details: "Ignite-CLI Support"},
	}

	// parentDeps is a map of module name to a list of module names that are disabled if the parent is disabled
	parentDeps = map[string][]string{
		// "cosmwasm": {spawn.AliasName("some-wasmd-child-feature")},
	}
)

const (
	FlagWalletPrefix = "bech32"
	FlagBinDaemon    = "bin"
	FlagDebugging    = "debug"
	FlagTokenDenom   = "denom"
	FlagGithubOrg    = "org"
	FlagDisabled     = "disable"
	FlagNoGit        = "skip-git"
	FlagBypassPrompt = "bypass-prompt"
)

func init() {
	features := make([]string, len(SupportedFeatures))
	for idx, feat := range SupportedFeatures {
		features[idx] = feat.ID
	}

	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().StringP(FlagBinDaemon, "b", "simd", "binary name")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
	newChain.Flags().String(FlagTokenDenom, "token", "bank token denomination")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, strings.Join(features, ","))
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
	newChain.Flags().Bool(FlagNoGit, false, "ignore git init")
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
		logger := GetLogger()

		projName := strings.ToLower(args[0])
		homeDir := "." + projName

		disabled, _ := cmd.Flags().GetStringSlice(FlagDisabled)
		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinDaemon)
		denom, _ := cmd.Flags().GetString(FlagTokenDenom)
		ignoreGitInit, _ := cmd.Flags().GetBool(FlagNoGit)
		githubOrg, _ := cmd.Flags().GetString(FlagGithubOrg)

		// Show a UI if the user did not specific to bypass it, or if nothing is disabled.
		bypassPrompt, _ := cmd.Flags().GetBool(FlagBypassPrompt)
		if len(disabled) == 0 && !bypassPrompt {
			items, err := selectItems(0, SupportedFeatures, true)
			if err != nil {
				logger.Error("Error selecting disabled", "err", err)
				return
			}
			disabled = items.NOTSlice()
		}

		disabled = CleanDisabled(disabled)

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

// CleanDisabled normalizes the names, removes any parent dependencies, and removes duplicates
func CleanDisabled(disabled []string) []string {
	for i, name := range disabled {
		// normalize disabled to standard aliases
		alias := spawn.AliasName(name)
		disabled[i] = alias

		// if we disable a feature which has disabled dependency, we need to disable those too
		if deps, ok := parentDeps[alias]; ok {
			// duplicates will arise, removed in the next step
			disabled = append(disabled, deps...)
		}
	}

	// remove duplicates
	dups := make(map[string]bool)
	for _, d := range disabled {
		dups[d] = true
	}

	disabled = []string{}
	for d := range dups {
		disabled = append(disabled, d)
	}

	return disabled
}

func normalizeWhitelistVarRun(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "binary":
		name = FlagBinDaemon
	case "disabled", "remove":
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
