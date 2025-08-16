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
		files, err := s.scanSinglePath(path, opts)
		if err != nil {
			return nil, fmt.Errorf("扫描路径 %s 失败: %v", path, err)
		}
		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}

// scanSinglePath 扫描单个路径
//
// 参数:
//   - path: 要扫描的路径
//   - opts: 扫描选项
//
// 返回:
//   - FileInfoList: 扫描到的文件信息列表
//   - error: 扫描过程中的错误
func (s *FileScanner) scanSinglePath(path string, opts ScanOptions) (FileInfoList, error) {
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
		return s.scanDirectory(absPath, absPath, opts)
	} else {
		fileInfo := s.buildFileInfo(pathInfo, absPath, absPath)
		return FileInfoList{fileInfo}, nil
	}
}

// scanDirectory 扫描目录
//
// 参数:
//   - dirPath: 目录路径
//   - rootDir: 根目录
//   - opts: 扫描选项
//
// 返回:
//   - FileInfoList: 扫描到的文件信息列表
//   - error: 扫描过程中的错误
func (s *FileScanner) scanDirectory(dirPath, rootDir string, opts ScanOptions) (FileInfoList, error) {
	var files FileInfoList

	// 如果只显示目录本身
	if opts.DirItself {
		pathInfo, err := s.getFileInfo(dirPath)
		if err != nil {
			return nil, err
		}
		fileInfo := s.buildFileInfo(pathInfo, dirPath, dirPath)
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
			subFiles, err := s.scanDirectory(absEntryPath, rootDir, opts)
			if err != nil {
				return nil, fmt.Errorf("递归扫描目录 %s 失败: %v", absEntryPath, err)
			}
			files = append(files, subFiles...)
		}

		// 添加当前文件信息
		info := s.buildFileInfo(fileInfo, absEntryPath, rootDir)
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

// buildFileInfo 构建文件信息
//
// 参数:
//   - fileInfo: 文件信息
//   - absPath: 绝对路径
//   - rootDir: 根目录
//
// 返回:
//   - FileInfo: 文件信息
func (s *FileScanner) buildFileInfo(fileInfo os.FileInfo, absPath string, rootDir string) FileInfo {
	// 确定显示名称
	var baseName string
	if listCmdRecursion.Get() {
		relPath, err := filepath.Rel(rootDir, absPath)
		if err != nil {
			baseName = absPath
		} else {
			baseName = relPath
		}
	} else {
		baseName = filepath.Base(absPath)
	}

	// 获取文件类型
	entryType := s.getEntryType(fileInfo)

	// 获取文件扩展名
	var fileExt string
	if strings.HasPrefix(baseName, ".") && len(baseName) > 1 {
		if strings.Contains(baseName[1:], ".") {
			fileExt = filepath.Ext(baseName)
		}
	} else if strings.Contains(baseName, ".") {
		fileExt = filepath.Ext(baseName)
	}

	// 获取符号链接目标
	var linkTargetPath string
	if entryType == types.SymlinkType {
		linkTargetPath, _ = os.Readlink(absPath)
		if linkTargetPath == "" {
			linkTargetPath = "?"
		}
	}

	// 获取文件所有者信息
	owner, group := common.GetFileOwner(absPath)

	return FileInfo{
		EntryType:      entryType,
		Name:           baseName,
		Size:           fileInfo.Size(),
		ModTime:        fileInfo.ModTime(),
		Perm:           fileInfo.Mode().Perm().String(),
		Owner:          owner,
		Group:          group,
		FileExt:        fileExt,
		LinkTargetPath: linkTargetPath,
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
			case ".exe", ".com", ".cmd", ".bat", ".ps1", ".psm1":
				return types.ExecutableType
			case ".lnk", ".url":
				return types.SymlinkType
			}
		case "linux", "darwin":
			if mode&0111 != 0 {
				return types.ExecutableType
			}
		}

		return types.FileType
	}

	return types.UnknownType
}
