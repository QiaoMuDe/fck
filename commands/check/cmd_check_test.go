package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/colorlib"
)

func TestCheckCmdMain_Integration(t *testing.T) {
	cl := colorlib.New()

	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	content := "test content for main"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建有效的校验文件
	validCheckFile := filepath.Join(tempDir, "valid.hash")
	validContent := `#md5#2024-01-01 10:00:00
c34652066a18513105ac1ab96fcbef8e ` + testFile
	err = os.WriteFile(validCheckFile, []byte(validContent), 0644)
	if err != nil {
		t.Fatalf("创建校验文件失败: %v", err)
	}

	// 测试文件校验功能（直接调用内部函数）
	parser := newHashFileParser(cl)
	hashMap, hashType, err := parser.parseFile(validCheckFile, false)
	if err != nil {
		t.Fatalf("解析校验文件失败: %v", err)
	}

	if len(hashMap) != 1 {
		t.Errorf("期望解析1个文件，实际解析了%d个", len(hashMap))
	}

	if hashType == "" {
		t.Errorf("哈希类型不应该为空")
	}

	// 测试校验器
	checker := newFileChecker(cl, hashType)
	err = checker.checkFiles(hashMap)
	if err != nil {
		t.Errorf("文件校验失败: %v", err)
	}
}

func TestCheckCmdMain_ErrorCases(t *testing.T) {
	cl := colorlib.New()

	tests := []struct {
		name        string
		setupFile   func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "校验文件不存在",
			setupFile: func() string {
				return "nonexistent.hash"
			},
			expectError: true,
			errorMsg:    "校验文件不存在",
		},
		{
			name: "无效的校验文件头",
			setupFile: func() string {
				tempDir := t.TempDir()
				invalidFile := filepath.Join(tempDir, "invalid.hash")
				content := "invalid header\nc34652066a18513105ac1ab96fcbef8e test.txt"
				_ = os.WriteFile(invalidFile, []byte(content), 0644)
				return invalidFile
			},
			expectError: true,
			errorMsg:    "校验文件头格式错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkFile := tt.setupFile()

			// 直接测试解析器
			parser := newHashFileParser(cl)
			_, _, err := parser.parseFile(checkFile, false)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有发生错误")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("错误消息不匹配，期望包含: %s, 实际: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误但发生了错误: %v", err)
			}
		})
	}
}

func TestParseFile_BackwardCompatibility(t *testing.T) {
	cl := colorlib.New()

	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "compat_test.txt")
	content := "backward compatibility test"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建校验文件
	checkFile := filepath.Join(tempDir, "compat.hash")
	checkContent := `#md5#2024-01-01 10:00:00
c34652066a18513105ac1ab96fcbef8e ` + testFile
	err = os.WriteFile(checkFile, []byte(checkContent), 0644)
	if err != nil {
		t.Fatalf("创建校验文件失败: %v", err)
	}

	tests := []struct {
		name        string
		isRelPath   bool
		expectError bool
		expectCount int
	}{
		{
			name:        "绝对路径模式",
			isRelPath:   false,
			expectError: false,
			expectCount: 1,
		},
		{
			name:        "相对路径模式",
			isRelPath:   true,
			expectError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用新的解析器
			parser := newHashFileParser(cl)
			hashMap, hashType, err := parser.parseFile(checkFile, tt.isRelPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有发生错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误但发生了错误: %v", err)
				return
			}

			if hashType == "" {
				t.Errorf("哈希类型不应该为空")
			}

			if len(hashMap) != tt.expectCount {
				t.Errorf("哈希映射数量不匹配，期望: %d, 实际: %d", tt.expectCount, len(hashMap))
			}
		})
	}
}
