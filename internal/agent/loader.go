package agent

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed agents.json
var builtinAgentsJSON []byte

// LoadBuiltinAgents 加载内置 agents
func LoadBuiltinAgents() error {
	var configs []CustomAgentConfig
	if err := json.Unmarshal(builtinAgentsJSON, &configs); err != nil {
		return fmt.Errorf("failed to parse builtin agents: %w", err)
	}

	for _, config := range configs {
		// 确保默认值已设置
		if config.FileFormat == "" {
			config.FileFormat = "md"
		}
		if config.Mapping == nil {
			config.Mapping = getDefaultMapping(config.FileFormat)
		}
		agent := NewCustomAgent(config)
		Register(agent)
	}

	return nil
}

func init() {
	// 自动加载内置 agents
	if err := LoadBuiltinAgents(); err != nil {
		panic(fmt.Sprintf("failed to load builtin agents: %v", err))
	}
}
