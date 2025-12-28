package export

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/status"
)

// DocumentSelector handles document selection via TUI
type DocumentSelector struct {
	projectPath string
	fs          afero.Fs
	parser      *status.Parser
}

// NewDocumentSelector creates a new document selector
func NewDocumentSelector(projectPath string, fs afero.Fs) *DocumentSelector {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	return &DocumentSelector{
		projectPath: projectPath,
		fs:          fs,
		parser:      status.NewParser(fs),
	}
}

// SelectDocuments shows TUI for document selection
// Returns the export configuration based on user choices
func (s *DocumentSelector) SelectDocuments(flagOutputPath string, flagTOC, flagStats, flagDepGraph bool) (*ExportConfig, error) {
	config := &ExportConfig{
		ProjectPath:      s.projectPath,
		GenerateTOC:      flagTOC,
		GenerateStats:    flagStats,
		GenerateDepGraph: flagDepGraph,
	}

	// Step 1: Select root documents
	rootDocs, err := s.selectRootDocuments()
	if err != nil {
		return nil, err
	}
	config.IncludeRoot = rootDocs

	// Step 2: Select features
	features, err := s.selectFeatures()
	if err != nil {
		return nil, err
	}
	config.IncludeFeatures = features

	// Step 3: Select additional options
	includeWorkflows, includeSpecs, err := s.selectAdditionalOptions(config)
	if err != nil {
		return nil, err
	}
	config.IncludeWorkflows = includeWorkflows
	config.IncludeSpecs = includeSpecs

	// Step 4: Confirm output path
	outputPath, err := s.confirmOutputPath(flagOutputPath)
	if err != nil {
		return nil, err
	}
	config.OutputPath = outputPath

	return config, nil
}

// selectRootDocuments shows multi-select for root documents
func (s *DocumentSelector) selectRootDocuments() ([]string, error) {
	// Define available root documents
	availableDocs := map[string]string{
		"background.md":   "Project background and context",
		"dependency.md":   "Dependency catalog",
		"deployment.md":   "Deployment and release documentation",
		"storage.md":      "Storage design",
		"tasks.md":        "Task tracking",
		"metrics.md":      "Observability metrics",
		"api/":            "API documentation",
		"blocker.md":      "Blockers documentation",
		"architecture.md": "Architecture documentation",
	}

	// Check which documents exist
	var options []string
	docMap := make(map[string]string)

	for doc, desc := range availableDocs {
		docPath := filepath.Join(s.projectPath, doc)
		exists, err := afero.Exists(s.fs, docPath)
		if err == nil && exists {
			option := fmt.Sprintf("%-18s - %s", doc, desc)
			options = append(options, option)
			docMap[option] = doc
		}
	}

	if len(options) == 0 {
		// No root documents found, skip this step
		return []string{}, nil
	}

	// Show multi-select
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select root documents to export (space to select, enter to confirm):",
		Options: options,
		Default: options, // Select all by default
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, fmt.Errorf("root document selection cancelled")
	}

	// Convert back to document names
	selectedDocs := make([]string, 0, len(selected))
	for _, option := range selected {
		if doc, ok := docMap[option]; ok {
			selectedDocs = append(selectedDocs, doc)
		}
	}

	return selectedDocs, nil
}

// selectFeatures shows multi-select for features
func (s *DocumentSelector) selectFeatures() ([]string, error) {
	// Parse all features
	features, err := s.parser.ParseFeaturesDir(s.projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse features: %w", err)
	}

	if len(features) == 0 {
		// No features found, skip this step
		return []string{}, nil
	}

	// Build options with status information
	options := make([]string, len(features))
	featureMap := make(map[string]string)

	for i, feature := range features {
		statusStr := strings.ReplaceAll(string(feature.Status), "_", " ")
		option := fmt.Sprintf("%-30s [%s]", feature.Name, statusStr)
		options[i] = option
		featureMap[option] = feature.Name
	}

	// Show multi-select
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select features to export (space to select, enter to confirm):",
		Options: options,
		Default: options, // Select all by default
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, fmt.Errorf("feature selection cancelled")
	}

	// Convert back to feature names
	selectedFeatures := make([]string, 0, len(selected))
	for _, option := range selected {
		if featureName, ok := featureMap[option]; ok {
			selectedFeatures = append(selectedFeatures, featureName)
		}
	}

	return selectedFeatures, nil
}

// selectAdditionalOptions shows multi-select for additional export options
func (s *DocumentSelector) selectAdditionalOptions(config *ExportConfig) (includeWorkflows, includeSpecs bool, err error) {
	// Only show workflow/spec options if features are selected
	if len(config.IncludeFeatures) == 0 {
		return false, false, nil
	}

	options := []string{
		"Include workflow diagrams (.mmd files)",
		"Include specification files (spec/)",
	}

	// Prepare default selections
	defaults := []string{}
	if s.hasWorkflows(config.IncludeFeatures) {
		defaults = append(defaults, options[0])
	}
	if s.hasSpecs(config.IncludeFeatures) {
		defaults = append(defaults, options[1])
	}

	// Show multi-select
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select additional content to include:",
		Options: options,
		Default: defaults,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return false, false, fmt.Errorf("options selection cancelled")
	}

	// Parse selections
	for _, s := range selected {
		if s == options[0] {
			includeWorkflows = true
		} else if s == options[1] {
			includeSpecs = true
		}
	}

	return includeWorkflows, includeSpecs, nil
}

// confirmOutputPath asks for output file path
func (s *DocumentSelector) confirmOutputPath(flagPath string) (string, error) {
	// Generate default path
	defaultPath := flagPath
	if defaultPath == "" {
		defaultPath = fmt.Sprintf("./archie-export-%s.md", time.Now().Format("2006-01-02"))
	}

	// Ask user for output path
	var outputPath string
	prompt := &survey.Input{
		Message: "Output file path:",
		Default: defaultPath,
	}

	if err := survey.AskOne(prompt, &outputPath, survey.WithValidator(survey.Required)); err != nil {
		return "", fmt.Errorf("output path selection cancelled")
	}

	// Check if file already exists
	exists, err := afero.Exists(s.fs, outputPath)
	if err == nil && exists {
		var overwrite bool
		confirmPrompt := &survey.Confirm{
			Message: fmt.Sprintf("File %s already exists. Overwrite?", outputPath),
			Default: false,
		}

		if err := survey.AskOne(confirmPrompt, &overwrite); err != nil {
			return "", fmt.Errorf("confirmation cancelled")
		}

		if !overwrite {
			// Add timestamp to avoid overwriting
			ext := filepath.Ext(outputPath)
			base := outputPath[:len(outputPath)-len(ext)]
			outputPath = fmt.Sprintf("%s-%d%s", base, time.Now().Unix(), ext)
		}
	}

	return outputPath, nil
}

// hasWorkflows checks if any of the features have workflow directories
func (s *DocumentSelector) hasWorkflows(features []string) bool {
	for _, feature := range features {
		workflowPath := filepath.Join(s.projectPath, "workflow", feature)
		exists, err := afero.DirExists(s.fs, workflowPath)
		if err == nil && exists {
			return true
		}
	}
	return false
}

// hasSpecs checks if any of the features have spec files
func (s *DocumentSelector) hasSpecs(features []string) bool {
	for _, feature := range features {
		specPath := filepath.Join(s.projectPath, "spec", feature+".spec.md")
		exists, err := afero.Exists(s.fs, specPath)
		if err == nil && exists {
			return true
		}
	}
	return false
}
