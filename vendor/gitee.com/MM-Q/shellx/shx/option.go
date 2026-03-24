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
//   - 如果命令已经执行过, 会 panic
//   - 如果目录不存在或不是目录, 会 panic
func (s *Shx) WithDir(dir string) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

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
	s.dir = dir

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
//   - 如果命令已经执行过, 会 panic
//   - 如果 key 为空, 则忽略
func (s *Shx) WithEnv(key, value string) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if key == "" {
		panic("environment variable key cannot be empty")
	}

	// 获取当前环境变量列表
	var envList []string
	s.env.Each(func(name string, vr expand.Variable) bool {
		// 剔除已存在的环境变量, 保留其他环境变量
		if name != key {
			envList = append(envList, fmt.Sprintf("%s=%s", name, vr.String()))
		}
		return true
	})

	// 添加新的环境变量
	envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	s.env = expand.ListEnviron(envList...)

	return s
}

// WithEnvMap 批量设置环境变量
//
// 参数:
//   - envs: 环境变量映射 (key-value)
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 如果命令已经执行过, 会 panic
func (s *Shx) WithEnvMap(envs map[string]string) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if len(envs) == 0 {
		panic("environment map cannot be empty")
	}

	// 获取当前环境变量列表
	envList := make([]string, 0)
	s.env.Each(func(name string, vr expand.Variable) bool {
		envList = append(envList, fmt.Sprintf("%s=%s", name, vr.String()))
		return true
	})

	// 添加新的环境变量 (会覆盖已有的)
	for key, value := range envs {
		if key != "" {
			envList = append(envList, fmt.Sprintf("%s=%s", key, value))
		}
	}

	s.env = expand.ListEnviron(envList...)
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
//   - 如果命令已经执行过, 会 panic
//   - 格式错误的项会被忽略
//   - 同名的变量, 后出现的会覆盖先出现的
func (s *Shx) WithEnvs(envs []string) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if len(envs) == 0 {
		panic("environment slice cannot be empty")
	}

	// 解析新环境变量为 map 用于去重和覆盖
	newEnvs := make(map[string]string, len(envs))
	for _, env := range envs {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && parts[0] != "" {
			newEnvs[parts[0]] = parts[1]
		}
	}

	// 合并: 保留旧变量, 新变量覆盖同名的
	var envList []string
	s.env.Each(func(name string, vr expand.Variable) bool {
		if _, exists := newEnvs[name]; !exists {
			envList = append(envList, fmt.Sprintf("%s=%s", name, vr.String()))
		}
		return true
	})

	// 添加新变量
	for key, value := range newEnvs {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	s.env = expand.ListEnviron(envList...)
	return s
}

// WithStdin 设置标准输入
//
// 参数:
//   - r: 输入读取器
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 如果命令已经执行过, 会 panic
func (s *Shx) WithStdin(r io.Reader) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	s.stdin = r
	return s
}

// WithStdout 设置标准输出
//
// 参数:
//   - w: 输出写入器
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 如果命令已经执行过, 会 panic
func (s *Shx) WithStdout(w io.Writer) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	s.stdout = w
	return s
}

// WithStderr 设置标准错误
//
// 参数:
//   - w: 错误输出写入器
//
// 返回:
//   - *Shx: 命令对象 (支持链式调用)
//
// 注意:
//   - 如果命令已经执行过, 会 panic
func (s *Shx) WithStderr(w io.Writer) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	s.stderr = w
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
//   - 如果命令已经执行过, 会 panic
//   - 如果 d <= 0, 则忽略 (不设置超时)
func (s *Shx) WithTimeout(d time.Duration) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	if d > 0 {
		s.timeout = d
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
//   - 如果命令已经执行过, 会 panic
//   - 设置的上下文会完全覆盖 WithTimeout 设置的超时
func (s *Shx) WithContext(ctx context.Context) *Shx {
	if s.executed.Load() {
		panic("shx has already been executed")
	}

	s.ctx = ctx
	return s
}
