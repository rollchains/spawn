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
func (ms msgServer) SendTx(ctx context.Context, msg *types.MsgSendTx) (*types.MsgSendTxResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	// panic("SendTx is unimplemented")

	// sender, err := sdk.AccAddressFromBech32(msg.Sender)
	// if err != nil {
	// 	return nil, err
	// }

	sequence, err := ms.sendPacket(
		ctx, msg.SourcePort, msg.SourceChannel, msg.Sender, msg.SomeData, msg.TimeoutTimestamp)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendTxResponse{
		Sequence: sequence,
	}, nil
}

func (ms msgServer) sendPacket(ctx context.Context, sourcePort, sourceChannel, sender, someData string, timeoutTimestamp uint64) (sequence uint64, err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	channel, found := ms.ChannelKeeper.GetChannel(sdkCtx, sourcePort, sourceChannel)
	if !found {
		return 0, fmt.Errorf("channel not found: port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := channel.GetCounterparty().GetPortID()
	destinationChannel := channel.GetCounterparty().GetChannelID()

	fmt.Printf("destinationPort: %s\n", destinationPort)
	fmt.Printf("destinationChannel: %s\n", destinationChannel)
	fmt.Printf("Channel Information: %+v\n", channel)

	// begin createOutgoingPacket logic
	// See spec for this logic: https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay
	channelCap, ok := ms.ScopedKeeper.GetCapability(sdkCtx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, fmt.Errorf("module does not own channel capability")
	}

	// TODO(future): lock up the nameservice name in an escrow.

	packetData := types.ExamplePacketData{
		Sender:   sender,
		SomeData: someData,
	}
	// packetDataBz := types.EncodePacketData(packetData)

	sequence, err = ms.ics4Wrapper.SendPacket(sdkCtx, channelCap, sourcePort, sourceChannel, clienttypes.ZeroHeight(), timeoutTimestamp, packetData.GetBytes())
	if err != nil {
		return 0, err
	}
	return sequence, nil
}
