package cmd

import (
	"fmt"
	"os"

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

// 用于根据路径按照类型匹配颜色然后输出
func printPathColor(path string, cl *colorlib.ColorLib) error {
	// 获取路径信息
	pathInfo, statErr := os.Stat(path)
	if statErr != nil {
		return fmt.Errorf("获取路径信息失败: %v", statErr)
	}

	// 根据路径类型设置颜色
	switch mode := pathInfo.Mode(); {
	// 目录
	case mode.IsDir():
		cl.Blue(path)
	// 文件
	case mode.IsRegular():
		cl.Green(path)
	// 符号链接
	case mode&os.ModeSymlink != 0:
		cl.Yellow(path)
	// 其他
	default:
		cl.Red(path)
	}

	return nil
}
