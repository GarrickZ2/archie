package project

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed project_structure.json
var projectStructureJSON []byte

// ProjectStructure 项目结构配置
type ProjectStructure struct {
	Files            []string          `json:"files"`
	Directories      []string          `json:"directories"`
	FilesWithContent map[string]string `json:"files_with_content"`
}

// LoadProjectStructure 加载项目结构配置
func LoadProjectStructure() (*ProjectStructure, error) {
	var structure ProjectStructure
	if err := json.Unmarshal(projectStructureJSON, &structure); err != nil {
		return nil, fmt.Errorf("failed to parse project structure: %w", err)
	}
	return &structure, nil
}
