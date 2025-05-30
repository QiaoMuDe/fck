package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
)

var ColorMap = map[string]map[string]bool{
	"lred": {
		".log":     true,
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
	},
	"lgreen": {
		".exe":   true,
		".dll":   true,
		".so":    true,
		".run":   true,
		".dylib": true,
		".apk":   true,
		".ipa":   true,
	},
	"lyellow": {
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
	},
	"lblue": {
		".go":   true,
		".py":   true,
		".js":   true,
		".ts":   true,
		".html": true,
		".css":  true,
		".java": true,
		".c":    true,
		".cpp":  true,
		".h":    true,
		".hpp":  true,
		".sh":   true,
		".bash": true,
		".zsh":  true,
		".bat":  true,
		".ps1":  true,
		".rb":   true,
		".rs":   true,
		".php":  true,
	},
	"lpurple": {
		".md5":    true,
		".sha1":   true,
		".sha256": true,
		".sha512": true,
		".hash":   true,
		".mod":    true,
		".sum":    true,
		".pem":    true,
		".key":    true,
		".crt":    true,
		".cer":    true,
	},
	"lcyan": {
		".zip":    true,
		".tar":    true,
		".gz":     true,
		".rar":    true,
		".7z":     true,
		".bz2":    true,
		".xz":     true,
		".jar":    true,
		".war":    true,
		".tgz":    true,
		".tar.gz": true,
		".deb":    true,
		".rpm":    true,
	},
	"lwhite": {
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
			case "lblue":
				cl.Lblue(fs)
			case "lyellow":
				cl.Lyellow(fs)
			case "lgreen":
				cl.Lgreen(fs)
			case "lred":
				cl.Lred(fs)
			case "lpurple":
				cl.Lpurple(fs)
			case "lcyan":
				cl.Lcyan(fs)
			case "lwhite":
				cl.Lwhite(fs)
			}
			return
		}
	}

	// 如果没有匹配的颜色，使用灰色输出
	cl.Gray(fs)
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
	// 目录 - 使用蓝色输出
	case mode.IsDir():
		cl.Blue(s)
	// 普通文件 - 使用绿色输出
	case mode.IsRegular():
		cl.Green(s)
	// 符号链接 - 使用黄色输出
	case mode&os.ModeSymlink != 0:
		cl.Yellow(s)
	// 其他类型文件 - 使用灰色输出
	default:
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
		cl.Blue(path) // 目录 - 使用蓝色输出
	case mode&os.ModeSymlink != 0:
		cl.Yellow(path) // 符号链接 - 使用黄色输出
	case mode&os.ModeDevice != 0:
		cl.Red(path) // 设备文件 - 使用红色输出
	case mode&os.ModeNamedPipe != 0:
		cl.Red(path) // 命名管道 - 使用红色输出
	case mode&os.ModeSocket != 0:
		cl.Red(path) // 套接字文件 - 使用红色输出
	case mode&os.ModeType == 0 && mode&0111 != 0:
		cl.Green(path) // 可执行文件 - 使用绿色输出
	case mode.IsRegular():
		printColoredFile(path, cl) // 普通文件 - 用于根据文件名的后缀输出
	default:
		cl.Gray(path) // 其他类型文件 - 使用灰色输出
	}

	return nil
}
