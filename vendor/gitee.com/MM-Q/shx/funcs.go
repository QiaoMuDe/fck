package shx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// FormatOptions 控制 shell 脚本格式化的行为选项
type FormatOptions struct {
	Indent           uint // 缩进空格数（0 表示使用 tab）
	SwitchCaseIndent bool // case 语句体是否缩进
	KeepComments     bool // 是否保留注释
	BinaryNextLine   bool // &&、|| 等二元操作符是否换行显示
	FunctionNextLine bool // 函数体 { 是否换行
	SpaceRedirects   bool // 重定向符前后是否加空格
	SingleLine       bool // 是否单行输出
	Minify           bool // 是否最小化输出（压缩模式）
}

// DefaultFormatOptions 返回默认格式化选项
//
// 返回：
//   - FormatOptions: 默认格式化选项
//
// 默认启用：
//   - 缩进: 4 空格
//   - case 语句缩进: 启用
//   - 注释保留: 启用
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Indent:           4,
		SwitchCaseIndent: true,
		KeepComments:     true,
	}
}

// FormatWithOptions 使用指定选项格式化 shell 命令字符串
//
// 参数：
//   - script: shell 命令字符串
//   - opts: 格式化选项
//
// 返回：
//   - string: 格式化后的字符串
//   - error: 解析错误或系统错误
//
// 示例：
//
//	formatted, err := shx.FormatWithOptions("for i in 1 2 3;do echo $i;done", shx.DefaultFormatOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(formatted)
func FormatWithOptions(script string, opts FormatOptions) (string, error) {
	// 构建解析器选项
	var parserOpts []syntax.ParserOption
	parserOpts = append(parserOpts, syntax.Variant(syntax.LangBash))
	if opts.KeepComments {
		parserOpts = append(parserOpts, syntax.KeepComments(true))
	}
	// 解析 shell 脚本为 AST
	parser := syntax.NewParser(parserOpts...)
	file, err := parser.Parse(strings.NewReader(script), "")
	if err != nil {
		return "", err
	}

	// 构建格式化器选项
	var printerOpts []syntax.PrinterOption
	if opts.Indent > 0 {
		// 添加缩进选项
		printerOpts = append(printerOpts, syntax.Indent(opts.Indent))
	}
	if opts.SwitchCaseIndent {
		// 添加 case 缩进选项
		printerOpts = append(printerOpts, syntax.SwitchCaseIndent(true))
	}
	if opts.BinaryNextLine {
		// 添加二元操作符换行选项
		printerOpts = append(printerOpts, syntax.BinaryNextLine(true))
	}
	if opts.FunctionNextLine {
		// 添加函数体 {} 换行选项
		printerOpts = append(printerOpts, syntax.FunctionNextLine(true))
	}
	if opts.SpaceRedirects {
		// 添加重定向符前后空格选项
		printerOpts = append(printerOpts, syntax.SpaceRedirects(true))
	}
	if opts.SingleLine {
		// 添加单行输出选项
		printerOpts = append(printerOpts, syntax.SingleLine(true))
	}
	if opts.Minify {
		// 添加最小化输出选项
		printerOpts = append(printerOpts, syntax.Minify(true))
	}

	// 使用构建好的选项创建格式化器并输出
	var buf bytes.Buffer
	printer := syntax.NewPrinter(printerOpts...)
	if err := printer.Print(&buf, file); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FormatScriptWithOptions 使用指定选项格式化 shell 脚本文件
//
// 参数：
//   - filePath: 脚本文件路径
//   - opts: 格式化选项
//
// 返回：
//   - string: 格式化后的脚本内容
//   - error: 解析错误或系统错误
//
// 示例：
//
//	formatted, err := shx.FormatScriptWithOptions("deploy.sh", shx.DefaultFormatOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(formatted)
func FormatScriptWithOptions(filePath string, opts FormatOptions) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return FormatWithOptions(string(data), opts)
}

// Format 使用默认选项格式化 shell 命令字符串
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
//   - 缩进: 4 空格
//   - case 语句缩进: 启用
//   - 注释保留: 启用
func Format(script string) (string, error) {
	return FormatWithOptions(script, DefaultFormatOptions())
}

// FormatScript 使用默认选项格式化 shell 脚本文件
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
	return FormatScriptWithOptions(filePath, DefaultFormatOptions())
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

// FindCmd 查找命令
//
// 增强版，在标准库 exec.LookPath 基础上增加了以下能力：
//   - 处理 Go 1.19+ 的 ErrDot 安全限制（当前目录程序）
//   - 返回绝对路径
//   - Windows 上检查可执行文件扩展名
//
// 参数:
//   - name: 命令名称
//
// 返回:
//   - string: 命令的绝对路径
//   - error: 错误信息
func FindCmd(name string) (string, error) {
	// 优先使用标准库 exec.LookPath 查找
	path, err := exec.LookPath(name)
	if err == nil {
		// 确保返回绝对路径
		if filepath.IsAbs(path) {
			return path, nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return abs, nil
	}

	// 处理 Go 1.19+ 的 ErrDot 错误（当前目录的程序）
	if errors.Is(err, exec.ErrDot) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		if isExecutable(abs) {
			return abs, nil
		}
		return "", fmt.Errorf("command %q in current directory is not executable", name)
	}

	// 其他错误（命令未找到等）直接返回
	return "", err
}

// FindCommandPath 查找单个命令的绝对路径
//
// 供其他包复用，只返回第一个匹配的路径
// 优先使用标准库 exec.LookPath, 处理 ErrDot 情况，找不到则返回空字符串
//
// 参数:
//   - name: 命令名称
//
// 返回:
//   - string: 命令的绝对路径，如果找不到则返回空字符串
func FindCommandPath(name string) string {
	path, err := FindCmd(name)
	if err != nil {
		return ""
	}
	return path
}

// windowsExts 定义 Windows 可执行文件扩展名集合
var windowsExts = map[string]bool{
	".exe": true,
	".bat": true,
	".cmd": true,
	".com": true,
}

// isExecutable 检查文件是否可执行
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - bool: 是否可执行
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Windows 上通过扩展名判断可执行性
	if runtime.GOOS == "windows" {
		ext := strings.ToLower(filepath.Ext(path))
		return windowsExts[ext] && !info.IsDir()
	}

	// Unix 类系统检查文件权限
	mode := info.Mode()
	return mode.IsRegular() && (mode.Perm()&0111 != 0)
}
