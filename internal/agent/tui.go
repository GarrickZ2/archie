package agent

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/GarrickZ2/archie/internal/ui"
)

const (
	addNewCustomAgentOption = "[+ Add new custom agent]"
)

// TUISelector TUI 选择器
type TUISelector struct {
	projectPath string
}

// NewTUISelector 创建 TUI 选择器
func NewTUISelector(projectPath string) *TUISelector {
	return &TUISelector{projectPath: projectPath}
}

// SelectAgent 选择 agent（带 TUI 界面）
func (s *TUISelector) SelectAgent() (selectedName string, isCustom bool, err error) {
	// 加载已有的自定义 agents（从用户根目录）
	if err := LoadAndRegister(); err != nil {
		// 如果加载失败（比如文件不存在），继续执行
		// 但不是错误，因为可能还没有自定义 agents
	}

	// 获取所有可用的 agents
	agents, err := GetAllAgents(s.projectPath)
	if err != nil {
		return "", false, fmt.Errorf("failed to get agents: %w", err)
	}

	// 构建选项列表
	options := make([]string, 0, len(agents)+1)
	optionToAgent := make(map[string]AgentInfo)

	for _, agent := range agents {
		// 添加勾选标记
		prefix := "  "
		if agent.IsInitialized {
			prefix = "✓ "
		}

		option := prefix + agent.DisplayName
		options = append(options, option)
		optionToAgent[option] = agent
	}

	// 添加 "Add new custom agent" 选项
	options = append(options, addNewCustomAgentOption)

	// Note: Welcome banner is now shown by the init command
	// Just show the agent selection prompt

	// 使用 survey 进行选择
	var selected string
	prompt := &survey.Select{
		Message: "Select Agent:",
		Options: options,
		Default: options[0],
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", false, fmt.Errorf("selection cancelled")
	}

	// 处理 "Add new custom agent" 选项
	if selected == addNewCustomAgentOption {
		config, err := s.promptCustomAgentConfig()
		if err != nil {
			return "", false, err
		}

		// 保存自定义 agent 配置到用户根目录
		store := NewCustomAgentStore(nil)
		if err := store.Add(config); err != nil {
			return "", false, fmt.Errorf("failed to save custom agent: %w", err)
		}

		// 注册新的自定义 agent
		customAgent := NewCustomAgent(config)
		Register(customAgent)

		return config.Name, true, nil
	}

	// 返回选中的 agent
	agentInfo := optionToAgent[selected]
	return agentInfo.Name, agentInfo.IsCustom, nil
}

// promptCustomAgentConfig 提示用户输入自定义 agent 配置
func (s *TUISelector) promptCustomAgentConfig() (CustomAgentConfig, error) {
	var config CustomAgentConfig

	// 1. Agent Name
	namePrompt := &survey.Input{
		Message: "Agent name:",
		Help:    "The name of your custom agent (e.g., 'cursor', 'my-agent')",
	}
	if err := survey.AskOne(namePrompt, &config.Name, survey.WithValidator(survey.Required)); err != nil {
		return config, err
	}

	// 2. Background Doc Name (默认 AGENTS.md)
	backgroundDocPrompt := &survey.Input{
		Message: "Agent doc name (press Enter for default: AGENTS.md):",
		Help:    "The documentation file that provides context to the agent",
		Default: "AGENTS.md",
	}
	if err := survey.AskOne(backgroundDocPrompt, &config.AgentDoc); err != nil {
		return config, err
	}
	// 如果用户没有输入，使用默认值
	if config.AgentDoc == "" {
		config.AgentDoc = "AGENTS.md"
	}

	// 3. Commands Directory
	commandsDirPrompt := &survey.Input{
		Message: "Commands folder:",
		Help:    "Directory where command prompts will be stored (e.g., '.cursor/rules/archie')",
	}
	if err := survey.AskOne(commandsDirPrompt, &config.CommandsDir, survey.WithValidator(survey.Required)); err != nil {
		return config, err
	}

	// 4. Sub-agents Directory (optional)
	subAgentsDirPrompt := &survey.Input{
		Message: "Sub-agents folder (optional, press Enter to skip):",
		Help:    "Directory for sub-agent prompts. Leave empty if not supported.",
	}
	if err := survey.AskOne(subAgentsDirPrompt, &config.SubAgentsDir); err != nil {
		return config, err
	}

	// 5. File Format
	fileFormatPrompt := &survey.Select{
		Message: "File format:",
		Help:    "Output format for command files (md = Markdown with frontmatter, toml = TOML)",
		Options: []string{"md", "toml"},
		Default: "md",
	}
	if err := survey.AskOne(fileFormatPrompt, &config.FileFormat); err != nil {
		return config, err
	}

	// 6. Set default mapping based on file format
	config.Mapping = getDefaultMapping(config.FileFormat)

	// 验证配置
	if err := ValidateCustomAgentConfig(config); err != nil {
		return config, err
	}

	return config, nil
}

// ConfirmReconfigure 确认是否重新配置
func (s *TUISelector) ConfirmReconfigure(agentName string) (bool, error) {
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Agent '%s' is already initialized. Do you want to re-configure it?", agentName),
		Default: false,
	}

	if err := survey.AskOne(prompt, &confirm); err != nil {
		return false, err
	}

	return confirm, nil
}

// CustomAgentManager 管理 custom agents
type CustomAgentManager struct {
	store *CustomAgentStore
}

// NewCustomAgentManager 创建管理器
func NewCustomAgentManager() *CustomAgentManager {
	return &CustomAgentManager{
		store: NewCustomAgentStore(nil),
	}
}

// ShowManagementUI 显示管理界面
func (m *CustomAgentManager) ShowManagementUI() error {
	for {
		// 加载所有 custom agents
		configs, err := m.store.Load()
		if err != nil {
			return fmt.Errorf("failed to load custom agents: %w", err)
		}

		// 构建选项列表
		options := []string{
			"[+ Add new custom agent]",
		}

		// 添加现有的 custom agents（仅非官方的）
		customConfigs := make([]CustomAgentConfig, 0)
		for _, config := range configs {
			if !config.Official {
				customConfigs = append(customConfigs, config)
				options = append(options, fmt.Sprintf("[-] Remove: %s", config.Name))
			}
		}

		options = append(options, "[x] Exit")

		// 显示选择界面
		var selected string
		prompt := &survey.Select{
			Message: "Manage Custom Agents:",
			Options: options,
		}

		if err := survey.AskOne(prompt, &selected); err != nil {
			return fmt.Errorf("selection cancelled")
		}

		// 处理选择
		switch selected {
		case "[+ Add new custom agent]":
			if err := m.addCustomAgent(); err != nil {
				ui.ShowError(fmt.Sprintf("Failed to add custom agent: %v", err))
				fmt.Println()
				continue
			}
			ui.ShowSuccess("Custom agent added successfully!")
			fmt.Println()

		case "[x] Exit":
			return nil

		default:
			// 删除 custom agent
			if err := m.removeCustomAgent(selected, customConfigs); err != nil {
				ui.ShowError(fmt.Sprintf("Failed to remove custom agent: %v", err))
				fmt.Println()
				continue
			}
			ui.ShowSuccess("Custom agent removed successfully!")
			fmt.Println()
		}
	}
}

// addCustomAgent 添加新的 custom agent
func (m *CustomAgentManager) addCustomAgent() error {
	var config CustomAgentConfig

	// 1. Agent Name
	namePrompt := &survey.Input{
		Message: "Agent name:",
		Help:    "The name of your custom agent (e.g., 'cursor', 'my-agent')",
	}
	if err := survey.AskOne(namePrompt, &config.Name, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// 2. Background Doc Name (默认 AGENTS.md)
	backgroundDocPrompt := &survey.Input{
		Message: "Agent doc name (press Enter for default: AGENTS.md):",
		Help:    "The documentation file that provides context to the agent",
		Default: "AGENTS.md",
	}
	if err := survey.AskOne(backgroundDocPrompt, &config.AgentDoc); err != nil {
		return err
	}
	if config.AgentDoc == "" {
		config.AgentDoc = "AGENTS.md"
	}

	// 3. Commands Directory
	commandsDirPrompt := &survey.Input{
		Message: "Commands folder:",
		Help:    "Directory where command prompts will be stored (e.g., '.cursor/rules/archie')",
	}
	if err := survey.AskOne(commandsDirPrompt, &config.CommandsDir, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// 4. Sub-agents Directory (optional)
	subAgentsDirPrompt := &survey.Input{
		Message: "Sub-agents folder (optional, press Enter to skip):",
		Help:    "Directory for sub-agent prompts. Leave empty if not supported.",
	}
	if err := survey.AskOne(subAgentsDirPrompt, &config.SubAgentsDir); err != nil {
		return err
	}

	// 5. File Format
	fileFormatPrompt := &survey.Select{
		Message: "File format:",
		Help:    "Output format for command files (md = Markdown with frontmatter, toml = TOML)",
		Options: []string{"md", "toml"},
		Default: "md",
	}
	if err := survey.AskOne(fileFormatPrompt, &config.FileFormat); err != nil {
		return err
	}

	// 6. Set default mapping based on file format
	config.Mapping = getDefaultMapping(config.FileFormat)

	// 验证配置
	if err := ValidateCustomAgentConfig(config); err != nil {
		return err
	}

	// 保存
	return m.store.Add(config)
}

// removeCustomAgent 删除 custom agent
func (m *CustomAgentManager) removeCustomAgent(selected string, configs []CustomAgentConfig) error {
	// 从选项中提取 agent 名称
	// 格式: "[-] Remove: agent-name"
	var agentName string
	for _, config := range configs {
		if selected == fmt.Sprintf("[-] Remove: %s", config.Name) {
			agentName = config.Name
			break
		}
	}

	if agentName == "" {
		return fmt.Errorf("invalid selection")
	}

	// 确认删除
	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to remove custom agent '%s'?", agentName),
		Default: false,
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("deletion cancelled")
	}

	// 删除
	return m.store.Remove(agentName)
}
