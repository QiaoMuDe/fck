# Package completion

包完成功能为Bash和PowerShell环境中的命令行自动补全提供了实现。此文件实现了自动补全的核心功能，包括标志补全、子命令补全和参数值补全，提升了命令行交互中的用户体验。

## Package Overview

- **Bash Shell 自动补全实现**  
  实现了 Bash Shell 环境下的命令行自动补全功能，生成 Bash 补全脚本，支持标志和子命令的智能补全。

- **PowerShell 自动补全实现**  
  实现了 PowerShell 环境下的命令行自动补全功能，生成 PowerShell 补全脚本，支持标志和子命令的智能补全。

- **自动补全内部实现**  
  包含自动补全功能的内部实现逻辑，提供补全算法和匹配策略等核心功能的底层支持。

## Constants

### Bash 相关常量

```go
const (
    BashFlagParamItem = "{{.ProgramName}}_flag_params[%q]=%q\n"  // 标志参数项格式
    BashEnumOptions   = "{{.ProgramName}}_enum_options[%q]=%q\n" // 枚举选项格式
)
```

### 默认参数配置

```go
const (
    // DefaultFlagParamsCapacity 预估的标志参数初始容量
    // 基于常见 CLI 工具分析，大多数工具的标志数量在 100-500 之间
    DefaultFlagParamsCapacity = 256

    // NamesPerItem 每个标志/命令的名称数量(长名+短名)
    NamesPerItem = 2

    // MaxTraverseDepth 命令树遍历的最大深度限制
    // 防止循环引用导致的无限递归，一般 CLI 工具很少超过 20 层
    MaxTraverseDepth = 50
)
```

### PowerShell 相关常量

```go
const (
    // 标志参数条目(含枚举选项)
    PwshFlagParamItem = "	@{ Context = \"{{.Context}}\"; Parameter = \"{{.Parameter}}\"; ParamType = \"{{.ParamType}}\"; ValueType = \"{{.ValueType}}\"; Options = @({{.Options}}) }"
    // 命令树条目
    PwshCmdTreeItem = "	@{ Context = \"{{.Context}}\"; Options = @({{.Options}}) }"
)
```

### 补全脚本模板

```go
const (
    // 优化的 Bash 补全模板 - 集成高性能模糊匹配功能
    BashFunctionHeader = `...`
)

const (
    // PowerShell 自动补全脚本头部
    PwshFunctionHeader = `...`
)
```

## Variables

### 注意事项

```go
var (
    // CompletionNotesCN 中文版本注意事项
    CompletionNotesCN = []string{
        "Windows 环境: 需要 PowerShell 5.1 或更高版本以支持 Register-ArgumentCompleter",
        "Linux 环境: 需要 bash 4.0 或更高版本以支持关联数组特性",
        "请确保您的环境满足上述版本要求，否则自动补全功能可能无法正常工作",
    }

    // CompletionNotesEN 英文版本注意事项
    CompletionNotesEN = []string{
        "Windows environment: Requires PowerShell 5.1 or higher to support Register-ArgumentCompleter",
        "Linux environment: Requires bash 4.0 or higher to support associative array features",
        "Please ensure your environment meets the above version requirements, otherwise the auto-completion feature may not work properly",
    }
)
```

### 示例用法

#### 中文示例

```go
var CompletionExamplesCN = []types.ExampleInfo{
    {Description: "Linux 环境 临时启用", Usage: "source <(%s --generate-shell-completion bash)"},
    {Description: "Linux 环境 永久启用(添加到 ~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},

    {Description: "Windows 环境 临时启用", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
    {Description: "Windows 环境 永久启用(添加到 PowerShell 配置文件)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
}
```

#### 英文示例

```go
var CompletionExamplesEN = []types.ExampleInfo{
    {Description: "Linux environment temporary activation", Usage: "source <(%s --generate-shell-completion bash)"},
    {Description: "Linux environment permanent activation (add to ~/.bashrc)", Usage: "echo \"source <(%s --generate-shell-completion bash)\" >> ~/.bashrc"},

    {Description: "Windows environment temporary activation", Usage: "%s --generate-shell-completion powershell | Out-String | Invoke-Expression"},
    {Description: "Windows environment permanent activation (add to PowerShell profile)", Usage: "echo \"%s --generate-shell-completion powershell | Out-String | Invoke-Expression\" >> $PROFILE"},
}
```

## Functions

### GenerateShellCompletion

```go
func GenerateShellCompletion(ctx *types.CmdContext, shellType string) (string, error)
```

- **功能描述**  
  生成 shell 自动补全脚本。
- **参数**  
  - `ctx`: 命令上下文。
  - `shellType`: shell 类型 ("bash", "pwsh", "powershell")。
- **返回值**  
  - `string`: 自动补全脚本。
  - `error`: 错误信息。

## Types

### FlagParam

```go
type FlagParam struct {
    CommandPath string   // 命令路径，如 "/cmd/subcmd"
    Name        string   // 标志名称(保留原始大小写)
    Type        string   // 参数需求类型: "required"|"optional"|"none"
    ValueType   string   // 参数值类型: "path"|"string"|"number"|"enum"|"bool" 等
    EnumOptions []string // 枚举类型的可选值列表
}
```

- **功能描述**  
  表示标志参数及其需求类型和值类型。