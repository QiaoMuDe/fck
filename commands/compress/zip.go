package compress

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// Zip 函数用于创建ZIP压缩文件
//
// 参数:
//   - zipFilePath: 生成的ZIP文件路径
//   - sourceDir: 需要压缩的源目录路径
//
// 返回值:
//   - error: 操作过程中遇到的错误
func Zip(zipFilePath string, sourceDir string) error {
	// 确保路径为绝对路径
	var absErr error
	if zipFilePath, absErr = ensureAbsPath(zipFilePath, "ZIP文件路径"); absErr != nil {
		return absErr
	}
	if sourceDir, absErr = ensureAbsPath(sourceDir, "源目录路径"); absErr != nil {
		return absErr
	}

	// 创建 ZIP 文件
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件失败: %w", err)
	}
	defer func() { _ = zipFile.Close() }()

	// 创建 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer func() { _ = zipWriter.Close() }()

	// 遍历目录并添加文件到 ZIP 包 (使用 WalkDir 提升性能)
	err = filepath.WalkDir(sourceDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			// 如果不存在则忽略
			if os.IsNotExist(err) {
				return nil
			}

			// 其他错误
			return fmt.Errorf("遍历目录时出错: %w", err)
		}

		// 获取相对路径，保留顶层目录
		headerName, err := filepath.Rel(filepath.Dir(sourceDir), path)
		if err != nil {
			return fmt.Errorf("获取相对路径失败: %w", err)
		}

		// 替换路径分隔符为正斜杠(ZIP 文件格式要求)
		headerName = filepath.ToSlash(headerName)

		// 根据文件类型处理
		switch {
		case entry.Type().IsRegular(): // 处理普通文件
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("获取文件信息失败: %w", err)
			}
			return processRegularFile(zipWriter, path, headerName, info)
		case entry.IsDir(): // 处理目录
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("获取目录信息失败: %w", err)
			}
			return processDirectory(zipWriter, headerName, info)
		case entry.Type()&fs.ModeSymlink != 0: // 处理符号链接
			return processSymlink(zipWriter, path, headerName, entry.Type())
		default: // 处理特殊文件
			return processSpecialFile(zipWriter, headerName, entry.Type())
		}
	})

	// 检查是否有错误发生
	if err != nil {
		return fmt.Errorf("打包目录到 ZIP 失败: %w", err)
	}

	return nil
}

// ensureAbsPath 确保路径为绝对路径，如果不是则转换为绝对路径
//
// 参数:
//   - path: 待检查的路径
//   - pathType: 路径类型描述（用于错误信息）
//
// 返回值:
//   - string: 绝对路径
//   - error: 转换过程中的错误
func ensureAbsPath(path, pathType string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("转换%s为绝对路径失败: %w", pathType, err)
	}
	return absPath, nil
}

// 缓冲区对象池，复用缓冲区减少内存分配
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024) // 默认32KB缓冲区
	},
}

// getBuffer 从对象池获取缓冲区
//
// 参数:
//   - size: 缓冲区大小
//
// 返回值:
//   - []byte: 获取到的缓冲区
func getBuffer(size int) []byte {
	buffer, ok := bufferPool.Get().([]byte)
	if !ok || len(buffer) < size {
		// 如果类型断言失败或池中的缓冲区太小，创建新的
		return make([]byte, size)
	}
	return buffer[:size]
}

// putBuffer 将缓冲区归还到对象池
//
// 参数:
//   - buffer: 要归还的缓冲区
//
// 说明:
//   - 该函数将缓冲区归还到对象池，以便后续复用。
//   - 只有容量不超过1MB的缓冲区才会被归还，以避免对象池占用过多内存。
func putBuffer(buffer []byte) {
	if cap(buffer) <= 1024*1024 { // 只回收不超过1MB的缓冲区
		//nolint:staticcheck // SA6002: 忽略装箱警告，对象池的性能收益远大于装箱开销
		bufferPool.Put(buffer)
	}
}

// processRegularFile 处理普通文件
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - path: string - 文件路径
//   - headerName: string - ZIP 文件中的文件名
//   - info: os.FileInfo - 文件信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processRegularFile(zipWriter *zip.Writer, path, headerName string, info os.FileInfo) error {
	// 创建文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件头失败: %w", err)
	}
	header.Name = headerName    // 设置文件名
	header.Method = zip.Deflate // 使用 Deflate 压缩算法

	// 创建 ZIP 写入器
	fileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("创建 ZIP 写入器失败: %w", err)
	}

	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}

	// 获取文件大小
	fileSize := info.Size()

	// 获取缓冲区大小并创建缓冲区
	bufferSize := getBufferSize(fileSize)
	buffer := getBuffer(bufferSize)
	defer putBuffer(buffer)

	// 复制文件内容到ZIP写入器
	_, err = io.CopyBuffer(fileWriter, file, buffer)

	// 立即关闭文件并检查错误
	closeErr := file.Close()
	if err != nil {
		return fmt.Errorf("写入 ZIP 文件失败: %w", err)
	}
	if closeErr != nil {
		return fmt.Errorf("关闭文件失败: %w", closeErr)
	}

	return nil
}

// processDirectory 处理目录
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - headerName: string - ZIP 文件中的目录名
//   - info: os.FileInfo - 目录信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processDirectory(zipWriter *zip.Writer, headerName string, info os.FileInfo) error {
	// 创建目录文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件头失败: %w", err)
	}
	// 设置目录名
	header.Name = headerName + "/" // 目录名后添加斜杠
	header.Method = zip.Store      // 使用不压缩的方法

	// 创建目录文件头
	if _, err := zipWriter.CreateHeader(header); err != nil {
		return fmt.Errorf("创建 ZIP 目录失败: %w", err)
	}
	return nil
}

// processSymlink 处理软链接
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - path: string - 软链接路径
//   - headerName: string - ZIP 文件中的软链接名
//   - mode: fs.FileMode - 文件模式
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSymlink(zipWriter *zip.Writer, path, headerName string, mode fs.FileMode) error {
	// 读取软链接目标
	target, err := os.Readlink(path)
	if err != nil {
		return fmt.Errorf("读取软链接目标失败: %w", err)
	}

	// 创建软链接文件头
	header := &zip.FileHeader{
		Name:   headerName,
		Method: zip.Store,
	}
	header.SetMode(mode)

	// 创建软链接文件
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("创建 ZIP 软链接失败: %w", err)
	}
	if _, err := writer.Write([]byte(target)); err != nil {
		return fmt.Errorf("写入软链接目标失败: %w", err)
	}
	return nil
}

// processSpecialFile 处理特殊文件类型
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - headerName: string - ZIP 文件中的特殊文件名
//   - mode: fs.FileMode - 文件模式
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSpecialFile(zipWriter *zip.Writer, headerName string, mode fs.FileMode) error {
	// 创建 ZIP 文件头
	header := &zip.FileHeader{
		Name:   headerName,
		Method: zip.Store,
	}
	header.SetMode(mode)

	// 创建 ZIP 文件写入器
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("创建 ZIP 特殊文件失败: %w", err)
	}
	if _, err := writer.Write([]byte{}); err != nil {
		return fmt.Errorf("写入特殊文件失败: %w", err)
	}
	return nil
}

// getBufferSize 根据文件大小动态设置缓冲区大小。该函数会根据传入的文件大小，
// 选择合适的缓冲区大小，以优化文件读写操作的性能。不同的文件大小范围对应不同的缓冲区大小。
//
// 参数:
//   - fileSize: 文件的大小，单位为字节，类型为 int64。
//
// 返回值:
//   - 缓冲区的大小，单位为字节，类型为 int。
func getBufferSize(fileSize int64) int {
	switch {
	// 当文件大小小于 512KB 时，设置缓冲区大小为 32KB
	case fileSize < 512*1024:
		return 32 * 1024
	// 当文件大小小于 1MB 时，设置缓冲区大小为 64KB
	case fileSize < 1*1024*1024:
		return 64 * 1024
	// 当文件大小小于 5MB 时，设置缓冲区大小为 128KB
	case fileSize < 5*1024*1024:
		return 128 * 1024
	// 当文件大小小于 10MB 时，设置缓冲区大小为 256KB
	case fileSize < 10*1024*1024:
		return 256 * 1024
	// 当文件大小小于 100MB 时，设置缓冲区大小为 512KB
	case fileSize < 100*1024*1024:
		return 512 * 1024
	// 当文件大小大于等于 100MB 时，设置缓冲区大小为 1MB
	default:
		return 1024 * 1024
	}
}
