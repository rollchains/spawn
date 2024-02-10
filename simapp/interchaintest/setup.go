package e2e

import (
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types" //spawntag:globalfee

	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	globalfee "github.com/reecepbcups/globalfee/x/globalfee/types"
	tokenfactory "github.com/reecepbcups/tokenfactory/x/tokenfactory/types"
)

var (
	VotingPeriod     = "15s"
	MaxDepositPeriod = "10s"

	Denom   = "mydenom"
	Name    = "appName"
	ChainID = "chainid-1"
	Binary  = "wasmd"

	Bech32 = "wasm"

	NumberVals         = 1
	NumberFullNodes    = 0
	GenesisFundsAmount = sdkmath.NewInt(1000_000000) // 1k tokens

	ChainImage = ibc.NewDockerImage("wasmd", "local", "1025:1025")

	GasCoin = sdk.NewDecCoinFromDec(Denom, sdkmath.LegacyMustNewDecFromStr("0.0")) // spawntag:globalfee

	DefaultGenesis = []cosmos.GenesisKV{
		// default
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", VotingPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", MaxDepositPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", Denom),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.amount", "1"),
		// poa: gov & testing account
		cosmos.NewGenesisKV("app_state.poa.params.admins", []string{"wasm10d07y265gmmuvt4z0w9aw880jnsr700js7zslc", "wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g"}),
		// globalfee: set minimum fee requirements
		cosmos.NewGenesisKV("app_state.globalfee.params.minimum_gas_prices", sdk.DecCoins{GasCoin}),
		// tokenfactory: set create cost in set denom or in gas usage.
		cosmos.NewGenesisKV("app_state.tokenfactory.params.denom_creation_fee", nil),
		cosmos.NewGenesisKV("app_state.tokenfactory.params.denom_creation_gas_consume", 1), // cost 1 gas to create a new denom

	}

	DefaultChainConfig = ibc.ChainConfig{
		Images: []ibc.DockerImage{
			ChainImage,
		},
		GasAdjustment: 1.5,
		ModifyGenesis: cosmos.ModifyGenesis(DefaultGenesis),
		EncodingConfig: func() *moduletestutil.TestEncodingConfig {
			cfg := cosmos.DefaultEncoding()
			// TODO: add encoding types here for the modules you want to use
			wasm.RegisterInterfaces(cfg.InterfaceRegistry)
			tokenfactory.RegisterInterfaces(cfg.InterfaceRegistry)
			globalfee.RegisterInterfaces(cfg.InterfaceRegistry)
			return &cfg
		}(),
		Type:           "cosmos",
		Name:           Name,
		ChainID:        ChainID,
		Bin:            Binary,
		Bech32Prefix:   Bech32,
		Denom:          Denom,
		CoinType:       "118",
		GasPrices:      "0" + Denom,
		TrustingPeriod: "336h",
	}

	DefaultChainSpec = interchaintest.ChainSpec{
		Name:          Name,
		ChainName:     Name,
		Version:       ChainImage.Version,
		ChainConfig:   DefaultChainConfig,
		NumValidators: &NumberVals,
		NumFullNodes:  &NumberFullNodes,
	}
)
