package types

import "cosmossdk.io/collections"

var (
	// ParamsKey saves the current module params.
	ParamsKey = collections.NewPrefix(0)
)

const (
	ModuleName = "example"

	// TODO: let's just only ues ModuleName for all these instead of aliasing? or is there some reflection reason to have this.
	RouterKey = ModuleName

	StoreKey = ModuleName

	QuerierRoute = ModuleName
)
