# Cat 命令统一查看器设计方案 V2

> **设计日期**: 2026-04-16  
> **版本**: v2.0  
> **状态**: 待评审

---

## 一、核心思想

**"统一读取 → 统一处理 → 高亮判断 → 输出"**

1. **统一读取**: 所有模式都读取整个文件到内存
2. **统一处理**: 根据 CLI 标志处理切片数据
3. **高亮判断**: 最后根据 `-H/--highlight` 决定是否渲染
4. **输出**: 高亮模式先渲染再打印，非高亮模式直接打印

---

## 二、标志处理规则

| 标志 | 说明 | 适用模式 |
|------|------|----------|
| `-n/--number` | 显示所有行号 | 高亮/非高亮 |
| `-b/--number-nonblank` | 显示非空行行号 | 高亮/非高亮 |
| `-E/--show-ends` | 显示行尾`$` | **仅非高亮** |
| `-T/--show-tabs` | 显示制表符为`^I` | **仅非高亮** |
| `-A/--show-all` | 等价于`-ET` | **仅非高亮** |
| `-u/--head` | 显示前N行 | 高亮/非高亮 |
| `-d/--tail` | 显示后N行 | 高亮/非高亮 |
| `-H/--highlight` | 启用语法高亮 | 控制渲染 |

**注意**: `-N/--show-newline` 标志已移除

---

## 三、处理流程

```
┌─────────────────────────────────────────────────────────────┐
│  1. 读取文件                                                  │
│     content, err := os.ReadFile(path)                        │
└───────────────────────┬─────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  2. 统一换行符                                               │
│     content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n")) │
└───────────────────────┬─────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  3. 按行分割                                                 │
│     lines := bytes.Split(content, []byte("\n"))              │
└───────────────────────┬─────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  4. 根据 head/tail 切片                                      │
│     if head > 0: lines = lines[:head]                        │
│     if tail > 0: lines = lines[len(lines)-tail:]             │
└───────────────────────┬─────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  5. 处理特殊显示标志（仅非高亮模式）                          │
│     -E: 在行尾添加 "$"                                       │
│     -T: 将 "\t" 替换为 "^I"                                  │
└───────────────────────┬─────────────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────────────┐
│  6. 判断是否高亮                                             │
│     if config.Highlight {                                    │
│         → 高亮渲染后输出                                      │
│     } else {                                                 │
│         → 直接输出                                            │
│     }                                                        │
└─────────────────────────────────────────────────────────────┘
```

---

## 四、核心数据结构

```go
// FileViewer 统一文件查看器
type FileViewer struct {
	config   *CatConfig
	hlConfig HighlightConfig
}

// NewFileViewer 创建文件查看器
func NewFileViewer(config *CatConfig) *FileViewer {
	return &FileViewer{
		config:   config,
		hlConfig: DefaultHighlightConfig(),
	}
}
```

---

## 五、核心方法实现

### 5.1 统一查看入口

```go
// View 查看文件（统一入口）
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - error: 处理错误
func (v *FileViewer) View(path string) error {
	// 1. 读取文件
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// 2. 统一换行符
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

	// 3. 按行分割
	allLines := bytes.Split(content, []byte("\n"))

	// 4. 根据 head/tail 切片
	lines := v.sliceLines(allLines)

	// 5. 处理特殊显示标志（仅非高亮模式）
	if !v.config.Highlight {
		lines = v.applySpecialFlags(lines)
	}

	// 6. 输出（根据是否高亮决定）
	if v.config.Highlight {
		return v.outputHighlighted(lines, path)
	}
	return v.outputPlain(lines)
}
```

### 5.2 切片处理

```go
// sliceLines 根据 head/tail 标志切片
//
// 参数:
//   - lines: 所有行
//
// 返回:
//   - [][]byte: 切片后的行
func (v *FileViewer) sliceLines(lines [][]byte) [][]byte {
	// head 模式
	if v.config.HeadLines > 0 {
		if v.config.HeadLines < len(lines) {
			return lines[:v.config.HeadLines]
		}
		return lines
	}

	// tail 模式
	if v.config.TailLines > 0 {
		if v.config.TailLines < len(lines) {
			start := len(lines) - v.config.TailLines
			return lines[start:]
		}
		return lines
	}

	// 全部
	return lines
}
```

### 5.3 特殊标志处理（仅非高亮模式）

```go
// applySpecialFlags 应用特殊显示标志
//
// 参数:
//   - lines: 行数据
//
// 返回:
//   - [][]byte: 处理后的行
//
// 说明:
//   - 仅在高亮模式为 false 时调用
//   - 处理 -E(行尾$), -T(制表符^I)
func (v *FileViewer) applySpecialFlags(lines [][]byte) [][]byte {
	if !v.config.ShowEnd && !v.config.ShowTabs {
		return lines
	}

	result := make([][]byte, len(lines))
	for i, line := range lines {
		var parts []string

		// 处理 -T: 制表符显示为 ^I
		content := string(line)
		if v.config.ShowTabs {
			content = strings.ReplaceAll(content, "\t", "^I")
		}
		parts = append(parts, content)

		// 处理 -E: 显示行尾 $
		if v.config.ShowEnd {
			parts = append(parts, "$")
		}

		result[i] = []byte(strings.Join(parts, ""))
	}

	return result
}
```

### 5.4 普通输出

```go
// outputPlain 普通模式输出
//
// 参数:
//   - lines: 行数据
//
// 返回:
//   - error: 输出错误
func (v *FileViewer) outputPlain(lines [][]byte) error {
	lineNum := 0

	for _, line := range lines {
		content := string(line)
		isBlank := len(strings.TrimSpace(content)) == 0

		// 处理行号 (-n 或 -b)
		if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
			lineNum++
			fmt.Printf("%6d\t", lineNum)
		}

		// 输出内容
		fmt.Println(content)
	}

	return nil
}
```

### 5.5 高亮输出

```go
// outputHighlighted 高亮模式输出
//
// 参数:
//   - lines: 行数据
//   - path: 文件路径（用于检测语言）
//
// 返回:
//   - error: 输出错误
func (v *FileViewer) outputHighlighted(lines [][]byte, path string) error {
	// 1. 合并内容
	content := bytes.Join(lines, []byte("\n"))

	// 2. 语法高亮
	hlContent, ok, err := highlightContent(content, path, v.hlConfig)
	if err != nil || !ok {
		// 高亮失败，降级为普通输出
		return v.outputPlain(lines)
	}

	// 3. 按行分割高亮后的内容
	hlLines := bytes.Split(hlContent, []byte("\n"))

	// 4. 输出行（带行号）
	lineNum := 0

	for _, line := range hlLines {
		content := string(line)
		// 移除 ANSI 码后判断是否为空行
		isBlank := len(strings.TrimSpace(stripANSI(content))) == 0

		// 处理行号 (-n 或 -b)
		if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
			lineNum++
			fmt.Printf("%6d\t", lineNum)
		}

		// 输出高亮内容
		fmt.Println(content)
	}

	return nil
}
```

---

## 六、CatConfig 更新

移除 `ShowNewline` 字段：

```go
// CatConfig cat 命令配置
type CatConfig struct {
	// CLI 参数
	Targets      []string // 目标文件列表
	ShowLineNum  bool     // -n 显示所有行号
	ShowNonBlank bool     // -b 显示非空行行号
	ShowEnd      bool     // -E 显示行尾$
	ShowTabs     bool     // -T 显示制表符为^I
	ShowAll      bool     // -A 等价于 -ET
	// ShowNewline  bool     // -N 显示换行符类型 [已移除]
	HeadLines    int      // --head 显示前N行 (0表示全部)
	TailLines    int      // --tail 显示后N行 (0表示全部)
	Quiet        bool     // -q 静默模式 (不显示错误信息)
	Text         bool     // -a, --text 强制将二进制文件视为文本处理
	IgnoreBinary bool     // -I, --ignore-binary 完全忽略二进制文件
	UseLess      bool     // -l, --less 使用分页器查看文件内容
	Highlight    bool     // -H, --highlight 启用语法高亮

	// 运行时
	LineCounter int // 行号计数器
}
```

---

## 七、CLI 标志更新

移除 `catShowNewline` 标志：

```go
var (
	catShowLineNum  *qflag.BoolFlag // -n 显示所有行号
	catShowNonBlank *qflag.BoolFlag // -b 显示非空行行号
	catShowEnd      *qflag.BoolFlag // -E 显示行尾$
	catShowTabs     *qflag.BoolFlag // -T 显示制表符为^I
	catShowAll      *qflag.BoolFlag // -A 等价于 -ET
	// catShowNewline  *qflag.BoolFlag // -N 显示换行符类型 [已移除]
	catHeadLines    *qflag.IntFlag  // -u 显示前N行
	catTailLines    *qflag.IntFlag  // -d 显示后N行
	catQuiet        *qflag.BoolFlag // -q 静默模式
	catText         *qflag.BoolFlag // -a, --text 强制文本模式
	catIgnoreBinary *qflag.BoolFlag // -I, --ignore-binary 忽略二进制文件
	catUseLess      *qflag.BoolFlag // -l, --less 使用分页器
	catHighlight    *qflag.BoolFlag // -H, --highlight 启用语法高亮
)
```

---

## 八、调用流程

### 8.1 重构后的 processFile

```go
func processFile(path string, config *CatConfig) error {
	// 1. 基础检查（目录、二进制）
	if err := checkFile(path, config); err != nil {
		return err
	}

	// 2. 创建统一查看器
	viewer := NewFileViewer(config)

	// 3. 统一查看（内部处理 head/tail/高亮等逻辑）
	return viewer.View(path)
}
```

### 8.2 完整调用链

```
CatCmdMain
  └── processFile(path, config)
        ├── checkFile(path, config)     // 检查目录、二进制
        ├── viewer := NewFileViewer(config)
        └── viewer.View(path)
              ├── os.ReadFile(path)      // 统一读取
              ├── bytes.ReplaceAll()     // 统一换行符
              ├── bytes.Split()          // 按行分割
              ├── sliceLines()           // head/tail 切片
              ├── applySpecialFlags()    // 特殊标志（仅非高亮）
              │
              ├── config.Highlight ?
              │     ├── YES: outputHighlighted()
              │     │           ├── highlightContent()  // 渲染
              │     │           └── 打印（带行号）
              │     │
              │     └── NO: outputPlain()
              │                 └── 打印（带行号）
```

---

## 九、文件变更

### 9.1 修改文件

- **`internal/cli/cat.go`**:
  - 移除 `catShowNewline` 标志定义
  - 移除 `-N/--show-newline` 标志注册
  - 更新 `runCat` 函数，不再传递 `ShowNewline`

- **`internal/commands/cat/cmd_cat.go`**:
  - 从 `CatConfig` 中移除 `ShowNewline` 字段
  - 简化 `processFile` 函数
  - 删除 `processHead`、`processTail`、`processLine`、`readLine`

- **`internal/commands/cat/syntax_highlight.go`** → **`viewer.go`**:
  - 删除 `HighlightViewer` 结构体
  - 创建 `FileViewer` 结构体
  - 实现 `View`、`sliceLines`、`applySpecialFlags`、`outputPlain`、`outputHighlighted` 方法
  - 保留 `highlightContent`、`shouldHighlight`、`stripANSI` 等辅助函数

### 9.2 代码结构

```
internal/commands/cat/
├── cmd_cat.go          # 简化后的主逻辑
├── viewer.go           # 统一查看器
├── highlight.go        # 高亮相关函数（可选，可从 viewer.go 分离）
└── ov_pager.go         # 分页器模式（不变）

internal/cli/
└── cat.go              # 移除 -N 标志
```

---

## 十、优势总结

| 方面 | 旧实现 | 新实现 |
|------|--------|--------|
| 读取逻辑 | 2 套（文件指针流式 / 内存读取） | 1 套（统一内存读取） |
| 切片逻辑 | 2 套（流式计数 / 数组切片） | 1 套（统一数组切片） |
| 行号处理 | 2 套（分散在多处） | 1 套（统一在 output 方法） |
| 特殊标志 | 2 套（processLine / 无） | 1 套（统一 applySpecialFlags） |
| 代码复杂度 | 高 | 低 |
| 维护性 | 差 | 好 |

---

## 十一、边缘案例

| 场景 | 处理 |
|------|------|
| 空文件 | `bytes.Split` 返回 `[[]byte{}]`，输出空 |
| 文件末尾无换行符 | 最后一行无换行符，正常处理 |
| head > 总行数 | 返回所有行 |
| tail > 总行数 | 返回所有行 |
| 高亮失败 | 降级为普通输出 |
| 不支持的语言 | `highlightContent` 返回 `ok=false`，降级为普通输出 |

---

## 十二、实现步骤

1. **P0**: 修改 `internal/cli/cat.go`，移除 `-N` 标志
2. **P1**: 修改 `internal/commands/cat/cmd_cat.go`，从 `CatConfig` 移除 `ShowNewline`
3. **P2**: 创建 `viewer.go`，实现 `FileViewer` 结构体和 `View` 方法
4. **P3**: 实现 `sliceLines` 方法
5. **P4**: 实现 `applySpecialFlags` 方法
6. **P5**: 实现 `outputPlain` 和 `outputHighlighted` 方法
7. **P6**: 修改 `cmd_cat.go`，使用 `FileViewer` 替换原有逻辑
8. **P7**: 删除 `cmd_cat.go` 中的旧方法
9. **P8**: 删除或重命名 `syntax_highlight.go`
10. **P9**: 测试验证

---

**总结**: 通过移除 `-N` 标志，统一换行符处理更简单，代码更清晰。统一查看器将高亮和非高亮模式的处理逻辑集中管理，维护更容易。
