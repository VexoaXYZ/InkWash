package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vexoa/inkwash/internal/config"
	"github.com/vexoa/inkwash/internal/update"
)

var (
	checkOnly bool
	forceUpdate bool
	rollback bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long: `Check for new versions of InkWash and optionally install them.

By default, this command will check for updates and prompt before installing.
Use --check to only check for updates without installing.
Use --force to skip the confirmation prompt and install immediately.`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates without installing")
	updateCmd.Flags().BoolVar(&forceUpdate, "force", false, "Force update without confirmation")
	updateCmd.Flags().BoolVar(&rollback, "rollback", false, "Rollback to the previous version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Handle rollback
	if rollback {
		return update.Rollback()
	}

	updater := update.NewUpdater(config.Version)

	// Check for updates
	fmt.Println("Checking for updates...")
	info, err := updater.CheckForUpdate()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	// Save the check time
	update.SaveUpdateCheckTime()

	if !info.Available {
		fmt.Printf("You are running the latest version (v%s)\n", info.CurrentVersion)
		return nil
	}

	// Display update information
	fmt.Printf("\nUpdate available!\n")
	fmt.Printf("Current version: v%s\n", info.CurrentVersion)
	fmt.Printf("Latest version:  v%s\n", info.LatestVersion)
	
	if info.ReleaseNotes != "" {
		fmt.Printf("\nRelease notes:\n%s\n", info.ReleaseNotes)
	}

	// If check only, stop here
	if checkOnly {
		return nil
	}

	// Confirm update
	if !forceUpdate {
		fmt.Print("\nDo you want to install this update? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Update cancelled")
			return nil
		}
	}

	// Perform update
	fmt.Println("\nInstalling update...")
	if err := updater.Update(info); err != nil {
		fmt.Fprintf(os.Stderr, "\nUpdate failed: %v\n", err)
		fmt.Println("\nTo update manually:")
		fmt.Printf("1. Download the latest version from: https://github.com/%s/releases/latest\n", update.GithubRepo)
		fmt.Printf("2. Replace your current binary with the downloaded file\n")
		fmt.Printf("3. Make sure the file is executable (chmod +x on Unix systems)\n")
		return err
	}

	fmt.Printf("\nUpdate successful! Please restart InkWash to use version v%s\n", info.LatestVersion)
	return nil
}