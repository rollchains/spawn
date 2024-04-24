package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "example-plugin",
	Short: "Info About the Plugin",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

func main() {
	rootCmd.AddCommand(addCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A command you can use to perform addition of 2 numbers!",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		num1, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error parsing the first number")
			os.Exit(1)
		}
		num2, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error parsing the second number")
			os.Exit(1)
		}

		fmt.Println("add called")
		fmt.Println("Performing the addition of the following numbers: ")
		fmt.Printf("Num1: %v\n", num1)
		fmt.Printf("Num2: %v\n", num2)
		fmt.Printf("Addition of those 2 numbers is: %v\n", num1+num2)
	},
}
