package main

import (
	"strings"

	localictypes "github.com/strangelove-ventures/interchaintest/local-interchain/interchain/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
)

func main() {
	// ProjectName:myproject Bech32Prefix:roll HomeDir:.myproject BinDaemon:rolld Denom:uroll GithubOrg:rollchains IgnoreGitInit:true DisabledModules:[poa staking] Logger:0xc001264330 isUsingICS:false

	c := localictypes.NewChainBuilder("myproject", "localchain-1", "rolld", "uroll", "roll").
		SetBlockTime("2000ms").
		SetDockerImage(ibc.NewDockerImage(strings.ToLower("myproject"), "local", "")).
		SetTrustingPeriod("336h").
		SetHostPortOverride(localictypes.BaseHostPortOverride()).
		SetDefaultSDKv47Genesis(2).
		SetICSConsumerLink("localethos-1")

	c.ICSVersionOverride = ibc.ICSConfig{
		ProviderVerOverride: "v5", // no migration needed since both are the same
		ConsumerVerOverride: "v5",
	}

	c.Genesis.Modify = []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", c.Denom),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.amount", "1"),
	}

	ethos := localictypes.NewChainBuilder("ethos", "localethos-1", "ethosd", "uethos", "cosmos").
		SetDebugging(true).
		SetBech32Prefix("ethos").
		SetDockerImage(ibc.NewDockerImage("ethos", "local", "1025:1025")).
		SetBlockTime("2000ms").
		SetDefaultSDKv47Genesis(2)

	ethos.Genesis.Modify = []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", "10s"),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", ethos.Denom),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.amount", "1"),
	}

	eth := localictypes.ChainEthereum()
	eth.SetConfigFileOverrides([]localictypes.ConfigFileOverrides{
		{
			Paths: testutil.Toml{
				"--load-state": "chains/avs-and-eigenlayer-deployed-anvil-state.json",
			},
		},
	})

	localictypes.NewChainsConfig(ethos, eth, c).SaveJSON("chains/ethos-ics.json")
}
