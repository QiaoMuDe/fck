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

// FileInfo 文件信息结构
type FileInfo struct {
	EntryType      string    // 文件类型
	Name           string    // 文件名
	Size           int64     // 文件大小
	ModTime        time.Time // 修改时间
	Perm           string    // 权限
	Owner          string    // 所有者
	Group          string    // 组
	FileExt        string    // 文件扩展名
	LinkTargetPath string    // 符号链接目标路径
}

// FileInfoList 文件信息列表类型
type FileInfoList []FileInfo
