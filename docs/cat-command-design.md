# cat 子命令设计方案

## 1. 功能需求分析

### 核心功能
显示文件内容（类似 Unix cat 命令）

### 扩展功能
- 显示行号 (`-n`)
- 显示非空行行号 (`-b`)
- 显示特殊字符 (`-A`, `-E`, `-T`)
- 显示文件头部N行 (`--head`)
- 显示文件尾部N行 (`--tail`)

---

## 2. 文件结构设计

```
internal/
├── commands/
│   └── cat/
│       ├── cmd_cat.go          # 主逻辑
│       └── utils.go            # 工具函数
└── cli/
    └── cat.go                  # CLI 定义
```

---

## 3. 配置结构体设计

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
    HeadLines    int      // --head 显示前N行（默认10，0表示全部）
    TailLines    int      // --tail 显示后N行（默认10，0表示全部）
    Quiet        bool     // -q 静默模式（不显示错误信息）
    
    // 运行时
    LineCounter  int      // 行号计数器
}
```

---

## 4. CLI 标志设计

| 长标志 | 短标志 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| `--number` | `-n` | bool | false | 显示所有行号 |
| `--number-nonblank` | `-b` | bool | false | 仅显示非空行行号 |
| `--show-ends` | `-E` | bool | false | 在每行末尾显示$ |
| `--show-tabs` | `-T` | bool | false | 将制表符显示为^I |
| `--show-all` | `-A` | bool | false | 等价于 -ET |
| `--head` | `-h` | int | 10 | 显示前N行（0表示全部） |
| `--tail` | `-t` | int | 10 | 显示后N行（0表示全部） |
| `--quiet` | `-q` | bool | false | 静默模式（不显示错误） |

**注意：**
- 同时指定 `-b` 和 `-n` 时，`-b` 优先
- `--head` 和 `--tail` 互斥，不能同时使用
- head/tail 默认值10遵循 Unix 惯例

---

## 5. 主函数流程

```go
// CatCmdMain 执行 cat 命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误（如果有）
func CatCmdMain(config CatConfig) error {
    // 1. 验证参数
    if len(config.Targets) == 0 {
        return fmt.Errorf("no files specified")
    }
    
    // 2. 处理标志冲突（-b 优先级高于 -n）
    if config.ShowNonBlank {
        config.ShowLineNum = false
    }
    
    // 3. 处理 --show-all
    if config.ShowAll {
        config.ShowEnd = true
        config.ShowTabs = true
    }
    
    // 4. 验证 head/tail 互斥
    if config.HeadLines > 0 && config.TailLines > 0 {
        return fmt.Errorf("cannot use --head and --tail together")
    }
    
    // 5. 处理每个文件
    for _, target := range config.Targets {
        if err := processFile(target, &config); err != nil {
            if !config.Quiet {
                return err
            }
        }
    }
    
    return nil
}
```

---

## 6. 核心处理逻辑

### 文件处理函数

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
    
    // 普通处理：逐行读取
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        processLine(line, config)
    }
    
    return scanner.Err()
}
```

### Head 处理函数

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
    scanner := bufio.NewScanner(file)
    lineCount := 0
    
    for scanner.Scan() {
        if lineCount >= config.HeadLines {
            break
        }
        line := scanner.Text()
        processLine(line, config)
        lineCount++
    }
    
    return scanner.Err()
}
```

### Tail 处理函数

```go
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
    ring := make([]string, config.TailLines)
    index := 0
    count := 0
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        ring[index] = scanner.Text()
        index = (index + 1) % config.TailLines
        if count < config.TailLines {
            count++
        }
    }
    
    if err := scanner.Err(); err != nil {
        return err
    }
    
    // 按顺序输出
    start := (index - count + config.TailLines) % config.TailLines
    for i := 0; i < count; i++ {
        processLine(ring[(start+i)%config.TailLines], config)
    }
    
    return nil
}
```

### 行处理函数

```go
// processLine 处理单行
//
// 参数:
//   - line: 行内容
//   - config: 命令配置
func processLine(line string, config *CatConfig) {
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
    
    // 处理行尾标记
    if config.ShowEnd {
        fmt.Print("$")
    }
    fmt.Println()
}
```

---

## 7. 使用示例

```bash
# 基础用法
fck cat file.txt

# 显示行号
fck cat -n file.txt

# 显示非空行行号
fck cat -b file.txt

# 显示特殊字符
fck cat -A file.txt

# 显示前10行（默认）
fck cat --head file.txt

# 显示前20行
fck cat --head 20 file.txt

# 显示后10行（默认）
fck cat --tail file.txt

# 显示后5行
fck cat -t 5 file.txt

# 多个文件
fck cat file1.txt file2.txt file3.txt

# 静默模式（不显示错误）
fck cat -q nonexistent.txt
```

---

## 8. 边缘情况处理

| 场景 | 处理方式 |
|------|----------|
| 文件不存在 | 返回错误（静默模式下不输出） |
| 是目录 | 返回错误 |
| 无权限 | 返回错误 |
| 二进制文件 | 正常输出（可能显示乱码） |
| 超大文件 | 流式处理，不占用大量内存 |
| 无行尾换行符 | 正常处理，不额外添加换行 |
| 多个文件 | 顺序处理，不添加分隔符 |
| 同时指定 Head/Tail | 返回错误（互斥） |
| Head/Tail 值为 0 | 显示全部行 |
| 空文件 | 不输出任何内容 |

---

## 9. 实现步骤

1. 创建 `internal/commands/cat/` 目录
2. 创建 `cmd_cat.go` 实现核心逻辑
3. 创建 `internal/cli/cat.go` 定义 CLI
4. 在 `root.go` 中注册命令
5. 编译测试
6. 功能验证

---

## 10. 代码规范检查清单

- [ ] 目录结构符合规范
- [ ] 命名符合规范（CatConfig, CatCmdMain, cat.go）
- [ ] 函数级注释完整
- [ ] 错误信息为英文
- [ ] 帮助信息为中文
- [ ] 跨平台兼容（Windows/Linux/macOS）
- [ ] 编译通过
- [ ] 功能测试通过
