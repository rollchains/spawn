package types

import (
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ ibcexported.PacketData = (*ExamplePacketData)(nil)

// GetPacketSender implements exported.PacketData.
func (epd *ExamplePacketData) GetPacketSender(sourcePortID string) string {
	panic("unimplemented")
}

// GetBytes is a helper for serialising
func (epd ExamplePacketData) GetBytes() []byte {
	// return sdk.MustSortJSON(mustProtoMarshalJSON(&ftpd))
	return nil
}
