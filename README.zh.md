# Archie（中文）

**Archie** 是一个以文档为核心的 AI CLI 工具，帮助工程师把模糊想法系统性地推进为**清晰、可执行的技术设计和 Spec**，并在过程中持续管理项目进度。

Archie 不是写代码的工具。
它是一个 **面向技术设计、规范沉淀和执行编排的系统**。

---

## 为什么需要 Archie？

很多项目在真正开始写代码前就已经失败了：

- 需求不清晰
- 设计分散在各处
- 决策只存在于聊天中
- 项目进度不可见、不可控

Archie 通过以下原则解决这些问题：

- 📄 **文档即事实源**
- 🔄 **严格的 Feature 状态机**
- 🧠 **AI 辅助的 Review / Design / Spec**
- 📦 **所有内容本地 Markdown 化**
- 🧩 **可组合的子 Agent（API / Workflow / Storage / Metrics / Tasks / Test Plans）**

---

## 核心思想

Archie 把 **Feature** 作为最小原子，并强制最小设计链路：

> **Feature → Workflow → Spec**

其他内容（API、存储、指标、任务、测试计划、发布）都是对这条主链路的支撑。

---

## 安装

```bash
go install github.com/GarrickZ2/archie@latest
```

### 验证安装

```bash
archie --help
```

---

## 如何使用 Archie

Archie 有**两种使用模式**：

### 模式 1: CLI 命令（终端）

在终端直接运行的工作空间管理命令：

| 命令 | 描述 |
|------|------|
| `archie init` | 初始化工作空间结构并安装 agent 命令 |
| `archie setup` | 交互式 TUI 编辑背景和管理 features |
| `archie status` | 显示项目状态和交互式 feature 浏览器 |
| `archie export` | 导出文档到单个 markdown 文件 |

### 模式 2: Agent 命令（编码助手）

通过对话式 slash 命令在**编码助手内**使用的 AI 驱动设计命令。

**支持的编码助手：**
- 🤖 **Claude Code** (`.claude/commands/`)
- 🤖 **Cursor** (`.cursor/commands/`)
- 🤖 **Windsurf** (`.windsurf/workflows/`)
- 🤖 **Gemini Code Assist** (`.gemini/commands/`)
- 🤖 **Qwen Code** (`.qwen/commands/`)
- 🤖 自定义 agents（通过 `archie custom-agent`）

**工作原理：**
1. 运行 `archie init` 安装 agent 命令文件
2. 在编码助手中打开项目（如 Claude Code、Cursor）
3. 使用 slash 命令 + 对话调用 Archie agents

**可用的 Agent 命令：**

| Slash 命令 | 描述 | 需要的状态 |
|-----------|------|-----------|
| `/init` | 初始化项目工作空间 | N/A |
| `/review` | 明确需求和依赖 | NOT_REVIEWED, UNDER_REVIEW, BLOCKED |
| `/design` | 生成设计产物（workflow、API、storage、metrics） | READY_FOR_DESIGN, UNDER_DESIGN |
| `/test-plan` | 生成全面测试计划 | DESIGNED, SPEC_READY, IMPLEMENTING, FINISHED |
| `/spec` | 生成编码就绪规范 | DESIGNED |
| `/plan` | 管理执行任务和时间线 | DESIGNED, SPEC_READY, IMPLEMENTING, FINISHED |
| `/revise` | 在任何阶段修改/创建/合并/拆分/删除 features | Any |
| `/ask` | 询问关于项目的问题 | Any |
| `/fix` | 修复验证错误或 schema 问题 | Any |

---

## 快速开始

### 步骤 1: 安装 Archie (CLI)

```bash
go install github.com/GarrickZ2/archie@latest
```

### 步骤 2: 初始化项目 (CLI)

```bash
mkdir my-project
cd my-project
archie init
```

这会创建：
- 工作空间文件结构（`features/`, `workflow/`, `spec/` 等）
- Schema 模板（`.archie/docs/schema/`）
- 编码助手的 agent 命令文件

### 步骤 3: 设置项目上下文 (CLI)

```bash
archie setup
```

交互式 TUI：
- 编辑 `background.md`（项目上下文、技术栈、硬性规则）
- 创建和管理 feature 文件

### 步骤 4: 在编码助手中使用 Agents

在 **Claude Code**、**Cursor**、**Windsurf** 或你喜欢的编码助手中打开项目。

**Claude Code 中的示例对话：**

```
你: /review user-authentication

Claude: 我来帮你审查 user-authentication feature。
       让我明确需求和依赖...

       [进行 AI 辅助审查，发现依赖，识别阻塞项，更新 feature 文档]

       ✅ 审查完成。状态 → READY_FOR_DESIGN

---

你: /design user-authentication

Claude: 我来生成设计产物。

       [创建工作流图、API 规范、存储 schema、指标定义、任务分解]

       ✅ 设计完成。已生成：
          - workflow/user-authentication/main.mmd
          - api/api.md (已更新)
          - storage.md (已更新)
          - metrics.md (已更新)

---

你: /test-plan user-authentication

Claude: 我来创建全面的测试计划。

       [生成测试用例矩阵（Unit/Integration/E2E）、
        失败场景、覆盖率思维导图、自动化策略]

       ✅ 测试计划完成。
          - testplan/user-authentication.md

---

你: /spec user-authentication

Claude: 我来生成编码就绪规范。

       [聚合所有设计产物为可执行 spec]

       ✅ Spec 就绪。状态 → SPEC_READY
          - spec/user-authentication.spec.md
```

### 步骤 5: 监控进度 (CLI)

```bash
archie status
```

交互式 TUI 显示：
- 整体项目健康度
- Feature 状态分布
- 被阻塞的 features
- 详细 feature 信息

### 步骤 6: 导出文档 (CLI)

```bash
archie export
```

生成单个 markdown 文件，包含：
- 选择的文档
- Feature 规范
- 目录、统计和依赖图

---

## Archie 能做什么？

使用 Archie 你可以：

- ✅ 从空项目或杂乱笔记初始化工程
- ✅ 系统化 Review Feature，澄清需求
- ✅ 设计 Workflow、API、存储和可观测性
- ✅ 生成全面的测试计划
- ✅ 生成可直接交付给 Coding Agent / SDE 的 Spec
- ✅ 管理 Blocker、Tasks 和项目推进
- ✅ 导出结构化设计文档
- ✅ 管理 feature 生命周期（创建/合并/拆分/删除）

一切通过 CLI 完成，支持交互式 TUI。

---

## 项目结构

```
.
├── background.md           # 项目上下文、技术栈、硬性规则
├── features/
│   ├── README.md          # Feature 注册表
│   └── <feature-key>.md   # Feature 定义
├── workflow/
│   └── <feature-key>/
│       ├── workflow.md    # 工作流叙述
│       └── main.mmd       # 主流程图（Mermaid）
├── spec/
│   └── <feature-key>.spec.md  # 编码就绪规范
├── testplan/
│   └── <feature-key>.md       # 测试用例清单
├── tasks.md               # 每个 feature 的执行任务
├── metrics.md            # 每个 feature 的可观测性
├── storage.md            # 每个 feature 的数据库 schema
├── api/
│   └── api.md            # API 索引和变更记录
├── dependency.md         # 依赖目录
└── deployment.md         # 发布说明和检查清单
```

---

## 命令参考

### CLI 命令

| 命令 | 描述 |
|------|------|
| `archie init` | 用 schema 模板初始化 Archie 工作空间 |
| `archie setup` | 交互式 TUI 编辑背景和管理 features |
| `archie status` | 显示项目状态和交互式 feature 浏览器 |
| `archie export` | 导出文档到单个 markdown 文件 |

### AI Agent 命令

| 命令 | 描述 | 允许的状态 |
|------|------|-----------|
| `archie init` | 初始化项目工作空间并规范化文档 | N/A |
| `archie review <key>` | 明确需求和依赖 | NOT_REVIEWED, UNDER_REVIEW, BLOCKED |
| `archie design <key>` | 生成设计产物（workflow、API、storage、metrics、tasks） | READY_FOR_DESIGN, UNDER_DESIGN |
| `archie test-plan <key>` | 生成全面测试计划 | DESIGNED, SPEC_READY, IMPLEMENTING, FINISHED |
| `archie spec <key>` | 生成编码就绪规范 | DESIGNED |
| `archie plan <key>` | 管理执行任务和时间线 | DESIGNED, SPEC_READY, IMPLEMENTING, FINISHED |
| `archie revise <key>` | 在任何阶段修改/创建/合并/拆分/删除 features | Any |
| `archie ask` | 询问关于项目的问题 | Any |
| `archie fix` | 修复验证错误或 schema 问题 | Any |

---

## Feature 状态机

```
NOT_REVIEWED → UNDER_REVIEW → READY_FOR_DESIGN → UNDER_DESIGN →
DESIGNED → SPEC_READY → IMPLEMENTING → FINISHED

特殊状态: BLOCKED（可以在任何阶段出现）
```

---

## 关键概念

### Commands（命令）
管理 feature 生命周期的高层编排器。

### Subagents（子 Agent）
专注能力单元：
- Workflow Designer（工作流设计器）
- API Designer（API 设计器）
- Storage Designer（存储设计器）
- Metrics Designer（指标设计器）
- Task Manager（任务管理器）
- Test Plan Generator（测试计划生成器）

### Templates（模板）
Archie 使用 schema 模板（`.archie/docs/schema/`），让团队在保持机器可解析的同时维护一致的文档结构。

### State Machine（状态机）
每个 feature 都遵循严格的状态推进，确保实现前的设计质量。

---

## 示例流程

```bash
# 1. 初始化项目
archie init

# 2. 设置背景上下文并创建 features
archie setup

# 3. 审查 feature 的需求
archie review checkout-discount

# 4. 设计 feature
archie design checkout-discount

# 5. 生成测试计划
archie test-plan checkout-discount

# 6. 创建编码规范
archie spec checkout-discount

# 7. 管理执行
archie plan checkout-discount

# 8. 检查状态
archie status

# 9. 导出文档
archie export
```

---

## 高级用法

### 使用 `revise` 管理 Feature

```bash
# 创建新 feature
archie revise --create payment-gateway

# 合并两个 features
archie revise --merge user-login,user-signup --into user-auth

# 拆分 feature
archie revise --split checkout-flow

# 删除 feature
archie revise --delete duplicate-feature

# 更改 feature 状态
archie revise --change checkout-discount --status IMPLEMENTING
```

### 状态监控

```bash
# 交互式状态浏览器
archie status

# 紧凑状态报告
archie status --compact
```

### 文档导出

```bash
# 导出所有选项
archie export

# 导出到指定文件
archie export -o docs/design.md

# 导出时不包含目录
archie export --no-toc

# 导出时不包含统计
archie export --no-stats
```

---

## 适合谁？

- 后端 / 全栈工程师
- Tech Lead
- 初创团队
- 基础设施 / 平台团队
- 讨厌混乱设计文档的人

---

## 设计哲学

- **本地优先**：所有内容存储在你的文件系统中
- **Markdown 原生**：人类可读，可版本控制
- **可预测 > 魔法**：可预测的、明确的行为
- **显式 > 隐式**：没有隐藏状态，清晰的契约
- **人始终掌控决策**：AI 辅助，你决定

---

## 当前状态

Archie 仍处于早期设计和开发阶段。

欢迎反馈、讨论和共建。

---

## 许可证

MIT License - 详见 LICENSE 文件。

---

## 贡献

欢迎贡献！请：
1. Fork 仓库
2. 创建 feature 分支
3. 提交带有清晰描述的 pull request

如有 bug 或功能请求，请开 issue。
