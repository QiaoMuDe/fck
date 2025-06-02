package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
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
		listInfos.SortByFileNameDesc()
	} else if *listCmdSortByTime && *listCmdReverseSort {
		// -t 为 true, -r 为 false, 则按文件修改时间升序排序
		listInfos.SortByFileNameAsc()
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

	// // 检查是否同时指定了 -f 和 -L 选项
	// if *listCmdFileOnly && *listCmdSymlink {
	// 	return errors.New("不能同时指定 -f 和 -L 选项")
	// }

	// // 检查是否同时指定了 -f 和 -ro 选项
	// if *listCmdFileOnly && *listCmdReadOnly {
	// 	return errors.New("不能同时指定 -f 和 -ro 选项")
	// }

	// // 检查是否同时指定了 -f 和 -ho 选项
	// if *listCmdFileOnly && *listCmdHiddenOnly {
	// 	return errors.New("不能同时指定 -f 和 -ho 选项")
	// }

	return nil
}

// list命令的默认运行函数
func listCmdDefault(cl *colorlib.ColorLib, lfs globals.ListInfos) error {
	// 获取终端宽度
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("获取终端宽度时发生了错误: %v", err)
	}

	// 提取文件名切片
	fileNames := lfs.GetFileNames()

	// 使用 go-pretty 库进行格式化输出
	// 计算最长文件名的长度
	maxWidth := text.LongestLineLen(strings.Join(fileNames, "\n"))

	// 动态计算每行可以容纳的列数
	columns := width / (maxWidth + 4) // 每列增加4个字符的间距 +4 是为了增加文件名之间的间距 和下面的pad函数里的maxWidth+4保持一致
	if columns == 0 {
		columns = 1 // 至少显示一列
	}

	// 构建多列输出
	for i := 0; i < len(lfs); i += columns {
		// 计算当前行要显示的文件数
		end := i + columns
		if end > len(lfs) {
			end = len(lfs)
		}

		// 打印当前行的所有文件
		for j := i; j < end; j++ {
			// 使用 go-pretty 的 Pad 函数对齐文件名
			var paddedFilename string
			if *listCmdQuoteNames {
				paddedFilename = text.Pad(fmt.Sprintf("%q", lfs[j].Name), maxWidth+4, ' ')
			} else {
				paddedFilename = text.Pad(lfs[j].Name, maxWidth+4, ' ')
			}

			// 检查是否启用颜色输出
			if *listCmdColor {
				fmt.Print(getColorString(lfs[j], paddedFilename, cl))
			} else {
				fmt.Print(paddedFilename)
			}

		}
		fmt.Println() // 换行
	}

	return nil
}

// listCmdLong 函数用于以长格式输出文件信息。
func listCmdLong(cl *colorlib.ColorLib, ifs globals.ListInfos) error {
	for _, info := range ifs {
		// 获取文件权限字符串
		infoPerm := getPermString(cl, info)

		// 启用颜色输出
		if *listCmdColor {
			// 类型
			infoType := getColorString(info, info.EntryType, cl)

			// 文件大小
			infoSize := cl.Syellow(humanReadableSize(info.Size, 1))

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

			// 根据是否显示用户组信息输出不同格式
			if *listCmdShowUserGroup {
				// 长格式输出, 格式: 类型-所有者-组-其他用户 所属用户  所属组 文件大小 修改时间 文件名
				fmt.Printf("%s%s  %-4s %-4s  %18s  %-12s  %-10s\n", infoType, infoPerm, info.Owner, info.Group, infoSize, infoModTime, infoName)
			} else {
				// 长格式输出, 格式: 类型-所有者-组-其他用户 文件大小 修改时间 文件名
				fmt.Printf("%s%s  %18s  %-12s  %-10s\n", infoType, infoPerm, infoSize, infoModTime, infoName)
			}

			continue
		}

		// 不启用颜色输出 (默认) 根据是否显示用户组信息输出不同格式
		if *listCmdShowUserGroup {
			// 判断是否启用引号
			if *listCmdQuoteNames {
				fmt.Printf("%s%s  %-4s %4s  %8s  %-12s  %-10q\n",
					info.EntryType,
					infoPerm,
					info.Owner,
					info.Group,
					humanReadableSize(info.Size, 1),
					info.ModTime.Format("2006-01-02 15:04:05"),
					info.Name,
				)
			} else {
				fmt.Printf("%s%s  %-4s %4s  %8s  %-12s  %-10s\n",
					info.EntryType,
					infoPerm,
					info.Owner,
					info.Group,
					humanReadableSize(info.Size, 1),
					info.ModTime.Format("2006-01-02 15:04:05"),
					info.Name,
				)
			}
		} else {
			// 判断是否启用引号
			if *listCmdQuoteNames {
				fmt.Printf("%s%s  %7s  %-12s  %-10q\n",
					info.EntryType,
					infoPerm,
					humanReadableSize(info.Size, 1),
					info.ModTime.Format("2006-01-02 15:04:05"),
					info.Name,
				)
			} else {
				fmt.Printf("%s%s  %7s  %-12s  %-10s\n",
					info.EntryType,
					infoPerm,
					humanReadableSize(info.Size, 1),
					info.ModTime.Format("2006-01-02 15:04:05"),
					info.Name,
				)
			}
		}
	}
	return nil
}

// getPermString 函数用于获取文件或目录的权限字符串表示。
func getPermString(cl *colorlib.ColorLib, info globals.ListInfo) (infoPerm string) {
	// 根据标志输出颜色化的权限字符串
	if *listCmdColor {
		for i := 0; i < len(info.Perm); i++ {
			switch i {
			case 0:
				continue
			case 1:
				infoPerm += cl.Syellow(string(info.Perm[i]))
			case 2:
				infoPerm += cl.Sred(string(info.Perm[i]))
			case 3:
				infoPerm += cl.Sgreen(string(info.Perm[i]))
			case 4:
				infoPerm += cl.Syellow(string(info.Perm[i]))
			case 5:
				infoPerm += cl.Sred(string(info.Perm[i]))
			case 6:
				infoPerm += cl.Sgreen(string(info.Perm[i]))
			case 7:
				infoPerm += cl.Syellow(string(info.Perm[i]))
			case 8:
				infoPerm += cl.Sred(string(info.Perm[i]))
			case 9:
				infoPerm += cl.Sgreen(string(info.Perm[i]))
			}
		}
	} else {
		// 只保留后面的9位
		infoPerm = info.Perm[1:]
	}

	return infoPerm
}

// getFileInfos 函数用于获取指定路径下所有文件和目录的详细信息。
// 参数 path 表示要获取信息的目录路径。
// 返回值为一个 globals.ListInfo 类型的切片，包含了每个文件和目录的详细信息，以及可能出现的错误。
func getFileInfos(path string) (globals.ListInfos, error) {
	// 初始化一个用于存储文件信息的切片
	var infos globals.ListInfos

	// 获取指定路径的文件信息
	pathInfo, statErr := os.Stat(path)
	if statErr != nil {
		return nil, handleError(path, statErr)
	}

	// 检查初始路径是否应该被跳过
	if shouldSkipFile(pathInfo.Name(), pathInfo.IsDir(), pathInfo, true) {
		infos = make(globals.ListInfos, 0)
		return infos, nil
	}

	// 根据是否为目录进行处理
	if pathInfo.IsDir() {
		// 如果设置了 -D 标志，只列出目录本身
		if *listCmdDirEctory {
			// 检查文件是否应该被跳过
			if shouldSkipFile(pathInfo.Name(), pathInfo.IsDir(), pathInfo, false) {
				infos = make(globals.ListInfos, 0)
				return infos, nil
			}

			// 获取目录的绝对路径
			absPath, absErr := filepath.Abs(path)
			if absErr != nil {
				return nil, fmt.Errorf("获取目录 %s 的绝对路径时发生了错误: %v", path, absErr)
			}

			// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
			info := buildFileInfo(pathInfo, absPath)

			// 将构建好的目录信息添加到切片中
			infos = append(infos, info)

			return infos, nil
		}

		// 读取目录下的文件
		files, readDirErr := os.ReadDir(path)
		if readDirErr != nil {
			return nil, handleError(path, readDirErr)
		}

		// 遍历目录下的每个文件和目录
		for _, file := range files {
			// 获取文件的详细信息，如大小、修改时间等
			fileInfo, statErr := file.Info()
			if statErr != nil {
				return nil, fmt.Errorf("获取文件 %s 的信息时发生了错误: %v", file.Name(), statErr)
			}

			// 检查文件是否应该被跳过
			if shouldSkipFile(file.Name(), file.IsDir(), fileInfo, false) {
				continue
			}

			// 如果设置了-R标志且当前是目录, 则递归处理子目录
			if *listCmdRecursion && file.IsDir() {
				// 先保留子目录本身
				subDirPath := filepath.Join(path, file.Name())
				subInfos, err := getFileInfos(subDirPath)
				if err != nil {
					return nil, fmt.Errorf("递归处理目录 %s 时出错: %v", subDirPath, err)
				}
				// 先添加子目录，再添加其内容
				infos = append(infos, subInfos...)
			}

			// 获取文件的绝对路径
			absPath, absErr := filepath.Abs(file.Name())
			if absErr != nil {
				return nil, fmt.Errorf("获取文件 %s 的绝对路径时发生了错误: %v", file.Name(), absErr)
			}

			// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
			info := buildFileInfo(fileInfo, absPath)

			// 将构建好的目录信息添加到切片中
			infos = append(infos, info)
		}
	} else {
		// 检查文件是否应该被跳过
		if shouldSkipFile(pathInfo.Name(), pathInfo.IsDir(), pathInfo, false) {
			infos = make(globals.ListInfos, 0)
			return infos, nil
		}

		// 如果 path 是一个普通文件, 获取其绝对路径
		absPath, absErr := filepath.Abs(path)
		if absErr != nil {
			return nil, fmt.Errorf("获取文件 %s 的绝对路径时发生了错误: %v", path, absErr)
		}

		// 构建一个 globals.ListInfo 结构体，存储目录的详细信息
		info := buildFileInfo(pathInfo, absPath)

		// 将构建好的目录信息添加到切片中
		infos = append(infos, info)
	}

	// 返回存储文件信息的切片和可能出现的错误
	return infos, nil
}

// 根据文件的模式判断条目类型
func getEntryType(fileInfo os.FileInfo) string {
	mode := fileInfo.Mode()
	switch {
	case mode.IsDir():
		// 如果是目录，条目类型标记为 'd' 或 'dir' 或 'directory'
		return "d"
	case mode.IsRegular():
		// 如果是普通文件，条目类型标记为 'f' 或 'file'
		return "f"
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
	case mode&0111 != 0:
		// 如果是可执行文件，条目类型标记为 'x' 或 'executable'
		// 使用 mode&0111 来判断文件是否可执行
		return "x"
	case fileInfo.Size() == 0:
		// 如果是空文件或目录，条目类型标记为 'e' 或 'empty'
		return "e"
	default:
		// 其他类型，条目类型标记为 '?'
		return "?"
	}
}

// GetUserAndGroup 获取指定 UID 和 GID 的用户名和组名
// 参数:
//
//	uid - 用户ID
//	gid - 组ID
//
// 返回:
//
//	用户名和组名， 如果在非Linux系统或解析失败时返回 "?"
func GetUserAndGroup(uid, gid string) (string, string) {
	// 检查操作系统
	if runtime.GOOS != "linux" {
		return "?", "?"
	}

	// 获取用户名
	userName := "?"
	if u, err := user.LookupId(uid); err == nil {
		userName = u.Username
	}

	// 获取组名
	groupName := "?"
	if g, err := user.LookupGroupId(gid); err == nil {
		groupName = g.Name
	}

	return userName, groupName
}

// shouldSkipFile 函数用于依据命令行参数和文件属性，判断是否需要跳过指定的文件或目录。
// 该函数会综合考虑多个命令行选项，如显示所有文件（-a）、仅显示隐藏文件（-ho）、
// 仅显示文件（-f）和仅显示目录（-d）等，结合文件或目录的名称及类型，做出跳过与否的决策。
// 参数:
// - name: 表示文件或目录的名称，用于判断是否为隐藏文件。
// - isDir: 布尔值，用于标识当前条目是否为目录。
// 返回值:
// - 布尔类型，若满足跳过条件则返回 true，否则返回 false。
func shouldSkipFile(name string, isDir bool, fileInfo os.FileInfo, main bool) bool {
	// 场景 1: 未启用 -a 选项，且当前文件或目录为隐藏文件时，应跳过该条目
	// -a 选项用于显示所有文件，若未启用该选项，隐藏文件默认不显示
	if !*listCmdAll && isHidden(name) {
		return true
	}
	// 场景 2: 同时启用 -a 和 -ho 选项，但当前文件或目录并非隐藏文件时，应跳过该条目
	// -a 选项用于显示所有文件，-ho 选项用于仅显示隐藏文件，两者同时启用时，非隐藏文件需跳过
	if *listCmdAll && *listCmdHiddenOnly && !isHidden(name) {
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
		if *listCmdReadOnly && !isReadOnly(fileInfo.Name()) {
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

// func getRecursiveFileInfos(path string) (globals.ListInfos, error) {
// 	infos, err := getFileInfos(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if *listCmdRecursion {
// 		for _, info := range infos {
// 			if info.EntryType == "d" {
// 				subPath := filepath.Join(path, info.Name)
// 				subInfos, err := getRecursiveFileInfos(subPath)
// 				if err != nil {
// 					return nil, err
// 				}
// 				infos = append(infos, subInfos...)
// 			}
// 		}
// 	}
// 	return infos, nil
// }
