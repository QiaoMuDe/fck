// Package watch 实现了命令监控功能。
// 该包提供了周期性执行命令并显示输出结果的功能，
// 支持间隔设置、次数限制、差异高亮、清屏、颜色输出等特性，类似于 Linux 的 watch 命令。
package watch

import (
	"bytes"
	"errors"
	"time"
)

// DiffMode 差异高亮模式
type DiffMode int

const (
	// DiffModeNone 无差异高亮
	DiffModeNone DiffMode = iota
	// DiffModeLine 行级差异
	DiffModeLine
	// DiffModeWord 词级差异
	DiffModeWord
	// DiffModeCumulative 累积差异（显示所有变化过的行）
	DiffModeCumulative
)

// LineDiffStatus 行的差异状态
type LineDiffStatus int

const (
	// LineDiffSame 无变化
	LineDiffSame LineDiffStatus = iota
	// LineDiffChanged 已变化
	LineDiffChanged
	// LineDiffAdded 新增行
	LineDiffAdded
	// LineDiffRemoved 删除行
	LineDiffRemoved
)

// LineDiff 带有差异状态的行
type LineDiff struct {
	Content string         // 行内容
	Status  LineDiffStatus // 差异状态
}

// WatchConfig watch 命令配置
type WatchConfig struct {
	Command     string        // 要执行的命令
	Interval    time.Duration // 执行间隔
	NoHeader    bool          // 不显示标题栏
	NoColor     bool          // 禁用颜色
	ClearScreen bool          // 每次执行前清屏
	Diff        DiffMode      // 差异高亮模式
	MaxCount    int           // 最大执行次数，-1 表示无限
	ExitOnError bool          // 出错时退出
	Timeout     time.Duration // 单次执行超时
	Precise     bool          // 精确计时模式
	Quiet       bool          // 静默模式（不显示输出、不计算差异）
}

// Validate 验证配置有效性
func (c *WatchConfig) Validate() error {
	if c.Command == "" {
		return errors.New("command is empty")
	}
	if c.MaxCount < -1 {
		return errors.New("maxCount must be >= -1 (-1 or 0 means unlimited)")
	}
	if c.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}
	return nil
}

// ExecutionResult 命令执行结果
type ExecutionResult struct {
	Stdout   string        // 标准输出
	Stderr   string        // 标准错误
	ExitCode int           // 退出码
	Duration time.Duration // 执行耗时
}

// 最大输出限制（10MB）
const maxOutputSize = 10 * 1024 * 1024

// 默认终端宽度
const defaultWidth = 80

// limitedWriter 带大小限制的写入器
type limitedWriter struct {
	buf       *bytes.Buffer
	maxSize   int
	truncated bool
}

// newLimitedWriter 创建带大小限制的写入器
func newLimitedWriter(maxSize int) *limitedWriter {
	return &limitedWriter{
		buf:     &bytes.Buffer{},
		maxSize: maxSize,
	}
}

// Write 实现 io.Writer 接口
func (w *limitedWriter) Write(p []byte) (n int, err error) {
	if w.truncated {
		return len(p), nil
	}

	// 检查是否会超过限制
	if w.buf.Len()+len(p) > w.maxSize {
		// 只写入剩余空间
		remaining := w.maxSize - w.buf.Len()
		if remaining > 0 {
			w.buf.Write(p[:remaining])
		}
		w.truncated = true
		// 添加截断提示
		w.buf.WriteString("\n... [output truncated due to size limit] ...\n")
		return len(p), nil
	}

	return w.buf.Write(p)
}

// String 返回缓冲区内容
func (w *limitedWriter) String() string {
	return w.buf.String()
}
