package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"gitee.com/MM-Q/fck/internal/utils"
)

// TailConfig tail 命令配置
type TailConfig struct {
	Targets       []string      // 目标文件列表
	Lines         int           // -n, 行数（默认10）
	Bytes         int64         // -c, 字节数（与 -n 互斥）
	Follow        bool          // -f, 实时追踪
	Quiet         bool          // -q, 不显示文件名标题
	Verbose       bool          // -v, 总是显示文件名标题
	FromStdin     bool          // 是否从标准输入读取
	SleepInterval time.Duration // 轮询间隔（默认100ms）
}

// TailFile tail 文件状态追踪
type TailFile struct {
	Path   string
	File   *os.File
	Reader *bufio.Reader
	Size   int64
}

// TailCmdMain tail 命令主入口
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func TailCmdMain(config TailConfig) error {
	// 防御性检查：确保行数至少为1（当不使用字节模式时）
	if config.Bytes == 0 && config.Lines <= 0 {
		config.Lines = 10
	}

	// 默认轮询间隔
	if config.SleepInterval == 0 {
		config.SleepInterval = 100 * time.Millisecond
	}

	// 优先处理管道输入
	if utils.IsStdinPipe() {
		return tailStdin(config, os.Stdin)
	}

	// 检查文件参数
	if len(config.Targets) == 0 {
		return fmt.Errorf("no file specified")
	}

	// 实时追踪模式
	if config.Follow {
		return tailFollowFiles(config)
	}

	// 普通模式
	return tailFiles(config)
}

// tailFiles 普通模式处理多个文件
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 错误
func tailFiles(config TailConfig) error {
	showHeader := !config.Quiet && (config.Verbose || len(config.Targets) > 1)

	for i, path := range config.Targets {
		if i > 0 {
			fmt.Println()
		}

		if showHeader {
			fmt.Printf("==> %s <==\n", path)
		}

		if err := tailFile(path, config); err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", path, err)
		}
	}

	return nil
}

// tailFile 处理单个文件
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回值:
//   - error: 错误
func tailFile(path string, config TailConfig) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if config.Bytes > 0 {
		return tailByBytes(file, config.Bytes)
	}
	return tailByLines(file, config.Lines)
}

// tailByLines 按行读取（环形缓冲区）
//
// 参数:
//   - file: 文件
//   - n: 行数
//
// 返回值:
//   - error: 错误
func tailByLines(file *os.File, n int) error {
	ring := make([]string, n)
	index := 0
	count := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ring[index] = scanner.Text()
		index = (index + 1) % n
		count++
	}

	// 计算起始位置
	start := 0
	if count >= n {
		start = index
	}

	// 输出
	linesToPrint := n
	if count < n {
		linesToPrint = count
	}

	for i := 0; i < linesToPrint; i++ {
		fmt.Println(ring[(start+i)%n])
	}

	return scanner.Err()
}

// tailByBytes 按字节读取
//
// 参数:
//   - file: 文件
//   - n: 字节数
//
// 返回值:
//   - error: 错误
func tailByBytes(file *os.File, n int64) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	start := int64(0)
	if size > n {
		start = size - n
	}

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, file)
	return err
}

// tailStdin 从标准输入读取
//
// 参数:
//   - config: 命令配置
//   - stdin: 标准输入
//
// 返回值:
//   - error: 错误
func tailStdin(config TailConfig, stdin io.Reader) error {
	if config.Bytes > 0 {
		// 标准输入不支持字节模式，需要缓存
		data, err := io.ReadAll(stdin)
		if err != nil {
			return err
		}
		start := int64(0)
		if int64(len(data)) > config.Bytes {
			start = int64(len(data)) - config.Bytes
		}
		_, err = os.Stdout.Write(data[start:])
		return err
	}

	// 行模式使用环形缓冲区
	ring := make([]string, config.Lines)
	index := 0
	count := 0

	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		ring[index] = scanner.Text()
		index = (index + 1) % config.Lines
		count++
	}

	start := 0
	if count >= config.Lines {
		start = index
	}

	linesToPrint := config.Lines
	if count < config.Lines {
		linesToPrint = count
	}

	for i := 0; i < linesToPrint; i++ {
		fmt.Println(ring[(start+i)%config.Lines])
	}

	return scanner.Err()
}

// tailFollowFiles 实时追踪多个文件
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 错误
func tailFollowFiles(config TailConfig) error {
	// 打开所有文件
	files := make([]*TailFile, 0, len(config.Targets))
	for _, path := range config.Targets {
		tf, err := openTailFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", path, err)
			continue
		}
		files = append(files, tf)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files to follow")
	}

	// 显示文件名标题
	showHeader := !config.Quiet && len(files) > 1

	// 先显示初始内容
	for _, tf := range files {
		if showHeader {
			fmt.Printf("==> %s <==\n", tf.Path)
		}
		if err := tailFile(tf.Path, TailConfig{Lines: config.Lines, Bytes: config.Bytes}); err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", tf.Path, err)
		}
	}

	// 进入追踪模式
	ticker := time.NewTicker(config.SleepInterval)
	defer ticker.Stop()

	for range ticker.C {
		for _, tf := range files {
			if err := followFile(tf, showHeader); err != nil {
				// 文件可能被删除或重命名，尝试重新打开
				if newTf, err := reopenTailFile(tf); err == nil {
					*tf = *newTf
				}
			}
		}
	}

	return nil
}

// followFile 追踪单个文件的新内容
//
// 参数:
//   - tf: 文件追踪状态
//   - showHeader: 是否显示文件名标题
//
// 返回值:
//   - error: 错误
func followFile(tf *TailFile, showHeader bool) error {
	stat, err := tf.File.Stat()
	if err != nil {
		return err
	}

	newSize := stat.Size()
	if newSize < tf.Size {
		// 文件被截断或替换
		_, _ = tf.File.Seek(0, io.SeekStart)
		tf.Size = 0
		if showHeader {
			fmt.Printf("\n==> %s <==\n", tf.Path)
		}
	} else if newSize == tf.Size {
		// 没有新内容
		return nil
	}

	// 读取新内容
	for {
		line, err := tf.Reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if showHeader && tf.Size == 0 {
			fmt.Printf("\n==> %s <==\n", tf.Path)
		}

		fmt.Print(line)
		tf.Size += int64(len(line))
	}

	return nil
}

// openTailFile 打开文件用于追踪
//
// 参数:
//   - path: 文件路径
//
// 返回值:
//   - *TailFile: 文件追踪状态
//   - error: 错误
func openTailFile(path string) (*TailFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	// 定位到文件末尾
	_, _ = file.Seek(0, io.SeekEnd)

	return &TailFile{
		Path:   path,
		File:   file,
		Reader: bufio.NewReader(file),
		Size:   stat.Size(),
	}, nil
}

// reopenTailFile 重新打开文件
//
// 参数:
//   - tf: 原文件追踪状态
//
// 返回值:
//   - *TailFile: 新的文件追踪状态
//   - error: 错误
func reopenTailFile(tf *TailFile) (*TailFile, error) {
	_ = tf.File.Close()
	return openTailFile(tf.Path)
}
