# color 包文档

```go
package color // import "gitee.com/MM-Q/color"
```

Package color 是一个 ANSI 颜色包，用于向标准输出输出彩色或 SGR 定义的文本。 API 可以通过多种方式使用，选择适合你的方式即可。

## 快速开始

使用简单且默认的辅助函数，配合预定义的前景色：

```go
// 直接打印字符串（自动追加换行符）
color.Cyan("以青色打印文本。")
color.Red("我们有红色")
color.Yellow("也有黄色！")
color.Magenta("还有很多其他颜色 ..")

// 使用格式化版本（不自动追加换行符，需手动添加）
color.Bluef("以蓝色打印 %s。\n", "文本")

// 高亮色
color.HiGreen("亮绿色。")
color.HiBlack("亮黑色就是灰色..")
color.HiWhite("闪亮的白色！")

// 格式化高亮色
color.HiGreenf("亮绿色: %s\n", "信息")
```

## 自定义颜色组合

然而，有时需要自定义颜色组合。以下是一些创建自定义颜色对象 并使用每个独立颜色对象的打印函数的示例。

```go
// 创建一个新的颜色对象
c := color.New(color.FgCyan).Add(color.Underline)
c.Println("打印带下划线的青色文本。")

// 或者直接添加到 New() 中
d := color.New(color.FgCyan, color.Bold)
d.Printf("这也打印粗体青色 %s\n", "！")

// 混合前景色和背景色，创建新的组合！
red := color.New(color.FgRed)

boldRed := red.Add(color.Bold)
boldRed.Println("这将打印粗体红色文本。")

whiteBackground := red.Add(color.BgWhite)
whiteBackground.Println("红色文本配白色背景。")

// 使用你自己的 io.Writer 输出
color.New(color.FgBlue).Fprintln(myWriter, "蓝色！")

blue := color.New(color.FgBlue)
blue.Fprint(myWriter, "这将打印蓝色文本。")
```

## 创建 PrintXxx 函数

你可以创建 PrintXxx 函数来进一步简化：

```go
// 创建自定义打印函数以方便使用
red := color.New(color.FgRed).PrintfFunc()
red("警告")
red("错误：%s", err)

// 混合多个属性
notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
notice("不要忘记这个...")
```

## 使用 FprintXxx 函数

你也可以使用 FprintXxx 函数传入你自己的 io.Writer：

```go
blue := color.New(FgBlue).FprintfFunc()
blue(myWriter, "重要通知：%s", stars)

// 混合多个属性
success := color.New(color.Bold, color.FgGreen).FprintlnFunc()
success(myWriter, "不要忘记这个...")
```

## 创建 SprintXxx 函数

或者创建 SprintXxx 函数来将字符串与其他非彩色字符串混合：

```go
yellow := New(FgYellow).SprintFunc()
red := New(FgRed).SprintFunc()

fmt.Printf("这是一个 %s，这是一个 %s。\n", yellow("警告"), red("错误"))

info := New(FgWhite, BgGreen).SprintFunc()
fmt.Printf("这个 %s 太棒了！\n", info("包"))
```

## Windows 支持

Windows 支持默认启用。所有 Print 函数都能按预期工作。 但是，仅对于 color.SprintXXX 函数，用户应该使用 fmt.FprintXXX 并将输出设置为 color.Output：

```go
fmt.Fprintf(color.Output, "Windows 支持：%s", color.SGreen("通过"))

info := New(FgWhite, BgGreen).SprintFunc()
fmt.Fprintf(color.Output, "这个 %s 太棒了！\n", info("包"))
```

## 与现有代码集成

可以与现有代码一起使用。只需使用 Set() 方法将标准输出设置为给定参数。 这样就不需要重写现有代码。

```go
// 使用便捷的标准颜色。
color.Set(color.FgYellow)

fmt.Println("现有文本现在将显示为黄色")
fmt.Printf("这个也是 %s\n", "黄色")

color.Unset() // 不要忘记取消设置

// 你可以混合参数
color.Set(color.FgMagenta, color.Bold)
defer color.Unset() // 在你的函数中使用它

fmt.Println("所有文本现在将显示为粗体洋红色。")
```

## 禁用颜色输出

可能会有需要禁用颜色输出的情况（例如将应用程序的标准输出管道传输到其他地方）。 `Color` 支持全局和单个颜色定义禁用颜色。例如，假设你有一个 CLI 应用程序 和一个 `--no-color` 布尔标志。你可以轻松禁用颜色输出：

```go
var flagNoColor = flag.Bool("no-color", false, "禁用颜色输出")

if *flagNoColor {
    color.NoColor = true // 禁用彩色输出
}
```

你也可以通过将 NO_COLOR 环境变量设置为任何值来禁用颜色。

它还支持单个颜色定义（本地）。你可以随时禁用/启用颜色输出：

```go
c := color.New(color.FgCyan)
c.Println("打印青色文本")

c.DisableColor()
c.Println("这将不打印任何颜色")

c.EnableColor()
c.Println("这又打印青色了...")
```

---

## 变量

```go
var (
    // NoColor 定义输出是否着色。它根据 stdout 的文件描述符是否指向终端
    // 动态设置为 false 或 true。如果设置了 NO_COLOR 环境变量（无论其值是什么），
    // 它也会被设置为 true。这是一个全局选项, 影响所有颜色。
    // 如需对每个颜色块进行更多控制，请单独使用 DisableColor() 方法。
    NoColor = noColorIsSet() || os.Getenv("TERM") == "dumb" || !stdoutIsTerminal()

    // Output 定义打印函数的标准输出。默认使用 stdOut()。
    Output = stdOut()

    // Error 定义打印函数的标准错误输出。默认使用 stdErr()。
    Error = stdErr()
)
```

---

## 便捷函数

包提供了大量便捷函数用于快速打印彩色文本。这些函数分为两类：

### 包级便捷函数

直接使用包名调用，无需创建 Color 对象：

```go
// 打印函数（自动追加换行符，仅接受字符串参数）
color.Red("错误信息")
color.Green("操作完成")
color.Blue("提示信息")

// 格式化打印函数（不自动追加换行符）
color.Greenf("成功: %s\n", "操作完成")

// 返回字符串函数（仅接受字符串参数）
errMsg := color.SRed("错误")
successMsg := color.SGreen("成功")

// 格式化返回字符串函数
errMsgf := color.SRedf("错误: %s", "详情")
```

**可用的颜色函数：**

| 标准色 | 高亮色 | 说明 |
|--------|--------|------|
| `Black` / `SBlack` | `HiBlack` / `SHiBlack` | 黑色/高亮黑色 |
| `Red` / `SRed` | `HiRed` / `SHiRed` | 红色/高亮红色 |
| `Green` / `SGreen` | `HiGreen` / `SHiGreen` | 绿色/高亮绿色 |
| `Yellow` / `SYellow` | `HiYellow` / `SHiYellow` | 黄色/高亮黄色 |
| `Blue` / `SBlue` | `HiBlue` / `SHiBlue` | 蓝色/高亮蓝色 |
| `Magenta` / `SMagenta` | `HiMagenta` / `SHiMagenta` | 洋红色/高亮洋红色 |
| `Cyan` / `SCyan` | `HiCyan` / `SHiCyan` | 青色/高亮青色 |
| `White` / `SWhite` | `HiWhite` / `SHiWhite` | 白色/高亮白色 |
| `Gray` / `SGray` | - | 灰色（HiBlack 的别名） |

**函数命名规则：**
- `Xxx(message)` - 打印彩色文本并自动追加换行符（仅接受字符串参数）
- `Xxxf(format, a...)` - 打印格式化彩色文本（不自动追加换行符）
- `SXxx(message)` - 返回彩色字符串（仅接受字符串参数）
- `SXxxf(format, a...)` - 返回格式化彩色字符串

> 💡 **提示**：不带 `f` 后缀的方法只接受字符串参数，会自动追加换行符；带 `f` 后缀的方法支持格式化字符串，不自动追加换行符。

### 全局实例方法

通过 `color.G()` 或 `color.GetGlobal()` 获取全局实例，支持链式配置：

```go
// 使用默认配置
color.G().Red("错误信息")
color.G().SGreen("成功")

// 自定义配置
color.G().SetBold(true).SetUnderline(true).Red("粗体下划线红色")

// 临时禁用颜色
color.G().SetNoColor(true).Blue("这行不会显示颜色")
```

**全局实例特性：**
- 支持样式配置（加粗、斜体、下划线等）
- 支持自定义输出目标
- 线程安全
- 可动态启用/禁用颜色

---

## 函数

### Unset

```go
func Unset()
```

Unset 重置所有转义属性并清除输出。 通常在 Set() 之后调用。

---

## 类型

### Attribute

```go
type Attribute int
```

Attribute 定义单个 SGR（Select Graphic Rendition）代码， 用于控制终端文本的显示属性，如颜色、样式等。

#### 基础文本样式属性

这些属性控制文本的显示样式，如加粗、斜体、下划线等。

```go
const (
    Reset        Attribute = iota // 重置所有属性为默认值
    Bold                          // 加粗文本
    Faint                         // 轻淡/暗淡文本（降低亮度）
    Italic                        // 斜体文本
    Underline                     // 下划线文本
    BlinkSlow                     // 慢速闪烁（每秒少于 150 次）
    BlinkRapid                    // 快速闪烁（每分钟 150 次以上）
    ReverseVideo                  // 反显（前景色和背景色互换）
    Concealed                     // 隐蔽/隐藏文本（与背景色相同）
    CrossedOut                    // 删除线文本
)
```

#### 重置属性

这些属性用于取消对应的文本样式。

```go
const (
    ResetBold      Attribute = iota + 22 // 重置加粗
    ResetItalic                          // 重置斜体
    ResetUnderline                       // 重置下划线
    ResetBlinking                        // 重置闪烁

    ResetReversed   // 重置反显
    ResetConcealed  // 重置蔽蔽
    ResetCrossedOut // 重置删除线
)
```

#### 前景文本颜色（标准 8 色）

这些颜色在大多数终端中都受支持。

```go
const (
    FgBlack   Attribute = iota + 30 // 黑色前景
    FgRed                           // 红色前景
    FgGreen                         // 绿色前景
    FgYellow                        // 黄色前景
    FgBlue                          // 蓝色前景
    FgMagenta                       // 洋红色/品红色前景
    FgCyan                          // 青色前景
    FgWhite                         // 白色前景
    FgGray = FgHiBlack              // 灰色前景（FgHiBlack 的别名）
)
```

#### 前景高亮文本颜色（高亮 8 色）

这些是高亮版本的标准前景色，比标准色更明亮。

```go
const (
    FgHiBlack   Attribute = iota + 90 // 高亮黑色前景（灰色）
    FgHiRed                           // 高亮红色前景
    FgHiGreen                         // 高亮绿色前景
    FgHiYellow                        // 高亮黄色前景
    FgHiBlue                          // 高亮蓝色前景
    FgHiMagenta                       // 高亮洋红色前景
    FgHiCyan                          // 高亮青色前景
    FgHiWhite                         // 高亮白色前景
)
```

#### 背景文本颜色（标准 8 色）

这些颜色用于设置文本的背景色。

```go
const (
    BgBlack   Attribute = iota + 40 // 黑色背景
    BgRed                           // 红色背景
    BgGreen                         // 绿色背景
    BgYellow                        // 黄色背景
    BgBlue                          // 蓝色背景
    BgMagenta                       // 洋红色/品红色背景
    BgCyan                          // 青色背景
    BgWhite                         // 白色背景
    BgGray = BgHiBlack              // 灰色背景（BgHiBlack 的别名）
)
```

#### 背景高亮文本颜色（高亮 8 色）

这些是高亮版本的背景色，比标准色更明亮。

```go
const (
    BgHiBlack   Attribute = iota + 100 // 高亮黑色背景（灰色）
    BgHiRed                            // 高亮红色背景
    BgHiGreen                          // 高亮绿色背景
    BgHiYellow                         // 高亮黄色背景
    BgHiBlue                           // 高亮蓝色背景
    BgHiMagenta                        // 高亮洋红色背景
    BgHiCyan                           // 高亮青色背景
    BgHiWhite                          // 高亮白色背景
)
```

---

### Color

```go
type Color struct {
    // Has unexported fields.
}
```

Color 定义一个由 SGR 参数定义的自定义颜色对象。

#### BgRGB

```go
func BgRGB(r, g, b int) *Color
```

BgRGB 返回一个新的 24 位 RGB 背景色。

**参数:**
- `r`: 红色分量 (0-255)
- `g`: 绿色分量 (0-255)
- `b`: 蓝色分量 (0-255)

**返回值:**
- `*Color`: 配置好的颜色对象

---

#### New

```go
func New(value ...Attribute) *Color
```

New 返回一个新创建的颜色对象。

**参数:**
- `value`: 任意数量的 SGR 参数。

**返回值:**
- `*Color`: 新创建的颜色对象。

---

#### RGB

```go
func RGB(r, g, b int) *Color
```

RGB 返回一个新的 24 位 RGB 前景色。

**参数:**
- `r`: 红色分量 (0-255)
- `g`: 绿色分量 (0-255)
- `b`: 蓝色分量 (0-255)

**返回值:**
- `*Color`: 配置好的颜色对象

---

#### Set

```go
func Set(p ...Attribute) *Color
```

Set 立即设置给定的 SGR 参数。 将使用给定的 SGR 参数更改输出颜色，直到调用 color.Unset() 为止。

**参数:**
- `p`: 任意数量的 SGR 参数

**返回值:**
- `*Color`: 配置好的颜色对象

---

#### Add

```go
func (c *Color) Add(value ...Attribute) *Color
```

Add 用于链式添加 SGR 参数。 可以使用任意数量的参数进行组合并创建自定义颜色对象。

**参数:**
- `value`: 任意数量的 SGR 参数

**返回值:**
- `*Color`: 当前颜色对象，支持链式调用

**示例:**

```go
c.Add(FgRed, Underline).Println("红色下划线文本")
```

---

#### AddBgRGB

```go
func (c *Color) AddBgRGB(r, g, b int) *Color
```

AddBgRGB 用于链式添加背景 RGB SGR 参数。 可以使用任意数量的参数进行组合并创建自定义颜色对象。

**参数:**
- `r`: 红色分量 (0-255)
- `g`: 绿色分量 (0-255)
- `b`: 蓝色分量 (0-255)

**返回值:**
- `*Color`: 当前颜色对象，支持链式调用

**示例:**

```go
c.AddBgRGB(255, 128, 0).Println("橙色背景")
```

---

#### AddRGB

```go
func (c *Color) AddRGB(r, g, b int) *Color
```

AddRGB 用于链式添加前景 RGB SGR 参数。 可以使用任意数量的参数进行组合并创建自定义颜色对象。

**参数:**
- `r`: 红色分量 (0-255)
- `g`: 绿色分量 (0-255)
- `b`: 蓝色分量 (0-255)

**返回值:**
- `*Color`: 当前颜色对象，支持链式调用

**示例:**

```go
c.AddRGB(255, 128, 0).Println("橙色文本")
```

---

#### DisableColor

```go
func (c *Color) DisableColor()
```

DisableColor 禁用颜色输出。 可用于在不更改任何现有代码的情况下禁用颜色输出，例如配合 "--no-color" 标志使用。 要重新启用，请使用 EnableColor() 方法。

---

#### EnableColor

```go
func (c *Color) EnableColor()
```

EnableColor 启用颜色输出。 与 DisableColor() 一起使用。如果颜色未被禁用，此方法没有副作用。

---

#### Equals

```go
func (c *Color) Equals(c2 *Color) bool
```

Equals 比较两种颜色是否相等。

**参数:**
- `c2`: 要比较的另一个颜色对象

**返回值:**
- `bool`: 如果两种颜色相等则返回 true

---

#### Fprint

```go
func (c *Color) Fprint(w io.Writer, a ...interface{}) (n int, err error)
```

Fprint 使用其操作数的默认格式进行格式化并写入 w。 当操作数都不是字符串时，在它们之间添加空格。

**参数:**
- `w`: 目标写入器
- `a`: 要格式化的操作数

**返回值:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误（如果有）

> **注意:** 在 Windows 上，如果 w 是 *os.File 类型，用户应该用 colorable.NewColorable() 包装 w。

---

#### FprintFunc

```go
func (c *Color) FprintFunc() func(w io.Writer, a ...interface{})
```

FprintFunc 返回一个新函数，该函数使用 color.Fprint() 将传入的参数打印为彩色。

**返回值:**
- `func(w io.Writer, a ...interface{})`: 打印函数

---

#### Fprintf

```go
func (c *Color) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
```

Fprintf 根据格式说明符进行格式化并写入 w。

**参数:**
- `w`: 目标写入器
- `format`: 格式字符串
- `a`: 要格式化的操作数

**返回值:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误（如果有）

> **注意:** 在 Windows 上，如果 w 是 *os.File 类型，用户应该用 colorable.NewColorable() 包装 w。

---

#### FprintfFunc

```go
func (c *Color) FprintfFunc() func(w io.Writer, format string, a ...interface{})
```

FprintfFunc 返回一个新函数，该函数使用 color.Fprintf() 将传入的参数打印为彩色。

**返回值:**
- `func(w io.Writer, format string, a ...interface{})`: 打印函数

---

#### Fprintln

```go
func (c *Color) Fprintln(w io.Writer, a ...interface{}) (n int, err error)
```

Fprintln 使用其操作数的默认格式进行格式化并写入 w，并在末尾添加换行符。

**参数:**
- `w`: 目标写入器
- `a`: 要格式化的操作数

**返回值:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误（如果有）

> **注意:** 在 Windows 上，如果 w 是 *os.File 类型，用户应该用 colorable.NewColorable() 包装 w。

---

#### FprintlnFunc

```go
func (c *Color) FprintlnFunc() func(w io.Writer, a ...interface{})
```

FprintlnFunc 返回一个新函数，该函数使用 color.Fprintln() 将传入的参数打印为彩色。

**返回值:**
- `func(w io.Writer, a ...interface{})`: 打印函数

---

#### Print

```go
func (c *Color) Print(a ...interface{}) (n int, err error)
```

Print 使用其操作数的默认格式进行格式化并写入标准输出。 当操作数都不是字符串时，在它们之间添加空格。

**参数:**
- `a`: 要格式化的操作数

**返回值:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误（如果有）

---

#### PrintFunc

```go
func (c *Color) PrintFunc() func(a ...interface{})
```

PrintFunc 返回一个新函数，该函数使用 color.Print() 将传入的参数打印为彩色。

**返回值:**
- `func(a ...interface{})`: 打印函数

---

#### Printf

```go
func (c *Color) Printf(format string, a ...interface{}) (n int, err error)
```

Printf 根据格式说明符进行格式化并写入标准输出。

**参数:**
- `format`: 格式字符串
- `a`: 要格式化的操作数

**返回值:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误（如果有）

---

#### PrintfFunc

```go
func (c *Color) PrintfFunc() func(format string, a ...interface{})
```

PrintfFunc 返回一个新函数，该函数使用 color.Printf() 将传入的参数打印为彩色。

**返回值:**
- `func(format string, a ...interface{})`: 打印函数

---

#### Println

```go
func (c *Color) Println(a ...interface{}) (n int, err error)
```

Println 使用其操作数的默认格式进行格式化并写入标准输出，并在末尾添加换行符。

**参数:**
- `a`: 要格式化的操作数

**返回值:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误（如果有）

---

#### PrintlnFunc

```go
func (c *Color) PrintlnFunc() func(a ...interface{})
```

PrintlnFunc 返回一个新函数，该函数使用 color.Println() 将传入的参数打印为彩色。

**返回值:**
- `func(a ...interface{})`: 打印函数

---

#### Sprint

```go
func (c *Color) Sprint(a ...interface{}) string
```

Sprint 使用其操作数的默认格式进行格式化并返回结果字符串。 当操作数都不是字符串时，在它们之间添加空格。

**参数:**
- `a`: 要格式化的操作数

**返回值:**
- `string`: 格式化后的字符串

---

#### SprintFunc

```go
func (c *Color) SprintFunc() func(a ...interface{}) string
```

SprintFunc 返回一个新函数，该函数使用 color.Sprint() 将传入的参数格式化为彩色字符串。

**返回值:**
- `func(a ...interface{}) string`: 格式化函数

---

#### Sprintf

```go
func (c *Color) Sprintf(format string, a ...interface{}) string
```

Sprintf 根据格式说明符进行格式化并返回结果字符串。

**参数:**
- `format`: 格式字符串
- `a`: 要格式化的操作数

**返回值:**
- `string`: 格式化后的字符串

---

#### SprintfFunc

```go
func (c *Color) SprintfFunc() func(format string, a ...interface{}) string
```

SprintfFunc 返回一个新函数，该函数使用 color.Sprintf() 将传入的参数格式化为彩色字符串。

**返回值:**
- `func(format string, a ...interface{}) string`: 格式化函数

---

#### Sprintln

```go
func (c *Color) Sprintln(a ...interface{}) string
```

Sprintln 使用其操作数的默认格式进行格式化并返回结果字符串，并在末尾添加换行符。

**参数:**
- `a`: 要格式化的操作数

**返回值:**
- `string`: 格式化后的字符串

---

#### SprintlnFunc

```go
func (c *Color) SprintlnFunc() func(a ...interface{}) string
```

SprintlnFunc 返回一个新函数，该函数使用 color.Sprintln() 将传入的参数格式化为彩色字符串。

**返回值:**
- `func(a ...interface{}) string`: 格式化函数

---

#### UnsetWriter

```go
func (c *Color) UnsetWriter(w io.Writer)
```

UnsetWriter 使用给定的 io.Writer 重置所有转义属性并清除输出。 通常在 SetWriter() 之后调用。

**参数:**
- `w`: 目标写入器

---

## 全局实例

全局实例提供了一种便捷的方式来使用颜色输出，无需手动创建 Color 对象。全局实例具有独立的配置管理，支持样式设置（加粗、下划线等），并且是线程安全的。

### 使用示例

```go
// 使用全局实例直接打印
color.GetGlobal().Red("错误信息")
color.GetGlobal().Green("成功信息")

// 设置样式配置
color.GetGlobal().SetBold(true).SetUnderline(true)
color.GetGlobal().Blue("带下划线的蓝色粗体文本")

// 获取带颜色的字符串（不打印）
redText := color.GetGlobal().SRed("红色文本")
```

---

### GlobalColor

```go
type GlobalColor struct {
    // Has unexported fields.
}
```

GlobalColor 是全局颜色实例的类型，包含独立的配置和输出设置，与普通的 Color 实例完全分离。

---

### StyleConfig

```go
type StyleConfig struct {
    NoColor    bool      // 是否禁用颜色
    Bold       bool      // 加粗
    Underline  bool      // 下划线
    Italic     bool      // 斜体
    Blink      bool      // 闪烁
    Faint      bool      // 暗淡
    CrossedOut bool      // 删除线
    Output     io.Writer // 输出目标
}
```

StyleConfig 定义颜色样式配置。

#### Clone

```go
func (s *StyleConfig) Clone() *StyleConfig
```

Clone 创建样式配置的深拷贝，返回一个新的 StyleConfig 实例，复制当前配置的所有字段。

**注意:** Output 字段是接口类型，只会复制引用，不会深拷贝底层的写入器。

**返回值:**
- `*StyleConfig`: 新的样式配置实例

---

### 全局实例函数

#### GetGlobal

```go
func GetGlobal() *GlobalColor
```

GetGlobal 返回全局颜色实例。首次调用时会自动初始化全局实例。

**返回值:**
- `*GlobalColor`: 全局颜色实例

---

#### G

```go
func G() *GlobalColor
```

G 是 GetGlobal 的快捷方式，返回全局颜色实例。使用更短的函数名，方便频繁调用。

**返回值:**
- `*GlobalColor`: 全局颜色实例

**示例:**

```go
// 使用 G() 快捷方式
c := color.G()
c.Red("红色文字")
c.Info("信息日志")

// 等价于
c := color.GetGlobal()
c.Red("红色文字")
```

---

#### ResetGlobal

```go
func ResetGlobal()
```

ResetGlobal 重置全局实例到默认状态。这会重新创建全局实例并应用默认配置。

---

### GlobalColor 方法

GlobalColor 提供以下配置方法，均支持链式调用：

| 方法 | 说明 |
|------|------|
| `SetConfig(config)` | 设置样式配置（自动克隆） |
| `GetConfig()` | 获取当前配置 |
| `GetConfigClone()` | 获取配置的克隆副本 |
| `SetOutput(w)` | 设置输出目标 |
| `SetNoColor(bool)` | 启用/禁用颜色 |
| `SetBold(bool)` | 启用/禁用加粗 |
| `SetUnderline(bool)` | 启用/禁用下划线 |
| `SetItalic(bool)` | 启用/禁用斜体 |
| `SetBlink(bool)` | 启用/禁用闪烁 |
| `SetFaint(bool)` | 启用/禁用暗淡 |
| `SetCrossedOut(bool)` | 启用/禁用删除线 |

**颜色方法：**

与包级便捷函数类似，GlobalColor 提供相同命名的颜色方法：

```go
// 打印方法（自动追加换行符）
g.Black("文本")
g.Red("错误: %s", err)
g.Green("成功")

// 返回字符串方法
str := g.SRed("红色文本")
str := g.SGreen("绿色%s", "文本")
```

支持的颜色与包级函数相同：Black, Red, Green, Yellow, Blue, Magenta, Cyan, White, Gray 及其高亮版本（HiXxx）和字符串返回版本（SXxx, SHiXxx）。

---

