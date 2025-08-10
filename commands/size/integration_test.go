package size

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/colorlib"
)

// 在测试开始前初始化命令
func init() {
	// 确保size命令被初始化
	if sizeCmd == nil {
		InitSizeCmd()
	}
}

// TestIntegrationSizeCmdMain 集成测试 - 测试主函数的完整流程
func TestIntegrationSizeCmdMain(t *testing.T) {
	// 跳过集成测试，除非明确要求运行
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 创建临时目录和文件
	tempDir := t.TempDir()

	// 创建测试文件
	testFiles := map[string]string{
		"small.txt":  "small content",
		"medium.txt": strings.Repeat("medium content ", 100),
		"large.txt":  strings.Repeat("large content ", 1000),
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
	}

	// 创建子目录
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("创建子目录失败: %v", err)
	}

	subFile := filepath.Join(subDir, "sub.txt")
	if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
		t.Fatalf("创建子文件失败: %v", err)
	}

	// 切换到临时目录
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(tempDir)

	// 创建ColorLib实例用于测试
	cl := colorlib.NewColorLib()

	// 注意：这里需要模拟命令行参数，但由于依赖全局变量sizeCmd，
	// 实际的集成测试可能需要重构代码以支持依赖注入

	t.Run("测试默认行为", func(t *testing.T) {
		// 由于SizeCmdMain依赖全局变量sizeCmd.Args()，
		// 这个测试需要额外的设置来模拟命令行参数
		t.Skip("需要重构以支持依赖注入")
	})

	// 测试辅助函数
	t.Run("测试expandPath集成", func(t *testing.T) {
		paths, err := expandPath("*.txt")
		if err != nil {
			t.Errorf("expandPath 失败: %v", err)
		}
		if len(paths) != 3 {
			t.Errorf("期望找到3个txt文件，实际找到%d个", len(paths))
		}
	})

	t.Run("测试addPathToList集成", func(t *testing.T) {
		var itemList items
		addPathToList("small.txt", &itemList, cl)

		if len(itemList) != 1 {
			t.Errorf("期望itemList长度为1，实际为%d", len(itemList))
		}

		if itemList[0].Name != "small.txt" {
			t.Errorf("期望文件名为small.txt，实际为%s", itemList[0].Name)
		}
	})

	// 重定向输出进行测试
	t.Run("测试printSizeTable输出", func(t *testing.T) {
		// 创建测试数据
		testItems := items{
			{Name: "test1.txt", Size: "1.2 KB"},
			{Name: "test2.txt", Size: "2.5 KB"},
		}

		// 由于printSizeTable直接输出到stdout，
		// 这里只能测试它不会panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("printSizeTable panic: %v", r)
			}
		}()

		printSizeTable(testItems, cl)
	})
}

// TestRealFileSystem 真实文件系统测试
func TestRealFileSystem(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过真实文件系统测试")
	}

	// 测试当前目录
	t.Run("当前目录", func(t *testing.T) {
		size, err := getPathSize(".")
		if err != nil {
			t.Errorf("获取当前目录大小失败: %v", err)
		}
		if size <= 0 {
			t.Errorf("当前目录大小应该大于0，实际为%d", size)
		}
		t.Logf("当前目录大小: %s", humanReadableSize(size, 2))
	})

	// 测试Go源文件
	t.Run("Go源文件", func(t *testing.T) {
		goFiles, err := expandPath("*.go")
		if err != nil {
			t.Logf("没有找到Go文件: %v", err)
			return
		}

		for _, file := range goFiles {
			size, err := getPathSize(file)
			if err != nil {
				t.Errorf("获取文件%s大小失败: %v", file, err)
				continue
			}
			t.Logf("文件%s大小: %s", file, humanReadableSize(size, 2))
		}
	})
}

// TestErrorHandling 错误处理测试
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{"不存在的文件", "/path/that/does/not/exist", true},
		{"空路径", "", true},
		{"无效路径字符", "\x00invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getPathSize(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("getPathSize(%s) 期望返回错误，但没有错误", tt.path)
				} else {
					t.Logf("getPathSize(%s) 正确返回错误: %v", tt.path, err)
				}
			} else {
				if err != nil {
					t.Errorf("getPathSize(%s) 返回意外错误: %v", tt.path, err)
				}
			}
		})
	}
}

// TestConcurrentAccess 并发访问测试
func TestConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过并发测试")
	}

	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "concurrent.txt")
	content := strings.Repeat("concurrent test content ", 1000)
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 并发测试
	const numGoroutines = 10
	const numIterations = 100

	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIterations; j++ {
				_, err := getPathSize(testFile)
				if err != nil {
					results <- err
					return
				}
			}
			results <- nil
		}()
	}

	// 收集结果
	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			t.Errorf("并发测试失败: %v", err)
		}
	}
}
