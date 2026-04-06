package xargs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"gitee.com/MM-Q/shellx/shx"
)

// XargsConfig xargs命令配置
type XargsConfig struct {
	// 输入配置
	Delimiter string // 分隔符
	NullMode  bool   // \0 分隔模式
	ArgFile   string // 参数文件

	// 批量控制
	MaxArgs  int // 每批最大参数
	MaxLines int // 每批最大行数
	MaxChars int // 命令最大长度

	// 执行配置
	ReplaceStr   string // 占位符（如 {}）
	ReplaceDelim string // 自定义占位符字符串
	MaxProcs     int    // 并行数
	NoRunIfEmpty bool   // 空输入不执行
	Interactive  bool   // 确认模式
	Verbose      bool   // 显示命令
	ExitOnError  bool   // 出错停止

	// 目标命令
	Command     string   // 要执行的命令
	CommandArgs []string // 命令固定参数
}

// XargsStats 执行统计
type XargsStats struct {
	Executed int // 执行次数
	Failed   int // 失败次数
}

// XargsCmdMain 执行xargs命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误
func XargsCmdMain(config XargsConfig) error {
	// 读取参数
	args, err := readArgs(config)
	if err != nil {
		return err
	}

	// 空输入处理
	if len(args) == 0 {
		if config.NoRunIfEmpty {
			return nil
		}
		return fmt.Errorf("没有输入参数")
	}

	// 按规则分批
	batches := splitBatches(args, config)

	// 执行批次
	stats := &XargsStats{}
	if config.MaxProcs > 1 {
		err = runParallel(batches, config, stats)
	} else {
		err = runSequential(batches, config, stats)
	}

	if err != nil {
		return err
	}

	// 如果有失败，返回错误
	if stats.Failed > 0 {
		return fmt.Errorf("执行完成，%d 次成功，%d 次失败", stats.Executed-stats.Failed, stats.Failed)
	}

	return nil
}

// readArgs 读取参数
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - []string: 参数列表
//   - error: 读取错误
func readArgs(config XargsConfig) ([]string, error) {
	var reader *bufio.Reader

	// 从文件或标准输入读取
	if config.ArgFile != "" {
		file, err := os.Open(config.ArgFile)
		if err != nil {
			return nil, fmt.Errorf("打开文件失败: %w", err)
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	var args []string
	delimiter := config.Delimiter

	if config.NullMode {
		delimiter = "\x00"
	}

	if delimiter == "" {
		// 默认分隔符：空格、换行、制表
		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("读取输入失败: %w", err)
		}
	} else {
		// 自定义分隔符
		data, err := reader.ReadString('\x00')
		if err != nil && err.Error() != "EOF" {
			return nil, fmt.Errorf("读取输入失败: %w", err)
		}
		// 移除末尾的 \x00
		data = strings.TrimSuffix(data, "\x00")
		if data != "" {
			parts := strings.Split(data, delimiter)
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part != "" {
					args = append(args, part)
				}
			}
		}
	}

	return args, nil
}

// splitBatches 将参数分批
//
// 参数:
//   - args: 参数列表
//   - config: 命令配置
//
// 返回:
//   - [][]string: 分批后的参数
func splitBatches(args []string, config XargsConfig) [][]string {
	var batches [][]string
	var currentBatch []string

	// 确定每批大小
	batchSize := config.MaxArgs
	if batchSize <= 0 {
		batchSize = len(args) // 默认全部一批
	}

	for i, arg := range args {
		currentBatch = append(currentBatch, arg)

		// 检查是否达到批次限制
		shouldSplit := false

		// 按参数个数
		if config.MaxArgs > 0 && len(currentBatch) >= config.MaxArgs {
			shouldSplit = true
		}

		// 按行数
		if config.MaxLines > 0 && len(currentBatch) >= config.MaxLines {
			shouldSplit = true
		}

		// 按命令长度
		if config.MaxChars > 0 {
			cmdLen := calculateCmdLen(config, currentBatch)
			if cmdLen > config.MaxChars {
				// 如果当前批次已超限，且不止一个参数，移除最后一个
				if len(currentBatch) > 1 {
					currentBatch = currentBatch[:len(currentBatch)-1]
					batches = append(batches, currentBatch)
					currentBatch = []string{arg}
				} else {
					shouldSplit = true
				}
			}
		}

		// 到达批次末尾或最后一个参数
		if shouldSplit || i == len(args)-1 {
			if len(currentBatch) > 0 {
				batches = append(batches, currentBatch)
				currentBatch = nil
			}
		}
	}

	// 处理剩余的参数
	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batches
}

// calculateCmdLen 计算命令长度
//
// 参数:
//   - config: 命令配置
//   - batch: 当前批次参数
//
// 返回:
//   - int: 命令长度
func calculateCmdLen(config XargsConfig, batch []string) int {
	length := len(config.Command)
	for _, arg := range config.CommandArgs {
		length += len(arg) + 1 // +1 for space
	}
	for _, arg := range batch {
		length += len(arg) + 1 // +1 for space
	}
	return length
}

// runSequential 顺序执行批次
//
// 参数:
//   - batches: 参数批次
//   - config: 命令配置
//   - stats: 执行统计
//
// 返回:
//   - error: 执行错误
func runSequential(batches [][]string, config XargsConfig, stats *XargsStats) error {
	for _, batch := range batches {
		if err := executeBatch(batch, config, stats); err != nil {
			if config.ExitOnError {
				return err
			}
			// 继续执行下一批
		}
	}
	return nil
}

// runParallel 并行执行批次
//
// 参数:
//   - batches: 参数批次
//   - config: 命令配置
//   - stats: 执行统计
//
// 返回:
//   - error: 执行错误
func runParallel(batches [][]string, config XargsConfig, stats *XargsStats) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxProcs)
	var mu sync.Mutex
	var firstErr error

	for _, batch := range batches {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(b []string) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			if err := executeBatch(b, config, stats); err != nil {
				mu.Lock()
				if config.ExitOnError && firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}(batch)
	}

	wg.Wait()
	return firstErr
}

// executeBatch 执行单个批次
//
// 参数:
//   - batch: 批次参数
//   - config: 命令配置
//   - stats: 执行统计
//
// 返回:
//   - error: 执行错误
func executeBatch(batch []string, config XargsConfig, stats *XargsStats) error {
	stats.Executed++

	// 构建命令
	cmdStr := buildCommand(batch, config)

	// 显示命令
	if config.Verbose {
		fmt.Fprintln(os.Stderr, cmdStr)
	}

	// 确认模式
	if config.Interactive {
		fmt.Fprintf(os.Stderr, "执行? (y/n): ")
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			return fmt.Errorf("读取确认失败: %w", err)
		}
		if response != "y" && response != "Y" {
			return nil
		}
	}

	// 执行命令
	if err := shx.RunToTerminal(cmdStr); err != nil {
		stats.Failed++
		return fmt.Errorf("执行失败: %w", err)
	}

	return nil
}

// buildCommand 构建命令字符串
//
// 参数:
//   - batch: 批次参数
//   - config: 命令配置
//
// 返回:
//   - string: 命令字符串
func buildCommand(batch []string, config XargsConfig) string {
	// 确定占位符
	placeholder := config.ReplaceStr
	if placeholder == "" {
		placeholder = config.ReplaceDelim
	}
	if placeholder == "" {
		placeholder = "{}"
	}

	// 替换模式
	if config.ReplaceStr != "" || config.ReplaceDelim != "" {
		// 每个参数单独替换，执行多次
		var cmds []string
		for _, arg := range batch {
			cmd := config.Command
			// 替换命令中的占位符
			cmd = strings.ReplaceAll(cmd, placeholder, arg)
			// 替换固定参数中的占位符
			for _, fixedArg := range config.CommandArgs {
				replacedArg := strings.ReplaceAll(fixedArg, placeholder, arg)
				cmd += " " + replacedArg
			}
			cmds = append(cmds, cmd)
		}
		return strings.Join(cmds, " && ")
	}

	// 默认模式：所有参数追加到命令后
	cmd := config.Command
	for _, arg := range config.CommandArgs {
		cmd += " " + arg
	}
	for _, arg := range batch {
		cmd += " " + arg
	}
	return cmd
}
