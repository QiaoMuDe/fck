package xargs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"gitee.com/MM-Q/shellx/shx"
)

// XargsConfig 配置结构体
type XargsConfig struct {
	MaxArgs      int
	MaxProcs     int
	Replace      string
	NoRunIfEmpty bool
	MaxLines     int
	Verbose      bool
	Delimiter    string
	UseNull      bool
	Eof          string
	TargetCmd    string
	InitialArgs  []string
}

// XargsStats 执行统计（可选）
type XargsStats struct {
	Executed int
	Errors   int
}

// XargsCmdMain 主函数
func XargsCmdMain(config XargsConfig, stdin io.Reader) error {
	// 根据 -0 选项设置分隔符
	delimiter := config.Delimiter
	if config.UseNull {
		delimiter = "\x00"
	}

	// 读取标准输入
	inputs, err := readInputs(stdin, delimiter, config.Eof)
	if err != nil {
		return fmt.Errorf("读取输入失败: %w", err)
	}

	// 检查空输入
	if len(inputs) == 0 && config.NoRunIfEmpty {
		return nil
	}

	// 限制执行次数
	if config.MaxLines > 0 && len(inputs) > config.MaxLines {
		inputs = inputs[:config.MaxLines]
	}

	// 执行命令
	return executeXargs(inputs, config)
}

// readInputs 从 reader 读取输入，按分隔符分割
func readInputs(r io.Reader, delimiter string, eof string) ([]string, error) {
	var inputs []string
	scanner := bufio.NewScanner(r)

	// 自定义分隔函数
	if delimiter != "\n" {
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := strings.Index(string(data), delimiter); i >= 0 {
				return i + len(delimiter), data[0:i], nil
			}
			if atEOF {
				return len(data), data, nil
			}
			return 0, nil, nil
		})
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行（除非使用 -0 选项）
		if line == "" && delimiter != "\x00" {
			continue
		}

		// 检查结束标志
		if eof != "" && line == eof {
			break
		}

		inputs = append(inputs, line)
	}

	return inputs, scanner.Err()
}

// executeXargs 执行 xargs 核心逻辑
func executeXargs(inputs []string, config XargsConfig) error {
	// 批量分组
	batches := batchInputs(inputs, config.MaxArgs)

	// 串行执行
	if config.MaxProcs == 1 {
		for _, batch := range batches {
			if err := executeBatch(batch, config); err != nil {
				return err
			}
		}
		return nil
	}

	// 并行执行
	return executeParallel(batches, config)
}

// batchInputs 将输入分批
func batchInputs(inputs []string, batchSize int) [][]string {
	if batchSize <= 0 {
		batchSize = 1
	}

	var batches [][]string
	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		batches = append(batches, inputs[i:end])
	}

	return batches
}

// executeBatch 执行一批命令
func executeBatch(batch []string, config XargsConfig) error {
	// 有占位符：每行单独执行（逐行替换）
	if config.Replace != "" {
		for _, item := range batch {
			args := make([]string, len(config.InitialArgs))
			copy(args, config.InitialArgs)

			// 替换占位符
			for i, arg := range args {
				args[i] = strings.ReplaceAll(arg, config.Replace, item)
			}

			// 如果参数列表为空，使用替换值作为参数
			if len(args) == 0 {
				args = append(args, item)
			}

			// 打印命令（如果启用 verbose）
			if config.Verbose {
				fmt.Printf("+ %s %s\n", config.TargetCmd, strings.Join(args, " "))
			}

			// 使用 shx 执行命令
			cmd := shx.NewArgs(config.TargetCmd, args...).
				WithStdout(os.Stdout).
				WithStderr(os.Stderr).
				WithStdin(os.Stdin)

			if err := cmd.Exec(); err != nil {
				return err
			}
		}
		return nil
	}

	// 无占位符：所有参数追加，执行一次命令
	args := make([]string, len(config.InitialArgs))
	copy(args, config.InitialArgs)
	args = append(args, batch...)

	// 打印命令（如果启用 verbose）
	if config.Verbose {
		fmt.Printf("+ %s %s\n", config.TargetCmd, strings.Join(args, " "))
	}

	// 使用 shx 执行命令
	cmd := shx.NewArgs(config.TargetCmd, args...).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin)

	return cmd.Exec()
}

// executeParallel 并行执行
func executeParallel(batches [][]string, config XargsConfig) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxProcs)
	errChan := make(chan error, len(batches))

	for _, batch := range batches {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(b []string) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			if err := executeBatch(b, config); err != nil {
				errChan <- err
			}
		}(batch)
	}

	wg.Wait()
	close(errChan)

	// 检查错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
