// Package qflag 提供对标准库flag的封装，自动实现长短标志，并默认绑定-h/--help标志打印帮助信息。
// 用户可通过Cmd.Help字段自定义帮助内容，支持直接赋值字符串或从文件加载。
// 该包是一个功能强大的命令行参数解析库，支持子命令、多种数据类型标志、参数验证等高级特性。
package qflag

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	"gitee.com/MM-Q/qflag/cmd"
	"gitee.com/MM-Q/qflag/flags"
)

/*
项目地址: https://gitee.com/MM-Q/qflag
*/

var (
	// qCommandLine 全局默认Command实例（保持原名，与标准库flag对齐）
	qCommandLine *cmd.Cmd

	// qCommandLineOnce 确保全局默认Cmd实例只被初始化一次
	qCommandLineOnce sync.Once
)

// getQCommandLine 获取全局默认命令实例（延迟初始化）
// 内部函数, 确保QCommandLine只会被初始化一次, 线程安全
func getQCommandLine() *cmd.Cmd {
	qCommandLineOnce.Do(func() {
		// 使用一致的命令名生成逻辑
		cmdName := "myapp"
		if len(os.Args) > 0 {
			cmdName = filepath.Base(os.Args[0])
		}

		// 创建全局默认Cmd实例
		qCommandLine = cmd.NewCmd(cmdName, "", flag.ExitOnError)
	})
	return qCommandLine
}

// ================================================================================
// 操作方法 - 解析与管理 (17个)
// ================================================================================

// Parse 完整解析命令行参数（含子命令处理）
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 自动检测并解析子命令及其参数（若存在）
//  3. 验证枚举类型标志的有效性
//
// 参数：
//   - args: 原始命令行参数切片（包含可能的子命令及参数）
//
// 返回值：
//   - err: 解析过程中遇到的错误（如标志格式错误、子命令解析失败等）
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次（线程安全）
//   - 若检测到子命令，会将剩余参数传递给子命令的Parse方法
//   - 处理内置标志执行逻辑
func Parse() error {
	return getQCommandLine().Parse(os.Args[1:])
}

// ParseFlagsOnly 仅解析当前命令的标志参数（忽略子命令）
// 主要功能：
//  1. 解析当前命令的长短标志及内置标志
//  2. 验证枚举类型标志的有效性
//  3. 明确忽略所有子命令及后续参数
//
// 参数：
//   - args: 原始命令行参数切片（子命令及后续参数会被忽略）
//
// 返回值：
//   - err: 解析过程中遇到的错误（如标志格式错误等）
//
// 注意事项：
//   - 每个Cmd实例仅会被解析一次（线程安全）
//   - 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数
//   - 处理内置标志逻辑
func ParseFlagsOnly() error {
	return getQCommandLine().ParseFlagsOnly(os.Args[1:])
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
	return getQCommandLine().AddSubCmd(subCmds...)
}

// Args 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数切片。
// 非标志参数是指命令行中未被识别为标志的参数
//
// 返回值:
//   - []string: 包含所有非标志参数的字符串切片。
func Args() []string {
	return getQCommandLine().Args()
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
	return getQCommandLine().Arg(i)
}

// NArg 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数的数量
//
// 返回值:
//   - int: 非标志参数的数量。
func NArg() int {
	return getQCommandLine().NArg()
}

// NFlag 获取全局默认命令实例 `QCommandLine` 解析后已定义和使用的标志的数量
//
// 返回值:
//   - int: 标志的数量。
func NFlag() int {
	return getQCommandLine().NFlag()
}

// PrintHelp 输出全局默认命令实例 `QCommandLine` 的帮助信息。
// 帮助信息通常包含命令的名称、可用的标志及其描述等内容。
func PrintHelp() {
	getQCommandLine().PrintHelp()
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
	return getQCommandLine().FlagExists(name)
}

// Name 获取全局默认命令实例 `QCommandLine` 的名称
//
// 返回值:
//   - 优先返回长名称, 如果长名称不存在则返回短名称
func Name() string {
	return getQCommandLine().Name()
}

// LongName 获取命令长名称
func LongName() string {
	return getQCommandLine().LongName()
}

// ShortName 获取命令短名称
func ShortName() string {
	return getQCommandLine().ShortName()
}

// SubCmds 获取所有已注册的子命令列表
func SubCmds() []*cmd.Cmd {
	return getQCommandLine().SubCmds()
}

// SubCmdMap 获取所有已注册的子命令映射
func SubCmdMap() map[string]*cmd.Cmd {
	return getQCommandLine().SubCmdMap()
}

// CmdExists 检查子命令是否存在
//
// 参数:
//   - cmdName: 子命令名称
//
// 返回:
//   - bool: 子命令是否存在
func CmdExists(cmdName string) bool {
	return getQCommandLine().CmdExists(cmdName)
}

// IsParsed 检查命令行参数是否已解析
//
// 返回:
//   - bool: 是否已解析
func IsParsed() bool {
	return getQCommandLine().IsParsed()
}

// FlagRegistry 获取标志注册表
//
// 返回值:
//   - *flags.FlagRegistry: 标志注册表
func FlagRegistry() *flags.FlagRegistry {
	return getQCommandLine().FlagRegistry()
}

// ================================================================================
// Get 方法 - 获取配置信息
// ================================================================================

// GetVersion 获取全局默认命令的版本信息
//
// 返回值：
//   - string: 版本信息字符串。
func GetVersion() string {
	return getQCommandLine().GetVersion()
}

// GetDescription 获取命令描述信息
func GetDescription() string {
	return getQCommandLine().GetDescription()
}

// GetNotes 获取所有备注信息
func GetNotes() []string {
	return getQCommandLine().GetNotes()
}

// GetUseChinese 获取是否使用中文
// 该函数用于获取当前命令行标志是否使用中文
//
// 返回值:
//   - bool: 如果使用中文, 则返回true; 否则返回false。
func GetUseChinese() bool {
	return getQCommandLine().GetUseChinese()
}

// GetExamples 获取示例信息
// 该函数用于获取命令行标志的示例信息列表
//
// 返回值:
//   - []cmd.ExampleInfo: 示例信息列表，每个元素为 ExampleInfo 类型。
func GetExamples() []cmd.ExampleInfo {
	return getQCommandLine().GetExamples()
}

// GetHelp 返回全局默认命令实例 `QCommandLine` 的帮助信息
//
// 返回值:
//   - string: 命令行帮助信息。
func GetHelp() string {
	return getQCommandLine().GetHelp()
}

// GetLogoText 获取全局默认命令实例 `QCommandLine` 的 logo 文本
//
// 返回值:
//   - string: 配置的 logo 文本。
func GetLogoText() string {
	return getQCommandLine().GetLogoText()
}

// GetUsageSyntax 获取全局默认命令实例 `QCommandLine` 的用法信息
//
// 返回值:
//   - string: 命令行用法信息。
func GetUsageSyntax() string {
	return getQCommandLine().GetUsageSyntax()
}

// GetModuleHelps 获取模块帮助信息
//
// 返回值:
//   - string: 模块帮助信息。
func GetModuleHelps() string {
	return getQCommandLine().GetModuleHelps()
}

// ================================================================================
// Set 方法 - 设置配置信息(14个)
// ================================================================================

// SetVersion 为全局默认命令设置版本信息
//
// 参数说明：
//   - version: 版本信息字符串，用于标识命令的版本。
func SetVersion(version string) {
	getQCommandLine().SetVersion(version)
}

// SetVersionf 为全局默认命令设置版本信息
//
// 参数说明：
//   - format: 格式化字符串，用于标识命令的版本。
//   - args: 可变参数列表，用于替换格式化字符串中的占位符。
func SetVersionf(format string, args ...any) {
	getQCommandLine().SetVersionf(format, args...)
}

// SetDescription 设置命令描述信息
func SetDescription(desc string) {
	getQCommandLine().SetDescription(desc)
}

// SetUseChinese 设置是否使用中文
// 该函数用于设置当前命令行标志是否使用中文
//
// 参数:
//   - useChinese: 如果使用中文,则传入true;否则传入false。
func SetUseChinese(useChinese bool) {
	getQCommandLine().SetUseChinese(useChinese)
}

// AddNote 添加注意事项
// 该函数用于添加命令行标志的注意事项，这些注意事项将在命令行帮助信息中显示
//
// 参数:
//   - note: 注意事项内容，字符串类型。
func AddNote(note string) {
	getQCommandLine().AddNote(note)
}

// AddNotes 添加注意事项
// 该函数用于添加命令行标志的注意事项，这些注意事项将在命令行帮助信息中显示
//
// 参数:
//   - notes: 注意事项内容，字符串切片，每个元素为一个注意事项。
func AddNotes(notes []string) {
	getQCommandLine().AddNotes(notes)
}

// AddExample 添加示例
// 该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示
//
// 参数:
//   - desc: 示例描述，字符串类型。
//   - usage: 示例用法，字符串类型。
func AddExample(desc, usage string) {
	getQCommandLine().AddExample(desc, usage)
}

// AddExamples 添加示例
// 该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示
//
// 参数:
//   - examples: 示例列表，每个元素为 ExampleInfo 类型。
func AddExamples(examples []cmd.ExampleInfo) {
	getQCommandLine().AddExamples(examples)
}

// SetHelp 配置全局默认命令实例 `QCommandLine` 的帮助信息
//
// 参数:
//   - help: 新的帮助信息，字符串类型。
func SetHelp(help string) {
	getQCommandLine().SetHelp(help)
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
	getQCommandLine().SetUsageSyntax(usageSyntax)
}

// SetLogoText 配置全局默认命令实例 `QCommandLine` 的 logo 文本
//
// 参数:
//   - logoText: 配置的 logo 文本，字符串类型。
func SetLogoText(logoText string) {
	getQCommandLine().SetLogoText(logoText)
}

// SetModuleHelps 配置模块帮助信息
//
// 参数:
//   - moduleHelps: 模块帮助信息，字符串类型。
func SetModuleHelps(moduleHelps string) {
	getQCommandLine().SetModuleHelps(moduleHelps)
}

// SetExitOnBuiltinFlags 设置是否在解析内置参数时退出
// 默认情况下为true，当解析到内置参数时，QFlag将退出程序
//
// 参数:
//   - exit: 是否退出
func SetExitOnBuiltinFlags(exit bool) {
	getQCommandLine().SetExitOnBuiltinFlags(exit)
}

// SetEnableCompletion 设置是否启用自动完成功能
//
// 参数:
//   - enable: 是否启用自动完成功能
func SetEnableCompletion(enable bool) {
	getQCommandLine().SetEnableCompletion(enable)
}

// ================================================================================
// 链式调用方法 - 用于构建器模式，提供更流畅的API体验 (14个)
// ================================================================================

// WithDescription 设置命令描述（链式调用）
//
// 参数:
//   - desc: 命令描述
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithDescription(desc string) *cmd.Cmd {
	return getQCommandLine().WithDescription(desc)
}

// WithVersion 设置版本信息（链式调用）
//
// 参数:
//   - version: 版本信息
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithVersion(version string) *cmd.Cmd {
	return getQCommandLine().WithVersion(version)
}

// WithVersionf 设置版本信息（链式调用，支持格式化）
//
// 参数:
//   - format: 版本信息格式字符串
//   - args: 格式化参数
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithVersionf(format string, args ...any) *cmd.Cmd {
	return getQCommandLine().WithVersionf(format, args...)
}

// WithUseChinese 设置是否使用中文帮助信息（链式调用）
//
// 参数:
//   - useChinese: 是否使用中文帮助信息
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithUseChinese(useChinese bool) *cmd.Cmd {
	return getQCommandLine().WithUseChinese(useChinese)
}

// WithUsageSyntax 设置自定义命令用法（链式调用）
//
// 参数:
//   - usageSyntax: 自定义命令用法
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithUsageSyntax(usageSyntax string) *cmd.Cmd {
	return getQCommandLine().WithUsageSyntax(usageSyntax)
}

// WithLogoText 设置logo文本（链式调用）
//
// 参数:
//   - logoText: logo文本字符串
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithLogoText(logoText string) *cmd.Cmd {
	return getQCommandLine().WithLogoText(logoText)
}

// WithHelp 设置用户自定义命令帮助信息（链式调用）
//
// 参数:
//   - help: 用户自定义命令帮助信息
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithHelp(help string) *cmd.Cmd {
	return getQCommandLine().WithHelp(help)
}

// WithNote 添加备注信息到命令（链式调用）
//
// 参数:
//   - note: 备注信息
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithNote(note string) *cmd.Cmd {
	return getQCommandLine().WithNote(note)
}

// WithNotes 添加备注信息切片到命令（链式调用）
//
// 参数:
//   - notes: 备注信息列表
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithNotes(notes []string) *cmd.Cmd {
	return getQCommandLine().WithNotes(notes)
}

// WithExample 为命令添加使用示例（链式调用）
//
// 参数:
//   - desc: 示例描述
//   - usage: 示例用法
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithExample(desc, usage string) *cmd.Cmd {
	return getQCommandLine().WithExample(desc, usage)
}

// WithExamples 添加使用示例列表到命令（链式调用）
//
// 参数:
//   - examples: 示例信息列表，每个元素为 ExampleInfo 类型。
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithExamples(examples []cmd.ExampleInfo) *cmd.Cmd {
	return getQCommandLine().WithExamples(examples)
}

// WithExitOnBuiltinFlags 设置是否在解析内置参数时退出（链式调用）
//
// 参数:
//   - exit: 是否退出
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithExitOnBuiltinFlags(exit bool) *cmd.Cmd {
	return getQCommandLine().WithExitOnBuiltinFlags(exit)
}

// WithEnableCompletion 设置是否启用自动补全（链式调用）
//
// 参数:
//   - enable: true表示启用补全,false表示禁用
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithEnableCompletion(enable bool) *cmd.Cmd {
	return getQCommandLine().WithEnableCompletion(enable)
}

// WithModuleHelps 设置自定义模块帮助信息（链式调用）
//
// 参数:
//   - moduleHelps: 自定义模块帮助信息
//
// 返回值:
//   - *cmd.Cmd: 返回命令实例，支持链式调用
func WithModuleHelps(moduleHelps string) *cmd.Cmd {
	return getQCommandLine().WithModuleHelps(moduleHelps)
}
