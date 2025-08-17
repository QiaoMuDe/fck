// Package find 实现了文件查找的并发搜索功能。
// 该文件提供了多线程并发搜索器，用于提高大目录结构的搜索性能。
package find

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// ConcurrentSearcher 负责并发搜索协调
type ConcurrentSearcher struct {
	searcher   *FileSearcher // 基础搜索器
	maxWorkers int           // 最大并发worker数量
}

// NewConcurrentSearcher 创建新的并发搜索器
//
// 参数:
//   - searcher: 基础搜索器
//   - maxWorkers: 最大并发worker数量
//
// 返回:
//   - ConcurrentSearcher: 并发搜索器对象
func NewConcurrentSearcher(searcher *FileSearcher, maxWorkers int) *ConcurrentSearcher {
	// 检查最大并发worker数量, 默认值为 CPU 核心数
	if maxWorkers <= 0 {
		maxWorkers = runtime.NumCPU()
	}
	return &ConcurrentSearcher{
		searcher:   searcher,
		maxWorkers: maxWorkers,
	}
}

// SearchConcurrent 执行并发搜索
//
// 参数:
//   - findPath: 搜索路径
//
// 返回:
//   - error: 错误信息
func (cs *ConcurrentSearcher) SearchConcurrent(findPath string) error {
	var wg sync.WaitGroup
	pathChan := make(chan string, cs.maxWorkers*100) // 减少缓冲区大小
	errorChan := make(chan error, 50)                // 错误通道

	// 启动worker goroutines处理路径
	for i := 0; i < cs.maxWorkers; i++ {
		wg.Add(1)
		go cs.worker(&wg, pathChan, errorChan)
	}

	// 主逻辑通过 WalkDir 将路径发送到 channel
	walkDirErr := filepath.WalkDir(findPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return cs.handleWalkError(p, err)
		}

		// 跳过findPath本身
		if p == findPath {
			return nil
		}

		// 检查深度限制
		if cs.exceedsMaxDepth(p, findPath) {
			return filepath.SkipDir
		}

		// 检查隐藏项
		if cs.shouldSkipHidden(p) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 检查符号链接循环
		if d.Type()&os.ModeSymlink != 0 {
			if cs.searcher.isSymlinkLoop(p) {
				return filepath.SkipDir
			}
		}

		// 发送路径到 channel
		select {
		case pathChan <- p:
		default:
			// 如果通道满了，直接处理
			return cs.processPathDirect(p, d)
		}

		return nil
	})

	// 遍历完成, 关闭路径通道
	close(pathChan)

	// 等待所有 worker 完成
	wg.Wait()

	// 关闭错误通道
	close(errorChan)

	// 处理错误
	return cs.handleErrors(walkDirErr, errorChan)
}

// worker 处理路径的工作协程
//
// 参数:
//   - wg: WaitGroup 对象
//   - pathChan: 路径通道
//   - errorChan: 错误通道
func (cs *ConcurrentSearcher) worker(wg *sync.WaitGroup, pathChan <-chan string, errorChan chan<- error) {
	defer func() {
		// 捕获 panic
		if r := recover(); r != nil {
			select {
			case errorChan <- fmt.Errorf("worker goroutine 发生 panic: %v\nstack trace: %s", r, debug.Stack()):
			default:
			}
		}
		wg.Done()
	}()

	// 读取通道处理路径
	for path := range pathChan {
		if err := cs.processPath(path); err != nil {
			select {
			case errorChan <- err:
			default:
			}
		}
	}
}

// processPath 处理单个路径
//
// 参数:
//   - path: 路径
//
// 返回:
//   - error: 错误信息
func (cs *ConcurrentSearcher) processPath(path string) error {
	// 检查是否存在
	entry, statErr := os.Stat(path)
	if statErr != nil {
		return cs.handleWalkError(path, statErr)
	}

	// 构建 DirEntryWrapper 包装器
	dirEntry := &types.DirEntryWrapper{
		NameVal:  entry.Name(),
		IsDirVal: entry.IsDir(),
		ModeVal:  entry.Mode(),
	}

	// 通过 processEntry 处理路径
	processErr := cs.searcher.processEntry(dirEntry, path)
	if processErr != nil {
		if processErr == filepath.SkipDir {
			// 如果是 SkipDir, 跳过该目录即可
			return nil
		}

		if errors.Is(processErr, os.ErrPermission) {
			cs.searcher.config.Cl.PrintErrorf("路径 %s 权限不足, 已跳过\n", path)
			return nil
		}

		return fmt.Errorf("处理路径失败 %s: %v", path, processErr)
	}

	return nil
}

// processPathDirect 直接处理路径（当通道满时使用）
//
// 参数:
//   - path: 路径
//   - d: 目录项
//
// 返回:
//   - error: 错误信息
func (cs *ConcurrentSearcher) processPathDirect(path string, d os.DirEntry) error {
	return cs.searcher.processEntry(d, path)
}

// handleWalkError 便捷处理错误
//
// 参数:
//   - path: 路径
//   - err: 错误信息
//
// 返回:
//   - error: 错误信息
func (cs *ConcurrentSearcher) handleWalkError(path string, err error) error {
	// 忽略不存在的路径错误
	if os.IsNotExist(err) {
		return nil
	}

	// 忽略权限错误
	if os.IsPermission(err) {
		cs.searcher.config.Cl.PrintErrorf("路径 %s 因权限不足已跳过\n", path)
		return nil
	}

	return fmt.Errorf("访问 %s 时出错: %s", path, err)
}

// exceedsMaxDepth 检查是否超过最大深度
//
// 参数:
//   - path: 路径
//   - findPath: 查找路径
//
// 返回:
//   - bool: 是否超过最大深度
func (cs *ConcurrentSearcher) exceedsMaxDepth(path, findPath string) bool {
	depth := strings.Count(path[len(findPath):], string(filepath.Separator)) // 计算当前路径的深度
	return findCmdMaxDepth.Get() >= 0 && depth > findCmdMaxDepth.Get()
}

// shouldSkipHidden 检查是否应该跳过隐藏文件
//
// 参数:
//   - path: 路径
//
// 返回:
//   - bool: 是否应该跳过隐藏文件
func (cs *ConcurrentSearcher) shouldSkipHidden(path string) bool {
	return !findCmdHidden.Get() && common.IsHidden(path)
}

// handleErrors 处理所有错误
//
// 参数:
//   - walkDirErr: 遍历目录错误
//   - errorChan: 错误通道
//
// 返回:
//   - error: 错误信息
func (cs *ConcurrentSearcher) handleErrors(walkDirErr error, errorChan <-chan error) error {
	// 优先检查遍历过程中是否遇到错误
	if walkDirErr != nil {
		if os.IsPermission(walkDirErr) {
			return fmt.Errorf("权限不足, 无法访问某些目录: %v", walkDirErr)
		} else if os.IsNotExist(walkDirErr) {
			return fmt.Errorf("路径不存在: %v", walkDirErr)
		}
		return fmt.Errorf("遍历目录时出错: %v", walkDirErr)
	}

	// 收集并分类错误（限制最多显示5个不同错误）
	errorMap := make(map[string]error)
	errorCount := 0
	for err := range errorChan {
		if err != nil && errorCount < 5 {
			if _, exists := errorMap[err.Error()]; !exists {
				errorMap[err.Error()] = err
				errorCount++
			}
		}
	}

	// 合并错误并添加统计信息
	if len(errorMap) > 0 {
		var combinedErr error
		var errorList []string
		for _, err := range errorMap {
			combinedErr = errors.Join(combinedErr, err)
			errorList = append(errorList, fmt.Sprintf("\t- %s", err))
		}
		return fmt.Errorf("共发现%d类错误:\n%s", len(errorMap), strings.Join(errorList, "\n"))
	}

	return nil
}
