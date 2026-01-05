package status

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// FeatureStatus 定义 feature 状态
type FeatureStatus string

const (
	StatusNotReviewed    FeatureStatus = "NOT_REVIEWED"
	StatusUnderReview    FeatureStatus = "UNDER_REVIEW"
	StatusBlocked        FeatureStatus = "BLOCKED"
	StatusReadyForDesign FeatureStatus = "READY_FOR_DESIGN"
	StatusUnderDesign    FeatureStatus = "UNDER_DESIGN"
	StatusDesigned       FeatureStatus = "DESIGNED"
	StatusSpecReady      FeatureStatus = "SPEC_READY"
	StatusImplementing   FeatureStatus = "IMPLEMENTING"
	StatusFinished       FeatureStatus = "FINISHED"
	StatusUnknown        FeatureStatus = "UNKNOWN"
)

// AllStatuses 所有可能的状态（按流程顺序）
var AllStatuses = []FeatureStatus{
	StatusNotReviewed,
	StatusUnderReview,
	StatusBlocked,
	StatusReadyForDesign,
	StatusUnderDesign,
	StatusDesigned,
	StatusSpecReady,
	StatusImplementing,
	StatusFinished,
}

// Feature 表示一个 feature 及其状态信息
type Feature struct {
	Name         string
	Status       FeatureStatus
	Owner        string
	LastUpdated  string
	Reason       string
	FilePath     string
	Dependencies map[string]string // feature-key -> reason
}

// Parser 解析 feature 文件
type Parser struct {
	fs afero.Fs
}

// NewParser 创建解析器
func NewParser(fs afero.Fs) *Parser {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &Parser{fs: fs}
}

// ParseFeaturesDir 解析 features 目录下的所有 feature 文件
func (p *Parser) ParseFeaturesDir(projectPath string) ([]Feature, error) {
	featuresPath := filepath.Join(projectPath, "features")

	// 检查目录是否存在
	exists, err := afero.DirExists(p.fs, featuresPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check features directory: %w", err)
	}
	if !exists {
		return []Feature{}, nil
	}

	// 读取目录中的文件
	files, err := afero.ReadDir(p.fs, featuresPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read features directory: %w", err)
	}

	var features []Feature
	for _, file := range files {
		// 跳过目录
		if file.IsDir() {
			continue
		}

		// 只处理 .md 文件
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(featuresPath, file.Name())
		feature, err := p.ParseFeatureFile(filePath)
		if err != nil {
			// 如果解析失败，创建一个带有 UNKNOWN 状态的 feature
			featureName := strings.TrimSuffix(file.Name(), ".md")
			feature = Feature{
				Name:     featureName,
				Status:   StatusUnknown,
				FilePath: filePath,
			}
		}

		features = append(features, feature)
	}

	return features, nil
}

// ParseFeatureFile 解析单个 feature 文件
func (p *Parser) ParseFeatureFile(filePath string) (Feature, error) {
	feature := Feature{
		FilePath:     filePath,
		Name:         strings.TrimSuffix(filepath.Base(filePath), ".md"),
		Status:       StatusUnknown,
		Dependencies: make(map[string]string),
	}

	file, err := p.fs.Open(filePath)
	if err != nil {
		return feature, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inStatusSection := false
	inDependenciesSection := false

	// 正则表达式匹配 Status section 中的字段
	statusRegex := regexp.MustCompile(`^-\s*Value:\s*(.+)$`)
	ownerRegex := regexp.MustCompile(`^-\s*Owner:\s*(.*)$`)
	dateRegex := regexp.MustCompile(`^-\s*Last Updated:\s*(.*)$`)
	reasonRegex := regexp.MustCompile(`^-\s*Reason:\s*(.*)$`)

	// 匹配 Feature Dependencies: - `feature-key`: [Reason]
	dependencyRegex := regexp.MustCompile(`^-\s*` + "`" + `([^` + "`" + `]+)` + "`" + `\s*:\s*(.*)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 检测 ## Status section
		if strings.HasPrefix(line, "## Status") {
			inStatusSection = true
			inDependenciesSection = false
			continue
		}

		// 检测 ## Feature Dependencies section
		if strings.HasPrefix(line, "## Feature Dependencies") {
			inStatusSection = false
			inDependenciesSection = true
			continue
		}

		// 如果遇到下一个 section，停止解析当前 section
		if strings.HasPrefix(line, "##") {
			inStatusSection = false
			inDependenciesSection = false
			continue
		}

		// 在 Status section 中解析字段
		if inStatusSection {
			if matches := statusRegex.FindStringSubmatch(line); len(matches) > 1 {
				feature.Status = FeatureStatus(strings.TrimSpace(matches[1]))
			} else if matches := ownerRegex.FindStringSubmatch(line); len(matches) > 1 {
				feature.Owner = strings.TrimSpace(matches[1])
			} else if matches := dateRegex.FindStringSubmatch(line); len(matches) > 1 {
				feature.LastUpdated = strings.TrimSpace(matches[1])
			} else if matches := reasonRegex.FindStringSubmatch(line); len(matches) > 1 {
				feature.Reason = strings.TrimSpace(matches[1])
			}
		}

		// 在 Feature Dependencies section 中解析依赖
		if inDependenciesSection {
			if matches := dependencyRegex.FindStringSubmatch(line); len(matches) > 2 {
				featureKey := strings.TrimSpace(matches[1])
				reason := strings.TrimSpace(matches[2])
				// 过滤掉模板占位符
				if featureKey != "<feature-key>" && featureKey != "" {
					feature.Dependencies[featureKey] = reason
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return feature, fmt.Errorf("error reading file: %w", err)
	}

	return feature, nil
}

// IsValidStatus 检查状态是否有效
func IsValidStatus(status FeatureStatus) bool {
	for _, s := range AllStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// GetStatusProgress 获取状态在流程中的进度（0-100）
func GetStatusProgress(status FeatureStatus) int {
	for i, s := range AllStatuses {
		if s == status {
			// BLOCKED 不计入进度
			if status == StatusBlocked {
				return 0
			}
			// 计算进度百分比（排除 BLOCKED）
			totalSteps := len(AllStatuses) - 1 // 减去 BLOCKED
			currentStep := i
			if i > 2 { // BLOCKED 之后的状态要减1
				currentStep--
			}
			return (currentStep * 100) / (totalSteps - 1)
		}
	}
	return 0
}

// GetStatusColor 获取状态对应的颜色代码
func GetStatusColor(status FeatureStatus) string {
	switch status {
	case StatusBlocked:
		return "\033[31m" // Red
	case StatusNotReviewed, StatusUnknown:
		return "\033[90m" // Gray
	case StatusUnderReview, StatusReadyForDesign:
		return "\033[33m" // Yellow
	case StatusUnderDesign, StatusDesigned:
		return "\033[36m" // Cyan
	case StatusSpecReady, StatusImplementing:
		return "\033[34m" // Blue
	case StatusFinished:
		return "\033[32m" // Green
	default:
		return "\033[0m" // Reset
	}
}

// IsOld 检查 LastUpdated 是否超过指定天数
func (f *Feature) IsOld(days int) bool {
	if f.LastUpdated == "" || f.LastUpdated == "YYYY-MM-DD" {
		return false
	}

	// 尝试解析日期
	date, err := time.Parse("2006-01-02", f.LastUpdated)
	if err != nil {
		return false
	}

	daysSince := time.Since(date).Hours() / 24
	return int(daysSince) > days
}
