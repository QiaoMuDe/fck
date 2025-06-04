package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
)

// ColorMap 定义了不同颜色对应的文件后缀名映射
var ColorMap = map[string]map[string]bool{
	"red": {
		".exe":   true, // Windows可执行文件
		".bin":   true, // 二进制可执行文件
		".app":   true, // macOS应用程序
		".msi":   true, // Windows安装程序
		".dll":   true, // Windows动态链接库
		".so":    true, // Linux共享库
		".run":   true, // Linux可执行文件
		".dylib": true, // macOS动态库
		".sh":    true, // Shell脚本
		".bash":  true, // Bash脚本
		".zsh":   true, // Zsh脚本
		".bat":   true, // Windows批处理脚本
		".ps1":   true, // PowerShell脚本
		".wasm":  true, // WebAssembly文件
	},
	"yellow": {
		".md":            true, // Markdown文件
		".json":          true, // JSON文件
		".jsonl":         true, // JSON Lines文件
		".xml":           true, // XML文件
		".yaml":          true, // YAML文件
		".yml":           true, // YAML配置文件
		".toml":          true, // TOML配置文件
		".ini":           true, // INI配置文件
		".conf":          true, // 配置文件
		".cfg":           true, // 配置文件
		".properties":    true, // Java属性文件
		".env":           true, // 环境变量文件
		".gitignore":     true, // Git忽略文件
		".gitattributes": true, // Git属性配置文件
		".gitmodules":    true, // Git子模块配置文件
		".gitkeep":       true, // Git保留文件
		".gitconfig":     true, // Git配置文件
		".git":           true, // Git版本控制系统文件
		".svn":           true, // Subversion版本控制系统文件
		".hg":            true, // Mercurial版本控制系统文件
		".bzr":           true, // Bazaar版本控制系统文件
		".lock":          true, // 锁文件
		".lockfile":      true, // 锁文件
		".bak":           true, // 备份文件
		".bakup":         true, // 备份文件
		".tmp":           true, // 临时文件
		".swp":           true, // Vim交换文件
		".swo":           true, // Vim交换文件
	},
	"green": {
		".c":       true, // C源文件
		".cpp":     true, // C++源文件
		".h":       true, // C头文件
		".hpp":     true, // C++头文件
		".go":      true, // Go源文件
		".py":      true, // Python脚本
		".class":   true, // Java字节码文件
		".js":      true, // JavaScript文件
		".ts":      true, // TypeScript文件
		".html":    true, // HTML文件
		".css":     true, // CSS样式表
		".java":    true, // Java源文件
		".pyd":     true, // Python字节码文件
		".pyc":     true, // Python编译文件
		".pyo":     true, // Python编译文件
		".pyw":     true, // Python脚本
		".rb":      true, // Ruby脚本
		".rs":      true, // Rust脚本
		".php":     true, // PHP脚本
		".swift":   true, // Swift脚本
		".kotlin":  true, // Kotlin脚本
		".scala":   true, // Scala脚本
		".elm":     true, // Elm脚本
		".lua":     true, // Lua脚本
		".pl":      true, // Perl脚本
		".perl":    true, // Perl脚本
		".r":       true, // R语言脚本
		".vbs":     true, // VBScript文件
		".psm1":    true, // PowerShell模块文件
		".s":       true, // 汇编语言源文件
		".o":       true, // 编译后的目标文件
		".a":       true, // 静态库文件
		".lib":     true, // 静态库文件
		".sql":     true, // SQL脚本文件
		".db":      true, // 数据库文件
		".db3":     true, // SQLite数据库文件
		".sqlite3": true, // SQLite数据库文件
		".sqlite":  true, // SQLite数据库文件
		".dbf":     true, // DBF文件
		".mdb":     true, // Access数据库文件
		".accdb":   true, // Access数据库文件
		".dbx":     true, // DBX文件
		".pdb":     true, // 调试信息文件
	},
	"purple": {
		".zip":     true, // ZIP压缩文件
		".tar":     true, // TAR压缩文件
		".gz":      true, // GZIP压缩文件
		".bz2":     true, // BZIP2压缩文件
		".xz":      true, // XZ压缩文件
		".7z":      true, // 7Z压缩文件
		".rar":     true, // RAR压缩文件
		".tar.gz":  true, // TAR.GZ压缩文件
		".tar.bz2": true, // TAR.BZ2压缩文件
		".tar.xz":  true, // TAR.XZ压缩文件
		".tgz":     true, // TAR.GZ压缩文件
		".tbz2":    true, // TAR.BZ2压缩文件
		".txz":     true, // TAR.XZ压缩文件
		".jar":     true, // Java Archive文件
		".war":     true, // Web Archive文件
		".ear":     true, // Enterprise Archive文件
		".apk":     true, // Android应用包
		".ipa":     true, // iOS应用包
	},
	"cyan": {
		".lnk": true, // Windows快捷方式文件
	},
}

// 定义全局常量的颜色映射
var PermissionColorMap = map[int]string{
	1: "green",  // 所有者-读-绿色
	2: "yellow", // 所有者-写-黄色
	3: "red",    // 所有者-执行-红色
	4: "green",  // 组-读-绿色
	5: "yellow", // 组-写-黄色
	6: "red",    // 组-执行-红色
	7: "green",  // 其他-读-绿色
	8: "yellow", // 其他-写-黄色
	9: "red",    // 其他-执行-红色
}

// printColoredFile 根据文件后缀名以不同颜色输出文件路径
func printColoredFile(fs string, cl *colorlib.ColorLib) {
	// 获取文件后缀名
	fileExt := strings.ToLower(filepath.Ext(fs))

	for color, extensions := range ColorMap {
		if extensions[fileExt] {
			switch color {
			case "yellow":
				// 检查是否包含目录分割符
				if strings.Contains(fs, string(os.PathSeparator)) {
					// 把路径分割成目录和文件名
					dir, file := filepath.Split(fs)

					fmt.Println(cl.Scyan(dir) + cl.Syellow(file))
				} else {
					cl.Yellow(fs) // 如果没有目录分割符，则直接输出文件名
				}
			case "green":
				// 检查是否包含目录分割符
				if strings.Contains(fs, string(os.PathSeparator)) {
					// 把路径分割成目录和文件名
					dir, file := filepath.Split(fs)

					fmt.Println(cl.Scyan(dir) + cl.Sgreen(file))
				} else {
					cl.Green(fs) // 如果没有目录分割符，则直接输出文件名
				}
			case "red":
				// 检查是否包含目录分割符
				if strings.Contains(fs, string(os.PathSeparator)) {
					// 把路径分割成目录和文件名
					dir, file := filepath.Split(fs)

					fmt.Println(cl.Scyan(dir) + cl.Sred(file))
				} else {
					cl.Red(fs) // 如果没有目录分割符，则直接输出文件名
				}
			case "purple":
				// 检查是否包含目录分割符
				if strings.Contains(fs, string(os.PathSeparator)) {
					// 把路径分割成目录和文件名
					dir, file := filepath.Split(fs)

					fmt.Println(cl.Scyan(dir) + cl.Spurple(file))
				} else {
					cl.Purple(fs) // 如果没有目录分割符，则直接输出文件名
				}
			case "cyan":
				// 检查是否包含目录分割符
				if strings.Contains(fs, string(os.PathSeparator)) {
					// 把路径分割成目录和文件名
					dir, file := filepath.Split(fs)

					fmt.Println(cl.Scyan(dir) + cl.Scyan(file))
				} else {
					cl.Cyan(fs) // 如果没有目录分割符，则直接输出文件名
				}
			}
			return
		}
	}

	// 如果没有匹配的文件后缀名, 则使用白色来渲染字符串
	// 检查是否包含目录分割符
	if strings.Contains(fs, string(os.PathSeparator)) {
		// 把路径分割成目录和文件名
		dir, file := filepath.Split(fs)

		fmt.Println(cl.Scyan(dir) + cl.Swhite(file))
	} else {
		cl.White(fs) // 如果没有目录分割符，则直接输出文件名
	}
}

// printStringColor 根据路径类型以不同颜色输出字符串
// 参数:
//
//	path: 要检查的路径(用于获取文件类型信息)
//	s: 要输出的字符串内容
//	cl: colorlib.ColorLib实例, 用于彩色输出
//
// 返回值:
//
//	error: 如果获取路径信息失败则返回错误, 否则返回nil
func printStringColor(path string, s string, cl *colorlib.ColorLib) error {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(path)
	if statErr != nil {
		return fmt.Errorf("获取路径信息失败: %v", statErr)
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		cl.Blue(s)

	case mode.IsRegular():
		// 普通文件 - 使用绿色输出
		cl.Green(s)

	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		cl.Scyan(s)

	default:
		// 其他类型文件 - 使用白色输出
		cl.White(s)
	}

	return nil
}

// printPathColor 根据路径类型以不同颜色输出路径字符串
// 参数:
//
//	path: 要检查的路径(用于获取文件类型信息)
//	cl: colorlib.ColorLib实例, 用于彩色输出
//
// 返回值:
//
//	error: 如果获取路径信息失败则返回错误, 否则返回nil
func printPathColor(path string, cl *colorlib.ColorLib) error {
	// 获取路径信息
	pathInfo, statErr := os.Lstat(path)
	if statErr != nil {
		return fmt.Errorf("获取路径信息失败: %v", statErr)
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	case mode.IsDir():
		// 目录 - 使用蓝色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Sblue(file))
		} else {
			cl.Blue(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeSymlink != 0:
		// 符号链接 - 使用青色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Scyan(file))
		} else {
			cl.Cyan(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用黄色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Syellow(file))
		} else {
			cl.Yellow(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用黄色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Syellow(file))
		} else {
			cl.Yellow(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用黄色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Syellow(file))
		} else {
			cl.Yellow(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeType == 0 && mode&0111 != 0:
		// 可执行文件 - 使用红色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)
			fmt.Println(cl.Scyan(dir) + cl.Sred(file))
		} else {
			cl.Red(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode.IsRegular():
		printColoredFile(path, cl) // 普通文件 - 用于根据文件名的后缀输出
	default:
		// 其他类型文件 - 使用白色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Swhite(file))
		} else {
			cl.White(path) // 如果没有目录分割符，则直接输出文件名
		}
	}

	return nil
}

// getColorString 函数的作用是根据传入的文件信息、路径字符串以及颜色库实例，返回带有相应颜色的路径字符串。
// 参数:
// info: 包含文件类型和文件后缀名等信息的 globals.ListInfo 结构体实例。
// pF: 要处理的路径字符串。
// cl: 用于彩色输出的 colorlib.ColorLib 实例。
// 返回值:
// colorString: 经过颜色处理后的路径字符串。
func getColorString(info globals.ListInfo, pF string, cl *colorlib.ColorLib) (colorString string) {
	// 依据文件的类型来确定输出的颜色
	switch info.EntryType {
	case globals.DirType:
		// 若文件类型为目录，则使用蓝色来渲染字符串
		colorString = cl.Sblue(pF)
	case globals.SymlinkType:
		// 若文件类型为符号链接，则使用青色来渲染字符串
		colorString = cl.Scyan(pF)
	case globals.ExecutableType:
		// 若文件类型为可执行文件，则使用红色来渲染字符串
		colorString = cl.Sred(pF)
	case globals.SocketType, globals.PipeType, globals.BlockDeviceType, globals.CharDeviceType:
		// 若文件类型为套接字、管道、块设备、字符设备，则使用黄色来渲染字符串
		colorString = cl.Syellow(pF)
	case globals.EmptyType:
		// 若文件类型为空文件, 则使用黑色来渲染字符串
		colorString = cl.Sblack(pF)
	case globals.FileType:
		// 若文件类型为普通文件，则根据文件后缀名来确定颜色
		for color, extensions := range ColorMap {
			if extensions[info.FileExt] {
				switch color {
				case "yellow":
					colorString = cl.Syellow(pF)
				case "green":
					colorString = cl.Sgreen(pF)
				case "red":
					colorString = cl.Sred(pF)
				case "purple":
					colorString = cl.Spurple(pF)
				case "cyan":
					colorString = cl.Scyan(pF)
				}
				return colorString
			}
		}

		// 若没有找到匹配的文件后缀名, 则使用白色来渲染字符串
		colorString = cl.Swhite(pF)
	default:
		// 对于未匹配的类型，使用灰色来渲染字符串
		colorString = cl.Sgray(pF)
	}

	return colorString
}
