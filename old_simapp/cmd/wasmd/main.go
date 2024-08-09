package main

import (
	"os"

	"cosmossdk.io/log"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/rollchains/spawn/simapp/app"
)

func main() {
	rootCmd := NewRootCmd()
	rootCmd.AddCommand(AddConsumerSectionCmd(app.DefaultNodeHome)) // spawntag:ics

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
