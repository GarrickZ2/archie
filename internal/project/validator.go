package project

import (
	"fmt"

	"github.com/spf13/afero"
)

// PathValidator 定义路径验证接口
type PathValidator interface {
	Validate(path string) error
}

// DefaultPathValidator 默认的路径验证器实现
type DefaultPathValidator struct {
	fs afero.Fs
}

// NewPathValidator 创建路径验证器
func NewPathValidator(fs afero.Fs) PathValidator {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	return &DefaultPathValidator{fs: fs}
}

// Validate 验证路径是否可用于初始化项目
// 路径不存在 -> 允许（后续会创建）
// 路径存在且为空目录 -> 允许
// 路径存在但非空 -> 拒绝
// 路径存在但不是目录 -> 拒绝
func (v *DefaultPathValidator) Validate(path string) error {
	exists, err := afero.DirExists(v.fs, path)
	if err != nil {
		return fmt.Errorf("检查路径失败: %w", err)
	}

	// 路径不存在，允许（后续会创建）
	if !exists {
		// 检查路径是否是文件
		fileExists, err := afero.Exists(v.fs, path)
		if err != nil {
			return fmt.Errorf("检查路径失败: %w", err)
		}
		if fileExists {
			return fmt.Errorf("路径已存在但不是目录: %s", path)
		}
		return nil
	}

	// 路径存在，检查是否为空
	entries, err := afero.ReadDir(v.fs, path)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("目录不为空: %s", path)
	}

	return nil
}
