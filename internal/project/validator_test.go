package project

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestPathValidator_Validate_NonExistentPath(t *testing.T) {
	// 使用内存文件系统进行测试
	fs := afero.NewMemMapFs()
	validator := NewPathValidator(fs)

	// 测试不存在的路径（应该通过）
	err := validator.Validate("/test/project")
	assert.NoError(t, err, "不存在的路径应该通过验证")
}

func TestPathValidator_Validate_EmptyDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewPathValidator(fs)

	// 创建空目录
	fs.MkdirAll("/test/empty", 0755)

	// 测试空目录（应该通过）
	err := validator.Validate("/test/empty")
	assert.NoError(t, err, "空目录应该通过验证")
}

func TestPathValidator_Validate_NonEmptyDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewPathValidator(fs)

	// 创建非空目录
	fs.MkdirAll("/test/nonempty", 0755)
	afero.WriteFile(fs, "/test/nonempty/file.txt", []byte("content"), 0644)

	// 测试非空目录（应该失败）
	err := validator.Validate("/test/nonempty")
	assert.Error(t, err, "非空目录应该验证失败")
	assert.Contains(t, err.Error(), "不为空", "错误消息应该包含'不为空'")
}

func TestPathValidator_Validate_PathIsFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	validator := NewPathValidator(fs)

	// 创建文件（而不是目录）
	fs.MkdirAll("/test", 0755)
	afero.WriteFile(fs, "/test/file.txt", []byte("content"), 0644)

	// 测试路径是文件（应该失败）
	err := validator.Validate("/test/file.txt")
	assert.Error(t, err, "路径是文件应该验证失败")
	assert.Contains(t, err.Error(), "不是目录", "错误消息应该包含'不是目录'")
}

func TestNewPathValidator_WithNilFS(t *testing.T) {
	// 测试传入 nil 会使用默认的 OsFs
	validator := NewPathValidator(nil)
	assert.NotNil(t, validator, "应该返回非 nil 的验证器")

	// 验证返回的类型
	_, ok := validator.(*DefaultPathValidator)
	assert.True(t, ok, "应该返回 DefaultPathValidator 类型")
}
