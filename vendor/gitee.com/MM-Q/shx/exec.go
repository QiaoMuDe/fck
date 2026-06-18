package shx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Exec 执行命令
//
// 返回:
//   - error: 执行过程中的错误, 不包含退出码错误
//
// 示例:
//
//	err := shx.New("echo hello").Exec()
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *Shx) Exec() error {
	ctx := s.buildContext()
	return s.execWithContext(ctx, nil, nil)
}

// ExecOutput 执行命令并返回输出
//
// 返回:
//   - []byte: 命令输出 (stdout 和 stderr 合并)
//   - error: 执行过程中的错误
//
// 注意:
//   - 内部会自动捕获 stdout 和 stderr
//   - 如果需要区分 stdout 和 stderr, 请使用 WithStdout 和 WithStderr 自定义
func (s *Shx) ExecOutput() ([]byte, error) {
	var buf bytes.Buffer
	ctx := s.buildContext()
	err := s.execWithContext(ctx, &buf, &buf)
	return buf.Bytes(), err
}

// ExecContext 在指定上下文中执行命令
//
// 参数:
//   - ctx: 上下文 (用于取消执行)
//
// 返回:
//   - error: 执行过程中的错误
//
// 注意:
//   - 此方法会覆盖之前通过 WithContext 设置的上下文
//   - 此方法不受 WithTimeout 影响 (上下文的超时优先)
func (s *Shx) ExecContext(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}

	return s.execWithContext(ctx, nil, nil)
}

// ExecContextOutput 在指定上下文中执行并返回输出
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - []byte: 命令输出
//   - error: 执行过程中的错误
func (s *Shx) ExecContextOutput(ctx context.Context) ([]byte, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}

	var buf bytes.Buffer
	err := s.execWithContext(ctx, &buf, &buf)
	return buf.Bytes(), err
}

// buildContext 构建执行上下文
//
// 严格优先级:
// 1. 用户通过 WithContext 设置的上下文 (完全覆盖 WithTimeout 设置)
// 2. 用户通过 WithTimeout 设置的超时 (创建新上下文)
// 3. 默认背景上下文
func (s *Shx) buildContext() context.Context {
	// 如果之前有cancel函数，先调用清理资源
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	// 严格优先级：用户上下文完全覆盖超时设置
	if s.Ctx != nil {
		return s.Ctx
	}

	// 只有在没有设置用户上下文时才使用超时
	if s.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
		s.cancel = cancel
		return ctx
	}

	return context.Background()
}

// execWithContext 使用指定上下文执行命令
//
// 参数:
//   - ctx: 执行上下文
//
// 返回:
//   - error: 执行错误
func (s *Shx) execWithContext(ctx context.Context, stdout, stderr io.Writer) error {
	// 确保在退出时调用 cancel 函数
	if s.cancel != nil {
		defer s.cancel()
	}

	// 分支: 从脚本文件解析 vs 从命令字符串解析
	var file *syntax.File
	var err error

	if s.scriptFile != "" {
		// 从脚本文件解析
		f, openErr := os.Open(s.scriptFile)
		if openErr != nil {
			return fmt.Errorf("open script file failed: %w", openErr)
		}
		defer func() { _ = f.Close() }()
		file, err = s.parser.Parse(f, s.scriptFile)

	} else {
		// 检查空命令
		if strings.TrimSpace(s.raw) == "" {
			return fmt.Errorf("command cannot be empty")
		}
		// 从命令字符串解析
		file, err = s.parser.Parse(bytes.NewReader([]byte(s.raw)), "")
	}

	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// 创建执行器
	runner, err := s.buildRunner(stdout, stderr)
	if err != nil {
		return err
	}

	// 执行命令
	err = runner.Run(ctx, file)
	return handleError(err, s.displayName(), s.Timeout)
}

// displayName 返回用于错误信息显示的标识
func (s *Shx) displayName() string {
	if s.scriptFile != "" {
		return s.scriptFile
	}
	return s.raw
}

// buildRunner 构建执行器
//
// 返回:
//   - *interp.Runner: 执行器
//   - error: 构建错误
func (s *Shx) buildRunner(stdout, stderr io.Writer) (*interp.Runner, error) {
	out := stdout
	if out == nil {
		out = s.Stdout
	}
	errOut := stderr
	if errOut == nil {
		errOut = s.Stderr
	}

	opts := []interp.RunnerOption{
		interp.Env(s.Env),
		interp.Dir(s.Dir),
		interp.StdIO(s.Stdin, out, errOut),
	}

	return interp.New(opts...)
}
