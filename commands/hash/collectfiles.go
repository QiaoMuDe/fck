package hash

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
)

// collectFiles 函数用于收集指定路径下的所有文件
//
// 参数:
//   - targetPath: 要收集文件的路径
//   - recursive: 是否递归处理目录
//   - cl: ColorLib 实例，用于彩色输出
//
// 返回值：
//   - []string: 收集到的文件路径列表
//   - error: 如果发生错误，则返回错误信息
func collectFiles(targetPath string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	// 检查路径是否包含通配符
	if strings.ContainsAny(targetPath, "*?[]{}") {
		return collectGlobFiles(targetPath, recursive, cl)
	}

	return collectSinglePath(targetPath, recursive, cl)
}

// collectGlobFiles 处理包含通配符的路径
//
//	参数:
//	 - pattern: 包含通配符的路径模式
//	 - recursive: 是否递归处理目录
//	 - cl: ColorLib 实例，用于彩色输出
//
//	返回值：
//	 - []string: 收集到的文件路径列表
//	 - error: 如果发生错误，则返回错误信息
func collectGlobFiles(pattern string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	matchedPaths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("路径无效: %w", err)
	}

	if len(matchedPaths) == 0 {
		return nil, fmt.Errorf("没有找到匹配的文件")
	}

	// 预估容量
	var allFiles []string

	for _, path := range matchedPaths {
		if shouldSkipHidden(path) {
			continue
		}

		pathFiles, err := collectSinglePath(path, recursive, cl)
		if err != nil {
			// 对于通配符匹配，如果是目录相关的错误，只打印警告而不中断整个过程
			if isDirectorySkipError(err) {
				cl.PrintWarn(err.Error())
				continue
			}
			// 其他错误仍然返回
			return nil, err
		}

		allFiles = append(allFiles, pathFiles...)
	}

	return allFiles, nil
}

// collectSinglePath 处理单个路径(文件或目录)
//
// 参数:
//   - targetPath: 要收集文件的路径
//   - recursive: 是否递归处理目录
//   - cl: ColorLib 实例，用于彩色输出
//
// 返回:
//   - []item: 包含文件信息的切片
//   - error: 错误信息，如果发生错误则返回非nil值
func collectSinglePath(targetPath string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, wrapStatError(err, targetPath)
	}

	if shouldSkipHidden(targetPath) {
		return nil, fmt.Errorf("跳过隐藏项: %s", targetPath)
	}

	if info.IsDir() {
		return handleDirectory(targetPath, recursive, cl)
	}

	// 普通文件
	return []string{targetPath}, nil
}

// handleDirectory 处理目录
//
// 参数:
//   - dirPath: 要处理的目录路径
//   - recursive: 是否递归处理目录
//   - cl: ColorLib 实例，用于彩色输出
//
// 返回:
//   - []string: 包含文件路径的切片
//   - error: 错误信息，如果发生错误则返回非nil值
func handleDirectory(dirPath string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	if !recursive {
		return nil, fmt.Errorf("跳过目录: %s 请使用 -r 选项以递归方式处理", dirPath)
	}

	return walkDir(dirPath, recursive, cl)
}

// shouldSkipHidden 检查是否应该跳过隐藏文件/目录
//
// 参数:
//   - path: 要检查的路径
//
// 返回:
//   - bool: 如果应该跳过隐藏项，则返回true；否则返回false
func shouldSkipHidden(path string) bool {
	return !hashCmdHidden.Get() && common.IsHidden(path)
}

// wrapStatError 统一处理 os.Stat 错误
//
// 参数:
//   - err: 错误对象
//   - path: 错误路径
//
// 返回:
//   - error: 处理后的错误对象
func wrapStatError(err error, path string) error {
	if os.IsPermission(err) {
		return fmt.Errorf("权限不足: %s", path)
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", path)
	}
	return fmt.Errorf("无法获取文件信息: %w", err)
}

// walkDir 函数用于根据递归标志遍历指定目录并收集文件列表
//
// 参数:
//   - dirPath: 要遍历的目录路径
//   - recursive: 是否递归遍历子目录
//   - cl: ColorLib 实例，用于彩色输出
//
// 返回:
//   - []string: 包含文件路径的切片
//   - error: 错误信息，如果发生错误则返回非nil值
func walkDir(dirPath string, recursive bool, cl *colorlib.ColorLib) ([]string, error) {
	var files []string

	if !recursive {
		return walkDirNonRecursive(dirPath, cl)
	}

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if shouldSkipHidden(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, wrapWalkError(err, dirPath)
	}

	return files, nil
}

// walkDirNonRecursive 非递归遍历目录
//
// 参数:
//   - dirPath: 要遍历的目录路径
//   - cl: ColorLib 实例，用于彩色输出
//
// 返回:
//   - []string: 包含文件路径的切片
//   - error: 错误信息，如果发生错误则返回非nil值
func walkDirNonRecursive(dirPath string, cl *colorlib.ColorLib) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if shouldSkipHidden(entry.Name()) {
			continue
		}

		if entry.IsDir() {
			cl.PrintWarnf("跳过目录: %s 请使用 -r 选项以递归方式处理\n", entry.Name())
			continue
		}

		files = append(files, filepath.Join(dirPath, entry.Name()))
	}

	return files, nil
}

// wrapWalkError 统一处理遍历错误
//
// 参数:
//   - err: 错误对象
//   - dirPath: 错误路径
//
// 返回值:
//   - error: 处理后的错误对象
func wrapWalkError(err error, dirPath string) error {
	if os.IsPermission(err) {
		return fmt.Errorf("权限不足: %s", dirPath)
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", dirPath)
	}
	return fmt.Errorf("遍历目录失败: %w", err)
}

// isDirectorySkipError 检查是否是目录跳过错误
//
// 参数:
//   - err: 要检查的错误
//
// 返回值:
//   - bool: 如果是目录跳过错误则返回true，否则返回false
func isDirectorySkipError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "跳过目录") || strings.Contains(errStr, "跳过隐藏项")
}
