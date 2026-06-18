// Package shck 提供 Shell 脚本语法检查功能
package shck

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gitee.com/MM-Q/color"

	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/go-kit/term"
	"gitee.com/MM-Q/shx"
)

// ShckConfig shck 命令配置
type ShckConfig struct {
	Files []string // 位置参数（文件路径或通配符）
}

// ShckCmdMain shck 命令主入口
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func ShckCmdMain(config ShckConfig) error {
	// 优先处理标准输入（管道模式）
	if term.IsStdinPipe() {
		return checkStdin()
	}

	// 展开通配符获取实际文件列表
	files, err := fs.ExpandFiles(config.Files)
	if err != nil {
		return fmt.Errorf("expand files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files specified")
	}

	// 检查每个文件的语法
	var errs []string
	for _, file := range files {
		if err := checkFile(file); err != nil {
			errs = append(errs, fmt.Sprintf("\t✗ %s: %s", file, err.Error()))
		} else {
			_, _ = color.New(color.FgGreen, color.Bold).Printf("✓ %s — syntax is valid\n", file)
		}
	}

	// 如果有错误，返回聚合错误
	if len(errs) > 0 {
		return fmt.Errorf("syntax check failed:\n%s", strings.Join(errs, "\n"))
	}

	return nil
}

// checkStdin 检查标准输入的语法
func checkStdin() error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("read stdin: %w", err)
	}

	if err := shx.CheckSyntax(string(data)); err != nil {
		return err
	}

	_, _ = color.New(color.FgGreen, color.Bold).Println("✓ syntax is valid")
	return nil
}

// checkFile 检查单个文件的语法
func checkFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return shx.CheckSyntax(string(data))
}
