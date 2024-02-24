package e2e

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/strangelove-ventures/poa"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	numPOAVals = 2
)

func TestPOA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	// setup base chain
	chains := interchaintest.CreateChainWithConfig(t, numPOAVals, NumberFullNodes, Name, ChainImage.Version, DefaultChainConfig)
	chain := chains[0].(*cosmos.CosmosChain)

	enableBlockDB := false
	ctx, _, _, _ := interchaintest.BuildInitialChain(t, chains, enableBlockDB)

	// setup accounts
	acc0, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, "acc0", AccMnemonic, GenesisFundsAmount, chain)
	if err != nil {
		t.Fatal(err)
	}
	acc1, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, "acc1", Acc1Mnemonic, GenesisFundsAmount, chain)
	if err != nil {
		t.Fatal(err)
	}

	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), GenesisFundsAmount, chain)
	incorrectUser := users[0]

	// get validator operator addresses
	vals, err := chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	require.NoError(t, err)
	require.Equal(t, len(vals), numPOAVals)

	validators := make([]string, len(vals))
	for i, v := range vals {
		validators[i] = v.OperatorAddress
	}

	// === Test Cases ===
	testStakingDisabled(t, ctx, chain, validators, acc0, acc1)
	testGovernance(t, ctx, chain, acc0, validators)
	testPowerErrors(t, ctx, chain, validators, incorrectUser, acc0)
	testPending(t, ctx, chain, acc0)
	testRemoveValidator(t, ctx, chain, validators, acc0)
	testUpdatePOAParams(t, ctx, chain, validators, acc0, incorrectUser)
}

func testUpdatePOAParams(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, validators []string, acc0, incorrectUser ibc.Wallet) {
	var tx sdk.TxResponse
	var err error

	t.Run("fail: update-params message from a non authorized user", func(t *testing.T) {
		tx, err = POAUpdateParams(t, ctx, chain, incorrectUser, []string{incorrectUser.FormattedAddress()}, true)
		if err != nil {
			t.Fatal(err)
		}
		txRes, err := chain.GetTransaction(tx.TxHash)
		require.NoError(t, err)
		fmt.Printf("%+v", txRes)
		require.Contains(t, txRes.RawLog, poa.ErrNotAnAuthority.Error())
	})

	t.Run("fail: update staking params from a non authorized user", func(t *testing.T) {
		tx, err = POAUpdateStakingParams(t, ctx, chain, incorrectUser, stakingtypes.DefaultParams())
		if err != nil {
			t.Fatal(err)
		}
		txRes, err := chain.GetTransaction(tx.TxHash)
		require.NoError(t, err)
		fmt.Printf("%+v", txRes)
		require.EqualValues(t, txRes.Code, 3)

		sp, err := chain.StakingQueryParams(ctx)
		require.NoError(t, err)
		fmt.Printf("%+v", sp)
	})

	t.Run("success: update staking params from an authorized user.", func(t *testing.T) {
		stakingparams := stakingtypes.DefaultParams()
		tx, err = POAUpdateStakingParams(t, ctx, chain, acc0, stakingtypes.Params{
			UnbondingTime:     stakingparams.UnbondingTime,
			MaxValidators:     10,
			MaxEntries:        stakingparams.MaxEntries,
			HistoricalEntries: stakingparams.HistoricalEntries,
			BondDenom:         stakingparams.BondDenom,
			MinCommissionRate: stakingparams.MinCommissionRate,
		})
		if err != nil {
			t.Fatal(err)
		}

		txRes, err := chain.GetTransaction(tx.TxHash)
		require.NoError(t, err)
		fmt.Printf("%+v", txRes)
		require.EqualValues(t, txRes.Code, 0)

		sp, err := chain.StakingQueryParams(ctx)
		require.NoError(t, err)
		fmt.Printf("%+v", sp)
		require.EqualValues(t, sp.MaxValidators, 10)
	})

	t.Run("success: update-params message from an authorized user.", func(t *testing.T) {
		govModule := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
		randAcc := "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a"

		newAdmins := []string{acc0.FormattedAddress(), govModule, randAcc, incorrectUser.FormattedAddress()}
		tx, err = POAUpdateParams(t, ctx, chain, acc0, newAdmins, true)
		if err != nil {
			t.Fatal(err)
		}
		txRes, err := chain.GetTransaction(tx.TxHash)
		require.NoError(t, err)
		fmt.Printf("%+v", txRes)
		require.EqualValues(t, txRes.Code, 0)

		p := GetPOAParams(t, ctx, chain)
		for _, admin := range newAdmins {
			require.Contains(t, p.Admins, admin)
		}
	})

	t.Run("success: gov proposal update", func(t *testing.T) {
		govModule := "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
		randAcc := "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a"

		updatedParams := []cosmos.ProtoMessage{
			&poa.MsgUpdateParams{
				Sender: govModule,
				Params: poa.Params{
					Admins: []string{acc0.FormattedAddress(), govModule, randAcc},
				},
			},
		}

		proposal, err := chain.BuildProposal(updatedParams, "UpdateParams", "params", "ipfs://CID", fmt.Sprintf(`50%s`, chain.Config().Denom), govModule, false)
		require.NoError(t, err, "error building proposal")

		txProp, err := chain.SubmitProposal(ctx, incorrectUser.KeyName(), proposal)
		t.Log("txProp", txProp)
		require.NoError(t, err, "error submitting proposal")

		err = chain.VoteOnProposalAllValidators(ctx, txProp.ProposalID, "YES")
		require.NoError(t, err, "error voting on proposal")

		height, err := chain.Height(ctx)
		require.NoError(t, err, "failed to get height")

		// conmvert txProp.ProposalID to a uint64
		propId, err := strconv.ParseUint(txProp.ProposalID, 10, 64)
		require.NoError(t, err, "failed to convert proposalID to uint64")

		resp, _ := cosmos.PollForProposalStatusV1(ctx, chain, height, height+10, propId, govv1.StatusPassed)
		t.Log("PollForProposalStatusV8 resp", resp)
		require.EqualValues(t, govv1.StatusPassed, resp.Status, "proposal status did not change to passed in expected number of blocks")

		for _, admin := range GetPOAParams(t, ctx, chain).Admins {
			require.NotEqual(t, admin, incorrectUser.FormattedAddress())
		}
	})

}

func testRemoveValidator(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, validators []string, acc0 ibc.Wallet) {
	t.Log("\n===== TEST REMOVE VALIDATOR =====")
	powerOne := int64(9_000_000_000_000)
	powerTwo := int64(2_500_000)

	POASetPower(t, ctx, chain, acc0, validators[0], powerOne, "--unsafe")
	res, err := POASetPower(t, ctx, chain, acc0, validators[1], powerTwo, "--unsafe")
	require.NoError(t, err)
	fmt.Printf("%+v", res)

	// decode res.TxHash into a TxResponse
	txRes, err := chain.GetTransaction(res.TxHash)
	require.NoError(t, err)
	fmt.Printf("%+v", txRes)

	if err := testutil.WaitForBlocks(ctx, 2, chain); err != nil {
		t.Fatal(err)
	}

	vals, err := chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%d", powerOne), vals[0].Tokens)
	require.Equal(t, fmt.Sprintf("%d", powerTwo), vals[1].Tokens)

	// validate the validators both have a conesnsus-power of /1_000_000
	p1 := GetPOAConsensusPower(t, ctx, chain, vals[0].OperatorAddress)
	require.EqualValues(t, powerOne/1_000_000, p1) // = 9000000
	p2 := GetPOAConsensusPower(t, ctx, chain, vals[1].OperatorAddress)
	require.EqualValues(t, powerTwo/1_000_000, p2) // = 2

	// remove the 2nd validator (lower power)
	POARemove(t, ctx, chain, acc0, validators[1])

	// allow the poa.BeginBlocker to update new status
	if err := testutil.WaitForBlocks(ctx, 5, chain); err != nil {
		t.Fatal(err)
	}

	vals, err = chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%d", powerOne), vals[0].Tokens)
	require.Equal(t, "0", vals[1].Tokens)
	require.Equal(t, 1, vals[1].Status) // 1 = unbonded

	// validate the validator[1] has no consensus power
	require.EqualValues(t, 0, GetPOAConsensusPower(t, ctx, chain, vals[1].OperatorAddress))
}

func testStakingDisabled(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, validators []string, acc0, acc1 ibc.Wallet) {
	t.Log("\n===== TEST STAKING DISABLED =====")

	err := chain.GetNode().StakingDelegate(ctx, acc0.KeyName(), validators[0], "1stake")
	require.Error(t, err)
	require.Contains(t, err.Error(), poa.ErrStakingActionNotAllowed.Error())

	granter := acc1
	grantee := acc0

	// Grant grantee (acc0) the ability to delegate from granter (acc1)
	// res, err := ExecuteAuthzGrantMsg(t, ctx, chain, granter, grantee, "/cosmos.staking.v1beta1.MsgDelegate")
	// res, err := chain.GetNode().AuthzGrant(ctx, granter, grantee.FormattedAddress(), "/cosmos.staking.v1beta1.MsgDelegate")
	res, err := chain.GetNode().AuthzGrant(ctx, granter, grantee.FormattedAddress(), "generic", "--msg-type", "/cosmos.staking.v1beta1.MsgDelegate")
	require.NoError(t, err)
	require.EqualValues(t, res.Code, 0)

	// Generate nested message
	nested := []string{"tx", "staking", "delegate", validators[0], "1stake"}
	nestedCmd := TxCommandBuilder(ctx, chain, nested, granter.FormattedAddress())

	// Execute nested message via a wrapped Exec
	_, err = chain.GetNode().AuthzExec(ctx, grantee, nestedCmd)
	require.Error(t, err)
	require.Contains(t, err.Error(), poa.ErrStakingActionNotAllowed.Error())
}

func testPending(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, acc0 ibc.Wallet) {
	t.Log("\n===== TEST PENDING =====")

	_, err := POACreatePendingValidator(t, ctx, chain, acc0, "pl3Q8OQwtC7G2dSqRqsUrO5VZul7l40I+MKUcejqRsg=", "testval", "0.10", "0.25", "0.05")
	require.NoError(t, err)

	pv := GetPOAPending(t, ctx, chain)
	require.Equal(t, 1, len(pv))
	require.Equal(t, "0", pv[0].Tokens)
	require.Equal(t, "1", pv[0].MinSelfDelegation)

	_, err = POARemovePending(t, ctx, chain, acc0, pv[0].OperatorAddress)
	require.NoError(t, err)

	// validate it was removed
	pv = GetPOAPending(t, ctx, chain)
	require.Equal(t, 0, len(pv))
}

func testGovernance(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, acc0 ibc.Wallet, validators []string) {
	t.Log("\n===== TEST GOVERNANCE =====")
	// ibc.ChainConfig key: app_state.poa.params.admins must contain the governance address.
	propID := SubmitGovernanceProposalForValidatorChanges(t, ctx, chain, acc0, validators[0], 1_234_567, true)
	// ValidatorVote(t, ctx, chain, propID, cosmos.ProposalVoteYes, 25)

	// vote on proposal
	err := chain.VoteOnProposalAllValidators(ctx, propID, "YES")
	require.NoError(t, err)

	// validate the validator[0] was set to 1_234_567
	// val := GetValidators(t, ctx, chain).Validators[0]
	vals, err := chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	require.NoError(t, err)

	val := vals[0]

	require.Equal(t, val.Tokens, "1234567")
	p := GetPOAConsensusPower(t, ctx, chain, val.OperatorAddress)
	require.EqualValues(t, 1_234_567/1_000_000, p)
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
