# Package common

Package common 提供了 Windows 系统特定的文件属性检查功能。该文件实现了 Windows 平台下的隐藏文件检测、只读属性检查等系统相关功能。

Package common 提供了跨模块共享的通用工具函数和实用程序。该文件包含文件哈希计算、路径处理、错误处理、颜色输出等常用功能。

## CONSTANTS

### 字节单位定义

```go
const (
	Byte = 1 << (10 * iota) // 1 字节
	KB                      // 千字节 (1024 B)
	MB                      // 兆字节 (1024 KB)
)
```

## FUNCTIONS

### Checksum

Checksum 计算文件哈希值。

```go
func Checksum(filePath string, hashFunc func() hash.Hash) (string, error)
```

- 参数：
  - `filePath`：文件路径
  - `hashFunc`：哈希函数构造器
- 返回：
  - `string`：文件的十六进制哈希值
  - `error`：错误信息，如果计算失败
- 注意：
  - 根据文件大小动态分配缓冲区以提高性能
  - 支持任何实现`hash.Hash`接口的哈希算法
  - 使用`io.CopyBuffer`进行高效的文件读取和哈希计算

### ChecksumProgress

ChecksumProgress 计算文件哈希值（带进度条）。

```go
func ChecksumProgress(filePath string, hashFunc func() hash.Hash) (string, error)
```

- 参数：
  - `filePath`：文件路径
  - `hashFunc`：哈希函数构造器
- 返回：
  - `string`：文件的十六进制哈希值
  - `error`：错误信息，如果计算失败
- 注意：
  - 根据文件大小动态分配缓冲区以提高性能
  - 支持任何实现`hash.Hash`接口的哈希算法
  - 使用`io.CopyBuffer`进行高效的文件读取和哈希计算

### CompileRegex

CompileRegex 编译正则表达式。

```go
func CompileRegex(pattern string) (*regexp.Regexp, error)
```

- 参数：
  - `pattern`：正则表达式模式字符串
- 返回：
  - `*regexp.Regexp`：编译后的正则表达式对象，如果模式为空则返回`nil`
  - `error`：编译错误信息，如果编译失败
- 注意：
  - 仅在模式不为空时进行编译
  - 空模式会返回`nil`而不是错误

### GetFileOwner

GetFileOwner 用于Windows环境下的占位函数。

```go
func GetFileOwner(filePath string) (string, string)
```

- 参数：
  - `filePath`：文件路径
- 返回：
  - `string`：文件所有者的用户名
  - `string`：文件所有者的组名
- 注意：
  - 该函数在Windows环境下始终返回`?`问号。

### GetLast8Chars

GetLast8Chars 获取输入字符串的最后8个字符。

```go
func GetLast8Chars(s string) string
```

- 参数：
  - `s`：输入字符串
- 返回：
  - `string`：字符串的最后8个字符，如果字符串长度不足8个字符则返回原字符串

### HandleError

HandleError 处理路径检查时出现的错误。

```go
func HandleError(path string, err error) error
```

- 参数：
  - `path`：当前正在检查的路径（文件路径或目录路径）
  - `err`：在检查路径过程中产生的错误对象
- 返回：
  - `error`：包含更具描述性错误信息的新错误对象
- 注意：
  - 根据不同的错误类型生成对应的错误提示信息
  - 支持的错误类型：无效字符、权限错误、路径不存在等
  - 便于调用者定位和处理问题

### IsDriveRoot

IsDriveRoot 检查路径是否是盘符根目录。

```go
func IsDriveRoot(path string) bool
```

- 参数：
  - `path`：路径
- 返回：
  - `bool`：是否是盘符根目录

### IsHidden

IsHidden 判断Windows文件或目录是否为隐藏。

```go
func IsHidden(path string) bool
```

- 参数：
  - `path`：文件或目录路径
- 返回：
  - `bool`：是否为隐藏

### IsHiddenWindows

IsHiddenWindows 检查Windows文件是否为隐藏。

```go
func IsHiddenWindows(path string) bool
```

- 参数：
  - `path`：文件路径
- 返回：
  - `bool`：是否为隐藏文件

### IsReadOnly

IsReadOnly 判断Windows文件或目录是否为只读。

```go
func IsReadOnly(path string) bool
```

- 参数：
  - `path`：文件或目录的路径
- 返回：
  - `bool`：文件或目录是否为只读

### IsSystemFileOrDir

IsSystemFileOrDir 检查文件或目录是否是系统文件或特殊目录。

```go
func IsSystemFileOrDir(name string) bool
```

- 参数：
  - `name`：文件或目录名称
- 返回：
  - `bool`：是否为系统文件或特殊目录
- 注意：
  - 检查预定义的系统文件和目录列表
  - 包括Windows系统文件如`pagefile.sys`、`$RECYCLE.BIN`等

### RegexBuilder

RegexBuilder 构建正则表达式模式字符串。

```go
func RegexBuilder(pattern string, isRegex, wholeWord, caseSensitive bool) string
```

- 参数：
  - `pattern`：原始匹配模式
  - `isRegex`：是否启用正则表达式模式（`true`时不转义特殊字符）
  - `wholeWord`：是否启用全字匹配（在模式前后添加`^`和`$`）
  - `caseSensitive`：是否区分大小写（`false`时添加`(?i)`标志）
- 返回：
  - `string`：构建后的正则表达式字符串
- 注意：
  - 非正则模式下会转义特殊字符
  - 全字匹配会在模式前后添加`^`和`$`
  - 不区分大小写时会添加`(?i)`标志

### SprintStringColor

SprintStringColor 根据路径类型以不同颜色输出字符串。

```go
func SprintStringColor(p string, s string, cl *colorlib.ColorLib) string
```

- 参数：
  - `p`：要检查的路径，用于获取文件类型信息
  - `s`：要着色的字符串内容
  - `cl`：`colorlib.ColorLib`实例，用于彩色输出
- 返回：
  - `string`：根据路径类型以不同颜色返回的字符串

### WriteFileHeader

WriteFileHeader 写入文件头信息。

```go
func WriteFileHeader(file *os.File, hashType string, timestampFormat string) error
```

- 参数：
  - `file`：要写入的文件对象
  - `hashType`：哈希类型标识
  - `timestampFormat`：时间戳格式
- 返回：
  - `error`：错误信息，如果写入失败
- 注意：
  - 文件头格式为：`#hashType#timestamp`

### GetFileColorByExtension

GetFileColorByExtension 根据文件扩展名返回相应颜色（统一的按扩展名配色规则）。该函数会处理 macOS 系统文件、无扩展名的特殊配置文件，并基于扩展名集合返回相应的颜色化字符串。

```go
func GetFileColorByExtension(ext, path string, cl *colorlib.ColorLib) string
```

- 参数：
  - `ext`：文件扩展名（建议以包含点的形式传入，如“.go”，内部会统一转小写）
  - `path`：文件路径（用于提取文件名以及处理特殊文件名）
  - `cl`：`colorlib.ColorLib` 实例，用于输出对应颜色的字符串
- 返回：
  - `string`：根据规则着色后的整条路径字符串
- 规则说明：
  - 特殊的 macOS 系统文件（如 `.DS_Store`、`.localized`、以 `._` 开头）统一使用灰色
  - 无扩展名但属于特殊配置文件的名称统一使用黄色
  - 扩展名命中预定义集合时，返回：
    - 绿色：代码/脚本/可执行相关扩展
    - 黄色：配置/日志等
    - 红色：数据/压缩包等
    - 紫色：库文件/编译产物等
  - 未命中任何集合时，使用白色