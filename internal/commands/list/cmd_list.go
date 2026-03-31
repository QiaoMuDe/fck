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
	"gitee.com/MM-Q/fck/internal/types"
	common "gitee.com/MM-Q/fck/internal/utils"
)

// ListConfig 列表配置结构体
type ListConfig struct {
	Args          []string // 命令行参数
	All           bool     // 显示所有文件和目录
	Color         bool     // 启用颜色输出
	SortByTime    bool     // 按修改时间排序
	SortBySize    bool     // 按文件大小排序
	SortByName    bool     // 按文件名排序
	DirItself     bool     // 列出目录本身
	LongFormat    bool     // 长格式显示
	ReverseSort   bool     // 反向排序
	QuoteNames    bool     // 引用名称
	Recursion     bool     // 递归列出目录
	ShowUserGroup bool     // 显示用户和组
	TableStyle    string   // 表格样式
	Type          string   // 文件类型
	Icon          bool     // 显示文件图标
	DisableIndex  bool     // 禁用索引
}

// ListCmdMain list 命令主函数
//
// 参数:
//   - cl: 颜色库
//   - config: 列表配置
//
// 返回:
//   - error: 错误信息
func ListCmdMain(cl *colorlib.ColorLib, config ListConfig) error {
	if err := validateArgs(config); err != nil {
		return err
	}

	cl.SetColor(config.Color)

	paths := getPaths(config.Args)
	expandedPaths, err := expandPaths(paths, cl, config)
	if err != nil {
		return err
	}

	isMultiPath := shouldGroupByPath(paths, config)

	scanner := NewFileScanner()
	files, err := scanner.ScanWithOriginalPaths(paths, expandedPaths, getScanOptions(config))
	if err != nil {
		return err
	}

	processor := NewFileProcessor()
	processed := processor.Process(files, getProcessOptions(isMultiPath, config))

	formatter := NewFileFormatter(cl)
	return formatter.Render(processed, getFormatOptions(isMultiPath, config))
}

// getPaths 获取路径列表
//
// 参数:
//   - args: 命令行参数
//
// 返回:
//   - []string: 路径列表
func getPaths(args []string) []string {
	paths := args

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

// expandPaths 展开路径列表（处理通配符）
//
// 参数:
//   - paths: 原始路径列表
//   - cl: 颜色库
//   - config: 列表配置
//
// 返回:
//   - []string: 展开后的路径列表
//   - error: 错误信息
func expandPaths(paths []string, cl *colorlib.ColorLib, config ListConfig) ([]string, error) {
	var expandedPaths []string

	for _, path := range paths {
		path = filepath.Clean(path)

		isWildcardPath := strings.ContainsAny(path, "*?[]")

		if isWildcardPath {
			matches, err := filepath.Glob(path)
			if err != nil {
				cl.PrintErrorf("路径模式错误 %q: %v\n", path, err)
				continue
			}

			if len(matches) == 0 {
				cl.PrintWarnf("通配符路径未匹配到任何文件: %s\n", path)
				continue
			}

			for _, match := range matches {
				if config.All || !common.IsHidden(match) {
					expandedPaths = append(expandedPaths, match)
				}
			}

			continue
		}

		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				cl.PrintWarnf("路径不存在: %s\n", path)
			} else {
				cl.PrintErrorf("无法访问路径 %q: %v\n", path, err)
			}
			continue
		}

		if config.All || !common.IsHidden(path) {
			expandedPaths = append(expandedPaths, path)
		}
	}

	seen := make(map[string]bool)
	uniquePaths := make([]string, 0, len(expandedPaths))
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
// 参数:
//   - config: 列表配置
//
// 返回:
//   - ScanOptions: 扫描选项
func getScanOptions(config ListConfig) ScanOptions {
	var fileTypes []string

	if config.Type != types.FindTypeAll {
		fileTypes = []string{config.Type}
	}

	return ScanOptions{
		Recursive:  config.Recursion,
		ShowHidden: config.All,
		FileTypes:  fileTypes,
		DirItself:  config.DirItself,
	}
}

// getProcessOptions 获取处理选项
//
// 参数:
//   - isMultiPath: 是否为多路径场景
//   - config: 列表配置
//
// 返回:
//   - ProcessOptions: 处理选项
func getProcessOptions(isMultiPath bool, config ListConfig) ProcessOptions {
	var sortBy string

	if config.SortByTime {
		sortBy = "time"
	} else if config.SortBySize {
		sortBy = "size"
	} else if config.SortByName {
		sortBy = "name"
	} else {
		sortBy = "name"
	}

	return ProcessOptions{
		SortBy:      sortBy,
		Reverse:     config.ReverseSort,
		GroupByDir:  config.Recursion,
		GroupByPath: isMultiPath,
		IsMultiPath: isMultiPath,
	}
}

// validateArgs 验证参数
//
// 参数:
//   - config: 列表配置
//
// 返回:
//   - error: 错误信息
func validateArgs(config ListConfig) error {
	if config.SortBySize && config.SortByTime {
		return errors.New("不能同时指定 -s 和 -t 选项")
	}

	if config.SortBySize && config.SortByName {
		return errors.New("不能同时指定 -s 和 -n 选项")
	}

	if config.SortByTime && config.SortByName {
		return errors.New("不能同时指定 -t 和 -n 选项")
	}

	if (config.Type == types.FindTypeHiddenShort || config.Type == types.FindTypeHidden) && !config.All {
		return fmt.Errorf("必须指定 -a 选项才能使用 -type hidden 选项")
	}

	return nil
}

// shouldGroupByPath 判断是否应该按路径分组
//
// 参数:
//   - originalPaths: 原始路径列表
//   - config: 列表配置
//
// 返回:
//   - bool: 是否应该按路径分组
func shouldGroupByPath(originalPaths []string, config ListConfig) bool {
	if config.Recursion {
		return false
	}

	// 多个原始路径（如: ls dir1 dir2）才需要分组
	// 通配符展开的路径不需要分组，直接列出条目即可
	if len(originalPaths) > 1 {
		return true
	}

	return false
}

// getFormatOptions 获取格式化选项
//
// 参数:
//   - isMultiPath: 是否为多路径场景
//   - config: 列表配置
//
// 返回:
//   - FormatOptions: 格式化选项
func getFormatOptions(isMultiPath bool, config ListConfig) FormatOptions {
	shouldGroup := config.Recursion || isMultiPath

	return FormatOptions{
		LongFormat:    config.LongFormat,
		UseColor:      config.Color,
		TableStyle:    config.TableStyle,
		QuoteNames:    config.QuoteNames,
		ShowUserGroup: config.ShowUserGroup,
		ShouldGroup:   shouldGroup,
		DisableIndex:  config.DisableIndex,
		Recursion:     config.Recursion,
		Icon:          config.Icon,
	}
}
