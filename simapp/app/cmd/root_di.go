package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/legacy"
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"cosmossdk.io/runtime/v2"
	serverv2 "cosmossdk.io/server/v2"

	"github.com/rollchains/gordian/gcosmos/gccodec"
	"github.com/rollchains/gordian/gcosmos/gserver"

	"cosmossdk.io/server/v2/cometbft"
	"cosmossdk.io/x/auth/tx"
	authtxconfig "cosmossdk.io/x/auth/tx/config"
	"cosmossdk.io/x/auth/types"
	"github.com/rollchains/spawn/simapp" // TODO: rename me to just `github.com/rollchains/myunit`

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/std"
)

func NewRootCmdWithServer[T transaction.Tx](
	makeComponent func(cc client.Context) serverv2.ServerComponent[T],
) *cobra.Command {
	return newRootCmd[T](makeComponent)
}

var cmd *cobra.Command

// NewRootCmd creates a new root command for simd. It is called once in the main function.
func NewRootCmd[T transaction.Tx]() *cobra.Command {
	// <spawntag:cometbft
	cmd = NewRootCmdWithServer(func(cc client.Context) serverv2.ServerComponent[T] {
		return cometbft.New[T](
			&genericTxDecoder[T]{cc.TxConfig},
			cometbft.DefaultServerOptions[T](),
		)
	})
	// spawntag:cometbft>
	// <spawntag:gordian
	ctx := context.Background()
	srvCtx := server.NewDefaultContext()

	cmd = NewRootCmdWithServer(func(cc client.Context) serverv2.ServerComponent[transaction.Tx] {
		ctx = context.WithValue(ctx, client.ClientContextKey, client.Context{
			ChainID: "gcosmos", // TODO:
			HomeDir: simapp.DefaultNodeHome,
		})

		log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		codec := gccodec.NewTxDecoder(cc.TxConfig)
		c, err := gserver.NewComponent(ctx, log, codec, cc.Codec)
		if err != nil {
			panic(err)
		}
		return c
	})
	if err := server.SetCmdServerContext(cmd, srvCtx); err != nil {
		panic(err)
	}
	// spawntag:gordian>
	return cmd
}

func newRootCmd[T transaction.Tx](
	makeComponent func(cc client.Context) serverv2.ServerComponent[T],
) *cobra.Command {
	var (
		autoCliOpts   autocli.AppOptions
		moduleManager *runtime.MM[T]
		clientCtx     client.Context
	)

	if err := depinject.Inject(
		depinject.Configs(
			simapp.AppConfig(),
			depinject.Supply(log.NewNopLogger()),
			depinject.Provide(
				codec.ProvideInterfaceRegistry,
				codec.ProvideAddressCodec,
				codec.ProvideProtoCodec,
				codec.ProvideLegacyAmino,
				ProvideClientContext,
			),
			depinject.Invoke(
				std.RegisterInterfaces,
				std.RegisterLegacyAminoCodec,
			),
		),
		&autoCliOpts,
		&moduleManager,
		&clientCtx,
	); err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:           "simd",
		Short:         "simulation app",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx = clientCtx.WithCmdContext(cmd.Context())
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			customClientTemplate, customClientConfig := initClientConfig()
			clientCtx, err = config.CreateClientConfig(clientCtx, customClientTemplate, customClientConfig)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(clientCtx, cmd); err != nil {
				return err
			}

			return nil
		},
	}

	initRootCmd[T](rootCmd, clientCtx, moduleManager, makeComponent)
	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}

func ProvideClientContext(
	appCodec codec.Codec,
	interfaceRegistry codectypes.InterfaceRegistry,
	txConfigOpts tx.ConfigOptions,
	legacyAmino legacy.Amino,
	addressCodec address.Codec,
	validatorAddressCodec address.ValidatorAddressCodec,
	consensusAddressCodec address.ConsensusAddressCodec,
) client.Context {
	var err error

	amino, ok := legacyAmino.(*codec.LegacyAmino)
	if !ok {
		panic("legacy.Amino must be an *codec.LegacyAmino instance for legacy ClientContext")
	}

	clientCtx := client.Context{}.
		WithCodec(appCodec).
		WithInterfaceRegistry(interfaceRegistry).
		WithLegacyAmino(amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithAddressCodec(addressCodec).
		WithValidatorAddressCodec(validatorAddressCodec).
		WithConsensusAddressCodec(consensusAddressCodec).
		WithHomeDir(simapp.DefaultNodeHome).
		WithViper("") // uses by default the binary name as prefix

	// Read the config to overwrite the default values with the values from the config file
	customClientTemplate, customClientConfig := initClientConfig()
	clientCtx, err = config.CreateClientConfig(clientCtx, customClientTemplate, customClientConfig)
	if err != nil {
		panic(err)
	}

	// textual is enabled by default, we need to re-create the tx config grpc instead of bank keeper.
	txConfigOpts.TextualCoinMetadataQueryFn = authtxconfig.NewGRPCCoinMetadataQueryFn(clientCtx)
	txConfig, err := tx.NewTxConfigWithOptions(clientCtx.Codec, txConfigOpts)
	if err != nil {
		panic(err)
	}
	clientCtx = clientCtx.WithTxConfig(txConfig)

	return clientCtx
}
