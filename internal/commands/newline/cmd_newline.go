// Package newline 提供文件换行符检测和转换功能
package newline

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/fs"
)

// ListTypes 列出支持的换行符类型
func ListTypes() error {
	for _, t := range types.TargetNewlines {
		fmt.Println(t)
	}
	return nil
}

// CmdMain newline 命令主入口
//
// 参数:
//   - config: 转换配置
//
// 返回值:
//   - error: 执行错误（如果有）
func CmdMain(config Config) error {
	// 展开文件列表（支持通配符）
	files, err := fs.ExpandFiles(config.Files)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no files to process")
	}

	// 设置文件总数（用于控制输出格式）
	config.FileCount = len(files)

	// 统计信息
	stats := struct {
		total   int
		success int
		skipped int
		failed  int
	}{}

	// 处理每个文件
	for _, file := range files {
		result, err := ConvertFile(file, config)
		stats.total++

		if err != nil {
			stats.failed++
			if !config.Quiet {
				_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", file, err)
			}
			continue
		}

		if result.Success {
			if config.ToNewline == types.NewlineNone {
				// 检测模式不算成功转换
				stats.skipped++
			} else if result.FromType == result.ToType {
				stats.skipped++
			} else {
				stats.success++
				if !config.Quiet {
					fmt.Printf("%s: %s → %s (%d lines)\n",
						file, result.FromType, result.ToType, result.Lines)
				}
			}
		}
	}

	// 输出统计信息（转换模式、非静默模式且多个文件）
	if !config.Quiet && config.ToNewline != types.NewlineNone && stats.total > 1 {
		fmt.Printf("\nTotal: %d files\n", stats.total)
		if stats.success > 0 {
			fmt.Printf("Success: %d files\n", stats.success)
		}
		if stats.skipped > 0 {
			fmt.Printf("Skipped: %d files\n", stats.skipped)
		}
		if stats.failed > 0 {
			fmt.Printf("Failed: %d files\n", stats.failed)
		}
	}

	if stats.failed > 0 {
		return fmt.Errorf("conversion completed with %d errors", stats.failed)
	}

	return nil
}
