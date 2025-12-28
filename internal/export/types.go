package export

import "time"

// DocumentType represents types of documents in the project
type DocumentType string

const (
	DocTypeFeature  DocumentType = "feature"
	DocTypeWorkflow DocumentType = "workflow"
	DocTypeSpec     DocumentType = "spec"
)

// ExportConfig defines what to export
type ExportConfig struct {
	ProjectPath      string
	OutputPath       string
	IncludeRoot      []string // Root-level doc types to include
	IncludeFeatures  []string // Feature keys to include
	IncludeWorkflows bool     // Include workflow diagrams
	IncludeSpecs     bool     // Include spec files
	GenerateTOC      bool     // Generate table of contents
	GenerateDepGraph bool     // Generate dependency graph
	GenerateStats    bool     // Generate status statistics
}

// ExportedDocument represents a collected document
type ExportedDocument struct {
	Type     DocumentType
	Name     string
	Content  string
	Level    int // Heading level for TOC
	Children []*ExportedDocument
}

// ExportResult contains the result of export operation
type ExportResult struct {
	OutputPath    string
	DocumentCount int
	FeatureCount  int
	TotalSize     int64
	GeneratedAt   time.Time
	Warnings      []string
}

// CollectionWarning represents a warning during document collection
type CollectionWarning struct {
	Path   string
	Reason string
}
