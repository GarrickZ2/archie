package agent

import (
	"fmt"
	"strings"

	"github.com/GarrickZ2/archie/resources"
	"gopkg.in/yaml.v3"
)

// Formatter 格式化器接口
type Formatter interface {
	Format(template *resources.CommandTemplate, mapping map[string]string) (string, error)
}

// MDFormatter MD 格式化器
type MDFormatter struct{}

// Format 将 CommandTemplate 格式化为 Markdown（带 frontmatter）
func (f *MDFormatter) Format(template *resources.CommandTemplate, mapping map[string]string) (string, error) {
	// 构建 frontmatter
	frontmatter := make(map[string]interface{})
	content := template.Content

	// 处理 mapping：找到 content 字段的映射
	// 默认情况下，content 字段映射到 [CONTENT]（正文）
	contentMappedToContent := false
	for yamlKey, targetKey := range mapping {
		if targetKey == "[CONTENT]" {
			if yamlKey == "content" {
				// content 字段映射到正文，使用 template.Content
				contentMappedToContent = true
			} else {
				// 其他字段映射到正文，从 metadata 中获取
				if val, ok := template.Metadata[yamlKey]; ok {
					if strVal, ok := val.(string); ok {
						content = strVal
					}
				}
			}
			continue
		}
		// 其他字段映射到 frontmatter
		if yamlKey == "content" {
			// content 字段不映射到 [CONTENT]，跳过（不应该发生，但安全处理）
			continue
		}
		// 从 metadata 中获取值并映射到 frontmatter
		if val, ok := template.Metadata[yamlKey]; ok {
			frontmatter[targetKey] = val
		}
	}

	// 如果没有找到 content 映射，默认使用 template.Content
	if !contentMappedToContent {
		// 检查是否有其他字段映射到 [CONTENT]
		hasContentMapping := false
		for _, targetKey := range mapping {
			if targetKey == "[CONTENT]" {
				hasContentMapping = true
				break
			}
		}
		if !hasContentMapping {
			// 默认：content 字段就是正文
			content = template.Content
		}
	}

	// 添加未映射的 metadata 字段到 frontmatter
	for key, val := range template.Metadata {
		// 检查这个字段是否已经在 frontmatter 中（通过 mapping）
		mapped := false
		for yamlKey, targetKey := range mapping {
			if yamlKey == key {
				if targetKey == "[CONTENT]" {
					// 这个字段映射到正文，跳过
					mapped = true
				} else {
					// 这个字段已经映射到 frontmatter，跳过
					mapped = true
				}
				break
			}
		}
		if !mapped {
			frontmatter[key] = val
		}
	}

	// 生成 frontmatter YAML
	var frontmatterYAML string
	if len(frontmatter) > 0 {
		fmBytes, err := yaml.Marshal(frontmatter)
		if err != nil {
			return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
		}
		frontmatterYAML = strings.TrimSpace(string(fmBytes))
	}

	// 组合结果
	var result strings.Builder
	if frontmatterYAML != "" {
		result.WriteString("---\n")
		result.WriteString(frontmatterYAML)
		result.WriteString("\n---\n\n")
	}
	result.WriteString(content)

	return result.String(), nil
}

// TOMLFormatter TOML 格式化器
type TOMLFormatter struct{}

// Format 将 CommandTemplate 格式化为 TOML
func (f *TOMLFormatter) Format(template *resources.CommandTemplate, mapping map[string]string) (string, error) {
	var result strings.Builder

	// 找到正文字段的映射（通常是 "prompt"）
	contentKey := ""
	for yamlKey, targetKey := range mapping {
		if targetKey == "prompt" {
			contentKey = yamlKey
			break
		}
	}
	if contentKey == "" {
		// 默认使用 content 字段
		contentKey = "content"
	}

	// 构建 TOML 键值对
	tomlFields := make(map[string]interface{})

	// 首先处理 content 字段（正文字段）
	if contentKey == "content" {
		// 找到 content 映射到的目标字段
		targetContentKey := "prompt" // 默认
		for yamlKey, targetKey := range mapping {
			if yamlKey == "content" {
				targetContentKey = targetKey
				break
			}
		}
		tomlFields[targetContentKey] = template.Content
	} else {
		// contentKey 是 metadata 中的某个字段
		if val, ok := template.Metadata[contentKey]; ok {
			if strVal, ok := val.(string); ok {
				targetKey := mapping[contentKey]
				if targetKey == "" {
					targetKey = "prompt"
				}
				tomlFields[targetKey] = strVal
			}
		}
	}

	// 处理其他映射的字段
	for yamlKey, targetKey := range mapping {
		if yamlKey == contentKey {
			continue // 已经处理
		}
		if yamlKey == "content" {
			continue // content 已经处理
		}
		if val, ok := template.Metadata[yamlKey]; ok {
			tomlFields[targetKey] = val
		}
	}

	// 添加未映射的 metadata 字段
	for key, val := range template.Metadata {
		if key == contentKey {
			continue
		}
		// 检查是否已经在 mapping 中
		mapped := false
		for yamlKey := range mapping {
			if yamlKey == key {
				mapped = true
				break
			}
		}
		if !mapped {
			tomlFields[key] = val
		}
	}

	// 生成 TOML 内容
	// 注意：这里使用简单的格式，复杂的 TOML 可能需要专门的库
	for key, val := range tomlFields {
		switch v := val.(type) {
		case string:
			// 转义字符串中的特殊字符
			escaped := escapeTOMLString(v)
			result.WriteString(fmt.Sprintf("%s = %s\n", key, escaped))
		case int, int64:
			result.WriteString(fmt.Sprintf("%s = %d\n", key, v))
		case float64:
			result.WriteString(fmt.Sprintf("%s = %f\n", key, v))
		case bool:
			result.WriteString(fmt.Sprintf("%s = %v\n", key, v))
		default:
			// 对于复杂类型，转换为字符串
			result.WriteString(fmt.Sprintf("%s = %q\n", key, fmt.Sprintf("%v", v)))
		}
	}

	return result.String(), nil
}

// escapeTOMLString 转义 TOML 字符串
func escapeTOMLString(s string) string {
	// 如果包含换行符或特殊字符，使用多行字符串
	if strings.Contains(s, "\n") || strings.Contains(s, `"`) || strings.Contains(s, `\`) {
		// 使用三引号多行字符串
		escaped := strings.ReplaceAll(s, `"""`, `\"\"\"`)
		return fmt.Sprintf(`"""%s"""`, escaped)
	}
	// 简单字符串，转义引号
	escaped := strings.ReplaceAll(s, `"`, `\"`)
	return fmt.Sprintf(`"%s"`, escaped)
}

// GetFormatter 根据文件格式获取格式化器
func GetFormatter(fileFormat string) (Formatter, error) {
	switch fileFormat {
	case "md":
		return &MDFormatter{}, nil
	case "toml":
		return &TOMLFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", fileFormat)
	}
}
