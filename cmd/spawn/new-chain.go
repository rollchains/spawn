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
		{ID: "ibc-packetforward", IsSelected: true, Details: "Packet forwarding (for IBC)"},
		{ID: "cosmwasm", IsSelected: false, Details: "Cosmos smart contracts"},
		{ID: "wasm-light-client", IsSelected: false, Details: "08 Wasm Light Client"},
		{ID: "ignite-cli", IsSelected: false, Details: "Ignite-CLI Support"},
	}

	dependencies = map[string][]string{
		// 08wasm light client depends on cosmwasm (or not?)
		// "cosmwasm": {spawn.AliasName("wasm-light-client")},
	}
)

const (
	FlagWalletPrefix = "bech32"
	FlagBinDaemon    = "bin"
	FlagDebugging    = "debug"
	FlagTokenDenom   = "denom"
	FlagGithubOrg    = "org"
	FlagDisabled     = "disable"
	FlagEnabled      = "enable"
	FlagNoGit        = "skip-git"
	FlagBypassPrompt = "bypass-prompt"
)

var (
	disabledByDefault = []string{}
)

func init() {
	showcaseOnFeatures := []string{}
	for _, feat := range SupportedFeatures {
		if !feat.IsSelected {
			disabledByDefault = append(disabledByDefault, feat.ID)
		} else {
			showcaseOnFeatures = append(showcaseOnFeatures, feat.ID)
		}
	}

	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().StringP(FlagBinDaemon, "b", "simd", "binary name")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
	newChain.Flags().String(FlagTokenDenom, "token", "bank token denomination")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, "disable: "+strings.Join(showcaseOnFeatures, ","))
	newChain.Flags().StringSlice(FlagEnabled, []string{}, "enable : "+strings.Join(disabledByDefault, ","))
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
		enabled, _ := cmd.Flags().GetStringSlice(FlagEnabled)

		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinDaemon)
		denom, _ := cmd.Flags().GetString(FlagTokenDenom)
		ignoreGitInit, _ := cmd.Flags().GetBool(FlagNoGit)
		githubOrg, _ := cmd.Flags().GetString(FlagGithubOrg)

		// Show a UI if the user did not specific to bypass it, or if nothing is disabled. So they get to see what is to be picked from.
		bypassPrompt, _ := cmd.Flags().GetBool(FlagBypassPrompt)
		useUI := len(disabled) == 0 && len(enabled) == 0 && !bypassPrompt
		if useUI {
			items, err := selectItems(0, SupportedFeatures, true)
			if err != nil {
				logger.Error("Error selecting disabled", "err", err)
				return
			}
			disabled = items.NOTSlice()
		} else {
			// Auto disable features that are not off by default (duplicates are fine)
			// This is not done for the UI since that is set by the user for all directly.
			disabled = append(disabled, disabledByDefault...)
		}

		for i, name := range disabled {
			// normalize disabled to standard aliases
			alias := spawn.AliasName(name)
			disabled[i] = alias

			// if we disable a feature which has disabled dependency, we need to disable those too
			if deps, ok := dependencies[alias]; ok {
				// duplicates will arise, will be removed layer
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

		// 2) If a feature is enabled, remove it from the disabled slice
		for _, name := range enabled {
			alias := spawn.AliasName(name)

			for i, d := range disabled {
				if d == alias {
					// the user could have disabled a feature and enabled too.
					// if so, we go with disabled for it (thus not breaking early)
					disabled = append(disabled[:i], disabled[i+1:]...)
				}
			}
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
	case "enabled":
		name = FlagEnabled
	case "bypass", "skip", "force", "prompt-bypass", "bypass-ui", "no-ui":
		name = FlagBypassPrompt
	case "token", "denomination", "coin":
		name = FlagTokenDenom
	case "no-git", "ignore-git":
		name = FlagNoGit
	}

	return pflag.NormalizedName(name)
}
