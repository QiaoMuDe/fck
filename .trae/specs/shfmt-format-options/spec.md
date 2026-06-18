# shfmt 格式化配置选项 Spec

## Why

当前 `shfmt` 命令使用 `shx.Format()`（即 `DefaultFormatOptions()`）进行格式化，用户无法控制缩进、注释保留、操作符换行等格式行为。shx 库提供了 `FormatWithOptions()` 和 `FormatOptions` 结构体，支持丰富的格式化选项，需要在 CLI 层暴露这些能力。

## What Changes

### CLI 新增标志

在 `internal/cli/shfmt.go` 中新增以下标志：

| 标志 | 别名 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--indent` | `-i` | `uint` | `2` | 缩进空格数（`0` 表示使用 tab） |
| `--switch-case-indent` | — | `bool` | `true` | case 语句体是否缩进 |
| `--keep-comments` | — | `bool` | `true` | 是否保留注释 |
| `--binary-next-line` | — | `bool` | `false` | `&&`、`||` 等二元操作符是否换行显示 |
| `--function-next-line` | — | `bool` | `false` | 函数体 `{` 是否换行 |
| `--space-redirects` | — | `bool` | `false` | 重定向符前后是否加空格 |
| `--single-line` | — | `bool` | `false` | 是否单行输出 |
| `--minify` | `-m` | `bool` | `false` | 是否最小化输出（压缩模式） |

### 业务层变更

`internal/commands/shfmt/cmd_shfmt.go`：

- `ShfmtConfig` 结构体新增 `FormatOptions shx.FormatOptions` 字段
- `ShfmtCmdMain` / `processFile` / `processStdin` 中将 `shx.Format()` 替换为 `shx.FormatWithOptions(data, config.FormatOptions)`

### 默认值

当用户未指定任何格式化标志时，使用 `shx.DefaultFormatOptions()` 作为默认值，保持与现有行为一致。

## Impact

- Affected specs: shfmt CLI 接口、shfmt 业务逻辑
- Affected code:
  - `internal/cli/shfmt.go` — 新增 8 个标志注册
  - `internal/commands/shfmt/cmd_shfmt.go` — ShfmtConfig 新增 FormatOptions 字段，调用 FormatWithOptions
- 无新增依赖

## ADDED Requirements

### Requirement 1: 自定义缩进

#### Scenario: 使用 tab 缩进
- **WHEN** 用户执行 `fck shfmt -i 0 script.sh`
- **THEN** 使用 tab 缩进格式化（`Indent=0` 表示 tab）

#### Scenario: 使用 4 空格缩进
- **WHEN** 用户执行 `fck shfmt -i 4 script.sh`
- **THEN** 使用 4 空格缩进格式化

### Requirement 2: 控制 case 语句缩进

#### Scenario: 禁用 case 语句体缩进
- **WHEN** 用户执行 `fck shfmt --switch-case-indent=false script.sh`
- **THEN** case 语句体不缩进

### Requirement 3: 控制注释保留

#### Scenario: 移除注释
- **WHEN** 用户执行 `fck shfmt --keep-comments=false script.sh`
- **THEN** 格式化后移除所有注释

### Requirement 4: 操作符换行

#### Scenario: 二元操作符换行
- **WHEN** 用户执行 `fck shfmt --binary-next-line script.sh`
- **THEN** `&&`、`||` 等操作符在换行后显示

#### Scenario: 函数体换行
- **WHEN** 用户执行 `fck shfmt --function-next-line script.sh`
- **THEN** 函数体 `{` 在新行上

### Requirement 5: 重定向符空格

#### Scenario: 重定向符前后加空格
- **WHEN** 用户执行 `fck shfmt --space-redirects script.sh`
- **THEN** 重定向符前后添加空格（`cmd > file` 而非 `cmd>file`）

### Requirement 6: 单行/最小化输出

#### Scenario: 单行输出
- **WHEN** 用户执行 `fck shfmt --single-line script.sh`
- **THEN** 格式化后脚本为单行

#### Scenario: 最小化输出
- **WHEN** 用户执行 `fck shfmt -m script.sh`
- **THEN** 格式化后移除所有多余空格和换行（压缩模式）

### Requirement 7: 多个选项组合

#### Scenario: 组合多个标志
- **WHEN** 用户执行 `fck shfmt -i 4 -m -w script.sh`
- **THEN** 使用 4 空格缩进 + 最小化模式原地写入

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）
