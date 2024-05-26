package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rollchains/spawn/spawn"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type disabledCase struct {
	Name          string
	Disabled      []string
	NotContainAny []string
}

func TestDisabledGeneration(t *testing.T) {
	allButStaking := make([]string, 0, len(spawn.AllFeatures)-1)
	for _, f := range spawn.AllFeatures {
		if f != spawn.Staking {
			allButStaking = append(allButStaking, f)
		}
	}

	disabledCases := []disabledCase{
		{
			// by default ICS is used & staking is disabled ()
			Name:          "default",
			Disabled:      []string{"staking"},
			NotContainAny: []string{"StakingKeeper", "POAKeeper"},
		},
		{
			Name:     "everything but staking",
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

	for _, c := range disabledCases {
		c := c
		execTest(t, c)
	}
}

func TestGenerationForSingleModuleRemoval(t *testing.T) {
	// single module removal
	disabledCases := make([]disabledCase, 0, len(spawn.AllFeatures))

	for _, f := range spawn.AllFeatures {
		normalizedName := strings.ReplaceAll("remove"+f, "-", "")

		disabledCases = append(disabledCases, disabledCase{
			Name:     normalizedName,
			Disabled: []string{f},
		})
	}

	for _, c := range disabledCases {
		c := c
		execTest(t, c)
	}
}

func execTest(t *testing.T, c disabledCase) {
	name := "spawnunittest" + c.Name
	dc := c.Disabled

	fmt.Println("=====\ndisabled cases", name, dc)

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		// dirPath := path.Join(cwd, name)

		require.NoError(t, os.RemoveAll(name))

		fSys := afero.NewMemMapFs()
		fSys.MkdirAll(name, 0777)
		fSys.Chmod(name, 0777)

		cfg := spawn.NewChainConfig{
			ProjectName:     name,
			Bech32Prefix:    "cosmos",
			HomeDir:         "." + name,
			BinDaemon:       RandStringBytes(6) + "d",
			Denom:           "token" + RandStringBytes(3),
			GithubOrg:       RandStringBytes(15),
			IgnoreGitInit:   true, // mem fs issues, no need
			DisabledModules: dc,
			Logger:          Logger,
			FileSystem:      fSys,
		}
		cfg.Run(false)

		AssertValidGeneration(t, dc, c.NotContainAny, fSys, name)

		require.NoError(t, os.RemoveAll(name))
	})
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func AssertValidGeneration(t *testing.T, dc []string, notContainAny []string, fSys afero.Fs, dir string) {
	fileCount := 0

	err := afero.Walk(fSys, dir, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileCount++

		if filepath.Ext(p) == ".go" {
			base := path.Base(p)

			f, err := fSys.Open(p)
			require.NoError(t, err, fmt.Sprintf("can't open %s", base))

			require.NoError(t, fSys.Chmod(p, 0777))

			bz, err := afero.ReadAll(f)
			require.NoError(t, err, fmt.Sprintf("can't read %s", base))

			// ensure no disabled modules are present
			for _, text := range notContainAny {
				text := text
				require.NotContains(t, string(bz), text, fmt.Sprintf("disabled module %s found in %s (%s)", text, base, p))
			}
		}

		return nil
	})

	require.NoError(t, err, fmt.Sprintf("error walking directory for disabled: %v", dc))
	require.Greater(t, fileCount, 1, fmt.Sprintf("no files found in %s", fSys.Name()))
}
