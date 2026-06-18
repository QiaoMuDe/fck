# SHX 工具链功能增强 Spec

## Why

对已有的 3 个 SHX 工具链子命令（shfmt/shck/shx）进行功能增强，补齐 CI/脚本场景下的关键缺口：shck 缺静默模式、shfmt 缺列表模式、shx 退出码未透传。

## What Changes

1. **shck** — 新增 `-q`/`--quiet` 静默模式：成功时不输出，仅显示错误
2. **shfmt** — 新增 `-l`/`--list` 列表模式：仅列出格式化前后有差异的文件，不修改文件
3. **shx** — 退出码透传：命令执行失败时，fck 进程退出码与被执行命令的退出码一致

## Impact

- Affected specs: shfmt 命令、shck 命令、shx 命令
- Affected code:
  - `internal/cli/shck.go`
  - `internal/commands/shck/cmd_shck.go`
  - `internal/cli/shfmt.go`
  - `internal/commands/shfmt/cmd_shfmt.go`
  - `internal/cli/shx.go`
  - `internal/commands/shx/cmd_shx.go`

## ADDED Requirements

### Requirement: shck -q 静默模式

#### Scenario: 单文件正确，静默模式
- **WHEN** 用户执行 `fck shck -q valid.sh`
- **THEN** 不输出任何内容，退出码为 0

#### Scenario: 单文件错误，静默模式
- **WHEN** 用户执行 `fck shck -q invalid.sh`
- **THEN** 仅输出错误信息（含行号列号），退出码非 0

#### Scenario: 批量文件，部分错误，静默模式
- **WHEN** 用户执行 `fck shck -q src/*.sh`
- **THEN** 仅输出有语法错误的文件错误信息，正确文件不输出任何内容，退出码非 0

#### Scenario: 管道输入，正确，静默模式
- **WHEN** 用户执行 `cat valid.sh | fck shck -q`
- **THEN** 不输出任何内容，退出码为 0

### Requirement: shfmt -l 列表模式

#### Scenario: 单文件需要格式化
- **WHEN** 用户执行 `fck shfmt -l script.sh` 且 script.sh 格式与格式化后有差异
- **THEN** 输出 `script.sh`（文件路径），不修改文件

#### Scenario: 单文件已经格式化
- **WHEN** 用户执行 `fck shfmt -l script.sh` 且 script.sh 已经是格式化状态
- **THEN** 不输出任何内容，不修改文件

#### Scenario: 批量模式，列举需要格式化的文件
- **WHEN** 用户执行 `fck shfmt -l *.sh`
- **THEN** 逐个检查每个文件，仅输出格式化后有差异的文件路径（每行一个）

#### Scenario: -l 与 -w 组合
- **WHEN** 用户执行 `fck shfmt -l -w script.sh`
- **THEN** -l 优先级更高，只列表不写入，忽略 -w

### Requirement: shx 退出码透传

#### Scenario: 命令成功执行
- **WHEN** 用户执行 `fck shx "exit 0"`
- **THEN** 进程退出码为 0

#### Scenario: 命令失败，非零退出码
- **WHEN** 用户执行 `fck shx "exit 1"`
- **THEN** 不输出额外错误信息，进程退出码为 1

#### Scenario: 命令失败，特殊退出码
- **WHEN** 用户执行 `fck shx "exit 42"`
- **THEN** 不输出额外错误信息，进程退出码为 42

#### Scenario: 执行超时
- **WHEN** 用户执行 `fck shx -t 1s "sleep 10"`
- **THEN** 输出超时错误信息，进程退出码非 0（由 qflag 框架决定）

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）
