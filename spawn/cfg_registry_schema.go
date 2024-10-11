package spawn

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/rollchains/spawn/spawn/types"
)

const (
	DefaultWebsite                   = "https://example.com"
	DefaultDiscord                   = "https://discord.gg/your-discord"
	DefaultEmail                     = "example@example.com"
	DefaultLogoPNG                   = "https://raw.githubusercontent.com/cosmos/chain-registry/master/cosmoshub/images/atom.png"
	DefaultLogoSVG                   = "https://raw.githubusercontent.com/cosmos/chain-registry/master/cosmoshub/images/atom.svg"
	DefaultDescription               = "A short description of your project"
	DefaultChainID                   = "localchain-1"
	DefaultNetworkType               = "testnet" // or mainnet
	DefaultSlip44CoinType            = 118
	DefaultChainRegistrySchema       = "https://raw.githubusercontent.com/cosmos/chain-registry/master/chain.schema.json"
	DefaultChainRegistryAssetsSchema = "https://github.com/cosmos/chain-registry/blob/master/assetlist.schema.json"
	DefaultThemeHexColor             = "#FF2D00"
)

var caser = cases.Title(language.English)

func (cfg NewChainConfig) ChainRegistryFile() types.ChainRegistryFormat {
	// TODO: update as needed
	DefaultSDKVersion := "0.50"
	DefaultTendermintVersion := "0.38"
	DefaultIBCGoVersion := "8"
	DefaultCosmWasmVersion := ""
	if cfg.IsFeatureEnabled(CosmWasm) {
		DefaultCosmWasmVersion = "0.50"
	}
	DefaultConsensus := "tendermint" // TODO: gordian in the future on gen

	return types.ChainRegistryFormat{
		Schema:       DefaultChainRegistrySchema,
		ChainName:    cfg.ProjectName,
		ChainType:    "cosmos",
		Status:       "live",
		Website:      DefaultWebsite,
		NetworkType:  DefaultNetworkType,
		PrettyName:   caser.String(cfg.ProjectName),
		ChainID:      DefaultChainID,
		Bech32Prefix: cfg.Bech32Prefix,
		DaemonName:   cfg.BinDaemon,
		NodeHome:     cfg.NodeHome(),
		KeyAlgos:     []string{"secp256k1"},
		Slip44:       DefaultSlip44CoinType,
		Fees: types.Fees{
			FeeTokens: []types.FeeTokens{
				{
					Denom:            cfg.Denom,
					FixedMinGasPrice: 0,
					LowGasPrice:      0,
					AverageGasPrice:  0.025,
					HighGasPrice:     0.04,
				},
			},
		},
		Codebase: types.Codebase{
			GitRepo:            "https://" + cfg.GithubPath(),
			RecommendedVersion: "v1.0.0",
			CompatibleVersions: []string{"v0.9.0"},
			CosmosSdkVersion:   DefaultSDKVersion,
			Consensus: types.Consensus{
				Type:    DefaultConsensus,
				Version: DefaultTendermintVersion,
			},
			CosmwasmVersion: DefaultCosmWasmVersion,
			CosmwasmEnabled: cfg.IsFeatureEnabled(CosmWasm),
			IbcGoVersion:    DefaultIBCGoVersion,
			IcsEnabled:      []string{"ics20-1"},
			Genesis: types.Genesis{
				Name:       "v1",
				GenesisURL: fmt.Sprintf("https://%s/%s", cfg.GithubPath(), "networks/raw/main/genesis.json"),
			},
			Versions: []types.Versions{
				{
					Name:            "v1.0.0",
					Tag:             "v1.0.0",
					Height:          0,
					NextVersionName: "v2",
				},
			},
		},
		Staking: types.Staking{
			StakingTokens: []types.StakingTokens{
				{
					Denom: cfg.Denom,
				},
			},
			LockDuration: types.LockDuration{
				Time: "1814400s", // 21 days
			},
		},
		Images: []types.Images{
			{
				Png: DefaultLogoPNG,
				Theme: types.Theme{
					PrimaryColorHex: DefaultThemeHexColor,
				},
			},
		},
		Peers: types.Peers{},
		Apis: types.Apis{
			RPC: []types.RPC{
				{
					Address:  "http://127.0.0.1:26657",
					Provider: "localhost",
				},
			},
			Rest: []types.Rest{
				{
					Address:  "http://127.0.0.1:1317",
					Provider: "localhost",
				},
			},
		},
		Explorers: []types.Explorers{
			{
				Kind:        "cosmos",
				URL:         "https://example.com",
				TxPage:      "https://example.com/tx",
				AccountPage: "https://example.com/account",
			},
		},
		Keywords: []string{"cosmos", "spawn"},
	}
}

// The ICS MetadataFile is similar to this.
func (cfg NewChainConfig) ChainRegistryAssetsFile() types.ChainRegistryAssetsList {
	display := strings.TrimPrefix(strings.ToUpper(cfg.Denom), "U")

	return types.ChainRegistryAssetsList{
		Schema:    DefaultChainRegistryAssetsSchema,
		ChainName: cfg.ProjectName,
		Assets: []types.Assets{
			{
				Description: "The native token of " + cfg.ProjectName,
				DenomUnits: []types.DenomUnits{
					{
						Denom:    cfg.Denom, // utoken
						Exponent: 0,
					},
					{
						Denom:    display, // TOKEN
						Exponent: 6,
					},
				},
				Base:    cfg.Denom, // utoken
				Name:    fmt.Sprintf("%s %s", cfg.ProjectName, display),
				Display: strings.ToLower(display), // token
				Symbol:  display,                  // TOKEN
				LogoURIs: types.LogoURIs{
					Png: DefaultLogoPNG,
					Svg: DefaultLogoSVG,
				},
				Images: []types.ImagesAssetLists{
					{
						Png: DefaultLogoPNG,
						Svg: DefaultLogoSVG,
						Theme: types.Theme{
							PrimaryColorHex: DefaultThemeHexColor,
						},
					},
				},
				Socials: types.Socials{
					Website: DefaultWebsite,
					Twitter: "https://x.com/cosmoshub",
				},
			},
		},
	}
}
