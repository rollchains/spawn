package spawn

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

var (
	RunCheckInterval   = 24 * time.Hour
	howToInstallBinary = map[string]string{
		"local-ic": "git clone https://github.com/strangelove-ventures/interchaintest.git ictd --depth=1 -b __VERSION__ && cd ictd/local-interchain && make install && cd ../.. && rm -rf ictd",
		"spawn":    "git clone https://github.com/rollchains/spawn.git --depth=1 -b __VERSION__ && cd spawn && make install && cd .. && rm -rf spawn",
	}
	BinaryToGithubAPI = map[string]string{
		"local-ic": "https://api.github.com/repos/strangelove-ventures/interchaintest/releases",
		"spawn":    "https://api.github.com/repos/rollchains/spawn/releases",
	}
)

type (
	Release struct {
		Id          int64   `json:"id"`
		Name        string  `json:"name"`
		TagName     string  `json:"tag_name"`
		PublishedAt string  `json:"published_at"`
		Assets      []Asset `json:"assets"`

		Prerelease bool `json:"prerelease"`
		Draft      bool `json:"draft"`
	}
	Asset struct {
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
)

func GetLatestGithubReleases(apiRepoURL string) ([]Release, error) {
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

func GetLocalVersion(logger *slog.Logger, binName, latestVer string) string {
	loc := WhereIsBinInstalled(binName)
	if loc == "" {
		// WhereIsBinInstalled loggers error already
		return ""
	}

	output, err := ExecCommandWithOutput(loc, "version")
	if err != nil {
		// typically old spawn / local-ic versions
		logger.Error("Error calling version command", "bin", binName+" version", "err", err)
		return "v0.0.0"
	}

	out := string(output)

	if i := strings.Index(out, "-"); i != -1 {
		out = semver.Canonical(out[:i])
	}

	if out == "" {
		logger.Debug("Could not parse version", "output", out, "setting to", "v0.0.0")
		out = "v0.0.0"
	}

	// strip out \n and \r
	out = strings.Trim(out, "\n\r")
	return out
}

func WhereIsBinInstalled(binName string) string {
	// looks for relative paths as well as the global $PATH for binaries
	for _, path := range []string{binName, path.Join("bin", binName), path.Join("local-interchain", binName)} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	if path, err := exec.LookPath(binName); err == nil {
		return path
	}

	return ""
}

// OutOfDateCheckLog logs & returns true if it is out of date.
func OutOfDateCheckLog(logger *slog.Logger, binName, current, latest string) bool {

	if current == "" {
		logger.Error("Could not find "+binName+" binary", "install", GetInstallMsg(howToInstallBinary[binName], latest))
		return true
	}

	isOutOfDate := semver.Compare(current, latest) < 0

	if isOutOfDate {
		logger.Error(
			"New "+binName+" available",
			"current", current,
			"latest", latest,
			"install", GetInstallMsg(howToInstallBinary[binName], latest),
		)
	}
	return isOutOfDate
}

func GetInstallMsg(msg, latestVer string) string {
	return strings.ReplaceAll(msg, "__VERSION__", latestVer)
}

// DoOutdatedNotificationRunCheck returns true if it is time to run the check again.
// It saves the last run time to a file for future runs.
func DoOutdatedNotificationRunCheck(logger *slog.Logger) bool {
	now := time.Now()

	lastCheckFile, err := GetLatestVersionCheckFile(logger)
	if err != nil {
		return false
	}

	lastCheck, err := os.ReadFile(lastCheckFile)
	if err != nil {
		logger.Error("Error reading last check file", "err", err)
		return false
	}

	lastCheckTime, err := time.Parse(time.RFC3339, string(lastCheck))
	if err != nil {
		// i.e. empty file
		lastCheckTime = time.Unix(0, 0)
	}

	diff := now.Sub(lastCheckTime)
	isTime := diff > RunCheckInterval

	if isTime {
		if err := WriteLastTimeToFile(logger, lastCheckFile, now); err != nil {
			logger.Error("Error writing last check file", "err", err)
			return false
		}
	}

	return isTime
}

func WriteLastTimeToFile(logger *slog.Logger, lastCheckFile string, t time.Time) error {
	return os.WriteFile(lastCheckFile, []byte(t.Format(time.RFC3339)), 0666)
}

// GetLatestVersionCheckFile grabs the check file used to determine when to run the version check.
func GetLatestVersionCheckFile(logger *slog.Logger) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Error getting user home dir", "err", err)
		return "", err
	}

	spawnDir := path.Join(home, ".spawn")
	if _, err := os.Stat(spawnDir); os.IsNotExist(err) {
		if err := os.Mkdir(spawnDir, 0755); err != nil {
			logger.Error("Error creating home spawn dir", "err", err)
			return "", err
		}
	}

	lastCheckFile := path.Join(spawnDir, "last_ver_check.txt")
	if _, err := os.Stat(lastCheckFile); os.IsNotExist(err) {
		if _, err := os.Create(lastCheckFile); err != nil {
			logger.Error("Error creating last check file", "err", err)
			return "", err
		}

		epoch := time.Unix(0, 0)
		if err := os.WriteFile(lastCheckFile, []byte(epoch.Format(time.RFC3339)), 0666); err != nil {
			logger.Error("Error writing last check file", "err", err)
			return "", err
		}
	}

	return lastCheckFile, nil
}
