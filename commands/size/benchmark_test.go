package size

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// BenchmarkHumanReadableSizeVariousSizes 测试不同大小的格式化性能
func BenchmarkHumanReadableSizeVariousSizes(b *testing.B) {
	testCases := []struct {
		name string
		size int64
	}{
		{"Bytes", 512},
		{"KB", 1536},
		{"MB", 1572864},
		{"GB", 1610612736},
		{"TB", 1649267441664},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				humanReadableSize(tc.size, 2)
			}
		})
	}
}

// BenchmarkExpandPath 测试路径展开性能
func BenchmarkExpandPath(b *testing.B) {
	// 创建临时目录和多个文件
	tempDir := b.TempDir()

	// 创建100个测试文件
	for i := 0; i < 100; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("file%03d.txt", i))
		if err := os.WriteFile(filename, []byte("test"), 0644); err != nil {
			b.Fatalf("创建测试文件失败: %v", err)
		}
	}

	// 切换到临时目录
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := expandPath("*.txt")
		if err != nil {
			b.Fatalf("expandPath 失败: %v", err)
		}
	}
}

// BenchmarkGetPathSizeFile 测试单个文件大小计算性能
func BenchmarkGetPathSizeFile(b *testing.B) {
	tempDir := b.TempDir()

	// 创建不同大小的测试文件
	testCases := []struct {
		name string
		size int
	}{
		{"1KB", 1024},
		{"10KB", 10240},
		{"100KB", 102400},
		{"1MB", 1048576},
	}

	for _, tc := range testCases {
		filename := filepath.Join(tempDir, tc.name+".txt")
		content := strings.Repeat("a", tc.size)
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			b.Fatalf("创建测试文件失败: %v", err)
		}

		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := getPathSize(filename)
				if err != nil {
					b.Fatalf("getPathSize 失败: %v", err)
				}
			}
		})
	}
}

// BenchmarkGetPathSizeDirectory 测试目录大小计算性能
func BenchmarkGetPathSizeDirectory(b *testing.B) {
	tempDir := b.TempDir()

	// 创建包含多个文件的目录结构
	testCases := []struct {
		name      string
		fileCount int
		fileSize  int
	}{
		{"10Files_1KB", 10, 1024},
		{"50Files_1KB", 50, 1024},
		{"100Files_1KB", 100, 1024},
		{"10Files_10KB", 10, 10240},
	}

	for _, tc := range testCases {
		subDir := filepath.Join(tempDir, tc.name)
		if err := os.Mkdir(subDir, 0755); err != nil {
			b.Fatalf("创建子目录失败: %v", err)
		}

		content := strings.Repeat("a", tc.fileSize)
		for i := 0; i < tc.fileCount; i++ {
			filename := filepath.Join(subDir, fmt.Sprintf("file%03d.txt", i))
			if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
				b.Fatalf("创建测试文件失败: %v", err)
			}
		}

		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := getPathSize(subDir)
				if err != nil {
					b.Fatalf("getPathSize 失败: %v", err)
				}
			}
		})
	}
}

// BenchmarkGetPathSizeDeepDirectory 测试深层目录结构性能
func BenchmarkGetPathSizeDeepDirectory(b *testing.B) {
	tempDir := b.TempDir()

	// 创建深层目录结构
	currentDir := tempDir
	for depth := 0; depth < 10; depth++ {
		currentDir = filepath.Join(currentDir, fmt.Sprintf("level%d", depth))
		if err := os.Mkdir(currentDir, 0755); err != nil {
			b.Fatalf("创建深层目录失败: %v", err)
		}

		// 在每层创建一个文件
		filename := filepath.Join(currentDir, fmt.Sprintf("file%d.txt", depth))
		content := strings.Repeat("a", 1024) // 1KB
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			b.Fatalf("创建深层文件失败: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := getPathSize(tempDir)
		if err != nil {
			b.Fatalf("getPathSize 失败: %v", err)
		}
	}
}

// BenchmarkMemoryUsage 内存使用测试
func BenchmarkMemoryUsage(b *testing.B) {
	tempDir := b.TempDir()

	// 创建大量小文件
	for i := 0; i < 1000; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("small%04d.txt", i))
		if err := os.WriteFile(filename, []byte("small file content"), 0644); err != nil {
			b.Fatalf("创建小文件失败: %v", err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs() // 报告内存分配情况

	for i := 0; i < b.N; i++ {
		_, err := getPathSize(tempDir)
		if err != nil {
			b.Fatalf("getPathSize 失败: %v", err)
		}
	}
}
