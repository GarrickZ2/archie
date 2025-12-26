package export

import (
	"fmt"
	"regexp"
	"strings"
)

// TOCGenerator generates table of contents
type TOCGenerator struct{}

// TOCEntry represents an entry in the table of contents
type TOCEntry struct {
	Level   int
	Title   string
	Anchor  string
	Children []*TOCEntry
}

// NewTOCGenerator creates a new TOC generator
func NewTOCGenerator() *TOCGenerator {
	return &TOCGenerator{}
}

// Generate generates TOC from content
func (g *TOCGenerator) Generate(content string) string {
	entries := g.parseHeaders(content)

	if len(entries) == 0 {
		return ""
	}

	var toc strings.Builder
	toc.WriteString("## Table of Contents\n\n")
	toc.WriteString(g.formatTOC(entries, 0))

	return toc.String()
}

// parseHeaders parses markdown headers from content
func (g *TOCGenerator) parseHeaders(content string) []*TOCEntry {
	lines := strings.Split(content, "\n")
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	var entries []*TOCEntry

	for _, line := range lines {
		matches := headerRegex.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) == 3 {
			level := len(matches[1])
			title := strings.TrimSpace(matches[2])

			// Skip the main title (# Archie Project Export)
			if level == 1 {
				continue
			}

			entry := &TOCEntry{
				Level:  level,
				Title:  title,
				Anchor: g.generateAnchor(title),
			}

			entries = append(entries, entry)
		}
	}

	return entries
}

// generateAnchor generates anchor link from title
func (g *TOCGenerator) generateAnchor(title string) string {
	// Remove markdown formatting
	title = strings.ReplaceAll(title, "**", "")
	title = strings.ReplaceAll(title, "*", "")
	title = strings.ReplaceAll(title, "`", "")

	// Convert to lowercase
	anchor := strings.ToLower(title)

	// Replace spaces with hyphens
	anchor = strings.ReplaceAll(anchor, " ", "-")

	// Remove special characters except hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	anchor = reg.ReplaceAllString(anchor, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	anchor = reg.ReplaceAllString(anchor, "-")

	// Trim hyphens from start and end
	anchor = strings.Trim(anchor, "-")

	return anchor
}

// formatTOC formats TOC entries as markdown list
func (g *TOCGenerator) formatTOC(entries []*TOCEntry, baseIndent int) string {
	var content strings.Builder

	for _, entry := range entries {
		// Calculate indentation (level 2 = 0 indent, level 3 = 2 spaces, etc.)
		indent := (entry.Level - 2) * 2
		if indent < 0 {
			indent = 0
		}

		// Add indentation
		content.WriteString(strings.Repeat(" ", indent))

		// Add list item
		content.WriteString(fmt.Sprintf("- [%s](#%s)\n", entry.Title, entry.Anchor))
	}

	return content.String()
}
