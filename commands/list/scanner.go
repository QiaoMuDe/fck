// Package list 实现了文件系统扫描功能。
// 该文件提供了文件和目录的扫描、过滤、类型识别等核心功能，支持递归扫描和多种文件类型过滤。
package list

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
)

// FileScanner 文件扫描器
type FileScanner struct {
	cache map[string]os.FileInfo // 缓存文件信息
	mutex sync.RWMutex           // 缓存锁
}

// NewFileScanner 创建新的文件扫描器
func NewFileScanner() *FileScanner {
	return &FileScanner{
		cache: make(map[string]os.FileInfo),
	}
}

// Scan 扫描指定路径的文件
//
// 参数:
//   - paths: 要扫描的路径列表
//   - opts: 扫描选项
//
// 返回:
//   - FileInfoList: 扫描到的文件信息列表
//   - error: 扫描过程中的错误
func (s *FileScanner) Scan(paths []string, opts ScanOptions) (FileInfoList, error) {
	var allFiles FileInfoList

	for _, path := range paths {
		files, err := s.scanSinglePathWithOriginal(path, path, opts)
		if err != nil {
			return nil, fmt.Errorf("扫描路径 %s 失败: %v", path, err)
		}
		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}

// ScanWithOriginalPaths 扫描指定路径的文件（保持原始路径和展开路径的对应关系）
//
// 参数:
//   - originalPaths: 用户输入的原始路径列表
//   - expandedPaths: 展开后的路径列表
//   - opts: 扫描选项
//
// 返回:
//   - FileInfoList: 扫描到的文件信息列表
//   - error: 扫描过程中的错误
func (s *FileScanner) ScanWithOriginalPaths(originalPaths, expandedPaths []string, opts ScanOptions) (FileInfoList, error) {
	var allFiles FileInfoList

	// 创建原始路径到展开路径的映射
	pathMapping := s.createPathMapping(originalPaths, expandedPaths)

	for _, expandedPath := range expandedPaths {
		// 找到对应的原始路径
		originalPath := s.findOriginalPath(expandedPath, pathMapping)

		// 检查是否为通配符展开的目录
		isWildcardDir := strings.ContainsAny(originalPath, "*?[]")
		if isWildcardDir {
			if info, err := os.Stat(expandedPath); err == nil && info.IsDir() {
				// 通配符展开的目录：扫描目录内容，但保持原始路径为目录路径
				files, err := s.scanDirectoryWithOriginal(expandedPath, expandedPath, expandedPath, opts)
				if err != nil {
					return nil, fmt.Errorf("扫描目录 %s 失败: %v", expandedPath, err)
				}
				allFiles = append(allFiles, files...)
				continue
			}
		}

		files, err := s.scanSinglePathWithOriginal(expandedPath, originalPath, opts)
		if err != nil {
			return nil, fmt.Errorf("扫描路径 %s 失败: %v", expandedPath, err)
		}
		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}

// createPathMapping 创建原始路径到展开路径的映射关系
func (s *FileScanner) createPathMapping(originalPaths, expandedPaths []string) map[string]string {
	mapping := make(map[string]string)

	// 如果原始路径和展开路径数量相同，则一一对应
	if len(originalPaths) == len(expandedPaths) {
		for i, expandedPath := range expandedPaths {
			mapping[expandedPath] = originalPaths[i]
		}
		return mapping
	}

	// 处理通配符展开的情况
	for _, originalPath := range originalPaths {
		if strings.ContainsAny(originalPath, "*?[]") {
			// 通配符路径，找到所有匹配的展开路径
			matches, _ := filepath.Glob(originalPath)
			for _, match := range matches {
				// 找到在expandedPaths中的对应项
				for _, expandedPath := range expandedPaths {
					if match == expandedPath {
						mapping[expandedPath] = originalPath
					}
				}
			}
		} else {
			// 非通配符路径，直接映射
			mapping[originalPath] = originalPath
		}
	}

	return mapping
}

// findOriginalPath 根据展开路径找到对应的原始路径
func (s *FileScanner) findOriginalPath(expandedPath string, mapping map[string]string) string {
	if originalPath, exists := mapping[expandedPath]; exists {
		return originalPath
	}
	// 如果找不到映射，返回展开路径本身
	return expandedPath
}

// scanSinglePathWithOriginal 扫描单个路径（保存原始路径）
//
// 参数:
//   - path: 要扫描的路径
//   - originalPath: 用户输入的原始路径
//   - opts: 扫描选项
//
// 返回:
//   - FileInfoList: 扫描到的文件信息列表
//   - error: 扫描过程中的错误
func (s *FileScanner) scanSinglePathWithOriginal(path string, originalPath string, opts ScanOptions) (FileInfoList, error) {
	// 清理路径
	path = filepath.Clean(path)

	// 获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %v", err)
	}

	// 检查是否为系统文件或目录
	if common.IsSystemFileOrDir(filepath.Base(absPath)) {
		return nil, fmt.Errorf("不能列出系统文件或目录: %s", absPath)
	}

	// 获取文件信息
	pathInfo, err := s.getFileInfo(absPath)
	if err != nil {
		return nil, common.HandleError(absPath, err)
	}

	// 检查是否应该跳过
	if s.shouldSkipFile(absPath, pathInfo.IsDir(), pathInfo, true, opts) {
		return FileInfoList{}, nil
	}

	// 根据是否为目录进行处理
	if pathInfo.IsDir() {
		return s.scanDirectoryWithOriginal(absPath, absPath, originalPath, opts)
	} else {
		fileInfo := s.buildFileInfoWithOriginal(pathInfo, absPath, originalPath)
		return FileInfoList{fileInfo}, nil
	}
}

// scanDirectoryWithOriginal 扫描目录（保存原始路径）
//
// 参数:
//   - dirPath: 目录路径
//   - rootDir: 根目录
//   - originalPath: 用户输入的原始路径
//   - opts: 扫描选项
//
// 返回:
//   - FileInfoList: 扫描到的文件信息列表
//   - error: 扫描过程中的错误
func (s *FileScanner) scanDirectoryWithOriginal(dirPath, rootDir string, originalPath string, opts ScanOptions) (FileInfoList, error) {
	var files FileInfoList

	// 如果只显示目录本身
	if opts.DirItself {
		pathInfo, err := s.getFileInfo(dirPath)
		if err != nil {
			return nil, err
		}
		fileInfo := s.buildFileInfoWithOriginal(pathInfo, dirPath, originalPath)
		return FileInfoList{fileInfo}, nil
	}

	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, common.HandleError(dirPath, err)
	}

	// 处理目录中的每个条目
	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		absEntryPath, err := filepath.Abs(entryPath)
		if err != nil {
			continue
		}

		// 检查是否为系统文件
		if common.IsSystemFileOrDir(filepath.Base(absEntryPath)) {
			continue
		}

		// 获取文件信息
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		// 检查是否应该跳过
		if s.shouldSkipFile(absEntryPath, entry.IsDir(), fileInfo, false, opts) {
			continue
		}

		// 如果是递归模式且当前是目录
		if opts.Recursive && entry.IsDir() {
			subFiles, err := s.scanDirectoryWithOriginal(absEntryPath, rootDir, originalPath, opts)
			if err != nil {
				return nil, fmt.Errorf("递归扫描目录 %s 失败: %v", absEntryPath, err)
			}
			files = append(files, subFiles...)
		}

		// 添加当前文件信息
		info := s.buildFileInfoWithOriginal(fileInfo, absEntryPath, originalPath)
		files = append(files, info)
	}

	return files, nil
}

// getFileInfo 获取文件信息（带缓存）
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - os.FileInfo: 文件信息
//   - error: 获取过程中的错误
func (s *FileScanner) getFileInfo(path string) (os.FileInfo, error) {
	// 命中缓存
	s.mutex.RLock()
	if info, exists := s.cache[path]; exists {
		s.mutex.RUnlock()
		return info, nil
	}
	s.mutex.RUnlock()

	// 未命中缓存, 获取并缓存文件信息
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	s.mutex.Lock()
	s.cache[path] = info
	s.mutex.Unlock()

	return info, nil
}

// shouldSkipFile 判断是否应该跳过文件
//
// 参数:
//   - path: 文件路径
//   - isDir: 是否为目录
//   - fileInfo: 文件信息
//   - isMain: 是否为处理目录本身
//   - opts: 扫描选项
//
// 返回:
//   - bool: 是否应该跳过
func (s *FileScanner) shouldSkipFile(path string, isDir bool, fileInfo os.FileInfo, isMain bool, opts ScanOptions) bool {
	// 隐藏文件检查
	if !opts.ShowHidden && common.IsHidden(path) {
		return true
	}

	// 如果显示隐藏文件但指定了隐藏文件类型过滤
	if opts.ShowHidden && len(opts.FileTypes) > 0 {
		for _, fileType := range opts.FileTypes {
			if (fileType == types.FindTypeHiddenShort || fileType == types.FindTypeHidden) && !common.IsHidden(path) {
				return true
			}
		}
	}

	// 只有不是处理目录本身时，才进行类型过滤
	if !isMain && len(opts.FileTypes) > 0 {
		for _, fileType := range opts.FileTypes {
			switch fileType {
			case types.FindTypeFileShort, types.FindTypeFile:
				if isDir {
					return true
				}
			case types.FindTypeDirShort, types.FindTypeDir:
				if !isDir {
					return true
				}
			case types.FindTypeSymlinkShort, types.FindTypeSymlink:
				if fileInfo.Mode()&os.ModeSymlink == 0 {
					return true
				}
			case types.FindTypeReadonly, types.FindTypeReadonlyShort:
				if !common.IsReadOnly(path) {
					return true
				}
			}
		}
	}

	return false
}

// buildFileInfoWithOriginal 构建文件信息（包含原始路径）
//
// 参数:
//   - fileInfo: 文件信息
//   - absPath: 绝对路径
//   - originalPath: 用户指定的原始路径
//
// 返回:
//   - FileInfo: 文件信息
func (s *FileScanner) buildFileInfoWithOriginal(fileInfo os.FileInfo, absPath string, originalPath string) FileInfo {
	// 确定显示名称
	var baseName string
	if listCmdRecursion.Get() {
		baseName = absPath // 递归模式下显示绝对路径
	} else {
		baseName = filepath.Base(absPath) // 非递归模式下显示文件名
	}

	// 获取文件类型
	entryType := s.getEntryType(fileInfo)

	// 获取文件扩展名
	var fileExt string
	if strings.HasPrefix(baseName, ".") && len(baseName) > 1 {
		if strings.Contains(baseName[1:], ".") {
			fileExt = filepath.Ext(baseName) // 隐藏文件扩展名
		}

	} else if strings.Contains(baseName, ".") {
		fileExt = filepath.Ext(baseName) // 普通文件扩展名
	}

	// 获取符号链接目标
	var linkTargetPath string
	if entryType == types.SymlinkType {
		linkTargetPath, _ = os.Readlink(absPath)
		if linkTargetPath == "" {
			linkTargetPath = "?" // 读取失败或空链接时返回问号
		}
	}

	// 获取文件所有者信息
	owner, group := common.GetFileOwner(absPath)

	return FileInfo{
		EntryType:      entryType,                       // 文件类型
		Name:           baseName,                        // 显示名称
		Path:           absPath,                         // 绝对路径
		OriginalPath:   originalPath,                    // 用户指定的原始路径
		Size:           fileInfo.Size(),                 // 文件大小
		ModTime:        fileInfo.ModTime(),              // 修改时间
		Perm:           fileInfo.Mode().Perm().String(), // 权限信息
		Owner:          owner,                           // 文件所有者
		Group:          group,                           // 文件所属组
		FileExt:        fileExt,                         // 文件扩展名
		LinkTargetPath: linkTargetPath,                  // 符号链接目标
	}
}

// getEntryType 获取文件类型
//
// 参数:
//   - fileInfo: 文件信息
//
// 返回:
//   - string: 文件类型
func (s *FileScanner) getEntryType(fileInfo os.FileInfo) string {
	mode := fileInfo.Mode()

	// 检查符号链接
	if mode&os.ModeSymlink != 0 {
		return types.SymlinkType
	}

	// 检查目录
	if mode.IsDir() {
		return types.DirType
	}

	// 检查特殊文件类型
	if mode&os.ModeSocket != 0 {
		return types.SocketType
	}
	if mode&os.ModeNamedPipe != 0 {
		return types.PipeType
	}
	if mode&os.ModeDevice != 0 {
		if mode&os.ModeCharDevice != 0 {
			return types.CharDeviceType
		}
		return types.BlockDeviceType
	}

	// 检查普通文件
	if mode.IsRegular() {
		// 空文件
		if fileInfo.Size() == 0 {
			return types.EmptyType
		}

		// 可执行文件检查
		switch runtime.GOOS {
		case "windows":
			ext := strings.ToLower(filepath.Ext(fileInfo.Name()))
			switch ext {
			case ".exe", ".bat":
				return types.ExecutableType
			case ".lnk", ".url":
				return types.SymlinkType
			}

		default:
			if mode&0111 != 0 {
				return types.ExecutableType
			}
		}

		return types.FileType
	}

	return types.UnknownType
}
