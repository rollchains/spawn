package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/rollchains/spawn/simapp/x/ibcmodule/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// SendTx implements types.MsgServer.
func (ms msgServer) SendExampleTx(ctx context.Context, msg *types.MsgSendExampleTx) (*types.MsgSendExampleTxResponse, error) {
	sequence, err := ms.sendPacket(
		ctx, msg.SourcePort, msg.SourceChannel, msg.Sender, msg.SomeData, msg.TimeoutTimestamp)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendExampleTxResponse{
		Sequence: sequence,
	}, nil
}

func (ms msgServer) sendPacket(ctx context.Context, sourcePort, sourceChannel, sender, someData string, timeoutTimestamp uint64) (sequence uint64, err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	channelCap, ok := ms.ScopedKeeper.GetCapability(sdkCtx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, fmt.Errorf("module does not own channel capability")
	}

	packetData := types.ExamplePacketData{
		Sender:   sender,
		SomeData: someData,
	}

	sequence, err = ms.ics4Wrapper.SendPacket(sdkCtx, channelCap, sourcePort, sourceChannel, clienttypes.ZeroHeight(), timeoutTimestamp, packetData.MustGetBytes())
	if err != nil {
		return 0, err
	}
	return sequence, nil
}
