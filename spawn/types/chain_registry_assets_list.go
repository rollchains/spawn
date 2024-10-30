package types

import (
	"encoding/json"
	"os"
)

// Update: Images -> ImagesAssetLists
// - Fix: Assets.Images -> new type ImagesAssetLists

type ChainRegistryAssetsList struct {
	Schema    string   `json:"$schema"`
	ChainName string   `json:"chain_name"`
	Assets    []Assets `json:"assets"`
}
type DenomUnits struct {
	Denom    string `json:"denom"`
	Exponent int    `json:"exponent"`
}
type LogoURIs struct {
	Png string `json:"png"`
	Svg string `json:"svg"`
}
type ImagesAssetLists struct {
	Png   string `json:"png"`
	Svg   string `json:"svg"`
	Theme Theme  `json:"theme,omitempty"`
}
type Assets struct {
	Description string             `json:"description"`
	DenomUnits  []DenomUnits       `json:"denom_units"`
	Base        string             `json:"base"`
	Name        string             `json:"name"`
	Display     string             `json:"display"`
	Symbol      string             `json:"symbol"`
	LogoURIs    LogoURIs           `json:"logo_URIs"`
	Images      []ImagesAssetLists `json:"images"`
	Socials     Socials            `json:"socials,omitempty"`
}

type Socials struct {
	Website string `json:"website,omitempty"`
	Twitter string `json:"twitter,omitempty"`
}

func (v ChainRegistryAssetsList) SaveJSON(loc string) error {
	bz, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(loc, bz, 0666)
}
