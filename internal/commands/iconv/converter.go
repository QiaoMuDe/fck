// Package iconv 提供文件编码检测和转换功能
package iconv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/fs"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// ConvertFile 转换单个文件编码
//
// 参数:
//   - srcPath: 源文件路径
//   - config: 转换配置
//
// 返回值:
//   - ConversionResult: 转换结果
//   - error: 执行错误（如果有）
func ConvertFile(srcPath string, config IconvConfig) (ConversionResult, error) {
	result := ConversionResult{FilePath: srcPath}

	// 1. 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		result.Error = err
		return result, err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	// 2. 读取样本用于编码检测
	sample := make([]byte, types.GrepEncodingCheckSize)
	n, err := srcFile.Read(sample)
	if err != nil && err != io.EOF {
		result.Error = err
		return result, err
	}
	sample = sample[:n]
	_, err = srcFile.Seek(0, 0)
	if err != nil {
		result.Error = err
		return result, err
	}

	// 3. 确定源编码
	fromEncoding := config.FromEncoding
	if fromEncoding == types.EncodingAuto {
		detected, _, err := DetectEncoding(sample)
		if err != nil || detected == "" {
			// 检测失败，使用默认编码 UTF-8
			detected = types.EncodingUTF8
			//confidence = 0
		}
		fromEncoding = detected            // 检测到的编码
		result.FromEncoding = fromEncoding // 记录检测到的编码

		// 目标编码为 none：只检测并打印
		if config.ToEncoding == types.EncodingNone {
			result.Success = true
			if !config.Quiet {
				fmt.Printf("%s: %s\n", srcPath, fromEncoding)
			}
			return result, nil
		}
	}
	result.FromEncoding = fromEncoding

	// 4. 如果源编码和目标编码相同，跳过
	if strings.EqualFold(fromEncoding, config.ToEncoding) {
		result.Success = true
		result.ToEncoding = config.ToEncoding
		if !config.Quiet {
			fmt.Printf("%s: %s (no conversion needed)\n", srcPath, fromEncoding)
		}
		return result, nil
	}

	// 5. 确定目标路径
	var dstPath string
	if config.Write {
		// 原地写入：使用临时文件，后续重命名
		dstPath = srcPath + types.IconvTempExt
	} else {
		// 非原地写入
		if config.Output == "." {
			// 默认：当前目录下生成 原文件名.converted
			dstPath = srcPath + types.IconvConvertedExt
		} else {
			// 检查 -o 指定的是目录还是文件
			outputInfo, err := os.Stat(config.Output)
			if err == nil && outputInfo.IsDir() {
				// 是目录：目录下生成 原文件名.converted
				baseName := filepath.Base(srcPath)
				dstPath = filepath.Join(config.Output, baseName+types.IconvConvertedExt)
			} else if err == nil {
				// 是文件路径：直接写入指定路径
				dstPath = config.Output
			} else {
				// 路径不存在，报错
				result.Error = fmt.Errorf("output path does not exist: %s", config.Output)
				return result, result.Error
			}
		}
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		result.Error = err
		return result, err
	}

	// 6. 执行转换
	bytesRead, bytesWritten, err := doConvert(srcFile, dstFile, fromEncoding, config.ToEncoding)
	_ = dstFile.Close()

	result.BytesRead = bytesRead
	result.BytesWritten = bytesWritten

	if err != nil {
		_ = os.Remove(dstPath)
		result.Error = err
		return result, err
	}

	// 7. 处理输出
	if config.Write {
		// 原地写入：先关闭源文件（Windows 需要）
		_ = srcFile.Close()

		// 备份原文件，重命名临时文件
		if config.Backup {
			backupPath := srcPath + types.IconvBackupExt
			_ = fs.MoveEx(srcPath, backupPath, true)
		} else {
			_ = os.Remove(srcPath)
		}
		err = fs.MoveEx(dstPath, srcPath, true)
	}

	result.Success = err == nil
	result.ToEncoding = config.ToEncoding
	return result, err
}

// doConvert 执行实际的编码转换
//
// 参数:
//   - src: 源数据读取器
//   - dst: 目标数据写入器
//   - from: 源编码
//   - to: 目标编码
//
// 返回值:
//   - int64: 读取字节数
//   - int64: 写入字节数
//   - error: 转换错误（如果有）
func doConvert(src io.Reader, dst io.Writer, from, to string) (int64, int64, error) {
	// 获取解码器和编码器
	decoder := getEncoding(from).NewDecoder()
	encoder := getEncoding(to).NewEncoder()

	// 创建转换链：src → decoder → encoder → dst
	reader := transform.NewReader(src, decoder)
	writer := transform.NewWriter(dst, encoder)

	// 执行转换
	bytesRead, err := io.Copy(writer, reader)
	if err != nil {
		return bytesRead, 0, err
	}

	// 刷新编码器缓冲区
	err = writer.Close()

	// 获取实际写入字节数（需要额外统计）
	// 简化处理：返回读取字节数作为参考
	return bytesRead, bytesRead, err
}

// getEncoding 根据名称获取编码
//
// 参数:
//   - name: 编码名称
//
// 返回值:
//   - encoding.Encoding: 编码对象
func getEncoding(name string) encoding.Encoding {
	switch strings.ToUpper(name) {
	case types.EncodingUTF8:
		return encoding.Nop // UTF-8 无需转换
	case types.EncodingUTF8BOM:
		return unicode.UTF8BOM
	case types.EncodingUTF16:
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case types.EncodingUTF16LE:
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case types.EncodingUTF16BE:
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case types.EncodingGBK:
		return simplifiedchinese.GBK
	case types.EncodingGB2312:
		return simplifiedchinese.GB18030 // GB18030 兼容 GB2312
	case types.EncodingGB18030:
		return simplifiedchinese.GB18030
	case types.EncodingBig5:
		return traditionalchinese.Big5
	case types.EncodingLatin1:
		return charmap.ISO8859_1
	case types.EncodingShiftJIS:
		return japanese.ShiftJIS
	case types.EncodingEUCKR:
		return korean.EUCKR
	default:
		return encoding.Nop
	}
}
