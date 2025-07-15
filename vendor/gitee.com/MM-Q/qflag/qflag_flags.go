package qflag

import (
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// String 为全局默认命令创建一个字符串类型的命令行标志。
// 该函数会调用全局默认命令实例的 String 方法，为命令行添加一个支持长短标志的字符串参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *flags.StringFlag: 指向新创建的字符串标志对象的指针。
func String(longName, shortName, defValue, usage string) *flags.StringFlag {
	return QCommandLine.String(longName, shortName, defValue, usage)
}

// Int 为全局默认命令创建一个整数类型的命令行标志。
// 该函数会调用全局默认命令实例的 Int 方法，为命令行添加一个支持长短标志的整数参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *flags.IntFlag: 指向新创建的整数标志对象的指针。
func Int(longName, shortName string, defValue int, usage string) *flags.IntFlag {
	return QCommandLine.Int(longName, shortName, defValue, usage)
}

// Bool 为全局默认命令创建一个布尔类型的命令行标志。
// 该函数会调用全局默认命令实例的 Bool 方法，为命令行添加一个支持长短标志的布尔参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *flags.BoolFlag: 指向新创建的布尔标志对象的指针。
func Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag {
	return QCommandLine.Bool(longName, shortName, defValue, usage)
}

// Float64 为全局默认命令创建一个浮点数类型的命令行标志。
// 该函数会调用全局默认命令实例的 Float64 方法，为命令行添加一个支持长短标志的浮点数参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --name 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *flags.FloatFlag: 指向新创建的浮点数标志对象的指针。
func Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag {
	return QCommandLine.Float64(longName, shortName, defValue, usage)
}

// StringVar 函数的作用是将一个字符串类型的命令行标志绑定到全局默认命令的 `StringFlag` 指针上。
// 借助全局默认命令实例 `QCommandLine` 的 `StringVar` 方法，为命令行添加支持长短标志的字符串参数，
// 并将该参数与传入的 `StringFlag` 指针关联，以便后续获取和使用该标志的值。
//
// 参数值：
//   - f: 指向 `StringFlag` 的指针，用于存储和管理该字符串类型命令行标志的相关信息，包括当前值、默认值等。
//   - name: 命令行标志的长名称，在命令行中需以 `--name` 的格式使用。
//   - shortName: 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途。
func StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string) {
	QCommandLine.StringVar(f, longName, shortName, defValue, usage)
}

// IntVar 函数的作用是将整数类型的命令行标志绑定到全局默认命令的 `IntFlag` 指针上。
// 它借助全局默认命令实例 `QCommandLine` 的 `IntVar` 方法，为命令行添加支持长短标志的整数参数，
// 并将该参数与传入的 `IntFlag` 指针建立关联，方便后续对该标志的值进行获取和使用。
//
// 参数值：
//   - f: 指向 `IntFlag` 类型的指针，此指针用于存储和管理整数类型命令行标志的各类信息，
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
func IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string) {
	QCommandLine.IntVar(f, longName, shortName, defValue, usage)
}

// BoolVar 函数的作用是将布尔类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。
// 它会调用全局默认命令实例的 `BoolVar` 方法，为命令行添加一个支持长短和短标志的布尔参数，
// 并将该参数与传入的 `BoolFlag` 指针建立关联，后续可以通过该指针获取和使用该标志的值。
//
// 参数值：
//   - f: 指向 `BoolFlag` 类型的指针，用于存储和管理布尔类型命令行标志的相关信息，如当前值、默认值等。
//   - longName: 标志的长名称，在命令行中以 `--name` 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时，会使用此默认值。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示给用户，解释该标志的用途。
func BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string) {
	QCommandLine.BoolVar(f, longName, shortName, defValue, usage)
}

// Float64Var 为全局默认命令绑定一个浮点数类型的命令行标志到指定的 `FloatFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Float64Var` 方法，为命令行添加支持长短标志的浮点数参数，
// 并将该参数与传入的 `FloatFlag` 指针关联，以便后续获取和使用该标志的值。
//
// 参数值：
//   - f: 指向 `FloatFlag` 的指针，用于存储和管理该浮点数类型命令行标志的相关信息，包括当前值、默认值等。
//   - longName: 命令行标志的长名称，在命令行中需以 `--name` 的格式使用。
//   - shortName: 命令行标志的短名称，在命令行中需以 `-shortName` 的格式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户，用于解释该标志的用途。
func Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string) {
	QCommandLine.Float64Var(f, longName, shortName, defValue, usage)
}

// Enum 为全局默认命令定义一个枚举类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Enum` 方法，为命令行添加支持长短标志的枚举类型参数，
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 `--name` 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 `-shortName` 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//   - enumValues: 枚举值的集合，用于指定标志可接受的取值范围。
//
// 返回值：
//   - *flags.EnumFlag: 指向新创建的枚举类型标志对象的指针。
func Enum(longName, shortName string, defValue string, usage string, enumValues []string) *flags.EnumFlag {
	return QCommandLine.Enum(longName, shortName, defValue, usage, enumValues)
}

// EnumVar 为全局默认命令将一个枚举类型的命令行标志绑定到指定的 `EnumFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `EnumVar` 方法，为命令行添加支持长短标志的枚举类型参数，
//
// 参数值：
//   - f: 指向 `EnumFlag` 类型的指针，此指针用于存储和管理枚举类型命令行标志的各类信息，
//     如当前标志的值、默认值等。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--name` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
//   - enumValues: 枚举值的集合，用于指定标志可接受的取值范围。
func EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string) {
	QCommandLine.EnumVar(f, longName, shortName, defValue, usage, enumValues)
}

// Duration 为全局默认命令定义一个时间间隔类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Duration` 方法，为命令行添加支持长短标志的时间间隔类型参数，
//
// 参数值：
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
//
// 返回值：
//   - *flags.DurationFlag: 指向新创建的时间间隔类型标志对象的指针。
func Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag {
	return QCommandLine.Duration(longName, shortName, defValue, usage)
}

// DurationVar 为全局默认命令将一个时间间隔类型的命令行标志绑定到指定的 `DurationFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `DurationVar` 方法，为命令行添加支持长短标志的时间间隔类型参数，
//
// 参数值：
//   - f: 指向 `DurationFlag` 类型的指针，此指针用于存储和管理时间间隔类型命令行标志的各类信息，
//     如当前标志的值、默认值等。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
func DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string) {
	QCommandLine.DurationVar(f, longName, shortName, defValue, usage)
}

// Slice 为全局默认命令定义一个字符串切片类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Slice` 方法，为命令行添加支持长短标志的字符串切片类型参数，
//
// 参数值：
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
//
// 返回值：
//   - *flags.SliceFlag: 指向新创建的字符串切片类型标志对象的指针。
func Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag {
	return QCommandLine.Slice(longName, shortName, defValue, usage)
}

// SliceVar 为全局默认命令将一个字符串切片类型的命令行标志绑定到指定的 `SliceFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `SliceVar` 方法，为命令行添加支持长短标志的字符串切片类型参数，
//
// 参数值：
//   - f: 指向要绑定的 `SliceFlag` 对象的指针。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 该命令行标志的默认值。当用户在命令行中未指定该标志时，会采用此默认值。该值会被复制一份，避免外部修改影响内部状态。
//   - usage: 该命令行标志的帮助说明信息，在显示帮助信息时会呈现给用户，用以解释该标志的具体用途。
func SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string) {
	QCommandLine.SliceVar(f, longName, shortName, defValue, usage)
}

// Int64 为全局默认命令定义一个64位整数类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64` 方法，为命令行添加支持长短标志的64位整数类型参数，
//
// 参数值：
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 命令行标志的默认值。
//   - usage: 命令行标志的用法说明。
//
// 返回值：
//   - *flags.Int64Flag: 指向新创建的64位整数类型标志对象的指针。
func Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag {
	return QCommandLine.Int64(longName, shortName, defValue, usage)
}

// Int64Var 函数创建一个64位整数类型标志，并将其绑定到指定的 `Int64Flag` 指针
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Int64Var` 方法，为命令行添加支持长短标志的64位整数类型参数，
//
// 参数值：
//   - f: 指向要绑定的 `Int64Flag` 对象的指针。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 命令行标志的默认值。
//   - usage: 命令行标志的用法说明。
func Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string) {
	QCommandLine.Int64Var(f, longName, shortName, defValue, usage)
}

// Uint16 为全局默认命令定义一个无符号16位整数类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Uint16` 方法，为命令行添加支持长短标志的无符号16位整数类型参数，
//
// 参数值：
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 命令行标志的默认值。
//   - usage: 命令行标志的用法说明。
//
// 返回值：
//   - *flags.Uint16Flag: 指向新创建的无符号16位整数类型标志对象的指针。
func Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag {
	return QCommandLine.Uint16(longName, shortName, defValue, usage)
}

// Uint16Var 函数创建一个无符号16位整数类型标志，并将其绑定到指定的 `Uint16Flag` 指针
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Uint16Var` 方法，为命令行添加支持长短标志的无符号16位整数类型参数，
//
// 参数值：
//   - f: 指向要绑定的 `Uint16Flag` 对象的指针。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 命令行标志的默认值。
//   - usage: 命令行标志的用法说明。
func Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string) {
	QCommandLine.Uint16Var(f, longName, shortName, defValue, usage)
}

// Time 为全局默认命令定义一个时间类型的命令行标志。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `Time` 方法，为命令行添加支持长短标志的时间类型参数，
//
// 参数值：
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 命令行标志的默认值。
//   - usage: 命令行标志的用法说明。
//
// 返回值：
//   - *flags.TimeFlag: 指向新创建的时间类型标志对象的指针。
func Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag {
	return QCommandLine.Time(longName, shortName, defValue, usage)
}

// TimeVar 为全局默认命令定义一个时间类型的命令行标志，并将其绑定到指定的 `TimeFlag` 指针。
// 该函数会调用全局默认命令实例 `QCommandLine` 的 `TimeVar` 方法，为命令行添加支持长短标志的时间类型参数，
//
// 参数值：
//   - f: 指向要绑定的 `TimeFlag` 对象的指针。
//   - longName: 命令行标志的长名称，在命令行中使用时需遵循 `--longName` 的格式。
//   - shortName: 命令行标志的短名称，在命令行中使用时需遵循 `-shortName` 的格式。
//   - defValue: 命令行标志的默认值。
//   - usage: 命令行标志的用法说明。
func TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string) {
	QCommandLine.TimeVar(f, longName, shortName, defValue, usage)
}

// Map 为全局默认命令创建一个键值对类型的命令行标志。
// 该函数会调用全局默认命令实例的 Map 方法，为命令行添加一个支持长短标志的键值对参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *flags.MapFlag: 指向新创建的键值对标志对象的指针。
func Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag {
	return QCommandLine.Map(longName, shortName, defValue, usage)
}

// MapVar 为全局默认命令将一个键值对类型的命令行标志绑定到指定的 MapFlag 指针。
// 该函数会调用全局默认命令实例的 MapVar 方法，为命令行添加支持长短标志的键值对参数，
// 并将该参数与传入的 MapFlag 指针关联，以便后续获取和使用该标志的值。
//
// 参数值：
//   - f: 指向 MapFlag 的指针，用于存储和管理该键值对类型命令行标志的相关信息。
//   - longName: 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
//   - shortName: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
func MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string) {
	QCommandLine.MapVar(f, longName, shortName, defValue, usage)
}

// Path 为全局默认命令创建一个路径类型的命令行标志。
// 该函数会调用全局默认命令实例的 Path 方法，为命令行添加一个支持长短标志的路径参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 标志的默认值，当命令行未指定该标志时使用。
//   - usage: 标志的帮助说明信息，用于在显示帮助信息时展示。
//
// 返回值：
//   - *flags.PathFlag: 指向新创建的路径标志对象的指针。
func Path(longName, shortName string, defValue string, usage string) *flags.PathFlag {
	return QCommandLine.Path(longName, shortName, defValue, usage)
}

// PathVar 为全局默认命令将一个路径类型的命令行标志绑定到指定的 PathFlag 指针。
// 该函数会调用全局默认命令实例的 PathVar 方法，为命令行添加支持长短标志的路径参数，
// 并将该参数与传入的 PathFlag 指针关联，以便后续获取和使用该标志的值。
//
// 参数值：
//   - f: 指向 PathFlag 的指针，用于存储和管理该路径类型命令行标志的相关信息。
//   - longName: 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
//   - shortName: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
func PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string) {
	QCommandLine.PathVar(f, longName, shortName, defValue, usage)
}

// Uint32 为全局默认命令创建一个无符号32位整数类型的命令行标志。
// 该函数会调用全局默认命令实例的 Uint32 方法，为命令行添加一个支持长短标志的无符号32位整数类型参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
//
// 返回值：
//   - *flags.Uint32Flag: 指向新创建的无符号32位整数标志对象的指针。
func Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag {
	return QCommandLine.Uint32(longName, shortName, defValue, usage)
}

// Uint32Var 创建并绑定一个无符号32位整数标志。
//
// 参数值：
//   - f: 指向要绑定的标志对象的指针。
//   - longName: 标志的完整名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
func Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string) {
	QCommandLine.Uint32Var(f, longName, shortName, defValue, usage)
}

// Uint64 为全局默认命令创建一个无符号64位整数类型的命令行标志。
// 该函数会调用全局默认命令实例的 Uint64 方法，为命令行添加一个支持长短标志的无符号64位整数类型参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -s 的形式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
//
// 返回值：
//   - *flags.Uint64Flag: 指向新创建的无符号64位整数标志对象的指针。
func Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag {
	return QCommandLine.Uint64(longName, shortName, defValue, usage)
}

// Uint64Var 为全局默认命令将一个无符号64位整数类型的命令行标志绑定到指定的 Uint64Flag 指针。
// 该函数会调用全局默认命令实例的 Uint64Var 方法，为命令行添加支持长短标志的无符号64位整数类型参数，
// 并将参数值绑定到指定的 Uint64Flag 指针变量中。
//
// 参数值：
//   - f: 指向 Uint64Flag 的指针，用于存储和管理该无符号64位整数类型命令行标志的相关信息。
//   - longName: 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
//   - shortName: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
func Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string) {
	QCommandLine.Uint64Var(f, longName, shortName, defValue, usage)
}

// IP4 为全局默认命令创建一个 IPv4 地址类型的命令行标志。
// 该函数会调用全局默认命令实例的 IP4 方法，为命令行添加一个支持长短标志的 IPv4 地址类型参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
//
// 返回值：
//   - *flags.IP4Flag: 指向新创建的 IPv4 地址标志对象的指针。
func IP4(longName, shortName string, defValue string, usage string) *flags.IP4Flag {
	return QCommandLine.IP4(longName, shortName, defValue, usage)
}

// IP4Var 为全局默认命令将一个 IPv4 地址类型的命令行标志绑定到指定的 IP4Flag 指针。
// 该函数会调用全局默认命令实例的 IP4Var 方法，为命令行添加支持长短标志的 IPv4 地址类型参数，
// 并将参数值绑定到指定的 IP4Flag 指针变量中。
//
// 参数值：
//   - f: 指向 IP4Flag 的指针，用于存储和管理该 IPv4 地址类型命令行标志的相关信息。
//   - longName: 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
//   - shortName: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
func IP4Var(f *flags.IP4Flag, longName, shortName string, defValue string, usage string) {
	QCommandLine.IP4Var(f, longName, shortName, defValue, usage)
}

// IP6 为全局默认命令创建一个 IPv6 地址类型的命令行标志。
// 该函数会调用全局默认命令实例的 IP6 方法，为命令行添加一个支持长短标志的 IPv6 地址类型参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
//
// 返回值：
// - *flags.IP6Flag: 指向新创建的 IPv6 地址标志对象的指针。
func IP6(longName, shortName string, defValue string, usage string) *flags.IP6Flag {
	return QCommandLine.IP6(longName, shortName, defValue, usage)
}

// IP6Var 为全局默认命令将一个 IPv6 地址类型的命令行标志绑定到指定的 IP6Flag 指针。
// 该函数会调用全局默认命令实例的 IP6Var 方法，为命令行添加支持长短标志的 IPv6 地址类型参数，
// 并将参数值绑定到指定的 IP6Flag 指针变量中。
//
// 参数值：
//   - f: 指向 IP6Flag 的指针，用于存储和管理该 IPv6 地址类型命令行标志的相关信息。
//   - longName: 命令行标志的长名称，在命令行中需以 --longName 的格式使用。
//   - shortName: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值。
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
func IP6Var(f *flags.IP6Flag, longName, shortName string, defValue string, usage string) {
	QCommandLine.IP6Var(f, longName, shortName, defValue, usage)
}

// URL 为全局默认命令创建一个 URL 地址类型的命令行标志。
// 该函数会调用全局默认命令实例的 URL 方法，为命令行添加一个支持长短标志的 URL 地址类型参数。
//
// 参数值：
//   - longName: 标志的长名称，在命令行中以 --longName 的形式使用。
//   - shortName: 标志的短名称，在命令行中以 -shortName 的形式使用。
//   - defValue: 该命令行标志的默认值，当用户在命令行中未指定该标志时，会使用此默认值
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户。
//
// 返回值：
//   - *flags.URLFlag: 指向新创建的 URL 地址标志对象的指针。
func URL(longName, shortName string, defValue string, usage string) *flags.URLFlag {
	return QCommandLine.URL(longName, shortName, defValue, usage)
}

// URLVar 为全局默认命令将一个 URL 地址类型的命令行标志绑定到指定的 URLFlag 指针。
// 该函数会调用全局默认命令实例的 URLVar 方法，为命令行添加支持长短标志的 URL 地址类型参数，
// 并将参数值绑定到指定的 URLFlag 指针变量中。
//
// 参数值：
//   - f: 指向 URLFlag 的指针，用于存储和管理该 URL 地址类型命令行标志的相关信息。
//   - longName: 命令行标志的长名称，在命令行中需以 --longName 的格式使用
//   - shortName: 命令行标志的短名称，在命令行中需以 -shortName 的格式使用
//   - defValue: 命令行标志的默认值
//   - usage: 该命令行标志的帮助说明信息，会在显示帮助信息时展示给用户
func URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string) {
	QCommandLine.URLVar(f, longName, shortName, defValue, usage)
}
