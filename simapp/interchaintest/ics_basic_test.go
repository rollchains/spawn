package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcconntypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"

	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/relayer"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// This tests Cosmos Interchain Security, spinning up a provider and a single consumer chain.
func TestICSBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	ctx := context.Background()

	vals := 1
	fNodes := 0

	providerVer := "v5.0.0-rc0"

	// Chain Factory
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name: "ics-provider", Version: providerVer,
			NumValidators: &vals, NumFullNodes: &fNodes,
			ChainConfig: ibc.ChainConfig{
				GasAdjustment:  1.5,
				TrustingPeriod: "336h",
			},
		},
		{
			Name:          Name,
			Version:       "local",
			ChainName:     ChainID,
			NumValidators: &vals,
			NumFullNodes:  &fNodes,
			ChainConfig: ibc.ChainConfig{
				Images: []ibc.DockerImage{
					ChainImage,
				},
				GasAdjustment:  1.5,
				EncodingConfig: GetEncodingConfig(),
				Type:           "cosmos",
				Name:           Name,
				ChainID:        ChainID,
				Bin:            Binary,
				Bech32Prefix:   Bech32,
				Denom:          Denom,
				CoinType:       "118",
				GasPrices:      "0" + Denom,
				TrustingPeriod: "336h",
				// InterchainSecurityConfig: ibc.ICSConfig{
				// 	ProviderVerOverride: providerVer,
				// 	ConsumerVerOverride: "v4.1.0",
				// },
			},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	provider, consumer := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(RelayerRepo, RelayerVersion, "100:1000"),
		relayer.StartupFlags("--block-history", "200"),
	).Build(t, client, network)

	// Prep Interchain
	const ibcPath = "ics-path"
	ic := interchaintest.NewInterchain().
		AddChain(provider).
		AddChain(consumer).
		AddRelayer(r, "relayer").
		AddProviderConsumerLink(interchaintest.ProviderConsumerLink{
			Provider: provider,
			Consumer: consumer,
			Relayer:  r,
			Path:     ibcPath,
		})

	// Log location
	f, err := interchaintest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)

	// Reporter/logs
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	// Build interchain
	err = ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: false,
	})
	require.NoError(t, err, "failed to build interchain")

	require.NoError(t, provider.FinishICSProviderSetup(ctx, r, eRep, ibcPath))

	amt := math.NewInt(10_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", amt, consumer, provider)
	consumerUser, providerUser := users[0], users[1]

	t.Run("validate funding", func(t *testing.T) {
		bal, err := consumer.BankQueryBalance(ctx, consumerUser.FormattedAddress(), consumer.Config().Denom)
		require.NoError(t, err)
		require.EqualValues(t, amt, bal)

		bal, err = provider.BankQueryBalance(ctx, providerUser.FormattedAddress(), provider.Config().Denom)
		require.NoError(t, err)
		require.EqualValues(t, amt, bal)
	})

	t.Run("validate consumer action executed", func(t *testing.T) {
		bal, err := consumer.BankQueryBalance(ctx, consumerUser.FormattedAddress(), consumer.Config().Denom)
		require.NoError(t, err)
		require.EqualValues(t, amt, bal)
	})

	t.Run("provider -> consumer IBC transfer", func(t *testing.T) {
		providerChannelInfo, err := r.GetChannels(ctx, eRep, provider.Config().ChainID)
		require.NoError(t, err)

		channelID, err := getTransferChannel(providerChannelInfo)
		require.NoError(t, err, providerChannelInfo)

		consumerChannelInfo, err := r.GetChannels(ctx, eRep, consumer.Config().ChainID)
		require.NoError(t, err)

		consumerChannelID, err := getTransferChannel(consumerChannelInfo)
		require.NoError(t, err, consumerChannelInfo)

		dstAddress := consumerUser.FormattedAddress()
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

		consumerBal, err := consumer.BankQueryBalance(ctx, consumerUser.FormattedAddress(), dstIbcDenom)
		require.NoError(t, err)
		require.EqualValues(t, sendAmt, consumerBal)
	})
}

func getTransferChannel(channels []ibc.ChannelOutput) (string, error) {
	for _, channel := range channels {
		if channel.PortID == "transfer" && channel.State == ibcconntypes.OPEN.String() {
			return channel.ChannelID, nil
		}
	}

	return "", fmt.Errorf("no open transfer channel found")
}
