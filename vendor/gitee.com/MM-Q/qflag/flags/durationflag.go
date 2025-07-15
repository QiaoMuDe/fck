package flags

import (
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/qflag/qerr"
)

// DurationFlag 时间间隔类型标志结构体
// 继承BaseFlag[time.Duration]泛型结构体,实现Flag接口
type DurationFlag struct {
	BaseFlag[time.Duration]
	mu sync.Mutex // 互斥锁, 用于保护并发访问
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *DurationFlag) Type() FlagType { return FlagTypeDuration }

// Set 实现flag.Value接口, 解析并设置时间间隔值
//
// 参数:
//   - value: 待设置的值
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *DurationFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查空值
	if value == "" {
		return qerr.NewValidationError("duration cannot be empty")
	}

	// 将单位转换为小写, 确保解析的准确性
	lowercaseValue := strings.ToLower(value)

	// 解析时间间隔字符串
	duration, err := time.ParseDuration(lowercaseValue)
	if err != nil {
		return qerr.NewValidationErrorf("invalid duration format: %v (valid units: ns/us/ms/s/m/h)", err)
	}

	// 检查负值（可选）
	if duration < 0 {
		return qerr.NewValidationError("negative duration not allowed")
	}

	// 调用基类方法设置值
	return f.BaseFlag.Set(duration)
}

// String 实现flag.Value接口, 返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
func (f *DurationFlag) String() string {
	return f.Get().String()
}
