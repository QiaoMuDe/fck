// Package types 定义了 fck 工具中使用的所有数据类型、常量和配置结构。
// 该文件包含哈希算法映射、文件类型常量、表格样式配置、查找类型定义等核心类型定义。
package types

import (
	_ "embed"
	"os"
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

	// 校验文件模式
	ChecksumModePortable = "PORTABLE"
	ChecksumModeLocal    = "LOCAL"

	// 行缓冲区初始大小 (64KB)
	InitialBufferSize = 64 * 1024

	// 行缓冲区默认最大大小 (10MB)
	DefaultMaxBufferSize = 10 * 1024 * 1024

	// Sed 命令备份文件后缀
	SedBackupSuffix = ".bak"

	// Sed 命令临时文件前缀模式
	SedTempFilePattern = ".sed-*.tmp"

	// GrepEncodingCheckSize UTF-8 编码检测样本大小 (8KB)
	GrepEncodingCheckSize = 8 * 1024
)

// 编码格式常量定义
const (
	EncodingNone     = "none" // 只检测不转换
	EncodingAuto     = "auto"
	EncodingUTF8     = "UTF-8"
	EncodingUTF8BOM  = "UTF-8-BOM"
	EncodingUTF16    = "UTF-16"
	EncodingUTF16LE  = "UTF-16LE"
	EncodingUTF16BE  = "UTF-16BE"
	EncodingGBK      = "GBK"
	EncodingGB2312   = "GB2312"
	EncodingGB18030  = "GB18030"
	EncodingBig5     = "Big5"
	EncodingLatin1   = "Latin1"
	EncodingShiftJIS = "ShiftJIS"
	EncodingEUCKR    = "EUC-KR"
)

// SupportedEncodings 支持的编码格式列表（用于 -f）
var SupportedEncodings = []string{
	EncodingAuto, EncodingUTF8, EncodingUTF8BOM, EncodingUTF16, EncodingUTF16LE, EncodingUTF16BE,
	EncodingGBK, EncodingGB2312, EncodingGB18030, EncodingBig5, EncodingLatin1, EncodingShiftJIS, EncodingEUCKR,
}

// TargetEncodings 目标编码格式列表（用于 -t，包含 none）
var TargetEncodings = []string{
	EncodingNone, EncodingUTF8, EncodingUTF8BOM, EncodingUTF16, EncodingUTF16LE, EncodingUTF16BE,
	EncodingGBK, EncodingGB2312, EncodingGB18030, EncodingBig5, EncodingLatin1, EncodingShiftJIS, EncodingEUCKR,
}

// 虚拟哈希表条目
type VirtualHashEntry struct {
	// 真实路径
	RealPath string

	// 哈希值
	Hash string
}

// 虚拟哈希表
type VirtualHashMap map[string]VirtualHashEntry

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
