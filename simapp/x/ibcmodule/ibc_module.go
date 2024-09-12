package ibcmodule

import (
	"context"
	"fmt"
	"strings"

	"github.com/rollchains/spawn/simapp/x/ibcmodule/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/rollchains/spawn/simapp/x/ibcmodule/types"
)

var _ porttypes.IBCModule = (*ExampleIBCModule)(nil)

// ExampleIBCModule implements all the callbacks
// that modules must define as specified in ICS-26
type ExampleIBCModule struct {
	keeper keeper.Keeper
}

// NewExampleIBCModule creates a new IBCModule given the keeper and underlying application.
func NewExampleIBCModule(k keeper.Keeper) ExampleIBCModule {
	return ExampleIBCModule{
		keeper: k,
	}
}

func (im ExampleIBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	if strings.TrimSpace(version) == "" {
		version = types.Version
	}

	// if order != channeltypes.UNORDERED {
	// 	return "", fmt.Errorf("invalid channel order; expected UNORDERED")
	// }

	if counterparty.PortId != types.ModuleName {
		return "", fmt.Errorf("invalid counterparty port ID; expected %s, got %s", types.ModuleName, counterparty.PortId)
	}

	// OpenInit must claim the channelCapability that IBC passes into the callback
	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", fmt.Errorf("failed to claim capability: %w", err)
	}

	return version, nil
}

func (im ExampleIBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID, channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (version string, err error) {
	// OpenTry must claim the channelCapability that IBC passes into the callback
	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return "", err
	}

	if counterpartyVersion != types.Version {
		fmt.Println("invalid counterparty version, proposing current app version", "counterpartyVersion", counterpartyVersion, "version", types.Version)
		return types.Version, nil // TODO: err here?
	}

	return types.Version, nil
}

func (im ExampleIBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID, channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	if counterpartyVersion != types.Version {
		return fmt.Errorf("invalid counterparty version: expected %s, got %s", types.Version, counterpartyVersion)
	}
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface.
func (im ExampleIBCModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface.
func (im ExampleIBCModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return fmt.Errorf("channel close is disabled for this module")
}

// OnChanCloseConfirm implements the IBCModule interface.
func (im ExampleIBCModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface.
func (im ExampleIBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	logger := im.keeper.Logger(ctx)
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	var data types.ExamplePacketData
	var ackErr error
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		ackErr = fmt.Errorf("cannot unmarshal example packet data: %v", err)
		logger.Error(fmt.Sprintf("%s sequence %d", ackErr.Error(), packet.Sequence))
		ack = channeltypes.NewErrorAcknowledgement(ackErr)
	}

	// only attempt the application logic if the packet data was successfully decoded
	if ack.Success() {
		// TODO: perform your logic here
		err := im.handleOnRecvLogic(ctx, data)
		if err != nil {
			ack = channeltypes.NewErrorAcknowledgement(err)
			ackErr = err
			logger.Error(fmt.Sprintf("%s sequence %d", ackErr.Error(), packet.Sequence))
		} else {
			logger.Info("successfully handled example packet", "sequence", packet.Sequence)
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, data.Sender),
			sdk.NewAttribute("some_data", data.SomeData),
			sdk.NewAttribute("ack_success", fmt.Sprintf("%t", ack.Success())),
		),
	)

	return ack
}

func (im ExampleIBCModule) handleOnRecvLogic(ctx context.Context, data types.ExamplePacketData) error {
	v, err := im.keeper.ExampleStore.Get(ctx)
	if err != nil {
		return err
	}

	err = im.keeper.ExampleStore.Set(ctx, v+1)
	if err != nil {
		return err
	}

	return nil
}

// OnAcknowledgementPacket implements the IBCModule interface.
func (im ExampleIBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return fmt.Errorf("cannot unmarshal example packet acknowledgement: %v", err)
	}

	var data types.ExamplePacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return fmt.Errorf("cannot unmarshal example packet data: %v", err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, data.Sender),
			sdk.NewAttribute("some_data", data.SomeData),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute("success", string(resp.Result)),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute("error", resp.Error),
			),
		)
	}

	return nil
}

// OnTimeoutPacket implements the IBCMiddleware interface.
func (im ExampleIBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var data types.ExamplePacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return fmt.Errorf("cannot unmarshal example packet data: %v", err)
	}

	// Handle timeout logic here as necessary (i.e. refunds for example) or nothing at all.

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"timeout",
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute("sender", data.Sender),
		),
	)

	return nil
}
