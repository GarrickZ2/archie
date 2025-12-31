package agent

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/afero"
)

// AgentState 记录 agent 的初始化状态
type AgentState struct {
	AgentName     string    `json:"agent_name"`
	IsCustom      bool      `json:"is_custom"`
	InitializedAt time.Time `json:"initialized_at"`
}

// ProjectState 项目状态
type ProjectState struct {
	InitializedAgents []AgentState `json:"initialized_agents"`
}

// StateManager 状态管理器
type StateManager struct {
	fs afero.Fs
}

// NewStateManager 创建状态管理器
func NewStateManager(fs afero.Fs) *StateManager {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &StateManager{fs: fs}
}

// Load 加载项目状态
func (m *StateManager) Load(projectPath string) (*ProjectState, error) {
	statePath := m.getStatePath(projectPath)

	// 检查文件是否存在
	exists, err := afero.Exists(m.fs, statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to check state file: %w", err)
	}
	if !exists {
		return &ProjectState{InitializedAgents: []AgentState{}}, nil
	}

	// 读取文件
	data, err := afero.ReadFile(m.fs, statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	// 反序列化
	var state ProjectState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &state, nil
}

// Save 保存项目状态
func (m *StateManager) Save(projectPath string, state *ProjectState) error {
	statePath := m.getStatePath(projectPath)

	// 确保目录存在
	dir := filepath.Dir(statePath)
	if err := m.fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 序列化
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// 写入文件
	if err := afero.WriteFile(m.fs, statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// MarkInitialized 标记 agent 为已初始化
func (m *StateManager) MarkInitialized(projectPath string, agentName string, isCustom bool) error {
	state, err := m.Load(projectPath)
	if err != nil {
		return err
	}

	// 检查是否已存在
	for i, a := range state.InitializedAgents {
		if a.AgentName == agentName {
			// 更新时间戳
			state.InitializedAgents[i].InitializedAt = time.Now()
			return m.Save(projectPath, state)
		}
	}

	// 添加新记录
	state.InitializedAgents = append(state.InitializedAgents, AgentState{
		AgentName:     agentName,
		IsCustom:      isCustom,
		InitializedAt: time.Now(),
	})

	return m.Save(projectPath, state)
}

// IsInitialized 检查 agent 是否已初始化
func (m *StateManager) IsInitialized(projectPath string, agentName string) (bool, error) {
	state, err := m.Load(projectPath)
	if err != nil {
		return false, err
	}

	for _, a := range state.InitializedAgents {
		if a.AgentName == agentName {
			return true, nil
		}
	}

	return false, nil
}

// GetInitializedAgents 获取所有已初始化的 agents
func (m *StateManager) GetInitializedAgents(projectPath string) ([]AgentState, error) {
	state, err := m.Load(projectPath)
	if err != nil {
		return nil, err
	}

	return state.InitializedAgents, nil
}

// getStatePath 获取状态文件路径
func (m *StateManager) getStatePath(projectPath string) string {
	return filepath.Join(projectPath, ".archie", "state.json")
}

// AgentInfo agent 信息（用于 TUI 显示）
type AgentInfo struct {
	Name          string
	IsCustom      bool
	IsOfficial    bool
	IsInitialized bool
	DisplayName   string
}

// GetAllAgents 获取所有可用的 agents（包括已注册和自定义）
func GetAllAgents(projectPath string) ([]AgentInfo, error) {
	stateManager := NewStateManager(nil)
	state, err := stateManager.Load(projectPath)
	if err != nil {
		return nil, err
	}

	// 创建已初始化 agents 的映射
	initializedMap := make(map[string]bool)
	for _, a := range state.InitializedAgents {
		initializedMap[a.AgentName] = true
	}

	// 获取所有已注册的 agents
	var infos []AgentInfo
	for _, agent := range List() {
		name := agent.Name()
		isCustom := false
		isOfficial := false

		// 检查是否为 CustomAgent（可能是官方或自定义）
		if ca, ok := agent.(*CustomAgent); ok {
			isCustom = ca.IsCustom()
			isOfficial = ca.IsOfficial()
		}

		// 构建显示名称
		displayName := name
		if isCustom {
			displayName = fmt.Sprintf("%s (custom)", name)
		}

		infos = append(infos, AgentInfo{
			Name:          name,
			IsCustom:      isCustom,
			IsOfficial:    isOfficial,
			IsInitialized: initializedMap[name],
			DisplayName:   displayName,
		})
	}

	// 按名字首字母排序
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return infos, nil
}
