package qflag

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// generateHelpInfo 生成命令帮助信息
// cmd: 当前命令
// 返回值: 命令帮助信息
func generateHelpInfo(cmd *Cmd) string {
	// 处理根命令
	if cmd == nil {
		return ""
	}

	// 确保内置标志已初始化
	cmd.initBuiltinFlags()

	// 如果用户指定了自定义帮助信息则优先返回
	if cmd.userInfo.help != "" {
		return cmd.userInfo.help
	}

	// 根据语言选择模板实例
	var tpl HelpTemplate
	if cmd.GetUseChinese() {
		tpl = ChineseTemplate // 中文模板实例
	} else {
		tpl = EnglishTemplate // 英文模板实例
	}

	// 使用buffer提高字符串拼接性能
	var buf bytes.Buffer

	// 写入命令名称和描述
	writeCommandHeader(cmd, tpl, &buf)

	// 写入Logo信息
	writeLogoText(cmd, &buf)

	// 写入用法说明
	writeUsageLine(cmd, tpl, &buf)

	// 写入选项信息
	writeOptions(cmd, tpl, &buf)

	// 写入子命令信息
	writeSubCmds(cmd, tpl, &buf)

	// 写入自定义模块帮助信息
	writeModuleHelps(cmd, &buf)

	// 写入示例信息
	writeExamples(cmd, tpl, &buf)

	// 写入注意事项
	writeNotes(cmd, tpl, &buf)

	return buf.String()
}

// writeModuleHelps 写入自定义模块帮助信息
// cmd: 当前命令
// buf: 输出缓冲区
func writeModuleHelps(cmd *Cmd, buf *bytes.Buffer) {
	// 如果存在自定义模块帮助信息，则写入
	if cmd.GetModuleHelps() != "" {
		buf.WriteString("\n" + cmd.GetModuleHelps() + "\n")
	}
}

// writeLogoText 写入Logo信息
// cmd: 当前命令
// buf: 输出缓冲区
func writeLogoText(cmd *Cmd, buf *bytes.Buffer) {
	// 如果配置了Logo文本, 则写入
	if cmd.userInfo.logoText != "" {
		buf.WriteString(cmd.GetLogoText() + "\n")
	}
}

// writeCommandHeader 写入命令名称和描述
// cmd: 当前命令
// tpl: 模板实例
// buf: 输出缓冲区
func writeCommandHeader(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 如果子命令不为空, 则使用带有短名称的模板
	if cmd.ShortName() != "" {
		fmt.Fprintf(buf, tpl.CmdNameWithShort, cmd.LongName(), cmd.ShortName())
	} else {
		fmt.Fprintf(buf, tpl.CmdName, cmd.LongName())
	}

	// 如果描述不为空, 则写入描述
	if cmd.userInfo.description != "" {
		fmt.Fprintf(buf, tpl.CmdDescription, cmd.userInfo.description)
	}
}

// writeUsageLine 写入用法说明
// cmd: 当前命令
// tpl: 模板实例
// buf: 输出缓冲区
func writeUsageLine(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 使用模板中的用法说明前缀
	usageLinePrefix := tpl.UsagePrefix
	var usageLine string

	// 优先使用用户自定义用法
	if cmd.userInfo.usage != "" {
		usageLine = usageLinePrefix + cmd.userInfo.usage
	} else {
		// 获取命令的完整路径
		fullCmdPath := getFullCommandPath(cmd)
		usageLine = usageLinePrefix + fullCmdPath

		// 如果是主命令(父命令为nil)，使用全局选项模板
		if cmd.parentCmd == nil {
			// 添加子命令部分
			if len(cmd.subCmds) > 0 {
				usageLine += tpl.UseageGlobalOptions
				usageLine += tpl.UsageSubCmd
			}

			// 添加选项部分
			usageLine += tpl.UseageInfoWithOptions

		} else {
			// 子命令，使用子命令选项模板
			// 如果存在子命令，则添加子命令用法
			if len(cmd.subCmds) > 0 {
				usageLine += tpl.UsageSubCmd
			}

			// 添加选项部分
			usageLine += tpl.UseageInfoWithOptions
		}
	}

	buf.WriteString(usageLine)
}

// writeOptions 写入选项信息
// cmd: 当前命令
// tpl: 模板实例
// buf: 输出缓冲区
func writeOptions(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 获取所有标志信息并排序
	flags := collectFlags(cmd)
	// 如果没有标志，不显示选项部分
	if len(flags) == 0 {
		return
	}
	// 写入选项头
	buf.WriteString(tpl.OptionsHeader)
	sortFlags(flags)

	// 计算描述信息对齐位置
	maxWidth, err := calculateMaxWidth(flags)
	if err != nil {
		// 如果计算宽度出错，则使用默认宽度
		maxWidth = 30
	}
	descStartPos := maxWidth + 5 // 增加5个空格作为间距

	for _, flag := range flags {
		// 格式化选项部分
		optPart := ""

		// 如果标志有短标志, 则使用模板1, 否则使用模板2
		if flag.shortFlag != "" {
			// --longFlag, -shortFlag <type>
			optPart = fmt.Sprintf(tpl.Option1, flag.longFlag, flag.shortFlag, flag.typeStr)
		} else {
			// --longFlag <type>
			optPart = fmt.Sprintf(tpl.Option2, flag.longFlag, flag.typeStr)
		}

		// 计算选项部分需要的填充空格
		padding := descStartPos - len(optPart)
		if padding < 1 {
			padding = 1
		}

		// 格式化整行输出
		fmt.Fprintf(buf, tpl.OptionDefault,
			optPart, padding, "", flag.usage, flag.defValue)
	}
}

// collectFlags 收集所有标志信息
func collectFlags(cmd *Cmd) []FlagInfo {
	var flags []FlagInfo

	// 遍历所有标志, 收集标志信息
	for _, f := range cmd.flagRegistry.allFlags {
		flag := f
		flags = append(flags, FlagInfo{
			longFlag:  flag.GetLongName(),
			shortFlag: flag.GetShortName(),
			usage:     flag.GetUsage(),
			defValue:  fmt.Sprintf("%v", flag.GetDefault()),
			typeStr:   flagTypeToString(flag.GetFlagType()),
		})
	}

	return flags
}

// sortFlags 按短标志字母顺序排序，有短标志的选项优先
func sortFlags(flags []FlagInfo) {
	sort.Slice(flags, func(i, j int) bool {
		a, b := flags[i], flags[j]
		return sortWithShortNamePriority(
			a.shortFlag != "",
			b.shortFlag != "",
			a.longFlag,
			b.longFlag,
			a.shortFlag,
			b.shortFlag,
		)
	})
}

// sortWithShortNamePriority 通用排序函数
// 排序优先级: 1.有短名称的优先 2.按长名称字母序 3.短名称字母序
// aHasShort: a是否有短名称
// bHasShort: b是否有短名称
// aName: a的长名称
// bName: b的长名称
// aShort: a的短名称
// bShort: b的短名称
func sortWithShortNamePriority(aHasShort, bHasShort bool, aName, bName, aShort, bShort string) bool {
	// 1. 有短名称的优先
	if aHasShort != bHasShort {
		return aHasShort
	}

	// 2. 按长名称字母顺序排序
	if aName != bName {
		return aName < bName
	}

	// 3. 都有短名称则按短名称字母顺序排序
	return aShort < bShort
}

// writeSubCmds 写入子命令信息
// cmd: 当前命令
// tpl: 模板实例
// buf: 输出缓冲区
func writeSubCmds(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 没有子命令则返回
	if len(cmd.subCmds) == 0 {
		return
	}

	// 添加子命令标题
	buf.WriteString(tpl.SubCmdsHeader)

	// 排序子命令：
	// 1. 按长命令名首字母排序
	// 2. 有短命令名的优先
	// 3. 只有长命令名的排最后
	sortedSubCmds := make([]*Cmd, len(cmd.subCmds))

	// 拷贝子命令到临时切片
	copy(sortedSubCmds, cmd.subCmds)

	// 排序子命令
	sort.Slice(sortedSubCmds, func(i, j int) bool {
		a, b := sortedSubCmds[i], sortedSubCmds[j]
		return sortWithShortNamePriority(
			a.ShortName() != "",
			b.ShortName() != "",
			a.LongName(),
			b.LongName(),
			a.ShortName(),
			b.ShortName(),
		)
	})

	// 计算最大命令名长度用于对齐
	maxNameLen := 0
	for _, subCmd := range sortedSubCmds {
		nameLen := len(subCmd.fs.Name())
		// 如果子命令有短名称，则计算短名称长度
		if subCmd.ShortName() != "" {
			nameLen += len(subCmd.ShortName()) + 5 // 添加5个空格, 保持命令对齐
		}
		// 更新最大命令名长度
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
	}

	// 生成对齐的子命令信息
	for _, subCmd := range sortedSubCmds {
		namePart := subCmd.fs.Name()
		if subCmd.ShortName() != "" {
			namePart = fmt.Sprintf("%s, %s", subCmd.fs.Name(), subCmd.ShortName())
		}

		// 格式化输出，确保描述信息对齐
		fmt.Fprintf(buf, "  %-*s\t%s\n", maxNameLen, namePart, subCmd.userInfo.description)
	}
}

// calculateMaxWidth 计算最大标志名称宽度用于对齐
func calculateMaxWidth(flags []FlagInfo) (int, error) {
	// 如果没有标志，则返回0
	if len(flags) == 0 {
		return 0, nil
	}

	maxWidth := 0
	for _, flag := range flags {
		var nameLength int
		if flag.shortFlag != "" {
			// 格式: "--longFlag, -shortFlag <type>"
			nameLength = len(flag.longFlag) + len(flag.shortFlag) + len(flag.typeStr) + 8 // 2('--') + 1(' ') + 2('<type>') + 2(', ') + 1('-')
		} else {
			// 格式: "--longFlag <type>"
			nameLength = len(flag.longFlag) + len(flag.typeStr) + 5 // 2('--') + 1(' ') + 2('<type>')
		}

		// 如果名称长度大于最大宽度，则更新最大宽度
		if nameLength > maxWidth {
			maxWidth = nameLength
		}
	}

	return maxWidth, nil
}

// writeExamples 写入示例信息
// cmd: 当前命令
// tpl: 模板实例
// buf: 输出缓冲区
func writeExamples(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 如果没有示例信息，则返回
	examples := cmd.GetExamples()
	if len(examples) == 0 {
		return
	}

	// 添加示例信息标题
	buf.WriteString(tpl.ExamplesHeader)

	// 遍历添加示例信息
	for i, example := range examples {
		// 格式化示例信息
		fmt.Fprintf(buf, tpl.ExampleItem, i+1, example.Description, example.Usage)

		// 如果不是最后一个示例，添加空行
		if i < len(examples)-1 {
			fmt.Fprintln(buf)
		}
	}
}

// writeNotes 写入注意事项
func writeNotes(cmd *Cmd, tpl HelpTemplate, buf *bytes.Buffer) {
	// 如果没有注意事项，则返回
	notes := cmd.GetNotes()
	if len(notes) == 0 {
		return
	}

	// 添加注意事项标题
	buf.WriteString(tpl.NotesHeader)

	// 遍历添加注意事项
	for i, note := range notes {
		fmt.Fprintf(buf, tpl.NoteItem, i+1, note)
	}
}

// hasCycle 检测命令间是否存在循环引用
// 采用深度优先搜索(DFS)算法，通过访问标记避免重复检测
// 参数:
//
//	parent: 当前命令
//	child: 待添加的子命令
//
// 返回值:
//
//	如果存在循环引用则返回true
func hasCycle(parent, child *Cmd) bool {
	if parent == nil || child == nil {
		return false
	}

	visited := make(map[*Cmd]bool)
	return dfs(parent, child, visited)
}

// dfs 深度优先搜索检测循环引用
func dfs(target, current *Cmd, visited map[*Cmd]bool) bool {
	// 如果已访问过当前节点，直接返回避免无限循环
	if visited[current] {
		return false
	}
	visited[current] = true

	// 找到目标节点，存在循环引用
	if current == target {
		return true
	}

	// 递归检查所有子命令
	for _, subCmd := range current.subCmds {
		if dfs(target, subCmd, visited) {
			return true
		}
	}

	// 检查父命令链
	if current.parentCmd != nil {
		return dfs(target, current.parentCmd, visited)
	}

	return false
}

// joinErrors 将错误切片合并为单个错误，并去除重复错误
func joinErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}
	if len(errors) == 1 {
		return errors[0]
	}

	// 使用map去重
	uniqueErrors := make(map[string]error)
	for _, err := range errors {
		errStr := err.Error()
		if _, exists := uniqueErrors[errStr]; !exists {
			uniqueErrors[errStr] = err
		}
	}

	// 构建错误信息
	var b strings.Builder
	b.WriteString(fmt.Sprintf("A total of %d unique errors:\n", len(uniqueErrors)))
	i := 1
	for _, err := range uniqueErrors {
		b.WriteString(fmt.Sprintf("  %d. %v\n", i, err))
		i++
	}

	// 使用常量格式字符串，将错误信息作为参数传入
	return fmt.Errorf("Merged error message:\n%s", b.String())
}

// getFullCommandPath 递归构建完整的命令路径，从根命令到当前命令
func getFullCommandPath(cmd *Cmd) string {
	if cmd.parentCmd == nil {
		return cmd.fs.Name()
	}
	return getFullCommandPath(cmd.parentCmd) + " " + cmd.fs.Name()
}

// GetExecutablePath 获取程序的绝对安装路径
// 如果无法通过 os.Executable 获取路径,则使用 os.Args[0] 作为替代
// 返回：程序的绝对路径字符串
func GetExecutablePath() string {
	// 尝试使用 os.Executable 获取可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		// 如果 os.Executable 报错,使用 os.Args[0] 作为替代
		exePath = os.Args[0]
	}
	// 使用 filepath.Abs 确保路径是绝对路径
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		// 如果 filepath.Abs 报错,直接返回原始路径
		return exePath
	}
	return absPath
}

// flagTypeToString 将FlagType转换为字符串
func flagTypeToString(flagType FlagType) string {
	switch flagType {
	case FlagTypeInt:
		return "<int>"
	case FlagTypeString:
		return "<string>"
	case FlagTypeBool:
		// 布尔类型没有参数类型字符串
		return ""
	case FlagTypeFloat:
		return "<float>"
	case FlagTypeEnum:
		return "<enum>"
	case FlagTypeDuration:
		return "<duration>"
	default:
		return "<unknown>"
	}
}
