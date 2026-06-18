// Package shfmt 实现 Shell 脚本格式化功能
// 支持标准输入和文件两种模式，可原地写入并备份原文件
package shfmt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/go-kit/term"
	"gitee.com/MM-Q/shx"
)

// ShfmtConfig shfmt 命令配置
type ShfmtConfig struct {
	Write         bool              // -w: 原地写入
	Backup        bool              // -b: 写入前备份
	List          bool              // -l: 列表模式
	Files         []string          // 位置参数（文件路径或通配符）
	FormatOptions shx.FormatOptions // 格式化选项
}

// ShfmtCmdMain 执行 shfmt 命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误
func ShfmtCmdMain(config ShfmtConfig) error {
	// 优先处理标准输入（管道模式）
	if term.IsStdinPipe() {
		return processStdin(config.FormatOptions)
	}

	// 展开通配符，获取实际文件列表
	files, err := fs.ExpandFiles(config.Files)
	if err != nil {
		return fmt.Errorf("failed to expand file patterns: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files matched")
	}

	// 逐文件处理
	for _, file := range files {
		if err := processFile(file, config, config.List); err != nil {
			return fmt.Errorf("error processing %s: %w", file, err)
		}
	}

	return nil
}

// processStdin 处理标准输入
// 从 stdin 读取 Shell 脚本内容，格式化后输出到 stdout
//
// 参数:
//   - opts: 格式化选项
//
// 返回:
//   - error: 处理错误
func processStdin(opts shx.FormatOptions) error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	formatted, err := shx.FormatWithOptions(string(data), opts)
	if err != nil {
		return fmt.Errorf("failed to format script: %w", err)
	}

	fmt.Print(formatted)
	return nil
}

// processFile 处理单个文件
// 读取文件内容，格式化后根据配置选择输出到 stdout、写回文件或列出差异
//
// 参数:
//   - file: 文件路径
//   - config: 命令配置
//   - listMode: 列表模式标志（仅列出需要格式化的文件路径）
//
// 返回:
//   - error: 处理错误
func processFile(file string, config ShfmtConfig, listMode bool) error {
	// 读取文件内容
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 格式化
	formatted, err := shx.FormatWithOptions(string(data), config.FormatOptions)
	if err != nil {
		return fmt.Errorf("failed to format script: %w", err)
	}

	// 列表模式：比较差异，有变化则输出文件路径（只读，不修改文件）
	if listMode {
		if string(data) != formatted {
			fmt.Println(file)
		}
		return nil
	}

	// 预览模式：输出到 stdout
	if !config.Write {
		fmt.Print(formatted)
		return nil
	}

	// 原地写入模式
	// 如果需要备份，先创建备份
	if config.Backup {
		backupPath := file + ".bak"
		if err := fs.CopyEx(file, backupPath, true, false); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// 原子写入：先写临时文件再 rename
	dir := filepath.Dir(file)
	tmpFile, err := os.CreateTemp(dir, "shfmt-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// 写入格式化内容
	if _, err := tmpFile.WriteString(formatted); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// 关闭临时文件
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// 原子替换原文件
	if err := fs.MoveEx(tmpPath, file, true, false); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to replace file: %w", err)
	}

	_, _ = color.New(color.FgGreen, color.Bold).Printf("✓ %s — formatted\n", file)
	return nil
}
