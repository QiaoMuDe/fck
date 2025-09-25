// Package shellx 命令执行核心模块
// 本文件定义了 Command 结构体及其相关方法，是 shellx 包的核心实现。
//
// Command 结构体采用一体化设计，集配置、构建、执行于一体，支持：
//   - 链式配置：WithWorkDir、WithEnv、WithTimeout、WithContext 等
//   - 同步执行：Exec、ExecOutput、ExecStdout、ExecResult
//   - 异步执行：ExecAsync、Wait
//   - 进程控制：Kill、Signal、IsRunning、GetPID
//   - 状态管理：IsExecuted（确保命令只执行一次）
//   - 延迟构建：exec.Cmd 对象在执行时才创建，确保超时控制精确
//
// 提供完整的命令执行解决方案，支持多种执行模式和丰富的配置选项。
package shellx

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Command 命令对象 - 集配置、构建、执行于一体
type Command struct {
	// 基本命令配置
	shellType ShellType // shell类型
	raw       string    // 原始命令字符串
	name      string    // 命令名
	args      []string  // 命令参数

	// 执行环境配置
	dir    string    // 工作目录
	envs   []string  // 环境变量
	stdin  io.Reader // 标准输入
	stdout io.Writer // 标准输出
	stderr io.Writer // 标准错误输出

	// 上下文和超时配置
	userCtx context.Context // 用户设置的上下文
	timeout time.Duration   // 超时时间

	// 执行状态和控制
	execCmd *exec.Cmd          // 真正的exec.Cmd对象（延迟创建）
	cancel  context.CancelFunc // 超时上下文的取消函数
	execOne atomic.Bool        // 确保只执行一次
	mu      sync.RWMutex       // 保护配置字段的并发安全
}

// ############################################
// 构造函数
// ############################################

// NewCmd 创建新的命令对象 (数组方式 - 可变参数)
//
// 参数：
//   - name: 命令名
//   - args: 命令参数列表
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 默认通过shell执行, 可以通过WithShell方法指定shell类型
//   - 默认为ShellDef1, 根据操作系统自动选择shell(Windows系统默认为cmd, 其他系统默认为sh)
//   - 默认继承父进程的环境变量, 可以通过WithEnv方法设置环境变量
func NewCmd(name string, args ...string) *Command {
	if name == "" {
		panic("name cannot be empty")
	}

	return &Command{
		name:      name,
		args:      args,
		envs:      os.Environ(), // 默认继承父进程的环境变量
		shellType: ShellDef1,    // 默认根据操作系统自动选择shell
		mu:        sync.RWMutex{},
	}
}

// NewCmds 创建新的命令对象 (数组方式 - 切片参数)
//
// 参数：
//   - cmdArgs: 命令参数列表，第一个元素为命令名，后续元素为参数
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 默认通过shell执行, 可以通过WithShell方法指定shell类型
//   - 默认为ShellDef1, 根据操作系统自动选择shell(Windows系统默认为cmd, 其他系统默认为sh)
//   - 默认继承父进程的环境变量, 可以通过WithEnv方法设置环境变量
func NewCmds(cmdArgs []string) *Command {
	if len(cmdArgs) == 0 {
		panic("cmdArgs cannot be empty")
	}

	name := cmdArgs[0] // 第一个元素为命令名
	args := []string{} // 后续元素为参数
	if len(cmdArgs) > 1 {
		args = cmdArgs[1:]
	}

	return NewCmd(name, args...)
}

// NewCmdStr 创建新的命令对象 (字符串方式)
//
// 参数：
//   - cmdStr: 命令字符串
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 默认通过shell执行, 可以通过WithShell方法指定shell类型
//   - 默认为ShellDef1, 根据操作系统自动选择shell(Windows系统默认为cmd, 其他系统默认为sh)
//   - 默认继承父进程的环境变量, 可以通过WithEnv方法设置环境变量
func NewCmdStr(cmdStr string) *Command {
	cmds := ParseCmd(cmdStr) // 使用命令解析器解析命令字符串
	cmd := NewCmds(cmds)
	cmd.raw = cmdStr // 保存原始命令字符串
	return cmd
}

// ############################################
// 配置方法
// ############################################

// WithWorkDir 设置命令的工作目录
//
// 参数：
//   - dir: 命令的工作目录
//
// 返回：
//   - *Command: 命令对象
func (c *Command) WithWorkDir(dir string) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	if dir == "" {
		return c
	}

	info, statErr := os.Lstat(dir)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			panic(fmt.Sprintf("dir %s does not exist", dir))
		}

		panic(fmt.Sprintf("stat %s failed: %v", dir, statErr))
	}
	if !info.IsDir() {
		panic(fmt.Sprintf("dir %s is not a directory", dir))
	}

	c.dir = dir
	return c
}

// WithEnv 设置命令的环境变量
//
// 参数：
//   - key: 环境变量的键
//   - value: 环境变量的值
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 该方法会验证key是否为空, 如果为空则忽略。
//   - 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.
func (c *Command) WithEnv(key, value string) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.envs == nil {
		c.envs = os.Environ()
	}

	if key != "" {
		c.envs = append(c.envs, fmt.Sprintf("%s=%s", key, value))
	}
	return c
}

// WithEnvs 批量设置命令的环境变量
//
// 参数：
//   - envs: []string类型，环境变量列表，每个元素为"key=value"格式
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 该方法会验证环境变量格式，只添加验证通过的环境变量。
//   - 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.
func (c *Command) WithEnvs(envs []string) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(envs) == 0 {
		return c
	}

	if c.envs == nil {
		c.envs = os.Environ()
	}

	// 验证环境变量格式，只添加验证通过的环境变量
	validEnvs := make([]string, 0, len(envs))
	for _, env := range envs {
		if err := validateEnvVar(env); err == nil {
			validEnvs = append(validEnvs, env)
		}
	}

	c.envs = append(c.envs, validEnvs...)
	return c
}

// WithTimeout 设置命令的超时时间(便捷方式)
//
// 参数：
//   - timeout: time.Duration类型，命令执行的超时时间
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 该方法会验证超时时间是否小于等于0, 如果小于等于0则忽略。
//   - 该超时时间优先级低于上下文设置的超时时间.
func (c *Command) WithTimeout(timeout time.Duration) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 超时大于0时才设置
	if timeout > 0 {
		c.timeout = timeout
	}
	return c
}

// WithContext 设置命令的上下文
//
// 参数：
//   - ctx: context.Context类型，用于取消命令执行和超时控制
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 该方法会验证上下文是否为空，如果为空则panic.
//   - 该上下文会覆盖之前设置的超时时间.
func (c *Command) WithContext(ctx context.Context) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ctx == nil {
		panic("context cannot be nil")
	}
	c.userCtx = ctx
	return c
}

// WithStdin 设置命令的标准输入
//
// 参数：
//   - stdin: io.Reader类型，用于提供命令的标准输入
//
// 返回：
//   - *Command: 命令对象
func (c *Command) WithStdin(stdin io.Reader) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stdin == nil {
		panic("stdin cannot be nil")
	}
	c.stdin = stdin
	return c
}

// WithStdout 设置命令的标准输出
//
// 参数：
//   - stdout: io.Writer类型，用于接收命令的标准输出
//
// 返回：
//   - *Command: 命令对象
func (c *Command) WithStdout(stdout io.Writer) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stdout == nil {
		panic("stdout cannot be nil")
	}
	c.stdout = stdout
	return c
}

// WithStderr 设置命令的标准错误输出
//
// 参数：
//   - stderr: io.Writer类型，用于接收命令的标准错误输出
//
// 返回：
//   - *Command: 命令对象
func (c *Command) WithStderr(stderr io.Writer) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stderr == nil {
		panic("stderr cannot be nil")
	}
	c.stderr = stderr
	return c
}

// WithShell 设置命令的shell类型
//
// 参数：
//   - shell: ShellType类型，表示要使用的shell类型
//
// 返回：
//   - *Command: 命令对象
func (c *Command) WithShell(shell ShellType) *Command {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.shellType = shell
	return c
}

// ############################################
// 属性获取方法
// ############################################

// ShellType 获取shell类型
//
// 返回:
//   - ShellType: shell类型
func (c *Command) ShellType() ShellType {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.shellType
}

// Raw 获取原始命令字符串
//
// 返回:
//   - string: 原始命令字符串
func (c *Command) Raw() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.raw
}

// Name 获取命令名称
//
// 返回:
//   - string: 命令名称
func (c *Command) Name() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.name
}

// Args 获取命令参数列表
//
// 返回:
//   - []string: 命令参数列表
func (c *Command) Args() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tempArgs := make([]string, len(c.args))
	copy(tempArgs, c.args)
	return tempArgs
}

// CmdStr 获取命令字符串
//
// 返回:
//   - string: 命令字符串
func (c *Command) CmdStr() string {
	if c.execCmd == nil {
		return c.raw
	} else {
		return c.execCmd.String()
	}
}

// WorkDir 获取命令执行的工作目录
//
// 返回:
//   - string: 命令执行目录
func (c *Command) WorkDir() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.dir
}

// Env 获取命令环境变量列表
//
// 返回:
//   - []string: 命令环境变量列表
func (c *Command) Env() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	tempEnv := make([]string, len(c.envs))
	copy(tempEnv, c.envs)
	return tempEnv
}

// Timeout 获取命令执行超时时间
//
// 返回:
//   - time.Duration: 命令执行超时时间
func (c *Command) Timeout() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.timeout
}

// ############################################
// 执行方法
// ############################################

// Exec 执行命令(阻塞)
//
// 返回:
//   - error: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型
func (c *Command) Exec() error {
	if !c.execOne.CompareAndSwap(false, true) {
		return ErrAlreadyExecuted
	}

	// 执行时才构建真正的exec.Cmd
	c.buildExecCmd()

	// 确保资源清理
	defer c.cleanup()

	err := c.execCmd.Run()
	return judgeError(err, c)
}

// ExecOutput 执行命令并返回合并后的输出(阻塞)
//
// 返回:
//   - []byte: 命令输出
//   - error: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型
//
// 注意:
//   - 由于需要捕获默认的stdout和stderr合并输出, 内部已经设置了WithStdout(os.Stdout)和WithStderr(os.Stderr)
func (c *Command) ExecOutput() ([]byte, error) {
	if !c.execOne.CompareAndSwap(false, true) {
		return nil, ErrAlreadyExecuted
	}

	// 执行时才构建真正的exec.Cmd
	c.buildExecCmd()

	// 确保资源清理
	defer c.cleanup()

	output, err := c.execCmd.CombinedOutput()
	return output, judgeError(err, c)
}

// ExecStdout 执行命令并返回标准输出(阻塞)
//
// 返回:
//   - []byte: 标准输出
//   - error: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型
func (c *Command) ExecStdout() ([]byte, error) {
	if !c.execOne.CompareAndSwap(false, true) {
		return nil, ErrAlreadyExecuted
	}

	// 执行时才构建真正的exec.Cmd
	c.buildExecCmd()

	// 确保资源清理
	defer c.cleanup()

	output, err := c.execCmd.Output()
	return output, judgeError(err, c)
}

// ExecResult 执行命令并返回完整的执行结果(阻塞)
//
// 使用示例:
//
//	result, err := cmd.ExecResult()
//	if err != nil {
//	    if IsTimeoutError(err) {
//	        log.Printf("Command timeout: %v", err)
//	    } else if IsCanceledError(err) {
//	        log.Printf("Command canceled: %v", err)
//	    } else {
//	        log.Printf("Command failed: %v", err)
//	    }
//	    return
//	}
//	// 处理成功情况
//	fmt.Println(string(result.Output()))
//
// 返回:
//   - *Result: 执行结果对象, 包含输出、时间、退出码等信息
//   - error: 执行过程中的错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型
func (c *Command) ExecResult() (*Result, error) {
	if !c.execOne.CompareAndSwap(false, true) {
		return nil, ErrAlreadyExecuted
	}

	// 执行时才构建真正的exec.Cmd
	c.buildExecCmd()

	// 确保资源清理
	defer c.cleanup()

	// 命令执行开始时间
	startTime := time.Now()

	// 执行命令
	output, err := c.execCmd.CombinedOutput()

	// 命令执行结束时间
	endTime := time.Now()

	// 获取命令的退出码
	var exitCode int
	if err != nil {
		exitCode = -1
	}
	// 创建Result对象
	result := &Result{
		startTime: startTime,              // 命令开始时间
		endTime:   endTime,                // 命令结束时间
		duration:  endTime.Sub(startTime), // 命令执行时间
		output:    output,                 // 命令输出
		success:   err == nil,             // 命令是否执行成功
		exitCode:  exitCode,               // 命令退出码
	}

	return result, judgeError(err, c)
}

// ExecAsync 异步执行命令(非阻塞)
//
// 返回:
//   - error: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型
func (c *Command) ExecAsync() error {
	if !c.execOne.CompareAndSwap(false, true) {
		return ErrAlreadyExecuted
	}

	// 执行时才构建真正的exec.Cmd
	c.buildExecCmd()

	err := c.execCmd.Start()
	return judgeError(err, c)
}

// Wait 等待命令执行完成(仅在异步执行时有效)
//
// 返回:
//   - error: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型
func (c *Command) Wait() error {
	if c.execCmd == nil {
		return ErrNotStarted
	}

	err := c.execCmd.Wait()

	// 清理资源
	c.cleanup()

	return judgeError(err, c)
}

// Cmd 获取底层的 exec.Cmd 对象
//
// 返回:
//   - *exec.Cmd: 底层的 exec.Cmd 对象
func (c *Command) Cmd() *exec.Cmd {
	if c.execCmd == nil {
		c.buildExecCmd() // 如果还没构建，先构建
	}
	return c.execCmd
}

// Kill 杀死当前命令的进程
//
// 返回:
//   - error: 错误信息
func (c *Command) Kill() error {
	if c.execCmd == nil || c.execCmd.Process == nil {
		return ErrNoProcess
	}
	return c.execCmd.Process.Kill()
}

// Signal 向当前进程发送信号
//
// 参数:
//   - sig: 信号类型
//
// 返回:
//   - error: 错误信息
func (c *Command) Signal(sig os.Signal) error {
	if c.execCmd == nil || c.execCmd.Process == nil {
		return ErrNoProcess
	}
	return c.execCmd.Process.Signal(sig)
}

// IsRunning 检查进程是否还在运行
//
// 返回:
//   - bool: 是否在运行
func (c *Command) IsRunning() bool {
	if c.execCmd == nil || c.execCmd.Process == nil {
		return false
	}

	if c.execCmd.ProcessState != nil {
		return false // 进程已结束
	}

	// 尝试发送信号0检查进程是否存在
	err := c.execCmd.Process.Signal(syscall.Signal(0))
	return err == nil
}

// GetPID 获取进程ID
//
// 返回:
//   - int: 进程ID, 如果进程不存在返回0
func (c *Command) GetPID() int {
	if c.execCmd == nil || c.execCmd.Process == nil {
		return 0
	}
	return c.execCmd.Process.Pid
}

// IsExecuted 检查命令是否已经执行过
//
// 返回:
//   - bool: 是否已执行
func (c *Command) IsExecuted() bool {
	return c.execOne.Load()
}

// getEffectiveTimeout 获取有效的超时时间
// 优先使用用户上下文的超时，其次使用设置的超时时间
func (c *Command) getEffectiveTimeout() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 如果有用户上下文且有截止时间，计算剩余时间
	if c.userCtx != nil {
		if deadline, ok := c.userCtx.Deadline(); ok {
			remaining := time.Until(deadline) // 计算剩余时间
			if remaining > 0 {
				return remaining
			}
		}
	}

	// 否则返回设置的超时时间
	return c.timeout
}
