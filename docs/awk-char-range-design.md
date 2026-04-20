# awk 命令字符范围提取功能设计方案

## 1. 设计目标

为 `fck awk` 命令添加字符范围提取功能，使其能够替代 Linux `cut -c` 命令的功能。

## 2. 方案概述

新增 `-c, --chars` 参数用于字符范围提取，与现有的 `-f, --field` 参数互斥。

## 3. 使用示例

```bash
# 提取第1-10个字符
fck awk -c 1-10 file.txt

# 提取第5个字符到行尾
fck awk -c 5- file.txt

# 提取前20个字符
fck awk -c -20 file.txt

# 提取多个范围（第1、5-10、15个字符）
fck awk -c 1,5-10,15 file.txt

# 与正则匹配组合使用
fck awk -p "ERROR" -c 1-50 log.txt

# 显示行号 + 字符提取
fck awk -n -c 1-20 file.txt

# 自定义输出分隔符
fck awk -c 1-5,10-15 -O" | " file.txt
```

## 4. 数据结构设计

### 4.1 新增类型定义

```go
// CharRange 字符范围定义
type CharRange struct {
    Start int  // 起始位置（1-based，包含）
    End   int  // 结束位置（1-based，包含），-1表示到行尾
}

// IsValid 检查范围是否有效
func (cr CharRange) IsValid() bool {
    if cr.Start < 1 {
        return false
    }
    if cr.End != -1 && cr.End < cr.Start {
        return false
    }
    return true
}

// Extract 从字符串中提取字符范围
func (cr CharRange) Extract(s string) string {
    runes := []rune(s) // 支持Unicode
    start := cr.Start - 1 // 转换为0-based
    
    if start >= len(runes) {
        return ""
    }
    
    end := cr.End
    if end == -1 || end > len(runes) {
        end = len(runes)
    }
    
    return string(runes[start:end])
}
```

### 4.2 修改 AwkConfig

```go
// AwkConfig 命令配置
type AwkConfig struct {
    Pattern     string      // 正则匹配模式
    Fields      []int       // 字段索引列表（0=整行, -1=NF, 1+=正常字段）
    FieldSep    string      // 输入分隔符
    OutputSep   string      // 输出分隔符
    ShowLineNum bool        // 显示行号
    Files       []string    // 输入文件列表
    CharRanges  []CharRange // 字符范围列表（新增）
}

// IsFieldMode 是否为字段提取模式
func (c AwkConfig) IsFieldMode() bool {
    return len(c.Fields) > 0
}

// IsCharMode 是否为字符提取模式
func (c AwkConfig) IsCharMode() bool {
    return len(c.CharRanges) > 0
}
```

## 5. CLI 参数设计

### 5.1 新增参数

```go
var (
    awkPattern   *qflag.StringFlag   // -p, --pattern  正则匹配模式
    awkField     *qflag.IntSliceFlag // -f, --field    字段索引
    awkChars     *qflag.StringFlag   // -c, --chars    字符范围（新增）
    awkSeparator *qflag.StringFlag   // -F, --separator 输入分隔符
    awkOutputSep *qflag.StringFlag   // -O, --output-separator 输出分隔符
    awkLineNum   *qflag.BoolFlag     // -n, --line-number 显示行号
)
```

### 5.2 互斥组配置

```go
cmdOpts := &qflag.CmdOpts{
    Desc:       "文本字段处理工具",
    UseChinese: true,
    Examples: map[string]string{
        // 现有示例...
        "提取前10个字符":      "fck awk -c 1-10 file.txt",
        "提取第5个字符到行尾":  "fck awk -c 5- file.txt",
        "提取多个字符范围":     "fck awk -c 1,5-10,15 file.txt",
    },
    Notes: []string{
        "字段索引从1开始, 0表示整行, -1表示最后一个字段",
        "字符范围格式: N(单个字符), N-M(范围), N-(从N到结尾), -N(前N个)",
        "-f 和 -c 参数互斥，不能同时使用",
        "支持从管道读取输入, 此时不需要指定文件参数",
        "正则匹配使用 Go 正则语法",
        "默认分隔符为空格, 支持任意字符串作为分隔符",
    },
    MutexGroups: []qflag.MutexGroup{
        {
            Name:      "extract",
            Flags:     []string{"field", "chars"},
            AllowNone: true, // 允许都不指定（输出整行），但不能同时指定
        },
    },
}
```

## 6. 字符范围解析逻辑

```go
// ParseCharRanges 解析字符范围字符串
// 支持格式: "1-10,15,20-,5"
func ParseCharRanges(s string) ([]CharRange, error) {
    var ranges []CharRange
    parts := strings.Split(s, ",")
    
    for _, part := range parts {
        part = strings.TrimSpace(part)
        if part == "" {
            continue
        }
        
        cr, err := parseSingleRange(part)
        if err != nil {
            return nil, fmt.Errorf("invalid range %q: %w", part, err)
        }
        ranges = append(ranges, cr)
    }
    
    return ranges, nil
}

// parseSingleRange 解析单个范围
// 支持: "5" -> {5,5}, "1-10" -> {1,10}, "5-" -> {5,-1}, "-10" -> {1,10}
func parseSingleRange(s string) (CharRange, error) {
    // 处理 "-N" 格式（前N个字符）
    if strings.HasPrefix(s, "-") {
        n, err := strconv.Atoi(s[1:])
        if err != nil || n < 1 {
            return CharRange{}, fmt.Errorf("invalid range")
        }
        return CharRange{Start: 1, End: n}, nil
    }
    
    // 处理 "N-" 格式（从N到结尾）
    if strings.HasSuffix(s, "-") {
        n, err := strconv.Atoi(s[:len(s)-1])
        if err != nil || n < 1 {
            return CharRange{}, fmt.Errorf("invalid range")
        }
        return CharRange{Start: n, End: -1}, nil
    }
    
    // 处理 "N-M" 格式
    if strings.Contains(s, "-") {
        parts := strings.SplitN(s, "-", 2)
        start, err1 := strconv.Atoi(parts[0])
        end, err2 := strconv.Atoi(parts[1])
        if err1 != nil || err2 != nil || start < 1 || end < 1 || end < start {
            return CharRange{}, fmt.Errorf("invalid range")
        }
        return CharRange{Start: start, End: end}, nil
    }
    
    // 处理单个数字 "N"
    n, err := strconv.Atoi(s)
    if err != nil || n < 1 {
        return CharRange{}, fmt.Errorf("invalid range")
    }
    return CharRange{Start: n, End: n}, nil
}
```

## 7. 输出行处理修改

```go
// outputLine 输出行
func outputLine(nr int, line string, fields []string, nf int, config AwkConfig) error {
    var parts []string

    // 添加行号
    if config.ShowLineNum {
        parts = append(parts, fmt.Sprintf("%d", nr))
    }

    // 字符提取模式
    if config.IsCharMode() {
        for _, cr := range config.CharRanges {
            value := cr.Extract(line)
            parts = append(parts, value)
        }
    } else {
        // 字段提取模式（原有逻辑）
        for _, idx := range config.Fields {
            var value string
            switch idx {
            case 0: // 整行
                value = line
            case -1: // 最后一个字段
                if nf > 0 {
                    value = fields[nf-1]
                }
            default: // 正常字段（1-based）
                if idx > 0 && idx <= nf {
                    value = fields[idx-1]
                }
            }
            parts = append(parts, value)
        }
    }

    // 输出结果
    fmt.Println(strings.Join(parts, config.OutputSep))
    return nil
}
```

## 8. 边界情况处理

| 场景 | 处理方式 |
|------|----------|
| 行长度不足 | 返回能提取的部分，不报错 |
| 中文字符 | 使用 `[]rune` 按字符而非字节处理 |
| 空行 | 输出空字符串 |
| 范围重叠 | 按用户指定顺序输出，不去重 |
| 超大范围 | 自动截断到行尾 |

## 9. 测试用例

```go
// 测试解析
TestParseCharRanges("1-10", []CharRange{{1, 10}})
TestParseCharRanges("5-", []CharRange{{5, -1}})
TestParseCharRanges("-20", []CharRange{{1, 20}})
TestParseCharRanges("1,5-10,15", []CharRange{{1,1}, {5,10}, {15,15}})

// 测试提取
line := "Hello, 世界!"
CharRange{1, 5}.Extract(line)   // "Hello"
CharRange{8, 9}.Extract(line)   // "世界"
CharRange{5, -1}.Extract(line)  // "o, 世界!"
```

## 10. 实现步骤

1. 修改 `types.go`：添加 `CharRange` 类型和修改 `AwkConfig`
2. 修改 `cmd_awk.go`：
   - 添加 `ParseCharRanges` 函数
   - 修改 `outputLine` 函数支持字符模式
3. 修改 `awk.go`（CLI）：
   - 添加 `-c, --chars` 参数
   - 配置互斥组
   - 更新帮助信息和示例
4. 添加单元测试
5. 编译验证

## 11. 与 Linux cut 对比

| 功能 | Linux cut | fck awk |
|------|-----------|---------|
| 提取字段 | `cut -d: -f1,3` | `fck awk -F: -f 1,3` |
| 提取字符 | `cut -c 1-10` | `fck awk -c 1-10` |
| 多个范围 | `cut -c 1-5,10-15` | `fck awk -c 1-5,10-15` |
| 从N到结尾 | `cut -c 5-` | `fck awk -c 5-` |
| 正则过滤 | 不支持 | `fck awk -p "pattern" -c 1-10` |
| 显示行号 | 不支持 | `fck awk -n -c 1-10` |

## 12. 注意事项

1. **Unicode 支持**：使用 `[]rune` 而非 `[]byte` 处理，确保中文字符正确
2. **性能考虑**：字符提取不需要分割字段，比字段模式更快
3. **互斥验证**：确保 `-f` 和 `-c` 不会同时被指定
4. **向后兼容**：不指定 `-f` 或 `-c` 时保持原有行为（输出整行）
