package keeper_test

import (
	"testing"

	apiv1 "github.com/rollchains/simapp/api/example/v1"
	"github.com/stretchr/testify/require"
)

func TestORM(t *testing.T) {
	f := SetupTest(t)

	dt := f.k.OrmDB.ExampleDataTable()
	acc := []byte("test_acc")
	amt := uint64(7)

	err := dt.Insert(f.ctx, &apiv1.ExampleData{
		Account: acc,
		Amount:  amt,
	})
	require.NoError(t, err)

	d, err := dt.Has(f.ctx, []byte("test_acc"))
	require.NoError(t, err)
	require.True(t, d)

	// get it
	res, err := dt.Get(f.ctx, []byte("test_acc"))
	require.NoError(t, err)
	require.NotNil(t, res)
	require.EqualValues(t, amt, res.Amount)
}
