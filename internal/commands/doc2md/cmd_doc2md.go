package doc2md

import (
	"bytes"
	"fmt"
	"io"
	"os"

	markitdown "gitee.com/MM-Q/doc2md"
	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/go-kit/term"
)

// Doc2mdCmdMain 执行 doc2md 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func Doc2mdCmdMain(config Doc2mdConfig) error {
	// 创建转换引擎并应用选项
	engine := buildEngine(config.KeepDataURIs)

	// 优先判断 stdin 管道输入（与 fck 其他命令的管道范式一致）
	if term.IsStdinPipe() {
		return convertStdin(engine, config)
	}

	// 无管道输入时，处理文件
	if len(config.Files) == 0 {
		return fmt.Errorf("no input file specified, and no piped stdin detected")
	}

	// 设置输出文件时一次只能转换一个输入文件，避免静默覆盖
	if config.Output != "" && len(config.Files) > 1 {
		return fmt.Errorf("only one input file can be converted at a time when using -o/--output")
	}

	// 展开通配符并逐个转换
	allFiles, err := fs.Expand(config.Files)
	if err != nil {
		return err
	}
	if len(allFiles) == 0 {
		return fmt.Errorf("no matching input files")
	}

	for _, file := range allFiles {
		result, err := engine.ConvertFile(file)
		if err != nil {
			return formatConvertError(file, err)
		}
		if err := writeOutput(result, config.Output); err != nil {
			return err
		}
	}

	return nil
}

// buildEngine 创建 doc2md 转换引擎
//
// 参数:
//   - keepDataURIs: 是否保留完整 base64 图片数据 URI
//
// 返回值:
//   - *markitdown.MarkItDown: 转换引擎实例
func buildEngine(keepDataURIs bool) *markitdown.MarkItDown {
	if keepDataURIs {
		return markitdown.New(markitdown.WithKeepDataURIs(true))
	}
	return markitdown.New()
}

// convertStdin 处理 stdin 管道输入
//
// 参数:
//   - engine: 转换引擎
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func convertStdin(engine *markitdown.MarkItDown, config Doc2mdConfig) error {
	// 管道输入必须显式提供扩展名提示，否则无法识别格式
	if config.Extension == DocExtNone {
		return fmt.Errorf("stdin piped input requires -x/--extension to specify the document extension (e.g.: -x .docx)")
	}

	// 读取 stdin 全部内容到内存（ConvertReader 要求 io.ReadSeeker 且转换器会 Seek 重置）
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read standard input: %w", err)
	}
	reader := bytes.NewReader(data)

	info := markitdown.StreamInfo{
		Extension: config.Extension,
		MIMEType:  config.MIMEType,
		Charset:   config.Charset,
	}

	result, err := engine.ConvertReader(reader, info)
	if err != nil {
		return formatConvertError("<stdin>", err)
	}

	return writeOutput(result, config.Output)
}

// writeOutput 输出转换结果，未指定输出文件时打印到 stdout
//
// 参数:
//   - result: 转换结果
//   - output: 输出文件路径（为空时打印到 stdout）
//
// 返回值:
//   - error: 输出错误
func writeOutput(result *markitdown.DocumentConverterResult, output string) error {
	if output == "" {
		fmt.Println(result.Markdown)
		return nil
	}
	return os.WriteFile(output, []byte(result.Markdown), 0o644)
}

// formatConvertError 将 doc2md 错误转换为可读的英文错误信息
//
// 参数:
//   - source: 输入来源描述（文件路径或 <stdin>）
//   - err: 原始错误
//
// 返回值:
//   - error: 格式化后的错误
func formatConvertError(source string, err error) error {
	if markitdown.IsUnsupportedFormat(err) {
		return fmt.Errorf("%s: unsupported document format (%v)", source, err)
	}
	return fmt.Errorf("%s: conversion failed: %w", source, err)
}
