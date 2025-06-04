package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
)

func listCmdMain(cl *colorlib.ColorLib, cmd *flag.FlagSet) error {
	// 获取命令行参数
	listPath := cmd.Arg(0)

	// 如果没有指定路径, 则默认为当前目录
	if listPath == "" {
		listPath = "."
	}

	// 检查list命令的参数是否合法
	if err := checkListCmdArgs(); err != nil {
		return err
	}

	// 获取文件信息切片
	listInfos, getErr := getFileInfos(listPath)
	if getErr != nil {
		return fmt.Errorf("获取文件信息时发生了错误: %v", getErr)
	}

	// 检查切片是否为空
	if len(listInfos) == 0 {
		return nil
	}

	// 根据命令行参数排序文件信息切片
	if *listCmdSortByTime && !*listCmdReverseSort {
		// -t 为 true, -r 为 true, 则按文件修改时间降序排序
		listInfos.SortByModTimeDesc()
	} else if *listCmdSortByTime && *listCmdReverseSort {
		// -t 为 true, -r 为 false, 则按文件修改时间升序排序
		listInfos.SortByModTimeAsc()
	} else if *listCmdSortBySize && !*listCmdReverseSort {
		// -s 为 true, -r 为 true, 则按文件大小降序排序
		listInfos.SortByFileSizeDesc()
	} else if *listCmdSortBySize && *listCmdReverseSort {
		// -s 为 true, -r 为 false, 则按文件大小升序排序
		listInfos.SortByFileSizeAsc()
	} else if *listCmdSortByName && *listCmdReverseSort {
		// -n 为 true, -r 为 true, 则按文件名升序排序
		listInfos.SortByFileNameAsc()
	} else {
		// 默认按文件名降序排序
		listInfos.SortByFileNameDesc()
	}

	// 如果启用了长格式输出, 则调用 listCmdLong 函数
	if *listCmdLongFormat {
		if err := listCmdLong(cl, listInfos); err != nil {
			return err
		}
		return nil
	}

	// list命令的默认运行函数
	if defaultErr := listCmdDefault(cl, listInfos); defaultErr != nil {
		return defaultErr
	}

	return nil
}

// checkListCmdArgs 检查list命令的参数是否合法
func checkListCmdArgs() error {
	// 检查是否同时指定了 -s 和 -t 选项
	if *listCmdSortBySize && *listCmdSortByTime {
		return errors.New("不能同时指定 -s 和 -t 选项")
	}

	// 检查是否同时指定了 -s 和 -n 选项
	if *listCmdSortBySize && *listCmdSortByName {
		return errors.New("不能同时指定 -s 和 -n 选项")
	}

	// 检查是否同时指定了 -t 和 -n 选项
	if *listCmdSortByTime && *listCmdSortByName {
		return errors.New("不能同时指定 -t 和 -n 选项")
	}

	// 检查是否同时指定了 -f 和 -d 选项
	if *listCmdFileOnly && *listCmdDirOnly {
		return errors.New("不能同时指定 -f 和 -d 选项")
	}

	// 如果指定了-ho检查是否指定-a
	if *listCmdHiddenOnly && !*listCmdAll {
		return errors.New("必须指定 -a 选项才能使用 -ho 选项")
	}

	return nil
}

// list命令的默认运行函数
func listCmdDefault(cl *colorlib.ColorLib, lfs globals.ListInfos) error {
	// 获取终端宽度
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("获取终端宽度时发生了错误: %v", err)
	}

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 动态计算每行可以容纳的列数
	fileNames := make([]string, len(lfs))
	for i, info := range lfs {
		fileNames[i] = info.Name
	}
	maxWidth := text.LongestLineLen(strings.Join(fileNames, "\n"))
	columns := width / (maxWidth + 2) // 每列增加4个字符的间距
	if columns == 0 {
		columns = 1 // 至少显示一列
	}

	// 设置表格样式
	if *listCmdTheme != "" {
		// 根据主题设置表格样式
		switch *listCmdTheme {
		case "dark":
			t.SetStyle(table.StyleColoredDark)
		case "light":
			t.SetStyle(table.StyleColoredBright)
		case "d":
			t.SetStyle(table.StyleColoredDark)
		case "l":
			t.SetStyle(table.StyleColoredBright)
		}
	}

	// 设置表格Name列的对齐方式为左对齐
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Name", Align: text.AlignLeft},
	})

	// 构建多列输出
	for i := 0; i < len(lfs); i += columns {
		// 计算当前行要显示的文件数
		end := i + columns
		if end > len(lfs) {
			end = len(lfs)
		}

		// 创建一个临时表格行
		row := make([]interface{}, end-i)
		for j := i; j < end; j++ {
			var paddedFilename string
			if *listCmdQuoteNames {
				paddedFilename = fmt.Sprintf("%q", lfs[j].Name)
			} else {
				paddedFilename = lfs[j].Name
			}

			// 检查是否启用颜色输出
			if *listCmdColor {
				paddedFilename = getColorString(lfs[j], paddedFilename, cl)
			}

			row[j-i] = paddedFilename
		}

		// 添加行到表格
		t.AppendRow(row)
	}

	// 输出表格
	t.Render()

	return nil
}

// listCmdLong 函数用于以长格式输出文件信息。
func listCmdLong(cl *colorlib.ColorLib, ifs globals.ListInfos) error {
	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	if *listCmdShowUserGroup {
		t.AppendHeader(table.Row{"Type", "Perm", "Owner", "Group", "Size", "Unit", "ModTime", "Name"})
	} else {
		t.AppendHeader(table.Row{"Type", "Perm", "Size", "Unit", "ModTime", "Name"})
	}

	for _, info := range ifs {
		// 获取文件权限字符串
		infoPerm := FormatPermissionString(cl, info)

		// 启用颜色输出
		if *listCmdColor {
			// 类型
			infoType := getColorString(info, info.EntryType, cl)

			// 文件大小 和 单位
			infoSize, infoSizeUnit := humanSize(info.Size)

			// 渲染文件大小
			infoSize = cl.Syellow(infoSize)

			// 渲染文件大小单位
			infoSizeUnit = cl.Syellow(infoSizeUnit)

			// 修改时间
			infoModTime := cl.Sgreen(info.ModTime.Format("2006-01-02 15:04:05"))

			// 文件名
			var infoName string
			// 检查是否启用引号
			if *listCmdQuoteNames {
				infoName = getColorString(info, fmt.Sprintf("%q", info.Name), cl)
			} else {
				infoName = getColorString(info, info.Name, cl)
			}

			// 添加行到表格
			if *listCmdShowUserGroup {
				t.AppendRow(table.Row{infoType, infoPerm, info.Owner, info.Group, infoSize, infoSizeUnit, infoModTime, infoName})
			} else {
				t.AppendRow(table.Row{infoType, infoPerm, infoSize, infoSizeUnit, infoModTime, infoName})
			}
		} else {
			// 不启用颜色输出 (默认)
			infoSize, infoSizeUnit := humanSize(info.Size)
			if *listCmdShowUserGroup {
				// 判断是否启用引号
				if *listCmdQuoteNames {
					t.AppendRow(table.Row{info.EntryType, infoPerm, info.Owner, info.Group, infoSize, infoSizeUnit, info.ModTime.Format("2006-01-02 15:04:05"), fmt.Sprintf("%q", info.Name)})
				} else {
					t.AppendRow(table.Row{info.EntryType, infoPerm, info.Owner, info.Group, infoSize, infoSizeUnit, info.ModTime.Format("2006-01-02 15:04:05"), info.Name})
				}
			} else {
				// 判断是否启用引号
				if *listCmdQuoteNames {
					t.AppendRow(table.Row{info.EntryType, infoPerm, infoSize, infoSizeUnit, info.ModTime.Format("2006-01-02 15:04:05"), fmt.Sprintf("%q", info.Name)})
				} else {
					t.AppendRow(table.Row{info.EntryType, infoPerm, infoSize, infoSizeUnit, info.ModTime.Format("2006-01-02 15:04:05"), info.Name})
				}
			}
		}
	}

	// 设置列的对齐方式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Size", Align: text.AlignRight},
	})

	// 设置表格样式
	if *listCmdTheme != "" {
		// 根据主题设置表格样式
		switch *listCmdTheme {
		case "dark":
			t.SetStyle(table.StyleColoredDark)
		case "light":
			t.SetStyle(table.StyleColoredBright)
		case "d":
			t.SetStyle(table.StyleColoredDark)
		case "l":
			t.SetStyle(table.StyleColoredBright)
		}
	}

	// 输出表格
	t.Render()

	return nil
}

// FormatPermissionString 根据颜色模式格式化权限字符串
func FormatPermissionString(cl *colorlib.ColorLib, info globals.ListInfo) (formattedPerm string) {
	if *listCmdColor {
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
		// 如果不启用颜色输出，直接返回权限字符串（后9位）
		formattedPerm = info.Perm[1:]
	}

	return formattedPerm
}

// getFileInfos 函数用于获取指定路径下所有文件和目录的详细信息。
// 参数 path 表示要获取信息的目录路径。
// 返回值为一个 globals.ListInfo 类型的切片，包含了每个文件和目录的详细信息，以及可能出现的错误。
func getFileInfos(p string) (globals.ListInfos, error) {
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
		if *listCmdDirEctory {
			if isSystemFileOrDir(filepath.Base(absPath)) {
				return nil, fmt.Errorf("不能列出系统文件或目录: %s", absPath)
			}

			// 检查文件是否应该被跳过
			if shouldSkipFile(absPath, pathInfo.IsDir(), pathInfo, false) {
				infos = make(globals.ListInfos, 0)
				return infos, nil
			}

			// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
			info := buildFileInfo(pathInfo, absPath)

			// 将构建好的目录信息添加到切片中
			infos = append(infos, info)

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
			if *listCmdRecursion && file.IsDir() {
				// 使用goroutine并行处理子目录
				type result struct {
					infos globals.ListInfos
					err   error
				}
				resultChan := make(chan result)

				go func() {
					subInfos, err := getFileInfos(absFilePath)
					resultChan <- result{subInfos, err}
				}()

				res := <-resultChan
				if res.err != nil {
					return nil, fmt.Errorf("递归处理目录 %s 时出错: %v", absFilePath, res.err)
				}

				// 将子目录中的文件信息添加到切片中
				infos = append(infos, res.infos...)
			}

			// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
			info := buildFileInfo(fileInfo, absFilePath)

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
		info := buildFileInfo(pathInfo, absPath)

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
//	  "l" - 符号链接
//	  "s" - 套接字
//	  "p" - 命名管道
//	  "b" - 块设备
//	  "c" - 字符设备
//	  "x" - 可执行文件
//	  "e" - 空文件
//	  "?" - 未知类型
func getEntryType(fileInfo os.FileInfo) string {
	mode := fileInfo.Mode()
	switch {
	case mode.IsDir():
		// 如果是目录，条目类型标记为 'd' 或 'dir' 或 'directory'
		return "d"
	case mode.IsRegular():
		// 如果是普通文件，检查条目类型标记可能为 'f' 或 'file', 'x' 或 'executable'
		if mode&0111 != 0 {
			return "x" // 可执行文件
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(fileInfo.Name()))
		switch ext {
		case ".exe":
			return "x" // 可执行文件
		default:
			return "f" // 普通文件
		}
	case mode&os.ModeSymlink != 0:
		// 如果是符号链接，条目类型标记为 'l' 或 'symlink'
		return "l"
	case mode&os.ModeSocket != 0:
		// 如果是套接字，条目类型标记为 's' 或 'socket'
		return "s"
	case mode&os.ModeNamedPipe != 0:
		// 如果是命名管道 (FIFO)，条目类型标记为 'p' 或 'pipe'
		return "p"
	case mode&os.ModeDevice != 0:
		// 如果是块设备，条目类型标记为 'b' 或 'block-device'
		return "b"
	case mode&os.ModeCharDevice != 0:
		// 如果是字符设备，条目类型标记为 'c' 或 'char-device'
		return "c"
	case fileInfo.Size() == 0:
		// 如果是空文件或目录，条目类型标记为 'e' 或 'empty'
		return "e"
	default:
		// 其他类型，条目类型标记为 '?'
		return "?"
	}
}

// shouldSkipFile 函数用于依据命令行参数和文件属性，判断是否需要跳过指定的文件或目录。
// 该函数会综合考虑多个命令行选项，如显示所有文件（-a）、仅显示隐藏文件（-ho）、
// 仅显示文件（-f）和仅显示目录（-d）等，结合文件或目录的名称及类型，做出跳过与否的决策。
// 参数:
// - p: 表示文件或目录的名称，用于判断是否为隐藏文件。
// - isDir: 布尔值，用于标识当前条目是否为目录。
// 返回值:
// - 布尔类型，若满足跳过条件则返回 true，否则返回 false。
func shouldSkipFile(p string, isDir bool, fileInfo os.FileInfo, main bool) bool {
	// 场景 1: 未启用 -a 选项，且当前文件或目录为隐藏文件时，应跳过该条目
	// -a 选项用于显示所有文件，若未启用该选项，隐藏文件默认不显示
	if !*listCmdAll && isHidden(p) {
		return true
	}

	// 场景 2: 同时启用 -a 和 -ho 选项，但当前文件或目录并非隐藏文件时，应跳过该条目
	// -a 选项用于显示所有文件，-ho 选项用于仅显示隐藏文件，两者同时启用时，非隐藏文件需跳过
	if *listCmdAll && *listCmdHiddenOnly && !isHidden(p) {
		return true
	}

	// 场景 3: 启用 -f 选项，且当前条目为目录时，应跳过该条目
	// -f 选项用于仅显示文件，若当前条目为目录，则不符合要求，需跳过
	if !main {
		if *listCmdFileOnly && isDir {
			return true
		}
	}

	// 场景 4: 启用 -d 选项，且当前条目不是目录时，应跳过该条目
	// -d 选项用于仅显示目录，若当前条目不是目录，则不符合要求，需跳过
	if !main {
		if *listCmdDirOnly && !isDir {
			return true
		}
	}

	// 场景 5: 启用 -L 选项，且当前条目不是软链接时，应跳过该条目
	// 如果设置了-L标志且当前不是软链接, 则跳过
	if !main {
		if *listCmdSymlink && (fileInfo.Mode()&os.ModeSymlink == 0) {
			return true
		}
	}

	// 场景 6: 启用 -ro 选项，且当前文件不是只读时，应跳过该条目
	// 如果设置了-ro标志且当前文件不是只读, 则跳过
	if !main {
		if *listCmdReadOnly && !isReadOnly(p) {
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
func buildFileInfo(fileInfo os.FileInfo, absPath string) globals.ListInfo {
	// 从文件的绝对路径中提取出文件名，作为文件的显示名称
	baseName := filepath.Base(absPath)

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

	// 检查文件名中是否包含点号（.），如果包含则尝试提取文件的扩展名
	if strings.Contains(baseName, ".") {
		// 调用 filepath.Ext 函数从文件名中提取文件的扩展名
		fileExt = filepath.Ext(baseName)
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

	// 先将转换后的大小和单位拼接成一个字符串
	sizeF := fmt.Sprintf("%.0f", sizeFloat)

	// 如果转换后的大小为 0, 则返回 "0B"
	if sizeF == "0.00" {
		return "0", "B"
	}

	// 如果转换后的大小为 0, 则返回 "0"
	if sizeF == "0" {
		return "0", "B"
	}

	// 去除小数部分末尾的 .00 或 .0
	sizeF = strings.TrimSuffix(sizeF, ".00")
	sizeF = strings.TrimSuffix(sizeF, ".0")

	// 去除小数点部分末尾的0
	if strings.Contains(sizeF, ".") {
		sizeF = strings.TrimRight(sizeF, "0")
	}

	return sizeF, unit
}
