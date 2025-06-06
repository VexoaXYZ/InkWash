package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vexoa/inkwash/internal/config"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print detailed version information about InkWash.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("InkWash %s\n", config.Version)
		fmt.Printf("  Commit: %s\n", config.Commit)
		fmt.Printf("  Built: %s\n", config.Date)
		fmt.Printf("  Go: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}