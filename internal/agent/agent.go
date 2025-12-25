package agent

import "context"

// Agent 定义 Code Agent 的抽象接口
// 支持不同的 AI 平台（Claude Code, Cursor, etc.）
type Agent interface {
	// Name 返回 agent 的名称
	Name() string

	// PathConfig 返回该 agent 的路径配置
	PathConfig() PathConfig

	// Commands 返回所有格式化后的 command prompts
	// key: 文件名（如 "design.md" 或 "design.toml"）
	// value: 格式化后的文件内容
	Commands() map[string]string

	// SubAgents 返回所有 sub-agent prompts（如果支持）
	// key: 文件名（如 "api-designer.md"）
	// value: 文件内容
	SubAgents() map[string]string

	// SupportsSubAgents 返回是否支持 sub-agents
	SupportsSubAgents() bool
}

// PathConfig 定义 agent 的安装路径配置
type PathConfig struct {
	// CommandsDir command prompts 的安装目录
	// 例如：".claude/commands/archie"
	CommandsDir string

	// SubAgentsDir sub-agent prompts 的安装目录（如果支持）
	// 例如：".claude/agents/archie"
	SubAgentsDir string
}

// SetupConfig 定义 setup 配置
type SetupConfig struct {
	ProjectPath string       // 项目路径
	AgentType   string       // agent 类型
	Options     SetupOptions // 可选配置
}

// SetupOptions 可选的 setup 配置
type SetupOptions struct {
	IncludeExamples bool                   // 是否包含示例
	EnabledCommands []string               // 启用的命令（空则全部启用）
	CustomConfig    map[string]interface{} // 自定义配置
}

// Setupper 定义 agent setup 的接口
type Setupper interface {
	// Setup 在项目中 setup agent 配置
	Setup(ctx context.Context, config SetupConfig) error
}
