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
  - 支持内置标志自动处理

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
  - 处理内置标志
  - 不处理子命令路由
  - 使用defer确保命令状态和参数在函数返回时被设置