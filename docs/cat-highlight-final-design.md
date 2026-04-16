# Cat 命令语法高亮重构最终设计方案

> **设计日期**: 2026-04-16  
> **版本**: Final  
> **状态**: 待实现

---

## 一、核心思想

**"读取全部 → 统一换行符 → 切片 → 高亮选中部分 → 输出"**

---

## 二、设计方案

### 2.1 统一的高亮查看器

```go
// HighlightViewer 语法高亮查看器
type HighlightViewer struct {
    config   *CatConfig
    hlConfig HighlightConfig
}

// NewHighlightViewer 创建高亮查看器
func NewHighlightViewer(config *CatConfig) *HighlightViewer {
    return &HighlightViewer{
        config:   config,
        hlConfig: DefaultHighlightConfig(),
    }
}

// ViewAll 查看全部内容
func (v *HighlightViewer) ViewAll(path string) error

// ViewHead 查看前N行
func (v *HighlightViewer) ViewHead(path string, n int) error

// ViewTail 查看后N行
func (v *HighlightViewer) ViewTail(path string, n int) error
```

### 2.2 核心实现

```go
// View 查看文件指定范围（先切片，后高亮）
//
// 参数:
//   - path: 文件路径
//   - startLine: 起始行号（0-based，包含）
//   - endLine: 结束行号（不包含，-1 表示到末尾）
//
// 返回:
//   - error: 处理错误
func (v *HighlightViewer) View(path string, startLine, endLine int) error {
    // 1. 读取文件
    content, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }

    // 2. 统一换行符为 \n（处理 Windows \r\n）
    content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

    // 3. 按行分割
    allLines := bytes.Split(content, []byte("\n"))
    totalLines := len(allLines)

    // 4. 计算实际范围
    actualStart, actualEnd := v.calculateRange(totalLines, startLine, endLine)

    // 5. 提取需要显示的行
    selectedLines := allLines[actualStart:actualEnd]

    // 6. 合并为待高亮内容
    contentToHighlight := bytes.Join(selectedLines, []byte("\n"))

    // 7. 语法高亮（仅对选中部分）
    hlContent, ok, hlErr := highlightContent(contentToHighlight, path, v.hlConfig)
    if hlErr != nil || !ok {
        // 高亮失败或不支持，降级为普通输出
        return v.outputPlain(selectedLines, actualStart)
    }

    // 8. 输出高亮后的内容
    return v.outputHighlighted(hlContent, actualStart)
}

// calculateRange 计算实际行范围
//
// 参数:
//   - total: 总行数
//   - start: 请求的起始行（0-based）
//   - end: 请求的结束行（-1 表示到末尾）
//
// 返回:
//   - actualStart: 实际起始行
//   - actualEnd: 实际结束行（不包含）
func (v *HighlightViewer) calculateRange(total, start, end int) (actualStart, actualEnd int) {
    actualStart = start
    if actualStart < 0 {
        actualStart = 0
    }
    if actualStart > total {
        actualStart = total
    }

    actualEnd = end
    if actualEnd < 0 || actualEnd > total {
        actualEnd = total
    }

    if actualEnd < actualStart {
        actualEnd = actualStart
    }

    return actualStart, actualEnd
}

// outputPlain 普通输出（无高亮）
func (v *HighlightViewer) outputPlain(lines [][]byte, startLineNum int) error {
    v.config.LineCounter = startLineNum

    for _, lineBytes := range lines {
        line := string(lineBytes)
        isBlank := len(strings.TrimSpace(line)) == 0

        // 处理行号
        if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
            v.config.LineCounter++
            fmt.Printf("%6d\t", v.config.LineCounter)
        }

        fmt.Println(line)
    }

    return nil
}

// outputHighlighted 高亮输出
func (v *HighlightViewer) outputHighlighted(hlContent []byte, startLineNum int) error {
    // 按行分割高亮后的内容
    lines := bytes.Split(hlContent, []byte("\n"))
    v.config.LineCounter = startLineNum

    for _, lineBytes := range lines {
        line := string(lineBytes)
        isBlank := len(strings.TrimSpace(stripANSI(line))) == 0

        // 处理行号
        if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
            v.config.LineCounter++
            fmt.Printf("%6d\t", v.config.LineCounter)
        }

        fmt.Println(line)
    }

    return nil
}
```

### 2.3 语义化包装方法

```go
// ViewAll 查看全部内容
func (v *HighlightViewer) ViewAll(path string) error {
    return v.View(path, 0, -1)
}

// ViewHead 查看前N行
func (v *HighlightViewer) ViewHead(path string, n int) error {
    if n <= 0 {
        return nil
    }
    return v.View(path, 0, n)
}

// ViewTail 查看后N行
func (v *HighlightViewer) ViewTail(path string, n int) error {
    if n <= 0 {
        return nil
    }

    // 读取文件获取总行数
    content, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }

    // 统一换行符
    content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
    allLines := bytes.Split(content, []byte("\n"))
    totalLines := len(allLines)

    startLine := totalLines - n
    if startLine < 0 {
        startLine = 0
    }

    return v.View(path, startLine, -1)
}
```

---

## 三、调用流程

### 3.1 重构后的 processFile

```go
func processFile(path string, config *CatConfig) error {
    // ... 前置检查（二进制、目录等） ...

    // 高亮模式
    if config.Highlight {
        viewer := NewHighlightViewer(config)

        switch {
        case config.HeadLines > 0:
            return viewer.ViewHead(path, config.HeadLines)
        case config.TailLines > 0:
            return viewer.ViewTail(path, config.TailLines)
        default:
            return viewer.ViewAll(path)
        }
    }

    // 普通模式（保持原有逻辑）
    if config.HeadLines > 0 {
        return processHead(file, config)
    }
    if config.TailLines > 0 {
        return processTail(file, config)
    }

    // 普通逐行处理 ...
}
```

### 3.2 流程对比

| 模式 | 流程 |
|------|------|
| **ViewAll** | 读取 → 统一换行符 → 分割 → 全部行 → 高亮全部 → 输出 |
| **ViewHead** | 读取 → 统一换行符 → 分割 → 取前N行 → 高亮N行 → 输出 |
| **ViewTail** | 读取 → 统一换行符 → 分割 → 取后N行 → 高亮N行 → 输出 |

---

## 四、文件变更

### 4.1 删除文件
- `highlight_processor.go`

### 4.2 修改文件
- `syntax_highlight.go`: 添加 `HighlightViewer` 结构体和方法
- `cmd_cat.go`: 使用 `HighlightViewer` 替换原有三个高亮函数调用

---

## 五、关键代码位置

### 5.1 syntax_highlight.go 新增内容

```go
package cat

import (
    "bytes"
    "fmt"
    "os"
    "strings"
)

// HighlightViewer 语法高亮查看器
type HighlightViewer struct {
    config   *CatConfig
    hlConfig HighlightConfig
}

// NewHighlightViewer 创建高亮查看器
func NewHighlightViewer(config *CatConfig) *HighlightViewer {
    return &HighlightViewer{
        config:   config,
        hlConfig: DefaultHighlightConfig(),
    }
}

// View 查看文件指定范围
func (v *HighlightViewer) View(path string, startLine, endLine int) error {
    // ... 实现见 2.2 节 ...
}

// ViewAll 查看全部内容
func (v *HighlightViewer) ViewAll(path string) error {
    return v.View(path, 0, -1)
}

// ViewHead 查看前N行
func (v *HighlightViewer) ViewHead(path string, n int) error {
    if n <= 0 {
        return nil
    }
    return v.View(path, 0, n)
}

// ViewTail 查看后N行
func (v *HighlightViewer) ViewTail(path string, n int) error {
    if n <= 0 {
        return nil
    }

    content, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }

    content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
    allLines := bytes.Split(content, []byte("\n"))
    totalLines := len(allLines)

    startLine := totalLines - n
    if startLine < 0 {
        startLine = 0
    }

    return v.View(path, startLine, -1)
}

// calculateRange 计算实际行范围
func (v *HighlightViewer) calculateRange(total, start, end int) (actualStart, actualEnd int) {
    actualStart = start
    if actualStart < 0 {
        actualStart = 0
    }
    if actualStart > total {
        actualStart = total
    }

    actualEnd = end
    if actualEnd < 0 || actualEnd > total {
        actualEnd = total
    }

    if actualEnd < actualStart {
        actualEnd = actualStart
    }

    return actualStart, actualEnd
}

// outputPlain 普通输出
func (v *HighlightViewer) outputPlain(lines [][]byte, startLineNum int) error {
    v.config.LineCounter = startLineNum

    for _, lineBytes := range lines {
        line := string(lineBytes)
        isBlank := len(strings.TrimSpace(line)) == 0

        if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
            v.config.LineCounter++
            fmt.Printf("%6d\t", v.config.LineCounter)
        }

        fmt.Println(line)
    }

    return nil
}

// outputHighlighted 高亮输出
func (v *HighlightViewer) outputHighlighted(hlContent []byte, startLineNum int) error {
    lines := bytes.Split(hlContent, []byte("\n"))
    v.config.LineCounter = startLineNum

    for _, lineBytes := range lines {
        line := string(lineBytes)
        isBlank := len(strings.TrimSpace(stripANSI(line))) == 0

        if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
            v.config.LineCounter++
            fmt.Printf("%6d\t", v.config.LineCounter)
        }

        fmt.Println(line)
    }

    return nil
}
```

### 5.2 cmd_cat.go 修改内容

```go
func processFile(path string, config *CatConfig) error {
    // ... 前置检查 ...

    // 高亮模式
    if config.Highlight {
        viewer := NewHighlightViewer(config)

        switch {
        case config.HeadLines > 0:
            return viewer.ViewHead(path, config.HeadLines)
        case config.TailLines > 0:
            return viewer.ViewTail(path, config.TailLines)
        default:
            return viewer.ViewAll(path)
        }
    }

    // ... 普通模式 ...
}
```

---

## 六、边缘情况处理

| 场景 | 处理方式 |
|------|----------|
| 空文件 | `allLines` 为 `[[]]`，正常处理 |
| head N > 总行数 | `calculateRange` 限制为 total |
| tail N > 总行数 | `startLine` 计算为 0，输出全部 |
| 高亮失败 | 降级为 `outputPlain` |
| 不支持的语言 | `highlightContent` 返回 `ok=false`，降级 |
| Windows 换行符 | `ReplaceAll` 统一为 `\n` |
| Mac 旧格式 `\r` | 暂不处理（现代 Mac 用 `\n`） |

---

## 七、实现步骤

1. **P0**: 在 `syntax_highlight.go` 中添加 `HighlightViewer` 结构体
2. **P1**: 实现 `View`、`calculateRange`、`outputPlain`、`outputHighlighted`
3. **P2**: 实现 `ViewAll`、`ViewHead`、`ViewTail`
4. **P3**: 修改 `cmd_cat.go`，替换高亮处理逻辑
5. **P4**: 删除 `highlight_processor.go`
6. **P5**: 测试验证

---

## 八、优势总结

| 方面 | 旧实现 | 新实现 |
|------|--------|--------|
| 代码重复 | 3 个函数 80% 重复 | 1 个核心方法 + 3 个包装方法 |
| 维护性 | 差 | 好 |
| 换行符处理 | 未统一 | 统一处理 `\r\n` |
| 高亮范围 | 全部高亮 | 只高亮需要显示的部分 |
| 扩展性 | 差 | 好（可添加 ViewRange 等） |

---

**总结**: 通过统一的高亮查看器，将重复代码合并，同时处理了 Windows 换行符问题，只高亮需要显示的内容，代码更简洁易维护。
