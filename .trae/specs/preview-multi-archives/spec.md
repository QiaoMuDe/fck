# Preview 命令支持通配符与批量预览 Spec

## Why

当前 `preview/pv` 命令仅支持单个压缩包预览（取 `args[0]`），用户无法使用通配符批量预览多个压缩包，需要手动逐个查看。

## What Changes

- **CLI 层**：`internal/cli/preview.go` 调用 `fs.ExpandFiles` 展开通配符，传递多个路径到命令层
- **命令层**：`internal/commands/preview/cmd_preview.go` 将 `PackPath string` 改为 `PackPaths []string`，循环处理每个压缩包，多个时自动显示文件名标题

## Impact

- Affected specs: 压缩包预览
- Affected code:
  - `internal/cli/preview.go` — 通配符展开，多参数传递
  - `internal/commands/preview/cmd_preview.go` — 支持多压缩包批量预览，多文件时打印分隔标题

## ADDED Requirements

### Requirement: 批量预览支持

Preview 命令 SHALL 支持通过通配符匹配多个压缩包并逐个预览。

#### Scenario: 通配符匹配多个文件
- **WHEN** 用户输入 `fck preview -i *.zip`
- **THEN** 展开通配符匹配到多个 .zip 文件
- **THEN** 逐个打印每个压缩包的信息，之间用 `==> filename <==` 分隔

#### Scenario: 单个文件（向后兼容）
- **WHEN** 用户输入 `fck preview archive.zip`
- **THEN** 行为与当前一致，不显示文件名标题

#### Scenario: 管道输入
- **WHEN** 用户通过管道输入
- **THEN** 行为保持与当前一致（当前不支持管道）

#### Scenario: 无匹配文件
- **WHEN** 通配符无匹配且没有参数
- **THEN** 返回错误 `"missing archive path"`

#### Scenario: 通配符全部无匹配
- **WHEN** 通配符匹配后列表为空
- **THEN** 返回错误 `"no archive files found"`
