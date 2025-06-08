package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vexoa/inkwash/internal/config"
	"github.com/vexoa/inkwash/internal/update"
)

var (
	skipUpdateCheck bool
	rootCmd = &cobra.Command{
		Use:   "inkwash",
		Short: "A CLI tool for creating FiveM servers instantly",
		Long: `InkWash is a powerful CLI tool that helps you create and manage FiveM servers
with a clean, optimized setup. It removes unnecessary files and provides a 
production-ready server configuration out of the box.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Skip update check for certain commands or if flag is set
			if skipUpdateCheck || cmd.Name() == "update" || cmd.Name() == "version" {
				return
			}

			// Check if we should check for updates (once per day)
			if update.ShouldCheckForUpdate() {
				checkForUpdatesInBackground()
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(updateCmd)
	
	// Add global flags
	rootCmd.PersistentFlags().BoolVar(&skipUpdateCheck, "skip-update-check", false, "Skip automatic update check")
}

func checkForUpdatesInBackground() {
	// Run update check in background to not block the main command
	go func() {
		updater := update.NewUpdater(config.Version)
		info, err := updater.CheckForUpdate()
		if err != nil {
			return // Silently fail
		}

		// Save the check time
		update.SaveUpdateCheckTime()

		if info.Available {
			// Print update notification to stderr so it doesn't interfere with command output
			fmt.Fprintf(os.Stderr, "\n[UPDATE] New version available: v%s (current: v%s)\n", info.LatestVersion, info.CurrentVersion)
			fmt.Fprintf(os.Stderr, "Run 'inkwash update' to upgrade\n\n")
		}
	}()
}