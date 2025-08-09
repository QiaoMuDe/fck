package size

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// 用于存储在输出时的项
type item struct {
	Name string
	Size string
}

// 用于存储在输出时的项列表
type items []item

// SizeCmdMain 是 size 子命令的主函数
//
// 参数:
//   - cl: 用于打印输出的 ColorLib 对象
//
// 返回:
// - error: 如果发生错误，返回错误信息，否则返回 nil
func SizeCmdMain(cl *colorlib.ColorLib) error {
	// 获取指定的路径
	targetPaths := sizeCmd.Args()

	// 如果没有指定路径, 则默认计算当前目录下每个项目的大小
	if len(targetPaths) == 0 {
		return fmt.Errorf("请指定要计算大小路径")
	}

	// 根据sizeCmdColor设置颜色模式
	if sizeCmdColor.Get() {
		cl.NoColor.Store(false)
	} else {
		cl.NoColor.Store(true)
	}

	// 新建一个表格项列表
	var itemList items

	// 遍历路径
	for _, targetPath := range targetPaths {
		// 清理路径
		targetPath = filepath.Clean(targetPath)

		// 如果路径包含通配符, 则使用通配符匹配路径
		if strings.Contains(targetPath, "*") {
			filePaths, err := filepath.Glob(targetPath)
			if err != nil {
				cl.PrintErrf("通配符匹配失败: %v\n", err)
				continue
			}

			if len(filePaths) == 0 {
				cl.PrintErr("没有找到匹配的文件")
				continue
			}

			// 计算每个匹配路径的大小
			for _, filePath := range filePaths {
				// 如果是隐藏文件且未启用-H选项, 则跳过
				if !sizeCmdHidden.Get() {
					if common.IsHidden(filePath) {
						continue
					}
				}

				// 获取文件大小
				size, err := getPathSize(filePath)
				if err != nil {
					cl.PrintErrf("计算大小时出现错误: %s\n", err)
				}

				// 添加到 items 数组中
				itemList = append(itemList, item{
					Name: filePath,
					Size: humanReadableSize(size, 2),
				})
				continue
			}

			// 打印输出
			printSizeTable(itemList, cl)
			return nil
		}

		// 如果是隐藏文件且未启用-H选项, 则跳过
		if !sizeCmdHidden.Get() {
			if common.IsHidden(targetPath) {
				continue
			}
		}

		// 获取文件信息
		info, statErr := os.Lstat(targetPath)
		if statErr != nil {
			cl.PrintErrf("获取文件信息失败: 路径 %s 错误: %v\n", targetPath, statErr)
			continue
		}

		// 如果是文件, 则直接计算大小
		if !info.IsDir() {
			// 添加到 items 数组中
			itemList = append(itemList, item{
				Name: targetPath,
				Size: humanReadableSize(info.Size(), 2),
			})
			// 打印输出
			printSizeTable(itemList, cl)
			continue
		}

		// 如果是目录, 则递归计算大小
		size, getErr := getPathSize(targetPath)
		if getErr != nil {
			cl.PrintErrf("计算目录大小失败: 路径 %s 错误: %v\n", targetPath, getErr)
			continue
		}

		// 添加到 items 数组中
		itemList = append(itemList, item{
			Name: targetPath,
			Size: humanReadableSize(size, 2),
		})

		// 打印输出
		printSizeTable(itemList, cl)
		continue
	}

	return nil
}

// getPathSize 获取路径大小
func getPathSize(path string) (int64, error) {
	// 如果是隐藏文件且未启用-H选项，则跳过
	if !sizeCmdHidden.Get() {
		if common.IsHidden(path) {
			return 0, nil
		}
	}

	// 获取文件信息
	info, statErr := os.Lstat(path)
	if statErr != nil {
		// 检查是否为权限错误
		if os.IsPermission(statErr) {
			// 如果是权限错误, 则忽略该错误并给出更友好的提示信息
			return 0, fmt.Errorf("权限不足: 路径 %s", path)
		}
		// 检查是否为文件不存在错误
		if os.IsNotExist(statErr) {
			// 如果是文件不存在错误, 则忽略该错误并给出更友好的提示信息
			return 0, fmt.Errorf("文件不存在: 路径 %s", path)
		}

		// 如果获取文件信息失败, 返回错误
		return 0, fmt.Errorf("获取文件信息失败: 路径 %s 错误: %v", path, statErr)
	}

	// 如果不是目录, 则直接返回文件大小
	if !info.IsDir() {
		return info.Size(), nil
	}

	// 定义总大小 (使用原子类型)
	var totalSize atomic.Int64

	// 文件信息结构体
	type fileInfo struct {
		path string
		size int64
	}

	// 创建并发任务池
	var wg sync.WaitGroup
	fileChan := make(chan fileInfo, 10000) // 缓冲区大小为10000
	errChan := make(chan error, 100)       // 错误通道大小为100

	// 启动worker goroutines
	workers := sizeCmdJob.Get()
	if workers == -1 {
		// 如果未指定并发数量, 则根据CPU核心数自动设置 (每个核心2个worker)
		workers = runtime.NumCPU() * 2
	}
	if workers <= 0 {
		// 如果并发数量小于等于0, 则使用单线程执行
		workers = 1
	}
	if workers > 20 {
		// 如果并发数量大于20, 则限制为20
		workers = 20
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				totalSize.Add(file.size) // 原子操作累加文件大小
			}
		}()
	}

	// 遍历目录 (使用更高效的 WalkDir)
	walkErr := filepath.WalkDir(path, func(filePath string, dirEntry fs.DirEntry, err error) error {
		// 如果遍历目录遇到权限错误, 则忽略该错误并给出更友好的提示信息
		if err != nil {
			// 跳过不存在的文件或目录
			if os.IsNotExist(err) {
				return nil
			}

			// 如果遍历目录失败, 返回错误
			return fmt.Errorf("遍历目录失败: 路径 %s 错误: %v", filePath, err)
		}

		// 如果是指定的路径, 则跳过
		if filePath == path {
			return nil
		}

		// 如果是隐藏文件且未启用-H选项，则跳过
		if !sizeCmdHidden.Get() {
			if common.IsHidden(filePath) {
				return nil
			}
		}

		// 把文件信息发送到通道, 由worker处理
		if !dirEntry.IsDir() {
			// 获取文件信息
			info, err := dirEntry.Info()
			if err != nil {
				if len(errChan) < cap(errChan) {
					errChan <- fmt.Errorf("获取文件信息失败: 路径 %s 错误: %v", filePath, err)
				}
				return nil
			}
			fileChan <- fileInfo{path: filePath, size: info.Size()}
		}
		return nil
	})

	close(fileChan) // 关闭文件通道
	wg.Wait()       // 等待所有worker完成
	close(errChan)  // 关闭错误通道

	if walkErr != nil {
		// 如果遍历目录遇到权限错误, 则忽略该错误并给出更友好的提示信息
		if os.IsPermission(walkErr) {
			// 如果遍历目录失败, 返回错误
			return 0, fmt.Errorf("权限不足: 路径 %s", path)
		}

		// 检查是否为文件不存在错误
		if os.IsNotExist(walkErr) {
			// 如果是文件不存在错误, 则忽略该错误并给出更友好的提示信息
			return 0, fmt.Errorf("文件不存在: 路径 %s", path)
		}

		// 如果遍历目录失败, 返回错误
		return 0, fmt.Errorf("遍历目录失败: %v", walkErr)
	}

	// 检查并去重错误，最多显示5类错误
	errorMap := make(map[string]bool)
	errorCount := 0
	var otherErrors []error
	for err := range errChan {
		errStr := err.Error()
		if !errorMap[errStr] {
			errorMap[errStr] = true
			errorCount++
			if strings.Contains(errStr, "遍历目录失败") {
				// walkErr是目录遍历错误
				if walkErr == nil {
					walkErr = err
				} else {
					walkErr = fmt.Errorf("%v; %v", walkErr, err)
				}
			} else if errorCount <= 5 {
				// 其他类型错误
				otherErrors = append(otherErrors, err)
			}
		}
	}
	// 组合所有错误
	if len(otherErrors) > 0 {
		if walkErr != nil {
			walkErr = fmt.Errorf("%v; 其他错误: %v", walkErr, otherErrors)
		} else {
			walkErr = fmt.Errorf("其他错误: %v", otherErrors)
		}
	}
	if errorCount > 5 {
		walkErr = fmt.Errorf("%v (还有%d个类似错误...)", walkErr, errorCount-5)
	}

	// 获取最终的总大小
	finalSize := totalSize.Load()

	// 根据是否有错误返回不同内容
	if walkErr != nil {
		// 有错误时返回错误和计算得到的大小
		return finalSize, fmt.Errorf("计算完成但遇到错误: %v (总大小: %s)", walkErr, humanReadableSize(finalSize, 2))
	}
	// 无错误时只返回大小
	return finalSize, nil
}

// 预定义单位和对应的阈值
var (
	units = []string{"B", "KB", "MB", "GB", "TB", "PB"}
	// 预计算阈值，避免重复计算
	thresholds = []float64{
		1,                                // B
		1024,                             // KB
		1024 * 1024,                      // MB
		1024 * 1024 * 1024,               // GB
		1024 * 1024 * 1024 * 1024,        // TB
		1024 * 1024 * 1024 * 1024 * 1024, // PB
	}
)

// humanReadableSize 函数用于将字节大小转换为可读的字符串格式
//
// 参数:
//   - size: 需要转换的字节大小
//   - fn: 保留的小数位数
//
// 返回:
//   - 字节大小的文字表示
func humanReadableSize(size int64, fn int) string {
	// 处理0值情况 - 提前返回
	if size == 0 {
		return "0 B"
	}

	sizeFloat := float64(size)

	// 使用循环找到合适的单位，从大到小遍历
	unitIndex := 0
	for i := len(thresholds) - 1; i > 0; i-- {
		if sizeFloat >= thresholds[i] {
			unitIndex = i
			sizeFloat /= thresholds[i]
			break
		}
	}

	// 统一的格式化逻辑
	var formatted string
	if sizeFloat < 10 && unitIndex > 0 {
		// 小于10且不是字节时，使用指定小数位数（至少1位）
		decimals := fn
		if decimals < 1 {
			decimals = 1
		}
		formatted = fmt.Sprintf("%."+fmt.Sprintf("%d", decimals)+"f", sizeFloat)
		// 移除末尾的.0
		formatted = strings.TrimSuffix(formatted, ".0")
	} else {
		// 大于等于10或者是字节时，使用整数格式
		formatted = fmt.Sprintf("%.0f", sizeFloat)
	}

	return fmt.Sprintf("%s %s", formatted, units[unitIndex])
}

// 打印文件大小表格到控制台
func printSizeTable(its items, cl *colorlib.ColorLib) {
	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头, 只有在-ts为none时才设置表头
	if sizeCmdTableStyle.Get() != "none" {
		t.AppendHeader(table.Row{"Size", "Name"})
	}

	// 遍历items数组, 添加行
	for i := range its {
		// 大小列使用白色
		colorSize := cl.Swhite(its[i].Size)

		// 文件名列根据类型着色
		colorName := common.SprintStringColor(its[i].Name, its[i].Name, cl)

		// 添加行
		t.AppendRow(table.Row{colorSize, colorName})
	}

	// 设置列的对齐方式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Size", Align: text.AlignRight}, // 文件大小 - 右对齐
		{Name: "Name", Align: text.AlignLeft},  // 文件名 - 左对齐
	})

	// 根据-ts的值设置表格样式
	if style, ok := types.TableStyleMap[sizeCmdTableStyle.Get()]; ok {
		t.SetStyle(style)
	}

	// 输出表格
	t.Render()
}
