package main

import (
	"testing"

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

			res := CleanDisabled(tc.disabled)
			if !tc.panics {
				require.EqualValues(t, tc.expected, res, "expected %v, got %v", tc.expected, res)
			}
		})
	}
}
