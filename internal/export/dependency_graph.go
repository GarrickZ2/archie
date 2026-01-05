package export

import (
	"fmt"
	"strings"

	"github.com/GarrickZ2/archie/internal/status"
)

// DependencyGraphGenerator generates dependency relationship graphs
type DependencyGraphGenerator struct {
	features []*status.FeatureDetail
}

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	Key          string
	Dependencies []string // Features that this feature depends on
}

// NewDependencyGraphGenerator creates a new dependency graph generator
func NewDependencyGraphGenerator(features []*status.FeatureDetail) *DependencyGraphGenerator {
	return &DependencyGraphGenerator{
		features: features,
	}
}

// Generate generates the Mermaid graph syntax
func (g *DependencyGraphGenerator) Generate() string {
	if len(g.features) == 0 {
		return ""
	}

	// Build dependency graph
	graph := g.buildGraph()

	if len(graph) == 0 {
		return ""
	}

	var content strings.Builder

	content.WriteString("## Dependency Graph\n\n")
	content.WriteString(g.formatMermaid(graph))

	// Add legend
	content.WriteString("\n**Legend:**\n\n")
	content.WriteString("- `→` indicates dependency (feature depends on target)\n")

	featureList := make([]string, 0, len(graph))
	for key := range graph {
		featureList = append(featureList, key)
	}
	if len(featureList) > 0 {
		content.WriteString(fmt.Sprintf("- **Features**: %s\n", strings.Join(featureList, ", ")))
	}

	content.WriteString("\n")

	return content.String()
}

// buildGraph builds the dependency node graph
func (g *DependencyGraphGenerator) buildGraph() map[string]*DependencyNode {
	graph := make(map[string]*DependencyNode)

	// Build nodes from features
	for _, feature := range g.features {
		node := &DependencyNode{
			Key:          feature.Key,
			Dependencies: []string{},
		}

		// Add feature dependencies
		for depKey := range feature.FeatureDependencies {
			if depKey != "" && depKey != "<feature-key>" {
				node.Dependencies = append(node.Dependencies, depKey)
			}
		}

		// Only add to graph if it has dependencies
		if len(node.Dependencies) > 0 {
			graph[feature.Key] = node
		}
	}

	return graph
}

// formatMermaid formats the graph as Mermaid syntax
func (g *DependencyGraphGenerator) formatMermaid(graph map[string]*DependencyNode) string {
	var content strings.Builder

	content.WriteString("```mermaid\n")
	content.WriteString("graph LR\n")

	// Add edges
	for key, node := range graph {
		// Format feature key for Mermaid (replace hyphens with underscores)
		featureID := strings.ReplaceAll(key, "-", "_")

		// Add dependencies: dependency -> this feature (arrow points to dependent)
		for _, dependency := range node.Dependencies {
			dependencyID := strings.ReplaceAll(dependency, "-", "_")
			content.WriteString(fmt.Sprintf("    %s[\"%s\"] --> %s[\"%s\"]\n",
				dependencyID, dependency, featureID, key))
		}
	}

	content.WriteString("```\n")

	return content.String()
}
