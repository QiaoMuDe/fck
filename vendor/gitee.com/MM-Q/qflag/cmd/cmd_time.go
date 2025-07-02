package cmd

import (
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// Time 添加时间类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 时间标志对象指针
func (c *Cmd) Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag {
	f := &flags.TimeFlag{}
	c.TimeVar(f, longName, shortName, defValue, usage)
	return f
}

// TimeVar 绑定时间类型标志到指针并内部注册Flag对象
//
// 参数依次为: 时间标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为nil
	if f == nil {
		panic("TimeFlag pointer cannot be nil")
	}

	// 参数校验（复用公共函数）
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值
	currentTime := new(time.Time)
	*currentTime = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentTime); initErr != nil {
		panic(initErr)
	}

	// 创建FlagMeta对象
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
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
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}
