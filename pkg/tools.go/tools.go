package tools

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

// 字节单位定义
const (
	Byte = 1 << (10 * iota) // 1 字节
	KB                      // 千字节 (1024 B)
	MB                      // 兆字节 (1024 KB)
	GB                      // 吉字节 (1024 MB)
	TB                      // 太字节 (1024 GB)
)

// 计算文件哈希值的函数
func Checksum(filePath string, hashFunc func() hash.Hash) (string, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("文件不存在或无法访问: %v", err)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建哈希对象
	hash := hashFunc()

	// 根据文件大小动态分配缓冲区
	fileSize := fileInfo.Size()
	bufferSize := calculateBufferSize(fileSize)
	buffer := make([]byte, bufferSize)

	// 使用 io.CopyBuffer 进行高效复制并计算哈希
	if _, err := io.CopyBuffer(hash, file, buffer); err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 返回哈希值的十六进制表示
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 根据文件大小计算最佳缓冲区大小
func calculateBufferSize(fileSize int64) int {
	switch {
	case fileSize < 32*KB: // 小于 32KB 的文件使用 32KB 缓冲区
		return int(32 * KB)
	case fileSize < 128*KB: // 32KB-128KB 使用 64KB 缓冲区
		return int(32 * KB)
	case fileSize < 512*KB: // 128KB-512KB 使用 128KB 缓冲区
		return int(64 * KB)
	case fileSize < 1*MB: // 512KB-1MB 使用 256KB 缓冲区
		return int(128 * KB)
	case fileSize < 4*MB: // 1MB-4MB 使用 512KB 缓冲区
		return int(256 * KB)
	case fileSize < 16*MB: // 4MB-16MB 使用 1MB 缓冲区
		return int(512 * KB)
	case fileSize < 64*MB: // 16MB-64MB 使用 2MB 缓冲区
		return int(1 * MB)
	default: // 大于 64MB 的文件使用 4MB 缓冲区
		return int(2 * MB)
	}
}
