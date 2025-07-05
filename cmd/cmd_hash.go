package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

// 字节单位定义
const (
	Byte = 1 << (10 * iota) // 1 字节
	KB                      // 千字节 (1024 B)
	MB                      // 兆字节 (1024 KB)
)

// hashCmdMain 是 hash 子命令的主函数
func hashCmdMain(cl *colorlib.ColorLib) error {
	// 获取指定的路径
	targetPaths := hashCmd.Args()

	// 如果没有指定路径，则打印错误信息并退出
	if len(targetPaths) == 0 {
		return fmt.Errorf("请指定要计算哈希值的路径")
	}

	// 检查指定的哈希算法是否有效
	hashType, ok := globals.SupportedAlgorithms[hashCmdType.Get()]
	if !ok {
		return fmt.Errorf("在校验哈希值时，哈希算法 %s 无效", hashCmdType.Get())
	}

	// 设置并发任务数
	if hashCmdJob.Get() == -1 {
		// 使用CPU核数*2
		if setErr := hashCmdJob.Set(fmt.Sprint(runtime.NumCPU() * 2)); setErr != nil {
			return fmt.Errorf("设置并发任务数失败: %v", setErr)
		}
	}
	if hashCmdJob.Get() <= 0 {
		// 最小并发数为1
		if setErr := hashCmdJob.Set(fmt.Sprint(1)); setErr != nil {
			return fmt.Errorf("设置并发任务数失败: %v", setErr)
		}
	}
	if hashCmdJob.Get() > 20 {
		// 最大并发数为20
		if setErr := hashCmdJob.Set(fmt.Sprint(20)); setErr != nil {
			return fmt.Errorf("设置并发任务数失败: %v", setErr)
		}
	}

	// 创建一个可以被外部取消的上下文（支持指定取消原因）
	ctx, cancel := context.WithCancelCause(context.Background())

	// 启动一个 goroutine，在用户按下 Ctrl+C 时取消上下文
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		// 使用自定义错误作为取消原因
		cancel(fmt.Errorf("用户中断操作"))
		cl.PrintWarn("已取消所有任务")
		os.Exit(1) // 退出程序
	}()

	// 检查是否需要写入文件
	if hashCmdWrite.Get() {
		cl.PrintOk("正在将哈希值写入文件，请稍候...")
	}

	// 遍历路径
	for _, targetPath := range targetPaths {
		// 清理路径
		targetPath = filepath.Clean(targetPath)

		// 定义一个切片来存储文件列表
		var files []string

		// 获取文件列表
		files, err := collectFiles(targetPath, hashCmdRecursion.Get(), cl)
		if err != nil {
			// 如果收集文件失败，则打印错误信息并退出
			cl.PrintErrf("在校验哈希值时，收集文件失败: %v\n", err)
			continue
		}

		// 检查文件列表是否为空
		if len(files) == 0 {
			cl.PrintWarnf("在校验哈希值时，路径 %s 没有找到任何文件\n", targetPath)
			continue
		}

		// 执行哈希任务
		errors := hashRunTasks(ctx, files, hashType)

		// 检查是否有错误发生
		if len(errors) > 0 {
			// 使用map去重错误信息，避免重复打印相同的错误信息
			// 定义一个字符串到布尔值的映射，用于存储已经出现过的错误信息
			errorMap := make(map[string]bool)
			// 遍历错误列表
			for _, err := range errors {
				// 将错误对象转换为字符串
				errStr := err.Error()
				// 检查该错误信息是否已经在映射中
				if !errorMap[errStr] {
					// 如果不在映射中，将该错误信息添加到映射中，表示已经出现过
					errorMap[errStr] = true
					// 打印错误信息
					cl.PrintErr(errStr)
				}
			}
		} else {
			// 打印成功信息
			if hashCmdWrite.Get() {
				cl.PrintOkf("已将哈希值写入文件 %s, 共处理 %d 个文件\n", globals.OutputFileName, len(files))
			}
		}
	}

	return nil
}

// hashRunTasks 执行哈希值校验任务，支持并发控制和错误处理
func hashRunTasks(ctx context.Context, files []string, hashType func() hash.Hash) []error {
	// 创建带取消原因的上下文
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil) // 正常完成时不设置取消原因

	// 根据 hashCmdJob.Get() 的值创建一个工作池
	jobPool := make(chan struct{}, hashCmdJob.Get())

	// 创建错误收集通道
	errors := make(chan error, len(files))

	// 创建一个等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 添加一个标志，确保只取消一次
	var once sync.Once

	// 检查是否需要写入文件
	var fileWrite *os.File
	if hashCmdWrite.Get() {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			errors <- fmt.Errorf("打开文件 %s 失败: %v", globals.OutputFileName, err)
			return []error{fmt.Errorf("打开文件 %s 失败: %v", globals.OutputFileName, err)}
		}
		defer fileWrite.Close()

		// 写入文件头
		if err := writeFileHeader(fileWrite, hashCmdType.Get(), globals.TimestampFormat); err != nil {
			errors <- fmt.Errorf("写入文件头 %s 失败: %v", globals.OutputFileName, err)
			return []error{fmt.Errorf("写入文件头 %s 失败: %v", globals.OutputFileName, err)}
		}
	}

fileLoop:
	// 遍历文件列表并启动任务
	for _, file := range files {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			errors <- ctx.Err()
			break fileLoop // 跳出循环，不再处理其他文件
		default:
			// 向工作池发送一个信号，表示有一个任务开始
			jobPool <- struct{}{}
			// 启动一个 goroutine 执行任务
			wg.Add(1)
			go func(file string) {
				// 使用 recover 防止 panic 导致程序崩溃
				defer func() {
					wg.Done() // 任务完成，减少等待组计数
					<-jobPool // 从工作池接收一个信号，表示任务完成

					// 恢复可能的 panic
					if r := recover(); r != nil {
						err := fmt.Errorf("处理文件 %s 时发生 panic: %v", file, r)
						errors <- err
						once.Do(func() { cancel(err) }) // 只取消一次
					}
				}()

				// 检查文件路径是否存在
				if _, err := os.Lstat(file); err != nil {
					errors <- fmt.Errorf("文件 %s 不存在或无法访问: %v", file, err)
					once.Do(func() { cancel(fmt.Errorf("文件 %s 不存在或无法访问: %v", file, err)) }) // 只取消一次
					return
				}

				// 执行任务并收集错误
				if err := hashTask(ctx, file, hashType, fileWrite); err != nil {
					// 发生错误时取消所有其他任务
					if errors != nil {
						errors <- err
						once.Do(func() { cancel(err) }) // 只取消一次
					}
				}
			}(file)
		}
	}

	// 等待所有任务完成，然后关闭错误通道
	go func() {
		wg.Wait()
		close(errors)
	}()

	// 收集所有错误
	var resultErrors []error
	for err := range errors {
		resultErrors = append(resultErrors, err)
	}

	// 如果上下文被取消，将取消原因添加到错误列表
	if ctx.Err() != nil {
		if cause := context.Cause(ctx); cause != nil {
			resultErrors = append(resultErrors, fmt.Errorf("任务被取消: %w", cause))
		} else {
			resultErrors = append(resultErrors, ctx.Err())
		}
	}

	return resultErrors
}

// hashTask 处理单个文件的哈希计算，支持上下文取消
func hashTask(ctx context.Context, filePath string, hashType func() hash.Hash, fileWrite *os.File) error {
	// 检查上下文是否已取消
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 检查文件如果是软链接，则跳过
	if fileInfo, err := os.Lstat(filePath); err == nil {
		if fileInfo.Mode()&fs.ModeSymlink != 0 {
			return nil
		}
	}

	// 计算文件哈希值
	hashValue, checkErr := checksum(filePath, hashType)
	if checkErr != nil {
		return fmt.Errorf("计算文件 %s 哈希值失败: %v", filePath, checkErr)
	}

	// 检查上下文是否已取消（计算哈希值可能需要一些时间）
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 根据 hashCmdWrite.Get() 标志决定是否写入文件
	if hashCmdWrite.Get() {
		// 写入到 globals.OutputFileName 文件
		if fileWrite != nil {
			_, err := fileWrite.WriteString(fmt.Sprintf("%s\t%q\n", hashValue, filePath))
			if err != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, err)
			}
		}
	} else {
		fmt.Printf("%s\t%s\n", hashValue, filePath)
	}

	return nil
}

// 计算文件哈希值的函数
func checksum(filePath string, hashFunc func() hash.Hash) (string, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("文件不存在或无法访问: %v", err)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建哈希对象
	hash := hashFunc()

	// 根据文件大小动态分配缓冲区
	fileSize := fileInfo.Size()
	bufferSize := calculateBufferSize(fileSize)
	buffer := make([]byte, bufferSize)

	// 使用 io.CopyBuffer 进行高效复制并计算哈希
	if _, err := io.CopyBuffer(hash, file, buffer); err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 返回哈希值的十六进制表示
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// calculateBufferSize 根据文件大小计算最佳缓冲区大小
// 参数 fileSize 是文件大小(字节)
// 返回值为计算出的缓冲区大小(字节)
func calculateBufferSize(fileSize int64) int {
	switch {
	case fileSize < 32*KB: // 小于 32KB 的文件使用 32KB 缓冲区
		return int(32 * KB)
	case fileSize < 128*KB: // 32KB-128KB 使用 64KB 缓冲区
		return int(64 * KB)
	case fileSize < 512*KB: // 128KB-512KB 使用 128KB 缓冲区
		return int(128 * KB)
	case fileSize < 1*MB: // 512KB-1MB 使用 256KB 缓冲区
		return int(256 * KB)
	case fileSize < 4*MB: // 1MB-4MB 使用 512KB 缓冲区
		return int(512 * KB)
	case fileSize < 16*MB: // 4MB-16MB 使用 1MB 缓冲区
		return int(1 * MB)
	case fileSize < 64*MB: // 16MB-64MB 使用 2MB 缓冲区
		return int(2 * MB)
	default: // 大于 64MB 的文件使用 4MB 缓冲区
		return int(4 * MB)
	}
}

// walkDir 函数用于根据递归标志遍历指定目录并收集文件列表。
// 参数 dirPath 是要遍历的目录路径。
// 参数 recursive 表示是否进行递归遍历，如果为 true 则遍历目录及其子目录，否则只遍历当前目录。
// 参数 cl 是一个颜色库实例，用于打印警告信息。
// 返回值为收集到的文件路径切片和可能出现的错误。
func walkDir(dirPath string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	// 初始化一个字符串切片，用于存储收集到的文件路径
	var files []string

	// 判断是否开启递归模式
	if !recursive {
		// 非递归模式：只获取目录下的直接文件
		// 读取指定目录下的所有条目
		dir, err := os.ReadDir(dirPath)
		// 检查读取目录是否失败
		if err != nil {
			// 如果失败，返回 nil 和包装后的错误信息
			return nil, fmt.Errorf("读取目录失败: %w", err)
		}

		// 遍历目录中的所有条目
		for _, entry := range dir {
			// 如果是隐藏项并且不允许隐藏项，则跳过该条目
			if !hashCmdHidden.Get() && isHidden(entry.Name()) {
				continue
			}

			// 判断当前条目是否为目录
			if entry.IsDir() {
				// 如果是目录，打印警告信息并跳过该目录
				cl.PrintWarnf("跳过目录：%s, 请使用 -r 选项以递归方式处理\n", entry.Name())
				// 继续遍历下一个条目
				continue
			}

			// 如果是文件，将其完整路径添加到文件列表中
			files = append(files, filepath.Join(dirPath, entry.Name()))
		}

		// 返回收集到的文件列表和 nil 错误
		return files, nil
	}

	// 递归模式：遍历目录及其子目录中的所有文件
	// 使用 filepath.WalkDir 函数递归遍历目录
	walkDirErr := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		// 检查遍历过程中是否出现错误
		if err != nil {
			// 如果出现错误，直接返回该错误
			return err
		}

		// 如果是隐藏项并且不允许隐藏项，则跳过该条目
		if !hashCmdHidden.Get() && isHidden(path) {
			// 如果是隐藏目录，跳过整个目录
			if d.IsDir() {
				return filepath.SkipDir
			}

			// 如果是隐藏文件，跳过该文件
			return nil
		}

		// 如果不是目录，则将其完整路径添加到文件列表中
		if !d.IsDir() {
			files = append(files, path)
		}

		// 继续遍历下一个条目
		return nil
	})

	// 检查 filepath.WalkDir 函数是否返回错误
	if walkDirErr != nil {
		// 如果是权限错误
		if os.IsPermission(walkDirErr) {
			// 返回错误信息，表示权限不足
			return nil, fmt.Errorf("权限不足: %s", dirPath)
		}
		// 如果是文件不存在错误
		if os.IsNotExist(walkDirErr) {
			// 返回错误信息，表示文件不存在
			return nil, fmt.Errorf("文件不存在: %s", dirPath)
		}

		// 如果出现错误，返回 nil 和包装后的错误信息
		return nil, fmt.Errorf("遍历目录失败: %w", walkDirErr)
	}

	// 返回收集到的文件列表和 nil 错误
	return files, nil
}

// collectFiles 函数用于收集指定路径下的所有文件。该路径可以是普通路径（文件或目录），也可以是包含通配符的路径。
// 参数 targetPath 是要处理的目标路径，可能包含通配符。
// 参数 recursive 表示是否递归遍历目录，如果为 true 则会递归遍历目录及其子目录，否则只遍历当前目录。
// 参数 cl 是一个颜色库实例，用于打印警告信息。
// 返回值为收集到的文件路径切片和可能出现的错误。
func collectFiles(targetPath string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	// 初始化一个字符串切片，用于存储收集到的文件路径
	var files []string

	// 检查路径是否包含通配符，通配符包括 *、?、[] 和 {}
	if strings.ContainsAny(targetPath, "*?[]{}") {
		// 处理包含通配符的路径，使用 filepath.Glob 函数获取匹配的文件或目录列表
		matchedFiles, err := filepath.Glob(targetPath)
		if err != nil {
			// 如果使用 filepath.Glob 函数时出现错误，返回错误信息，表示路径无效
			return nil, fmt.Errorf("路径无效: %w", err)
		}

		// 检查是否有匹配的文件或目录
		if len(matchedFiles) == 0 {
			// 如果没有找到匹配的文件或目录，返回错误信息
			return nil, fmt.Errorf("没有找到匹配的文件")
		}

		// 遍历所有匹配的文件或目录
		for _, file := range matchedFiles {
			// 获取文件或目录的信息
			info, statErr := os.Stat(file)
			if statErr != nil {
				// 检查是否是权限错误
				if os.IsPermission(statErr) {
					// 如果是权限错误，返回错误信息
					return nil, fmt.Errorf("权限不足: %s", file)
				}

				// 检查是否为文件不存在错误
				if os.IsNotExist(statErr) {
					// 如果是文件不存在错误，返回错误信息
					return nil, fmt.Errorf("文件不存在: %s", file)
				}

				// 如果无法获取文件或目录的信息，返回错误信息
				return nil, fmt.Errorf("无法获取文件信息: %w", statErr)
			}

			// 如果是隐藏项且不允许隐藏项, 则跳过该目录
			if !hashCmdHidden.Get() && isHidden(file) {
				continue
			}

			// 判断是否为目录
			if info.IsDir() {
				// 如果是目录且不允许递归目录
				if !recursive {
					cl.PrintWarnf("跳过目录：%s, 请使用 -r 选项以递归方式处理\n", file)
					continue
				}

				// 如果是目录，调用 walkDir 函数遍历该目录并收集文件
				filesInDir, err := walkDir(file, recursive, cl)
				if err != nil {
					// 如果遍历目录时出现错误，返回该错误
					return nil, err
				}

				// 将目录中的文件添加到文件列表中
				files = append(files, filesInDir...)

				// 继续遍历下一个条目
				continue
			}

			// 普通文件，添加到文件列表中
			files = append(files, file)
		}

		// 返回收集到的文件列表和 nil 错误
		return files, nil
	}

	// 获取路径的信息
	info, err := os.Stat(targetPath)
	if err != nil {
		// 如果无法获取路径的信息，返回错误信息
		return nil, fmt.Errorf("无法获取路径信息: %w", err)
	}

	// 如果是隐藏项且不允许隐藏项, 则跳过该路径
	if !hashCmdHidden.Get() && isHidden(targetPath) {
		return nil, fmt.Errorf("跳过隐藏项: %s", targetPath)
	}

	// 判断是否为目录
	if info.IsDir() {
		// 如果是目录且不允许递归目录
		if !recursive {
			return nil, fmt.Errorf("跳过目录：%s, 请使用 -r 选项以递归方式处理", targetPath)
		}

		// 如果是目录，调用 walkDir 函数遍历该目录并收集文件
		filesInDir, err := walkDir(targetPath, recursive, cl)
		if err != nil {
			// 如果遍历目录时出现错误，返回该错误
			return nil, err
		}

		// 将目录中的文件添加到文件列表中
		files = append(files, filesInDir...)

		// 返回收集到的文件列表和 nil 错误
		return files, nil
	}

	// 如果是普通文件，将其路径添加到文件列表中
	files = []string{targetPath}

	// 返回收集到的文件列表和 nil 错误
	return files, nil
}
