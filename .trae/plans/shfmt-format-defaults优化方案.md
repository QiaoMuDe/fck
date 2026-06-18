# shfmt 格式化选项默认值优化方案

## 方案

在 `init()` 注册标志时，通过 `shx.DefaultFormatOptions()` 获取默认值设置标志的默认值，消除硬编码。

### 变更

**`internal/cli/shfmt.go`** — 仅改此文件：
- `init()` 中：`def := shx.DefaultFormatOptions()`，三个标志用 `def.Indent`、`def.SwitchCaseIndent`、`def.KeepComments` 作为默认值
- Notes 中移除硬编码的 "2 空格缩进" 描述，改为引用库函数

### 验证
- [x] `go build ./...` 编译通过
- [x] `golangci-lint` 无新增错误
