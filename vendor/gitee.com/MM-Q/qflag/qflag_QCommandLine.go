package qflag

import (
	"flag"
	"os"
	"path/filepath"
	"time"
)

// QCommandLine 全局默认Cmd实例
var QCommandLine *Cmd

// QCommandLineInterface 定义了全局默认命令行接口，提供统一的命令行参数管理功能
// 该接口封装了命令行程序的常用操作，包括标志添加、参数解析和帮助信息展示
type QCommandLineInterface interface {
	// 元数据操作方法
	LongName() string                  // 获取命令长名称
	ShortName() string                 // 获取命令短名称
	GetDescription() string            // 获取命令描述信息
	SetDescription(desc string)        // 设置命令描述信息
	GetHelp() string                   // 获取命令帮助信息
	SetHelp(help string)               // 设置命令帮助信息
	SetUsageSyntax(usageSyntax string) // 设置命令用法格式
	GetUsageSyntax() string            // 获取命令用法格式
	GetUseChinese() bool               // 获取是否使用中文帮助信息
	SetUseChinese(useChinese bool)     // 设置是否使用中文帮助信息
	AddSubCmd(subCmd *Cmd)             // 添加子命令，子命令会继承父命令的上下文
	SubCmds() []*Cmd                   // 获取所有已注册的子命令列表
	Parse() error                      // 解析命令行参数，自动处理标志和子命令
	ParseFlagsOnly() error             // 解析命令行参数，仅处理标志，不处理子命令
	Args() []string                    // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string                  // 获取指定索引的非标志参数，索引越界返回空字符串
	NArg() int                         // 获取非标志参数的数量
	NFlag() int                        // 获取已解析的标志数量
	PrintHelp()                        // 打印命令帮助信息
	FlagExists(name string) bool       // 检查指定名称的标志是否存在(支持长/短名称)
	AddNote(note string)               // 添加一个注意事项
	GetNotes() []string                // 获取所有备注信息
	AddExample(e ExampleInfo)          // 添加一个示例信息
	GetExamples() []ExampleInfo        // 获取示例信息列表
	SetVersion(version string)         // 设置版本信息
	GetVersion() string                // 获取版本信息
	SetLogoText(logoText string)       // 设置logo文本
	GetLogoText() string               // 获取logo文本
	SetModuleHelps(moduleHelps string) // 设置自定义模块帮助信息
	GetModuleHelps() string            // 获取自定义模块帮助信息

	// 添加标志方法
	String(longName, shortName, defValue, usage string) *StringFlag                                // 添加字符串类型标志
	Int(longName, shortName string, defValue int, usage string) *IntFlag                           // 添加整数类型标志
	Bool(longName, shortName string, defValue bool, usage string) *BoolFlag                        // 添加布尔类型标志
	Float(longName, shortName string, defValue float64, usage string) *FloatFlag                   // 添加浮点数类型标志
	Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag       // 添加时间间隔类型标志
	Enum(longName, shortName string, defValue string, usage string, enumValues []string) *EnumFlag // 添加枚举类型标志
	Slice(longName, shortName string, defValue []string, usage string) *SliceFlag                  // 添加字符串切片类型标志

	// 绑定变量方法
	StringVar(f *StringFlag, longName, shortName, defValue, usage string)                                // 绑定字符串标志到指定变量
	IntVar(f *IntFlag, longName, shortName string, defValue int, usage string)                           // 绑定整数标志到指定变量
	BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string)                        // 绑定布尔标志到指定变量
	FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string)                   // 绑定浮点数标志到指定变量
	DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string)       // 绑定时间间隔类型标志到指定变量
	EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string) // 绑定枚举标志到指定变量
	SliceVar(f *SliceFlag, longName, shortName string, defValue []string, usage string)                  // 绑定字符串切片标志到指定变量
}

// 在包初始化时创建全局默认Cmd实例
func init() {
	// 处理可能的空os.Args情况
	if len(os.Args) == 0 {
		// 如果os.Args为空,则创建一个新的Cmd对象,命令行参数为"myapp",短名字为"a",错误处理方式为ExitOnError
		QCommandLine = NewCmd("myapp", "", flag.ExitOnError)
	} else {
		// 如果os.Args不为空,则创建一个新的Cmd对象,命令行参数为filepath.Base(os.Args[0]),错误处理方式为ExitOnError
		QCommandLine = NewCmd(filepath.Base(os.Args[0]), "", flag.ExitOnError)
	}
}

// SetVersion 为全局默认命令设置版本信息
// 参数说明：
//   - version: 版本信息字符串，用于标识命令的版本。
func SetVersion(version string) {
	QCommandLine.SetVersion(version)
}

// GetVersion 获取全局默认命令的版本信息
// 返回值：
//   - string: 版本信息字符串。
func GetVersion() string {
	return QCommandLine.GetVersion()
}

// String 为全局默认命令创建一个字符串类型的命令行标志。
// 该函数会调用全局默认命令实例的 String 方法，为命令行添加一个支持长短标志的字符串参数。
// 参数说明：
//   - name: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *StringFlag: 指向新创建的字符串标志对象的指针。
func String(longName, shortName, defValue, usage string) *StringFlag {
	return QCommandLine.String(longName, shortName, defValue, usage)
}

// Int 为全局默认命令创建一个整数类型的命令行标志。
// 该函数会调用全局默认命令实例的 Int 方法，为命令行添加一个支持长短标志的整数参数。
// 参数说明：
//   - name: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *IntFlag: 指向新创建的整数标志对象的指针。
func Int(longName, shortName string, defValue int, usage string) *IntFlag {
	return QCommandLine.Int(longName, shortName, defValue, usage)
}

// Bool 为全局默认命令创建一个布尔类型的命令行标志。
// 该函数会调用全局默认命令实例的 Bool 方法，为命令行添加一个支持长短标志的布尔参数。
// 参数说明：
//   - name: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *BoolFlag: 指向新创建的布尔标志对象的指针。
func Bool(longName, shortName string, defValue bool, usage string) *BoolFlag {
	return QCommandLine.Bool(longName, shortName, defValue, usage)
}

// Float 为全局默认命令创建一个浮点数类型的命令行标志。
// 该函数会调用全局默认命令实例的 Float 方法，为命令行添加一个支持长短标志的浮点数参数。
// 参数说明：
//   - name: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *FloatFlag: 指向新创建的浮点数标志对象的指针。
func Float(longName, shortName string, defValue float64, usage string) *FloatFlag {
	return QCommandLine.Float(longName, shortName, defValue, usage)
}

// StringVar 函数的作用是将一个字符串类型的命令行标志绑定到全局默认命令的 `StringFlag` 指针上。
// 借助全局默认命令实例 `QCommandLine` 的 `StringVar` 方法，为命令行添加支持长短标志的字符串参数，
// 并将该参数与传入的 `StringFlag` 指针关联，以便后续获取和使用该标志的值。
// 参数说明：
// - f: 指向 `StringFlag` 的指针，用于存储和管理该字符串类型命令行标志的相关信息，包括当前值、默认值等。
// - name: 命令行标志的长名称，在命令行中需以 `--name` 的格式使用。
// - shortName: 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用。
// - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
// - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途。
func StringVar(f *StringFlag, longName, shortName, defValue, usage string) {
	QCommandLine.StringVar(f, longName, shortName, defValue, usage)
}

// IntVar 函数的作用是将整数类型的命令行标志绑定到全局默认命令的 `IntFlag` 指针上。
// 它借助全局默认命令实例 `QCommandLine` 的 `IntVar` 方法，为命令行添加支持长短标志的整数参数，
// 并将该参数与传入的 `IntFlag` 指针建立关联，方便后续对该标志的值进行获取和使用。
//
// 参数说明：
//   - f: 指向 `IntFlag` 类型的指针，此指针用于存储和管理整数类型命令行标志的各类信息，
//     如当前标志的值、默认值等。
//   - name: 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
func IntVar(f *IntFlag, longName, shortName string, defValue int, usage string) {
	QCommandLine.IntVar(f, longName, shortName, defValue, usage)
}

// BoolVar 函数的作用是将布尔类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。
// 它会调用全局默认命令实例的 `BoolVar` 方法，为命令行添加一个支持长短和短标志的布尔参数，
// 并将该参数与传入的 `BoolFlag` 指针建立关联，后续可以通过该指针获取和使用该标志的值。
// 参数说明：
// - f: 指向 `BoolFlag` 类型的指针，用于存储和管理布尔类型命令行标志的相关信息，如当前值、默认值等。
// - name: 标志的长名称，在命令行中以 `--name` 的形式使用。
// - shortName: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
// - defValue: 标志的默认值，当命令行未指定该标志时，会使用此默认值。
// - usage: 标志的帮助说明信息，用于在显示帮助信息时展示给用户，解释该标志的用途。
func BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string) {
	QCommandLine.BoolVar(f, longName, shortName, defValue, usage)
}

// FloatVar 为全局默认命令绑定一个浮点数类型的命令行标志到指定的 `FloatFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `FloatVar` 方法，为命令行添加支持长短标志的浮点数参数，
// 并将该参数与传入的 `FloatFlag` 指针关联，以便后续获取和使用该标志的值。
// 参数说明：
// - f: 指向 `FloatFlag` 的指针，用于存储和管理该浮点数类型命令行标志的相关信息，包括当前值、默认值等。
// - name: 命令行标志的长名称，在命令行中需以 `--name` 的格式使用。
// - shortName: 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用。
// - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
// - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途。
func FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string) {
	QCommandLine.FloatVar(f, longName, shortName, defValue, usage)
}

// Enum 为全局默认命令定义一个枚举类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Enum` 方法，为命令行添加支持长短标志的枚举类型参数，
// 参数说明：
// - name: 标志的长名称，在命令行中以 `--name` 的形式使用。
// - shortName: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
// - defValue: 标志的默认值，当命令行未指定该标志时使用。
// - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
// - enumValues: 枚举值的集合，用于指定标志可接受的取值范围。
//
// 返回值：
// - *EnumFlag: 指向新创建的枚举类型标志对象的指针。
func Enum(longName, shortName string, defValue string, usage string, enumValues []string) *EnumFlag {
	return QCommandLine.Enum(longName, shortName, defValue, usage, enumValues)
}

// EnumVar 为全局默认命令将一个枚举类型的命令行标志绑定到指定的 `EnumFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `EnumVar` 方法，为命令行添加支持长短标志的枚举类型参数，
// 参数说明：
//   - f: 指向 `EnumFlag` 类型的指针，此指针用于存储和管理枚举类型命令行标志的各类信息，
//     如当前标志的值、默认值等。
//   - name: 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
//   - enumValues: 枚举值的集合，用于指定标志可接受的取值范围。
func EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string) {
	QCommandLine.EnumVar(f, longName, shortName, defValue, usage, enumValues)
}

// Duration 为全局默认命令定义一个时间间隔类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Duration` 方法，为命令行添加支持长短标志的时间间隔类型参数，
// 参数说明：
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
//
// 返回值：
//   - *DurationFlag: 指向新创建的时间间隔类型标志对象的指针。
func Duration(longName, shortName string, defValue time.Duration, usage string) *DurationFlag {
	return QCommandLine.Duration(longName, shortName, defValue, usage)
}

// DurationVar 为全局默认命令将一个时间间隔类型的命令行标志绑定到指定的 `DurationFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `DurationVar` 方法，为命令行添加支持长短标志的时间间隔类型参数，
// 参数说明：
//   - f: 指向 `DurationFlag` 类型的指针，此指针用于存储和管理时间间隔类型命令行标志的各类信息，
//     如当前标志的值、默认值等。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
func DurationVar(f *DurationFlag, longName, shortName string, defValue time.Duration, usage string) {
	QCommandLine.DurationVar(f, longName, shortName, defValue, usage)
}

// Slice 为全局默认命令定义一个字符串切片类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Slice` 方法，为命令行添加支持长短标志的字符串切片类型参数，
// 参数说明：
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
//
// 返回值：
//   - *SliceFlag: 指向新创建的字符串切片类型标志对象的指针。
func Slice(longName, shortName string, defValue []string, usage string) *SliceFlag {
	return QCommandLine.Slice(longName, shortName, defValue, usage)
}

// SliceVar 为全局默认命令将一个字符串切片类型的命令行标志绑定到指定的 `SliceFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `SliceVar` 方法，为命令行添加支持长短标志的字符串切片类型参数，
// 参数说明：
//   - f: 指向要绑定的 `SliceFlag` 对象的指针。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
func SliceVar(f *SliceFlag, longName, shortName string, defValue []string, usage string) {
	QCommandLine.SliceVar(f, longName, shortName, defValue, usage)
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

// AddSubCmd 向全局默认命令实例 `QCommandLine` 添加一个或多个子命令。
// 该函数会调用全局默认命令实例的 `AddSubCmd` 方法，支持批量添加子命令。
// 在添加过程中，会检查子命令是否为 `nil` 以及是否存在循环引用，若有异常则返回错误信息。
// 参数:
//   - subCmds: 可变参数，接收一个或多个 `*Cmd` 类型的子命令实例。
//
// 返回值:
//   - error: 若添加子命令过程中出现错误（如子命令为 `nil` 或存在循环引用），则返回错误信息；否则返回 `nil`。
func AddSubCmd(subCmds ...*Cmd) error {
	return QCommandLine.AddSubCmd(subCmds...)
}

// Args 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数切片。
// 非标志参数是指命令行中未被识别为标志的参数。
// 返回值:
//   - []string: 包含所有非标志参数的字符串切片。
func Args() []string {
	return QCommandLine.Args()
}

// Arg 获取全局默认命令实例 `QCommandLine` 解析后的指定索引位置的非标志参数。
// 索引从 0 开始，若索引超出非标志参数切片的范围，将返回空字符串。
// 参数:
//   - i: 非标志参数的索引位置，从 0 开始计数。
//
// 返回值:
//   - string: 指定索引位置的非标志参数；若索引越界，则返回空字符串。
func Arg(i int) string {
	return QCommandLine.Arg(i)
}

// NArg 获取全局默认命令实例 `QCommandLine` 解析后的非标志参数的数量。
// 返回值:
//   - int: 非标志参数的数量。
func NArg() int {
	return QCommandLine.NArg()
}

// NFlag 获取全局默认命令实例 `QCommandLine` 解析后已定义和使用的标志的数量。
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

// FlagExists 检查全局默认命令实例 `QCommandLine` 中是否存在指定名称的标志。
// 该函数会调用全局默认命令实例的 `FlagExists` 方法，用于检查命令行中是否存在指定名称的标志。
// 参数:
//   - name: 要检查的标志名称，可以是长名称或短名称。
//
// 返回值:
//   - bool: 若存在指定名称的标志，则返回 `true`；否则返回 `false`。
func FlagExists(name string) bool {
	return QCommandLine.FlagExists(name)
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
func SubCmds() []*Cmd {
	return QCommandLine.SubCmds()
}

// GetUseChinese 获取是否使用中文
// 该函数用于获取当前命令行标志是否使用中文。
// 返回值:
//   - bool: 如果使用中文，则返回true；否则返回false。
func GetUseChinese() bool {
	return QCommandLine.GetUseChinese()
}

// SetUseChinese 设置是否使用中文
// 该函数用于设置当前命令行标志是否使用中文。
// 参数:
//   - useChinese: 如果使用中文，则传入true；否则传入false。
func SetUseChinese(useChinese bool) {
	QCommandLine.SetUseChinese(useChinese)
}

// AddNote 添加注意事项
// 该函数用于添加命令行标志的注意事项，这些注意事项将在命令行帮助信息中显示。
// 参数:
//   - note: 注意事项内容，字符串类型。
func AddNote(note string) {
	QCommandLine.AddNote(note)
}

// AddExample 添加示例
// 该函数用于添加命令行标志的示例，这些示例将在命令行帮助信息中显示。
// 参数:
//   - e: 示例信息，ExampleInfo 类型。
func AddExample(e ExampleInfo) {
	QCommandLine.AddExample(e)
}

// GetExamples 获取示例信息
// 该函数用于获取命令行标志的示例信息列表。
// 返回值:
//   - []ExampleInfo: 示例信息列表，每个元素为 ExampleInfo 类型。
func GetExamples() []ExampleInfo {
	return QCommandLine.GetExamples()
}

// GetHelp 返回全局默认命令实例 `QCommandLine` 的帮助信息。
// 返回值:
//   - string: 命令行帮助信息。
func GetHelp() string {
	return QCommandLine.GetHelp()
}

// SetHelp 配置全局默认命令实例 `QCommandLine` 的帮助信息。
// 参数:
//   - help: 新的帮助信息，字符串类型。
func SetHelp(help string) {
	QCommandLine.SetHelp(help)
}

// SetUsageSyntax 配置全局默认命令实例 `QCommandLine` 的用法信息。
// 参数:
//   - usage: 新的用法信息，字符串类型。
//
// 示例:
//
//	qflag.SetUsageSyntax("Usage: qflag [options]")
func SetUsageSyntax(usageSyntax string) {
	QCommandLine.SetUsageSyntax(usageSyntax)
}

// GetUsageSyntax 获取全局默认命令实例 `QCommandLine` 的用法信息。
// 返回值:
//   - string: 命令行用法信息。
func GetUsageSyntax() string {
	return QCommandLine.GetUsageSyntax()
}

// SetLogoText 配置全局默认命令实例 `QCommandLine` 的 logo 文本。
// 默认 logo 文本为：
//
//	qflag v1.0.0
//	https://github.com/qinguoyi/qflag
//
// 参数:
//   - logoText: 配置的 logo 文本，字符串类型。
func SetLogoText(logoText string) {
	QCommandLine.SetLogoText(logoText)
}

// GetLogoText 获取全局默认命令实例 `QCommandLine` 的 logo 文本。
// 返回值:
//   - string: 配置的 logo 文本。
func GetLogoText() string {
	return QCommandLine.GetLogoText()
}

// SetModuleHelps 配置模块帮助信息。
// 参数:
//   - moduleHelps: 模块帮助信息，字符串类型。
func SetModuleHelps(moduleHelps string) {
	QCommandLine.SetModuleHelps(moduleHelps)
}

// GetModuleHelps 获取模块帮助信息。
// 返回值:
//   - string: 模块帮助信息。
func GetModuleHelps() string {
	return QCommandLine.GetModuleHelps()
}
