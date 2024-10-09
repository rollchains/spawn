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

func ExecCommandWithOutput(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
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
	if err := os.Chdir(".."); err != nil {
		cfg.Logger.Error("Error changing back to original directory", "err", err)
	}
}

func (cfg *NewChainConfig) MakeModTidy() {
	cfg.Logger.Info("Running `go mod tidy`, this may take a minute on the first run...")
	if err := os.Chdir(cfg.ProjectName); err != nil {
		cfg.Logger.Error("Error changing to project directory", "err", err)
	}
	if err := ExecCommand("make", "mod-tidy"); err != nil {
		cfg.Logger.Error("Error running `make mod-tidy`", "err", err)
	}
	if err := os.Chdir(".."); err != nil {
		cfg.Logger.Error("Error changing back to original directory", "err", err)
	}
}
