package hash

import (
	"context"
	"fmt"
	"hash"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// HashCmdMain 是 hash 子命令的主函数
func HashCmdMain(cl *colorlib.ColorLib) error {
	// 获取指定的路径
	targetPaths := hashCmd.Args()

	// 如果没有指定路径，则打印错误信息并退出
	if len(targetPaths) == 0 {
		return fmt.Errorf("请指定要计算哈希值的路径")
	}

	// 检查指定的哈希算法是否有效
	hashType, ok := types.SupportedAlgorithms[hashCmdType.Get()]
	if !ok {
		return fmt.Errorf("在校验哈希值时，哈希算法 %s 无效", hashCmdType.Get())
	}

	// 创建一个可以被外部取消的上下文（支持指定取消原因）
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil) // 正常完成时不设置取消原因

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
				cl.PrintOkf("已将哈希值写入文件 %s, 共处理 %d 个文件\n", types.OutputFileName, len(files))
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

	// 根据 CPU 核心数创建工作池
	jobPool := make(chan struct{}, runtime.NumCPU())

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
		fileWrite, err = os.OpenFile(types.OutputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			errors <- fmt.Errorf("打开文件 %s 失败: %v", types.OutputFileName, err)
			return []error{fmt.Errorf("打开文件 %s 失败: %v", types.OutputFileName, err)}
		}
		defer func() {
			if err := fileWrite.Close(); err != nil {
				errors <- fmt.Errorf("close file failed: %v", err)
			}
		}()

		// 写入文件头
		if err := common.WriteFileHeader(fileWrite, hashCmdType.Get(), types.TimestampFormat); err != nil {
			errors <- fmt.Errorf("写入文件头 %s 失败: %v", types.OutputFileName, err)
			return []error{fmt.Errorf("写入文件头 %s 失败: %v", types.OutputFileName, err)}
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
	hashValue, checkErr := common.Checksum(filePath, hashType)
	if checkErr != nil {
		return fmt.Errorf("计算文件 %s 哈希值失败: %v", filePath, checkErr)
	}

	// 检查上下文是否已取消（计算哈希值可能需要一些时间）
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 根据 hashCmdWrite.Get() 标志决定是否写入文件
	if hashCmdWrite.Get() {
		// 写入到 types.OutputFileName 文件
		if fileWrite != nil {
			_, err := fmt.Fprintf(fileWrite, "%s\t%q\n", hashValue, filePath)
			if err != nil {
				return fmt.Errorf("写入文件 %s 失败: %v", types.OutputFileName, err)
			}
		}
	} else {
		fmt.Printf("%s\t%s\n", hashValue, filePath)
	}

	return nil
}
