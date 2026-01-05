package status

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// FeatureDetail 包含 feature 的完整详细信息
type FeatureDetail struct {
	// 基本信息
	Key      string
	FilePath string

	// Status section
	Status      FeatureStatus
	Owner       string
	LastUpdated string
	Reason      string

	// Summary section
	OneLiner   string
	Background string
	UserStory  string

	// Scope section
	InScope  []string
	OutScope []string

	// Requirements
	Requirements    []string
	NonRequirements []string

	// Feature Dependencies (features that should be designed before this one)
	FeatureDependencies map[string]string // feature-key -> reason

	// Acceptance Criteria
	AcceptanceCriteria []string

	// Design Constraints
	DesignConstraints []string

	// Design Artifacts
	APIDesign      string
	StorageDesign  string
	WorkflowDesign string
	MetricsDesign  string
	TasksDesign    string

	// Spec
	SpecLocation  string
	SpecReadiness string

	// Related Records
	Blockers string

	// Changelog
	Changelog []string
}

// DetailParser 解析 feature 详细信息
type DetailParser struct {
	fs afero.Fs
}

// NewDetailParser 创建详细信息解析器
func NewDetailParser(fs afero.Fs) *DetailParser {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &DetailParser{fs: fs}
}

// ParseFeatureDetail 解析单个 feature 的完整详细信息
func (p *DetailParser) ParseFeatureDetail(projectPath, featureKey string) (*FeatureDetail, error) {
	filePath := filepath.Join(projectPath, "features", featureKey+".md")

	// 检查文件是否存在
	exists, err := afero.Exists(p.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to check if file exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("feature file not found: %s", filePath)
	}

	detail := &FeatureDetail{
		Key:                 featureKey,
		FilePath:            filePath,
		Status:              StatusUnknown,
		FeatureDependencies: make(map[string]string),
	}

	file, err := p.fs.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentSection := ""
	currentSubSection := ""

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// 检测主要 sections (## XXX)
		if strings.HasPrefix(trimmedLine, "## ") {
			currentSection = strings.TrimPrefix(trimmedLine, "## ")
			currentSubSection = ""
			continue
		}

		// 检测子 sections (### XXX)
		if strings.HasPrefix(trimmedLine, "### ") {
			currentSubSection = strings.TrimPrefix(trimmedLine, "### ")
			continue
		}

		// 解析各个 section 的内容
		p.parseSectionContent(detail, currentSection, currentSubSection, trimmedLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return detail, nil
}

// parseSectionContent 解析各个 section 的内容
func (p *DetailParser) parseSectionContent(detail *FeatureDetail, section, subSection, line string) {
	if line == "" {
		return
	}

	switch section {
	case "Status":
		p.parseStatusSection(detail, line)
	case "Summary":
		p.parseSummarySection(detail, line)
	case "Scope":
		p.parseScopeSection(detail, subSection, line)
	case "Requirements":
		if strings.HasPrefix(line, "- R") || strings.HasPrefix(line, "-") {
			detail.Requirements = append(detail.Requirements, strings.TrimPrefix(line, "- "))
		}
	case "Non-Requirements":
		if strings.HasPrefix(line, "- NR") || strings.HasPrefix(line, "-") {
			detail.NonRequirements = append(detail.NonRequirements, strings.TrimPrefix(line, "- "))
		}
	case "Feature Dependencies":
		p.parseFeatureDependenciesSection(detail, line)
	case "Acceptance Criteria":
		if strings.HasPrefix(line, "- AC") || strings.HasPrefix(line, "-") {
			detail.AcceptanceCriteria = append(detail.AcceptanceCriteria, strings.TrimPrefix(line, "- "))
		}
	case "Design Constraints":
		if strings.HasPrefix(line, "- ") {
			detail.DesignConstraints = append(detail.DesignConstraints, strings.TrimPrefix(line, "- "))
		}
	case "Design Artifacts":
		p.parseDesignArtifacts(detail, line)
	case "Spec":
		p.parseSpecSection(detail, line)
	case "Related Records":
		p.parseRelatedRecords(detail, line)
	case "Changelog":
		if strings.HasPrefix(line, "- ") {
			detail.Changelog = append(detail.Changelog, strings.TrimPrefix(line, "- "))
		}
	}
}

// parseStatusSection 解析 Status section
func (p *DetailParser) parseStatusSection(detail *FeatureDetail, line string) {
	if strings.HasPrefix(line, "- Value:") {
		detail.Status = FeatureStatus(strings.TrimSpace(strings.TrimPrefix(line, "- Value:")))
	} else if strings.HasPrefix(line, "- Owner:") {
		detail.Owner = strings.TrimSpace(strings.TrimPrefix(line, "- Owner:"))
	} else if strings.HasPrefix(line, "- Last Updated:") {
		detail.LastUpdated = strings.TrimSpace(strings.TrimPrefix(line, "- Last Updated:"))
	} else if strings.HasPrefix(line, "- Reason:") {
		detail.Reason = strings.TrimSpace(strings.TrimPrefix(line, "- Reason:"))
	}
}

// parseSummarySection 解析 Summary section
func (p *DetailParser) parseSummarySection(detail *FeatureDetail, line string) {
	if strings.HasPrefix(line, "- One-liner:") {
		detail.OneLiner = strings.TrimSpace(strings.TrimPrefix(line, "- One-liner:"))
	} else if strings.HasPrefix(line, "- Background / Motivation:") {
		detail.Background = strings.TrimSpace(strings.TrimPrefix(line, "- Background / Motivation:"))
	} else if strings.HasPrefix(line, "- User story / Use case:") {
		detail.UserStory = strings.TrimSpace(strings.TrimPrefix(line, "- User story / Use case:"))
	}
}

// parseScopeSection 解析 Scope section
func (p *DetailParser) parseScopeSection(detail *FeatureDetail, subSection, line string) {
	if !strings.HasPrefix(line, "- ") {
		return
	}

	item := strings.TrimPrefix(line, "- ")
	if subSection == "In Scope" {
		detail.InScope = append(detail.InScope, item)
	} else if subSection == "Out of Scope" {
		detail.OutScope = append(detail.OutScope, item)
	}
}

// parseFeatureDependenciesSection 解析 Feature Dependencies section
// Format: - `feature-key`: [Reason]
func (p *DetailParser) parseFeatureDependenciesSection(detail *FeatureDetail, line string) {
	if !strings.HasPrefix(line, "- ") {
		return
	}

	// 提取反引号包裹的 feature-key
	item := strings.TrimPrefix(line, "- ")

	// 查找第一个反引号对
	start := strings.Index(item, "`")
	if start == -1 {
		return
	}

	end := strings.Index(item[start+1:], "`")
	if end == -1 {
		return
	}

	featureKey := item[start+1 : start+1+end]

	// 查找冒号后的原因
	colonIdx := strings.Index(item[start+1+end:], ":")
	if colonIdx == -1 {
		return
	}

	reason := strings.TrimSpace(item[start+1+end+colonIdx+1:])

	// 过滤模板占位符
	if featureKey != "<feature-key>" && featureKey != "" {
		detail.FeatureDependencies[featureKey] = reason
	}
}

// parseDesignArtifacts 解析 Design Artifacts section
func (p *DetailParser) parseDesignArtifacts(detail *FeatureDetail, line string) {
	if strings.HasPrefix(line, "- API:") {
		detail.APIDesign = strings.TrimSpace(strings.TrimPrefix(line, "- API:"))
	} else if strings.HasPrefix(line, "- Storage:") {
		detail.StorageDesign = strings.TrimSpace(strings.TrimPrefix(line, "- Storage:"))
	} else if strings.HasPrefix(line, "- Workflow:") {
		detail.WorkflowDesign = strings.TrimSpace(strings.TrimPrefix(line, "- Workflow:"))
	} else if strings.HasPrefix(line, "- Metrics:") {
		detail.MetricsDesign = strings.TrimSpace(strings.TrimPrefix(line, "- Metrics:"))
	} else if strings.HasPrefix(line, "- Tasks:") {
		detail.TasksDesign = strings.TrimSpace(strings.TrimPrefix(line, "- Tasks:"))
	}
}

// parseSpecSection 解析 Spec section
func (p *DetailParser) parseSpecSection(detail *FeatureDetail, line string) {
	if strings.HasPrefix(line, "- Location:") {
		detail.SpecLocation = strings.TrimSpace(strings.TrimPrefix(line, "- Location:"))
	} else if strings.HasPrefix(line, "- Readiness:") {
		detail.SpecReadiness = strings.TrimSpace(strings.TrimPrefix(line, "- Readiness:"))
	}
}

// parseRelatedRecords 解析 Related Records section
func (p *DetailParser) parseRelatedRecords(detail *FeatureDetail, line string) {
	if strings.HasPrefix(line, "- Blockers:") {
		detail.Blockers = strings.TrimSpace(strings.TrimPrefix(line, "- Blockers:"))
	}
}

// DetailDisplay 负责展示详细的 feature 信息
type DetailDisplay struct {
	detail *FeatureDetail
}

// NewDetailDisplay 创建详细信息展示器
func NewDetailDisplay(detail *FeatureDetail) *DetailDisplay {
	return &DetailDisplay{detail: detail}
}

// Show 展示详细的 feature 信息
func (d *DetailDisplay) Show() {
	fmt.Println()
	d.showHeader()
	fmt.Println()
	d.showStatus()
	fmt.Println()
	d.showSummary()

	if len(d.detail.InScope) > 0 || len(d.detail.OutScope) > 0 {
		fmt.Println()
		d.showScope()
	}

	if len(d.detail.Requirements) > 0 {
		fmt.Println()
		d.showRequirements()
	}

	if len(d.detail.NonRequirements) > 0 {
		fmt.Println()
		d.showNonRequirements()
	}

	if len(d.detail.FeatureDependencies) > 0 {
		fmt.Println()
		d.showFeatureDependencies()
	}

	if len(d.detail.AcceptanceCriteria) > 0 {
		fmt.Println()
		d.showAcceptanceCriteria()
	}

	if len(d.detail.DesignConstraints) > 0 {
		fmt.Println()
		d.showDesignConstraints()
	}

	if d.hasDesignArtifacts() {
		fmt.Println()
		d.showDesignArtifacts()
	}

	if d.detail.SpecLocation != "" || d.detail.SpecReadiness != "" {
		fmt.Println()
		d.showSpec()
	}

	if d.detail.Blockers != "" {
		fmt.Println()
		d.showRelatedRecords()
	}

	if len(d.detail.Changelog) > 0 {
		fmt.Println()
		d.showChangelog()
	}

	fmt.Println()
}

// showHeader 显示标题
func (d *DetailDisplay) showHeader() {
	title := fmt.Sprintf("📄 Feature: %s", d.detail.Key)
	separator := strings.Repeat("═", 70)

	fmt.Println(ColorCyan + "  ╔" + separator + "╗" + ColorReset)
	padding := (70 - len(title)) / 2
	fmt.Print(ColorCyan + "  ║ " + ColorReset)
	fmt.Print(strings.Repeat(" ", padding))
	fmt.Print(ColorBold + title + ColorReset)
	fmt.Print(strings.Repeat(" ", 70-len(title)-padding))
	fmt.Println(ColorCyan + " ║" + ColorReset)
	fmt.Println(ColorCyan + "  ╚" + separator + "╝" + ColorReset)
}

// showStatus 显示状态信息
func (d *DetailDisplay) showStatus() {
	fmt.Println(ColorBold + "  📊 Status" + ColorReset)
	fmt.Println()

	statusColor := GetStatusColor(d.detail.Status)
	statusName := strings.ReplaceAll(string(d.detail.Status), "_", " ")

	fmt.Printf("  %s%-15s%s %s%s%s\n", ColorDim, "Current:", ColorReset, statusColor+ColorBold, statusName, ColorReset)

	if d.detail.Owner != "" && d.detail.Owner != "YYYY-MM-DD" {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "Owner:", ColorReset, d.detail.Owner)
	}

	if d.detail.LastUpdated != "" && d.detail.LastUpdated != "YYYY-MM-DD" {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "Last Updated:", ColorReset, d.detail.LastUpdated)
	}

	if d.detail.Reason != "" {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "Reason:", ColorReset, d.detail.Reason)
	}

	// 显示进度条
	progress := GetStatusProgress(d.detail.Status)
	if progress > 0 {
		fmt.Println()
		d.showProgressBar(progress)
	}
}

// showProgressBar 显示进度条
func (d *DetailDisplay) showProgressBar(progress int) {
	barWidth := 50
	filled := (progress * barWidth) / 100

	barColor := ColorGreen
	if progress < 30 {
		barColor = ColorRed
	} else if progress < 70 {
		barColor = ColorYellow
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Printf("  %sProgress:%s %s%s%s %s%d%%%s\n",
		ColorDim, ColorReset,
		barColor, bar, ColorReset,
		ColorBold, progress, ColorReset)
}

// showSummary 显示概要信息
func (d *DetailDisplay) showSummary() {
	fmt.Println(ColorBold + "  📝 Summary" + ColorReset)
	fmt.Println()

	if d.detail.OneLiner != "" {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "One-liner:", ColorReset, d.detail.OneLiner)
	}

	if d.detail.Background != "" {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "Background:", ColorReset, d.detail.Background)
	}

	if d.detail.UserStory != "" {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "User Story:", ColorReset, d.detail.UserStory)
	}
}

// showScope 显示范围
func (d *DetailDisplay) showScope() {
	fmt.Println(ColorBold + "  🎯 Scope" + ColorReset)
	fmt.Println()

	if len(d.detail.InScope) > 0 {
		fmt.Println(ColorGreen + "  ✓ In Scope:" + ColorReset)
		for _, item := range d.detail.InScope {
			if item != "" && item != "..." {
				fmt.Printf("    • %s\n", item)
			}
		}
		fmt.Println()
	}

	if len(d.detail.OutScope) > 0 {
		fmt.Println(ColorRed + "  ✗ Out of Scope:" + ColorReset)
		for _, item := range d.detail.OutScope {
			if item != "" && item != "..." {
				fmt.Printf("    • %s\n", item)
			}
		}
	}
}

// showRequirements 显示需求
func (d *DetailDisplay) showRequirements() {
	fmt.Println(ColorBold + "  ✅ Requirements" + ColorReset)
	fmt.Println()

	for _, req := range d.detail.Requirements {
		if req != "" && req != "R1:" && req != "R2:" {
			fmt.Printf("  • %s\n", req)
		}
	}
}

// showNonRequirements 显示非需求
func (d *DetailDisplay) showNonRequirements() {
	fmt.Println(ColorBold + "  ⛔ Non-Requirements" + ColorReset)
	fmt.Println()

	for _, nonReq := range d.detail.NonRequirements {
		if nonReq != "" && nonReq != "NR1:" && nonReq != "NR2:" {
			fmt.Printf("  • %s\n", nonReq)
		}
	}
}

// showFeatureDependencies 显示 feature 依赖关系
func (d *DetailDisplay) showFeatureDependencies() {
	fmt.Println(ColorBold + "  🔗 Feature Dependencies" + ColorReset)
	fmt.Println(ColorDim + "  (Recommended features to be designed before this one)" + ColorReset)
	fmt.Println()

	if len(d.detail.FeatureDependencies) > 0 {
		for featureKey, reason := range d.detail.FeatureDependencies {
			fmt.Printf("    • %s%s%s: %s\n", ColorBold, featureKey, ColorReset, reason)
		}
	}
}

// showAcceptanceCriteria 显示验收标准
func (d *DetailDisplay) showAcceptanceCriteria() {
	fmt.Println(ColorBold + "  🎯 Acceptance Criteria" + ColorReset)
	fmt.Println()

	for _, ac := range d.detail.AcceptanceCriteria {
		if ac != "" && ac != "AC1:" && ac != "AC2:" {
			fmt.Printf("  • %s\n", ac)
		}
	}
}

// showDesignConstraints 显示设计约束
func (d *DetailDisplay) showDesignConstraints() {
	fmt.Println(ColorBold + "  ⚙️  Design Constraints" + ColorReset)
	fmt.Println()

	for _, constraint := range d.detail.DesignConstraints {
		if constraint != "" && constraint != "..." {
			fmt.Printf("  • %s\n", constraint)
		}
	}
}

// hasDesignArtifacts 检查是否有设计制品
func (d *DetailDisplay) hasDesignArtifacts() bool {
	return d.detail.APIDesign != "" ||
		d.detail.StorageDesign != "" ||
		d.detail.WorkflowDesign != "" ||
		d.detail.MetricsDesign != "" ||
		d.detail.TasksDesign != ""
}

// showDesignArtifacts 显示设计制品
func (d *DetailDisplay) showDesignArtifacts() {
	fmt.Println(ColorBold + "  🎨 Design Artifacts" + ColorReset)
	fmt.Println()

	artifacts := []struct {
		label string
		value string
		icon  string
	}{
		{"API", d.detail.APIDesign, "🔌"},
		{"Storage", d.detail.StorageDesign, "💾"},
		{"Workflow", d.detail.WorkflowDesign, "🔄"},
		{"Metrics", d.detail.MetricsDesign, "📊"},
		{"Tasks", d.detail.TasksDesign, "✓"},
	}

	for _, artifact := range artifacts {
		if artifact.value != "" && !strings.Contains(artifact.value, "<feature-key>") {
			fmt.Printf("  %s %s%-12s%s %s\n",
				artifact.icon,
				ColorDim, artifact.label+":", ColorReset,
				artifact.value)
		}
	}
}

// showSpec 显示规格说明
func (d *DetailDisplay) showSpec() {
	fmt.Println(ColorBold + "  📋 Specification" + ColorReset)
	fmt.Println()

	if d.detail.SpecLocation != "" && !strings.Contains(d.detail.SpecLocation, "<feature-key>") {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "Location:", ColorReset, d.detail.SpecLocation)
	}

	if d.detail.SpecReadiness != "" {
		readinessColor := ColorGray
		switch d.detail.SpecReadiness {
		case "READY":
			readinessColor = ColorGreen
		case "DRAFT":
			readinessColor = ColorYellow
		case "NONE":
			readinessColor = ColorRed
		}
		fmt.Printf("  %s%-15s%s %s%s%s\n",
			ColorDim, "Readiness:", ColorReset,
			readinessColor, d.detail.SpecReadiness, ColorReset)
	}
}

// showRelatedRecords 显示相关记录
func (d *DetailDisplay) showRelatedRecords() {
	fmt.Println(ColorBold + "  📚 Related Records" + ColorReset)
	fmt.Println()

	if d.detail.Blockers != "" && !strings.Contains(d.detail.Blockers, "<feature-key>") {
		fmt.Printf("  %s%-15s%s %s\n", ColorDim, "Blockers:", ColorReset, d.detail.Blockers)
	}
}

// showChangelog 显示变更日志
func (d *DetailDisplay) showChangelog() {
	fmt.Println(ColorBold + "  📅 Changelog" + ColorReset)
	fmt.Println()

	for _, entry := range d.detail.Changelog {
		if entry != "" && !strings.HasPrefix(entry, "YYYY-MM-DD") {
			fmt.Printf("  • %s\n", entry)
		}
	}
}
