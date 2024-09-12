package keeper

import (
	"fmt"

	"github.com/rollchains/spawn/simapp/x/ibcmodule/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the modules state from a specified GenesisState.
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	k.SetPort(ctx, types.PortID)

	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, types.PortID) {
		// module binds to the port on InitChain
		// and claims the returned capability
		if err := k.BindPort(ctx, types.PortID); err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}
}

// ExportGenesis exports the modules state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	port := k.GetPort(ctx)
	return &types.GenesisState{
		PortId: port,
	}
}
