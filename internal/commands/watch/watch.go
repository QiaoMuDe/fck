package watch

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/shellx/shx"
	"golang.org/x/term"
)

// Executor 命令执行器
type Executor struct {
	command string        // 要执行的命令
	timeout time.Duration // 超时时间
}

// NewExecutor 创建新的执行器
func NewExecutor(command string, timeout time.Duration) *Executor {
	return &Executor{command: command, timeout: timeout}
}

// Run 执行命令并返回结果
func (e *Executor) Run(ctx context.Context) (*ExecutionResult, error) {
	start := time.Now()

	// 捕获输出（带大小限制）
	stdoutWriter := newLimitedWriter(maxOutputSize)
	stderrWriter := newLimitedWriter(maxOutputSize / 10) // stderr 限制 1MB

	// 使用 shellx/shx 执行命令（纯 Go 实现，跨平台一致）
	sh := shx.New(e.command).
		WithTimeout(e.timeout).
		WithContext(ctx).
		WithStdout(stdoutWriter).
		WithStderr(stderrWriter)

	err := sh.Exec()
	duration := time.Since(start)

	stdout := stdoutWriter.String()
	stderr := stderrWriter.String()

	// 合并 stdout 和 stderr（stderr 在前，类似 Linux watch）
	combinedOutput := stdout
	if stderr != "" {
		if stdout != "" {
			combinedOutput = stderr + "\n" + stdout
		} else {
			combinedOutput = stderr
		}
	}

	result := &ExecutionResult{
		Stdout:   combinedOutput,
		Stderr:   stderr,
		ExitCode: 0,
		Duration: duration,
	}

	// 解析退出码
	if exitErr, ok := err.(shx.ExitStatus); ok {
		result.ExitCode = int(exitErr.Code)
		// ExitStatus 是正常返回，不是错误
		return result, nil
	}
	if err != nil {
		// 检查是否是超时错误
		if ctx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("command timed out after %v", e.timeout)
		}
		return result, fmt.Errorf("command execution failed: %w", err)
	}

	return result, nil
}

// Scheduler 执行调度器
type Scheduler struct {
	interval     time.Duration // 执行间隔
	precise      bool          // 精确计时模式
	lastRunTime  time.Time     // 上次执行时间
	lastDuration time.Duration // 上次执行耗时
}

// NewScheduler 创建新的调度器
func NewScheduler(interval time.Duration, precise bool) *Scheduler {
	return &Scheduler{interval: interval, precise: precise}
}

// RecordRun 记录本次执行信息
func (s *Scheduler) RecordRun(duration time.Duration) {
	s.lastRunTime = time.Now()
	s.lastDuration = duration
}

// NextWait 计算下次等待时间
func (s *Scheduler) NextWait() time.Duration {
	if !s.precise {
		return s.interval
	}
	// 精确模式：补偿执行耗时
	if s.lastRunTime.IsZero() {
		return s.interval
	}
	elapsed := time.Since(s.lastRunTime)
	wait := s.interval - elapsed
	if wait < 0 {
		return 0
	}
	return wait
}

// DiffHighlighter 差异高亮器
type DiffHighlighter struct {
	mode       DiffMode           // 差异模式
	lastOutput string             // 上次输出
	cumulative map[int]bool       // 累积变化行记录
	cl         *colorlib.ColorLib // 颜色库
}

// NewDiffHighlighter 创建新的差异高亮器
func NewDiffHighlighter(mode DiffMode, cl *colorlib.ColorLib) *DiffHighlighter {
	return &DiffHighlighter{
		mode:       mode,
		cumulative: make(map[int]bool),
		cl:         cl,
	}
}

// Diff 对比当前输出与上次输出，返回带高亮的结果
func (d *DiffHighlighter) Diff(current string) string {
	if d.mode == DiffModeNone {
		return current
	}
	oldLines := strings.Split(d.lastOutput, "\n")
	newLines := strings.Split(current, "\n")
	diffLines := d.computeLineDiff(oldLines, newLines)
	result := d.renderDiff(diffLines)
	d.lastOutput = current
	return result
}

// computeLineDiff 计算行级差异
func (d *DiffHighlighter) computeLineDiff(oldLines, newLines []string) []LineDiff {
	result := make([]LineDiff, len(newLines))
	for i, line := range newLines {
		status := LineDiffSame
		if i >= len(oldLines) {
			status = LineDiffAdded
		} else if line != oldLines[i] {
			status = LineDiffChanged
		}
		result[i] = LineDiff{Content: line, Status: status}
	}
	// 累积模式：标记所有历史变化过的行
	if d.mode == DiffModeCumulative {
		for i := range result {
			if result[i].Status != LineDiffSame {
				d.cumulative[i] = true
			}
			if d.cumulative[i] {
				result[i].Status = LineDiffChanged
			}
		}
	}
	return result
}

// renderDiff 渲染带高亮的差异
func (d *DiffHighlighter) renderDiff(lines []LineDiff) string {
	if d.cl == nil {
		// 无颜色模式，使用前缀标记
		var result []string
		for _, line := range lines {
			prefix := "  "
			switch line.Status {
			case LineDiffChanged:
				prefix = "| "
			case LineDiffAdded:
				prefix = "+ "
			}
			result = append(result, prefix+line.Content)
		}
		return strings.Join(result, "\n")
	}
	// 使用颜色高亮
	var result []string
	for _, line := range lines {
		switch line.Status {
		case LineDiffChanged:
			result = append(result, d.cl.SbrightYellow(line.Content))
		case LineDiffAdded:
			result = append(result, d.cl.SbrightGreen(line.Content))
		default:
			result = append(result, line.Content)
		}
	}
	return strings.Join(result, "\n")
}

// OutputManager 输出管理器
type OutputManager struct {
	noHeader    bool               // 不显示标题栏
	useColor    bool               // 使用颜色
	clearScreen bool               // 清屏
	cl          *colorlib.ColorLib // 颜色库
}

// NewOutputManager 创建新的输出管理器
func NewOutputManager(noHeader, noColor, clearScreen bool, cl *colorlib.ColorLib) *OutputManager {
	return &OutputManager{
		noHeader:    noHeader,
		useColor:    !noColor && cl != nil,
		clearScreen: clearScreen,
		cl:          cl,
	}
}

// Clear 清屏
func (o *OutputManager) Clear() {
	if !o.clearScreen {
		return
	}
	// ANSI 清屏序列：光标移动到左上角 + 清屏 + 清除滚动缓冲区
	fmt.Print("\033[H\033[2J\033[3J")
	if runtime.GOOS == "windows" {
		fmt.Print("\r")
	}
}

// PrintHeader 打印标题栏
func (o *OutputManager) PrintHeader(interval time.Duration, command string) {
	if o.noHeader {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	intervalStr := formatInterval(interval)
	// 格式: Every Xs: command                    时间
	prefix := fmt.Sprintf("Every %s: %s", intervalStr, command)

	// 获取终端宽度，默认 80
	terminalWidth := o.getTerminalWidth()
	padding := terminalWidth - len(prefix) - len(timestamp)
	if padding < 1 {
		padding = 1
	}
	header := prefix + strings.Repeat(" ", padding) + timestamp
	if o.useColor && o.cl != nil {
		header = o.cl.Scyan(header)
	}
	fmt.Println(header)
	fmt.Println()
}

// getTerminalWidth 获取终端宽度，失败时返回默认值 80
func (o *OutputManager) getTerminalWidth() int {
	// 检查 stdout 是否为终端
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return defaultWidth
	}

	// 获取终端尺寸
	width, _, err := term.GetSize(fd)
	if err != nil {
		return defaultWidth
	}

	// 限制最小宽度，避免异常值
	if width < 40 {
		return defaultWidth
	}

	return width
}

// PrintOutput 打印命令输出
func (o *OutputManager) PrintOutput(output string) {
	fmt.Print(output)
	// 确保输出以换行符结尾
	if !strings.HasSuffix(output, "\n") {
		fmt.Println()
	}
}

// PrintError 打印错误信息
func (o *OutputManager) PrintError(err error) {
	if o.useColor && o.cl != nil {
		fmt.Println(o.cl.SbrightRed(fmt.Sprintf("error: %v", err)))
	} else {
		fmt.Printf("error: %v\n", err)
	}
}

// formatInterval 格式化间隔时间显示
func formatInterval(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d.Seconds() == float64(int(d.Seconds())) {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

// Runner watch 命令运行器
type Runner struct {
	config      *WatchConfig     // 配置
	executor    *Executor        // 执行器
	scheduler   *Scheduler       // 调度器
	highlighter *DiffHighlighter // 差异高亮器
	output      *OutputManager   // 输出管理器
}

// NewRunner 创建新的 watch 运行器
func NewRunner(config *WatchConfig, cl *colorlib.ColorLib) *Runner {
	return &Runner{
		config:      config,
		executor:    NewExecutor(config.Command, config.Timeout),
		scheduler:   NewScheduler(config.Interval, config.Precise),
		highlighter: NewDiffHighlighter(config.Diff, cl),
		output:      NewOutputManager(config.NoHeader, config.NoColor, config.ClearScreen, cl),
	}
}

// Run 执行 watch 命令
func (r *Runner) Run(ctx context.Context) error {
	// 验证配置
	if err := r.config.Validate(); err != nil {
		return err
	}

	// 设置信号监听
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-sigChan
		cancel()
	}()

	executionCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// 检查执行次数限制
		if r.config.MaxCount > 0 && executionCount >= r.config.MaxCount {
			return nil
		}
		executionCount++

		// 清屏（在命令执行前）
		r.output.Clear()
		// 打印标题栏
		r.output.PrintHeader(r.config.Interval, r.config.Command)

		// 执行命令
		result, err := r.executor.Run(ctx)

		if err != nil {
			// 静默模式下不打印错误，但 ExitOnError 仍然有效
			if !r.config.Quiet {
				r.output.PrintError(err)
			}
			if r.config.ExitOnError {
				return err
			}
		} else if !r.config.Quiet {
			// 非静默模式：处理差异高亮和输出
			output := result.Stdout
			if r.config.Diff != DiffModeNone {
				output = r.highlighter.Diff(output)
			}
			// 输出结果
			r.output.PrintOutput(output)
		}

		// 记录执行时间（用于精确计时）
		r.scheduler.RecordRun(result.Duration)

		// 最后一次不等待
		if r.config.MaxCount > 0 && executionCount >= r.config.MaxCount {
			return nil
		}

		// 等待下次执行
		waitTime := r.scheduler.NextWait()
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(waitTime):
		}
	}
}

// WatchCmdMain 执行 watch 命令（向后兼容的入口函数）
func WatchCmdMain(config WatchConfig) error {
	cl := colorlib.NewColorLib()
	runner := NewRunner(&config, cl)
	return runner.Run(context.Background())
}
