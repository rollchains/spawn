package app

import (
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

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

// <spawntag:staking
func (app *ChainApp) GetStakingKeeper() *stakingkeeper.Keeper {
	return app.StakingKeeper
}

// spawntag:staking>

func (app *ChainApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *ChainApp) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}

// <spawntag:ethos-ics
// TODO:
// func (app *ChainApp) GetTestAccountKeeper() ethostestutil.TestAccountKeeper {
// 	return app.AccountKeeper
// }

// func (app *ChainApp) GetTestBankKeeper() ethostestutil.TestBankKeeper {
// 	return app.BankKeeper
// }

// func (app *ChainApp) GetTestSlashingKeeper() ethostestutil.TestSlashingKeeper {
// 	return app.SlashingKeeper
// }

// spawntag:ethos-ics>
