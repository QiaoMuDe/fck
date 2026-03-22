package types

import (
	"fmt"
	"time"
)

// 常见时间格式常量
var (
	// RFC3339 RFC3339 格式 (2006-01-02T15:04:05Z07:00)
	TimeFormatRFC3339 = time.RFC3339

	// RFC3339Nano RFC3339 纳秒格式 (2006-01-02T15:04:05.999999999Z07:00)
	TimeFormatRFC3339Nano = time.RFC3339Nano

	// RFC1123 RFC1123 格式 (Mon, 02 Jan 2006 15:04:05 MST)
	TimeFormatRFC1123 = time.RFC1123

	// RFC1123Z RFC1123 带时区格式 (Mon, 02 Jan 2006 15:04:05 -0700)
	TimeFormatRFC1123Z = time.RFC1123Z

	// RFC822 RFC822 格式 (02 Jan 06 15:04 MST)
	TimeFormatRFC822 = time.RFC822

	// RFC822Z RFC822 带时区格式 (02 Jan 06 15:04 -0700)
	TimeFormatRFC822Z = time.RFC822Z

	// Kitchen 厨房格式 (3:04PM)
	TimeFormatKitchen = time.Kitchen

	// Stamp 简单时间戳格式 (Jan _2 15:04:05)
	TimeFormatStamp = time.Stamp

	// StampMilli 毫秒时间戳格式 (Jan _2 15:04:05.000)
	TimeFormatStampMilli = time.StampMilli

	// StampMicro 微秒时间戳格式 (Jan _2 15:04:05.000000)
	TimeFormatStampMicro = time.StampMicro

	// StampNano 纳秒时间戳格式 (Jan _2 15:04:05.000000000)
	TimeFormatStampNano = time.StampNano

	// DateOnly 日期格式 (2006-01-02)
	TimeFormatDateOnly = "2006-01-02"

	// TimeOnly 时间格式 (15:04:05)
	TimeFormatTimeOnly = "15:04:05"

	// DateTime 日期时间格式 (2006-01-02 15:04:05)
	TimeFormatDateTime = "2006-01-02 15:04:05"

	// DateTimeSlash 斜杠分隔的日期时间格式 (01/02/2006 15:04:05)
	TimeFormatDateTimeSlash = "01/02/2006 15:04:05"

	// DateTimeCompact 紧凑日期时间格式 (20060102150405)
	TimeFormatDateTimeCompact = "20060102150405"

	// ISO8601 ISO8601 格式 (2006-01-02T15:04:05Z)
	TimeFormatISO8601 = "2006-01-02T15:04:05Z"

	// ISO8601Nano ISO8601 纳秒格式 (2006-01-02T15:04:05.999999999Z)
	TimeFormatISO8601Nano = "2006-01-02T15:04:05.999999999Z"
)

// CommonTimeFormats 常见时间格式列表, 按优先级排序
var CommonTimeFormats = []string{
	TimeFormatRFC3339,
	TimeFormatRFC3339Nano,
	TimeFormatISO8601,
	TimeFormatISO8601Nano,
	TimeFormatDateTime,
	TimeFormatDateOnly,
	TimeFormatTimeOnly,
	TimeFormatRFC1123,
	TimeFormatRFC1123Z,
	TimeFormatDateTimeSlash,
	TimeFormatStamp,
	TimeFormatStampMilli,
	TimeFormatStampMicro,
	TimeFormatStampNano,
	TimeFormatRFC822,
	TimeFormatRFC822Z,
	TimeFormatKitchen,
	TimeFormatDateTimeCompact,
}

// ParseTimeWithFormats 尝试使用多种格式解析时间字符串
//
// 参数:
//   - value: 要解析的时间字符串
//   - formats: 要尝试的时间格式列表, 按优先级排序
//
// 返回值:
//   - time.Time: 解析后的时间
//   - string: 使用的时间格式
//   - error: 如果解析失败返回错误
//
// 功能说明:
//   - 按给定格式列表顺序尝试解析
//   - 返回第一个成功解析的时间和格式
//   - 如果所有格式都失败, 返回错误
func ParseTimeWithFormats(value string, formats []string) (time.Time, string, error) {
	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			return t, format, nil
		}
	}
	return time.Time{}, "", NewError("INVALID_TIME",
		fmt.Sprintf("unable to parse time value: %s with any known format", value),
		nil)
}

// ParseTimeWithCommonFormats 尝试使用常见格式解析时间字符串
//
// 函数功能:
//   - 尝试使用常见格式解析时间字符串
//
// 参数:
//   - value: 要解析的时间字符串
//
// 返回值:
//   - time.Time: 解析后的时间
//   - string: 使用的时间格式
//   - error: 如果解析失败返回错误
//
// 说明:
//   - 使用常见时间格式列表进行解析
//   - 返回第一个成功解析的时间和格式
func ParseTimeWithCommonFormats(value string) (time.Time, string, error) {
	return ParseTimeWithFormats(value, CommonTimeFormats)
}
