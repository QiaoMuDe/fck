package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Expand 展开文件路径列表
// 将包含通配符的模式展开为具体路径，返回所有匹配的文件和目录
//
// 参数:
//   - patterns: 文件路径模式列表，支持 *、?、[] 等通配符
//
// 返回值:
//   - []string: 展开后的路径列表（去重），包含文件和目录
//   - error: 模式语法错误时返回错误
//
// 行为说明:
//   - 通配符匹配成功: 返回匹配的所有路径（文件和目录）
//   - 通配符无匹配: 保留原模式（不报错，由调用者处理）
//   - 具体路径: 直接保留
//
// 示例:
//
//	paths, err := fs.Expand([]string{"*.go"})                    // [main.go utils.go]
//	paths, err := fs.Expand([]string{"src/*"})                   // [src/main.go src/utils src/pkg]
//	paths, err := fs.Expand([]string{"config.yaml"})             // [config.yaml]（原样保留）
//	paths, err := fs.Expand([]string{"*.notexist"})              // [*.notexist]（无匹配时保留）
func Expand(patterns []string) ([]string, error) {
	// 检查输入是否为空
	if len(patterns) == 0 {
		return []string{}, nil
	}

	// 初始化去重集合和结果切片
	// 预估容量：每个 pattern 平均匹配 2 个路径
	seen := make(map[string]bool, len(patterns)*2)
	result := make([]string, 0, len(patterns)*2)

	// 遍历模式列表，展开每个模式
	for _, pattern := range patterns {
		// 展开当前模式
		matches, err := expandPattern(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
		}

		// 过滤去重，只保留未被添加到结果列表的路径
		for _, match := range matches {
			// 清洗路径，去除 ./ 前缀和多余的分隔符
			cleanMatch := filepath.Clean(match)
			if !seen[cleanMatch] {
				seen[cleanMatch] = true
				result = append(result, cleanMatch)
			}
		}
	}

	return result, nil
}

// ExpandFiles 展开文件路径列表，只返回文件（排除目录）
// 与 Expand 类似，但自动过滤掉目录路径
//
// 参数:
//   - patterns: 文件路径模式列表，支持 *、?、[] 等通配符
//
// 返回值:
//   - []string: 展开后的文件路径列表（去重），只包含文件
//   - error: 模式语法错误时返回错误
//
// 示例:
//
//	files, err := fs.ExpandFiles([]string{"*.go"})               // [main.go utils.go]
//	files, err := fs.ExpandFiles([]string{"src/*"})              // [src/main.go]（排除 src/utils 目录）
func ExpandFiles(patterns []string) ([]string, error) {
	paths, err := Expand(patterns)
	if err != nil {
		return nil, err
	}

	// 过滤掉目录路径，只保留文件路径
	// 预估容量：每个 pattern 平均匹配 2 个路径
	files := make([]string, 0, len(paths))
	for _, path := range paths {
		if !isDir(path) {
			files = append(files, path)
		}
	}

	return files, nil
}

// expandPattern 展开单个模式
// 内部使用 filepath.Glob 实现
//
// 参数:
//   - pattern: 文件路径模式
//
// 返回值:
//   - []string: 展开后的路径列表
//   - error: 模式语法错误时返回错误
func expandPattern(pattern string) ([]string, error) {
	// 使用 filepath.Glob
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// 无匹配时保留原路径
	if len(matches) == 0 {
		return []string{pattern}, nil
	}

	return matches, nil
}

// isDir 检查路径是否为目录
//
// 参数:
//   - path: 要检查的路径
//
// 返回值:
//   - bool: 是否为目录
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ExpandPattern 展开单个路径模式
// 与 Expand 类似，但只接收单个模式，返回所有匹配的路径（文件和目录）
//
// 参数:
//   - pattern: 文件路径模式，支持 *、?、[] 等通配符
//
// 返回值:
//   - []string: 展开后的路径列表，包含文件和目录
//   - error: 模式语法错误时返回错误
//
// 示例:
//
//	paths, err := fs.ExpandPattern("*.go")                 // [main.go utils.go]
//	paths, err := fs.ExpandPattern("src/*")                // [src/main.go src/utils src/pkg]
//	paths, err := fs.ExpandPattern("config.yaml")          // [config.yaml]（原样保留）
//	paths, err := fs.ExpandPattern("*.notexist")           // [*.notexist]（无匹配时保留）
func ExpandPattern(pattern string) ([]string, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	matches, err := expandPattern(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}

	return matches, nil
}
