# Package list

Package list 实现了文件列表显示命令的主要逻辑。该文件包含 list 子命令的入口函数，负责参数验证、路径处理、文件扫描和格式化输出。

Package list 实现了文件列表显示的颜色输出功能。该文件提供了根据文件类型和扩展名进行彩色显示的功能，支持普通模式和开发者模式两种配色方案。

Package list 定义了 list 子命令的命令行标志和参数配置。该文件包含所有 list 命令支持的选项，如排序方式、显示格式、过滤条件等。

Package list 实现了文件列表的格式化输出功能。该文件提供了表格格式和网格格式两种显示方式，支持颜色输出、权限显示、文件大小格式化等功能。

Package list 定义了 list 子命令使用的数据模型和结构体。该文件包含文件信息结构体、扫描选项、处理选项、格式化选项等核心数据类型定义。

Package list 实现了文件列表的数据处理功能。该文件提供了文件列表的排序处理，支持按名称、时间、大小等多种方式排序。

Package list 实现了文件系统扫描功能。该文件提供了文件和目录的扫描、过滤、类型识别等核心功能，支持递归扫描和多种文件类型过滤。

## FUNCTIONS

### GetColorString

根据文件信息返回带有相应颜色的路径字符串

```go
func GetColorString(devColor bool, info FileInfo, pF string, cl *colorlib.ColorLib) string
```

参数：
- `devColor`: 是否使用开发环境配色方案
- `info`: 包含文件类型和文件后缀名等信息的 FileInfo 结构体实例
- `pF`: 要处理的路径字符串
- `cl`: 用于彩色输出的 colorlib.ColorLib 实例

返回：
- `string`: 经过颜色处理后的路径字符串

注意：
- 如果启用了开发环境模式，将使用开发环境配色方案
- 支持 Windows、macOS、Linux 平台的特殊文件类型处理

### InitListCmd

初始化 list 命令

```go
func InitListCmd() *cmd.Cmd
```

### ListCmdMain

list 命令主函数

```go
func ListCmdMain(cl *colorlib.ColorLib) error
```

参数：
- `cl`: 颜色库

返回：
- `error`: 错误信息

## TYPES

### FileFormatter

文件格式化器

```go
type FileFormatter struct {
	// Has unexported fields.
}
```

### NewFileFormatter

创建新的文件格式化器

```go
func NewFileFormatter(cl *colorlib.ColorLib) *FileFormatter
```

### Render

渲染文件列表

```go
func (f *FileFormatter) Render(files FileInfoList, opts FormatOptions) error
```

参数：
- `files`: 文件列表
- `opts`: 格式选项

返回：
- `error`: 错误

### FileInfo

list 子命令用于存储文件信息的结构体

```go
type FileInfo struct {
	Name           string    // 文件名 - BaseName
	Path           string    // 文件路径 - 绝对路径
	EntryType      string    // 类型 - 文件/目录/软链接
	Size           int64     // 大小 - 字节数
	ModTime        time.Time // 修改时间 - time.Time
	Perm           string    // 权限 - 类型-所有者-组-其他用户
	Owner          string    // 所属用户 - windows环境为?
	Group          string    // 所属组 - windows环境为?
	FileExt        string    // 扩展名 - 扩展名
	LinkTargetPath string    // 如果是软链接，则是指向的文件路径，否则为空字符串
}
```

### FileInfoList

文件信息列表类型

```go
type FileInfoList []FileInfo
```

### FileProcessor

文件数据处理器

```go
type FileProcessor struct{}
```

### NewFileProcessor

创建新的文件处理器

```go
func NewFileProcessor() *FileProcessor
```

### Process

处理文件列表：过滤 -> 排序 -> 分组

```go
func (p *FileProcessor) Process(files FileInfoList, opts ProcessOptions) FileInfoList
```

参数：
- `files`: 文件列表
- `opts`: 处理选项

返回：
- `FileInfoList`: 处理后的文件列表

### FileScanner

文件扫描器

```go
type FileScanner struct {
	// Has unexported fields.
}
```

### NewFileScanner

创建新的文件扫描器

```go
func NewFileScanner() *FileScanner
```

### Scan

扫描指定路径的文件

```go
func (s *FileScanner) Scan(paths []string, opts ScanOptions) (FileInfoList, error)
```

参数：
- `paths`: 要扫描的路径列表
- `opts`: 扫描选项

返回：
- `FileInfoList`: 扫描到的文件信息列表
- `error`: 扫描过程中的错误

### FormatOptions

格式化选项

```go
type FormatOptions struct {
	LongFormat    bool   // 是否长格式显示
	UseColor      bool   // 是否使用颜色
	DevColor      bool   // 是否使用开发者颜色
	TableStyle    string // 表格样式
	QuoteNames    bool   // 是否引用文件名
	ShowUserGroup bool   // 是否显示用户组
}
```

### ProcessOptions

处理选项

```go
type ProcessOptions struct {
	SortBy     string // 排序方式: "name", "time", "size"
	Reverse    bool   // 是否反向排序
	GroupByDir bool   // 是否按目录分组
}
```

### ScanOptions

扫描选项

```go
type ScanOptions struct {
	Recursive  bool     // 是否递归扫描
	ShowHidden bool     // 是否显示隐藏文件
	FileTypes  []string // 文件类型过滤
	DirItself  bool     // 是否只显示目录本身
}
```
