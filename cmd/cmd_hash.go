package cmd

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
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
	GB                      // 吉字节 (1024 MB)
	TB                      // 太字节 (1024 GB)
)

// hashCmdMain 是 hash 子命令的主函数
func hashCmdMain(cmd *flag.FlagSet, cl *colorlib.ColorLib) error {
	// 获取指定的路径
	targetPath := cmd.Arg(0)

	// 如果没有指定路径，则打印错误信息并退出
	if targetPath == "" {
		return fmt.Errorf("在校验哈希值时，必须指定一个路径")
	}

	// 检查并发数是否有效
	if *hashCmdJob <= 0 {
		return fmt.Errorf("在校验哈希值时，并发数必须大于 0")
	}

	// 清理路径
	targetPath = filepath.Clean(targetPath)

	// 定义一个切片来存储文件列表
	var files []string

	// 获取文件列表
	files, err := collectFiles(targetPath, *hashCmdRecursion, cl)
	if err != nil {
		return fmt.Errorf("在校验哈希值时，收集文件失败: %v", err)
	}

	// 检查文件列表是否为空
	if len(files) == 0 {
		return fmt.Errorf("在校验哈希值时，路径 %s 没有找到任何文件", targetPath)
	}

	// 检查指定的哈希算法是否有效
	hashType, ok := globals.SupportedAlgorithms[*hashCmdType]
	if !ok {
		return fmt.Errorf("在校验哈希值时，哈希算法 %s 无效", *hashCmdType)
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
	}()

	// 执行哈希任务
	errors := hashRunTasks(ctx, files, hashType)

	// 检查是否有错误发生
	if len(errors) > 0 {
		var errText string
		// 打印错误信息
		for _, err := range errors {
			if errText != err.Error() {
				cl.PrintErr(err.Error())
			}
		}
	} else {
		// 打印成功信息
		if *hashCmdWrite {
			cl.PrintOk(fmt.Sprintf("校验哈希值完成，共处理 %d 个文件, 并将哈希值写入文件 %s", len(files), globals.OutputFileName))
		}
	}

	return nil
}

// hashRunTasks 执行哈希值校验任务，支持并发控制和错误处理
func hashRunTasks(ctx context.Context, files []string, hashType func() hash.Hash) []error {
	// 创建带取消原因的上下文
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil) // 正常完成时不设置取消原因

	// 根据 *hashCmdJob 的值创建一个工作池
	jobPool := make(chan struct{}, *hashCmdJob)

	// 创建错误收集通道
	errors := make(chan error, len(files))

	// 创建一个等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 添加一个标志，确保只取消一次
	var once sync.Once

	// 检查是否需要写入文件
	var fileWrite *os.File
	if *hashCmdWrite {
		var err error
		// 打开文件以写入
		fileWrite, err = os.OpenFile(globals.OutputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			errors <- fmt.Errorf("打开文件 %s 失败: %v", globals.OutputFileName, err)
			return []error{fmt.Errorf("打开文件 %s 失败: %v", globals.OutputFileName, err)}
		}
		defer fileWrite.Close()
	}

	// 遍历文件列表并启动任务
	for _, file := range files {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			errors <- ctx.Err()
			break
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
func hashTask(ctx context.Context, filepath string, hashType func() hash.Hash, fileWrite *os.File) error {
	// 检查上下文是否已取消
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 计算文件哈希值
	hashValue, checkErr := checksum(filepath, hashType)
	if checkErr != nil {
		return fmt.Errorf("计算文件 %s 哈希值失败: %v", filepath, checkErr)
	}

	// 检查上下文是否已取消（计算哈希值可能需要一些时间）
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 根据 *hashCmdWrite 标志决定是否写入文件
	if *hashCmdWrite {
		// 写入到 globals.OutputFileName 文件
		if fileWrite != nil {
			_, err := fileWrite.WriteString(fmt.Sprintf("%s\t%q\n", hashValue, filepath))
			if err != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", globals.OutputFileName, err)
			}
		}
	} else {
		fmt.Printf("%s\t%s\n", hashValue, filepath)
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

// 根据文件大小计算最佳缓冲区大小
func calculateBufferSize(fileSize int64) int {
	switch {
	case fileSize < 32*KB: // 小于 32KB 的文件使用 32KB 缓冲区
		return int(32 * KB)
	case fileSize < 128*KB: // 32KB-128KB 使用 64KB 缓冲区
		return int(32 * KB)
	case fileSize < 512*KB: // 128KB-512KB 使用 128KB 缓冲区
		return int(64 * KB)
	case fileSize < 1*MB: // 512KB-1MB 使用 256KB 缓冲区
		return int(128 * KB)
	case fileSize < 4*MB: // 1MB-4MB 使用 512KB 缓冲区
		return int(256 * KB)
	case fileSize < 16*MB: // 4MB-16MB 使用 1MB 缓冲区
		return int(512 * KB)
	case fileSize < 64*MB: // 16MB-64MB 使用 2MB 缓冲区
		return int(1 * MB)
	default: // 大于 64MB 的文件使用 4MB 缓冲区
		return int(2 * MB)
	}
}

// collectFiles 收集指定路径下的所有文件，根据recursive标志决定是否递归
func collectFiles(targetPath string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	var files []string

	// 检查路径是否包含通配符
	if strings.ContainsAny(targetPath, "*?[]{}") {
		// 处理包含通配符的路径
		matchedFiles, err := filepath.Glob(targetPath)
		if err != nil {
			return nil, fmt.Errorf("路径无效: %w", err)
		}

		if len(matchedFiles) == 0 {
			return nil, fmt.Errorf("没有找到匹配的文件")
		}

		for _, file := range matchedFiles {
			if info, err := os.Stat(file); err != nil {
				return nil, fmt.Errorf("无法获取文件信息: %w", err)
			} else if info.IsDir() { // 如果是目录，跳过
				cl.PrintWarnf("跳过目录：%s", file)
				continue
			} else if !info.IsDir() { // 如果是文件，添加到文件列表
				files = append(files, file)
				continue
			}
		}
	} else {
		// 处理普通路径（可能是文件或目录）
		info, err := os.Stat(targetPath)
		if err != nil {
			return nil, fmt.Errorf("无法获取路径信息: %w", err)
		}

		if info.IsDir() {
			// 如果是目录，根据递归标志决定处理方式
			if recursive {
				// 递归模式：遍历目录及其子目录中的所有文件
				err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					// 将非目录文件添加到文件列表
					if !d.IsDir() {
						files = append(files, path)
					}
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("遍历目录失败: %w", err)
				}
			} else {
				// 非递归模式：只获取目录下的直接文件
				dir, err := os.ReadDir(targetPath)
				if err != nil {
					return nil, fmt.Errorf("读取目录失败: %w", err)
				}
				// 遍历目录中的所有文件
				for _, entry := range dir {
					if entry.IsDir() {
						cl.PrintWarnf("跳过目录：%s, 请使用 -r 选项以递归方式处理", entry.Name())
						continue
					}

					// 将文件添加到文件列表
					files = append(files, filepath.Join(targetPath, entry.Name()))
				}
			}
		} else {
			// 如果是文件，直接添加到文件列表
			files = []string{targetPath}
		}
	}

	return files, nil
}
