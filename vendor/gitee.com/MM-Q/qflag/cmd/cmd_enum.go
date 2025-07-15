package cmd

import "gitee.com/MM-Q/qflag/flags"

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
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("EnumFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
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
