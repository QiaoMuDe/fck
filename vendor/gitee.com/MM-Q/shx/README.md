<div align="center">

# Shx

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org) [![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](https://gitee.com/MM-Q/shx/blob/master/LICENSE) [![Gitee](https://img.shields.io/badge/Gitee-Repository-red?style=for-the-badge&logo=gitee)](https://gitee.com/MM-Q/shx)

**一个基于 mvdan.cc/sh/v3 的纯 Go Shell 命令执行库**

[🏠 仓库地址](https://gitee.com/MM-Q/shx) • [📖 API 文档](APIDOC.md)

</div>

---

## 项目简介

Shx 是一个基于 [mvdan.cc/sh/v3](https://mvdan.cc/sh/v3) 的纯 Go Shell 命令执行库，不依赖系统 shell，提供跨平台一致的命令执行体验。

- **纯 Go 实现** — 不依赖系统 shell，Windows/Linux/macOS 行为一致
- **Bash 方言默认** — 支持 `[[ ]]`、`function`、`select` 等 Bash 特有语法
- **脚本文件执行** — 原生支持执行 `.sh` 脚本文件
- **链式调用 API** — 流畅的方法链，支持工作目录、环境变量、超时等配置

---

## 安装

```bash
go get gitee.com/MM-Q/shx@latest
```

Go 版本要求：1.25.0 或更高

---

## 核心特性

| 特性 | 说明 |
|------|------|
| **纯 Go 实现** | 基于 mvdan.cc/sh/v3，不依赖系统 shell |
| **Bash 方言默认** | 默认使用 Bash 方言解析器，支持 `[[ ]]`、`function`、`select` 等 Bash 特有语法 |
| **脚本文件执行** | 原生支持执行 `.sh` 脚本文件 |
| **链式调用** | 支持流畅的方法链配置 |
| **超时控制** | 支持上下文超时和超时参数 |
| **输入输出重定向** | 灵活的标准输入输出配置 |
| **退出码检测** | `IsExitStatus` 可提取命令退出码 |
| **结构体字面量** | 导出字段支持直接通过结构体配置 |

---

## 使用示例

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gitee.com/MM-Q/shx"
)

func main() {
    // ---- 便捷函数 ----
    shx.Run("echo hello")                        // 简单执行
    shx.RunToTerminal("ls -la")                  // 输出到终端

    output, _ := shx.Out("date")                 // 获取输出
    fmt.Println(string(output))

    shx.RunWith("sleep 10", 5*time.Second)       // 超时执行

    // ---- 链式配置 ----
    out, err := shx.New("echo hello").
        WithTimeout(5 * time.Second).
        WithDir("/tmp").
        WithEnv("FOO", "bar").
        ExecOutput()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))

    // ---- 使用上下文 ----
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = shx.New("long-command").WithContext(ctx).Exec()

    // ---- 自定义输入输出 ----
    var stdout, stderr bytes.Buffer
    stdin := strings.NewReader("input")
    err = shx.New("cat").
        WithStdin(stdin).
        WithStdout(&stdout).
        WithStderr(&stderr).
        Exec()

    // ---- 执行脚本文件 ----
    shx.RunScript("deploy.sh")
    output, err = shx.OutScript("build.sh")

    out, err = shx.NewScript("test.sh").
        WithDir("/project").
        WithEnv("MODE", "ci").
        WithTimeout(30 * time.Second).
        ExecOutput()

    // ---- 检查退出码 ----
    err = shx.Run("exit 5")
    if code, ok := shx.IsExitStatus(err); ok {
        fmt.Printf("Exit code: %d\n", code)
    }
}
```

---

## 注意事项

- Shx 对象的配置方法（WithXxx）不是并发安全的，不要在多个 goroutine 中并发配置
- mvdan/sh 是同步执行的，不提供异步 API，如需异步请使用 goroutine 包装
- 不支持进程控制（无 PID、Kill、Signal），只能通过 context 取消
- 默认使用 Bash 方言解析（`syntax.LangBash`），支持 `[[ ]]`、`function`、`select` 等 Bash 特有语法

---

<div align="center">

**如果这个项目对您有帮助，请给它一个 ⭐ Star！**

[⬆顶部](#shx)

</div>
