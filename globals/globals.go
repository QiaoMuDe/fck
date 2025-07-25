package globals

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	_ "embed"
	"hash"
	"os"
	"regexp"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/colorlib"
	"github.com/jedib0t/go-pretty/v6/table"
)

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

// 虚拟哈希表条目
type VirtualHashEntry struct {
	// 真实路径
	RealPath string

	// 哈希值
	Hash string
}

// 虚拟哈希表
type VirtualHashMap map[string]VirtualHashEntry

var (
	// 支持的哈希算法列表
	SupportedAlgorithms = map[string]func() hash.Hash{
		"md5":    md5.New,
		"sha1":   sha1.New,
		"sha256": sha256.New,
		"sha512": sha512.New,
	}

	// 禁止输入的路径map
	ForbiddenPaths = map[string]bool{
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
)

// list子命令用于存储文件信息的结构体
type ListInfo struct {
	// 文件名 - BaseName
	Name string

	// 文件路径 - 绝对路径
	Path string

	// 类型 - 文件/目录/软链接
	EntryType string

	// 大小 - 字节数
	Size int64

	// 修改时间 - time.Time
	ModTime time.Time

	// 权限 - 类型-所有者-组-其他用户
	Perm string

	// 所属用户 - windows环境为?
	Owner string

	// 所属组 - windows环境为?
	Group string

	// 扩展名 - 扩展名
	FileExt string

	// 如果是软链接，则是指向的文件路径，否则为空字符串
	LinkTargetPath string
}

// list子命令用于存储文件信息的结构体切片
type ListInfos []ListInfo

// 用于存储listInfos的结构体
type ListInfosMap struct {
	ListInfos ListInfos `json:"listInfos"` // 存储文件信息的结构体切片
	FileNames []string  `json:"fileNames"` // 存储文件名的字符串切片
}

// SortByFileNameAsc 按照文件名升序排序
func (lis ListInfos) SortByFileNameAsc() {
	sort.Slice(lis, func(i, j int) bool {
		return lis[i].Name < lis[j].Name
	})
}

// SortByFileSizeAsc 按照文件大小升序排序
func (lis ListInfos) SortByFileSizeAsc() {
	sort.Slice(lis, func(i, j int) bool {
		return lis[i].Size < lis[j].Size
	})
}

// SortByModTimeAsc 按照修改时间升序排序
func (lis ListInfos) SortByModTimeAsc() {
	sort.Slice(lis, func(i, j int) bool {
		return lis[i].ModTime.Before(lis[j].ModTime)
	})
}

// SortByFileNameDesc 按照文件名降序排序
func (lis ListInfos) SortByFileNameDesc() {
	sort.Slice(lis, func(i, j int) bool {
		return lis[i].Name > lis[j].Name
	})
}

// SortByFileSizeDesc 按照文件大小降序排序
func (lis ListInfos) SortByFileSizeDesc() {
	sort.Slice(lis, func(i, j int) bool {
		return lis[i].Size > lis[j].Size
	})
}

// SortByModTimeDesc 按照修改时间降序排序
func (lis ListInfos) SortByModTimeDesc() {
	sort.Slice(lis, func(i, j int) bool {
		return lis[i].ModTime.After(lis[j].ModTime)
	})
}

// GetFileNames 获取文件名列表
func (lis ListInfos) GetFileNames() []string {
	names := make([]string, len(lis))
	for i, file := range lis {
		names[i] = file.Name
	}
	return names
}

// 定义文件类型标识符常量
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

// Table样式映射表
var TableStyleMap = map[string]table.Style{
	"default": table.StyleDefault,                    // 默认样式
	"l":       table.StyleLight,                      // 浅色样式
	"r":       table.StyleRounded,                    // 圆角样式
	"bd":      table.StyleBold,                       // 粗体样式
	"cb":      table.StyleColoredBright,              // 彩色亮色样式
	"cd":      table.StyleColoredDark,                // 彩色暗色样式
	"db":      table.StyleDouble,                     // 双线样式
	"cbb":     table.StyleColoredBlackOnBlueWhite,    // 黑色背景蓝色字体样式
	"cbc":     table.StyleColoredBlackOnCyanWhite,    // 青色背景蓝色字体样式
	"cbg":     table.StyleColoredBlackOnGreenWhite,   // 绿色背景蓝色字体样式
	"cbm":     table.StyleColoredBlackOnMagentaWhite, // 紫色背景蓝色字体样式
	"cby":     table.StyleColoredBlackOnYellowWhite,  // 黄色背景蓝色字体样式
	"cbr":     table.StyleColoredBlackOnRedWhite,     // 红色背景蓝色字体样式
	"cwb":     table.StyleColoredBlueWhiteOnBlack,    // 蓝色背景白色字体样式
	"ccw":     table.StyleColoredCyanWhiteOnBlack,    // 青色背景白色字体样式
	"cgw":     table.StyleColoredGreenWhiteOnBlack,   // 绿色背景白色字体样式
	"cmw":     table.StyleColoredMagentaWhiteOnBlack, // 紫色背景白色字体样式
	"crw":     table.StyleColoredRedWhiteOnBlack,     // 红色背景白色字体样式
	"cyw":     table.StyleColoredYellowWhiteOnBlack,  // 黄色背景白色字体样式
	"none":    StyleNone,                             // 禁用样式
}

// Table样式切片
var TableStyles = []string{
	"default", // 默认样式
	"l",       // 浅色样式
	"r",       // 圆角样式
	"bd",      // 粗体样式
	"cb",      // 彩色亮色样式
	"cd",      // 彩色暗色样式
	"db",      // 双线样式
	"cbb",     // 黑色背景蓝色字体样式
	"cbc",     // 青色背景蓝色字体样式
	"cbg",     // 绿色背景蓝色字体样式
	"cbm",     // 紫色背景蓝色字体样式
	"cby",     // 黄色背景蓝色字体样式
	"cbr",     // 红色背景蓝色字体样式
	"cwb",     // 蓝色背景白色字体样式
	"ccw",     // 青色背景白色字体样式
	"cgw",     // 绿色背景白色字体样式
	"cmw",     // 紫色背景白色字体样式
	"crw",     // 红色背景白色字体样式
	"cyw",     // 黄色背景白色字体样式
	"none",    // 禁用样式
}

// 定义禁用样式
var StyleNone = table.Style{
	Box: table.BoxStyle{
		PaddingLeft:      " ", // 左边框
		PaddingRight:     " ", // 右边框
		MiddleHorizontal: " ", // 水平线
		MiddleVertical:   " ", // 垂直线
		TopLeft:          " ", // 左上角
		TopRight:         " ", // 右上角
		BottomLeft:       " ", // 左下角
		BottomRight:      " ", // 右下角
	},
}

// 包装 os.DirEntry 以便复用 processFindCmd
type DirEntryWrapper struct {
	NameVal  string
	IsDirVal bool
	ModeVal  os.FileMode
}

func (d *DirEntryWrapper) Name() string               { return d.NameVal }
func (d *DirEntryWrapper) IsDir() bool                { return d.IsDirVal }
func (d *DirEntryWrapper) Type() os.FileMode          { return d.ModeVal }
func (d *DirEntryWrapper) Info() (os.FileInfo, error) { return nil, nil }

// 查找类型常量定义
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

// 限制查找的参数切片
var TypeLimits = []string{
	// 查找所有类型
	FindTypeAll,

	// 只查找文件
	FindTypeFile,
	// 只查找文件-短参数
	FindTypeFileShort,

	// 只查找目录
	FindTypeDir,
	// 只查找目录-短参数
	FindTypeDirShort,

	// 只查找软链接
	FindTypeSymlink,
	// 只查找软链接-短参数
	FindTypeSymlinkShort,

	// 只查找只读文件
	FindTypeReadonly,
	// 只查找只读文件-短参数
	FindTypeReadonlyShort,

	// 只查找隐藏文件或目录
	FindTypeHidden,
	// 只查找隐藏文件或目录-短参数
	FindTypeHiddenShort,

	// // 只查找空文件或目录
	// FindTypeEmpty,
	// // 只查找空文件或目录-短参数
	// FindTypeEmptyShort,

	// // 只查找可执行文件
	// FindTypeExecutable,
	// // 只查找可执行文件-短参数
	// FindTypeExecutableShort,

	// // 只查找socket文件(套接字)
	// FindTypeSocket,
	// // 只查找socket文件-短参数
	// FindTypeSocketShort,

	// // 只查找管道文件
	// FindTypePipe,
	// // 只查找管道文件-短参数
	// FindTypePipeShort,

	// // 只查找块设备文件
	// FindTypeBlock,
	// // 只查找块设备文件-短参数
	// FindTypeBlockShort,

	// // 只查找字符设备文件
	// FindTypeChar,
	// // 只查找字符设备文件-短参数
	// FindTypeCharShort,

	// // 只查找追加模式的文件
	// FindTypeAppend,
	// // 只查找追加模式的文件-短参数
	// FindTypeAppendShort,

	// // 只查找非追加模式的文件
	// FindTypeNonAppend,
	// // 只查找非追加模式的文件-短参数
	// FindTypeNonAppendShort,

	// // 只查找为独占模式的文件
	// FindTypeExclusive,
	// // 只查找为独占模式的文件-短参数
	// FindTypeExclusiveShort,
}

// 定义find子命令限制查找的参数
var FindLimits = map[string]bool{
	FindTypeAll:             true, // 查找所有类型
	FindTypeFile:            true, // 只查找文件
	FindTypeFileShort:       true, // 只查找文件-短参数
	FindTypeDir:             true, // 只查找目录
	FindTypeDirShort:        true, // 只查找目录-短参数
	FindTypeSymlink:         true, // 只查找软链接
	FindTypeSymlinkShort:    true, // 只查找软链接-短参数
	FindTypeReadonly:        true, // 只查找只读文件
	FindTypeReadonlyShort:   true, // 只查找只读文件-短参数
	FindTypeHidden:          true, // 只查找隐藏文件或目录
	FindTypeHiddenShort:     true, // 只查找隐藏文件或目录-短参数
	FindTypeEmpty:           true, // 只查找空文件或目录
	FindTypeEmptyShort:      true, // 只查找空文件或目录-短参数
	FindTypeExecutable:      true, // 只查找可执行文件
	FindTypeExecutableShort: true, // 只查找可执行文件-短参数
	FindTypeSocket:          true, // 只查找socket文件(套接字)
	FindTypeSocketShort:     true, // 只查找socket文件-短参数
	FindTypePipe:            true, // 只查找管道文件
	FindTypePipeShort:       true, // 只查找管道文件-短参数
	FindTypeBlock:           true, // 只查找块设备文件
	FindTypeBlockShort:      true, // 只查找块设备文件-短参数
	FindTypeChar:            true, // 只查找字符设备文件
	FindTypeCharShort:       true, // 只查找字符设备文件-短参数
	FindTypeAppend:          true, // 只查找追加模式的文件
	FindTypeAppendShort:     true, // 只查找追加模式的文件-短参数
	FindTypeNonAppend:       true, // 只查找非追加模式的文件
	FindTypeNonAppendShort:  true, // 只查找非追加模式的文件-短参数
	FindTypeExclusive:       true, // 只查找为独占模式的文件
	FindTypeExclusiveShort:  true, // 只查找为独占模式的文件-短参数
}

// IsValidFindType 检查给定的类型参数是否有效
// 参数:
//   - typeStr: 要检查的类型字符串
//
// 返回值:
//   - bool: 如果类型有效返回true, 否则返回false
func IsValidFindType(typeStr string) bool {
	_, ok := FindLimits[typeStr]
	return ok
}

// GetSupportedFindTypes 获取所有支持的查找类型列表
// 返回值:
//   - []string: 包含所有支持类型的字符串切片
func GetSupportedFindTypes() []string {
	types := make([]string, 0, len(FindLimits))
	for t := range FindLimits {
		types = append(types, t)
	}
	return types
}

// 定义Windows可执行文件扩展名map
var WindowsExecutableExts = map[string]bool{
	".exe":  true, // 可执行文件
	".bat":  true, // 批处理文件
	".cmd":  true, // 命令文件
	".ps1":  true, // PowerShell脚本文件
	".psm1": true, // PowerShell模块文件
	".msi":  true, // Windows安装程序
}

// 定义windows系统软链接或快捷方式扩展名map
var WindowsSymlinkExts = map[string]bool{
	".lnk": true, // 快捷方式
	".url": true, // 链接文件
}

// 自定义fck命令的logo
var FckHelpLogo = `    ________      ________          ___  __       
   |\  _____\    |\   ____\        |\  \|\  \     
   \ \  \__/     \ \  \___|        \ \  \/  /|_   
    \ \   __\     \ \  \            \ \   ___  \  
     \ \  \_|      \ \  \____        \ \  \\ \  \ 
      \ \__\        \ \_______\       \ \__\\ \__\
       \|__|         \|_______|        \|__| \|__|
                   FCK CLI                        
`

// FindConfig 用于封装find命令的配置参数和共享资源
// 避免函数参数过多难以管理
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
