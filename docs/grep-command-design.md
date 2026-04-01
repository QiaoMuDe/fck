# Grep 子命令设计方案

## 1. 功能概述

实现类 Unix grep 功能的文本搜索工具，专注于单文件内容过滤。目录遍历和多文件处理建议结合 `fck find -exec` 或管道使用。

## 2. 核心功能

### 2.1 基础搜索
- 在单个文件中搜索匹配正则表达式或固定字符串的行
- 支持从标准输入读取
- 支持限制最大匹配数

### 2.2 匹配模式
- 固定字符串匹配（默认，更轻量高效）
- 正则表达式匹配（-E, --regexp）
- 忽略大小写（-i, --ignore-case）

### 2.3 输出控制
- 显示行号（-n, --line-number）
- 显示文件名（-H 启用，格式：文件名:行号:内容）
- 显示不匹配的行（-v, --invert-match）
- 显示匹配计数（-c, --count）
- 最多匹配数量限制（-m, --max-count）
- 显示匹配前后上下文（-B, -A, -C）

## 3. 命令行接口设计

### 3.1 基本用法
```bash
fck grep [options] <pattern> [file]
```

**说明：**
- 不指定 file 时从标准输入读取
- 单文件场景，不涉及目录遍历

### 3.2 标志设计

| 短标志 | 长标志 | 类型 | 默认值 | 说明 |
|-------|-------|------|--------|------|
| -i | --ignore-case | bool | false | 忽略大小写 |
| -v | --invert-match | bool | false | 显示不匹配的行 |
| -n | --line-number | bool | false | 显示行号 |
| -c | --count | bool | false | 只显示匹配行数 |
| -H | --with-filename | bool | false | 显示文件名（格式：文件名:行号:内容） |
| -E | --regexp | bool | false | 使用正则表达式匹配（默认固定字符串） |
| -m | --max-count | int | 0 | 最大匹配行数（0表示无限制） |
| -B | --before-context | int | 0 | 显示匹配行前N行 |
| -A | --after-context | int | 0 | 显示匹配行后N行 |
| -C | --context | int | 0 | 显示匹配行前后N行（等价于 -B N -A N） |

### 3.3 输出格式

**默认（仅匹配内容）：**
```
匹配内容
```

**-n（显示行号）：**
```
行号:匹配内容
```

**-H（显示文件名）：**
```
文件名:行号:匹配内容
```

**-c（仅计数）：**
```
计数
```

### 3.4 示例

```bash
# 基础搜索
fck grep "error" log.txt

# 忽略大小写
fck grep -i "error" log.txt

# 显示行号
fck grep -n "func" main.go

# 显示文件名和行号
fck grep -nH "TODO" main.go

# 只显示计数
fck grep -c "http" access.log

# 显示不匹配的行
fck grep -v "^#" config.txt

# 正则表达式匹配
fck grep -E "^error|fatal" log.txt

# 整词匹配（使用 \b）
fck grep -E "\bmain\b" *.go

# 整行匹配（使用 ^$）
fck grep -E "^hello$" file.txt

# 限制匹配数量
fck grep -m 10 "http" access.log

# 多个模式（使用正则 |）
fck grep "error|warning" log.txt

# 结合 find 使用（推荐方式）
fck find . -name "*.go" -exec fck grep -n "TODO" {} \;

# 结合管道使用
cat log.txt | fck grep -i "warning"
```

## 4. 技术实现

### 4.1 文件结构
```
internal/
├── commands/grep/
│   └── cmd_grep.go      # 核心逻辑
└── cli/
    └── grep.go          # CLI 接口
```

### 4.2 核心数据结构

```go
// GrepConfig grep 命令配置
type GrepConfig struct {
    // CLI 参数
    Pattern       string   // 搜索模式
    Target        string   // 目标文件（空表示标准输入）
    IgnoreCase    bool     // -i 忽略大小写
    InvertMatch   bool     // -v 反向匹配
    LineNumber    bool     // -n 显示行号
    Count         bool     // -c 只显示计数
    WithFilename  bool     // -H 显示文件名
    Regexp        bool     // -E 使用正则表达式
    MaxCount      int      // -m 最大匹配数
    BeforeContext int      // -B 前文行数
    AfterContext  int      // -A 后文行数
    Context       int      // -C 上下文行数

    // 运行时
    compiledPattern *regexp.Regexp // 编译后的正则
    matchCount       int              // 当前匹配计数
    filename         string           // 实际文件名（用于输出）
}
```

### 4.3 核心算法

#### 4.3.1 模式编译
```go
func compilePattern(pattern string, config *GrepConfig) (*regexp.Regexp, error) {
    // 非正则模式：转义特殊字符
    if !config.Regexp {
        pattern = regexp.QuoteMeta(pattern)
    }

    // 编译正则
    flags := ""
    if config.IgnoreCase {
        flags = "(?i)"
    }

    return regexp.Compile(flags + pattern)
}
```

#### 4.3.2 匹配逻辑
```go
func matchLine(line string, re *regexp.Regexp) bool {
    return re.MatchString(line)
}
```

#### 4.3.3 主处理流程
```go
func GrepCmdMain(config GrepConfig) error {
    // 1. 编译模式
    re, err := compilePattern(config.Pattern, &config)
    if err != nil {
        return fmt.Errorf("invalid pattern: %w", err)
    }
    config.compiledPattern = re
    
    // 2. 确定输入源
    var reader io.Reader
    if config.Target == "" {
        reader = os.Stdin
        config.filename = "(standard input)"
    } else {
        file, err := os.Open(config.Target)
        if err != nil {
            return err
        }
        defer file.Close()
        reader = file
        config.filename = config.Target
    }
    
    // 3. 扫描处理
    scanner := bufio.NewScanner(reader)
    lineNum := 0
    
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        
        // 匹配检查
        matched := matchLine(line, config.compiledPattern)
        if config.InvertMatch {
            matched = !matched
        }
        
        if matched {
            config.matchCount++
            
            // 达到最大匹配数则停止
            if config.MaxCount > 0 && config.matchCount > config.MaxCount {
                break
            }
            
            // 输出结果
            if err := outputResult(line, lineNum, &config); err != nil {
                return err
            }
        }
    }
    
    // 4. 计数模式输出
    if config.Count {
        fmt.Println(config.matchCount)
    }
    
    return scanner.Err()
}
```

#### 4.3.4 输出格式化
```go
func outputResult(line string, lineNum int, config *GrepConfig) error {
    // 计数模式不输出具体内容
    if config.Count {
        return nil
    }
    
    var parts []string
    
    // 文件名
    if config.WithFilename {
        parts = append(parts, config.filename)
    }
    
    // 行号
    if config.LineNumber {
        parts = append(parts, strconv.Itoa(lineNum))
    }
    
    // 内容
    parts = append(parts, line)
    
    fmt.Println(strings.Join(parts, ":"))
    return nil
}
```

### 4.4 性能考虑
- 使用 `bufio.Scanner` 逐行读取，内存占用低
- 正则预编译，避免重复编译
- 达到 MaxCount 时立即终止，避免无效读取

## 5. 错误处理

| 错误场景 | 处理方式 |
|---------|---------|
| 无搜索模式 | 返回错误 "no pattern specified" |
| 正则表达式无效 | 返回编译错误 |
| 文件不存在 | 返回错误 "file not found" |
| 目标是目录 | 返回错误 "is a directory, use 'find -exec' instead" |

## 6. 与 Unix grep 的差异

1. **单文件限制**：不直接支持多文件和递归，推荐结合 `find -exec` 使用
2. **文件名显示**：默认不显示，需显式使用 `-H`
3. **输出格式**：统一使用 `:` 分隔（文件名:行号:内容）
4. **不支持**：-l/-L（单文件无意义）、-r（递归）、--include/exclude（文件过滤）

## 7. 实现步骤

### P0（核心功能）
- [ ] 基础正则匹配
- [ ] 标准输入/文件读取
- [ ] -i, -v, -n, -c 选项
- [ ] -m 最大匹配数

### P1（增强匹配）
- [ ] -w 整词匹配
- [ ] -x 整行匹配
- [ ] -F 固定字符串
- [ ] -e 多模式

### P2（输出控制）
- [ ] -H 文件名显示

## 8. 测试用例建议

1. **基础匹配**：普通字符串、正则元字符、中文
2. **选项组合**：-i + -v, -n + -c, -H + -n 等
3. **边界情况**：空文件、无匹配、全匹配、最大匹配数限制
4. **标准输入**：管道输入测试
5. **错误处理**：无效正则、目录输入、不存在文件
