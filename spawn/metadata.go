package spawn

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

const (
	DefaultWebsite                   = "https://example.com"
	DefaultLogo                      = "https://raw.githubusercontent.com/cosmos/chain-registry/master/cosmoshub/images/atom.png"
	DefaultLogoSVG                   = "https://raw.githubusercontent.com/cosmos/chain-registry/master/cosmoshub/images/atom.svg"
	DefaultDescription               = "A short description of your project"
	DefaultChainID                   = "newchain-1"
	DefaultNetworkType               = "testnet" // or mainnet
	DefaultSlip44CoinType            = 118
	DefaultChainRegistrySchema       = "https://raw.githubusercontent.com/cosmos/chain-registry/master/chain.schema.json"
	DefaultChainRegistryAssetsSchema = "https://github.com/cosmos/chain-registry/blob/master/assetlist.schema.json"
	DefaultThemeHexColor             = "#FF2D00"
)

func (cfg *NewChainConfig) MetadataFile() MetadataFile {
	now := time.Now().UTC()
	now = now.Round(time.Minute)

	mf := MetadataFile{
		Token: TokenMeta{
			DisplayDenom:  strings.ToUpper(cfg.Denom),
			Denom:         cfg.Denom,
			Decimals:      6,
			Inflation:     "0.10",
			InitialSupply: "1000000000000000000",
			MaxSupply:     "1000000000000000000",
		},
		Project: ProjectMeta{
			Github:           cfg.GithubPath(),
			TargetLaunchDate: now,
			Logo:             DefaultLogo,
			Website:          DefaultWebsite,
			Description:      DefaultDescription,
			ShortDescription: DefaultDescription,
			Whitepaper:       "https://example.com/whitepaper.pdf",
			Contact: ContactMeta{
				Email:    "",
				Telegram: "",
				Twitter:  "",
				Name:     "",
				Discord:  "",
			},
		},
		ICS: ICSMeta{},
	}

	if cfg.IsFeatureEnabled(InterchainSecurity) {
		mf.ICS = ICSMeta{
			SpawnTime: now,
			Title:     cfg.BinDaemon,
			Summary:   DefaultDescription + " ( in .md format)",
			ChainID:   DefaultChainID,
			InitialHeight: ICSClientTypes{
				RevisionHeight: 0,
				RevisionNumber: 1,
			},
			UnbondingPeriod:                   21 * 24 * time.Hour.Nanoseconds(), // 21 days
			CcvTimeoutPeriod:                  28 * 24 * time.Hour.Nanoseconds(), // 28 days
			TransferTimeoutPeriod:             1 * time.Hour.Nanoseconds(),       // 1 hour (matches stride-1 and neutron-1)
			ConsumerRedistributionFraction:    "0.75",
			BlocksPerDistributionTransmission: 1_000,
			HistoricalEntries:                 10_000,
			GenesisHash:                       "",
			BinaryHash:                        "",
			DistributionTransmissionChannel:   "",
			TopN:                              0,
			ValidatorsPowerCap:                0,
			ValidatorSetCap:                   25,
			Allowlist:                         []any{},
			Denylist:                          []any{},
		}
	}

	return mf
}

type MetadataFile struct {
	ICS     ICSMeta     `json:"ics"`
	Token   TokenMeta   `json:"token"`
	Project ProjectMeta `json:"project"`
}

func (mf MetadataFile) SaveJSON(loc string) error {
	var bz []byte
	var err error

	if mf.ICS.IsZero() {
		// Non-ICS chains would save 0 state to file despite IsZero() being true.
		// Hacky override with a new type to save instead.
		type MetadataFileBare struct {
			Token   TokenMeta   `json:"token"`
			Project ProjectMeta `json:"project"`
		}

		bz, err = json.MarshalIndent(MetadataFileBare{
			Token:   mf.Token,
			Project: mf.Project,
		}, "", "  ")
	} else {
		bz, err = json.MarshalIndent(mf, "", "  ")
		if err != nil {
			return err
		}
	}

	return os.WriteFile(loc, bz, 0644)
}

type ICSClientTypes struct {
	// IBC clienttypes.Height just without omitempty

	// the revision that the client is currently on
	RevisionNumber uint64 `protobuf:"varint,1,opt,name=revision_number,json=revisionNumber,proto3" json:"revision_number"`
	// the height within the given revision
	RevisionHeight uint64 `protobuf:"varint,2,opt,name=revision_height,json=revisionHeight,proto3" json:"revision_height"`
}

type ICSMeta struct {
	SpawnTime                         time.Time      `json:"spawn_time"`
	Title                             string         `json:"title"`
	Summary                           string         `json:"summary"`
	ChainID                           string         `json:"chain_id"`
	InitialHeight                     ICSClientTypes `json:"initial_height"`
	UnbondingPeriod                   int64          `json:"unbonding_period"`
	CcvTimeoutPeriod                  int64          `json:"ccv_timeout_period"`
	TransferTimeoutPeriod             int64          `json:"transfer_timeout_period"`
	ConsumerRedistributionFraction    string         `json:"consumer_redistribution_fraction"`
	BlocksPerDistributionTransmission int            `json:"blocks_per_distribution_transmission"`
	HistoricalEntries                 int            `json:"historical_entries"`
	GenesisHash                       string         `json:"genesis_hash"`
	BinaryHash                        string         `json:"binary_hash"`
	DistributionTransmissionChannel   string         `json:"distribution_transmission_channel"`
	TopN                              int            `json:"top_N"`
	ValidatorsPowerCap                int            `json:"validators_power_cap"`
	ValidatorSetCap                   int            `json:"validator_set_cap"`
	Allowlist                         []any          `json:"allowlist"`
	Denylist                          []any          `json:"denylist"`
}

// impl IsZero on ICSMeta
func (ics ICSMeta) IsZero() bool {
	// can't compare []any, so checking most defaults
	return (ics.Title == "" &&
		ics.Summary == "" &&
		ics.ChainID == "")
}

type TokenMeta struct {
	DisplayDenom  string `json:"display_denom"`
	Denom         string `json:"denom"`
	Decimals      int    `json:"decimals"`
	Inflation     string `json:"inflation"`
	InitialSupply string `json:"initial_supply"`
	MaxSupply     string `json:"max_supply"`
}

type ProjectMeta struct {
	Github           string      `json:"github"`
	TargetLaunchDate time.Time   `json:"target_launch_date"`
	Logo             string      `json:"logo"`
	Website          string      `json:"website"`
	Description      string      `json:"description"`
	ShortDescription string      `json:"short_description"`
	Whitepaper       string      `json:"whitepaper"`
	Contact          ContactMeta `json:"contact"`
}

type ContactMeta struct {
	Email    string `json:"email"`
	Telegram string `json:"telegram"`
	Twitter  string `json:"twitter"`
	Name     string `json:"name"`
	Discord  string `json:"discord"`
}
