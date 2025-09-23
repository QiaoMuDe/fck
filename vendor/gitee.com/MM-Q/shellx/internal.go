// Package shellx 内部实现模块
// 本文件包含 Command 结构体的内部实现方法，包括：
//   - buildExecCmd: 延迟构建 exec.Cmd 对象，支持上下文和超时控制
//   - cleanup: 资源清理函数，确保上下文取消函数被正确调用
//   - getCmdStr: 命令字符串获取函数，支持原始字符串和参数拼接
//
// 这些方法为 Command 的核心功能提供底层支持。
package shellx

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// buildExecCmd 在执行时构建真正的exec.Cmd对象
//
// 注意:
//   - 该方法会根据上下文和超时时间来创建exec.Cmd对象.
//   - 如果上下文设置了超时时间, 则会忽略超时参数.
func (c *Command) buildExecCmd() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.execCmd != nil {
		return // 已经构建过了
	}

	// 根据实际情况选择创建方式，避免不必要的上下文使用
	if c.userCtx != nil {
		// 用户设置了上下文，使用CommandContext(忽略timeout)
		if c.shellType != ShellNone {
			cmdStr := c.getCmdStr()
			c.execCmd = exec.CommandContext(c.userCtx, c.shellType.String(), c.shellType.shellFlags(), cmdStr)
		} else {
			c.execCmd = exec.CommandContext(c.userCtx, c.name, c.args...)
		}

	} else if c.timeout > 0 {
		// 只设置了超时，创建超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		c.cancel = cancel // 保存cancel函数用于资源清理
		c.userCtx = ctx   // 将内部创建的上下文保存到userCtx，方便错误判断

		if c.shellType != ShellNone {
			cmdStr := c.getCmdStr()
			c.execCmd = exec.CommandContext(ctx, c.shellType.String(), c.shellType.shellFlags(), cmdStr)
		} else {
			c.execCmd = exec.CommandContext(ctx, c.name, c.args...)
		}

	} else {
		// 都没有设置，使用普通的Command(不带上下文)
		if c.shellType != ShellNone {
			cmdStr := c.getCmdStr()
			c.execCmd = exec.Command(c.shellType.String(), c.shellType.shellFlags(), cmdStr)
		} else {
			c.execCmd = exec.Command(c.name, c.args...)
		}
	}

	// 设置exec.Cmd的其他属性
	c.execCmd.Dir = c.dir       // 设置工作目录
	c.execCmd.Env = c.envs      // 设置环境变量
	c.execCmd.Stdin = c.stdin   // 设置标准输入
	c.execCmd.Stdout = c.stdout // 设置标准输出
	c.execCmd.Stderr = c.stderr // 设置标准错误输出
}

// cleanup 清理资源
func (c *Command) cleanup() {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
}

// getCmdStr 获取命令字符串
//
// 参数:
//   - c: 命令对象
//
// 返回:
//   - string: 命令字符串
func (c *Command) getCmdStr() string {
	if c == nil {
		return ""
	}

	// 构建命令字符串
	if c.raw != "" {
		return c.raw
	} else if len(c.args) == 0 {
		return c.name
	} else {
		return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
	}
}
