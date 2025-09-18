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

	// 新增：判断是否为多路径场景
	isMultiPath := shouldGroupByPath(paths, expandedPaths)

	// 4. 扫描文件
	scanner := NewFileScanner()
	files, err := scanner.ScanWithOriginalPaths(paths, expandedPaths, getScanOptions())
	if err != nil {
		return err
	}

	// 5. 处理数据 - 传递多路径标识
	processor := NewFileProcessor()
	processed := processor.Process(files, getProcessOptions(isMultiPath))

	// 6. 格式化输出
	formatter := NewFileFormatter(cl)
	return formatter.Render(processed, getFormatOptions(isMultiPath))
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
// 参数:
//   - isMultiPath: 是否为多路径场景
//
// 返回:
//   - ProcessOptions: 处理选项
func getProcessOptions(isMultiPath bool) ProcessOptions {
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
		SortBy:      sortBy,                   // 排序方式
		Reverse:     listCmdReverseSort.Get(), // 是否反向排序
		GroupByDir:  listCmdRecursion.Get(),   // 原有的递归分组
		GroupByPath: isMultiPath,              // 新增的路径分组
		IsMultiPath: isMultiPath,              // 多路径标识
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

	return nil
}

// shouldGroupByPath 判断是否应该按路径分组
// 优化版本：减少不必要的循环和提前返回
//
// 参数:
//   - originalPaths: 原始路径列表
//   - expandedPaths: 展开后的路径列表
//
// 返回:
//   - bool: 是否应该按路径分组
func shouldGroupByPath(originalPaths, expandedPaths []string) bool {
	// 如果已经是递归模式，不需要额外分组
	if listCmdRecursion.Get() {
		return false
	}

	// 情况1：用户明确指定了多个路径参数（不是通配符展开的）
	if len(originalPaths) > 1 {
		return true
	}

	// 情况2：单个原始路径，通配符展开的情况
	if len(originalPaths) == 1 && len(expandedPaths) > 1 {
		// 检查原始路径是否包含通配符
		hasWildcard := strings.ContainsAny(originalPaths[0], "*?[]")
		if hasWildcard {
			// 通配符展开的情况：检查是否包含目录
			// 如果包含目录，需要分组显示（文件在当前目录组，目录单独分组）
			for _, path := range expandedPaths {
				if info, err := os.Stat(path); err == nil && info.IsDir() {
					return true
				}
			}
		}
		return false
	}

	return false
}

// getFormatOptions 获取格式化选项
//
// 参数:
//   - isMultiPath: 是否为多路径场景
//
// 返回:
//   - FormatOptions: 格式化选项
func getFormatOptions(isMultiPath bool) FormatOptions {
	// 预先计算是否应该分组，避免在渲染时重复判断
	shouldGroup := listCmdRecursion.Get() || isMultiPath

	return FormatOptions{
		LongFormat:    listCmdLongFormat.Get(),    // 是否长格式
		UseColor:      listCmdColor.Get(),         // 是否使用颜色
		TableStyle:    listCmdTableStyle.Get(),    // 表格样式
		QuoteNames:    listCmdQuoteNames.Get(),    // 是否引用名称
		ShowUserGroup: listCmdShowUserGroup.Get(), // 是否显示用户和组
		ShouldGroup:   shouldGroup,                // 是否应该分组显示
	}
}
