package main

import (
	"fmt"
	"go/format"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rollchains/spawn/spawn"
	"github.com/stretchr/testify/require"
)

type DisabledCase struct {
	Name          string
	Disabled      []string
	NotContainAny []string
}

func TestDisabledGeneration(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// proof-of-authority,tokenfactory,globalfee,ibc-packetforward,ibc-ratelimit,cosmwasm,wasm-light-client,interchain-security,ignite-cli

	disabledCases := []DisabledCase{
		// {
		// 	Name:     "mix1",
		// 	Disabled: []string{"globalfee", "wasmlc", "ignite"},
		// },
		// {
		// 	Name:     "ibcmix",
		// 	Disabled: []string{"packetforward", "ibc-rate-limit"},
		// },
		// {
		// 	Name:     "cwmix",
		// 	Disabled: []string{"cosmwasm", "cosmwasm", "wasm-light-client"},
		// },
	}

	for _, f := range spawn.AllFeatures {
		normalizedName := strings.ReplaceAll("remove"+f, "-", "")

		disabledCases = append(disabledCases, DisabledCase{
			Name:     normalizedName,
			Disabled: []string{f},
		})
	}

	for _, c := range disabledCases {
		name := "spawnunittest" + c.Name
		dc := c.Disabled

		fmt.Println("=====\ndisabled cases", name, dc)

		t.Run(name, func(t *testing.T) {
			dirPath := path.Join(cwd, name)

			require.NoError(t, os.RemoveAll(name))

			cfg := spawn.NewChainConfig{
				ProjectName:     name,
				Bech32Prefix:    "cosmos",
				HomeDir:         ".projName",
				BinDaemon:       "simd",
				Denom:           "token",
				GithubOrg:       "rollchains",
				IgnoreGitInit:   false,
				DisabledModules: dc,
				Logger:          slog.New(slog.NewJSONHandler(os.Stdout, nil)),
			}
			cfg.Run(false)

			assetValidGeneration(t, dirPath, dc, c.NotContainAny)

			require.NoError(t, os.RemoveAll(name))
		})
	}
}

func TestDisabledFuzzer(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	fmt.Println("=====\ndisabled fuzzer", cwd)

}

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

func assetValidGeneration(t *testing.T, dirPath string, dc []string, notContainAny []string) {
	fileCount := 0
	err := filepath.WalkDir(dirPath, func(p string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileCount++

		if filepath.Ext(p) == ".go" {
			base := path.Base(p)

			f, err := os.ReadFile(p)
			require.NoError(t, err, fmt.Sprintf("can't read %s", base))

			// ensure no disabled modules are present
			for _, text := range notContainAny {
				require.NotContains(t, string(f), text, fmt.Sprintf("disabled module %s found in %s", text, base))
			}

			_, err = format.Source(f)
			require.NoError(t, err, fmt.Sprintf("format issue: %v. using disabled: %v", base, dc))
		}

		return nil
	})
	require.NoError(t, err, fmt.Sprintf("error walking directory for disabled: %v", dc))
	require.Greater(t, fileCount, 1, fmt.Sprintf("no files found in %s", dirPath))
}
