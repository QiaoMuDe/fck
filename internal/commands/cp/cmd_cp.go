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
	Recursive   bool
	Sources     []string
	Target      string
}

// CpCmdMain 主函数
func CpCmdMain(config CpConfig) error {
	if len(config.Sources) == 0 {
		return fmt.Errorf("no source files specified to copy")
	}

	if config.Target == "" {
		return fmt.Errorf("no destination path specified")
	}

	for _, source := range config.Sources {
		if err := copyItem(source, config.Target, config.Force, config.Interactive, config.Verbose, config.Recursive); err != nil {
			return err
		}
	}

	return nil
}

// copyItem 复制文件或目录
//
// 参数:
//   - src: 源文件或目录路径
//   - dst: 目标文件或目录路径
//   - force: 是否强制覆盖已存在的文件
//   - interactive: 是否交互式覆盖（覆盖前提示）
//   - verbose: 是否显示复制的文件/目录
//   - recursive: 是否递归复制目录
//
// 返回值:
//   - error: 复制失败时返回错误，否则返回 nil
func copyItem(src, dst string, force, interactive, verbose bool, recursive bool) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	// 检查是否需要递归复制目录
	if srcInfo.IsDir() && !recursive {
		return fmt.Errorf("%s is a directory, use -r to copy recursively", src)
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

	if err := fs.CopyEx(src, dst, force, verbose); err != nil {
		return fmt.Errorf("copy failed: '%s' to '%s': %v", src, dst, err)
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
