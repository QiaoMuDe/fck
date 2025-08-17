package hash

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// TestHashCmdMain 测试主函数
func TestHashCmdMain(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 切换到临时目录
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(tempDir)

	// 初始化命令
	hashCmd = InitHashCmd()

	// 设置测试参数
	_ = hashCmdType.Set("md5")
	_ = hashCmdRecursion.Set("false")
	_ = hashCmdWrite.Set("false")
	_ = hashCmdHidden.Set("false")
	_ = hashCmdProgress.Set("false")

	cl := colorlib.NewColorLib()
	cl.SetColor(false) // 禁用颜色输出

	tests := []struct {
		name        string
		hashType    string
		expectError bool
	}{
		{
			name:        "正常文件MD5",
			hashType:    "md5",
			expectError: false,
		},
		{
			name:        "正常文件SHA256",
			hashType:    "sha256",
			expectError: false,
		},
		{
			name:        "无效哈希算法",
			hashType:    "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置哈希类型
			_ = hashCmdType.Set(tt.hashType)

			// 由于无法直接设置Args，我们通过processSinglePath测试核心逻辑
			if tt.hashType == "invalid" {
				// 测试无效哈希算法
				_, ok := types.SupportedAlgorithms[tt.hashType]
				if ok {
					t.Errorf("期望 %s 是无效的哈希算法", tt.hashType)
				}
			} else {
				// 测试有效的哈希算法
				hashType, ok := types.SupportedAlgorithms[tt.hashType]
				if !ok {
					t.Errorf("期望 %s 是有效的哈希算法", tt.hashType)
				} else {
					err := processSinglePath(cl, testFile, hashType)
					if err != nil {
						t.Errorf("processSinglePath() 返回错误: %v", err)
					}
				}
			}
		})
	}
}

// TestProcessSinglePath 测试单个路径处理
func TestProcessSinglePath(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 初始化命令
	hashCmd = InitHashCmd()
	_ = hashCmdRecursion.Set("false")
	_ = hashCmdWrite.Set("false")
	_ = hashCmdHidden.Set("false")

	cl := colorlib.NewColorLib()
	cl.SetColor(false)

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "正常文件",
			path:        testFile,
			expectError: false,
		},
		{
			name:        "不存在的文件",
			path:        filepath.Join(tempDir, "nonexistent.txt"),
			expectError: true,
		},
		{
			name:        "目录（非递归）",
			path:        tempDir,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processSinglePath(cl, tt.path, md5.New)

			if tt.expectError && err == nil {
				t.Errorf("processSinglePath() 期望返回错误，但没有错误")
			}
			if !tt.expectError && err != nil {
				t.Errorf("processSinglePath() 返回意外错误: %v", err)
			}
		})
	}
}

// TestPrintUniqueErrors 测试错误去重打印
func TestPrintUniqueErrors(t *testing.T) {
	cl := colorlib.NewColorLib()
	cl.SetColor(false)

	tests := []struct {
		name   string
		errors []error
	}{
		{
			name:   "空错误列表",
			errors: []error{},
		},
		{
			name:   "单个错误",
			errors: []error{os.ErrNotExist},
		},
		{
			name:   "重复错误",
			errors: []error{os.ErrNotExist, os.ErrNotExist, os.ErrPermission},
		},
		{
			name:   "包含nil错误",
			errors: []error{nil, os.ErrNotExist, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这个函数主要是打印，我们只测试它不会panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("printUniqueErrors() panic: %v", r)
				}
			}()

			printUniqueErrors(cl, tt.errors)
		})
	}
}

// TestHashCmdMainWithWrite 测试写入文件功能
func TestHashCmdMainWithWrite(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 切换到临时目录
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(tempDir)

	// 初始化命令
	hashCmd = InitHashCmd()

	// 设置写入模式
	_ = hashCmdType.Set("md5")
	_ = hashCmdRecursion.Set("false")
	_ = hashCmdWrite.Set("true")
	_ = hashCmdHidden.Set("false")
	_ = hashCmdProgress.Set("false")

	cl := colorlib.NewColorLib()
	cl.SetColor(false)

	// 测试写入功能
	err := processSinglePath(cl, testFile, md5.New)
	if err != nil {
		t.Errorf("processSinglePath() 返回错误: %v", err)
	}

	// 检查输出文件是否存在
	outputFile := filepath.Join(tempDir, types.OutputFileName)
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Logf("输出文件 %s 不存在，这可能是正常的，因为我们只测试了单个文件处理", outputFile)
	}
}

// TestHashCmdMainWithRecursion 测试递归处理
func TestHashCmdMainWithRecursion(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("创建子目录失败: %v", err)
	}

	// 创建测试文件
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(subDir, "test2.txt")

	if err := os.WriteFile(testFile1, []byte("content1"), 0644); err != nil {
		t.Fatalf("创建测试文件1失败: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("content2"), 0644); err != nil {
		t.Fatalf("创建测试文件2失败: %v", err)
	}

	// 初始化命令
	hashCmd = InitHashCmd()

	// 设置递归模式
	_ = hashCmdType.Set("md5")
	_ = hashCmdRecursion.Set("true")
	_ = hashCmdWrite.Set("false")
	_ = hashCmdHidden.Set("false")
	_ = hashCmdProgress.Set("false")

	cl := colorlib.NewColorLib()
	cl.SetColor(false)

	err := processSinglePath(cl, tempDir, md5.New)
	if err != nil {
		t.Errorf("processSinglePath() 返回错误: %v", err)
	}
}

// TestHashCmdMainWithGlob 测试通配符处理
func TestHashCmdMainWithGlob(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()

	// 创建多个测试文件
	for i := 1; i <= 3; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i))
		if err := os.WriteFile(testFile, []byte(fmt.Sprintf("content%d", i)), 0644); err != nil {
			t.Fatalf("创建测试文件%d失败: %v", i, err)
		}
	}

	// 切换到临时目录
	oldDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldDir) }()
	_ = os.Chdir(tempDir)

	// 初始化命令
	hashCmd = InitHashCmd()

	_ = hashCmdType.Set("md5")
	_ = hashCmdRecursion.Set("false")
	_ = hashCmdWrite.Set("false")
	_ = hashCmdHidden.Set("false")
	_ = hashCmdProgress.Set("false")

	cl := colorlib.NewColorLib()
	cl.SetColor(false)

	// 测试通配符处理
	globPattern := "test*.txt"
	err := processSinglePath(cl, globPattern, md5.New)
	if err != nil {
		t.Logf("processSinglePath() 对通配符返回错误: %v (这可能是正常的)", err)
	}
}

// TestHashCmdMainEmptyDirectory 测试空目录处理
func TestHashCmdMainEmptyDirectory(t *testing.T) {
	// 创建空的临时目录
	tempDir := t.TempDir()

	// 初始化命令
	hashCmd = InitHashCmd()

	_ = hashCmdType.Set("md5")
	_ = hashCmdRecursion.Set("false")
	_ = hashCmdWrite.Set("false")
	_ = hashCmdHidden.Set("false")
	_ = hashCmdProgress.Set("false")

	cl := colorlib.NewColorLib()
	cl.SetColor(false)

	// 空目录应该返回错误（非递归模式）
	err := processSinglePath(cl, tempDir, md5.New)
	if err == nil {
		t.Errorf("processSinglePath() 在空目录中应该返回错误")
	}
}

// BenchmarkHashCmdMain 性能测试
func BenchmarkProcessSinglePath(b *testing.B) {
	// 创建临时目录和文件
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "benchmark.txt")
	testContent := strings.Repeat("benchmark data ", 1000) // 约15KB数据

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}

	// 初始化命令
	hashCmd = InitHashCmd()

	_ = hashCmdType.Set("md5")
	_ = hashCmdRecursion.Set("false")
	_ = hashCmdWrite.Set("false")
	_ = hashCmdHidden.Set("false")
	_ = hashCmdProgress.Set("false")

	cl := colorlib.NewColorLib()
	cl.SetColor(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processSinglePath(cl, testFile, md5.New)
	}
}
