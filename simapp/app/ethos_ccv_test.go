package app

import (
	"encoding/json"

	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	cometbfttypes "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	providerApp "github.com/ethos-works/ethos/ethos-chain/testapps/provider"

	"testing"

	"github.com/ethos-works/ethos/ethos-chain/tests/integration"
	icstestingutils "github.com/ethos-works/ethos/ethos-chain/testutil/ibc_testing"
	consumertypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/types"
	"github.com/stretchr/testify/suite"
)

var (
	ccvSuite *integration.CCVTestSuite
)

func init() {
	sdk.SetAddrCacheEnabled(false)

	// Pass in concrete app types that implement the interfaces defined in https://github.com/cosmos/interchain-security/testutil/integration/interfaces.go
	// IMPORTANT: the concrete app types passed in as type parameters here must match the
	// concrete app types returned by the relevant app initers.
	ccvSuite = integration.NewCCVTestSuite[*providerApp.App, *ChainApp](
		// Pass in ibctesting.AppIniters for gaia (provider) and consumer.
		icstestingutils.ProviderAppIniter, SetupValSetAppIniter, []string{})
}

func TestCCVTestSuite(t *testing.T) {
	// Run tests
	suite.Run(t, ccvSuite)
}

// SetupValSetAppIniter is a simple wrapper for ICS e2e tests to satisfy interface
func SetupValSetAppIniter(initValUpdates []cometbfttypes.ValidatorUpdate) icstestingutils.AppIniter {
	return SetupTestingApp(initValUpdates)
}

// SetupTestingApp initializes the IBC-go testing application
func SetupTestingApp(initValUpdates []cometbfttypes.ValidatorUpdate) func() (ibctesting.TestingApp, map[string]json.RawMessage) {
	return func() (ibctesting.TestingApp, map[string]json.RawMessage) {
		db := dbm.NewMemDB()
		emptyAppOpts := make(simtestutil.AppOptionsMap, 0)

		testApp := NewChainApp(
			log.NewNopLogger(), db, nil, false, emptyAppOpts,
			[]wasmkeeper.Option{}, // spawntag:wasm
			baseapp.SetChainID(chainID),
		)
		encoding := testApp.AppCodec()

		// we need to set up a TestInitChainer where we can redefine MaxBlockGas in ConsensusParamsKeeper
		testApp.SetInitChainer(testApp.InitChainer)
		// and then we manually init baseapp and load states
		testApp.LoadLatestVersion()

		genesisState := testApp.DefaultGenesis()

		// NOTE ibc-go/v8/testing.SetupWithGenesisValSet requires a staking module
		// genesisState or it panics. Feed a minimum one.
		genesisState[stakingtypes.ModuleName] = encoding.MustMarshalJSON(
			&stakingtypes.GenesisState{
				Params: stakingtypes.Params{BondDenom: sdk.DefaultBondDenom},
			},
		)

		var consumerGenesis ccvtypes.ConsumerGenesisState
		encoding.MustUnmarshalJSON(genesisState[consumertypes.ModuleName], &consumerGenesis)
		consumerGenesis.Provider.InitialValSet = initValUpdates
		consumerGenesis.Params.Enabled = true
		genesisState[consumertypes.ModuleName] = encoding.MustMarshalJSON(&consumerGenesis)

		return testApp, genesisState
	}
}
