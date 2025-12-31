package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/afero"

	"github.com/GarrickZ2/archie/internal/ui"
	"github.com/GarrickZ2/archie/resources"
)

const (
	backgroundOption = "Background Documentation"
	featureOption    = "Feature Management"
	exitOption       = "[x] Exit"
	newMarker        = " ✨"
)

// Manager 管理项目文档的 setup
type Manager struct {
	projectPath string
	fs          afero.Fs
}

// NewManager 创建 setup manager
func NewManager(projectPath string) *Manager {
	return &Manager{
		projectPath: projectPath,
		fs:          afero.NewOsFs(),
	}
}

// ShowMainUI 显示主界面
func (m *Manager) ShowMainUI() error {
	for {
		// 检查 background.md 状态
		backgroundPath := filepath.Join(m.projectPath, "background.md")
		backgroundEmpty, err := IsFileEmpty(m.fs, backgroundPath)
		if err != nil {
			return fmt.Errorf("failed to check background.md: %w", err)
		}

		// 检查 features 文件夹状态
		featuresEmpty, err := m.isFeaturesEmpty()
		if err != nil {
			return fmt.Errorf("failed to check features folder: %w", err)
		}

		// 构建选项列表
		backgroundLabel := backgroundOption
		if backgroundEmpty {
			backgroundLabel += newMarker
		}

		featureLabel := featureOption
		if featuresEmpty {
			featureLabel += newMarker
		}

		options := []string{
			backgroundLabel,
			featureLabel,
			exitOption,
		}

		// 显示选择界面
		var selected string
		prompt := &survey.Select{
			Message: "Project Setup:",
			Options: options,
		}

		if err := survey.AskOne(prompt, &selected); err != nil {
			return fmt.Errorf("selection cancelled")
		}

		// 处理选择
		switch {
		case strings.HasPrefix(selected, backgroundOption):
			if err := m.handleBackground(); err != nil {
				ui.ShowError(fmt.Sprintf("Background setup failed: %v", err))
				continue
			}

		case strings.HasPrefix(selected, featureOption):
			if err := m.handleFeatures(); err != nil {
				ui.ShowError(fmt.Sprintf("Feature management failed: %v", err))
				continue
			}

		case selected == exitOption:
			return nil
		}
	}
}

// handleBackground 处理 background 文档
func (m *Manager) handleBackground() error {
	backgroundPath := filepath.Join(m.projectPath, "background.md")

	// 检查文件是否为空
	isEmpty, err := IsFileEmpty(m.fs, backgroundPath)
	if err != nil {
		return err
	}

	// 如果为空，先用模板初始化
	if isEmpty {
		template, err := resources.GetSchemaTemplate("background.md")
		if err != nil {
			return fmt.Errorf("failed to load background template: %w", err)
		}

		if err := afero.WriteFile(m.fs, backgroundPath, []byte(template), 0644); err != nil {
			return fmt.Errorf("failed to write background.md: %w", err)
		}

		ui.ShowInfo("Initialized background.md with template")
	}

	// 在编辑器中打开
	return OpenInEditor(backgroundPath)
}

// handleFeatures 处理 feature 管理
func (m *Manager) handleFeatures() error {
	for {
		featuresPath := filepath.Join(m.projectPath, "features")

		// 确保 features 目录存在
		if err := m.fs.MkdirAll(featuresPath, 0755); err != nil {
			return fmt.Errorf("failed to create features directory: %w", err)
		}

		// 获取所有 feature 文件
		features, err := m.listFeatures()
		if err != nil {
			return err
		}

		// 构建选项列表
		options := []string{"[+ Add new feature]"}

		for _, feature := range features {
			label := feature.Name
			if feature.IsEmpty {
				label += newMarker
			}
			options = append(options, label)
		}

		options = append(options, "[<- Back]")

		// 显示选择界面
		var selected string
		prompt := &survey.Select{
			Message: "Feature Management:",
			Options: options,
		}

		if err := survey.AskOne(prompt, &selected); err != nil {
			return fmt.Errorf("selection cancelled")
		}

		// 处理选择
		switch {
		case selected == "[+ Add new feature]":
			if err := m.addNewFeature(); err != nil {
				ui.ShowError(fmt.Sprintf("Failed to add feature: %v", err))
				continue
			}

		case selected == "[<- Back]":
			return nil

		default:
			// 打开已有的 feature
			featureName := strings.TrimSuffix(selected, newMarker)
			featureName = strings.TrimSpace(featureName)
			if err := m.openFeature(featureName); err != nil {
				ui.ShowError(fmt.Sprintf("Failed to open feature: %v", err))
				continue
			}
		}
	}
}

// addNewFeature 添加新的 feature
func (m *Manager) addNewFeature() error {
	// 提示输入 feature 名称
	var featureName string
	namePrompt := &survey.Input{
		Message: "Feature name:",
		Help:    "Enter the feature name (will be converted to kebab-case)",
	}

	if err := survey.AskOne(namePrompt, &featureName, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// 转换为 kebab-case
	featureKey := ToKebabCase(featureName)
	if featureKey == "" {
		return fmt.Errorf("invalid feature name")
	}

	// 创建文件路径
	featurePath := filepath.Join(m.projectPath, "features", featureKey+".md")

	// 检查文件是否已存在
	exists, err := afero.Exists(m.fs, featurePath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("feature '%s' already exists", featureKey)
	}

	// 加载模板
	template, err := resources.GetSchemaTemplate("feature.md")
	if err != nil {
		return fmt.Errorf("failed to load feature template: %w", err)
	}

	// 替换模板中的 feature-key
	content := strings.ReplaceAll(template, "<feature-key>", featureKey)

	// 写入文件
	if err := afero.WriteFile(m.fs, featurePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create feature file: %w", err)
	}

	ui.ShowSuccess(fmt.Sprintf("Created feature: %s", featureKey))

	// 在编辑器中打开
	return OpenInEditor(featurePath)
}

// openFeature 打开已有的 feature
func (m *Manager) openFeature(featureName string) error {
	featurePath := filepath.Join(m.projectPath, "features", featureName+".md")

	// 检查文件是否为空
	isEmpty, err := IsFileEmpty(m.fs, featurePath)
	if err != nil {
		return err
	}

	// 如果为空，用模板初始化
	if isEmpty {
		template, err := resources.GetSchemaTemplate("feature.md")
		if err != nil {
			return fmt.Errorf("failed to load feature template: %w", err)
		}

		// 替换模板中的 feature-key
		content := strings.ReplaceAll(template, "<feature-key>", featureName)

		if err := afero.WriteFile(m.fs, featurePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write feature file: %w", err)
		}

		ui.ShowInfo(fmt.Sprintf("Initialized %s with template", featureName))
	}

	// 在编辑器中打开
	return OpenInEditor(featurePath)
}

// isFeaturesEmpty 检查 features 文件夹是否为空（完全没有 feature 文件）
func (m *Manager) isFeaturesEmpty() (bool, error) {
	featuresPath := filepath.Join(m.projectPath, "features")

	exists, err := afero.DirExists(m.fs, featuresPath)
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}

	files, err := afero.ReadDir(m.fs, featuresPath)
	if err != nil {
		return false, err
	}

	// 检查是否有任何 .md 文件
	mdCount := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
			mdCount++
		}
	}

	return mdCount == 0, nil
}

// FeatureInfo feature 信息
type FeatureInfo struct {
	Name    string
	IsEmpty bool
}

// listFeatures 列出所有 features
func (m *Manager) listFeatures() ([]FeatureInfo, error) {
	featuresPath := filepath.Join(m.projectPath, "features")

	files, err := afero.ReadDir(m.fs, featuresPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []FeatureInfo{}, nil
		}
		return nil, err
	}

	var features []FeatureInfo
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		// 获取 feature 名称（去除 .md 扩展名）
		name := strings.TrimSuffix(file.Name(), ".md")

		// 检查文件是否为空
		filePath := filepath.Join(featuresPath, file.Name())
		isEmpty, err := IsFileEmpty(m.fs, filePath)
		if err != nil {
			return nil, err
		}

		features = append(features, FeatureInfo{
			Name:    name,
			IsEmpty: isEmpty,
		})
	}

	return features, nil
}
