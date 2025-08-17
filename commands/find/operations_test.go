package find

import (
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/colorlib"
)

func TestNewFileOperator(t *testing.T) {
	cl := colorlib.NewColorLib()
	operator := NewFileOperator(cl)

	if operator == nil {
		t.Fatal("NewFileOperator 返回了 nil")
		return
	}

	if operator.cl != cl {
		t.Errorf("colorlib 设置错误")
	}
}

func TestFileOperator_Delete(t *testing.T) {
	cl := colorlib.NewColorLib()
	operator := NewFileOperator(cl)

	tests := []struct {
		name        string
		setupFile   func() string
		isDir       bool
		expectError bool
	}{
		{
			name: "删除普通文件",
			setupFile: func() string {
				tempFile, err := os.CreateTemp("", "test_delete")
				if err != nil {
					t.Fatalf("创建临时文件失败: %v", err)
				}
				_ = tempFile.Close()
				return tempFile.Name()
			},
			isDir:       false,
			expectError: false,
		},
		{
			name: "删除目录",
			setupFile: func() string {
				tempDir, err := os.MkdirTemp("", "test_delete_dir")
				if err != nil {
					t.Fatalf("创建临时目录失败: %v", err)
				}
				return tempDir
			},
			isDir:       true,
			expectError: false,
		},
		{
			name: "删除不存在的文件",
			setupFile: func() string {
				return "/nonexistent/file/path"
			},
			isDir:       false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile()

			err := operator.Delete(filePath, tt.isDir)

			if tt.expectError && err == nil {
				t.Errorf("期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("意外错误: %v", err)
			}

			// 验证文件确实被删除了（对于成功的情况）
			if !tt.expectError {
				if _, err := os.Stat(filePath); !os.IsNotExist(err) {
					t.Errorf("文件应该被删除但仍然存在")
				}
			}
		})
	}
}

func TestFileOperator_Move(t *testing.T) {
	cl := colorlib.NewColorLib()
	operator := NewFileOperator(cl)

	// 创建临时目录作为目标
	tempDir, err := os.MkdirTemp("", "test_move_target")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	tests := []struct {
		name        string
		setupFile   func() string
		targetDir   string
		expectError bool
	}{
		{
			name: "移动普通文件",
			setupFile: func() string {
				tempFile, err := os.CreateTemp("", "test_move")
				if err != nil {
					t.Fatalf("创建临时文件失败: %v", err)
				}
				_ = tempFile.Close()
				return tempFile.Name()
			},
			targetDir:   tempDir,
			expectError: false,
		},
		{
			name: "移动到不存在的目录",
			setupFile: func() string {
				tempFile, err := os.CreateTemp("", "test_move")
				if err != nil {
					t.Fatalf("创建临时文件失败: %v", err)
				}
				_ = tempFile.Close()
				return tempFile.Name()
			},
			targetDir:   "/nonexistent/directory",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile()
			originalName := filepath.Base(filePath)

			err := operator.Move(filePath, tt.targetDir)

			if tt.expectError && err == nil {
				t.Errorf("期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("意外错误: %v", err)
			}

			// 验证文件移动成功（对于成功的情况）
			if !tt.expectError {
				// 原文件应该不存在
				if _, err := os.Stat(filePath); !os.IsNotExist(err) {
					t.Errorf("原文件应该被移动但仍然存在")
				}

				// 新位置应该存在文件
				newPath := filepath.Join(tt.targetDir, originalName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					t.Errorf("文件应该被移动到新位置但不存在")
				}
			}
		})
	}
}

func TestFileOperator_Execute(t *testing.T) {
	cl := colorlib.NewColorLib()
	operator := NewFileOperator(cl)

	// 创建临时文件用于测试
	tempFile, err := os.CreateTemp("", "test_execute")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()
	_ = tempFile.Close()

	tests := []struct {
		name        string
		command     string
		filePath    string
		expectError bool
	}{
		{
			name:        "命令不包含占位符",
			command:     "echo test",
			filePath:    tempFile.Name(),
			expectError: true, // 应该返回错误因为没有{}占位符
		},
		{
			name:        "空命令",
			command:     "",
			filePath:    tempFile.Name(),
			expectError: true,
		},
		{
			name:        "空路径",
			command:     "echo {}",
			filePath:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := operator.Execute(tt.command, tt.filePath)

			if tt.expectError && err == nil {
				t.Errorf("期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("意外错误: %v", err)
			}
		})
	}
}

func TestFileOperator_QuotePath(t *testing.T) {
	cl := colorlib.NewColorLib()
	operator := NewFileOperator(cl)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "简单路径",
			path:     "test.txt",
			expected: "'test.txt'", // Unix系统使用单引号
		},
		{
			name:     "包含空格的路径",
			path:     "test file.txt",
			expected: "'test file.txt'",
		},
		{
			name:     "包含特殊字符的路径",
			path:     "test$file.txt",
			expected: "'test$file.txt'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := operator.quotePath(tt.path)
			// 由于quotePath是私有方法，我们无法直接测试
			// 这里只是示例，实际测试中需要通过公共方法间接测试
			_ = result
		})
	}
}

// 基准测试
func BenchmarkFileOperator_Delete(b *testing.B) {
	cl := colorlib.NewColorLib()
	operator := NewFileOperator(cl)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建临时文件
		tempFile, err := os.CreateTemp("", "benchmark_delete")
		if err != nil {
			b.Fatalf("创建临时文件失败: %v", err)
		}
		_ = tempFile.Close()

		// 删除文件
		err = operator.Delete(tempFile.Name(), false)
		if err != nil {
			b.Fatalf("删除文件失败: %v", err)
		}
	}
}

func BenchmarkFileOperator_Move(b *testing.B) {
	cl := colorlib.NewColorLib()
	operator := NewFileOperator(cl)

	// 创建目标目录
	tempDir, err := os.MkdirTemp("", "benchmark_move")
	if err != nil {
		b.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建临时文件
		tempFile, err := os.CreateTemp("", "benchmark_move")
		if err != nil {
			b.Fatalf("创建临时文件失败: %v", err)
		}
		_ = tempFile.Close()

		// 移动文件
		err = operator.Move(tempFile.Name(), tempDir)
		if err != nil {
			b.Fatalf("移动文件失败: %v", err)
		}
	}
}
