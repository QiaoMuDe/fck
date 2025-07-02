package cmd

import (
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象
//
// 参数依次为: 时间间隔标志指针、长标志名、短标志、默认值、帮助说明
func (c *Cmd) DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string) {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	// 检查指针是否为空
	if f == nil {
		panic("DurationFlag pointer cannot be nil")
	}

	// 参数校验
	if validateErr := c.validateFlag(longName, shortName); validateErr != nil {
		panic(validateErr)
	}

	// 初始化默认值(值类型)
	currentDuration := new(time.Duration)
	*currentDuration = defValue

	// 初始化Flag对象
	if initErr := f.Init(longName, shortName, defValue, usage, currentDuration); initErr != nil {
		panic(initErr)
	}

	// 绑定短标志
	if shortName != "" {
		c.fs.DurationVar(f.GetPointer(), shortName, defValue, usage)
	}

	// 绑定长标志
	if longName != "" {
		c.fs.DurationVar(f.GetPointer(), longName, defValue, usage)
	}

	// 创建并注册标志元数据
	meta := &flags.FlagMeta{
		Flag: f, // 添加标志对象 - Flag对象
	}

	// 注册标志元数据
	if registerErr := c.flagRegistry.RegisterFlag(meta); registerErr != nil {
		panic(registerErr)
	}
}

// Duration 添加时间间隔类型标志, 返回标志对象指针
//
// 参数依次为: 长标志名、短标志、默认值、帮助说明
//
// 返回值: 时间间隔标志对象指针
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag {
	f := &flags.DurationFlag{}
	c.DurationVar(f, longName, shortName, defValue, usage)
	return f
}
