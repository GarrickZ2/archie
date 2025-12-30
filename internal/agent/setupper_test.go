package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "expand tilde path",
			input:    "~/.codex/prompts",
			expected: filepath.Join(homeDir, ".codex/prompts"),
			wantErr:  false,
		},
		{
			name:     "expand tilde with file",
			input:    "~/.codex/prompts/test.md",
			expected: filepath.Join(homeDir, ".codex/prompts/test.md"),
			wantErr:  false,
		},
		{
			name:     "no expansion for relative path",
			input:    ".codex/prompts",
			expected: ".codex/prompts",
			wantErr:  false,
		},
		{
			name:     "no expansion for absolute path",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("expandPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	projectPath := "/path/to/project"

	tests := []struct {
		name     string
		basePath string
		relPath  string
		expected string
		wantErr  bool
	}{
		{
			name:     "resolve relative path",
			basePath: projectPath,
			relPath:  ".codex/prompts/test.md",
			expected: filepath.Join(projectPath, ".codex/prompts/test.md"),
			wantErr:  false,
		},
		{
			name:     "resolve tilde path (should ignore basePath)",
			basePath: projectPath,
			relPath:  "~/.codex/prompts/test.md",
			expected: filepath.Join(homeDir, ".codex/prompts/test.md"),
			wantErr:  false,
		},
		{
			name:     "resolve tilde directory path",
			basePath: projectPath,
			relPath:  "~/.codex/prompts",
			expected: filepath.Join(homeDir, ".codex/prompts"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolvePath(tt.basePath, tt.relPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolvePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("resolvePath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestResolvePathIntegration(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	projectPath := "/tmp/test-project"

	// Test case 1: Relative path should be joined with project path
	relPath := filepath.Join(".claude/commands", "design.md")
	fullPath, err := resolvePath(projectPath, relPath)
	if err != nil {
		t.Errorf("resolvePath() error = %v", err)
	}
	expected := filepath.Join(projectPath, ".claude/commands/design.md")
	if fullPath != expected {
		t.Errorf("resolvePath() for relative path: got %v, want %v", fullPath, expected)
	}

	// Test case 2: Tilde path should expand to home directory
	tildeRelPath := filepath.Join("~/.codex/prompts", "design.md")
	fullPath, err = resolvePath(projectPath, tildeRelPath)
	if err != nil {
		t.Errorf("resolvePath() error = %v", err)
	}
	expected = filepath.Join(homeDir, ".codex/prompts/design.md")
	if fullPath != expected {
		t.Errorf("resolvePath() for tilde path: got %v, want %v", fullPath, expected)
	}

	// Verify that tilde path does not include project path
	if strings.Contains(fullPath, projectPath) {
		t.Errorf("resolvePath() for tilde path should not contain project path, got %v", fullPath)
	}
}
