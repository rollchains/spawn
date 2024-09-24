package spawn_test

import (
	"testing"

	"github.com/rollchains/spawn/spawn"
	"github.com/rollchains/spawn/spawn/types"
	"github.com/stretchr/testify/require"
)

func TestDisabled(t *testing.T) {
	type tcase struct {
		name     string
		disabled []string
		expected []string
		panics   bool
	}

	testCases := []tcase{
		{
			name:     "same",
			disabled: []string{"poa", "cosmwasm"},
			expected: []string{"poa", "cosmwasm"},
		},
		{
			name:     "remove poa duplicate",
			disabled: []string{"poa", "cosmwasm", "poa"},
			expected: []string{"poa", "cosmwasm"},
		},
		{
			name:     "remove poa duplicate",
			disabled: []string{"poa", "cosmwasm", "poa"},
			expected: []string{"poa", "cosmwasm"},
		},
		{
			name:     "panic due to invalid disabled feature",
			disabled: []string{"poa", "whatiamnotreal", "cosmwasm"},
			panics:   true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.panics {
						t.Errorf("expected no panic, but got %v", r)
					}
				}
			}()

			res := spawn.RemoveDuplicates(tc.disabled)
			if !tc.panics {
				require.Len(t, res, len(tc.expected))

				// ensure every element within tc.expected is in res (ignore order)
				found := make(map[string]bool)
				for _, e := range tc.expected {
					found[e] = false
				}

				for _, r := range res {
					if _, ok := found[r]; ok {
						found[r] = true
					}
				}

				for k, v := range found {
					require.True(t, v, "expected %s to be found in res", k)
				}
			}
		})
	}
}

func TestNormalizedNames(t *testing.T) {
	type tcase struct {
		name           string
		disabled       []string
		expected       []string
		parentDepPairs map[string][]string
		panics         bool
	}

	testCases := []tcase{
		{
			name:     "normal for both",
			disabled: []string{spawn.POA, spawn.CosmWasm},
			expected: []string{spawn.POA, spawn.CosmWasm},
		},
		{
			name:     "fix-all",
			disabled: []string{"proof-of-authority", "cw"},
			expected: []string{spawn.POA, spawn.CosmWasm},
		},
		{
			name:     "remove duplicate",
			disabled: []string{"proof-of-authority", "proof-of-authority"},
			expected: []string{spawn.POA},
		},
		{
			name:     "incorrect",
			disabled: []string{"notanoption"},
			panics:   true,
		},
		{
			name:     "incorrect with allowed",
			disabled: []string{spawn.POA, "notanoption"},
			panics:   true,
		},
		{
			name:     "remove staking and POA due to parentDeps",
			disabled: []string{spawn.POS},
			expected: []string{spawn.POA, spawn.POS},
			parentDepPairs: map[string][]string{
				spawn.POS: {spawn.POA},
			},
		},
		{
			name:     "remove what is expected",
			disabled: []string{spawn.POA, spawn.POS},
			expected: []string{spawn.POA, spawn.POS},
			parentDepPairs: map[string][]string{
				spawn.POS: {spawn.POA},
			},
		},
		{
			name:     "remove ics",
			disabled: []string{spawn.InterchainSecurity, spawn.POA},
			expected: []string{spawn.InterchainSecurity, spawn.POA},
			parentDepPairs: map[string][]string{
				spawn.POS: {spawn.POA},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.panics {
						t.Errorf("expected no panic, but got %v", r)
					}
				}
			}()

			res := spawn.NormalizeDisabledNames(tc.disabled, tc.parentDepPairs)
			if !tc.panics {
				require.Len(t, res, len(tc.expected))

				// ensure every element within tc.expected is in res (ignore order)
				found := make(map[string]bool)
				for _, e := range tc.expected {
					found[e] = false
				}

				for _, r := range res {
					if _, ok := found[r]; ok {
						found[r] = true
					}
				}

				for k, v := range found {
					require.True(t, v, "expected %s to be found in res", k)
				}
			}
		})
	}
}

type cfgCase struct {
	Desc string
	Cfg  spawn.NewChainConfig
	Err  error
}

func NewCfgCase(desc string, cfg spawn.NewChainConfig, expectedErr error) cfgCase {
	return cfgCase{
		Desc: desc,
		Cfg:  cfg,
		Err:  expectedErr,
	}
}

const (
	proj  = "myproject"
	bech  = "cosmos"
	home  = ".app"
	bin   = "appd"
	denom = "token"
	org   = "myorg"
)

func goodCfg() spawn.NewChainConfig {
	return spawn.NewChainConfig{
		ProjectName:  proj,
		Bech32Prefix: bech,
		HomeDir:      home,
		BinDaemon:    bin,
		Denom:        denom,
		GithubOrg:    org,
	}
}

func TestBadConfigInputs(t *testing.T) {
	chainCases := []cfgCase{
		NewCfgCase("valid config", goodCfg(), nil),
		NewCfgCase("no github org", goodCfg().WithOrg(""), types.ErrCfgEmptyOrg),
		NewCfgCase("no project name", goodCfg().WithProjectName(""), types.ErrCfgEmptyProject),
		NewCfgCase("project special chars -", goodCfg().WithProjectName("my-project"), types.ErrCfgProjSpecialChars),
		NewCfgCase("project special chars /", goodCfg().WithProjectName("my/project"), types.ErrCfgProjSpecialChars),
		NewCfgCase("binary name to short len 1", goodCfg().WithBinDaemon("a"), types.ErrCfgBinTooShort),
		NewCfgCase("success: binary name len 2", goodCfg().WithBinDaemon("ad"), nil),
		NewCfgCase("token denom too short len 1", goodCfg().WithDenom("a"), types.ErrCfgDenomTooShort),
		NewCfgCase("token denom too short len 2", goodCfg().WithDenom("ab"), types.ErrCfgDenomTooShort),
		NewCfgCase("success: token denom special chars", goodCfg().WithDenom("my-cool/token"), nil),
		NewCfgCase("success: token denom 3", goodCfg().WithDenom("abc"), nil),
		NewCfgCase("home dir too short", goodCfg().WithHomeDir("."), types.ErrCfgHomeDirTooShort),
		NewCfgCase("success: home dir valid", goodCfg().WithHomeDir(".a"), nil),
		NewCfgCase("bech32 prefix to short", goodCfg().WithBech32Prefix(""), types.ErrCfgEmptyBech32),
		NewCfgCase("bech32 not alpha", goodCfg().WithBech32Prefix("c919"), types.ErrCfgBech32Alpha),
		NewCfgCase("bech32 not alpha", goodCfg().WithBech32Prefix("1"), types.ErrCfgBech32Alpha),
		NewCfgCase("bech32 not alpha", goodCfg().WithBech32Prefix("---"), types.ErrCfgBech32Alpha),
		NewCfgCase("success: bech32 prefix", goodCfg().WithBech32Prefix("c"), nil),
	}

	for _, c := range chainCases {
		c := c

		t.Run(c.Desc, func(t *testing.T) {

			err := c.Cfg.Validate()
			if c.Err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), c.Err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestChainRegistry(t *testing.T) {
	cfg := goodCfg()
	cr := cfg.ChainRegistryFile()
	require.Equal(t, cfg.ProjectName, cr.ChainName)
	require.Equal(t, bech, cr.Bech32Prefix)
	require.Equal(t, bin, cr.DaemonName)
	require.Equal(t, denom, cr.Fees.FeeTokens[0].Denom)
}
