// Package hash 实现了文件哈希计算命令的主要逻辑。
// 该文件包含 hash 子命令的入口函数，负责参数验证、文件收集和哈希计算任务的执行。
package hash

import (
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/internal/types"
)

// HashConfig 哈希配置结构体
type HashConfig struct {
	TargetPaths []string // 目标路径列表
	Type        string   // 哈希算法类型
	Recursion   bool     // 递归处理目录
	Write       bool     // 写入文件
	Hidden      bool     // 处理隐藏文件
	Progress    bool     // 显示进度条
	Local       bool     // 本地模式
	BasePath    string   // 基准路径
}

// HashCmdMain 执行哈希命令
//
// 参数:
//   - cl: 颜色库实例
//   - config: 哈希配置
//
// 返回值:
//   - error: 哈希计算过程中可能发生的错误
func HashCmdMain(cl *colorlib.ColorLib, config HashConfig) error {
	// 未指定路径时默认匹配当前目录所有文件
	targetPaths := config.TargetPaths
	if len(targetPaths) == 0 {
		targetPaths = []string{"*"}
	}

	if config.Write {
		cl.Green("Generating checksum file...")
	}

	// 逐个处理目标路径，单个路径失败不影响其他路径
	for _, targetPath := range targetPaths {
		if err := processSinglePath(cl, filepath.Clean(targetPath), config); err != nil {
			cl.Redf("处理路径 %s 时发生错误: %v\n", targetPath, err)
		}
	}

	return nil
}

// processSinglePath 处理单个路径
//
// 参数:
//   - cl: 颜色库对象
//   - targetPath: 目标路径
//   - config: 哈希配置
//
// 返回:
//   - error: 错误信息
func processSinglePath(cl *colorlib.ColorLib, targetPath string, config HashConfig) error {
	// 设置隐藏文件处理标志，供 collectFiles 使用
	setHiddenFlag(config.Hidden)
	files, err := collectFiles(targetPath, config.Recursion, cl)
	if err != nil {
		return fmt.Errorf("收集文件失败: %w", err)
	}

	if len(files) == 0 {
		cl.Yellowf("路径 %s 没有找到任何文件\n", targetPath)
		return nil
	}

	// 便携模式需要将绝对路径转换为相对路径，便于跨机器使用
	if config.Write && !config.Local {
		files, err = convertToRelativePaths(files, config.BasePath)
		if err != nil {
			return fmt.Errorf("转换相对路径失败: %w", err)
		}
	}

	processed, errors := hashRunTasksRefactored(files, config.Type, config)

	if len(errors) > 0 {
		printUniqueErrors(cl, errors)
	} else if config.Write {
		cl.Greenf("Checksum saved: %s (%d files)\n", types.OutputFileName, processed)
	}

	return nil
}

// convertToRelativePaths 将文件路径转换为相对路径
//
// 参数:
//   - files: 文件路径列表
//   - basePath: 基准路径
//
// 返回:
//   - []string: 转换后的相对路径列表
//   - error: 错误信息
func convertToRelativePaths(files []string, basePath string) ([]string, error) {
	// 未指定基准路径时使用当前工作目录
	if basePath == "" {
		var err error
		basePath, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("获取当前工作目录失败: %v", err)
		}
	}

	var relativePaths []string
	for _, file := range files {
		var relPath string

		if !filepath.IsAbs(file) {
			relPath = file
		} else {
			var err error
			relPath, err = filepath.Rel(basePath, file)
			if err != nil {
				return nil, fmt.Errorf("无法转换路径 %s: %v", file, err)
			}
		}

		// 统一使用正斜杠，确保跨平台兼容性
		relPath = filepath.ToSlash(relPath)
		relativePaths = append(relativePaths, relPath)
	}
	return relativePaths, nil
}

// printUniqueErrors 去重并打印错误信息
//
// 参数:
//   - cl: 颜色库对象
//   - errors: 错误列表
func printUniqueErrors(cl *colorlib.ColorLib, errors []error) {
	if len(errors) == 0 {
		return
	}

	seen := make(map[string]struct{}, len(errors))

	for _, err := range errors {
		if err == nil {
			continue
		}

		errStr := err.Error()
		if _, exists := seen[errStr]; !exists {
			seen[errStr] = struct{}{}
			cl.Red(errStr)
		}
	}
}
