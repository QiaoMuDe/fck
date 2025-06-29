package cmd

import (
	"gitee.com/MM-Q/qflag/flags"
)

// String 添加字符串类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 字符串标志对象指针
func (c *Cmd) String(longName, shortName, defValue, usage string) *flags.StringFlag {
	f := &flags.StringFlag{}
	c.StringVar(f, longName, shortName, defValue, usage)
	return f
}

// StringVar 绑定字符串类型标志到指针并内部注册Flag对象
//
// 参数依次为: 字符串标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string) {
	// 检查指针是否为nil
	if f == nil {
		panic("StringFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 显式初始化当前值的默认值
	currentStr := new(string)
	*currentStr = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentStr); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.StringVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.StringVar(f.GetPointer(), longName, defValue, usage)
	}

	// 注册Flag对象
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}
