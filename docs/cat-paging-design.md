# Cat 命令分页查看功能设计方案

> **设计目标**: 为 cat 命令添加类似 Linux `less` 的分页查看功能
> **设计日期**: 2026-04-15
> **参考实现**: less, more, bat

---

## 一、需求分析

### 1.1 功能目标
- 支持分页查看大文件内容，避免一次性输出到终端
- 提供类似 `less` 的交互式导航体验
- 保持与现有 cat 命令的兼容性
- 支持管道输入的分页查看

### 1.2 使用场景
```bash
# 分页查看大文件
fck cat -p large.log

# 管道输入分页
fck cat file.txt | fck cat -p

# 结合其他选项使用
fck cat -p -n large.txt    # 分页显示并带行号
fck cat -p -u 1000 log.txt # 从第1000行开始分页
```

---

## 二、CLI 接口设计

### 2.1 新增标志

| 标志 | 长格式 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| `-p` | `--paging` | bool | false | 启用分页模式 |
| `-S` | `--chop-long-lines` | bool | false | 截断长行（不换行） |
| `-W` | `--wrap` | bool | true | 自动换行（与 -S 互斥） |
|      | `--no-init` | bool | false | 不清屏（用于管道） |

### 2.2 交互式快捷键

| 按键 | 功能 |
|------|------|
| `Space` / `f` / `PageDown` | 向下翻一页 |
| `b` / `PageUp` | 向上翻一页 |
| `j` / `↓` | 向下滚动一行 |
| `k` / `↑` | 向上滚动一行 |
| `g` / `Home` | 跳到文件开头 |
| `G` / `End` | 跳到文件结尾 |
| `/` | 搜索（向下） |
| `?` | 搜索（向上） |
| `n` | 下一个匹配 |
| `N` | 上一个匹配 |
| `q` / `Q` | 退出分页器 |
| `h` / `?` | 显示帮助 |

---

## 三、核心架构设计

### 3.1 模块结构

```
internal/commands/cat/
├── cmd_cat.go          # 主入口（已有）
├── pager.go            # 分页器核心逻辑（新增）
├── screen.go           # 终端控制（新增）
└── search.go           # 搜索功能（新增）
```

### 3.2 核心类型定义

```go
// PagerConfig 分页器配置
type PagerConfig struct {
    ChopLongLines bool   // -S: 截断长行
    Wrap          bool   // -W: 自动换行
    NoInit        bool   // 不清屏
    StartLine     int    // 起始行号
}

// Pager 分页器状态
type Pager struct {
    config      PagerConfig
    lines       []string      // 文件内容（或窗口缓冲区）
    totalLines  int           // 总行数
    currentLine int           // 当前顶部行号
    screenLines int           // 屏幕可显示行数
    screenCols  int           // 屏幕列数
    filePath    string        // 文件路径
    searchTerm  string        // 当前搜索词
    matches     []int         // 匹配行索引
    currentMatch int          // 当前匹配位置
}

// Screen 终端控制接口
type Screen struct {
    origMode    *term.State  // 原始终端状态
    width       int
    height      int
}
```

### 3.3 核心流程

```
┌─────────────────────────────────────────────────────────────┐
│  1. 检测分页标志 (-p)                                        │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  2. 初始化分页器                                             │
│     - 获取终端尺寸                                           │
│     - 设置原始模式（禁用行缓冲、不回显）                      │
│     - 清屏（除非 --no-init）                                  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  3. 加载文件内容                                             │
│     - 小文件（< 10MB）：全部加载到内存                        │
│     - 大文件（>= 10MB）：使用窗口缓冲区，按需读取              │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  4. 主循环                                                   │
│     ┌─────────────┐                                         │
│     │ 渲染当前页   │                                         │
│     └──────┬──────┘                                         │
│            │                                                │
│            ▼                                                │
│     ┌─────────────┐                                         │
│     │ 等待用户输入 │                                         │
│     └──────┬──────┘                                         │
│            │                                                │
│            ▼                                                │
│     ┌─────────────┐    退出命令?    ┌──────────┐           │
│     │ 处理命令    │ ──────────────> │ 恢复终端 │           │
│     └─────────────┘                 │ 退出程序 │           │
│            │                        └──────────┘           │
│            │ 翻页/滚动命令                                    │
│            ▼                                                │
│     ┌─────────────┐                                         │
│     │ 更新位置    │ ───────────────────────────────────────>│
│     └─────────────┘                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 四、详细实现方案

### 4.1 终端控制 (screen.go)

```go
package cat

import (
    "os"
    "golang.org/x/term"
)

// Screen 终端屏幕控制
type Screen struct {
    origState *term.State
    width     int
    height    int
}

// InitScreen 初始化终端屏幕
// 
// 功能:
//   - 保存原始终端状态
//   - 切换到原始模式（禁用行缓冲、不回显）
//   - 获取终端尺寸
//
// 返回:
//   - *Screen: 屏幕控制对象
//   - error: 初始化错误
func InitScreen() (*Screen, error) {
    // 检查是否为终端
    if !term.IsTerminal(int(os.Stdin.Fd())) {
        return nil, fmt.Errorf("not a terminal")
    }
    
    // 保存原始状态
    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        return nil, err
    }
    
    // 获取终端尺寸
    width, height, err := term.GetSize(int(os.Stdout.Fd()))
    if err != nil {
        term.Restore(int(os.Stdin.Fd()), oldState)
        return nil, err
    }
    
    return &Screen{
        origState: oldState,
        width:     width,
        height:    height,
    }, nil
}

// Restore 恢复终端原始状态
func (s *Screen) Restore() error {
    if s.origState != nil {
        return term.Restore(int(os.Stdin.Fd()), s.origState)
    }
    return nil
}

// Clear 清屏
func (s *Screen) Clear() {
    fmt.Print("\x1b[2J\x1b[H")  // ANSI 清屏并移动光标到左上角
}

// MoveCursor 移动光标到指定位置
func (s *Screen) MoveCursor(row, col int) {
    fmt.Printf("\x1b[%d;%dH", row, col)
}

// HideCursor 隐藏光标
func (s *Screen) HideCursor() {
    fmt.Print("\x1b[?25l")
}

// ShowCursor 显示光标
func (s *Screen) ShowCursor() {
    fmt.Print("\x1b[?25h")
}
```

### 4.2 分页器核心 (pager.go)

```go
package cat

import (
    "bufio"
    "fmt"
    "os"
)

// Pager 分页查看器
type Pager struct {
    screen      *Screen
    config      CatConfig
    pagerConfig PagerConfig
    
    // 内容存储
    lines       []string    // 所有行（小文件）或窗口缓冲区（大文件）
    totalLines  int         // 总行数
    
    // 显示状态
    topLine     int         // 当前显示区域的顶部行号（从1开始）
    screenRows  int         // 可用行数（减去状态行）
    screenCols  int         // 可用列数
    
    // 搜索状态
    searchTerm   string
    matches      []int      // 匹配行号
    currentMatch int        // 当前匹配索引
}

// NewPager 创建分页器
func NewPager(screen *Screen, config CatConfig, pconfig PagerConfig) *Pager {
    rows := screen.height - 1  // 保留一行用于状态栏
    if rows < 1 {
        rows = 1
    }
    
    return &Pager{
        screen:      screen,
        config:      config,
        pagerConfig: pconfig,
        screenRows:  rows,
        screenCols:  screen.width,
        topLine:     1,
    }
}

// LoadFile 加载文件内容
// 
// 策略:
//   - 文件 < 10MB: 全部加载到内存
//   - 文件 >= 10MB: 建立行索引，按需读取
func (p *Pager) LoadFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // 获取文件大小
    info, err := file.Stat()
    if err != nil {
        return err
    }
    
    // 小文件：全部加载
    if info.Size() < 10*1024*1024 {
        return p.loadSmallFile(file)
    }
    
    // 大文件：建立索引
    return p.loadLargeFile(file, info.Size())
}

// loadSmallFile 加载小文件到内存
func (p *Pager) loadSmallFile(file *os.File) error {
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        p.lines = append(p.lines, scanner.Text())
    }
    p.totalLines = len(p.lines)
    return scanner.Err()
}

// loadLargeFile 为大文件建立行索引
func (p *Pager) loadLargeFile(file *os.File, size int64) error {
    // TODO: 实现大文件支持
    // 1. 建立行偏移索引
    // 2. 使用窗口缓冲区
    return p.loadSmallFile(file) // 暂时使用小文件方式
}

// Run 运行分页器主循环
func (p *Pager) Run() error {
    if !p.pagerConfig.NoInit {
        p.screen.Clear()
    }
    
    for {
        // 渲染当前页
        if err := p.render(); err != nil {
            return err
        }
        
        // 读取用户输入
        cmd, err := p.readCommand()
        if err != nil {
            return err
        }
        
        // 处理命令
        switch cmd {
        case 'q', 'Q':
            return nil
        case ' ':  // 空格 - 下一页
            p.pageDown()
        case 'b', 'B':  // 上一页
            p.pageUp()
        case 'j', 0x1b5b42:  // j 或 下箭头
            p.scrollDown()
        case 'k', 0x1b5b41:  // k 或 上箭头
            p.scrollUp()
        case 'g':  // 跳到开头
            p.topLine = 1
        case 'G':  // 跳到结尾
            p.jumpToEnd()
        case '/':  // 搜索
            p.handleSearch()
        case 'n':  // 下一个匹配
            p.nextMatch()
        case 'N':  // 上一个匹配
            p.prevMatch()
        case 'h', '?':  // 帮助
            p.showHelp()
        }
    }
}

// render 渲染当前页面
func (p *Pager) render() error {
    p.screen.Clear()
    
    // 计算显示范围
    endLine := p.topLine + p.screenRows - 1
    if endLine > p.totalLines {
        endLine = p.totalLines
    }
    
    // 渲染每一行
    for i := p.topLine - 1; i < endLine && i < len(p.lines); i++ {
        line := p.lines[i]
        
        // 处理行号显示
        if p.config.ShowLineNum {
            fmt.Printf("%6d  ", i+1)
        }
        
        // 处理长行
        if p.pagerConfig.ChopLongLines && len(line) > p.screenCols-7 {
            line = line[:p.screenCols-7]
        }
        
        fmt.Println(line)
    }
    
    // 渲染状态栏
    p.renderStatusBar(endLine)
    
    return nil
}

// renderStatusBar 渲染底部状态栏
func (p *Pager) renderStatusBar(endLine int) {
    p.screen.MoveCursor(p.screen.height, 1)
    
    percent := 0
    if p.totalLines > 0 {
        percent = endLine * 100 / p.totalLines
    }
    
    status := fmt.Sprintf("%s lines %d-%d/%d (%d%%) ", 
        p.config.Targets[0], p.topLine, endLine, p.totalLines, percent)
    
    if p.searchTerm != "" {
        status += fmt.Sprintf("[search: %s] ", p.searchTerm)
    }
    
    // 反色显示状态栏
    fmt.Printf("\x1b[7m%-*s\x1b[0m", p.screenCols, status)
}

// readCommand 读取单个按键命令
func (p *Pager) readCommand() (byte, error) {
    var buf [3]byte
    n, err := os.Stdin.Read(buf[:])
    if err != nil {
        return 0, err
    }
    
    // 处理转义序列（方向键等）
    if n == 3 && buf[0] == 0x1b && buf[1] == 0x5b {
        // 返回组合值用于识别方向键
        return buf[0]<<16 | buf[1]<<8 | buf[2], nil
    }
    
    return buf[0], nil
}

// 导航方法
func (p *Pager) pageDown() {
    p.topLine += p.screenRows
    if p.topLine > p.totalLines {
        p.topLine = p.totalLines
    }
}

func (p *Pager) pageUp() {
    p.topLine -= p.screenRows
    if p.topLine < 1 {
        p.topLine = 1
    }
}

func (p *Pager) scrollDown() {
    if p.topLine < p.totalLines {
        p.topLine++
    }
}

func (p *Pager) scrollUp() {
    if p.topLine > 1 {
        p.topLine--
    }
}

func (p *Pager) jumpToEnd() {
    p.topLine = p.totalLines - p.screenRows + 1
    if p.topLine < 1 {
        p.topLine = 1
    }
}
```

### 4.3 搜索功能 (search.go)

```go
package cat

import (
    "fmt"
    "regexp"
    "strings"
)

// handleSearch 处理搜索输入
func (p *Pager) handleSearch() {
    // 在状态栏显示搜索提示
    p.screen.MoveCursor(p.screen.height, 1)
    fmt.Print("/")
    
    // 读取搜索词（简化实现）
    // 实际应该实现行编辑功能
    term := p.readLine()
    if term == "" {
        return
    }
    
    p.searchTerm = term
    p.performSearch()
}

// performSearch 执行搜索
func (p *Pager) performSearch() {
    p.matches = nil
    p.currentMatch = -1
    
    // 简单字符串匹配（可扩展为正则）
    for i, line := range p.lines {
        if strings.Contains(line, p.searchTerm) {
            p.matches = append(p.matches, i)
            // 跳转到第一个匹配
            if p.currentMatch == -1 {
                p.currentMatch = 0
                p.topLine = i + 1
                if p.topLine > p.totalLines-p.screenRows+1 {
                    p.topLine = p.totalLines - p.screenRows + 1
                    if p.topLine < 1 {
                        p.topLine = 1
                    }
                }
            }
        }
    }
}

// nextMatch 跳转到下一个匹配
func (p *Pager) nextMatch() {
    if len(p.matches) == 0 {
        return
    }
    p.currentMatch++
    if p.currentMatch >= len(p.matches) {
        p.currentMatch = 0
    }
    p.topLine = p.matches[p.currentMatch] + 1
}

// prevMatch 跳转到上一个匹配
func (p *Pager) prevMatch() {
    if len(p.matches) == 0 {
        return
    }
    p.currentMatch--
    if p.currentMatch < 0 {
        p.currentMatch = len(p.matches) - 1
    }
    p.topLine = p.matches[p.currentMatch] + 1
}

// readLine 从终端读取一行（简化版）
func (p *Pager) readLine() string {
    // 简化实现：读取直到回车
    var result strings.Builder
    buf := make([]byte, 1)
    
    for {
        n, _ := os.Stdin.Read(buf)
        if n == 0 {
            continue
        }
        
        if buf[0] == '\r' || buf[0] == '\n' {
            break
        }
        
        if buf[0] == 0x7f {  // Backspace
            str := result.String()
            if len(str) > 0 {
                result.Reset()
                result.WriteString(str[:len(str)-1])
                fmt.Print("\b \b")
            }
            continue
        }
        
        result.WriteByte(buf[0])
        fmt.Print(string(buf[0]))
    }
    
    return result.String()
}
```

### 4.4 帮助界面

```go
// showHelp 显示帮助信息
func (p *Pager) showHelp() {
    p.screen.Clear()
    
    help := `
┌─────────────────────────────────────────────────────────────┐
│                        CAT PAGER HELP                       │
├─────────────────────────────────────────────────────────────┤
│  Navigation                                                 │
│    Space, f, PageDown    Next page                          │
│    b, PageUp             Previous page                      │
│    j, Down arrow         Scroll down one line               │
│    k, Up arrow           Scroll up one line                 │
│    g, Home               Go to beginning of file            │
│    G, End                Go to end of file                  │
│                                                             │
│  Search                                                     │
│    /pattern              Search forward                     │
│    ?pattern              Search backward (if implemented)   │
│    n                     Next match                         │
│    N                     Previous match                     │
│                                                             │
│  Other                                                      │
│    h, ?                  Show this help                     │
│    q                     Quit pager                         │
└─────────────────────────────────────────────────────────────┘
`
    fmt.Print(help)
    
    // 等待任意键继续
    buf := make([]byte, 1)
    os.Stdin.Read(buf)
}
```

---

## 五、与现有代码集成

### 5.1 修改 cmd_cat.go

```go
// CatConfig 添加分页配置
type CatConfig struct {
    // ... 现有字段 ...
    
    // 分页模式
    Paging        bool   // -p, --paging 启用分页
    ChopLongLines bool   // -S, --chop-long-lines 截断长行
    NoInit        bool   // --no-init 不清屏
}

// CatCmdMain 修改主函数
func CatCmdMain(config CatConfig) error {
    // ... 现有验证代码 ...
    
    // 分页模式
    if config.Paging {
        return runPagerMode(config)
    }
    
    // ... 现有普通模式代码 ...
}

// runPagerMode 运行分页模式
func runPagerMode(config CatConfig) error {
    // 初始化终端
    screen, err := InitScreen()
    if err != nil {
        return fmt.Errorf("cannot initialize pager: %w", err)
    }
    defer screen.Restore()
    
    // 创建分页器
    pconfig := PagerConfig{
        ChopLongLines: config.ChopLongLines,
        Wrap:          !config.ChopLongLines,
        NoInit:        config.NoInit,
    }
    
    pager := NewPager(screen, config, pconfig)
    
    // 加载文件
    for _, target := range config.Targets {
        if err := pager.LoadFile(target); err != nil {
            return err
        }
    }
    
    // 运行分页器
    return pager.Run()
}
```

### 5.2 修改 cli/cat.go

```go
// 添加新标志
var (
    // ... 现有标志 ...
    catPaging       *qflag.BoolFlag // -p, --paging 分页模式
    catChopLongLines *qflag.BoolFlag // -S, --chop-long-lines 截断长行
    catNoInit       *qflag.BoolFlag // --no-init 不清屏
)

func init() {
    // ... 现有初始化 ...
    
    // 分页标志
    catPaging = CatCmd.Bool("paging", "p", "启用分页查看模式", false)
    catChopLongLines = CatCmd.Bool("chop-long-lines", "S", "截断长行不换行", false)
    catNoInit = CatCmd.Bool("no-init", "", "不清屏（用于管道）", false)
    
    // 更新 Examples
    cmdOpts.Examples["分页查看大文件"] = fmt.Sprintf("%s cat -p large.log", qflag.Root.Name())
    cmdOpts.Examples["分页查看并截断长行"] = fmt.Sprintf("%s cat -p -S wide.txt", qflag.Root.Name())
}

func runCat(cmd qflag.Command) error {
    config := cat.CatConfig{
        // ... 现有配置 ...
        Paging:        catPaging.Get(),
        ChopLongLines: catChopLongLines.Get(),
        NoInit:        catNoInit.Get(),
    }
    
    return cat.CatCmdMain(config)
}
```

---

## 六、大文件优化方案

### 6.1 行索引策略

对于超过 10MB 的大文件，采用以下策略：

```go
// LargeFileIndex 大文件索引
type LargeFileIndex struct {
    file       *os.File
    lineOffsets []int64   // 每行的文件偏移量
    totalLines int
}

// buildIndex 建立行索引
func (idx *LargeFileIndex) buildIndex() error {
    reader := bufio.NewReader(idx.file)
    offset := int64(0)
    
    for {
        line, err := reader.ReadBytes('\n')
        if len(line) > 0 {
            idx.lineOffsets = append(idx.lineOffsets, offset)
            idx.totalLines++
            offset += int64(len(line))
        }
        
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }
    
    return nil
}

// getLine 获取指定行内容
func (idx *LargeFileIndex) getLine(lineNum int) (string, error) {
    if lineNum < 1 || lineNum > idx.totalLines {
        return "", fmt.Errorf("line number out of range")
    }
    
    offset := idx.lineOffsets[lineNum-1]
    _, err := idx.file.Seek(offset, 0)
    if err != nil {
        return "", err
    }
    
    reader := bufio.NewReader(idx.file)
    line, err := reader.ReadString('\n')
    if err != nil && err != io.EOF {
        return "", err
    }
    
    return strings.TrimSuffix(line, "\n"), nil
}
```

### 6.2 窗口缓冲区

```go
// WindowBuffer 窗口缓冲区
type WindowBuffer struct {
    fileIndex   *LargeFileIndex
    buffer      []string  // 当前窗口的行内容
    startLine   int       // 缓冲区起始行号
    capacity    int       // 缓冲区容量（行数）
}

// loadWindow 加载窗口内容
func (wb *WindowBuffer) loadWindow(centerLine int) error {
    halfCapacity := wb.capacity / 2
    wb.startLine = centerLine - halfCapacity
    if wb.startLine < 1 {
        wb.startLine = 1
    }
    
    wb.buffer = wb.buffer[:0]
    
    for i := 0; i < wb.capacity; i++ {
        lineNum := wb.startLine + i
        if lineNum > wb.fileIndex.totalLines {
            break
        }
        
        line, err := wb.fileIndex.getLine(lineNum)
        if err != nil {
            return err
        }
        wb.buffer = append(wb.buffer, line)
    }
    
    return nil
}
```

---

## 七、测试计划

### 7.1 功能测试

| 测试项 | 预期结果 |
|--------|----------|
| 分页查看小文件 | 正常显示，可翻页 |
| 分页查看大文件 | 使用索引，内存占用低 |
| 搜索功能 | 高亮匹配，可跳转 |
| 长行处理 | -S 截断，默认换行 |
| 管道输入 | 支持 `cat file \| cat -p` |
| 终端恢复 | 退出后终端状态正常 |

### 7.2 边界测试

- 空文件
- 单行文件
- 无换行符的文件
- 包含控制字符的文件
- 终端尺寸变化

---

## 八、风险评估

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 终端兼容性 | 高 | 使用标准 ANSI 序列，提供回退 |
| 大文件性能 | 中 | 使用索引和窗口缓冲区 |
| 信号处理 | 中 | 捕获 SIGINT，确保恢复终端 |
| Windows 支持 | 低 | 使用 golang.org/x/term 跨平台 |

---

## 九、实现优先级

1. **P0 - 核心功能**
   - 基本分页显示
   - 翻页导航（Space, b）
   - 行滚动（j, k）
   - 退出（q）

2. **P1 - 增强功能**
   - 首尾跳转（g, G）
   - 搜索（/, n, N）
   - 长行处理（-S）

3. **P2 - 优化**
   - 大文件索引
   - 帮助界面
   - 管道输入支持

---

**设计完成时间**: 2026-04-15  
**设计者**: Claude Code
