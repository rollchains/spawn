package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/rollchains/spawn/spawn"
)

const (
	FlagLocationPath = "print-location"
)

func init() {
	LocalICCmd.Flags().Bool(FlagLocationPath, false, "print the location of local-ic binary")

}

// ---
// make install && ICTEST_HOME=./simapp spawn local-ic start testnet
// make install && cd simapp && spawn local-ic start testnet
// ---
var LocalICCmd = &cobra.Command{
	Use:   "local-ic",
	Short: "Local Interchain",
	Long:  "Wrapper for Local Interchain. Download with `make get-localic`",
	Example: `  - spawn local-ic chains
  - spawn local-ic start testnet
  - spawn local-ic interact localcosmos-1 query 'bank balances cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr'`,
	Run: func(cmd *cobra.Command, args []string) {
		debugBinaryLoc, _ := cmd.Flags().GetBool(FlagLocationPath)

		logger := GetLogger()

		loc := spawn.WhereIsBinInstalled("local-ic")
		if debugBinaryLoc {
			logger.Debug("local-ic binary", "location", loc)
			return
		}

		if err := os.Chmod(loc, 0755); err != nil {
			logger.Error("Error setting local-ic permissions", "err", err)
		}

		// set to use the current dir if it is not overridden
		if os.Getenv("ICTEST_HOME") == "" {
			if err := os.Setenv("ICTEST_HOME", "."); err != nil {
				logger.Error("Error setting ICTEST_HOME", "err", err)
			}
		}

		if err := spawn.ExecCommand(loc, args...); err != nil {
			logger.Error("Error calling local-ic", "err", err)
		}
	},
}
