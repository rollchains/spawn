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
		// 08wasm light client depends on cosmwasm
		"cosmwasm": {spawn.AliasName("wasm-light-client")},
	}
)

const (
	FlagWalletPrefix   = "bech32"
	FlagBinDaemon      = "bin"
	FlagDebugging      = "debug"
	FlagTokenDenom     = "denom"
	FlagGithubOrg      = "org"
	FlagDisabled       = "disable"
	FlagEnabled        = "enable"
	FlagNoGit          = "skip-git"
	FlagBypassPrompt   = "bypass-prompt"
	FlagIgniteCLIOptIn = "ignite-cli"
)

func init() {

	defaultOffFeatures := []string{}
	defaultOnFeatures := []string{}
	for _, feat := range SupportedFeatures {
		if !feat.IsSelected {
			defaultOffFeatures = append(defaultOffFeatures, feat.ID)
		} else {
			defaultOnFeatures = append(defaultOnFeatures, feat.ID)
		}
	}

	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().StringP(FlagBinDaemon, "b", "simd", "binary name")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
	newChain.Flags().String(FlagTokenDenom, "token", "bank token denomination")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, "disable features: "+strings.Join(defaultOnFeatures, ","))
	newChain.Flags().StringSlice(FlagEnabled, []string{}, "enable: "+strings.Join(defaultOffFeatures, ","))
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
		enabled, _ := cmd.Flags().GetStringSlice(FlagEnabled)
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

		// normalize disabled to standard aliases
		for i, name := range disabled {
			disabled[i] = spawn.AliasName(name)
		}

		// iterate through disabled and remove any which are in enabled
		for _, name := range disabled {
			for i, enabled := range enabled {
				if name == enabled {
					disabled = append(disabled[:i], disabled[i+1:]...)
				}
			}
		}

		// if we disable a feature which has disabled dependency, we need to disable those too
		for _, name := range disabled {
			if deps, ok := dependencies[name]; ok {
				disabled = append(disabled, deps...)
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
			EnabledModules:  enabled,
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
