package spawn

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"
)

func TestLoadingValues(t *testing.T) {
	// fmt.Println(DefaultSDKVersion)
	// fmt.Println(DefaultCosmWasmVersion)
	// fmt.Println(DefaultTendermintVersion)
	// fmt.Println(DefaultIBCGoVersion)
	require.True(t, semver.Compare(DefaultSDKVersion, "0.50.0") >= 0)
	require.True(t, semver.Compare(DefaultCosmWasmVersion, "0.50.0") >= 0)
	require.True(t, semver.Compare(DefaultTendermintVersion, "0.38.0") >= 0)
	require.True(t, semver.Compare(DefaultIBCGoVersion, "8.2.0") >= 0)
}
