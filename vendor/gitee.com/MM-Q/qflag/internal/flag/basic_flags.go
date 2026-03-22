package flag

import (
	"strconv"

	"gitee.com/MM-Q/qflag/internal/types"
)

// StringFlag 字符串标志
//
// StringFlag 用于处理字符串类型的命令行参数。
// 它接受任何字符串值, 包括空字符串。
type StringFlag struct {
	*BaseFlag[string]
}

// NewStringFlag 创建字符串标志
//
// 参数:
//   - longName: 长选项名, 如 "output"
//   - shortName: 短选项名, 如 "o"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *StringFlag: 字符串标志实例
func NewStringFlag(longName, shortName, desc string, default_ string) *StringFlag {
	return &StringFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeString, longName, shortName, desc, default_),
	}
}

// Set 设置字符串标志的值
//
// 参数:
//   - value: 要设置的字符串值
//
// 返回值:
//   - error: 如果验证失败返回错误
//
// 注意事项:
//   - 字符串标志接受任何字符串值, 包括空字符串
//   - 空字符串不经过验证器
//   - 非空字符串先验证，验证通过后再设置值
func (f *StringFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 处理空字符串（不验证）
	if value == "" {
		*f.value = value
		f.isSet = true
		return nil
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(value); err != nil {
			return err
		}
	}

	// 设置标志值并标记为已设置
	*f.value = value
	f.isSet = true

	return nil
}

// BoolFlag 布尔标志
//
// BoolFlag 用于处理布尔类型的命令行参数。
// 它接受多种布尔值表示形式, 包括 "true", "false", "1", "0", "t", "f", "TRUE", "FALSE" 等。
type BoolFlag struct {
	*BaseFlag[bool]
}

// NewBoolFlag 创建布尔标志
//
// 参数:
//   - longName: 长选项名, 如 "verbose"
//   - shortName: 短选项名, 如 "v"
//   - desc: 标志描述
//   - default_: 默认值
//
// 返回值:
//   - *BoolFlag: 布尔标志实例
func NewBoolFlag(longName, shortName, desc string, default_ bool) *BoolFlag {
	return &BoolFlag{
		BaseFlag: NewBaseFlag(types.FlagTypeBool, longName, shortName, desc, default_),
	}
}

// Set 设置布尔标志的值
//
// 参数:
//   - value: 要设置的字符串值
//
// 返回值:
//   - error: 如果解析失败或验证失败返回错误
//
// 注意事项:
//   - 支持标准库的布尔值解析格式
//   - 空字符串会被解析为true (这是Go flag包的标准行为)
//   - 支持的值: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
//   - 空字符串直接返回true，不经过验证
//   - 其他值先解析，然后验证，最后设置值
func (f *BoolFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 处理空字符串, 这是Go flag包中布尔标志的特殊行为, 空字符串会被解析为true
	if value == "" {
		*f.value = true
		f.isSet = true
		return nil
	}

	// 解析布尔值
	b, err := strconv.ParseBool(value)
	if err != nil {
		return types.WrapParseError(err, "bool", value)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(b); err != nil {
			return err
		}
	}

	// 设置标志值并标记为已设置
	*f.value = b
	f.isSet = true

	return nil
}
