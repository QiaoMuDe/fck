# Cat 命令语法高亮重构设计方案 v2

> **设计日期**: 2026-04-16  
> **版本**: v2.0  
> **状态**: 待评审

---

## 一、核心思想变更

### v1 方案（问题）
```
读取文件 → 语法高亮（整个文件）→ 按行分割 → 选择范围 → 输出
```
**问题**: 大文件时，高亮整个文件性能差

### v2 方案（优化）
```
读取文件 → 按行分割 → 选择范围（前几行/后几行/全部）→ 语法高亮（仅选中部分）→ 输出
```
**优势**: 只高亮需要显示的部分，性能更好

---

## 二、新设计方案

### 2.1 统一的高亮查看器

```go
// HighlightViewer 语法高亮查看器
type HighlightViewer struct {
    config     *CatConfig
    hlConfig   HighlightConfig
}

// NewHighlightViewer 创建高亮查看器
func NewHighlightViewer(config *CatConfig) *HighlightViewer {
    return &HighlightViewer{
        config:   config,
        hlConfig: DefaultHighlightConfig(),
    }
}

// View 查看文件（支持范围选择）
//
// 参数:
//   - path: 文件路径
//   - startLine: 起始行号（0-based）
//   - endLine: 结束行号（-1 表示到末尾，不包含）
//
// 返回:
//   - error: 处理错误
//
// 使用方式:
//   - 查看全部: viewer.View(path, 0, -1)
//   - 查看前N行: viewer.View(path, 0, N)
//   - 查看后N行: viewer.View(path, total-N, -1)
func (v *HighlightViewer) View(path string, startLine, endLine int) error

// ViewAll 查看全部内容
func (v *HighlightViewer) ViewAll(path string) error

// ViewHead 查看前N行
func (v *HighlightViewer) ViewHead(path string, n int) error

// ViewTail 查看后N行
func (v *HighlightViewer) ViewTail(path string, n int) error
```

### 2.2 核心实现

```go
// View 查看文件（先切片，后高亮）
func (v *HighlightViewer) View(path string, startLine, endLine int) error {
    // 1. 读取文件
    content, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }

    // 2. 按行分割
    allLines := bytes.Split(content, []byte("\n"))
    totalLines := len(allLines)

    // 3. 计算实际范围
    actualStart, actualEnd := v.calculateRange(totalLines, startLine, endLine)

    // 4. 提取需要显示的行
    selectedLines := allLines[actualStart:actualEnd]

    // 5. 合并为待高亮内容
    contentToHighlight := bytes.Join(selectedLines, []byte("\n"))

    // 6. 语法高亮（仅对选中部分）
    hlContent, ok, hlErr := highlightContent(contentToHighlight, path, v.hlConfig)
    if hlErr != nil || !ok {
        // 高亮失败或不支持，降级为普通输出
        return v.outputPlain(selectedLines, actualStart)
    }

    // 7. 输出高亮后的内容
    return v.outputHighlighted(hlContent, actualStart)
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
    // ...
}
```

### 3.2 流程对比

| 模式 | 流程 |
|------|------|
| **查看全部** | 读取 → 分割 → 全部行 → 高亮全部 → 输出 |
| **查看前N行** | 读取 → 分割 → 取前N行 → 高亮N行 → 输出 |
| **查看后N行** | 读取 → 分割 → 取后N行 → 高亮N行 → 输出 |

---

## 四、文件变更

### 4.1 删除文件
- `highlight_processor.go`

### 4.2 修改文件
- `syntax_highlight.go`: 添加 `HighlightViewer` 结构体和方法
- `cmd_cat.go`: 使用 `HighlightViewer` 替换原有三个高亮函数

---

## 五、优势对比

| 方面 | v1（先高亮后切片） | v2（先切片后高亮） |
|------|-------------------|-------------------|
| 大文件性能 | 差（高亮整个文件） | 好（只高亮部分） |
| 代码复杂度 | 低 | 低 |
| 内存占用 | 高 | 低 |
| 维护性 | 好 | 好 |

---

## 六、边缘情况

| 场景 | 处理 |
|------|------|
| 空文件 | 直接返回 |
| head N > 总行数 | 输出全部 |
| tail N > 总行数 | 输出全部 |
| 高亮失败 | 降级为普通输出 |
| 不支持的语言 | 降级为普通输出 |

---

**总结**: v2 方案通过 "先切片，后高亮" 的流程，避免了高亮无用内容，特别适合大文件的 head/tail 查看场景。
