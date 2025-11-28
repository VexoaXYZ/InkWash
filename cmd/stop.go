package cmd

import (
	"fmt"
	"os"

	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/server"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <server-name>",
	Short: "Stop a FiveM server",
	Long:  `Stop a running FiveM server by name.`,
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

		// Check if running
		if !pm.IsRunning(srv) {
			fmt.Printf("Server '%s' is not running\n", serverName)
			return
		}

		// Stop server
		fmt.Printf("Stopping server '%s' (PID: %d)...\n", serverName, srv.PID)

		if err := pm.Stop(srv); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to stop server: %v\n", err)
			os.Exit(1)
		}

		// Update registry
		if err := reg.Update(*srv); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to update registry: %v\n", err)
		}

		fmt.Printf("✓ Server '%s' stopped successfully\n", serverName)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
