package cmd

import (
	"fmt"
	"os"

	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/server"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <server-name>",
	Short: "Start a FiveM server",
	Long:  `Start a FiveM server by name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		// Load registry
		reg, err := registry.NewRegistry(registry.GetRegistryPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load registry: %v\n", err)
			os.Exit(1)
		}

		// Get server
		srv, err := reg.Get(serverName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Server '%s' not found\n", serverName)
			os.Exit(1)
		}

		// Create process manager
		pm := server.NewProcessManager()

		// Check if already running
		if pm.IsRunning(srv) {
			fmt.Printf("Server '%s' is already running (PID: %d)\n", serverName, srv.PID)
			return
		}

		// Start server
		fmt.Printf("Starting server '%s'...\n", serverName)

		if err := pm.Start(srv); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to start server: %v\n", err)
			os.Exit(1)
		}

		// Update registry
		if err := reg.Update(*srv); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to update registry: %v\n", err)
		}

		fmt.Printf("✓ Server '%s' started successfully (PID: %d)\n", serverName, srv.PID)
		fmt.Printf("\nView logs:\n")
		fmt.Printf("  inkwash logs %s\n", serverName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
