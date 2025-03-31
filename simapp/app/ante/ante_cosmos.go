package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	sdkvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evmoscosmosante "github.com/cosmos/evm/ante/cosmos" // spawntag:evm
	evmante "github.com/cosmos/evm/ante/evm"            // spawntag:evm
	evmtypes "github.com/cosmos/evm/x/vm/types"

	sdkmath "cosmossdk.io/math"
	circuitante "cosmossdk.io/x/circuit/ante"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	consumerdemocracy "github.com/cosmos/interchain-security/v5/app/consumer-democracy"
	ccvdemocracyante "github.com/cosmos/interchain-security/v5/app/consumer-democracy/ante"
	ccvconsumerante "github.com/cosmos/interchain-security/v5/app/consumer/ante"
	poaante "github.com/strangelove-ventures/poa/ante"
)

// newCosmosAnteHandler creates the default ante handler for Cosmos transactions
func NewCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	poaDoGenTxRateValidation := false
	poaRateFloor := sdkmath.LegacyMustNewDecFromStr("0.10")
	poaRateCeil := sdkmath.LegacyMustNewDecFromStr("0.50")

	return sdk.ChainAnteDecorators(
		// <spawntag:evm
		evmoscosmosante.NewRejectMessagesDecorator(), // reject MsgEthereumTxs
		evmoscosmosante.NewAuthzLimiterDecorator( // disable the Msg types that cannot be included on an authz.MsgExec msgs field
			sdk.MsgTypeURL(&evmtypes.MsgEthereumTx{}),
			sdk.MsgTypeURL(&sdkvesting.MsgCreateVestingAccount{}),
		),
		// spawntag:evm>
		ante.NewSetUpContextDecorator(),
		ccvconsumerante.NewMsgFilterDecorator(options.ConsumerKeeper),
		ccvconsumerante.NewDisabledModulesDecorator("/cosmos.evidence", "/cosmos.slashing"),
		ccvdemocracyante.NewForbiddenProposalsDecorator(consumerdemocracy.IsProposalWhitelisted, consumerdemocracy.IsModuleWhiteList),
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		evmoscosmosante.NewMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper), // spawntag:evm
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		evmante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper), // spawntag:evm
		poaante.NewPOADisableStakingDecorator(),
		poaante.NewPOADisableWithdrawDelegatorRewards(),
		poaante.NewCommissionLimitDecorator(poaDoGenTxRateValidation, poaRateFloor, poaRateCeil),
	)
}
