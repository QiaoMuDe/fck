# 命令通配符支持 Spec

## Why

当前 `head`、`hash`、`size`、`touch`、`truncate`、`iconv`、`json` 7 个命令在 CLI 层直接传递 `cmd.Args()` 到命令层，未对通配符进行展开，用户无法使用 `*.log`、`*.txt` 等模式匹配文件。

这些命令的共同点：已经支持多个文件参数，命令层内部本身就能处理多文件，缺的只是在 CLI 层做一次通配符展开。

## What Changes

每个文件改动模式一致：
1. 导入 `gitee.com/MM-Q/go-kit/fs`
2. 调用 `fs.ExpandFiles(args)` 展开通配符（只保留文件，过滤目录）
3. 将展开后的结果传给命令层 config

## Impact

- Affected code（7 个 CLI 文件，仅 CLI 层改动，命令层无需修改）：
  - `internal/cli/head.go`
  - `internal/cli/hash.go`
  - `internal/cli/size.go`
  - `internal/cli/touch.go`
  - `internal/cli/truncate.go`
  - `internal/cli/iconv.go`
  - `internal/cli/json.go`

## ADDED Requirements

### Requirement: 通配符展开

这 7 个命令 SHALL 在 CLI 层接收位置参数后，先调用 `fs.ExpandFiles` 通配符展开，再将展开结果传递给命令层。

#### Scenario: 正常匹配
- **WHEN** 用户输入 `fck head *.txt`
- **THEN** `*.txt` 展开为匹配的文件列表
- **THEN** 命令层接收展开后的文件列表正常处理

#### Scenario: 无匹配项
- **WHEN** 通配符表达式无匹配文件
- **THEN** `fs.ExpandFiles` 返回原表达式字符串
- **THEN** 命令层自行处理（如 `os.Open` 失败时返回文件不存在错误）

#### Scenario: 混合参数
- **WHEN** 用户输入 `fck hash file.txt *.iso`
- **THEN** 通配符展开后与具体路径合并
