package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vexoa/inkwash/internal/fivem"
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
		creator := fivem.NewServerCreator()
		
		config := fivem.ServerConfig{
			Name:     serverName,
			Path:     serverPath,
			Template: template,
		}

		if err := creator.Create(config); err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		fmt.Printf("✅ Successfully created FiveM server '%s' at %s\n", serverName, serverPath)
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