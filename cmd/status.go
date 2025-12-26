package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/GarrickZ2/archie/internal/status"
	"github.com/GarrickZ2/archie/internal/ui"
)

var (
	compactFlag bool
)

var statusCmd = &cobra.Command{
	Use:   "status [feature-key]",
	Short: "Show project status and progress",
	Long: `Display a comprehensive status report for all features in the project.

This command:
- Without arguments: Shows overall project status with all features
- With feature-key: Shows detailed information for a specific feature

Overall report:
- Parses all feature files in the features/ directory
- Extracts status information from each feature
- Shows overall progress, status distribution, and key insights
- Highlights blocked features that need attention
- Identifies stale features that haven't been updated recently

Detailed feature view:
- Shows complete feature information in a structured TUI format
- Displays status, summary, scope, requirements, dependencies, and more
- Provides a comprehensive view of a single feature

Statuses tracked:
  NOT_REVIEWED → UNDER_REVIEW → READY_FOR_DESIGN → UNDER_DESIGN →
  DESIGNED → SPEC_READY → IMPLEMENTING → FINISHED

Special status:
  BLOCKED - Features that are blocked and need attention`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&compactFlag, "compact", "c", false, "Show compact status report")
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Get current directory as project path
	projectPath, err := os.Getwd()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get current directory: %v", err))
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// If a feature key is provided, show detailed view for that feature
	if len(args) == 1 {
		return showFeatureDetail(projectPath, args[0])
	}

	// Otherwise, show overall status report
	return showOverallStatus(projectPath)
}

// showFeatureDetail 显示单个 feature 的详细信息
func showFeatureDetail(projectPath, featureKey string) error {
	detailParser := status.NewDetailParser(nil)
	detail, err := detailParser.ParseFeatureDetail(projectPath, featureKey)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to parse feature: %v", err))
		return fmt.Errorf("failed to parse feature: %w", err)
	}

	// Display detailed feature information
	display := status.NewDetailDisplay(detail)
	display.Show()

	return nil
}

// showOverallStatus 显示整体状态报告
func showOverallStatus(projectPath string) error {
	// Parse all features
	parser := status.NewParser(nil)
	features, err := parser.ParseFeaturesDir(projectPath)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to parse features: %v", err))
		return fmt.Errorf("failed to parse features: %w", err)
	}

	// Check if there are any features
	if len(features) == 0 {
		ui.ShowInfo("No features found in the features/ directory")
		fmt.Println()
		fmt.Println("💡 Tip: Use 'archie setup' to create and manage features")
		return nil
	}

	// Aggregate status information
	aggregator := status.NewAggregator(features)
	summary := aggregator.Aggregate()

	// Display status report
	display := status.NewDisplay(summary)

	if compactFlag {
		display.ShowCompact()
	} else {
		display.Show()
	}

	return nil
}
