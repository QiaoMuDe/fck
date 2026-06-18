# Shell 脚本处理子命令（shfmt/shck/shx）Spec

## Why

当前 fck 缺少对 Shell 脚本的处理能力。用户在日常工作中经常需要格式化 Shell 脚本、检查语法错误、以及执行 Shell 命令或脚本。借助 `gitee.com/MM-Q/shx` 库（已引入但尚未使用），可以快速实现这三个常用功能，补齐 fck 在脚本处理方面的空白。

结合 `gitee.com/MM-Q/go-kit/fs` 库的 `fs.ExpandFiles()` 通配符展开能力，支持批量处理多个脚本文件。

## What Changes

### 新增 3 个子命令

1. **shfmt** — Shell 脚本格式化（基于 `shx.Format()` / `shx.FormatScript()`）
2. **shck** — Shell 脚本语法检查（基于 `shx.CheckSyntax()` / `shx.CheckScriptSyntax()`）
3. **shx** — Shell 命令/脚本执行（基于 `shx.Run()` / `shx.Out()` / `shx.RunScript()` 等）

### 架构变更

- `internal/cli/root.go` — 注册 3 个新子命令到 `SubCmds`
- `internal/cli/shfmt.go` — shfmt 的 CLI 接口定义
- `internal/cli/shck.go` — shck 的 CLI 接口定义
- `internal/cli/shx.go` — shx 的 CLI 接口定义
- `internal/commands/shfmt/cmd_shfmt.go` — shfmt 的业务逻辑
- `internal/commands/shck/cmd_shck.go` — shck 的业务逻辑
- `internal/commands/shx/cmd_shx.go` — shx 的业务逻辑

## Impact

- Affected specs: CLI 接口规范、命令注册规范
- Affected code: `internal/cli/root.go`、新增 6 个文件
- 依赖变更：`gitee.com/MM-Q/shx` 和 `gitee.com/MM-Q/go-kit/fs` 均在 go.mod 中，无需新增依赖

## ADDED Requirements

### Requirement 1: shfmt 命令 — Shell 脚本格式化

#### Scenario: 从文件读取并格式化
- **WHEN** 用户执行 `fck shfmt script.sh`
- **THEN** 格式化后的脚本内容输出到 stdout

#### Scenario: 从管道读取并格式化
- **WHEN** 用户执行 `cat script.sh | fck shfmt`
- **THEN** 格式化后的脚本内容输出到 stdout

#### Scenario: 原地写入
- **WHEN** 用户执行 `fck shfmt -w script.sh`
- **THEN** 直接修改原文件，格式化后的内容写回 script.sh

#### Scenario: 原地写入并备份
- **WHEN** 用户执行 `fck shfmt -w -b script.sh`
- **THEN** 先创建 script.sh.bak 备份，再写入格式化内容

#### Scenario: 通配符匹配批量格式化
- **WHEN** 用户执行 `fck shfmt -w *.sh`
- **THEN** 使用 `fs.ExpandFiles()` 展开通配符，逐个格式化所有匹配的 .sh 文件

#### Scenario: 递归匹配批量格式化
- **WHEN** 用户执行 `fck shfmt -w **/*.sh`
- **THEN** 使用 `fs.ExpandFiles()` 递归展开，格式化当前目录及子目录下所有 .sh 文件

### Requirement 2: shck 命令 — Shell 脚本语法检查

#### Scenario: 语法正确
- **WHEN** 用户执行 `fck shck valid.sh`
- **THEN** 输出 `✓ syntax is valid`，退出码为 0

#### Scenario: 语法错误
- **WHEN** 用户执行 `fck shck invalid.sh`
- **THEN** 输出语法错误信息（行号、列号、错误描述），退出码为非 0

#### Scenario: 从管道读取检查
- **WHEN** 用户执行 `cat script.sh | fck shck`
- **THEN** 对管道输入的脚本内容进行语法检查

#### Scenario: 通配符匹配批量检查
- **WHEN** 用户执行 `fck shck src/*.sh`
- **THEN** 使用 `fs.ExpandFiles()` 展开通配符，逐个列出每个文件的检查结果

### Requirement 3: shx 命令 — Shell 命令/脚本执行

#### Scenario: 执行命令字符串
- **WHEN** 用户执行 `fck shx "echo hello"`
- **THEN** 执行命令并将输出打印到终端

#### Scenario: 执行脚本文件
- **WHEN** 用户执行 `fck shx script.sh`
- **THEN** 执行指定的 .sh 脚本文件

#### Scenario: 设置超时
- **WHEN** 用户执行 `fck shx -t 5s "sleep 10"`
- **THEN** 5 秒后超时退出，返回超时错误

#### Scenario: 设置工作目录
- **WHEN** 用户执行 `fck shx -d /tmp "pwd"`
- **THEN** 在 /tmp 目录下执行 pwd

#### Scenario: 设置环境变量
- **WHEN** 用户执行 `fck shx -e FOO=bar "echo $FOO"`
- **THEN** 输出 `bar`

#### Scenario: 从管道传递 stdin
- **WHEN** 用户执行 `echo "hello" | fck shx "cat"`
- **THEN** 管道内容作为命令的 stdin 输入

## REMOVED Requirements

（无）
