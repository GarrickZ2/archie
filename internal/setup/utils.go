package setup

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

// GetEditor 获取用户配置的编辑器，优先使用 $EDITOR，默认为 vim
func GetEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	return editor
}

// OpenInEditor 在编辑器中打开文件
func OpenInEditor(filePath string) error {
	editor := GetEditor()
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsFileEmpty 检查文件是否为空或仅包含空白字符
func IsFileEmpty(fs afero.Fs, filePath string) (bool, error) {
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}

	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return false, err
	}

	// 去除所有空白字符后检查是否为空
	trimmed := strings.TrimSpace(string(content))
	return trimmed == "", nil
}

// ToKebabCase 将字符串转换为 kebab-case 格式
// 例如: "My Feature Name" -> "my-feature-name"
//       "myFeatureName" -> "my-feature-name"
func ToKebabCase(s string) string {
	// 先处理 camelCase: myFeatureName -> my-feature-name
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	s = re.ReplaceAllString(s, "${1}-${2}")

	// 转换为小写
	s = strings.ToLower(s)

	// 将所有非字母数字字符替换为连字符
	re = regexp.MustCompile("[^a-z0-9]+")
	s = re.ReplaceAllString(s, "-")

	// 移除首尾的连字符
	s = strings.Trim(s, "-")

	// 将多个连续的连字符替换为单个
	re = regexp.MustCompile("-+")
	s = re.ReplaceAllString(s, "-")

	return s
}
