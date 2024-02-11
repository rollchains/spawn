package spawn

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (cfg *NewChainConfig) GitInitNewProjectRepo() {
	if err := ExecCommand("git", "init", cfg.ProjectName, "--quiet"); err != nil {
		fmt.Println("Error initializing git:", err)
	}
	if err := os.Chdir(cfg.ProjectName); err != nil {
		fmt.Println("Error changing to project directory:", err)
	}
	if err := ExecCommand("git", "add", "."); err != nil {
		fmt.Println("Error adding files to git:", err)
	}
	if err := ExecCommand("git", "commit", "-m", "initial commit", "--quiet"); err != nil {
		fmt.Println("Error committing initial files:", err)
	}
}
