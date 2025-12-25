package clone

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/project"
	"github.com/GarrickZ2/archie/internal/ui"
)

// CloneManager orchestrates the entire clone workflow
type CloneManager struct {
	fs               afero.Fs
	validator        *CloneValidator
	strategySelector *StrategySelector
	copyManager      *CopyManager
	replicator       *AgentReplicator
	initializer      project.Initializer
}

// CloneConfig defines the configuration for cloning
type CloneConfig struct {
	SourcePath string
	TargetPath string
}

// CloneResult contains the result of the clone operation
type CloneResult struct {
	Strategy            CloneStrategy
	CopyResult          *CopyResult
	ReplicationResult   *ReplicationResult
	InitializationError error
}

// NewCloneManager creates a new clone manager
func NewCloneManager(fs afero.Fs) *CloneManager {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	return &CloneManager{
		fs:          fs,
		validator:   NewCloneValidator(fs),
		copyManager: NewCopyManager(fs),
		replicator:  NewAgentReplicator(fs),
		initializer: project.NewInitializer(nil),
	}
}

// Clone executes the complete clone workflow
func (m *CloneManager) Clone(config CloneConfig) (*CloneResult, error) {
	result := &CloneResult{}

	validationResult, err := m.validator.Validate(config.SourcePath, config.TargetPath)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println()

	m.strategySelector = NewStrategySelector(validationResult.SourcePath, m.fs)

	strategy, itemsToCopy, err := m.strategySelector.SelectStrategy()
	if err != nil {
		return nil, fmt.Errorf("strategy selection failed: %w", err)
	}

	result.Strategy = strategy
	itemsToCopy = m.strategySelector.FilterExistingItems(itemsToCopy)

	fmt.Println()

	copyOpts := CopyOptions{
		SourcePath: validationResult.SourcePath,
		TargetPath: validationResult.TargetPath,
		Items:      itemsToCopy,
		Overwrite:  false,
	}

	copyResult, err := m.copyManager.Copy(copyOpts)
	if err != nil {
		return nil, fmt.Errorf("copy operation failed: %w", err)
	}

	result.CopyResult = copyResult

	if len(copyResult.Errors) > 0 {
		ui.ShowInfo(fmt.Sprintf("⚠️  %d errors during copying", len(copyResult.Errors)))
		for _, err := range copyResult.Errors {
			fmt.Printf("  • %v\n", err)
		}
		fmt.Println()
	}

	if err := m.initializer.Initialize(validationResult.TargetPath); err != nil {
		result.InitializationError = err
		return nil, fmt.Errorf("initialization failed: %w", err)
	}

	if len(validationResult.SourceAgents) > 0 {
		replicationOpts := ReplicationOptions{
			SourcePath:   validationResult.SourcePath,
			TargetPath:   validationResult.TargetPath,
			SourceAgents: validationResult.SourceAgents,
		}

		replicationResult, err := m.replicator.Replicate(replicationOpts)
		if err != nil {
			ui.ShowInfo(fmt.Sprintf("⚠️  Agent replication errors: %v", err))
		}

		result.ReplicationResult = replicationResult

		if len(replicationResult.Errors) > 0 {
			for _, err := range replicationResult.Errors {
				fmt.Printf("  • %v\n", err)
			}
			fmt.Println()
		}

		if len(replicationResult.SkippedAgents) > 0 {
			ui.ShowInfo(fmt.Sprintf("⚠️  Skipped %d unavailable agent(s): %s",
				len(replicationResult.SkippedAgents),
				strings.Join(replicationResult.SkippedAgents, ", ")))
			fmt.Println("  💡 Tip: Add missing agents with: archie custom-agent")
			fmt.Println()
		}
	}

	m.showCloneSuccessBox(result)

	return result, nil
}

// showCloneSuccessBox displays a success message box
func (m *CloneManager) showCloneSuccessBox(result *CloneResult) {
	content := []string{
		fmt.Sprintf("Strategy: %s%s%s", ui.ColorBold, result.Strategy, ui.ColorReset),
		"",
	}

	if len(result.CopyResult.CopiedFiles) > 0 {
		content = append(content, ui.ColorGreen+"Files:"+ui.ColorReset)
		for _, file := range result.CopyResult.CopiedFiles {
			content = append(content, fmt.Sprintf("  %s", file))
		}
		content = append(content, "")
	}

	if len(result.CopyResult.CopiedDirs) > 0 {
		content = append(content, ui.ColorGreen+"Directories:"+ui.ColorReset)
		for _, dir := range result.CopyResult.CopiedDirs {
			content = append(content, fmt.Sprintf("  %s/", dir))
		}
		content = append(content, "")
	}

	if result.ReplicationResult != nil && len(result.ReplicationResult.ReplicatedAgents) > 0 {
		content = append(content, ui.ColorGreen+"Agents:"+ui.ColorReset)
		for _, agentName := range result.ReplicationResult.ReplicatedAgents {
			content = append(content, fmt.Sprintf("  %s", agentName))
		}
		content = append(content, "")
	}

	content = append(content, ui.ColorBrightYellow+"Next Steps:"+ui.ColorReset)
	content = append(content, "  Review background.md and start designing")

	ui.PrintBox("Clone Complete!", content)
}
