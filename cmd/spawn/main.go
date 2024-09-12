package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/rollchains/spawn/spawn"
	"github.com/spf13/cobra"
)

var (
	// Set in the makefile ld_flags on compile
	SpawnVersion = ""
	LogLevelFlag = "log-level"
	rootCmd      = &cobra.Command{
		Use:   "spawn",
		Short: "Entry into the Interchain | Contact us: support@rollchains.com",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: false,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func NewRootCmd() *cobra.Command {
	return rootCmd
}

func main() {
	outOfDateChecker()

	rootCmd.AddCommand(newChain)
	rootCmd.AddCommand(LocalICCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of spawn",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(SpawnVersion)
		},
	})
	rootCmd.AddCommand(ModuleCmd())
	rootCmd.AddCommand(ProtoServiceGenerate())
	rootCmd.AddCommand(DocsCmd)

	applyPluginCmds()

	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error while executing your CLI. Err: %v\n", err)
		os.Exit(1)
	}
}

func GetLogger() *slog.Logger {
	w := os.Stderr

	logLevel := parseLogLevelFromFlags()

	slog.SetDefault(slog.New(
		// TODO: Windows support colored logs: https://github.com/mattn/go-colorable `tint.NewHandler(colorable.NewColorable(w), nil)`
		tint.NewHandler(w, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.Kitchen,
			// Enables colors only if the terminal supports it
			NoColor: !isatty.IsTerminal(w.Fd()),
		}),
	))

	return slog.Default()
}

// outOfDateChecker checks if binaries are up to date and logs if they are not.
// if not, it will prompt the user every command they run with spawn until they update.
// else, it will wait 24h+ before checking again.
func outOfDateChecker() {
	logger := GetLogger()

	if !spawn.DoOutdatedNotificationRunCheck(logger) {
		return
	}

	for _, program := range []string{"local-ic", "spawn"} {
		releases, err := spawn.GetLatestGithubReleases(spawn.BinaryToGithubAPI[program])
		if err != nil {
			logger.Error("Error getting latest local-ic releases", "err", err)
			return
		}
		latest := releases[0].TagName

		current := spawn.GetLocalVersion(logger, program, latest)
		if spawn.OutOfDateCheckLog(logger, program, current, latest) {
			// write check to -24h from now to spam the user until it's resolved.

			file, err := spawn.GetLatestVersionCheckFile(logger)
			if err != nil {
				return
			}

			if err := spawn.WriteLastTimeToFile(logger, file, time.Now().Add(-spawn.RunCheckInterval)); err != nil {
				logger.Error("Error writing last check file", "err", err)
				return
			}

			return
		}
	}
}

func parseLogLevelFromFlags() slog.Level {
	logLevel, err := rootCmd.PersistentFlags().GetString(LogLevelFlag)
	if err != nil {
		return slog.LevelInfo
	}

	var lvl slog.Level

	switch strings.ToLower(logLevel) {
	case "debug", "d", "dbg":
		lvl = slog.LevelDebug
	case "info", "i", "inf":
		lvl = slog.LevelInfo
	case "warn", "w", "wrn":
		lvl = slog.LevelWarn
	case "error", "e", "err", "fatal", "f", "ftl":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	return lvl
}
