package resources

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// Embed all command templates folder
//
//go:embed commands
var commandsFS embed.FS

// Embed all subagent templates folder
//
//go:embed subagents
var subagentsFS embed.FS

// Embed AGENTS.md from agents/
//
//go:embed agents/AGENTS.md
var AgentsMdContent string

// Embed docs folder
//go:embed docs
var docsFS embed.FS

// CommandTemplate 表示一个命令模板的结构
type CommandTemplate struct {
	// Metadata 包含所有非 content 的字段
	Metadata map[string]interface{} `yaml:",inline"`
	// Content 是命令的正文内容
	Content string `yaml:"content"`
}

// LoadCommands loads all command templates from the embedded commands folder
// Returns a map where key is filename (without extension) and value is CommandTemplate
func LoadCommands() (map[string]*CommandTemplate, error) {
	return loadCommandTemplatesFromFS(commandsFS, "commands")
}

// loadCommandTemplatesFromFS loads YAML command templates from embedded FS
func loadCommandTemplatesFromFS(fsys embed.FS, rootDir string) (map[string]*CommandTemplate, error) {
	templates := make(map[string]*CommandTemplate)

	err := fs.WalkDir(fsys, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-yaml files
		if d.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		// Read file content
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Parse YAML
		var template CommandTemplate
		if err := yaml.Unmarshal(content, &template); err != nil {
			return fmt.Errorf("failed to parse YAML file %s: %w", path, err)
		}

		// Use relative path as key (remove root directory prefix and extension)
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}
		// Remove .yaml extension
		key := relPath[:len(relPath)-5] // Remove ".yaml"

		templates[key] = &template
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load templates from %s: %w", rootDir, err)
	}

	return templates, nil
}

// LoadSubAgents loads all subagent templates from the embedded subagents folder
// Returns a map where key is filename and value is file content
func LoadSubAgents() (map[string]string, error) {
	return loadTemplatesFromFS(subagentsFS, "subagents")
}

// loadTemplatesFromFS is a helper function to load all .md files from an embedded FS
// This is kept for backward compatibility with subagents
func loadTemplatesFromFS(fsys embed.FS, rootDir string) (map[string]string, error) {
	templates := make(map[string]string)

	err := fs.WalkDir(fsys, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// Read file content
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Use relative path as key (remove root directory prefix)
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		templates[relPath] = string(content)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load templates from %s: %w", rootDir, err)
	}

	return templates, nil
}

// GetSchemaTemplate 获取指定的 schema 模板内容
// templateName: "background.md", "feature.md" 等
func GetSchemaTemplate(templateName string) (string, error) {
	path := filepath.Join("docs", "schema", templateName)
	content, err := fs.ReadFile(docsFS, path)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templateName, err)
	}
	return string(content), nil
}

// CopyDocsToProject copies the embedded docs folder to .archie/docs/ in the project
func CopyDocsToProject(projectPath string, filesystem afero.Fs) error {
	targetDir := filepath.Join(projectPath, ".archie", "docs")

	// Create target directory
	if err := filesystem.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	// Walk through embedded docs folder
	err := fs.WalkDir(docsFS, "docs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root "docs" directory itself
		if path == "docs" {
			return nil
		}

		// Calculate relative path (remove "docs/" prefix)
		relPath, err := filepath.Rel("docs", path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			// Create directory
			return filesystem.MkdirAll(targetPath, 0755)
		}

		// Read file content from embedded FS
		content, err := fs.ReadFile(docsFS, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write file to target location
		if err := afero.WriteFile(filesystem, targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy docs: %w", err)
	}

	return nil
}
