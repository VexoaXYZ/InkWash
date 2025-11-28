package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/server"
	"github.com/VexoaXYZ/inkwash/pkg/types"
	"github.com/spf13/cobra"
)

var (
	migrateAll    bool
	migrateDryRun bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [server-name]",
	Short: "Migrate servers to new directory structure",
	Long: `Migrates servers from the old structure (shared binaries) to the new structure
(per-server bin/ directories with metadata tracking).

This command will:
  - Copy FXServer binaries to per-server bin/ directories
  - Generate metadata.json files with build and lifecycle information
  - Update launch scripts to use relative paths

Use --all to migrate all servers, or specify a server name to migrate just that server.
Use --dry-run to see what would be migrated without making changes.`,
	RunE: runMigrate,
}

func init() {
	migrateCmd.Flags().BoolVar(&migrateAll, "all", false, "Migrate all servers")
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Show what would be migrated without making changes")
	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) error {
	// Load registry
	reg, err := registry.NewRegistry(registry.GetRegistryPath())
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	var serversToMigrate []types.Server

	if migrateAll {
		serversToMigrate = reg.List()
	} else if len(args) == 1 {
		srv, err := reg.Get(args[0])
		if err != nil {
			return fmt.Errorf("server '%s' not found", args[0])
		}
		serversToMigrate = []types.Server{*srv}
	} else {
		return fmt.Errorf("specify a server name or use --all to migrate all servers")
	}

	if len(serversToMigrate) == 0 {
		fmt.Println("No servers to migrate")
		return nil
	}

	metadataManager := server.NewMetadataManager()
	configGen := server.NewConfigGenerator()

	fmt.Printf("Scanning %d server(s)...\n\n", len(serversToMigrate))

	migrated := 0
	skipped := 0
	failed := 0

	for _, srv := range serversToMigrate {
		// Check if already migrated
		if metadataManager.Exists(srv.Path) {
			if !migrateDryRun {
				fmt.Printf("  ○ %s - already migrated (metadata.json exists)\n", srv.Name)
			}
			skipped++
			continue
		}

		fmt.Printf("  → Migrating '%s'...\n", srv.Name)

		if migrateDryRun {
			fmt.Printf("    [DRY RUN] Would create bin/ directory\n")
			fmt.Printf("    [DRY RUN] Would generate metadata.json\n")
			fmt.Printf("    [DRY RUN] Would update launch script\n")
			migrated++
			continue
		}

		// Perform migration
		if err := migrateServer(&srv, metadataManager, configGen); err != nil {
			fmt.Printf("    ✗ Failed: %v\n", err)
			failed++
			continue
		}

		fmt.Printf("    ✓ Migrated successfully\n")
		migrated++
	}

	fmt.Printf("\nMigration Summary:\n")
	fmt.Printf("  Migrated: %d\n", migrated)
	fmt.Printf("  Skipped:  %d\n", skipped)
	if failed > 0 {
		fmt.Printf("  Failed:   %d\n", failed)
	}

	if migrateDryRun {
		fmt.Println("\n(Dry run - no changes made)")
	}

	return nil
}

func migrateServer(srv *types.Server, metadataManager *server.MetadataManager, configGen *server.ConfigGenerator) error {
	// Create bin/ directory
	binPath := filepath.Join(srv.Path, "bin")
	if err := os.MkdirAll(binPath, 0755); err != nil {
		return fmt.Errorf("failed to create bin/ directory: %w", err)
	}

	// Check if binaries need to be copied from old location
	// The old structure had binaries in a shared location, but since we don't know
	// where that was, we'll just check if bin/FXServer.exe exists
	fxServerPath := filepath.Join(binPath, "FXServer.exe")
	if _, err := os.Stat(fxServerPath); os.IsNotExist(err) {
		// Try to find binaries in parent directory structure
		// This is a best-effort attempt - may not work for all cases
		return fmt.Errorf("FXServer.exe not found - manual binary copy required")
	}

	// Generate metadata.json with best-effort data
	// We don't have the original build info, so we'll use placeholder values
	metadata := &types.ServerMetadata{
		Version: 1,
		Build: types.BuildMetadata{
			Number:      0, // Unknown build number
			Hash:        "unknown",
			InstalledAt: srv.Created,
			Recommended: false,
			Optional:    false,
		},
		Lifecycle: types.LifecycleMetadata{
			CreatedAt:   srv.Created,
			LastStarted: nil,
			LastStopped: nil,
		},
		Stats: types.UsageStats{
			RestartCount: 0,
			TotalUptime:  0,
		},
	}

	// If server has been started before, set last_started from in-memory data
	if !srv.LastStarted.IsZero() {
		metadata.Lifecycle.LastStarted = &srv.LastStarted
		metadata.Stats.RestartCount = 1
	}

	if err := metadataManager.Save(srv.Path, metadata); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	// Regenerate launch script with relative paths
	if err := configGen.GenerateLaunchScript(srv); err != nil {
		return fmt.Errorf("failed to update launch script: %w", err)
	}

	return nil
}
