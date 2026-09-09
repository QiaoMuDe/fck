# 新增 doc2md 文档转 Markdown 命令 Spec

## Why

fck 目前缺少文档转换能力，无法把 Word/Excel/PPT/PDF 等办公文档快速转换为 Markdown。作者自研的 `gitee.com/MM-Q/doc2md` 库已提供完整的「文档 → Markdown」转换能力（12 类格式），但在 fck 中尚无对应命令入口。新增 `doc2md` 子命令将办公文档转换能力整合进 fck，复用项目现有的命令分层架构。

## What Changes

- 新增 `internal/commands/doc2md/` 业务实现目录
- 新增 `internal/cli/doc2md.go` 命令定义层文件（qflag 参数声明 + run 函数）
- 在 `internal/cli/root.go` 注册 `Doc2mdCmd` 到 `SubCmds`
- 新增外部依赖 `gitee.com/MM-Q/doc2md`（自研库，非第三方）
- 支持输入：本地文件路径 + stdin 管道（`-x` 提示扩展名），默认输出到 stdout，`-o` 指定输出文件

## Impact

- Affected specs: `internal/commands/#` 新增 doc2md 业务子目录
- Affected code:
  - `internal/commands/doc2md/cmd_doc2md.go` — 主业务逻辑（new + 引擎创建 + 输入分流 + 转换 + 输出）
  - `internal/commands/doc2md/types.go` — `Doc2mdConfig` 配置结构
  - `internal/cli/doc2md.go` — CLI 标志定义与命令注册
  - `internal/cli/root.go` — 添加 `Doc2mdCmd` 到 `SubCmds`
  - `go.mod` — 新增 `gitee.com/MM-Q/doc2md` 依赖

## ADDED Requirements

### Requirement 1: 命令注册

系统 SHALL 提供 `doc2md` 顶层子命令，并注册到根命令的 `SubCmds`，遵循项目固有「命令定义层 + 业务实现层」分层范式。

#### Scenario: 命令可见
- **WHEN** 用户执行 `fck --help`
- **THEN** 帮助列表中出现 `doc2md` 命令及其描述
- **WHEN** 用户执行 `fck doc2md --help`
- **THEN** 显示该命令的参数说明与示例

### Requirement 2: 文件路径转换

系统 SHALL 将指定本地文档转换为 Markdown 并输出。

#### Scenario: 转换单个文件
- **WHEN** 用户执行 `fck doc2md input.docx`
- **THEN** 将 input.docx 转换为 Markdown，输出到 stdout
- **AND** 转换失败时返回非零退出码并打印错误信息

#### Scenario: 通配符展开
- **WHEN** 用户执行 `fck doc2md docs/*.docx`
- **THEN** 使用 `go-kit/fs` 的 `Expand` 展开通配符后逐个转换

### Requirement 3: stdin 管道输入

系统 SHALL 支持从标准输入读取文档内容，兼容 fck 统一的管道范式。

#### Scenario: 管道输入
- **WHEN** 用户的 stdin 被管道输入（`cat input.docx | fck doc2md -x .docx`）且无文件参数
- **THEN** 读取 stdin 内容并结合 `-x` 提示格式进行转换
- **WHEN** stdin 为管道但未提供 `-x/-extension`
- **THEN** 给出明确错误提示，要求提供扩展名

### Requirement 4: 输出行为

系统 SHALL 默认将转换结果输出到 stdout，并支持 `-o` 将结果写入文件。

#### Scenario: 默认 stdout
- **WHEN** 用户未提供 `-o`
- **THEN** 转换后的 Markdown 打印到标准输出

#### Scenario: 写入文件
- **WHEN** 用户提供 `-o output.md`
- **THEN** 转换结果写入该文件，不打印到 stdout

### Requirement 5: 转换选项

系统 SHALL 提供与 doc2md CLI 一致的辅助选项，用于控制转换过程。

#### Scenario: 扩展名/MIME/字符集提示
- **WHEN** 用户提供 `-x .docx` / `-m <mime>` / `-c gbk`
- **THEN** 该提示被传递到 `markitdown.StreamInfo` 参与格式识别

#### Scenario: 保留完整图片数据 URI
- **WHEN** 用户提供 `--keep-data-uris`
- **THEN** 转换引擎启用 `markitdown.WithKeepDataURIs(true)`，输出保留完整 base64 图片数据 URI

### Requirement 6: 错误处理

系统 SHALL 对不可识别格式与转换失败给出可读的错误信息。

#### Scenario: 不支持格式
- **WHEN** 输入文件的格式 doc2md 无法识别
- **THEN** 输出「不支持此文档格式」类错误信息并返回非零退出码

#### Scenario: 转换失败
- **WHEN** doc2md 抛错（如文件损坏）
- **THEN** 错误信息透传，返回非零退出码

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）