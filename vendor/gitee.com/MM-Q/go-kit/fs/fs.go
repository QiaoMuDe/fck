package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetDefaultBinPath 返回默认bin路径
// 用于获取Go程序的默认bin路径，采用多级回退策略确保总能返回有效路径
//
// 返回:
//   - string: 默认bin路径，优先级为GOPATH/bin > 用户主目录/go/bin > 当前工作目录/bin
func GetDefaultBinPath() string {
	// 1. 优先使用GOPATH/bin
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return filepath.Join(gopath, "bin")
	}

	// 2. 尝试获取用户主目录/go/bin
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, "go", "bin")
	}

	// 3. 使用当前工作目录/bin（保底策略）
	if currentDir, err := os.Getwd(); err == nil {
		return filepath.Join(currentDir, "bin")
	}

	// 所有获取失败时返回相对路径（理论上不会执行到此处）
	return filepath.Join(".", "bin")
}

// GetUserHomeDir 获取用户家目录
// 用于获取用户家目录路径，提供多级降级策略确保总能返回有效路径
//
// 返回:
//   - string: 用户家目录路径，失败时依次降级为工作目录或当前目录
func GetUserHomeDir() string {
	// 尝试获取用户家目录
	homeDir, err := os.UserHomeDir()

	// 先判断是否成功获取家目录
	if err == nil {
		// 成功获取时，确保返回绝对路径
		absHome, absErr := filepath.Abs(homeDir)
		if absErr == nil {
			return absHome
		}
		// 如果转换绝对路径失败，直接返回原始家目录
		return homeDir
	}

	// 家目录获取失败，尝试获取当前工作目录
	wd, wdErr := os.Getwd()
	if wdErr == nil {
		return wd
	}

	// 所有路径获取都失败时，返回当前目录"."作为最后的保底
	return "."
}

// GetExecutablePath 获取程序的绝对安装路径
// 用于获取当前可执行文件的绝对路径，提供多级降级策略确保总能返回路径
//
// 返回:
//   - string: 程序的绝对路径，失败时降级为相对路径
func GetExecutablePath() string {
	// 尝试使用 os.Executable 获取可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		// 如果 os.Executable 报错，使用 os.Args[0] 作为替代
		exePath = os.Args[0]
	}
	// 使用 filepath.Abs 确保路径是绝对路径
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		// 如果 filepath.Abs 报错，直接返回原始路径
		return exePath
	}
	return absPath
}

// walkDir 遍历目录并收集文件列表
// 用于根据递归标志遍历指定目录，收集所有文件路径
//
// 参数:
//   - dirPath: 要遍历的目录路径
//   - recursive: 是否递归遍历子目录
//
// 返回:
//   - []string: 收集到的文件路径切片
//   - error: 遍历失败时返回错误
func walkDir(dirPath string, recursive bool) ([]string, error) {
	// 快速失败：检查路径是否为空
	if dirPath == "" {
		return nil, fmt.Errorf("directory path cannot be empty")
	}

	var files []string

	// 递归模式
	if recursive {
		err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
			// 快速失败：遇到错误立即返回
			if err != nil {
				return err
			}
			// 快速跳过：只处理文件
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		// 快速返回错误
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory %q: %w", dirPath, err)
		}
		return files, nil
	}

	// 非递归模式
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %q: %w", dirPath, err)
	}

	// 快速返回：如果目录为空
	if len(entries) == 0 {
		return []string{}, nil
	}

	// 预分配容量以提高性能
	files = make([]string, 0, len(entries))
	for _, entry := range entries {
		// 快速跳过：只处理文件
		if !entry.IsDir() {
			files = append(files, filepath.Join(dirPath, entry.Name()))
		}
	}

	return files, nil
}

// Collect 收集指定路径下的所有文件
// 用于收集文件或目录中的文件，支持通配符匹配和递归遍历
//
// 参数:
//   - targetPath: 目标路径，支持通配符(*?[]{})
//   - recursive: 是否递归遍历目录
//
// 返回:
//   - []string: 收集到的文件路径切片
//   - error: 收集失败时返回错误
func Collect(targetPath string, recursive bool) ([]string, error) {
	// 快速失败：检查路径是否为空
	if targetPath == "" {
		return nil, fmt.Errorf("target path cannot be empty")
	}

	// 快速路由：根据是否包含通配符选择处理方式
	if strings.ContainsAny(targetPath, "*?[]{}") {
		return collectGlobFiles(targetPath, recursive)
	}

	return collectSinglePath(targetPath, recursive)
}

// collectGlobFiles 处理包含通配符的路径模式并收集匹配的文件
// 使用filepath.Glob匹配通配符模式，然后收集所有匹配路径中的文件
//
// 参数:
//   - pattern: 包含通配符的路径模式（如 "*.go", "dir/*", "**/*.txt"）
//   - recursive: 当匹配到目录时，是否递归遍历子目录
//
// 返回:
//   - []string: 所有匹配文件的路径切片
//   - error: 模式无效、无匹配文件或处理过程中的错误
func collectGlobFiles(pattern string, recursive bool) ([]string, error) {
	// 快速失败：检查模式是否为空
	if pattern == "" {
		return nil, fmt.Errorf("glob pattern cannot be empty")
	}

	// 匹配通配符
	matchedFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid path pattern %q: %w", pattern, err)
	}

	// 快速返回：没有匹配的文件
	if len(matchedFiles) == 0 {
		return nil, fmt.Errorf("no matching files found for pattern %q", pattern)
	}

	var files []string
	for _, file := range matchedFiles {
		pathFiles, err := collectSinglePath(file, recursive)
		// 快速失败：遇到错误立即返回
		if err != nil {
			return nil, err
		}
		files = append(files, pathFiles...)
	}

	return files, nil
}

// collectSinglePath 处理单个具体路径，可以是文件或目录
// 如果是文件则直接返回该文件路径，如果是目录则调用walkDir遍历
//
// 参数:
//   - path: 要处理的具体路径（不包含通配符）
//   - recursive: 当路径为目录时，是否递归遍历子目录
//
// 返回:
//   - []string: 文件路径切片，单个文件返回包含该文件的切片，目录返回其中所有文件
//   - error: 路径不存在、无权限访问或遍历过程中的错误
func collectSinglePath(path string, recursive bool) ([]string, error) {
	// 快速失败：检查路径是否为空
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get path info for %q: %w", path, err)
	}

	// 快速返回：如果是文件，直接返回
	if !info.IsDir() {
		return []string{path}, nil
	}

	// 如果是目录，进行遍历
	return walkDir(path, recursive)
}

// wrapPathError 包装路径相关错误，提供统一的错误处理
//
// 参数:
//   - err: 原始错误
//   - path: 路径
//   - operation: 操作描述
//
// 返回:
//   - error: 包装后的错误
func wrapPathError(err error, path, operation string) error {
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist when %s: %s", operation, path)
	}
	if os.IsPermission(err) {
		return fmt.Errorf("permission denied when %s path '%s': %w", operation, path, err)
	}
	return fmt.Errorf("error when %s path '%s': %w", operation, path, err)
}

// GetSize 获取文件或目录的大小
// 用于计算文件或目录的总字节数，目录会递归计算所有普通文件的大小
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - int64: 文件或目录的总大小(字节)
//   - error: 路径不存在或访问失败时返回错误
func GetSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, wrapPathError(err, path, "accessing")
	}

	// 如果是普通文件，直接返回文件大小
	if info.Mode().IsRegular() {
		return info.Size(), nil
	}

	// 如果不是目录，提前返回 0(符号链接等特殊文件)
	if !info.IsDir() {
		return 0, nil
	}

	// 如果是目录，遍历计算总大小
	var totalSize int64
	walkDirErr := filepath.WalkDir(path, func(walkPath string, entry os.DirEntry, err error) error {
		if err != nil {
			// 对于不存在的文件，忽略并继续遍历
			if os.IsNotExist(err) {
				return nil
			}
			return wrapPathError(err, walkPath, "accessing")
		}

		// 只计算普通文件的大小
		if entry.Type().IsRegular() {
			fileInfo, err := entry.Info()
			if err != nil {
				// 文件在遍历过程中被删除，忽略
				if os.IsNotExist(err) {
					return nil
				}
				return wrapPathError(err, walkPath, "getting file info")
			}

			// 只计算普通文件的大小
			totalSize += fileInfo.Size()
		}

		return nil
	})

	if walkDirErr != nil {
		return 0, fmt.Errorf("failed to walk directory: %w", walkDirErr)
	}

	return totalSize, nil
}
