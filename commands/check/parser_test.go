package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/colorlib"
)

func TestHashFileParser_ParseFile(t *testing.T) {
	cl := colorlib.New()
	parser := newHashFileParser(cl)

	// 创建临时测试目录
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		fileContent string
		isRelPath   bool
		expectError bool
		errorMsg    string
		expectCount int
	}{
		{
			name: "有效的校验文件",
			fileContent: `#md5#2024-01-01 10:00:00
c34652066a18513105ac1ab96fcbef8e test.txt
da39a3ee5e6b4b0d3255bfef95601890afd80709 empty.txt`,
			isRelPath:   false,
			expectError: false,
			expectCount: 2,
		},
		{
			name: "带注释和空行的校验文件",
			fileContent: `#sha256#2024-01-01 10:00:00
# 这是注释
c34652066a18513105ac1ab96fcbef8e test.txt

# 另一个注释
da39a3ee5e6b4b0d3255bfef95601890afd80709 empty.txt`,
			isRelPath:   false,
			expectError: false,
			expectCount: 2,
		},
		{
			name: "相对路径模式",
			fileContent: `#md5#2024-01-01 10:00:00
c34652066a18513105ac1ab96fcbef8e root/subdir/test.txt`,
			isRelPath:   true,
			expectError: false,
			expectCount: 1,
		},
		{
			name:        "空文件",
			fileContent: ``,
			isRelPath:   false,
			expectError: true,
			errorMsg:    "校验文件为空",
		},
		{
			name: "无效的文件头",
			fileContent: `invalid header
c34652066a18513105ac1ab96fcbef8e test.txt`,
			isRelPath:   false,
			expectError: true,
			errorMsg:    "校验文件头格式错误, 格式应为 #hashType#timestamp",
		},
		{
			name: "不支持的哈希算法",
			fileContent: `#unsupported#2024-01-01 10:00:00
c34652066a18513105ac1ab96fcbef8e test.txt`,
			isRelPath:   false,
			expectError: true,
			errorMsg:    "不支持的哈希算法: unsupported",
		},
		{
			name: "没有有效内容",
			fileContent: `#md5#2024-01-01 10:00:00
# 只有注释

`,
			isRelPath:   false,
			expectError: true,
			errorMsg:    "没有找到有效的校验文件内容",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时文件
			testFile := filepath.Join(tempDir, "test_"+tt.name+".hash")
			err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("创建测试文件失败: %v", err)
			}

			// 执行解析
			hashMap, hashType, err := parser.parseFile(testFile, tt.isRelPath)

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
				return
			}

			if hashType == "" {
				t.Errorf("哈希类型不应该为空")
			}

			if len(hashMap) != tt.expectCount {
				t.Errorf("哈希映射数量不匹配，期望: %d, 实际: %d", tt.expectCount, len(hashMap))
			}

			// 验证相对路径处理
			if tt.isRelPath && tt.expectCount > 0 {
				for virtualPath := range hashMap {
					// 检查是否包含虚拟根目录（考虑路径分隔符的差异）
					if strings.Contains(virtualPath, "ROOTDIR") {
						t.Logf("相对路径处理正确，虚拟路径: %s", virtualPath)
					} else {
						t.Errorf("相对路径模式下应该包含虚拟根目录，实际: %s", virtualPath)
					}
				}
			}
		})
	}
}

func TestHashFileParser_ProcessRelativePath(t *testing.T) {
	cl := colorlib.New()
	parser := newHashFileParser(cl)

	tests := []struct {
		name        string
		filePath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "正常的相对路径",
			filePath:    "root/subdir/file.txt",
			expectError: false,
		},
		{
			name:        "单级路径",
			filePath:    "root/file.txt",
			expectError: false,
		},
		{
			name:        "空路径",
			filePath:    "",
			expectError: true,
			errorMsg:    "无效的文件路径",
		},
		{
			name:        "只有根目录",
			filePath:    "root",
			expectError: true,
			errorMsg:    "获取相对路径时出错",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			virtualPath, realPath, err := parser.processRelativePath(tt.filePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有发生错误")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("错误消息不匹配，期望包含: %s, 实际: %s", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("不期望错误但发生了错误: %v", err)
				return
			}

			// 检查虚拟路径是否包含虚拟根目录（考虑路径分隔符的差异）
			if !strings.Contains(virtualPath, "ROOTDIR") {
				t.Errorf("虚拟路径应该包含虚拟根目录，实际: %s", virtualPath)
			}

			if realPath != tt.filePath {
				t.Errorf("真实路径不匹配，期望: %s, 实际: %s", tt.filePath, realPath)
			}
		})
	}
}

func TestHashFileParser_FileNotExists(t *testing.T) {
	cl := colorlib.New()
	parser := newHashFileParser(cl)

	// 测试不存在的文件
	_, _, err := parser.parseFile("nonexistent.hash", false)
	if err == nil {
		t.Errorf("期望错误但没有发生错误")
		return
	}

	expectedMsg := "校验文件不存在: nonexistent.hash"
	if err.Error() != expectedMsg {
		t.Errorf("错误消息不匹配，期望: %s, 实际: %s", expectedMsg, err.Error())
	}
}
