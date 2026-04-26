# term

```go
import "gitee.com/MM-Q/go-kit/term"
```

## 基础菜单功能

本文件提供基础菜单功能，适用于简单的数字选项场景。 依赖 menu.go 中的类型和错误定义。

主要功能:
  - 简单的数字选项 (1, 2, 3...)
  - 自动编号选项
  - 支持自定义样式
  - 支持循环显示
  - 严格的输入验证 (使用 strconv.Atoi)

### 使用示例:

#### 示例1: 基本使用

```go
options := []string{
    "查看列表",
    "添加项目",
    "删除项目",
    "退出程序",
}

choice, err := term.ShowBasicMenuLine("主菜单", options, "请选择 (1-4): ")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("你选择了: %s\n", options[choice])
```

#### 示例2: 循环菜单

```go
options := []string{
    "查看列表",
    "添加项目",
    "删除项目",
    "退出程序",
}

err := term.ShowBasicMenuLoopLine("主菜单", options, "请选择 (1-4): ", func(index int) bool {
    switch index {
    case 0:
        fmt.Println("查看列表...")
        showList()
    case 1:
        fmt.Println("添加项目...")
        addItem()
    case 2:
        fmt.Println("删除项目...")
        deleteItem()
    case 3:
        fmt.Println("退出程序")
        return false  // 退出循环
    }
    return true  // 继续循环
})

if err != nil {
    log.Fatal(err)
}
```

#### 示例3: 自定义样式

```go
options := []string{
    "查看列表",
    "添加项目",
}

style := &term.MenuStyle{
    Prefix:    ">> ",
    Separator: " | ",
    Indent:    2,
    ShowTitle: true,
}

menuText := term.RenderBasicMenu("主菜单", options, style)
fmt.Println(menuText)
```

#### 示例4: 嵌套菜单

```go
mainOptions := []string{
    "用户管理",
    "系统设置",
    "退出",
}

err := term.ShowBasicMenuLoopLine("主菜单", mainOptions, "请选择: ", func(index int) bool {
    switch index {
    case 0:
        showUserMenu()
    case 1:
        showSettingsMenu()
    case 2:
        return false  // 退出
    }
    return true
})
```

#### 示例5: 错误处理

```go
options := []string{
    "查看列表",
    "添加项目",
}

choice, err := term.ShowBasicMenuLine("主菜单", options, "请选择: ")
if err != nil {
    if menuErr, ok := err.(*term.MenuError); ok {
        fmt.Printf("菜单错误: %s\n", menuErr.Message)
    } else {
        fmt.Printf("未知错误: %v\n", err)
    }
    return
}
fmt.Printf("你选择了: %s\n", options[choice])
```

## 结构化菜单功能

本文件提供结构化菜单功能，适用于需要自定义菜单键、退出选项、默认值等复杂场景。

主要功能:
  - 支持自定义菜单键 (如 "1", "a", "add")
  - 支持默认值
  - 支持退出选项
  - 支持自定义样式
  - 支持循环显示

### 使用示例:

#### 示例1: 基本使用

```go
menu := &term.Menu{
    Title: "主菜单",
    Items: []term.MenuItem{
        {Key: "1", Value: "查看列表"},
        {Key: "2", Value: "添加项目"},
        {Key: "3", Value: "删除项目"},
    },
    Prompt:   "请选择 (1-3): ",
    Default:  "1",
    AllowExit: true,
    ExitKey:  "q",
    ExitText: "退出",
}

choice, err := term.ShowMenuLine(menu)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("你选择了: %s\n", choice)
```

#### 示例2: 字母键菜单

```go
menu := &term.Menu{
    Title: "操作菜单",
    Items: []term.MenuItem{
        {Key: "a", Value: "添加项目"},
        {Key: "d", Value: "删除项目"},
        {Key: "v", Value: "查看列表"},
    },
    Prompt:   "请选择: ",
    AllowExit: true,
    ExitKey:  "q",
    ExitText: "退出",
}

choice, _ := term.ShowMenuLine(menu)
fmt.Printf("你选择了: %s\n", choice)
```

#### 示例3: 循环菜单

```go
menu := &term.Menu{
    Title: "主菜单",
    Items: []term.MenuItem{
        {Key: "1", Value: "查看列表"},
        {Key: "2", Value: "添加项目"},
    },
    Prompt:   "请选择: ",
    AllowExit: true,
    ExitKey:  "q",
    ExitText: "退出",
}

err := term.ShowMenuLoopLine(menu, func(key string) bool {
    switch key {
    case "1":
        fmt.Println("查看列表...")
    case "2":
        fmt.Println("添加项目...")
    case "q":
        return false  // 退出循环
    }
    return true  // 继续循环
})
```

#### 示例4: 自定义样式

```go
menu := &term.Menu{
    Title: "主菜单",
    Items: []term.MenuItem{
        {Key: "1", Value: "查看列表"},
    },
    Prompt: "请选择: ",
}

style := &term.MenuStyle{
    Prefix:    ">> ",
    Separator: " | ",
    Indent:    2,
    ShowTitle: true,
}

menuText := term.RenderMenu(menu, style)
fmt.Println(menuText)
```

#### 示例5: 嵌套菜单

```go
mainMenu := &term.Menu{
    Title: "主菜单",
    Items: []term.MenuItem{
        {Key: "1", Value: "用户管理"},
        {Key: "2", Value: "系统设置"},
    },
    Prompt:   "请选择: ",
    AllowExit: true,
    ExitKey:  "q",
    ExitText: "退出",
}

term.ShowMenuLoopLine(mainMenu, func(key string) bool {
    switch key {
    case "1":
        showUserMenu()
    case "2":
        showSettingsMenu()
    case "q":
        return false
    }
    return true
})
```

## Variables

```go
var (
    ErrEmptyMenu      = &MenuError{Code: "EMPTY_MENU", Message: "菜单不能为空"}
    ErrInvalidChoice  = &MenuError{Code: "INVALID_CHOICE", Message: "无效的选择"}
    ErrDuplicateKey   = &MenuError{Code: "DUPLICATE_KEY", Message: "菜单项键重复"}
    ErrInvalidDefault = &MenuError{Code: "INVALID_DEFAULT", Message: "默认值无效"}
)
```

预定义错误

## Functions

### func Confirm

```go
func Confirm(reader io.Reader, prompt string, defVal bool) bool
```

Confirm 确认框，支持默认值

**参数:**
  - reader: 输入源（通常为 os.Stdin）
  - prompt: 确认提示信息
  - defVal: 默认值 (true = 回车默认Yes)

**返回:**
  - bool: 用户确认结果

### func ConfirmLine

```go
func ConfirmLine(prompt string, defVal bool) bool
```

ConfirmLine 便捷函数，使用 os.Stdin 进行确认

**参数:**
  - prompt: 确认提示信息
  - defVal: 默认值 (true = 回车默认Yes)

**返回:**
  - bool: 用户确认结果

### func Read

```go
func Read(reader io.Reader, prompt string) string
```

Read 从指定的 reader 读取一行输入

**参数:**
  - reader: 输入源（通常为 os.Stdin）
  - prompt: 输入提示信息

**返回:**
  - string: 输入的字符串（已去首尾空格）

### func ReadFloat

```go
func ReadFloat(reader io.Reader, prompt string) (float64, error)
```

ReadFloat 从指定的 reader 读取浮点数输入

**参数:**
  - reader: 输入源（通常为 os.Stdin）
  - prompt: 输入提示信息

**返回:**
  - float64: 输入的浮点数
  - error: 输入为空或格式错误时返回错误

### func ReadFloatLine

```go
func ReadFloatLine(prompt string) (float64, error)
```

ReadFloatLine 便捷函数，使用 os.Stdin 读取浮点数输入

**参数:**
  - prompt: 输入提示信息

**返回:**
  - float64: 输入的浮点数
  - error: 输入为空或格式错误时返回错误

### func ReadFloatLineWithDef

```go
func ReadFloatLineWithDef(prompt string, def float64) (float64, error)
```

ReadFloatLineWithDef 便捷函数，使用 os.Stdin 读取带默认值的浮点数输入

**参数:**
  - prompt: 输入提示信息
  - def: 默认值

**返回:**
  - float64: 输入的浮点数（空输入时返回默认值）
  - error: 格式错误时返回错误

### func ReadFloatWithDef

```go
func ReadFloatWithDef(reader io.Reader, prompt string, def float64) (float64, error)
```

ReadFloatWithDef 从指定的 reader 读取浮点数输入，带默认值

**参数:**
  - reader: 输入源（通常为 os.Stdin）
  - prompt: 输入提示信息
  - def: 默认值

**返回:**
  - float64: 输入的浮点数（空输入时返回默认值）
  - error: 格式错误时返回错误

### func ReadInt

```go
func ReadInt(reader io.Reader, prompt string) (int, error)
```

ReadInt 从指定的 reader 读取整数输入

**参数:**
  - reader: 输入源（通常为 os.Stdin）
  - prompt: 输入提示信息

**返回:**
  - int: 输入的整数
  - error: 输入为空或格式错误时返回错误

### func ReadIntLine

```go
func ReadIntLine(prompt string) (int, error)
```

ReadIntLine 便捷函数，使用 os.Stdin 读取整数输入

**参数:**
  - prompt: 输入提示信息

**返回:**
  - int: 输入的整数
  - error: 输入为空或格式错误时返回错误

### func ReadIntLineWithDef

```go
func ReadIntLineWithDef(prompt string, def int) (int, error)
```

ReadIntLineWithDef 便捷函数，使用 os.Stdin 读取带默认值的整数输入

**参数:**
  - prompt: 输入提示信息
  - def: 默认值

**返回:**
  - int: 输入的整数（空输入时返回默认值）
  - error: 格式错误时返回错误

### func ReadIntWithDef

```go
func ReadIntWithDef(reader io.Reader, prompt string, def int) (int, error)
```

ReadIntWithDef 从指定的 reader 读取整数输入，带默认值

**参数:**
  - reader: 输入源（通常为 os.Stdin）
  - prompt: 输入提示信息
  - def: 默认值

**返回:**
  - int: 输入的整数（空输入时返回默认值）
  - error: 格式错误时返回错误

### func ReadLine

```go
func ReadLine(prompt string) string
```

ReadLine 便捷函数，使用 os.Stdin 读取一行输入

**参数:**
  - prompt: 输入提示信息

**返回:**
  - string: 输入的字符串（已去首尾空格）

### func ReadLineWithDef

```go
func ReadLineWithDef(prompt, def string) string
```

ReadLineWithDef 便捷函数，使用 os.Stdin 读取带默认值的输入

**参数:**
  - prompt: 输入提示信息
  - def: 默认值

**返回:**
  - string: 输入的字符串（已去首尾空格），空输入时返回默认值

### func ReadPassword

```go
func ReadPassword(prompt string) (string, error)
```

ReadPassword 读取密码输入，不显示回显

**参数:**
  - prompt: 输入提示信息

**返回:**
  - string: 输入的密码（已去首尾空格）
  - error: 错误信息

### func ReadWithDef

```go
func ReadWithDef(reader io.Reader, prompt, def string) string
```

ReadWithDef 带默认值的输入

**参数:**
  - reader: 输入源（通常为 os.Stdin）
  - prompt: 输入提示信息
  - def: 默认值

**返回:**
  - string: 输入的字符串（已去首尾空格），空输入时返回默认值

### func RenderBasicMenu

```go
func RenderBasicMenu(title string, options []string, style *MenuStyle) string
```

RenderBasicMenu 渲染基础菜单为字符串

**参数:**
  - title: 菜单标题
  - options: 选项列表
  - style: 菜单样式 (可选, 使用默认样式)

**返回:**
  - string: 渲染后的菜单字符串

### func RenderMenu

```go
func RenderMenu(menu *Menu, style *MenuStyle) string
```

RenderMenu 渲染菜单为字符串

**参数:**
  - menu: 菜单配置
  - style: 菜单样式 (可选, 使用默认样式)

**返回:**
  - string: 渲染后的菜单字符串

### func ShowBasicMenu

```go
func ShowBasicMenu(title string, options []string, reader io.Reader, prompt string) (int, error)
```

ShowBasicMenu 显示基础菜单并获取用户选择

**参数:**
  - title: 菜单标题
  - options: 选项列表
  - reader: 输入源
  - prompt: 提示信息

**返回:**
  - int: 用户选择的索引 (从 0 开始)
  - error: 错误信息

### func ShowBasicMenuLine

```go
func ShowBasicMenuLine(title string, options []string, prompt string) (int, error)
```

ShowBasicMenuLine 便捷函数, 使用 os.Stdin 显示菜单

**参数:**
  - title: 菜单标题
  - options: 选项列表
  - prompt: 提示信息

**返回:**
  - int: 用户选择的索引 (从 0 开始)
  - error: 错误信息

### func ShowBasicMenuLoop

```go
func ShowBasicMenuLoop(title string, options []string, reader io.Reader, prompt string, handler func(index int) bool) error
```

ShowBasicMenuLoop 循环显示基础菜单, 直到用户选择退出

**参数:**
  - title: 菜单标题
  - options: 选项列表
  - reader: 输入源
  - prompt: 提示信息
  - handler: 处理函数, 返回 true 继续循环, false 退出

**返回:**
  - error: 错误信息

### func ShowBasicMenuLoopLine

```go
func ShowBasicMenuLoopLine(title string, options []string, prompt string, handler func(index int) bool) error
```

ShowBasicMenuLoopLine 便捷函数, 使用 os.Stdin 循环显示菜单

**参数:**
  - title: 菜单标题
  - options: 选项列表
  - prompt: 提示信息
  - handler: 处理函数, 返回 true 继续循环, false 退出

**返回:**
  - error: 错误信息

### func ShowMenu

```go
func ShowMenu(menu *Menu, reader io.Reader) (string, error)
```

ShowMenu 显示结构化菜单并获取用户选择

**参数:**
  - menu: 菜单配置
  - reader: 输入源

**返回:**
  - string: 用户选择的键
  - error: 错误信息

### func ShowMenuLine

```go
func ShowMenuLine(menu *Menu) (string, error)
```

ShowMenuLine 便捷函数, 使用 os.Stdin 显示菜单

**参数:**
  - menu: 菜单配置

**返回:**
  - string: 用户选择的键
  - error: 错误信息

### func ShowMenuLoop

```go
func ShowMenuLoop(menu *Menu, reader io.Reader, handler func(key string) bool) error
```

ShowMenuLoop 循环显示菜单, 直到用户选择退出

**参数:**
  - menu: 菜单配置
  - reader: 输入源
  - handler: 处理函数, 返回 true 继续循环, false 退出

**返回:**
  - error: 错误信息

### func ShowMenuLoopLine

```go
func ShowMenuLoopLine(menu *Menu, handler func(key string) bool) error
```

ShowMenuLoopLine 便捷函数, 使用 os.Stdin 循环显示菜单

**参数:**
  - menu: 菜单配置
  - handler: 处理函数, 返回 true 继续循环, false 退出

**返回:**
  - error: 错误信息

### func ValidateMenu

```go
func ValidateMenu(menu *Menu) error
```

ValidateMenu 验证菜单配置

**参数:**
  - menu: 菜单配置

**返回:**
  - error: 验证错误

## Types

### type Menu

```go
type Menu struct {
    Title     string     // 菜单标题, 如 "主菜单"
    Items     []MenuItem // 菜单项列表 (如 ["1. 查看", "2. 添加"])
    Prompt    string     // 提示信息, 如 "请输入选择"
    Default   string     // 默认值 (可选), 如 "1"
    AllowExit bool       // 是否允许退出, 如 true
    ExitKey   string     // 退出键 (如 "q", "0"), 如 "q"
    ExitText  string     // 退出选项描述 (如 "退出"), 如 "退出"
}
```

Menu 菜单配置 (结构化菜单)

### type MenuError

```go
type MenuError struct {
    Code    string // 错误代码
    Message string // 错误消息
    Detail  string // 详细信息
}
```

MenuError 菜单错误类型

#### func NewMenuError

```go
func NewMenuError(code, message, detail string) *MenuError
```

NewMenuError 创建菜单错误

#### func (*MenuError) Error

```go
func (e *MenuError) Error() string
```

Error 实现 error 接口

### type MenuItem

```go
type MenuItem struct {
    Key   string // 选项键 (如 "1", "a", "add")
    Value string // 选项描述 (如 "查看列表")
}
```

MenuItem 菜单项

### type MenuStyle

```go
type MenuStyle struct {
    Prefix    string // 选项前缀 (如 "1. ", "[a] ")
    Separator string // 分隔符 (如 " - ", ": ")
    Indent    int    // 缩进空格数
    ShowTitle bool   // 是否显示标题
}
```

MenuStyle 菜单样式

#### func GetDefaultMenuStyle

```go
func GetDefaultMenuStyle() *MenuStyle
```

GetDefaultMenuStyle 获取默认菜单样式

**返回:**
  - *MenuStyle: 默认菜单样式

## 终端工具函数

### func IsStdinPipe

```go
func IsStdinPipe() bool
```

IsStdinPipe 检测标准输入是否为管道或文件重定向 用于区分程序是从终端直接接收输入，还是通过管道(\|)或文件重定向(<)接收输入，常用于判断是否需要显示交互式提示

**注意:**
  - 如果获取 stdin 状态失败，该函数返回 false（按终端处理）
  - 如需获取错误详情，请使用 IsStdinPipeWithError

**返回:**
  - `bool`: true 表示 stdin 是管道或文件重定向; false 表示是终端输入或出错

**示例:**

```go
if term.IsStdinPipe() {
    // 从管道读取数据
    data, _ := io.ReadAll(os.Stdin)
} else {
    // 交互式终端，显示提示
    fmt.Print("请输入: ")
}
```

### func IsStdinPipeWithError

```go
func IsStdinPipeWithError() (bool, error)
```

IsStdinPipeWithError 检测标准输入是否为管道或文件重定向（高级版本） 与 IsStdinPipe 功能相同，但会返回详细的错误信息，适用于需要精确错误处理的场景

**返回:**
  - `bool`: true 表示 stdin 是管道或文件重定向; false 表示是终端输入
  - `error`: 获取 stdin 状态时发生的错误，nil 表示成功

**示例:**

```go
isPipe, err := term.IsStdinPipeWithError()
if err != nil {
    log.Printf("无法检测输入类型: %v", err)
    return
}
if isPipe {
    // 处理管道输入
} else {
    // 处理交互式输入
}
```

### func GetSafeTerminalWidth

```go
func GetSafeTerminalWidth() int
```

GetSafeTerminalWidth 获取当前终端的宽度（字符列数） 该函数会尝试多种方式获取终端宽度，并确保返回的值在合理范围内:
1. 首先检查 COLUMNS 环境变量
2. 其次尝试通过系统调用获取终端尺寸
3. 如果都失败，返回默认值

**返回:**
  - `int`: 终端宽度（字符列数），范围在 [40, 1200] 之间

**示例:**

```go
width := term.GetSafeTerminalWidth()
fmt.Printf("当前终端宽度: %d 列\n", width)
```