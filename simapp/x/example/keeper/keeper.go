package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"

	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/strangelove-ventures/simapp/x/example/types"
)

type Keeper struct {
	cdc                   codec.BinaryCodec
	validatorAddressCodec addresscodec.Codec

	logger log.Logger

	// state management
	Schema collections.Schema
	Params collections.Item[types.Params]
}

// NewKeeper creates a new poa Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	validatorAddressCodec addresscodec.Codec,
	logger log.Logger,
) Keeper {
	logger = logger.With(log.ModuleKey, "x/"+types.ModuleName)

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:                   cdc,
		validatorAddressCodec: validatorAddressCodec,
		logger:                logger,

		// Stores
		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k Keeper) Logger() log.Logger {
	return k.logger
}

// ExportGenesis exports the module's state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params: params,
	}
}
