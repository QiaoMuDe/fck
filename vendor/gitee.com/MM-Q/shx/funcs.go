package shx

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"mvdan.cc/sh/v3/syntax"
)

// Run 执行命令
//
// 参数：
//   - cmd: 命令字符串
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.Run("echo hello")
func Run(cmd string) error {
	return New(cmd).Exec()
}

// RunToTerminal 执行命令并输出到终端
//
// 参数：
//   - cmd: 命令字符串
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.RunToTerminal("echo hello")
func RunToTerminal(cmd string) error {
	return New(cmd).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Exec()
}

// Out 执行并获取输出
//
// 参数：
//   - cmd: 命令字符串
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	output, err := shx.Out("ls -la")
func Out(cmd string) ([]byte, error) {
	return New(cmd).ExecOutput()
}

// RunWith 超时执行
//
// 参数：
//   - cmd: 命令字符串
//   - timeout: 超时时间
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.RunWith("sleep 10", 5*time.Second)
func RunWith(cmd string, timeout time.Duration) error {
	return New(cmd).WithTimeout(timeout).Exec()
}

// OutWith 超时执行并获取输出
//
// 参数：
//   - cmd: 命令字符串
//   - timeout: 超时时间
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	output, err := shx.OutWith("sleep 5", 10*time.Second)
func OutWith(cmd string, timeout time.Duration) ([]byte, error) {
	return New(cmd).WithTimeout(timeout).ExecOutput()
}

// RunWithIO 使用自定义输入输出执行
//
// 参数：
//   - cmd: 命令字符串
//   - stdin: 标准输入
//   - stdout: 标准输出
//   - stderr: 标准错误
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	var buf bytes.Buffer
//	err := shx.RunWithIO("cat", strings.NewReader("hello"), &buf, os.Stderr)
func RunWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) error {
	return New(cmd).WithStdin(stdin).WithStdout(stdout).WithStderr(stderr).Exec()
}

// OutWithIO 使用自定义输入输出执行并获取输出
//
// 参数：
//   - cmd: 命令字符串
//   - stdin: 标准输入
//   - stdout: 标准输出
//   - stderr: 标准错误
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	var buf bytes.Buffer
//	output, err := shx.OutWithIO("cat", strings.NewReader("hello"), &buf, os.Stderr)
func OutWithIO(cmd string, stdin io.Reader, stdout, stderr io.Writer) ([]byte, error) {
	return New(cmd).WithStdin(stdin).WithStdout(stdout).WithStderr(stderr).ExecOutput()
}

// RunScript 执行 bash 脚本文件
//
// 参数：
//   - filePath: 脚本文件路径
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.RunScript("deploy.sh")
func RunScript(filePath string) error {
	return NewScript(filePath).Exec()
}

// RunScriptToTerminal 执行 bash 脚本文件并输出到终端
//
// 参数：
//   - filePath: 脚本文件路径
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.RunScriptToTerminal("deploy.sh")
func RunScriptToTerminal(filePath string) error {
	return NewScript(filePath).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Exec()
}

// OutScript 执行 bash 脚本文件并获取输出
//
// 参数：
//   - filePath: 脚本文件路径
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	output, err := shx.OutScript("deploy.sh")
func OutScript(filePath string) ([]byte, error) {
	return NewScript(filePath).ExecOutput()
}

// RunScriptWith 超时执行 bash 脚本文件
//
// 参数：
//   - filePath: 脚本文件路径
//   - timeout: 超时时间
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	err := shx.RunScriptWith("long_task.sh", 30*time.Second)
func RunScriptWith(filePath string, timeout time.Duration) error {
	return NewScript(filePath).WithTimeout(timeout).Exec()
}

// OutScriptWith 超时执行 bash 脚本文件并获取输出
//
// 参数：
//   - filePath: 脚本文件路径
//   - timeout: 超时时间
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	output, err := shx.OutScriptWith("build.sh", 60*time.Second)
func OutScriptWith(filePath string, timeout time.Duration) ([]byte, error) {
	return NewScript(filePath).WithTimeout(timeout).ExecOutput()
}

// RunScriptWithIO 使用自定义输入输出执行 bash 脚本文件
//
// 参数：
//   - filePath: 脚本文件路径
//   - stdin: 标准输入
//   - stdout: 标准输出
//   - stderr: 标准错误
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	var buf bytes.Buffer
//	err := shx.RunScriptWithIO("script.sh", strings.NewReader("input"), &buf, os.Stderr)
func RunScriptWithIO(filePath string, stdin io.Reader, stdout, stderr io.Writer) error {
	return NewScript(filePath).WithStdin(stdin).WithStdout(stdout).WithStderr(stderr).Exec()
}

// OutScriptWithIO 使用自定义输入输出执行 bash 脚本文件并获取输出
//
// 参数：
//   - filePath: 脚本文件路径
//   - stdin: 标准输入
//   - stdout: 标准输出
//   - stderr: 标准错误
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	var buf bytes.Buffer
//	output, err := shx.OutScriptWithIO("script.sh", strings.NewReader("input"), &buf, os.Stderr)
func OutScriptWithIO(filePath string, stdin io.Reader, stdout, stderr io.Writer) ([]byte, error) {
	return NewScript(filePath).WithStdin(stdin).WithStdout(stdout).WithStderr(stderr).ExecOutput()
}

// RunCtxScript 使用上下文执行 bash 脚本文件
//
// 参数：
//   - ctx: 上下文
//   - filePath: 脚本文件路径
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	err := shx.RunCtxScript(ctx, "long_task.sh")
func RunCtxScript(ctx context.Context, filePath string) error {
	return NewScript(filePath).WithContext(ctx).Exec()
}

// OutCtxScript 使用上下文执行 bash 脚本文件并获取输出
//
// 参数：
//   - ctx: 上下文
//   - filePath: 脚本文件路径
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	output, err := shx.OutCtxScript(ctx, "build.sh")
func OutCtxScript(ctx context.Context, filePath string) ([]byte, error) {
	return NewScript(filePath).WithContext(ctx).ExecOutput()
}

// RunCtx 使用上下文执行
//
// 参数：
//   - ctx: 上下文
//   - cmd: 命令字符串
//
// 返回：
//   - error: 执行错误
//
// 示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	err := shx.RunCtx(ctx, "sleep 10")
func RunCtx(ctx context.Context, cmd string) error {
	return New(cmd).WithContext(ctx).Exec()
}

// OutCtx 使用上下文执行并获取输出
//
// 参数：
//   - ctx: 上下文
//   - cmd: 命令字符串
//
// 返回：
//   - []byte: 命令输出
//   - error: 执行错误
//
// 示例：
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	output, err := shx.OutCtx(ctx, "ls -la")
func OutCtx(ctx context.Context, cmd string) ([]byte, error) {
	return New(cmd).WithContext(ctx).ExecOutput()
}

// CheckSyntax 检查 shell 命令字符串的语法
//
// 参数：
//   - script: shell 命令字符串
//
// 返回：
//   - error: 语法错误时返回错误信息
//
// 示例：
//
//	err := shx.CheckSyntax("echo hello")
//	if err != nil {
//	    fmt.Println(err)
//	}
func CheckSyntax(script string) error {
	parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
	_, err := parser.Parse(strings.NewReader(script), "")
	if err == nil {
		return nil
	}
	return convertSyntaxError(err)
}

// CheckScriptSyntax 检查 shell 脚本文件的语法
//
// 参数：
//   - filePath: 脚本文件路径
//
// 返回：
//   - error: 语法正确或文件不存在时返回错误信息
//
// 示例：
//
//	err := shx.CheckScriptSyntax("deploy.sh")
//	if err != nil {
//	    fmt.Println(err)
//	}
func CheckScriptSyntax(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return CheckSyntax(string(data))
}

// Format 格式化 shell 命令字符串
//
// 参数：
//   - script: shell 命令字符串
//
// 返回：
//   - string: 格式化后的字符串
//   - error: 解析错误或系统错误
//
// 示例：
//
//	formatted, err := shx.Format("for i in 1 2 3;do echo $i;done")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(formatted)
//
// 默认格式：
//   - 缩进: 2 空格
//   - case 语句缩进: 启用
func Format(script string) (string, error) {
	parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
	file, err := parser.Parse(strings.NewReader(script), "")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	printer := syntax.NewPrinter(
		syntax.Indent(2),              // 缩进 2 空格
		syntax.SwitchCaseIndent(true), // 启用 case 语句缩进
	)
	if err := printer.Print(&buf, file); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FormatScript 格式化 shell 脚本文件
//
// 参数：
//   - filePath: 脚本文件路径
//
// 返回：
//   - string: 格式化后的脚本内容
//   - error: 解析错误或系统错误
//
// 示例：
//
//	formatted, err := shx.FormatScript("deploy.sh")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(formatted)
func FormatScript(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return Format(string(data))
}

// convertSyntaxError 将 mvdan 的语法错误转为自定义 SyntaxError 返回
func convertSyntaxError(err error) error {
	// 尝试 *syntax.ParseError（指针类型）
	var parsePtr *syntax.ParseError
	if errors.As(err, &parsePtr) {
		return &SyntaxError{
			File:    parsePtr.Filename,
			Line:    int(parsePtr.Pos.Line()),
			Column:  int(parsePtr.Pos.Col()),
			Message: parsePtr.Text,
		}
	}

	// 尝试 syntax.ParseError（值类型）
	if parseVal, ok := err.(syntax.ParseError); ok {
		return &SyntaxError{
			File:    parseVal.Filename,
			Line:    int(parseVal.Pos.Line()),
			Column:  int(parseVal.Pos.Col()),
			Message: parseVal.Text,
		}
	}

	// 非语法错误，原样返回
	return err
}
