package chmod

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// ChmodConfig 配置结构体
type ChmodConfig struct {
	Mode      string
	Targets   []string
	Recursive bool
	Verbose   bool
	Changes   bool
	Quiet     bool
}

// ChmodStats 操作统计
type ChmodStats struct {
	Modified int
	Skipped  int
	Errors   int
}

// ChmodCmdMain 主函数
func ChmodCmdMain(config ChmodConfig) error {
	if len(config.Targets) == 0 {
		return fmt.Errorf("未指定要修改权限的目标")
	}

	if config.Mode == "" {
		return fmt.Errorf("未指定权限模式")
	}

	mode, err := parseMode(config.Mode)
	if err != nil {
		return err
	}

	stats := &ChmodStats{}

	for _, target := range config.Targets {
		err := changeMode(target, config, mode, stats)
		if err != nil {
			stats.Errors++
			if !config.Quiet {
				return err
			}
		}
	}

	if config.Verbose {
		fmt.Printf("操作完成: %d 个修改", stats.Modified)
		if stats.Skipped > 0 {
			fmt.Printf(", %d 个跳过", stats.Skipped)
		}
		if stats.Errors > 0 {
			fmt.Printf(", %d 个错误", stats.Errors)
		}
		fmt.Println()
	}

	return nil
}

// changeMode 修改权限
func changeMode(path string, config ChmodConfig, mode os.FileMode, stats *ChmodStats) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if config.Quiet {
				return nil
			}
			return fmt.Errorf("路径不存在: %s", path)
		}
		return fmt.Errorf("访问路径失败: %w", err)
	}

	currentMode := info.Mode().Perm()

	if currentMode == mode {
		if config.Verbose {
			fmt.Printf("跳过: %s (权限已为 %04o)\n", path, currentMode)
		}
		stats.Skipped++
		return nil
	}

	if err := os.Chmod(path, mode); err != nil {
		return fmt.Errorf("修改权限失败: %w", err)
	}

	if config.Verbose || config.Changes {
		fmt.Printf("修改权限: %s (%04o -> %04o)\n", path, currentMode, mode)
	}

	stats.Modified++

	if info.IsDir() && config.Recursive {
		return filepath.WalkDir(path, func(walkPath string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				if config.Quiet {
					return nil
				}
				return fmt.Errorf("访问路径失败: %w", walkErr)
			}

			if walkPath == path {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				if config.Quiet {
					return nil
				}
				return fmt.Errorf("获取文件信息失败: %w", err)
			}

			if info.Mode()&os.ModeSymlink != 0 {
				if config.Verbose {
					fmt.Printf("跳过: %s (符号链接)\n", walkPath)
				}
				return nil
			}

			currentMode := info.Mode().Perm()

			if currentMode != mode {
				if err := os.Chmod(walkPath, mode); err != nil {
					if config.Quiet {
						return nil
					}
					return fmt.Errorf("修改权限失败: %w", err)
				}

				if config.Verbose || config.Changes {
					fmt.Printf("修改权限: %s (%04o -> %04o)\n", walkPath, currentMode, mode)
				}
				stats.Modified++
			} else {
				if config.Verbose {
					fmt.Printf("跳过: %s (权限已为 %04o)\n", walkPath, currentMode)
				}
				stats.Skipped++
			}

			return nil
		})
	}

	return nil
}

// parseMode 解析权限模式
func parseMode(modeStr string) (os.FileMode, error) {
	mode, err := strconv.ParseInt(modeStr, 8, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的权限模式: %s", modeStr)
	}

	if mode < 0 || mode > 0777 {
		return 0, fmt.Errorf("权限模式超出范围 (0-0777): %s", modeStr)
	}

	return os.FileMode(mode), nil
}
