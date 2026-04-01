package cp

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitee.com/MM-Q/go-kit/fs"
)

// CpConfig 配置结构体
type CpConfig struct {
	Force       bool
	Interactive bool
	Verbose     bool
	Sources     []string
	Target      string
}

// CpStats 操作统计
type CpStats struct {
	Copied int
	Errors int
}

// CpCmdMain 主函数
func CpCmdMain(config CpConfig) error {
	if len(config.Sources) == 0 {
		return fmt.Errorf("no source files specified to copy")
	}

	if config.Target == "" {
		return fmt.Errorf("no destination path specified")
	}

	stats := &CpStats{}

	for _, source := range config.Sources {
		err := copyItem(source, config.Target, config.Force, config.Interactive, config.Verbose, stats)
		if err != nil {
			stats.Errors++
			return err
		}
	}

	if config.Verbose {
		fmt.Printf("operation completed: %d copies\n", stats.Copied)
		if stats.Errors > 0 {
			fmt.Printf(", %d errors", stats.Errors)
		}
		fmt.Println()
	}

	return nil
}

// copyItem 复制文件或目录
func copyItem(src, dst string, force, interactive, verbose bool, stats *CpStats) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	dstInfo, err := os.Stat(dst)
	if err == nil {
		if dstInfo.IsDir() {
			dst = fmt.Sprintf("%s%c%s", dst, os.PathSeparator, srcInfo.Name())
		} else {
			if interactive {
				if !confirmOverwrite(dst) {
					if verbose {
						fmt.Printf("skipped: %s (user cancelled)\n", dst)
					}
					return nil
				}
			} else if !force {
				return fmt.Errorf("destination file already exists: %s", dst)
			}
		}
	}

	if err := fs.CopyEx(src, dst, force); err != nil {
		return fmt.Errorf("copy failed: '%s' to '%s': %v", src, dst, err)
	}

	if verbose {
		fmt.Printf("copy: %s -> %s\n", src, dst)
	}

	stats.Copied++
	return nil
}

// confirmOverwrite 交互式覆盖确认
func confirmOverwrite(dst string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("overwrite %s? (y/n): ", dst)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
