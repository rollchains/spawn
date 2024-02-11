package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"gitub.com/strangelove-ventures/spawn/spawn"
)

var (
	LocalICDefaultVersion = "v8.1.0"
	LocalICURL            = "https://github.com/strangelove-ventures/interchaintest/releases/download/" + LocalICDefaultVersion + "/local-ic"
)

const (
	FlagVersionOverride = "version"
	FlagForceDownload   = "download"
	FlagLocationPath    = "print-location"
)

func init() {
	LocalICCmd.Flags().String(FlagVersionOverride, LocalICDefaultVersion, "change the local-ic version to use")
	LocalICCmd.Flags().Bool(FlagForceDownload, false, "force download of local-ic")
	LocalICCmd.Flags().Bool(FlagLocationPath, false, "print the location of local-ic binary")
}

// ---
// make install && ICTEST_HOME=./simapp spawn local-ic start testnet
// make install && cd simapp && spawn local-ic start testnet
// ---
// TODO: Do something like `curl https://get.ignite.com/cli! | bash`? just with windows support for path
var LocalICCmd = &cobra.Command{
	Use:   "local-ic",
	Short: "Local Interchain",
	Long:  fmt.Sprintf("Download Local Interchain from %s", LocalICURL),
	// Args:  cobra.
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetString(FlagVersionOverride)
		forceDownload, _ := cmd.Flags().GetBool(FlagForceDownload)
		debugBinaryLoc, _ := cmd.Flags().GetBool(FlagLocationPath)

		loc := whereIsLocalICInstalled()
		if (forceDownload || loc == "") && version != "" {
			downloadBin(version)
			loc = whereIsLocalICInstalled()
		}

		if debugBinaryLoc {
			fmt.Println(loc)
			return
		}

		if err := os.Chmod(loc, 0755); err != nil {
			fmt.Println("Error setting local-ic permissions:", err)
		}

		// set to use the current dir if it is not overrriden
		if os.Getenv("ICTEST_HOME") == "" {
			if err := os.Setenv("ICTEST_HOME", "."); err != nil {
				fmt.Println("Error setting ICTEST_HOME:", err)
			}
		}

		if err := spawn.ExecCommand(loc, args...); err != nil {
			fmt.Println("Error calling local-ic:", err)
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

func downloadBin(version string) error {
	file := "local-ic"

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if version != "" && version != LocalICDefaultVersion {
		if version[0] != 'v' {
			version = "v" + version
		}

		LocalICURL = strings.ReplaceAll(LocalICURL, LocalICDefaultVersion, version)
	}

	dir := path.Join(currentDir, "bin")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	filePath := path.Join(dir, file)

	if err := downloadWithProgress(filePath, LocalICURL); err != nil {
		return err
	}

	if err := os.Chmod(file, 0755); err != nil {
		return err
	}

	fmt.Printf("✅ Local Interchain Downloaded to %s\n", filePath)
	return nil
}

func downloadWithProgress(destinationPath, downloadUrl string) error {
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	f, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription("⏳ Downloading Local-Interchain..."),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(0*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	io.Copy(io.MultiWriter(f, bar), resp.Body)
	return nil
}
