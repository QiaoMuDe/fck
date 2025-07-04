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
