package flag

import (
	"strconv"

	"gitee.com/MM-Q/qflag/internal/types"
)

// 平台相关的整数位数
const (
	// IntSize 当前平台上int类型的位数
	// 在32位系统上为32, 在64位系统上为64
	IntSize = strconv.IntSize
	// UintSize 当前平台上uint类型的位数, 与int相同
	// 在32位系统上为32, 在64位系统上为64
	UintSize = strconv.IntSize
)

// IntFlag 整数标志
//
// IntFlag 用于处理整数类型的命令行参数。
// 使用平台相关的int类型, 在32位系统上为32位整数, 在64位系统上为64位整数。
//
// 注意事项:
//   - 支持正数和负数
//   - 支持十进制格式
//   - 超出平台int范围会返回错误
type IntFlag struct {
	*BaseFlag[int]
}

// NewIntFlag 创建整数标志
//
// 参数:
//   - longName: 长选项名, 如 "count"
//   - shortName: 短选项名, 如 "c"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *IntFlag: 整数标志实例
func NewIntFlag(longName, shortName, desc string, default_ int) *IntFlag {
	return &IntFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeInt, longName, shortName, desc, default_),
	}
}

// Set 设置整数标志的值
//
// 参数:
//   - value: 要设置的整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseInt 解析字符串
//   - 使用平台相关的位数(IntSize)
//   - 如果值超出平台int范围, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *IntFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_INT", "int value cannot be empty", nil)
	}

	n, err := strconv.ParseInt(value, 10, IntSize)
	if err != nil {
		return types.WrapParseError(err, "int", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(int(n)); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = int(n)
	f.isSet = true

	return nil
}

// Int64Flag 64位整数标志
//
// Int64Flag 用于处理64位整数类型的命令行参数。
// 在所有平台上都使用固定的64位整数, 提供一致的行为。
//
// 注意事项:
//   - 支持正数和负数
//   - 支持十进制格式
//   - 范围: -9,223,372,036,854,775,808 到 9,223,372,036,854,775,807
type Int64Flag struct {
	*BaseFlag[int64]
}

// NewInt64Flag 创建64位整数标志
//
// 参数:
//   - longName: 长选项名, 如 "timestamp"
//   - shortName: 短选项名, 如 "ts"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *Int64Flag: 64位整数标志实例
func NewInt64Flag(longName, shortName, desc string, default_ int64) *Int64Flag {
	return &Int64Flag{
		BaseFlag: NewBaseFlag(types.FlagTypeInt64, longName, shortName, desc, default_),
	}
}

// Set 设置64位整数标志的值
//
// 参数:
//   - value: 要设置的64位整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseInt 解析字符串
//   - 固定使用64位精度
//   - 如果值超出64位整数范围, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *Int64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_INT64", "int64 value cannot be empty", nil)
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return types.WrapParseError(err, "int64", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(n); err != nil {
			return err
		}
	}

	*f.value = n
	f.isSet = true

	return nil
}

// UintFlag 无符号整数标志
//
// UintFlag 用于处理无符号整数类型的命令行参数。
// 使用平台相关的uint类型, 在32位系统上为32位无符号整数, 在64位系统上为64位无符号整数。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 超出平台uint范围会返回错误
type UintFlag struct {
	*BaseFlag[uint]
}

// NewUintFlag 创建新的无符号整数标志
//
// 参数:
//   - longName: 长选项名, 如 "port"
//   - shortName: 短选项名, 如 "p"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *UintFlag: 无符号整数标志实例
func NewUintFlag(longName, shortName, desc string, default_ uint) *UintFlag {
	return &UintFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeUint, longName, shortName, desc, default_),
	}
}

// Set 设置无符号整数标志的值
//
// 参数:
//   - value: 要设置的无符号整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseUint 解析字符串
//   - 使用平台相关的位数(UintSize)
//   - 如果值超出平台uint范围或为负数, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *UintFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_UINT", "uint value cannot be empty", nil)
	}

	n, err := strconv.ParseUint(value, 10, UintSize)
	if err != nil {
		return types.WrapParseError(err, "uint", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(uint(n)); err != nil {
			return err
		}
	}

	*f.value = uint(n)
	f.isSet = true

	return nil
}

// Uint8Flag 8位无符号整数标志
//
// Uint8Flag 用于处理8位无符号整数类型的命令行参数。
// 适用于处理字节值、小范围计数器等场景。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 255
type Uint8Flag struct {
	*BaseFlag[uint8]
}

// NewUint8Flag 创建新的8位无符号整数标志
//
// 参数:
//   - longName: 长选项名, 如 "byte"
//   - shortName: 短选项名, 如 "b"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *Uint8Flag: 8位无符号整数标志实例
func NewUint8Flag(longName, shortName, desc string, default_ uint8) *Uint8Flag {
	return &Uint8Flag{
		BaseFlag: NewBaseFlag(types.FlagTypeUint8, longName, shortName, desc, default_),
	}
}

// Set 设置8位无符号整数标志的值
//
// 参数:
//   - value: 要设置的8位无符号整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseUint 解析字符串
//   - 固定使用8位精度
//   - 如果值超出0-255范围或为负数, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *Uint8Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_UINT8", "uint8 value cannot be empty", nil)
	}

	n, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return types.WrapParseError(err, "uint8", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(uint8(n)); err != nil {
			return err
		}
	}

	*f.value = uint8(n)
	f.isSet = true

	return nil
}

// Uint16Flag 16位无符号整数标志
//
// Uint16Flag 用于处理16位无符号整数类型的命令行参数。
// 适用于处理端口号、短范围计数器等场景。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 65,535
type Uint16Flag struct {
	*BaseFlag[uint16]
}

// NewUint16Flag 创建新的16位无符号整数标志
//
// 参数:
//   - longName: 长选项名, 如 "port"
//   - shortName: 短选项名, 如 "p"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *Uint16Flag: 16位无符号整数标志实例
func NewUint16Flag(longName, shortName, desc string, default_ uint16) *Uint16Flag {
	return &Uint16Flag{
		BaseFlag: NewBaseFlag(types.FlagTypeUint16, longName, shortName, desc, default_),
	}
}

// Set 设置16位无符号整数标志的值
//
// 参数:
//   - value: 要设置的16位无符号整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseUint 解析字符串
//   - 固定使用16位精度
//   - 如果值超出0-65,535范围或为负数, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *Uint16Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_UINT16", "uint16 value cannot be empty", nil)
	}

	n, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return types.WrapParseError(err, "uint16", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(uint16(n)); err != nil {
			return err
		}
	}

	*f.value = uint16(n)
	f.isSet = true

	return nil
}

// Uint32Flag 32位无符号整数标志
//
// Uint32Flag 用于处理32位无符号整数类型的命令行参数。
// 适用于处理IP地址、大范围计数器等场景。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 4,294,967,295
type Uint32Flag struct {
	*BaseFlag[uint32]
}

// NewUint32Flag 创建新的32位无符号整数标志
//
// 参数:
//   - longName: 长选项名, 如 "ip"
//   - shortName: 短选项名, 如 "i"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *Uint32Flag: 32位无符号整数标志实例
func NewUint32Flag(longName, shortName, desc string, default_ uint32) *Uint32Flag {
	return &Uint32Flag{
		BaseFlag: NewBaseFlag(types.FlagTypeUint32, longName, shortName, desc, default_),
	}
}

// Set 设置32位无符号整数标志的值
//
// 参数:
//   - value: 要设置的32位无符号整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseUint 解析字符串
//   - 固定使用32位精度
//   - 如果值超出0-4,294,967,295范围或为负数, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *Uint32Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_UINT32", "uint32 value cannot be empty", nil)
	}

	n, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return types.WrapParseError(err, "uint32", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(uint32(n)); err != nil {
			return err
		}
	}

	*f.value = uint32(n)
	f.isSet = true

	return nil
}

// Uint64Flag 64位无符号整数标志
//
// Uint64Flag 用于处理64位无符号整数类型的命令行参数。
// 在所有平台上都使用固定的64位无符号整数, 提供一致的行为。
//
// 注意事项:
//   - 只支持非负数
//   - 支持十进制格式
//   - 范围: 0 到 18,446,744,073,709,551,615
type Uint64Flag struct {
	*BaseFlag[uint64]
}

// NewUint64Flag 创建新的64位无符号整数标志
//
// 参数:
//   - longName: 长选项名, 如 "id"
//   - shortName: 短选项名, 如 "i"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *Uint64Flag: 64位无符号整数标志实例
func NewUint64Flag(longName, shortName, desc string, default_ uint64) *Uint64Flag {
	return &Uint64Flag{
		BaseFlag: NewBaseFlag(types.FlagTypeUint64, longName, shortName, desc, default_),
	}
}

// Set 设置64位无符号整数标志的值
//
// 参数:
//   - value: 要设置的64位无符号整数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseUint 解析字符串
//   - 固定使用64位精度
//   - 如果值超出0-18,446,744,073,709,551,615范围或为负数, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *Uint64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_UINT64", "uint64 value cannot be empty", nil)
	}

	n, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return types.WrapParseError(err, "uint64", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(n); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = n
	f.isSet = true

	return nil
}

// Float64Flag 64位浮点数标志
//
// Float64Flag 用于处理64位浮点数类型的命令行参数。
// 支持整数、小数和科学计数法表示的数值。
//
// 注意事项:
//   - 支持正数和负数
//   - 支持十进制格式和科学计数法
//   - 支持特殊值: NaN、+Inf、-Inf
//   - 精度遵循IEEE 754双精度浮点数标准
type Float64Flag struct {
	*BaseFlag[float64]
}

// NewFloat64Flag 创建新的64位浮点数标志
//
// 参数:
//   - longName: 长选项名, 如 "ratio"
//   - shortName: 短选项名, 如 "r"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *Float64Flag: 64位浮点数标志实例
func NewFloat64Flag(longName, shortName, desc string, default_ float64) *Float64Flag {
	return &Float64Flag{
		BaseFlag: NewBaseFlag(types.FlagTypeFloat64, longName, shortName, desc, default_),
	}
}

// Set 设置64位浮点数标志的值
//
// 参数:
//   - value: 要设置的64位浮点数字符串
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 使用 strconv.ParseFloat 解析字符串
//   - 固定使用64位精度
//   - 支持十进制格式和科学计数法
//   - 如果值格式无效, 返回解析错误
//   - 先解析，然后验证，最后设置值
func (f *Float64Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return types.NewError("INVALID_FLOAT64", "float64 value cannot be empty", nil)
	}

	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return types.WrapParseError(err, "float64", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(n); err != nil {
			return err
		}
	}

	// 设置值并标记为已设置
	*f.value = n
	f.isSet = true

	return nil
}
