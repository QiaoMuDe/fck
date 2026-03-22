package flag

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/internal/utils"
)

// DurationFlag 持续时间标志
//
// DurationFlag 用于处理时间间隔类型的命令行参数。
// 支持Go标准库time.ParseDuration所支持的所有格式, 如 "300ms", "-1.5h", "2h45m" 等。
//
// 支持的格式:
//   - "ns": 纳秒
//   - "us" (或 "µs"): 微秒
//   - "ms": 毫秒
//   - "s": 秒
//   - "m": 分钟
//   - "h": 小时
//
// 注意事项:
//   - 支持负数表示负时间间隔
//   - 支持小数表示部分时间单位
//   - 可以组合多个单位, 如 "1h30m"
type DurationFlag struct {
	*BaseFlag[time.Duration]
}

// NewDurationFlag 创建新的持续时间标志
//
// 参数:
//   - longName: 长选项名, 如 "timeout"
//   - shortName: 短选项名, 如 "t"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *DurationFlag: 持续时间标志实例
func NewDurationFlag(longName, shortName, desc string, default_ time.Duration) *DurationFlag {
	return &DurationFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeDuration, longName, shortName, desc, default_),
	}
}

// Set 设置持续时间标志的值
//
// 参数:
//   - value: 要设置的时间间隔字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 time.ParseDuration 解析字符串
//   - 支持所有Go标准库支持的时间格式
//   - 如果值无法解析为时间间隔, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *DurationFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_DURATION", "duration value cannot be empty", nil)
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return types.WrapParseError(err, "duration", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(d); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = d
	f.isSet = true

	return nil
}

// TimeFlag 时间标志
//
// TimeFlag 用于处理时间类型的命令行参数。
// 支持自动检测多种常见时间格式, 也支持指定特定格式进行解析。
//
// 特性:
//   - 自动检测常见时间格式
//   - 支持自定义格式解析
//   - 记录当前使用的格式
//   - 线程安全的格式存储
//
// 常见支持格式:
//   - RFC3339: "2006-01-02T15:04:05Z07:00"
//   - RFC1123: "Mon, 02 Jan 2006 15:04:05 MST"
//   - 日期格式: "2006-01-02", "2006/01/02"
//   - 时间格式: "15:04:05", "15:04"
//   - 其他常见格式
type TimeFlag struct {
	*BaseFlag[time.Time]
	// 当前使用的格式
	currentFormat string
}

// NewTimeFlag 创建新的时间标志
//
// 参数:
//   - longName: 长选项名, 如 "start-time"
//   - shortName: 短选项名, 如 "s"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *TimeFlag: 时间标志实例
func NewTimeFlag(longName, shortName, desc string, default_ time.Time) *TimeFlag {
	return &TimeFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeTime, longName, shortName, desc, default_),
	}
}

// Set 使用常见格式自动解析时间
//
// 参数:
//   - value: 要设置的时间字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 types.ParseTimeWithCommonFormats 自动检测格式
//   - 成功解析后会记录使用的格式
//   - 如果无法匹配任何格式, 返回错误
//   - 先解析，然后验证，最后设置值和格式
func (f *TimeFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_TIME", "time value cannot be empty", nil)
	}

	t, format, err := types.ParseTimeWithCommonFormats(value)
	if err != nil {
		return err
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(t); err != nil {
			return err
		}
	}

	// 设置值和格式以及标记为已设置
	*f.value = t
	f.currentFormat = format
	f.isSet = true

	return nil
}

// SetWithFormat 使用指定格式解析时间
//
// 参数:
//   - value: 要设置的时间字符串
//   - format: 时间格式字符串, 遵循Go的time.Format布局
//
// 返回值:
//   - error: 如果解析失败返回错误
//
// 注意事项:
//   - 使用 time.Parse 按指定格式解析
//   - 成功解析后会更新当前使用的格式
//   - 格式字符串必须遵循Go的time.Format布局规则
func (f *TimeFlag) SetWithFormat(value, format string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	t, err := time.Parse(format, value)
	if err != nil {
		return types.WrapParseError(err, "time", value)
	}

	*f.value = t
	f.currentFormat = format
	f.isSet = true

	return nil
}

// GetFormat 获取当前使用的时间格式
//
// 返回值:
//   - string: 当前使用的时间格式, 如果未设置则返回空字符串
//
// 注意事项:
//   - 此方法是线程安全的
//   - 返回的格式可用于格式化其他时间值
func (f *TimeFlag) GetFormat() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.currentFormat
}

// FormatTime 使用当前解析格式格式化时间
//
// 参数:
//   - t: 要格式化的时间值
//
// 返回值:
//   - string: 格式化后的时间字符串
//
// 注意事项:
//   - 如果当前格式为空, 使用RFC3339格式
//   - 此方法是线程安全的
func (f *TimeFlag) FormatTime(t time.Time) string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.currentFormat == "" {
		return t.Format(types.TimeFormatRFC3339)
	}
	return t.Format(f.currentFormat)
}

// String 返回格式化的时间字符串
//
// 返回值:
//   - string: 格式化的时间字符串
//
// 注意事项:
//   - 如果当前格式为空, 使用RFC3339格式
//   - 此方法是线程安全的
//   - 实现了fmt.Stringer接口
func (f *TimeFlag) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if f.value == nil {
		return ""
	}
	if f.currentFormat == "" {
		return f.value.Format(types.TimeFormatRFC3339)
	}
	return f.value.Format(f.currentFormat)
}

// SizeFlag 大小标志 (支持KB、MB、GB等单位)
//
// SizeFlag 用于处理大小类型的命令行参数, 支持多种大小单位。
// 可以解析带有单位的大小值, 并将其转换为字节数。
//
// 支持的单位:
//   - B/b: 字节
//   - KB/kb/K/k: 千字节 (1024字节)
//   - MB/mb/M/m: 兆字节 (1024^2字节)
//   - GB/gb/G/g: 吉字节 (1024^3字节)
//   - TB/tb/T/t: 太字节 (1024^4字节)
//   - PB/pb/P/p: 拍字节 (1024^5字节)
//   - KiB/kib: 二进制千字节 (1024字节)
//   - MiB/mib: 二进制兆字节 (1024^2字节)
//   - GiB/gib: 二进制吉字节 (1024^3字节)
//   - TiB/tib: 二进制太字节 (1024^4字节)
//   - PiB/pib: 二进制拍字节 (1024^5字节)
//
// 注意事项:
//   - 支持小数, 如 "1.5MB"
//   - 不支持负数
//   - 默认单位为字节(B)
//   - 大小写不敏感
type SizeFlag struct {
	*BaseFlag[int64]
}

// NewSizeFlag 创建新的大小标志
//
// 参数:
//   - longName: 长选项名, 如 "max-size"
//   - shortName: 短选项名, 如 "s"
//   - desc: 标志描述
//   - default_: 默认值(以字节为单位)
//
// 返回值:
//   - *SizeFlag: 大小标志实例
func NewSizeFlag(longName, shortName, desc string, default_ int64) *SizeFlag {
	return &SizeFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeSize, longName, shortName, desc, default_),
	}
}

// Set 设置大小标志的值
//
// 参数:
//   - value: 要设置的大小字符串, 可包含单位
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 支持多种大小单位, 大小写不敏感
//   - 支持小数值, 如 "1.5MB"
//   - 不支持负数
//   - 如果未指定单位, 默认为字节(B)
//   - 如果值超出int64范围, 返回错误
//   - 先解析，然后验证，最后设置值
func (f *SizeFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 移除可能的空格
	value = strings.TrimSpace(value)

	// 检查是否包含单位
	if value == "" {
		return types.NewError("INVALID_SIZE", "size value cannot be empty", nil)
	}

	// 提取数字部分和单位部分
	var numStr, unit string
	for i, r := range value {
		// 跳过非数字、点号的字符 (不允许负号)
		if (r >= '0' && r <= '9') || r == '.' {
			continue
		}

		numStr = value[:i]                // 数字部分
		unit = strings.ToUpper(value[i:]) // 单位部分
		break
	}

	// 如果没有找到单位, 整个值都是数字
	if numStr == "" {
		numStr = value
		unit = "B" // 默认单位为字节
	}

	// 解析数字部分
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return types.WrapParseError(err, "size", value)
	}

	// 检查是否为负数
	if num < 0 {
		return types.NewError("INVALID_SIZE", "size value cannot be negative", nil)
	}

	// 根据单位转换大小
	var size float64
	switch unit {
	case "B", "b", "":
		size = num
	case "KB", "kb", "K", "k":
		size = num * float64(types.KB)
	case "MB", "mb", "M", "m":
		size = num * float64(types.MB)
	case "GB", "gb", "G", "g":
		size = num * float64(types.GB)
	case "TB", "tb", "T", "t":
		size = num * float64(types.TB)
	case "PB", "pb", "P", "p":
		size = num * float64(types.PB)
	case "KIB", "kib":
		size = num * float64(types.KIB)
	case "MIB", "mib":
		size = num * float64(types.MIB)
	case "GIB", "gib":
		size = num * float64(types.GIB)
	case "TIB", "tib":
		size = num * float64(types.TIB)
	case "PIB", "pib":
		size = num * float64(types.PIB)
	default:
		return types.NewError("INVALID_SIZE", fmt.Sprintf("unknown size unit: %s", unit), nil)
	}

	// 检查是否超出int64范围
	if size > float64(1<<63-1) {
		return types.NewError("INVALID_SIZE", fmt.Sprintf("size value too large: %s", value), nil)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(int64(size)); err != nil {
			return err
		}
	}

	// 设置值
	*f.value = int64(size)
	f.isSet = true

	return nil
}

// String 返回格式化的大小字符串
//
// 返回值:
//   - string: 格式化的大小字符串, 如 "1024B"、"2.5MB" 等
//
// 注意事项:
//   - 此方法是线程安全的
//   - 实现了fmt.Stringer接口
func (f *SizeFlag) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return utils.FormatSize(*f.value)
}
