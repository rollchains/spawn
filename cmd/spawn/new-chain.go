package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"

	"gitub.com/strangelove-ventures/spawn/spawn"
)

var (
	SupportedFeatures = []string{"tokenfactory", "poa", "globalfee", "wasm", "icahost", "icacontroller"}
)

const (
	FlagWalletPrefix = "bech32"
	FlagBinaryName   = "bin"
	FlagDebugging    = "debug"
	FlagTokenDenom   = "denom"
	FlagGithubOrg    = "org"
	FlagDisabled     = "disable"
	FlagNoGit        = "no-git"
)

func init() {
	newChain.Flags().String(FlagWalletPrefix, "cosmos", "chain wallet bech32 prefix")
	newChain.Flags().String(FlagBinaryName, "appd", "binary name")
	newChain.Flags().Bool(FlagDebugging, false, "enable debugging")
	newChain.Flags().StringSlice(FlagDisabled, []string{}, "disable features: "+strings.Join(SupportedFeatures, ", "))
	newChain.Flags().String(FlagTokenDenom, "stake", "token denom")
	newChain.Flags().Bool(FlagNoGit, false, "git init base")
	newChain.Flags().String(FlagGithubOrg, "rollchains", "github organization")
}

var newChain = &cobra.Command{
	Use:   "new-chain [project-name]",
	Short: "Create a new project",
	Example: fmt.Sprintf(
		`spawn new rollchain --%s=cosmos --%s=appd --%s=token --%s=tokenfactory,poa,globalfee`,
		FlagWalletPrefix, FlagBinaryName, FlagTokenDenom, FlagDisabled,
	),
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"new", "init"},
	Run: func(cmd *cobra.Command, args []string) {
		projName := strings.ToLower(args[0])
		appName := cases.Title(language.AmericanEnglish).String(projName) + "App"

		walletPrefix, _ := cmd.Flags().GetString(FlagWalletPrefix)
		binName, _ := cmd.Flags().GetString(FlagBinaryName)
		denom, _ := cmd.Flags().GetString(FlagTokenDenom)
		debug, _ := cmd.Flags().GetBool(FlagDebugging)
		disabled, _ := cmd.Flags().GetStringSlice(FlagDisabled)
		ignoreGitInit, _ := cmd.Flags().GetBool(FlagNoGit)
		githubOrg, _ := cmd.Flags().GetString(FlagGithubOrg)

		cfg := &spawn.NewChainConfig{
			ProjectName:  projName,
			Bech32Prefix: walletPrefix,
			AppName:      appName,
			AppDirName:   "." + projName,
			BinaryName:   binName,
			TokenDenom:   denom,
			Debugging:    debug,
			GithubOrg:    githubOrg,

			GitInitOnCreate: !ignoreGitInit,

			// by default everything is on, then we remove what the user wants to disable
			DisabledFeatures: disabled,
		}
		if err := cfg.Validate(); err != nil {
			fmt.Println("Error validating config:", err)
			return
		}

		cfg.NewChain()
		cfg.AnnounceSuccessfulBuild()
	},
}
