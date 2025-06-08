package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vexoa/inkwash/internal/services"
)

var (
	serverName string
	serverPath string
	template   string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new FiveM server",
	Long:  `Create a new FiveM server with optimized configuration and cleaned setup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create service container
		container := services.NewContainer()
		
		// Create context
		ctx := context.Background()
		
		// Create server using the new service
		server, err := container.ServerService.CreateServer(ctx, serverName, serverPath, template)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		fmt.Printf("✅ Server created successfully with ID: %s\n", server.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&serverName, "name", "n", "", "Server name (required)")
	createCmd.Flags().StringVarP(&serverPath, "path", "p", ".", "Path where server will be created")
	createCmd.Flags().StringVarP(&template, "template", "t", "default", "Server template to use")

	createCmd.MarkFlagRequired("name")
}