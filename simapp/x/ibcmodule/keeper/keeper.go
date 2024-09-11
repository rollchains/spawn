package keeper

import (
	"github.com/rollchains/spawn/simapp/x/ibcmodule/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/log"

	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	channelkeeper "github.com/cosmos/ibc-go/v8/modules/core/04-channel/keeper"
	portkeeper "github.com/cosmos/ibc-go/v8/modules/core/05-port/keeper"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// Keeper defines the middleware keeper.
type Keeper struct {
	storeKey         storetypes.StoreKey
	cdc              codec.BinaryCodec
	msgServiceRouter *baseapp.MsgServiceRouter
	Schema           collections.Schema

	ChannelKeeper channelkeeper.Keeper // TODO: this or use ics4Wrapper?
	PortKeeper    *portkeeper.Keeper
	ScopedKeeper  capabilitykeeper.ScopedKeeper

	// used to send the packet, usually the IBC channel keeper.
	ics4Wrapper porttypes.ICS4Wrapper
}

// NewKeeper creates a new swap Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	msgServiceRouter *baseapp.MsgServiceRouter,
	ics4Wrapper porttypes.ICS4Wrapper,
	scopedKeeper capabilitykeeper.ScopedKeeper,
	portKeeper *portkeeper.Keeper,
	channelKeeper channelkeeper.Keeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:              cdc,
		msgServiceRouter: msgServiceRouter,
		ics4Wrapper:      ics4Wrapper,

		storeKey:      storetypes.NewKVStoreKey(types.StoreKey), // TODO: remove me
		ChannelKeeper: channelKeeper,
		PortKeeper:    portKeeper,
		ScopedKeeper:  scopedKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// WithICS4Wrapper sets the ICS4Wrapper. This function may be used after
// the keeper's creation to set the middleware which is above this module
// in the IBC application stack.
func (k *Keeper) WithICS4Wrapper(wrapper porttypes.ICS4Wrapper) {
	k.ics4Wrapper = wrapper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+ibcexported.ModuleName+"-"+types.ModuleName)
}
