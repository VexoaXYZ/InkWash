package cmd

import (
	"fmt"
	"os"

	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/server"
	"github.com/VexoaXYZ/inkwash/internal/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all FiveM servers",
	Long:  `List all registered FiveM servers with their status.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load registry
		reg, err := registry.NewRegistry(registry.GetRegistryPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load registry: %v\n", err)
			os.Exit(1)
		}

		servers := reg.List()

		if len(servers) == 0 {
			fmt.Println("No servers found")
			fmt.Println("\nCreate a server:")
			fmt.Println("  inkwash create <server-name>")
			return
		}

		// Create process manager to check status
		pm := server.NewProcessManager()

		fmt.Printf("\n%s\n\n", ui.RenderHeader("SERVERS"))

		for _, srv := range servers {
			// Check actual running status
			isRunning := pm.IsRunning(&srv)

			// Status indicator
			var status string
			if isRunning {
				status = ui.RenderStatusRunning(srv.Status())
			} else {
				status = ui.RenderStatusStopped(srv.Status())
			}

			fmt.Printf("  %s  %s\n", status, ui.RenderAccent(srv.Name))
			fmt.Printf("      %s\n", ui.RenderMuted("Port: "+fmt.Sprint(srv.Port)))
			fmt.Printf("      %s\n", ui.RenderPath(srv.Path))

			if isRunning {
				// Get memory usage
				mem, err := pm.GetMemoryUsage(&srv)
				if err == nil {
					memGB := float64(mem) / 1024 / 1024 / 1024
					fmt.Printf("      %s\n", ui.RenderMuted(fmt.Sprintf("RAM: %.2f GB", memGB)))
				}
			}

			fmt.Println()
		}

		fmt.Printf("Total: %d server(s)\n\n", len(servers))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
