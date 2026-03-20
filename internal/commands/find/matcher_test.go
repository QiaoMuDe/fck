package find

import (
	"testing"
)

func TestNewPatternMatcher(t *testing.T) {
	tests := []struct {
		name      string
		cacheSize int
	}{
		{
			name:      "正常缓存大小",
			cacheSize: 100,
		},
		{
			name:      "零缓存大小",
			cacheSize: 0,
		},
		{
			name:      "负数缓存大小",
			cacheSize: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewPatternMatcher(tt.cacheSize)
			if matcher == nil {
				t.Errorf("NewPatternMatcher返回了nil")
			}
		})
	}
}

func TestPatternMatcher_GetRegex(t *testing.T) {
	matcher := NewPatternMatcher(10)

	tests := []struct {
		name        string
		pattern     string
		expectError bool
		expectNil   bool
	}{
		{
			name:        "空模式",
			pattern:     "",
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "有效正则表达式",
			pattern:     "test.*\\.go$",
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "无效正则表达式",
			pattern:     "[invalid",
			expectError: true,
			expectNil:   false,
		},
		{
			name:        "简单文本模式",
			pattern:     "hello",
			expectError: false,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, err := matcher.GetRegex(tt.pattern)

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

func TestPatternMatcher_CacheEfficiency(t *testing.T) {
	matcher := NewPatternMatcher(5)
	pattern := "test.*\\.go$"

	// 第一次获取，应该编译并缓存
	regex1, err := matcher.GetRegex(pattern)
	if err != nil {
		t.Fatalf("第一次获取失败: %v", err)
	}

	// 第二次获取，应该从缓存中获取
	regex2, err := matcher.GetRegex(pattern)
	if err != nil {
		t.Fatalf("第二次获取失败: %v", err)
	}

	// 验证返回的是同一个对象（缓存命中）
	if regex1 != regex2 {
		t.Errorf("缓存未命中，返回了不同的正则表达式对象")
	}
}

func BenchmarkPatternMatcher_GetRegex(b *testing.B) {
	matcher := NewPatternMatcher(100)
	pattern := "test.*\\.go$"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := matcher.GetRegex(pattern)
		if err != nil {
			b.Fatalf("获取正则表达式失败: %v", err)
		}
	}
}

func BenchmarkPatternMatcher_CacheHit(b *testing.B) {
	matcher := NewPatternMatcher(100)
	pattern := "test.*\\.go$"

	// 预热缓存
	_, _ = matcher.GetRegex(pattern)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := matcher.GetRegex(pattern)
		if err != nil {
			b.Fatalf("缓存命中失败: %v", err)
		}
	}
}
