package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

var _ ibcexported.PacketData = (*ExamplePacketData)(nil)

// GetPacketSender implements exported.PacketData.
func (epd *ExamplePacketData) GetPacketSender(sourcePortID string) string {
	return epd.Sender
}

// GetBytes returns the sorted JSON encoding of the packet data.
func (epd ExamplePacketData) GetBytes() ([]byte, error) {
	bz, err := codec.ProtoMarshalJSON(&epd, ModuleCdc.InterfaceRegistry())
	if err != nil {
		return nil, err
	}

	return sdk.MustSortJSON(bz), nil
}

// GetBytes must return the sorted JSON encoding of the packet data.
func (epd ExamplePacketData) MustGetBytes() []byte {
	bz, err := epd.GetBytes()
	if err != nil {
		panic(err)
	}
	return bz
}
