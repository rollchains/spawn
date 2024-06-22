package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"

	"github.com/rollchains/spawn/spawn"
)

const (
	FlagLocationPath  = "print-location"
	LocalICReleaseAPI = "https://api.github.com/repos/strangelove-ventures/interchaintest/releases"
)

func init() {
	LocalICCmd.Flags().Bool(FlagLocationPath, false, "print the location of local-ic binary")

	if !doRunCheck() {
		return
	}

	// TODO: make this non blocking (?) or save to a file for last run (check daily)
	curVer := LocalICCurrentVersion()
	if !strings.HasPrefix(curVer, "v") {
		fmt.Println("Local-IC version not found. Please run `make get-localic`")
		return
	}

	// in a go func run LocalICVersionCheck
	// LocalICVersionCheck()
	releases, err := GetLatestReleases(LocalICReleaseAPI)
	if err != nil {
		GetLogger().Error("Error getting latest local-ic releases", "err", err)
		return
	}

	if curVer != "" {
		curVer = strings.Split(curVer, "-")[0]
	}

	// TODO: latest of the v8 release line?
	latest := releases[0]

	if semver.Compare(curVer, latest.TagName) < 0 {
		GetLogger().Info(
			"New Local-IC version available",
			"latest", latest.TagName,
			"current", curVer,
		)
	}
}

// ---
// make install && ICTEST_HOME=./simapp spawn local-ic start testnet
// make install && cd simapp && spawn local-ic start testnet
// ---
var LocalICCmd = &cobra.Command{
	Use:   "local-ic",
	Short: "Local Interchain",
	Long:  "Wrapper for Local Interchain. Download with `make get-localic`",
	Example: `  - spawn local-ic chains
  - spawn local-ic start testnet
  - spawn local-ic interact localcosmos-1 query 'bank balances cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr'`,
	Run: func(cmd *cobra.Command, args []string) {
		debugBinaryLoc, _ := cmd.Flags().GetBool(FlagLocationPath)

		logger := GetLogger()

		loc := whereIsLocalICInstalled()
		if loc == "" {
			logger.Error("local-ic not found. Please run `make get-localic`")
			return
		}

		if debugBinaryLoc {
			logger.Debug("local-ic binary", "location", loc)
			return
		}

		if err := os.Chmod(loc, 0755); err != nil {
			logger.Error("Error setting local-ic permissions", "err", err)
		}

		// set to use the current dir if it is not overridden
		if os.Getenv("ICTEST_HOME") == "" {
			if err := os.Setenv("ICTEST_HOME", "."); err != nil {
				logger.Error("Error setting ICTEST_HOME", "err", err)
			}
		}

		if err := spawn.ExecCommand(loc, args...); err != nil {
			logger.Error("Error calling local-ic", "err", err)
		}
	},
}

func whereIsLocalICInstalled() string {
	for _, path := range []string{"local-ic", path.Join("bin", "local-ic"), path.Join("local-interchain", "localic")} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	if path, err := exec.LookPath("local-ic"); err == nil {
		return path
	}

	return ""
}

func LocalICCurrentVersion() string {
	// whereIsLocalICInstalled
	loc := whereIsLocalICInstalled()
	if loc == "" {
		// TODO: local-ic download?
		// TODO: change get-localic to take in a version?
		GetLogger().Error("local-ic not found. Please run `make get-localic`")
		return ""
	}

	output, err := spawn.ExecCommandWithOutput(loc, "version")
	if err != nil {
		fmt.Println("Error calling local-ic", err)
		return ""
	}

	return string(output)
}

func GetLatestReleases(apiRepoURL string) ([]Release, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, apiRepoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// parse response
	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}

type Release struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	TagName     string  `json:"tag_name"`
	PublishedAt string  `json:"published_at"`
	Assets      []Asset `json:"assets"`

	Prerelease bool `json:"prerelease"`
	Draft      bool `json:"draft"`
}

type Asset struct {
	URL      string `json:"url"`
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	Uploader struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}

// lastRunCheck returns true if it is time to run the check again.
// It saves the last run time to a file for future runs.
func doRunCheck() bool {
	logger := GetLogger()
	now := time.Now()

	// open the home dir + spawn
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Error getting user home dir", "err", err)
		return false
	}

	spawnDir := path.Join(home, ".spawn")
	if _, err := os.Stat(spawnDir); os.IsNotExist(err) {
		if err := os.Mkdir(spawnDir, 0755); err != nil {
			logger.Error("Error creating home spawn dir", "err", err)
			return false
		}
	}

	lastCheckFile := path.Join(spawnDir, "last_ver_check.txt")
	if _, err := os.Stat(lastCheckFile); os.IsNotExist(err) {
		if _, err := os.Create(lastCheckFile); err != nil {
			logger.Error("Error creating last check file", "err", err)
			return false
		}

		epoch := time.Unix(0, 0)
		if err := os.WriteFile(lastCheckFile, []byte(epoch.Format(time.RFC3339)), 0644); err != nil {
			logger.Error("Error writing last check file", "err", err)
			return false
		}
	}

	lastCheck, err := os.ReadFile(lastCheckFile)
	if err != nil {
		logger.Error("Error reading last check file", "err", err)
		return false
	}

	lastCheckTime, err := time.Parse(time.RFC3339, string(lastCheck))
	if err != nil {
		logger.Error("Error parsing last check time", "err", err)
		return false
	}

	diff := now.Sub(lastCheckTime)
	isTime := diff > 24*time.Hour

	if isTime {
		if err := os.WriteFile(lastCheckFile, []byte(now.Format(time.RFC3339)), 0644); err != nil {
			logger.Error("Error writing last check file", "err", err)
			return false
		}
	}

	return isTime
}
