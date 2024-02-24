package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

type GetCountResponse struct {
	// {"data":{"count":0}}
	Data *GetCountObj `json:"data"`
}

type GetCountObj struct {
	Count int64 `json:"count"`
}

func TestCosmWasmIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Base setup
	chains := interchaintest.CreateChainWithConfig(t, 1, 0, Name, ChainImage.Version, DefaultChainConfig)
	ctx, ic, _, _ := interchaintest.BuildInitialChain(t, chains, false)

	require.NotNil(t, ic)
	require.NotNil(t, ctx)

	chain := chains[0].(*cosmos.CosmosChain)

	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), GenesisFundsAmount, chain)
	user := users[0]

	StdExecute(t, ctx, chain, user)
	subMsg(t, ctx, chain, user)

	t.Cleanup(func() {
		_ = ic.Close()
	})
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

func subMsg(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet) {
	// ref: https://github.com/CosmWasm/wasmd/issues/1735

	// === execute a contract sub message ===
	_, senderContractAddr := SetupContract(t, ctx, chain, user.KeyName(), "contracts/cw721_base.wasm.gz", fmt.Sprintf(`{"name":"NFT #00001", "symbol":"nft-test-#00001", "minter":"%s"}`, user.FormattedAddress()))
	_, receiverContractAddr := SetupContract(t, ctx, chain, user.KeyName(), "contracts/cw721_receiver.wasm.gz", `{}`)

	// mint a token
	res, err := chain.ExecuteContract(ctx, user.KeyName(), senderContractAddr, fmt.Sprintf(`{"mint":{"token_id":"00000", "owner":"%s"}}`, user.FormattedAddress()), "--fees", "10000"+chain.Config().Denom)
	fmt.Println("First", res)
	require.NoError(t, err)

	// this purposely will fail with the current, we are just validating the messsage is not unknown.
	// sub message of unknown means the `wasmkeeper.WithMessageHandlerDecorator` is not setup properly.
	fail := "ImZhaWwi"
	res2, err := chain.ExecuteContract(ctx, user.KeyName(), senderContractAddr, fmt.Sprintf(`{"send_nft": { "contract": "%s", "token_id": "00000", "msg": "%s" }}`, receiverContractAddr, fail), "--fees", "10000"+chain.Config().Denom)
	require.NoError(t, err)
	fmt.Println("Second", res2)
	require.NotEqualValues(t, wasmtypes.ErrUnknownMsg.ABCICode(), res2.Code)
	require.NotContains(t, res2.RawLog, "unknown message from the contract")

	success := "InN1Y2NlZWQi"
	res3, err := chain.ExecuteContract(ctx, user.KeyName(), senderContractAddr, fmt.Sprintf(`{"send_nft": { "contract": "%s", "token_id": "00000", "msg": "%s" }}`, receiverContractAddr, success), "--fees", "10000"+chain.Config().Denom, "--amount", "10000"+chain.Config().Denom)
	require.NoError(t, err)
	fmt.Println("Third", res3)
	require.EqualValues(t, 0, res3.Code)
	require.NotContains(t, res3.RawLog, "unknown message from the contract")
}

// TODO: use internal functions now instead of these
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

// func ExecuteMsgWithFee(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet, contractAddr, amount, feeCoin, message string) {
// 	// amount is #utoken

// 	// There has to be a way to do this in ictest?
// 	cmd := []string{
// 		"junod", "tx", "wasm", "execute", contractAddr, message,
// 		"--node", chain.GetRPCAddress(),
// 		"--home", chain.HomeDir(),
// 		"--chain-id", chain.Config().ChainID,
// 		"--from", user.KeyName(),
// 		"--gas", "500000",
// 		"--fees", feeCoin,
// 		"--keyring-dir", chain.HomeDir(),
// 		"--keyring-backend", keyring.BackendTest,
// 		"-y",
// 	}

// 	if amount != "" {
// 		cmd = append(cmd, "--amount", amount)
// 	}

// 	_, _, err := chain.Exec(ctx, cmd, nil)
// 	require.NoError(t, err)

// 	if err := testutil.WaitForBlocks(ctx, 2, chain); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func ExecuteMsgWithFeeReturn(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet, contractAddr, amount, feeCoin, message string) (*sdk.TxResponse, error) {
// 	// amount is #utoken

// 	// There has to be a way to do this in ictest? (there is, use node.ExecTx)
// 	cmd := []string{
// 		"wasm", "execute", contractAddr, message,
// 		"--output", "json",
// 		"--node", chain.GetRPCAddress(),
// 		"--home", chain.HomeDir(),
// 		"--gas", "500000",
// 		"--fees", feeCoin,
// 		"--keyring-dir", chain.HomeDir(),
// 	}

// 	if amount != "" {
// 		cmd = append(cmd, "--amount", amount)
// 	}

// 	node := chain.GetNode()

// 	txHash, err := node.ExecTx(ctx, user.KeyName(), cmd...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// convert stdout into a TxResponse
// 	txRes, err := chain.GetTransaction(txHash)
// 	return txRes, err
// }
