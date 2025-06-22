package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/globals"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// prepareFileNames 准备文件名列表，处理引号和颜色
// 参数:
//
//	infos - 文件信息列表
//	quoteNames - 是否对文件名添加引号
//	useColor - 是否启用颜色输出
//	cl - 颜色库实例
//
// 返回:
//
//	格式化后的文件名列表
func prepareFileNames(infos globals.ListInfos, quoteNames bool, useColor bool, cl *colorlib.ColorLib) []string {
	fileNames := make([]string, len(infos))
	for i, info := range infos {
		_, file := filepath.Split(info.Name)
		var paddedFilename string
		if quoteNames {
			paddedFilename = fmt.Sprintf("%q", file)
		} else {
			paddedFilename = file
		}

		// 检查是否启用颜色输出
		if useColor {
			paddedFilename = getColorString(info, paddedFilename, cl)
		}

		fileNames[i] = paddedFilename
	}
	return fileNames
}

// renderDefaultTable 创建并渲染默认格式的多列表格
// 参数:
//
//	fileNames - 文件名列表
//	terminalWidth - 终端宽度
//	tableStyle - 表格样式名称
//
// 返回:
//
//	错误信息，如果有
func renderDefaultTable(fileNames []string, terminalWidth int, tableStyle string) error {
	// 动态计算每行可以容纳的列数
	maxWidth := text.LongestLineLen(strings.Join(fileNames, "\n"))

	// 每列增加2个字符的间距
	columns := terminalWidth / (maxWidth + 2)
	if columns == 0 {
		columns = 1 // 至少显示一列
	}

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 构建多列输出
	for i := 0; i < len(fileNames); i += columns {
		// 计算当前行的结束位置
		end := i + columns

		// 确保结束位置不超过文件名列表长度
		if end > len(fileNames) {
			end = len(fileNames)
		}

		// 创建一个临时表格行
		row := make([]interface{}, end-i)
		for j := i; j < end; j++ {
			row[j-i] = fileNames[j]
		}

		// 添加行到表格
		t.AppendRow(row)
	}

	// 设置表格样式
	if tableStyle != "" {
		// 根据表格样式名称设置样式
		if style, ok := globals.TableStyleMap[tableStyle]; ok {
			t.SetStyle(style)
		}
	}

	// 输出表格
	t.Render()

	return nil
}

// renderFileTable 创建并渲染文件信息表格
// 参数:
//
//	cl - 颜色库实例
//	infos - 要显示的文件信息列表
//
// 返回:
//
//	错误信息，如果有
func renderFileTable(cl *colorlib.ColorLib, infos globals.ListInfos) error {
	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	if listCmdTableStyle.Get() != "none" {
		if listCmdShowUserGroup.Get() {
			t.AppendHeader(table.Row{"Type", "Perm", "Owner", "Group", "Size", "Unit", "ModTime", "Name"})
		} else {
			t.AppendHeader(table.Row{"Type", "Perm", "Size", "Unit", "ModTime", "Name"})
		}
	}

	// 添加文件信息行
	for _, info := range infos {
		// 获取文件名（仅文件名，不含路径）
		_, fileName := filepath.Split(info.Name)

		// 获取文件权限字符串
		infoPerm := FormatPermissionString(cl, info)

		// 类型
		infoType := getColorString(info, info.EntryType, cl)

		// 文件大小 和 单位
		infoSize, infoSizeUnit := humanSize(info.Size)

		// 渲染文件大小
		cl.NoBold.Store(true)
		infoSize = cl.Syellow(infoSize)

		// 渲染文件大小单位
		infoSizeUnit = cl.Syellow(infoSizeUnit)

		// 修改时间 禁用加粗
		infoModTime := cl.Sblue(info.ModTime.Format("2006-01-02 15:04:05"))
		cl.NoBold.Store(false)

		// 文件名
		var infoName string
		// 检查是否启用引号
		if listCmdQuoteNames.Get() {
			// 检查是否为软链接
			if info.EntryType == globals.SymlinkType {
				infoName = getColorString(info, fmt.Sprintf("%q -> %q", fileName, info.LinkTargetPath), cl)
			} else {
				infoName = getColorString(info, fmt.Sprintf("%q", fileName), cl)
			}
		} else {
			// 检查是否为软链接
			if info.EntryType == globals.SymlinkType {
				infoName = getColorString(info, fmt.Sprintf("%s -> %s", fileName, info.LinkTargetPath), cl)
			} else {
				infoName = getColorString(info, fileName, cl)
			}
		}

		// 添加行到表格
		if listCmdShowUserGroup.Get() {
			t.AppendRow(table.Row{infoType, infoPerm, info.Owner, info.Group, infoSize, infoSizeUnit, infoModTime, infoName})
		} else {
			t.AppendRow(table.Row{infoType, infoPerm, infoSize, infoSizeUnit, infoModTime, infoName})
		}
	}

	// 设置列的对齐方式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Size", Align: text.AlignRight},     // 文件大小 - 右对齐
		{Name: "Type", Align: text.AlignCenter},    // 文件类型 - 居中对齐
		{Name: "Owner", Align: text.AlignCenter},   // 所有者 - 居中对齐
		{Name: "Group", Align: text.AlignCenter},   // 组 - 居中对齐
		{Name: "Perm", Align: text.AlignLeft},      // 权限 - 左对齐
		{Name: "ModTime", Align: text.AlignCenter}, // 修改时间 - 居中对齐
		{Name: "Unit", Align: text.AlignCenter},    // 单位 - 居中对齐
		{Name: "Name", Align: text.AlignLeft},      // 文件名 - 左对齐
	})

	// 设置表格样式
	if listCmdTableStyle.Get() != "" {
		// 根据-ts的值设置表格样式
		if style, ok := globals.TableStyleMap[listCmdTableStyle.Get()]; ok {
			t.SetStyle(style)
		}
	}

	// 输出表格
	t.Render()

	return nil
}
