package project

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitializer_Initialize_NewDirectory(t *testing.T) {
	// 使用内存文件系统
	fs := afero.NewMemMapFs()

	initializer := NewInitializer(&Config{
		FileSystem: fs,
	})

	// 初始化新项目
	err := initializer.Initialize("/test/project")
	assert.NoError(t, err, "初始化应该成功")

	// 验证关键文件是否创建
	tests := []string{
		"/test/project/architecture.md",
		"/test/project/background.md",
		"/test/project/api/api.md",
		"/test/project/rpc/rpc.md",
		"/test/project/storage/storage.md",
		"/test/project/storage/storage.sql",
		"/test/project/workflow/workflow.md",
		"/test/project/workflow/example/example.md",
		"/test/project/workflow/example/example.mmd",
		"/test/project/features.md",
		"/test/project/tasks.md",
		"/test/project/metrics.md",
		"/test/project/deployment.md",
		"/test/project/dependency.md",
		"/test/project/blocker.md",
	}

	for _, path := range tests {
		exists, err := afero.Exists(fs, path)
		assert.NoError(t, err, "检查文件存在性时不应出错: %s", path)
		assert.True(t, exists, "文件应该存在: %s", path)
	}

	// 验证目录是否创建（assets 标记为 __DIR__）
	exists, err := afero.DirExists(fs, "/test/project/assets")
	assert.NoError(t, err)
	assert.True(t, exists, "assets 目录应该存在")
}

func TestInitializer_Initialize_NonEmptyDirectory(t *testing.T) {
	fs := afero.NewMemMapFs()

	// 预先创建非空目录
	fs.MkdirAll("/test/project", 0755)
	afero.WriteFile(fs, "/test/project/existing.txt", []byte("test"), 0644)

	initializer := NewInitializer(&Config{
		FileSystem: fs,
	})

	// 初始化应该失败
	err := initializer.Initialize("/test/project")
	assert.Error(t, err, "非空目录初始化应该失败")
	assert.Contains(t, err.Error(), "not empty", "错误消息应该包含'not empty'")
}

func TestNewInitializer_WithNilConfig(t *testing.T) {
	// 测试传入 nil config 使用默认配置
	initializer := NewInitializer(nil)
	assert.NotNil(t, initializer, "应该返回非 nil 的初始化器")

	// 验证返回的类型
	_, ok := initializer.(*DefaultInitializer)
	assert.True(t, ok, "应该返回 DefaultInitializer 类型")
}

func TestNewInitializer_WithPartialConfig(t *testing.T) {
	fs := afero.NewMemMapFs()

	// 只提供部分配置
	initializer := NewInitializer(&Config{
		FileSystem: fs,
		// Validator 为 nil，应该使用默认值
	})

	assert.NotNil(t, initializer, "应该返回非 nil 的初始化器")

	// 验证可以正常初始化
	err := initializer.Initialize("/test/partial")
	assert.NoError(t, err, "使用部分配置初始化应该成功")
}
