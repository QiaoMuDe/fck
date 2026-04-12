# Package completion

```go
import "gitee.com/MM-Q/qflag/internal/completion"
```

## 包介绍

Package completion Bash Shell 自动补全实现 本文件实现了Bash Shell环境下的命令行自动补全功能, 生成Bash补全脚本,
支持标志和子命令的智能补全。

Package completion 自动补全内部实现 本文件包含了自动补全功能的内部实现逻辑, 提供补全算法、 匹配策略等核心功能的底层支持。

Package completion PowerShell 自动补全实现 本文件实现了PowerShell环境下的命令行自动补全功能,
生成PowerShell补全脚本, 支持标志和子命令的智能补全。

## CONSTANTS

### const (
	BashFlagParamItem = "{{.ProgramName}}_flag_params[%q]=%q\n"  // 标志参数项格式
	BashEnumOptions   = "{{.ProgramName}}_enum_options[%q]=%q\n" // 枚举选项格式
)

### const (
	// 标志参数条目(含枚举选项)
	PwshFlagParamItem = "	@{ Context = \"{{.Context}}\"; Parameter = \"{{.Parameter}}\"; ParamType = \"{{.ParamType}}\"; ValueType = \"{{.ValueType}}\"; Options = @({{.Options}}) }"
	// 命令树条目
	PwshCmdTreeItem = "	@{ Context = \"{{.Context}}\"; Options = @({{.Options}}) }"
)

## FUNCTIONS

### func GenAndPrint(cmd types.Command, shellType string)

```go
func GenAndPrint(cmd types.Command, shellType string)
```

GenAndPrint 生成并打印补全脚本

参数:
  - cmd: 要生成补全脚本的命令
  - shellType: Shell类型 (bash, pwsh, powershell)

### func Generate(cmd types.Command, shellType string) (string, error)

```go
func Generate(cmd types.Command, shellType string) (string, error)
```

Generate 生成补全脚本

参数:
  - cmd: 要生成补全脚本的命令
  - shellType: Shell类型 (bash, pwsh, powershell)

返回值:
  - string: 生成的补全脚本
  - error: 生成失败时返回错误

## TYPES

### type FlagParam struct

```go
type FlagParam struct {
	CommandPath string   // 命令路径, 如 "/cmd/subcmd"
	Name        string   // 标志名称(保留原始大小写)
	Type        string   // 参数需求类型: "required"|"optional"|"none"
	ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool"等
	EnumOptions []string // 枚举类型的可选值列表
}
```

FlagParam 表示标志参数及其需求类型和值类型

### type ContextResult struct

```go
type ContextResult struct {
    Context         string   // 上下文路径，如 "/server/start/"
    Command         string   // 当前命令名
    Depth           int      // 嵌套深度
    CurrentCmd      string   // 当前命令名称
    CurrentDesc     string   // 当前命令描述
    SubCommands     []string // 可用子命令列表
    Flags           []string // 可用标志列表（长短名称）
    IsFlagContext   bool     // 是否处于标志上下文
    FlagsStartIndex int      // 标志开始的位置（-1 表示无）
    ParentContext   string   // 父上下文路径
}
```

ContextResult 表示上下文计算结果，包含当前命令信息和可用选项

---

## FUNCTIONS

### func HandleDynamicComplete(root types.Command, instruction string, params []string) error

```go
func HandleDynamicComplete(root types.Command, instruction string, params []string) error
```

HandleDynamicComplete 处理 __complete 子命令的路由

参数:
  - root: 根命令实例，用于查询命令树
  - instruction: 指令名称（fuzzy, context, candidates, enum, all）
  - params: 指令参数列表

返回值:
  - error: 处理错误

### func HandleContext(root types.Command, args []string) error

```go
func HandleContext(root types.Command, args []string) error
```

HandleContext 处理 context 指令，计算并输出上下文路径

参数:
  - root: 根命令实例
  - args: 命令行参数（子命令名称列表）

返回值:
  - error: 处理错误

### func CalculateContext(root types.Command, tokens []string, cursorPos int) *ContextResult

```go
func CalculateContext(root types.Command, tokens []string, cursorPos int) *ContextResult
```

CalculateContext 计算当前上下文路径

参数:
  - root: 根命令实例
  - tokens: 命令行令牌列表（包含程序名）
  - cursorPos: 光标位置

返回值:
  - *ContextResult: 上下文计算结果

### func HandleCandidates(root types.Command, args []string) error

```go
func HandleCandidates(root types.Command, args []string) error
```

HandleCandidates 处理 candidates 指令，输出候选选项列表

参数:
  - root: 根命令实例
  - args: [context]

返回值:
  - error: 处理错误

### func GetCandidates(root types.Command, context string) ([]string, error)

```go
func GetCandidates(root types.Command, context string) ([]string, error)
```

GetCandidates 获取候选选项（供程序内部使用）

参数:
  - root: 根命令实例
  - context: 上下文路径

返回值:
  - []string: 候选选项列表
  - error: 处理错误

### func HandleEnum(root types.Command, args []string) error

```go
func HandleEnum(root types.Command, args []string) error
```

HandleEnum 处理 enum 指令，输出枚举值列表

参数:
  - root: 根命令实例
  - args: [context, flag-name]

返回值:
  - error: 处理错误

### func GetEnumValues(root types.Command, context string, flagName string) ([]string, error)

```go
func GetEnumValues(root types.Command, context string, flagName string) ([]string, error)
```

GetEnumValues 获取枚举值（供程序内部使用）

参数:
  - root: 根命令实例
  - context: 上下文路径
  - flagName: 标志名称

返回值:
  - []string: 枚举值列表
  - error: 处理错误