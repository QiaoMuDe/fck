// cmd_genhelp 命令行帮助信息生成器
package cmd

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// flagInfo 标志信息结构体
// 用于存储命令行标志的元数据，包括长名称、短名称、用法说明和默认值
type flagInfo struct {
	longFlag  string // 长标志名称
	shortFlag string // 短标志名称
	usage     string // 用法说明
	defValue  string // 默认值
	typeStr   string // 参数类型字符串
}

// ExampleInfo 示例信息结构体
// 用于存储命令的使用示例，包括描述和示例内容
type ExampleInfo struct {
	Description string // 示例描述
	Usage       string // 示例使用方式
}

// HelpTemplate 帮助信息模板结构体
type HelpTemplate struct {
	CmdName              string // 命令名称模板
	CmdNameWithShort     string // 命令名称带短名称模板
	CmdDescription       string // 命令描述模板
	UsagePrefix          string // 用法说明前缀模板
	UsageSubCmd          string // 用法说明子命令模板
	UsageInfoWithOptions string // 带选项的用法说明信息模板
	UsageGlobalOptions   string // 全局选项部分
	OptionsHeader        string // 选项头部模板
	Option1              string // 选项模板(带短选项)
	Option2              string // 选项模板(无短选项)
	Option3              string // 选项模板(无长选项)
	OptionDefault        string // 选项模板的默认值
	SubCmdsHeader        string // 子命令头部模板
	SubCmd               string // 子命令模板
	SubCmdWithShort      string // 子命令带短名称模板
	NotesHeader          string // 注意事项头部模板
	NoteItem             string // 注意事项项模板
	DefaultNote          string // 默认注意事项
	ExamplesHeader       string // 示例信息头部模板
	ExampleItem          string // 示例信息项模板
}

// 英文模板实例
var EnglishTemplate = HelpTemplate{
	CmdName:              "Name: %s\n\n",
	UsagePrefix:          "Usage: ",                                                                                                                             // 命令名称模板
	UsageSubCmd:          " [subcmd]",                                                                                                                           // 命令名称模板
	UsageInfoWithOptions: " [options]\n\n",                                                                                                                      // 带选项的用法说明信息模板
	UsageGlobalOptions:   " [global options]",                                                                                                                   // 全局选项部分
	CmdNameWithShort:     "Name: %s, %s\n\n",                                                                                                                    // 命令名称带短名称模板
	CmdDescription:       "Desc: %s\n\n",                                                                                                                        // 命令描述模板
	OptionsHeader:        "Options:\n",                                                                                                                          // 选项头部模板
	Option1:              "  --%s, -%s %s",                                                                                                                      // 选项模板(带短选项)
	Option2:              "  --%s %s",                                                                                                                           // 选项模板(无短选项)
	Option3:              "  -%s %s",                                                                                                                            // 新增：仅短选项的模板
	OptionDefault:        "%s%*s%s (default: %s)\n",                                                                                                             // 选项模板默认值
	SubCmdsHeader:        "\nSubCmds:\n",                                                                                                                        // 子命令头部模板
	SubCmd:               "  %s\t%s\n",                                                                                                                          // 子命令模板
	SubCmdWithShort:      "  %s, %s\t%s\n",                                                                                                                      // 子命令模板(带短选项)
	NotesHeader:          "\nNotes:\n",                                                                                                                          // 注意事项头部模板
	NoteItem:             "  %d. %s\n",                                                                                                                          // 注意事项模板
	DefaultNote:          "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.", // 默认注意事项
	ExamplesHeader:       "\nExamples:\n",                                                                                                                       // 示例信息头部模板
	ExampleItem:          "  %d. %s\n     %s\n",                                                                                                                 // 序号、描述、用法
}

// 中文模板实例
var ChineseTemplate = HelpTemplate{
	CmdName:              "名称: %s\n\n",                  // 命令名称模板
	UsagePrefix:          "用法: ",                        // 用法说明前缀模板
	UsageSubCmd:          " [子命令]",                      // 用法说明子命令模板
	UsageInfoWithOptions: " [选项]\n\n",                   // 带选项的用法说明信息模板
	UsageGlobalOptions:   " [全局选项]",                     // 全局选项部分
	CmdNameWithShort:     "名称: %s, %s\n\n",              // 命令名称带短名称模板
	CmdDescription:       "描述: %s\n\n",                  // 命令描述模板
	OptionsHeader:        "选项:\n",                       // 选项头部模板
	Option1:              "  --%s, -%s %s",              // 选项模板(带短选项)
	Option2:              "  --%s %s",                   // 选项模板(无短选项)
	Option3:              "  -%s %s",                    // 新增：仅短选项的模板
	OptionDefault:        "%s%*s%s (默认值: %s)\n",         // 选项模板默认值
	SubCmdsHeader:        "\n子命令:\n",                    // 子命令头部模板
	SubCmd:               "  %s\t%s\n",                  // 子命令模板
	SubCmdWithShort:      "  %s, %s\t%s\n",              // 子命令模板(带短选项)
	NotesHeader:          "\n注意事项:\n",                   //注意事项头部模板
	NoteItem:             "  %d、%s\n",                   //注意事项模板
	DefaultNote:          "当长选项和短选项同时使用时，最后指定的选项将优先生效。", //默认注意事项
	ExamplesHeader:       "\n示例:\n",                     // 示例信息头部模板
	ExampleItem:          "  %d、%s\n     %s\n",          // 序号、描述、用法
}

// generateHelpInfo 生成帮助信息
//
// 参数:
//   - cmd: 当前命令
//
// 返回值:
//   - string: 帮助信息
func generateHelpInfo(cmd *Cmd) string {
	// 处理根命令
	if cmd == nil {
		return ""
	}

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
	// 修改后的命令名称显示逻辑
	if cmd.LongName() != "" && cmd.ShortName() != "" {
		// 同时有长短名称
		fmt.Fprintf(buf, tpl.CmdNameWithShort, cmd.LongName(), cmd.ShortName())
	} else if cmd.LongName() != "" {
		// 只有长名称
		fmt.Fprintf(buf, tpl.CmdName, cmd.LongName())
	} else {
		// 只有短名称
		fmt.Fprintf(buf, tpl.CmdName, cmd.ShortName())
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
	if cmd.userInfo.usageSyntax != "" {
		usageLine = usageLinePrefix + cmd.userInfo.usageSyntax + "\n"
	} else {
		// 获取命令的完整路径
		fullCmdPath := getFullCommandPath(cmd)
		usageLine = usageLinePrefix + fullCmdPath

		// 如果是主命令(父命令为nil)，使用全局选项模板
		if cmd.parentCmd == nil {
			// 添加子命令部分
			if len(cmd.subCmds) > 0 {
				usageLine += tpl.UsageGlobalOptions
				usageLine += tpl.UsageSubCmd
			}

			// 添加选项部分
			usageLine += tpl.UsageInfoWithOptions

		} else {
			// 子命令，使用子命令选项模板
			// 如果存在子命令，则添加子命令用法
			if len(cmd.subCmds) > 0 {
				usageLine += tpl.UsageSubCmd
			}

			// 添加选项部分
			usageLine += tpl.UsageInfoWithOptions
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

	// 遍历所有标志
	for _, flag := range flags {
		// 格式化选项部分
		optPart := ""

		// 根据标志生成选项部分
		if flag.longFlag != "" && flag.shortFlag != "" {
			// 同时有长短选项
			optPart = fmt.Sprintf(tpl.Option1, flag.longFlag, flag.shortFlag, flag.typeStr)
		} else if flag.longFlag != "" {
			// 只有长选项
			optPart = fmt.Sprintf(tpl.Option2, flag.longFlag, flag.typeStr)
		} else {
			// 只有短选项
			optPart = fmt.Sprintf(tpl.Option3, flag.shortFlag, flag.typeStr)
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
func collectFlags(cmd *Cmd) []flagInfo {
	var flagInfos []flagInfo

	// 遍历所有标志, 收集标志信息
	for _, f := range cmd.flagRegistry.GetAllFlagMetas() {
		flag := f // 获取标志类型

		// 收集默认值
		defValue := fmt.Sprintf("%v", flag.GetDefault())

		// 对Duration类型进行特殊格式化
		if flag.GetFlagType() == flags.FlagTypeDuration {
			if duration, ok := flag.GetDefault().(time.Duration); ok {
				defValue = duration.String() // 获取时间间隔标志的默认值的字符串表示
			}
		}

		// 创建标志元数据
		flagInfos = append(flagInfos, flagInfo{
			longFlag:  flag.GetLongName(),                         // 长标志名
			shortFlag: flag.GetShortName(),                        // 短标志
			usage:     flag.GetUsage(),                            // 使用说明
			defValue:  defValue,                                   // 默认值
			typeStr:   flags.FlagTypeToString(flag.GetFlagType()), // 标志类型字符串
		})
	}

	return flagInfos
}

// sortFlags 按短标志字母顺序排序，有短标志的选项优先
func sortFlags(flags []flagInfo) {
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

	// 获取子命令列表
	sortedSubCmds := cmd.SubCmds()

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
		// 构建子命令名称
		var namePart string

		// 根据子命令名称和短名称生成名称部分
		if subCmd.LongName() != "" && subCmd.ShortName() != "" {
			// 如果长短名称都存在, 则同时显示
			namePart = fmt.Sprintf("%s, %s", subCmd.LongName(), subCmd.ShortName())
		} else if subCmd.LongName() != "" {
			// 如果只有长短名称中的一个存在, 则只显示一个
			namePart = subCmd.LongName()
		} else if subCmd.ShortName() != "" {
			// 如果只有长短名称中的一个存在, 则只显示一个
			namePart = subCmd.ShortName()
		} else {
			// 如果长短名称都不存在, 则使用默认的注册名
			namePart = subCmd.fs.Name()
		}

		// 格式化输出，确保描述信息对齐
		fmt.Fprintf(buf, "  %-*s\t%s\n", maxNameLen, namePart, subCmd.userInfo.description)
	}
}

// calculateMaxWidth 计算最大标志名称宽度用于对齐
func calculateMaxWidth(flags []flagInfo) (int, error) {
	// 如果没有标志，则返回0
	if len(flags) == 0 {
		return 0, nil
	}

	maxWidth := 0
	for _, flag := range flags {
		var nameLength int
		if flag.longFlag != "" && flag.shortFlag != "" {
			// 同时有长短选项: --longFlag, -shortFlag <type>
			nameLength = len(flag.longFlag) + len(flag.shortFlag) + len(flag.typeStr) + 8 // 2('--') + 1(' ') + 2('<type>') + 2(', ') + 1('-')
		} else if flag.longFlag != "" {
			// 只有长选项: --longFlag <type>
			nameLength = len(flag.longFlag) + len(flag.typeStr) + 5 // 2('--') + 1(' ') + 2('<type>')
		} else {
			// 只有短选项: -shortFlag <type>
			nameLength = len(flag.shortFlag) + len(flag.typeStr) + 4 // 1('-') + 1(' ') + 2('<type>')
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

// getFullCommandPath 递归构建完整的命令路径，从根命令到当前命令
func getFullCommandPath(cmd *Cmd) string {
	if cmd.parentCmd == nil {
		return cmd.fs.Name()
	}
	return getFullCommandPath(cmd.parentCmd) + " " + cmd.fs.Name()
}
