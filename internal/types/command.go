package types

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
var FindTypeLimits = []string{
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

	// 只查找空文件或目录
	FindTypeEmpty,
	// 只查找空文件或目录-短参数
	FindTypeEmptyShort,

	// 只查找可执行文件
	FindTypeExecutable,
	// 只查找可执行文件-短参数
	FindTypeExecutableShort,

	// 只查找socket文件(套接字)
	FindTypeSocket,
	// 只查找socket文件-短参数
	FindTypeSocketShort,

	// 只查找管道文件
	FindTypePipe,
	// 只查找管道文件-短参数
	FindTypePipeShort,

	// 只查找块设备文件
	FindTypeBlock,
	// 只查找块设备文件-短参数
	FindTypeBlockShort,

	// 只查找字符设备文件
	FindTypeChar,
	// 只查找字符设备文件-短参数
	FindTypeCharShort,

	// 只查找追加模式的文件
	FindTypeAppend,
	// 只查找追加模式的文件-短参数
	FindTypeAppendShort,

	// 只查找非追加模式的文件
	FindTypeNonAppend,
	// 只查找非追加模式的文件-短参数
	FindTypeNonAppendShort,

	// 只查找为独占模式的文件
	FindTypeExclusive,
	// 只查找为独占模式的文件-短参数
	FindTypeExclusiveShort,
}

// 限制查找的参数切片
var ListTypeLimits = []string{
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

// DNS 查询类型常量定义
const (
	// DNS A 记录类型
	DNSTypeA = "a"
	// DNS AAAA 记录类型
	DNSTypeAAAA = "aaaa"
	// DNS MX 记录类型
	DNSTypeMX = "mx"
	// DNS NS 记录类型
	DNSTypeNS = "ns"
	// DNS TXT 记录类型
	DNSTypeTXT = "txt"
	// DNS CNAME 记录类型
	DNSTypeCNAME = "cname"
	// DNS PTR 记录类型
	DNSTypePTR = "ptr"
	// DNS SOA 记录类型
	DNSTypeSOA = "soa"
	// DNS SRV 记录类型
	DNSTypeSRV = "srv"
	// DNS ANY 记录类型
	DNSTypeANY = "any"
)

// DNSQueryTypes DNS 查询类型切片
var DNSQueryTypes = []string{
	DNSTypeA,
	DNSTypeAAAA,
	DNSTypeMX,
	DNSTypeNS,
	DNSTypeTXT,
	DNSTypeCNAME,
	DNSTypePTR,
	DNSTypeSOA,
	DNSTypeSRV,
	DNSTypeANY,
}

// GetDNSAnyTypes 获取 ANY 查询包含的所有记录类型（排除 ANY 本身）
func GetDNSAnyTypes() []string {
	types := make([]string, 0, len(DNSQueryTypes)-1)
	for _, t := range DNSQueryTypes {
		if t != DNSTypeANY {
			types = append(types, t)
		}
	}
	return types
}

// TCP 扫描输出格式常量
const (
	// TCPScanFormatDefault 默认简洁格式
	TCPScanFormatDefault = "def"
	// TCPScanFormatTable 表格格式
	TCPScanFormatTable = "table"
	// TCPScanFormatJSON JSON格式
	TCPScanFormatJSON = "json"
	// TCPScanFormatCSV CSV格式
	TCPScanFormatCSV = "csv"
)

// TCPScanFormatOptions TCP扫描输出格式选项列表
var TCPScanFormatOptions = []string{
	TCPScanFormatDefault,
	TCPScanFormatTable,
	TCPScanFormatJSON,
	TCPScanFormatCSV,
}
