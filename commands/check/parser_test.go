package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
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
			hashMap, hashFunc, err := parser.parseFile(testFile, tt.isRelPath)

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

			if hashFunc == nil {
				t.Errorf("哈希函数不应该为nil")
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

func TestHashFileParser_ParseHeader(t *testing.T) {
	cl := colorlib.New()
	parser := newHashFileParser(cl)

	tests := []struct {
		name        string
		header      string
		content     string
		expectError bool
		errorMsg    string
		expectAlgo  string
	}{
		{
			name:        "有效的MD5头",
			header:      "#md5#2024-01-01 10:00:00",
			content:     "c34652066a18513105ac1ab96fcbef8e test.txt",
			expectError: false,
			expectAlgo:  "md5",
		},
		{
			name:        "有效的SHA1头",
			header:      "#sha1#2024-01-01 10:00:00",
			content:     "da39a3ee5e6b4b0d3255bfef95601890afd80709 test.txt",
			expectError: false,
			expectAlgo:  "sha1",
		},
		{
			name:        "有效的SHA256头",
			header:      "#sha256#2024-01-01 10:00:00",
			content:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 test.txt",
			expectError: false,
			expectAlgo:  "sha256",
		},
		{
			name:        "有效的SHA512头",
			header:      "#sha512#2024-01-01 10:00:00",
			content:     "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e test.txt",
			expectError: false,
			expectAlgo:  "sha512",
		},
		{
			name:        "空的哈希类型#01",
			header:      "##2024-01-01 10:00:00",
			expectError: true,
			errorMsg:    "校验文件头格式错误, 格式应为 #hashType#timestamp",
		},
		{
			name:        "不支持的哈希算法",
			header:      "#blake2#2024-01-01 10:00:00",
			expectError: true,
			errorMsg:    "不支持的哈希算法: blake2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时文件
			tempDir := t.TempDir()
			testFile := filepath.Join(tempDir, "test_header.hash")
			fileContent := tt.header
			if tt.content != "" {
				fileContent += "\n" + tt.content
			}
			err := os.WriteFile(testFile, []byte(fileContent), 0644)
			if err != nil {
				t.Fatalf("创建测试文件失败: %v", err)
			}

			// 执行解析
			_, hashFunc, err := parser.parseFile(testFile, false)

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

			if hashFunc == nil {
				t.Errorf("哈希函数不应该为nil")
			}

			// 验证哈希函数是否正确
			if tt.expectAlgo != "" {
				expectedFunc, exists := types.SupportedAlgorithms[tt.expectAlgo]
				if !exists {
					t.Fatalf("测试配置错误: 不支持的算法 %s", tt.expectAlgo)
				}

				// 通过比较哈希结果来验证函数是否正确
				testData := []byte("test")
				expected := expectedFunc()
				expected.Write(testData)
				expectedSum := expected.Sum(nil)

				actual := hashFunc()
				actual.Write(testData)
				actualSum := actual.Sum(nil)

				if len(expectedSum) != len(actualSum) {
					t.Errorf("哈希函数不匹配，期望长度: %d, 实际长度: %d", len(expectedSum), len(actualSum))
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
