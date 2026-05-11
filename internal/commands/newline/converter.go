// Package newline 提供文件换行符检测和转换功能
package newline

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/go-kit/fs"
)

// ConversionResult 转换结果
type ConversionResult struct {
	FilePath string
	FromType string
	ToType   string
	Success  bool
	Error    error
	Lines    int // 文件行数
}

// Config 换行符转换配置
type Config struct {
	Files     []string
	ToNewline string
	Write     bool
	Backup    bool
	Output    string
	Quiet     bool
	FileCount int // 文件总数（用于控制输出格式）
}

// ConvertFile 转换单个文件的换行符
//
// 参数:
//   - srcPath: 源文件路径
//   - config: 转换配置
//
// 返回值:
//   - ConversionResult: 转换结果
//   - error: 执行错误（如果有）
func ConvertFile(srcPath string, config Config) (ConversionResult, error) {
	result := ConversionResult{FilePath: srcPath}

	// 1. 读取文件
	data, err := os.ReadFile(srcPath)
	if err != nil {
		result.Error = err
		return result, err
	}

	// 2. 检测当前换行符
	detection := Detect(data)
	result.FromType = detection.Type

	// 3. 目标为 none：只检测
	if config.ToNewline == types.NewlineNone {
		result.Success = true
		result.Lines = detection.TotalLines()
		if !config.Quiet {
			// 单文件只输出类型，多文件输出路径:类型 (行号)
			if config.FileCount == 1 {
				fmt.Println(detection.Type)
			} else {
				fmt.Printf("%s: %s\n", srcPath, detection.String())
			}
		}
		return result, nil
	}

	// 4. 确定目标类型
	targetType := strings.ToUpper(config.ToNewline)
	result.ToType = targetType

	// 5. 已经是目标格式
	if detection.Type == targetType {
		result.Success = true
		result.Lines = detection.TotalLines()
		if !config.Quiet {
			fmt.Printf("%s: %s (no conversion needed)\n", srcPath, detection.Type)
		}
		return result, nil
	}

	// 6. 执行转换
	converted := Convert(data, config.ToNewline)

	// 7. 确定目标路径
	var dstPath string
	if config.Write {
		// 原地写入：使用临时文件
		dstPath = srcPath + types.NewlineTempExt
	} else {
		// 非原地写入
		if config.Output == "." {
			// 默认：当前目录下生成 原文件名.unix 或 .win
			dstPath = srcPath + getConvertedExt(config.ToNewline)
		} else {
			// 检查 -o 指定的是目录还是文件
			outputInfo, err := os.Stat(config.Output)
			if err == nil && outputInfo.IsDir() {
				// 是目录：目录下生成 原文件名.unix/.win
				baseName := filepath.Base(srcPath)
				dstPath = filepath.Join(config.Output, baseName+getConvertedExt(config.ToNewline))
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

	// 8. 写入转换后的内容
	err = os.WriteFile(dstPath, converted, 0644)
	if err != nil {
		result.Error = err
		return result, err
	}

	// 9. 处理原地写入
	if config.Write {
		// 备份原文件，重命名临时文件
		if config.Backup {
			backupPath := srcPath + types.NewlineBackupExt
			_ = fs.MoveEx(srcPath, backupPath, true)
		} else {
			_ = os.Remove(srcPath)
		}
		err = fs.MoveEx(dstPath, srcPath, true)
	}

	result.Success = err == nil
	if result.Success {
		result.Lines = detection.TotalLines()
	}
	return result, err
}

// Convert 转换换行符
//
// 参数:
//   - data: 原始数据
//   - targetType: 目标换行符类型 (lf/crlf/cr)
//
// 返回值:
//   - []byte: 转换后的数据
func Convert(data []byte, targetType string) []byte {
	// 先统一转换为 LF
	normalized := normalizeToLF(data)

	switch targetType {
	case types.NewlineLF:
		return normalized
	case types.NewlineCRLF:
		return bytes.ReplaceAll(normalized, []byte{'\n'}, []byte("\r\n"))
	case types.NewlineCR:
		return bytes.ReplaceAll(normalized, []byte{'\n'}, []byte{'\r'})
	default:
		return data
	}
}

// normalizeToLF 将所有换行符统一为 LF
func normalizeToLF(data []byte) []byte {
	// 先将 CRLF 转为 LF
	result := bytes.ReplaceAll(data, []byte("\r\n"), []byte{'\n'})
	// 再将 CR 转为 LF
	result = bytes.ReplaceAll(result, []byte{'\r'}, []byte{'\n'})
	return result
}

// getConvertedExt 根据目标类型返回转换后的文件后缀
func getConvertedExt(targetType string) string {
	switch targetType {
	case types.NewlineLF:
		return ".unix"
	case types.NewlineCRLF:
		return ".win"
	case types.NewlineCR:
		return ".mac"
	default:
		return ".converted"
	}
}
