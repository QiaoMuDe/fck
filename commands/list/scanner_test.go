package list

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestNewFileScanner(t *testing.T) {
	scanner := NewFileScanner()

	if scanner == nil {
		t.Fatal("NewFileScanner() 返回 nil")
		return
	}

	if scanner.cache == nil {
		t.Error("NewFileScanner() 缓存未初始化")
	}
}

func TestFileScanner_Scan(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建测试文件和目录
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.txt")
	hiddenFile := filepath.Join(tempDir, ".hidden")
	testDir := filepath.Join(tempDir, "subdir")

	// 创建文件
	for _, file := range []string{testFile1, testFile2, hiddenFile} {
		f, err := os.Create(file)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("关闭测试文件失败: %v", err)
		}
	}

	// 创建目录
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	// 在子目录中创建文件
	subFile := filepath.Join(testDir, "subfile.txt")
	f, err := os.Create(subFile)
	if err != nil {
		t.Fatalf("创建子目录文件失败: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("关闭子目录文件失败: %v", err)
	}

	scanner := NewFileScanner()

	tests := []struct {
		name        string
		paths       []string
		opts        ScanOptions
		expectError bool
		minFiles    int
	}{
		{
			name:  "扫描单个目录",
			paths: []string{tempDir},
			opts: ScanOptions{
				Recursive:  false,
				ShowHidden: false,
				FileTypes:  nil,
				DirItself:  false,
			},
			expectError: false,
			minFiles:    3, // test1.txt, test2.txt, subdir (不包括隐藏文件)
		},
		{
			name:  "递归扫描",
			paths: []string{tempDir},
			opts: ScanOptions{
				Recursive:  true,
				ShowHidden: false,
				FileTypes:  nil,
				DirItself:  false,
			},
			expectError: false,
			minFiles:    4, // test1.txt, test2.txt, subdir, subfile.txt
		},
		{
			name:  "显示隐藏文件",
			paths: []string{tempDir},
			opts: ScanOptions{
				Recursive:  false,
				ShowHidden: true,
				FileTypes:  nil,
				DirItself:  false,
			},
			expectError: false,
			minFiles:    4, // test1.txt, test2.txt, .hidden, subdir
		},
		{
			name:  "只显示文件",
			paths: []string{tempDir},
			opts: ScanOptions{
				Recursive:  false,
				ShowHidden: false,
				FileTypes:  []string{types.FindTypeFile},
				DirItself:  false,
			},
			expectError: false,
			minFiles:    2, // test1.txt, test2.txt
		},
		{
			name:  "只显示目录",
			paths: []string{tempDir},
			opts: ScanOptions{
				Recursive:  false,
				ShowHidden: false,
				FileTypes:  []string{types.FindTypeDir},
				DirItself:  false,
			},
			expectError: false,
			minFiles:    1, // subdir
		},
		{
			name:  "显示目录本身",
			paths: []string{tempDir},
			opts: ScanOptions{
				Recursive:  false,
				ShowHidden: false,
				FileTypes:  nil,
				DirItself:  true,
			},
			expectError: false,
			minFiles:    1, // tempDir 本身
		},
		{
			name:        "扫描不存在的路径",
			paths:       []string{filepath.Join(tempDir, "nonexistent")},
			opts:        ScanOptions{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := scanner.Scan(tt.paths, tt.opts)

			if tt.expectError && err == nil {
				t.Errorf("Scan() 期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Scan() 意外错误 = %v", err)
			}

			if !tt.expectError && len(files) < tt.minFiles {
				t.Errorf("Scan() 返回文件数 = %v, 期望至少 %v", len(files), tt.minFiles)
			}
		})
	}
}

func TestFileScanner_GetFileInfo(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	if closeErr := f.Close(); closeErr != nil {
		t.Fatalf("关闭测试文件失败: %v", closeErr)
	}

	scanner := NewFileScanner()

	// 第一次调用 - 应该从文件系统获取
	info1, err := scanner.getFileInfo(testFile)
	if err != nil {
		t.Errorf("getFileInfo() 错误 = %v", err)
	}

	if info1 == nil {
		t.Error("getFileInfo() 返回 nil")
	}

	// 第二次调用 - 应该从缓存获取
	info2, err := scanner.getFileInfo(testFile)
	if err != nil {
		t.Errorf("getFileInfo() 缓存调用错误 = %v", err)
	}

	if info2 == nil {
		t.Error("getFileInfo() 缓存调用返回 nil")
	}

	// 验证缓存是否工作
	if info1.Name() != info2.Name() {
		t.Error("getFileInfo() 缓存未正确工作")
	}

	// 测试不存在的文件
	_, err = scanner.getFileInfo(filepath.Join(tempDir, "nonexistent.txt"))
	if err == nil {
		t.Error("getFileInfo() 对不存在文件应该返回错误")
	}
}

func TestFileScanner_ShouldSkipFile(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	normalFile := filepath.Join(tempDir, "normal.txt")
	hiddenFile := filepath.Join(tempDir, ".hidden")

	// 创建文件
	for _, file := range []string{normalFile, hiddenFile} {
		f, err := os.Create(file)
		if err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("关闭测试文件失败: %v", err)
		}
	}

	scanner := NewFileScanner()

	// 获取文件信息
	normalInfo, _ := os.Stat(normalFile)
	hiddenInfo, _ := os.Stat(hiddenFile)

	tests := []struct {
		name     string
		path     string
		isDir    bool
		fileInfo os.FileInfo
		isMain   bool
		opts     ScanOptions
		expected bool
	}{
		{
			name:     "普通文件，不显示隐藏文件",
			path:     normalFile,
			isDir:    false,
			fileInfo: normalInfo,
			isMain:   false,
			opts:     ScanOptions{ShowHidden: false},
			expected: false,
		},
		{
			name:     "隐藏文件，不显示隐藏文件",
			path:     hiddenFile,
			isDir:    false,
			fileInfo: hiddenInfo,
			isMain:   false,
			opts:     ScanOptions{ShowHidden: false},
			expected: true,
		},
		{
			name:     "隐藏文件，显示隐藏文件",
			path:     hiddenFile,
			isDir:    false,
			fileInfo: hiddenInfo,
			isMain:   false,
			opts:     ScanOptions{ShowHidden: true},
			expected: false,
		},
		{
			name:     "普通文件，只显示目录",
			path:     normalFile,
			isDir:    false,
			fileInfo: normalInfo,
			isMain:   false,
			opts:     ScanOptions{FileTypes: []string{types.FindTypeDir}},
			expected: true,
		},
		{
			name:     "目录，只显示文件",
			path:     tempDir,
			isDir:    true,
			fileInfo: normalInfo, // 这里用任意文件信息
			isMain:   false,
			opts:     ScanOptions{FileTypes: []string{types.FindTypeFile}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.shouldSkipFile(tt.path, tt.isDir, tt.fileInfo, tt.isMain, tt.opts)

			if result != tt.expected {
				t.Errorf("shouldSkipFile() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestFileScanner_BuildFileInfo(t *testing.T) {
	// 创建临时测试文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	if closeErr := f.Close(); closeErr != nil {
		t.Fatalf("关闭测试文件失败: %v", closeErr)
	}

	// 写入一些内容
	content := "test content"
	if writeErr := os.WriteFile(testFile, []byte(content), 0644); writeErr != nil {
		t.Fatalf("写入测试文件失败: %v", writeErr)
	}

	scanner := NewFileScanner()
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("获取文件信息失败: %v", err)
	}

	// 初始化命令标志
	InitListCmd()

	result := scanner.buildFileInfo(fileInfo, testFile, tempDir)

	// 验证结果
	if result.Name != "test.txt" {
		t.Errorf("buildFileInfo().Name = %v, 期望 %v", result.Name, "test.txt")
	}

	if result.Size != int64(len(content)) {
		t.Errorf("buildFileInfo().Size = %v, 期望 %v", result.Size, len(content))
	}

	if result.EntryType != types.FileType {
		t.Errorf("buildFileInfo().EntryType = %v, 期望 %v", result.EntryType, types.FileType)
	}

	if result.FileExt != ".txt" {
		t.Errorf("buildFileInfo().FileExt = %v, 期望 %v", result.FileExt, ".txt")
	}
}

func TestFileScanner_GetEntryType(t *testing.T) {
	// 创建临时测试文件和目录
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testDir := filepath.Join(tempDir, "subdir")
	emptyFile := filepath.Join(tempDir, "empty.txt")

	// 创建普通文件
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	if _, writeErr := f.WriteString("content"); writeErr != nil {
		t.Fatalf("写入测试文件失败: %v", writeErr)
	}
	if closeErr := f.Close(); closeErr != nil {
		t.Fatalf("关闭测试文件失败: %v", closeErr)
	}

	// 创建空文件
	f2, err := os.Create(emptyFile)
	if err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}
	if err := f2.Close(); err != nil {
		t.Fatalf("关闭空文件失败: %v", err)
	}

	// 创建目录
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}

	scanner := NewFileScanner()

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "普通文件",
			filePath: testFile,
			expected: types.FileType,
		},
		{
			name:     "空文件",
			filePath: emptyFile,
			expected: types.EmptyType,
		},
		{
			name:     "目录",
			filePath: testDir,
			expected: types.DirType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo, err := os.Stat(tt.filePath)
			if err != nil {
				t.Fatalf("获取文件信息失败: %v", err)
			}

			result := scanner.getEntryType(fileInfo)

			if result != tt.expected {
				t.Errorf("getEntryType() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

// 基准测试
func BenchmarkFileScanner_Scan(b *testing.B) {
	// 创建临时测试目录
	tempDir := b.TempDir()

	// 创建多个测试文件
	for i := 0; i < 100; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i))
		f, err := os.Create(testFile)
		if err != nil {
			b.Fatalf("创建测试文件失败: %v", err)
		}
		if err := f.Close(); err != nil {
			b.Fatalf("关闭测试文件失败: %v", err)
		}
	}

	scanner := NewFileScanner()
	opts := ScanOptions{
		Recursive:  false,
		ShowHidden: false,
		FileTypes:  nil,
		DirItself:  false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.Scan([]string{tempDir}, opts)
		if err != nil {
			b.Fatalf("扫描失败: %v", err)
		}
	}
}

func BenchmarkFileScanner_GetFileInfo(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// 创建测试文件
	f, err := os.Create(testFile)
	if err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}
	if err := f.Close(); err != nil {
		b.Fatalf("关闭测试文件失败: %v", err)
	}

	scanner := NewFileScanner()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.getFileInfo(testFile)
		if err != nil {
			b.Fatalf("获取文件信息失败: %v", err)
		}
	}
}
