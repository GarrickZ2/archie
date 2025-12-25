package project

import (
	"fmt"
	"path/filepath"

	"github.com/GarrickZ2/archie/resources"
	"github.com/spf13/afero"
)

// Initializer 定义项目初始化接口
type Initializer interface {
	Initialize(targetPath string) error
}

// Config 初始化器配置
type Config struct {
	FileSystem afero.Fs      // 文件系统实现
	Validator  PathValidator // 路径验证器
}

// DefaultInitializer 默认的项目初始化器实现
type DefaultInitializer struct {
	fs        afero.Fs
	validator PathValidator
}

// NewInitializer 创建项目初始化器
// cfg 为 nil 时将使用默认配置
func NewInitializer(cfg *Config) Initializer {
	if cfg == nil {
		cfg = &Config{}
	}

	// 提供默认实现
	if cfg.FileSystem == nil {
		cfg.FileSystem = afero.NewOsFs()
	}
	if cfg.Validator == nil {
		cfg.Validator = NewPathValidator(cfg.FileSystem)
	}

	return &DefaultInitializer{
		fs:        cfg.FileSystem,
		validator: cfg.Validator,
	}
}

// Initialize initializes a new project
// targetPath is the target folder path
// If folder doesn't exist, create it
// If exists and contains .archie, continue (create missing files)
// If exists but no .archie, return error
func (d *DefaultInitializer) Initialize(targetPath string) error {
	// Step 1: Check directory status
	hasArchie, err := d.checkDirectory(targetPath)
	if err != nil {
		return err
	}

	// Step 2: Ensure target directory exists
	if err := d.fs.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Step 3: If directory exists but is not an archie project, return error
	dirExists, _ := afero.DirExists(d.fs, targetPath)
	if dirExists && !hasArchie {
		entries, err := afero.ReadDir(d.fs, targetPath)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}
		// Allow if directory is empty or only contains .archie
		if len(entries) > 0 && !(len(entries) == 1 && entries[0].Name() == ".archie") {
			return fmt.Errorf("directory is not empty and not an archie project (missing .archie folder)")
		}
	}

	// Step 4: Create project structure (skip existing files)
	if err := d.createProjectStructure(targetPath); err != nil {
		return err
	}

	// Step 5: Copy docs to .archie/docs/
	if err := resources.CopyDocsToProject(targetPath, d.fs); err != nil {
		return fmt.Errorf("failed to copy docs: %w", err)
	}

	return nil
}

// createProjectStructure creates project file structure (using project_structure.json)
func (d *DefaultInitializer) createProjectStructure(targetPath string) error {
	// Load project structure configuration
	structure, err := LoadProjectStructure()
	if err != nil {
		return fmt.Errorf("failed to load project structure: %w", err)
	}

	// 1. Create directories
	for _, dir := range structure.Directories {
		fullPath := filepath.Join(targetPath, dir)
		if err := d.fs.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	// 2. Create files (AGENTS.md uses embedded content, others are empty)
	for _, file := range structure.Files {
		fullPath := filepath.Join(targetPath, file)

		// AGENTS.md: force overwrite; others: skip if exists
		isAgentsDoc := file == "AGENTS.md"
		if !isAgentsDoc {
			// Check if file already exists
			exists, err := afero.Exists(d.fs, fullPath)
			if err != nil {
				return fmt.Errorf("failed to check file %s: %w", fullPath, err)
			}
			if exists {
				continue // Skip existing files
			}
		}

		// Create parent directory
		dir := filepath.Dir(fullPath)
		if err := d.fs.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Determine file content
		var content string
		if isAgentsDoc {
			content = GetAgentsMdContent()
		}

		// Create file (AGENTS.md overwrites, others are empty)
		if err := afero.WriteFile(d.fs, fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", fullPath, err)
		}
	}

	// 3. Create files with content
	for file, content := range structure.FilesWithContent {
		fullPath := filepath.Join(targetPath, file)

		// Check if file already exists
		exists, err := afero.Exists(d.fs, fullPath)
		if err != nil {
			return fmt.Errorf("failed to check file %s: %w", fullPath, err)
		}
		if exists {
			continue // Skip existing files
		}

		// Create parent directory
		dir := filepath.Dir(fullPath)
		if err := d.fs.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Create file
		if err := afero.WriteFile(d.fs, fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", fullPath, err)
		}
	}

	return nil
}

// checkDirectory checks directory status
// Returns: (hasArchie, error)
func (d *DefaultInitializer) checkDirectory(targetPath string) (bool, error) {
	// Check if directory exists
	exists, err := afero.DirExists(d.fs, targetPath)
	if err != nil {
		return false, err
	}

	// If directory doesn't exist, return false
	if !exists {
		return false, nil
	}

	// Check if .archie folder exists
	archieDir := filepath.Join(targetPath, ".archie")
	hasArchie, err := afero.DirExists(d.fs, archieDir)
	if err != nil {
		return false, err
	}

	return hasArchie, nil
}
