package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
	"golang.org/x/term"
)

func listCmdMain(cl *colorlib.ColorLib) error {
	// 获取所有命令行参数
	paths := listCmd.Args()

	// 如果没有指定路径, 则默认为当前目录
	if len(paths) == 0 {
		if runtime.GOOS == "windows" {
			// 获取当前目录
			dir, pwdErr := os.Getwd()
			if pwdErr != nil {
				// 如果获取当前目录失败，则使用 "."
				paths = []string{"."}
			} else {
				// 如果获取当前目录成功，则使用当前目录
				paths = []string{dir}
			}
		} else {
			// 非Windows系统下，使用 "."
			paths = []string{"."}
		}
	}

	// 检查list命令的参数是否合法
	if err := checkListCmdArgs(); err != nil {
		return err
	}

	// 根据listCmdColor.Get()选项启用颜色输出
	if listCmdColor.Get() {
		cl.NoColor.Store(false)
	} else {
		cl.NoColor.Store(true)
	}

	// 展开通配符并收集所有路径
	var expandedPaths []string
	for _, path := range paths {
		// 清理路径
		path = filepath.Clean(path)

		// 处理路径中的通配符
		matches, err := filepath.Glob(path)
		if err != nil {
			cl.PrintErrf("路径模式错误 %q: %v\n", path, err)
			continue
		}

		// 如果路径模式没有匹配任何文件，则打印错误信息
		if len(matches) == 0 {
			cl.PrintWarnf("该路径下为空或不是一个有效路径: %s\n", path)
			continue
		}

		// 过滤掉应该跳过的文件
		for _, match := range matches {
			// 如果 -a=true ，则显示所有文件，包括隐藏文件, 如果 -a=false ，则仅显示非隐藏文件
			if listCmdAll.Get() || !isHidden(match) {
				expandedPaths = append(expandedPaths, match)
			}
		}
	}

	// 去重路径
	seen := make(map[string]bool)
	uniquePaths := []string{} // 用于存储去重后的路径
	for _, p := range expandedPaths {
		if !seen[p] {
			seen[p] = true
			uniquePaths = append(uniquePaths, p)
		}
	}

	// 分离处理文件和目录
	var fileInfos globals.ListInfos
	for _, path := range uniquePaths {
		pathInfo, statErr := os.Stat(path)
		if statErr != nil {
			cl.PrintErrf("获取路径信息失败 %q: %v\n", path, statErr)
			continue
		}

		// 收集文件信息 指定路径和root路径
		mu := &sync.Mutex{}
		infos, getErr := getFileInfos(path, path, mu)
		if getErr != nil {
			cl.PrintErrf("获取文件信息失败 %q: %v\n", path, getErr)
			continue
		}

		if pathInfo.IsDir() {
			// 根据命令行参数排序目录中的文件信息
			sortFileInfos(infos, listCmdSortByTime.Get(), listCmdSortBySize.Get(), listCmdSortByName.Get(), listCmdReverseSort.Get())

			// 只在处理多个项目时打印目录路径
			if len(uniquePaths) > 1 {
				if listCmdLongFormat.Get() {
					// 获取绝对路径
					absPath, absErr := filepath.Abs(path)
					if absErr != nil {
						absPath = path
					}

					// 打印目录标题
					cl.Bluef("%s: \n", absPath)
				} else {
					// 打印目录标题
					cl.Bluef("%s: \n", filepath.Base(path))
				}
			}

			// 处理目录，单独生成表格
			if listCmdLongFormat.Get() {
				if err := listCmdLong(cl, infos); err != nil {
					cl.PrintErrf("处理目录 %q 时出错: %v\n", path, err)
				}
			} else {
				if err := listCmdDefault(cl, infos); err != nil {
					cl.PrintErrf("处理目录 %q 时出错: %v\n", path, err)
				}
			}
			// 只在处理多个项目时打印空行分隔
			if len(uniquePaths) > 1 {
				fmt.Println()
			}
		} else {
			// 收集文件信息
			fileInfos = append(fileInfos, infos...)
		}
	}

	// 处理所有文件
	if len(fileInfos) > 0 {
		// 打印文件组标题（当前目录）
		currentDir, err := os.Getwd()
		if err != nil {
			currentDir = "."
		}
		cl.Bluef("%s: \n", currentDir)

		// 根据命令行参数排序文件信息切片
		sortFileInfos(fileInfos, listCmdSortByTime.Get(), listCmdSortBySize.Get(), listCmdSortByName.Get(), listCmdReverseSort.Get())

		if listCmdLongFormat.Get() {
			if err := listCmdLong(cl, fileInfos); err != nil {
				return err
			}
		} else {
			if err := listCmdDefault(cl, fileInfos); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkListCmdArgs 检查list命令的参数是否合法
func checkListCmdArgs() error {
	// 检查是否同时指定了 -s 和 -t 选项
	if listCmdSortBySize.Get() && listCmdSortByTime.Get() {
		return errors.New("不能同时指定 -s 和 -t 选项")
	}

	// 检查是否同时指定了 -s 和 -n 选项
	if listCmdSortBySize.Get() && listCmdSortByName.Get() {
		return errors.New("不能同时指定 -s 和 -n 选项")
	}

	// 检查是否同时指定了 -t 和 -n 选项
	if listCmdSortByTime.Get() && listCmdSortByName.Get() {
		return errors.New("不能同时指定 -t 和 -n 选项")
	}

	// 检查是否同时指定了 -f 和 -d 选项
	if listCmdFileOnly.Get() && listCmdDirOnly.Get() {
		return errors.New("不能同时指定 -f 和 -d 选项")
	}

	// 如果指定了-ho检查是否指定-a
	if listCmdHiddenOnly.Get() && !listCmdAll.Get() {
		return errors.New("必须指定 -a 选项才能使用 -ho 选项")
	}

	// 检查是否-ts的表格样式是否为合法值
	if listCmdTableStyle.Get() != "" {
		if _, ok := globals.TableStyleMap[listCmdTableStyle.Get()]; !ok {
			return fmt.Errorf("无效的表格样式: %s", listCmdTableStyle.Get())
		}
	}

	// 检查是否同时指定了 -c 和 --dev-color
	if listCmdDevColor.Get() && !listCmdColor.Get() {
		return fmt.Errorf("如果要使用 -%s, 必须要先启用 -%s", listCmdDevColor.ShortName(), listCmdColor.ShortName())
	}

	return nil
}

// list命令的默认运行函数
// listCmdDefault 函数用于以默认格式输出文件信息，支持递归目录分组显示和多列表格布局。
func listCmdDefault(cl *colorlib.ColorLib, lfs globals.ListInfos) error {
	// 获取终端宽度
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("获取终端宽度时发生了错误: %v", err)
	}

	// 如果启用了递归，按目录分组显示
	if listCmdRecursion.Get() {
		// 按目录分组文件
		dirFiles := make(map[string][]globals.ListInfo)
		for _, info := range lfs {
			dir, _ := filepath.Split(info.Name)
			if dir == "" {
				dir = "."
			}
			dirFiles[dir] = append(dirFiles[dir], info)
		}

		// 排序目录
		dirs := make([]string, 0, len(dirFiles))
		for dir := range dirFiles {
			dirs = append(dirs, dir)
		}
		sort.Strings(dirs)

		// 处理每个目录
		for _, dir := range dirs {
			infos := dirFiles[dir]
			fileNames := prepareFileNames(infos, listCmdQuoteNames.Get(), listCmdColor.Get(), cl)

			// 打印目录名
			cl.Bluef("%s:\n", dir)

			// 创建并渲染表格
			if err := renderDefaultTable(fileNames, width, listCmdTableStyle.Get()); err != nil {
				return fmt.Errorf("渲染目录 %s 表格失败: %v", dir, err)
			}
			fmt.Println() // 目录间空行
		}

		return nil
	}

	// 非递归模式，直接输出所有文件
	fileNames := prepareFileNames(lfs, listCmdQuoteNames.Get(), listCmdColor.Get(), cl)

	// 创建并渲染表格
	return renderDefaultTable(fileNames, width, listCmdTableStyle.Get())
}

// listCmdLong 函数用于以长格式输出文件信息，支持递归目录分组显示。
func listCmdLong(cl *colorlib.ColorLib, ifs globals.ListInfos) error {
	// 如果启用了递归，按目录分组显示
	if listCmdRecursion.Get() {
		// 按目录分组文件
		dirFiles := make(map[string][]globals.ListInfo)
		for _, info := range ifs {
			dir, _ := filepath.Split(info.Name)
			if dir == "" {
				dir = "."
			}
			dirFiles[dir] = append(dirFiles[dir], info)
		}

		// 排序目录
		dirs := make([]string, 0, len(dirFiles))
		for dir := range dirFiles {
			dirs = append(dirs, dir)
		}
		sort.Strings(dirs)

		// 处理每个目录
		for _, dir := range dirs {
			infos := dirFiles[dir]

			// 打印目录名
			cl.Bluef("%s:\n", dir)

			// 创建并渲染表格
			if err := renderFileTable(cl, infos); err != nil {
				return fmt.Errorf("渲染目录 %s 表格失败: %v", dir, err)
			}
			fmt.Println() // 目录间空行
		}
		return nil
	}

	// 非递归模式，直接渲染表格
	return renderFileTable(cl, ifs)
}

// FormatPermissionString 根据颜色模式格式化权限字符串
func FormatPermissionString(cl *colorlib.ColorLib, info globals.ListInfo) (formattedPerm string) {
	// 检查权限字符串长度是否有效 (至少10个字符, 如 "-rwxr-xr-x")
	if len(info.Perm) < 10 {
		return "???-???-???"
	}

	if listCmdColor.Get() {
		for i := 1; i < len(info.Perm); i++ { // 跳过第一个字符（通常是文件类型标志）
			colorName := PermissionColorMap[i] // 获取当前字符的颜色名称
			switch colorName {
			case "green":
				formattedPerm += cl.Sgreen(string(info.Perm[i]))
			case "yellow":
				formattedPerm += cl.Syellow(string(info.Perm[i]))
			case "red":
				formattedPerm += cl.Sred(string(info.Perm[i]))
			default:
				formattedPerm += string(info.Perm[i]) // 如果颜色未定义，直接添加
			}
		}
	} else {
		// 如果不启用颜色输出, 直接返回权限字符串(后9位)
		formattedPerm = info.Perm[1:]
	}

	return formattedPerm
}

// getFileInfos 函数用于获取指定路径下所有文件和目录的详细信息。
// 参数 path 表示要获取信息的目录路径。
// 返回值为一个 globals.ListInfo 类型的切片，包含了每个文件和目录的详细信息，以及可能出现的错误。
func getFileInfos(p string, rootDir string, mu *sync.Mutex) (globals.ListInfos, error) {
	// 初始化一个用于存储文件信息的切片
	var infos globals.ListInfos

	// 清理路径
	p = filepath.Clean(p)

	// 获取绝对路径
	absPath, absErr := filepath.Abs(p)
	if absErr != nil {
		return nil, fmt.Errorf("获取文件 %s 的绝对路径时发生了错误: %v", p, absErr)
	}

	// 检查是否为系统文件或目录
	if isSystemFileOrDir(filepath.Base(absPath)) {
		return nil, fmt.Errorf("不能列出系统文件或目录: %s", absPath)
	}

	// 获取指定路径的文件信息
	pathInfo, statErr := os.Stat(absPath)
	if statErr != nil {
		return nil, handleError(absPath, statErr)
	}

	// 检查初始路径是否应该被跳过
	if shouldSkipFile(absPath, pathInfo.IsDir(), pathInfo, true) {
		infos = make(globals.ListInfos, 0)
		return infos, nil
	}

	// 根据是否为目录进行处理
	if pathInfo.IsDir() {
		// 如果设置了 -D 标志，只列出目录本身
		if listCmdDirEctory.Get() {
			if isSystemFileOrDir(filepath.Base(absPath)) {
				return nil, fmt.Errorf("不能列出系统文件或目录: %s", absPath)
			}

			// 检查文件是否应该被跳过
			if shouldSkipFile(absPath, pathInfo.IsDir(), pathInfo, false) {
				infos = make(globals.ListInfos, 0)
				return infos, nil
			}

			// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
			// 由于 buildFileInfo 函数需要三个参数，这里添加 rootDir 参数，由于当前处理的是目录本身，rootDir 可以设为 absPath
			info := buildFileInfo(pathInfo, absPath, absPath)

			// 使用互斥锁保护切片操作
			mu.Lock()
			infos = append(infos, info)
			mu.Unlock()

			return infos, nil
		}

		// 读取目录下的文件
		files, readDirErr := os.ReadDir(absPath)
		if readDirErr != nil {
			return nil, handleError(absPath, readDirErr)
		}

		// 遍历目录下的每个文件和目录
		for _, file := range files {
			absFilePath, absErr := filepath.Abs(filepath.Join(absPath, file.Name()))
			if absErr != nil {
				return nil, fmt.Errorf("获取文件 %s 的绝对路径时发生了错误: %v", file.Name(), absErr)
			}

			// 检查是否为系统文件或目录
			if isSystemFileOrDir(filepath.Base(absFilePath)) {
				continue
			}

			// 获取文件的详细信息，如大小、修改时间等
			fileInfo, statErr := file.Info()
			if statErr != nil {
				return nil, handleError(absFilePath, statErr)
			}

			// 检查文件是否应该被跳过
			if shouldSkipFile(absFilePath, file.IsDir(), fileInfo, false) {
				continue
			}

			// 如果设置了-R标志且当前是目录, 则递归处理子目录
			if listCmdRecursion.Get() && file.IsDir() {
				// 使用带并发控制的goroutine处理子目录
				type result struct {
					infos globals.ListInfos
					err   error
				}
				resultChan := make(chan result)

				// 获取CPU核心数并设置最大并发数为核心数*2
				maxConcurrency := runtime.NumCPU() * 2
				sem := make(chan struct{}, maxConcurrency)

				go func() {
					sem <- struct{}{}        // 获取信号量
					defer func() { <-sem }() // 释放信号量

					subInfos, err := getFileInfos(absFilePath, rootDir, mu)
					resultChan <- result{subInfos, err}
				}()

				res := <-resultChan
				if res.err != nil {
					return nil, fmt.Errorf("递归处理目录 %s 时出错: %v", absFilePath, res.err)
				}

				// 使用互斥锁保护切片操作
				mu.Lock()
				infos = append(infos, res.infos...)
				mu.Unlock()
			}

			// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
			info := buildFileInfo(fileInfo, absFilePath, rootDir)

			// 将构建好的目录信息添加到切片中
			infos = append(infos, info)
		}
	} else {
		// 检查是否为系统文件或目录
		if isSystemFileOrDir(filepath.Base(absPath)) {
			return nil, fmt.Errorf("不能列出系统文件或目录: %s", absPath)
		}

		// 检查文件是否应该被跳过
		if shouldSkipFile(absPath, pathInfo.IsDir(), pathInfo, false) {
			infos = make(globals.ListInfos, 0)
			return infos, nil
		}

		// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
		// 由于 buildFileInfo 函数需要三个参数，这里第三个参数使用 absPath 作为 rootDir
		info := buildFileInfo(pathInfo, absPath, absPath)

		// 将构建好的目录信息添加到切片中
		infos = append(infos, info)
	}

	// 返回存储文件信息的切片和可能出现的错误
	return infos, nil
}

// getEntryType 根据文件模式返回对应的类型标识符
// 参数:
//
//	fileInfo - 文件信息对象
//
// 返回值:
//
//	string - 文件类型标识符:
//	  "d" - 目录
//	  "f" - 普通文件
//	  "l" - 符号链接(软链接/快捷方式)
//	  "s" - 套接字
//	  "p" - 命名管道
//	  "b" - 块设备
//	  "c" - 字符设备
//	  "x" - 可执行文件
//	  "e" - 空文件
//	  "?" - 未知类型
func getEntryType(fileInfo os.FileInfo) string {
	// 获取文件模式
	mode := fileInfo.Mode()

	// 检查是否是符号链接
	if mode&os.ModeSymlink != 0 {
		return globals.SymlinkType
	}

	// 检查是否是目录
	if mode.IsDir() {
		return globals.DirType
	}

	// 检查是否是套接字
	if mode&os.ModeSocket != 0 {
		return globals.SocketType
	}

	// 检查是否是命名管道
	if mode&os.ModeNamedPipe != 0 {
		return globals.PipeType
	}

	// 块设备和字符设备
	if mode&os.ModeDevice != 0 {
		if mode&os.ModeCharDevice != 0 {
			return globals.CharDeviceType
		}
		return globals.BlockDeviceType
	}

	// 检查是否是普通文件
	if mode.IsRegular() {
		// 检查是否是空文件
		if fileInfo.Size() == 0 {
			return globals.EmptyType
		}

		// 平台特定的可执行文件判断
		switch runtime.GOOS {
		case "windows":
			// Windows可执行文件扩展名
			ext := strings.ToLower(filepath.Ext(fileInfo.Name()))
			switch ext {
			case ".exe", ".com", ".cmd", ".bat", ".ps1", ".psm1":
				return globals.ExecutableType
			case ".lnk", ".url":
				return globals.SymlinkType
			}
		case "linux", "darwin":
			// Unix-like系统可执行文件判断
			if mode&0111 != 0 {
				return globals.ExecutableType
			}
		}

		// 默认为普通文件
		return globals.FileType
	}

	// 其他类型
	return globals.UnknownType
}

// shouldSkipFile 函数用于依据命令行参数和文件属性，判断是否需要跳过指定的文件或目录。
// 该函数会综合考虑多个命令行选项，如显示所有文件（-a）、仅显示隐藏文件（-ho）、
// 仅显示文件（-f）和仅显示目录（-d）等，结合文件或目录的名称及类型，做出跳过与否的决策。
// 参数:
// - p: 表示文件或目录的名称，用于判断是否为隐藏文件。
// - isDir: 布尔值，用于标识当前条目是否为目录。
// - fileInfo: os.FileInfo 类型，包含文件或目录的详细信息，如权限、大小等。
// - main: 布尔值，用于标识当前是否为处理目录本身，而非其子目录。
// 返回值:
// - 布尔类型，若满足跳过条件则返回 true，否则返回 false。
func shouldSkipFile(p string, isDir bool, fileInfo os.FileInfo, main bool) bool {
	// 场景 1: 未启用 -a 选项，且当前文件或目录为隐藏文件时，应跳过该条目
	// -a 选项用于显示所有文件，若未启用该选项，隐藏文件默认不显示
	if !listCmdAll.Get() && isHidden(p) {
		return true
	}

	// 场景 2: 同时启用 -a 和 -ho 选项，但当前文件或目录并非隐藏文件时，应跳过该条目
	// -a 选项用于显示所有文件，-ho 选项用于仅显示隐藏文件，两者同时启用时，非隐藏文件需跳过
	if listCmdAll.Get() && listCmdHiddenOnly.Get() && !isHidden(p) {
		return true
	}

	// 场景 3: 启用 -f 选项，且当前条目为目录时，且未启用仅显示目录选项，应跳过该条目
	// -f 选项用于仅显示文件，若当前条目为目录，则不符合要求，需跳过
	if !main {
		if listCmdFileOnly.Get() && isDir && !listCmdDirEctory.Get() {
			return true
		}
	}

	// 场景 4: 启用 -d 选项，且当前条目不是目录时，且未启用仅显示文件选项，应跳过该条目
	// -d 选项用于仅显示目录，若当前条目不是目录，则不符合要求，需跳过
	if !main {
		if listCmdDirOnly.Get() && !isDir && !listCmdDirOnly.Get() {
			return true
		}
	}

	// 场景 5: 启用 -L 选项，且当前条目不是软链接时，应跳过该条目
	// 如果设置了-L标志且当前不是软链接, 则跳过
	if !main {
		if listCmdSymlink.Get() && (fileInfo.Mode()&os.ModeSymlink == 0) {
			return true
		}
	}

	// 场景 6: 启用 -ro 选项，且当前文件不是只读时，应跳过该条目
	// 如果设置了-ro标志且当前文件不是只读, 则跳过
	if !main {
		if listCmdReadOnly.Get() && !isReadOnly(p) {
			return true
		}
	}

	// 若不满足上述任何跳过条件，则不跳过该条目
	return false
}

// handleError 函数的作用是针对检查指定路径时出现的错误进行处理，会根据不同的错误类型生成对应的错误提示信息。
// 参数 path 代表当前正在检查的路径，该路径可以是文件路径或者目录路径。
// 参数 err 是在检查路径过程中产生的错误对象，函数会依据这个错误对象的类型来决定返回何种错误信息。
// 返回值为一个新的错误对象，其中包含了更具描述性的错误信息，方便调用者定位和处理问题。
func handleError(path string, err error) error {
	// 可以增加对os.ErrInvalid的错误处理
	if errors.Is(err, os.ErrInvalid) {
		return fmt.Errorf("路径 %s 包含无效字符: %v", path, err)
	}

	// 检查错误是否为权限错误(os.ErrPermission)，若是，则返回包含路径信息和原错误信息的权限错误提示。
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("检查路径 %s 时发生了权限错误: %v", path, err)
	}
	// 检查错误是否为路径不存在错误(os.ErrNotExist)，若是，则返回表明该目录不存在的错误提示。
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("目录 %s 不存在", path)
	}
	// 若错误既不是权限错误也不是路径不存在错误，返回一个通用的错误提示，包含路径信息和原错误信息。
	return fmt.Errorf("检查路径 %s 时发生了错误: %v", path, err)
}

// buildFileInfo 函数的作用是根据传入的文件信息和文件的绝对路径，构建一个 globals.ListInfo 结构体对象。
// 该结构体对象包含了文件的各类详细信息，方便后续对文件信息进行统一管理和使用。
// 参数说明：
// - fileInfo: os.FileInfo 类型，包含了文件的基本属性，如文件大小、修改时间、权限模式等。
// - absPath: 文件的绝对路径字符串, 用于提取文件名以及获取文件所属的用户和组信息。
// 返回值：
// - 返回一个 globals.ListInfo 结构体实例，该实例封装了文件的完整信息。
func buildFileInfo(fileInfo os.FileInfo, absPath string, rootDir string) globals.ListInfo {
	// 从文件的绝对路径中提取出文件名，作为文件的显示名称
	var baseName string
	// 如果启用了递归显示，则使用相对路径作为文件名
	if listCmdRecursion.Get() {
		relPath, err := filepath.Rel(rootDir, absPath)
		if err != nil {
			baseName = absPath // 如果无法获取相对路径，则使用绝对路径作为文件名
		} else {
			baseName = relPath // 如果成功获取相对路径，则使用相对路径作为文件名
		}
	} else {
		baseName = filepath.Base(absPath) // 如果未启用递归显示，则直接使用文件名
	}

	// 调用 getEntryType 函数，根据文件的模式判断文件的条目类型，如目录、普通文件、符号链接等
	entryType := getEntryType(fileInfo)

	// 从 fileInfo 中获取文件的大小，单位为字节
	fileSize := fileInfo.Size()

	// 从 fileInfo 中获取文件的最后一次修改时间
	fileModTime := fileInfo.ModTime()

	// 从 fileInfo 的模式中提取文件的权限信息，并将其转换为字符串形式
	filePermStr := fileInfo.Mode().Perm().String()

	// 调用 getFileOwner 函数，根据文件的绝对路径获取文件所属的用户和组信息
	u, g := getFileOwner(absPath)

	// 初始化文件的扩展名变量，默认为空字符串
	var fileExt string

	// 检查文件名是否以点开头，如果是，则从第二个字符开始检查是否包含点号
	if strings.HasPrefix(baseName, ".") && len(baseName) > 1 {
		if strings.Contains(baseName[1:], ".") {
			// 调用 filepath.Ext 函数从文件名中提取文件的扩展名
			fileExt = filepath.Ext(baseName)
		}
	} else if strings.Contains(baseName, ".") {
		// 调用 filepath.Ext 函数从文件名中提取文件的扩展名
		fileExt = filepath.Ext(baseName)
	}
	/*
		- .bashrc → 不提取扩展名（fileExt为空）
		- .vimrc.bak → 提取 ".bak" 作为扩展名
		- document.txt → 提取 ".txt" 作为扩展名
		- Makefile → 不提取扩展名
	*/

	// 检查是否为符号链接，如果是，则获取符号链接指向的目标文件的信息
	var linkTargetPath string
	var linkGetErr error
	if entryType == globals.SymlinkType {
		// 获取符号链接指向的目标文件的绝对路径
		linkTargetPath, linkGetErr = os.Readlink(absPath)
		if linkGetErr != nil {
			linkTargetPath = "?"
		}
	}

	// 构建并返回一个 globals.ListInfo 结构体实例，将前面获取到的文件信息填充到结构体中
	return globals.ListInfo{
		// 文件的条目类型，如 'd' 表示目录，'f' 表示普通文件等
		EntryType: entryType,
		// 文件的显示名称
		Name: baseName,
		// 文件的大小，单位为字节
		Size: fileSize,
		// 文件的最后一次修改时间
		ModTime: fileModTime,
		// 文件的权限信息字符串
		Perm: filePermStr,
		// 文件所属的用户名称
		Owner: u,
		// 文件所属的组名称
		Group: g,
		// 文件的扩展名，如果没有则为空字符串
		FileExt: fileExt,
		// 符号链接指向的目标文件的绝对路径，如果不是符号链接则为空字符串
		LinkTargetPath: linkTargetPath,
	}
}

// sortFileInfos 根据排序参数对文件信息进行排序
// 参数:
//
//	infos - 要排序的文件信息切片
//	sortByTime - 是否按修改时间排序
//	sortBySize - 是否按文件大小排序
//	sortByName - 是否按文件名排序
//	reverse - 是否反转排序顺序
func sortFileInfos(infos globals.ListInfos, sortByTime, sortBySize, sortByName, reverse bool) {
	// 若同时指定按时间排序且不反转排序顺序，则按修改时间降序排序
	if sortByTime && !reverse {
		infos.SortByModTimeDesc()
		// 若同时指定按时间排序且反转排序顺序，则按修改时间升序排序
	} else if sortByTime && reverse {
		infos.SortByModTimeAsc()
		// 若同时指定按文件大小排序且不反转排序顺序，则按文件大小降序排序
	} else if sortBySize && !reverse {
		infos.SortByFileSizeDesc()
		// 若同时指定按文件大小排序且反转排序顺序，则按文件大小升序排序
	} else if sortBySize && reverse {
		infos.SortByFileSizeAsc()
		// 若同时指定按文件名排序且反转排序顺序，则按文件名升序排序
	} else if sortByName && reverse {
		infos.SortByFileNameAsc()
		// 默认为按文件名降序排序, 如果仅指定反转排序顺序，则按文件名升序排序
	} else if reverse {
		infos.SortByFileNameAsc()
	} else {
		// 其他情况，按文件名降序排序
		infos.SortByFileNameDesc()
	}
}

// humanSize 函数用于将字节大小转换为可读的字符串格式
// 该函数接收一个 int64 类型的字节大小参数, 返回两个字符串: 第一个字符串表示转换后的大小, 第二个字符串表示单位
func humanSize(size int64) (string, string) {
	// 定义存储字节单位的切片, 按照从小到大的顺序排列
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	// 定义字节单位之间的换算基数, 这里使用 1024 作为二进制换算标准
	base := float64(1024)

	// 用于存储最终选择的合适单位
	var unit string
	// 将传入的 int64 类型的字节大小转换为 float64 类型, 方便后续计算
	sizeFloat := float64(size)

	// 根据字节大小选择最合适的单位
	// 如果字节大小小于 1024B, 则直接使用 B 作为单位
	if sizeFloat < base {
		unit = units[0]
		// 如果字节大小小于 1024KB, 则使用 KB 作为单位, 并将字节大小除以 1024 转换为 KB
	} else if sizeFloat < base*base {
		unit = units[1]
		sizeFloat /= base
		// 如果字节大小小于 1024MB, 则使用 MB 作为单位, 并将字节大小除以 1024*1024 转换为 MB
	} else if sizeFloat < base*base*base {
		unit = units[2]
		sizeFloat /= base * base
		// 如果字节大小小于 1024GB, 则使用 GB 作为单位, 并将字节大小除以 1024*1024*1024 转换为 GB
	} else if sizeFloat < base*base*base*base {
		unit = units[3]
		sizeFloat /= base * base * base
		// 如果字节大小小于 1024TB, 则使用 TB 作为单位, 并将字节大小除以 1024*1024*1024*1024 转换为 TB
	} else if sizeFloat < base*base*base*base*base {
		unit = units[4]
		sizeFloat /= base * base * base * base
		// 否则使用 PB 作为单位, 并将字节大小除以 1024*1024*1024*1024*1024 转换为 PB
	} else {
		unit = units[5]
		sizeFloat /= base * base * base * base * base
	}

	// 根据数值大小选择合适精度
	var sizeF string
	if sizeFloat < 10 {
		sizeF = fmt.Sprintf("%.1f", sizeFloat) // 小于10时保留1位小数
	} else {
		sizeF = fmt.Sprintf("%.0f", sizeFloat) // 大于等于10时取整
	}

	// 处理特殊情况: 10.0 -> 10
	sizeF = strings.TrimSuffix(sizeF, ".0")

	// 处理0值情况
	if sizeF == "0" || sizeF == "0.0" {
		return "0", "B"
	}

	return sizeF, unit
}
