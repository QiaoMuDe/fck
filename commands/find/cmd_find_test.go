package find

import (
	"strings"
	"testing"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

func TestCreateFindConfig(t *testing.T) {
	initTestFlags() // 确保标志变量已初始化

	tests := []struct {
		name        string
		setupFlags  func()
		expectError bool
		validate    func(*types.FindConfig) bool
	}{
		{
			name: "默认配置",
			setupFlags: func() {
				// 重置所有标志为默认值
				_ = findCmdRegex.Set("false")
				_ = findCmdWholeWord.Set("false")
				_ = findCmdCase.Set("false")
				_ = findCmdName.Set("")
				_ = findCmdExcludeName.Set("")
				_ = findCmdPath.Set("")
				_ = findCmdExcludePath.Set("")
			},
			expectError: false,
			validate: func(config *types.FindConfig) bool {
				return !config.IsRegex && !config.WholeWord && !config.CaseSensitive
			},
		},
		{
			name: "正则表达式模式",
			setupFlags: func() {
				_ = findCmdRegex.Set("true")
				_ = findCmdName.Set("test.*\\.go$")
				_ = findCmdCase.Set("true")
			},
			expectError: false,
			validate: func(config *types.FindConfig) bool {
				return config.IsRegex && config.CaseSensitive && config.NameRegex != nil
			},
		},
		{
			name: "无效正则表达式",
			setupFlags: func() {
				_ = findCmdRegex.Set("true")
				_ = findCmdName.Set("[invalid")
			},
			expectError: true,
			validate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFlags()
			cl := colorlib.NewColorLib()

			config, err := createFindConfig(cl)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				return
			}

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if tt.validate != nil && !tt.validate(config) {
				t.Errorf("配置验证失败")
			}
		})
	}
}

func TestCompileRegexPattern(t *testing.T) {
	initTestFlags() // 确保标志变量已初始化

	tests := []struct {
		name          string
		pattern       string
		isRegex       bool
		wholeWord     bool
		caseSensitive bool
		expectError   bool
		expectNil     bool
	}{
		{
			name:          "空模式",
			pattern:       "",
			isRegex:       false,
			wholeWord:     false,
			caseSensitive: true,
			expectError:   false,
			expectNil:     true,
		},
		{
			name:          "简单文本模式",
			pattern:       "test",
			isRegex:       false,
			wholeWord:     false,
			caseSensitive: true,
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "全词匹配",
			pattern:       "test",
			isRegex:       false,
			wholeWord:     true,
			caseSensitive: true,
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "不区分大小写",
			pattern:       "Test",
			isRegex:       false,
			wholeWord:     false,
			caseSensitive: false,
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "有效正则表达式",
			pattern:       "test.*\\.go$",
			isRegex:       true,
			wholeWord:     false,
			caseSensitive: true,
			expectError:   false,
			expectNil:     false,
		},
		{
			name:          "无效正则表达式",
			pattern:       "[invalid",
			isRegex:       true,
			wholeWord:     false,
			caseSensitive: true,
			expectError:   true,
			expectNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, err := compileRegexPattern(tt.pattern, tt.isRegex, tt.wholeWord, tt.caseSensitive)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				return
			}

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if tt.expectNil && regex != nil {
				t.Errorf("期望返回nil但返回了正则表达式对象")
			}

			if !tt.expectNil && regex == nil {
				t.Errorf("期望返回正则表达式对象但返回了nil")
			}
		})
	}
}

func TestProcessExtensions(t *testing.T) {
	initTestFlags() // 确保标志变量已初始化

	tests := []struct {
		name       string
		extensions []string
		validate   func(*types.FindConfig) bool
	}{
		{
			name:       "空扩展名列表",
			extensions: []string{},
			validate: func(config *types.FindConfig) bool {
				// 检查map是否为空
				count := 0
				config.FindExtSliceMap.Range(func(key, value interface{}) bool {
					count++
					return true
				})
				return count == 0
			},
		},
		{
			name:       "带点的扩展名",
			extensions: []string{".go", ".txt"},
			validate: func(config *types.FindConfig) bool {
				val1, ok1 := config.FindExtSliceMap.Load(".go")
				val2, ok2 := config.FindExtSliceMap.Load(".txt")
				return ok1 && ok2 && val1.(bool) && val2.(bool)
			},
		},
		{
			name:       "不带点的扩展名",
			extensions: []string{"go", "txt"},
			validate: func(config *types.FindConfig) bool {
				val1, ok1 := config.FindExtSliceMap.Load(".go")
				val2, ok2 := config.FindExtSliceMap.Load(".txt")
				return ok1 && ok2 && val1.(bool) && val2.(bool)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟设置扩展名 - 重置后一次性设置所有扩展名
			findCmdExt.Reset()
			if len(tt.extensions) > 0 {
				// 将所有扩展名连接成一个字符串，用逗号分隔
				extStr := strings.Join(tt.extensions, ",")
				_ = findCmdExt.Set(extStr)
			}

			config := &types.FindConfig{}
			err := processExtensions(config)

			if err != nil {
				t.Errorf("意外错误: %v", err)
				return
			}

			if !tt.validate(config) {
				// 打印调试信息
				t.Logf("扩展名列表: %v", tt.extensions)
				config.FindExtSliceMap.Range(func(key, value interface{}) bool {
					t.Logf("存储的扩展名: %s = %v", key, value)
					return true
				})
				t.Errorf("扩展名处理验证失败")
			}
		})
	}
}

// 基准测试
func BenchmarkCreateFindConfig(b *testing.B) {
	initTestFlags() // 确保标志变量已初始化
	cl := colorlib.NewColorLib()

	// 设置一些标志
	_ = findCmdRegex.Set("true")
	_ = findCmdName.Set("test.*\\.go$")
	_ = findCmdCase.Set("true")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := createFindConfig(cl)
		if err != nil {
			b.Fatalf("配置创建失败: %v", err)
		}
	}
}

func BenchmarkCompileRegexPattern(b *testing.B) {
	pattern := "test.*\\.go$"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := compileRegexPattern(pattern, true, false, true)
		if err != nil {
			b.Fatalf("正则表达式编译失败: %v", err)
		}
	}
}
