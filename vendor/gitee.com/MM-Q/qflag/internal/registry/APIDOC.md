# Package registry 

```go
import "gitee.com/MM-Q/qflag/internal/registry"
```

Package registry 提供标志和命令的注册表实现

registry 包基于泛型注册表提供了专门的标志和命令注册表实现。 主要组件: 
  - FlagRegistryImpl: 标志注册表, 支持标志的注册、查找和别名管理
  - CmdRegistryImpl: 命令注册表, 支持命令的注册、查找和别名管理

特性: 
  - 类型安全的注册表接口
  - 支持长名称和短名称查找
  - 别名管理功能
  - 统一的错误处理

---

## FUNCTIONS

### func NewCmdRegistry() types.CmdRegistry

```go
func NewCmdRegistry() types.CmdRegistry
```

NewCmdRegistry 创建新的命令注册表实例

**返回值:**
  - types.CmdRegistry: 命令注册表接口实例

**功能说明: **
  - 创建泛型注册表实例
  - 包装为CmdRegistryImpl
  - 返回接口类型, 隐藏实现细节

---

### func NewFlagRegistry() types.FlagRegistry

```go
func NewFlagRegistry() types.FlagRegistry
```

NewFlagRegistry 创建新的标志注册表实例

**返回值:**
  - types.FlagRegistry: 标志注册表接口实例

**功能说明: **
  - 创建泛型注册表实例
  - 包装为FlagRegistryImpl
  - 返回接口类型, 隐藏实现细节

---

### func NewRegistry[T any]() *registry[T]

```go
func NewRegistry[T any]() *registry[T]
```

NewRegistry 创建新的泛型注册表实例

**返回值:**
  - *registry[T]: 初始化完成的注册表实例

**功能说明: **
  - 初始化空的存储映射
  - 创建名称索引映射
  - 初始化ID计数器
  - 准备接收注册项

---

## TYPES

### type CmdRegistryImpl struct

```go
type CmdRegistryImpl struct {
    *registry[types.Command]
}
```

CmdRegistryImpl 命令注册表实现

CmdRegistryImpl 是 types.CmdRegistry 接口的具体实现, 基于泛型注册表 提供了命令的注册、查找、别名管理等功能。

**特性: **
  - 支持长名称和短名称查找
  - 支持别名添加和管理
  - 提供完整的命令生命周期管理
  - 继承泛型注册表的所有基础功能

#### func (r *CmdRegistryImpl) Clear()

```go
func (r *CmdRegistryImpl) Clear()
```

Clear 清空注册表中的所有命令

**功能说明: **
  - 移除所有命令
  - 重置注册表到初始状态
  - 释放相关内存

#### func (r *CmdRegistryImpl) Count() int

```go
func (r *CmdRegistryImpl) Count() int
```

Count 获取注册表中的命令数量

**返回值:**
  - int: 命令总数

**功能说明: **
  - 返回当前注册的命令数量
  - 时间复杂度O(1)

#### func (r *CmdRegistryImpl) Get(name string) (types.Command, bool)

```go
func (r *CmdRegistryImpl) Get(name string) (types.Command, bool)
```

Get 根据名称获取命令

**参数:**
  - name: 命令名称

**返回值:**
  - types.Command: 找到的命令
  - bool: 是否找到, true表示找到

**功能说明: **
  - 支持长名称查找
  - 直接委托给底层注册表

#### func (r *CmdRegistryImpl) Has(name string) bool

```go
func (r *CmdRegistryImpl) Has(name string) bool
```

Has 检查指定名称的命令是否存在

**参数:**
  - name: 要检查的命令名称

**返回值:**
  - bool: 是否存在, true表示存在

**功能说明: **
  - 快速存在性检查
  - 不返回命令本身, 提高效率

#### func (r *CmdRegistryImpl) List() []types.Command

```go
func (r *CmdRegistryImpl) List() []types.Command
```

List 获取所有注册的命令列表

**返回值:**
  - []types.Command: 所有命令的切片

**功能说明: **
  - 返回注册表中所有命令
  - 顺序不确定, 取决于map遍历顺序

#### func (r *CmdRegistryImpl) Range(f func(name string, cmd types.Command) bool)

```go
func (r *CmdRegistryImpl) Range(f func(name string, cmd types.Command) bool)
```

Range 遍历注册表中的所有命令

**参数:**
  - f: 遍历函数, 接收名称和命令, 返回是否继续遍历

**功能说明: **
  - 按注册顺序遍历所有命令
  - 支持提前终止遍历
  - 遍历过程中修改注册表可能导致不确定行为

#### func (r *CmdRegistryImpl) Register(cmd types.Command) error

```go
func (r *CmdRegistryImpl) Register(cmd types.Command) error
```

Register 注册新命令到注册表

**参数:**
  - cmd: 要注册的命令, 不能为nil

**返回值:**
  - error: 注册失败时返回错误

**错误情况: **
  - 命令为nil: 返回 INVALID_COMMAND 错误
  - 命令名称为空: 返回 INVALID_NAME 错误
  - 命令名称已存在: 返回 ErrAlreadyExists 错误
  - 长名称和短名称都为空: 返回 INVALID_NAME 错误

**功能说明: **
  - 验证命令有效性
  - 提取命令名称
  - 调用底层注册方法 (支持长名称和短名称) 
  - 命令对象只存储一次, 长名称和短名称都指向同一个对象

#### func (r *CmdRegistryImpl) Unregister(name string) error

```go
func (r *CmdRegistryImpl) Unregister(name string) error
```

Unregister 从注册表中移除指定命令

**参数:**
  - name: 要移除的命令名称

**返回值:**
  - error: 移除失败时返回错误

**错误情况: **
  - 命令不存在: 返回 ErrCmdNotFound 错误

**功能说明: **
  - 调用底层移除方法
  - 自动清理相关索引

---

### type FlagRegistryImpl struct

```go
type FlagRegistryImpl struct {
    *registry[types.Flag]
}
```

FlagRegistryImpl 标志注册表实现

FlagRegistryImpl 是 types.FlagRegistry 接口的具体实现, 基于泛型注册表 提供了标志的注册、查找、别名管理等功能。

**特性: **
  - 支持长名称和短名称查找
  - 支持别名添加和管理
  - 提供完整的标志生命周期管理
  - 继承泛型注册表的所有基础功能

#### func (r *FlagRegistryImpl) Clear()

```go
func (r *FlagRegistryImpl) Clear()
```

Clear 清空注册表中的所有标志

**功能说明: **
  - 移除所有标志
  - 重置注册表到初始状态
  - 释放相关内存

#### func (r *FlagRegistryImpl) Count() int

```go
func (r *FlagRegistryImpl) Count() int
```

Count 获取注册表中的标志数量

**返回值:**
  - int: 标志总数

**功能说明: **
  - 返回当前注册的标志数量
  - 时间复杂度O(1)

#### func (r *FlagRegistryImpl) Get(name string) (types.Flag, bool)

```go
func (r *FlagRegistryImpl) Get(name string) (types.Flag, bool)
```

Get 根据名称获取标志

**参数:**
  - name: 标志名称

**返回值:**
  - types.Flag: 找到的标志
  - bool: 是否找到, true表示找到

**功能说明: **
  - 支持长名称查找
  - 直接委托给底层注册表

#### func (r *FlagRegistryImpl) Has(name string) bool

```go
func (r *FlagRegistryImpl) Has(name string) bool
```

Has 检查指定名称的标志是否存在

**参数:**
  - name: 要检查的标志名称

**返回值:**
  - bool: 是否存在, true表示存在

**功能说明: **
  - 快速存在性检查
  - 不返回标志本身, 提高效率

#### func (r *FlagRegistryImpl) List() []types.Flag

```go
func (r *FlagRegistryImpl) List() []types.Flag
```

List 获取所有注册的标志列表

**返回值:**
  - []types.Flag: 所有标志的切片

**功能说明: **
  - 返回注册表中所有标志
  - 顺序不确定, 取决于map遍历顺序

#### func (r *FlagRegistryImpl) Range(f func(name string, flag types.Flag) bool)

```go
func (r *FlagRegistryImpl) Range(f func(name string, flag types.Flag) bool)
```

Range 遍历注册表中的所有标志

**参数:**
  - f: 遍历函数, 接收名称和标志, 返回是否继续遍历

**功能说明: **
  - 按注册顺序遍历所有标志
  - 支持提前终止遍历
  - 遍历过程中修改注册表可能导致不确定行为

#### func (r *FlagRegistryImpl) Register(flag types.Flag) error

```go
func (r *FlagRegistryImpl) Register(flag types.Flag) error
```

Register 注册新标志到注册表

**参数:**
  - flag: 要注册的标志, 不能为nil

**返回值:**
  - error: 注册失败时返回错误

**错误情况: **
  - 标志为nil: 返回 INVALID_FLAG 错误
  - 标志名称为空: 返回 INVALID_NAME 错误
  - 标志名称已存在: 返回 ErrAlreadyExists 错误
  - 长名称和短名称都为空: 返回 INVALID_NAME 错误

**功能说明: **
  - 验证标志有效性
  - 提取标志名称
  - 调用底层注册方法 (支持长名称和短名称) 
  - 标志对象只存储一次, 长名称和短名称都指向同一个对象

#### func (r *FlagRegistryImpl) Unregister(name string) error

```go
func (r *FlagRegistryImpl) Unregister(name string) error
```

Unregister 从注册表中移除指定标志

**参数:**
  - name: 要移除的标志名称

**返回值:**
  - error: 移除失败时返回错误

**错误情况: **
  - 标志不存在: 返回 ErrFlagNotFound 错误

**功能说明: **
  - 调用底层移除方法
  - 自动清理相关索引