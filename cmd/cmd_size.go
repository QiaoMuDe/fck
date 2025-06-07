package cmd

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
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

// sizeCmdMain 是 size 子命令的主函数
func sizeCmdMain(sizeCmd *flag.FlagSet, cl *colorlib.ColorLib) error {
	// 获取指定的路径
	targetPaths := sizeCmd.Args()

	// 如果没有指定路径, 则默认计算当前目录下每个项目的大小
	if len(targetPaths) == 0 {
		targetPaths = []string{"*"}
	}

	// 检查表格的样式是否有效
	if *sizeCmdTableStyle != "" {
		if _, ok := globals.TableStyleMap[*sizeCmdTableStyle]; !ok {
			return fmt.Errorf("无效的表格样式: %s", *sizeCmdTableStyle)
		}
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
				size, err := getPathSize(filePath)
				if err != nil {
					cl.PrintErrf("计算大小时出现错误: %s\n", err)
				}

				// 如果没启用表格输出, 则直接打印结果
				if *sizeCmdTableStyle == "" {
					if *sizeCmdColor {
						if err := printStringColor(filePath, fmt.Sprintf("%-15s\t%s", humanReadableSize(size, 2), filePath), cl); err != nil {
							cl.PrintErrf("输出路径时出错: %s\n", err)
							continue
						}
					} else {
						fmt.Printf("%-15s\t%s\n", humanReadableSize(size, 2), filePath)
						continue
					}
				} else {
					// 添加到 items 数组中
					itemList = append(itemList, item{
						Name: filePath,
						Size: humanReadableSize(size, 2),
					})
				}
			}

			// 打印输出
			if *sizeCmdTableStyle != "" {
				printSizeTable(itemList, cl)
			}
			return nil
		}

		// 如果是文件, 则直接计算大小
		if info, err := os.Stat(targetPath); err != nil {
			cl.PrintErrf("获取文件信息失败: 路径 %s 错误: %v\n", targetPath, err)
			continue
		} else if !info.IsDir() {
			// 如果没启用表格输出, 则直接打印结果
			if *sizeCmdTableStyle == "" {
				// 根据是否启用颜色打印结果
				if *sizeCmdColor {
					if err := printStringColor(targetPath, fmt.Sprintf("%-15s\t%s", humanReadableSize(info.Size(), 2), targetPath), cl); err != nil {
						cl.PrintErrf("输出路径时出错: %s\n", err)
						continue
					}
				} else {
					fmt.Printf("%-15s\t%s\n", humanReadableSize(info.Size(), 2), targetPath)
					continue
				}
			} else {
				// 添加到 items 数组中
				itemList = append(itemList, item{
					Name: targetPath,
					Size: humanReadableSize(info.Size(), 2),
				})
				// 打印输出
				printSizeTable(itemList, cl)
				continue
			}
		}

		// 如果是目录, 则递归计算大小
		size, err := getPathSize(targetPath)
		if err != nil {
			cl.PrintErrf("计算目录大小失败: 路径 %s 错误: %v\n", targetPath, err)
			continue
		}

		// 如果没启用表格输出, 则直接打印结果
		if *sizeCmdTableStyle == "" {
			// 根据是否启用颜色打印结果
			if *sizeCmdColor {
				if err := printStringColor(targetPath, fmt.Sprintf("%-15s\t%s", humanReadableSize(size, 2), targetPath), cl); err != nil {
					cl.PrintErrf("输出路径时出错: %s\n", err)
					continue
				}
			} else {
				fmt.Printf("%-15s\t%s\n", humanReadableSize(size, 2), targetPath)
				continue
			}
		} else {
			// 添加到 items 数组中
			itemList = append(itemList, item{
				Name: targetPath,
				Size: humanReadableSize(size, 2),
			})

			// 打印输出
			printSizeTable(itemList, cl)
			continue
		}
	}

	return nil
}

// getPathSize 获取路径大小
func getPathSize(path string) (int64, error) {
	// 获取文件信息
	info, statErr := os.Stat(path)
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

	// 定义总大小
	var totalSize int64

	// 创建并发任务池
	var wg sync.WaitGroup
	var mu sync.Mutex
	fileChan := make(chan string, 10000) // 缓冲区大小为10000
	errChan := make(chan error, 100)     // 错误通道大小为100

	// 启动worker goroutines
	workers := *sizeCmdJob
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
			for filePath := range fileChan {
				fileInfo, err := os.Lstat(filePath)
				if err != nil {
					if len(errChan) < cap(errChan) {
						errChan <- fmt.Errorf("获取文件信息失败: 路径 %s 错误: %v", filePath, err)
					}
					continue
				}

				mu.Lock()
				totalSize += fileInfo.Size()
				mu.Unlock()
			}
		}()
	}

	// 遍历目录
	walkErr := filepath.Walk(path, func(filePath string, fileInfo fs.FileInfo, err error) error {
		// 如果遍历目录遇到权限错误, 则忽略该错误并给出更友好的提示信息
		if err != nil {
			// 跳过不存在的文件或目录
			if os.IsNotExist(err) {
				return nil
			}

			// 如果遍历目录失败, 返回错误
			return fmt.Errorf("遍历目录失败: 路径 %s 错误: %v", filePath, err)
		}

		// 如果是文件或者不是根目录, 则发送到worker处理
		if !fileInfo.IsDir() || (fileInfo.IsDir() && filePath != path) {
			fileChan <- filePath
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

	// 根据是否有错误返回不同内容
	if walkErr != nil {
		// 有错误时返回错误和计算得到的大小
		return totalSize, fmt.Errorf("计算完成但遇到错误: %v (总大小: %s)", walkErr, humanReadableSize(totalSize, 2))
	}
	// 无错误时只返回大小
	return totalSize, nil
}

// humanReadableSize 函数用于将字节大小转换为可读的字符串格式
// 该函数接收一个 int64 类型的字节大小参数, 返回一个表示该大小的可读字符串
func humanReadableSize(size int64, fn int) string {
	// 定义存储字节单位的切片, 按照从小到大的顺序排列
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	// 定义字节单位之间的换算基数, 这里使用 1024 作为二进制换算标准
	base := float64(1024)

	// 用于存储最终选择的合适单位
	var unit string
	// 将传入的 int64 类型的字节大小转换为 float64 类型, 方便后续计算
	sizeFloat := float64(size)

	// 根据字节大小选择最合适的单位
	// 如果字节大小小于 1024B, 则直接使用 B 作为单位
	if sizeFloat < base {
		unit = units[0]
		// 如果字节大小小于 1024KB, 则使用 KB 作为单位, 并将字节大小除以 1024 转换为 KB
	} else if sizeFloat < base*base {
		unit = units[1]
		sizeFloat /= base
		// 如果字节大小小于 1024MB, 则使用 MB 作为单位, 并将字节大小除以 1024*1024 转换为 MB
	} else if sizeFloat < base*base*base {
		unit = units[2]
		sizeFloat /= base * base
		// 如果字节大小小于 1024GB, 则使用 GB 作为单位, 并将字节大小除以 1024*1024*1024 转换为 GB
	} else if sizeFloat < base*base*base*base {
		unit = units[3]
		sizeFloat /= base * base * base
		// 如果字节大小小于 1024TB, 则使用 TB 作为单位, 并将字节大小除以 1024*1024*1024*1024 转换为 TB
	} else if sizeFloat < base*base*base*base*base {
		unit = units[4]
		sizeFloat /= base * base * base * base
		// 否则使用 PB 作为单位, 并将字节大小除以 1024*1024*1024*1024*1024 转换为 PB
	} else {
		unit = units[5]
		sizeFloat /= base * base * base * base * base
	}

	// 构造格式化字符串
	format := fmt.Sprintf("%%.%df", fn)

	// 使用构造的格式化字符串进行格式化
	sizeF := fmt.Sprintf(format, sizeFloat)

	// 如果转换后的大小为 0, 则返回 "0B"
	if sizeF == "0.00" {
		return "0 B"
	}

	// 去除小数部分末尾的 .00 或 .0
	sizeF = strings.TrimSuffix(sizeF, ".00")
	sizeF = strings.TrimSuffix(sizeF, ".0")

	// 去除小数点部分末尾的0
	if strings.Contains(sizeF, ".") {
		sizeF = strings.TrimRight(sizeF, "0")
	}

	// 先将转换后的大小和单位拼接成一个字符串
	result := fmt.Sprintf("%s %s", sizeF, unit)

	return result
}

// 打印文件大小表格到控制台
func printSizeTable(its items, cl *colorlib.ColorLib) {
	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头, 只有在-ts为none时才设置表头
	if *sizeCmdTableStyle != "" && *sizeCmdTableStyle != "none" {
		t.AppendHeader(table.Row{"Size", "Name"})
	}

	for i := range its {
		// 添加行
		if *sizeCmdColor {
			colorSize, sizeErr := SprintStringColor(its[i].Name, its[i].Size, cl)
			if sizeErr != nil {
				colorSize = its[i].Size
			}

			colorName, nameErr := SprintStringColor(its[i].Name, its[i].Name, cl)
			if nameErr != nil {
				colorName = its[i].Name
			}

			t.AppendRow(table.Row{colorSize, colorName})

		} else {
			t.AppendRow(table.Row{its[i].Size, its[i].Name})
		}
	}

	// 设置列的对齐方式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Size", Align: text.AlignRight}, // 文件大小 - 右对齐
		{Name: "Name", Align: text.AlignLeft},  // 文件名 - 左对齐
	})

	// 设置表格样式
	if *sizeCmdTableStyle != "" {
		// 根据-ts的值设置表格样式
		if style, ok := globals.TableStyleMap[*sizeCmdTableStyle]; ok {
			t.SetStyle(style)
		}
	}

	// 输出表格
	t.Render()
}
