package size

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestHelper 测试辅助结构体
type TestHelper struct {
	t       *testing.T
	tempDir string
	cleanup []func()
}

// NewTestHelper 创建新的测试辅助器
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t:       t,
		tempDir: t.TempDir(),
		cleanup: make([]func(), 0),
	}
}

// CreateFile 创建测试文件
func (h *TestHelper) CreateFile(name string, content string) string {
	filePath := filepath.Join(h.tempDir, name)

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		h.t.Fatalf("创建目录失败: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		h.t.Fatalf("创建文件%s失败: %v", name, err)
	}

	return filePath
}

// CreateDir 创建测试目录
func (h *TestHelper) CreateDir(name string) string {
	dirPath := filepath.Join(h.tempDir, name)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		h.t.Fatalf("创建目录%s失败: %v", name, err)
	}
	return dirPath
}

// CreateFileWithSize 创建指定大小的测试文件
func (h *TestHelper) CreateFileWithSize(name string, size int) string {
	content := strings.Repeat("a", size)
	return h.CreateFile(name, content)
}

// CreateDirectoryStructure 创建复杂的目录结构
func (h *TestHelper) CreateDirectoryStructure() map[string]string {
	structure := map[string]string{
		"root.txt":            "root file content",
		"dir1/file1.txt":      "file1 content",
		"dir1/file2.txt":      "file2 content",
		"dir1/subdir/sub.txt": "sub file content",
		"dir2/large.txt":      strings.Repeat("large content ", 1000),
		"dir2/small.txt":      "small",
		".hidden/hidden.txt":  "hidden content",
	}

	paths := make(map[string]string)
	for relativePath, content := range structure {
		fullPath := h.CreateFile(relativePath, content)
		paths[relativePath] = fullPath
	}

	return paths
}

// GetTempDir 获取临时目录路径
func (h *TestHelper) GetTempDir() string {
	return h.tempDir
}

// AddCleanup 添加清理函数
func (h *TestHelper) AddCleanup(fn func()) {
	h.cleanup = append(h.cleanup, fn)
}

// Cleanup 执行所有清理函数
func (h *TestHelper) Cleanup() {
	for _, fn := range h.cleanup {
		fn()
	}
}

// AssertFileExists 断言文件存在
func (h *TestHelper) AssertFileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		h.t.Errorf("文件不存在: %s", path)
	}
}

// AssertFileSize 断言文件大小
func (h *TestHelper) AssertFileSize(path string, expectedSize int64) {
	info, err := os.Stat(path)
	if err != nil {
		h.t.Errorf("获取文件信息失败: %v", err)
		return
	}

	if info.Size() != expectedSize {
		h.t.Errorf("文件大小不匹配: 期望 %d, 实际 %d", expectedSize, info.Size())
	}
}

// TestWithHelper 使用测试辅助器的示例测试
func TestWithHelper(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// 创建测试文件
	file1 := helper.CreateFile("test1.txt", "Hello World")
	file2 := helper.CreateFileWithSize("test2.txt", 1024)

	// 创建目录结构
	structure := helper.CreateDirectoryStructure()

	// 测试文件是否正确创建
	helper.AssertFileExists(file1)
	helper.AssertFileExists(file2)
	helper.AssertFileSize(file2, 1024)

	// 测试目录结构
	for _, path := range structure {
		helper.AssertFileExists(path)
	}

	// 测试 getPathSize 函数
	size1, err := getPathSize(file1)
	if err != nil {
		t.Errorf("getPathSize 失败: %v", err)
	}
	if size1 != 11 { // "Hello World" = 11 bytes
		t.Errorf("文件大小不正确: 期望 11, 实际 %d", size1)
	}

	// 测试目录大小计算
	dirSize, err := getPathSize(helper.GetTempDir())
	if err != nil {
		t.Errorf("计算目录大小失败: %v", err)
	}
	if dirSize <= 0 {
		t.Errorf("目录大小应该大于0, 实际 %d", dirSize)
	}

	t.Logf("临时目录大小: %s", humanReadableSize(dirSize, 2))
}

// BenchmarkTestHelper 测试辅助器性能测试
func BenchmarkTestHelper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		helper := NewTestHelper(&testing.T{})
		helper.CreateFile("bench.txt", "benchmark content")
		helper.CreateFileWithSize("large.txt", 10240)
		// 注意：在基准测试中通常不需要调用 Cleanup，因为使用的是 t.TempDir()
	}
}

// TestHelperConcurrency 测试辅助器并发安全性
func TestHelperConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过并发测试")
	}

	helper := NewTestHelper(t)
	defer helper.Cleanup()

	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					results <- fmt.Errorf("goroutine %d panic: %v", id, r)
					return
				}
				results <- nil
			}()

			// 每个goroutine创建不同的文件
			filename := fmt.Sprintf("concurrent_%d.txt", id)
			content := fmt.Sprintf("content from goroutine %d", id)

			path := helper.CreateFile(filename, content)
			helper.AssertFileExists(path)

			// 测试文件大小计算
			size, err := getPathSize(path)
			if err != nil {
				results <- fmt.Errorf("goroutine %d: getPathSize 失败: %v", id, err)
				return
			}

			if size != int64(len(content)) {
				results <- fmt.Errorf("goroutine %d: 大小不匹配", id)
				return
			}
		}(i)
	}

	// 收集结果
	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			t.Errorf("并发测试失败: %v", err)
		}
	}
}
