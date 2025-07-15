package cmd

import (
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// DurationVar 绑定时间间隔类型标志到指针并内部注册Flag对象
//
// 参数值:
//   - f: *flags.DurationFlag - 时间间隔标志对象指针
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: time.Duration - 默认值
//   - usage: string - 帮助说明
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
	if initErr := f.Init(longName, shortName, usage, currentDuration); initErr != nil {
		panic(initErr)
	}

	// 绑定长短标志
	if shortName != "" {
		c.fs.Var(f, shortName, usage)
	}
	if longName != "" {
		c.fs.Var(f, longName, usage)
	}

	// 注册标志元数据
	if registerErr := c.flagRegistry.RegisterFlag(&flags.FlagMeta{Flag: f}); registerErr != nil {
		panic(registerErr)
	}
}

// Duration 添加时间间隔类型标志, 返回标志对象指针
//
// 参数值:
//   - longName: string - 长标志名
//   - shortName: string - 短标志
//   - defValue: time.Duration - 默认值
//   - usage: string - 帮助说明
//
// 返回值:
//   - *flags.DurationFlag - 时间间隔标志对象指针
func (c *Cmd) Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag {
	f := &flags.DurationFlag{}
	c.DurationVar(f, longName, shortName, defValue, usage)
	return f
}
