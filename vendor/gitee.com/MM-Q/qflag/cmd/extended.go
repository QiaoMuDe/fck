// Package cmd 扩展标志类型支持
// 本文件提供了Cmd结构体的扩展标志创建方法，包括枚举、时间间隔、切片、时间、映射等高级类型标志的创建和绑定功能。
package cmd

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// =============================================================================
// 枚举类型标志
// =============================================================================

// Enum 添加枚举类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: string - 默认值
//   - usage: string - 帮助说明
//   - options: []string - 限制该标志取值的枚举值切片
//
// 返回值:
//   - *flags.EnumFlag - 枚举标志对象指针
func (c *Cmd) Enum(longName, shortName string, defValue string, usage string, options []string) *flags.EnumFlag {
	f := &flags.EnumFlag{}
	c.EnumVar(f, longName, shortName, defValue, usage, options)
	return f
}

// EnumVar 绑定枚举类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.EnumFlag - 枚举标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: string - 默认值
//   - usage: string - 帮助说明
//   - options: []string - 限制该标志取值的枚举值切片
func (c *Cmd) EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, options []string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("EnumFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 初始化枚举值
	if options == nil {
		options = make([]string, 0)
	}

	// 调用枚举专用Init方法
	if initErr := f.Init(longName, shortName, defValue, usage, options); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// =============================================================================
// 无符号整数类型标志
// =============================================================================

// Uint16Var 绑定16位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 16位无符号整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 参数校验
	if f == nil {
		panic("Uint16Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 显式初始化
	currentUint16 := new(uint16)
	*currentUint16 = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentUint16); initErr != nil {
		panic(initErr)
	}

	// 注册到flagSet
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint16 添加16位无符号整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Uint16Flag: 16位无符号整数标志对象指针
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag {
	f := &flags.Uint16Flag{}
	c.Uint16Var(f, longName, shortName, defValue, usage)
	return f
}

// Uint32Var 绑定32位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 32位无符号整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 校验标志指针
	if f == nil {
		panic("Uint32Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 设置标志指针
	currentUint32 := new(uint32)
	*currentUint32 = defValue

	// 绑定默认值
	if initErr := f.Init(longName, shortName, usage, currentUint32); initErr != nil {
		panic(initErr)
	}

	// 绑定标志到指针
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint32 添加32位无符号整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Uint32Flag: 32位无符号整数标志对象指针
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag {
	f := &flags.Uint32Flag{}
	c.Uint32Var(f, longName, shortName, defValue, usage)
	return f
}

// Uint64Var 绑定64位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 64位无符号整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("Uint64Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 显式初始化
	currentUint64 := new(uint64)
	*currentUint64 = defValue

	// 注册flag
	if initErr := f.Init(longName, shortName, usage, currentUint64); initErr != nil {
		panic(initErr)
	}

	// 注册到flagSet
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint64 添加64位无符号整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Uint64Flag: 64位无符号整数标志对象指针
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag {
	f := &flags.Uint64Flag{}
	c.Uint64Var(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 时间类型标志
// =============================================================================

// Time 添加时间类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值(时间表达式, 如"now", "zero", "1h", "2006-01-02")
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.TimeFlag: 时间标志对象指针
//
// 支持的默认值格式:
//   - "now" 或 "" : 当前时间
//   - "zero" : 零时间 (time.Time{})
//   - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
//   - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
//   - RFC3339等标准格式
func (c *Cmd) Time(longName, shortName string, defValue string, usage string) *flags.TimeFlag {
	f := &flags.TimeFlag{}
	c.TimeVar(f, longName, shortName, defValue, usage)
	return f
}

// TimeVar 绑定时间类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 时间标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值(时间表达式, 如"now", "zero", "1h", "2006-01-02")
//   - usage: 帮助说明
//
// 支持的默认值格式:
//   - "now" 或 "" : 当前时间
//   - "zero" : 零时间 (time.Time{})
//   - "1h", "30m", "-2h" : 相对时间（基于当前时间的偏移）
//   - "2006-01-02", "2006-01-02 15:04:05" : 绝对时间格式
//   - RFC3339等标准格式
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue string, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("TimeFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage); initErr != nil {
		panic(initErr)
	}

	// 绑定短标志
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.DurationFlag - 时间间隔标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: time.Duration - 默认值
//   - usage: string - 帮助说明
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("DurationFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 初始化默认值(值类型)
	currentDuration := new(time.Duration)
	*currentDuration = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentDuration); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册标志元数据
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Duration 添加时间间隔类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: time.Duration - 默认值
//   - usage: string - 帮助说明
//
// 返回值:
//   - *flags.DurationFlag - 时间间隔标志对象指针
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag {
	f := &flags.DurationFlag{}
	c.DurationVar(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 键值对类型标志
// =============================================================================

// MapVar 绑定键值对类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 键值对标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("MapFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 如果默认值为nil，则初始化为空map
	if defValue == nil {
		defValue = map[string]string{}
	}

	// 初始化值
	currentMap := new(map[string]string)
	*currentMap = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentMap); initErr != nil {
		panic(initErr)
	}

	// 设置默认分隔符
	f.SetDelimiters(flags.FlagSplitComma, flags.FlagKVEqual)

	// 绑定短标志
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Map 添加键值对类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.MapFlag: 键值对标志对象指针
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag {
	f := &flags.MapFlag{}
	c.MapVar(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 切片类型标志
// =============================================================================

// Slice 绑定字符串切片类型标志并内部注册Flag对象
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.SliceFlag: 字符串切片标志对象指针
func (c *Cmd) Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag {
	f := &flags.SliceFlag{}
	c.SliceVar(f, longName, shortName, defValue, usage)
	return f
}

// SliceVar 绑定字符串切片类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 字符串切片标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("SliceFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 确保默认值不为空
	if defValue == nil {
		defValue = make([]string, 0)
	}

	// 初始化Flag对象字段
	if initErr := f.Init(longName, shortName, defValue, usage); initErr != nil {
		panic(initErr)
	}

	// 绑定短标志
	if shortName != "" {
		c.ctx.FlagSet.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.ctx.FlagSet.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}
