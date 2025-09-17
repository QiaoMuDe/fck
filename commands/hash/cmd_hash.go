// Package hash 实现了文件哈希计算命令的主要逻辑。
// 该文件包含 hash 子命令的入口函数，负责参数验证、文件收集和哈希计算任务的执行。
package hash

import (
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// HashCmdMain 是 hash 子命令的主函数
func HashCmdMain(cl *colorlib.ColorLib) error {
	// 验证输入参数
	targetPaths := hashCmd.Args()
	if len(targetPaths) == 0 {
		targetPaths = []string{"*"} // 默认处理当前目录
	}

	// 显示写入文件提示
	if hashCmdWrite.Get() {
		cl.PrintOk("正在将哈希值写入文件，请稍候...")
	}

	// 遍历所有目标路径
	for _, targetPath := range targetPaths {
		if err := processSinglePath(cl, filepath.Clean(targetPath), hashCmdType.Get()); err != nil {
			// 记录错误但继续处理其他路径
			cl.PrintErrorf("处理路径 %s 时发生错误: %v\n", targetPath, err)
		}
	}

	return nil
}

// processSinglePath 处理单个路径
//
// 参数:
//   - cl: 颜色库对象
//   - targetPath: 目标路径
//   - hashType: 哈希算法类型
//
// 返回:
//   - error: 错误信息
func processSinglePath(cl *colorlib.ColorLib, targetPath string, hashType string) error {
	// 收集文件
	files, err := collectFiles(targetPath, hashCmdRecursion.Get(), cl)
	if err != nil {
		return fmt.Errorf("收集文件失败: %w", err)
	}

	// 检查文件列表是否为空
	if len(files) == 0 {
		cl.PrintWarnf("路径 %s 没有找到任何文件\n", targetPath)
		return nil
	}

	// 如果是便携模式且需要写入文件，转换为相对路径
	if hashCmdWrite.Get() && !hashCmdLocal.Get() {
		files, err = convertToRelativePaths(files)
		if err != nil {
			return fmt.Errorf("转换相对路径失败: %w", err)
		}
	}

	// 执行哈希任务
	errors := hashRunTasksRefactored(files, hashType)

	// 处理执行结果
	if len(errors) > 0 {
		printUniqueErrors(cl, errors)
	} else if hashCmdWrite.Get() {
		cl.PrintOkf("已将哈希值写入文件 %s, 共处理 %d 个文件\n", types.OutputFileName, len(files))
	}

	return nil
}

// convertToRelativePaths 将文件路径转换为相对路径
//
// 参数:
//   - files: 文件路径列表
//
// 返回:
//   - []string: 转换后的相对路径列表
//   - error: 错误信息
func convertToRelativePaths(files []string) ([]string, error) {
	// 获取基准路径
	basePath := hashCmdBasePath.Get()
	if basePath == "" {
		var err error
		basePath, err = os.Getwd() // 默认使用当前工作目录
		if err != nil {
			return nil, fmt.Errorf("获取当前工作目录失败: %v", err)
		}
	}

	var relativePaths []string
	for _, file := range files {
		// 将绝对路径转换为相对于basePath的相对路径
		relPath, err := filepath.Rel(basePath, file)
		if err != nil {
			return nil, fmt.Errorf("无法转换路径 %s: %v", file, err)
		}
		// 统一使用正斜杠作为分隔符（跨平台兼容）
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
			cl.PrintError(errStr)
		}
	}
}
