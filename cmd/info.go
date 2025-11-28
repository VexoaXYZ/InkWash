package cmd

import (
	"fmt"
	"time"

	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/server"
	"github.com/VexoaXYZ/inkwash/pkg/types"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <server-name>",
	Short: "Display detailed information about a server",
	Long:  `Shows build information, lifecycle events, and usage statistics for a server.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	serverName := args[0]

	// Load registry
	reg, err := registry.NewRegistry(registry.GetRegistryPath())
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Get server
	srv, err := reg.Get(serverName)
	if err != nil {
		return fmt.Errorf("server '%s' not found", serverName)
	}

	// Load metadata
	metadataManager := server.NewMetadataManager()
	metadata, err := metadataManager.Load(srv.Path)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	// Display server info
	fmt.Printf("\n%s\n", bold("SERVER INFORMATION"))
	fmt.Printf("  Name:     %s\n", srv.Name)
	fmt.Printf("  Path:     %s\n", srv.Path)
	fmt.Printf("  Port:     %d\n", srv.Port)
	fmt.Printf("  Status:   %s\n", getStatusString(srv))

	// Display build info
	fmt.Printf("\n%s\n", bold("BUILD"))
	fmt.Printf("  Number:      %d\n", metadata.Build.Number)
	fmt.Printf("  Hash:        %s\n", metadata.Build.Hash)
	fmt.Printf("  Installed:   %s\n", formatTime(metadata.Build.InstalledAt))
	fmt.Printf("  Type:        %s\n", getBuildType(metadata.Build.Recommended, metadata.Build.Optional))

	// Display lifecycle info
	fmt.Printf("\n%s\n", bold("LIFECYCLE"))
	fmt.Printf("  Created:      %s (%s)\n",
		formatTime(metadata.Lifecycle.CreatedAt),
		formatRelativeTime(metadata.Lifecycle.CreatedAt))

	if metadata.Lifecycle.LastStarted != nil {
		fmt.Printf("  Last Started: %s (%s)\n",
			formatTime(*metadata.Lifecycle.LastStarted),
			formatRelativeTime(*metadata.Lifecycle.LastStarted))
	} else {
		fmt.Printf("  Last Started: Never\n")
	}

	if metadata.Lifecycle.LastStopped != nil {
		fmt.Printf("  Last Stopped: %s (%s)\n",
			formatTime(*metadata.Lifecycle.LastStopped),
			formatRelativeTime(*metadata.Lifecycle.LastStopped))
	} else {
		fmt.Printf("  Last Stopped: Never\n")
	}

	// Display usage stats
	fmt.Printf("\n%s\n", bold("USAGE STATISTICS"))
	fmt.Printf("  Restart Count: %d\n", metadata.Stats.RestartCount)
	fmt.Printf("  Total Uptime:  %s\n", formatDuration(metadata.Stats.TotalUptime))

	fmt.Println()
	return nil
}

func getStatusString(srv *types.Server) string {
	if srv.IsRunning() {
		return fmt.Sprintf("Running (PID: %d)", srv.PID)
	}
	return "Stopped"
}

func getBuildType(recommended, optional bool) string {
	if recommended {
		return "Recommended"
	}
	if optional {
		return "Optional"
	}
	return "Standard"
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}

func formatRelativeTime(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	if days < 30 {
		return fmt.Sprintf("%d days ago", days)
	}

	months := days / 30
	if months == 1 {
		return "1 month ago"
	}
	if months < 12 {
		return fmt.Sprintf("%d months ago", months)
	}

	years := months / 12
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string
	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}
	if seconds > 0 || len(parts) == 0 {
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) == 2 {
		return parts[0] + " " + parts[1]
	}
	// For 3+ parts, show first two
	return parts[0] + " " + parts[1]
}

func bold(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}
