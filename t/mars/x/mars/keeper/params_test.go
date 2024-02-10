package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "mars/testutil/keeper"
	"mars/x/mars/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.MarsKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
