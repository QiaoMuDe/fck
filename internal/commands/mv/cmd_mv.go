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

// MvStats 操作统计
type MvStats struct {
	Moved  int
	Errors int
}

// MvCmdMain 主函数
func MvCmdMain(config MvConfig) error {
	if len(config.Sources) == 0 {
		return fmt.Errorf("未指定要移动的源文件")
	}

	if config.Target == "" {
		return fmt.Errorf("未指定目标路径")
	}

	stats := &MvStats{}

	for _, source := range config.Sources {
		err := moveItem(source, config.Target, config.Force, config.Interactive, config.Verbose, stats)
		if err != nil {
			stats.Errors++
			return err
		}
	}

	if config.Verbose {
		fmt.Printf("操作完成: %d 个移动", stats.Moved)
		if stats.Errors > 0 {
			fmt.Printf(", %d 个错误", stats.Errors)
		}
		fmt.Println()
	}

	return nil
}

// moveItem 移动文件或目录
func moveItem(src, dst string, force, interactive, verbose bool, stats *MvStats) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	dstInfo, err := os.Stat(dst)
	if err == nil {
		if dstInfo.IsDir() {
			dst = fmt.Sprintf("%s%c%s", dst, os.PathSeparator, srcInfo.Name())
		} else {
			if interactive {
				if !confirmOverwrite(dst) {
					if verbose {
						fmt.Printf("跳过: %s (用户取消)\n", dst)
					}
					return nil
				}
			} else if !force {
				return fmt.Errorf("目标文件已存在: %s", dst)
			}
		}
	}

	if err := fs.MoveEx(src, dst, force); err != nil {
		return fmt.Errorf("移动 '%s' 到 '%s' 失败: %v", src, dst, err)
	}

	if verbose {
		fmt.Printf("移动: %s -> %s\n", src, dst)
	}

	stats.Moved++
	return nil
}

// confirmOverwrite 交互式覆盖确认
func confirmOverwrite(dst string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("覆盖 %s? (y/n): ", dst)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
