package globals

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	_ "embed"
	"hash"
	"sort"
	"time"
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
	}
)

// hash 子命令帮助信息
//
//go:embed help/help_hash.txt
var HashHelp string

// size 子命令帮助信息
//
//go:embed help/help_size.txt
var SizeHelp string

// check 子命令帮助信息
//
//go:embed help/help_diff.txt
var DiffHelp string

// find 子命令帮助信息
//
//go:embed help/help_find.txt
var FindHelp string

// fck 主命令帮助信息
//
//go:embed help/help.txt
var FckHelp string

// list 子命令帮助信息
//
//go:embed help/help_list.txt
var ListHelp string

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
