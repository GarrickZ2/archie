package agent

import (
	"fmt"
	"sync"
)

var (
	registry = &Registry{
		agents: make(map[string]Agent),
	}
)

// Registry agent 注册表
type Registry struct {
	mu     sync.RWMutex
	agents map[string]Agent
}

// Register 注册一个 agent
func Register(agent Agent) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.agents[agent.Name()] = agent
}

// Get 获取指定名称的 agent
func Get(name string) (Agent, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	agent, ok := registry.agents[name]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", name)
	}
	return agent, nil
}

// List 列出所有已注册的 agent
func List() []Agent {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	agents := make([]Agent, 0, len(registry.agents))
	for _, agent := range registry.agents {
		agents = append(agents, agent)
	}
	return agents
}

// Names 返回所有 agent 的名称
func Names() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	names := make([]string, 0, len(registry.agents))
	for name := range registry.agents {
		names = append(names, name)
	}
	return names
}
