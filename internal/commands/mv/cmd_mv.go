package mv

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitee.com/MM-Q/go-kit/fs"
)

// MvConfig 配置结构体
type MvConfig struct {
	Force       bool
	Interactive bool
	Verbose     bool
	Sources     []string
	Target      string
}

// MvCmdMain 主函数
func MvCmdMain(config MvConfig) error {
	if len(config.Sources) == 0 {
		return fmt.Errorf("no source files specified")
	}

	if config.Target == "" {
		return fmt.Errorf("no target path specified")
	}

	for _, source := range config.Sources {
		if err := moveItem(source, config.Target, config.Force, config.Interactive, config.Verbose); err != nil {
			return err
		}
	}

	return nil
}

// moveItem 移动文件或目录
func moveItem(src, dst string, force, interactive, verbose bool) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("get source file info failed: %w", err)
	}

	dstInfo, err := os.Stat(dst)
	if err == nil {
		if dstInfo.IsDir() {
			dst = fmt.Sprintf("%s%c%s", dst, os.PathSeparator, srcInfo.Name())
		} else {
			if interactive {
				if !confirmOverwrite(dst) {
					if verbose {
						fmt.Printf("skipped: %s (user cancel)\n", dst)
					}
					return nil
				}
			} else if !force {
				return fmt.Errorf("target file exists and is not a directory: %s", dst)
			}
		}
	}

	if err := fs.MoveEx(src, dst, force, verbose); err != nil {
		return fmt.Errorf("move '%s' to '%s' failed: %v", src, dst, err)
	}

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
