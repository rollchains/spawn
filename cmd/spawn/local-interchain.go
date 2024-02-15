package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"

	"gitub.com/strangelove-ventures/spawn/spawn"
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
	// Args:  cobra.
	Run: func(cmd *cobra.Command, args []string) {
		debugBinaryLoc, _ := cmd.Flags().GetBool(FlagLocationPath)

		loc := whereIsLocalICInstalled()
		if loc == "" {
			fmt.Println("local-ic not found. Please run `make get-localic`")
			return
		}

		if debugBinaryLoc {
			fmt.Println(loc)
			return
		}

		if err := os.Chmod(loc, 0755); err != nil {
			fmt.Println("Error setting local-ic permissions:", err)
		}

		// set to use the current dir if it is not overrriden
		if os.Getenv("ICTEST_HOME") == "" {
			if err := os.Setenv("ICTEST_HOME", "."); err != nil {
				fmt.Println("Error setting ICTEST_HOME:", err)
			}
		}

		if err := spawn.ExecCommand(loc, args...); err != nil {
			fmt.Println("Error calling local-ic:", err)
		}
	},
}

func whereIsLocalICInstalled() string {
	for _, path := range []string{"local-ic", path.Join("bin", "local-ic"), path.Join("local-interchain", "localic")} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	if path, err := exec.LookPath("local-ic"); err == nil {
		return path
	}

	return ""
}
