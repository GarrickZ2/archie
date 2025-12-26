package export

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/ui"
)

// ExportManager orchestrates the entire export workflow
type ExportManager struct {
	fs          afero.Fs
	projectPath string

	// Configuration from flags
	flagOutputPath string
	flagTOC        bool
	flagStats      bool
	flagDepGraph   bool

	// Components
	selector  *DocumentSelector
	collector *DocumentCollector
	merger    *DocumentMerger
}

// NewExportManager creates a new export manager
func NewExportManager(projectPath string, fs afero.Fs) *ExportManager {
	if fs == nil {
		fs = afero.NewOsFs()
	}

	return &ExportManager{
		fs:           fs,
		projectPath:  projectPath,
		flagTOC:      true,
		flagStats:    true,
		flagDepGraph: true,
		selector:     NewDocumentSelector(projectPath, fs),
		collector:    NewDocumentCollector(projectPath, fs),
		merger:       NewDocumentMerger(),
	}
}

// SetFlags sets command-line flags
func (m *ExportManager) SetFlags(outputPath string, toc, stats, depGraph bool) {
	m.flagOutputPath = outputPath
	m.flagTOC = toc
	m.flagStats = stats
	m.flagDepGraph = depGraph
}

// Export executes the export workflow
func (m *ExportManager) Export() (*ExportResult, error) {
	// Step 1: Validate project structure
	ui.ShowStep(1, 5, "Validating project...")
	if err := m.validateProject(); err != nil {
		return nil, fmt.Errorf("project validation failed: %w", err)
	}

	// Step 2: Get user selections via TUI
	ui.ShowStep(2, 5, "Selecting documents to export...")
	config, err := m.selector.SelectDocuments(m.flagOutputPath, m.flagTOC, m.flagStats, m.flagDepGraph)
	if err != nil {
		return nil, fmt.Errorf("selection failed: %w", err)
	}

	// Step 3: Collect documents
	ui.ShowStep(3, 5, "Collecting documents...")
	documents, warnings, err := m.collector.Collect(config)
	if err != nil {
		return nil, fmt.Errorf("collection failed: %w", err)
	}

	// Show warnings if any
	if len(warnings) > 0 {
		ui.ShowInfo(fmt.Sprintf("Found %d warnings during collection", len(warnings)))
	}

	// Step 4: Merge documents
	ui.ShowStep(4, 5, "Merging and formatting...")
	mergedContent, err := m.merger.Merge(config, documents)
	if err != nil {
		return nil, fmt.Errorf("merge failed: %w", err)
	}

	// Step 5: Write output file
	ui.ShowStep(5, 5, "Writing output file...")
	if err := afero.WriteFile(m.fs, config.OutputPath, []byte(mergedContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write output file: %w", err)
	}

	// Calculate statistics
	fileInfo, _ := m.fs.Stat(config.OutputPath)
	var totalSize int64
	if fileInfo != nil {
		totalSize = fileInfo.Size()
	}

	// Convert warnings to strings
	warningStrings := make([]string, len(warnings))
	for i, w := range warnings {
		warningStrings[i] = fmt.Sprintf("%s: %s", w.Path, w.Reason)
	}

	return &ExportResult{
		OutputPath:    config.OutputPath,
		DocumentCount: len(documents),
		FeatureCount:  len(config.IncludeFeatures),
		TotalSize:     totalSize,
		GeneratedAt:   time.Now(),
		Warnings:      warningStrings,
	}, nil
}

// validateProject checks if the current directory is a valid archie project
func (m *ExportManager) validateProject() error {
	// Check for at least one root document
	rootDocs := []string{"background.md", "dependency.md", "storage.md", "api", ".archie", "features"}
	hasAnyDoc := false

	for _, doc := range rootDocs {
		docPath := filepath.Join(m.projectPath, doc)
		exists, err := afero.Exists(m.fs, docPath)
		if err == nil && exists {
			hasAnyDoc = true
			break
		}
	}

	if !hasAnyDoc {
		return fmt.Errorf("not an archie project: no archie documents found in %s", m.projectPath)
	}

	return nil
}
