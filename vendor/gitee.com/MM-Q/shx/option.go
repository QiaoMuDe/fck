package shx

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"mvdan.cc/sh/v3/expand"
)

// WithDir 设置工作目录
//
// 参数:
//   - dir: 工作目录路径
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 如果目录不存在或不是目录, 会 panic
func (s *Shx) WithDir(dir string) *Shx {
	// 处理空目录
	if dir == "" {
		dir = "."
	}

	// 验证目录是否存在
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("directory %s does not exist", dir))
		}
		panic(fmt.Sprintf("stat %s failed: %v", dir, err))
	}
	if !info.IsDir() {
		panic(fmt.Sprintf("%s is not a directory", dir))
	}

	// 设置工作目录
	s.Dir = dir

	return s
}

// WithEnv 设置环境变量
//
// 参数:
//   - key: 环境变量名
//   - value: 环境变量值
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 如果 key 为空, 会 panic
func (s *Shx) WithEnv(key, value string) *Shx {
	if key == "" {
		panic("environment variable key cannot be empty")
	}
	s.mergeEnv([]string{fmt.Sprintf("%s=%s", key, value)})
	return s
}

// WithEnvs 批量设置环境变量
//
// 参数:
//   - envs: 环境变量切片, 每个元素格式为 "key=value"
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 格式错误的项会被忽略
//   - 同名的变量, 后出现的会覆盖先出现的
func (s *Shx) WithEnvs(envs []string) *Shx {
	if len(envs) == 0 {
		panic("environment slice cannot be empty")
	}
	s.mergeEnv(envs)
	return s
}

// mergeEnv 合并环境变量到当前环境
//
// 参数:
//   - newEnvs: 新环境变量切片, 每个元素格式为 "key=value"
func (s *Shx) mergeEnv(newEnvs []string) {
	newMap := make(map[string]string, len(newEnvs))
	for _, env := range newEnvs {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && parts[0] != "" {
			newMap[parts[0]] = parts[1]
		}
	}

	var envList []string
	s.Env.Each(func(name string, vr expand.Variable) bool {
		if _, exists := newMap[name]; !exists {
			envList = append(envList, fmt.Sprintf("%s=%s", name, vr.String()))
		}
		return true
	})

	for key, value := range newMap {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	s.Env = expand.ListEnviron(envList...)
}

// WithStdin 设置标准输入
//
// 参数:
//   - r: 输入读取器
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
func (s *Shx) WithStdin(r io.Reader) *Shx {
	s.Stdin = r
	return s
}

// WithStdout 设置标准输出
//
// 参数:
//   - w: 输出写入器
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
func (s *Shx) WithStdout(w io.Writer) *Shx {
	s.Stdout = w
	return s
}

// WithStderr 设置标准错误
//
// 参数:
//   - w: 错误输出写入器
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
func (s *Shx) WithStderr(w io.Writer) *Shx {
	s.Stderr = w
	return s
}

// WithTimeout 设置超时时间
//
// 参数:
//   - d: 超时时间
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 如果 d <= 0, 则忽略 (不设置超时)
func (s *Shx) WithTimeout(d time.Duration) *Shx {
	if d > 0 {
		s.Timeout = d
	}

	return s
}

// WithContext 设置上下文
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 设置的上下文会完全覆盖 WithTimeout 设置的超时
func (s *Shx) WithContext(ctx context.Context) *Shx {
	s.Ctx = ctx
	return s
}
