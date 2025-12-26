package export

import (
	"fmt"
	"strings"

	"github.com/GarrickZ2/archie/internal/status"
)

// StatisticsGenerator generates status statistics for features
type StatisticsGenerator struct {
	features []status.Feature
	summary  *status.Summary
}

// NewStatisticsGenerator creates a new statistics generator
func NewStatisticsGenerator(features []status.Feature) *StatisticsGenerator {
	aggregator := status.NewAggregator(features)
	summary := aggregator.Aggregate()

	return &StatisticsGenerator{
		features: features,
		summary:  summary,
	}
}

// Generate generates the complete statistics markdown
func (g *StatisticsGenerator) Generate() string {
	if len(g.features) == 0 {
		return ""
	}

	var content strings.Builder

	content.WriteString("## Status Statistics\n\n")

	// Overall Progress
	content.WriteString("### Overall Progress\n\n")
	content.WriteString(g.formatProgressBar())
	content.WriteString("\n\n")

	// Status Distribution Table
	content.WriteString("### Status Distribution\n\n")
	content.WriteString(g.formatStatusTable())
	content.WriteString("\n")

	// Key Insights
	insights := g.summary.GetTopInsights()
	if len(insights) > 0 {
		content.WriteString("### Key Insights\n\n")
		content.WriteString(g.formatInsights(insights))
		content.WriteString("\n")
	}

	return content.String()
}

// formatProgressBar generates an ASCII progress bar
func (g *StatisticsGenerator) formatProgressBar() string {
	progress := g.summary.OverallProgress
	barWidth := 50
	filled := (progress * barWidth) / 100

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	return fmt.Sprintf("%s **%d%% Complete**\n", bar, progress)
}

// formatStatusTable generates the status distribution table
func (g *StatisticsGenerator) formatStatusTable() string {
	var table strings.Builder

	table.WriteString("| Status            | Count | Percentage |\n")
	table.WriteString("|-------------------|-------|------------|\n")

	totalCount := g.summary.TotalFeatures
	if totalCount == 0 {
		totalCount = 1 // Avoid division by zero
	}

	// Show all statuses that have at least one feature
	for _, st := range status.AllStatuses {
		count := g.summary.StatusCounts[st]
		if count > 0 {
			percentage := (count * 100) / totalCount
			statusName := strings.ReplaceAll(string(st), "_", " ")
			table.WriteString(fmt.Sprintf("| %-17s | %-5d | %-10s |\n",
				statusName, count, fmt.Sprintf("%d%%", percentage)))
		}
	}

	// Total row
	table.WriteString(fmt.Sprintf("| **%-15s** | **%-3d** | **%-8s** |\n",
		"Total", g.summary.TotalFeatures, "100%"))

	return table.String()
}

// formatInsights formats the key insights
func (g *StatisticsGenerator) formatInsights(insights []string) string {
	var content strings.Builder

	for _, insight := range insights {
		content.WriteString(fmt.Sprintf("- %s\n", insight))
	}

	return content.String()
}
