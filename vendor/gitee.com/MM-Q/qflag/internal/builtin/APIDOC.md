# Package builtin 

```go
import "gitee.com/MM-Q/qflag/internal/builtin"
```

---

## TYPES

### type BuiltinFlagManager struct

```go
type BuiltinFlagManager struct {
    // Has unexported fields.
}
```

BuiltinFlagManager 内置标志管理器

BuiltinFlagManager 负责管理所有内置标志的注册和处理。 它维护一个处理器映射表, 根据标志类型找到对应的处理器。

#### func NewBuiltinFlagManager() *BuiltinFlagManager

```go
func NewBuiltinFlagManager() *BuiltinFlagManager
```

NewBuiltinFlagManager 创建内置标志管理器

**返回值:**
  - *BuiltinFlagManager: 内置标志管理器实例

**功能说明: **
  - 初始化处理器映射表和标志名映射表
  - 注册默认的内置标志处理器

#### func (m *BuiltinFlagManager) HandleBuiltinFlags(cmd types.Command) error

```go
func (m *BuiltinFlagManager) HandleBuiltinFlags(cmd types.Command) error
```

HandleBuiltinFlags 处理内置标志

**参数:**
  - cmd: 要处理标志的命令

**返回值:**
  - error: 处理失败时返回错误

**功能说明: **
  - 遍历命令的所有标志, 检查是否是内置标志
  - 如果是内置标志且被设置, 则执行对应的处理器
  - 处理器通常会执行操作并退出程序

#### func (m *BuiltinFlagManager) RegisterBuiltinFlags(cmd types.Command) error

```go
func (m *BuiltinFlagManager) RegisterBuiltinFlags(cmd types.Command) error
```

RegisterBuiltinFlags 注册内置标志

**参数:**
  - cmd: 要注册标志的命令

**返回值:**
  - error: 注册失败时返回错误

**功能说明: **
  - 遍历所有处理器, 检查是否应该注册对应的标志
  - 根据命令的语言设置使用相应的描述信息
  - 创建并注册标志到命令中

#### func (m *BuiltinFlagManager) RegisterHandler(handler types.BuiltinFlagHandler)

```go
func (m *BuiltinFlagManager) RegisterHandler(handler types.BuiltinFlagHandler)
```

RegisterHandler 注册内置标志处理器

**参数:**
  - handler: 要注册的处理器

**功能说明: **
  - 将处理器添加到处理器映射表
  - 注册处理器的标志名映射
  - 支持长名称和短名称的映射

---

### type CompletionHandler struct

```go
type CompletionHandler struct{}
```

CompletionHandler 补全标志处理器

CompletionHandler 负责处理补全标志 (--completion) 当用户指定补全标志时, 会生成对应的Shell自动补全脚本。

#### func (h *CompletionHandler) Handle(cmd types.Command) error

```go
func (h *CompletionHandler) Handle(cmd types.Command) error
```

Handle 处理补全标志

**参数:**
  - cmd: 要处理的命令

**返回值:**
  - error: 处理失败时返回错误

**功能说明: **
  - 从命令行参数获取Shell类型
  - 生成对应的补全脚本
  - 输出脚本并退出程序

#### func (h *CompletionHandler) ShouldRegister(cmd types.Command) bool

```go
func (h *CompletionHandler) ShouldRegister(cmd types.Command) bool
```

ShouldRegister 判断是否应该注册补全标志

**参数:**
  - cmd: 要检查的命令

**返回值:**
  - bool: 总是返回true

**功能说明: **
  - 补全标志总是注册, 因为所有命令都应该支持补全功能
  - 补全标志只在根命令中注册

#### func (h *CompletionHandler) Type() types.BuiltinFlagType

```go
func (h *CompletionHandler) Type() types.BuiltinFlagType
```

Type 返回标志类型

**返回值:**
  - types.BuiltinFlagType: CompletionFlag

#### func (h *CompletionHandler) ShouldSkipRegistration(cmd types.Command) bool

```go
func (h *CompletionHandler) ShouldSkipRegistration(cmd types.Command) bool
```

ShouldSkipRegistration 判断是否应该跳过注册

**参数:**
  - cmd: 要检查的命令

**返回值:**
  - bool: 如果标志已存在返回true, 否则返回false

**功能说明:**
  - 检查补全标志的长名称是否已存在
  - 补全标志没有短名称, 只检查长名称
  - 避免重复注册标志

**实现细节:**
  - 使用 cmd.GetFlag() 检查标志是否存在
  - 如果标志已存在, 返回true跳过注册
  - 否则返回false继续注册

---

### type HelpHandler struct

```go
type HelpHandler struct{}
```

HelpHandler 帮助标志处理器

HelpHandler 负责处理帮助标志 (-h/--help) 。 当用户指定帮助标志时, 会打印命令的帮助信息并退出程序。

#### func (h *HelpHandler) Handle(cmd types.Command) error

```go
func (h *HelpHandler) Handle(cmd types.Command) error
```

Handle 处理帮助标志

**参数:**
  - cmd: 要处理的命令

**返回值:**
  - error: 处理失败时返回错误

**功能说明: **
  - 打印命令的帮助信息
  - 使用状态码0退出程序

#### func (h *HelpHandler) ShouldRegister(cmd types.Command) bool

```go
func (h *HelpHandler) ShouldRegister(cmd types.Command) bool
```

ShouldRegister 判断是否应该注册帮助标志

**参数:**
  - cmd: 要检查的命令

**返回值:**
  - bool: 总是返回true

**功能说明: **
  - 帮助标志总是注册, 因为所有命令都应该支持帮助功能

#### func (h *HelpHandler) Type() types.BuiltinFlagType

```go
func (h *HelpHandler) Type() types.BuiltinFlagType
```

Type 返回标志类型

**返回值:**
  - types.BuiltinFlagType: HelpFlag

#### func (h *HelpHandler) ShouldSkipRegistration(cmd types.Command) bool

```go
func (h *HelpHandler) ShouldSkipRegistration(cmd types.Command) bool
```

ShouldSkipRegistration 判断是否应该跳过注册

**参数:**
  - cmd: 要检查的命令

**返回值:**
  - bool: 如果标志已存在返回true, 否则返回false

**功能说明:**
  - 检查帮助标志的长名称和短名称是否已存在
  - 避免重复注册标志
  - 支持重复解析场景

**实现细节:**
  - 使用 cmd.GetFlag() 检查长名称和短名称
  - 如果任一名称已存在, 返回true跳过注册
  - 否则返回false继续注册

---

### type VersionHandler struct

```go
type VersionHandler struct{}
```

VersionHandler 版本标志处理器

VersionHandler 负责处理版本标志 (-v/--version) 。 当用户指定版本标志时, 会打印命令的版本信息并退出程序。

#### func (h *VersionHandler) Handle(cmd types.Command) error

```go
func (h *VersionHandler) Handle(cmd types.Command) error
```

Handle 处理版本标志

**参数:**
  - cmd: 要处理的命令

**返回值:**
  - error: 处理失败时返回错误

**功能说明: **
  - 打印命令的版本信息
  - 使用状态码0退出程序

#### func (h *VersionHandler) ShouldRegister(cmd types.Command) bool

```go
func (h *VersionHandler) ShouldRegister(cmd types.Command) bool
```

ShouldRegister 判断是否应该注册版本标志

**参数:**
  - cmd: 要检查的命令

**返回值:**
  - bool: 如果设置了版本信息返回true, 否则返回false

**功能说明: **
  - 只有在命令设置了版本信息时才注册版本标志
  - 版本标志只在根命令中注册

#### func (h *VersionHandler) Type() types.BuiltinFlagType

```go
func (h *VersionHandler) Type() types.BuiltinFlagType
```

Type 返回标志类型

**返回值:**
  - types.BuiltinFlagType: VersionFlag

#### func (h *VersionHandler) ShouldSkipRegistration(cmd types.Command) bool

```go
func (h *VersionHandler) ShouldSkipRegistration(cmd types.Command) bool
```

ShouldSkipRegistration 判断是否应该跳过注册

**参数:**
  - cmd: 要检查的命令

**返回值:**
  - bool: 如果标志已存在返回true, 否则返回false

**功能说明:**
  - 检查版本标志的长名称和短名称是否已存在
  - 避免重复注册标志
  - 支持重复解析场景

**实现细节:**
  - 使用 cmd.GetFlag() 检查长名称和短名称
  - 如果任一名称已存在, 返回true跳过注册
  - 否则返回false继续注册

---