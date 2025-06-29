package cmd

import "gitee.com/MM-Q/qflag/flags"

// Int64Var 绑定64位整数类型标志到指针并内部注册Flag对象
//
// 参数依次为: 64位整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("Int64Flag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentInt64 := new(int64)
	*currentInt64 = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentInt64); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.Int64Var(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Int64Var(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Int64 添加64位整数类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 64位整数标志对象指针
func (c *Cmd) Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag {
	f := &flags.Int64Flag{}
	c.Int64Var(f, longName, shortName, defValue, usage)
	return f
}
