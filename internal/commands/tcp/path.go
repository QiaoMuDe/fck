package tcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathType 路径类型
type PathType int

const (
	PathTypeFile     PathType = iota // 单个文件
	PathTypeDir                      // 目录
	PathTypeWildcard                 // 通配符模式
)

// PathInfo 路径信息
type PathInfo struct {
	Type  PathType // 路径类型
	Paths []string // 匹配的文件路径列表
}

// resolvePath 解析路径，支持文件、目录和通配符
//
// 参数:
//   - path: 路径字符串
//   - maxFileSize: 最大文件大小限制（字节）
//
// 返回值:
//   - *PathInfo: 路径信息
//   - error: 解析错误
func resolvePath(path string, maxFileSize int64) (*PathInfo, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	// 检查是否包含通配符
	if strings.Contains(path, "*") || strings.Contains(path, "?") {
		return resolveWildcard(path, maxFileSize)
	}

	// 获取文件信息
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if info.IsDir() {
		// 目录：读取目录下所有文件（非递归）
		return resolveDir(path, maxFileSize)
	}

	// 单个文件
	// 检查文件大小
	if info.Size() > maxFileSize {
		return nil, fmt.Errorf("file too large: %s (%d bytes, max %d bytes)", path, info.Size(), maxFileSize)
	}

	return &PathInfo{
		Type:  PathTypeFile,
		Paths: []string{path},
	}, nil
}

// resolveDir 解析目录，返回目录下所有文件（非递归）
//
// 参数:
//   - dir: 目录路径
//   - maxFileSize: 最大文件大小限制（字节）
//
// 返回值:
//   - *PathInfo: 路径信息
//   - error: 解析错误
func resolveDir(dir string, maxFileSize int64) (*PathInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dir, entry.Name())
			// 检查文件大小
			info, err := entry.Info()
			if err != nil {
				continue // 跳过无法获取信息的文件
			}
			if info.Size() > maxFileSize {
				continue // 跳过大文件
			}
			files = append(files, filePath)
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in directory: %s", dir)
	}

	return &PathInfo{
		Type:  PathTypeDir,
		Paths: files,
	}, nil
}

// resolveWildcard 解析通配符模式
//
// 参数:
//   - pattern: 通配符模式
//   - maxFileSize: 最大文件大小限制（字节）
//
// 返回值:
//   - *PathInfo: 路径信息
//   - error: 解析错误
func resolveWildcard(pattern string, maxFileSize int64) (*PathInfo, error) {
	// 统一路径分隔符为系统分隔符（Windows 上 filepath.Glob 需要反斜杠）
	normalizedPattern := filepath.FromSlash(pattern)
	matches, err := filepath.Glob(normalizedPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid wildcard pattern: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files match pattern: %s", pattern)
	}

	// 过滤掉目录和大文件
	var files []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err == nil && !info.IsDir() {
			// 检查文件大小
			if info.Size() > maxFileSize {
				continue // 跳过大文件
			}
			files = append(files, match)
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files match pattern: %s", pattern)
	}

	return &PathInfo{
		Type:  PathTypeWildcard,
		Paths: files,
	}, nil
}
