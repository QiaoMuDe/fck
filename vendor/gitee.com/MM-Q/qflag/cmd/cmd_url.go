package cmd

import "gitee.com/MM-Q/qflag/flags"

// URLVar 绑定URL类型标志到指针并内部注册Flag对象
// 参数依次为: URL标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) URLVar(f *flags.URLFlag, longName, shortName string, defValue string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("URLFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化标志值
	currentURL := new(string)
	*currentURL = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentURL); initErr != nil {
		panic(initErr)
	}

	// 注册标志
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

// URL 添加URL类型标志, 返回标志对象指针
// 参数依次为: 长标志名、短标志、默认值、帮助说明
// 返回值: URL标志对象指针
func (c *Cmd) URL(longName, shortName string, defValue string, usage string) *flags.URLFlag {
	f := &flags.URLFlag{}
	c.URLVar(f, longName, shortName, defValue, usage)
	return f
}
