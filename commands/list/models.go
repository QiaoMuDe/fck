// Package list 定义了 list 子命令使用的数据模型和结构体。
// 该文件包含文件信息结构体、扫描选项、处理选项、格式化选项等核心数据类型定义。
package list

import (
	"time"
)

// timeFormat 时间格式化字符串
const timeFormat = "2006-01-02 15:04:05"

// 全局符号链接箭头常量
const symlinkArrow = " -> "

// ScanOptions 扫描选项
type ScanOptions struct {
	Recursive  bool     // 是否递归扫描
	ShowHidden bool     // 是否显示隐藏文件
	FileTypes  []string // 文件类型过滤
	DirItself  bool     // 是否只显示目录本身
}

// ProcessOptions 处理选项
type ProcessOptions struct {
	SortBy      string // 排序方式: "name", "time", "size"
	Reverse     bool   // 是否反向排序
	GroupByDir  bool   // 是否按目录分组 (原有的递归分组)
	GroupByPath bool   // 是否按路径分组 (新增：用于多路径/通配符场景)
	IsMultiPath bool   // 是否为多路径场景 (新增：标识符)
}

// FormatOptions 格式化选项
type FormatOptions struct {
	LongFormat    bool   // 是否长格式显示
	UseColor      bool   // 是否使用颜色
	TableStyle    string // 表格样式
	QuoteNames    bool   // 是否引用文件名
	ShowUserGroup bool   // 是否显示用户组
	ShouldGroup   bool   // 是否应该分组显示 (新增：避免重复判断)
	DisableIndex  bool   // 是否禁用索引
}

// list子命令用于存储文件信息的结构体
type FileInfo struct {
	Name           string    // 文件名 - BaseName
	Path           string    // 文件路径 - 绝对路径
	OriginalPath   string    // 原始路径 - 用户指定的路径（用于分组显示）
	EntryType      EntryType // 类型 - 文件/目录/软链接
	Size           int64     // 大小 - 字节数
	ModTime        time.Time // 修改时间 - time.Time
	Perm           string    // 权限 - 类型-所有者-组-其他用户
	Owner          string    // 所属用户 - windows环境为?
	Group          string    // 所属组 - windows环境为?
	FileExt        string    // 扩展名 - 扩展名
	LinkTargetPath string    // 如果是软链接，则是指向的文件路径，否则为空字符串
}

// FileInfoList 文件信息列表类型
type FileInfoList []FileInfo

// 定义全局常量的颜色映射
var permissionColorMap = map[int]colorType{
	1: colorTypeGreen,  // 所有者-读-绿色
	2: colorTypeYellow, // 所有者-写-黄色
	3: colorTypeRed,    // 所有者-执行-红色
	4: colorTypeGreen,  // 组-读-绿色
	5: colorTypeYellow, // 组-写-黄色
	6: colorTypeRed,    // 组-执行-红色
	7: colorTypeGreen,  // 其他-读-绿色
	8: colorTypeYellow, // 其他-写-黄色
	9: colorTypeRed,    // 其他-执行-红色
}

type colorType uint8

const (
	colorTypeGreen  colorType = iota // 绿色
	colorTypeYellow                  // 黄色
	colorTypeRed                     // 红色
)

// EntryType 定义文件类型
type EntryType string

// 定义文件类型标识符常量
const (
	DirType         EntryType = "d" // 目录类型
	SymlinkType     EntryType = "l" // 符号链接类型
	SocketType      EntryType = "s" // 套接字类型
	PipeType        EntryType = "p" // 管道类型
	BlockDeviceType EntryType = "b" // 块设备类型
	CharDeviceType  EntryType = "c" // 字符设备类型
	ExecutableType  EntryType = "x" // 可执行文件类型
	EmptyType       EntryType = "e" // 空文件类型
	FileType        EntryType = "f" // 普通文件类型
	UnknownType     EntryType = "?" // 未知类型
)

// String 返回文件类型对应的字符串
func (e EntryType) String() string {
	return string(e)
}
