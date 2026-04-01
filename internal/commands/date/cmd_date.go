package date

import (
	"fmt"
	"strconv"
	"time"
)

var formatAliases = map[string]string{
	"iso":      "2006-01-02T15:04:05Z07:00",
	"date":     "2006-01-02",
	"time":     "15:04:05",
	"datetime": "2006-01-02 15:04:05",
	"compact":  "2006-01-02T15:04:05",
	"cn-date":  "2006年01月02日",
	"cn-time":  "15时04分05秒",
	"cn":       "2006年01月02日 15时04分05秒",
	"us-date":  "01/02/2006",
	"us-time":  "03:04:05 PM",
	"us":       "01/02/2006 03:04:05 PM",
	"eu-date":  "02.01.2006",
	"eu-time":  "15:04:05",
	"eu":       "02.01.2006 15:04:05",
	"log":      "2006-01-02 15:04:05.000",
	"filename": "20060102_150405",
	"num":      "20060102150405",
	"dnum":     "20060102",
	"tnum":     "150405",
}

type DateConfig struct {
	Format    string
	Timestamp string
	Timezone  string
	UTC       bool
	Unix      bool
}

func DateCmdMain(config DateConfig) error {
	var t time.Time
	var err error

	// 解析时间戳
	if config.Timestamp != "" {
		t, err = parseTimestamp(config.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}
	} else {
		t = time.Now()
	}

	// 设置时区
	loc, err := getLocation(config)
	if err != nil {
		return fmt.Errorf("failed to set timezone location: %w", err)
	}

	// 转换时区
	t = t.In(loc)

	// 如果是Unix时间戳，直接输出
	if config.Unix {
		fmt.Println(t.Unix())
		return nil
	}

	// 格式化输出
	format := getFormat(config)
	fmt.Println(t.Format(format))

	return nil
}

// parseTimestamp 解析时间戳字符串为time.Time类型
//
// 参数:
//   - timestamp: 时间戳字符串
//
// 返回值:
//   - time.Time: 解析后的时间对象
//   - error: 解析失败时返回的错误
func parseTimestamp(timestamp string) (time.Time, error) {
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp format: %w", err)
	}

	if ts > 1e12 {
		return time.UnixMilli(ts), nil
	}

	return time.Unix(ts, 0), nil
}

// getLocation 获取时区
//
// 参数:
//   - config: 配置结构体
//
// 返回值:
//   - *time.Location: 时区对象
//   - error: 加载时区失败时返回的错误
func getLocation(config DateConfig) (*time.Location, error) {
	if config.UTC {
		return time.UTC, nil
	}

	if config.Timezone != "" {
		return time.LoadLocation(config.Timezone)
	}

	return time.Local, nil
}

// getFormat 获取格式化字符串
//
// 参数:
//   - config: 配置结构体
//
// 返回值:
//   - string: 格式化字符串
//   - error: 格式化字符串不存在时返回的错误
func getFormat(config DateConfig) string {
	if config.Format == "" {
		return "2006-01-02 15:04:05"
	}

	if format, ok := formatAliases[config.Format]; ok {
		return format
	}

	return config.Format
}
