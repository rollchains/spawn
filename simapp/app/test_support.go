package app

import (
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	ibctesting "github.com/cosmos/ibc-go/v8/testing/types"
	ethostestutil "github.com/ethos-works/ethos/ethos-chain/testutil/integration"      // spawntag:ethos-ics
	ibcconsumerkeeper "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/keeper" // spawntag:ethos-ics

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
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

// spawntag:ethos-ics>

// <spawntag:ics
func (app *ChainApp) GetConsumerKeeper() ibcconsumerkeeper.Keeper {
	return app.ConsumerKeeper
}

// spawntag:ics>

func (app *ChainApp) GetStakingKeeper() ibctesting.StakingKeeper {
	return app.StakingKeeper
}

func (app *ChainApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *ChainApp) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}
