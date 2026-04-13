# Grep 多文件处理设计方案

## 当前现状

当前 `grep` 命令支持以下模式：

```bash
# 单文件搜索
fck grep "pattern" file.txt

# 标准输入
fck grep "pattern" < file.txt

# 递归搜索目录
fck grep "pattern" ./dir -r
fck grep "pattern" ./dir -r --include "*.go"
```

**不支持**：直接指定多个文件或通配符
```bash
# 不支持
fck grep "pattern" file1.txt file2.txt file3.txt
fck grep "pattern" *.go
```

## 目标

支持多文件处理，与 Linux grep 行为一致：

```bash
# 多个指定文件
fck grep "pattern" file1.txt file2.txt file3.txt

# 通配符（shell 展开）
fck grep "pattern" *.go

# 混合使用
fck grep "pattern" file1.txt *.go ./dir/*.js
```

## 设计方案

### 1. CLI 层修改 (internal/cli/grep.go)

**核心思路**：
- 接受多个文件参数（`args[1:]`）
- 展开通配符
- 遍历处理每个文件
- 多文件时自动显示文件名（类似 `-H` 行为）

```go
func runGrep(cmd qflag.Command) error {
    args := cmd.Args()

    // 检查参数
    if len(args) < 1 {
        return fmt.Errorf("no pattern specified")
    }

    pattern := args[0]
    
    // 收集文件参数（从 args[1] 开始）
    var fileArgs []string
    if len(args) >= 2 {
        fileArgs = args[1:]
    }

    // 展开通配符
    var targets []string
    for _, arg := range fileArgs {
        matches, err := filepath.Glob(arg)
        if err != nil {
            fmt.Printf("error globbing %s: %v\n", arg, err)
            continue
        }
        if len(matches) == 0 {
            targets = append(targets, arg)
        } else {
            targets = append(targets, matches...)
        }
    }

    // 准备基础配置
    baseConfig := grep.GrepConfig{
        Pattern:       pattern,
        IgnoreCase:    grepIgnoreCase.Get(),
        InvertMatch:   grepInvertMatch.Get(),
        LineNumber:    grepLineNumber.Get(),
        Count:         grepCount.Get(),
        WithFilename:  grepWithFilename.Get(),
        Regexp:        grepRegexp.Get(),
        MaxCount:      grepMaxCount.Get(),
        BeforeContext: grepBeforeContext.Get(),
        AfterContext:  grepAfterContext.Get(),
        Context:       grepContext.Get(),
        NoColor:       grepNoColor.Get(),
        Recursive:     grepRecursive.Get(),
        FollowSymlink: grepFollowSymlink.Get(),
        Include:       grepInclude.Get(),
        Exclude:       grepExclude.Get(),
        IncludeDir:    grepIncludeDir.Get(),
        ExcludeDir:    grepExcludeDir.Get(),
        NoMessages:    grepNoMessages.Get(),
        Hidden:        grepHidden.Get(),
        Text:          grepText.Get(),
        IgnoreBinary:  grepIgnoreBinary.Get(),
    }

    // 无文件参数：从 stdin 读取
    if len(targets) == 0 {
        baseConfig.Target = ""
        return grep.GrepCmdMain(baseConfig)
    }

    // 单文件：保持原有行为
    if len(targets) == 1 {
        baseConfig.Target = targets[0]
        return grep.GrepCmdMain(baseConfig)
    }

    // 多目标：循环处理（支持文件和目录共存）
    var hasError bool
    for _, target := range targets {
        config := baseConfig
        config.Target = target
        // 多文件时强制显示文件名
        config.WithFilename = true

        if err := grep.GrepCmdMain(config); err != nil {
            fmt.Printf("error processing %s: %v\n", target, err)
            hasError = true
        }
    }

    if hasError {
        return fmt.Errorf("some files failed to process")
    }
    return nil
}
```

### 2. Core 层修改 (internal/commands/grep/cmd_grep.go)

**核心修改**：
- `GrepCmdMain` 接收单个目标（文件或目录）
- 根据类型和 `-r` 标志决定处理方式
- 每个目标独立计数（`matchCount`）

**修改点 1：支持文件和目录混合处理**

```go
func GrepCmdMain(config GrepConfig) error {
    // ... 原有验证和编译逻辑 ...

    // 重置匹配计数（每个目标独立计数）
    config.matchCount = 0
    config.lastPrintedFile = ""

    // 未指定目标：从 stdin 读取
    if config.Target == "" {
        config.filename = "(standard input)"
        return processFile(os.Stdin, &config)
    }

    // 检查目标类型
    info, err := os.Stat(config.Target)
    if err != nil {
        if config.NoMessages {
            return nil
        }
        if os.IsNotExist(err) {
            return fmt.Errorf("file not found: %s", config.Target)
        }
        return fmt.Errorf("cannot access file: %w", err)
    }

    // 根据类型处理
    if info.IsDir() {
        // 目录：需要 -r 标志
        if !config.Recursive && !config.FollowSymlink {
            if config.NoMessages {
                return nil
            }
            return fmt.Errorf("is a directory, use -r or -R for recursive search: %s", config.Target)
        }
        return processPathRecursive(config.Target, &config)
    }

    // 普通文件：直接处理
    return processSingleFileWithError(config.Target, &config)
}
```

**修改点 2：多文件输出格式**

当前已有两种输出格式：
- `printLineSingle`：单文件模式
- `printLineRecursive`：递归模式

**新增多文件模式**：类似递归模式，但不需要空行分隔

```go
// printLineMulti 多文件模式输出行
// 格式: 文件名:行号:内容（多文件时）
func printLineMulti(highlightedLine string, lineNum int, separator string, config *GrepConfig) {
    if config.Count {
        return
    }

    var result string
    
    // 文件名（蓝色）
    result += formatPathWithColor(config.filename, config.NoColor)
    
    // 行号（青色）
    if config.LineNumber {
        result += ":"
        if !config.NoColor {
            result += ColorCyan + strconv.Itoa(lineNum) + ColorReset
        } else {
            result += strconv.Itoa(lineNum)
        }
    }
    
    // 内容
    result += separator + highlightedLine
    fmt.Println(result)
}
```

**修改点 3：printLine 函数增加多文件判断**

```go
func printLine(line string, lineNum int, separator string, config *GrepConfig) {
    if config.Count {
        return
    }

    highlightedLine := line
    if separator == ":" {
        highlightedLine = highlightLine(line, config)
    }

    switch {
    case config.Recursive || config.FollowSymlink:
        printLineRecursive(highlightedLine, lineNum, separator, config)
    case config.multiFile:  // 新增：多文件模式
        printLineMulti(highlightedLine, lineNum, separator, config)
    default:
        printLineSingle(highlightedLine, lineNum, separator, config)
    }
}
```

**修改点 4：GrepConfig 增加多文件标记**

```go
type GrepConfig struct {
    // ... 原有字段 ...
    
    multiFile bool // 是否为多文件模式（自动设置）
}
```

### 3. 输出格式对比

| 模式 | 格式示例 |
|------|----------|
| 单文件 | `line content` 或 `10:line content` |
| 多文件 | `file.txt:10:line content` |
| 递归 | `file.txt`（单独一行）<br>`10:line content` |

### 4. 边界情况处理

#### 4.1 文件不存在
- 打印错误到 stderr
- 继续处理其他文件
- 最后返回错误码

#### 4.2 目录参数（无 -r）
- 打印错误："is a directory, use -r or -R for recursive search"
- 继续处理其他文件

#### 4.3 通配符无匹配
- 保留原参数
- 让 core 层处理文件不存在错误

#### 4.4 与 -r 标志共存
- 多文件和 `-r` 可以同时使用
- 自动识别文件类型：
  - 普通文件：直接处理
  - 目录：如果设置了 `-r` 则递归处理，否则报错
- 示例：`grep pattern file1.txt ./dir -r`
  - file1.txt：直接处理
  - ./dir：递归处理

### 5. 实现步骤

1. **修改 CLI 层** (`internal/cli/grep.go`)
   - 收集多个文件参数
   - 展开通配符
   - 循环调用 `GrepCmdMain`

2. **修改 Core 层** (`internal/commands/grep/cmd_grep.go`)
   - 添加 `multiFile` 标记
   - 重置 `matchCount` 和 `lastPrintedFile`
   - 添加 `printLineMulti` 函数
   - 修改 `printLine` 支持多文件模式

3. **测试验证**
   - 单文件（向后兼容）
   - 多文件
   - 通配符
   - 错误处理

### 6. 代码变更范围

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/cli/grep.go` | 修改 | 支持多文件参数和通配符 |
| `internal/commands/grep/cmd_grep.go` | 修改 | 多文件输出格式，重置计数器 |

### 7. 示例用法

```bash
# 基础多文件
fck grep "pattern" a.txt b.txt c.txt

# 通配符
fck grep "pattern" *.go

# 混合
fck grep "pattern" file1.txt *.go ./dir/*.js

# 带选项
fck grep "pattern" -n -i *.go
fck grep "pattern" -c *.txt

# 与 -r 共存（文件直接处理，目录递归处理）
fck grep "pattern" file1.txt ./dir -r
```

### 8. 与现有功能对比

| 功能 | 当前实现 | 多文件支持后 |
|------|----------|--------------|
| 单文件 | `grep pattern file` | 保持不变 |
| 标准输入 | `grep pattern` | 保持不变 |
| 递归 | `grep pattern dir -r` | 保持不变 |
| 多文件 | 不支持 | `grep pattern file1 file2` |
| 通配符 | 不支持 | `grep pattern *.go` |

## 注意事项

1. **向后兼容**：单文件用法完全保持兼容
2. **自动显示文件名**：多文件时自动启用 `-H` 行为
3. **与 -r 共存**：多文件和递归模式可以同时使用，自动识别文件类型
4. **独立计数**：每个文件的匹配计数独立
5. **错误聚合**：多文件时收集所有错误，最后统一返回
