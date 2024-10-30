package main

import (
	"fmt"
	"os"

	clientv2helpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/core/transaction"
	serverv2 "cosmossdk.io/server/v2"
	"github.com/rollchains/spawn/simapp/app/cmd"
)

func main() {
	homeDir := "/tmp/foo"

	rootCmd := cmd.NewRootCmd[transaction.Tx](homeDir)
	if err := serverv2.Execute(rootCmd, clientv2helpers.EnvPrefix, homeDir); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
