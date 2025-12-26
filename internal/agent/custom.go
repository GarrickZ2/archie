package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// CustomAgentConfig 自定义 agent 的配置
type CustomAgentConfig struct {
	Name         string            `json:"name"`
	AgentDoc     string            `json:"background_doc"`
	CommandsDir  string            `json:"commands_dir"`
	SubAgentsDir string            `json:"sub_agents_dir"`
	Official     bool              `json:"official"` // 是否为官方 agent
	FileFormat   string            `json:"file_format"` // "md" 或 "toml"
	Mapping      map[string]string `json:"mapping"`    // YAML 字段到目标格式的映射
}

// CustomAgent 自定义 agent 实现
type CustomAgent struct {
	config CustomAgentConfig
}

// NewCustomAgent 创建自定义 agent
func NewCustomAgent(config CustomAgentConfig) *CustomAgent {
	// 设置默认值
	if config.FileFormat == "" {
		config.FileFormat = "md"
	}
	if config.Mapping == nil {
		config.Mapping = getDefaultMapping(config.FileFormat)
	}
	return &CustomAgent{config: config}
}

// getDefaultMapping 获取默认的 mapping 配置
func getDefaultMapping(fileFormat string) map[string]string {
	switch fileFormat {
	case "toml":
		return map[string]string{
			"description": "description",
			"content":     "prompt",
		}
	case "md":
		fallthrough
	default:
		return map[string]string{
			"content": "[CONTENT]",
		}
	}
}

// Name 返回 agent 名称
func (a *CustomAgent) Name() string {
	return a.config.Name
}

// PathConfig 返回路径配置
func (a *CustomAgent) PathConfig() PathConfig {
	return PathConfig{
		CommandsDir:  a.config.CommandsDir,
		SubAgentsDir: a.config.SubAgentsDir,
	}
}

// Commands 返回格式化后的命令提示
func (a *CustomAgent) Commands() map[string]string {
	// 获取格式化后的命令
	formatted, err := GetFormattedCommands(a.config.Name)
	if err != nil {
		// 如果格式化失败，返回空 map
		return make(map[string]string)
	}
	return formatted
}

// SubAgents 返回子 agent 提示
func (a *CustomAgent) SubAgents() map[string]string {
	if a.config.SubAgentsDir == "" {
		return nil
	}

	return GetSubAgents()
}

// SupportsSubAgents 返回是否支持子 agents
func (a *CustomAgent) SupportsSubAgents() bool {
	return a.config.SubAgentsDir != ""
}

// IsCustom 标识这是一个自定义 agent（非官方）
func (a *CustomAgent) IsCustom() bool {
	return !a.config.Official
}

// IsOfficial 标识这是一个官方 agent
func (a *CustomAgent) IsOfficial() bool {
	return a.config.Official
}

// CustomAgentStore 自定义 agent 存储
type CustomAgentStore struct {
	fs afero.Fs
}

// NewCustomAgentStore 创建存储实例
func NewCustomAgentStore(fs afero.Fs) *CustomAgentStore {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &CustomAgentStore{fs: fs}
}

// getGlobalConfigPath 获取全局配置文件路径（用户根目录）
func getGlobalConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".archie", "custom_agents.json"), nil
}

// Save 保存自定义 agent 配置到用户根目录
func (s *CustomAgentStore) Save(configs []CustomAgentConfig) error {
	configPath, err := getGlobalConfigPath()
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := s.fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 序列化配置
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := afero.WriteFile(s.fs, configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Load 加载自定义 agent 配置从用户根目录
func (s *CustomAgentStore) Load() ([]CustomAgentConfig, error) {
	configPath, err := getGlobalConfigPath()
	if err != nil {
		return nil, err
	}

	// 检查文件是否存在
	exists, err := afero.Exists(s.fs, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}
	if !exists {
		return []CustomAgentConfig{}, nil
	}

	// 读取文件
	data, err := afero.ReadFile(s.fs, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 反序列化
	var configs []CustomAgentConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return configs, nil
}

// Add 添加新的自定义 agent
func (s *CustomAgentStore) Add(config CustomAgentConfig) error {
	// 加载现有配置
	configs, err := s.Load()
	if err != nil {
		return err
	}

	// 检查是否已存在同名 agent
	for _, c := range configs {
		if c.Name == config.Name {
			return fmt.Errorf("custom agent '%s' already exists", config.Name)
		}
	}

	// 添加新配置
	configs = append(configs, config)

	// 保存
	return s.Save(configs)
}

// Remove 删除自定义 agent
func (s *CustomAgentStore) Remove(agentName string) error {
	// 加载现有配置
	configs, err := s.Load()
	if err != nil {
		return err
	}

	// 查找并删除
	found := false
	newConfigs := make([]CustomAgentConfig, 0, len(configs))
	for _, c := range configs {
		if c.Name == agentName {
			found = true
			continue
		}
		newConfigs = append(newConfigs, c)
	}

	if !found {
		return fmt.Errorf("custom agent '%s' not found", agentName)
	}

	// 保存
	return s.Save(newConfigs)
}

// LoadAndRegister 加载并注册所有自定义 agents（从用户根目录）
func LoadAndRegister() error {
	store := NewCustomAgentStore(nil)
	configs, err := store.Load()
	if err != nil {
		return err
	}

	for _, config := range configs {
		agent := NewCustomAgent(config)
		Register(agent)
	}

	return nil
}

// GetConfigFilePath 获取配置文件路径（用于检查）
func GetConfigFilePath() (string, error) {
	return getGlobalConfigPath()
}

// ValidateCustomAgentConfig 验证自定义 agent 配置
func ValidateCustomAgentConfig(config CustomAgentConfig) error {
	if config.Name == "" {
		return fmt.Errorf("agent name cannot be empty")
	}
	if config.AgentDoc == "" {
		config.AgentDoc = "AGENTS.md" // 设置默认值
	}
	if config.CommandsDir == "" {
		return fmt.Errorf("commands directory cannot be empty")
	}
	// 验证 file_format
	if config.FileFormat != "" && config.FileFormat != "md" && config.FileFormat != "toml" {
		return fmt.Errorf("file_format must be 'md' or 'toml'")
	}
	// 验证 mapping 必须包含正文字段的映射
	if config.Mapping != nil {
		hasContent := false
		for _, target := range config.Mapping {
			if target == "[CONTENT]" || target == "prompt" {
				hasContent = true
				break
			}
		}
		if !hasContent {
			return fmt.Errorf("mapping must include a content field mapping ([CONTENT] for md or 'prompt' for toml)")
		}
	}
	// SubAgentsDir 可选，所以不检查
	return nil
}
