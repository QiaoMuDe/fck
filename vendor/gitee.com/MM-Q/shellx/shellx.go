// Package shellx 提供了一个功能完善、易于使用的Go语言shell命令执行库。
//
// 本库基于Go标准库的os/exec包进行封装，提供了更加友好的API和丰富的功能，
// 支持同步和异步命令执行、输入输出重定向、精确超时控制、上下文管理、
// 多种shell类型支持等功能，并提供类型安全的API和友好的链式调用接口。
//
// 主要特性：
//   - 一体化设计：Command集配置、构建、执行于一体，无需Build()方法
//   - 支持三种命令创建方式：NewCmd(可变参数)、NewCmds(切片)、NewCmdStr(字符串解析)
//   - 丰富的便捷函数：Exec、ExecStr、ExecOut、ExecOutStr及其带超时版本
//   - 链式调用API，支持流畅的方法链
//   - 精确超时控制：延迟构建exec.Cmd，确保超时计时精确
//   - 完整的错误处理和类型安全
//   - 支持多种shell类型（sh、bash、cmd、powershell、pwsh等）
//   - 同步和异步执行支持
//   - 命令执行状态管理和进程控制
//   - 输入输出重定向和环境变量设置
//   - 上下文取消和优先级控制
//   - 并发安全的设计
//   - 跨平台兼容（Windows、Linux、macOS）
//
// 核心组件：
//   - Command: 命令对象，集配置、构建、执行于一体
//   - Result: 命令执行结果，包含输出、错误、时间等信息
//   - ShellType: Shell类型枚举，支持多种shell
//
// 基本用法：
//
//	import "gitee.com/MM-Q/shellx"
//
//	// 方式1：使用可变参数创建命令（无需Build）
//	err := shellx.NewCmd("ls", "-la").
//		WithWorkDir("/tmp").
//		WithTimeout(30 * time.Second).
//		WithShell(shellx.ShellBash).
//		Exec()
//
//	// 方式2：使用字符串创建命令
//	output, err := shellx.NewCmdStr(`echo "hello world"`).
//		WithEnv("MY_VAR", "value").
//		ExecOutput()
//
//	// 方式3：获取完整结果
//	result, err := shellx.NewCmd("git", "status").
//		WithTimeout(10 * time.Second).
//		ExecResult()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Exit Code: %d\n", result.Code())
//	fmt.Printf("Success: %t\n", result.Success())
//	fmt.Printf("Duration: %v\n", result.Duration())
//	fmt.Printf("Output: %s\n", result.Output())
//
// 便捷函数用法：
//
//	// 基础执行函数
//	err := shellx.Exec("ls", "-la")                    // 执行命令，输出到控制台
//	err := shellx.ExecStr("echo hello")                // 字符串方式执行
//	output, err := shellx.ExecOut("ls", "-la")         // 执行并返回输出
//	output, err := shellx.ExecOutStr("echo hello")     // 字符串方式执行并返回输出
//
//	// 带超时的执行函数
//	err := shellx.ExecT(5*time.Second, "sleep", "10")                    // 5秒超时
//	err := shellx.ExecStrT(3*time.Second, "ping google.com")       // 字符串方式，3秒超时
//	output, err := shellx.ExecOutT(2*time.Second, "curl", "example.com") // 返回输出，2秒超时
//	output, err := shellx.ExecOutStrT(1*time.Second, "date")             // 字符串方式，返回输出，1秒超时
//
// 高级用法：
//
//	// 设置标准输入输出
//	var stdout, stderr bytes.Buffer
//	stdin := strings.NewReader("input data")
//
//	err := shellx.NewCmd("cat").
//		WithStdin(stdin).
//		WithStdout(&stdout).
//		WithStderr(&stderr).
//		Exec()
//
//	// 使用上下文控制（优先级高于WithTimeout）
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	err := shellx.NewCmd("long-running-command").
//		WithContext(ctx).
//		WithTimeout(5*time.Second).  // 这个会被忽略
//		Exec()
//
//	// 异步执行和进程控制
//	cmd := shellx.NewCmd("sleep", "100")
//	err := cmd.ExecAsync()
//	pid := cmd.GetPID()
//	isRunning := cmd.IsRunning()
//	cmd.Kill() // 或 cmd.Signal(syscall.SIGTERM)
//	err = cmd.Wait()
//
// 命令解析：
//
//	// 支持复杂的命令字符串解析，包括引号处理
//	cmd := shellx.NewCmdStr(`git commit -m "Initial commit with 'quotes'"`)
//	// 解析结果：["git", "commit", "-m", "Initial commit with 'quotes'"]
//
// Shell类型：
//
//	// 支持多种shell类型
//	shellx.ShellSh         // sh shell
//	shellx.ShellBash       // bash shell
//	shellx.ShellCmd        // Windows cmd
//	shellx.ShellPowerShell // Windows PowerShell
//	shellx.ShellPwsh       // PowerShell Core
//	shellx.ShellNone       // 直接执行，不使用shell
//	shellx.ShellDef1    // 根据操作系统自动选择
//
// 注意事项：
//   - 每个Command对象只能执行一次，重复执行会返回错误
//   - Command是并发安全的，可以在多个goroutine中安全使用
//   - 命令执行会继承父进程的环境变量，可通过WithEnv添加额外变量
//   - 超时控制在执行时创建上下文，确保计时精确
//   - 异步执行需要调用Wait()等待完成或使用Kill()终止
package shellx
