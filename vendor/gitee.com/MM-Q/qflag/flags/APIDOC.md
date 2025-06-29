# Package flags

`flags` 定义了所有标志类型的通用接口和基础标志结构体。

## CONSTANTS

定义标志的分隔符常量：

```go
const (
    FlagSplitComma     = "," // 逗号
    FlagSplitSemicolon = ";" // 分号
    FlagSplitPipe      = "|" // 竖线
    FlagKVColon        = ":" // 冒号
    FlagKVEqual        = "=" // 等号
)
```

定义非法字符集常量，防止非法的标志名称：

```go
const InvalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
```

## VARIABLES

内置标志名称：

```go
var (
    HelpFlagName            = "help"    // 帮助标志名称
    HelpFlagShortName       = "h"       // 帮助标志短名称
    ShowInstallPathFlagName = "sip"     // 显示安装路径标志名称
    VersionFlagLongName     = "version" // 版本标志名称
    VersionFlagShortName    = "v"       // 版本标志短名称
)
```

`Flag` 支持的标志分隔符切片：

```go
var FlagSplitSlice = []string{
    FlagSplitComma,     // 逗号
    FlagSplitSemicolon, // 分号
    FlagSplitPipe,      // 竖线
    FlagKVColon,        // 冒号
}
```

## FUNCTIONS

`FlagTypeToString` 将 `FlagType` 转换为字符串：

```go
func FlagTypeToString(flagType FlagType) string {
    switch flagType {
    case FlagTypeInt:
        return "<int>"
    case FlagTypeInt64:
        return "<int64>"
    case FlagTypeUint16:
        return "<uint16>"
    case FlagTypeString:
        return "<string>"
    case FlagTypeBool:
        // 布尔类型没有参数类型字符串
        return ""
    case FlagTypeFloat64:
        return "<float64>"
    case FlagTypeEnum:
        return "<enum>"
    case FlagTypeDuration:
        return "<duration>"
    case FlagTypeTime:
        return "<time>"
    case FlagTypeMap:
        return "<map>"
    case FlagTypePath:
        return "<path>"
    default:
        return "<unknown>"
    }
}
```

## TYPES

### BaseFlag

泛型基础标志结构体，封装所有标志的通用字段和方法：

```go
type BaseFlag[T any] struct {
    longName    string     // 长标志名称
    shortName   string     // 短标志字符
    defValue    T          // 默认值
    usage       string     // 帮助说明
    value       *T         // 标志值指针
    mu          sync.Mutex // 并发访问锁
    validator   Validator  // 验证器接口
    initialized bool       // 标志是否已初始化
    isSet       bool       // 标志是否已被设置值
}
```

- `Get` 获取标志的实际值，优先级：已设置的值 > 默认值，线程安全：

```go
func (f *BaseFlag[T]) Get() T {
    f.mu.Lock()
    defer f.mu.Unlock()

    // 如果标志值不为空，则返回标志值
    if f.value != nil {
        return *f.value
    }

    // 否则返回默认值
    return f.defValue
}
```

- `GetDefault` 获取标志的默认值：

```go
func (f *BaseFlag[T]) GetDefault() T { return f.defValue }
```

- `GetDefaultAny` 获取标志的默认值（`any` 类型）：

```go
func (f *BaseFlag[T]) GetDefaultAny() any { return f.defValue }
```

- `GetPointer` 返回标志值的指针，注意获取指针过程受锁保护，但直接修改指针指向的值仍会绕过验证机制，多线程环境下修改时需额外同步措施，建议优先使用 `Set()` 方法：

```go
func (f *BaseFlag[T]) GetPointer() *T {
    f.mu.Lock()
    defer f.mu.Unlock()
    return f.value
}
```

- `Init` 初始化标志的元数据和值指针，无需显式调用，仅在创建标志对象时自动调用：

```go
func (f *BaseFlag[T]) Init(longName, shortName string, defValue T, usage string, value *T) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // 检查是否已初始化
    if f.initialized {
        return fmt.Errorf("flag %s/%s already initialized", f.shortName, f.longName)
    }

    // 检查长短标志是否同时为空
    if longName == "" && shortName == "" {
        return fmt.Errorf("longName and shortName cannot both be empty")
    }

    // 验证值指针（避免后续空指针异常）
    if value == nil {
        return fmt.Errorf("value pointer cannot be nil")
    }

    f.longName = longName   // 初始化长标志名
    f.shortName = shortName // 初始化短标志名
    f.defValue = defValue   // 初始化默认值
    f.usage = usage         // 初始化标志用途
    f.value = value         // 初始化值指针
    f.initialized = true    // 设置初始化完成标志

    return nil
}
```

- `IsSet` 判断标志是否已被设置值，返回值：`true` 表示已设置值，`false` 表示未设置：

```go
func (f *BaseFlag[T]) IsSet() bool {
    f.mu.Lock()
    defer f.mu.Unlock()
    return f.isSet
}
```

- `LongName` 获取标志的长名称：

```go
func (f *BaseFlag[T]) LongName() string { return f.longName }
```

- `Reset` 将标志值重置为默认值，线程安全：

```go
func (f *BaseFlag[T]) Reset() {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.value = nil   // 重置为未设置状态，下次 Get() 会返回默认值
    f.isSet = false // 重置设置状态
}
```

- `Set` 设置标志的值：

```go
func (f *BaseFlag[T]) Set(value T) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // 创建一个副本
    v := value

    // 设置标志值前先进行验证
    if f.validator != nil {
        if err := f.validator.Validate(v); err != nil {
            return fmt.Errorf("invalid value for %s: %w", f.longName, err)
        }
    }

    // 设置标志值
    f.value = &v

    // 标志已设置
    f.isSet = true

    return nil
}
```

- `SetValidator` 设置标志的验证器：

```go
func (f *BaseFlag[T]) SetValidator(validator Validator) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.validator = validator
}
```

- `ShortName` 获取标志的短名称：

```go
func (f *BaseFlag[T]) ShortName() string { return f.shortName }
```

- `String` 返回标志的字符串表示：

```go
func (f *BaseFlag[T]) String() string {
    return fmt.Sprint(f.Get())
}
```

- `Usage` 获取标志的用法说明：

```go
func (f *BaseFlag[T]) Usage() string { return f.usage }
```

### BoolFlag

布尔类型标志结构体，继承 `BaseFlag[bool]` 泛型结构体，实现 `Flag` 接口：

```go
type BoolFlag struct {
    BaseFlag[bool]
}

func (f *BoolFlag) Type() FlagType { return FlagTypeBool }
```

### DurationFlag

时间间隔类型标志结构体，继承 `BaseFlag[time.Duration]` 泛型结构体，实现 `Flag` 接口：

```go
type DurationFlag struct {
    BaseFlag[time.Duration]
}

// Set 实现 flag.Value 接口，解析并设置时间间隔值
func (f *DurationFlag) Set(value string) error {
    // 检查空值
    if value == "" {
        return fmt.Errorf("duration cannot be empty")
    }

    // 将单位转换为小写，确保解析的准确性
    lowercaseValue := strings.ToLower(value)

    // 解析时间间隔字符串
    duration, err := time.ParseDuration(lowercaseValue)
    if err != nil {
        return fmt.Errorf("invalid duration format: %v (valid units: ns/us/ms/s/m/h)", err)
    }

    // 检查负值（可选）
    if duration < 0 {
        return fmt.Errorf("negative duration not allowed")
    }

    // 调用基类方法设置值
    return f.BaseFlag.Set(duration)
}

// String 实现 flag.Value 接口，返回当前值的字符串表示
func (f *DurationFlag) String() string {
    return f.Get().String()
}

func (f *DurationFlag) Type() FlagType { return FlagTypeDuration }
```

### EnumFlag

枚举类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，增加枚举特有的选项验证：

```go
type EnumFlag struct {
    BaseFlag[string]
    optionMap map[string]bool // 枚举值映射
}

// Init 初始化枚举类型标志，无需显式调用，仅在创建标志对象时自动调用
func (f *EnumFlag) Init(longName, shortName string, defValue string, usage string, options []string) error {
    // 初始化枚举值
    if options == nil {
        options = make([]string, 0)
    }

    // 1. 初始化基类字段
    valuePtr := new(string)

    // 默认值小写处理
    *valuePtr = strings.ToLower(defValue)

    // 调用基类方法初始化字段
    if err := f.BaseFlag.Init(longName, shortName, defValue, usage, valuePtr); err != nil {
        return err
    }

    // 2. 初始化枚举 optionMap（仅在 Init 阶段修改，无需额外锁）
    f.optionMap = make(map[string]bool)
    for _, opt := range options {
        if opt == "" {
            return fmt.Errorf("enum option cannot be empty")
        }
        f.optionMap[strings.ToLower(opt)] = true
    }

    // 3. 验证默认值有效性
    if len(options) > 0 && !f.optionMap[strings.ToLower(defValue)] {
        return fmt.Errorf("default value '%s' not in enum options %v", defValue, options)
    }

    return nil
}

// IsCheck 检查枚举值是否有效
func (f *EnumFlag) IsCheck(value string) error {
    // 如果枚举 map 为空，则不需要检查
    if len(f.optionMap) == 0 {
        return nil
    }

    // 转换为小写
    value = strings.ToLower(value)

    // 检查值是否在枚举 map 中
    if _, valid := f.optionMap[value]; !valid {
        var options []string
        for k := range f.optionMap {
            // 添加枚举值
            options = append(options, k)
        }
        return fmt.Errorf("invalid enum value '%s', options are %v", value, options)
    }
    return nil
}

// Set 实现 flag.Value 接口，解析并设置枚举值
func (f *EnumFlag) Set(value string) error {
    // 先验证值是否有效
    if err := f.IsCheck(value); err != nil {
        return err
    }
    // 调用基类方法设置值
    return f.BaseFlag.Set(value)
}

// String 实现 flag.Value 接口，返回当前值的字符串表示
func (f *EnumFlag) String() string { return f.Get() }

func (f *EnumFlag) Type() FlagType { return FlagTypeEnum }
```

### Flag

所有标志类型的通用接口，定义了标志的元数据访问方法：

```go
type Flag interface {
    LongName() string   // 获取标志的长名称
    ShortName() string  // 获取标志的短名称
    Usage() string      // 获取标志的用法
    Type() FlagType     // 获取标志类型
    GetDefaultAny() any // 获取标志的默认值（any 类型）
    String() string     // 获取标志的字符串表示
    IsSet() bool        // 判断标志是否已设置值
    Reset()             // 重置标志值为默认值
}
```

### FlagMeta

统一存储标志的完整元数据：

```go
type FlagMeta struct {
    Flag Flag // 标志对象
}

func (m *FlagMeta) GetDefault() any { return m.Flag.GetDefaultAny() }
func (m *FlagMeta) GetFlag() Flag { return m.Flag }
func (m *FlagMeta) GetFlagType() FlagType { return m.Flag.Type() }
func (m *FlagMeta) GetLongName() string { return m.Flag.LongName() }
func (m *FlagMeta) GetShortName() string { return m.Flag.ShortName() }
func (m *FlagMeta) GetUsage() string { return m.Flag.Usage() }
```

### FlagMetaInterface

标志元数据接口，定义了标志元数据的获取方法：

```go
type FlagMetaInterface interface {
    GetFlagType() FlagType // 获取标志类型
    GetFlag() Flag         // 获取标志对象
    GetLongName() string   // 获取标志的长名称
    GetShortName() string  // 获取标志的短名称
    GetUsage() string      // 获取标志的用法描述
    GetDefault() any       // 获取标志的默认值
    GetValue() any         // 获取标志的当前值
}
```

### FlagRegistry

集中管理所有标志元数据及索引：

```go
type FlagRegistry struct {
    mu       sync.RWMutex         // 并发访问锁
    byLong   map[string]*FlagMeta // 按长名称索引
    byShort  map[string]*FlagMeta // 按短名称索引
    allFlags []*FlagMeta          // 所有标志元数据列表
}

// 创建一个空的标志注册表
func NewFlagRegistry() *FlagRegistry {
    return &FlagRegistry{
        mu:       sync.RWMutex{},
        byLong:   make(map[string]*FlagMeta),
        byShort:  make(map[string]*FlagMeta),
        allFlags: make([]*FlagMeta, 0),
    }
}

// GetAllFlags 获取所有标志元数据列表
func (r *FlagRegistry) GetAllFlags() []*FlagMeta {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.allFlags
}

// GetByLong 通过长标志名称查找对应的标志元数据
func (r *FlagRegistry) GetByLong(longName string) (*FlagMeta, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    meta, exists := r.byLong[longName]
    return meta, exists
}

// GetByName 通过标志名称查找标志元数据
func (r *FlagRegistry) GetByName(name string) (*FlagMeta, bool) {
    // 先尝试按长名称查找
    if meta, exists := r.GetByLong(name); exists {
        return meta, exists
    }

    // 再尝试按短名称查找
    if meta, exists := r.GetByShort(name); exists {
        return meta, exists
    }

    // 未找到
    return nil, false
}

// GetByShort 通过短标志名称查找对应的标志元数据
func (r *FlagRegistry) GetByShort(shortName string) (*FlagMeta, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    meta, exists := r.byShort[shortName]
    return meta, exists
}

// GetLongFlags 获取长标志映射
func (r *FlagRegistry) GetLongFlags() map[string]*FlagMeta {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.byLong
}

// GetShortFlags 获取短标志映射
func (r *FlagRegistry) GetShortFlags() map[string]*FlagMeta {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.byShort
}

// RegisterFlag 注册一个新的标志元数据到注册表中
func (r *FlagRegistry) RegisterFlag(meta *FlagMeta) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 检查长短标志是否都为空
    if meta.GetLongName() == "" && meta.GetShortName() == "" {
        return fmt.Errorf("flag must have at least one name")
    }

    // 检查长标志是否已存在
    if meta.GetLongName() != "" {
        if _, exists := r.byLong[meta.GetLongName()]; exists {
            return fmt.Errorf("long flag %s already exists", meta.GetLongName())
        }
    }

    // 检查短标志是否已存在
    if meta.GetShortName() != "" {
        if _, exists := r.byShort[meta.GetShortName()]; exists {
            return fmt.Errorf("short flag %s already exists", meta.GetShortName())
        }
    }

    // 添加长标志索引
    if meta.GetLongName() != "" {
        r.byLong[meta.GetLongName()] = meta
    }

    // 添加短标志索引
    if meta.GetShortName() != "" {
        r.byShort[meta.GetShortName()] = meta
    }

    // 添加到所有标志列表
    r.allFlags = append(r.allFlags, meta)

    return nil
}
```

### FlagRegistryInterface

标志注册表接口，定义了标志元数据的增删改查操作：

```go
type FlagRegistryInterface interface {
    GetAllFlags() []*FlagMeta                      // 获取所有标志元数据列表
    GetLongFlags() map[string]*FlagMeta            // 获取长标志映射
    GetShortFlags() map[string]*FlagMeta           // 获取短标志映射
    RegisterFlag(meta *FlagMeta) error             // 注册一个新的标志元数据到注册表中
    GetByLong(longName string) (*FlagMeta, bool)   // 通过长标志名称查找对应的标志元数据
    GetByShort(shortName string) (*FlagMeta, bool) // 通过短标志名称查找对应的标志元数据
    GetByName(name string) (*FlagMeta, bool)       // 通过标志名称查找标志元数据
}
```

### FlagType

标志类型：

```go
type FlagType int

const (
    FlagTypeInt      FlagType = iota + 1 // 整数类型
    FlagTypeInt64                        // 64 位整数类型
    FlagTypeUint16                       // 16 位无符号整数类型
    FlagTypeString                       // 字符串类型
    FlagTypeBool                         // 布尔类型
    FlagTypeFloat64                      // 64 位浮点数类型
    FlagTypeEnum                         // 枚举类型
    FlagTypeDuration                     // 时间间隔类型
    FlagTypeSlice                        // 切片类型
    FlagTypeTime                         // 时间类型
    FlagTypeMap                          // 映射类型
    FlagTypePath                         // 路径类型
)
```

### Float64Flag

浮点型标志结构体，继承 `BaseFlag[float64]` 泛型结构体，实现 `Flag` 接口：

```go
type Float64Flag struct {
    BaseFlag[float64]
}

func (f *Float64Flag) Type() FlagType { return FlagTypeFloat64 }
```

### Int64Flag

64 位整数类型标志结构体，继承 `BaseFlag[int64]` 泛型结构体，实现 `Flag` 接口：

```go
type Int64Flag struct {
    BaseFlag[int64]
}

// SetRange 设置 64 位整数的有效范围
func (f *Int64Flag) SetRange(min, max int64) {
    validator := &validator.IntRangeValidator64{Min: min, Max: max}
    f.SetValidator(validator)
}

func (f *Int64Flag) Type() FlagType { return FlagTypeInt64 }
```

### IntFlag

整数类型标志结构体，继承 `BaseFlag[int]` 泛型结构体，实现 `Flag` 接口：

```go
type IntFlag struct {
    BaseFlag[int]
}

// SetRange 设置整数的有效范围
func (f *IntFlag) SetRange(min, max int) {
    validator := &validator.IntRangeValidator{Min: min, Max: max}
    f.SetValidator(validator)
}

func (f *IntFlag) Type() FlagType { return FlagTypeInt }
```

### MapFlag

键值对类型标志结构体，继承 `BaseFlag[map[string]string]` 泛型结构体，实现 `Flag` 接口：

```go
type MapFlag struct {
    BaseFlag[map[string]string]
    keyDelimiter   string     // 键值对之间的分隔符
    valueDelimiter string     // 键和值之间的分隔符
    mu             sync.Mutex // 互斥锁
    IgnoreCase     bool       // 是否忽略键的大小写
}

// Set 实现 flag.Value 接口，解析并设置键值对
func (f *MapFlag) Set(value string) error {
    if value == "" {
        return fmt.Errorf("map value cannot be empty")
    }

    // 获取当前值
    current := f.Get()
    if current == nil {
        current = make(map[string]string)
    }

    f.mu.Lock()
    defer f.mu.Unlock()

    // 使用键分隔符分割多个键值对
    pairs := strings.Split(value, f.keyDelimiter)
    for _, pair := range pairs {
        // 使用值分隔符分割键和值
        kv := strings.SplitN(pair, f.valueDelimiter, 2)

        // 检查键值对是否包含两个部分
        if len(kv) != 2 {
            return fmt.Errorf("invalid key-value pair format: %s", pair)
        }

        // 去除键和值的前后空格
        key := strings.TrimSpace(kv[0])
        val := strings.TrimSpace(kv[1])

        // 如果需要忽略大小写，则将键转换为小写
        if f.IgnoreCase {
            key = strings.ToLower(key)
        }

        // 检查键和值是否为空
        if key == "" {
            return fmt.Errorf("empty key in key-value pair: %s", pair)
        }
        if val == "" {
            return fmt.Errorf("empty value in key-value pair: %s", pair)
        }

        // 更新当前值
        current[key] = val
    }

    return f.BaseFlag.Set(current)
}

// SetDelimiters 设置键值对分隔符
func (f *MapFlag) SetDelimiters(keyDelimiter, valueDelimiter string) {
    f.mu.Lock()
    defer f.mu.Unlock()
    if keyDelimiter == "" {
        keyDelimiter = FlagSplitComma // 默认使用逗号
    }
    if valueDelimiter == "" {
        valueDelimiter = FlagKVEqual // 默认使用等号
    }

    // 设置分隔符
    f.keyDelimiter = keyDelimiter
    f.valueDelimiter = valueDelimiter
}

// SetIgnoreCase 设置是否忽略键的大小写
func (f *MapFlag) SetIgnoreCase(enable bool) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.IgnoreCase = enable
}

// String 实现 flag.Value 接口，返回当前值的字符串表示
func (f *MapFlag) String() string {
    m := f.Get()
    if m == nil {
        return ""
    }
    var parts []string
    for k, v := range m {
        parts = append(parts, fmt.Sprintf("%s%s%s", k, f.valueDelimiter, v))
    }
    return strings.Join(parts, f.keyDelimiter)
}

func (f *MapFlag) Type() FlagType { return FlagTypeMap }
```

### PathFlag

路径类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，实现 `Flag` 接口：

```go
type PathFlag struct {
    BaseFlag[string]
}

// Init 初始化路径标志
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error {
    // 初始化路径标志值指针
    valuePtr := new(string)

    // 规范化默认路径为绝对路径
    absDefValue, err := filepath.Abs(defValue)
    if err != nil {
        return fmt.Errorf("failed to normalize default path: %v", err)
    }

    // 设置默认值
    *valuePtr = absDefValue

    // 调用基类方法初始化
    if err := f.BaseFlag.Init(longName, shortName, defValue, usage, valuePtr); err != nil {
        return err
    }

    // 设置路径验证器
    f.SetValidator(&validator.PathValidator{})
    return nil
}

// Set 实现 flag.Value 接口，解析并验证路径
func (f *PathFlag) Set(value string) error {
    if value == "" {
        return fmt.Errorf("path cannot be empty")
    }

    // 规范化路径为绝对路径
    absPath, err := filepath.Abs(value)
    if err != nil {
        return fmt.Errorf("failed to get absolute path: %v", err)
    }

    // 调用基类方法设置值（会触发验证器验证）
    return f.BaseFlag.Set(absPath)
}

// String 实现 flag.Value 接口，返回当前值的字符串表示
func (f *PathFlag) String() string { return f.Get() }

func (f *PathFlag) Type() FlagType { return FlagTypePath }
```

### SliceFlag

切片类型标志结构体，继承 `BaseFlag[[]string]` 泛型结构体，实现 `Flag` 接口：

```go
type SliceFlag struct {
    BaseFlag[[]string]            // 基类
    delimiters         []string   // 分隔符
    mu                 sync.Mutex // 锁
    SkipEmpty          bool       // 是否跳过空元素
}

// Clear 清空切片所有元素
func (f *SliceFlag) Clear() {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.value = &[]string{}
}

// Contains 检查切片是否包含指定元素
func (f *SliceFlag) Contains(element string) bool {
    current := f.Get()
    for _, item := range current {
        if item == element {
            return true
        }
    }
    return false
}

// GetDelimiters 获取当前分隔符列表
func (f *SliceFlag) GetDelimiters() []string {
    f.mu.Lock()
    defer f.mu.Unlock()
    res := make([]string, len(f.delimiters))
    copy(res, f.delimiters)
    return res
}

// Init 初始化切片类型标志
func (f *SliceFlag) Init(longName, shortName string, defValue []string, usage string) error {
    if defValue == nil {
        defValue = make([]string, 0)
    }

    valueCopy := make([]string, len(defValue))
    copy(valueCopy, defValue)
    valuePtr := &valueCopy

    if err := f.BaseFlag.Init(longName, shortName, defValue, usage, valuePtr); err != nil {
        return err
    }

    f.SetDelimiters(FlagSplitSlice)

    return nil
}

// Len 获取切片长度
func (f *SliceFlag) Len() int {
    f.mu.Lock()
    defer f.mu.Unlock()
    if f.value == nil {
        return 0
    }
    return len(*f.value)
}

// Remove 从切片中移除指定元素
func (f *SliceFlag) Remove(element string) error {
    current := f.Get()
    newSlice := []string{}
    for _, item := range current {
        if item != element {
            newSlice = append(newSlice, item)
        }
    }
    return f.BaseFlag.Set(newSlice)
}

// Set 实现 flag.Value 接口，解析并设置切片值
func (f *SliceFlag) Set(value string) error {
    if value == "" {
        return fmt.Errorf("slice cannot be empty")
    }

    current := f.Get()
    elements := []string{}

    f.mu.Lock()
    defer f.mu.Unlock()

    found := false
    for _, delimiter := range f.delimiters {
        if strings.Contains(value, delimiter) {
            elements = strings.Split(value, delimiter)
            for i, e := range elements {
                elements[i] = strings.TrimSpace(e)
            }
            found = true
            break
        }
    }

    if !found {
        elements = []string{strings.TrimSpace(value)}
    }

    if f.SkipEmpty {
        filtered := make([]string, 0, len(elements))
        for _, e := range elements {
            if e != "" {
                filtered = append(filtered, e)
            }
        }
        elements = filtered
    }

    newValues := make([]string, 0, len(current)+len(elements))
    newValues = append(newValues, current...)
    newValues = append(newValues, elements...)

    return f.BaseFlag.Set(newValues)
}

// SetDelimiters 设置切片解析的分隔符列表
func (f *SliceFlag) SetDelimiters(delimiters []string) {
    f.mu.Lock()
    defer f.mu.Unlock()

    if len(delimiters) == 0 {
        delimiters = FlagSplitSlice
    }

    f.delimiters = delimiters
}

// SetSkipEmpty 设置是否跳过空元素
func (f *SliceFlag) SetSkipEmpty(skip bool) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.SkipEmpty = skip
}

// Sort 对切片进行排序
func (f *SliceFlag) Sort() error {
    current := f.Get()
    sort.Strings(current)
    return f.BaseFlag.Set(current)
}

// String 实现 flag.Value 接口，返回当前值的字符串表示
func (f *SliceFlag) String() string {
    return strings.Join(f.Get(), ",")
}

func (f *SliceFlag) Type() FlagType { return FlagTypeSlice }
```

### StringFlag

字符串类型标志结构体，继承 `BaseFlag[string]` 泛型结构体，实现 `Flag` 接口：

```go
type StringFlag struct {
    BaseFlag[string]
}

// Contains 检查字符串是否包含指定子串
func (f *StringFlag) Contains(substr string) bool {
    return strings.Contains(f.Get(), substr)
}

// Len 获取字符串标志的长度
func (f *StringFlag) Len() int {
    return len(f.Get())
}

// String 返回带引号的字符串值
func (f *StringFlag) String() string {
    return fmt.Sprintf("%q", f.Get())
}

// ToLower 将字符串标志值转换为小写
func (f *StringFlag) ToLower() string {
    return strings.ToLower(f.Get())
}

// ToUpper 将字符串标志值转换为大写
func (f *StringFlag) ToUpper() string {
    return strings.ToUpper(f.Get())
}

func (f *StringFlag) Type() FlagType { return FlagTypeString }
```

### TimeFlag

时间类型标志结构体，继承 `BaseFlag[time.Time]` 泛型结构体，实现 `Flag` 接口：

```go
type TimeFlag struct {
    BaseFlag[time.Time]
    outputFormat string // 自定义输出格式
}

// Set 实现 flag.Value 接口，解析并设置时间值
func (f *TimeFlag) Set(value string) error {
    var t time.Time
    var err error

    for _, format := range supportedTimeFormats {
        t, err = time.Parse(format, value)
        if err == nil {
            break
        }
    }

    if err != nil {
        return fmt.Errorf("invalid time format: %v (supported formats include %v)", err, supportedTimeFormats)
    }

    return f.BaseFlag.Set(t)
}

// SetOutputFormat 设置时间输出格式
func (f *TimeFlag) SetOutputFormat(format string) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.outputFormat = format
}

// String 实现 flag.Value 接口，返回当前时间的字符串表示
func (f *TimeFlag) String() string {
    f.mu.Lock()
    defer f.mu.Unlock()
    value := f.value
    format := f.outputFormat
    if format != "" {
        return value.Format(format)
    }
    return value.Format(time.RFC3339) // 默认格式
}

func (f *TimeFlag) Type() FlagType { return FlagTypeTime }
```

### TypedFlag

所有标志类型的通用接口，定义了标志的元数据访问方法和默认值访问方法：

```go
type TypedFlag[T any] interface {
    Flag                    // 继承标志接口
    GetDefault() T          // 获取标志的具体类型默认值
    Get() T                 // 获取标志的具体类型值
    GetPointer() *T         // 获取标志值的指针
    Set(T) error            // 设置标志的具体类型值
    SetValidator(Validator) // 设置标志的验证器
}
```

### Uint16Flag

16 位无符号整数类型标志结构体，继承 `BaseFlag[uint16]` 泛型结构体，实现 `Flag` 接口：

```go
type Uint16Flag struct {
    BaseFlag[uint16]
}

// Set 实现 flag.Value 接口，解析并设置 16 位无符号整数值
func (f *Uint16Flag) Set(value string) error {
    num, err := strconv.ParseUint(value, 10, 16)
    if err != nil {
        return fmt.Errorf("invalid uint16 value: %v", err)
    }
    val := uint16(num)
    return f.BaseFlag.Set(val)
}

// String 实现 flag.Value 接口，返回当前值的字符串表示
func (f *Uint16Flag) String() string {
    return fmt.Sprint(f.Get())
}

func (f *Uint16Flag) Type() FlagType { return FlagTypeUint16 }
```

### Validator

验证器接口，所有自定义验证器需实现此接口：

```go
type Validator interface {
    // Validate 验证参数值是否合法
    Validate(value any) error
}
```