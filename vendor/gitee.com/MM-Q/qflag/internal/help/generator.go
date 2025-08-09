// Package help 帮助信息生成器
// 本文件实现了命令行帮助信息的自动生成功能，包括标志列表、用法说明、
// 子命令信息等帮助内容的格式化和输出。
package help

import (
	"bytes"

	"gitee.com/MM-Q/qflag/internal/types"
)

// GenerateHelp 生成帮助信息
// 纯函数设计，不依赖任何结构体状态
func GenerateHelp(ctx *types.CmdContext) string {
	// 检查上下文是否为空
	if ctx == nil {
		return ""
	}

	// 检查是否配置了自定义帮助信息
	if ctx.Config.Help != "" {
		return ctx.Config.Help
	}

	// 根据语言选择模板实例
	var tpl HelpTemplate
	if ctx.Config.UseChinese {
		tpl = ChineseTemplate // 中文模板实例
	} else {
		tpl = EnglishTemplate // 英文模板实例
	}

	// 使用buffer提高字符串拼接性能
	var buf bytes.Buffer

	// 写入命令名称和描述
	writeCommandHeader(ctx, tpl, &buf)

	// 写入Logo信息
	writeLogoText(ctx, &buf)

	// 写入用法说明
	writeUsageLine(ctx, tpl, &buf)

	// 写入选项信息
	writeOptions(ctx, tpl, &buf)

	// 写入子命令信息
	writeSubCmds(ctx, tpl, &buf)

	// 写入自定义模块帮助信息
	writeModuleHelps(ctx, &buf)

	// 写入示例信息
	writeExamples(ctx, tpl, &buf)

	// 写入注意事项
	writeNotes(ctx, tpl, &buf)

	return buf.String()
}

// flagInfo 标志信息结构体
// 用于存储命令行标志的元数据，包括长名称、短名称、用法说明和默认值
type flagInfo struct {
	longFlag  string // 长标志名称
	shortFlag string // 短标志名称
	usage     string // 用法说明
	defValue  string // 默认值
	typeStr   string // 参数类型字符串
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
