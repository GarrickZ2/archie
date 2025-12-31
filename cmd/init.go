package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/GarrickZ2/archie/internal/agent"
	"github.com/GarrickZ2/archie/internal/project"
	"github.com/GarrickZ2/archie/internal/ui"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new technical design documentation project",
	Long: `Initialize a new technical design documentation project in the current directory.

This command will:
1. Bootstrap the project structure
2. Let you select and configure a Code Agent (claude-code or custom)
3. Generate agent-specific commands and sub-agents

If the directory is not empty, you'll be prompted for confirmation.`,
	Args: cobra.NoArgs,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Show welcome banner with ARCHIE logo
	ui.ShowWelcomeBanner()

	// Use current directory
	targetPath, err := os.Getwd()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get current directory: %v", err))
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Initialize project structure
	initializer := project.NewInitializer(nil)
	if err := initializer.Initialize(targetPath); err != nil {
		ui.ShowError(fmt.Sprintf("Initialization failed: %v", err))
		return fmt.Errorf("project initialization failed: %w", err)
	}

	// Agent selection and setup via TUI
	selector := agent.NewTUISelector(targetPath)

	// Select agent
	agentName, isCustom, err := selector.SelectAgent()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Agent selection failed: %v", err))
		return fmt.Errorf("agent selection failed: %w", err)
	}

	fmt.Println()

	// Check if already initialized
	stateManager := agent.NewStateManager(nil)
	isInitialized, err := stateManager.IsInitialized(targetPath, agentName)
	if err != nil {
		return fmt.Errorf("failed to check agent state: %w", err)
	}

	// Confirm reconfigure if needed
	if isInitialized {
		shouldReconfigure, err := selector.ConfirmReconfigure(agentName)
		if err != nil {
			return err
		}
		if !shouldReconfigure {
			return nil
		}
		fmt.Println()
	}

	// Setup agent
	selectedAgent, err := agent.Get(agentName)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get agent: %v", err))
		return fmt.Errorf("failed to get agent: %w", err)
	}

	setupper := agent.NewSetupper(nil)
	setupConfig := agent.SetupConfig{
		ProjectPath: targetPath,
		AgentType:   agentName,
		Options: agent.SetupOptions{
			IncludeExamples: true,
		},
	}

	if err := setupper.Setup(ctx, setupConfig); err != nil {
		ui.ShowError(fmt.Sprintf("Agent setup failed: %v", err))
		return fmt.Errorf("agent setup failed: %w", err)
	}

	// Mark as initialized
	if err := stateManager.MarkInitialized(targetPath, agentName, isCustom); err != nil {
		return fmt.Errorf("failed to update state: %w", err)
	}

	fmt.Println()

	// Show success message with all initialized agents
	showSuccessBox(targetPath, agentName, selectedAgent.PathConfig())

	return nil
}

// showSuccessBox displays a nicely formatted success message
func showSuccessBox(projectPath string, currentAgentName string, currentPathConfig agent.PathConfig) {
	// Get all initialized agents
	stateManager := agent.NewStateManager(nil)
	initializedAgents, err := stateManager.GetInitializedAgents(projectPath)
	if err != nil {
		// Fallback to showing only current agent if error
		initializedAgents = []agent.AgentState{}
	}

	// Sort agents by name alphabetically
	sort.Slice(initializedAgents, func(i, j int) bool {
		return initializedAgents[i].AgentName < initializedAgents[j].AgentName
	})

	content := []string{}

	// Show all initialized agents
	if len(initializedAgents) > 0 {
		content = append(content, ui.ColorGreen+"Initialized Agents:"+ui.ColorReset)
		for _, ag := range initializedAgents {
			agentType := "Official"
			if ag.IsCustom {
				agentType = "Custom"
			}

			// Highlight current agent
			if ag.AgentName == currentAgentName {
				content = append(content, fmt.Sprintf("  %s%s%s (%s) %s← current%s",
					ui.ColorBold, ag.AgentName, ui.ColorReset, agentType,
					ui.ColorBrightYellow, ui.ColorReset))
			} else {
				content = append(content, fmt.Sprintf("  %s (%s)", ag.AgentName, agentType))
			}
		}
		content = append(content, "")
	}

	content = append(content, ui.ColorGreen+"Generated:"+ui.ColorReset)
	content = append(content, fmt.Sprintf("  Commands: %s/", currentPathConfig.CommandsDir))

	if currentPathConfig.SubAgentsDir != "" {
		content = append(content, fmt.Sprintf("  Sub-Agents: %s/", currentPathConfig.SubAgentsDir))
	}

	content = append(content, "")
	content = append(content, ui.ColorBrightYellow+"Next Steps:"+ui.ColorReset)
	content = append(content, "  Review generated files and start designing")

	ui.PrintBox("Setup Complete!", content)
}
