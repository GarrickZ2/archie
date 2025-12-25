package cmd

import (
	"github.com/spf13/cobra"

	"github.com/GarrickZ2/archie/internal/agent"
)

var customAgentCmd = &cobra.Command{
	Use:   "custom-agent",
	Short: "Manage custom agents",
	Long: `Manage custom agents stored in ~/.archie/custom_agents.json.

This command provides a TUI interface to:
1. Add new custom agents
2. Remove existing custom agents
3. List all custom agents

Custom agents are stored globally in your home directory,
so they can be reused across all projects.`,
	Args: cobra.NoArgs,
	RunE: runCustomAgent,
}

func init() {
	rootCmd.AddCommand(customAgentCmd)
}

func runCustomAgent(cmd *cobra.Command, args []string) error {
	manager := agent.NewCustomAgentManager()
	return manager.ShowManagementUI()
}
