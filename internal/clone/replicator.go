package clone

import (
	"context"
	"fmt"

	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/agent"
)

// AgentReplicator replicates agent setup from source to target
type AgentReplicator struct {
	fs           afero.Fs
	stateManager *agent.StateManager
	setupper     agent.Setupper
}

// ReplicationOptions defines options for agent replication
type ReplicationOptions struct {
	SourcePath   string
	TargetPath   string
	SourceAgents []agent.AgentState
}

// ReplicationResult contains the result of agent replication
type ReplicationResult struct {
	ReplicatedAgents []string
	SkippedAgents    []string
	Errors           []error
}

// NewAgentReplicator creates a new agent replicator
func NewAgentReplicator(fs afero.Fs) *AgentReplicator {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &AgentReplicator{
		fs:           fs,
		stateManager: agent.NewStateManager(fs),
		setupper:     agent.NewSetupper(fs),
	}
}

// Replicate replicates agent setup from source to target
func (r *AgentReplicator) Replicate(opts ReplicationOptions) (*ReplicationResult, error) {
	result := &ReplicationResult{
		ReplicatedAgents: []string{},
		SkippedAgents:    []string{},
		Errors:           []error{},
	}

	ctx := context.Background()

	for _, agentState := range opts.SourceAgents {
		// Check if agent is available in registry
		available, err := r.verifyAgentAvailability(agentState.AgentName)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to verify agent %s: %w", agentState.AgentName, err))
			continue
		}

		if !available {
			result.SkippedAgents = append(result.SkippedAgents, agentState.AgentName)
			continue
		}

		// Replicate agent
		if err := r.replicateAgent(ctx, agentState, opts.TargetPath); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to replicate agent %s: %w", agentState.AgentName, err))
			continue
		}

		result.ReplicatedAgents = append(result.ReplicatedAgents, agentState.AgentName)
	}

	return result, nil
}

// verifyAgentAvailability checks if an agent is available in the registry
func (r *AgentReplicator) verifyAgentAvailability(agentName string) (bool, error) {
	_, err := agent.Get(agentName)
	if err != nil {
		// Agent not found
		return false, nil
	}
	return true, nil
}

// replicateAgent replicates a single agent setup
func (r *AgentReplicator) replicateAgent(ctx context.Context, agentState agent.AgentState, targetPath string) error {
	// Setup agent in target project
	setupConfig := agent.SetupConfig{
		ProjectPath: targetPath,
		AgentType:   agentState.AgentName,
		Options: agent.SetupOptions{
			IncludeExamples: true,
		},
	}

	if err := r.setupper.Setup(ctx, setupConfig); err != nil {
		return fmt.Errorf("agent setup failed: %w", err)
	}

	// Mark as initialized in target project
	if err := r.stateManager.MarkInitialized(targetPath, agentState.AgentName, agentState.IsCustom); err != nil {
		return fmt.Errorf("failed to mark agent as initialized: %w", err)
	}

	return nil
}
