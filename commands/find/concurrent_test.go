package find

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestNewConcurrentSearcher(t *testing.T) {
	// 创建基础搜索器
	config := &types.FindConfig{
		Cl:         colorlib.NewColorLib(),
		MatchCount: &atomic.Int64{},
	}
	matcher := NewPatternMatcher(100)
	operator := NewFileOperator(config.Cl)
	searcher := NewFileSearcher(config, matcher, operator)

	tests := []struct {
		name            string
		maxWorkers      int
		expectedWorkers int
	}{
		{
			name:            "正常worker数量",
			maxWorkers:      4,
			expectedWorkers: 4,
		},
		{
			name:            "零worker数量",
			maxWorkers:      0,
			expectedWorkers: runtime.NumCPU(),
		},
		{
			name:            "负数worker数量",
			maxWorkers:      -1,
			expectedWorkers: runtime.NumCPU(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := NewConcurrentSearcher(searcher, tt.maxWorkers)

			if cs == nil {
				t.Errorf("NewConcurrentSearcher返回了nil")
				return
			}

			if cs.maxWorkers != tt.expectedWorkers {
				t.Errorf("worker数量错误，期望 %d，得到 %d", tt.expectedWorkers, cs.maxWorkers)
			}

			if cs.searcher != searcher {
				t.Errorf("搜索器设置错误")
			}
		})
	}
}

func TestConcurrentSearcher_SearchConcurrent(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "concurrent_search_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 创建测试文件结构
	testFiles := []string{
		"test1.go",
		"test2.txt",
		"subdir/test3.go",
		"subdir/test4.txt",
		"subdir/nested/test5.go",
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
		maxWorkers  int
		expectError bool
	}{
		{
			name: "基本并发搜索",
			setupConfig: func() *types.FindConfig {
				return &types.FindConfig{
					Cl:         colorlib.NewColorLib(),
					MatchCount: &atomic.Int64{},
				}
			},
			maxWorkers:  2,
			expectError: false,
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
			maxWorkers:  4,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.setupConfig()
			matcher := NewPatternMatcher(100)
			operator := NewFileOperator(config.Cl)
			searcher := NewFileSearcher(config, matcher, operator)
			concurrentSearcher := NewConcurrentSearcher(searcher, tt.maxWorkers)

			err := concurrentSearcher.SearchConcurrent(tempDir)

			if tt.expectError && err == nil {
				t.Errorf("期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("意外错误: %v", err)
			}
		})
	}
}

func TestConcurrentSearcher_Performance(t *testing.T) {
	// 创建大量测试文件的临时目录
	tempDir, err := os.MkdirTemp("", "performance_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 创建多个文件
	fileCount := 100
	for i := 0; i < fileCount; i++ {
		fileName := filepath.Join(tempDir, "file"+fmt.Sprintf("%d", i)+".txt")
		if err := os.WriteFile(fileName, []byte("test content"), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
	}

	config := &types.FindConfig{
		Cl:         colorlib.NewColorLib(),
		MatchCount: &atomic.Int64{},
	}
	matcher := NewPatternMatcher(100)
	operator := NewFileOperator(config.Cl)
	searcher := NewFileSearcher(config, matcher, operator)

	// 测试不同的worker数量
	workerCounts := []int{1, 2, 4, 8}

	for _, workers := range workerCounts {
		t.Run(fmt.Sprintf("workers_%d", workers), func(t *testing.T) {
			concurrentSearcher := NewConcurrentSearcher(searcher, workers)

			start := time.Now()
			err := concurrentSearcher.SearchConcurrent(tempDir)
			duration := time.Since(start)

			if err != nil {
				t.Errorf("搜索失败: %v", err)
			}

			t.Logf("Worker数量: %d, 耗时: %v", workers, duration)
		})
	}
}

// 基准测试
func BenchmarkConcurrentSearcher_SearchConcurrent(b *testing.B) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	if err != nil {
		b.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// 创建测试文件
	for i := 0; i < 50; i++ {
		fileName := filepath.Join(tempDir, fmt.Sprintf("file%d.txt", i))
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
	concurrentSearcher := NewConcurrentSearcher(searcher, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.MatchCount.Store(0) // 重置计数器
		err := concurrentSearcher.SearchConcurrent(tempDir)
		if err != nil {
			b.Fatalf("搜索失败: %v", err)
		}
	}
}
