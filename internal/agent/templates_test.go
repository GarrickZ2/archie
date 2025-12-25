package agent

import (
	"testing"
)

func TestGetCommands(t *testing.T) {
	commands := GetCommands()

	if len(commands) == 0 {
		t.Error("GetCommands() returned empty map")
	}

	t.Logf("GetCommands() returned %d commands", len(commands))

	// Verify all commands have non-empty content
	for filename, template := range commands {
		if template == nil {
			t.Errorf("Command %s has nil template", filename)
			continue
		}
		if template.Content == "" {
			t.Errorf("Command %s has empty content", filename)
		}
	}
}

func TestGetSubAgents(t *testing.T) {
	subAgents := GetSubAgents()

	if len(subAgents) == 0 {
		t.Error("GetSubAgents() returned empty map")
	}

	t.Logf("GetSubAgents() returned %d subagents", len(subAgents))

	// Verify all subagents have non-empty content
	for filename, content := range subAgents {
		if content == "" {
			t.Errorf("SubAgent %s has empty content", filename)
		}
	}
}
