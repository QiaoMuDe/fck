package check

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestFileChecker_CheckFiles(t *testing.T) {
	cl := colorlib.New()
	checker := newFileChecker(cl, md5.New)

	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.txt")
	testFile3 := filepath.Join(tempDir, "nonexistent.txt")

	content1 := "test content 1"
	content2 := "test content 2"

	err := os.WriteFile(testFile1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	err = os.WriteFile(testFile2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 计算正确的哈希值
	hash1 := fmt.Sprintf("%x", md5.Sum([]byte(content1)))
	hash2 := fmt.Sprintf("%x", md5.Sum([]byte(content2)))
	wrongHash := "00000000000000000000000000000000"

	tests := []struct {
		name        string
		hashMap     types.VirtualHashMap
		expectError bool
	}{
		{
			name: "所有文件匹配",
			hashMap: types.VirtualHashMap{
				testFile1: types.VirtualHashEntry{
					RealPath: testFile1,
					Hash:     hash1,
				},
				testFile2: types.VirtualHashEntry{
					RealPath: testFile2,
					Hash:     hash2,
				},
			},
			expectError: false,
		},
		{
			name: "部分文件不匹配",
			hashMap: types.VirtualHashMap{
				testFile1: types.VirtualHashEntry{
					RealPath: testFile1,
					Hash:     hash1,
				},
				testFile2: types.VirtualHashEntry{
					RealPath: testFile2,
					Hash:     wrongHash,
				},
			},
			expectError: false,
		},
		{
			name: "文件不存在",
			hashMap: types.VirtualHashMap{
				testFile3: types.VirtualHashEntry{
					RealPath: testFile3,
					Hash:     hash1,
				},
			},
			expectError: false,
		},
		{
			name:        "空的哈希映射",
			hashMap:     types.VirtualHashMap{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checker.checkFiles(tt.hashMap)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有发生错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误但发生了错误: %v", err)
			}
		})
	}
}

func TestFileChecker_Worker(t *testing.T) {
	cl := colorlib.New()
	checker := newFileChecker(cl, md5.New)

	// 创建临时测试文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "worker_test.txt")
	content := "worker test content"

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 计算正确的哈希值
	expectedHash := fmt.Sprintf("%x", md5.Sum([]byte(content)))

	// 创建通道
	jobs := make(chan types.VirtualHashEntry, 1)
	results := make(chan checkResult, 1)

	// 发送任务
	jobs <- types.VirtualHashEntry{
		RealPath: testFile,
		Hash:     expectedHash,
	}
	close(jobs)

	// 启动worker
	var wg sync.WaitGroup
	wg.Add(1)
	go checker.worker(jobs, results, &wg)

	// 等待完成
	wg.Wait()
	close(results)

	// 检查结果
	result := <-results
	if result.err != nil {
		t.Errorf("worker处理出错: %v", result.err)
	}

	if result.filePath != testFile {
		t.Errorf("文件路径不匹配，期望: %s, 实际: %s", testFile, result.filePath)
	}

	if result.expectedHash != expectedHash {
		t.Errorf("期望哈希不匹配，期望: %s, 实际: %s", expectedHash, result.expectedHash)
	}

	if result.actualHash != expectedHash {
		t.Errorf("实际哈希不匹配，期望: %s, 实际: %s", expectedHash, result.actualHash)
	}
}

func TestFileChecker_WorkerWithNonexistentFile(t *testing.T) {
	cl := colorlib.New()
	checker := newFileChecker(cl, md5.New)

	// 创建通道
	jobs := make(chan types.VirtualHashEntry, 1)
	results := make(chan checkResult, 1)

	// 发送不存在文件的任务
	jobs <- types.VirtualHashEntry{
		RealPath: "nonexistent_file.txt",
		Hash:     "somehash",
	}
	close(jobs)

	// 启动worker
	var wg sync.WaitGroup
	wg.Add(1)
	go checker.worker(jobs, results, &wg)

	// 等待完成
	wg.Wait()
	close(results)

	// 检查结果数量 - 不存在的文件应该被跳过，不会发送到results通道
	resultCount := 0
	for range results {
		resultCount++
	}

	if resultCount != 0 {
		t.Errorf("不存在的文件应该被跳过，但收到了%d个结果", resultCount)
	}
}

func TestFileChecker_PrintSummary(t *testing.T) {
	cl := colorlib.New()
	checker := newFileChecker(cl, md5.New)

	tests := []struct {
		name       string
		processed  int
		mismatched int
		errors     int
		total      int
	}{
		{
			name:       "全部通过",
			processed:  5,
			mismatched: 0,
			errors:     0,
			total:      5,
		},
		{
			name:       "部分失败",
			processed:  5,
			mismatched: 2,
			errors:     1,
			total:      5,
		},
		{
			name:       "全部失败",
			processed:  3,
			mismatched: 2,
			errors:     1,
			total:      3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这个测试主要验证函数不会panic
			// 实际的输出验证需要捕获stdout，这里简化处理
			checker.printSummary(tt.processed, tt.mismatched, tt.errors, tt.total)
		})
	}
}

func TestNewFileChecker(t *testing.T) {
	cl := colorlib.New()
	checker := newFileChecker(cl, md5.New)

	if checker == nil {
		t.Errorf("NewFileChecker 返回了 nil")
		return
	}

	if checker.cl != cl {
		t.Errorf("ColorLib 不匹配")
	}

	if checker.hashFunc == nil {
		t.Errorf("哈希函数不应该为 nil")
	}

	if checker.maxWorkers <= 0 {
		t.Errorf("最大工作线程数应该大于 0")
	}
}

func TestFileChecker_CollectResults(t *testing.T) {
	cl := colorlib.New()
	checker := newFileChecker(cl, md5.New)

	// 创建结果通道
	results := make(chan checkResult, 3)

	// 发送测试结果
	results <- checkResult{
		filePath:     "test1.txt",
		expectedHash: "hash1",
		actualHash:   "hash1",
		err:          nil,
	}
	results <- checkResult{
		filePath:     "test2.txt",
		expectedHash: "hash2",
		actualHash:   "different_hash",
		err:          nil,
	}
	results <- checkResult{
		filePath:     "test3.txt",
		expectedHash: "hash3",
		actualHash:   "",
		err:          fmt.Errorf("计算哈希失败"),
	}
	close(results)

	// 收集结果
	err := checker.collectResults(results, 3)
	if err != nil {
		t.Errorf("collectResults返回了错误: %v", err)
	}
}
