package cmd

import "gitee.com/MM-Q/qflag/flags"

// Uint16Var 绑定16位无符号整数类型标志到指针并内部注册Flag对象
//
// 参数依次为: 16位无符号整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 参数校验
	if f == nil {
		panic("Uint16Flag pointer cannot be nil")
	}

	// 通用校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
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
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint16 添加16位无符号整数类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 16位无符号整数标志对象指针
func (c *Cmd) Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag {
	f := &flags.Uint16Flag{}
	c.Uint16Var(f, longName, shortName, defValue, usage)
	return f
}

// Uint32Var 绑定32位无符号整数类型标志到指针并内部注册Flag对象
// 参数依次为: 32位无符号整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Uint32Var(f *flags.Uint32Flag, longName, shortName string, defValue uint32, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 校验标志指针
	if f == nil {
		panic("Uint32Flag pointer cannot be nil")
	}

	// 通用标志校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
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
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint32 添加32位无符号整数类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: 32位无符号整数标志对象指针
func (c *Cmd) Uint32(longName, shortName string, defValue uint32, usage string) *flags.Uint32Flag {
	f := &flags.Uint32Flag{}
	c.Uint32Var(f, longName, shortName, defValue, usage)
	return f
}

// Uint64Var 绑定64位无符号整数类型标志到指针并内部注册Flag对象
// 参数依次为: 64位无符号整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Uint64Var(f *flags.Uint64Flag, longName, shortName string, defValue uint64, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("Uint64Flag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
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
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册到flagRegistry
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Uint64 添加64位无符号整数类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: 64位无符号整数标志对象指针
func (c *Cmd) Uint64(longName, shortName string, defValue uint64, usage string) *flags.Uint64Flag {
	f := &flags.Uint64Flag{}
	c.Uint64Var(f, longName, shortName, defValue, usage)
	return f
}
