package list

import (
	"fmt"
	"testing"
	"time"
)

func TestNewFileProcessor(t *testing.T) {
	processor := NewFileProcessor()

	if processor == nil {
		t.Error("NewFileProcessor() 返回 nil")
	}
}

func TestFileProcessor_Process(t *testing.T) {
	processor := NewFileProcessor()

	// 创建测试文件列表
	files := FileInfoList{
		{
			Name:      "file3.txt",
			Size:      300,
			ModTime:   time.Now().Add(-3 * time.Hour),
			EntryType: FileType,
		},
		{
			Name:      "file1.txt",
			Size:      100,
			ModTime:   time.Now().Add(-1 * time.Hour),
			EntryType: FileType,
		},
		{
			Name:      "file2.txt",
			Size:      200,
			ModTime:   time.Now().Add(-2 * time.Hour),
			EntryType: FileType,
		},
	}

	tests := []struct {
		name     string
		files    FileInfoList
		opts     ProcessOptions
		expected []string // 期望的文件名顺序
	}{
		{
			name:  "按名称排序",
			files: files,
			opts: ProcessOptions{
				SortBy:     "name",
				Reverse:    false,
				GroupByDir: false,
			},
			expected: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
		{
			name:  "按名称反向排序",
			files: files,
			opts: ProcessOptions{
				SortBy:     "name",
				Reverse:    true,
				GroupByDir: false,
			},
			expected: []string{"file3.txt", "file2.txt", "file1.txt"},
		},
		{
			name:  "按大小排序",
			files: files,
			opts: ProcessOptions{
				SortBy:     "size",
				Reverse:    false,
				GroupByDir: false,
			},
			expected: []string{"file3.txt", "file2.txt", "file1.txt"}, // 大到小
		},
		{
			name:  "按大小反向排序",
			files: files,
			opts: ProcessOptions{
				SortBy:     "size",
				Reverse:    true,
				GroupByDir: false,
			},
			expected: []string{"file1.txt", "file2.txt", "file3.txt"}, // 小到大
		},
		{
			name:  "按时间排序",
			files: files,
			opts: ProcessOptions{
				SortBy:     "time",
				Reverse:    false,
				GroupByDir: false,
			},
			expected: []string{"file1.txt", "file2.txt", "file3.txt"}, // 新到旧
		},
		{
			name:  "按时间反向排序",
			files: files,
			opts: ProcessOptions{
				SortBy:     "time",
				Reverse:    true,
				GroupByDir: false,
			},
			expected: []string{"file3.txt", "file2.txt", "file1.txt"}, // 旧到新
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.Process(tt.files, tt.opts)

			if len(result) != len(tt.expected) {
				t.Errorf("Process() 返回长度 = %v, 期望 %v", len(result), len(tt.expected))
				return
			}

			for i, expectedName := range tt.expected {
				if result[i].Name != expectedName {
					t.Errorf("Process()[%d].Name = %v, 期望 %v", i, result[i].Name, expectedName)
				}
			}
		})
	}
}

func TestFileProcessor_Sort(t *testing.T) {
	processor := NewFileProcessor()

	// 创建测试文件列表
	files := FileInfoList{
		{Name: "c.txt", Size: 300, ModTime: time.Now().Add(-1 * time.Hour)},
		{Name: "a.txt", Size: 100, ModTime: time.Now().Add(-3 * time.Hour)},
		{Name: "b.txt", Size: 200, ModTime: time.Now().Add(-2 * time.Hour)},
	}

	tests := []struct {
		name     string
		files    FileInfoList
		opts     ProcessOptions
		expected []string
	}{
		{
			name:  "默认排序（按名称）",
			files: files,
			opts: ProcessOptions{
				SortBy:  "unknown", // 未知排序方式，应该默认为名称
				Reverse: false,
			},
			expected: []string{"a.txt", "b.txt", "c.txt"},
		},
		{
			name:  "空文件列表",
			files: FileInfoList{},
			opts: ProcessOptions{
				SortBy:  "name",
				Reverse: false,
			},
			expected: []string{},
		},
		{
			name:  "单个文件",
			files: FileInfoList{{Name: "single.txt"}},
			opts: ProcessOptions{
				SortBy:  "name",
				Reverse: false,
			},
			expected: []string{"single.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.sort(tt.files, tt.opts)

			if len(result) != len(tt.expected) {
				t.Errorf("sort() 返回长度 = %v, 期望 %v", len(result), len(tt.expected))
				return
			}

			for i, expectedName := range tt.expected {
				if result[i].Name != expectedName {
					t.Errorf("sort()[%d].Name = %v, 期望 %v", i, result[i].Name, expectedName)
				}
			}
		})
	}
}

func TestFileProcessor_SortByName(t *testing.T) {
	processor := NewFileProcessor()

	files := FileInfoList{
		{Name: "Z.txt"},
		{Name: "a.txt"},
		{Name: "B.txt"},
		{Name: "c.txt"},
	}

	// 正向排序
	processor.sortByName(files, false)
	expected := []string{"a.txt", "B.txt", "c.txt", "Z.txt"}

	for i, expectedName := range expected {
		if files[i].Name != expectedName {
			t.Errorf("sortByName() 正向排序[%d].Name = %v, 期望 %v", i, files[i].Name, expectedName)
		}
	}

	// 反向排序
	processor.sortByName(files, true)
	expectedReverse := []string{"Z.txt", "c.txt", "B.txt", "a.txt"}

	for i, expectedName := range expectedReverse {
		if files[i].Name != expectedName {
			t.Errorf("sortByName() 反向排序[%d].Name = %v, 期望 %v", i, files[i].Name, expectedName)
		}
	}
}

func TestFileProcessor_SortBySize(t *testing.T) {
	processor := NewFileProcessor()

	files := FileInfoList{
		{Name: "small.txt", Size: 100},
		{Name: "large.txt", Size: 300},
		{Name: "medium.txt", Size: 200},
	}

	// 正向排序（大到小）
	processor.sortBySize(files, false)
	expected := []string{"large.txt", "medium.txt", "small.txt"}

	for i, expectedName := range expected {
		if files[i].Name != expectedName {
			t.Errorf("sortBySize() 正向排序[%d].Name = %v, 期望 %v", i, files[i].Name, expectedName)
		}
	}

	// 反向排序（小到大）
	processor.sortBySize(files, true)
	expectedReverse := []string{"small.txt", "medium.txt", "large.txt"}

	for i, expectedName := range expectedReverse {
		if files[i].Name != expectedName {
			t.Errorf("sortBySize() 反向排序[%d].Name = %v, 期望 %v", i, files[i].Name, expectedName)
		}
	}
}

func TestFileProcessor_SortByTime(t *testing.T) {
	processor := NewFileProcessor()

	now := time.Now()
	files := FileInfoList{
		{Name: "old.txt", ModTime: now.Add(-3 * time.Hour)},
		{Name: "new.txt", ModTime: now.Add(-1 * time.Hour)},
		{Name: "medium.txt", ModTime: now.Add(-2 * time.Hour)},
	}

	// 正向排序（新到旧）
	processor.sortByTime(files, false)
	expected := []string{"new.txt", "medium.txt", "old.txt"}

	for i, expectedName := range expected {
		if files[i].Name != expectedName {
			t.Errorf("sortByTime() 正向排序[%d].Name = %v, 期望 %v", i, files[i].Name, expectedName)
		}
	}

	// 反向排序（旧到新）
	processor.sortByTime(files, true)
	expectedReverse := []string{"old.txt", "medium.txt", "new.txt"}

	for i, expectedName := range expectedReverse {
		if files[i].Name != expectedName {
			t.Errorf("sortByTime() 反向排序[%d].Name = %v, 期望 %v", i, files[i].Name, expectedName)
		}
	}
}

// 基准测试
func BenchmarkFileProcessor_Process(b *testing.B) {
	processor := NewFileProcessor()

	// 创建大量测试文件
	files := make(FileInfoList, 1000)
	now := time.Now()
	for i := 0; i < 1000; i++ {
		files[i] = FileInfo{
			Name:    fmt.Sprintf("file%04d.txt", i),
			Size:    int64(i * 100),
			ModTime: now.Add(time.Duration(i) * time.Minute),
		}
	}

	opts := ProcessOptions{
		SortBy:     "name",
		Reverse:    false,
		GroupByDir: false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.Process(files, opts)
	}
}

func BenchmarkFileProcessor_SortByName(b *testing.B) {
	processor := NewFileProcessor()

	// 创建测试文件列表
	files := make(FileInfoList, 1000)
	for i := 0; i < 1000; i++ {
		files[i] = FileInfo{
			Name: fmt.Sprintf("file%04d.txt", 1000-i), // 逆序创建
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每次测试前重新创建副本
		testFiles := make(FileInfoList, len(files))
		copy(testFiles, files)
		processor.sortByName(testFiles, false)
	}
}

func BenchmarkFileProcessor_SortBySize(b *testing.B) {
	processor := NewFileProcessor()

	// 创建测试文件列表
	files := make(FileInfoList, 1000)
	for i := 0; i < 1000; i++ {
		files[i] = FileInfo{
			Name: fmt.Sprintf("file%04d.txt", i),
			Size: int64(1000 - i), // 逆序大小
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每次测试前重新创建副本
		testFiles := make(FileInfoList, len(files))
		copy(testFiles, files)
		processor.sortBySize(testFiles, false)
	}
}

func BenchmarkFileProcessor_SortByTime(b *testing.B) {
	processor := NewFileProcessor()

	// 创建测试文件列表
	files := make(FileInfoList, 1000)
	now := time.Now()
	for i := 0; i < 1000; i++ {
		files[i] = FileInfo{
			Name:    fmt.Sprintf("file%04d.txt", i),
			ModTime: now.Add(time.Duration(1000-i) * time.Minute), // 逆序时间
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每次测试前重新创建副本
		testFiles := make(FileInfoList, len(files))
		copy(testFiles, files)
		processor.sortByTime(testFiles, false)
	}
}
