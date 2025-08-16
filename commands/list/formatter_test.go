package list

import (
	"testing"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestNewFileFormatter(t *testing.T) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	if formatter == nil {
		t.Fatal("NewFileFormatter() 返回 nil")
	}

	if formatter.colorLib == nil {
		t.Error("NewFileFormatter() colorLib 未设置")
	}
}

func TestFileFormatter_GetSafeTerminalWidth(t *testing.T) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	width := formatter.getSafeTerminalWidth()

	// 验证返回的宽度在合理范围内
	if width < 20 || width > 500 {
		t.Errorf("getSafeTerminalWidth() = %v, 期望在 20-500 之间", width)
	}
}

func TestFileFormatter_HumanSize(t *testing.T) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	tests := []struct {
		name         string
		size         int64
		expectedSize string
		expectedUnit string
	}{
		{
			name:         "零字节",
			size:         0,
			expectedSize: "0",
			expectedUnit: "B",
		},
		{
			name:         "小于1KB",
			size:         512,
			expectedSize: "512",
			expectedUnit: "B",
		},
		{
			name:         "1KB",
			size:         1024,
			expectedSize: "1",
			expectedUnit: "KB",
		},
		{
			name:         "1.5KB",
			size:         1536,
			expectedSize: "1.5",
			expectedUnit: "KB",
		},
		{
			name:         "1MB",
			size:         1024 * 1024,
			expectedSize: "1",
			expectedUnit: "MB",
		},
		{
			name:         "1GB",
			size:         1024 * 1024 * 1024,
			expectedSize: "1",
			expectedUnit: "GB",
		},
		{
			name:         "1TB",
			size:         1024 * 1024 * 1024 * 1024,
			expectedSize: "1",
			expectedUnit: "TB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size, unit := formatter.humanSize(tt.size)

			if size != tt.expectedSize {
				t.Errorf("humanSize(%v) size = %v, 期望 %v", tt.size, size, tt.expectedSize)
			}

			if unit != tt.expectedUnit {
				t.Errorf("humanSize(%v) unit = %v, 期望 %v", tt.size, unit, tt.expectedUnit)
			}
		})
	}
}

func TestFileFormatter_PrepareFileNames(t *testing.T) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	files := FileInfoList{
		{
			Name:      "test1.txt",
			EntryType: types.FileType,
			Size:      100,
			ModTime:   time.Now(),
			FileExt:   ".txt",
		},
		{
			Name:      "test2.log",
			EntryType: types.FileType,
			Size:      200,
			ModTime:   time.Now(),
			FileExt:   ".log",
		},
	}

	tests := []struct {
		name     string
		files    FileInfoList
		opts     FormatOptions
		expected []string
	}{
		{
			name:  "普通文件名",
			files: files,
			opts: FormatOptions{
				QuoteNames: false,
				UseColor:   false,
			},
			expected: []string{"test1.txt", "test2.log"},
		},
		{
			name:  "引用文件名",
			files: files,
			opts: FormatOptions{
				QuoteNames: true,
				UseColor:   false,
			},
			expected: []string{"\"test1.txt\"", "\"test2.log\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.prepareFileNames(tt.files, tt.opts)

			if len(result) != len(tt.expected) {
				t.Errorf("prepareFileNames() 返回长度 = %v, 期望 %v", len(result), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("prepareFileNames()[%d] = %v, 期望 %v", i, result[i], expected)
				}
			}
		})
	}
}

func TestFileFormatter_GetMaxWidth(t *testing.T) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	tests := []struct {
		name      string
		fileNames []string
		expected  int
	}{
		{
			name:      "单个短文件名",
			fileNames: []string{"test.txt"},
			expected:  8, // "test.txt" 长度
		},
		{
			name:      "多个文件名",
			fileNames: []string{"short.txt", "very_long_filename.txt"},
			expected:  22, // "very_long_filename.txt" 长度
		},
		{
			name:      "空列表",
			fileNames: []string{},
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.getMaxWidth(tt.fileNames)

			if result != tt.expected {
				t.Errorf("getMaxWidth() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestFileFormatter_FormatPermissionString(t *testing.T) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	// 初始化命令标志
	InitListCmd()

	tests := []struct {
		name     string
		fileInfo FileInfo
		expected string
	}{
		{
			name: "普通文件权限",
			fileInfo: FileInfo{
				Perm: "-rw-r--r--",
			},
			expected: "rw-r--r--", // 去掉第一个字符
		},
		{
			name: "权限字符串太短",
			fileInfo: FileInfo{
				Perm: "short",
			},
			expected: "?-?-?",
		},
		{
			name: "目录权限",
			fileInfo: FileInfo{
				Perm: "drwxr-xr-x",
			},
			expected: "rwxr-xr-x", // 去掉第一个字符
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.formatPermissionString(tt.fileInfo)

			// 如果没有启用颜色，应该返回纯文本
			if len(result) == 0 {
				t.Errorf("formatPermissionString() 返回空字符串")
			}

			// 基本长度检查
			if tt.fileInfo.Perm == "short" && result != "?-?-?" {
				t.Errorf("formatPermissionString() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestFileFormatter_Render(t *testing.T) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	// 初始化命令标志
	InitListCmd()

	files := FileInfoList{
		{
			Name:      "test.txt",
			EntryType: types.FileType,
			Size:      100,
			ModTime:   time.Now(),
			Perm:      "-rw-r--r--",
			Owner:     "user",
			Group:     "group",
			FileExt:   ".txt",
		},
	}

	tests := []struct {
		name        string
		files       FileInfoList
		opts        FormatOptions
		expectError bool
	}{
		{
			name:  "网格格式渲染",
			files: files,
			opts: FormatOptions{
				LongFormat: false,
				UseColor:   false,
			},
			expectError: false,
		},
		{
			name:  "表格格式渲染",
			files: files,
			opts: FormatOptions{
				LongFormat: true,
				UseColor:   false,
			},
			expectError: false,
		},
		{
			name:        "空文件列表",
			files:       FileInfoList{},
			opts:        FormatOptions{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatter.Render(tt.files, tt.opts)

			if tt.expectError && err == nil {
				t.Errorf("Render() 期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Render() 意外错误 = %v", err)
			}
		})
	}
}

// 基准测试
func BenchmarkFileFormatter_HumanSize(b *testing.B) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	sizes := []int64{
		0,
		1024,
		1024 * 1024,
		1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, size := range sizes {
			formatter.humanSize(size)
		}
	}
}

func BenchmarkFileFormatter_PrepareFileNames(b *testing.B) {
	cl := colorlib.New()
	formatter := NewFileFormatter(cl)

	// 创建大量测试文件
	files := make(FileInfoList, 1000)
	now := time.Now()
	for i := 0; i < 1000; i++ {
		files[i] = FileInfo{
			Name:      "test_file_" + string(rune(i)) + ".txt",
			EntryType: types.FileType,
			Size:      int64(i * 100),
			ModTime:   now,
			FileExt:   ".txt",
		}
	}

	opts := FormatOptions{
		QuoteNames: false,
		UseColor:   false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatter.prepareFileNames(files, opts)
	}
}
