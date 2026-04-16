# Cat 命令语法高亮重构设计方案

> **设计日期**: 2026-04-16  
> **版本**: v1.0  
> **状态**: 待评审

---

## 一、问题分析

### 1.1 当前实现的问题

当前代码存在三个几乎相同逻辑的高亮处理函数：

1. `processFileWithHighlight` - 整文件高亮
2. `processHeadWithHighlight` - 前N行高亮
3. `processTailWithHighlight` - 后N行高亮

**重复代码**:
- 都读取整个文件内容
- 都调用 `highlightContent` 进行高亮
- 都按行分割处理
- 都处理行号逻辑

### 1.2 核心问题

语法高亮需要完整上下文，所以必须：
1. 读取整个文件
2. 一次性高亮
3. 然后选择需要的行输出

当前实现中，head/tail 模式重复了高亮逻辑，只是输出行范围不同。

---

## 二、重构方案

### 2.1 核心思想

**"先高亮，后切片"** - 统一的高亮流程：

```
读取文件 → 语法高亮 → 按行分割 → 选择范围 → 输出
```

### 2.2 统一的高亮查看器

```go
// HighlightViewer 语法高亮查看器
type HighlightViewer struct {
    config     *CatConfig
    hlConfig   HighlightConfig
}

// View 查看文件（支持范围选择）
//
// 参数:
//   - path: 文件路径
//   - startLine: 起始行号（0-based，-1 表示从头开始）
//   - endLine: 结束行号（-1 表示到末尾）
//
// 返回:
//   - error: 处理错误
//
// 使用方式:
//   - 查看全部: viewer.View(path, -1, -1)
//   - 查看前N行: viewer.View(path, 0, N)
//   - 查看后N行: viewer.View(path, -N, -1) 或 viewer.View(path, total-N, -1)
func (v *HighlightViewer) View(path string, startLine, endLine int) error
```

### 2.3 简化后的调用流程

```go
// processFile 中的高亮分支
if config.Highlight {
    viewer := NewHighlightViewer(config)
    
    if config.HeadLines > 0 {
        // 查看前N行
        return viewer.View(path, 0, config.HeadLines)
    }
    
    if config.TailLines > 0 {
        // 查看后N行（viewer 内部计算起始行）
        return viewer.View(path, -config.TailLines, -1)
    }
    
    // 查看全部
    return viewer.View(path, -1, -1)
}
```

### 2.4 内部实现

```go
// HighlightViewer 实现

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

// View 查看文件（支持范围选择）
func (v *HighlightViewer) View(path string, startLine, endLine int) error {
    // 1. 读取文件
    content, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }

    // 2. 语法高亮（一次性）
    hlContent, ok, hlErr := highlightContent(content, path, v.hlConfig)
    if hlErr != nil {
        if !v.config.Quiet {
            fmt.Fprintf(os.Stderr, "warn: syntax highlight failed for %s: %v\n", path, hlErr)
        }
        // 降级为普通处理
        return v.viewPlain(content, startLine, endLine)
    }
    if !ok {
        // 不支持高亮，降级为普通处理
        return v.viewPlain(content, startLine, endLine)
    }

    // 3. 按行分割
    lines := bytes.Split(hlContent, []byte("\n"))
    totalLines := len(lines)

    // 4. 计算实际范围
    actualStart, actualEnd := v.calculateRange(totalLines, startLine, endLine)

    // 5. 重置行号计数器
    v.config.LineCounter = actualStart

    // 6. 输出指定范围
    for i := actualStart; i < actualEnd && i < totalLines; i++ {
        line := string(lines[i])
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

// calculateRange 计算实际行范围
//
// 参数:
//   - total: 总行数
//   - start: 请求的起始行（-1 表示从头，负数表示从末尾倒数）
//   - end: 请求的结束行（-1 表示到末尾）
//
// 返回:
//   - actualStart: 实际起始行（0-based）
//   - actualEnd: 实际结束行（不包含）
//
// 示例:
//   - total=100, start=-1, end=-1  → (0, 100)  全部
//   - total=100, start=0, end=10   → (0, 10)   前10行
//   - total=100, start=-10, end=-1 → (90, 100) 后10行
func (v *HighlightViewer) calculateRange(total, start, end int) (actualStart, actualEnd int) {
    // 处理 start
    if start < 0 {
        // 从末尾计算
        actualStart = total + start
        if actualStart < 0 {
            actualStart = 0
        }
    } else {
        actualStart = start
    }

    // 处理 end
    if end < 0 || end > total {
        actualEnd = total
    } else {
        actualEnd = end
    }

    // 边界检查
    if actualStart > total {
        actualStart = total
    }
    if actualEnd < actualStart {
        actualEnd = actualStart
    }

    return actualStart, actualEnd
}

// viewPlain 无高亮模式查看（降级处理）
func (v *HighlightViewer) viewPlain(content []byte, startLine, endLine int) error {
    lines := bytes.Split(content, []byte("\n"))
    totalLines := len(lines)

    actualStart, actualEnd := v.calculateRange(totalLines, startLine, endLine)
    v.config.LineCounter = actualStart

    for i := actualStart; i < actualEnd && i < totalLines; i++ {
        line := string(lines[i])
        isBlank := len(strings.TrimSpace(line)) == 0

        if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
            v.config.LineCounter++
            fmt.Printf("%6d\t", v.config.LineCounter)
        }

        fmt.Println(line)
    }

    return nil
}
```

---

## 三、文件结构变更

### 3.1 删除文件

- `highlight_processor.go`（整个文件删除）

### 3.2 修改文件

**`cmd_cat.go`**:

```go
// 修改 processFile 中的高亮分支
if config.Highlight {
    viewer := NewHighlightViewer(config)
    
    switch {
    case config.HeadLines > 0:
        return viewer.View(path, 0, config.HeadLines)
    case config.TailLines > 0:
        // 先获取总行数，或让 viewer 内部处理
        return viewer.ViewTail(path, config.TailLines)
    default:
        return viewer.View(path, -1, -1)
    }
}
```

**`syntax_highlight.go`**:

添加 `HighlightViewer` 结构体和相关方法。

---

## 四、API 设计对比

### 4.1 旧 API（重复）

```go
processFileWithHighlight(file, path, config)      // 整文件
processHeadWithHighlight(file, path, config)      // 前N行
processTailWithHighlight(file, path, config)      // 后N行
```

**问题**:
- 3 个函数，逻辑重复 80%
- 每个函数都独立读取文件、高亮、分割
- 维护困难

### 4.2 新 API（统一）

```go
viewer := NewHighlightViewer(config)

viewer.View(path, -1, -1)      // 整文件
viewer.View(path, 0, N)        // 前N行
viewer.View(path, -N, -1)      // 后N行（或提供 ViewTail 辅助方法）
```

**优点**:
- 1 个结构体，1 个核心方法
- 文件读取和高亮只做一次
- 通过参数控制范围，逻辑清晰
- 易于扩展（如添加 ViewRange 查看中间行）

---

## 五、辅助方法设计

为简化调用，可以提供语义化的包装方法：

```go
// ViewAll 查看全部内容
func (v *HighlightViewer) ViewAll(path string) error {
    return v.View(path, -1, -1)
}

// ViewHead 查看前N行
func (v *HighlightViewer) ViewHead(path string, n int) error {
    return v.View(path, 0, n)
}

// ViewTail 查看后N行
func (v *HighlightViewer) ViewTail(path string, n int) error {
    // 需要知道总行数，或传入负数让 calculateRange 处理
    return v.View(path, -n, -1)
}
```

---

## 六、调用示例

### 6.1 重构后的 processFile

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

### 6.2 与分页器模式对比

```go
// 日常查看模式（新）
if config.Highlight {
    viewer := NewHighlightViewer(config)
    return viewer.ViewAll(path)  // 高亮后输出到 stdout
}

// 分页器模式（已有）
if config.Highlight {
    // 高亮后交给 ov 处理
    hlContent, _, _ := highlightContent(content, path, hlConfig)
    doc.ControlReader(bytes.NewReader(hlContent), nil)
}
```

---

## 七、实现优先级

1. **P0**: 在 `syntax_highlight.go` 中添加 `HighlightViewer` 结构体
2. **P1**: 实现 `View` 核心方法和 `calculateRange` 辅助方法
3. **P2**: 添加 `ViewAll`、`ViewHead`、`ViewTail` 语义化方法
4. **P3**: 修改 `cmd_cat.go`，替换原有的三个高亮处理函数
5. **P4**: 删除 `highlight_processor.go`
6. **P5**: 测试验证

---

## 八、边缘情况处理

| 场景 | 处理方式 |
|------|----------|
| 空文件 | 直接返回，无输出 |
| 单行文件 | 正常高亮输出 |
| head N > 总行数 | 输出全部内容 |
| tail N > 总行数 | 输出全部内容 |
| 高亮失败 | 降级为普通文本输出 |
| 不支持的语言 | 降级为普通文本输出 |
| 文件读取失败 | 返回错误 |

---

**总结**: 通过 "先高亮，后切片" 的统一流程，将三个重复函数合并为一个灵活的查看器，代码量减少约 60%，维护性大幅提升。
