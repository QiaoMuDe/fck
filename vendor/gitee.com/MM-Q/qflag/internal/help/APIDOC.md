# Package help

**Import Path:** `gitee.com/MM-Q/qflag/internal/help`

Package help 帮助信息生成器。本包实现了命令行帮助信息的自动生成功能，包括标志列表、用法说明、子命令信息等帮助内容的格式化和输出。

## 功能模块

- **帮助信息生成器** - 实现了命令行帮助信息的自动生成功能
- **帮助信息排序和组织** - 实现了帮助信息的排序和组织功能，包括标志排序、子命令排序等
- **测试辅助工具** - 提供了帮助信息模块的测试辅助函数和工具
- **帮助信息输出和格式化** - 实现了帮助信息的输出和格式化功能，支持多种输出格式和样式

## 目录

- [常量](#常量)
- [变量](#变量)
- [函数](#函数)
- [类型](#类型)

## 常量

```go
const (
    // 默认最大宽度，当计算失败时使用
    DefaultMaxWidth = 30

    // 描述信息与选项之间的间距
    DescriptionPadding = 5

    // 子命令名称分隔符长度 (", " 的长度)
    SubCmdSeparatorLen = 2

    // 子命令对齐额外空格数
    SubCmdAlignSpaces = 5

    // 最小填充空格数
    MinPadding = 1
)
```

帮助信息格式化常量，用于控制帮助信息的布局和格式。

## 变量

### ChineseTemplate

```go
var ChineseTemplate = HelpTemplate{
    CmdName:              "名称: %s\n\n",
    UsagePrefix:          "用法: ",
    UsageSubCmd:          " [子命令]",
    UsageInfoWithOptions: " [选项]\n\n",
    UsageGlobalOptions:   " [全局选项]",
    CmdNameWithShort:     "名称: %s, %s\n\n",
    CmdDescription:       "描述: %s\n\n",
    OptionsHeader:        "选项:\n",
    Option1:              "  -%s, --%s %s",
    Option2:              "  --%s %s",
    Option3:              "  -%s %s",
    OptionDefault:        "%s%*s%s (默认值: %s)\n",
    SubCmdsHeader:        "\n子命令:\n",
    SubCmd:               "  %s\t%s\n",
    SubCmdWithShort:      "  %s, %s\t%s\n",
    NotesHeader:          "\n注意事项:\n",
    NoteItem:             "  %d、%s\n",
    DefaultNote:          "当长选项和短选项同时使用时，最后指定的选项将优先生效。",
    ExamplesHeader:       "\n示例:\n",
    ExampleItem:          "  %d、%s\n     %s\n",
}
```

中文模板实例，提供中文环境下的帮助信息格式化模板。

### EnglishTemplate

```go
var EnglishTemplate = HelpTemplate{
    CmdName:              "Name: %s\n\n",
    UsagePrefix:          "Usage: ",
    UsageSubCmd:          " [subcmd]",
    UsageInfoWithOptions: " [options]\n\n",
    UsageGlobalOptions:   " [global options]",
    CmdNameWithShort:     "Name: %s, %s\n\n",
    CmdDescription:       "Desc: %s\n\n",
    OptionsHeader:        "Options:\n",
    Option1:              "  -%s, --%s %s",
    Option2:              "  --%s %s",
    Option3:              "  -%s %s",
    OptionDefault:        "%s%*s%s (default: %s)\n",
    SubCmdsHeader:        "\nSubCmds:\n",
    SubCmd:               "  %s\t%s\n",
    SubCmdWithShort:      "  %s, %s\t%s\n",
    NotesHeader:          "\nNotes:\n",
    NoteItem:             "  %d. %s\n",
    DefaultNote:          "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.",
    ExamplesHeader:       "\nExamples:\n",
    ExampleItem:          "  %d. %s\n     %s\n",
}
```

英文模板实例，提供英文环境下的帮助信息格式化模板。

## 函数

### GenerateHelp

```go
func GenerateHelp(ctx *types.CmdContext) string
```

GenerateHelp 生成帮助信息。纯函数设计，不依赖任何结构体状态。

**参数:**
- `ctx`: 命令上下文，包含生成帮助信息所需的所有数据

**返回值:**
- `string`: 格式化后的帮助信息字符串

**特点:**
- 纯函数设计，无副作用
- 支持中英文模板
- 自动格式化选项和子命令
- 支持自定义模板

## 类型

### HelpTemplate

```go
type HelpTemplate struct {
    CmdName              string // 命令名称模板
    CmdNameWithShort     string // 命令名称带短名称模板
    CmdDescription       string // 命令描述模板
    UsagePrefix          string // 用法说明前缀模板
    UsageSubCmd          string // 用法说明子命令模板
    UsageInfoWithOptions string // 带选项的用法说明信息模板
    UsageGlobalOptions   string // 全局选项部分
    OptionsHeader        string // 选项头部模板
    Option1              string // 选项模板(带短选项)
    Option2              string // 选项模板(无短选项)
    Option3              string // 选项模板(无长选项)
    OptionDefault        string // 选项模板的默认值
    SubCmdsHeader        string // 子命令头部模板
    SubCmd               string // 子命令模板
    SubCmdWithShort      string // 子命令带短名称模板
    NotesHeader          string // 注意事项头部模板
    NoteItem             string // 注意事项项模板
    DefaultNote          string // 默认注意事项
    ExamplesHeader       string // 示例信息头部模板
    ExampleItem          string // 示例信息项模板
}
```

HelpTemplate 帮助信息模板结构体，定义了帮助信息各个部分的格式化模板。

**用途:**
- 定义帮助信息的显示格式
- 支持国际化（中英文模板）
- 提供灵活的自定义格式化选项

### NamedItem

```go
type NamedItem interface {
    GetLongName() string
    GetShortName() string
}
```

NamedItem 表示具有长名称和短名称的项目接口。

**方法:**

#### GetLongName

```go
GetLongName() string
```

获取项目的长名称。

**返回值:**
- `string`: 项目的长名称

#### GetShortName

```go
GetShortName() string
```

获取项目的短名称。

**返回值:**
- `string`: 项目的短名称

**实现者:**
- 命令行选项（flags）
- 子命令（subcommands）
- 其他具有长短名称的命令行元素

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag/internal/help"
    "gitee.com/MM-Q/qflag/internal/types"
)

func main() {
    // 创建命令上下文
    ctx := &types.CmdContext{
        // ... 设置命令信息
    }
    
    // 生成帮助信息
    helpText := help.GenerateHelp(ctx)
    fmt.Println(helpText)
}
```

### 自定义模板

```go
// 创建自定义模板
customTemplate := help.HelpTemplate{
    CmdName:       "命令: %s\n\n",
    UsagePrefix:   "使用方法: ",
    OptionsHeader: "可用选项:\n",
    // ... 其他模板字段
}

// 在生成帮助信息时使用自定义模板
// (需要在 GenerateHelp 函数中支持模板参数)
```

## 设计特点

1. **纯函数设计** - `GenerateHelp` 函数不依赖任何全局状态
2. **模板化** - 支持中英文模板，易于国际化
3. **灵活格式化** - 通过常量控制布局和间距
4. **接口抽象** - `NamedItem` 接口提供统一的命名项目抽象

## 注意事项

- 帮助信息的生成依赖于 `types.CmdContext` 中的数据完整性
- 模板字符串使用 Go 的 `fmt` 包格式化语法
- 常量值的修改会影响所有帮助信息的显示效果
- 支持的最大宽度和填充设置确保了良好的显示效果
