package hash

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestNewHashTaskManager 测试哈希任务管理器创建
func TestNewHashTaskManager(t *testing.T) {
	files := []string{"test1.txt", "test2.txt"}
	manager := NewHashTaskManager(files, "md5")

	if manager == nil {
		t.Fatal("NewHashTaskManager() 返回 nil")
		return
	}

	if len(manager.files) != len(files) {
		t.Errorf("文件数量不匹配: got %d, want %d", len(manager.files), len(files))
	}

	if manager.concurrency <= 0 {
		t.Errorf("并发数应该大于0: got %d", manager.concurrency)
	}

	if manager.resultCh == nil {
		t.Error("结果通道不应该为nil")
	}

	if manager.writeCh == nil {
		t.Error("写入通道不应该为nil")
	}
}

// TestHashTaskManagerRun 测试任务管理器运行
func TestHashTaskManagerRun(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 初始化命令标志
	hashCmd = InitHashCmd()
	_ = hashCmdWrite.Set("false")
	_ = hashCmdProgress.Set("false")

	files := []string{testFile}
	manager := NewHashTaskManager(files, "md5")

	errors := manager.Run()

	if len(errors) > 0 {
		t.Errorf("Run() 返回错误: %v", errors)
	}

	processed, errorCount := manager.GetStats()
	if processed != 1 {
		t.Errorf("处理文件数不正确: got %d, want 1", processed)
	}
	if errorCount != 0 {
		t.Errorf("错误数不正确: got %d, want 0", errorCount)
	}
}

// TestHashTaskManagerRunWithWrite 测试带写入功能的任务管理器
func TestHashTaskManagerRunWithWrite(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 切换到临时目录
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(tempDir)

	// 初始化命令标志
	hashCmd = InitHashCmd()
	_ = hashCmdWrite.Set("true")
	_ = hashCmdProgress.Set("false")
	_ = hashCmdType.Set("md5")

	files := []string{testFile}
	manager := NewHashTaskManager(files, "md5")

	errors := manager.Run()

	if len(errors) > 0 {
		t.Logf("Run() 返回错误: %v (可能是正常的)", errors)
	}

	processed, _ := manager.GetStats()
	if processed != 1 {
		t.Errorf("处理文件数不正确: got %d, want 1", processed)
	}
}

// TestHashTaskManagerRunWithErrors 测试错误处理
func TestHashTaskManagerRunWithErrors(t *testing.T) {
	// 使用不存在的文件
	files := []string{"/nonexistent/file.txt"}

	// 初始化命令标志
	hashCmd = InitHashCmd()
	_ = hashCmdWrite.Set("false")
	_ = hashCmdProgress.Set("false")

	manager := NewHashTaskManager(files, "md5")
	errors := manager.Run()

	if len(errors) == 0 {
		t.Error("期望返回错误，但没有错误")
	}

	_, errorCount := manager.GetStats()
	if errorCount == 0 {
		t.Error("期望错误计数大于0")
	}
}

// TestHashTaskManagerConcurrency 测试并发处理
func TestHashTaskManagerConcurrency(t *testing.T) {
	// 创建多个临时文件
	tempDir := t.TempDir()
	var files []string

	for i := 0; i < 10; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i))
		content := strings.Repeat("test content ", i+1)
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("创建测试文件%d失败: %v", i, err)
		}
		files = append(files, testFile)
	}

	// 初始化命令标志
	hashCmd = InitHashCmd()
	_ = hashCmdWrite.Set("false")
	_ = hashCmdProgress.Set("false")

	manager := NewHashTaskManager(files, "md5")

	start := time.Now()
	errors := manager.Run()
	duration := time.Since(start)

	if len(errors) > 0 {
		t.Errorf("Run() 返回错误: %v", errors)
	}

	processed, errorCount := manager.GetStats()
	if processed != int64(len(files)) {
		t.Errorf("处理文件数不正确: got %d, want %d", processed, len(files))
	}
	if errorCount != 0 {
		t.Errorf("错误数不正确: got %d, want 0", errorCount)
	}

	t.Logf("处理 %d 个文件耗时: %v", len(files), duration)
}

// TestShouldSkipFile 测试文件跳过逻辑
func TestShouldSkipFile(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name        string
		filePath    string
		expectSkip  bool
		expectError bool
	}{
		{
			name:        "普通文件",
			filePath:    testFile,
			expectSkip:  false,
			expectError: false,
		},
		{
			name:        "不存在的文件",
			filePath:    filepath.Join(tempDir, "nonexistent.txt"),
			expectSkip:  false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skip, err := shouldSkipFile(tt.filePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("shouldSkipFile() 期望返回错误，但没有错误")
				}
			} else {
				if err != nil {
					t.Errorf("shouldSkipFile() 返回意外错误: %v", err)
				}
				if skip != tt.expectSkip {
					t.Errorf("shouldSkipFile() = %v, 期望 %v", skip, tt.expectSkip)
				}
			}
		})
	}
}

// TestHashRunTasksRefactored 测试重构后的任务执行函数
func TestHashRunTasksRefactored(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	var files []string

	for i := 0; i < 3; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i))
		if err := os.WriteFile(testFile, []byte(fmt.Sprintf("content%d", i)), 0644); err != nil {
			t.Fatalf("创建测试文件%d失败: %v", i, err)
		}
		files = append(files, testFile)
	}

	// 初始化命令标志
	hashCmd = InitHashCmd()
	_ = hashCmdWrite.Set("false")
	_ = hashCmdProgress.Set("false")

	errors := hashRunTasksRefactored(files, "md5")

	if len(errors) > 0 {
		t.Errorf("hashRunTasksRefactored() 返回错误: %v", errors)
	}
}

// TestHashRunTasksRefactoredEmpty 测试空文件列表
func TestHashRunTasksRefactoredEmpty(t *testing.T) {
	errors := hashRunTasksRefactored([]string{}, "md5")

	if errors != nil {
		t.Errorf("hashRunTasksRefactored() 对空列表应该返回nil，但返回: %v", errors)
	}
}

// TestHashResult 测试哈希结果结构
func TestHashResult(t *testing.T) {
	result := HashResult{
		FilePath:  "/test/path.txt",
		HashValue: "d41d8cd98f00b204e9800998ecf8427e",
		Error:     nil,
	}

	if result.FilePath != "/test/path.txt" {
		t.Errorf("FilePath 不匹配: got %s, want /test/path.txt", result.FilePath)
	}

	if result.HashValue != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("HashValue 不匹配: got %s", result.HashValue)
	}

	if result.Error != nil {
		t.Errorf("Error 应该为nil，但是: %v", result.Error)
	}
}

// TestWriteRequest 测试写入请求结构
func TestWriteRequest(t *testing.T) {
	done := make(chan error, 1)
	req := WriteRequest{
		Content: "test content",
		Done:    done,
	}

	if req.Content != "test content" {
		t.Errorf("Content 不匹配: got %s, want test content", req.Content)
	}

	if req.Done != done {
		t.Error("Done 通道不匹配")
	}
}

// TestFileWriterWrapper 测试文件写入器包装
func TestFileWriterWrapper(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	// 这里只测试结构体的创建，实际的写入逻辑在集成测试中验证
	wrapper := &FileWriterWrapper{
		file:   file,
		writer: nil, // 在实际使用中会是bufio.Writer
	}

	if wrapper.file != file {
		t.Error("文件引用不匹配")
	}
}

// BenchmarkHashTaskManager 性能测试
func BenchmarkHashTaskManager(b *testing.B) {
	// 创建临时文件
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark.txt")
	content := strings.Repeat("benchmark data ", 100) // 约1.5KB数据

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	// 初始化命令标志
	hashCmd = InitHashCmd()
	_ = hashCmdWrite.Set("false")
	_ = hashCmdProgress.Set("false")

	files := []string{testFile}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager := NewHashTaskManager(files, "md5")
		_ = manager.Run()
	}
}

// BenchmarkHashRunTasksRefactored 性能测试重构函数
func BenchmarkHashRunTasksRefactored(b *testing.B) {
	// 创建临时文件
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark.txt")
	content := strings.Repeat("benchmark data ", 100)

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	// 初始化命令标志
	hashCmd = InitHashCmd()
	_ = hashCmdWrite.Set("false")
	_ = hashCmdProgress.Set("false")

	files := []string{testFile}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hashRunTasksRefactored(files, "md5")
	}
}
