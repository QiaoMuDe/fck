package cmd

import "gitee.com/MM-Q/qflag/flags"

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
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("SliceFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
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
		c.fs.Var(f, shortName, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}
