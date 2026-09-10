# 枚举帮助信息接入 qflag.EnumHelp 替换计划

## Summary
将项目中手写多行方框枚举帮助（A 类）统一替换为 qflag 库便捷函数 `qflag.EnumHelp(desc, options, indent)`。涉及 10 处枚举：7 处 `table-style`、pack 的 `compression`、pack/unpack 的 `progress-style`。描述文案集中抽取为 `types` 共享切片，避免 7 处 table-style 重复维护。

## Current State Analysis
- qflag v0.5.21 已提供 `qflag.EnumHelp(desc string, options []string, indent string) string`，支持 "值" 与 "值: 描述" 两种形态，按终端显示宽度对齐方框。
- doc2md.go 已接入该函数，作现有范式参考。
- 盘点 `internal/cli/` 全部 `.Enum(` 调用，分三类：
  - **A 类【可替换】**：手写多行 `[值]   - 描述` 方框列表。
    - 7 处 `table-style`：df.go:30、ifconfig.go:29、list.go:53、port.go:35、proc.go:39、size.go:30、wc.go:32（allowedValues 均来自 `types.TableStyles`）
    - pack.go:39 `compression`（`types.SupportedCompressionLevels`）
    - pack.go:47、unpack.go:35 `progress-style`（`types.SupportedProgressStyles`）
  - **B/C 类【不替换】**：`fmt.Sprintf("%v", list)` 单行或简短文本（curl.go request、dns.go type/format、hash.go type、md.go style、json.go type、proc.go sort、iconv.go from/to、newline.go to），改动会改变既有简洁样式。
  - **别名形态【不替换】**：find.go type、list.go type，`[f | file]` 为长短别名，allowedValues 含长短两份，EnumHelp 逐个 `[值]` 展示不适用。

## Proposed Changes

### 1. `internal/types/format.go`：新增 `TableStyleOptions`
新增共享切片，值+描述与现有手写 desc 文案一致（20 项）：
```go
// TableStyleOptions 表格样式的帮助选项（值: 描述），供 qflag.EnumHelp 生成 enum 帮助
var TableStyleOptions = []string{
    "def:  默认样式",
    "l:    浅色样式",
    "r:    圆角样式",
    "bd:   粗体样式",
    "cb:   亮色彩色样式",
    "cd:   暗色彩色样式",
    "db:   双线样式",
    "cbb:  黑色背景蓝色字体",
    "cbc:  青色背景蓝色字体",
    "cbg:  绿色背景蓝色字体",
    "cbm:  紫色背景蓝色字体",
    "cby:  黄色背景蓝色字体",
    "cbr:  红色背景蓝色字体",
    "cwb:  蓝色背景白色字体",
    "ccw:  青色背景白色字体",
    "cgw:  绿色背景白色字体",
    "cmw:  紫色背景白色字体",
    "crw:  红色背景白色字体",
    "cyw:  黄色背景白色字体",
    "none: 禁用边框样式",
}
```
> 注：`qflag.EnumHelp` 会自动按方框内最大值对齐，故 desc 后不再需要手写对齐空格；值均为字母/短码，无冒号冲突。

### 2. `internal/types/compress_type.go`：新增 `CompressionLevelOptions` 与 `ProgressStyleOptions`
```go
// CompressionLevelOptions 压缩级别的帮助选项（值: 描述）
var CompressionLevelOptions = []string{
    "default: 默认压缩级别",
    "none:    不压缩",
    "fast:    快速压缩",
    "best:    最佳压缩",
    "huffman: huffman 压缩",
}

// ProgressStyleOptions 进度条样式的帮助选项（值: 描述）
var ProgressStyleOptions = []string{
    "text:    文本样式",
    "default: 默认样式",
    "unicode: unicode 样式",
    "ascii:   ascii 样式",
}
```

### 3. 替换 7 处 `table-style` 手写 desc
逐个将 `df.go`、`ifconfig.go`、`list.go`、`port.go`、`proc.go`、`size.go`、`wc.go` 中 `"指定表格样式，支持以下选项：\n" + 多行方框` 整段，替换为：
```go
qflag.EnumHelp("指定表格样式，支持以下选项:", types.TableStyleOptions, "\t\t\t\t\t")
```
- 保留原默认值与 `types.TableStyles` 参数不变（如 `"none", types.TableStyles`）。
- 各文件 `internal/types` 均已 import（现有 `types.TableStyles` 已引用）。

### 4. 替换 pack.go `compression`
`pack.go` L39-44 手写段 → `qflag.EnumHelp("压缩级别，支持以下选项:", types.CompressionLevelOptions, "\t\t\t\t\t")`，保留默认值 `types.CompressionLevelDefault` 与 `types.SupportedCompressionLevels`。

### 5. 替换 pack.go / unpack.go `progress-style`
`pack.go` L47-51、`unpack.go` L35-39 手写段 → `qflag.EnumHelp("进度条样式，支持以下选项:", types.ProgressStyleOptions, "\t\t\t\t\t")`，保留原默认值（`types.ProgressStyleAscii`）与 `types.SupportedProgressStyles`。

> 每个 Enum 调用处 desc 末尾不再手写换行；`EnumHelp` 返回含规整换行且末行接 default 值。

## Assumptions & Decisions
- **范围限定 A 类**：仅替换手写多行方框帮助；find/list 的 `--type`（别名形态）与 B/C 单行内联（Sprintf）保持现状，避免破坏既有简洁样式。
- **描述文案统一**：table-style 各文件手写文案（如 df 的"亮色彩色样式/暗色彩色样式"）存在细微差异，统一收敛为 `TableStyleOptions` 中的一套文案。
- **描述集中抽为 types 切片**：7 处 table-style 复用同一 `TableStyleOptions`，避免就地重复 20 行；`options` 仅用于生成帮助文案，`allowedValues` 仍沿用现有切片（值为真实合法输入）。
- **不加别名字段、不重构 find/list type**：超出本任务范围。

## Verification
1. `go build ./...` 编译通过（需联网下载依赖，沙箱受限时在本地或允许沙箱外执行）。
2. `go vet ./internal/cli/ ./internal/types/` 无新增告警。
3. 运行 `fck df --help`、`fck list --help`、`fck pack --help`、`fck unpack --help`，抽查 `table-style`、`compression`、`progress-style` 帮助区块：方框值左对齐、`- 描述` 对齐、`(default: xxx)` 接末行、无多余空行，样式与其他已用 EnumHelp 的命令一致。
4. 确认 find/list `--type` 帮助输出未受影响（应保持原样）。