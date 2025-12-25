package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "archie",
	Short: "Technical design documentation project initialization tool",
	Long: `Archie is a CLI tool for initializing technical design documentation projects.
It creates a standardized file structure to help teams quickly start writing technical design documents.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Future: add global flags here
	// rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
