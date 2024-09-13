package main_test

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	main "github.com/rollchains/spawn/cmd/spawn"
	"github.com/rollchains/spawn/spawn"
	"github.com/stretchr/testify/require"
)

func TestModuleGeneration(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	cfg := spawn.NewChainConfig{
		ProjectName:     "default",
		Bech32Prefix:    "cosmos",
		HomeDir:         ".default",
		BinDaemon:       main.RandStringBytes(6) + "d",
		Denom:           "token" + main.RandStringBytes(3),
		GithubOrg:       main.RandStringBytes(15),
		IgnoreGitInit:   false,
		DisabledModules: []string{"explorer"},
		Logger:          main.Logger,
	}

	type mc struct {
		Name           string
		Args           []string
		OutputContains string
	}

	mcs := []mc{
		{
			Name: "iibcmid",
			Args: []string{"new", "myibc", "--ibc-module"},
		},
		{
			Name: "iibcmod",
			Args: []string{"new", "myibcmw", "--ibc-middleware"},
		},
		{
			Name: "standard",
			Args: []string{"new", "standard"},
		},
	}

	for _, c := range mcs {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			name := "spawnmoduleunittest" + c.Name

			cfg.ProjectName = name
			cfg.HomeDir = "." + name
			fmt.Println("=====\nName", name)

			dirPath := path.Join(cwd, name)
			require.NoError(t, os.RemoveAll(name))

			require.NoError(t, cfg.ValidateAndRun(false), "failed to generate proper chain")

			// move to new repo
			require.NoError(t, os.Chdir(dirPath))

			cmd := main.ModuleCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(c.Args)
			cmd.Execute()
			// out, err := io.ReadAll(b)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// TODO: this is not being read from. Fix.
			// require.Contains(t, string(out), c.OutputContains, "output: "+string(out))

			// validate the go source is good
			main.AssertValidGeneration(t, dirPath, nil, nil, cfg)

			require.NoError(t, os.Chdir(cwd))
			require.NoError(t, os.RemoveAll(name))
		})
	}
}
