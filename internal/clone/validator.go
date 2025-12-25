package clone

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/agent"
)

// CloneValidator validates source and target paths for cloning
type CloneValidator struct {
	fs afero.Fs
}

// ValidationResult contains the validation result and loaded metadata
type ValidationResult struct {
	SourcePath   string
	TargetPath   string
	SourceAgents []agent.AgentState
}

// NewCloneValidator creates a new validator
func NewCloneValidator(fs afero.Fs) *CloneValidator {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &CloneValidator{fs: fs}
}

// ValidateSource validates the source path
// Returns error if:
// - Path doesn't exist
// - Path is not a directory
// - .archie/ folder doesn't exist in source
func (v *CloneValidator) ValidateSource(sourcePath string) error {
	// Check if path exists
	exists, err := afero.Exists(v.fs, sourcePath)
	if err != nil {
		return fmt.Errorf("failed to check source path: %w", err)
	}
	if !exists {
		return fmt.Errorf("source path does not exist: %s", sourcePath)
	}

	// Check if path is a directory
	isDir, err := afero.IsDir(v.fs, sourcePath)
	if err != nil {
		return fmt.Errorf("failed to check if source is directory: %w", err)
	}
	if !isDir {
		return fmt.Errorf("source path is not a directory: %s", sourcePath)
	}

	// Check if .archie/ exists
	archieDir := filepath.Join(sourcePath, ".archie")
	archieExists, err := afero.DirExists(v.fs, archieDir)
	if err != nil {
		return fmt.Errorf("failed to check .archie directory: %w", err)
	}
	if !archieExists {
		return fmt.Errorf("source is not an archie project (missing .archie folder): %s", sourcePath)
	}

	return nil
}

// ValidateTarget validates the target path
// Returns error if:
// - Source and target are the same (after resolving absolute paths)
func (v *CloneValidator) ValidateTarget(sourcePath, targetPath string) error {
	// Resolve absolute paths
	absSource, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to resolve source path: %w", err)
	}

	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	// Check if source and target are the same
	if absSource == absTarget {
		return fmt.Errorf("source and target cannot be the same directory")
	}

	// Note: We don't check if target is empty or has .archie/ because
	// the Initializer will handle that validation (same as init command)

	return nil
}

// Validate performs full validation of source and target
// Also loads source project's initialized agents
func (v *CloneValidator) Validate(sourcePath, targetPath string) (*ValidationResult, error) {
	// Validate source
	if err := v.ValidateSource(sourcePath); err != nil {
		return nil, err
	}

	// Validate target
	if err := v.ValidateTarget(sourcePath, targetPath); err != nil {
		return nil, err
	}

	// Resolve absolute paths for result
	absSource, _ := filepath.Abs(sourcePath)
	absTarget, _ := filepath.Abs(targetPath)

	// Load source project's agent states
	stateManager := agent.NewStateManager(v.fs)
	sourceAgents, err := stateManager.GetInitializedAgents(absSource)
	if err != nil {
		// Non-fatal: If we can't load agent states, continue with empty list
		// This might happen if state.json doesn't exist or is malformed
		sourceAgents = []agent.AgentState{}
	}

	return &ValidationResult{
		SourcePath:   absSource,
		TargetPath:   absTarget,
		SourceAgents: sourceAgents,
	}, nil
}
