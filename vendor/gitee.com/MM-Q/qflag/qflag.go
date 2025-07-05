// Package qflag 提供对标准库flag的封装, 自动实现长短标志, 并默认绑定-h/--help标志打印帮助信息。
// 用户可通过Cmd.Help字段自定义帮助内容, 支持直接赋值字符串或从文件加载。
package qflag

import (
	"os"
	"time"

	"gitee.com/MM-Q/qflag/cmd"
	"gitee.com/MM-Q/qflag/flags"
)

/*
项目地址: https://gitee.com/MM-Q/qflag
*/

// QCommandLineInterface 定义了全局默认命令行接口，提供统一的命令行参数管理功能
// 该接口封装了命令行程序的常用操作，包括标志添加、参数解析和帮助信息展示
type QCommandLineInterface interface {
	// 元数据操作方法
	Name() string                             // 获取命令名称
	LongName() string                         // 获取命令长名称
	ShortName() string                        // 获取命令短名称
	GetDescription() string                   // 获取命令描述信息
	SetDescription(desc string)               // 设置命令描述信息
	GetHelp() string                          // 获取命令帮助信息
	SetHelp(help string)                      // 设置命令帮助信息
	LoadHelp(filepath string) error           // 从指定文件加载帮助信息
	SetUsageSyntax(usageSyntax string)        // 设置命令用法格式
	GetUsageSyntax() string                   // 获取命令用法格式
	GetUseChinese() bool                      // 获取是否使用中文帮助信息
	SetUseChinese(useChinese bool)            // 设置是否使用中文帮助信息
	AddSubCmd(subCmd *cmd.Cmd)                // 添加子命令，子命令会继承父命令的上下文
	SubCmds() []*cmd.Cmd                      // 获取所有已注册的子命令列表
	Args() []string                           // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string                         // 获取指定索引的非标志参数，索引越界返回空字符串
	NArg() int                                // 获取非标志参数的数量
	NFlag() int                               // 获取已解析的标志数量
	PrintHelp()                               // 打印命令帮助信息
	FlagExists(name string) bool              // 检查指定名称的标志是否存在(支持长/短名称)
	AddNote(note string)                      // 添加一个注意事项
	GetNotes() []string                       // 获取所有备注信息
	AddExample(e cmd.ExampleInfo)             // 添加一个示例信息
	GetExamples() []cmd.ExampleInfo           // 获取示例信息列表
	SetVersion(version string)                // 设置版本信息
	GetVersion() string                       // 获取版本信息
	SetLogoText(logoText string)              // 设置logo文本
	GetLogoText() string                      // 获取logo文本
	SetModuleHelps(moduleHelps string)        // 设置自定义模块帮助信息
	GetModuleHelps() string                   // 获取自定义模块帮助信息
	SetExitOnBuiltinFlags(exit bool) *cmd.Cmd // 设置是否在处理内置标志时退出
	SetDisableBuiltinFlags(disable bool) *Cmd // 设置是否禁用内置标志注册
	CmdExists(cmdName string) bool            // 检查指定名称的命令是否存在

	// 标志解析方法
	Parse() error          // 解析命令行参数，自动处理标志和子命令
	ParseFlagsOnly() error // 解析命令行参数，仅处理标志，不处理子命令
	IsParsed() bool        // 检查是否已解析命令行参数

	// 添加标志方法
	String(longName, shortName, defValue, usage string) *flags.StringFlag                                // 添加字符串类型标志
	Int(longName, shortName string, defValue int, usage string) *flags.IntFlag                           // 添加整数类型标志
	Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag                        // 添加布尔类型标志
	Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag               // 添加浮点数类型标志
	Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag       // 添加时间间隔类型标志
	Enum(longName, shortName string, defValue string, usage string, enumValues []string) *flags.EnumFlag // 添加枚举类型标志
	Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag                  // 添加字符串切片类型标志
	Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag                     // 添加64位整型类型标志
	Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag                  // 添加无符号16位整型类型标志
	Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag                   // 添加时间类型标志
	Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag             // 添加Map标志
	Path(longName, shortName string, defValue string, usage string) *flags.PathFlag                      // 添加路径标志
	Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag                  // 添加无符号32位整型类型标志
	Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag                  // 添加无符号64位整型类型标志
	IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag                        // 添加IPv4地址标志
	IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag                        // 添加IPv6地址标志
	URL(longName, shortName string, defValue string, usage string) *flags.URLFlag                        // 添加URL标志

	// 绑定变量方法
	StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)                                // 绑定字符串标志到指定变量
	IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)                           // 绑定整数标志到指定变量
	BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)                        // 绑定布尔标志到指定变量
	Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)               // 绑定浮点数标志到指定变量
	DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)       // 绑定时间间隔类型标志到指定变量
	EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string) // 绑定枚举标志到指定变量
	SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)                  // 绑定字符串切片标志到指定变量
	Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)                     // 绑定64位整型标志到指定变量
	Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)                  // 绑定16位无符号整型标志到指定变量
	TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)                   // 绑定时间类型标志到指定变量
	MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)             // 绑定字符串映射标志到指定变量
	PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)                      // 绑定路径标志到指定变量
	Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string)                  // 绑定无符号32位整型标志到指定变量
	Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string)                  // 绑定无符号64位整型标志到指定变量
	IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string)                        // 绑定IPv4地址标志到指定变量
	IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string)                        // 绑定IPv6地址标志到指定变量
	URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string)                        // 绑定URL标志到指定变量
}

// SetVersion 为全局默认命令设置版本信息
//
// 参数说明：
//   - version: 版本信息字符串，用于标识命令的版本。
func SetVersion(version string) {
	QCommandLine.SetVersion(version)
}

// GetVersion 获取全局默认命令的版本信息
//
// 返回值：
//   - string: 版本信息字符串。
func GetVersion() string {
	return QCommandLine.GetVersion()
}

// Parse 完整解析命令行参数（含子命令处理）
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 自动检测并解析子命令及其参数（若存在）
//  3. 验证枚举类型标志的有效性
//
// 参数：
//
//	args: 原始命令行参数切片（包含可能的子命令及参数）
//
// 返回值：
//
//	解析过程中遇到的错误（如标志格式错误、子命令解析失败等）
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次（线程安全）
//   - 若检测到子命令，会将剩余参数传递给子命令的Parse方法
//   - 处理内置标志执行逻辑
func Parse() error {
	return QCommandLine.Parse(os.Args[1:])
}

// ParseFlagsOnly 仅解析当前命令的标志参数（忽略子命令）
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 验证枚举类型标志的有效性
//  3. 明确忽略所有子命令及后续参数
//
// 参数：
//
//	args: 原始命令行参数切片（子命令及后续参数会被忽略）
//
// 返回值：
//
//	解析过程中遇到的错误（如标志格式错误等）
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次（线程安全）
//   - 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数
//   - 处理内置标志逻辑
func ParseFlagsOnly() error {
	return QCommandLine.ParseFlagsOnly(os.Args[1:])
}

// AddSubCmd 向全局默认命令实例 `QCommandLine` 添加一个或多个子命令
// 该函数会调用全局默认命令实例的 `AddSubCmd` 方法，支持批量添加子命令
// 在添加过程中，会检查子命令是否为 `nil` 以及是否存在循环引用，若有异常则返回错误信息
//
// 参数:
//   - subCmds: 可变参数，接收一个或多个 `*Cmd` 类型的子命令实例
//
// 返回值:
//   - error: 若添加子命令过程中出现错误（如子命令为 `nil` 或存在循环引用），则返回错误信息；否则返回 `nil`。
func AddSubCmd(subCmds ...*cmd.Cmd) error {
	return QCommandLine.AddSubCmd(subCmds...)
}

// Args 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数切片。
// 非标志参数是指命令行中未被识别为标志的参数
//
// 返回值:
//   - []string: 包含所有非标志参数的字符串切片。
func Args() []string {
	return QCommandLine.Args()
}

// Arg 获取全局默认命令实例 `QCommandLine` 解析后的指定索引位置的非标志参数
// 索引从 0 开始，若索引超出非标志参数切片的范围，将返回空字符串
//
// 参数:
//   - i: 非标志参数的索引位置，从 0 开始计数
//
// 返回值:
//   - string: 指定索引位置的非标志参数；若索引越界，则返回空字符串
func Arg(i int) string {
	return QCommandLine.Arg(i)
}

// NArg 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数的数量
//
// 返回值:
//   - int: 非标志参数的数量。
func NArg() int {
	return QCommandLine.NArg()
}

// NFlag 获取全局默认命令实例 `QCommandLine` 解析后已定义和使用的标志的数量
//
// 返回值:
//   - int: 标志的数量。
func NFlag() int {
	return QCommandLine.NFlag()
}

// PrintHelp 输出全局默认命令实例 `QCommandLine` 的帮助信息。
// 帮助信息通常包含命令的名称、可用的标志及其描述等内容。
func PrintHelp() {
	QCommandLine.PrintHelp()
}

// FlagExists 检查全局默认命令实例 `QCommandLine` 中是否存在指定名称的标志
// 该函数会调用全局默认命令实例的 `FlagExists` 方法，用于检查命令行中是否存在指定名称的标志
//
// 参数:
//   - name: 要检查的标志名称，可以是长名称或短名称。
//
// 返回值:
//   - bool: 若存在指定名称的标志，则返回 `true`；否则返回 `false`。
func FlagExists(name string) bool {
	return QCommandLine.FlagExists(name)
}

// Name 获取全局默认命令实例 `QCommandLine` 的名称
//
// 返回值:
//   - 优先返回长名称, 如果长名称不存在则返回短名称
func Name() string {
	return QCommandLine.Name()
}

// LongName 获取命令长名称
func LongName() string {
	return QCommandLine.LongName()
}

// ShortName 获取命令短名称
func ShortName() string {
	return QCommandLine.ShortName()
}

// GetDescription 获取命令描述信息
func GetDescription() string {
	return QCommandLine.GetDescription()
}

// SetDescription 设置命令描述信息
func SetDescription(desc string) {
	QCommandLine.SetDescription(desc)
}

// GetNotes 获取所有备注信息
func GetNotes() []string {
	return QCommandLine.GetNotes()
}

// SubCmds 获取所有已注册的子命令列表
func SubCmds() []*cmd.Cmd {
	return QCommandLine.SubCmds()
}

// SubCmdMap 获取所有已注册的子命令映射
func SubCmdMap() map[string]*cmd.Cmd {
	return QCommandLine.SubCmdMap()
}

// GetUseChinese 获取是否使用中文
// 该函数用于获取当前命令行标志是否使用中文
//
// 返回值:
//   - bool: 如果使用中文, 则返回true; 否则返回false。
func GetUseChinese() bool {
	return QCommandLine.GetUseChinese()
}

// SetUseChinese 设置是否使用中文
// 该函数用于设置当前命令行标志是否使用中文
//
// 参数:
//   - useChinese: 如果使用中文,则传入true;否则传入false。
func SetUseChinese(useChinese bool) {
	QCommandLine.SetUseChinese(useChinese)
}

// AddNote 添加注意事项
// 该函数用于添加命令行标志的注意事项，这些注意事项将在命令行帮助信息中显示
//
// 参数:
//   - note: 注意事项内容，字符串类型。
func AddNote(note string) {
	QCommandLine.AddNote(note)
}

// AddExample 添加示例
// 该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示
//
// 参数:
//   - e: 示例信息，ExampleInfo 类型。
func AddExample(e cmd.ExampleInfo) {
	QCommandLine.AddExample(e)
}

// GetExamples 获取示例信息
// 该函数用于获取命令行标志的示例信息列表
//
// 返回值:
//   - []ExampleInfo: 示例信息列表，每个元素为 ExampleInfo 类型。
func GetExamples() []cmd.ExampleInfo {
	return QCommandLine.GetExamples()
}

// GetHelp 返回全局默认命令实例 `QCommandLine` 的帮助信息
//
// 返回值:
//   - string: 命令行帮助信息。
func GetHelp() string {
	return QCommandLine.GetHelp()
}

// SetHelp 配置全局默认命令实例 `QCommandLine` 的帮助信息
//
// 参数:
//   - help: 新的帮助信息，字符串类型。
func SetHelp(help string) {
	QCommandLine.SetHelp(help)
}

// LoadHelp 从文件中加载帮助信息
//
// 参数:
//   - filepath: 文件路径，字符串类型。
//
// 返回值:
//   - error: 如果加载失败，则返回错误信息；否则返回 nil。
//
// 示例:
//
//	qflag.LoadHelp("help.txt")
func LoadHelp(filepath string) error {
	return QCommandLine.LoadHelp(filepath)
}

// SetUsageSyntax 配置全局默认命令实例 `QCommandLine` 的用法信息
//
// 参数:
//   - usage: 新的用法信息，字符串类型。
//
// 示例:
//
//	qflag.SetUsageSyntax("Usage: qflag [options]")
func SetUsageSyntax(usageSyntax string) {
	QCommandLine.SetUsageSyntax(usageSyntax)
}

// GetUsageSyntax 获取全局默认命令实例 `QCommandLine` 的用法信息
//
// 返回值:
//   - string: 命令行用法信息。
func GetUsageSyntax() string {
	return QCommandLine.GetUsageSyntax()
}

// SetLogoText 配置全局默认命令实例 `QCommandLine` 的 logo 文本
//
// 参数:
//   - logoText: 配置的 logo 文本，字符串类型。
func SetLogoText(logoText string) {
	QCommandLine.SetLogoText(logoText)
}

// GetLogoText 获取全局默认命令实例 `QCommandLine` 的 logo 文本
//
// 返回值:
//   - string: 配置的 logo 文本。
func GetLogoText() string {
	return QCommandLine.GetLogoText()
}

// SetModuleHelps 配置模块帮助信息
//
// 参数:
//   - moduleHelps: 模块帮助信息，字符串类型。
func SetModuleHelps(moduleHelps string) {
	QCommandLine.SetModuleHelps(moduleHelps)
}

// GetModuleHelps 获取模块帮助信息
//
// 返回值:
//   - string: 模块帮助信息。
func GetModuleHelps() string {
	return QCommandLine.GetModuleHelps()
}

// SetExitOnBuiltinFlags 设置是否在解析内置参数时退出
// 默认情况下为true，当解析到内置参数时，QFlag将退出程序
//
// 参数:
//   - exit: 是否退出
//
// 返回值:
//   - *cmd.Cmd: 当前命令对象
func SetExitOnBuiltinFlags(exit bool) *cmd.Cmd {
	QCommandLine.SetExitOnBuiltinFlags(exit)
	return QCommandLine
}

// SetDisableBuiltinFlags 设置是否禁用内置参数
// 默认情况下为false，当设置为true时，QFlag将忽略内置参数
//
// 参数:
//   - disable: 是否禁用
//
// 返回值:
//   - *cmd.Cmd: 当前命令对象
func SetDisableBuiltinFlags(disable bool) *cmd.Cmd {
	QCommandLine.SetDisableBuiltinFlags(disable)
	return QCommandLine
}

// CmdExists 检查子命令是否存在
//
// 参数:
//   - cmdName: 子命令名称
//
// 返回:
//   - bool: 子命令是否存在
func CmdExists(cmdName string) bool {
	return QCommandLine.CmdExists(cmdName)
}

// IsParsed 检查命令行参数是否已解析
//
// 返回:
//   - bool: 是否已解析
func IsParsed() bool {
	return QCommandLine.IsParsed()
}
