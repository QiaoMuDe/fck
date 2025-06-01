package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"gitee.com/MM-Q/colorlib"
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

	// 读取目录下的文件
	files, readDirErr := os.ReadDir(listPath)
	if readDirErr != nil {
		return fmt.Errorf("读取目录 %s 时发生了错误: %v", listPath, readDirErr)
	}

	// 获取终端宽度
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("获取终端宽度时发生了错误: %v", err)
	}

	// 提取文件名并排序
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	sort.Strings(filenames) // 按文件名排序

	// 使用 go-pretty 库进行格式化输出
	// 计算最长文件名的长度
	maxWidth := text.LongestLineLen(strings.Join(filenames, "\n"))

	// 动态计算每行可以容纳的列数
	columns := width / (maxWidth + 2) // 每列增加2个字符的间距
	if columns == 0 {
		columns = 1 // 至少显示一列
	}

	// 构建多列输出
	for i := 0; i < len(filenames); i += columns {
		// 计算当前行要显示的文件数
		end := i + columns
		if end > len(filenames) {
			end = len(filenames)
		}

		// 打印当前行的所有文件
		for j := i; j < end; j++ {
			// 使用 go-pretty 的 Pad 函数对齐文件名
			paddedFilename := text.Pad(filenames[j], maxWidth+1, ' ')
			fmt.Print(getColoredPath(paddedFilename, cl))
		}
		fmt.Println()
	}

	return nil
}
