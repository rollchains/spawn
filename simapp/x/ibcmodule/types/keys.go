package types

const (
	// ModuleName defines the name of module.
	ModuleName = "ibcmodule"

	// PortID defines the port ID that module module binds to.
	PortID = ModuleName

	// Version defines the current version the IBC module supports
	Version = ModuleName + "-1"

	StoreKey = ModuleName

	EventTypePacket = "example_data_packet"
)
