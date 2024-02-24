package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/strangelove-ventures/poa"
	module "github.com/strangelove-ventures/simapp/x/example"
	"github.com/strangelove-ventures/simapp/x/example/keeper"
	"github.com/strangelove-ventures/simapp/x/example/types"
)

var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
	minttypes.ModuleName:           {authtypes.Minter},
	govtypes.ModuleName:            {authtypes.Burner},
}

type testFixture struct {
	suite.Suite

	ctx         sdk.Context
	k           keeper.Keeper
	msgServer   types.MsgServer
	queryServer types.QueryServer
	appModule   *module.AppModule

	accountkeeper authkeeper.AccountKeeper
	bankkeeper    bankkeeper.BaseKeeper
	mintkeeper    mintkeeper.Keeper

	addrs      []sdk.AccAddress
	govModAddr string
}

func SetupTest(t *testing.T, baseValShares int64) *testFixture {
	t.Helper()
	f := new(testFixture)
	// require := require.New(t)

	// Base setup
	logger := log.NewTestLogger(t)
	encCfg := moduletestutil.MakeTestEncodingConfig()

	f.govModAddr = authtypes.NewModuleAddress(govtypes.ModuleName).String()
	f.addrs = simtestutil.CreateIncrementalAccounts(3)

	key := storetypes.NewKVStoreKey(poa.ModuleName)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	f.ctx = testCtx.Ctx

	// Register SDK modules.
	// registerBaseSDKModules(f, encCfg, storeService, logger, require)

	// Setup POA Keeper.
	f.k = keeper.NewKeeper(encCfg.Codec, storeService, logger, f.govModAddr)
	f.msgServer = keeper.NewMsgServerImpl(f.k)
	f.queryServer = keeper.NewQuerier(f.k)
	f.appModule = module.NewAppModule(encCfg.Codec, f.k)

	// register interfaces
	registerModuleInterfaces(encCfg)

	// Setup initial keeper states
	// require.NoError(f.accountkeeper.AccountNumber.Set(f.ctx, 1))
	// f.accountkeeper.SetModuleAccount(f.ctx, f.accountkeeper.GetModuleAccount(f.ctx, minttypes.ModuleName))
	// f.mintkeeper.InitGenesis(f.ctx, f.accountkeeper, minttypes.DefaultGenesisState())

	// f.createBaseStakingValidators(t, baseValShares)

	return f
}

func registerModuleInterfaces(encCfg moduletestutil.TestEncodingConfig) {
	authtypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	stakingtypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	types.RegisterInterfaces(encCfg.InterfaceRegistry)
}
