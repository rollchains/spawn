package keeper

import (
	"mars/x/mars/types"
)

var _ types.QueryServer = Keeper{}
