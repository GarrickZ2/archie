package agent

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/resources"
)

// DefaultSetupper 默认的 setupper 实现
type DefaultSetupper struct {
	fs afero.Fs
}

// NewSetupper 创建 setupper
func NewSetupper(fs afero.Fs) Setupper {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &DefaultSetupper{fs: fs}
}

// Setup sets up agent configuration in the project
func (s *DefaultSetupper) Setup(ctx context.Context, config SetupConfig) error {
	// Get agent
	agent, err := Get(config.AgentType)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}

	pathConfig := agent.PathConfig()

	// 1. Get agent configuration (for background doc)
	var agentDoc string
	if customAgent, ok := agent.(*CustomAgent); ok {
		agentDoc = customAgent.config.AgentDoc
	} else {
		// For non-CustomAgent, default to AGENTS.md
		agentDoc = "AGENTS.md"
	}

	// 2. Force overwrite agent doc (AGENTS.md or custom doc) if specified
	if agentDoc != "" {
		agentDocPath := filepath.Join(config.ProjectPath, agentDoc)
		if err = s.writeFile(agentDocPath, resources.AgentsMdContent); err != nil {
			return fmt.Errorf("failed to sync to %s: %w", agentDoc, err)
		}
	}

	// 3. Install commands to {commands_dir}/
	commands := agent.Commands()
	for filename, content := range commands {
		relPath := filepath.Join(pathConfig.CommandsDir, filename)
		fullPath := filepath.Join(config.ProjectPath, relPath)

		if err = s.writeFile(fullPath, content); err != nil {
			return err
		}
	}

	// 4. Install sub-agents to {sub_agents_dir}/ (if supported)
	if agent.SupportsSubAgents() {
		subAgents := agent.SubAgents()
		for filename, content := range subAgents {
			relPath := filepath.Join(pathConfig.SubAgentsDir, filename)
			fullPath := filepath.Join(config.ProjectPath, relPath)

			if err = s.writeFile(fullPath, content); err != nil {
				return err
			}
		}
	}

	return nil
}

// writeFile writes a single file
func (s *DefaultSetupper) writeFile(fullPath, content string) error {
	// Create parent directory
	dir := filepath.Dir(fullPath)
	if err := s.fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := afero.WriteFile(s.fs, fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return nil
}
