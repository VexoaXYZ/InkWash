package cmd

import (
	"fmt"
	"os"

	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/ui/wizard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert GTA5 mods to FiveM resources",
	Long:  `Convert GTA5 mods from gta5-mods.com to FiveM resources using the convert.cfx.rs service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load registry
		reg, err := registry.NewRegistry(registry.GetRegistryPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load registry: %v\n", err)
			os.Exit(1)
		}

		// Check if any servers exist
		servers := reg.List()
		if len(servers) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No servers found. Please create a server first.\n")
			fmt.Println("\nCreate a server:")
			fmt.Println("  inkwash create <server-name>")
			os.Exit(1)
		}

		// Create and run wizard
		wizardModel := wizard.NewConvertWizard(reg)
		p := tea.NewProgram(wizardModel, tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Check if wizard completed successfully
		if m, ok := finalModel.(*wizard.ConvertWizardModel); ok {
			if !m.Completed() {
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
