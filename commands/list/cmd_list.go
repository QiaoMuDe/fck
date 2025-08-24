// Package list 实现了文件列表显示命令的主要逻辑。
// 该文件包含 list 子命令的入口函数，负责参数验证、路径处理、文件扫描和格式化输出。
package list

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// ListCmdMain list 命令主函数
//
// 参数:
//   - cl: 颜色库
//
// 返回:
//   - error: 错误信息
func ListCmdMain(cl *colorlib.ColorLib) error {
	// 1. 参数验证
	if err := validateArgs(); err != nil {
		return err
	}

	// 2. 配置颜色输出
	cl.SetColor(listCmdColor.Get())

	// 3. 获取和处理路径
	paths := getPaths()
	expandedPaths, err := expandPaths(paths, cl)
	if err != nil {
		return err
	}

	// 4. 扫描文件
	scanner := NewFileScanner()
	files, err := scanner.Scan(expandedPaths, getScanOptions())
	if err != nil {
		return err
	}

	// 5. 处理数据
	processor := NewFileProcessor()
	processed := processor.Process(files, getProcessOptions())

	// 6. 格式化输出
	formatter := NewFileFormatter(cl)
	return formatter.Render(processed, getFormatOptions())
}

// getPaths 获取路径列表
//
// 返回:
//   - []string: 路径列表
func getPaths() []string {
	paths := listCmd.Args() // 获取命令行参数

	// 如果没有指定路径, 则默认为当前目录
	if len(paths) == 0 {
		if runtime.GOOS == "windows" {
			if dir, err := os.Getwd(); err == nil {
				return []string{dir}
			}
		}
		return []string{"."}
	}

	return paths
}

// expandPaths 展开通配符路径
//
// 参数:
//   - paths: 路径列表
//   - cl: 颜色库
//
// 返回:
//   - []string: 展开后的路径列表
//   - error: 错误信息
func expandPaths(paths []string, cl *colorlib.ColorLib) ([]string, error) {
	var expandedPaths []string // 展开后的路径列表

	// 扫描匹配的路径
	for _, path := range paths {
		// 清理路径
		path = filepath.Clean(path)

		// 判断是否为通配符路径
		isWildcardPath := strings.ContainsAny(path, "*?[]")

		if isWildcardPath {
			// 处理通配符路径
			matches, err := filepath.Glob(path)
			if err != nil {
				cl.PrintErrorf("路径模式错误 %q: %v\n", path, err)
				continue
			}

			// 如果路径模式没有匹配任何文件，则打印错误信息
			if len(matches) == 0 {
				cl.PrintWarnf("通配符路径未匹配到任何文件: %s\n", path)
				continue
			}

			// 过滤隐藏文件: 默认不显示隐藏文件
			for _, match := range matches {
				if listCmdAll.Get() || !common.IsHidden(match) {
					expandedPaths = append(expandedPaths, match)
				}
			}

			continue
		}

		// 处理普通路径
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				cl.PrintWarnf("路径不存在: %s\n", path)
			} else {
				cl.PrintErrorf("无法访问路径 %q: %v\n", path, err)
			}
			continue
		}

		// 检查是否为隐藏文件: 默认不显示隐藏文件
		if listCmdAll.Get() || !common.IsHidden(path) {
			expandedPaths = append(expandedPaths, path)
		}
	}

	// 删除重复的路径
	seen := make(map[string]bool)
	uniquePaths := make([]string, 0, len(expandedPaths)) // 去重后的路径列表
	for _, p := range expandedPaths {
		if !seen[p] {
			seen[p] = true
			uniquePaths = append(uniquePaths, p)
		}
	}

	return uniquePaths, nil
}

// getScanOptions 获取扫描选项
//
// 返回:
//   - ScanOptions: 扫描选项
func getScanOptions() ScanOptions {
	var fileTypes []string // 文件类型

	// 获取限制类型
	if listCmdType.Get() != types.FindTypeAll {
		fileTypes = []string{listCmdType.Get()}
	}

	return ScanOptions{
		Recursive:  listCmdRecursion.Get(), // 是否递归扫描
		ShowHidden: listCmdAll.Get(),       // 显示隐藏文件
		FileTypes:  fileTypes,              // 文件类型
		DirItself:  listCmdDirItself.Get(), // 是否包括当前目录
	}
}

// getProcessOptions 获取处理选项
//
// 返回:
//   - ProcessOptions: 处理选项
func getProcessOptions() ProcessOptions {
	var sortBy string // 排序方式

	if listCmdSortByTime.Get() {
		sortBy = "time" // 按照时间排序
	} else if listCmdSortBySize.Get() {
		sortBy = "size" // 按照大小排序
	} else if listCmdSortByName.Get() {
		sortBy = "name" // 按照名称排序
	} else {
		sortBy = "name" // 默认按名称排序
	}

	return ProcessOptions{
		SortBy:     sortBy,                   // 排序方式
		Reverse:    listCmdReverseSort.Get(), // 是否反向排序
		GroupByDir: listCmdRecursion.Get(),   // 是否按目录分组
	}
}

// validateArgs 验证参数
//
// 返回:
//   - error: 错误
func validateArgs() error {
	// 检查是否同时指定了 -s 和 -t 选项
	if listCmdSortBySize.Get() && listCmdSortByTime.Get() {
		return errors.New("不能同时指定 -s 和 -t 选项")
	}

	// 检查是否同时指定了 -s 和 -n 选项
	if listCmdSortBySize.Get() && listCmdSortByName.Get() {
		return errors.New("不能同时指定 -s 和 -n 选项")
	}

	// 检查是否同时指定了 -t 和 -n 选项
	if listCmdSortByTime.Get() && listCmdSortByName.Get() {
		return errors.New("不能同时指定 -t 和 -n 选项")
	}

	// 如果指定了-ho检查是否指定-a
	if (listCmdType.Get() == types.FindTypeHiddenShort || listCmdType.Get() == types.FindTypeHidden) && !listCmdAll.Get() {
		return fmt.Errorf("必须指定 %s 选项才能使用 %s 选项", listCmdAll.Name(), listCmdType.Name())
	}

	// 检查是否同时指定了 -c 和 --dev-color
	if listCmdDevColor.Get() && !listCmdColor.Get() {
		return fmt.Errorf("如果要使用 -%s, 必须要先启用 -%s", listCmdDevColor.ShortName(), listCmdColor.ShortName())
	}

	return nil
}

// getFormatOptions 获取格式化选项
//
// 返回:
//   - FormatOptions: 格式化选项
func getFormatOptions() FormatOptions {
	return FormatOptions{
		LongFormat:    listCmdLongFormat.Get(),    // 是否长格式
		UseColor:      listCmdColor.Get(),         // 是否使用颜色
		DevColor:      listCmdDevColor.Get(),      // 是否使用开发者颜色
		TableStyle:    listCmdTableStyle.Get(),    // 表格样式
		QuoteNames:    listCmdQuoteNames.Get(),    // 是否引用名称
		ShowUserGroup: listCmdShowUserGroup.Get(), // 是否显示用户和组
	}
}
