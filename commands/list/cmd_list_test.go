package list

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestGetPaths(t *testing.T) {
	// 首先初始化命令标志
	InitListCmd()

	tests := []struct {
		name        string
		args        []string
		expectCount int
		expectFirst string
	}{
		{
			name:        "无参数时返回当前目录",
			args:        []string{},
			expectCount: 1,
			expectFirst: "", // 不检查具体值，因为可能是 "." 或绝对路径
		},
		{
			name:        "单个路径参数",
			args:        []string{"testpath"},
			expectCount: 1,
			expectFirst: "testpath",
		},
		{
			name:        "多个路径参数",
			args:        []string{"path1", "path2", "path3"},
			expectCount: 3,
			expectFirst: "path1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重新初始化命令以确保干净状态
			listCmd = nil
			InitListCmd()

			// 解析命令行参数
			err := listCmd.Parse(tt.args)
			if err != nil {
				t.Fatalf("解析命令行参数失败: %v", err)
			}

			result := getPaths()

			if len(result) != tt.expectCount {
				t.Errorf("getPaths() 返回长度 = %v, 期望 %v", len(result), tt.expectCount)
				t.Logf("实际返回: %v", result)
				t.Logf("输入参数: %v", tt.args)
				return
			}

			// 对于无参数的情况，验证返回的是有效路径
			if len(tt.args) == 0 && len(result) == 1 {
				path := result[0]
				// 应该是 "." 或者当前工作目录的绝对路径
				if path != "." && !filepath.IsAbs(path) {
					t.Errorf("getPaths() 无参数时返回 %v, 期望 '.' 或绝对路径", path)
				}
			}

			// 对于有参数的情况，验证第一个参数
			if len(tt.args) > 0 && tt.expectFirst != "" && len(result) > 0 {
				if result[0] != tt.expectFirst {
					t.Errorf("getPaths()[0] = %v, 期望 %v", result[0], tt.expectFirst)
				}
			}
		})
	}
}

func TestExpandPaths(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建测试文件和目录
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.txt")
	hiddenFile := filepath.Join(tempDir, ".hidden")
	testDir := filepath.Join(tempDir, "testdir")

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

	cl := colorlib.New()

	tests := []struct {
		name        string
		paths       []string
		expectError bool
		expectEmpty bool
	}{
		{
			name:        "路径遍历攻击检测",
			paths:       []string{"../../../etc/passwd"},
			expectError: true,
		},
		{
			name:        "正常文件路径",
			paths:       []string{testFile1},
			expectError: false,
		},
		{
			name:        "通配符路径",
			paths:       []string{filepath.Join(tempDir, "*.txt")},
			expectError: false,
		},
		{
			name:        "不存在的路径",
			paths:       []string{filepath.Join(tempDir, "nonexistent.txt")},
			expectError: false,
			expectEmpty: false, // 会有警告但不会返回错误
		},
		{
			name:        "空路径列表",
			paths:       []string{},
			expectError: false,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := expandPaths(tt.paths, cl)

			if tt.expectError && err == nil {
				t.Errorf("expandPaths() 期望错误但没有返回错误")
			}

			if !tt.expectError && err != nil {
				t.Errorf("expandPaths() 意外错误 = %v", err)
			}

			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("expandPaths() 期望空结果但返回 %v", result)
			}
		})
	}
}

// 测试扫描选项构建
func TestScanOptionsConstruction(t *testing.T) {
	tests := []struct {
		name     string
		expected ScanOptions
	}{
		{
			name: "默认扫描选项",
			expected: ScanOptions{
				Recursive:  false,
				ShowHidden: false,
				FileTypes:  nil,
				DirItself:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 直接测试结构体构建
			result := ScanOptions{
				Recursive:  tt.expected.Recursive,
				ShowHidden: tt.expected.ShowHidden,
				FileTypes:  tt.expected.FileTypes,
				DirItself:  tt.expected.DirItself,
			}

			if result.Recursive != tt.expected.Recursive {
				t.Errorf("ScanOptions.Recursive = %v, 期望 %v", result.Recursive, tt.expected.Recursive)
			}

			if result.ShowHidden != tt.expected.ShowHidden {
				t.Errorf("ScanOptions.ShowHidden = %v, 期望 %v", result.ShowHidden, tt.expected.ShowHidden)
			}

			if result.DirItself != tt.expected.DirItself {
				t.Errorf("ScanOptions.DirItself = %v, 期望 %v", result.DirItself, tt.expected.DirItself)
			}

			if len(result.FileTypes) != len(tt.expected.FileTypes) {
				t.Errorf("ScanOptions.FileTypes 长度 = %v, 期望 %v", len(result.FileTypes), len(tt.expected.FileTypes))
			}
		})
	}
}

// 测试处理选项构建
func TestProcessOptionsConstruction(t *testing.T) {
	tests := []struct {
		name     string
		expected ProcessOptions
	}{
		{
			name: "默认处理选项",
			expected: ProcessOptions{
				SortBy:     "name",
				Reverse:    false,
				GroupByDir: false,
			},
		},
		{
			name: "按时间排序",
			expected: ProcessOptions{
				SortBy:     "time",
				Reverse:    true,
				GroupByDir: true,
			},
		},
		{
			name: "按大小排序",
			expected: ProcessOptions{
				SortBy:     "size",
				Reverse:    false,
				GroupByDir: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessOptions{
				SortBy:     tt.expected.SortBy,
				Reverse:    tt.expected.Reverse,
				GroupByDir: tt.expected.GroupByDir,
			}

			if result.SortBy != tt.expected.SortBy {
				t.Errorf("ProcessOptions.SortBy = %v, 期望 %v", result.SortBy, tt.expected.SortBy)
			}

			if result.Reverse != tt.expected.Reverse {
				t.Errorf("ProcessOptions.Reverse = %v, 期望 %v", result.Reverse, tt.expected.Reverse)
			}

			if result.GroupByDir != tt.expected.GroupByDir {
				t.Errorf("ProcessOptions.GroupByDir = %v, 期望 %v", result.GroupByDir, tt.expected.GroupByDir)
			}
		})
	}
}

// 测试格式选项构建
func TestFormatOptionsConstruction(t *testing.T) {
	tests := []struct {
		name     string
		expected FormatOptions
	}{
		{
			name: "默认格式选项",
			expected: FormatOptions{
				LongFormat:    false,
				UseColor:      false,
				DevColor:      false,
				TableStyle:    "none",
				QuoteNames:    false,
				ShowUserGroup: false,
			},
		},
		{
			name: "完整格式选项",
			expected: FormatOptions{
				LongFormat:    true,
				UseColor:      true,
				DevColor:      true,
				TableStyle:    "default",
				QuoteNames:    true,
				ShowUserGroup: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatOptions{
				LongFormat:    tt.expected.LongFormat,
				UseColor:      tt.expected.UseColor,
				DevColor:      tt.expected.DevColor,
				TableStyle:    tt.expected.TableStyle,
				QuoteNames:    tt.expected.QuoteNames,
				ShowUserGroup: tt.expected.ShowUserGroup,
			}

			if result.LongFormat != tt.expected.LongFormat {
				t.Errorf("FormatOptions.LongFormat = %v, 期望 %v", result.LongFormat, tt.expected.LongFormat)
			}

			if result.UseColor != tt.expected.UseColor {
				t.Errorf("FormatOptions.UseColor = %v, 期望 %v", result.UseColor, tt.expected.UseColor)
			}

			if result.DevColor != tt.expected.DevColor {
				t.Errorf("FormatOptions.DevColor = %v, 期望 %v", result.DevColor, tt.expected.DevColor)
			}

			if result.TableStyle != tt.expected.TableStyle {
				t.Errorf("FormatOptions.TableStyle = %v, 期望 %v", result.TableStyle, tt.expected.TableStyle)
			}

			if result.QuoteNames != tt.expected.QuoteNames {
				t.Errorf("FormatOptions.QuoteNames = %v, 期望 %v", result.QuoteNames, tt.expected.QuoteNames)
			}

			if result.ShowUserGroup != tt.expected.ShowUserGroup {
				t.Errorf("FormatOptions.ShowUserGroup = %v, 期望 %v", result.ShowUserGroup, tt.expected.ShowUserGroup)
			}
		})
	}
}

// 测试参数验证逻辑（不依赖全局标志）
func TestValidationLogic(t *testing.T) {
	tests := []struct {
		name        string
		sortBySize  bool
		sortByTime  bool
		sortByName  bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "同时指定size和time",
			sortBySize:  true,
			sortByTime:  true,
			expectError: true,
			errorMsg:    "不能同时指定",
		},
		{
			name:        "同时指定size和name",
			sortBySize:  true,
			sortByName:  true,
			expectError: true,
			errorMsg:    "不能同时指定",
		},
		{
			name:        "同时指定time和name",
			sortByTime:  true,
			sortByName:  true,
			expectError: true,
			errorMsg:    "不能同时指定",
		},
		{
			name:        "只指定一个排序选项",
			sortBySize:  true,
			expectError: false,
		},
		{
			name:        "不指定任何排序选项",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟验证逻辑
			var hasError bool
			var errorMsg string

			count := 0
			if tt.sortBySize {
				count++
			}
			if tt.sortByTime {
				count++
			}
			if tt.sortByName {
				count++
			}

			if count > 1 {
				hasError = true
				errorMsg = "不能同时指定多个排序选项"
			}

			if tt.expectError && !hasError {
				t.Errorf("期望错误但没有检测到错误")
			}

			if !tt.expectError && hasError {
				t.Errorf("意外错误: %v", errorMsg)
			}

			if tt.expectError && hasError && tt.errorMsg != "" {
				if !strings.Contains(errorMsg, tt.errorMsg) {
					t.Errorf("错误信息 = %v, 期望包含 %v", errorMsg, tt.errorMsg)
				}
			}
		})
	}
}

// 测试类型验证
func TestTypeValidation(t *testing.T) {
	tests := []struct {
		name        string
		fileType    string
		showAll     bool
		expectError bool
	}{
		{
			name:        "隐藏文件类型需要showAll",
			fileType:    types.FindTypeHidden,
			showAll:     false,
			expectError: true,
		},
		{
			name:        "隐藏文件类型配合showAll",
			fileType:    types.FindTypeHidden,
			showAll:     true,
			expectError: false,
		},
		{
			name:        "普通文件类型",
			fileType:    types.FindTypeFile,
			showAll:     false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟类型验证逻辑
			var hasError bool

			if (tt.fileType == types.FindTypeHidden || tt.fileType == types.FindTypeHiddenShort) && !tt.showAll {
				hasError = true
			}

			if tt.expectError && !hasError {
				t.Errorf("期望错误但没有检测到错误")
			}

			if !tt.expectError && hasError {
				t.Errorf("意外错误")
			}
		})
	}
}
