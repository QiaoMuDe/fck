package cmd

import "gitee.com/MM-Q/qflag/flags"

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
