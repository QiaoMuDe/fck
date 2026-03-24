// Package shx 提供了基于 mvdan.cc/sh/v3 的纯 Go shell 命令执行功能。
//
// 本包是 ShellX 的子包, 提供了与主包相似的 API 风格, 但具有更好的跨平台一致性。
// 它使用 mvdan.cc/sh/v3 进行命令解析和执行, 不依赖系统 shell。
//
// 主要特性:
//   - 纯 Go 实现, 不依赖系统 shell
//   - 更好的跨平台一致性 (Windows/Linux/macOS 行为一致)
//   - 链式调用 API, 支持流畅的方法链
//   - 支持超时控制和上下文取消
//   - 最小并发保护 (使用 atomic.Bool 防止重复执行)
//
// 基本用法:
//
//	import "gitee.com/MM-Q/shellx/shx"
//
//	// 简单执行
//	err := shx.Exec("echo hello world")
//
//	// 获取输出
//	output, err := shx.Output("ls -la")
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
// 注意事项:
//   - Shx 对象的配置方法 (WithXxx) 不是并发安全的, 不要在多个 goroutine 中并发配置
//   - 每个 Shx 对象只能执行一次, 重复执行会返回错误
//   - mvdan/sh 是同步执行的, 不提供异步 API, 如需异步请使用 goroutine 包装
//   - 不支持进程控制 (无 PID、Kill、Signal) , 只能通过 context 取消
package shx

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

// Shx 表示一个待执行的 shell 命令
type Shx struct {
	// 命令配置
	raw    string         // 原始命令字符串
	parser *syntax.Parser // 语法解析器 (可自定义)

	// 执行环境
	dir    string         // 工作目录
	env    expand.Environ // 环境变量
	stdin  io.Reader      // 标准输入
	stdout io.Writer      // 标准输出
	stderr io.Writer      // 标准错误

	// 上下文和超时
	ctx     context.Context    // 用户上下文
	timeout time.Duration      // 超时时间
	cancel  context.CancelFunc // 超时上下文的取消函数

	// 执行状态 (使用 atomic.Bool 实现最小并发保护)
	executed atomic.Bool // 是否已执行
}

// ExitStatus 包装退出状态错误
type ExitStatus struct {
	Code uint8
}

// Error 实现 error 接口
func (e ExitStatus) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}
