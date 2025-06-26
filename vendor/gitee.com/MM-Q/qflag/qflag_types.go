package qflag

// 定义非法字符集常量, 防止非法的标志名称
const invalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"

// FlagInfo 标志信息结构体
// 用于存储命令行标志的元数据，包括长名称、短名称、用法说明和默认值
type FlagInfo struct {
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

// 标志类型
type FlagType int

const (
	FlagTypeInt      FlagType = iota + 1 // 整数类型
	FlagTypeString                       // 字符串类型
	FlagTypeBool                         // 布尔类型
	FlagTypeFloat                        // 浮点数类型
	FlagTypeEnum                         // 枚举类型
	FlagTypeDuration                     // 时间间隔类型
	FlagTypeSlice                        // 切片类型
	FlagTypeMap                          // 映射类型
)

// 内置标志名称
var (
	helpFlagName            = "help"    // 帮助标志名称
	helpFlagShortName       = "h"       // 帮助标志短名称
	showInstallPathFlagName = "sip"     // 显示安装路径标志名称
	versionFlagLongName     = "version" // 版本标志名称
	versionFlagShortName    = "v"       // 版本标志短名称
)

// HelpTemplate 帮助信息模板结构体
type HelpTemplate struct {
	CmdName               string // 命令名称模板
	CmdNameWithShort      string // 命令名称带短名称模板
	CmdDescription        string // 命令描述模板
	UsagePrefix           string // 用法说明前缀模板
	UsageSubCmd           string // 用法说明子命令模板
	UseageInfoWithOptions string // 带选项的用法说明信息模板
	UseageGlobalOptions   string // 全局选项部分
	OptionsHeader         string // 选项头部模板
	Option1               string // 选项模板(带短选项)
	Option2               string // 选项模板(无短选项)
	Option3               string // 选项模板(无长选项)
	OptionDefault         string // 选项模板的默认值
	SubCmdsHeader         string // 子命令头部模板
	SubCmd                string // 子命令模板
	SubCmdWithShort       string // 子命令带短名称模板
	NotesHeader           string // 注意事项头部模板
	NoteItem              string // 注意事项项模板
	DefaultNote           string // 默认注意事项
	ExamplesHeader        string // 示例信息头部模板
	ExampleItem           string // 示例信息项模板
}

// 英文模板实例
var EnglishTemplate = HelpTemplate{
	CmdName:               "Name: %s\n\n",
	UsagePrefix:           "Usage: ",                                                                                                                             // 命令名称模板
	UsageSubCmd:           " [subcmd]",                                                                                                                           // 命令名称模板
	UseageInfoWithOptions: " [options]\n\n",                                                                                                                      // 带选项的用法说明信息模板
	UseageGlobalOptions:   " [global options]",                                                                                                                   // 全局选项部分
	CmdNameWithShort:      "Name: %s, %s\n\n",                                                                                                                    // 命令名称带短名称模板
	CmdDescription:        "Desc: %s\n\n",                                                                                                                        // 命令描述模板
	OptionsHeader:         "Options:\n",                                                                                                                          // 选项头部模板
	Option1:               "  --%s, -%s %s",                                                                                                                      // 选项模板(带短选项)
	Option2:               "  --%s %s",                                                                                                                           // 选项模板(无短选项)
	Option3:               "  -%s %s",                                                                                                                            // 新增：仅短选项的模板
	OptionDefault:         "%s%*s%s (default: %s)\n",                                                                                                             // 选项模板默认值
	SubCmdsHeader:         "\nSubCmds:\n",                                                                                                                        // 子命令头部模板
	SubCmd:                "  %s\t%s\n",                                                                                                                          // 子命令模板
	SubCmdWithShort:       "  %s, %s\t%s\n",                                                                                                                      // 子命令模板(带短选项)
	NotesHeader:           "\nNotes:\n",                                                                                                                          // 注意事项头部模板
	NoteItem:              "  %d. %s\n",                                                                                                                          // 注意事项模板
	DefaultNote:           "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.", // 默认注意事项
	ExamplesHeader:        "\nExamples:\n",                                                                                                                       // 示例信息头部模板
	ExampleItem:           "  %d. %s\n     %s\n",                                                                                                                 // 序号、描述、用法
}

// 中文模板实例
var ChineseTemplate = HelpTemplate{
	CmdName:               "名称: %s\n\n",                  // 命令名称模板
	UsagePrefix:           "用法: ",                        // 用法说明前缀模板
	UsageSubCmd:           " [子命令]",                      // 用法说明子命令模板
	UseageInfoWithOptions: " [选项]\n\n",                   // 带选项的用法说明信息模板
	UseageGlobalOptions:   " [全局选项]",                     // 全局选项部分
	CmdNameWithShort:      "名称: %s, %s\n\n",              // 命令名称带短名称模板
	CmdDescription:        "描述: %s\n\n",                  // 命令描述模板
	OptionsHeader:         "选项:\n",                       // 选项头部模板
	Option1:               "  --%s, -%s %s",              // 选项模板(带短选项)
	Option2:               "  --%s %s",                   // 选项模板(无短选项)
	Option3:               "  -%s %s",                    // 新增：仅短选项的模板
	OptionDefault:         "%s%*s%s (默认值: %s)\n",         // 选项模板默认值
	SubCmdsHeader:         "\n子命令:\n",                    // 子命令头部模板
	SubCmd:                "  %s\t%s\n",                  // 子命令模板
	SubCmdWithShort:       "  %s, %s\t%s\n",              // 子命令模板(带短选项)
	NotesHeader:           "\n注意事项:\n",                   //注意事项头部模板
	NoteItem:              "  %d、%s\n",                   //注意事项模板
	DefaultNote:           "当长选项和短选项同时使用时，最后指定的选项将优先生效。", //默认注意事项
	ExamplesHeader:        "\n示例:\n",                     // 示例信息头部模板
	ExampleItem:           "  %d、%s\n     %s\n",          // 序号、描述、用法
}

// 定义标志的分隔符常量
const (
	FlagSplitComma     = "," // 逗号
	FlagSplitSemicolon = ";" // 分号
	FlagSplitPipe      = "|" // 竖线
	FlagKVColon        = ":" // 冒号
	FlagKVEqual        = "=" // 等号
)

// Flag支持的标志分隔符切片
var FlagSplitSlice = []string{
	FlagSplitComma,     // 逗号
	FlagSplitSemicolon, // 分号
	FlagSplitPipe,      // 竖线
	FlagKVColon,        // 冒号
}
