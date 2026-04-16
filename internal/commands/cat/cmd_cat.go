package cat

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"gitee.com/MM-Q/fck/internal/utils"
)

// CatConfig cat 命令配置
type CatConfig struct {
	// CLI 参数
	Targets      []string // 目标文件列表
	ShowLineNum  bool     // -n 显示所有行号
	ShowNonBlank bool     // -b 显示非空行行号
	ShowEnd      bool     // -E 显示行尾$
	ShowTabs     bool     // -T 显示制表符为^I
	ShowAll      bool     // -A 等价于 -ET
	ShowNewline  bool     // -N 显示换行符类型
	HeadLines    int      // --head 显示前N行 (0表示全部)
	TailLines    int      // --tail 显示后N行 (0表示全部)
	Quiet        bool     // -q 静默模式 (不显示错误信息)
	Text         bool     // -a, --text 强制将二进制文件视为文本处理
	IgnoreBinary bool     // -I, --ignore-binary 完全忽略二进制文件
	UseLess      bool     // -l, --less 使用分页器查看文件内容

	// 运行时
	LineCounter int // 行号计数器
}

// lineInfo 存储行内容和换行符信息
type lineInfo struct {
	content string
	newline string
}

// CatCmdMain 执行 cat 命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误 (如果有)
func CatCmdMain(config CatConfig) error {
	// 1. 验证参数
	if len(config.Targets) == 0 {
		return fmt.Errorf("no files specified")
	}

	// 2. 如果使用分页器模式，直接调用分页器查看
	if config.UseLess {
		return runOVMode(config)
	}

	// 3. 处理标志冲突 (-b 优先级高于 -n)
	if config.ShowNonBlank {
		config.ShowLineNum = false
	}

	// 4. 处理 --show-all
	if config.ShowAll {
		config.ShowEnd = true
		config.ShowTabs = true
	}

	// 5. 验证 head/tail 互斥
	if config.HeadLines > 0 && config.TailLines > 0 {
		return fmt.Errorf("cannot use --head (-u) and --tail (-d) together")
	}

	// 6. 处理每个文件
	for _, target := range config.Targets {
		if err := processFile(target, &config); err != nil {
			if !config.Quiet {
				return err
			}
		}
	}

	return nil
}

// readLine 读取一行，返回内容、换行符标记和错误
//
// 参数:
//   - reader: bufio.Reader
//
// 返回:
//   - content: 行内容（不含换行符）
//   - newline: 换行符标记（"\n"、"\r\n" 或空字符串）
//   - err: 错误信息
func readLine(reader *bufio.Reader) (content, newline string, err error) {
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", "", err
	}

	// 检测换行符类型
	switch {
	case strings.HasSuffix(line, "\r\n"):
		// Windows 格式 (CR+LF)
		return line[:len(line)-2], "\\r\\n", err

	case strings.HasSuffix(line, "\n"):
		// Unix 格式 (LF)
		return line[:len(line)-1], "\\n", err

	default:
		// 无换行符（最后一行）
		return line, "", err
	}
}

// processFile 处理单个文件
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果有)
func processFile(path string, config *CatConfig) error {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	// 获取文件信息 (用于判断是否是目录)
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info %s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}

	// 二进制文件检测 (除非强制文本模式)
	if !config.Text {
		isBinary, err := utils.IsBinaryFile(file)
		if err != nil {
			return fmt.Errorf("cannot detect file type for %s: %w", path, err)
		}

		// 处理二进制文件
		if isBinary {
			// -I 模式：静默跳过
			if config.IgnoreBinary {
				return nil
			}

			// 默认行为：输出提示并跳过
			if !config.Quiet {
				fmt.Printf("bin file %s matches\n", path)
			}
			return nil
		}
	}

	// 根据 head/tail 选项处理
	if config.HeadLines > 0 {
		return processHead(file, config)
	}

	if config.TailLines > 0 {
		return processTail(file, config)
	}

	// 普通处理：使用 bufio.Reader 逐行读取
	reader := bufio.NewReader(file)
	for {
		line, newline, err := readLine(reader)
		if err != nil && err != io.EOF {
			return err
		}

		processLine(line, newline, config)

		if err == io.EOF {
			break
		}
	}

	return nil
}

// processHead 处理文件前N行
//
// 参数:
//   - file: 已打开的文件
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果有)
func processHead(file *os.File, config *CatConfig) error {
	reader := bufio.NewReader(file)
	lineCount := 0

	// 循环读取每一行
	for lineCount < config.HeadLines {
		line, newline, err := readLine(reader)
		if err != nil && err != io.EOF {
			return err
		}

		processLine(line, newline, config)
		lineCount++

		if err == io.EOF {
			break
		}
	}

	return nil
}

// processTail 处理文件后N行
//
// 参数:
//   - file: 已打开的文件
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果有)
func processTail(file *os.File, config *CatConfig) error {
	// 使用环形缓冲区存储最后N行
	ring := make([]lineInfo, config.TailLines)
	index := 0
	count := 0

	reader := bufio.NewReader(file)
	for {
		line, newline, err := readLine(reader)
		if err != nil && err != io.EOF {
			return err
		}

		ring[index] = lineInfo{content: line, newline: newline}
		index = (index + 1) % config.TailLines
		if count < config.TailLines {
			count++
		}

		if err == io.EOF {
			break
		}
	}

	// 按顺序输出
	start := (index - count + config.TailLines) % config.TailLines
	for i := 0; i < count; i++ {
		info := ring[(start+i)%config.TailLines]
		processLine(info.content, info.newline, config)
	}

	return nil
}

// processLine 处理单行
//
// 参数:
//   - line: 行内容
//   - newline: 换行符标记
//   - config: 命令配置
func processLine(line, newline string, config *CatConfig) {
	isBlank := len(strings.TrimSpace(line)) == 0

	// 处理行号
	if config.ShowLineNum || (config.ShowNonBlank && !isBlank) {
		config.LineCounter++
		fmt.Printf("%6d\t", config.LineCounter)
	}

	// 处理特殊字符显示
	if config.ShowTabs {
		line = strings.ReplaceAll(line, "\t", "^I")
	}

	// 输出内容
	fmt.Print(line)

	// 处理换行符显示
	if config.ShowNewline {
		fmt.Print(newline)
	}

	// 处理行尾标记
	if config.ShowEnd {
		fmt.Print("$")
	}
	fmt.Println()
}
