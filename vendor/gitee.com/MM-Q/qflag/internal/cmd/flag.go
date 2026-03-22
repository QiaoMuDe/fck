package cmd

import (
	"time"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/utils"
)

// Int 创建整数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.IntFlag: 新创建的整数标志
func (c *Cmd) Int(longName, shortName, description string, default_ int) *flag.IntFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewIntFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// String 创建字符串标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.StringFlag: 新创建的字符串标志
func (c *Cmd) String(longName, shortName, description string, default_ string) *flag.StringFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewStringFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Bool 创建布尔标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.BoolFlag: 新创建的布尔标志
func (c *Cmd) Bool(longName, shortName, description string, default_ bool) *flag.BoolFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewBoolFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Int64 创建64位整数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.Int64Flag: 新创建的64位整数标志
func (c *Cmd) Int64(longName, shortName, description string, default_ int64) *flag.Int64Flag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewInt64Flag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Uint 创建无符号整数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.UintFlag: 新创建的无符号整数标志
func (c *Cmd) Uint(longName, shortName, description string, default_ uint) *flag.UintFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewUintFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Uint8 创建8位无符号整数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.Uint8Flag: 新创建的8位无符号整数标志
func (c *Cmd) Uint8(longName, shortName, description string, default_ uint8) *flag.Uint8Flag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewUint8Flag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Uint16 创建16位无符号整数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.Uint16Flag: 新创建的16位无符号整数标志
func (c *Cmd) Uint16(longName, shortName, description string, default_ uint16) *flag.Uint16Flag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewUint16Flag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Uint32 创建32位无符号整数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.Uint32Flag: 新创建的32位无符号整数标志
func (c *Cmd) Uint32(longName, shortName, description string, default_ uint32) *flag.Uint32Flag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewUint32Flag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Uint64 创建64位无符号整数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.Uint64Flag: 新创建的64位无符号整数标志
func (c *Cmd) Uint64(longName, shortName, description string, default_ uint64) *flag.Uint64Flag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewUint64Flag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Float64 创建64位浮点数标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.Float64Flag: 新创建的64位浮点数标志
func (c *Cmd) Float64(longName, shortName, description string, default_ float64) *flag.Float64Flag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewFloat64Flag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Enum 创建枚举标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//   - allowedValues: 允许的枚举值列表
//
// 返回值:
//   - *flag.EnumFlag: 新创建的枚举标志
func (c *Cmd) Enum(longName, shortName, description, default_ string, allowedValues []string) *flag.EnumFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewEnumFlag(longName, shortName, description, default_, allowedValues)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Duration 创建持续时间标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.DurationFlag: 新创建的持续时间标志
func (c *Cmd) Duration(longName, shortName, description string, default_ time.Duration) *flag.DurationFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewDurationFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Time 创建时间标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.TimeFlag: 新创建的时间标志
func (c *Cmd) Time(longName, shortName, description string, default_ time.Time) *flag.TimeFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewTimeFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Size 创建大小标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.SizeFlag: 新创建的大小标志
func (c *Cmd) Size(longName, shortName, description string, default_ int64) *flag.SizeFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewSizeFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// StringSlice 创建字符串切片标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.StringSliceFlag: 新创建的字符串切片标志
func (c *Cmd) StringSlice(longName, shortName, description string, default_ []string) *flag.StringSliceFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewStringSliceFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// IntSlice 创建整数切片标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.IntSliceFlag: 新创建的整数切片标志
func (c *Cmd) IntSlice(longName, shortName, description string, default_ []int) *flag.IntSliceFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewIntSliceFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Int64Slice 创建64位整数切片标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.Int64SliceFlag: 新创建的64位整数切片标志
func (c *Cmd) Int64Slice(longName, shortName, description string, default_ []int64) *flag.Int64SliceFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewInt64SliceFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}

// Map 创建映射标志
//
// 参数:
//   - longName: 长标志名 (如 --long-name)
//   - shortName: 短标志名 (如 -s)
//   - description: 标志的描述信息
//   - default_: 标志的默认值
//
// 返回值:
//   - *flag.MapFlag: 新创建的映射标志
func (c *Cmd) Map(longName, shortName, description string, default_ map[string]string) *flag.MapFlag {
	if err := utils.ValidateFlagName(c, longName, shortName); err != nil {
		panic(err)
	}

	f := flag.NewMapFlag(longName, shortName, description, default_)
	if err := c.flagRegistry.Register(f); err != nil {
		panic(err)
	}
	return f
}
