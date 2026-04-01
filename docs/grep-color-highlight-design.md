# Grep 颜色高亮功能设计方案

## 1. 功能概述

为 grep 命令添加匹配关键字高亮显示功能，默认将匹配到的关键字渲染为红色，同时提供禁用颜色的选项。

## 2. 新增标志

| 长标志 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| --no-color | bool | false | 禁用颜色高亮 |

## 3. 实现方案

### 3.1 数据结构修改

**`internal/commands/grep/cmd_grep.go`**

```go
// GrepConfig 添加字段
type GrepConfig struct {
    // ... 现有字段 ...
    NoColor bool // --no-color 禁用颜色
    
    // 运行时
    compiledPattern *regexp.Regexp
    matchCount      int
    filename        string
}
```

### 3.2 颜色常量定义

```go
const (
    ColorRed    = "\033[31m" // 红色
    ColorReset  = "\033[0m"  // 重置
)
```

### 3.3 高亮函数实现

```go
// highlightLine 对匹配的关键字进行高亮
//
// 参数:
//   - line: 原始行内容
//   - config: 命令配置
//
// 返回:
//   - string: 高亮后的行内容
func highlightLine(line string, config *GrepConfig) string {
    // 禁用颜色或计数模式，直接返回
    if config.NoColor || config.Count {
        return line
    }

    if config.Regexp {
        return highlightRegex(line, config)
    }
    return highlightString(line, config)
}

// highlightRegex 正则模式高亮
func highlightRegex(line string, config *GrepConfig) string {
    // 使用 FindAllStringIndex 找到所有匹配位置
    matches := config.compiledPattern.FindAllStringIndex(line, -1)
    if matches == nil {
        return line
    }

    // 从后向前插入颜色码，避免位置偏移
    result := line
    for i := len(matches) - 1; i >= 0; i-- {
        start, end := matches[i][0], matches[i][1]
        matched := result[start:end]
        result = result[:start] + ColorRed + matched + ColorReset + result[end:]
    }
    return result
}

// highlightString 固定字符串模式高亮
func highlightString(line string, config *GrepConfig) string {
    pattern := config.Pattern
    if pattern == "" {
        return line
    }

    // 根据忽略大小写选项处理
    searchLine := line
    searchPattern := pattern
    if config.IgnoreCase {
        searchLine = strings.ToLower(line)
        searchPattern = strings.ToLower(pattern)
    }

    // 找到所有匹配位置（从后向前）
    var positions [][2]int
    start := 0
    for {
        idx := strings.Index(searchLine[start:], searchPattern)
        if idx == -1 {
            break
        }
        actualStart := start + idx
        actualEnd := actualStart + len(pattern)
        positions = append(positions, [2]int{actualStart, actualEnd})
        start = actualEnd
    }

    // 从后向前插入颜色码
    result := line
    for i := len(positions) - 1; i >= 0; i-- {
        start, end := positions[i][0], positions[i][1]
        matched := result[start:end]
        result = result[:start] + ColorRed + matched + ColorReset + result[end:]
    }
    return result
}
```

### 3.4 CLI 标志定义

**`internal/cli/grep.go`**

```go
var (
    // ... 现有标志 ...
    grepNoColor *qflag.BoolFlag // --no-color 禁用颜色
)

func init() {
    // ...
    grepNoColor = GrepCmd.Bool("no-color", "", "禁用颜色高亮", false)
    
    cmdOpts := &qflag.CmdOpts{
        // ...
        Notes: []string{
            "默认使用固定字符串匹配，更轻量高效",
            "使用 -E 启用正则表达式匹配",
            "默认高亮显示匹配关键字，使用 --no-color 禁用",
            "不指定 file 时从标准输入读取",
            "单文件场景，目录遍历请结合 find 使用",
        },
    }
}

func runGrep(cmd qflag.Command) error {
    // ...
    config := grep.GrepConfig{
        // ...
        NoColor: grepNoColor.Get(),
    }
    return grep.GrepCmdMain(config)
}
```

### 3.5 输出函数修改

**修改 `printLine` 函数，在输出前调用高亮：**

```go
func printLine(line string, lineNum int, separator string, config *GrepConfig) {
    if config.Count {
        return
    }

    // 高亮匹配内容
    highlightedLine := highlightLine(line, config)

    var parts []string
    if config.WithFilename {
        parts = append(parts, config.filename)
    }
    if config.LineNumber {
        parts = append(parts, strconv.Itoa(lineNum))
    }
    parts = append(parts, highlightedLine)

    // 输出逻辑...
    fmt.Println(result)
}
```

## 4. 边界情况处理

### 4.1 重叠匹配
- 正则模式：`FindAllStringIndex` 自动处理非重叠匹配
- 固定字符串：顺序查找，天然非重叠

### 4.2 特殊字符
- 颜色码本身不会参与匹配（高亮在匹配之后）
- 行号、文件名不高亮

### 4.3 性能考虑
- 从后向前插入颜色码，避免位置偏移计算
- 计数模式跳过颜色处理
- 禁用颜色时直接返回原字符串

## 5. 示例输出

```bash
# 默认（带颜色）
$ fck grep "error" log.txt
[红色]error[重置]: something went wrong

# 禁用颜色
$ fck grep --no-color "error" log.txt
error: something went wrong

# 多个匹配
$ fck grep -i "error" log.txt
[红色]Error[重置]: something went wrong
[红色]error[重置]: another [红色]error[重置]
```

## 6. 实现步骤

1. [ ] 添加 `NoColor` 字段到 `GrepConfig`
2. [ ] 实现 `highlightLine` 高亮函数
3. [ ] 添加 `--no-color` 标志到 CLI
4. [ ] 修改 `printLine` 调用高亮
5. [ ] 测试各种场景（正则/固定字符串、忽略大小写、多个匹配）
