package hash

import (
	"testing"
)

// TestInitHashCmd 测试哈希命令初始化
func TestInitHashCmd(t *testing.T) {
	cmd := InitHashCmd()

	if cmd == nil {
		t.Fatal("InitHashCmd() 返回 nil")
	}

	// 测试命令基本属性
	if cmd.Name() != "hash" {
		t.Errorf("命令名称不正确: got %s, want hash", cmd.Name())
	}

	if cmd.ShortName() != "h" {
		t.Errorf("命令短名称不正确: got %s, want h", cmd.ShortName())
	}

	// 测试标志是否正确初始化
	if hashCmdType == nil {
		t.Error("hashCmdType 标志未初始化")
	}

	if hashCmdRecursion == nil {
		t.Error("hashCmdRecursion 标志未初始化")
	}

	if hashCmdWrite == nil {
		t.Error("hashCmdWrite 标志未初始化")
	}

	if hashCmdHidden == nil {
		t.Error("hashCmdHidden 标志未初始化")
	}

	if hashCmdProgress == nil {
		t.Error("hashCmdProgress 标志未初始化")
	}
}

// TestHashCmdFlags 测试哈希命令标志
func TestHashCmdFlags(t *testing.T) {
	// 重新初始化以确保干净状态
	hashCmd = InitHashCmd()

	tests := []struct {
		name     string
		flagName string
		setValue string
		getValue func() string
	}{
		{
			name:     "type标志",
			flagName: "type",
			setValue: "sha256",
			getValue: func() string { return hashCmdType.Get() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置标志值
			if tt.flagName == "type" {
				err := hashCmdType.Set(tt.setValue)
				if err != nil {
					t.Errorf("设置 %s 标志失败: %v", tt.flagName, err)
				}
			}

			// 验证标志值
			if got := tt.getValue(); got != tt.setValue {
				t.Errorf("%s 标志值不正确: got %s, want %s", tt.flagName, got, tt.setValue)
			}
		})
	}
}

// TestHashCmdBoolFlags 测试布尔标志
func TestHashCmdBoolFlags(t *testing.T) {
	// 重新初始化以确保干净状态
	hashCmd = InitHashCmd()

	tests := []struct {
		name string
		flag interface {
			Set(string) error
			Get() bool
		}
		setValue     bool
		defaultValue bool
	}{
		{
			name:         "recursion标志",
			flag:         hashCmdRecursion,
			setValue:     true,
			defaultValue: false,
		},
		{
			name:         "write标志",
			flag:         hashCmdWrite,
			setValue:     true,
			defaultValue: false,
		},
		{
			name:         "hidden标志",
			flag:         hashCmdHidden,
			setValue:     true,
			defaultValue: false,
		},
		{
			name:         "progress标志",
			flag:         hashCmdProgress,
			setValue:     true,
			defaultValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试默认值
			if got := tt.flag.Get(); got != tt.defaultValue {
				t.Errorf("%s 默认值不正确: got %v, want %v", tt.name, got, tt.defaultValue)
			}

			// 设置新值
			var setValue string
			if tt.setValue {
				setValue = "true"
			} else {
				setValue = "false"
			}

			err := tt.flag.Set(setValue)
			if err != nil {
				t.Errorf("设置 %s 标志失败: %v", tt.name, err)
			}

			// 验证新值
			if got := tt.flag.Get(); got != tt.setValue {
				t.Errorf("%s 标志值不正确: got %v, want %v", tt.name, got, tt.setValue)
			}
		})
	}
}

// TestHashCmdTypeFlag 测试类型标志的有效值
func TestHashCmdTypeFlag(t *testing.T) {
	hashCmd = InitHashCmd()

	validTypes := []string{"md5", "sha1", "sha256", "sha512"}

	for _, validType := range validTypes {
		t.Run("valid_type_"+validType, func(t *testing.T) {
			err := hashCmdType.Set(validType)
			if err != nil {
				t.Errorf("设置有效类型 %s 失败: %v", validType, err)
			}

			if got := hashCmdType.Get(); got != validType {
				t.Errorf("类型值不正确: got %s, want %s", got, validType)
			}
		})
	}

	// 测试无效类型
	invalidTypes := []string{"invalid", "md4", "sha3"}

	for _, invalidType := range invalidTypes {
		t.Run("invalid_type_"+invalidType, func(t *testing.T) {
			err := hashCmdType.Set(invalidType)
			if err == nil {
				t.Errorf("设置无效类型 %s 应该失败，但成功了", invalidType)
			}
		})
	}
}

// TestHashCmdDefaultValues 测试默认值
func TestHashCmdDefaultValues(t *testing.T) {
	hashCmd = InitHashCmd()

	// 测试默认值
	if hashCmdType.Get() != "md5" {
		t.Errorf("type 默认值不正确: got %s, want md5", hashCmdType.Get())
	}

	if hashCmdRecursion.Get() != false {
		t.Errorf("recursion 默认值不正确: got %v, want false", hashCmdRecursion.Get())
	}

	if hashCmdWrite.Get() != false {
		t.Errorf("write 默认值不正确: got %v, want false", hashCmdWrite.Get())
	}

	if hashCmdHidden.Get() != false {
		t.Errorf("hidden 默认值不正确: got %v, want false", hashCmdHidden.Get())
	}

	if hashCmdProgress.Get() != false {
		t.Errorf("progress 默认值不正确: got %v, want false", hashCmdProgress.Get())
	}
}

// TestHashCmdMultipleInit 测试多次初始化
func TestHashCmdMultipleInit(t *testing.T) {
	// 第一次初始化
	cmd1 := InitHashCmd()
	if cmd1 == nil {
		t.Fatal("第一次 InitHashCmd() 返回 nil")
	}

	// 第二次初始化
	cmd2 := InitHashCmd()
	if cmd2 == nil {
		t.Fatal("第二次 InitHashCmd() 返回 nil")
	}

	// 验证两次初始化返回的是同一个对象（如果是单例模式）
	// 或者至少都是有效的命令对象
	if cmd1.Name() != cmd2.Name() {
		t.Errorf("多次初始化返回不同的命令名称: %s vs %s", cmd1.Name(), cmd2.Name())
	}
}

// TestHashCmdFlagInteraction 测试标志之间的交互
func TestHashCmdFlagInteraction(t *testing.T) {
	hashCmd = InitHashCmd()

	// 测试设置多个标志
	_ = hashCmdType.Set("sha256")
	_ = hashCmdRecursion.Set("true")
	_ = hashCmdWrite.Set("true")
	_ = hashCmdHidden.Set("true")
	_ = hashCmdProgress.Set("true")

	// 验证所有标志都正确设置
	if hashCmdType.Get() != "sha256" {
		t.Errorf("type 标志设置失败")
	}
	if !hashCmdRecursion.Get() {
		t.Errorf("recursion 标志设置失败")
	}
	if !hashCmdWrite.Get() {
		t.Errorf("write 标志设置失败")
	}
	if !hashCmdHidden.Get() {
		t.Errorf("hidden 标志设置失败")
	}
	if !hashCmdProgress.Get() {
		t.Errorf("progress 标志设置失败")
	}
}

// BenchmarkInitHashCmd 性能测试初始化
func BenchmarkInitHashCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = InitHashCmd()
	}
}

// BenchmarkHashCmdFlagSet 性能测试标志设置
func BenchmarkHashCmdFlagSet(b *testing.B) {
	hashCmd = InitHashCmd()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hashCmdType.Set("md5")
		_ = hashCmdRecursion.Set("false")
		_ = hashCmdWrite.Set("false")
		_ = hashCmdHidden.Set("false")
		_ = hashCmdProgress.Set("false")
	}
}
