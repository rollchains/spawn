package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcconntypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestBasicChain(t *testing.T) {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		&ProviderChain, // spawntag:ics
		&DefaultChainSpec,
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)
	provider := chains[1].(*cosmos.CosmosChain) // spawntag:ics

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	// <spawntag:ics
	// Relayer Factory
	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		relayer.CustomDockerImage(RelayerRepo, RelayerVersion, "100:1000"),
		relayer.StartupFlags("--block-history", "200"),
	).Build(t, client, network)

	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)

	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)
	// spawntag:ics>

	ic := interchaintest.NewInterchain().
		AddChain(chain)

	// <spawntag:ics
	ic = ic.AddChain(provider).
		AddRelayer(r, "relayer").
		AddProviderConsumerLink(interchaintest.ProviderConsumerLink{
			Provider: provider,
			Consumer: chain,
			Relayer:  r,
			Path:     ibcPath,
		})
	// spawntag:ics>

	require.NoError(t, ic.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// <spawntag:ics
	require.NoError(t, provider.FinishICSProviderSetup(ctx, r, eRep, ibcPath))
	// spawntag:ics>

	amt := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", amt,
		chain,
		provider, //spawntag:ics
	)
	user := users[0]
	providerUser := users[1] // spawntag:ics

	t.Run("validate funding", func(t *testing.T) {
		bal, err := chain.BankQueryBalance(ctx, user.FormattedAddress(), chain.Config().Denom)
		require.NoError(t, err)
		require.EqualValues(t, amt, bal)

		// <spawntag:ics
		bal, err = provider.BankQueryBalance(ctx, providerUser.FormattedAddress(), provider.Config().Denom)
		require.NoError(t, err)
		require.EqualValues(t, amt, bal)
		// spawntag:ics>
	})

	// <spawntag:ics
	t.Run("provider -> consumer IBC transfer", func(t *testing.T) {
		providerChannelInfo, err := r.GetChannels(ctx, eRep, provider.Config().ChainID)
		require.NoError(t, err)

		channelID, err := getTransferChannel(providerChannelInfo)
		require.NoError(t, err, providerChannelInfo)

		consumerChannelInfo, err := r.GetChannels(ctx, eRep, chain.Config().ChainID)
		require.NoError(t, err)

		consumerChannelID, err := getTransferChannel(consumerChannelInfo)
		require.NoError(t, err, consumerChannelInfo)

		dstAddress := user.FormattedAddress()
		sendAmt := math.NewInt(7)
		transfer := ibc.WalletAmount{
			Address: dstAddress,
			Denom:   provider.Config().Denom,
			Amount:  sendAmt,
		}

		tx, err := provider.SendIBCTransfer(ctx, channelID, providerUser.KeyName(), transfer, ibc.TransferOptions{})
		require.NoError(t, err)
		require.NoError(t, tx.Validate())

		require.NoError(t, r.Flush(ctx, eRep, ibcPath, channelID))

		srcDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom("transfer", consumerChannelID, provider.Config().Denom))
		dstIbcDenom := srcDenomTrace.IBCDenom()

		consumerBal, err := chain.BankQueryBalance(ctx, user.FormattedAddress(), dstIbcDenom)
		require.NoError(t, err)
		require.EqualValues(t, sendAmt, consumerBal)
	})
	// spawntag:ics>
}

// <spawntag:ics
func getTransferChannel(channels []ibc.ChannelOutput) (string, error) {
	for _, channel := range channels {
		if channel.PortID == "transfer" && channel.State == ibcconntypes.OPEN.String() {
			return channel.ChannelID, nil
		}
	}

	return "", fmt.Errorf("no open transfer channel found")
}

// spawntag:ics>
