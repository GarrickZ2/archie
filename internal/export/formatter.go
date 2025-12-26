package export

import (
	"fmt"
	"strings"
	"time"
)

// Formatter provides formatting utilities for export
type Formatter struct{}

// NewFormatter creates a new formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// AdjustHeadingLevels adjusts markdown heading levels
// Adds baseLevel to all existing headings
func (f *Formatter) AdjustHeadingLevels(content string, baseLevel int) string {
	if baseLevel <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if line is a heading
		if strings.HasPrefix(trimmed, "#") {
			// Count existing heading level
			level := 0
			for _, ch := range trimmed {
				if ch == '#' {
					level++
				} else {
					break
				}
			}

			// Get the title part (after the #'s)
			title := strings.TrimSpace(trimmed[level:])

			// Create new heading with adjusted level
			newLevel := baseLevel + level
			if newLevel > 6 {
				newLevel = 6 // Max heading level is 6
			}

			newLine := strings.Repeat("#", newLevel) + " " + title
			result = append(result, newLine)
		} else {
			// Not a heading, keep as is
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// SanitizeMarkdown cleans up markdown content
func (f *Formatter) SanitizeMarkdown(content string) string {
	// Remove extra blank lines (more than 2 consecutive)
	lines := strings.Split(content, "\n")
	var result []string
	blankCount := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			blankCount++
			if blankCount <= 2 {
				result = append(result, line)
			}
		} else {
			blankCount = 0
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// GenerateMetadata generates document header metadata
func (f *Formatter) GenerateMetadata(config *ExportConfig, docCount, featureCount int) string {
	var metadata strings.Builder

	metadata.WriteString("# Archie Project Export\n\n")

	// Generation info
	metadata.WriteString(fmt.Sprintf("**Generated:** %s\n\n",
		time.Now().Format("2006-01-02 15:04:05")))

	metadata.WriteString(fmt.Sprintf("**Project:** %s\n\n", config.ProjectPath))

	// Statistics
	metadata.WriteString(fmt.Sprintf("**Documents:** %d | **Features:** %d\n\n",
		docCount, featureCount))

	// Separator
	metadata.WriteString("---\n\n")

	return metadata.String()
}

// FormatSectionTitle formats a section title with proper casing
func (f *Formatter) FormatSectionTitle(title string) string {
	// Convert document types to nice titles
	titleMap := map[string]string{
		"background":   "Background",
		"dependency":   "Dependencies",
		"deployment":   "Deployment",
		"metrics":      "Metrics",
		"storage":      "Storage Design",
		"tasks":        "Tasks",
		"api":          "API Documentation",
		"faq":          "FAQ",
		"blocker":      "Blockers",
		"architecture": "Architecture",
		"feature":      "Feature",
		"workflow":     "Workflow",
		"spec":         "Specification",
	}

	if mapped, ok := titleMap[strings.ToLower(title)]; ok {
		return mapped
	}

	// Default: capitalize first letter of each word
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// CreateSeparator creates a visual separator
func (f *Formatter) CreateSeparator() string {
	return "\n---\n\n"
}
