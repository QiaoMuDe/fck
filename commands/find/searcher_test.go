package find

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestNewFileSearcher(t *testing.T) {
	config := &types.FindConfig{
		Cl:         colorlib.NewColorLib(),
		MatchCount: &atomic.Int64{},
	}
	matcher := NewPatternMatcher(100)
	operator := NewFileOperator(config.Cl)

	searcher := NewFileSearcher(config, matcher, operator)

	if searcher == nil {
		t.Fatal("NewFileSearcher 返回了 nil")
		return
	}

	if searcher.config != config {
		t.Errorf("配置设置错误")
	}

	if searcher.matcher != matcher {
		t.Errorf("匹配器设置错误")
	}

	if searcher.operator != operator {
		t.Errorf("操作器设置错误")
	}
}

func TestFileSearcher_Search(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "searcher_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 创建测试文件
	testFiles := []string{
		"test.go",
		"main.go",
		"readme.txt",
		"config.json",
		"subdir/nested.go",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)

		// 创建目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败: %v", err)
		}

		// 创建文件
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
	}

	tests := []struct {
		name        string
		setupConfig func() *types.FindConfig
		expectError bool
		expectCount int64
	}{
		{
			name: "搜索所有文件",
			setupConfig: func() *types.FindConfig {
				return &types.FindConfig{
					Cl:         colorlib.NewColorLib(),
					MatchCount: &atomic.Int64{},
				}
			},
			expectError: false,
			expectCount: 0, // 默认配置下可能不匹配任何文件，因为没有设置匹配条件
		},
		{
			name: "搜索Go文件",
			setupConfig: func() *types.FindConfig {
				config := &types.FindConfig{
					Cl:         colorlib.NewColorLib(),
					MatchCount: &atomic.Int64{},
				}
				config.FindExtSliceMap.Store(".go", true)
				return config
			},
			expectError: false,
			expectCount: 3, // 应该找到3个.go文件
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.setupConfig()
			matcher := NewPatternMatcher(100)
			operator := NewFileOperator(config.Cl)
			searcher := NewFileSearcher(config, matcher, operator)

			err := searcher.Search(tempDir)

			if tt.expectError && err == nil {
				t.Errorf("期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("意外错误: %v", err)
			}

			if !tt.expectError {
				actualCount := config.MatchCount.Load()
				if actualCount != tt.expectCount {
					t.Errorf("匹配数量错误，期望 %d，得到 %d", tt.expectCount, actualCount)
				}
			}
		})
	}
}

func TestFileSearcher_NeedFileInfo(t *testing.T) {
	initTestFlags() // 确保标志变量已初始化
	config := &types.FindConfig{
		Cl:         colorlib.NewColorLib(),
		MatchCount: &atomic.Int64{},
	}
	matcher := NewPatternMatcher(100)
	operator := NewFileOperator(config.Cl)
	searcher := NewFileSearcher(config, matcher, operator)

	// 测试默认情况下不需要文件信息
	if searcher.needFileInfo() {
		t.Errorf("默认情况下不应该需要文件信息")
	}

	// 模拟设置需要文件信息的标志
	_ = findCmdSize.Set("+100k")
	defer func() { _ = findCmdSize.Set("") }() // 清理

	if !searcher.needFileInfo() {
		t.Errorf("设置大小过滤后应该需要文件信息")
	}
}

// 基准测试
func BenchmarkFileSearcher_Search(b *testing.B) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "benchmark_searcher")
	if err != nil {
		b.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 创建测试文件
	for i := 0; i < 100; i++ {
		fileName := filepath.Join(tempDir, "file"+fmt.Sprintf("%d", i)+".txt")
		if err := os.WriteFile(fileName, []byte("test content"), 0644); err != nil {
			b.Fatalf("创建测试文件失败: %v", err)
		}
	}

	config := &types.FindConfig{
		Cl:         colorlib.NewColorLib(),
		MatchCount: &atomic.Int64{},
	}
	matcher := NewPatternMatcher(100)
	operator := NewFileOperator(config.Cl)
	searcher := NewFileSearcher(config, matcher, operator)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.MatchCount.Store(0) // 重置计数器
		err := searcher.Search(tempDir)
		if err != nil {
			b.Fatalf("搜索失败: %v", err)
		}
	}
}
