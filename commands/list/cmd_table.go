package list

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// prepareFileNames 准备文件名列表，处理引号和颜色
//
// 参数:
//   - infos: 文件信息列表
//   - quoteNames: 是否对文件名添加引号
//   - useColor: 是否启用颜色输出
//   - cl: 颜色库实例
//
// 返回:
//   - string[]: 处理后的文件名列表
func prepareFileNames(infos types.ListInfos, quoteNames bool, useColor bool, cl *colorlib.ColorLib) []string {
	// 创建一个文件名列表，长度与文件信息列表相同
	fileNames := make([]string, len(infos))
	// 遍历文件信息列表
	for i, info := range infos {
		// 获取文件名
		_, file := filepath.Split(info.Name)
		// 判断是否使用引号
		var paddedFilename string
		if quoteNames {
			// 使用引号
			paddedFilename = fmt.Sprintf("%q", file)
		} else {
			// 不使用引号
			paddedFilename = file
		}

		// 检查是否启用颜色输出
		if useColor {
			// 获取颜色字符串
			paddedFilename = common.GetColorString(listCmdDevColor.Get(), info, paddedFilename, cl)
		}

		// 将文件名添加到文件名列表中
		fileNames[i] = paddedFilename
	}
	// 返回文件名列表
	return fileNames
}

// renderDefaultTable 创建并渲染默认格式的多列表格
//
// 参数:
//   - fileNames: 文件名列表
//   - terminalWidth: 终端宽度
//   - tableStyle: 表格样式名称
//
// 返回:
//   - error: 如果发生错误，返回错误信息，否则返回 nil
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
		if style, ok := types.TableStyleMap[tableStyle]; ok {
			t.SetStyle(style)
		}
	}

	// 输出表格
	t.Render()

	return nil
}

// renderFileTable 创建并渲染文件信息表格
//
// 参数:
//   - cl: 颜色库实例
//   - infos: 要显示的文件信息列表
//
// 返回:
//   - error: 如果发生错误，返回错误信息，否则返回 nil
func renderFileTable(cl *colorlib.ColorLib, infos types.ListInfos) error {
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
	showUserGroup := listCmdShowUserGroup.Get() // 检查是否显示用户组
	quoteNames := listCmdQuoteNames.Get()       // 检查是否对文件名添加引号
	formatStr := "%s"                           // 默认格式化字符串
	if quoteNames {
		formatStr = "\"%s\"" // 如果需要引号，则使用添加引号的格式(这里不能使用%q会导致样式渲染问题)
	}
	arrow := " -> "                                  // 软链接箭头
	arrowColor := cl.Swhite(arrow)                   // 软链接箭头颜色
	linkFormat := formatStr + arrowColor + formatStr // 软链接格式化字符串
	nameFormat := formatStr                          // 文件名格式化字符串
	const timeFormat = "2006-01-02 15:04:05"         // 时间格式

	// 遍历文件信息列表
	for _, info := range infos {
		// 获取文件名（仅文件名，不含路径）
		_, fileName := filepath.Split(info.Name)

		// 获取文件权限字符串
		infoPerm := FormatPermissionString(cl, info)

		// 类型
		infoType := common.GetColorString(listCmdDevColor.Get(), info, info.EntryType, cl)

		// 文件大小 和 单位
		infoSize, infoSizeUnit := humanSize(info.Size)

		// 渲染文件大小
		cl.NoBold.Store(true)
		infoSize = cl.Syellow(infoSize)

		// 渲染文件大小单位
		infoSizeUnit = cl.Syellow(infoSizeUnit)

		// 修改时间 禁用加粗
		infoModTime := cl.Sblue(info.ModTime.Format(timeFormat))
		cl.NoBold.Store(false)

		// 文件名
		var infoName string

		// 检查是否为软链接
		if info.EntryType == types.SymlinkType {
			// 根据软连接路径是否存在设置颜色
			if _, statErr := os.Stat(info.LinkTargetPath); os.IsNotExist(statErr) {
				// 如果源文件不存在, 则软连接为红色
				linkPath := cl.Sred(fileName)
				// 源文件为灰色
				sourcePath := cl.Sgray(info.LinkTargetPath)
				// 组装字符串
				infoName = fmt.Sprintf(linkFormat, linkPath, sourcePath)
			} else {
				// 如果源文件存在, 则软连接为青色
				linkPath := cl.Scyan(fileName)
				// 源文件为根据类型设置颜色
				sourcePath := common.SprintStringColor(info.LinkTargetPath, info.LinkTargetPath, cl)
				// 组装字符串
				infoName = fmt.Sprintf(linkFormat, linkPath, sourcePath)
			}
		} else {
			infoName = common.GetColorString(listCmdDevColor.Get(), info, fmt.Sprintf(nameFormat, fileName), cl)
		}

		// 添加行到表格
		if showUserGroup {
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
		if style, ok := types.TableStyleMap[listCmdTableStyle.Get()]; ok {
			t.SetStyle(style)
		}
	}

	// 输出表格
	t.Render()

	return nil
}
