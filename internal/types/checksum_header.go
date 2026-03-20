package types

import (
	"fmt"
)

// ChecksumHeader 校验文件头信息结构体
type ChecksumHeader struct {
	HashType  string // 哈希类型 (md5, sha1, sha256等)
	Timestamp string // 生成时间戳
	Mode      string // 模式 (PORTABLE/LOCAL)
	BasePath  string // 基准路径 (仅LOCAL模式下使用)
}

// String 生成文件头字符串
func (h *ChecksumHeader) String() string {
	if h.Mode == ChecksumModeLocal && h.BasePath != "" {
		return fmt.Sprintf("#%s#%s#%s#%s\n", h.HashType, h.Timestamp, h.Mode, h.BasePath)
	}
	return fmt.Sprintf("#%s#%s#%s\n", h.HashType, h.Timestamp, h.Mode)
}

// IsPortableMode 判断是否为便携模式
func (h *ChecksumHeader) IsPortableMode() bool {
	return h.Mode == ChecksumModePortable || h.Mode == ""
}

// IsLocalMode 判断是否为本地模式
func (h *ChecksumHeader) IsLocalMode() bool {
	return h.Mode == ChecksumModeLocal
}
