package hash

import (
	"bufio"
	"context"
	"fmt"
	"hash"
	"io/fs"
	"os"
	"runtime"
	"sync"

	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// FileWriter 文件写入接口
type FileWriter interface {
	Write(content string) error
	Close() error
}

// safeFileWriter 安全的文件写入器
type safeFileWriter struct {
	file   *os.File
	writer *bufio.Writer
	mutex  sync.Mutex
}

// newSafeFileWriter 创建安全的文件写入器
func newSafeFileWriter(filename, hashType string) (*safeFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开文件 %s 失败: %w", filename, err)
	}

	writer := bufio.NewWriter(file)

	// 写入文件头
	if err := common.WriteFileHeader(file, hashType, types.TimestampFormat); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("写入文件头失败: %w", err)
	}

	return &safeFileWriter{
		file:   file,
		writer: writer,
	}, nil
}

// Write 线程安全的写入
func (w *safeFileWriter) Write(content string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	_, err := w.writer.WriteString(content)
	return err
}

// Close 关闭文件写入器
func (w *safeFileWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err := w.writer.Flush(); err != nil {
		_ = w.file.Close()
		return err
	}

	return w.file.Close()
}

// hashTaskRunner 封装任务执行逻辑
type hashTaskRunner struct {
	ctx         context.Context
	cancel      context.CancelCauseFunc
	concurrency int
	files       []string
	hashType    func() hash.Hash
	errors      chan error
	results     []error
	once        sync.Once
	fileWriter  *safeFileWriter
}

// hashRunTasks 执行哈希值校验任务，支持并发控制和错误处理
func hashRunTasks(files []string, hashType func() hash.Hash) []error {
	if len(files) == 0 {
		return nil
	}

	// 创建带取消原因的上下文
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	// 根据任务特性调整并发数（I/O密集型任务可以适当增加）
	concurrency := runtime.NumCPU() * 2
	if len(files) < concurrency {
		concurrency = len(files)
	}

	// 创建任务执行器
	runner := &hashTaskRunner{
		ctx:         ctx,
		cancel:      cancel,
		concurrency: concurrency,
		files:       files,
		hashType:    hashType,
		errors:      make(chan error, len(files)),    // 足够的缓冲区
		results:     make([]error, 0, len(files)/10), // 预估错误率
	}

	return runner.run()
}

// run 执行所有任务
func (r *hashTaskRunner) run() []error {
	// 初始化文件写入器
	if err := r.initFileWriter(); err != nil {
		return []error{err}
	}
	defer r.closeFileWriter()

	// 启动错误收集器
	go r.collectErrors()

	// 启动工作池
	r.runWorkerPool()

	// 等待所有任务完成
	close(r.errors)

	// 添加上下文取消错误
	r.addContextError()

	return r.results
}

// initFileWriter 初始化文件写入器
func (r *hashTaskRunner) initFileWriter() error {
	if !hashCmdWrite.Get() {
		return nil
	}

	writer, err := newSafeFileWriter(types.OutputFileName, hashCmdType.Get())
	if err != nil {
		return fmt.Errorf("初始化文件写入器失败: %w", err)
	}

	r.fileWriter = writer
	return nil
}

// runWorkerPool 运行工作池
func (r *hashTaskRunner) runWorkerPool() {
	var wg sync.WaitGroup
	jobPool := make(chan struct{}, r.concurrency)

	for _, file := range r.files {
		// 检查上下文是否已取消
		select {
		case <-r.ctx.Done():
			r.handleError(r.ctx.Err())
			goto cleanup
		default:
		}

		jobPool <- struct{}{}
		wg.Add(1)

		go func(filePath string) {
			defer func() {
				wg.Done()
				<-jobPool
				r.recoverPanic(filePath)
			}()

			r.processFile(filePath)
		}(file)
	}

cleanup:
	wg.Wait()
}

// processFile 处理单个文件
func (r *hashTaskRunner) processFile(filePath string) {
	var writer FileWriter
	if r.fileWriter != nil {
		writer = r.fileWriter
	}
	if err := hashTask(r.ctx, filePath, r.hashType, writer); err != nil {
		r.handleError(err)
	}
}

// handleError 统一错误处理
func (r *hashTaskRunner) handleError(err error) {
	select {
	case r.errors <- err:
	case <-r.ctx.Done():
		// 上下文已取消，不再发送错误
		return
	}
	r.once.Do(func() { r.cancel(err) })
}

// recoverPanic 恢复panic
func (r *hashTaskRunner) recoverPanic(filePath string) {
	if rec := recover(); rec != nil {
		err := fmt.Errorf("处理文件 %s 时发生 panic: %v", filePath, rec)
		r.handleError(err)
	}
}

// collectErrors 收集错误
func (r *hashTaskRunner) collectErrors() {
	for err := range r.errors {
		r.results = append(r.results, err)
	}
}

// addContextError 添加上下文错误
func (r *hashTaskRunner) addContextError() {
	if r.ctx.Err() != nil {
		if cause := context.Cause(r.ctx); cause != nil {
			r.results = append(r.results, fmt.Errorf("任务被取消: %w", cause))
		} else {
			r.results = append(r.results, r.ctx.Err())
		}
	}
}

// closeFileWriter 关闭文件写入器
func (r *hashTaskRunner) closeFileWriter() {
	if r.fileWriter != nil {
		if err := r.fileWriter.Close(); err != nil {
			r.results = append(r.results, fmt.Errorf("关闭文件写入器失败: %v", err))
		}
	}
}

// shouldSkipFile 检查是否应该跳过文件
func shouldSkipFile(filePath string) (bool, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return false, err
	}

	// 跳过软链接
	return fileInfo.Mode()&fs.ModeSymlink != 0, nil
}

// computeHashWithContext 支持上下文取消的哈希计算
func computeHashWithContext(ctx context.Context, filePath string, hashType func() hash.Hash) (string, error) {
	// 这里可以在 common.Checksum 中添加上下文支持
	// 或者使用带超时的上下文
	done := make(chan struct{})
	var hashValue string
	var err error

	go func() {
		defer close(done)
		hashValue, err = common.Checksum(filePath, hashType)
	}()

	select {
	case <-done:
		return hashValue, err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// outputResult 统一输出结果
func outputResult(hashValue, filePath string, writer FileWriter) error {
	output := formatHashOutput(hashValue, filePath)

	if writer != nil {
		return writer.Write(output)
	}

	fmt.Print(output)
	return nil
}

// formatHashOutput 统一输出格式
func formatHashOutput(hashValue, filePath string) string {
	return fmt.Sprintf("%s\t%q\n", hashValue, filePath)
}

// hashTask 处理单个文件的哈希计算，支持上下文取消
func hashTask(ctx context.Context, filePath string, hashType func() hash.Hash, writer FileWriter) error {
	// 检查上下文是否已取消
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// 检查文件状态（软链接检查）
	if skip, err := shouldSkipFile(filePath); err != nil {
		return fmt.Errorf("检查文件 %s 状态失败: %w", filePath, err)
	} else if skip {
		return nil // 跳过软链接
	}

	// 计算文件哈希值（支持上下文取消）
	hashValue, err := computeHashWithContext(ctx, filePath, hashType)
	if err != nil {
		return fmt.Errorf("计算文件 %s 哈希值失败: %w", filePath, err)
	}

	// 输出结果
	return outputResult(hashValue, filePath, writer)
}
