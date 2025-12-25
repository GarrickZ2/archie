package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/GarrickZ2/archie/internal/agent"
	"github.com/GarrickZ2/archie/internal/clone"
	"github.com/GarrickZ2/archie/internal/ui"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <source-path>",
	Short: "Clone an existing Archie project",
	Long: `Clone an existing Archie project to the current directory.

This command will:
1. Validate source project (must have .archie folder)
2. Let you select a clone strategy (context, light, full, or custom)
3. Copy selected files and directories from source
4. Initialize missing project structure
5. Replicate agent configuration from source

Clone Strategies:
  context  - Copy project context and structure (recommended)
             Includes: background.md, dependency.md, architecture.md,
                      storage.md, faq.md, api/

  light    - Copy only background documentation
             Includes: background.md

  full     - Copy everything from source project
             Includes: All standard files and directories

  custom   - Manually select which files and directories to copy
             You can choose from all standard project items

Examples:
  # Clone from a local project
  archie clone /path/to/source-project

  # Clone from a relative path
  archie clone ../another-project`,
	Args: cobra.ExactArgs(1),
	RunE: runClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	// Load custom agents (non-fatal if fails)
	if err := agent.LoadAndRegister(); err != nil {
		// Just log info, continue with built-in agents
		ui.ShowInfo("⚠️  Could not load custom agents: " + err.Error())
		fmt.Println()
	}

	// Get source path from args
	sourcePath := args[0]

	// Convert to absolute path
	absSourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Invalid source path: %v", err))
		return fmt.Errorf("invalid source path: %w", err)
	}

	// Get target path (current directory)
	targetPath, err := os.Getwd()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get current directory: %v", err))
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create clone manager
	manager := clone.NewCloneManager(nil)

	// Execute clone
	config := clone.CloneConfig{
		SourcePath: absSourcePath,
		TargetPath: targetPath,
	}

	_, err = manager.Clone(config)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Clone failed: %v", err))
		return err
	}

	return nil
}
