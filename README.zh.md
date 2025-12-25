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
- 🧩 **可组合的子 Agent（API / Workflow / Storage / Metrics / Tasks）**

---

## 核心思想

Archie 把 **Feature** 作为最小原子，并强制最小设计链路：

> **Feature → Workflow → Spec**

其他内容（API、存储、指标、任务、发布）都是对这条主链路的支撑。

---

## Archie 能做什么？

使用 Archie 你可以：

- 从空项目或杂乱笔记初始化工程
- 系统化 Review Feature，澄清需求
- 设计 Workflow、API、存储和可观测性
- 生成可直接交付给 Coding Agent / SDE 的 Spec
- 管理 Blocker、Tasks 和项目推进
- 导出结构化设计文档

一切通过 CLI 完成。

---

## 示例流程

```bash
archie init
archie review checkout-discount
archie design checkout-discount
archie spec checkout-discount
archie update-progress checkout-discount --note "开始实现"
```

---

## 项目结构（简化）

```
.
├── background.md
├── features/
│   └── checkout-discount.md
├── workflow/
│   └── checkout-discount/
├── spec/
│   └── checkout-discount.spec.md
├── tasks.md
├── metrics.md
└── api/
```

---

## 关键概念

### Commands（命令）
高层编排指令：
- `init`
- `review`
- `design`
- `spec`
- `update-progress`

### Subagents（子 Agent）
专注能力单元：
- Workflow Designer
- API Designer
- Storage Designer
- Metrics Designer
- Task Manager

### Templates（模板）
Archie 提供轻量模板系统，允许团队自定义文档结构，同时保持可解析性。

---

## 适合谁？

- 后端 / 全栈工程师
- Tech Lead
- 初创团队
- 基础设施 / 平台团队
- 讨厌混乱设计文档的人

---

## 设计哲学

- 本地优先
- Markdown 原生
- 可预测 > 魔法
- 显式 > 隐式
- 人始终掌控决策

---

## 当前状态

Archie 仍处于早期设计和开发阶段。

欢迎反馈、讨论和共建。
