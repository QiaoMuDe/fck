package flag

import (
	"fmt"
	"sync"

	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/internal/utils"
)

// BaseFlag 泛型基础标志结构体
//
// BaseFlag 是所有标志类型的基础结构, 使用泛型支持多种数据类型。
// 它提供了标志的基本功能, 包括名称管理、值存储、默认值处理、
// 环境变量绑定等。
//
// 线程安全:
//   - 所有公共方法都使用读写锁保护, 确保并发安全
//   - 读操作使用读锁, 写操作使用写锁
//
// 字段说明:
//   - mu: 读写锁, 保护并发访问
//   - longName: 长选项名称, 如 "--help"
//   - shortName: 短选项名称, 如 "-h"
//   - desc: 标志描述信息
//   - flagType: 标志类型枚举值
//   - value: 指向当前值的指针
//   - default_: 默认值
//   - isSet: 标志是否已被设置
//   - envVar: 关联的环境变量名
//   - validator: 验证器函数
type BaseFlag[T any] struct {
	mu        sync.RWMutex       // 读写锁
	value     *T                 // 当前值指针
	default_  T                  // 默认值
	isSet     bool               // 标志是否已被设置
	envVar    string             // 关联的环境变量名
	validator types.Validator[T] // 验证器函数

	// 不可变属性, 无需挂锁
	longName  string         // 长选项名称
	shortName string         // 短选项名称
	desc      string         // 标志描述信息
	flagType  types.FlagType // 标志类型枚举值
}

// NewBaseFlag 创建新的基础标志实例
//
// 参数:
//   - flagType: 标志类型枚举值
//   - longName: 长选项名, 如 "help"
//   - shortName: 短选项名, 如 "h"
//   - default_: 默认值
//
// 返回值:
//   - *BaseFlag[T]: 基础标志实例
//
// 注意事项:
//   - 此函数会初始化内部值指针, 并将默认值复制到值中
//   - 创建后的标志初始状态为未设置(isSet=false)
func NewBaseFlag[T any](flagType types.FlagType, longName, shortName, desc string, default_ T) *BaseFlag[T] {
	value := new(T)
	*value = default_ // 设置默认值

	return &BaseFlag[T]{
		longName:  longName,
		shortName: shortName,
		desc:      desc,
		flagType:  flagType,
		value:     value,
		default_:  default_,
	}
}

// Name 获取标志名称
//
// 返回值:
//   - string: 优先返回长名称, 如果长名称为空则返回短名称
func (f *BaseFlag[T]) Name() string {
	if f.longName != "" {
		return f.longName
	}
	return f.shortName
}

// LongName 获取标志的长名称
//
// 返回值:
//   - string: 长选项名称, 如 "help"
func (f *BaseFlag[T]) LongName() string {
	return f.longName
}

// ShortName 获取标志的短名称
//
// 返回值:
//   - string: 短选项名称, 如 "h"
func (f *BaseFlag[T]) ShortName() string {
	return f.shortName
}

// Desc 获取标志的描述信息
//
// 返回值:
//   - string: 标志描述文本
func (f *BaseFlag[T]) Desc() string {
	return f.desc
}

// Type 获取标志的类型
//
// 返回值:
//   - types.FlagType: 标志类型枚举值
func (f *BaseFlag[T]) Type() types.FlagType {
	return f.flagType
}

// Get 获取标志的当前值
//
// 返回值:
//   - T: 标志的当前值
func (f *BaseFlag[T]) Get() T {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return *f.value
}

// GetStr 获取标志当前值的字符串表示
//
// 返回值:
//   - string: 标志当前值的字符串表示
//
// 功能说明:
//   - 获取标志当前值的字符串表示
//   - 与String()方法不同, 此方法专注于值本身
//   - 用于内置标志处理中获取标志值
func (f *BaseFlag[T]) GetStr() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return fmt.Sprintf("%v", *f.value)
}

// GetDef 获取标志的默认值
//
// 返回值:
//   - any: 标志的默认值, 使用any类型以支持泛型
func (f *BaseFlag[T]) GetDef() any {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.default_
}

// IsSet 检查标志是否已被设置
//
// 返回值:
//   - bool: 如果标志已被设置返回true, 否则返回false
func (f *BaseFlag[T]) IsSet() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.isSet
}

// GetEnvVar 获取关联的环境变量名
//
// 返回值:
//   - string: 环境变量名, 如果未绑定则返回空字符串
func (f *BaseFlag[T]) GetEnvVar() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.envVar
}

// EnumValues 获取枚举类型的可选值
//
// 返回值:
//   - []string: 枚举类型的可选值列表
//
// 功能说明:
//   - 实现 Flag 接口的 EnumValues 方法
//   - 非枚举类型返回空切片
//   - 此方法是线程安全的
func (f *BaseFlag[T]) EnumValues() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 非枚举类型返回空切片
	if f.flagType != types.FlagTypeEnum {
		return []string{}
	}

	// 对于枚举类型, 具体实现应该重写此方法
	return []string{}
}

// Set 设置标志的值
//
// 参数:
//   - value: 要设置的字符串值
//
// 返回值:
//   - error: 如果设置失败返回错误
//
// 注意事项:
//   - 这是基础实现, 具体子类应该重写此方法实现自己的解析逻辑
//   - 基础实现仅返回nil, 不进行任何实际设置操作
func (f *BaseFlag[T]) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 这里仅作为示例, 具体每个子类实现自己的逻辑

	// 这里不进行实际的验证, 子类应该重写此方法实现自己的验证逻辑

	return nil
}

// Reset 重置标志为默认值
//
// 将标志的值重置为默认值, 并将isSet状态设置为false
func (f *BaseFlag[T]) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	*f.value = f.default_
	f.isSet = false
}

// String 返回标志的格式化名称
//
// 返回值:
//   - string: 格式化的标志名称, 用于显示
//
// 注意事项:
//   - 使用utils.FormatFlagName进行格式化
//   - 通常用于帮助信息显示
func (f *BaseFlag[T]) String() string {
	return utils.FormatFlagName(f.longName, f.shortName)
}

// BindEnv 绑定环境变量
//
// 参数:
//   - name: 环境变量名
//
// 注意事项:
//   - 绑定后, 标志可以从指定的环境变量读取值
//   - 环境变量的优先级低于命令行参数
func (f *BaseFlag[T]) BindEnv(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.envVar = name
}

// GetValuePtr 返回值指针, 用于注册到标准库 flag 包
//
// 返回值:
//   - *T: 指向标志值的指针
//
// 注意事项:
// 1. 此方法主要用于与标准库 flag 包集成, 不推荐在常规代码中使用
// 2. 返回的指针指向内部状态, 直接修改可能破坏线程安全
// 3. 仅应在程序初始化阶段 (标志注册时) 使用, 避免并发访问
// 4. 如需在多线程环境中访问标志值, 请使用 Get() 方法
func (f *BaseFlag[T]) GetValuePtr() *T {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.value
}

// SetValidator 设置验证器
//
// 参数:
//   - validator: 验证器函数
//
// 功能说明:
//   - 设置标志的验证器
//   - 如果之前已设置验证器，会被覆盖
//   - 验证器会在 Set 方法中解析完值后被调用
func (f *BaseFlag[T]) SetValidator(validator types.Validator[T]) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.validator = validator
}

// ClearValidator 清除验证器
//
// 功能说明:
//   - 移除标志的验证器
//   - 之后调用 Set 方法将不会进行验证
func (f *BaseFlag[T]) ClearValidator() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.validator = nil
}

// HasValidator 检查是否设置了验证器
//
// 返回值:
//   - bool: 是否设置了验证器
func (f *BaseFlag[T]) HasValidator() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.validator != nil
}
