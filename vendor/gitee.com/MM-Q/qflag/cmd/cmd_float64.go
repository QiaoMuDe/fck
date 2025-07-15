package cmd

import "gitee.com/MM-Q/qflag/flags"

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
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

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
	if initErr := f.Init(longName, shortName, usage, currentFloat); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}
