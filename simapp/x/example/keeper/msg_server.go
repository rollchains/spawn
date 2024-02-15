package keeper

import (
	"context"

	"github.com/strangelove-ventures/simapp/x/example/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (msgServer) UpdateParams(context.Context, *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	panic("unimplemented")
}
