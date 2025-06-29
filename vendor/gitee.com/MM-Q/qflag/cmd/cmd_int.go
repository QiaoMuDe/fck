package cmd

import "gitee.com/MM-Q/qflag/flags"

// IntVar 绑定整数类型标志到指针并内部注册Flag对象
//
// 参数依次为: 整数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("IntFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentInt := new(int)
	*currentInt = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentInt); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.IntVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.IntVar(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Int 添加整数类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: 整数标志对象指针
func (c *Cmd) Int(longName, shortName string, defValue int, usage string) *flags.IntFlag {
	f := &flags.IntFlag{}
	c.IntVar(f, longName, shortName, defValue, usage)
	return f
}
