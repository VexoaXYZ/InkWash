package cmd

import (
	"fmt"
	"os"

	"github.com/VexoaXYZ/inkwash/internal/cache"
	"github.com/VexoaXYZ/inkwash/internal/registry"
	"github.com/VexoaXYZ/inkwash/internal/server"
	"github.com/VexoaXYZ/inkwash/internal/ui/wizard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create [server-name]",
	Short: "Create a new FiveM server",
	Long: `Create a new FiveM server with interactive configuration.

If server name is provided, uses defaults for other options.
Otherwise, launches interactive wizard.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Launch interactive wizard
			cachePath := registry.GetDefaultCachePath()
			binaryCache, err := cache.NewBinaryCache(cachePath, viper.GetInt("cache.max_builds"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to initialize cache: %v\n", err)
				os.Exit(1)
			}

			registryPath := registry.GetRegistryPath()
			reg, err := registry.NewRegistry(registryPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to initialize registry: %v\n", err)
				os.Exit(1)
			}

			vaultPath := registry.GetDefaultConfigPath() + "/keys.enc"
			vault, err := cache.NewKeyVault(vaultPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to load key vault: %v\n", err)
				os.Exit(1)
			}

			installer := server.NewInstaller(binaryCache, reg)
			wizardModel := wizard.NewCreateWizard(installer, vault, reg)

			p := tea.NewProgram(wizardModel, tea.WithAltScreen())
			finalModel, err := p.Run()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Check completion
			if wm, ok := finalModel.(*wizard.CreateWizardModel); ok {
				if wm.Completed() {
					fmt.Printf("\nServer '%s' is ready!\n", wm.ServerName())
				}
			}
			return
		}

		serverName := args[0]

		// Get flags
		buildNumber, _ := cmd.Flags().GetInt("build")
		keyID, _ := cmd.Flags().GetString("key")
		port, _ := cmd.Flags().GetInt("port")
		installPath, _ := cmd.Flags().GetString("path")

		if installPath == "" {
			installPath = viper.GetString("defaults.install_path")
		}

		if port == 0 {
			port = viper.GetInt("defaults.port")
		}

		// Initialize systems
		cachePath := registry.GetDefaultCachePath()
		binaryCache, err := cache.NewBinaryCache(cachePath, viper.GetInt("cache.max_builds"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to initialize cache: %v\n", err)
			os.Exit(1)
		}

		registryPath := registry.GetRegistryPath()
		reg, err := registry.NewRegistry(registryPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to initialize registry: %v\n", err)
			os.Exit(1)
		}

		// Get license key
		var licenseKey string
		if keyID != "" {
			vaultPath := registry.GetDefaultConfigPath() + "/keys.enc"
			vault, err := cache.NewKeyVault(vaultPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to load key vault: %v\n", err)
				os.Exit(1)
			}

			key, err := vault.Get(keyID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: License key not found: %v\n", err)
				os.Exit(1)
			}

			licenseKey = key.Key
		}

		// Create installer
		installer := server.NewInstaller(binaryCache, reg)

		// Install with progress
		fmt.Printf("Creating server '%s'...\n\n", serverName)

		err = installer.Install(serverName, installPath, buildNumber, licenseKey, port, func(progress server.InstallProgress) {
			fmt.Printf("[%d/%d] %s", progress.CompletedSteps, progress.TotalSteps, progress.Step)

			if progress.DownloadSpeed > 0 {
				fmt.Printf(" (%.1f MB/s, ETA: %s)", progress.DownloadSpeed, progress.DownloadETA.Round(1))
			}

			fmt.Println()
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✓ Server '%s' created successfully!\n", serverName)
		fmt.Printf("\nStart your server:\n")
		fmt.Printf("  inkwash start %s\n", serverName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().IntP("build", "b", 17000, "FXServer build number")
	createCmd.Flags().StringP("key", "k", "", "License key ID from vault")
	createCmd.Flags().IntP("port", "p", 0, "Server port (default: 30120)")
	createCmd.Flags().String("path", "", "Installation path")
}
