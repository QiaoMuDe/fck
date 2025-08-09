// Package cmd 基础标志创建和管理功能
// 本文件提供了Cmd结构体的基础标志创建方法，包括字符串、整数、布尔、浮点数等基本类型标志的创建和绑定功能。
package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/flags"
)

// =============================================================================
// 布尔类型标志
// =============================================================================

// BoolVar 绑定布尔类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.BoolFlag - 布尔标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: bool - 默认值
//   - usage: string - 帮助说明
func (c *Cmd) BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("BoolFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 显式初始化
	currentBool := new(bool) // 创建当前值指针
	*currentBool = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentBool); initErr != nil {
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

// Bool 添加布尔类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: bool - 默认值
//   - usage: string - 帮助说明
//
// 返回值:
//   - *flags.BoolFlag - 布尔标志对象指针
func (c *Cmd) Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag {
	f := &flags.BoolFlag{}
	c.BoolVar(f, longName, shortName, defValue, usage)
	return f
}

// =============================================================================
// 字符串类型标志
// =============================================================================

// String 添加字符串类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.StringFlag: 字符串标志对象指针
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag {
	f := &flags.StringFlag{}
	c.StringVar(f, longName, shortName, defValue, usage)
	return f
}

// StringVar 绑定字符串类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 字符串标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("StringFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 显式初始化当前值的默认值
	currentStr := new(string)
	*currentStr = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentStr); initErr != nil {
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

// =============================================================================
// 64位浮点数类型标志
// =============================================================================

// Float64 添加浮点型标志, 返回标志对象指针
//
// 参数值:
//   - longName - 长标志名
//   - shortName - 短标志
//   - defValue - 默认值
//   - usage - 帮助说明
//
// 返回值:
//   - *flags.Float64Flag - 浮点型标志对象指针
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag {
	f := &flags.Float64Flag{}
	c.Float64Var(f, longName, shortName, defValue, usage)
	return f
}

// Float64Var 绑定浮点型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.Float64Flag - 浮点型标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: float64 - 默认值
//   - usage: string - 帮助说明
func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("FloatFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 显式初始化默认值
	currentFloat := new(float64) // 显式堆分配
	*currentFloat = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentFloat); initErr != nil {
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
// 整数类型标志
// =============================================================================

// IntVar 绑定整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("IntFlag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 初始化默认值
	currentInt := new(int)
	*currentInt = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentInt); initErr != nil {
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

// Int 添加整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.IntFlag: 整数标志对象指针
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag {
	f := &flags.IntFlag{}
	c.IntVar(f, longName, shortName, defValue, usage)
	return f
}

// Int64Var 绑定64位整数类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 64位整数标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string) {
	c.ctx.Mutex.Lock()
	defer c.ctx.Mutex.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("Int64Flag pointer cannot be nil")
	}

	// 检查标志是否为内置标志
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(longName); ok {
		panic(fmt.Errorf("flag long name %s is reserved", longName))
	}
	if ok := c.ctx.BuiltinFlags.IsBuiltinFlag(shortName); ok {
		panic(fmt.Errorf("flag short name %s is reserved", shortName))
	}

	// 初始化默认值
	currentInt64 := new(int64)
	*currentInt64 = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentInt64); initErr != nil {
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

// Int64 添加64位整数类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.Int64Flag: 64位整数标志对象指针
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag {
	f := &flags.Int64Flag{}
	c.Int64Var(f, longName, shortName, defValue, usage)
	return f
}
