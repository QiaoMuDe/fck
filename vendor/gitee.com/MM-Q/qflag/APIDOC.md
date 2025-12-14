# Package qflag

`import "gitee.com/MM-Q/qflag"`

Package qflag 根包统一导出入口。本文件用于将各子包的核心功能导出到根包，简化外部使用。通过类型别名和变量导出的方式，为用户提供统一的API接口。

Package qflag 提供对标准库flag的封装，自动实现长短标志，并默认绑定-h/--help标志打印帮助信息。用户可通过Cmd.Help字段自定义帮助内容，支持直接赋值字符串或从文件加载。该包是一个功能强大的命令行参数解析库，支持子命令、多种数据类型标志、参数验证等高级特性。

## Variables

### NewCmd

```go
var NewCmd = cmd.NewCmd
```

NewCmd 创建新的命令实例。

**参数:**
- longName: 命令的全称（如: ls, rm, mkdir 等）
- shortName: 命令的简称（如: l, r, m 等）
- errorHandling: 标志解析错误处理策略

**返回值:**
- `*cmd.Cmd`: 新创建的命令实例

**errorHandling可选值:**
- `flag.ContinueOnError`: 遇到错误时继续解析，并将错误返回
- `flag.ExitOnError`: 遇到错误时立即退出程序，并将错误返回
- `flag.PanicOnError`: 遇到错误时立即触发panic，并将错误返回

## Functions

### Parse

```go
func Parse() error
```

Parse 解析所有命令行参数，包括根命令和所有子命令的标志参数。

**返回:**
- `error`: 解析过程中遇到的错误，若成功则为 nil

### ParseFlagsOnly

```go
func ParseFlagsOnly() error
```

ParseFlagsOnly 解析根命令的所有标志参数，不包括子命令。

**返回:**
- `error`: 解析过程中遇到的错误，若成功则为 nil

## Types

### BoolFlag

```go
type BoolFlag = flags.BoolFlag
```

BoolFlag 导出flag包中的BoolFlag结构体。

### Cmd

```go
type Cmd = cmd.Cmd
```

Cmd 导出cmd包中的Cmd结构体。

### Root

```go
var Root *Cmd
```

Root 全局根命令实例，提供对全局标志和子命令的访问。用户可以通过 `qflag.Root.String()` 这样的方式直接创建全局标志。这是访问命令行功能的主要入口点，推荐优先使用。

### DurationFlag

```go
type DurationFlag = flags.DurationFlag
```

DurationFlag 导出flag包中的DurationFlag结构体。

### EnumFlag

```go
type EnumFlag = flags.EnumFlag
```

EnumFlag 导出flag包中的EnumFlag结构体。

### Flag

```go
type Flag = flags.Flag
```

Flag 导出flag包中的Flag结构体。

### Float64Flag

```go
type Float64Flag = flags.Float64Flag
```

Float64Flag 导出flag包中的Float64Flag结构体。

### Int64Flag

```go
type Int64Flag = flags.Int64Flag
```

Int64Flag 导出flag包中的Int64Flag结构体。

### Int64SliceFlag

```go
type Int64SliceFlag = flags.Int64SliceFlag
```

Int64SliceFlag 导出flag包中的Int64SliceFlag结构体。

### IntFlag

```go
type IntFlag = flags.IntFlag
```

IntFlag 导出flag包中的IntFlag结构体。

### IntSliceFlag

```go
type IntSliceFlag = flags.IntSliceFlag
```

IntSliceFlag 导出flag包中的IntSliceFlag结构体。

### MapFlag

```go
type MapFlag = flags.MapFlag
```

MapFlag 导出flag包中的MapFlag结构体。

### SizeFlag

```go
type SizeFlag = flags.SizeFlag
```

SizeFlag 导出flag包中的SizeFlag结构体。

### StringFlag

```go
type StringFlag = flags.StringFlag
```

StringFlag 导出flag包中的StringFlag结构体。

### StringSliceFlag

```go
type StringSliceFlag = flags.StringSliceFlag
```

StringSliceFlag 导出flag包中的StringSliceFlag结构体。

### TimeFlag

```go
type TimeFlag = flags.TimeFlag
```

TimeFlag 导出flag包中的TimeFlag结构体。

### Uint16Flag

```go
type Uint16Flag = flags.Uint16Flag
```

Uint16Flag 导出flag包中的Uint16Flag结构体。

### Uint32Flag

```go
type Uint32Flag = flags.Uint32Flag
```

Uint32Flag 导出flag包中的Uint32Flag结构体。

### Uint64Flag

```go
type Uint64Flag = flags.Uint64Flag
```

Uint64Flag 导出flag包中的Uint64Flag结构体。