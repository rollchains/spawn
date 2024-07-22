package main

import (
	"os/exec"
	"strings"
)

type CLIInfo struct {
	Description string
	Cmds        map[string]string
}

// Parses the CLI commands from a cobra binary to showcase values within the `plugins` subcommand.
func ParseCobraCLICmd(binAbsPath string) (CLIInfo, error) {
	output, err := exec.Command(binAbsPath).Output()
	if err != nil {
		return CLIInfo{}, err
	}

	sl := strings.Split(string(output), "\n")

	ci := CLIInfo{
		Cmds: make(map[string]string),
	}
	isAvailableCmds := false
	for idx, line := range sl {
		if idx == 0 {
			ci.Description = line
			if ci.Description == "" {
				ci.Description = "No description"
			}
		}

		// if line stars with `Available Commands:`, next lines are commands
		if strings.Contains(line, "Available Commands:") {
			isAvailableCmds = true
			continue
		}

		if isAvailableCmds {
			content := []string{}
			for _, item := range strings.Split(line, " ") {
				if item != "" {
					content = append(content, item)
				}
			}

			if len(content) >= 2 {
				ci.Cmds[content[0]] = strings.Join(content[1:], " ")
			}

			if line == "Flags:" {
				isAvailableCmds = false
				continue
			}
		}
	}
	return ci, nil
}
