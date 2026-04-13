package grep

import (
	"io/fs"
	"os"
	"path/filepath"
)

// processSingleFile 处理单个文件（递归模式使用）
// 根据 NoMessages 配置决定是否返回错误
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果 NoMessages 为 false)
func processSingleFile(path string, config *GrepConfig) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if config.NoMessages {
			return nil // 静默模式，跳过错误
		}
		if os.IsNotExist(err) {
			return nil // 递归模式下文件不存在不报错
		}
		return nil // 权限错误等，跳过
	}
	if fileInfo.IsDir() {
		return nil // 是目录，跳过
	}

	file, err := os.Open(path)
	if err != nil {
		if config.NoMessages {
			return nil // 静默模式，跳过错误
		}
		return nil // 无法打开，跳过
	}
	defer func() { _ = file.Close() }()

	// 处理二进制文件
	skip, err := handleBinaryFile(file, path, config)
	if err != nil {
		return nil // 递归模式下跳过错误
	}
	if skip {
		return nil
	}

	config.filename = path
	return processFile(file, config)
}

// processPathRecursive 递归处理目录
// 递归模式下强制显示文件路径
//
// 参数:
//   - path: 目录或文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果 NoMessages 为 false)
func processPathRecursive(path string, config *GrepConfig) error {
	// 递归模式强制显示文件名
	config.WithFilename = true

	fileInfo, err := os.Stat(path)
	if err != nil {
		if config.NoMessages {
			return nil // 静默模式，跳过错误
		}
		if os.IsNotExist(err) {
			return nil // 路径不存在，跳过
		}
		return nil // 权限错误等，跳过
	}

	// 如果是文件，直接处理
	if !fileInfo.IsDir() {
		return processSingleFile(path, config)
	}

	// 递归遍历目录
	visited := make(map[string]bool) // 用于检测符号链接循环
	return filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			// 权限错误等，根据 NoMessages 决定是否返回
			if config.NoMessages {
				return nil
			}
			return nil // 递归遍历中跳过错误
		}

		// 检查隐藏文件/目录（默认跳过）
		// 跳过根目录 "." 的检查
		if !config.Hidden && d.Name() != "." && isHidden(d.Name()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 处理符号链接
		if d.Type()&fs.ModeSymlink != 0 {
			if !config.FollowSymlink {
				return nil // 跳过符号链接
			}
			// 解析符号链接
			realPath, err := filepath.EvalSymlinks(filePath)
			if err != nil {
				return nil // 无法解析，跳过
			}
			// 检查循环
			if visited[realPath] {
				return nil // 循环引用，跳过
			}
			visited[realPath] = true

			// 获取符号链接指向的文件信息
			info, err := os.Stat(realPath)
			if err != nil {
				return nil
			}
			if info.IsDir() {
				// 是目录，继续遍历
				return nil
			}
			// 是文件，处理它
			return processSingleFile(realPath, config)
		}

		// 目录过滤
		if d.IsDir() {
			dirName := d.Name()
			if !matchDirFilters(dirName, config) {
				return filepath.SkipDir
			}
			return nil
		}

		// 文件过滤
		fileName := d.Name()
		if !matchFileFilters(fileName, config) {
			return nil
		}

		// 处理文件
		return processSingleFile(filePath, config)
	})
}

// matchFileFilters 检查文件是否匹配过滤条件
//
// 参数:
//   - fileName: 文件名
//   - config: 命令配置
//
// 返回:
//   - bool: 是否匹配
func matchFileFilters(fileName string, config *GrepConfig) bool {
	// 检查 exclude
	for _, pattern := range config.Exclude {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return false
		}
	}

	// 检查 include（如果有指定，必须匹配其中一个）
	if len(config.Include) > 0 {
		for _, pattern := range config.Include {
			if matched, _ := filepath.Match(pattern, fileName); matched {
				return true
			}
		}
		return false
	}

	return true
}

// matchDirFilters 检查目录是否匹配过滤条件
//
// 参数:
//   - dirName: 目录名
//   - config: 命令配置
//
// 返回:
//   - bool: 是否匹配
func matchDirFilters(dirName string, config *GrepConfig) bool {
	// 检查 exclude-dir
	for _, pattern := range config.ExcludeDir {
		if matched, _ := filepath.Match(pattern, dirName); matched {
			return false
		}
	}

	// 检查 include-dir（如果有指定，必须匹配其中一个）
	if len(config.IncludeDir) > 0 {
		for _, pattern := range config.IncludeDir {
			if matched, _ := filepath.Match(pattern, dirName); matched {
				return true
			}
		}
		return false
	}

	return true
}
