package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
)

// ColorMap 定义了不同颜色对应的文件后缀名映射
var ColorMap = map[string]map[string]bool{
	"red": {
		".bak":     true,
		".tmp":     true,
		".swp":     true,
		".swo":     true,
		".old":     true,
		".db":      true,
		".sql":     true,
		".sqlite":  true,
		".db3":     true,
		".sqlite3": true,
		".mdb":     true,
		".pdb":     true,
		".zip":     true,
		".tar":     true,
		".gz":      true,
		".rar":     true,
		".7z":      true,
		".bz2":     true,
		".xz":      true,
		".jar":     true,
		".war":     true,
		".tgz":     true,
		".tar.gz":  true,
		".deb":     true,
		".rpm":     true,
		".md5":     true,
		".sha1":    true,
		".sha256":  true,
		".sha512":  true,
		".hash":    true,
		".mod":     true,
		".sum":     true,
		".pem":     true,
		".key":     true,
		".crt":     true,
		".cer":     true,
	},
	"green": {
		".exe":   true,
		".dll":   true,
		".so":    true,
		".run":   true,
		".dylib": true,
		".apk":   true,
		".ipa":   true,
		".go":    true,
		".py":    true,
		".js":    true,
		".ts":    true,
		".html":  true,
		".css":   true,
		".java":  true,
		".c":     true,
		".cpp":   true,
		".h":     true,
		".hpp":   true,
		".sh":    true,
		".bash":  true,
		".zsh":   true,
		".bat":   true,
		".ps1":   true,
		".rb":    true,
		".rs":    true,
		".php":   true,
	},
	"yellow": {
		".txt":  true,
		".md":   true,
		".json": true,
		".xml":  true,
		".yaml": true,
		".yml":  true,
		".toml": true,
		".ini":  true,
		".conf": true,
		".cfg":  true,
		".doc":  true,
		".docx": true,
		".pdf":  true,
		".csv":  true,
		".xls":  true,
		".xlsx": true,
		".ppt":  true,
		".pptx": true,
		".rtf":  true,
		".log":  true,
		".lnk":  true,
	},
	"white": {
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".webp": true,
		".svg":  true,
		".ico":  true,
		".avif": true,
		".mp4":  true,
		".mkv":  true,
		".avi":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".mp3":  true,
		".wav":  true,
		".ogg":  true,
		".flac": true,
		".aac":  true,
		".m4a":  true,
	},
}

// printColoredFile 根据文件后缀名以不同颜色输出文件路径
func printColoredFile(fs string, cl *colorlib.ColorLib) {
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
			case "white":
				// 检查是否包含目录分割符
				if strings.Contains(fs, string(os.PathSeparator)) {
					// 把路径分割成目录和文件名
					dir, file := filepath.Split(fs)

					fmt.Println(cl.Scyan(dir) + cl.Swhite(file))
				} else {
					cl.White(fs) // 如果没有目录分割符，则直接输出文件名
				}
			}
			return
		}
	}

	// 检查是否包含目录分割符
	if strings.Contains(fs, string(os.PathSeparator)) {
		// 把路径分割成目录和文件名
		dir, file := filepath.Split(fs)

		fmt.Println(cl.Scyan(dir) + cl.Sgray(file))
	} else {
		cl.Gray(fs) // 如果没有目录分割符，则直接输出文件名
	}
}

// printStringColor 根据路径类型以不同颜色输出字符串
// 参数:
//
//	path: 要检查的路径(用于获取文件类型信息)
//	s: 要输出的字符串内容
//	cl: colorlib.ColorLib实例，用于彩色输出
//
// 返回值:
//
//	error: 如果获取路径信息失败则返回错误，否则返回nil
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
		// 符号链接 - 使用黄色输出
		cl.Yellow(s)

	default:
		// 其他类型文件 - 使用灰色输出
		cl.Gray(s)
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
		// 符号链接 - 使用黄色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Syellow(file))
		} else {
			cl.Yellow(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeDevice != 0:
		// 设备文件 - 使用红色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Sred(file))
		} else {
			cl.Red(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeNamedPipe != 0:
		// 命名管道 - 使用红色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Sred(file))
		} else {
			cl.Red(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeSocket != 0:
		// 套接字文件 - 使用红色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Sred(file))
		} else {
			cl.Red(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode&os.ModeType == 0 && mode&0111 != 0:
		// 可执行文件 - 使用绿色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)
			fmt.Println(cl.Scyan(dir) + cl.Sgreen(file))
		} else {
			cl.Green(path) // 如果没有目录分割符，则直接输出文件名
		}
	case mode.IsRegular():
		printColoredFile(path, cl) // 普通文件 - 用于根据文件名的后缀输出
	default:
		// 其他类型文件 - 使用灰色输出
		// 检查是否包含目录分割符
		if strings.Contains(path, string(os.PathSeparator)) {
			// 把路径分割成目录和文件名
			dir, file := filepath.Split(path)

			fmt.Println(cl.Scyan(dir) + cl.Sgray(file))
		} else {
			cl.Gray(path) // 如果没有目录分割符，则直接输出文件名
		}
	}

	return nil
}
