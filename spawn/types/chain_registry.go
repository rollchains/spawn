package types

import (
	"encoding/json"
	"os"
)

// TOOL: https://mholt.github.io/json-to-go/
// Manual Modifications:
// - Consensus: add omitempty
// - Codebase.Consensus: add omitempty
// - Peers: add omitempty
// - added: ChainRegistryFormat ChainType string

type ChainRegistryFormat struct {
	Schema       string      `json:"$schema"`
	ChainName    string      `json:"chain_name"`
	ChainType    string      `json:"chain_type"`
	Status       string      `json:"status"`
	Website      string      `json:"website"`
	NetworkType  string      `json:"network_type"`
	PrettyName   string      `json:"pretty_name"`
	ChainID      string      `json:"chain_id"`
	Bech32Prefix string      `json:"bech32_prefix"`
	DaemonName   string      `json:"daemon_name"`
	NodeHome     string      `json:"node_home"`
	KeyAlgos     []string    `json:"key_algos"`
	Slip44       int         `json:"slip44"`
	Fees         Fees        `json:"fees"`
	Staking      Staking     `json:"staking"`
	Codebase     Codebase    `json:"codebase"`
	Images       []Images    `json:"images"`
	Peers        Peers       `json:"peers"`
	Apis         Apis        `json:"apis"`
	Explorers    []Explorers `json:"explorers"`
	Keywords     []string    `json:"keywords"`
}
type FeeTokens struct {
	Denom            string  `json:"denom"`
	FixedMinGasPrice int     `json:"fixed_min_gas_price"`
	LowGasPrice      int     `json:"low_gas_price"`
	AverageGasPrice  float64 `json:"average_gas_price"`
	HighGasPrice     float64 `json:"high_gas_price"`
}
type Fees struct {
	FeeTokens []FeeTokens `json:"fee_tokens"`
}
type StakingTokens struct {
	Denom string `json:"denom"`
}
type LockDuration struct {
	Time string `json:"time"`
}
type Staking struct {
	StakingTokens []StakingTokens `json:"staking_tokens"`
	LockDuration  LockDuration    `json:"lock_duration"`
}
type Consensus struct {
	Type    string `json:"type,omitempty"`
	Version string `json:"version,omitempty"`
}
type Genesis struct {
	Name       string `json:"name"`
	GenesisURL string `json:"genesis_url"`
}
type Versions struct {
	Name               string    `json:"name"`
	Tag                string    `json:"tag"`
	Height             int       `json:"height"`
	NextVersionName    string    `json:"next_version_name"`
	Proposal           int       `json:"proposal,omitempty"`
	RecommendedVersion string    `json:"recommended_version,omitempty"`
	CompatibleVersions []string  `json:"compatible_versions,omitempty"`
	CosmosSdkVersion   string    `json:"cosmos_sdk_version,omitempty"`
	Consensus          Consensus `json:"consensus,omitempty"`
	CosmwasmVersion    string    `json:"cosmwasm_version,omitempty"`
	CosmwasmEnabled    bool      `json:"cosmwasm_enabled,omitempty"`
	IbcGoVersion       string    `json:"ibc_go_version,omitempty"`
	IcsEnabled         []string  `json:"ics_enabled,omitempty"`
}
type Codebase struct {
	GitRepo            string     `json:"git_repo"`
	RecommendedVersion string     `json:"recommended_version"`
	CompatibleVersions []string   `json:"compatible_versions"`
	CosmosSdkVersion   string     `json:"cosmos_sdk_version"`
	Consensus          Consensus  `json:"consensus,omitempty"`
	CosmwasmVersion    string     `json:"cosmwasm_version"`
	CosmwasmEnabled    bool       `json:"cosmwasm_enabled"`
	IbcGoVersion       string     `json:"ibc_go_version"`
	IcsEnabled         []string   `json:"ics_enabled"`
	Genesis            Genesis    `json:"genesis"`
	Versions           []Versions `json:"versions"`
}
type Theme struct {
	PrimaryColorHex string `json:"primary_color_hex"`
}
type Images struct {
	Png   string `json:"png"`
	Theme Theme  `json:"theme"`
}
type Seeds struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Provider string `json:"provider"`
}
type PersistentPeers struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Provider string `json:"provider"`
}
type Peers struct {
	Seeds           []Seeds           `json:"seeds,omitempty"`
	PersistentPeers []PersistentPeers `json:"persistent_peers,omitempty"`
}
type RPC struct {
	Address  string `json:"address"`
	Provider string `json:"provider"`
}
type Rest struct {
	Address  string `json:"address"`
	Provider string `json:"provider"`
}
type Apis struct {
	RPC  []RPC  `json:"rpc"`
	Rest []Rest `json:"rest"`
}
type Explorers struct {
	Kind        string `json:"kind"`
	URL         string `json:"url"`
	TxPage      string `json:"tx_page"`
	AccountPage string `json:"account_page"`
}

func (v ChainRegistryFormat) SaveJSON(loc string) error {
	bz, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(loc, bz, 0666)
}
