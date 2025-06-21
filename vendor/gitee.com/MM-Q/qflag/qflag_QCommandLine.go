package qflag

import (
	"flag"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// QCommandLine 全局默认Cmd实例
var QCommandLine *Cmd

// parseOnce 保证Parse函数只会执行一次
var parseOnce sync.Once

// QCommandLineInterface 定义了全局默认命令行接口，提供统一的命令行参数管理功能
// 该接口封装了命令行程序的常用操作，包括标志添加、参数解析和帮助信息展示
type QCommandLineInterface interface {
	String(longName, shortName, defValue, usage string) *StringFlag                                      // 添加字符串类型标志
	Int(longName, shortName string, defValue int, usage string) *IntFlag                                 // 添加整数类型标志
	Bool(longName, shortName string, defValue bool, usage string) *BoolFlag                              // 添加布尔类型标志
	Float(longName, shortName string, defValue float64, usage string) *FloatFlag                         // 添加浮点数类型标志
	Enum(longName, shortName string, defValue string, usage string, enumValues []string) *EnumFlag       // 添加枚举类型标志
	StringVar(f *StringFlag, longName, shortName, defValue, usage string)                                // 绑定字符串标志到指定变量
	IntVar(f *IntFlag, longName, shortName string, defValue int, usage string)                           // 绑定整数标志到指定变量
	BoolVar(f *BoolFlag, longName, shortName string, defValue bool, usage string)                        // 绑定布尔标志到指定变量
	FloatVar(f *FloatFlag, longName, shortName string, defValue float64, usage string)                   // 绑定浮点数标志到指定变量
	EnumVar(f *EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string) // 绑定枚举标志到指定变量
	Parse() error                                                                                        // 解析命令行参数，自动处理标志和子命令
	AddSubCmd(subCmd *Cmd)                                                                               // 添加子命令，子命令会继承父命令的上下文
	Args() []string                                                                                      // 获取所有非标志参数(未绑定到任何标志的参数)
	Arg(i int) string                                                                                    // 获取指定索引的非标志参数，索引越界返回空字符串
	NArg() int                                                                                           // 获取非标志参数的数量
	NFlag() int                                                                                          // 获取已解析的标志数量
	PrintUsage()                                                                                         // 打印命令使用说明到标准输出
	FlagExists(name string) bool                                                                         // 检查指定名称的标志是否存在(支持长/短名称)

	AddNote(note string)           // 添加一个注意事项
	GetUseChinese() bool           // 获取是否使用中文帮助信息
	SetUseChinese(useChinese bool) // 设置是否使用中文帮助信息
	AddExample(e ExampleInfo)      // 添加一个示例信息
	GetExamples() []ExampleInfo    // 获取示例信息列表

}

// 在包初始化时创建全局默认Cmd实例
func init() {
	// 处理可能的空os.Args情况
	if len(os.Args) == 0 {
		// 如果os.Args为空,则创建一个新的Cmd对象,命令行参数为"app",短名字为"a",错误处理方式为ExitOnError
		QCommandLine = NewCmd("app", "a", flag.ExitOnError)
	} else {
		// 如果os.Args不为空,则创建一个新的Cmd对象,命令行参数为filepath.Base(os.Args[0]),短名字为第一个字符,错误处理方式为ExitOnError
		longName := filepath.Base(os.Args[0])
		shortName := string(longName[0]) // 获取第一个字符作为短名称
		QCommandLine = NewCmd(longName, shortName, flag.ExitOnError)
	}
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

// Parse 解析命令行参数, 自动检查长短标志互斥, 并处理帮助标志 (全局默认命令)
// 该方法会自动处理以下情况:
// 1. 长短标志互斥检查
// 2. -h/--help 帮助标志处理
// 3. -sip/--show-install-path 安装路径标志处理
// 4. 子命令自动检测和参数传递(当第一个非标志参数匹配子命令名称时)
// 注意: 该方法保证每个Cmd实例只会解析一次
func Parse() error {
	var err error
	parseOnce.Do(func() {
		// 解析命令行参数
		err = QCommandLine.Parse(os.Args[1:])
	})
	return err
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

// PrintUsage 输出全局默认命令实例 `QCommandLine` 的使用说明信息。
// 使用说明信息通常包含命令的名称、可用的标志及其描述等内容。
func PrintUsage() {
	QCommandLine.PrintUsage()
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
