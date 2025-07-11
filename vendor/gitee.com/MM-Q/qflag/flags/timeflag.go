package flags

import (
	"sync"
	"time"

	"gitee.com/MM-Q/qflag/qerr"
)

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
}

// Type 返回标志类型
func (f *TimeFlag) Type() FlagType { return FlagTypeTime }

// Set 实现flag.Value接口, 解析并设置时间值
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
// 加锁保证outputFormat和value的并发安全访问
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
func (f *TimeFlag) SetOutputFormat(format string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.outputFormat = format
}
