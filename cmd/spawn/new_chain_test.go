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

type DisabledCase struct {
	Name          string
	Disabled      []string
	NotContainAny []string
}

func TestDisabledGeneration(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	allButStaking := make([]string, 0, len(spawn.AllFeatures)-1)
	for _, f := range spawn.AllFeatures {
		if f != spawn.Staking {
			allButStaking = append(allButStaking, f)
		}
	}

	// custom cases
	disabledCases := []DisabledCase{
		{
			Name:     "onlystaking",
			Disabled: allButStaking,
		},
		{
			Name:     "stdmix1",
			Disabled: []string{spawn.GlobalFee, spawn.Ignite, spawn.TokenFactory},
		},
		{
			Name:     "noibcaddons",
			Disabled: []string{spawn.PacketForward, spawn.IBCRateLimit},
		},
		{
			Name:     "nocw",
			Disabled: []string{spawn.CosmWasm, spawn.WasmLC},
		},
	}

	// single module removal
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
				HomeDir:         "." + name,
				BinDaemon:       RandStringBytes(6) + "d",
				Denom:           "token" + RandStringBytes(3),
				GithubOrg:       RandStringBytes(15),
				IgnoreGitInit:   false,
				DisabledModules: dc,
				Logger:          Logger,
			}
			cfg.Run(false)

			AssertValidGeneration(t, dirPath, dc, c.NotContainAny)

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

func AssertValidGeneration(t *testing.T, dirPath string, dc []string, notContainAny []string) {
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
