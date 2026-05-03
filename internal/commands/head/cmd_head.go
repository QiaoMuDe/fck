package head

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"gitee.com/MM-Q/go-kit/term"
)

// HeadConfig head 命令配置
type HeadConfig struct {
	Targets   []string // 目标文件列表
	Lines     int      // -n, 行数（默认10）
	Bytes     int64    // -c, 字节数（与 -n 互斥）
	Quiet     bool     // -q, 不显示文件名标题
	Verbose   bool     // -v, 总是显示文件名标题
	FromStdin bool     // 是否从标准输入读取
}

// HeadCmdMain head 命令主入口
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func HeadCmdMain(config HeadConfig) error {
	// 防御性检查：确保行数至少为1（当不使用字节模式时）
	if config.Bytes == 0 && config.Lines <= 0 {
		config.Lines = 10
	}

	// 优先处理管道输入
	if term.IsStdinPipe() {
		return headStdin(config, os.Stdin)
	}

	// 检查文件参数
	if len(config.Targets) == 0 {
		return fmt.Errorf("no file specified")
	}

	// 多个文件处理
	showHeader := !config.Quiet && (config.Verbose || len(config.Targets) > 1)

	for i, path := range config.Targets {
		if i > 0 {
			fmt.Println()
		}

		if showHeader {
			fmt.Printf("==> %s <==\n", path)
		}

		if err := headFile(path, config); err != nil {
			fmt.Fprintf(os.Stderr, "head: %s: %v\n", path, err)
		}
	}

	return nil
}

// headFile 处理单个文件
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回值:
//   - error: 错误
func headFile(path string, config HeadConfig) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if config.Bytes > 0 {
		return headByBytes(file, config.Bytes)
	}
	return headByLines(file, config.Lines)
}

// headByLines 按行读取
//
// 参数:
//   - file: 文件
//   - n: 行数
//
// 返回值:
//   - error: 错误
func headByLines(file *os.File, n int) error {
	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		count++
		if count >= n {
			break
		}
	}

	return scanner.Err()
}

// headByBytes 按字节读取
//
// 参数:
//   - file: 文件
//   - n: 字节数
//
// 返回值:
//   - error: 错误
func headByBytes(file *os.File, n int64) error {
	reader := io.LimitReader(file, n)
	_, err := io.Copy(os.Stdout, reader)
	return err
}

// headStdin 从标准输入读取
//
// 参数:
//   - config: 命令配置
//   - stdin: 标准输入
//
// 返回值:
//   - error: 错误
func headStdin(config HeadConfig, stdin io.Reader) error {
	if config.Bytes > 0 {
		reader := io.LimitReader(stdin, config.Bytes)
		_, err := io.Copy(os.Stdout, reader)
		return err
	}

	scanner := bufio.NewScanner(stdin)
	count := 0

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		count++
		if count >= config.Lines {
			break
		}
	}

	return scanner.Err()
}
