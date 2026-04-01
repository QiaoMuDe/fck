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
		return fmt.Errorf("no targets specified to operate")
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
		fmt.Printf("operation completed: %d created files, %d updated files", stats.Created, stats.Updated)
		if stats.Skipped > 0 {
			fmt.Printf(", %d skipped", stats.Skipped)
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
					fmt.Printf("skip: %s (does not exist)\n", path)
				}
				stats.Skipped++
				return nil
			}
			if err := createFile(path, config, stats); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("error checking file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("not a file: %s", path)
	}

	return updateFileTime(path, config, targetTime, stats)
}

// createFile 创建文件
func createFile(path string, config TouchConfig, stats *TouchStats) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	_ = file.Close()

	if config.Verbose {
		fmt.Printf("created file: %s\n", path)
	}

	stats.Created++
	return nil
}

// updateFileTime 更新文件时间
func updateFileTime(path string, config TouchConfig, targetTime time.Time, stats *TouchStats) error {
	var atime, mtime time.Time

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error checking file: %w", err)
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
		return fmt.Errorf("error updating file time: %w", err)
	}

	if config.Verbose {
		fmt.Printf("updated time: %s\n", path)
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
		return time.Time{}, fmt.Errorf("time format error, should be: YYYYMMDDHHMM.SS")
	}

	year, err := strconv.Atoi(timeStr[0:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("year parse error: %w", err)
	}

	month, err := strconv.Atoi(timeStr[4:6])
	if err != nil {
		return time.Time{}, fmt.Errorf("month parse error: %w", err)
	}

	day, err := strconv.Atoi(timeStr[6:8])
	if err != nil {
		return time.Time{}, fmt.Errorf("day parse error: %w", err)
	}

	hour, err := strconv.Atoi(timeStr[8:10])
	if err != nil {
		return time.Time{}, fmt.Errorf("hour parse error: %w", err)
	}

	minute, err := strconv.Atoi(timeStr[10:12])
	if err != nil {
		return time.Time{}, fmt.Errorf("minute parse error: %w", err)
	}

	second := 0
	if len(timeStr) >= 14 {
		second, err = strconv.Atoi(timeStr[12:14])
		if err != nil {
			return time.Time{}, fmt.Errorf("second parse error: %w", err)
		}
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local), nil
}
