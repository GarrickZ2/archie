package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/GarrickZ2/archie/internal/status"
	"github.com/GarrickZ2/archie/internal/ui"
)

var (
	compactFlag  bool
	overviewFlag bool
	featureFlag  string
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project status and progress (Interactive TUI)",
	Long: `Display a comprehensive status report for all features in the project.

This command launches an interactive TUI that allows you to:
- View overall project status with all features
- Browse and select individual features for detailed information

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
	Args: cobra.NoArgs,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&compactFlag, "compact", "c", false, "Show compact status report")
	statusCmd.Flags().BoolVarP(&overviewFlag, "overview", "o", false, "Show overview directly")
	statusCmd.Flags().StringVarP(&featureFlag, "feature", "f", "", "Show feature list or specific feature detail (feature-key or file path)")
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Get current directory as project path
	projectPath, err := os.Getwd()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get current directory: %v", err))
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Handle direct mode flags
	if overviewFlag {
		return showOverallStatus(projectPath)
	}

	if featureFlag != "" {
		// Show specific feature detail
		return showFeatureByPath(projectPath, featureFlag)
	}

	// -f without value shows feature list
	if cmd.Flags().Changed("feature") {
		return showFeatureListMenu(projectPath)
	}

	// Show interactive TUI menu
	return showInteractiveMenu(projectPath)
}

// showInteractiveMenu displays the main TUI menu
func showInteractiveMenu(projectPath string) error {
	// Main menu options
	options := []string{
		"Overview - Show overall project status",
		"Feature List - Browse and select individual features",
	}

	var choice string
	prompt := &survey.Select{
		Message: "What would you like to view?",
		Options: options,
	}

	if err := survey.AskOne(prompt, &choice); err != nil {
		return fmt.Errorf("menu selection cancelled")
	}

	// Handle menu choice
	switch choice {
	case options[0]: // Overview
		return showOverallStatus(projectPath)
	case options[1]: // Feature List
		return showFeatureListMenu(projectPath)
	}

	return nil
}

// showFeatureListMenu displays a list of features and allows selection
func showFeatureListMenu(projectPath string) error {
	// Parse all features
	parser := status.NewParser(nil)
	features, err := parser.ParseFeaturesDir(projectPath)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to parse features: %v", err))
		return fmt.Errorf("failed to parse features: %w", err)
	}

	if len(features) == 0 {
		ui.ShowInfo("No features found in the features/ directory")
		fmt.Println()
		fmt.Println("Tip: Use 'archie setup' to create and manage features")
		return nil
	}

	// Build feature options grouped by status
	optionsMap := make(map[string]string) // option text -> feature name
	var options []string

	// Group features by status
	statusGroups := make(map[status.FeatureStatus][]status.Feature)
	for _, feature := range features {
		statusGroups[feature.Status] = append(statusGroups[feature.Status], feature)
	}

	// Display order
	displayOrder := []status.FeatureStatus{
		status.StatusImplementing,
		status.StatusSpecReady,
		status.StatusDesigned,
		status.StatusUnderDesign,
		status.StatusReadyForDesign,
		status.StatusUnderReview,
		status.StatusNotReviewed,
		status.StatusBlocked,
		status.StatusFinished,
	}

	// Build options in order
	for _, s := range displayOrder {
		if featureList, ok := statusGroups[s]; ok && len(featureList) > 0 {
			statusName := strings.ReplaceAll(string(s), "_", " ")
			for _, feature := range featureList {
				option := fmt.Sprintf("%-35s [%s]", feature.Name, statusName)
				options = append(options, option)
				optionsMap[option] = feature.Name
			}
		}
	}

	// Show feature selection
	var choice string
	prompt := &survey.Select{
		Message:  "Select a feature to view details:",
		Options:  options,
		PageSize: 15,
	}

	if err := survey.AskOne(prompt, &choice); err != nil {
		return fmt.Errorf("feature selection cancelled")
	}

	// Get selected feature name
	featureName := optionsMap[choice]

	// Show feature detail
	return showFeatureDetail(projectPath, featureName)
}

// showFeatureByPath shows feature detail by feature-key or file path
func showFeatureByPath(projectPath, input string) error {
	// Extract feature key from various input formats
	featureKey := extractFeatureKey(input)
	return showFeatureDetail(projectPath, featureKey)
}

// extractFeatureKey extracts feature-key from file path or returns input as-is
func extractFeatureKey(input string) string {
	// Remove leading "./"
	input = strings.TrimPrefix(input, "./")

	// Remove "features/" prefix if present
	input = strings.TrimPrefix(input, "features/")

	// Remove ".md" suffix if present
	input = strings.TrimSuffix(input, ".md")

	return input
}

// showFeatureDetail 显示单个 feature 的详细信息
func showFeatureDetail(projectPath, featureKey string) error {
	detailParser := status.NewDetailParser(nil)
	detail, err := detailParser.ParseFeatureDetail(projectPath, featureKey)
	if err != nil {
		// Feature not found, provide helpful suggestions
		return handleFeatureNotFound(projectPath, featureKey)
	}

	// Display detailed feature information
	display := status.NewDetailDisplay(detail)
	display.Show()

	return nil
}

// handleFeatureNotFound 处理 feature 未找到的情况
func handleFeatureNotFound(projectPath, featureKey string) error {
	ui.ShowError(fmt.Sprintf("Feature '%s' not found", featureKey))
	fmt.Println()

	// Get all available features
	parser := status.NewParser(nil)
	features, err := parser.ParseFeaturesDir(projectPath)
	if err != nil {
		return fmt.Errorf("failed to list features: %w", err)
	}

	if len(features) == 0 {
		fmt.Println("No features found in the features/ directory")
		fmt.Println("   Use 'archie setup' to create and manage features")
		return nil
	}

	// Find similar features using edit distance
	// Allow up to 3 character differences for suggestions
	similarFeatures := status.FindSimilarFeatures(featureKey, features, 3)

	if len(similarFeatures) > 0 {
		fmt.Println("Did you mean:")
		fmt.Println()
		for i, match := range similarFeatures {
			if i >= 5 { // Show at most 5 suggestions
				break
			}
			fmt.Printf("   • %s\n", match.Name)
		}
		fmt.Println()
	}

	// Show all available features
	fmt.Println("Available features:")
	fmt.Println()

	// Group features by status for better organization
	statusGroups := make(map[status.FeatureStatus][]string)
	for _, feature := range features {
		statusGroups[feature.Status] = append(statusGroups[feature.Status], feature.Name)
	}

	// Display features organized by status
	displayOrder := []status.FeatureStatus{
		status.StatusImplementing,
		status.StatusSpecReady,
		status.StatusDesigned,
		status.StatusUnderDesign,
		status.StatusReadyForDesign,
		status.StatusUnderReview,
		status.StatusNotReviewed,
		status.StatusBlocked,
		status.StatusFinished,
	}

	for _, s := range displayOrder {
		if features, ok := statusGroups[s]; ok && len(features) > 0 {
			statusColor := status.GetStatusColor(s)
			statusName := strings.ReplaceAll(string(s), "_", " ")
			fmt.Printf("   %s%s%s\n", statusColor, statusName, "\033[0m")
			for _, name := range features {
				fmt.Printf("      • %s\n", name)
			}
			fmt.Println()
		}
	}

	fmt.Println("Use: archie status <feature-key> to view details")
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
		fmt.Println("Tip: Use 'archie setup' to create and manage features")
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
