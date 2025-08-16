package list

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fck/commands/internal/common"
	"gitee.com/MM-Q/fck/commands/internal/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/term"
)

// FileFormatter 文件格式化器
type FileFormatter struct {
	colorLib *colorlib.ColorLib
}

// NewFileFormatter 创建新的文件格式化器
func NewFileFormatter(cl *colorlib.ColorLib) *FileFormatter {
	return &FileFormatter{
		colorLib: cl,
	}
}

// Render 渲染文件列表
//
// 参数:
//   - files: 文件列表
//   - opts: 格式选项
//
// 返回:
//   - error: 错误
func (f *FileFormatter) Render(files FileInfoList, opts FormatOptions) error {
	if len(files) == 0 {
		return nil
	}

	// 根据是否递归分组显示
	if listCmdRecursion.Get() {
		return f.renderGrouped(files, opts)
	}

	// 直接渲染
	if opts.LongFormat {
		return f.renderTable(files, opts)
	}
	return f.renderGrid(files, opts)
}

// renderGrouped 按目录分组渲染
//
// 参数:
//   - files: 文件列表
//   - opts: 格式选项
//
// 返回:
//   - error: 错误
func (f *FileFormatter) renderGrouped(files FileInfoList, opts FormatOptions) error {
	// 按目录分组
	dirFiles := make(map[string]FileInfoList)
	for _, file := range files {
		dir, _ := filepath.Split(file.Name)
		if dir == "" {
			dir = "."
		}
		dirFiles[dir] = append(dirFiles[dir], file)
	}

	// 排序目录
	dirs := make([]string, 0, len(dirFiles))
	for dir := range dirFiles {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	// 渲染每个目录
	for i, dir := range dirs {
		if i > 0 {
			fmt.Println() // 目录间空行
		}

		// 打印目录标题
		f.colorLib.Bluef("%s:\n", dir)

		// 渲染目录内容
		dirFileList := dirFiles[dir]
		if opts.LongFormat {
			if err := f.renderTable(dirFileList, opts); err != nil {
				return fmt.Errorf("渲染目录 %s 表格失败: %v", dir, err)
			}
		} else {
			if err := f.renderGrid(dirFileList, opts); err != nil {
				return fmt.Errorf("渲染目录 %s 网格失败: %v", dir, err)
			}
		}
	}

	return nil
}

// renderGrid 渲染网格格式（默认格式）
//
// 参数:
//   - files: 文件列表
//   - opts: 格式选项
//
// 返回:
//   - error: 错误
func (f *FileFormatter) renderGrid(files FileInfoList, opts FormatOptions) error {
	// 获取终端宽度
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80 // 默认宽度
	}

	// 准备文件名列表
	fileNames := f.prepareFileNames(files, opts)

	// 计算列数
	maxWidth := f.getMaxWidth(fileNames)
	columns := width / (maxWidth + 2)
	if columns == 0 {
		columns = 1
	}

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 构建多列输出
	for i := 0; i < len(fileNames); i += columns {
		end := i + columns
		if end > len(fileNames) {
			end = len(fileNames)
		}

		row := make([]interface{}, end-i)
		for j := i; j < end; j++ {
			row[j-i] = fileNames[j]
		}

		t.AppendRow(row)
	}

	// 设置表格样式
	if opts.TableStyle != "" && opts.TableStyle != "none" {
		if style, ok := types.TableStyleMap[opts.TableStyle]; ok {
			t.SetStyle(style)
		}
	}

	t.Render()
	return nil
}

// renderTable 渲染表格格式（长格式）
//
// 参数:
//   - files: 文件列表
//   - opts: 格式选项
//
// 返回:
//   - error: 错误
func (f *FileFormatter) renderTable(files FileInfoList, opts FormatOptions) error {
	// 转换为原有格式以复用现有逻辑
	var listInfos types.ListInfos
	for _, file := range files {
		listInfo := types.ListInfo{
			EntryType:      file.EntryType,
			Name:           file.Name,
			Size:           file.Size,
			ModTime:        file.ModTime,
			Perm:           file.Perm,
			Owner:          file.Owner,
			Group:          file.Group,
			FileExt:        file.FileExt,
			LinkTargetPath: file.LinkTargetPath,
		}
		listInfos = append(listInfos, listInfo)
	}

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	if opts.TableStyle != "none" {
		if opts.ShowUserGroup {
			t.AppendHeader(table.Row{"Type", "Perm", "Owner", "Group", "Size", "Unit", "ModTime", "Name"})
		} else {
			t.AppendHeader(table.Row{"Type", "Perm", "Size", "Unit", "ModTime", "Name"})
		}
	}

	// 添加数据行
	for _, info := range listInfos {
		f.addTableRow(t, info, opts)
	}

	// 设置列对齐
	f.configureColumns(t)

	// 设置表格样式
	if opts.TableStyle != "" {
		if style, ok := types.TableStyleMap[opts.TableStyle]; ok {
			t.SetStyle(style)
		}
	}

	t.Render()
	return nil
}

// prepareFileNames 准备文件名列表
//
// 参数:
//   - files: 文件列表
//   - opts: 格式选项
//
// 返回:
//   - []string: 文件名列表
func (f *FileFormatter) prepareFileNames(files FileInfoList, opts FormatOptions) []string {
	fileNames := make([]string, len(files))

	// 遍历所有文件
	for i, file := range files {
		_, fileName := filepath.Split(file.Name)

		// 处理引号
		var name string
		if opts.QuoteNames {
			name = fmt.Sprintf("%q", fileName)
		} else {
			name = fileName
		}

		// 处理颜色
		if opts.UseColor {
			info := types.ListInfo{
				EntryType:      file.EntryType,      // 文件类型
				Name:           file.Name,           // 文件名
				Size:           file.Size,           // 文件大小
				ModTime:        file.ModTime,        // 修改时间
				Perm:           file.Perm,           // 权限
				Owner:          file.Owner,          // 所有者
				Group:          file.Group,          // 所属组
				FileExt:        file.FileExt,        // 文件扩展名
				LinkTargetPath: file.LinkTargetPath, // 符号链接目标路径
			}
			name = common.GetColorString(opts.DevColor, info, name, f.colorLib) // 添加颜色
		}

		fileNames[i] = name
	}

	return fileNames
}

// getMaxWidth 获取最大宽度
//
// 参数:
//   - fileNames: 文件名列表
//
// 返回:
//   - int: 最大宽度
func (f *FileFormatter) getMaxWidth(fileNames []string) int {
	return text.LongestLineLen(strings.Join(fileNames, "\n"))
}

// addTableRow 添加表格行
//
// 参数:
//   - t: 表格写入器
//   - info: 文件信息
//   - opts: 格式选项
func (f *FileFormatter) addTableRow(t table.Writer, info types.ListInfo, opts FormatOptions) {
	// 获取文件名
	_, fileName := filepath.Split(info.Name)

	// 格式化权限
	infoPerm := f.formatPermissionString(info)

	// 文件类型
	infoType := common.GetColorString(opts.DevColor, info, info.EntryType, f.colorLib)

	// 文件大小和单位
	infoSize, infoSizeUnit := f.humanSize(info.Size)
	f.colorLib.NoBold.Store(true)
	infoSize = f.colorLib.Syellow(infoSize)
	infoSizeUnit = f.colorLib.Syellow(infoSizeUnit)

	// 修改时间
	const timeFormat = "2006-01-02 15:04:05"
	infoModTime := f.colorLib.Sblue(info.ModTime.Format(timeFormat))
	f.colorLib.NoBold.Store(false)

	// 文件名处理
	var infoName string
	formatStr := "%s"
	if opts.QuoteNames {
		formatStr = "\"%s\""
	}

	// 符号链接特殊处理
	if info.EntryType == types.SymlinkType {
		arrow := " -> "                                  // 软链接箭头
		arrowColor := f.colorLib.Swhite(arrow)           // 软链接箭头颜色
		linkFormat := formatStr + arrowColor + formatStr // 软链接格式化字符串

		// 检查软连接目标是否存在
		if _, err := os.Stat(info.LinkTargetPath); os.IsNotExist(err) {
			linkPath := f.colorLib.Sred(fileName)                    // 软链接文件名颜色
			sourcePath := f.colorLib.Sgray(info.LinkTargetPath)      // 软链接目标路径颜色
			infoName = fmt.Sprintf(linkFormat, linkPath, sourcePath) // 软链接格式化字符串
		} else {
			linkPath := f.colorLib.Scyan(fileName)
			sourcePath := common.SprintStringColor(info.LinkTargetPath, info.LinkTargetPath, f.colorLib)
			infoName = fmt.Sprintf(linkFormat, linkPath, sourcePath)
		}
	} else {
		// 普通文件
		infoName = common.GetColorString(opts.DevColor, info, fmt.Sprintf(formatStr, fileName), f.colorLib)
	}

	// 添加行
	if opts.ShowUserGroup {
		t.AppendRow(table.Row{infoType, infoPerm, info.Owner, info.Group, infoSize, infoSizeUnit, infoModTime, infoName})
	} else {
		t.AppendRow(table.Row{infoType, infoPerm, infoSize, infoSizeUnit, infoModTime, infoName})
	}
}

// formatPermissionString 格式化权限字符串
//
// 参数:
//   - info: 文件信息
//
// 返回:
//   - string: 格式化后的权限字符串
func (f *FileFormatter) formatPermissionString(info types.ListInfo) string {
	if len(info.Perm) < 10 {
		return "?-?-?"
	}

	if !listCmdColor.Get() {
		return info.Perm[1:]
	}

	// 格式化权限字符串
	var formattedPerm string
	for i := 1; i < len(info.Perm); i++ {
		colorName := common.PermissionColorMap[i]
		switch colorName {
		case "green":
			formattedPerm += f.colorLib.Sgreen(string(info.Perm[i]))
		case "yellow":
			formattedPerm += f.colorLib.Syellow(string(info.Perm[i]))
		case "red":
			formattedPerm += f.colorLib.Sred(string(info.Perm[i]))
		default:
			formattedPerm += string(info.Perm[i])
		}
	}

	return formattedPerm
}

// configureColumns 配置列对齐
//
// 参数:
//   - t: 表格写入器
func (f *FileFormatter) configureColumns(t table.Writer) {
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Size", Align: text.AlignRight},
		{Name: "Type", Align: text.AlignCenter},
		{Name: "Owner", Align: text.AlignCenter},
		{Name: "Group", Align: text.AlignCenter},
		{Name: "Perm", Align: text.AlignLeft},
		{Name: "ModTime", Align: text.AlignCenter},
		{Name: "Unit", Align: text.AlignCenter},
		{Name: "Name", Align: text.AlignLeft},
	})
}

// humanSize 转换文件大小为可读格式
//
// 参数:
//   - size: 文件大小
//
// 返回:
//   - string: 可读的文件大小
//   - string: 单位
func (f *FileFormatter) humanSize(size int64) (string, string) {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	base := float64(1024)

	var unit string
	sizeFloat := float64(size)

	// 判断处理单位
	if sizeFloat < base {
		unit = units[0]
	} else if sizeFloat < base*base {
		unit = units[1]
		sizeFloat /= base
	} else if sizeFloat < base*base*base {
		unit = units[2]
		sizeFloat /= base * base
	} else if sizeFloat < base*base*base*base {
		unit = units[3]
		sizeFloat /= base * base * base
	} else if sizeFloat < base*base*base*base*base {
		unit = units[4]
		sizeFloat /= base * base * base * base
	} else {
		unit = units[5]
		sizeFloat /= base * base * base * base * base
	}

	var sizeStr string
	if sizeFloat < 10 {
		sizeStr = fmt.Sprintf("%.1f", sizeFloat) // 保留一位小数
	} else {
		sizeStr = fmt.Sprintf("%.0f", sizeFloat) // 不保留小数
	}

	// 处理0值情况
	sizeStr = strings.TrimSuffix(sizeStr, ".0")

	// 处理0.0值情况
	if sizeStr == "0" || sizeStr == "0.0" {
		return "0", "B"
	}

	return sizeStr, unit
}
