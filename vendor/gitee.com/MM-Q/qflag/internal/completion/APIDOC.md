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