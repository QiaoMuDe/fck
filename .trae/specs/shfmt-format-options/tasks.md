# Tasks

- [x] Task 1: 更新 `ShfmtConfig` 结构体 — 新增 `FormatOptions shx.FormatOptions` 字段
  - [x] 1.1 在 `internal/commands/shfmt/cmd_shfmt.go` 的 `ShfmtConfig` 中添加 `FormatOptions` 字段
  - [x] 1.2 当未设置时使用 `shx.DefaultFormatOptions()` 作为默认值（通过 CLI 层默认值保证）

- [x] Task 2: 替换格式化调用 — 将 `shx.Format()` 替换为 `shx.FormatWithOptions()`
  - [x] 2.1 修改 `processStdin()` 中的调用
  - [x] 2.2 修改 `processFile()` 中的调用

- [x] Task 3: 新增 CLI 标志 — 在 `internal/cli/shfmt.go` 中注册 8 个格式化选项标志
  - [x] 3.1 `--indent` / `-i` (uint, 默认 2)
  - [x] 3.2 `--switch-case-indent` (bool, 默认 true)
  - [x] 3.3 `--keep-comments` (bool, 默认 true)
  - [x] 3.4 `--binary-next-line` (bool, 默认 false)
  - [x] 3.5 `--function-next-line` (bool, 默认 false)
  - [x] 3.6 `--space-redirects` (bool, 默认 false)
  - [x] 3.7 `--single-line` (bool, 默认 false)
  - [x] 3.8 `--minify` / `-m` (bool, 默认 false)

- [x] Task 4: 更新 `runShfmt` — 将 CLI 标志值映射到 `ShfmtConfig.FormatOptions`
  - [x] 4.1 构建 `shx.FormatOptions` 实例
  - [x] 4.2 传递给 `shfmt.ShfmtCmdMain()`

- [x] Task 5: 编译验证 — 确保项目编译通过
  - [x] 5.1 `go build ./...` 编译无错误
  - [x] 5.2 `golangci-lint` lint 无错误

# Task Dependencies

- Task 3/4 依赖 Task 1 (配置结构体先行)
- Task 2 依赖 Task 1 (需要 FormatOptions 字段)
- Task 5 依赖所有前置 Task
