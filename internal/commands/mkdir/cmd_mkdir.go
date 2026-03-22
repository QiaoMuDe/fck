package mkdir

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// MkdirConfig 配置结构体
type MkdirConfig struct {
	Targets []string
	Parents bool
	Mode    int
	Verbose bool
}

// CreateStats 创建统计
type CreateStats struct {
	Directories int
	Skipped     int
	Errors      int
}

// MkdirCmdMain 主函数
func MkdirCmdMain(config MkdirConfig) error {
	if len(config.Targets) == 0 {
		return fmt.Errorf("未指定要创建的目录")
	}

	stats := &CreateStats{}

	for _, target := range config.Targets {
		err := createDirectory(target, config, stats)
		if err != nil {
			stats.Errors++
			return err
		}
	}

	if config.Verbose {
		fmt.Printf("创建完成: %d 个目录", stats.Directories)
		if stats.Skipped > 0 {
			fmt.Printf(", %d 个跳过", stats.Skipped)
		}
		fmt.Println()
	}

	return nil
}

// createDirectory 创建目录
func createDirectory(path string, config MkdirConfig, stats *CreateStats) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %w", err)
	}

	info, err := os.Stat(absPath)
	if err == nil {
		if info.IsDir() {
			if config.Verbose {
				fmt.Printf("跳过: %s (已存在)\n", path)
			}
			stats.Skipped++
			return nil
		}
		return fmt.Errorf("路径已存在且不是目录: %s", path)
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("检查路径失败: %w", err)
	}

	if config.Parents {
		if err := os.MkdirAll(absPath, os.FileMode(config.Mode)); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	} else {
		if err := os.Mkdir(absPath, os.FileMode(config.Mode)); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	if config.Verbose {
		fmt.Printf("创建目录: %s\n", path)
	}

	stats.Directories++
	return nil
}

// ParseMode 解析权限模式
func ParseMode(modeStr string) (int, error) {
	if modeStr == "" {
		return 0755, nil
	}

	mode, err := strconv.ParseInt(modeStr, 8, 64)
	if err != nil {
		return 0755, fmt.Errorf("无效的权限模式: %s", modeStr)
	}

	if mode < 0 || mode > 0777 {
		return 0755, fmt.Errorf("权限模式超出范围 (0-0777): %s", modeStr)
	}

	return int(mode), nil
}
