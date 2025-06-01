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
	if err := checkListCmdArgs(listPath); err != nil {
		return err
	}

	// 获取文件信息切片
	listInfos, getErr := getFileInfos(listPath)
	if getErr != nil {
		return fmt.Errorf("获取文件信息时发生了错误: %v", getErr)
	}

	// 根据命令行参数排序文件信息切片
	if *listCmdSortByTime && !*listCmdReverseSort {
		// -t 为 true, -r 为 false, 则按文件修改时间升序排序
		listInfos.SortByFileNameAsc()
	} else if *listCmdSortByTime && *listCmdReverseSort {
		// -t 为 true, -r 为 true, 则按文件修改时间降序排序
		listInfos.SortByFileNameDesc()
	} else if *listCmdSortBySize && !*listCmdReverseSort {
		// -s 为 true, -r 为 false, 则按文件大小升序排序
		listInfos.SortByFileSizeAsc()
	} else if *listCmdSortBySize && *listCmdReverseSort {
		// -s 为 true, -r 为 true, 则按文件大小降序排序
		listInfos.SortByFileSizeDesc()
	} else if *listCmdSortByName && *listCmdReverseSort {
		// -n 为 true, -r 为 false, 则按文件名降序排序
		listInfos.SortByFileNameDesc()
	} else {
		// -n 为 true, -r 为 true, 则按文件名升序排序
		listInfos.SortByFileNameAsc()
	}

	// 如果启用了长格式输出, 则调用 listCmdLong 函数
	if *listCmdLongFormat {
		if err := listCmdLong(cl, listInfos); err != nil {
			return err
		}
		return nil
	}

	// -q 双引号输出模式
	if *listCmdQuoteNames {
		if err := listCmdQuoteNamesMode(cl, listInfos); err != nil {
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

// 双引号输出模式
func listCmdQuoteNamesMode(cl *colorlib.ColorLib, lfs globals.ListInfos) error {
	// 每行输出3个文件
	for i := 0; i < len(lfs); i += 3 {
		// 检查是否超出了文件列表的范围
		if i >= len(lfs) {
			break
		}

		// 输出文件名
		fmt.Printf("\"%s\"", lfs[i].Name)

		// 检查是否还有下一个文件
		if i+1 < len(lfs) {
			fmt.Printf(" \"%s\"", lfs[i+1].Name)
		}

		if i+2 < len(lfs) {
			fmt.Printf(" \"%s\"", lfs[i+2].Name)
		}

		fmt.Println() // 换行
	}

	return nil
}

// checkListCmdArgs 检查list命令的参数是否合法
func checkListCmdArgs(listPath string) error {
	// 检查路径是否存在
	if _, err := os.Stat(listPath); err != nil {
		// 权限错误
		if os.IsPermission(err) {
			return fmt.Errorf("权限不足: 路径 %s", listPath)
		}

		// 路径不存在
		if os.IsNotExist(err) {
			return fmt.Errorf("路径不存在: %s", listPath)
		}

		// 其他错误
		return fmt.Errorf("检查路径时发生了错误: %v", err)
	}

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
	columns := width / (maxWidth + 1) // 每列增加2个字符的间距
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
			paddedFilename := text.Pad(lfs[j].Name, maxWidth+1, ' ')

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
			infoName := getColorString(info, info.Name, cl)

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
			fmt.Printf("%s%s  %-4s %4s  %8s  %-12s  %-10s\n",
				info.EntryType,
				infoPerm,
				info.Owner,
				info.Group,
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

	// 读取目录下的文件
	files, readDirErr := os.ReadDir(path)
	if readDirErr != nil {
		// 判断是否是权限错误
		if errors.Is(readDirErr, os.ErrPermission) {
			return nil, fmt.Errorf("读取目录 %s 时发生了权限错误: %v", path, readDirErr)
		}

		// 判断是否是不存在的目录错误
		if errors.Is(readDirErr, os.ErrNotExist) {
			return nil, fmt.Errorf("目录 %s 不存在", path)
		}

		// 判断是否是其他类型的错误
		if !errors.Is(readDirErr, os.ErrNotExist) && !errors.Is(readDirErr, os.ErrPermission) {
			return nil, fmt.Errorf("读取目录 %s 时发生了未知错误: %v", path, readDirErr)
		}

		// 如果不是已知的错误类型, 则返回原始错误
		return nil, fmt.Errorf("读取目录 %s 时发生了错误: %v", path, readDirErr)
	}

	// 遍历目录下的每个文件和目录
	for _, file := range files {
		// 跳过隐藏文件
		if !*listCmdAll {
			if isHidden(file.Name()) {
				continue
			}
		}

		// 获取文件的绝对路径
		absPath, absErr := filepath.Abs(file.Name())
		if absErr != nil {
			return nil, fmt.Errorf("获取文件 %s 的绝对路径时发生了错误: %v", file.Name(), absErr)
		}

		// 从绝对路径中提取文件名
		baseName := filepath.Base(absPath)

		// 获取文件的详细信息，如大小、修改时间等
		fileInfo, statErr := file.Info()
		if statErr != nil {
			return nil, fmt.Errorf("获取文件 %s 的信息时发生了错误: %v", file.Name(), statErr)
		}

		// 根据文件的模式判断条目类型
		entryType := getEntryType(fileInfo)

		// 获取文件的大小
		fileSize := fileInfo.Size()

		// 获取文件的最后修改时间
		fileModTime := fileInfo.ModTime()

		// 获取文件的权限信息
		filePerm := fileInfo.Mode().Perm()
		// 将权限信息转换为字符串形式
		filePermStr := filePerm.String()

		// 获取文件所属的用户和组
		u, g := getFileOwner(absPath)

		// 检查文件名是否包含.点, 如果包含尝试获取文件扩展名
		var fileExt string
		if strings.Contains(baseName, ".") {
			// 获取文件扩展名
			fileExt = filepath.Ext(baseName)
		}

		// 构建一个 globals.ListInfo 结构体，存储文件的详细信息
		info := globals.ListInfo{
			EntryType: entryType,   // 条目类型
			Name:      baseName,    // 文件名
			Size:      fileSize,    // 文件大小
			ModTime:   fileModTime, // 文件修改时间
			Perm:      filePermStr, // 文件权限
			Owner:     u,           // 文件所属用户
			Group:     g,           // 文件所属组
			FileExt:   fileExt,     // 文件扩展名
		}

		// 将构建好的文件信息添加到切片中
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
// 参数：
//
//	uid - 用户ID
//	gid - 组ID
//
// 返回：
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
