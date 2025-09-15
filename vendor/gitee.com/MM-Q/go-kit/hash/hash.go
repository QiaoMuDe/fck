package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/go-kit/pool"
	"github.com/schollz/progressbar/v3"
)

// 支持的哈希算法列表
var supportedAlgorithms = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}

// IsAlgorithmSupported 检查给定的哈希算法名称是否受支持。
// 匹配时会忽略算法名称的大小写。
//
// 参数:
//   - algorithm: 要检查的哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）。
//
// 返回:
//   - bool: 如果算法受支持则返回 true，否则返回 false。
func IsAlgorithmSupported(algorithm string) bool {
	// 如果算法名称为空，则返回 false
	if algorithm == "" {
		return false
	}

	_, ok := supportedAlgorithms[strings.ToLower(algorithm)]
	return ok
}

// getHashAlgorithm 根据算法名称获取对应的哈希函数构造器。
// 匹配时会忽略算法名称的大小写。
//
// 参数:
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）。
//
// 返回:
//   - func() hash.Hash: 对应的哈希函数构造器。
//   - error: 如果不支持该算法，则返回错误。
func getHashAlgorithm(algorithm string) (func() hash.Hash, error) {
	// 如果算法名称为空，则返回错误
	if algorithm == "" {
		return nil, fmt.Errorf("hash algorithm name cannot be empty")
	}

	algoFunc, ok := supportedAlgorithms[strings.ToLower(algorithm)]
	if !ok {
		return nil, fmt.Errorf("unsupported hash algorithm: %s", algorithm)
	}
	return algoFunc, nil
}

// checksumCore 核心哈希计算逻辑，支持可选的进度条显示
//
// 参数:
//   - filePath: 文件路径
//   - algorithm: 哈希算法名称
//   - showProgress: 是否显示进度条
//
// 返回:
//   - string: 文件的十六进制哈希值
//   - error: 错误信息，如果计算失败
func checksumCore(filePath, algorithm string, showProgress bool) (string, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("file does not exist or is inaccessible: %v", err)
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// 获取哈希函数构造器
	hashFunc, err := getHashAlgorithm(algorithm)
	if err != nil {
		return "", err
	}
	h := hashFunc()

	// 根据文件大小动态分配缓冲区，确保最小为1KB
	fileSize := fileInfo.Size()
	bufferSize := pool.CalculateBufferSize(fileSize)
	if bufferSize < pool.KB {
		bufferSize = pool.KB
	}
	buf := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buf) // 使用完毕后归还到对象池

	// 默认写入器为哈希函数
	var writer io.Writer = h

	// 如果需要显示进度条，则创建进度条
	if showProgress {
		bar := progressbar.NewOptions64(
			fileSize,                          // 进度条总长度
			progressbar.OptionClearOnFinish(), // 结束时清除进度条
			progressbar.OptionSetDescription(fmt.Sprintf("正在处理'%s'('%s')", filepath.Base(filePath), strings.ToUpper(algorithm))), // 显示描述
			progressbar.OptionSetElapsedTime(true),      // 显示已用时间
			progressbar.OptionSetPredictTime(true),      // 显示预计剩余时间
			progressbar.OptionSetRenderBlankState(true), // 在进度条完成之前显示空白状态
			progressbar.OptionShowBytes(true),           // 显示进度条传输的字节
			progressbar.OptionShowCount(),               // 显示当前进度的总和
			//progressbar.OptionShowElapsedTimeOnFinish(),        // 完成后显示已用时间
			progressbar.OptionSetTheme(progressbar.ThemeASCII), // ASCII 进度条主题(默认为 Unicode 进度条主题)
			progressbar.OptionFullWidth(),                      // 设置进度条为终端最大宽度
		)
		defer func() {
			_ = bar.Finish() // 完成进度条
			_ = bar.Close()  // 关闭进度条
		}()
		writer = io.MultiWriter(h, bar)
	}

	// 使用 io.CopyBuffer 进行高效复制并计算哈希
	if _, err := io.CopyBuffer(writer, file, buf); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Checksum 计算文件哈希值
//
// 参数:
//   - filePath: 文件路径
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）
//
// 返回:
//   - string: 文件的十六进制哈希值
//   - error: 错误信息，如果计算失败
//
// 注意:
//   - 根据文件大小动态分配缓冲区以提高性能
//   - 支持任何实现hash.Hash接口的哈希算法
//   - 使用io.CopyBuffer进行高效的文件读取和哈希计算
func Checksum(filePath string, algorithm string) (string, error) {
	return checksumCore(filePath, algorithm, false)
}

// ChecksumProgress 计算文件哈希值(带进度条)
//
// 参数:
//   - filePath: 文件路径
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）
//
// 返回:
//   - string: 文件的十六进制哈希值
//   - error: 错误信息，如果计算失败
//
// 注意:
//   - 根据文件大小动态分配缓冲区以提高性能
//   - 支持任何实现hash.Hash接口的哈希算法
//   - 使用io.CopyBuffer进行高效的文件读取和哈希计算
func ChecksumProgress(filePath string, algorithm string) (string, error) {
	return checksumCore(filePath, algorithm, true)
}

// HashData 计算内存数据哈希值
//
// 参数:
//   - data: 要计算哈希的字节数据
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）
//
// 返回:
//   - string: 数据的十六进制哈希值
//   - error: 错误信息，如果计算失败
//
// 注意:
//   - 直接在内存中计算，无需文件I/O，性能更高
//   - 支持任何大小的数据，包括空数据
//   - 使用标准库优化的hash实现，性能最佳
func HashData(data []byte, algorithm string) (string, error) {
	// 参数验证
	if data == nil {
		return "", fmt.Errorf("data cannot be nil")
	}

	// 获取哈希函数构造器
	hashFunc, err := getHashAlgorithm(algorithm)
	if err != nil {
		return "", err
	}
	h := hashFunc()

	// 直接写入所有数据 - 标准库已经优化过了
	if _, err := h.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data to hash: %v", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashString 计算字符串哈希值（便利函数）
//
// 参数:
//   - data: 要计算哈希的字符串
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）
//
// 返回:
//   - string: 字符串的十六进制哈希值
//   - error: 错误信息，如果计算失败
//
// 注意:
//   - 这是HashData的便利包装函数
//   - 内部将字符串转换为字节切片进行处理
//   - 适用于文本数据、配置字符串等场景
func HashString(data string, algorithm string) (string, error) {
	return HashData([]byte(data), algorithm)
}

// HashReader 计算io.Reader数据哈希值
//
// 参数:
//   - reader: 数据源读取器
//   - algorithm: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）
//
// 返回:
//   - string: 读取数据的十六进制哈希值
//   - error: 错误信息，如果计算失败
//
// 注意:
//   - 适用于流式数据处理，如网络数据、管道数据等
//   - 使用缓冲区进行高效读取，避免频繁的小块读取
//   - 会完全消费Reader中的数据
//   - 使用对象池优化内存分配
func HashReader(reader io.Reader, algorithm string) (string, error) {
	// 参数验证
	if reader == nil {
		return "", fmt.Errorf("reader cannot be nil")
	}

	// 获取哈希函数构造器
	hashFunc, err := getHashAlgorithm(algorithm)
	if err != nil {
		return "", err
	}
	h := hashFunc()

	// 从对象池获取缓冲区进行高效读取
	const bufferSize = 32 * 1024 // 32KB缓冲区，平衡内存使用和I/O效率
	buf := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buf)

	// 使用io.CopyBuffer进行高效复制和哈希计算
	if _, err := io.CopyBuffer(h, reader, buf); err != nil {
		return "", fmt.Errorf("failed to read data from reader: %v", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
