# cat 命令显示换行符设计方案

## 1. 功能需求

显示文件每行实际包含的换行符类型（`\n` 或 `\r\n`），帮助用户识别文件的换行符格式。

## 2. 显示效果

### Linux/Mac 文件 (LF `\n`)
```
hello world\n$
line two\n$
```

### Windows 文件 (CRLF `\r\n`)
```
hello world\r\n$
line two\r\n$
```

### 无换行符结尾的最后一行
```
last line without newline$
```

## 3. CLI 标志设计

| 长标志 | 短标志 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| `--show-newline` | `-N` | bool | false | 显示换行符类型（`\n` 或 `\r\n`） |

**注意：**
- 使用大写 `-N` 避免与 `-n`（显示行号）冲突
- 可与 `-E`（显示 `$`）同时使用

## 4. 配置结构体修改

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
    ShowNewline  bool     // -N 显示换行符类型
    HeadLines    int      // --head 显示前N行 (0表示全部)
    TailLines    int      // --tail 显示后N行 (0表示全部)
    Quiet        bool     // -q 静默模式 (不显示错误信息)

    // 运行时
    LineCounter int // 行号计数器
}
```

## 5. 核心实现方案（方案B）

### 5.1 新增辅助函数

```go
// readLine 读取一行，返回内容、换行符标记和错误
//
// 参数:
//   - reader: bufio.Reader
//
// 返回:
//   - content: 行内容（不含换行符）
//   - newline: 换行符标记（"\n"、"\r\n" 或空字符串）
//   - err: 错误信息
func readLine(reader *bufio.Reader) (content, newline string, err error) {
    line, err := reader.ReadString('\n')
    if err != nil && err != io.EOF {
        return "", "", err
    }

    // 检测换行符类型
    if strings.HasSuffix(line, "\r\n") {
        return line[:len(line)-2], "\\r\\n", err
    } else if strings.HasSuffix(line, "\n") {
        return line[:len(line)-1], "\\n", err
    }

    // 无换行符（最后一行）
    return line, "", err
}
```

### 5.2 修改 processLine 函数

```go
// processLine 处理单行
//
// 参数:
//   - line: 行内容
//   - newline: 换行符标记
//   - config: 命令配置
func processLine(line, newline string, config *CatConfig) {
    isBlank := len(strings.TrimSpace(line)) == 0

    // 处理行号
    if config.ShowLineNum || (config.ShowNonBlank && !isBlank) {
        config.LineCounter++
        fmt.Printf("%6d\t", config.LineCounter)
    }

    // 处理特殊字符显示
    if config.ShowTabs {
        line = strings.ReplaceAll(line, "\t", "^I")
    }

    // 输出内容
    fmt.Print(line)

    // 处理换行符显示
    if config.ShowNewline {
        fmt.Print(newline)
    }

    // 处理行尾标记
    if config.ShowEnd {
        fmt.Print("$")
    }
    fmt.Println()
}
```

### 5.3 修改 processFile 函数

```go
// processFile 处理单个文件
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误（如果有）
func processFile(path string, config *CatConfig) error {
    // 打开文件
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open file %s: %w", path, err)
    }
    defer file.Close()

    // 获取文件信息（用于判断是否是目录）
    info, err := file.Stat()
    if err != nil {
        return fmt.Errorf("failed to get file info %s: %w", path, err)
    }
    if info.IsDir() {
        return fmt.Errorf("%s is a directory", path)
    }

    // 根据 head/tail 选项处理
    if config.HeadLines > 0 {
        return processHead(file, config)
    }

    if config.TailLines > 0 {
        return processTail(file, config)
    }

    // 普通处理：使用 bufio.Reader 逐行读取
    reader := bufio.NewReader(file)
    for {
        line, newline, err := readLine(reader)
        if err != nil && err != io.EOF {
            return err
        }

        processLine(line, newline, config)

        if err == io.EOF {
            break
        }
    }

    return nil
}
```

### 5.4 修改 processHead 函数

```go
// processHead 处理文件前N行
//
// 参数:
//   - file: 已打开的文件
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误（如果有）
func processHead(file *os.File, config *CatConfig) error {
    reader := bufio.NewReader(file)
    lineCount := 0

    for {
        if lineCount >= config.HeadLines {
            break
        }

        line, newline, err := readLine(reader)
        if err != nil && err != io.EOF {
            return err
        }

        processLine(line, newline, config)
        lineCount++

        if err == io.EOF {
            break
        }
    }

    return nil
}
```

### 5.5 修改 processTail 函数

```go
// lineInfo 存储行内容和换行符信息
type lineInfo struct {
    content string
    newline string
}

// processTail 处理文件后N行
//
// 参数:
//   - file: 已打开的文件
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误（如果有）
func processTail(file *os.File, config *CatConfig) error {
    // 使用环形缓冲区存储最后N行
    ring := make([]lineInfo, config.TailLines)
    index := 0
    count := 0

    reader := bufio.NewReader(file)
    for {
        line, newline, err := readLine(reader)
        if err != nil && err != io.EOF {
            return err
        }

        ring[index] = lineInfo{content: line, newline: newline}
        index = (index + 1) % config.TailLines
        if count < config.TailLines {
            count++
        }

        if err == io.EOF {
            break
        }
    }

    // 按顺序输出
    start := (index - count + config.TailLines) % config.TailLines
    for i := 0; i < count; i++ {
        info := ring[(start+i)%config.TailLines]
        processLine(info.content, info.newline, config)
    }

    return nil
}
```

## 6. CLI 定义修改

```go
var (
    catShowLineNum  *qflag.BoolFlag // -n 显示所有行号
    catShowNonBlank *qflag.BoolFlag // -b 显示非空行行号
    catShowEnd      *qflag.BoolFlag // -E 显示行尾$
    catShowTabs     *qflag.BoolFlag // -T 显示制表符为^I
    catShowAll      *qflag.BoolFlag // -A 等价于 -ET
    catShowNewline  *qflag.BoolFlag // -N 显示换行符类型
    catHeadLines    *qflag.IntFlag  // --head 显示前N行
    catTailLines    *qflag.IntFlag  // --tail 显示后N行
    catQuiet        *qflag.BoolFlag // -q 静默模式
)

func init() {
    CatCmd = qflag.NewCmd("cat", "", qflag.ExitOnError)

    catShowLineNum = CatCmd.Bool("number", "n", "显示所有行号", false)
    catShowNonBlank = CatCmd.Bool("number-nonblank", "b", "仅显示非空行行号", false)
    catShowEnd = CatCmd.Bool("show-ends", "E", "在每行末尾显示$", false)
    catShowTabs = CatCmd.Bool("show-tabs", "T", "将制表符显示为^I", false)
    catShowAll = CatCmd.Bool("show-all", "A", "等价于 -ET", false)
    catShowNewline = CatCmd.Bool("show-newline", "N", "显示换行符类型（\\n 或 \\r\\n）", false)
    catHeadLines = CatCmd.Int("head", "", "显示前N行（0表示全部）", 0)
    catTailLines = CatCmd.Int("tail", "t", "显示后N行（0表示全部）", 0)
    catQuiet = CatCmd.Bool("quiet", "q", "静默模式（不显示错误）", false)

    // ... 其余代码不变
}

func runCat(cmd qflag.Command) error {
    config := cat.CatConfig{
        Targets:      cmd.Args(),
        ShowLineNum:  catShowLineNum.Get(),
        ShowNonBlank: catShowNonBlank.Get(),
        ShowEnd:      catShowEnd.Get(),
        ShowTabs:     catShowTabs.Get(),
        ShowAll:      catShowAll.Get(),
        ShowNewline:  catShowNewline.Get(),  // 新增
        HeadLines:    catHeadLines.Get(),
        TailLines:    catTailLines.Get(),
        Quiet:        catQuiet.Get(),
    }

    return cat.CatCmdMain(config)
}
```

## 7. 使用示例

```bash
# 显示换行符类型
fck cat -N file.txt

# 同时显示换行符和行尾标记
fck cat -NE file.txt

# 显示行号和换行符类型
fck cat -nN file.txt

# 显示前10行并显示换行符
fck cat --head 10 -N file.txt

# 显示后5行并显示换行符
fck cat -t 5 -N file.txt
```

## 8. 边缘情况处理

| 场景 | 处理方式 |
|------|----------|
| Linux 文件 (LF) | 显示 `\n` |
| Windows 文件 (CRLF) | 显示 `\r\n` |
| 旧 Mac 文件 (CR) | 显示 `\r`（如果存在） |
| 最后一行无换行符 | 不显示换行符标记 |
| 空文件 | 无输出 |
| 空行 | 只显示换行符标记 |

## 9. 实现步骤

1. 修改 `internal/commands/cat/cmd_cat.go`
   - 添加 `ShowNewline` 字段到 `CatConfig`
   - 添加 `readLine()` 辅助函数
   - 修改 `processLine()` 添加 `newline` 参数
   - 修改 `processFile()` 使用 `bufio.Reader`
   - 修改 `processHead()` 使用 `bufio.Reader`
   - 修改 `processTail()` 使用 `bufio.Reader` 和 `lineInfo`

2. 修改 `internal/cli/cat.go`
   - 添加 `catShowNewline` 变量
   - 添加 `--show-newline` / `-N` 标志
   - 在 `runCat()` 中传递 `ShowNewline`

3. 编译测试
4. 功能验证

## 10. 代码规范检查清单

- [ ] 函数级注释完整
- [ ] 错误信息为英文
- [ ] 帮助信息为中文
- [ ] 跨平台兼容
- [ ] 编译通过
- [ ] 功能测试通过
