package spawn_test

import (
	"testing"

	"github.com/rollchains/spawn/spawn"
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
			disabled: []string{"poa", "globalfee", "cosmwasm"},
			expected: []string{"poa", "globalfee", "cosmwasm"},
		},
		{
			name:     "remove poa duplicate",
			disabled: []string{"poa", "globalfee", "cosmwasm", "poa"},
			expected: []string{"poa", "globalfee", "cosmwasm"},
		},
		{
			name:     "remove poa and globalfee duplicate",
			disabled: []string{"poa", "globalfee", "cosmwasm", "poa", "globalfee"},
			expected: []string{"poa", "globalfee", "cosmwasm"},
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
			disabled: []string{spawn.POA, spawn.GlobalFee, spawn.CosmWasm},
			expected: []string{spawn.POA, spawn.GlobalFee, spawn.CosmWasm},
		},
		{
			name:     "fix-all",
			disabled: []string{"proof-of-authority", "global-fee", "cw"},
			expected: []string{spawn.POA, spawn.GlobalFee, spawn.CosmWasm},
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
			disabled: []string{spawn.Staking},
			expected: []string{spawn.POA, spawn.Staking},
			parentDepPairs: map[string][]string{
				spawn.Staking: {spawn.POA},
			},
		},
		{
			name:     "remove what is expected",
			disabled: []string{spawn.POA, spawn.Staking},
			expected: []string{spawn.POA, spawn.Staking},
			parentDepPairs: map[string][]string{
				spawn.Staking: {spawn.POA},
			},
		},
		{
			name:     "remove ics",
			disabled: []string{spawn.InterchainSecurity, spawn.POA},
			expected: []string{spawn.InterchainSecurity, spawn.POA},
			parentDepPairs: map[string][]string{
				spawn.Staking: {spawn.POA},
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
