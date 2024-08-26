package e2e

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	ibcconntypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
	poa "github.com/strangelove-ventures/poa"
	tokenfactory "github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"
)

var (
	VotingPeriod     = "15s"
	MaxDepositPeriod = "10s"

	Denom   = "mydenom"
	Name    = "appName"
	ChainID = "localchain-1"
	Binary  = "wasmd"
	Bech32  = "mybechprefix"

	NumberVals         = 1
	NumberFullNodes    = 0
	GenesisFundsAmount = sdkmath.NewInt(1000_000000) // 1k tokens

	ChainImage = ibc.NewDockerImage("wasmd", "local", "1025:1025")

	DefaultGenesis = []cosmos.GenesisKV{
		// default
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", VotingPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", MaxDepositPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.denom", Denom),
		cosmos.NewGenesisKV("app_state.gov.params.min_deposit.0.amount", "1"),
		// tokenfactory: set create cost in set denom or in gas usage.
		cosmos.NewGenesisKV("app_state.tokenfactory.params.denom_creation_fee", nil),
		cosmos.NewGenesisKV("app_state.tokenfactory.params.denom_creation_gas_consume", 1), // cost 1 gas to create a new denom

	}

	DefaultChainConfig = ibc.ChainConfig{
		Images: []ibc.DockerImage{
			ChainImage,
		},
		GasAdjustment:  1.5,
		ModifyGenesis:  cosmos.ModifyGenesis(DefaultGenesis), // spawntag:staking
		EncodingConfig: GetEncodingConfig(),
		Type:           "cosmos",
		Name:           Name,
		ChainID:        ChainID,
		Bin:            Binary,
		Bech32Prefix:   Bech32,
		Denom:          Denom,
		CoinType:       "118",
		GasPrices:      "0" + Denom,
		TrustingPeriod: "504h",
	}

	DefaultChainSpec = interchaintest.ChainSpec{
		Name:          Name,
		ChainName:     Name,
		Version:       ChainImage.Version,
		ChainConfig:   DefaultChainConfig,
		NumValidators: &NumberVals,
		NumFullNodes:  &NumberFullNodes,
	}

	SecondDefaultChainSpec = func() interchaintest.ChainSpec {
		SecondChainSpec := DefaultChainSpec
		SecondChainSpec.ChainID += "2"
		SecondChainSpec.Name += "2"
		SecondChainSpec.ChainName += "2"
		return SecondChainSpec
	}()

	// cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr - test_node.sh
	AccMnemonic  = "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry"
	Acc1Mnemonic = "wealth flavor believe regret funny network recall kiss grape useless pepper cram hint member few certain unveil rather brick bargain curious require crowd raise"

	RelayerRepo    = "ghcr.io/cosmos/relayer"
	RelayerVersion = "main"

	vals   = 1
	fNodes = 0

	// <spawntag:ics
	ProviderVer   = "v15.2.0"
	ProviderChain = interchaintest.ChainSpec{
		Name: "gaia", Version: ProviderVer,
		NumValidators: &vals, NumFullNodes: &fNodes,
		ChainConfig: ibc.ChainConfig{
			GasAdjustment:  3.0,
			TrustingPeriod: "504h",
			ModifyGenesis:  cosmos.ModifyGenesis([]cosmos.GenesisKV{
				// v4+ ICS provider required
				// cosmos.NewGenesisKV("app_state.provider.params.blocks_per_epoch", "1"),
			}),
		},
	}
	// spawntag:ics>
)

func GetEncodingConfig() *moduletestutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()
	// TODO: add encoding types here for the modules you want to use
	wasm.RegisterInterfaces(cfg.InterfaceRegistry)
	tokenfactory.RegisterInterfaces(cfg.InterfaceRegistry)
	poa.RegisterInterfaces(cfg.InterfaceRegistry)
	ccvconsumertypes.RegisterInterfaces(cfg.InterfaceRegistry)
	return &cfg
}

// Other Helpers
func ExecuteQuery(ctx context.Context, chain *cosmos.CosmosChain, cmd []string, i interface{}, extraFlags ...string) {
	flags := []string{
		"--node", chain.GetRPCAddress(),
		"--output=json",
	}
	flags = append(flags, extraFlags...)

	ExecuteExec(ctx, chain, cmd, i, flags...)
}
func ExecuteExec(ctx context.Context, chain *cosmos.CosmosChain, cmd []string, i interface{}, extraFlags ...string) {
	command := []string{chain.Config().Bin}
	command = append(command, cmd...)
	command = append(command, extraFlags...)
	fmt.Println(command)

	stdout, _, err := chain.Exec(ctx, command, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(stdout))
	if err := json.Unmarshal(stdout, &i); err != nil {
		fmt.Println(err)
	}
}

// Executes a command from CommandBuilder
func ExecuteTransaction(ctx context.Context, chain *cosmos.CosmosChain, cmd []string) (sdk.TxResponse, error) {
	var err error
	var stdout []byte

	stdout, _, err = chain.Exec(ctx, cmd, nil)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	if err := testutil.WaitForBlocks(ctx, 2, chain); err != nil {
		return sdk.TxResponse{}, err
	}

	var res sdk.TxResponse
	if err := json.Unmarshal(stdout, &res); err != nil {
		return res, err
	}

	return res, err
}

func TxCommandBuilder(ctx context.Context, chain *cosmos.CosmosChain, cmd []string, fromUser string, extraFlags ...string) []string {
	return TxCommandBuilderNode(ctx, chain.GetNode(), cmd, fromUser, extraFlags...)
}

func TxCommandBuilderNode(ctx context.Context, node *cosmos.ChainNode, cmd []string, fromUser string, extraFlags ...string) []string {
	command := []string{node.Chain.Config().Bin}
	command = append(command, cmd...)
	command = append(command, "--node", node.Chain.GetRPCAddress())
	command = append(command, "--home", node.HomeDir())
	command = append(command, "--chain-id", node.Chain.Config().ChainID)
	command = append(command, "--from", fromUser)
	command = append(command, "--keyring-backend", keyring.BackendTest)
	command = append(command, "--output=json")
	command = append(command, "--yes")

	gasFlag := false
	for _, flag := range extraFlags {
		if flag == "--gas" {
			gasFlag = true
		}
	}

	if !gasFlag {
		command = append(command, "--gas", "500000")
	}

	command = append(command, extraFlags...)
	return command
}

func getTransferChannel(channels []ibc.ChannelOutput) (string, error) {
	for _, channel := range channels {
		if channel.PortID == "transfer" && channel.State == ibcconntypes.OPEN.String() {
			return channel.ChannelID, nil
		}
	}

	return "", fmt.Errorf("no open transfer channel found: %+v", channels)
}
