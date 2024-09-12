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
	"testing"

	"github.com/rollchains/spawn/spawn"
	"github.com/stretchr/testify/require"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

func AllFeaturesButStaking() []string {
	allButStaking := make([]string, 0, len(spawn.AllFeatures)-1)
	for _, f := range spawn.AllFeatures {
		if f != spawn.POS {
			allButStaking = append(allButStaking, f)
		}
	}
	return allButStaking
}
