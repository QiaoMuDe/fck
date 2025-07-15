package cmd

import "gitee.com/MM-Q/qflag/flags"

// MapVar 绑定键值对类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: 键值对标志指针
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
func (c *Cmd) MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("MapFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentMap := new(map[string]string)
	*currentMap = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, usage, currentMap); initErr != nil {
		panic(initErr)
	}

	// 设置默认分隔符
	f.SetDelimiters(flags.FlagSplitComma, flags.FlagKVEqual)

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

// Map 添加键值对类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: 长标志名
//   - shortName: 短标志名
//   - defValue: 默认值
//   - usage: 帮助说明
//
// 返回值:
//   - *flags.MapFlag: 键值对标志对象指针
func (c *Cmd) Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag {
	f := &flags.MapFlag{}
	c.MapVar(f, longName, shortName, defValue, usage)
	return f
}
