package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

const LocalICURL = "https://github.com/strangelove-ventures/interchaintest/releases/download/v8.0.0/local-ic"

var DownloadLocalIC = &cobra.Command{
	Use:   "download",
	Short: fmt.Sprintf("Download LocalInterchain from %s", LocalICURL),
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		downloadBin()
	},
}

// TODO: Download & move to the users GOPATH automatically? (keep a copy in a build/ folder?)
func downloadBin() error {
	file := "local-ic"

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Printf("Downloading Local Interchain binary...")

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	downloadWithProgress(path.Join(currentDir, file), LocalICURL)

	// Make the binary executable
	if err := os.Chmod(file, 0755); err != nil {
		return err
	}

	fmt.Printf("Downloaded Local Interchain binary to %s\n", path.Join(currentDir, file))

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

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading local-ic",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	return nil
}
