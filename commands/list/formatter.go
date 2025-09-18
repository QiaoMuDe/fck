// Package list 实现了文件列表的格式化输出功能。
// 该文件提供了表格格式和网格格式两种显示方式，支持颜色输出、权限显示、文件大小格式化等功能。
package list

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

	// 使用预先计算的分组标识，避免重复判断
	if opts.ShouldGroup {
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
		var groupKey string

		if listCmdRecursion.Get() {
			// 递归模式：使用相对路径的目录部分
			dir, _ := filepath.Split(file.Name)
			if dir == "" {
				dir = "."
			}
			groupKey = dir

		} else {
			// 非递归模式：检查是否为通配符展开的情况
			if strings.ContainsAny(file.OriginalPath, "*?[]") {
				// 通配符展开的情况：根据文件类型分组
				if file.EntryType == types.DirType {
					// 目录：使用目录路径作为组键，显示目录内容
					groupKey = file.Path
				} else {
					// 文件：放在当前目录组
					groupKey = "."
				}

			} else {
				// 非通配符：使用用户指定的原始路径
				groupKey = file.OriginalPath
				if groupKey == "" {
					groupKey = "."
				}
			}
		}

		dirFiles[groupKey] = append(dirFiles[groupKey], file)
	}

	// 如果只有一个分组且不是递归模式，直接显示内容
	if len(dirFiles) == 1 && !listCmdRecursion.Get() {
		for _, fileList := range dirFiles {
			if opts.LongFormat {
				return f.renderTable(fileList, opts)
			}
			return f.renderGrid(fileList, opts)
		}
	}

	// 多分组显示
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

// renderGrid 渲染网格格式 (默认格式)
//
// 参数:
//   - files: 文件列表
//   - opts: 格式选项
//
// 返回:
//   - error: 错误
func (f *FileFormatter) renderGrid(files FileInfoList, opts FormatOptions) error {
	if len(files) == 0 {
		return nil
	}

	// 获取终端宽度
	width := f.getSafeTerminalWidth()

	// 准备文件名列表
	fileNames := f.prepareFileNames(files, opts)

	// 计算最优列数
	columns := f.calculateOptimalColumns(fileNames, width)

	// 创建表格
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// 构建多列输出
	for i := 0; i < len(fileNames); i += columns {
		end := i + columns
		if end > len(fileNames) {
			end = len(fileNames)
		}

		// 添加到行
		row := make([]any, end-i)
		for j := i; j < end; j++ {
			row[j-i] = fileNames[j]
		}

		t.AppendRow(row)
	}

	// 设置表格样式
	if opts.TableStyle != "" {
		if style, ok := types.TableStyleMap[opts.TableStyle]; ok {
			t.SetStyle(style)
		}
	}

	t.Render()
	return nil
}

// getSafeTerminalWidth 安全获取终端宽度
//
// 返回:
//   - int: 终端宽度
func (f *FileFormatter) getSafeTerminalWidth() int {
	defaultWidth := 80 // 默认宽度
	minWidth := 40     // 最小宽度
	maxWidth := 1200   // 最大宽度

	// 检查环境变量
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width >= minWidth && width <= maxWidth {
			return width
		}
	}

	// 检查是否为终端
	fd := os.Stdout.Fd()
	if fd > 1024 || !term.IsTerminal(int(fd)) {
		return defaultWidth
	}

	// 安全的类型转换和获取尺寸
	if fd <= uintptr(^uint(0)>>1) { // 确保不会溢出
		if width, _, err := term.GetSize(int(fd)); err == nil {
			if width >= minWidth && width <= maxWidth {
				return width
			}
		}
	}

	return defaultWidth
}

// renderTable 渲染表格格式 (长格式)
//
// 参数:
//   - files: 文件列表
//   - opts: 格式选项
//
// 返回:
//   - error: 错误
func (f *FileFormatter) renderTable(files FileInfoList, opts FormatOptions) error {
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
	for _, info := range files {
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
			info := FileInfo{
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
			name = GetColorString(info, name, f.colorLib) // 添加颜色
		}

		fileNames[i] = name
	}

	return fileNames
}

// getFileNameStatsOptimized 优化的文件名统计信息获取
// 使用一次遍历 + 采样优化，适用于大量文件的场景
//
// 参数:
//   - fileNames: 文件名列表
//
// 返回:
//   - median: 中位数宽度
//   - max: 最大宽度
func (f *FileFormatter) getFileNameStatsOptimized(fileNames []string) (median, max int) {
	if len(fileNames) == 0 {
		return 20, 20 // 默认宽度
	}

	// 采样优化：如果文件数量很大，使用采样来提高性能
	sampleSize := len(fileNames)
	step := 1

	// 当文件数量超过1000时, 进行采样
	if sampleSize > 1000 {
		sampleSize = 1000 // 最多采样1000个文件
		step = len(fileNames) / sampleSize
		if step == 0 {
			step = 1
		}
	}

	// 一次遍历计算所有统计信息
	lengths := make([]int, 0, sampleSize)
	maxWidth := 0

	for i := 0; i < len(fileNames); i += step {
		// 计算实际显示宽度(去除ANSI颜色代码)
		width := text.RuneWidthWithoutEscSequences(fileNames[i])
		lengths = append(lengths, width)

		// 同时记录最大宽度
		if width > maxWidth {
			maxWidth = width
		}
	}

	// 计算中位数
	sort.Ints(lengths)
	mid := len(lengths) / 2
	var medianWidth int

	if len(lengths)%2 == 0 {
		// 偶数个元素，取中间两个的平均值
		medianWidth = (lengths[mid-1] + lengths[mid]) / 2
	} else {
		// 奇数个元素，取中间值
		medianWidth = lengths[mid]
	}

	return medianWidth, maxWidth
}

// calculateOptimalColumns 计算最优列数
// 基于终端宽度和文件名长度分布动态计算最优列数
//
// 参数:
//   - fileNames: 文件名列表
//   - width: 终端宽度
//
// 返回:
//   - int: 最优列数
func (f *FileFormatter) calculateOptimalColumns(fileNames []string, width int) int {
	if len(fileNames) == 0 {
		return 1
	}

	// 一次性获取中位数和最大宽度，避免重复计算
	medianWidth, maxWidth := f.getFileNameStatsOptimized(fileNames)

	// 考虑表格边框和间距, 每列额外需要3个字符的空间
	columnSpacing := 3
	baseColumns := width / (medianWidth + columnSpacing)

	// 限制列数范围
	if baseColumns < 1 {
		return 1
	}

	// 动态计算最大列数，基于终端宽度和文件名长度分布
	// 1. 基础限制：宽终端允许更多列，窄终端限制更少列
	maxColumns := width / (medianWidth * 2)
	if maxColumns < 3 {
		maxColumns = 3 // 最少允许3列
	} else if maxColumns > 12 {
		maxColumns = 12 // 最多允许12列
	}

	// 2. 根据文件名长度分布调整
	// 如果最大宽度远大于中位数宽度，适当减少最大列数
	if maxWidth > medianWidth*3 {
		maxColumns = int(float64(maxColumns) * 0.8) // 减少20%的列数
	}

	// 3. 应用动态计算的最大列数限制
	if baseColumns > maxColumns {
		baseColumns = maxColumns
	}

	// 验证实际效果：检查是否有文件名会导致行过宽
	if maxWidth > width/baseColumns {
		// 如果最长的文件名在当前列数下会超出，适当减少列数
		safeColumns := width / (maxWidth + columnSpacing)
		if safeColumns > 0 && safeColumns < baseColumns {
			return safeColumns
		}
	}

	return baseColumns
}

// addTableRow 添加表格行
//
// 参数:
//   - t: 表格写入器
//   - info: 文件信息
//   - opts: 格式选项
func (f *FileFormatter) addTableRow(t table.Writer, info FileInfo, opts FormatOptions) {
	// 获取文件名
	_, fileName := filepath.Split(info.Name)

	// 格式化权限
	infoPerm := f.formatPermissionString(info)

	// 文件类型
	infoType := GetColorString(info, info.EntryType, f.colorLib)

	// 文件大小和单位
	infoSize, infoSizeUnit := f.humanSize(info.Size)
	f.colorLib.SetBold(false)
	infoSize = f.colorLib.Syellow(infoSize)
	infoSizeUnit = f.colorLib.Syellow(infoSizeUnit)

	// 修改时间
	const timeFormat = "2006-01-02 15:04:05"
	infoModTime := f.colorLib.Sblue(info.ModTime.Format(timeFormat))
	f.colorLib.SetBold(true)

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
		linkFormat := formatStr + arrowColor + formatStr // 软链接格式化占位符

		// 检查软连接目标是否存在
		if _, err := os.Stat(info.LinkTargetPath); os.IsNotExist(err) {
			linkPath := f.colorLib.Sred(fileName)               // 软链接文件名颜色
			sourcePath := f.colorLib.Sgray(info.LinkTargetPath) // 软链接目标路径颜色
			infoName = fmt.Sprintf(linkFormat, linkPath, sourcePath)
		} else {
			linkPath := f.colorLib.Scyan(fileName)
			sourcePath := common.SprintStringColor(info.LinkTargetPath, info.LinkTargetPath, f.colorLib)
			infoName = fmt.Sprintf(linkFormat, linkPath, sourcePath)
		}
	} else {
		// 普通文件
		infoName = GetColorString(info, fmt.Sprintf(formatStr, fileName), f.colorLib)
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
func (f *FileFormatter) formatPermissionString(info FileInfo) string {
	if len(info.Perm) < 10 {
		return "?-?-?"
	}

	if !listCmdColor.Get() {
		return info.Perm[1:]
	}

	// 格式化权限字符串
	var formattedPerm string
	for i := 1; i < len(info.Perm); i++ {
		colorType := permissionColorMap[i]
		switch colorType {
		case colorTypeGreen: // 绿色
			formattedPerm += f.colorLib.Sgreen(string(info.Perm[i]))

		case colorTypeYellow: // 黄色
			formattedPerm += f.colorLib.Syellow(string(info.Perm[i]))

		case colorTypeRed: // 红色
			formattedPerm += f.colorLib.Sred(string(info.Perm[i]))

		default: // 未知颜色类型，使用默认
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
