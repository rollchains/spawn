package main

import (
	"fmt"
	"go/format"
	"io/fs"
	"log/slog"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rollchains/spawn/spawn"
	"github.com/stretchr/testify/require"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func TestDisabledGeneration(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	allButStaking := make([]string, 0, len(spawn.AllFeatures)-1)
	for _, f := range spawn.AllFeatures {
		if f != spawn.POS {
			allButStaking = append(allButStaking, f)
		}
	}

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
			Disabled: allButStaking,
		},
		{
			Name:          "noibcaddons",
			Disabled:      []string{spawn.PacketForward, spawn.IBCRateLimit},
			NotContainAny: []string{"packetforward", "RatelimitKeeper"},
		},
		{
			Name:          "stdmix1",
			Disabled:      []string{spawn.GlobalFee, spawn.Ignite, spawn.TokenFactory},
			NotContainAny: []string{"TokenFactoryKeeper", "GlobalFeeKeeper"},
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

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func AssertValidGeneration(t *testing.T, dirPath string, dc []string, notContainAny []string, cfg spawn.NewChainConfig) {
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
				text := text
				require.NotContains(t, string(f), text, fmt.Sprintf("disabled module %s found in %s (%s) with config %+v", text, base, p, cfg))
			}

			_, err = format.Source(f)
			require.NoError(t, err, fmt.Sprintf("format issue: %v. using disabled: %v", base, dc))
		}

		return nil
	})
	require.NoError(t, err, fmt.Sprintf("error walking directory for disabled: %v", dc))
	require.Greater(t, fileCount, 1, fmt.Sprintf("no files found in %s", dirPath))
}
