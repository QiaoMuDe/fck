// Package size 实现了文件和目录大小计算功能。
// 该文件提供了计算文件、目录大小的核心功能，支持通配符路径展开、隐藏文件处理、进度显示和表格输出。
package size

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/schollz/progressbar/v3"
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
		targetPaths = []string{"*"} // 使用通配符匹配当前目录所有项目
	}

	// 根据sizeCmdColor设置颜色模式
	cl.SetColor(sizeCmdColor.Get())

	// 新建一个表格项列表
	var itemList items

	// 遍历路径
	for _, targetPath := range targetPaths {
		// 清理路径
		targetPath = filepath.Clean(targetPath)

		// 获取所有需要处理的路径
		pathsToProcess, err := expandPath(targetPath)
		if err != nil {
			cl.PrintErrorf("展开路径失败: %v\n", err)
			continue
		}

		// 处理每个路径
		for _, path := range pathsToProcess {
			addPathToList(path, &itemList, cl)
		}
	}

	// 统一打印所有结果
	if len(itemList) > 0 {
		printSizeTable(itemList, cl)
	}

	return nil
}

// expandPath 展开路径(处理通配符)
//
// 参数:
//   - path: 要展开的路径
//
// 返回:
//   - []string: 展开后的路径列表
//   - error: 如果发生错误，返回错误信息，否则返回 nil
func expandPath(path string) ([]string, error) {
	// 如果路径不包含通配符，则直接返回
	if !strings.Contains(path, "*") {
		return []string{path}, nil
	}

	// 使用 filepath.Glob 获取所有匹配的文件路径
	filePaths, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	// 检查是否有匹配的文件
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("没有找到匹配的文件: %s", path)
	}

	return filePaths, nil
}

// addPathToList 添加路径到列表(统一处理逻辑)
//
// 参数:
//   - path: 要添加的路径
//   - itemList: 存储路径的列表
//   - cl: 用于打印输出的 ColorLib 对象
func addPathToList(path string, itemList *items, cl *colorlib.ColorLib) {
	// 统一的隐藏文件检查
	if !sizeCmdHidden.Get() && common.IsHidden(path) {
		return
	}

	// 获取文件或目录大小
	size, err := getPathSize(path)
	if err != nil {
		cl.PrintErrorf("计算大小失败: %s - %v\n", path, err)
		return
	}

	// 添加到数组
	*itemList = append(*itemList, item{
		Name: path,
		Size: humanReadableSize(size, 2),
	})
}

// getPathSize 获取路径大小
//
// 参数:
//   - path: 要获取大小的路径
//
// 返回:
//   - int64: 路径大小
//   - error: 如果发生错误，返回错误信息，否则返回 nil
func getPathSize(path string) (int64, error) {
	// 统一的隐藏文件处理逻辑
	includeHidden := sizeCmdHidden.Get()
	if !includeHidden && common.IsHidden(path) {
		return 0, nil
	}

	// 获取文件信息
	info, err := os.Lstat(path)
	if err != nil {
		// 统一的错误处理
		switch {
		case os.IsPermission(err):
			return 0, fmt.Errorf("权限不足: 路径 %s", path)
		case os.IsNotExist(err):
			return 0, fmt.Errorf("文件不存在: 路径 %s", path)
		default:
			return 0, fmt.Errorf("获取文件信息失败: 路径 %s 错误: %v", path, err)
		}
	}

	// 如果不是目录, 则直接返回文件大小
	if !info.IsDir() {
		return info.Size(), nil
	}

	// 目录大小计算
	var totalSize int64
	var skippedFiles int

	// 创建进度条(只有目录才有进度条)
	bar := progressbar.NewOptions64(
		-1,                                // 总进度
		progressbar.OptionClearOnFinish(), // 完成后清除进度条
		progressbar.OptionSetDescription(filepath.Base(path)+" 计算中"), // 设置进度条描述
	)
	defer func() {
		// 完成进度条
		_ = bar.Finish()

		// 关闭进度条
		_ = bar.Close()
	}()

	// 遍历目录 (使用更高效的 WalkDir)
	walkErr := filepath.WalkDir(path, func(filePath string, dirEntry fs.DirEntry, err error) error {
		// 如果遍历目录遇到错误, 静默跳过并计数
		if err != nil {
			if !os.IsNotExist(err) {
				skippedFiles++
			}
			return nil
		}

		// 跳过根目录本身
		if filePath == path {
			return nil
		}

		// 统一的隐藏文件检查逻辑
		if !includeHidden && common.IsHidden(filePath) {
			return nil
		}

		// 只处理文件，累加文件大小
		if !dirEntry.IsDir() {
			fileInfo, infoErr := dirEntry.Info()
			if infoErr != nil {
				skippedFiles++
				return nil
			}

			fileSize := fileInfo.Size()
			totalSize += fileSize
			_ = bar.Add64(fileSize)
		}

		return nil
	})

	// 处理遍历错误
	if walkErr != nil {
		switch {
		case os.IsPermission(walkErr):
			return 0, fmt.Errorf("权限不足: 路径 %s", path)
		case os.IsNotExist(walkErr):
			return 0, fmt.Errorf("文件不存在: 路径 %s", path)
		default:
			return 0, fmt.Errorf("遍历目录失败: %v", walkErr)
		}
	}

	// 注意: 跳过的文件不作为错误处理，这是正常的操作行为
	// 如果需要了解跳过的文件数量，可以通过日志或其他方式输出
	return totalSize, nil
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
