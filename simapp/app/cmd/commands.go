package cmd

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"cosmossdk.io/client/v2/offchain"
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/log"
	runtimev2 "cosmossdk.io/runtime/v2"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/server/v2/api/grpc"
	"cosmossdk.io/server/v2/store"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/rollchains/spawn/simapp" // TODO: rename me to just `github.com/rollchains/myunit`

	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func newApp[T transaction.Tx](
	logger log.Logger, viper *viper.Viper,
) serverv2.AppI[T] {
	return serverv2.AppI[T](simapp.NewSimApp[T](logger, viper))
}

func initRootCmd[T transaction.Tx](
	rootCmd *cobra.Command,
	clientCtx client.Context,
	moduleManager *runtimev2.MM[T],
	makeComponent func(cc client.Context) serverv2.ServerComponent[T],
) {
	cfg := sdk.GetConfig()
	cfg.Seal()

	rootCmd.AddCommand(
		genutilcli.InitCmd(moduleManager),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
		// pruning.Cmd(newApp), // TODO add to comet server
		// snapshot.Cmd(newApp), // TODO add to comet server
	)

	// TODO: ?
	// logger, err := serverv2.NewLogger(viper.New(), rootCmd.OutOrStdout())
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to create logger: %v", err))
	// }

	// add keybase, auxiliary RPC, query, genesis, and tx child commands
	rootCmd.AddCommand(
		genesisCommand(moduleManager, appExport),
		queryCommand(),
		txCommand(),
		keys.Commands(),
		offchain.OffChain(),
	)

	// wire server commands
	if err := serverv2.AddCommands(
		rootCmd,
		newApp,
		serverv2.DefaultServerConfig(),
		// logger,
		makeComponent(clientCtx),
		grpc.New[T](),
		store.New[T](),
	); err != nil {
		panic(err)
	}
}

// genesisCommand builds genesis-related `simd genesis` command. Users may provide application specific commands as a parameter
func genesisCommand[T transaction.Tx](
	moduleManager *runtimev2.MM[T],
	appExport servertypes.AppExporter,
	cmds ...*cobra.Command,
) *cobra.Command {
	// compatAppExporter := func(logger log.Logger,
	// 	db corestore.KVStoreWithBatch,
	// 	traceWriter io.Writer,
	// 	height int64,
	// 	forZeroHeight bool,
	// 	jailAllowedAddrs []string,
	// 	opts servertypes.AppOptions,
	// 	modulesToExport []string) (servertypes.ExportedApp, error) {
	// 	viperAppOpts, ok := appOpts.(*viper.Viper)
	// 	if !ok {
	// 		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	// 	}
	// 	return appExport(logger, height, forZeroHeight, jailAllowedAddrs, viperAppOpts, modulesToExport)
	// }

	cmd := genutilcli.Commands(moduleManager.Modules()[genutiltypes.ModuleName].(genutil.AppModule), moduleManager, appExport)
	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	return cmd
}

// appExport creates a new simapp (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db corestore.KVStoreWithBatch,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	// overwrite the FlagInvCheckPeriod
	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)

	simApp := simapp.NewSimApp(logger, viperAppOpts)
	// if height != -1 {
	// 	simApp = simapp.NewSimApp(logger, viperAppOpts)

	// 	if err := simApp.LoadHeight(height); err != nil {
	// 		return servertypes.ExportedApp{}, err
	// 	}
	// } else {
	// 	simApp = simapp.NewSimApp(logger, viperAppOpts)
	// }

	return simApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
