package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/rollchains/spawn/spawn"
)

var (
	// SupportedFeatures is a list of all features that can be toggled.
	// - UI: uses the IsSelected
	// - CLI: all are enabled by default. Must opt out.

	ConsensusFeatures = items{
		// Consensus (only 1 per app)
		{ID: "proof-of-authority", IsSelected: true, IsConsensus: true, Details: "Proof-of-Authority consensus algorithm (permissioned network)"},
		{ID: "proof-of-stake", IsSelected: false, IsConsensus: true, Details: "Proof-of-Stake consensus algorithm (permissionless network)"},
		{ID: "interchain-security", IsSelected: false, IsConsensus: true, Details: "Cosmos Hub Interchain Security"},
		// {ID: "ics-ethos", IsSelected: false, IsConsensus: true, Details: "Interchain-Security with Ethos ETH restaking"},
	}

	SupportedFeatures = append(ConsensusFeatures, items{
		{ID: "tokenfactory", IsSelected: true, Details: "Native token minting, sending, and burning on the chain"},
		{ID: "ibc-packetforward", IsSelected: true, Details: "Packet forwarding"},
		{ID: "ibc-ratelimit", IsSelected: false, Details: "Thresholds for outflow as a percent of total channel value"},
		{ID: "cosmwasm", IsSelected: false, Details: "Cosmos smart contracts"},
		{ID: "wasm-light-client", IsSelected: false, Details: "08 Wasm Light Client"},
		{ID: "optimistic-execution", IsSelected: true, Details: "Pre-process blocks ahead of consensus request"},
		{ID: "block-explorer", IsSelected: false, Details: "Ping Pub Explorer"},
	}...)

	// parentDeps is a list of modules that are disabled if a parent module is disabled.
	// i.e. Without staking, POA is not possible as it depends on staking.
	parentDeps = map[string][]string{}
)

const (
	FlagWalletPrefix = "wallet-prefix"
	FlagBinDaemon    = "binary"
	FlagDebugging    = "debug"
	FlagTokenDenom   = "denom"
	FlagGithubOrg    = "org"
	FlagDisabled     = "disable"
	FlagConsensus    = "consensus"
	FlagNoGit        = "skip-git"
	FlagBypassPrompt = "bypass-prompt"
)

func init() {
	features := make([]string, 0)
	consensus := make([]string, 0)

	for _, feat := range SupportedFeatures {
		if feat.IsConsensus {
			consensus = append(consensus, feat.ID)
		} else {
			features = append(features, feat.ID)
		}
	}

	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain bech32 wallet prefix")
	newChain.Flags().StringP(FlagBinDaemon, "b", "simd", "binary name")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
	newChain.Flags().String(FlagTokenDenom, "token", "bank token denomination")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, strings.Join(features, ","))
	newChain.Flags().String(FlagConsensus, "", strings.Join(consensus, ",")) // must be set to nothing is nothing is set
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
	newChain.Flags().Bool(FlagNoGit, false, "ignore git init")
	newChain.Flags().Bool(FlagBypassPrompt, false, "bypass UI prompt")
	newChain.Flags().SetNormalizeFunc(normalizeWhitelistVarRun)
}

var newChain = &cobra.Command{
	Use:   "new-chain [project-name]",
	Short: "Create a new project",
	Example: fmt.Sprintf(
		`  - spawn new rollchain --consensus=proof-of-stake --%s=cosmos --%s=simd --%s=token --org=abcde
  - spawn new rollchain --consensus=proof-of-authority --%s=tokenfactory
  - spawn new rollchain --consensus=interchain-security --%s=cosmwasm --%s
  - spawn new rollchain --%s`,
		FlagWalletPrefix, FlagBinDaemon, FlagTokenDenom, FlagDisabled, FlagDisabled, FlagNoGit, FlagBypassPrompt,
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
		consensus, _ := cmd.Flags().GetString(FlagConsensus)

		bypassPrompt, _ := cmd.Flags().GetBool(FlagBypassPrompt)

		// Show a UI to select the consensus algorithm (POS, POA, ICS) if a custom one was not specified.
		if !bypassPrompt {
			if len(consensus) == 0 {
				text := "Consensus Selector (( enter to toggle ))"
				items, err := selectItems(text, 0, SupportedFeatures, false, true, true)
				if err != nil {
					logger.Error("Error selecting consensus", "err", err)
					return
				}
				consensus = items.String()
			}
		}
		if len(consensus) == 0 {
			consensus = spawn.POA // set the default if still none is provided
		}

		consensus = spawn.AliasName(consensus)
		logger.Debug("Consensus selected", "consensus", consensus)

		// Disable all consensus algorithms except the one selected.
		disabledConsensus := make([]string, 0)
		for _, feat := range ConsensusFeatures {
			name := spawn.AliasName(feat.ID)
			if name != consensus {
				// if consensus is proof-of-authority, allow proof of stake
				if consensus == spawn.POA && name == spawn.POS {
					continue
				} else if consensus == spawn.InterchainSecurity && name == spawn.POS {
					continue
				}

				disabledConsensus = append(disabledConsensus, name)
			}
		}
		logger.Debug("Disabled Consensuses", "disabled", disabledConsensus, "using", consensus)

		// Disable all features not selected
		// Show a UI if the user did not specific to bypass it, or if nothing is disabled.
		if len(disabled) == 0 && !bypassPrompt {
			text := "Feature Selector (( enter to toggle ))"
			items, err := selectItems(text, 0, SupportedFeatures, true, false, false)
			if err != nil {
				logger.Error("Error selecting disabled", "err", err)
				return
			}
			disabled = items.NOTSlice()

		}

		disabled = append(disabled, disabledConsensus...)
		disabled = spawn.NormalizeDisabledNames(disabled, parentDeps)

		logger.Debug("Disabled features final", "features", disabled)

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

		if err := cfg.ValidateAndRun(true); err != nil {
			logger.Error("Error creating new chain", "err", err)
			return
		}
	},
}

func normalizeWhitelistVarRun(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "bin", "daemon":
		name = FlagBinDaemon
	case "disabled", "remove":
		name = FlagDisabled
	case "bypass", "skip", "force", "prompt-bypass", "bypass-ui", "no-ui":
		name = FlagBypassPrompt
	case "token", "denomination", "coin":
		name = FlagTokenDenom
	case "no-git", "ignore-git":
		name = FlagNoGit
	case "bech32", "prefix", "wallet":
		name = FlagWalletPrefix
	case "organization", "namespace":
		name = FlagGithubOrg
	}

	return pflag.NormalizedName(name)
}
