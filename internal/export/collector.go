package export

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/status"
)

// DocumentCollector collects document content from filesystem
type DocumentCollector struct {
	fs           afero.Fs
	projectPath  string
	parser       *status.Parser
	detailParser *status.DetailParser
}

// NewDocumentCollector creates a new document collector
func NewDocumentCollector(projectPath string, fs afero.Fs) *DocumentCollector {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	return &DocumentCollector{
		fs:           fs,
		projectPath:  projectPath,
		parser:       status.NewParser(fs),
		detailParser: status.NewDetailParser(fs),
	}
}

// Collect collects all documents based on configuration
// Returns documents, warnings, and error
func (c *DocumentCollector) Collect(config *ExportConfig) ([]*ExportedDocument, []CollectionWarning, error) {
	var documents []*ExportedDocument
	var warnings []CollectionWarning

	// Collect root documents
	for _, rootDoc := range config.IncludeRoot {
		doc, err := c.collectRootDoc(rootDoc)
		if err != nil {
			warnings = append(warnings, CollectionWarning{
				Path:   rootDoc,
				Reason: err.Error(),
			})
			continue
		}
		documents = append(documents, doc)
	}

	// Collect features and their related documents
	for _, featureKey := range config.IncludeFeatures {
		// Collect feature document
		featureDoc, err := c.collectFeature(featureKey)
		if err != nil {
			warnings = append(warnings, CollectionWarning{
				Path:   fmt.Sprintf("features/%s.md", featureKey),
				Reason: err.Error(),
			})
			continue
		}
		documents = append(documents, featureDoc)

		// Collect workflow if requested
		if config.IncludeWorkflows {
			workflowDoc, err := c.collectWorkflow(featureKey)
			if err != nil {
				warnings = append(warnings, CollectionWarning{
					Path:   fmt.Sprintf("workflow/%s", featureKey),
					Reason: err.Error(),
				})
			} else if workflowDoc != nil {
				documents = append(documents, workflowDoc)
			}
		}

		// Collect spec if requested
		if config.IncludeSpecs {
			specDoc, err := c.collectSpec(featureKey)
			if err != nil {
				warnings = append(warnings, CollectionWarning{
					Path:   fmt.Sprintf("spec/%s.spec.md", featureKey),
					Reason: err.Error(),
				})
			} else if specDoc != nil {
				documents = append(documents, specDoc)
			}
		}
	}

	return documents, warnings, nil
}

// collectRootDoc collects a root-level document
func (c *DocumentCollector) collectRootDoc(name string) (*ExportedDocument, error) {
	filePath := filepath.Join(c.projectPath, name)

	// Check if it's a directory (like api/)
	isDir, err := afero.IsDir(c.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to check path: %w", err)
	}

	if isDir {
		// For directories, try to read a main file (like api.md in api/)
		mainFile := filepath.Join(filePath, strings.TrimSuffix(name, "/")+".md")
		content, err := afero.ReadFile(c.fs, mainFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory main file: %w", err)
		}

		return &ExportedDocument{
			Type:    DocumentType(strings.TrimSuffix(name, "/")),
			Name:    name,
			Content: string(content),
			Level:   2, // ##
		}, nil
	}

	// Read file content
	content, err := afero.ReadFile(c.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Determine document type
	docType := DocumentType(strings.TrimSuffix(name, ".md"))

	return &ExportedDocument{
		Type:    docType,
		Name:    name,
		Content: string(content),
		Level:   2, // ##
	}, nil
}

// collectFeature collects a feature document with all its sections
func (c *DocumentCollector) collectFeature(featureKey string) (*ExportedDocument, error) {
	// Parse feature detail
	detail, err := c.detailParser.ParseFeatureDetail(c.projectPath, featureKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feature: %w", err)
	}

	// Format feature content as structured markdown
	var content strings.Builder

	// Status section
	content.WriteString("### Status\n\n")
	statusStr := strings.ReplaceAll(string(detail.Status), "_", " ")
	content.WriteString(fmt.Sprintf("- **Current Status**: %s\n", statusStr))
	if detail.Owner != "" {
		content.WriteString(fmt.Sprintf("- **Owner**: %s\n", detail.Owner))
	}
	if detail.LastUpdated != "" && detail.LastUpdated != "YYYY-MM-DD" {
		content.WriteString(fmt.Sprintf("- **Last Updated**: %s\n", detail.LastUpdated))
	}
	if detail.Reason != "" {
		content.WriteString(fmt.Sprintf("- **Reason**: %s\n", detail.Reason))
	}
	content.WriteString("\n")

	// Summary section
	if detail.OneLiner != "" || detail.Background != "" || detail.UserStory != "" {
		content.WriteString("### Summary\n\n")
		if detail.OneLiner != "" {
			content.WriteString(fmt.Sprintf("**One-liner**: %s\n\n", detail.OneLiner))
		}
		if detail.Background != "" {
			content.WriteString(fmt.Sprintf("**Background**: %s\n\n", detail.Background))
		}
		if detail.UserStory != "" {
			content.WriteString(fmt.Sprintf("**User Story**: %s\n\n", detail.UserStory))
		}
	}

	// Scope section
	if len(detail.InScope) > 0 || len(detail.OutScope) > 0 {
		content.WriteString("### Scope\n\n")
		if len(detail.InScope) > 0 {
			content.WriteString("**In Scope:**\n\n")
			for _, item := range detail.InScope {
				if item != "" && item != "..." {
					content.WriteString(fmt.Sprintf("- %s\n", item))
				}
			}
			content.WriteString("\n")
		}
		if len(detail.OutScope) > 0 {
			content.WriteString("**Out of Scope:**\n\n")
			for _, item := range detail.OutScope {
				if item != "" && item != "..." {
					content.WriteString(fmt.Sprintf("- %s\n", item))
				}
			}
			content.WriteString("\n")
		}
	}

	// Requirements
	if len(detail.Requirements) > 0 {
		content.WriteString("### Requirements\n\n")
		for _, req := range detail.Requirements {
			if req != "" && !strings.HasPrefix(req, "R1:") && !strings.HasPrefix(req, "R2:") {
				content.WriteString(fmt.Sprintf("- %s\n", req))
			}
		}
		content.WriteString("\n")
	}

	// Acceptance Criteria
	if len(detail.AcceptanceCriteria) > 0 {
		content.WriteString("### Acceptance Criteria\n\n")
		for _, ac := range detail.AcceptanceCriteria {
			if ac != "" && !strings.HasPrefix(ac, "AC1:") && !strings.HasPrefix(ac, "AC2:") {
				content.WriteString(fmt.Sprintf("- %s\n", ac))
			}
		}
		content.WriteString("\n")
	}

	// Dependencies
	if len(detail.Upstreams) > 0 || len(detail.Downstreams) > 0 {
		content.WriteString("### Dependencies\n\n")
		if len(detail.Upstreams) > 0 {
			content.WriteString("**Upstreams:**\n\n")
			for name, reason := range detail.Upstreams {
				if name != "" && name != "<dependency-name>" {
					content.WriteString(fmt.Sprintf("- **%s**: %s\n", name, reason))
				}
			}
			content.WriteString("\n")
		}
		if len(detail.Downstreams) > 0 {
			content.WriteString("**Downstreams:**\n\n")
			for name, reason := range detail.Downstreams {
				if name != "" && name != "<dependency-name>" {
					content.WriteString(fmt.Sprintf("- **%s**: %s\n", name, reason))
				}
			}
			content.WriteString("\n")
		}
	}

	return &ExportedDocument{
		Type:    DocTypeFeature,
		Name:    featureKey,
		Content: content.String(),
		Level:   2, // ##
	}, nil
}

// collectWorkflow collects workflow documentation for a feature
func (c *DocumentCollector) collectWorkflow(featureKey string) (*ExportedDocument, error) {
	workflowDir := filepath.Join(c.projectPath, "workflow", featureKey)

	// Check if workflow directory exists
	exists, err := afero.DirExists(c.fs, workflowDir)
	if err != nil || !exists {
		return nil, fmt.Errorf("workflow directory not found")
	}

	var content strings.Builder

	// Try to read workflow.md
	workflowFile := filepath.Join(workflowDir, "workflow.md")
	workflowContent, err := afero.ReadFile(c.fs, workflowFile)
	if err == nil {
		content.WriteString(string(workflowContent))
		content.WriteString("\n\n")
	}

	// Collect all .mmd files
	files, err := afero.ReadDir(c.fs, workflowDir)
	if err == nil {
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".mmd" {
				mmdPath := filepath.Join(workflowDir, file.Name())
				mmdContent, err := afero.ReadFile(c.fs, mmdPath)
				if err == nil {
					content.WriteString(fmt.Sprintf("#### Diagram: %s\n\n", file.Name()))
					content.WriteString("```mermaid\n")
					content.WriteString(string(mmdContent))
					content.WriteString("\n```\n\n")
				}
			}
		}
	}

	if content.Len() == 0 {
		return nil, fmt.Errorf("no workflow content found")
	}

	return &ExportedDocument{
		Type:    DocTypeWorkflow,
		Name:    featureKey + "-workflow",
		Content: content.String(),
		Level:   3, // ###
	}, nil
}

// collectSpec collects specification document for a feature
func (c *DocumentCollector) collectSpec(featureKey string) (*ExportedDocument, error) {
	specPath := filepath.Join(c.projectPath, "spec", featureKey+".spec.md")

	// Check if spec file exists
	exists, err := afero.Exists(c.fs, specPath)
	if err != nil || !exists {
		return nil, fmt.Errorf("spec file not found")
	}

	// Read spec content
	content, err := afero.ReadFile(c.fs, specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	return &ExportedDocument{
		Type:    DocTypeSpec,
		Name:    featureKey + "-spec",
		Content: string(content),
		Level:   3, // ###
	}, nil
}
