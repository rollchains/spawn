package types

const (
	// ModuleName defines the name of module.
	ModuleName = "ibcmodule"

	// PortID defines the port ID that module module binds to.
	PortID = ModuleName

	// Version defines the current version the IBC module supports
	Version = ModuleName + "-1"

	// StoreKey is the store key string for the module.
	StoreKey = ModuleName

	// RouterKey is the message route for the module.
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the module.
	QuerierRoute = ModuleName

	EventTypePacket = "example_data_packet"
)
