// Package shx 实现了 Shell 命令/脚本执行功能
// 支持命令字符串执行和 .sh 脚本文件执行两种模式
package shx

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitee.com/MM-Q/go-kit/fs"
	"gitee.com/MM-Q/shx"
)

// ShxConfig 配置结构体
type ShxConfig struct {
	Timeout time.Duration // -t: 超时时间
	Dir     string        // -d: 工作目录
	Env     []string      // -e: 环境变量 (格式: KEY=VALUE)
	Cmd     string        // 命令字符串或脚本文件路径
	Stdin   bool          // 是否有管道 stdin
}

// ShxCmdMain 执行 Shell 命令或脚本
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
//
// 执行模式:
//   - 如果命令以 .sh 结尾且文件存在：使用脚本执行模式 (shx.NewScript)
//   - 否则：使用命令字符串执行模式 (shx.New)
func ShxCmdMain(config ShxConfig) error {
	// 命令或脚本为空时返回错误
	if config.Cmd == "" {
		return fmt.Errorf("command or script path cannot be empty")
	}

	// 确定执行方式
	var executor *shx.Shx
	if isScriptFile(config.Cmd) {
		executor = shx.NewScript(config.Cmd)
	} else {
		executor = shx.New(config.Cmd)
	}

	// 配置执行器
	if config.Timeout > 0 {
		executor = executor.WithTimeout(config.Timeout)
	}
	if config.Dir != "" {
		executor = executor.WithDir(config.Dir)
	}
	if len(config.Env) > 0 {
		executor = executor.WithEnvs(config.Env)
	}
	if config.Stdin {
		executor = executor.WithStdin(os.Stdin)
	}

	// 设置标准输出和标准错误输出到终端
	executor = executor.WithStdout(os.Stdout).WithStderr(os.Stderr)

	// 执行并返回结果（shx 会处理退出码等）
	return executor.Exec()
}

// isScriptFile 判断是否为脚本文件
//
// 参数:
//   - path: 文件路径
//
// 返回值:
//   - bool: 是否是以 .sh 结尾且存在的文件
func isScriptFile(path string) bool {
	if !strings.HasSuffix(strings.ToLower(path), ".sh") {
		return false
	}
	return fs.Exists(path) && fs.IsFile(path)
}
