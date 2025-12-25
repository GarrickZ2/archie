package resources

import (
	"testing"
)

func TestLoadCommands(t *testing.T) {
	commands, err := LoadCommands()
	if err != nil {
		t.Fatalf("LoadCommands() error = %v", err)
	}

	if len(commands) == 0 {
		t.Error("LoadCommands() returned empty map")
	}

	// Check for expected command files (without extension, as YAML files are loaded)
	expectedCommands := []string{
		"archie-design",
		"archie-review",
		"archie-fix",
		"archie-spec",
		"archie-init",
		"archie-compact",
		"archie-revise",
		"archie-update-progress",
	}

	for _, cmd := range expectedCommands {
		if template, ok := commands[cmd]; !ok {
			t.Errorf("Expected command %s not found", cmd)
		} else if template == nil {
			t.Errorf("Command %s has nil template", cmd)
		} else if template.Content == "" {
			t.Errorf("Command %s has empty content", cmd)
		}
	}

	t.Logf("Loaded %d commands", len(commands))
}

func TestLoadSubAgents(t *testing.T) {
	subAgents, err := LoadSubAgents()
	if err != nil {
		t.Fatalf("LoadSubAgents() error = %v", err)
	}

	if len(subAgents) == 0 {
		t.Error("LoadSubAgents() returned empty map")
	}

	// Check for expected subagent files
	expectedSubAgents := []string{
		"archie-api.md",
		"archie-workflow.md",
		"archie-storage.md",
		"archie-metrics.md",
		"archie-task.md",
	}

	for _, sa := range expectedSubAgents {
		if content, ok := subAgents[sa]; !ok {
			t.Errorf("Expected subagent %s not found", sa)
		} else if content == "" {
			t.Errorf("Subagent %s has empty content", sa)
		}
	}

	t.Logf("Loaded %d subagents", len(subAgents))
}
