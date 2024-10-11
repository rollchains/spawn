package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/rollchains/spawn/spawn"
	"github.com/stretchr/testify/require"
)

func TestDisabledGeneration(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	type disabledCase struct {
		Name          string
		Disabled      []string
		NotContainAny []string
	}

	// custom cases
	// NOTE: block-explorer is disabled for all cases.
	disabledCases := []disabledCase{
		{
			// by default ICS is used
			Name:          "default",
			Disabled:      []string{},
			NotContainAny: []string{"POAKeeper"},
		},
		{
			Name:     "everythingbutstaking",
			Disabled: AllFeaturesButStaking(),
		},
		{
			Name:          "noibcaddons",
			Disabled:      []string{spawn.PacketForward, spawn.IBCRateLimit},
			NotContainAny: []string{"packetforward", "RatelimitKeeper"},
		},
		{
			Name:          "stdmix1",
			Disabled:      []string{spawn.TokenFactory},
			NotContainAny: []string{"TokenFactoryKeeper"},
		},
		{
			Name:          "nocw",
			Disabled:      []string{spawn.CosmWasm, spawn.WasmLC},
			NotContainAny: []string{"wasmkeeper", "wasmtypes"},
		},
	}

	// single module removal
	for _, f := range spawn.AllFeatures {
		normalizedName := strings.ReplaceAll("remove"+f, "-", "")

		disabledCases = append(disabledCases, disabledCase{
			Name:     normalizedName,
			Disabled: []string{f, spawn.BlockExplorer},
		})
	}

	for _, c := range disabledCases {
		c := c
		name := "spawnunittest" + c.Name
		dc := append(c.Disabled, spawn.BlockExplorer)

		fmt.Println("=====\ndisabled cases", name, dc)

		t.Run(name, func(t *testing.T) {
			dirPath := path.Join(cwd, name)
			require.NoError(t, os.RemoveAll(name))

			cfg := spawn.NewChainConfig{
				ProjectName:     name,
				Bech32Prefix:    "cosmos",
				HomeDir:         "." + name,
				BinDaemon:       RandStringBytes(6) + "d",
				Denom:           "token" + RandStringBytes(3),
				GithubOrg:       RandStringBytes(15),
				IgnoreGitInit:   false,
				DisabledModules: dc,
				Logger:          Logger,
			}
			require.NoError(t, cfg.ValidateAndRun(false), "failed to generate proper chain")

			AssertValidGeneration(t, dirPath, dc, c.NotContainAny, cfg)
			require.NoError(t, os.RemoveAll(name))
		})
	}
}
