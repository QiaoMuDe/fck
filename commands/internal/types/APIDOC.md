# Package types

Package types 定义了 fck 工具中使用的所有数据类型、常量和配置结构。该文件包含哈希算法映射、文件类型常量、表格样式配置、查找类型定义等核心类型定义。

## CONSTANTS

### 文件相关常量

```go
const (
	// 输出哈希值的文件名
	OutputFileName = "checksum.hash"

	// 输出对比结果的文件名
	OutputCheckFileName = "check_dir.check"

	// 时间戳格式
	TimestampFormat = "2006-01-02 15:04:05"

	// 虚拟基准目录 /ROOTDIR
	VirtualRootDir = "/ROOTDIR"
)
```

### 文件类型标识符常量

```go
const (
	DirType         = "d" // 目录类型
	SymlinkType     = "l" // 符号链接类型
	SocketType      = "s" // 套接字类型
	PipeType        = "p" // 管道类型
	BlockDeviceType = "b" // 块设备类型
	CharDeviceType  = "c" // 字符设备类型
	ExecutableType  = "x" // 可执行文件类型
	EmptyType       = "e" // 空文件类型
	FileType        = "f" // 普通文件类型
	UnknownType     = "?" // 未知类型
)
```

### 查找类型常量定义

```go
const (
	// 查找所有类型
	FindTypeAll = "all"

	// 只查找文件
	FindTypeFile = "file"
	// 只查找文件-短参数
	FindTypeFileShort = "f"

	// 只查找目录
	FindTypeDir = "dir"
	// 只查找目录-短参数
	FindTypeDirShort = "d"

	// 只查找软链接
	FindTypeSymlink = "symlink"
	// 只查找软链接-短参数
	FindTypeSymlinkShort = "l"

	// 只查找只读文件
	FindTypeReadonly = "readonly"
	// 只查找只读文件-短参数
	FindTypeReadonlyShort = "r"

	// 只查找隐藏文件或目录
	FindTypeHidden = "hidden"
	// 只查找隐藏文件或目录-短参数
	FindTypeHiddenShort = "h"

	// 只查找空文件或目录
	FindTypeEmpty = "empty"
	// 只查找空文件或目录-短参数
	FindTypeEmptyShort = "e"

	// 只查找可执行文件
	FindTypeExecutable = "executable"
	// 只查找可执行文件-短参数
	FindTypeExecutableShort = "x"

	// 只查找socket文件(套接字)
	FindTypeSocket = "socket"
	// 只查找socket文件-短参数
	FindTypeSocketShort = "s"

	// 只查找管道文件
	FindTypePipe = "pipe"
	// 只查找管道文件-短参数
	FindTypePipeShort = "p"

	// 只查找块设备文件
	FindTypeBlock = "block"
	// 只查找块设备文件-短参数
	FindTypeBlockShort = "b"

	// 只查找字符设备文件
	FindTypeChar = "char"
	// 只查找字符设备文件-短参数
	FindTypeCharShort = "c"

	// 只查找追加模式的文件
	FindTypeAppend = "append"
	// 只查找追加模式的文件-短参数
	FindTypeAppendShort = "a"

	// 只查找非追加模式的文件
	FindTypeNonAppend = "nonappend"
	// 只查找非追加模式的文件-短参数
	FindTypeNonAppendShort = "n"

	// 只查找为独占模式的文件
	FindTypeExclusive = "exclusive"
	// 只查找为独占模式的文件-短参数
	FindTypeExclusiveShort = "u"
)
```

## VARIABLES

### 哈希算法映射

```go
var (
	// 支持的哈希算法列表
	SupportedAlgorithms = map[string]func() hash.Hash{
		"md5":    md5.New,
		"sha1":   sha1.New,
		"sha256": sha256.New,
		"sha512": sha512.New,
	}
)
```

### 禁止输入的路径map

```go
var ForbiddenPaths = map[string]bool{
	"./":              true,
	".":               true,
	"..":              true,
	"...":             true,
	"....":            true,
	"./.":             true,
	"./..":            true,
	"./...":           true,
	"./....":          true,
	"./*":             true,
	"./**":            true,
	"./../":           true,
	"./../../":        true,
	"../":             true,
	"../../":          true,
	"../../../":       true,
	"../../../../":    true,
	"../../../../../": true,
	"../*":            true,
	"../**":           true,
	"../../*":         true,
}
```

### 自定义fck命令的logo

```go
var FckHelpLogo = `
    ________      ________          ___  __       
   |\  _____\    |\   ____\        |\  \|\  \     
   \ \  \__/     \ \  \___|        \ \  \/  /|_   
    \ \   __\     \ \  \            \ \   ___  \  
     \ \  \_|      \ \  \____        \ \  \\ \  \ 
      \ \__\        \ \_______\       \ \__\\ \__\
       \|__|         \|_______|        \|__| \|__|
                   FCK CLI                        
`
```

### 查找类型限制

```go
var FindLimits = map[string]bool{
	FindTypeAll:             true,
	FindTypeFile:            true,
	FindTypeFileShort:       true,
	FindTypeDir:             true,
	FindTypeDirShort:        true,
	FindTypeSymlink:         true,
	FindTypeSymlinkShort:    true,
	FindTypeReadonly:        true,
	FindTypeReadonlyShort:   true,
	FindTypeHidden:          true,
	FindTypeHiddenShort:     true,
	FindTypeEmpty:           true,
	FindTypeEmptyShort:      true,
	FindTypeExecutable:      true,
	FindTypeExecutableShort: true,
	FindTypeSocket:          true,
	FindTypeSocketShort:     true,
	FindTypePipe:            true,
	FindTypePipeShort:       true,
	FindTypeBlock:           true,
	FindTypeBlockShort:      true,
	FindTypeChar:            true,
	FindTypeCharShort:       true,
	FindTypeAppend:          true,
	FindTypeAppendShort:     true,
	FindTypeNonAppend:       true,
	FindTypeNonAppendShort:  true,
	FindTypeExclusive:       true,
	FindTypeExclusiveShort:  true,
}
```

### 查找类型限制切片

```go
var FindTypeLimits = []string{
	FindTypeAll,
	FindTypeFile,
	FindTypeFileShort,
	FindTypeDir,
	FindTypeDirShort,
	FindTypeSymlink,
	FindTypeSymlinkShort,
	FindTypeReadonly,
	FindTypeReadonlyShort,
	FindTypeHidden,
	FindTypeHiddenShort,
	FindTypeEmpty,
	FindTypeEmptyShort,
	FindTypeExecutable,
	FindTypeExecutableShort,
	FindTypeSocket,
	FindTypeSocketShort,
	FindTypePipe,
	FindTypePipeShort,
	FindTypeBlock,
	FindTypeBlockShort,
	FindTypeChar,
	FindTypeCharShort,
	FindTypeAppend,
	FindTypeAppendShort,
	FindTypeNonAppend,
	FindTypeNonAppendShort,
	FindTypeExclusive,
	FindTypeExclusiveShort,
}
```

### 列表查找类型限制切片

```go
var ListTypeLimits = []string{
	FindTypeAll,
	FindTypeFile,
	FindTypeFileShort,
	FindTypeDir,
	FindTypeDirShort,
	FindTypeSymlink,
	FindTypeSymlinkShort,
	FindTypeReadonly,
	FindTypeReadonlyShort,
	FindTypeHidden,
	FindTypeHiddenShort,
}
```

### 禁用样式

```go
var StyleNone = table.Style{
	Box: table.BoxStyle{
		PaddingLeft:      " ",
		PaddingRight:     "  ",
		MiddleHorizontal: " ",
		MiddleVertical:   " ",
		TopLeft:          " ",
		TopRight:         " ",
		BottomLeft:       " ",
		BottomRight:      " ",
	},
}
```

### Table样式映射表

```go
var TableStyleMap = map[string]table.Style{
	"default": table.StyleDefault,
	"l":       table.StyleLight,
	"r":       table.StyleRounded,
	"bd":      table.StyleBold,
	"cb":      table.StyleColoredBright,
	"cd":      table.StyleColoredDark,
	"db":      table.StyleDouble,
	"cbb":     table.StyleColoredBlackOnBlueWhite,
	"cbc":     table.StyleColoredBlackOnCyanWhite,
	"cbg":     table.StyleColoredBlackOnGreenWhite,
	"cbm":     table.StyleColoredBlackOnMagentaWhite,
	"cby":     table.StyleColoredBlackOnYellowWhite,
	"cbr":     table.StyleColoredBlackOnRedWhite,
	"cwb":     table.StyleColoredBlueWhiteOnBlack,
	"ccw":     table.StyleColoredCyanWhiteOnBlack,
	"cgw":     table.StyleColoredGreenWhiteOnBlack,
	"cmw":     table.StyleColoredMagentaWhiteOnBlack,
	"crw":     table.StyleColoredRedWhiteOnBlack,
	"cyw":     table.StyleColoredYellowWhiteOnBlack,
	"none":    StyleNone,
}
```

### Table样式切片

```go
var TableStyles = []string{
	"default",
	"l",
	"r",
	"bd",
	"cb",
	"cd",
	"db",
	"cbb",
	"cbc",
	"cbg",
	"cbm",
	"cby",
	"cbr",
	"cwb",
	"ccw",
	"cgw",
	"cmw",
	"crw",
	"cyw",
	"none",
}
```

### Windows可执行文件扩展名map

```go
var WindowsExecutableExts = map[string]bool{
	".exe":  true,
	".bat":  true,
	".cmd":  true,
	".ps1":  true,
	".psm1": true,
	".msi":  true,
}
```

### Windows软链接或快捷方式扩展名map

```go
var WindowsSymlinkExts = map[string]bool{
	".lnk": true,
	".url": true,
}
```

## FUNCTIONS

### GetSupportedFindTypes

GetSupportedFindTypes 获取所有支持的查找类型列表。

```go
func GetSupportedFindTypes() []string
```

- 返回值：
  - `[]string`：包含所有支持类型的字符串切片。

### IsValidFindType

IsValidFindType 检查给定的类型参数是否有效。

```go
func IsValidFindType(typeStr string) bool
```

- 参数：
  - `typeStr`：要检查的类型字符串。
- 返回值：
  - `bool`：如果类型有效返回`true`，否则返回`false`。

## TYPES

### DirEntryWrapper

DirEntryWrapper 包装`os.DirEntry`以便复用`processFindCmd`。

```go
type DirEntryWrapper struct {
	NameVal  string
	IsDirVal bool
	ModeVal  os.FileMode
}
```

#### Info

```go
func (d *DirEntryWrapper) Info() (os.FileInfo, error)
```

#### IsDir

```go
func (d *DirEntryWrapper) IsDir() bool
```

#### Name

```go
func (d *DirEntryWrapper) Name() string
```

#### Type

```go
func (d *DirEntryWrapper) Type() os.FileMode
```

### FindConfig

FindConfig 用于封装find命令的配置参数和共享资源，避免函数参数过多难以管理。

```go
type FindConfig struct {
	Cl              *colorlib.ColorLib // 颜色库实例
	NameRegex       *regexp.Regexp     // 文件名匹配正则
	ExNameRegex     *regexp.Regexp     // 排除文件名正则
	PathRegex       *regexp.Regexp     // 路径匹配正则
	ExPathRegex     *regexp.Regexp     // 排除路径正则
	IsRegex         bool               // 是否启用正则匹配
	WholeWord       bool               // 是否全词匹配
	CaseSensitive   bool               // 是否区分大小写
	MatchCount      *atomic.Int64      // 匹配计数原子变量
	NamePattern     string             // 文件名匹配模式
	PathPattern     string             // 路径匹配模式
	ExNamePattern   string             // 排除文件名匹配模式
	ExPathPattern   string             // 排除路径匹配模式
	FindExtSliceMap sync.Map           // ext切片标志的映射
}
```

### VirtualHashEntry

VirtualHashEntry 虚拟哈希表条目。

```go
type VirtualHashEntry struct {
	// 真实路径
	RealPath string

	// 哈希值
	Hash string
}
```

### VirtualHashMap

VirtualHashMap 虚拟哈希表。

```go
type VirtualHashMap map[string]VirtualHashEntry
```