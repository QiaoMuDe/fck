// Package iconv 提供文件编码检测和转换功能
package iconv

import (
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/fck/internal/types"
)

// IconvConfig iconv 命令配置
type IconvConfig struct {
	// 输入参数
	Files        []string // 目标文件列表（支持通配符和多个文件）
	FromEncoding string   // 源编码（auto 表示自动检测）
	ToEncoding   string   // 目标编码

	// 输出选项
	Write  bool   // 原地写入（覆盖原文件）
	Backup bool   // 原地写入时创建备份
	Output string // 输出路径（目录或文件，默认当前目录）

	// 其他选项
	List      bool // 列出支持的编码
	Quiet     bool // 静默模式
	FileCount int  // 文件总数（用于控制输出格式）
}

// ConversionResult 转换结果
type ConversionResult struct {
	FilePath     string // 文件路径
	Success      bool   // 是否成功
	FromEncoding string // 检测/指定的源编码
	ToEncoding   string // 目标编码
	BytesRead    int64  // 读取字节数
	BytesWritten int64  // 写入字节数
	Error        error  // 错误信息（如果有）
}

// IconvCmdMain iconv 命令主入口
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误（如果有）
func IconvCmdMain(config IconvConfig) error {
	// 列出支持的编码
	if config.List {
		return ListEncodings()
	}

	// 检查文件参数
	if len(config.Files) == 0 {
		return fmt.Errorf("no input files specified")
	}

	// 展开文件列表（处理通配符）
	files, err := expandFileList(config.Files)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found matching the specified patterns")
	}

	// 设置文件总数（用于控制输出格式）
	config.FileCount = len(files)

	// 统计结果
	var stats struct {
		total   int // 总文件数
		success int // 成功转换文件数
		skipped int // 跳过文件数（检测模式或未转换）
		failed  int // 失败文件数（转换错误）
	}

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
			if config.ToEncoding == types.EncodingNone {
				// 检测模式不算成功转换
				stats.skipped++
			} else if result.FromEncoding == result.ToEncoding {
				stats.skipped++
			} else {
				stats.success++
				if !config.Quiet {
					fmt.Printf("%s: %s → %s (%d bytes)\n",
						file, result.FromEncoding, result.ToEncoding, result.BytesRead)
				}
			}
		}
	}

	// 输出统计信息（转换模式、非静默模式且多个文件）
	if !config.Quiet && config.ToEncoding != types.EncodingNone && stats.total > 1 {
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

// ListEncodings 列出支持的编码格式
//
// 返回值:
//   - error: 执行错误（如果有）
func ListEncodings() error {
	for _, enc := range types.TargetEncodings {
		fmt.Println(enc)
	}
	return nil
}

// expandFileList 展开文件列表（处理通配符）
//
// 参数:
//   - patterns: 文件路径模式列表
//
// 返回值:
//   - []string: 展开后的文件路径列表
//   - error: 执行错误（如果有）
func expandFileList(patterns []string) ([]string, error) {
	var files []string
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %s: %w", pattern, err)
		}

		if len(matches) == 0 {
			// 如果没有匹配到，检查是否是目录
			info, err := os.Stat(pattern)
			if err == nil && info.IsDir() {
				return nil, fmt.Errorf("%s is a directory, not a file", pattern)
			}
			// 保留原路径（可能是具体文件）
			matches = []string{pattern}
		}

		for _, match := range matches {
			// 跳过目录
			info, err := os.Stat(match)
			if err == nil && info.IsDir() {
				continue
			}
			if !seen[match] {
				seen[match] = true
				files = append(files, match)
			}
		}
	}

	return files, nil
}
