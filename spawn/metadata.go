package spawn

import (
	"encoding/json"
	"os"
)

type (
	MetadataFile struct {
		Display Display `json:"display"`
	}
	Display struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Links       Links   `json:"links"`
		Widget      *Widget `json:"widget,omitempty"`
	}
	Links struct {
		Logo       string `json:"logo"`
		Discord    string `json:"discord"`
		Email      string `json:"email"`
		Github     string `json:"github"`
		Telegram   string `json:"telegram"`
		Twitter    string `json:"twitter"`
		Website    string `json:"website"`
		Whitepaper string `json:"whitepaper"`
	}
	Widget struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		ButtonText  string `json:"buttonText,omitempty"`
		ButtonURL   string `json:"buttonUrl,omitempty"`
	}
)

func (cfg *NewChainConfig) MetadataFile() MetadataFile {
	mf := MetadataFile{
		Display: Display{
			Name:        cfg.ProjectName,
			Description: cfg.ProjectName + " is an Interchain blockchain.",
			Links: Links{
				Logo:       DefaultLogoPNG,
				Discord:    DefaultDiscord,
				Email:      DefaultEmail,
				Github:     "https://" + cfg.GithubPath(),
				Telegram:   "https://t.me/example",
				Twitter:    "https://twitter.com/example_account",
				Website:    DefaultWebsite,
				Whitepaper: "https://bitcoin.org/bitcoin.pdf",
			},
		}}

	if cfg.IsFeatureEnabled(InterchainSecurity) {
		mf.Display.Widget = &Widget{
			Title:       cfg.ProjectName,
			Description: cfg.ProjectName + " is a Interchain blockchain.",
			ButtonText:  "Learn More",
			ButtonURL:   DefaultWebsite,
		}
	}

	return mf
}

func (mf MetadataFile) SaveJSON(loc string) error {
	bz, err := json.MarshalIndent(mf, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(loc, bz, 0666)
}
