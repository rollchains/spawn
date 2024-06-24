package app

import (
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/client"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	ibctesting "github.com/cosmos/ibc-go/v8/testing/types"
	ethostestutil "github.com/ethos-works/ethos/ethos-chain/testutil/integration"      // spawntag:ethos-ics
	ibcconsumerkeeper "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/keeper" // spawntag:ethos-ics
)

func (app *ChainApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *ChainApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *ChainApp) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

func (app *ChainApp) GetTestEvidenceKeeper() evidencekeeper.Keeper {
	return app.EvidenceKeeper
}

// func (app *ChainApp) GetStakingKeeper() *stakingkeeper.Keeper { // ?spawntag:ethos-ics
func (app *ChainApp) GetStakingKeeper() ibctesting.StakingKeeper { // spawntag:ethos-ics
	return app.StakingKeeper
}

func (app *ChainApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *ChainApp) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}

// <spawntag:ethos-ics
func (app *ChainApp) GetTestAccountKeeper() ethostestutil.TestAccountKeeper {
	return app.AccountKeeper
}

func (app *ChainApp) GetTestBankKeeper() ethostestutil.TestBankKeeper {
	return app.BankKeeper
}

func (app *ChainApp) GetTestSlashingKeeper() ethostestutil.TestSlashingKeeper {
	return app.SlashingKeeper
}

func (app *ChainApp) GetConsumerKeeper() ibcconsumerkeeper.Keeper {
	panic("GetConsumerKeeper not implemented on app for ethos yet. TODO(reece)")
	return ibcconsumerkeeper.Keeper{}
	// return app.ConsumerKeeper
}

// GetTxConfig returns ChainApp's TxConfig
func (app *ChainApp) GetTxConfig() client.TxConfig {
	return app.TxConfig()
}

// spawntag:ethos-ics>
