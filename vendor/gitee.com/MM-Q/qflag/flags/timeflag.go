package flags

import (
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/qflag/qerr"
)

// =============================================================================
// 时间类型标志
// =============================================================================

// 支持的时间格式列表
var supportedTimeFormats = []string{
	time.RFC3339,          // 2006-01-02T15:04:05Z07:00
	"2006-01-02",          // 2006-01-02
	"2006-01-02 15:04",    // 2006-01-02 15:04
	"2006-01-02 15:04:05", // 2006-01-02 15:04:05
	time.ANSIC,            // Mon Jan _2 15:04:05 MST 2006
	time.UnixDate,         // Mon Jan _2 15:04:05 MST 2006
	time.RubyDate,         // Mon Jan 02 15:04:05 -0700 2006
	time.RFC822,           // 01 Jan 06 15:04 MST
	time.RFC822Z,          // 02 Jan 06 15:04 -0700
	time.RFC850,           // Monday, 02-Jan-06 15:04:05 MST
	time.RFC1123,          // Mon, 02 Jan 2006 15:04:05 MST
	time.RFC1123Z,         // Mon, 02 Jan 2006 15:04:05 -0700
	time.Kitchen,          // 3:04PM
	time.Stamp,            // Mon Jan _2 15:04:05
	time.StampMilli,       // Mon Jan _2 15:04:05.000
	time.StampMicro,       // Mon Jan _2 15:04:05.000000
	time.StampNano,        // Mon Jan _2 15:04:05.000000000
}

// TimeFlag 时间类型标志结构体
// 继承BaseFlag[time.Time]泛型结构体,实现Flag接口
type TimeFlag struct {
	BaseFlag[time.Time]
	outputFormat string     // 自定义输出格式
	mu           sync.Mutex // 保护outputFormat和value的并发访问
	initOnce     sync.Once  // 确保只初始化一次
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *TimeFlag) Type() FlagType { return FlagTypeTime }

// Set 实现flag.Value接口, 解析并设置时间值
//
// 参数:
//   - value: 待解析的时间字符串
//
// 返回值:
//   - error: 解析或验证失败时返回错误信息
func (f *TimeFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var t time.Time
	var err error

	// 尝试解析时间字符串
	for _, format := range supportedTimeFormats {
		t, err = time.Parse(format, value)
		if err == nil {
			break
		}
	}
	if err != nil {
		return qerr.NewValidationErrorf("invalid time format: %v (supported formats include %v)", err, supportedTimeFormats)
	}

	// 调用基类设置值
	return f.BaseFlag.Set(t)
}

// String 实现flag.Value接口, 返回当前时间的字符串表示
//
// 返回值:
//   - string: 格式化后的时间字符串
//
// 注意: 加锁保证outputFormat和value的并发安全访问
func (f *TimeFlag) String() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 获取当前值和输出格式
	value := f.Get()
	format := f.outputFormat

	// 如果设置了输出格式, 则使用该格式
	if format != "" {
		return value.Format(format)
	}
	return value.Format(time.RFC3339) // 默认格式
}

// SetOutputFormat 设置时间输出格式
//
// 参数:
//   - format: 时间格式化字符串
//
// 注意: 此方法线程安全
func (f *TimeFlag) SetOutputFormat(format string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.outputFormat = format
}

// Init 初始化时间类型标志（使用 sync.Once 确保只初始化一次）
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值（字符串格式，支持多种时间表达）
//   - usage: 帮助说明
//
// 返回值:
//   - error: 初始化错误信息
//
// 支持的默认值格式:
//   - "now" 或 "" : 当前时间
//   - "zero" : 零时间 (time.Time{})
//   - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
//   - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
//   - RFC3339等标准格式
//
// 注意: 重复调用此方法是安全的，后续调用将被忽略
func (f *TimeFlag) Init(longName, shortName string, defValue string, usage string) error {
	var initErr error
	f.initOnce.Do(func() {
		f.mu.Lock()
		defer f.mu.Unlock()

		// 解析字符串默认值为 time.Time
		parsedTime, err := f.parseTimeString(defValue)
		if err != nil {
			initErr = qerr.NewValidationErrorf("invalid default time value '%s': %v", defValue, err)
			return
		}

		// 创建时间值指针
		timePtr := new(time.Time)
		*timePtr = parsedTime

		// 调用基类的 Init 方法
		initErr = f.BaseFlag.Init(longName, shortName, usage, timePtr)
	})
	return initErr
}

// parseTimeString 解析时间字符串为 time.Time
//
// 参数:
//   - s: 时间字符串
//
// 返回值:
//   - time.Time: 解析后的时间
//   - error: 解析错误
func (f *TimeFlag) parseTimeString(s string) (time.Time, error) {
	// 处理特殊值
	switch strings.ToLower(strings.TrimSpace(s)) {
	// now 或 空字符串 表示当前时间
	case "now", "":
		return time.Now(), nil
	// zero 表示零时间
	case "zero":
		return time.Time{}, nil
	}

	// 尝试解析相对时间（如 "1h", "-30m", "2h30m"）
	if d, err := time.ParseDuration(s); err == nil {
		return time.Now().Add(d), nil
	}

	// 尝试解析绝对时间格式
	for _, format := range supportedTimeFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	// 如果都解析失败，返回错误
	return time.Time{}, qerr.NewValidationError("unsupported time format, supported formats: 'now', 'zero', duration (1h, 30m), or standard time formats")
}

// =============================================================================
// 时间间隔类型标志
// =============================================================================

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
