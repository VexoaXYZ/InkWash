package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <server-name>",
	Short: "View server logs",
	Long:  `View logs for a FiveM server.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		follow, _ := cmd.Flags().GetBool("follow")
		lines, _ := cmd.Flags().GetInt("lines")

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

		// Get log file path
		logPath := filepath.Join(srv.Path, "logs", "server.log")

		// Check if log exists
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Log file not found: %s\n", logPath)
			os.Exit(1)
		}

		// Open log file
		file, err := os.Open(logPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to open log: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		if follow {
			// TODO: Implement tail -f functionality
			fmt.Println("Follow mode not implemented yet")
			fmt.Println("Showing last lines instead...")
		}

		// Show last N lines
		scanner := bufio.NewScanner(file)
		var allLines []string

		for scanner.Scan() {
			allLines = append(allLines, scanner.Text())
		}

		// Print last N lines
		start := len(allLines) - lines
		if start < 0 {
			start = 0
		}

		for i := start; i < len(allLines); i++ {
			fmt.Println(allLines[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	logsCmd.Flags().IntP("lines", "n", 50, "Number of lines to show")
}
