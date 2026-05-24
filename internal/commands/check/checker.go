// Package check 实现了文件哈希校验功能。
// 该文件提供了并发文件校验器，用于验证文件的完整性和一致性。
package check

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"gitee.com/MM-Q/color"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/hash"
)

// fileChecker 文件校验器
type fileChecker struct {
	cl         *color.GlobalColor // 颜色库
	hashType   string             // 哈希算法
	maxWorkers int                // 最大并发数(默认: 逻辑处理器数量)
	verbose    bool               // 详细模式（显示校验通过的文件）
}

// newFileChecker 创建新的文件校验器
func newFileChecker(cl *color.GlobalColor, hashType string, verbose bool) *fileChecker {
	return &fileChecker{
		cl:         cl,
		hashType:   hashType,
		maxWorkers: runtime.NumCPU(),
		verbose:    verbose,
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
		c.cl.Yellow("no files to check")
		return nil
	}

	// 创建工作通道
	jobs := make(chan types.VirtualHashEntry, len(hashMap))
	results := make(chan checkResult, len(hashMap))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < c.maxWorkers; i++ {
		wg.Go(
			func() {
				c.worker(jobs, results)
			},
		)
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
func (c *fileChecker) worker(jobs <-chan types.VirtualHashEntry, results chan<- checkResult) {
	for entry := range jobs {
		result := checkResult{
			filePath:     entry.RealPath,
			expectedHash: entry.Hash,
		}

		// 检查文件是否存在，如果不存在则发送错误结果
		if _, err := os.Stat(entry.RealPath); err != nil {
			results <- checkResult{
				filePath: entry.RealPath,
				err:      err,
			}
			continue
		}

		// 解析哈希算法
		hashType, err := hash.ParseAlgorithm(c.hashType)
		if err != nil {
			results <- checkResult{
				filePath: entry.RealPath,
				err:      err,
			}
		}

		// 计算文件哈希
		actualHash, err := hash.Checksum(entry.RealPath, hashType)
		if err != nil {
			result.err = fmt.Errorf("failed to calculate checksum: %v", err)
		} else {
			result.actualHash = actualHash
		}

		results <- result
	}
}

// collectResults 收集校验结果
func (c *fileChecker) collectResults(results <-chan checkResult, totalFiles int) error {
	var (
		passedCount    int // 校验通过的文件数
		mismatchCount  int // 哈希不匹配的文件数
		notFoundCount  int // 文件不存在的文件数
		errorCount     int // 其他错误的文件数
		processedCount int // 总处理文件数
	)

	for result := range results {
		processedCount++

		if result.err != nil {
			// 检查是否是文件不存在错误
			if os.IsNotExist(result.err) ||
				strings.Contains(result.err.Error(), "不存在") ||
				strings.Contains(result.err.Error(), "not found") ||
				strings.Contains(result.err.Error(), "no such file") {
				c.cl.Yellowf("file %s not found, skip verification\n", result.filePath)
				notFoundCount++
			} else {
				c.cl.Redf("%s ✗ (error: %v)\n", result.filePath, result.err)
				errorCount++
			}
			continue
		}

		// 比较哈希值
		if result.actualHash != result.expectedHash {
			c.cl.Redf("%s ✗ (hash mismatch)\n", result.filePath)
			mismatchCount++
		} else {
			if c.verbose {
				c.cl.Greenf("%s ✓\n", result.filePath)
			}
			passedCount++
		}
	}

	// 输出校验结果统计
	c.printSummary(passedCount, mismatchCount, notFoundCount, errorCount, totalFiles)

	return nil
}

// printSummary 打印校验结果摘要
func (c *fileChecker) printSummary(passed, mismatched, notFound, errors, total int) {
	result := c.cl.SBlue("verification completed: ")
	result += c.cl.SGreenf("%d passed", passed)

	if mismatched > 0 {
		result += ", "
		result += c.cl.SRedf("%d failed", mismatched)
	}

	if notFound > 0 {
		result += ", "
		result += c.cl.SYellowf("%d not found", notFound)
	}

	if errors > 0 {
		result += ", "
		result += c.cl.SRedf("%d errors", errors)
	}

	result += c.cl.SWhitef(" (total: %d files)", total)
	fmt.Println(result)
}
