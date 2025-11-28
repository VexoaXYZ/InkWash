package cmd

import (
	"fmt"
	"os"

	"github.com/VexoaXYZ/inkwash/internal/cache"
	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/ui"
	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage license keys",
	Long:  `Manage FiveM license keys in encrypted vault.`,
}

var keyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new license key",
	Run: func(cmd *cobra.Command, args []string) {
		label, _ := cmd.Flags().GetString("label")
		key, _ := cmd.Flags().GetString("key")

		// Prompt for inputs if not provided
		if label == "" {
			fmt.Print("Enter a label for this key: ")
			fmt.Scanln(&label)
		}

		if key == "" {
			fmt.Print("Enter license key (cfxk_...): ")
			fmt.Scanln(&key)
		}

		// Load vault
		vaultPath := registry.GetDefaultConfigPath() + "/keys.enc"
		vault, err := cache.NewKeyVault(vaultPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load vault: %v\n", err)
			os.Exit(1)
		}

		// Add key
		id, err := vault.Add(label, key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to add key: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", ui.RenderSuccess("License key added"))
		fmt.Printf("ID: %s\n", id)
		fmt.Printf("Label: %s\n", label)
	},
}

var keyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all license keys",
	Run: func(cmd *cobra.Command, args []string) {
		// Load vault
		vaultPath := registry.GetDefaultConfigPath() + "/keys.enc"
		vault, err := cache.NewKeyVault(vaultPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load vault: %v\n", err)
			os.Exit(1)
		}

		keys := vault.List()

		if len(keys) == 0 {
			fmt.Println("No license keys found")
			fmt.Println("\nAdd a key:")
			fmt.Println("  inkwash key add")
			return
		}

		fmt.Printf("\n%s\n\n", ui.RenderHeader("LICENSE KEYS"))

		for _, key := range keys {
			fmt.Printf("  %s\n", ui.RenderAccent(key.Label))
			fmt.Printf("    ID:  %s\n", ui.RenderMuted(key.ID))
			fmt.Printf("    Key: %s\n", ui.RenderMuted(cache.MaskKey(key.Key)))
			fmt.Printf("    Created: %s\n\n", ui.RenderMuted(key.Created.Format("Jan 2, 2006")))
		}

		fmt.Printf("Total: %d key(s)\n\n", len(keys))
	},
}

var keyRemoveCmd = &cobra.Command{
	Use:   "remove <key-id>",
	Short: "Remove a license key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyID := args[0]

		// Load vault
		vaultPath := registry.GetDefaultConfigPath() + "/keys.enc"
		vault, err := cache.NewKeyVault(vaultPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to load vault: %v\n", err)
			os.Exit(1)
		}

		// Remove key
		if err := vault.Remove(keyID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to remove key: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", ui.RenderSuccess("License key removed"))
	},
}

func init() {
	rootCmd.AddCommand(keyCmd)

	keyCmd.AddCommand(keyAddCmd)
	keyCmd.AddCommand(keyListCmd)
	keyCmd.AddCommand(keyRemoveCmd)

	keyAddCmd.Flags().StringP("label", "l", "", "Label for the key")
	keyAddCmd.Flags().StringP("key", "k", "", "License key")
}
