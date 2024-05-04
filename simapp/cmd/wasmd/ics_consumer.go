package main

import (
	"encoding/json"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	types1 "github.com/cometbft/cometbft/abci/types"
	pvm "github.com/cometbft/cometbft/privval"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
	ccvtypes "github.com/cosmos/interchain-security/v5/x/ccv/types"
	"github.com/spf13/cobra"

	ibctypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ccvprovidertypes "github.com/cosmos/interchain-security/v5/x/ccv/provider/types"
)

// original credit: https://github.com/Stride-Labs/stride/blob/v22.0.0/cmd/consumer.go
func AddConsumerSectionCmd(nodeHome string) *cobra.Command {
	genesisMutator := NewDefaultGenesisIO()

	txCmd := &cobra.Command{
		Use:                        "add-consumer-section [chainID]",
		Args:                       cobra.ExactArgs(1),
		Short:                      "ONLY FOR TESTING PURPOSES! Modifies genesis so that chain can be started locally with one node.",
		SuggestionsMinimumDistance: 2,
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID := args[0]

			overrideHome, err := cmd.Flags().GetString(flags.FlagHome)
			if err != nil {
				return err
			}

			if overrideHome != "" {
				nodeHome = overrideHome
			}

			return genesisMutator.AlterConsumerModuleState(cmd, nodeHome, func(state *GenesisData, _ map[string]json.RawMessage) error {
				initialValset := []types1.ValidatorUpdate{}
				genesisState := CreateMinimalConsumerTestGenesis(chainID)
				serverCtx := server.GetServerContextFromCmd(cmd)
				config := serverCtx.Config
				config.SetRoot(nodeHome)

				privValidator := pvm.LoadFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
				pk, err := privValidator.GetPubKey()
				if err != nil {
					return err
				}
				sdkPublicKey, err := cryptocodec.FromCmtPubKeyInterface(pk)
				if err != nil {
					return err
				}
				tmProtoPublicKey, err := cryptocodec.ToCmtProtoPublicKey(sdkPublicKey)
				if err != nil {
					return err
				}

				initialValset = append(initialValset, types1.ValidatorUpdate{PubKey: tmProtoPublicKey, Power: 100})

				vals, err := tmtypes.PB2TM.ValidatorUpdates(initialValset)
				if err != nil {
					return errorsmod.Wrap(err, "could not convert val updates to validator set")
				}

				genesisState.Provider.InitialValSet = initialValset
				genesisState.Provider.ConsensusState.NextValidatorsHash = tmtypes.NewValidatorSet(vals).Hash()

				state.ConsumerModuleState = genesisState
				return nil
			})
		},
	}

	txCmd.Flags().String(flags.FlagHome, nodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(txCmd)

	return txCmd
}

type GenesisMutator interface {
	AlterConsumerModuleState(cmd *cobra.Command, callback func(state *GenesisData, appState map[string]json.RawMessage) error) error
}

type DefaultGenesisIO struct {
	DefaultGenesisReader
}

func NewDefaultGenesisIO() *DefaultGenesisIO {
	return &DefaultGenesisIO{DefaultGenesisReader: DefaultGenesisReader{}}
}

func (x DefaultGenesisIO) AlterConsumerModuleState(cmd *cobra.Command, homeDir string, callback func(state *GenesisData, appState map[string]json.RawMessage) error) error {
	g, err := x.ReadGenesis(cmd, homeDir)
	if err != nil {
		return err
	}
	if err := callback(g, g.AppState); err != nil {
		return err
	}
	if err := g.ConsumerModuleState.Validate(); err != nil {
		return err
	}
	clientCtx := client.GetClientContextFromCmd(cmd)
	consumerGenStateBz, err := clientCtx.Codec.MarshalJSON(g.ConsumerModuleState)
	if err != nil {
		return errorsmod.Wrap(err, "marshal consumer genesis state")
	}

	g.AppState[ccvconsumertypes.ModuleName] = consumerGenStateBz
	appStateJSON, err := json.Marshal(g.AppState)
	if err != nil {
		return errorsmod.Wrap(err, "marshal application genesis state")
	}

	g.GenDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(g.GenDoc, g.GenesisFile)
}

type DefaultGenesisReader struct{}

func (d DefaultGenesisReader) ReadGenesis(cmd *cobra.Command, homeDir string) (*GenesisData, error) {
	serverCtx := server.GetServerContextFromCmd(cmd)
	config := serverCtx.Config
	config.SetRoot(homeDir)

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	return NewGenesisData(
		genFile,
		genDoc,
		appState,
		nil,
	), nil
}

type GenesisData struct {
	GenesisFile         string
	GenDoc              *genutiltypes.AppGenesis
	AppState            map[string]json.RawMessage
	ConsumerModuleState *ccvtypes.ConsumerGenesisState
}

func NewGenesisData(genesisFile string, genDoc *genutiltypes.AppGenesis, appState map[string]json.RawMessage, consumerModuleState *ccvtypes.ConsumerGenesisState) *GenesisData {
	return &GenesisData{GenesisFile: genesisFile, GenDoc: genDoc, AppState: appState, ConsumerModuleState: consumerModuleState}
}

// This function creates consumer module genesis state that is used as starting point for modifications
// that allows the ICS chain to be started locally without having to start the provider chain and the relayer.
// It is also used in tests that are starting the chain node.
func CreateMinimalConsumerTestGenesis(chainID string) *ccvtypes.ConsumerGenesisState {
	genesisState := ccvtypes.DefaultConsumerGenesisState()
	genesisState.Params.Enabled = true
	genesisState.NewChain = true
	genesisState.Provider.ClientState = ccvprovidertypes.DefaultParams().TemplateClient
	genesisState.Provider.ClientState.ChainId = chainID
	genesisState.Provider.ClientState.LatestHeight = ibctypes.Height{RevisionNumber: 0, RevisionHeight: 1}
	trustPeriod, err := ccvtypes.CalculateTrustPeriod(genesisState.Params.UnbondingPeriod, ccvprovidertypes.DefaultTrustingPeriodFraction)
	if err != nil {
		panic("provider client trusting period error")
	}
	genesisState.Provider.ClientState.TrustingPeriod = trustPeriod
	genesisState.Provider.ClientState.UnbondingPeriod = genesisState.Params.UnbondingPeriod
	genesisState.Provider.ClientState.MaxClockDrift = ccvprovidertypes.DefaultMaxClockDrift
	genesisState.Provider.ConsensusState = &ibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
		Root:      commitmenttypes.MerkleRoot{Hash: []byte("dummy")},
	}

	return genesisState
}
