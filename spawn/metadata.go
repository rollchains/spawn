package spawn

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

func (cfg *NewChainConfig) MetadataFile() MetadataFile {
	now := time.Now()

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
			Logo:             "https://example.com/logo.png",
			Website:          "https://example.com",
			Description:      "A short description of your project",
			ShortDescription: "A short description of your project",
			Whitepaper:       "https://example.com/whitepaper.pdf",
			Contact: ContactMeta{
				Email:    "",
				Telegram: "",
				Twitter:  "",
				Name:     "",
				Discord:  "",
			},
		},
	}

	if cfg.isUsingICS {
		mf.ICS = ICSMeta{
			SpawnTime:   now,
			Title:       cfg.BinDaemon,
			Description: ".md description of your chain and all other relevant information",
			ChainID:     "newchain-1",
			InitialHeight: ICSInitialHeight{
				RevisionHeight: 0,
				RevisionNumber: 1,
			},
			UnbondingPeriod:                   86400000000000,
			CcvTimeoutPeriod:                  259200000000000,
			TransferTimeoutPeriod:             1800000000000,
			ConsumerRedistributionFraction:    "0.75",
			BlocksPerDistributionTransmission: 1000,
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
	ICS     ICSMeta     `json:"ics,omitempty"`
	Token   TokenMeta   `json:"token,omitempty"`
	Project ProjectMeta `json:"project,omitempty"`
}

func (mf MetadataFile) SaveJSON(loc string) error {
	bz, err := json.MarshalIndent(mf, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(loc, bz, 0644)
}

type ICSInitialHeight struct {
	RevisionHeight int `json:"revision_height"`
	RevisionNumber int `json:"revision_number"`
}

type ICSMeta struct {
	SpawnTime                         time.Time        `json:"spawn_time"`
	Title                             string           `json:"title"`
	Description                       string           `json:"description"`
	ChainID                           string           `json:"chain_id"`
	InitialHeight                     ICSInitialHeight `json:"initial_height"`
	UnbondingPeriod                   int64            `json:"unbonding_period"`
	CcvTimeoutPeriod                  int64            `json:"ccv_timeout_period"`
	TransferTimeoutPeriod             int64            `json:"transfer_timeout_period"`
	ConsumerRedistributionFraction    string           `json:"consumer_redistribution_fraction"`
	BlocksPerDistributionTransmission int              `json:"blocks_per_distribution_transmission"`
	HistoricalEntries                 int              `json:"historical_entries"`
	GenesisHash                       string           `json:"genesis_hash"`
	BinaryHash                        string           `json:"binary_hash"`
	DistributionTransmissionChannel   string           `json:"distribution_transmission_channel"`
	TopN                              int              `json:"top_N"`
	ValidatorsPowerCap                int              `json:"validators_power_cap"`
	ValidatorSetCap                   int              `json:"validator_set_cap"`
	Allowlist                         []any            `json:"allowlist"`
	Denylist                          []any            `json:"denylist"`
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
