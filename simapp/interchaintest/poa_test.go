package e2e

import (
	"context"
	"fmt"
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/strangelove-ventures/poa"
	"go.uber.org/zap/zaptest"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	numPOAVals int = 2
)

func TestPOA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	ctx := context.Background()
	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)
	client, network := interchaintest.DockerSetup(t)

	cs := &DefaultChainSpec
	cs.NumValidators = &numPOAVals
	cs.Env = []string{
		fmt.Sprintf("POA_ADMIN_ADDRESS=%s", "wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g"), // acc0 / admin
	}

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		cs,
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)

	// Setup Interchain
	ic := interchaintest.NewInterchain().
		AddChain(chain)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// setup accounts
	admin, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, "admin", AccMnemonic, GenesisFundsAmount, chain)
	require.NoError(t, err)

	acc0, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, "acc0", Acc1Mnemonic, GenesisFundsAmount, chain)
	require.NoError(t, err)

	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), GenesisFundsAmount, chain)
	incorrectUser := users[0]

	// get validator operator addresses
	vals, err := chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	require.NoError(t, err)

	validators := make([]string, len(vals))
	for i, v := range vals {
		validators[i] = v.OperatorAddress
	}
	require.Equal(t, len(validators), numPOAVals)

	// === Test Cases ===
	testStakingDisabled(t, ctx, chain, validators, acc0)
	testWithdrawDelegatorRewardsDisabled(t, ctx, chain, validators, acc0)
	testPowerErrors(t, ctx, chain, validators, incorrectUser, admin)
	testRemovePending(t, ctx, chain, admin)
}

func testWithdrawDelegatorRewardsDisabled(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, validators []string, acc0 ibc.Wallet) {
	t.Log("\n===== TEST WITHDRAW DELEGATOR REWARDS DISABLED =====")

	txRes, _ := WithdrawDelegatorRewards(t, ctx, chain, acc0, validators[0])
	require.Contains(t, txRes.RawLog, poa.ErrWithdrawDelegatorRewardsNotAllowed.Error())
}

func testStakingDisabled(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, validators []string, acc0 ibc.Wallet) {
	t.Log("\n===== TEST STAKING DISABLED =====")

	err := chain.GetNode().StakingDelegate(ctx, acc0.KeyName(), validators[0], "1stake")
	require.Error(t, err)
	require.Contains(t, err.Error(), poa.ErrStakingActionNotAllowed.Error())
}

func testRemovePending(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, admin ibc.Wallet) {
	t.Log("\n===== TEST PENDING =====")

	res, _ := POACreatePendingValidator(t, ctx, chain, admin, "pl3Q8OQwtC7G2dSqRqsUrO5VZul7l40I+MKUcejqRsg=", "testval", "0.10", "0.25", "0.05")
	require.EqualValues(t, 0, res.Code)

	require.NoError(t, testutil.WaitForBlocks(ctx, 2, chain))

	pv := GetPOAPending(t, ctx, chain)
	require.Equal(t, 1, len(pv))
	require.Equal(t, "0", pv[0].Tokens.String())
	require.Equal(t, "1", pv[0].MinSelfDelegation.String())

	res, _ = POARemovePending(t, ctx, chain, admin, pv[0].OperatorAddress)
	require.EqualValues(t, 0, res.Code)

	// validate it was removed
	pv = GetPOAPending(t, ctx, chain)
	require.Equal(t, 0, len(pv))
}

func testPowerErrors(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, validators []string, incorrectUser ibc.Wallet, admin ibc.Wallet) {
	t.Log("\n===== TEST POWER ERRORS =====")
	var res sdk.TxResponse
	var err error

	t.Run("fail: set-power message from a non authorized user", func(t *testing.T) {
		res, _ = POASetPower(t, ctx, chain, incorrectUser, validators[1], 1_000_000)
		res, err := chain.GetTransaction(res.TxHash)
		require.NoError(t, err)
		require.Contains(t, res.RawLog, poa.ErrNotAnAuthority.Error())
	})

	t.Run("fail: set-power message below minimum power requirement (self bond)", func(t *testing.T) {
		res, err = POASetPower(t, ctx, chain, admin, validators[0], 1)
		require.Error(t, err) // cli validate error
		require.Contains(t, err.Error(), poa.ErrPowerBelowMinimum.Error())
	})

	t.Run("fail: set-power message above 30%% without unsafe flag", func(t *testing.T) {
		res, _ = POASetPower(t, ctx, chain, admin, validators[0], 9_000_000_000_000_000)
		res, err := chain.GetTransaction(res.TxHash)
		require.NoError(t, err)
		require.Contains(t, res.RawLog, poa.ErrUnsafePower.Error())
	})
}
