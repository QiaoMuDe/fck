// Package list 定义了 list 子命令使用的数据模型和结构体。
// 该文件包含文件信息结构体、扫描选项、处理选项、格式化选项等核心数据类型定义。
package list

import (
	"time"
)

// ScanOptions 扫描选项
type ScanOptions struct {
	Recursive  bool     // 是否递归扫描
	ShowHidden bool     // 是否显示隐藏文件
	FileTypes  []string // 文件类型过滤
	DirItself  bool     // 是否只显示目录本身
}

// ProcessOptions 处理选项
type ProcessOptions struct {
	SortBy     string // 排序方式: "name", "time", "size"
	Reverse    bool   // 是否反向排序
	GroupByDir bool   // 是否按目录分组
}

// FormatOptions 格式化选项
type FormatOptions struct {
	LongFormat    bool   // 是否长格式显示
	UseColor      bool   // 是否使用颜色
	DevColor      bool   // 是否使用开发者颜色
	TableStyle    string // 表格样式
	QuoteNames    bool   // 是否引用文件名
	ShowUserGroup bool   // 是否显示用户组
}

// list子命令用于存储文件信息的结构体
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

// FileInfoList 文件信息列表类型
type FileInfoList []FileInfo

// 定义全局常量的颜色映射
var permissionColorMap = map[int]string{
	1: "green",  // 所有者-读-绿色
	2: "yellow", // 所有者-写-黄色
	3: "red",    // 所有者-执行-红色
	4: "green",  // 组-读-绿色
	5: "yellow", // 组-写-黄色
	6: "red",    // 组-执行-红色
	7: "green",  // 其他-读-绿色
	8: "yellow", // 其他-写-黄色
	9: "red",    // 其他-执行-红色
}
