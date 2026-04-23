// Package tee 实现了 tee 命令
// 从标准输入读取数据，同时输出到标准输出和指定的文件
package tee

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// TeeConfig 配置结构体
type TeeConfig struct {
	Append           bool     // 追加模式
	IgnoreInterrupts bool     // 忽略中断信号
	Files            []string // 输出文件列表
}

// TeeCmdMain 执行 tee 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func TeeCmdMain(config TeeConfig) error {
	// 如果没有文件参数，直接复制到 stdout
	if len(config.Files) == 0 {
		_, err := io.Copy(os.Stdout, os.Stdin)
		return err
	}

	// 打开所有输出文件
	files, err := openFiles(config.Files, config.Append)
	if err != nil {
		return err
	}
	defer closeFiles(files)

	// 设置信号处理（如需要）
	if config.IgnoreInterrupts {
		setupSignalHandler()
	}

	// 构建 writer 列表：stdout + 所有文件
	writers := make([]io.Writer, 0, len(files)+1)
	writers = append(writers, os.Stdout)
	for _, f := range files {
		writers = append(writers, f)
	}

	// 使用标准库 io.MultiWriter 同时写入所有目标
	_, err = io.Copy(io.MultiWriter(writers...), os.Stdin)
	return err
}

// openFiles 打开所有输出文件
//
// 参数:
//   - paths: 文件路径列表
//   - isAppend: 是否追加模式
//
// 返回值:
//   - []*os.File: 打开的文件句柄列表
//   - error: 打开错误
func openFiles(paths []string, isAppend bool) ([]*os.File, error) {
	flag := os.O_CREATE | os.O_WRONLY
	if isAppend {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	var files []*os.File
	for _, path := range paths {
		f, err := os.OpenFile(path, flag, 0644)
		if err != nil {
			closeFiles(files)
			return nil, fmt.Errorf("failed to open output file %q: %w", path, err)
		}
		files = append(files, f)
	}
	return files, nil
}

// closeFiles 关闭所有文件
//
// 参数:
//   - files: 文件句柄列表
func closeFiles(files []*os.File) {
	for _, f := range files {
		_ = f.Close()
	}
}

// setupSignalHandler 设置忽略 SIGINT
func setupSignalHandler() {
	signal.Ignore(syscall.SIGINT)
}
