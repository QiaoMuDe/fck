package cmd

import "gitee.com/MM-Q/qflag/flags"

// Float64 添加浮点型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 浮点型标志对象指针
func (c *Cmd) Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag {
	f := &flags.Float64Flag{}
	c.Float64Var(f, longName, shortName, defValue, usage)
	return f
}

// Float64Var 绑定浮点型标志到指针并内部注册Flag对象
//
// 参数依次为: 浮点数标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string) {
	// 检查指针是否为空
	if f == nil {
		panic("FloatFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化默认值
	currentFloat := new(float64) // 显式堆分配
	*currentFloat = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentFloat); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.Float64Var(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Float64Var(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}
