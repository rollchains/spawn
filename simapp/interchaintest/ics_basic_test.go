package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	interchaintest "github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
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
				GasAdjustment: 1.5,
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
				InterchainSecurityConfig: ibc.ICSConfig{
					ProviderVerOverride: providerVer,
					ConsumerVerOverride: "v4.1.0", // v5.0.0-rc0 & v4.1.0 are compatible
				},
			},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	provider, consumer := chains[0], chains[1]

	// Relayer Factory
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
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

	err = testutil.WaitForBlocks(ctx, 5, provider, consumer)
	require.NoError(t, err, "failed to wait for blocks")
}
