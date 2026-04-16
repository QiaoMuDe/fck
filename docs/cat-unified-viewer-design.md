# Cat 命令统一查看器设计方案

> **设计日期**: 2026-04-16  
> **版本**: v1.0  
> **状态**: 待评审

---

## 一、核心思想

**"一个查看器，内部根据高亮标志决定处理方式"**

将高亮模式和非高亮模式的处理逻辑统一到 `FileViewer` 中，对外提供统一的 `ViewAll`、`ViewHead`、`ViewTail` 接口，内部根据 `config.Highlight` 决定如何处理。

---

## 二、设计方案

### 2.1 统一查看器结构

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

### 2.2 对外接口

```go
// ViewAll 查看全部内容
func (v *FileViewer) ViewAll(path string) error

// ViewHead 查看前N行
func (v *FileViewer) ViewHead(path string, n int) error

// ViewTail 查看后N行
func (v *FileViewer) ViewTail(path string, n int) error
```

### 2.3 内部实现

每个对外接口内部根据 `config.Highlight` 分发到对应的私有方法：

```go
// ViewHead 查看前N行
func (v *FileViewer) ViewHead(path string, n int) error {
	if v.config.Highlight {
		return v.viewHeadHighlighted(path, n)
	}
	return v.viewHeadPlain(path, n)
}
```

---

## 三、详细实现

### 3.1 高亮模式实现

```go
// viewHeadHighlighted 高亮模式查看前N行
func (v *FileViewer) viewHeadHighlighted(path string, n int) error {
	// 1. 读取文件
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// 2. 统一换行符
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

	// 3. 按行分割并取前N行
	allLines := bytes.Split(content, []byte("\n"))
	if n > len(allLines) {
		n = len(allLines)
	}
	selectedLines := allLines[:n]

	// 4. 合并待高亮内容
	contentToHighlight := bytes.Join(selectedLines, []byte("\n"))

	// 5. 语法高亮
	hlContent, ok, hlErr := highlightContent(contentToHighlight, path, v.hlConfig)
	if hlErr != nil || !ok {
		// 高亮失败，降级为普通输出
		return v.outputPlain(selectedLines, 0)
	}

	// 6. 输出高亮内容
	return v.outputHighlighted(hlContent, 0)
}

// viewTailHighlighted 高亮模式查看后N行
func (v *FileViewer) viewTailHighlighted(path string, n int) error {
	// 1. 读取文件
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// 2. 统一换行符
	content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))

	// 3. 按行分割并取后N行
	allLines := bytes.Split(content, []byte("\n"))
	totalLines := len(allLines)
	startLine := totalLines - n
	if startLine < 0 {
		startLine = 0
	}
	selectedLines := allLines[startLine:]

	// 4. 合并待高亮内容
	contentToHighlight := bytes.Join(selectedLines, []byte("\n"))

	// 5. 语法高亮
	hlContent, ok, hlErr := highlightContent(contentToHighlight, path, v.hlConfig)
	if hlErr != nil || !ok {
		// 高亮失败，降级为普通输出
		return v.outputPlain(selectedLines, startLine)
	}

	// 6. 输出高亮内容
	return v.outputHighlighted(hlContent, startLine)
}

// viewAllHighlighted 高亮模式查看全部
func (v *FileViewer) viewAllHighlighted(path string) error {
	// 1. 读取文件
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// 2. 语法高亮
	hlContent, ok, hlErr := highlightContent(content, path, v.hlConfig)
	if hlErr != nil || !ok {
		// 高亮失败，降级为普通输出
		content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
		lines := bytes.Split(content, []byte("\n"))
		return v.outputPlain(lines, 0)
	}

	// 3. 输出高亮内容
	return v.outputHighlighted(hlContent, 0)
}
```

### 3.2 普通模式实现

```go
// viewHeadPlain 普通模式查看前N行
func (v *FileViewer) viewHeadPlain(path string, n int) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineCount := 0

	for lineCount < n {
		line, newline, err := v.readLine(reader)
		if err != nil && err != io.EOF {
			return err
		}

		v.outputLine(line, newline)
		lineCount++

		if err == io.EOF {
			break
		}
	}

	return nil
}

// viewTailPlain 普通模式查看后N行
func (v *FileViewer) viewTailPlain(path string, n int) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	// 使用环形缓冲区存储最后N行
	ring := make([]lineInfo, n)
	index := 0
	count := 0

	reader := bufio.NewReader(file)
	for {
		line, newline, err := v.readLine(reader)
		if err != nil && err != io.EOF {
			return err
		}

		ring[index] = lineInfo{content: line, newline: newline}
		index = (index + 1) % n
		if count < n {
			count++
		}

		if err == io.EOF {
			break
		}
	}

	// 按顺序输出
	start := (index - count + n) % n
	for i := 0; i < count; i++ {
		info := ring[(start+i)%n]
		v.outputLine(info.content, info.newline)
	}

	return nil
}

// viewAllPlain 普通模式查看全部
func (v *FileViewer) viewAllPlain(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, newline, err := v.readLine(reader)
		if err != nil && err != io.EOF {
			return err
		}

		v.outputLine(line, newline)

		if err == io.EOF {
			break
		}
	}

	return nil
}
```

### 3.3 辅助方法

```go
// readLine 读取一行
func (v *FileViewer) readLine(reader *bufio.Reader) (content, newline string, err error) {
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", "", err
	}

	switch {
	case strings.HasSuffix(line, "\r\n"):
		return line[:len(line)-2], "\r\n", err
	case strings.HasSuffix(line, "\n"):
		return line[:len(line)-1], "\n", err
	default:
		return line, "", err
	}
}

// outputLine 输出行（普通模式）
func (v *FileViewer) outputLine(line, newline string) {
	isBlank := len(strings.TrimSpace(line)) == 0

	// 处理行号
	if v.config.ShowLineNum || (v.config.ShowNonBlank && !isBlank) {
		v.config.LineCounter++
		fmt.Printf("%6d\t", v.config.LineCounter)
	}

	// 处理特殊字符显示
	if v.config.ShowTabs {
		line = strings.ReplaceAll(line, "\t", "^I")
	}

	// 输出内容
	fmt.Print(line)

	// 处理换行符显示
	if v.config.ShowNewline {
		fmt.Print(newline)
	}

	// 处理行尾标记
	if v.config.ShowEnd {
		fmt.Print("$")
	}
	fmt.Println()
}

// outputPlain 批量输出普通行（高亮降级时使用）
func (v *FileViewer) outputPlain(lines [][]byte, startLineNum int) error {
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

// outputHighlighted 批量输出高亮行
func (v *FileViewer) outputHighlighted(hlContent []byte, startLineNum int) error {
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

---

## 四、调用流程

### 4.1 重构后的 processFile

```go
func processFile(path string, config *CatConfig) error {
	// 1. 基础检查（目录、二进制）
	if err := checkFile(path, config); err != nil {
		return err
	}

	// 2. 创建统一查看器
	viewer := NewFileViewer(config)

	// 3. 根据模式查看
	switch {
	case config.HeadLines > 0:
		return viewer.ViewHead(path, config.HeadLines)
	case config.TailLines > 0:
		return viewer.ViewTail(path, config.TailLines)
	default:
		return viewer.ViewAll(path)
	}
}
```

### 4.2 流程对比

| 模式 | 旧流程 | 新流程 |
|------|--------|--------|
| 高亮 head | `viewer.ViewHead` (内部读取文件) | `viewer.ViewHead` → `viewHeadHighlighted` |
| 普通 head | `processHead` (使用已打开文件) | `viewer.ViewHead` → `viewHeadPlain` |
| 高亮 tail | `viewer.ViewTail` (内部读取文件) | `viewer.ViewTail` → `viewTailHighlighted` |
| 普通 tail | `processTail` (使用已打开文件) | `viewer.ViewTail` → `viewTailPlain` |

---

## 五、文件变更

### 5.1 修改文件
- `syntax_highlight.go`: 
  - 重命名为 `viewer.go` 或保持原名
  - 删除 `HighlightViewer`，改为 `FileViewer`
  - 整合普通模式和高亮模式的所有方法

- `cmd_cat.go`:
  - 简化 `processFile` 函数
  - 删除 `processHead`、`processTail`、`processLine`、`readLine`

### 5.2 代码结构

```
internal/commands/cat/
├── cmd_cat.go          # 简化后的主逻辑
├── viewer.go           # 统一查看器（原 syntax_highlight.go）
├── ov_pager.go         # 分页器模式（不变）
└── highlight.go        # 高亮相关函数（从 viewer.go 分离，可选）
```

---

## 六、优势总结

| 方面 | 旧实现 | 新实现 |
|------|--------|--------|
| 查看器数量 | 2 个（HighlightViewer + 普通函数） | 1 个（FileViewer） |
| 代码分布 | 分散在多个文件 | 集中在 viewer.go |
| 维护性 | 差（修改需改多处） | 好（修改只改一处） |
| 扩展性 | 难 | 易（添加模式只需添加方法） |
| 逻辑一致性 | 差（两种处理方式） | 好（统一入口，内部分发） |

---

## 七、实现步骤

1. **P0**: 在 `syntax_highlight.go` 中创建 `FileViewer` 结构体
2. **P1**: 实现高亮模式的三个私有方法
3. **P2**: 实现普通模式的三个私有方法
4. **P3**: 实现辅助方法（readLine、outputLine、outputPlain、outputHighlighted）
5. **P4**: 修改 `cmd_cat.go`，使用 `FileViewer` 替换原有逻辑
6. **P5**: 删除 `cmd_cat.go` 中的旧方法
7. **P6**: 测试验证

---

**总结**: 通过统一查看器，将高亮和非高亮模式的处理逻辑集中管理，代码更清晰、更易维护。
