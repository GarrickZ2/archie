package agent

import (
	"fmt"
	"log"

	"github.com/GarrickZ2/archie/resources"
)

// GetCommands returns all command templates as raw templates
// This function loads commands from the embedded FS
func GetCommands() map[string]*resources.CommandTemplate {
	commands, err := resources.LoadCommands()
	if err != nil {
		log.Printf("Warning: failed to load commands: %v", err)
		return make(map[string]*resources.CommandTemplate)
	}
	return commands
}

// GetFormattedCommands returns formatted commands for a specific agent
func GetFormattedCommands(agentName string) (map[string]string, error) {
	agent, err := Get(agentName)
	if err != nil {
		return nil, err
	}

	// Get agent config
	var fileFormat string
	var mapping map[string]string

	if customAgent, ok := agent.(*CustomAgent); ok {
		fileFormat = customAgent.config.FileFormat
		mapping = customAgent.config.Mapping
	} else {
		// For builtin agents, get from agents.json
		// This should have been loaded via LoadBuiltinAgents
		// Fallback to default
		fileFormat = "md"
		mapping = map[string]string{"content": "[CONTENT]"}
	}

	// Get formatter
	formatter, err := GetFormatter(fileFormat)
	if err != nil {
		return nil, err
	}

	// Load raw templates
	rawTemplates := GetCommands()
	formatted := make(map[string]string)

	// Format each template
	for filename, template := range rawTemplates {
		// Determine file extension based on format
		ext := ".md"
		if fileFormat == "toml" {
			ext = ".toml"
		}
		formattedFilename := filename + ext

		formattedContent, err := formatter.Format(template, mapping)
		if err != nil {
			return nil, fmt.Errorf("failed to format command %s: %w", filename, err)
		}

		formatted[formattedFilename] = formattedContent
	}

	return formatted, nil
}

// GetSubAgents returns all subagent templates
// This function loads subagents from the embedded FS
func GetSubAgents() map[string]string {
	subAgents, err := resources.LoadSubAgents()
	if err != nil {
		log.Printf("Warning: failed to load subagents: %v", err)
		return make(map[string]string)
	}
	return subAgents
}
