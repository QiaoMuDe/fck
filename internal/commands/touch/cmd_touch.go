package touch

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// TouchConfig 配置结构体
type TouchConfig struct {
	Targets  []string
	NoCreate bool
	Time     string
	Access   bool
	Modify   bool
	Verbose  bool
}

// TouchStats 操作统计
type TouchStats struct {
	Created int
	Updated int
	Skipped int
	Errors  int
}

// TouchCmdMain 主函数
func TouchCmdMain(config TouchConfig) error {
	if len(config.Targets) == 0 {
		return fmt.Errorf("未指定要操作的文件")
	}

	stats := &TouchStats{}
	targetTime, err := parseTime(config.Time)
	if err != nil {
		return err
	}

	for _, target := range config.Targets {
		err := touchFile(target, config, targetTime, stats)
		if err != nil {
			stats.Errors++
			return err
		}
	}

	if config.Verbose {
		fmt.Printf("操作完成: %d 个创建, %d 个更新", stats.Created, stats.Updated)
		if stats.Skipped > 0 {
			fmt.Printf(", %d 个跳过", stats.Skipped)
		}
		fmt.Println()
	}

	return nil
}

// touchFile 处理单个文件
func touchFile(path string, config TouchConfig, targetTime time.Time, stats *TouchStats) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if config.NoCreate {
				if config.Verbose {
					fmt.Printf("跳过: %s (不存在)\n", path)
				}
				stats.Skipped++
				return nil
			}
			if err := createFile(path, config, stats); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("访问文件失败: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("路径是目录，不是文件: %s", path)
	}

	return updateFileTime(path, config, targetTime, stats)
}

// createFile 创建文件
func createFile(path string, config TouchConfig, stats *TouchStats) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	_ = file.Close()

	if config.Verbose {
		fmt.Printf("创建文件: %s\n", path)
	}

	stats.Created++
	return nil
}

// updateFileTime 更新文件时间
func updateFileTime(path string, config TouchConfig, targetTime time.Time, stats *TouchStats) error {
	var atime, mtime time.Time

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	if config.Access && !config.Modify {
		atime = targetTime
		mtime = info.ModTime()
	} else if config.Modify && !config.Access {
		atime = info.ModTime()
		mtime = targetTime
	} else {
		atime = targetTime
		mtime = targetTime
	}

	if err := os.Chtimes(path, atime, mtime); err != nil {
		return fmt.Errorf("更新文件时间失败: %w", err)
	}

	if config.Verbose {
		fmt.Printf("更新时间: %s\n", path)
	}

	stats.Updated++
	return nil
}

// parseTime 解析时间字符串
func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Now(), nil
	}

	timeStr = strings.ReplaceAll(timeStr, ".", "")
	if len(timeStr) < 12 {
		return time.Time{}, fmt.Errorf("时间格式错误，应为 YYYYMMDDHHMM.SS")
	}

	year, err := strconv.Atoi(timeStr[0:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("年份解析失败: %w", err)
	}

	month, err := strconv.Atoi(timeStr[4:6])
	if err != nil {
		return time.Time{}, fmt.Errorf("月份解析失败: %w", err)
	}

	day, err := strconv.Atoi(timeStr[6:8])
	if err != nil {
		return time.Time{}, fmt.Errorf("日期解析失败: %w", err)
	}

	hour, err := strconv.Atoi(timeStr[8:10])
	if err != nil {
		return time.Time{}, fmt.Errorf("小时解析失败: %w", err)
	}

	minute, err := strconv.Atoi(timeStr[10:12])
	if err != nil {
		return time.Time{}, fmt.Errorf("分钟解析失败: %w", err)
	}

	second := 0
	if len(timeStr) >= 14 {
		second, err = strconv.Atoi(timeStr[12:14])
		if err != nil {
			return time.Time{}, fmt.Errorf("秒数解析失败: %w", err)
		}
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local), nil
}
