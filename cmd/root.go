package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inkwash",
	Short: "Premium FiveM server manager",
	Long: `Inkwash - A world-class CLI tool for managing FiveM servers.

Inkwash brings the polish of modern web applications to terminal interfaces.
Create, manage, and monitor FiveM servers with beautiful animations and
real-time metrics.`,
	// If no subcommand is provided, launch the interactive dashboard
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Launch interactive dashboard
		fmt.Println("Inkwash dashboard coming soon...")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/inkwash/config.yaml)")
	rootCmd.PersistentFlags().Bool("no-animations", false, "disable all animations")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".inkwash" (without extension).
		viper.AddConfigPath(home + "/.config/inkwash")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Set defaults
	viper.SetDefault("defaults.install_path", getDefaultInstallPath())
	viper.SetDefault("defaults.port", 30120)
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.max_builds", 3)
	viper.SetDefault("ui.theme", "purple")
	viper.SetDefault("ui.animations", "auto")
	viper.SetDefault("ui.refresh_interval", 2)
	viper.SetDefault("telemetry.enabled", true)
	viper.SetDefault("advanced.parallel_downloads", true)
	viper.SetDefault("advanced.download_chunks", 3)
	viper.SetDefault("advanced.log_level", "info")
}

func getDefaultInstallPath() string {
	if isWindows() {
		return "C:\\FXServer"
	}
	home, _ := os.UserHomeDir()
	return home + "/fxserver"
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}
