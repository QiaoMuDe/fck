package cat

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gitee.com/MM-Q/fck/internal/utils"
	goKitUtils "gitee.com/MM-Q/go-kit/utils"
)

// ContentSource 内容源接口
type ContentSource interface {
	Name() string          // 名称（文件名或 "stdin"）
	Read() ([]byte, error) // 读取内容
	Size() (int64, error)  // 获取大小（用于检查限制）
	IsBinary() (bool, error)
}

// FileSource 文件内容源
type FileSource struct {
	path    string
	maxSize int64 // 最大大小限制，0 表示无限制
}

// NewFileSource 创建文件内容源
//
// 参数:
//   - path: 文件路径
//   - maxSize: 最大大小限制（字节），0 表示无限制
//
// 返回:
//   - *FileSource: 文件内容源实例
func NewFileSource(path string, maxSize int64) *FileSource {
	return &FileSource{path: path, maxSize: maxSize}
}

// Name 返回文件名
//
// 返回:
//   - string: 文件路径
func (f *FileSource) Name() string {
	return f.path
}

// Read 读取文件内容（带大小限制）
//
// 返回:
//   - []byte: 文件内容
//   - error: 读取错误（包括大小超限错误）
func (f *FileSource) Read() ([]byte, error) {
	if f.maxSize > 0 {
		// 先检查文件大小
		size, err := f.Size()
		if err != nil {
			return nil, err
		}
		if size > f.maxSize {
			return nil, fmt.Errorf("content too large (%s), use -l flag or increase limit with -S flag", goKitUtils.FormatBytes(size))
		}
	}
	return os.ReadFile(f.path)
}

// Size 获取文件大小
//
// 返回:
//   - int64: 文件大小（字节）
//   - error: 获取大小错误
func (f *FileSource) Size() (int64, error) {
	info, err := os.Stat(f.path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// IsBinary 检测文件是否为二进制
//
// 返回:
//   - bool: 是否为二进制文件
//   - error: 检测错误
func (f *FileSource) IsBinary() (bool, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = file.Close()
	}()
	return utils.IsBinaryFile(file)
}

// StdinSource 管道内容源
type StdinSource struct {
	name    string
	maxSize int64 // 最大大小限制，0 表示无限制
}

// NewStdinSource 创建管道内容源
//
// 参数:
//   - name: 名称（为空则使用 "stdin"）
//   - maxSize: 最大大小限制（字节），0 表示无限制
//
// 返回:
//   - *StdinSource: 管道内容源实例
func NewStdinSource(name string, maxSize int64) *StdinSource {
	if name == "" {
		name = "stdin"
	}
	return &StdinSource{name: name, maxSize: maxSize}
}

// Name 返回名称
//
// 返回:
//   - string: 名称
func (s *StdinSource) Name() string {
	return s.name
}

// Read 从 stdin 读取内容（带大小限制）
//
// 返回:
//   - []byte: 读取的内容
//   - error: 读取错误（包括大小超限错误）
func (s *StdinSource) Read() ([]byte, error) {
	if s.maxSize > 0 {
		return readWithLimit(os.Stdin, s.maxSize)
	}
	return io.ReadAll(os.Stdin)
}

// readWithLimit 带大小限制的读取
//
// 参数:
//   - r: 读取器
//   - limit: 大小限制（字节）
//
// 返回:
//   - []byte: 读取的内容
//   - error: 读取错误（包括大小超限错误）
func readWithLimit(r io.Reader, limit int64) ([]byte, error) {
	var buf bytes.Buffer
	written, err := io.CopyN(&buf, r, limit+1)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if written > limit {
		return nil, fmt.Errorf("content too large (%s), use -l flag or increase limit with -S flag", goKitUtils.FormatBytes(limit))
	}
	return buf.Bytes(), nil
}

// Size 获取已读取内容的大小（需要先读取）
//
// 返回:
//   - int64: 始终返回 -1（管道大小未知）
//   - error: 始终返回 nil
func (s *StdinSource) Size() (int64, error) {
	return -1, nil
}

// IsBinary 检测内容是否为二进制（简单检测）
//
// 返回:
//   - bool: 始终返回 false（管道输入默认不视为二进制）
//   - error: 始终返回 nil
func (s *StdinSource) IsBinary() (bool, error) {
	// 管道输入默认不视为二进制，除非显式检测到
	return false, nil
}
