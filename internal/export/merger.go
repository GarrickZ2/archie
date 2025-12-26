package export

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/status"
)

// DocumentMerger merges all documents into single markdown file
type DocumentMerger struct {
	tocGenerator *TOCGenerator
	formatter    *Formatter
}

// NewDocumentMerger creates a new document merger
func NewDocumentMerger() *DocumentMerger {
	return &DocumentMerger{
		tocGenerator: NewTOCGenerator(),
		formatter:    NewFormatter(),
	}
}

// Merge merges all documents into final markdown content
func (m *DocumentMerger) Merge(config *ExportConfig, documents []*ExportedDocument) (string, error) {
	var finalContent strings.Builder

	// Step 1: Generate header with metadata
	header := m.formatter.GenerateMetadata(config, len(documents), len(config.IncludeFeatures))
	finalContent.WriteString(header)

	// Step 2: Build main content (without TOC first, we'll generate it later)
	var mainContent strings.Builder

	// Add status statistics if enabled and we have features
	if config.GenerateStats && len(config.IncludeFeatures) > 0 {
		statsContent := m.generateStatistics(config)
		if statsContent != "" {
			mainContent.WriteString(statsContent)
			mainContent.WriteString(m.formatter.CreateSeparator())
		}
	}

	// Add dependency graph if enabled and we have features
	if config.GenerateDepGraph && len(config.IncludeFeatures) > 0 {
		depGraphContent := m.generateDependencyGraph(config)
		if depGraphContent != "" {
			mainContent.WriteString(depGraphContent)
			mainContent.WriteString(m.formatter.CreateSeparator())
		}
	}

	// Step 3: Add root documents
	if len(config.IncludeRoot) > 0 {
		mainContent.WriteString("# Documentation\n\n")

		for _, doc := range documents {
			if doc.Type != DocTypeFeature && doc.Type != DocTypeWorkflow && doc.Type != DocTypeSpec {
				formatted := m.formatDocument(doc)
				mainContent.WriteString(formatted)
				mainContent.WriteString("\n\n")
			}
		}

		mainContent.WriteString(m.formatter.CreateSeparator())
	}

	// Step 4: Add features and their related documents
	if len(config.IncludeFeatures) > 0 {
		mainContent.WriteString("# Features\n\n")

		for _, featureKey := range config.IncludeFeatures {
			// Find feature document
			var featureDoc *ExportedDocument
			for _, doc := range documents {
				if doc.Type == DocTypeFeature && doc.Name == featureKey {
					featureDoc = doc
					break
				}
			}

			if featureDoc == nil {
				continue
			}

			// Add feature header
			mainContent.WriteString(fmt.Sprintf("## %s\n\n", featureKey))

			// Add feature content
			mainContent.WriteString(featureDoc.Content)
			mainContent.WriteString("\n")

			// Add workflow if exists
			if config.IncludeWorkflows {
				workflowName := featureKey + "-workflow"
				for _, doc := range documents {
					if doc.Type == DocTypeWorkflow && doc.Name == workflowName {
						mainContent.WriteString("### Workflow\n\n")
						mainContent.WriteString(doc.Content)
						mainContent.WriteString("\n")
						break
					}
				}
			}

			// Add spec if exists
			if config.IncludeSpecs {
				specName := featureKey + "-spec"
				for _, doc := range documents {
					if doc.Type == DocTypeSpec && doc.Name == specName {
						mainContent.WriteString("### Specification\n\n")
						// Adjust spec heading levels
						specContent := m.formatter.AdjustHeadingLevels(doc.Content, 1)
						mainContent.WriteString(specContent)
						mainContent.WriteString("\n")
						break
					}
				}
			}

			mainContent.WriteString("\n")
		}
	}

	// Step 5: Generate TOC if enabled
	fullContent := mainContent.String()
	if config.GenerateTOC {
		toc := m.tocGenerator.Generate(header + fullContent)
		if toc != "" {
			finalContent.WriteString(toc)
			finalContent.WriteString(m.formatter.CreateSeparator())
		}
	}

	// Add main content
	finalContent.WriteString(fullContent)

	// Sanitize final output
	return m.formatter.SanitizeMarkdown(finalContent.String()), nil
}

// formatDocument formats a single document
func (m *DocumentMerger) formatDocument(doc *ExportedDocument) string {
	var content strings.Builder

	// Add section title
	title := m.formatter.FormatSectionTitle(string(doc.Type))
	content.WriteString(fmt.Sprintf("## %s\n\n", title))

	// Add document content
	content.WriteString(doc.Content)

	return content.String()
}

// generateStatistics generates statistics content
func (m *DocumentMerger) generateStatistics(config *ExportConfig) string {
	// Load features to generate statistics
	parser := status.NewParser(afero.NewOsFs())
	features, err := parser.ParseFeaturesDir(config.ProjectPath)
	if err != nil || len(features) == 0 {
		return ""
	}

	// Filter features based on config
	var selectedFeatures []status.Feature
	for _, feature := range features {
		for _, selectedKey := range config.IncludeFeatures {
			if feature.Name == selectedKey {
				selectedFeatures = append(selectedFeatures, feature)
				break
			}
		}
	}

	if len(selectedFeatures) == 0 {
		return ""
	}

	statsGen := NewStatisticsGenerator(selectedFeatures)
	return statsGen.Generate()
}

// generateDependencyGraph generates dependency graph content
func (m *DocumentMerger) generateDependencyGraph(config *ExportConfig) string {
	// Load feature details to get dependencies
	detailParser := status.NewDetailParser(afero.NewOsFs())
	var featureDetails []*status.FeatureDetail

	for _, featureKey := range config.IncludeFeatures {
		detail, err := detailParser.ParseFeatureDetail(config.ProjectPath, featureKey)
		if err != nil {
			continue
		}
		featureDetails = append(featureDetails, detail)
	}

	if len(featureDetails) == 0 {
		return ""
	}

	depGraphGen := NewDependencyGraphGenerator(featureDetails)
	return depGraphGen.Generate()
}
