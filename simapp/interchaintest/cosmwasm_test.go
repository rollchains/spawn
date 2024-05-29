package e2e

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type GetCountResponse struct {
	// {"data":{"count":0}}
	Data *GetCountObj `json:"data"`
}

type GetCountObj struct {
	Count int64 `json:"count"`
}

func TestCosmWasmIntegration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)
	client, network := interchaintest.DockerSetup(t)

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		&DefaultChainSpec,
		&ProviderChain, // spawntag:ics
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)
	provider := chains[1].(*cosmos.CosmosChain) // spawntag:ics

	// <spawntag:ics
	// Relayer Factory
	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		relayer.CustomDockerImage(RelayerRepo, RelayerVersion, "100:1000"),
		relayer.StartupFlags("--block-history", "200"),
	).Build(t, client, network)
	// spawntag:ics>

	// Setup Interchain
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

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	require.NoError(t, provider.FinishICSProviderSetup(ctx, r, eRep, ibcPath)) // spawntag:ics

	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), GenesisFundsAmount, chain)
	user := users[0]

	StdExecute(t, ctx, chain, user)
}

func StdExecute(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet) (contractAddr string) {
	_, contractAddr = SetupContract(t, ctx, chain, user.KeyName(), "contracts/cw_template.wasm", `{"count":0}`)
	chain.ExecuteContract(ctx, user.KeyName(), contractAddr, `{"increment":{}}`, "--fees", "10000"+chain.Config().Denom)

	var res GetCountResponse
	err := SmartQueryString(t, ctx, chain, contractAddr, `{"get_count":{}}`, &res)
	require.NoError(t, err)

	require.Equal(t, int64(1), res.Data.Count)

	return contractAddr
}

func SmartQueryString(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, contractAddr, queryMsg string, res interface{}) error {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(queryMsg), &jsonMap); err != nil {
		t.Fatal(err)
	}
	err := chain.QueryContract(ctx, contractAddr, jsonMap, &res)
	return err
}

func SetupContract(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, keyname string, fileLoc string, message string, extraFlags ...string) (codeId, contract string) {
	codeId, err := chain.StoreContract(ctx, keyname, fileLoc)
	if err != nil {
		t.Fatal(err)
	}

	needsNoAdminFlag := true
	for _, flag := range extraFlags {
		if flag == "--admin" {
			needsNoAdminFlag = false
		}
	}

	contractAddr, err := chain.InstantiateContract(ctx, keyname, codeId, message, needsNoAdminFlag, extraFlags...)
	if err != nil {
		t.Fatal(err)
	}

	return codeId, contractAddr
}
