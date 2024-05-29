package e2e

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestIBCRateLimit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cs := &DefaultChainSpec
	cs.ModifyGenesis = cosmos.ModifyGenesis([]cosmos.GenesisKV{cosmos.NewGenesisKV("app_state.ratelimit.blacklisted_denoms", []string{cs.Denom})})

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		cs,
		&ProviderChain, // spawntag:ics
		// <spawntag:staking
		func() *interchaintest.ChainSpec {
			SecondChainSpec := &DefaultChainSpec
			// TODO: This should not be needed
			// blacklistedDenoms := []string{SecondChainSpec.Denom}
			// SecondChainSpec.ModifyGenesis = cosmos.ModifyGenesis(
			// 	append(DefaultGenesis,
			// 		cosmos.NewGenesisKV("app_state.ratelimit.blacklisted_denoms", blacklistedDenoms),
			// 	),
			// )
			SecondChainSpec.ChainID += "2"
			return SecondChainSpec
		}(),
		// spawntag:staking>
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)
	secondary := chains[1].(*cosmos.CosmosChain)

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)
	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t, zaptest.Level(zapcore.DebugLevel)),
		interchaintestrelayer.CustomDockerImage(RelayerRepo, RelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "200"),
	).Build(t, client, network)

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ic := interchaintest.NewInterchain().
		AddChain(chain).
		AddChain(secondary).
		AddRelayer(r, "relayer")

	// <spawntag:staking
	ic = ic.AddLink(interchaintest.InterchainLink{
		Chain1:  chain,
		Chain2:  secondary,
		Relayer: r,
		Path:    ibcPath,
	})
	// spawntag:staking>
	// <spawntag:ics
	ic = ic.AddProviderConsumerLink(interchaintest.ProviderConsumerLink{
		Provider: secondary,
		Consumer: chain,
		Relayer:  r,
		Path:     ibcPath,
	})
	// spawntag:ics>

	// Build interchain
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	}))

	require.NoError(t, secondary.FinishICSProviderSetup(ctx, r, eRep, ibcPath)) // spawntag:ics

	// Create and Fund User Wallets
	fundAmount := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", fundAmount, chain, secondary)
	userA, userB := users[0], users[1]

	userAInitial, err := chain.GetBalance(ctx, userA.FormattedAddress(), chain.Config().Denom)
	fmt.Println("userAInitial", userAInitial)
	require.NoError(t, err)
	require.True(t, userAInitial.Equal(fundAmount))

	// Get Channel ID
	aInfo, err := r.GetChannels(ctx, eRep, chain.Config().ChainID)
	require.NoError(t, err)
	aChannelID, err := getTransferChannel(aInfo)
	require.NoError(t, err)
	fmt.Println("aChannelID", aChannelID)

	// Send Transaction
	amountToSend := math.NewInt(1_000_000)
	dstAddress := userB.FormattedAddress()
	transfer := ibc.WalletAmount{
		Address: dstAddress,
		Denom:   chain.Config().Denom,
		Amount:  amountToSend,
	}

	// Validate transfer error occurs
	_, err = chain.SendIBCTransfer(ctx, aChannelID, userA.KeyName(), transfer, ibc.TransferOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "denom is blacklisted")
}
