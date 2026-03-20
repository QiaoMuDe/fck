package check

import (
	"testing"
)

func TestHashLineValidator_ValidateLine(t *testing.T) {
	validator := newHashLineValidator()

	tests := []struct {
		name        string
		line        string
		lineNum     int
		expectHash  string
		expectPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效的MD5哈希行",
			line:        "c34652066a18513105ac1ab96fcbef8e test.txt",
			lineNum:     1,
			expectHash:  "c34652066a18513105ac1ab96fcbef8e",
			expectPath:  "test.txt",
			expectError: false,
		},
		{
			name:        "有效的SHA1哈希行",
			line:        "da39a3ee5e6b4b0d3255bfef95601890afd80709 empty.txt",
			lineNum:     2,
			expectHash:  "da39a3ee5e6b4b0d3255bfef95601890afd80709",
			expectPath:  "empty.txt",
			expectError: false,
		},
		{
			name:        "有效的SHA256哈希行",
			line:        "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 file.txt",
			lineNum:     3,
			expectHash:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			expectPath:  "file.txt",
			expectError: false,
		},
		{
			name:        "带引号的文件路径",
			line:        `c34652066a18513105ac1ab96fcbef8e "test file.txt"`,
			lineNum:     4,
			expectHash:  "c34652066a18513105ac1ab96fcbef8e",
			expectPath:  "test file.txt",
			expectError: false,
		},
		{
			name:        "带空格的文件路径",
			line:        "c34652066a18513105ac1ab96fcbef8e test file with spaces.txt",
			lineNum:     5,
			expectHash:  "c34652066a18513105ac1ab96fcbef8e",
			expectPath:  "test file with spaces.txt",
			expectError: false,
		},
		{
			name:        "空行应该跳过",
			line:        "",
			lineNum:     6,
			expectHash:  "",
			expectPath:  "",
			expectError: false,
		},
		{
			name:        "注释行应该跳过",
			line:        "# 这是注释",
			lineNum:     7,
			expectHash:  "",
			expectPath:  "",
			expectError: false,
		},
		{
			name:        "格式错误-缺少文件路径",
			line:        "c34652066a18513105ac1ab96fcbef8e",
			lineNum:     8,
			expectError: true,
			errorMsg:    "第8行格式错误: 缺少哈希值或文件路径",
		},
		{
			name:        "格式错误-无效哈希值",
			line:        "invalid_hash test.txt",
			lineNum:     9,
			expectError: true,
			errorMsg:    "第9行哈希值格式无效: invalid_hash",
		},
		{
			name:        "格式错误-哈希值太短",
			line:        "abc123 test.txt",
			lineNum:     10,
			expectError: true,
			errorMsg:    "第10行哈希值格式无效: abc123",
		},
		{
			name:        "安全错误-路径过长",
			line:        "c34652066a18513105ac1ab96fcbef8e " + generateLongPath(5000),
			lineNum:     12,
			expectError: true,
			errorMsg:    "第12行文件路径过长",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, path, err := validator.validateLine(tt.line, tt.lineNum)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有发生错误")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					// 对于路径过长的情况，只检查错误消息的开头
					switch tt.name {
					case "安全错误-路径过长":
						if !contains(err.Error(), "第12行文件路径过长") {
							t.Errorf("错误消息不匹配，期望包含: %s, 实际: %s", "第12行文件路径过长", err.Error())
						}
					case "安全错误-路径遍历攻击":
						if !contains(err.Error(), tt.errorMsg) {
							t.Errorf("错误消息不匹配，期望包含: %s, 实际: %s", tt.errorMsg, err.Error())
						}
					default:
						t.Errorf("错误消息不匹配，期望: %s, 实际: %s", tt.errorMsg, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误但发生了错误: %v", err)
				return
			}

			if hash != tt.expectHash {
				t.Errorf("哈希值不匹配，期望: %s, 实际: %s", tt.expectHash, hash)
			}

			if path != tt.expectPath {
				t.Errorf("路径不匹配，期望: %s, 实际: %s", tt.expectPath, path)
			}
		})
	}
}

func TestHashLineValidator_CleanFilePath(t *testing.T) {
	validator := newHashLineValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "去除双引号",
			input:    `"test.txt"`,
			expected: "test.txt",
		},
		{
			name:     "替换双反斜杠",
			input:    `path\\to\\file.txt`,
			expected: `path\to\file.txt`,
		},
		{
			name:     "清理路径",
			input:    `./path/../file.txt`,
			expected: "file.txt",
		},
		{
			name:     "复合清理",
			input:    `"./path\\..\\file.txt"`,
			expected: "file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.cleanFilePath(tt.input)
			if result != tt.expected {
				t.Errorf("清理路径结果不匹配，期望: %s, 实际: %s", tt.expected, result)
			}
		})
	}
}

func TestHashLineValidator_ValidateFilePath(t *testing.T) {
	validator := newHashLineValidator()

	tests := []struct {
		name        string
		filePath    string
		lineNum     int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效路径",
			filePath:    "test.txt",
			lineNum:     1,
			expectError: false,
		},
		{
			name:        "有效的子目录路径",
			filePath:    "subdir/test.txt",
			lineNum:     2,
			expectError: false,
		},
		{
			name:        "空路径",
			filePath:    "",
			lineNum:     4,
			expectError: true,
			errorMsg:    "第4行文件路径为空",
		},
		{
			name:        "路径过长",
			filePath:    generateLongPath(5000),
			lineNum:     5,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateFilePath(tt.filePath, tt.lineNum)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有发生错误")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					// 对于路径过长的情况，只检查错误消息的开头
					switch tt.name {
					case "路径过长":
						if !contains(err.Error(), "第5行文件路径过长") {
							t.Errorf("错误消息不匹配，期望包含: %s, 实际: %s", "第5行文件路径过长", err.Error())
						}
					default:
						t.Errorf("错误消息不匹配，期望: %s, 实际: %s", tt.errorMsg, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误但发生了错误: %v", err)
			}
		})
	}
}

// 辅助函数
func generateLongPath(length int) string {
	path := ""
	for i := 0; i < length; i++ {
		path += "a"
	}
	return path + ".txt"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
