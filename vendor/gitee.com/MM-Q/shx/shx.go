// Package shx 提供了基于 mvdan.cc/sh/v3 的纯 Go shell 命令执行功能。
//
// 它使用 mvdan.cc/sh/v3 进行命令解析和执行, 不依赖系统 shell。
//
// 主要特性:
//   - 纯 Go 实现, 不依赖系统 shell
//   - 更好的跨平台一致性 (Windows/Linux/macOS 行为一致)
//   - 链式调用 API, 支持流畅的方法链
//   - 支持超时控制和上下文取消
//   - 支持执行 .sh 脚本文件
//
// 基本用法:
//
//	import "gitee.com/MM-Q/shx"
//
//	// 简单执行
//	err := shx.Run("echo hello world")
//
//	// 获取输出
//	output, err := shx.Out("ls -la")
//
//	// 链式配置
//	output, err := shx.New("echo hello").
//		WithTimeout(5 * time.Second).
//		WithDir("/tmp").
//		ExecOutput()
//
//	// 使用上下文
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	err := shx.New("long-running-command").WithContext(ctx).Exec()
//
//	// 执行脚本文件
//	err = shx.RunScript("deploy.sh")
//
// 注意事项:
//   - Shx 对象的配置方法 (WithXxx) 不是并发安全的, 不要在多个 goroutine 中并发配置
//   - mvdan/sh 是同步执行的, 不提供异步 API, 如需异步请使用 goroutine 包装
//   - 不支持进程控制 (无 PID、Kill、Signal) , 只能通过 context 取消
package shx

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// Shx 表示一个待执行的 shell 命令或脚本文件
//
// 支持两种输入方式：
//   - 命令字符串（通过 New/NewArgs/NewCmds 创建）
//   - bash 脚本文件（通过 NewScript 创建）
//
// 字段可通过结构体字面量或 WithXxx 方法配置：
//
//	cmd := &Shx{
//	    Dir:     "/tmp",
//	    Timeout: 5 * time.Second,
//	}
type Shx struct {
	// 导出字段（可通过结构体字面量或 WithXxx 方法配置）
	Env     expand.Environ  // 环境变量
	Dir     string          // 工作目录
	Timeout time.Duration   // 超时时间
	Ctx     context.Context // 上下文
	Stdin   io.Reader       // 标准输入
	Stdout  io.Writer       // 标准输出
	Stderr  io.Writer       // 标准错误

	// 内部字段（不对外暴露）
	raw        string             // 原始命令字符串
	parser     *syntax.Parser     // 语法解析器
	scriptFile string             // 脚本文件路径
	cancel     context.CancelFunc // 超时上下文的取消函数
}

// ExitStatus 包装退出状态错误
//
// 通过 Unwrap() 支持错误链，可使用 errors.Is(err, interp.ExitStatus(N)) 检测。
type ExitStatus struct {
	Code uint8
	err  error // 原始错误，用于错误链
}

// Error 实现 error 接口
func (e ExitStatus) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

// Unwrap 返回原始错误，支持 errors.Is/As 遍历错误链
func (e ExitStatus) Unwrap() error {
	return e.err
}

// New 从字符串创建命令
//
// 参数:
//   - cmdStr: 命令字符串
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.New("echo hello world")
//	cmd := shx.New("ls -la | grep .go")
func New(cmdStr string) *Shx {
	return newShx(cmdStr, "")
}

// NewArgs 从命令名和可变参数创建命令
//
// 参数:
//   - cmd: 命令名
//   - args: 可变参数列表
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.NewArgs("ls", "-la", "/tmp")
//	cmd := shx.NewArgs("git", "commit", "-m", "message")
func NewArgs(cmd string, args ...string) *Shx {
	parts := make([]string, 0, 1+len(args))
	parts = append(parts, cmd)
	parts = append(parts, args...)
	return newShx(strings.Join(parts, " "), "")
}

// NewCmds 从命令切片创建命令
//
// 参数:
//   - cmds: 命令切片，每个元素是一个完整的命令部分
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.NewCmds([]string{"ls", "-la", "|", "grep", ".go"})
//	cmd := shx.NewCmds([]string{"echo", "hello", ">", "output.txt"})
func NewCmds(cmds []string) *Shx {
	return newShx(strings.Join(cmds, " "), "")
}

// NewScript 从 bash 脚本文件创建命令
//
// 参数:
//   - filePath: 脚本文件路径，必须以 .sh 结尾
//
// 返回:
//   - *Shx: 命令对象
//
// 示例:
//
//	cmd := shx.NewScript("deploy.sh")
//	cmd := shx.NewScript("/path/to/script.sh")
//
// 注意:
//   - 如果 filePath 为空或不是 .sh 后缀, 会 panic
func NewScript(filePath string) *Shx {
	if filePath == "" {
		panic("script file path must not be empty")
	}
	if strings.ToLower(filepath.Ext(filePath)) != ".sh" {
		panic(fmt.Sprintf("script file path must end with .sh extension: %s", filePath))
	}
	return newShx("", filePath)
}

// newShx 创建 Shx 对象并初始化公共字段
func newShx(rawStr, scriptFile string) *Shx {
	return &Shx{
		raw:        rawStr,
		scriptFile: scriptFile,
		parser:     syntax.NewParser(syntax.Variant(syntax.LangBash)),
		Env:        expand.ListEnviron(os.Environ()...),
		Dir:        mustGetwd(),
	}
}

// mustGetwd 获取当前工作目录, 如果失败则返回空字符串
func mustGetwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// Raw 获取原始命令字符串
//
// 返回:
//   - string: 原始命令字符串
func (s *Shx) Raw() string {
	return s.raw
}
