// Package shellx 命令执行核心模块
// 本文件定义了 Command 结构体及其相关方法，是 shellx 包的核心实现。
//
// Command 结构体采用一体化设计，集配置、构建、执行于一体，支持：
//   - 链式配置：WithWorkDir、WithEnv、WithTimeout、WithContext 等
//   - 同步执行：Exec、ExecOutput、ExecStdout
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
	"sync/atomic"
	"syscall"
	"time"
)

// Command 命令对象 - 集配置、构建、执行于一体
//
// 注意事项:
//   - Command 对象的配置方法 (WithXxx) 不是并发安全的，不要在多个 goroutine 中并发配置
//   - 每个 Command 对象只能执行一次，重复执行会返回错误
//   - 执行方法是并发安全的，使用 atomic.Bool 防止重复执行
//   - 属性获取方法不是并发安全的，不要在多个 goroutine 中并发调用
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
		panic("name must not be empty")
	}

	return &Command{
		name:      name,
		args:      args,
		envs:      os.Environ(), // 默认继承父进程的环境变量
		shellType: ShellDef1,    // 默认根据操作系统自动选择shell
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
		panic("cmdArgs must not be empty")
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
	if cmdStr == "" {
		panic("cmdStr must not be empty")
	}

	cmds, err := SplitE(cmdStr) // 使用命令拆分器拆分命令字符串
	if err != nil {
		// 拆分失败时触发panic，快速提醒开发者
		panic(fmt.Sprintf("command split failed: %v, original command: %q", err, cmdStr))
	}

	if len(cmds) == 0 {
		panic(fmt.Sprintf("split command is empty, original command: %q", cmdStr))
	}

	// 调用NewCmds来创建命令对象，复用逻辑
	cmd := NewCmds(cmds)

	// 保存原始命令字符串
	cmd.raw = cmdStr

	return cmd
}

// ############################################
// 配置方法
// ############################################

// WithWorkDir 设置命令的工作目录
//
// 参数：
//   - dir: 工作目录路径
//
// 返回：
//   - *Command: 命令对象
//
// 注意:
//   - 该方法会验证目录是否存在，如果不存在会panic
//   - 空字符串表示使用当前目录
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithWorkDir(dir string) *Command {
	// 如果目录不为空，验证目录
	if dir != "" {
		// 检查目录是否存在
		info, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				panic(fmt.Sprintf("directory %s does not exist", dir))
			}
			panic(fmt.Sprintf("stat %s failed: %v", dir, err))
		}

		// 检查是否为目录
		if !info.IsDir() {
			panic(fmt.Sprintf("%s is not a directory", dir))
		}
	} else {
		// 空字符串表示使用当前目录
		dir = "."
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
//   - 该方法会验证key是否为空, 如果为空会panic。
//   - 该方法会验证环境变量格式，格式错误会panic。
//   - 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithEnv(key, value string) *Command {
	if key == "" {
		panic("environment variable key cannot be empty")
	}

	envStr := fmt.Sprintf("%s=%s", key, value)
	// 验证环境变量格式
	if err := validateEnvVar(envStr); err != nil {
		panic(fmt.Sprintf("environment variable format error: %v", err))
	}
	c.envs = append(c.envs, envStr)
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
//   - 该方法会验证环境变量格式，格式错误的项会被忽略。
//   - 无需添加系统环境变量os.Environ(), 系统环境变量会自动继承.
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithEnvs(envs []string) *Command {
	if len(envs) == 0 {
		return c
	}

	// 验证每个环境变量的格式
	validEnvs := make([]string, 0, len(envs))
	for _, env := range envs {
		if err := validateEnvVar(env); err != nil {
			panic(fmt.Sprintf("environment variable format error: %v", err))
		}
		validEnvs = append(validEnvs, env)
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
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithTimeout(timeout time.Duration) *Command {
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
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithContext(ctx context.Context) *Command {
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
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithStdin(stdin io.Reader) *Command {
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
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithStdout(stdout io.Writer) *Command {
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
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithStderr(stderr io.Writer) *Command {
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
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发配置
func (c *Command) WithShell(shell ShellType) *Command {
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
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
func (c *Command) ShellType() ShellType {
	return c.shellType
}

// Raw 获取原始命令字符串
//
// 返回:
//   - string: 原始命令字符串
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
func (c *Command) Raw() string {
	return c.raw
}

// Name 获取命令名称
//
// 返回:
//   - string: 命令名称
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
func (c *Command) Name() string {
	return c.name
}

// Args 获取命令参数列表
//
// 返回:
//   - []string: 命令参数列表
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
//   - 返回的是参数的副本，修改返回值不会影响原始对象
func (c *Command) Args() []string {
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
		return c.getCmdStr()

	} else {
		return c.execCmd.String()
	}
}

// WorkDir 获取命令执行的工作目录
//
// 返回:
//   - string: 命令执行目录
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
func (c *Command) WorkDir() string {
	return c.dir
}

// Env 获取命令环境变量列表
//
// 返回:
//   - []string: 命令环境变量列表
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
//   - 返回的是环境变量的副本，修改返回值不会影响原始对象
func (c *Command) Env() []string {
	tempEnv := make([]string, len(c.envs))
	copy(tempEnv, c.envs)
	return tempEnv
}

// Timeout 获取命令执行超时时间
//
// 返回:
//   - time.Duration: 命令执行超时时间
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
func (c *Command) Timeout() time.Duration {
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
	if err := c.buildExecCmd(); err != nil {
		return err
	}

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
	if err := c.buildExecCmd(); err != nil {
		return nil, err
	}

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
	if err := c.buildExecCmd(); err != nil {
		return nil, err
	}

	// 确保资源清理
	defer c.cleanup()

	output, err := c.execCmd.Output()
	return output, judgeError(err, c)
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
	if err := c.buildExecCmd(); err != nil {
		return err
	}

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

// WaitWithCode 等待命令执行完成并返回退出码(仅在异步执行时有效)
//
// 返回:
//   - int: 命令退出码(0表示成功，-1表示无法提取的执行错误，其他值表示命令返回的退出码)
//   - error: 错误信息，可通过 IsTimeoutError() 和 IsCanceledError() 判断错误类型
func (c *Command) WaitWithCode() (int, error) {
	if c.execCmd == nil {
		return -1, ErrNotStarted
	}

	err := c.execCmd.Wait()

	// 清理资源
	c.cleanup()

	// 获取命令的退出码
	exitCode := extractExitCode(err)

	return exitCode, judgeError(err, c)
}

// Cmd 获取底层的 exec.Cmd 对象
//
// 返回:
//   - *exec.Cmd: 底层的 exec.Cmd 对象
func (c *Command) Cmd() *exec.Cmd {
	if c.execCmd == nil {
		if err := c.buildExecCmd(); err != nil {
			// 对于 Cmd() 方法，我们可以选择 panic 或返回 nil
			// 这里选择 panic 以保持行为一致性
			panic(err)
		}
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
// 注意: 此方法提供基本的进程状态检查，可能不是100%准确，
// 特别是在Windows系统上可能存在权限问题。
// 对于精确的进程状态管理，建议使用Wait或WaitWithCode方法。
//
// 返回:
//   - bool: 是否在运行
func (c *Command) IsRunning() bool {
	if c.execCmd == nil || c.execCmd.Process == nil {
		return false
	}

	// 如果ProcessState不为nil，表示进程已经结束
	if c.execCmd.ProcessState != nil {
		return false
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
//
// 注意:
//   - 此方法不是并发安全的，不要在多个goroutine中并发调用
func (c *Command) getEffectiveTimeout() time.Duration {
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
