package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "example-plugin",
	Short: "Info About the spawn example-plugin",
}

func main() {
	rootCmd.AddCommand(AddCmd(), FlagTestCmd())

	// hides 'completion' command
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func FlagTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "flags-test",
		Short:   "Test using flags with a plugin",
		Example: `spawn plugin example-plugin flags-test -- --value 7`,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			myValue, _ := cmd.Flags().GetInt("value")
			fmt.Printf("my-value: %v", myValue)
		},
	}

	cmd.Flags().Int("value", 0, "A value you can set")

	return cmd
}

func AddCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "add",
		Short:   "A command you can use to perform addition of 2 numbers!",
		Example: `spawn plugin example-plugin add 1 2`,
		Args:    cobra.ExactArgs(2),
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
}
