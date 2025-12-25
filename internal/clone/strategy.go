package clone

import (
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/project"
)

// CloneStrategy defines the clone strategy
type CloneStrategy string

const (
	StrategyContext CloneStrategy = "context"
	StrategyLight   CloneStrategy = "light"
	StrategyFull    CloneStrategy = "full"
	StrategyCustom  CloneStrategy = "custom"
)

// StrategyDefinition defines files and directories for a strategy
type StrategyDefinition struct {
	Name        CloneStrategy
	Description string
	Files       []string
	Directories []string
}

// CloneItem represents a file or directory to be cloned
type CloneItem struct {
	Name string
	Type string // "file" or "directory"
}

// Predefined strategies
var strategies = map[CloneStrategy]StrategyDefinition{
	StrategyContext: {
		Name:        StrategyContext,
		Description: "Copy project context and structure (recommended)",
		Files:       []string{"background.md", "dependency.md", "architecture.md", "storage.md", "faq.md"},
		Directories: []string{"api"},
	},
	StrategyLight: {
		Name:        StrategyLight,
		Description: "Copy only background documentation",
		Files:       []string{"background.md"},
		Directories: []string{},
	},
	StrategyFull: {
		Name:        StrategyFull,
		Description: "Copy everything from source project",
		Files:       []string{"background.md", "architecture.md", "dependency.md", "deployment.md", "tasks.md", "metrics.md", "faq.md", "blocker.md", "storage.md"},
		Directories: []string{"api", "workflow", "spec", "features", "assets"},
	},
}

// StrategySelector handles clone strategy selection via TUI
type StrategySelector struct {
	sourcePath string
	fs         afero.Fs
}

// NewStrategySelector creates a new strategy selector
func NewStrategySelector(sourcePath string, fs afero.Fs) *StrategySelector {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &StrategySelector{
		sourcePath: sourcePath,
		fs:         fs,
	}
}

// SelectStrategy shows TUI for strategy selection
// Returns selected strategy and items to copy
func (s *StrategySelector) SelectStrategy() (CloneStrategy, []CloneItem, error) {
	// Build strategy options
	options := []string{
		fmt.Sprintf("context  - %s", strategies[StrategyContext].Description),
		fmt.Sprintf("light    - %s", strategies[StrategyLight].Description),
		fmt.Sprintf("full     - %s", strategies[StrategyFull].Description),
		fmt.Sprintf("custom   - Manually select items to copy"),
	}

	// Show selection UI
	var selected string
	prompt := &survey.Select{
		Message: "Select Clone Strategy:",
		Options: options,
		Default: options[0], // context is default
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", nil, fmt.Errorf("strategy selection cancelled")
	}

	// Parse selected strategy
	var strategy CloneStrategy
	switch {
	case selected == options[0]:
		strategy = StrategyContext
	case selected == options[1]:
		strategy = StrategyLight
	case selected == options[2]:
		strategy = StrategyFull
	case selected == options[3]:
		strategy = StrategyCustom
	}

	// Get items based on strategy
	var items []CloneItem
	var err error

	if strategy == StrategyCustom {
		items, err = s.promptCustomSelection()
		if err != nil {
			return "", nil, err
		}
	} else {
		items = strategyDefToItems(strategies[strategy])
	}

	return strategy, items, nil
}

// promptCustomSelection shows multi-select UI for custom item selection
func (s *StrategySelector) promptCustomSelection() ([]CloneItem, error) {
	// Load standard items from project_structure.json
	standardItems, err := s.loadStandardItems()
	if err != nil {
		return nil, fmt.Errorf("failed to load standard items: %w", err)
	}

	// Build options for multi-select
	options := make([]string, len(standardItems))
	itemMap := make(map[string]CloneItem)

	for i, item := range standardItems {
		option := item.Name
		if item.Type == "directory" {
			option += "/"
		}
		options[i] = option
		itemMap[option] = item
	}

	// Show multi-select UI
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select items to copy (space to select, enter to confirm):",
		Options: options,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, fmt.Errorf("selection cancelled")
	}

	// Convert selected options back to CloneItems
	selectedItems := make([]CloneItem, 0, len(selected))
	for _, option := range selected {
		if item, ok := itemMap[option]; ok {
			selectedItems = append(selectedItems, item)
		}
	}

	return selectedItems, nil
}

// loadStandardItems loads standard items from project_structure.json
// Excludes .archie directory as per user requirement
func (s *StrategySelector) loadStandardItems() ([]CloneItem, error) {
	// Load project structure
	structure, err := project.LoadProjectStructure()
	if err != nil {
		return nil, fmt.Errorf("failed to load project structure: %w", err)
	}

	var items []CloneItem

	// Add files
	for _, file := range structure.Files {
		items = append(items, CloneItem{
			Name: file,
			Type: "file",
		})
	}

	// Add directories (exclude .archie)
	for _, dir := range structure.Directories {
		if dir == ".archie" {
			continue // Skip .archie as per user requirement
		}
		items = append(items, CloneItem{
			Name: dir,
			Type: "directory",
		})
	}

	return items, nil
}

// strategyDefToItems converts a strategy definition to CloneItems
func strategyDefToItems(def StrategyDefinition) []CloneItem {
	var items []CloneItem

	for _, file := range def.Files {
		items = append(items, CloneItem{
			Name: file,
			Type: "file",
		})
	}

	for _, dir := range def.Directories {
		items = append(items, CloneItem{
			Name: dir,
			Type: "directory",
		})
	}

	return items
}

// GetStrategyDefinition returns the definition for a strategy
func GetStrategyDefinition(strategy CloneStrategy) StrategyDefinition {
	return strategies[strategy]
}

// FilterExistingItems filters items that exist in the source directory
func (s *StrategySelector) FilterExistingItems(items []CloneItem) []CloneItem {
	var existingItems []CloneItem

	for _, item := range items {
		itemPath := filepath.Join(s.sourcePath, item.Name)
		exists, err := afero.Exists(s.fs, itemPath)
		if err == nil && exists {
			existingItems = append(existingItems, item)
		}
	}

	return existingItems
}
