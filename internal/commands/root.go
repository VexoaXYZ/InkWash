package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "inkwash",
	Short: "A CLI tool for creating FiveM servers instantly",
	Long: `InkWash is a powerful CLI tool that helps you create and manage FiveM servers
with a clean, optimized setup. It removes unnecessary files and provides a 
production-ready server configuration out of the box.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}