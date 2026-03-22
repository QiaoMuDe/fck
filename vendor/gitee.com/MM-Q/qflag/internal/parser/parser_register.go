package parser

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// registerFlag 注册标志到FlagSet
//
// 参数:
//   - f: 要注册的标志
//
// 注意事项:
//   - 分别注册长名称和短名称
//   - 为每个名称创建独立的包装器
//   - 使用标志的描述和默认值
func (p *DefaultParser) registerFlag(f types.Flag) {
	longName := f.LongName()
	shortName := f.ShortName()
	description := f.Desc()

	// 注册长名称
	if longName != "" {
		p.flagSet.Var(newFlagValueWrapper(f), longName, description)
	}

	// 注册短名称
	if shortName != "" {
		p.flagSet.Var(newFlagValueWrapper(f), shortName, description)
	}
}

// flagValueWrapper 标志值包装器
//
// flagValueWrapper 是flag.Value接口的适配器, 用于将自定义标志类型
// 适配到Go标准库的flag包中。
//
// 作用:
//   - 实现flag.Value接口
//   - 将设置和获取操作委托给内部标志
//   - 支持非标准标志类型的注册
type flagValueWrapper struct {
	flag types.Flag // 内部标志实例
}

// newFlagValueWrapper 创建标志值包装器
//
// 参数:
//   - flag: 要包装的标志实例
//
// 返回值:
//   - *flagValueWrapper: 标志值包装器实例
//
// 功能说明:
//   - 创建新的包装器实例
//   - 封装标志实例以供标准库使用
func newFlagValueWrapper(flag types.Flag) *flagValueWrapper {
	return &flagValueWrapper{flag: flag}
}

// Set 设置标志值
//
// 参数:
//   - value: 要设置的值
//
// 返回值:
//   - error: 如果设置失败返回错误
//
// 注意事项:
//   - 委托给内部标志的Set方法
//   - 实现flag.Value接口
func (w *flagValueWrapper) Set(value string) error {
	return w.flag.Set(value)
}

// String 获取标志的字符串表示
//
// 返回值:
//   - string: 标志的字符串表示
//
// 注意事项:
//   - 委托给内部标志的String方法
//   - 实现flag.Value接口
func (w *flagValueWrapper) String() string {
	return w.flag.String()
}

// IsBoolFlag 检查是否是布尔标志
//
// 返回值:
//   - bool: 如果是布尔标志返回true, 否则返回false
//
// 注意事项:
//   - 委托给内部标志的Type方法
//   - 实现flag.BoolFlag接口
//   - 布尔标志在命令行中可以不指定值, 默认为true
func (w *flagValueWrapper) IsBoolFlag() bool {
	return w.flag.Type() == types.FlagTypeBool
}
