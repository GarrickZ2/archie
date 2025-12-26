package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/GarrickZ2/archie/internal/export"
	"github.com/GarrickZ2/archie/internal/ui"
)

var (
	outputPath string
	noTOC      bool
	noStats    bool
	noDepGraph bool
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export project documentation to a single markdown file",
	Long: `Export project documentation with interactive selection.

This command will:
1. Let you select which documentation to export
2. Let you select which features to include
3. Generate table of contents, statistics, and dependency graph
4. Merge everything into a single markdown file

Options can be customized through flags or interactive prompts.`,
	Args: cobra.NoArgs,
	RunE: runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: ./archie-export-YYYY-MM-DD.md)")
	exportCmd.Flags().BoolVar(&noTOC, "no-toc", false, "Skip table of contents generation")
	exportCmd.Flags().BoolVar(&noStats, "no-stats", false, "Skip status statistics")
	exportCmd.Flags().BoolVar(&noDepGraph, "no-dep-graph", false, "Skip dependency graph")
}

func runExport(cmd *cobra.Command, args []string) error {
	// Get current directory as project path
	projectPath, err := os.Getwd()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get current directory: %v", err))
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create export manager
	manager := export.NewExportManager(projectPath, nil)

	// Set flags from command line
	manager.SetFlags(outputPath, !noTOC, !noStats, !noDepGraph)

	// Execute export
	result, err := manager.Export()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Export failed: %v", err))
		return err
	}

	// Show success summary
	showExportSuccess(result)

	return nil
}

func showExportSuccess(result *export.ExportResult) {
	content := []string{
		fmt.Sprintf("Documents: %d", result.DocumentCount),
		fmt.Sprintf("Features: %d", result.FeatureCount),
		fmt.Sprintf("Size: %s", formatBytes(result.TotalSize)),
		"",
		fmt.Sprintf("Output: %s", result.OutputPath),
	}

	// Show warnings if any
	if len(result.Warnings) > 0 {
		content = append(content, "")
		content = append(content, fmt.Sprintf("Warnings: %d", len(result.Warnings)))
	}

	ui.PrintBox("Export Complete!", content)

	// Show warnings details
	if len(result.Warnings) > 0 {
		fmt.Println()
		ui.ShowInfo("Warnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
		fmt.Println()
	}
}

// formatBytes formats bytes to human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
