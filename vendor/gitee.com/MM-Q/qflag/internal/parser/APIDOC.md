# Package parser 

```go
import "gitee.com/MM-Q/qflag/internal/parser"
```

Package parser 提供命令行参数解析功能

parser 包实现了命令行参数解析的核心功能, 包括: 
  - 基于flag包的参数解析
  - 环境变量绑定
  - 子命令解析和路由
  - 标志验证 (快速失败模式) 
  - 内置标志自动处理
  - 智能纠错提示（子命令和标志建议）

---

## 智能纠错功能说明

### 子命令纠错

当用户输入的子命令不存在时, 解析器会自动推荐相似的子命令。

**触发条件:**
1. 命令有子命令
2. 第一个剩余参数不以"-"开头
3. 参数不匹配任何子命令名

**示例:**
```bash
$ myapp cnfig
myapp: 'cnfig' is not a valid command. See 'myapp --help'.

The most similar commands are
        config
```

### 标志纠错

当用户输入的标志不存在时, 解析器会自动推荐相似的标志。

**触发条件:**
1. 参数以"-"开头
2. 不在已注册标志列表中

**边界处理:**
- 遇到"--"停止扫描, 后面视为位置参数
- 遇到子命令名停止扫描, 后续标志由子命令处理

**示例:**
```bash
$ myapp config --verb
myapp: unknown flag: '--verb'

The most similar flags are
        --verbose
```

### 错误类型

解析器使用以下错误类型返回智能纠错信息:

- `*types.UnknownSubcommandError`: 未知子命令错误, 包含建议列表
- `*types.UnknownFlagError`: 未知标志错误, 包含建议列表

这些错误类型实现了 `error` 接口, 可以直接打印或格式化输出。

---

## FUNCTIONS

### func NewDefaultParser(errorHandling types.ErrorHandling) types.Parser

```go
func NewDefaultParser(errorHandling types.ErrorHandling) types.Parser
```

NewDefaultParser 创建默认解析器实例

**参数:**
  - errorHandling: 错误处理策略, 决定解析错误时的行为

**返回值:**
  - types.Parser: 解析器接口实例

---

## TYPES

### type DefaultParser struct

```go
type DefaultParser struct {
    // Has unexported fields.
}
```

DefaultParser 默认解析器实现

DefaultParser 是types.Parser接口的默认实现, 基于Go标准库的flag包。 它负责解析命令行参数、处理环境变量和路由子命令。

**特性:**
  - 支持所有标准标志类型
  - 支持环境变量绑定
  - 支持子命令解析和路由
  - 支持标志验证 (快速失败模式)
  - 支持互斥组验证
  - 支持必需组验证
  - 支持标志依赖关系验证 (互斥依赖和必需依赖)
  - 支持内置标志自动处理
  - 支持智能纠错提示

#### func (p *DefaultParser) Parse(cmd types.Command, args []string) error

```go
func (p *DefaultParser) Parse(cmd types.Command, args []string) error
```

Parse 解析命令行参数并处理子命令

**参数:**
  - cmd: 要解析的命令
  - args: 命令行参数列表

**返回值:**
  - error: 如果解析失败返回错误

**注意事项:**
  - 首先调用ParseOnly解析参数
  - 检查剩余参数是否为子命令
  - 如果是子命令, 递归解析子命令
  - 不执行子命令的运行函数
  - **智能纠错**: 如果子命令不存在且参数不以"-"开头, 返回带建议的错误

#### func (p *DefaultParser) ParseAndRoute(cmd types.Command, args []string) error

```go
func (p *DefaultParser) ParseAndRoute(cmd types.Command, args []string) error
```

ParseAndRoute 解析命令行参数、处理子命令并执行

**参数:**
  - cmd: 要解析的命令
  - args: 命令行参数列表

**返回值:**
  - error: 如果解析或执行失败返回错误

**注意事项:**
  - 首先调用ParseOnly解析参数
  - 检查剩余参数是否为子命令
  - 如果是子命令, 递归解析并执行子命令
  - 如果不是子命令, 执行当前命令的运行函数
  - 如果命令没有设置运行函数, 返回错误
  - **智能纠错**: 如果子命令不存在且参数不以"-"开头, 返回带建议的错误

#### func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error

```go
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error
```

ParseOnly 仅解析命令行参数, 不执行子命令路由

**参数:**
  - cmd: 要解析的命令
  - args: 命令行参数列表

**返回值:**
  - error: 如果解析失败返回错误

**注意事项:**
  - 重置所有标志到默认状态（避免重复解析时的遗留值）
  - 注册内置标志
  - 创建新的FlagSet实例进行解析
  - 注册命令的所有标志到FlagSet
  - 先解析命令行参数
  - 再加载环境变量 (仅在标志未被命令行参数设置时)
  - 验证互斥组规则 (如果配置了互斥组)
  - 验证必需组规则 (如果配置了必需组)
  - 验证标志依赖关系 (如果配置了标志依赖)
  - 处理内置标志
  - 不处理子命令路由
  - 使用defer确保命令状态和参数在函数返回时被设置
  - **智能纠错**: 在调用标准库解析前预扫描未知标志, 返回带建议的错误

---

### type SuggestionFinder struct

```go
type SuggestionFinder struct {
    maxSuggestions int // 最大建议数量
}
```

SuggestionFinder 智能纠错查找器

封装建议查找逻辑, 用于在子命令或标志输入错误时提供相似建议。

#### func NewSuggestionFinder(maxSuggestions int) *SuggestionFinder

```go
func NewSuggestionFinder(maxSuggestions int) *SuggestionFinder
```

NewSuggestionFinder 创建查找器

**参数:**
  - maxSuggestions: 最大建议数量

**返回值:**
  - *SuggestionFinder: 查找器实例

#### func (f *SuggestionFinder) FindForSubcommand(input string, cmd types.Command) []string

```go
func (f *SuggestionFinder) FindForSubcommand(input string, cmd types.Command) []string
```

FindForSubcommand 查找子命令建议

根据输入字符串, 在命令的所有子命令中查找相似的命令名(包括长短名称)。
隐藏命令(如内置补全命令)会被过滤掉。

**参数:**
  - input: 用户输入的错误子命令
  - cmd: 当前命令

**返回值:**
  - []string: 相似子命令列表(最多maxSuggestions个)

#### func (f *SuggestionFinder) FindForFlag(input string, cmd types.Command) []string

```go
func (f *SuggestionFinder) FindForFlag(input string, cmd types.Command) []string
```

FindForFlag 查找标志建议

根据输入字符串, 在命令的所有标志中查找相似的标志名(包括长短名称)。

**参数:**
  - input: 用户输入的错误标志
  - cmd: 当前命令

**返回值:**
  - []string: 相似标志列表(最多maxSuggestions个)
