package check

import (
	"fmt"
	"hash"
	"os"
	"runtime"
	"sync"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// fileChecker 文件校验器
type fileChecker struct {
	cl         *colorlib.ColorLib // 颜色库
	hashFunc   func() hash.Hash   // 哈希函数
	maxWorkers int                // 最大并发数(默认: 逻辑处理器数量)
}

// newFileChecker 创建新的文件校验器
func newFileChecker(cl *colorlib.ColorLib, hashFunc func() hash.Hash) *fileChecker {
	return &fileChecker{
		cl:         cl,
		hashFunc:   hashFunc,
		maxWorkers: runtime.NumCPU(),
	}
}

// checkResult 校验结果
type checkResult struct {
	filePath     string // 文件路径
	expectedHash string // 期望的哈希值
	actualHash   string // 实际的哈希值
	err          error  // 错误信息
}

// checkFiles 并发校验文件
func (c *fileChecker) checkFiles(hashMap types.VirtualHashMap) error {
	if len(hashMap) == 0 {
		c.cl.PrintWarnf("没有文件需要校验\n")
		return nil
	}

	// 创建工作通道
	jobs := make(chan types.VirtualHashEntry, len(hashMap))
	results := make(chan checkResult, len(hashMap))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < c.maxWorkers; i++ {
		wg.Add(1)
		go c.worker(jobs, results, &wg)
	}

	// 发送任务
	go func() {
		defer close(jobs)
		for _, entry := range hashMap {
			jobs <- entry
		}
	}()

	// 等待所有工作完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	return c.collectResults(results, len(hashMap))
}

// worker 工作协程
func (c *fileChecker) worker(jobs <-chan types.VirtualHashEntry, results chan<- checkResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for entry := range jobs {
		result := checkResult{
			filePath:     entry.RealPath,
			expectedHash: entry.Hash,
		}

		// 检查文件是否存在
		if _, err := os.Stat(entry.RealPath); err != nil {
			c.cl.PrintWarnf("文件 %s 不存在，跳过校验\n", entry.RealPath)
			continue
		}

		// 计算文件哈希
		actualHash, err := common.Checksum(entry.RealPath, c.hashFunc)
		if err != nil {
			result.err = fmt.Errorf("计算文件哈希失败: %v", err)
		} else {
			result.actualHash = actualHash
		}

		results <- result
	}
}

// collectResults 收集校验结果
func (c *fileChecker) collectResults(results <-chan checkResult, totalFiles int) error {
	var (
		errorCount     int
		mismatchCount  int
		processedCount int
	)

	for result := range results {
		processedCount++

		if result.err != nil {
			c.cl.Redf("%s ✗\n", result.filePath)
			errorCount++
			continue
		}

		// 比较哈希值
		if result.actualHash != result.expectedHash {
			c.cl.Redf("%s ✗\n", result.filePath)
			mismatchCount++
		} else {
			c.cl.Greenf("%s ✓\n", result.filePath)
		}
	}

	// 输出校验结果统计
	c.printSummary(processedCount, mismatchCount, errorCount, totalFiles)

	return nil
}

// printSummary 打印校验结果摘要
func (c *fileChecker) printSummary(processed, mismatched, errors, total int) {
	passedCount := processed - mismatched - errors
	c.cl.Bluef("完成: %d/%d 通过\n", passedCount, total)
}
