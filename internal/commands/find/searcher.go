// Package find 实现了文件查找的核心搜索逻辑。
// 该文件提供了文件搜索器，负责遍历目录、应用过滤条件、执行操作和输出结果。
package find

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	common "gitee.com/MM-Q/fck/internal/utils"
)

// FileSearcher 负责核心搜索逻辑
type FileSearcher struct {
	config   *types.FindConfig // 查找配置
	matcher  *PatternMatcher   // 模式匹配器
	operator *FileOperator     // 文件操作器
}

// NewFileSearcher 创建新的文件搜索器
func NewFileSearcher(config *types.FindConfig, matcher *PatternMatcher, operator *FileOperator) *FileSearcher {
	return &FileSearcher{
		config:   config,
		matcher:  matcher,
		operator: operator,
	}
}

// Search 执行文件搜索
//
// 参数:
//   - findPath: 查找路径
//
// 返回:
//   - error: 搜索错误（如果有）
func (s *FileSearcher) Search(findPath string) error {
	walkDirErr := filepath.WalkDir(findPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}

			if os.IsPermission(err) {
				if !s.config.Quiet {
					s.config.Cl.PrintErrorf("权限不足, 无法访问某些目录: %s\n", path)
				}
				return nil
			}

			return fmt.Errorf("访问时出错：%s", err)
		}

		if path == findPath {
			return nil
		}

		depth := strings.Count(path[len(findPath):], string(filepath.Separator))
		if s.config.MaxDepth >= 0 && depth > s.config.MaxDepth {
			return filepath.SkipDir
		}

		if entry.Type()&os.ModeSymlink != 0 {
			if s.isSymlinkLoop(path) {
				return filepath.SkipDir
			}
		}

		if processErr := s.processEntry(entry, path); processErr != nil {
			return processErr
		}

		return nil
	})

	if walkDirErr != nil {
		return fmt.Errorf("遍历目录时出错: %v", walkDirErr)
	}

	return nil
}

// processEntry 处理单个文件或目录条目
//
// 参数:
//   - entry: 文件或目录条目
//   - path: 文件或目录的路径
//
// 返回:
//   - error: 处理错误（如果有）
func (s *FileSearcher) processEntry(entry os.DirEntry, path string) error {
	if s.config.NamePattern != "" && s.config.PathPattern != "" && s.config.And {
		if s.matcher.MatchName(entry.Name(), s.config.NamePattern, s.config) && s.matcher.MatchPath(path, s.config.PathPattern, s.config) {
			if err := s.applyFilters(entry, path); err != nil {
				return err
			}
		}
		return nil
	}

	if s.config.NamePattern != "" && s.config.PathPattern != "" && s.config.Or {
		if s.matcher.MatchName(entry.Name(), s.config.NamePattern, s.config) || s.matcher.MatchPath(path, s.config.PathPattern, s.config) {
			if err := s.applyFilters(entry, path); err != nil {
				return err
			}
		}
		return nil
	}

	if s.config.NamePattern != "" {
		if s.matcher.MatchName(entry.Name(), s.config.NamePattern, s.config) {
			if err := s.applyFilters(entry, path); err != nil {
				return err
			}
		}
		return nil
	}

	if s.config.PathPattern != "" {
		if s.matcher.MatchPath(path, s.config.PathPattern, s.config) {
			if err := s.applyFilters(entry, path); err != nil {
				return err
			}
		}
		return nil
	}

	if err := s.applyFilters(entry, path); err != nil {
		return err
	}

	return nil
}

// applyFilters 应用所有筛选条件
//
// 参数:
//   - entry: 文件或目录条目
//   - path: 文件或目录的路径
//
// 返回:
//   - error: 筛选错误（如果有）
func (s *FileSearcher) applyFilters(entry os.DirEntry, path string) error {
	if !s.config.Hidden && common.IsHidden(path) {
		if entry.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// 如果指定了排除文件或目录名, 跳过匹配的文件或目录
	if s.config.ExNamePattern != "" && s.matcher.MatchName(entry.Name(), s.config.ExNamePattern, s.config) {
		if entry.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// 如果指定了排除路径, 跳过匹配的路径
	if s.config.ExPathPattern != "" && s.matcher.MatchPath(path, s.config.ExPathPattern, s.config) {
		if entry.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// 仅在需要文件元信息时才获取
	var cacheInfo fs.FileInfo
	var cacheErr error

	// 判断是否需要获取文件信息
	if s.needFileInfo() {
		cacheInfo, cacheErr = entry.Info()
		if cacheErr != nil {
			return nil
		}
	}

	// 预计算文件扩展名并缓存
	entryExt := filepath.Ext(entry.Name())

	// 应用类型筛选
	if !s.matchType(entry, path, entryExt, cacheInfo) {
		return nil
	}

	if s.config.SizePattern != "" && !s.matcher.MatchSize(cacheInfo.Size(), s.config.SizePattern) {
		return nil
	}

	if s.config.ModTimePattern != "" && !s.matcher.MatchTime(cacheInfo.ModTime(), s.config.ModTimePattern) {
		return nil
	}

	if len(s.config.ExtSlice) > 0 {
		if _, ok := s.config.FindExtSliceMap.Load(entryExt); !ok {
			return nil
		}
	}

	// 执行操作或输出结果
	return s.executeAction(entry, path)
}

// needFileInfo 判断是否需要获取文件信息
//
// 返回:
//   - bool: 是否需要获取文件信息
//
// 描述:
//  1. 如果需要查找空文件或目录, 则需要获取文件信息
//  2. 如果指定了修改时间或文件大小, 则需要获取文件信息
func (s *FileSearcher) needFileInfo() bool {
	return (s.config.Type == types.FindTypeEmpty || s.config.Type == types.FindTypeEmptyShort) ||
		s.config.ModTimePattern != "" ||
		s.config.SizePattern != ""
}

// matchType 检查文件类型是否匹配指定的查找类型
//
// 该方法根据用户指定的文件类型参数（-type）来判断当前文件或目录是否符合条件。
// 支持多种文件类型的匹配，包括普通文件、目录、符号链接、隐藏文件等。
//
// 参数:
//   - entry: 文件或目录条目，包含基本的文件信息
//   - path: 文件或目录的完整路径
//   - entryExt: 文件扩展名（已预计算，用于性能优化）
//   - cacheInfo: 文件元信息（仅在需要时获取，可能为nil）
//
// 返回:
//   - bool: 文件类型是否匹配用户指定的查找类型
//
// 支持的文件类型:
//   - f/file: 普通文件（非目录）
//   - d/dir: 目录
//   - l/symlink: 符号链接（Windows下通过扩展名判断，Unix下通过文件模式）
//   - h/hidden: 隐藏文件或目录
//   - r/readonly: 只读文件或目录
//   - e/empty: 空文件或空目录
//   - x/executable: 可执行文件（Windows下通过扩展名，Unix下通过权限位）
//   - s/socket: 套接字文件（仅Unix系统支持）
//   - p/pipe: 命名管道（仅Unix系统支持）
//   - b/block: 块设备文件（仅Unix系统支持）
//   - c/char: 字符设备文件（仅Unix系统支持）
//   - a/append: 仅追加文件（仅Unix系统支持）
//   - A/non-append: 非仅追加文件（仅Unix系统支持）
//   - E/exclusive: 独占文件（仅Unix系统支持）
//
// 注意事项:
//  1. Windows系统不支持Unix特有的文件类型（socket、pipe、block等）
//  2. 符号链接在Windows下通过扩展名判断，在Unix下通过文件模式判断
//  3. 可执行文件在Windows下通过扩展名判断，在Unix下通过权限位判断
//  4. 空目录检查需要读取目录内容，可能影响性能
//  5. 如果没有指定类型参数，默认匹配所有类型
func (s *FileSearcher) matchType(entry os.DirEntry, path, entryExt string, cacheInfo fs.FileInfo) bool {
	switch s.config.Type {
	case types.FindTypeFile, types.FindTypeFileShort:
		return !entry.IsDir()

	case types.FindTypeDir, types.FindTypeDirShort:
		return entry.IsDir()

	case types.FindTypeSymlink, types.FindTypeSymlinkShort:
		if runtime.GOOS == "windows" {
			return types.WindowsSymlinkExts[entryExt]
		}
		return entry.Type()&os.ModeSymlink != 0

	case types.FindTypeHidden, types.FindTypeHiddenShort:
		return common.IsHidden(path)

	case types.FindTypeReadonly, types.FindTypeReadonlyShort:
		return common.IsReadOnly(path)

	case types.FindTypeEmpty, types.FindTypeEmptyShort:
		if entry.IsDir() {
			// 对于目录，检查是否为空目录
			dirEntries, err := os.ReadDir(path)
			return err == nil && len(dirEntries) == 0
		}
		// 对于文件，检查文件大小是否为0
		return cacheInfo != nil && cacheInfo.Size() == 0

	case types.FindTypeExecutable, types.FindTypeExecutableShort: // x, executable
		// 匹配可执行文件
		if runtime.GOOS == "windows" {
			// Windows下通过扩展名判断（.exe, .bat, .cmd等）
			ext := strings.ToLower(entryExt)
			return types.WindowsExecutableExts[ext]
		}
		// Unix系统下检查是否为普通文件且具有执行权限
		return entry.Type().IsRegular() && entry.Type()&0111 != 0

	case types.FindTypeSocket, types.FindTypeSocketShort: // s, socket
		// 匹配套接字文件（仅Unix系统支持）
		if runtime.GOOS == "windows" {
			return false // Windows不支持
		}
		return entry.Type()&os.ModeSocket != 0

	case types.FindTypePipe, types.FindTypePipeShort: // p, pipe
		// 匹配命名管道（仅Unix系统支持）
		if runtime.GOOS == "windows" {
			return false // Windows不支持
		}
		return entry.Type()&os.ModeNamedPipe != 0

	case types.FindTypeBlock, types.FindTypeBlockShort: // b, block
		// 匹配块设备文件（仅Unix系统支持）
		if runtime.GOOS == "windows" {
			return false // Windows不支持
		}
		return entry.Type()&os.ModeDevice != 0

	case types.FindTypeChar, types.FindTypeCharShort: // c, char
		// 匹配字符设备文件（仅Unix系统支持）
		if runtime.GOOS == "windows" {
			return false // Windows不支持
		}
		return entry.Type()&os.ModeCharDevice != 0

	case types.FindTypeAppend, types.FindTypeAppendShort: // a, append
		// 匹配仅追加文件（仅Unix系统支持）
		if runtime.GOOS == "windows" {
			return false // Windows不支持
		}
		return entry.Type()&os.ModeAppend != 0

	case types.FindTypeNonAppend, types.FindTypeNonAppendShort: // A, non-append
		// 匹配非仅追加文件（仅Unix系统支持）
		if runtime.GOOS == "windows" {
			return false // Windows不支持
		}
		return entry.Type()&os.ModeAppend == 0

	case types.FindTypeExclusive, types.FindTypeExclusiveShort: // E, exclusive
		// 匹配独占文件（仅Unix系统支持）
		if runtime.GOOS == "windows" {
			return false // Windows不支持
		}
		return entry.Type()&os.ModeExclusive != 0
	}

	// 默认情况：如果没有指定类型参数，匹配所有类型
	return true
}

// executeAction 执行相应的操作
//
// 参数:
//   - entry: 文件或目录的DirEntry对象
//   - path: 文件或目录的完整路径
//
// 返回:
//   - error: 如果发生错误，则返回错误信息；否则返回nil
func (s *FileSearcher) executeAction(entry os.DirEntry, path string) error {
	if !s.config.Count {
		if s.config.Delete {
			if err := s.operator.Delete(path); err != nil {
				return err
			}
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if s.config.MovePath != "" {
			if err := s.operator.Move(path, s.config.MovePath); err != nil {
				return err
			}
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if s.config.ExecCmd != "" {
			if err := s.operator.Execute(s.config.ExecCmd, path); err != nil {
				return fmt.Errorf("执行-exec命令时发生了错误: %v", err)
			}
			return nil
		}
	}

	s.outputResult(path, entry)
	return nil
}

// outputResult 输出搜索结果
//
// 参数:
//   - path: 文件或目录的路径
//   - d: 匹配到的DirEntry对象
func (s *FileSearcher) outputResult(path string, d os.DirEntry) {
	s.config.MatchCount.Add(1)

	if s.config.Count {
		return
	}

	if s.config.FullPath {
		fullPath, pathErr := filepath.Abs(path)
		if pathErr != nil {
			fullPath = path
		}
		printPathColor(fullPath, s.config.Cl, d, true)
	} else {
		printPathColor(path, s.config.Cl, d, true)
	}
}

// isSymlinkLoop 检查符号链接是否存在循环
//
// 参数:
//   - path: 符号链接的路径
//
// 返回:
//   - bool: 如果符号链接存在循环, 则返回true; 否则返回false
func (s *FileSearcher) isSymlinkLoop(path string) bool {
	maxDepth := s.config.MaxDepthLimit
	visited := make(map[string]bool)
	currentPath := filepath.Clean(path)

	for depth := 0; depth < maxDepth; depth++ {
		// 检查是否已访问过当前路径
		if visited[currentPath] {
			return true
		}
		visited[currentPath] = true

		// 获取文件信息
		info, err := os.Lstat(currentPath)
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			return false
		}

		// 解析符号链接
		newPath, err := os.Readlink(currentPath)
		if err != nil {
			return false
		}

		// 处理相对路径
		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(filepath.Dir(currentPath), newPath)
		}
		currentPath = filepath.Clean(newPath)
	}

	return false // 达到最大深度仍未发现循环
}
