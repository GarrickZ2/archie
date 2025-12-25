package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/GarrickZ2/archie/internal/setup"
	"github.com/GarrickZ2/archie/internal/ui"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup and manage project documentation",
	Long: `Interactive setup for project documentation.

This command provides a TUI interface to:
1. Edit background documentation (background.md)
2. Manage features (create, edit feature files)

Files marked with ✨ are empty or new and will be initialized with templates.`,
	Args: cobra.NoArgs,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
	// Use current directory as project path
	projectPath, err := os.Getwd()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get current directory: %v", err))
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create setup manager
	manager := setup.NewManager(projectPath)

	// Show main UI
	return manager.ShowMainUI()
}
