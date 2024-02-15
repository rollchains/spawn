package spawn

import (
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
		cfg.Logger.Error("Error initializing git", "err", err)
	}
	if err := os.Chdir(cfg.ProjectName); err != nil {
		cfg.Logger.Error("Error changing to project directory", "err", err)
	}
	if err := ExecCommand("git", "add", "."); err != nil {
		cfg.Logger.Error("Error adding files to git", "err", err)
	}
	if err := ExecCommand("git", "commit", "-m", "initial commit", "--quiet"); err != nil {
		cfg.Logger.Error("Error committing initial files", "err", err)
	}
}
