package spawn

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/rollchains/spawn/spawn/types"
)

var caser = cases.Title(language.English)

func (cfg NewChainConfig) ChainRegistryFile() types.ChainRegistryFormat {
	DefaultSDKVersion := MustParseVersionFromGoMod("github.com/cosmos/cosmos-sdk", true)
	DefaultTendermintVersion := MustParseVersionFromGoMod("github.com/cometbft/cometbft", true)
	DefaultIBCGoVersion := MustParseVersionFromGoMod("github.com/cosmos/ibc-go/v8", true)

	DefaultCosmWasmVersion, err := ParseVersionFromGoMod("github.com/CosmWasm/wasmd", true)
	if err != nil {
		DefaultCosmWasmVersion = ""
	}

	return types.ChainRegistryFormat{
		Schema:       DefaultChainRegistrySchema,
		ChainName:    cfg.ProjectName,
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
			// TODO: versions should be fetched from the repo go.mod
			GitRepo:            "https://" + cfg.GithubPath(),
			RecommendedVersion: "v1.0.0",
			CompatibleVersions: []string{"v0.9.0"},
			CosmosSdkVersion:   DefaultSDKVersion,
			Consensus: types.Consensus{
				Type:    "tendermint", // TODO: gordian in the future on gen
				Version: DefaultTendermintVersion,
			},
			CosmwasmVersion: DefaultCosmWasmVersion,
			CosmwasmEnabled: cfg.IsCosmWasmEnabled(),
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
				Png: DefaultLogo,
				Theme: types.Theme{
					PrimaryColorHex: DefaultThemeHexColor,
				},
			},
		},
		Peers: types.Peers{},
		Apis: types.Apis{
			RPC: []types.RPC{
				{
					Address:  "tcp://127.0.0.1:26657",
					Provider: "localhost",
				},
			},
			Rest: []types.Rest{
				{
					Address:  "tcp://127.0.0.1:1317",
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
					Png: DefaultLogo,
					Svg: DefaultLogoSVG,
				},
				Images: []types.ImagesAssetLists{
					{
						Png: DefaultLogo,
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
