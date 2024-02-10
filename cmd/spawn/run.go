package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// generate base JSOn if not already off of the app (local-ic format)
// heighliner build latest.
// check if local-ic is found in path, if not, tell the user to download & move to their GOPATH (or do automatically)
// if local-ic is installed, call it here automatically
var BuildAppImage = &cobra.Command{
	Use:   "docker-build",
	Short: "Build Docker Image for your app",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building Docker Image for your app...")
		installHeighliner()
		fmt.Println("Building Local Docker Image...")
		buildLocalDockerImage()
	},
}

func installHeighliner() {
	binary := "heighliner"

	if err := exec.Command(binary).Run(); err != nil {
		fmt.Println("heighliner not found, installing...")

		if err := exec.Command("make", "get-heighliner").Run(); err != nil {
			panic(fmt.Sprintf("Error installing heighliner: %s", err))
		}

		fmt.Println("heighliner installed!")

		if err := os.RemoveAll("heighliner"); err != nil {
			panic(err)
		}
	}
}

func buildLocalDockerImage() {
	stdout, err := exec.Command("make", "local-image").Output()
	if err != nil {
		panic(fmt.Sprintf("Error building local image: %s. Make sure you are in your project repo when running this command.", err))
	}

	fmt.Println(string(stdout))
}
