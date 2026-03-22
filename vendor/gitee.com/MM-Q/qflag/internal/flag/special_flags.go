package flag

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// EnumFlag 枚举标志
//
// EnumFlag 用于处理枚举类型的命令行参数, 限制输入值必须在预定义的允许值列表中。
// 使用映射表(map)实现O(1)时间复杂度的值查找, 提高性能。
//
// 特性:
//   - 使用映射表进行快速值验证
//   - 不允许空字符串作为枚举值
//   - 默认值必须在允许值列表中
//   - 不允许设置空值
type EnumFlag struct {
	*BaseFlag[string]
	// 用于快速查找的映射表
	allowedMap map[string]bool
}

// NewEnumFlag 创建枚举标志
//
// 参数:
//   - longName: 长选项名, 如 "mode"
//   - shortName: 短选项名, 如 "m"
//   - desc: 标志描述
//   - default_: 默认值, 必须在允许值列表中
//   - allowedValues: 允许的枚举值列表, 不能为空且不能包含空字符串
//
// 返回值:
//   - *EnumFlag: 枚举标志实例
//
// 注意事项:
//   - 允许值列表不能为空
//   - 允许值列表中不能包含空字符串
//   - 默认值必须在允许值列表中
//   - 如果验证失败, 会panic
func NewEnumFlag(longName, shortName, desc, default_ string, allowedValues []string) *EnumFlag {
	// 检查允许值列表不能为空
	if len(allowedValues) == 0 {
		panic(types.NewError("EMPTY_ENUM_VALUES", "allowed values cannot be empty for enum flag", nil))
	}

	// 创建映射表用于快速查找
	allowedMap := make(map[string]bool, len(allowedValues))
	for _, value := range allowedValues {
		// 不允许空字符串作为枚举值
		if value == "" {
			panic(types.NewError("EMPTY_ENUM_VALUE", "empty string cannot be used as enum value", nil))
		}
		allowedMap[value] = true
	}

	// 检查默认值是否在允许值中
	if !allowedMap[default_] {
		panic(types.NewError("INVALID_DEFAULT_ENUM",
			fmt.Sprintf("default value '%s' must be in allowed values", default_), nil))
	}

	return &EnumFlag{
		BaseFlag:   NewBaseFlag(types.FlagTypeEnum, longName, shortName, desc, default_),
		allowedMap: allowedMap,
	}
}

// Set 设置枚举标志的值
//
// 参数:
//   - value: 要设置的字符串值
//
// 返回值:
//   - error: 如果值不在允许列表中或为空, 返回错误
//
// 注意事项:
//   - 不允许设置空值
//   - 使用映射表进行O(1)时间复杂度的值验证
//   - 错误消息会列出所有允许的值
//   - 先进行枚举值验证，然后调用用户验证器，最后设置值
func (f *EnumFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 不允许空值
	if value == "" {
		return types.NewError("EMPTY_ENUM_VALUE", "empty value not allowed for enum flag", nil)
	}

	// 使用映射表快速检查值是否在允许的枚举值中
	if !f.allowedMap[value] {
		// 值不在允许的枚举值中, 返回错误
		return types.NewError("INVALID_ENUM_VALUE",
			fmt.Sprintf("invalid enum value: %s, allowed values are: %s", value, strings.Join(f.getAllowedValues(), ", ")),
			nil)
	}

	// 验证（如果设置了验证器）
	if f.validator != nil {
		if err := f.validator(value); err != nil {
			return err
		}
	}

	// 设置值并标记已设置
	*f.value = value
	f.isSet = true

	return nil
}

// GetAllowedValues 获取允许的枚举值
//
// 返回值:
//   - []string: 允许的枚举值列表
//
// 注意事项:
//   - 返回的切片顺序可能不一致, 因为基于map的key生成
//   - 此方法是线程安全的
func (f *EnumFlag) GetAllowedValues() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.getAllowedValues()
}

// getAllowedValues 获取允许的枚举值 (内部方法, 不加锁)
//
// 返回值:
//   - []string: 允许的枚举值列表
//
// 注意事项:
//   - 这是内部方法, 调用者需要自己处理线程安全
//   - 返回的切片顺序可能不一致, 因为基于map的key生成
func (f *EnumFlag) getAllowedValues() []string {
	// 从map的key生成切片
	result := make([]string, 0, len(f.allowedMap))
	for value := range f.allowedMap {
		result = append(result, value)
	}
	return result
}

// IsAllowed 检查值是否在允许的枚举值中
//
// 参数:
//   - value: 要检查的值
//
// 返回值:
//   - bool: 如果值在允许列表中返回true, 否则返回false
//
// 注意事项:
//   - 此方法是线程安全的
//   - 使用映射表进行O(1)时间复杂度的查找
func (f *EnumFlag) IsAllowed(value string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.allowedMap[value]
}

// EnumValues 获取枚举类型的可选值
//
// 返回值:
//   - []string: 枚举类型的可选值列表
//
// 功能说明:
//   - 实现 Flag 接口的 EnumValues 方法
//   - 返回所有允许的枚举值
//   - 此方法是线程安全的
func (f *EnumFlag) EnumValues() []string {
	return f.GetAllowedValues()
}
