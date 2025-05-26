package cmd

import (
	"fmt"
	"os"
	"time"

	"gitee.com/MM-Q/colorlib"
)

// getLast8Chars 函数用于获取输入字符串的最后8个字符
func getLast8Chars(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return s
	}
	return s[len(s)-8:]
}

// printPathColor 根据路径类型以不同颜色输出路径字符串
// 参数:
//
//	path: 要检查的路径(用于获取文件类型信息)
//	s: 要输出的字符串内容
//	cl: colorlib.ColorLib实例，用于彩色输出
//
// 返回值:
//
//	error: 如果获取路径信息失败则返回错误，否则返回nil
func printPathColor(path string, s string, cl *colorlib.ColorLib) error {
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
	// 其他类型文件 - 使用红色输出
	default:
		cl.Red(s)
	}

	return nil
}

// writeFileHeader 写入文件头信息
func writeFileHeader(file *os.File, hashType string, timestampFormat string) error {
	// 获取当前时间
	now := time.Now()

	// 构造文件头内容
	header := fmt.Sprintf("#%s#%s\n", hashType, now.Format(timestampFormat))

	// 写入文件头
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("写入文件头失败: %v", err)
	}
	return nil
}
