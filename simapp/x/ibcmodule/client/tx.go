package cli

import (
	"time"

	"github.com/rollchains/spawn/simapp/x/ibcmodule/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// NewTxCmd creates and returns the tx command
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "ibcmodule",
		Short:                      "ibcmodule subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewSomeDataTxCmd(),
	)

	return cmd
}

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
)

var defaultTimeout = uint64((time.Duration(10) * time.Minute).Nanoseconds())

// NewSomeDataTxCmd
func NewSomeDataTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example-tx [src-port] [src-channel] [data]",
		Short: "Send a packet with some data",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress().String()
			srcPort := args[0]
			srcChannel := args[1]
			someData := args[2]

			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}

			if timeoutTimestamp != 0 {
				now := time.Now().UnixNano()
				timeoutTimestamp = uint64(now + time.Duration(1*time.Hour).Nanoseconds())
			}

			msg := &types.MsgSendExampleTx{
				Sender:           sender,
				SourcePort:       srcPort,
				SourceChannel:    srcChannel,
				TimeoutTimestamp: timeoutTimestamp,
				SomeData:         someData,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, defaultTimeout, "Packet timeout timestamp in nanoseconds from now. Default is 10 minutes.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
