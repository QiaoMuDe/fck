// Package hash 实现了文件哈希计算的并发任务管理功能。
// 该文件提供了哈希任务管理器，支持多线程并发计算文件哈希值，并可选择将结果写入文件。
package hash

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/hash"
)

// HashResult 哈希计算结果
type HashResult struct {
	FilePath  string // 文件路径
	HashValue string // 哈希值
	Error     error  // 错误信息
}

// WriteRequest 写入请求
type WriteRequest struct {
	Content string     // 要写入的内容
	Done    chan error // 完成通知通道
}

// HashTaskManager 哈希任务管理器
type HashTaskManager struct {
	// 配置参数
	config      HashConfig // 哈希配置
	files       []string   // 文件列表
	hashType    string     // 哈希类型
	concurrency int        // 并发数

	// 通道
	resultCh chan HashResult   // 哈希结果通道
	writeCh  chan WriteRequest // 写入请求通道

	// 控制
	ctx    context.Context         // 上下文
	cancel context.CancelCauseFunc // 上下文取消函数

	// 状态
	wg          sync.WaitGroup // 并发任务等待组
	writerWg    sync.WaitGroup // 写入协程等待组
	errors      []error        // 错误列表
	errorsMutex sync.Mutex     // 错误列表互斥锁

	// 统计
	processedCount atomic.Int64 // 已处理文件数
	errorCount     atomic.Int64 // 错误计数
}

// NewHashTaskManager 创建哈希任务管理器
//
// 参数:
//   - files: 文件列表
//   - hashType: 哈希类型
//
// 返回值:
//   - *HashTaskManager: 哈希任务管理器
func NewHashTaskManager(files []string, hashType string, config HashConfig) *HashTaskManager {
	ctx, cancel := context.WithCancelCause(context.Background())

	concurrency := runtime.NumCPU() * 2
	if concurrency > len(files) {
		concurrency = len(files)
	}

	if concurrency > 50 {
		concurrency = 50
	}

	return &HashTaskManager{
		config:      config,
		files:       files,
		hashType:    hashType,
		concurrency: concurrency,
		resultCh:    make(chan HashResult, concurrency*2),
		writeCh:     make(chan WriteRequest, 100),
		ctx:         ctx,
		cancel:      cancel,
		errors:      make([]error, 0),
	}
}

// Run 执行所有哈希任务
//
// 返回值:
//   - []error: 错误列表
func (m *HashTaskManager) Run() []error {
	defer m.cancel(nil)

	// 启动写入协程，负责将哈希结果写入文件
	if m.config.Write {
		m.writerWg.Go(
			func() {
				m.writerWorker()
			},
		)
	}

	// 启动结果收集协程，负责处理计算结果
	var resultWg sync.WaitGroup
	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		m.resultCollector()
	}()

	// 启动计算工作池
	m.startComputeWorkers()

	// 等待所有计算任务完成
	m.wg.Wait()
	close(m.resultCh)

	resultWg.Wait()

	// 关闭写入通道并等待写入协程结束
	if m.config.Write {
		close(m.writeCh)
		m.writerWg.Wait()
	}

	return m.errors
}

// startComputeWorkers 启动计算工作池
func (m *HashTaskManager) startComputeWorkers() {
	// 创建文件任务通道，用于分发文件路径给工作协程
	fileCh := make(chan string, m.concurrency)

	// 启动工作协程池，每个协程从通道中获取文件并计算哈希
	for i := 0; i < m.concurrency; i++ {
		m.wg.Go(
			func() {
				m.computeWorker(fileCh)
			},
		)
	}

	// 分发文件任务到通道
	go func() {
		defer close(fileCh)
		for _, file := range m.files {
			select {
			case fileCh <- file:
			case <-m.ctx.Done():
				return
			}
		}
	}()
}

// computeWorker 计算工作协程
//
// 参数:
//   - fileCh: 文件路径通道
func (m *HashTaskManager) computeWorker(fileCh <-chan string) {
	for {
		select {
		case filePath, ok := <-fileCh:
			if !ok {
				return // 通道已关闭
			}
			m.processFile(filePath)

		case <-m.ctx.Done():
			return // 上下文已取消
		}
	}
}

// processFile 处理单个文件
//
// 参数:
//   - filePath: 要处理的文件路径
func (m *HashTaskManager) processFile(filePath string) {
	defer func() {
		if r := recover(); r != nil {
			result := HashResult{
				FilePath: filePath,
				Error:    fmt.Errorf("panic when processing file %s: %v", filePath, r),
			}
			m.sendResult(result)
		}
	}()

	// 检查文件是否应跳过处理
	if skip, err := shouldSkipFile(filePath, m.config); err != nil {
		result := HashResult{
			FilePath: filePath,
			Error:    fmt.Errorf("failed to check file status %s: %w", filePath, err),
		}
		m.sendResult(result)
		return
	} else if skip {
		return
	}

	result := HashResult{
		FilePath: filePath,
	}

	// 解析哈希算法
	hashType, err := hash.ParseAlgorithm(m.hashType)
	if err != nil {
		result.Error = fmt.Errorf("failed to parse hash algorithm %s: %w", m.hashType, err)
		m.sendResult(result)
		return
	}

	// 根据配置选择哈希计算方式
	if m.config.Progress {
		result.HashValue, result.Error = hash.ChecksumProgress(filePath, hashType)
	} else {
		result.HashValue, result.Error = hash.Checksum(filePath, hashType)
	}

	m.sendResult(result)
}

// sendResult 发送计算结果
//
// 参数：
//   - result: 计算结果
func (m *HashTaskManager) sendResult(result HashResult) {
	select {
	case m.resultCh <- result:
	case <-m.ctx.Done():
		// 上下文已取消，记录错误
		m.addError(fmt.Errorf("failed to send result %s: task cancelled", result.FilePath))
	}
}

// resultCollector 结果收集协程
// 负责从结果通道读取哈希计算结果，输出到终端或发送写入请求
func (m *HashTaskManager) resultCollector() {
	for result := range m.resultCh {
		if result.Error != nil {
			m.addError(result.Error)
			m.errorCount.Add(1)
			continue
		}

		// 非写入模式下直接输出到终端
		if !m.config.Write {
			if m.config.Quote {
				fmt.Printf("%s\t%q\n", result.HashValue, result.FilePath)
			} else {
				fmt.Printf("%s\t%s\n", result.HashValue, result.FilePath)
			}
		}

		// 写入模式下发送写入请求
		if m.config.Write {
			var content string
			if m.config.Quote {
				content = fmt.Sprintf("%s\t%q\n", result.HashValue, result.FilePath)
			} else {
				content = fmt.Sprintf("%s\t%s\n", result.HashValue, result.FilePath)
			}
			m.requestWrite(content)
		}

		m.processedCount.Add(1)
	}
}

// requestWrite 请求写入
// 采用同步写入机制：发送请求后等待写入完成，确保写入顺序与计算顺序一致
//
// 参数:
//   - content: 要写入的内容
func (m *HashTaskManager) requestWrite(content string) {
	// 检查上下文是否已取消
	if m.ctx.Err() != nil {
		return
	}

	req := WriteRequest{
		Content: content,             // 要写入的内容
		Done:    make(chan error, 1), // 写入完成信号
	}

	// 尝试发送写入请求
	select {
	case m.writeCh <- req:
		// 成功发送，等待写入完成
		if err := <-req.Done; err != nil {
			m.addError(fmt.Errorf("failed to write content to file: %w", err))
		}
	case <-m.ctx.Done():
		// 上下文已取消
		return
	default:
		// 写入通道满或已关闭，记录警告
		m.addError(fmt.Errorf("write channel is full or closed, skip content write"))
	}
}

// writerWorker 写入工作协程
func (m *HashTaskManager) writerWorker() {
	// 初始化文件写入器
	wrapper, err := m.initFileWriter()
	if err != nil {
		m.addError(fmt.Errorf("failed to init file writer: %w", err))
		return
	}
	defer m.closeWriter(wrapper)

	// 处理写入请求
	for req := range m.writeCh {
		err := m.writeContent(wrapper, req.Content)
		req.Done <- err
	}
}

// initFileWriter 初始化文件写入器
//
// 返回值:
//   - *FileWriterWrapper: 文件写入器包装
//   - error: 错误信息，如果发生错误则返回非nil值
func (m *HashTaskManager) initFileWriter() (*FileWriterWrapper, error) {
	file, err := os.OpenFile(types.OutputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", types.OutputFileName, err)
	}

	if err := m.writeFileHeader(file, m.hashType); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("failed to write file header: %w", err)
	}

	return &FileWriterWrapper{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// writeFileHeader 写入文件头信息
//
// 参数:
//   - file: 要写入的文件对象
//   - hashType: 哈希类型标识
//
// 返回:
//   - error: 错误信息，如果写入失败
//
// 注意:
//   - 便携模式文件头格式为: #hashType#timestamp#PORTABLE
//   - 本地模式文件头格式为: #hashType#timestamp#LOCAL#basePath
func (m *HashTaskManager) writeFileHeader(file *os.File, hashType string) error {
	header := &types.ChecksumHeader{
		HashType:  hashType,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	if m.config.Local {
		header.Mode = types.ChecksumModeLocal
		header.BasePath = m.config.BasePath
		if header.BasePath == "" {
			var err error
			header.BasePath, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %v", err)
			}
		}
	} else {
		header.Mode = types.ChecksumModePortable
	}

	headerStr := header.String()

	if _, err := file.WriteString(headerStr); err != nil {
		return fmt.Errorf("failed to write file header: %v", err)
	}
	return nil
}

// writeContent 写入内容
//
// 参数:
//   - wrapper: 文件写入器包装
//   - content: 要写入的内容
//
// 返回值:
//   - error: 错误信息，如果发生错误则返回非nil值
func (m *HashTaskManager) writeContent(wrapper *FileWriterWrapper, content string) error {
	_, err := wrapper.writer.WriteString(content)
	return err
}

// FileWriterWrapper 文件写入器包装
type FileWriterWrapper struct {
	file   *os.File
	writer *bufio.Writer
}

// closeWriter 关闭写入器
//
// 参数:
//   - wrapper: 文件写入器包装
func (m *HashTaskManager) closeWriter(wrapper *FileWriterWrapper) {
	var errs []error

	if err := wrapper.writer.Flush(); err != nil {
		errs = append(errs, fmt.Errorf("failed to flush writer: %w", err))
	}

	if err := wrapper.file.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close file: %w", err))
	}

	if len(errs) > 0 {
		for _, err := range errs {
			m.addError(err)
		}
	}
}

// addError 线程安全地添加错误
//
// 参数:
//   - err: 错误信息
func (m *HashTaskManager) addError(err error) {
	m.errorsMutex.Lock()
	defer m.errorsMutex.Unlock()
	m.errors = append(m.errors, err)
}

// GetStats 获取统计信息
//
// 返回值:
//   - processed: 已处理的文件数
//   - errors: 错误数
func (m *HashTaskManager) GetStats() (processed, errors int64) {
	return m.processedCount.Load(), m.errorCount.Load()
}

// hashRunTasksRefactored 重构后的任务执行函数
//
// 参数:
//   - files: 文件列表
//   - hashType: 哈希类型函数
//
// 返回:
//   - int64: 实际处理的文件数
//   - []error: 错误列表
func hashRunTasksRefactored(files []string, hashType string, config HashConfig) (int64, []error) {
	if len(files) == 0 {
		return 0, nil
	}

	manager := NewHashTaskManager(files, hashType, config)
	errors := manager.Run()
	processed, _ := manager.GetStats()
	return processed, errors
}

// shouldSkipFile 检查是否应该跳过文件
//
// 参数:
//   - filePath: 文件路径
//   - config: 哈希配置
//
// 返回:
//   - bool: 如果应该跳过文件，则返回true；否则返回false
//   - error: 错误信息，如果发生错误则返回非nil值
func shouldSkipFile(filePath string, config HashConfig) (bool, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return false, err
	}

	// 跳过软链接
	if fileInfo.Mode()&fs.ModeSymlink != 0 {
		return true, nil
	}

	// 写入模式下跳过输出文件，避免校验文件本身被计入
	// 使用 filepath.Clean 清理路径后直接比较，零系统调用
	if config.Write {
		cleanedPath := filepath.Clean(filePath)
		if cleanedPath == types.OutputFileName {
			return true, nil
		}
	}

	return false, nil
}
