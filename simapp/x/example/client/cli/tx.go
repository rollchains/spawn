package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/strangelove-ventures/simapp/x/example/types"
)

// NewTxCmd returns a root CLI command handler for certain modules/Clock
// transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      types.ModuleName + " subcommands.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewRegisterClockContract(),
	)
	return txCmd
}

// NewRegisterClockContract returns a CLI command handler for registering a
// contract for the clock module.
func NewRegisterClockContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [some-value]",
		Short: "Update the params (must be submitted from the authority)",
		Long:  "Register a clock contract. Sender must be admin of the contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := cliCtx.GetFromAddress()

			someValue, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateParams{
				Authority: senderAddress.String(),
				Params: types.Params{
					SomeValue: someValue,
				},
			}

			if err := msg.Validate(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
