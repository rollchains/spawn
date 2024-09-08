package spawn

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type (
	ChainExplorerAsset struct {
		Base        string `json:"base"`
		Symbol      string `json:"symbol"`
		Exponent    string `json:"exponent"`
		CoingeckoId string `json:"coingecko_id"`
		Logo        string `json:"logo"`
	}

	Endpoint struct {
		Provider string `json:"provider"`
		Address  string `json:"address"`
	}

	ChainExplorer struct {
		ChainName  string               `json:"chain_name"`
		Api        []Endpoint           `json:"api"`
		Rpc        []Endpoint           `json:"rpc"`
		SdkVersion string               `json:"sdk_version"`
		CoinType   string               `json:"coin_type"`
		MinTxFee   string               `json:"min_tx_fee"`
		AddrPrefix string               `json:"addr_prefix"`
		Logo       string               `json:"logo"`
		ThemeColor string               `json:"theme_color"`
		Assets     []ChainExplorerAsset `json:"assets"`
	}
)

// hacky: pingpub does not have a docker file for some reason...
const dockerFile = `# docker build . -t pingpub:latest

FROM node:20-alpine

RUN apk add --no-cache yarn npm

WORKDIR /app

COPY . .

# install node_modules to the image
RUN yarn --ignore-engines

CMD [ "yarn", "--ignore-engines", "serve", "--host", "0.0.0.0" ]`

func NewEndpoint(provider, address string) Endpoint {
	return Endpoint{
		Provider: provider,
		Address:  address,
	}
}

func (cfg NewChainConfig) NewPingPubExplorer() error {
	if err := os.Chdir(cfg.ProjectName); err != nil {
		cfg.Logger.Error("chdir", "err", err)
	}
	if err := ExecCommand("git", "clone", "https://github.com/ping-pub/explorer.git", "--depth", "1"); err != nil {
		cfg.Logger.Error("git clone", "err", err)
	}
	if err := os.Chdir(".."); err != nil {
		cfg.Logger.Error("chdir", "err", err)
	}

	mainnet := path.Join(cfg.ProjectName, "explorer", "chains", "mainnet")
	cfg.clearDir(mainnet)
	cfg.clearDir(path.Join(cfg.ProjectName, "explorer", "chains", "testnet"))

	// Create JSON config file for explorer
	explorer := cfg.NewChainExplorerConfig()
	bz, err := json.MarshalIndent(explorer, "", "  ")
	if err != nil {
		cfg.Logger.Error("Error marshalling chain explorer config", "err", err)
	}

	err = os.WriteFile(path.Join(mainnet, fmt.Sprintf("%s.json", cfg.ProjectName)), bz, 0644)
	if err != nil {
		cfg.Logger.Error("Error writing chain explorer config", "err", err)
	}

	err = os.WriteFile(path.Join(cfg.ProjectName, "explorer", "Dockerfile"), []byte(dockerFile), 0644)
	if err != nil {
		cfg.Logger.Error("Error writing docker file", "err", err)
	}

	return nil
}

func (cfg NewChainConfig) NewChainExplorerConfig() ChainExplorer {
	logo := "https://img.freepik.com/free-vector/letter-s-box-logo-design-template_474888-3345.jpg?size=338&ext=jpg&ga=GA1.1.2008272138.1721260800&semt=ais_user"
	return ChainExplorer{
		ChainName:  cfg.ProjectName,
		Api:        []Endpoint{NewEndpoint("api.localhost", "http://127.0.0.1:1317")},
		Rpc:        []Endpoint{NewEndpoint("rpc.localhost", "http://127.0.0.1:26657")},
		SdkVersion: "0.50",
		CoinType:   "118",
		MinTxFee:   "800",
		AddrPrefix: cfg.Bech32Prefix,
		Logo:       logo,
		ThemeColor: "#001be7",
		Assets: []ChainExplorerAsset{
			{
				Base:        cfg.Denom,
				Symbol:      strings.ToUpper(cfg.Denom),
				Exponent:    "6",
				CoingeckoId: "",
				Logo:        logo,
			},
		},
	}
}

func (cfg NewChainConfig) clearDir(dirLoc ...string) {
	dir := path.Join(dirLoc...)

	if err := os.RemoveAll(dir); err != nil {
		cfg.Logger.Error("Error removing directory", "err", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		cfg.Logger.Error("Error creating directory", "err", err)
	}
}
