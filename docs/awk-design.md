# awk 子命令设计方案

> **设计日期**: 2026-04-18  
> **设计目标**: 跨平台简化版 awk，使用标志控制替代程序解析  
> **遵循规范**: qflag 命令开发规范

---

## 一、功能范围

### 1.1 支持的功能

| 功能类别 | 具体功能 | 实现方式 |
|----------|----------|----------|
| **模式匹配** | 正则匹配 | `-p, --pattern` 标志 |
| **字段处理** | 指定字段输出 | `-f, --field` 标志（可多次使用） |
| **打印** | 打印整行或指定字段 | 根据 `-f` 参数决定 |
| **内置变量** | 行号(NR)、字段数(NF) | 特殊字段索引支持 |

### 1.2 不支持的功能（简化）

- 算术运算
- 条件语句（if/for/while）
- 内置函数（length/substr/gsub等）
- 多动作
- BEGIN/END 块
- 复杂的字段条件表达式

---

## 二、CLI 设计

### 2.1 命令定义

```go
// internal/cli/awk.go
var AwkCmd *qflag.Cmd

var (
    awkPattern     *qflag.StringFlag  // -p, --pattern  正则匹配模式
    awkField       *qflag.IntSliceFlag // -f, --field    字段索引（可多次使用）
    awkSeparator   *qflag.StringFlag  // -F, --separator 输入分隔符
    awkOutputSep   *qflag.StringFlag  // -O, --output-separator 输出分隔符
    awkLineNum     *qflag.BoolFlag    // -n, --line-number 显示行号
)
```

### 2.2 标志详情

| 标志 | 短标志 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| `--pattern` | `-p` | string | "" | 正则匹配模式，匹配行才处理 |
| `--field` | `-f` | int[] | [] | 要输出的字段索引（1-based），可多次使用 |
| `--separator` | `-F` | string | " " | 输入字段分隔符 |
| `--output-separator` | `-O` | string | " " | 输出字段分隔符 |
| `--line-number` | `-n` | bool | false | 输出行号 |

### 2.3 特殊字段索引

| 索引值 | 含义 | 说明 |
|--------|------|------|
| `0` | 整行 | 等同于 `$0` |
| `-1` | NF | 最后一个字段（动态计算） |
| `1~N` | 第N个字段 | 正常字段索引 |

---

## 三、使用示例

### 3.1 基础用法

```bash
# 打印第1列（默认空格分隔）
fck awk -f 1 file.txt

# 打印第1和第3列
fck awk -f 1 -f 3 file.txt

# 指定分隔符，打印第1和第3列
fck awk -F: -f 1 -f 3 /etc/passwd

# 打印整行
fck awk -f 0 file.txt

# 打印最后一列
fck awk -f -1 file.txt
```

### 3.2 模式匹配

```bash
# 正则匹配，打印包含 error 的整行
fck awk -p "error" -f 0 log.txt

# 正则匹配，打印匹配行的第2列
fck awk -p "error" -f 2 log.txt

# 匹配以 # 开头的行（注释行）
fck awk -p "^#" -f 0 config.txt
```

### 3.3 行号输出

```bash
# 打印行号和第1列
fck awk -n -f 1 file.txt

# 打印行号、第1列和最后一列
fck awk -n -f 1 -f -1 file.txt
```

### 3.4 分隔符定制

```bash
# CSV 文件处理（逗号分隔）
fck awk -F, -f 1 -f 3 data.csv

# 制表符分隔
fck awk -F"\t" -f 2 file.tsv

# 自定义输出格式
fck awk -F: -O " -> " -f 1 -f 6 /etc/passwd
```

### 3.5 组合使用

```bash
# 匹配包含 root 的行，输出第1和第6列，用 " -> " 连接
fck awk -F: -p "root" -f 1 -f 6 -O " -> " /etc/passwd

# 显示行号，制表符分隔文件，输出第2和最后一列
fck awk -n -F"\t" -f 2 -f -1 file.tsv
```

---

## 四、目录结构

```
internal/
├── commands/
│   └── awk/
│       ├── types.go      # 配置结构体定义
│       └── cmd_awk.go    # 主逻辑实现
└── cli/
    └── awk.go            # CLI 定义
```

---

## 五、代码设计

### 5.1 配置结构体（types.go）

```go
package awk

// AwkConfig 命令配置
type AwkConfig struct {
    Pattern     string   // 正则匹配模式
    Fields      []int    // 字段索引列表（0=整行, -1=NF, 1+=正常字段）
    FieldSep    string   // 输入分隔符
    OutputSep   string   // 输出分隔符
    ShowLineNum bool     // 显示行号
    Files       []string // 输入文件列表
}

// AwkStats 执行统计
type AwkStats struct {
    ProcessedLines int // 处理行数
    MatchedLines   int // 匹配行数
    OutputLines    int // 输出行数
}
```

### 5.2 主逻辑（cmd_awk.go）

```go
package awk

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "regexp"
    "strings"
)

// AwkCmdMain 执行 awk 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func AwkCmdMain(config AwkConfig) error {
    // 编译正则表达式（如果提供了模式）
    var re *regexp.Regexp
    if config.Pattern != "" {
        var err error
        re, err = regexp.Compile(config.Pattern)
        if err != nil {
            return fmt.Errorf("invalid pattern: %w", err)
        }
    }

    stats := &AwkStats{}

    // 处理文件或标准输入
    if len(config.Files) == 0 {
        return processInput(os.Stdin, "-", config, re, stats)
    }

    for _, file := range config.Files {
        if err := processFile(file, config, re, stats); err != nil {
            return err
        }
    }

    return nil
}

// processFile 处理单个文件
func processFile(filename string, config AwkConfig, re *regexp.Regexp, stats *AwkStats) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("cannot open file %s: %w", filename, err)
    }
    defer file.Close()

    return processInput(file, filename, config, re, stats)
}

// processInput 处理输入流（支持大文件和管道）
//
// 使用 bufio.Reader 替代 Scanner，支持任意行长度
func processInput(input *os.File, name string, config AwkConfig, re *regexp.Regexp, stats *AwkStats) error {
    reader := bufio.NewReader(input)
    lineNum := 0

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                // 处理最后一行（可能没有换行符）
                if len(line) > 0 {
                    lineNum++
                    if err := processLine(lineNum, line, config, re, stats); err != nil {
                        return err
                    }
                }
                break
            }
            return fmt.Errorf("read error: %w", err)
        }

        // 去掉末尾换行符
        line = strings.TrimSuffix(line, "\n")
        line = strings.TrimSuffix(line, "\r") // Windows 兼容

        lineNum++
        if err := processLine(lineNum, line, config, re, stats); err != nil {
            return err
        }
    }

    return nil
}

// processLine 处理单行
func processLine(lineNum int, line string, config AwkConfig, re *regexp.Regexp, stats *AwkStats) error {
    stats.ProcessedLines++

    // 模式匹配检查
    if re != nil && !re.MatchString(line) {
        return nil
    }
    stats.MatchedLines++

    // 分割字段
    fields := strings.Split(line, config.FieldSep)
    nf := len(fields)

    // 输出处理
    if err := outputLine(lineNum, line, fields, nf, config); err != nil {
        return err
    }
    stats.OutputLines++

    return nil
}

// outputLine 输出行
func outputLine(nr int, line string, fields []string, nf int, config AwkConfig) error {
    var parts []string

    // 添加行号
    if config.ShowLineNum {
        parts = append(parts, fmt.Sprintf("%d", nr))
    }

    // 添加字段
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

    // 输出
    fmt.Println(strings.Join(parts, config.OutputSep))
    return nil
}
```

### 5.3 CLI 定义（awk.go）

```go
package cli

import (
    "fmt"

    "gitee.com/MM-Q/fck/internal/commands/awk"
    "gitee.com/MM-Q/qflag"
)

var AwkCmd *qflag.Cmd

var (
    awkPattern   *qflag.StringFlag   // -p, --pattern  正则匹配模式
    awkField     *qflag.IntSliceFlag // -f, --field    字段索引
    awkSeparator *qflag.StringFlag   // -F, --separator 输入分隔符
    awkOutputSep *qflag.StringFlag   // -O, --output-separator 输出分隔符
    awkLineNum   *qflag.BoolFlag     // -n, --line-number 显示行号
)

func init() {
    AwkCmd = qflag.NewCmd("awk", "awk", qflag.ExitOnError)

    awkPattern = AwkCmd.String("pattern", "p", "正则匹配模式，只处理匹配的行", "")
    awkField = AwkCmd.IntSlice("field", "f", "要输出的字段索引（1-based，0=整行，-1=最后一个字段，可多次使用）", nil)
    awkSeparator = AwkCmd.String("separator", "F", "输入字段分隔符", " ")
    awkOutputSep = AwkCmd.String("output-separator", "O", "输出字段分隔符", " ")
    awkLineNum = AwkCmd.Bool("line-number", "n", "显示行号", false)

    cmdOpts := &qflag.CmdOpts{
        Desc:       "文本字段处理工具",
        UseChinese: true,
        Examples: map[string]string{
            "打印第1列":                    "fck awk -f 1 file.txt",
            "打印第1和第3列":               "fck awk -f 1 -f 3 file.txt",
            "指定分隔符":                   "fck awk -F: -f 1 -f 3 /etc/passwd",
            "打印最后一列":                 "fck awk -f -1 file.txt",
            "正则匹配并打印":               "fck awk -p \"error\" -f 0 log.txt",
            "显示行号":                     "fck awk -n -f 1 file.txt",
            "CSV处理":                      "fck awk -F, -f 1 -f 3 data.csv",
            "制表符分隔":                   "fck awk -F\"\\t\" -f 2 file.tsv",
            "自定义输出格式":               "fck awk -F: -O\" -> \" -f 1 -f 6 /etc/passwd",
            "管道输入":                     "cat file.txt | fck awk -f 1",
        },
        Notes: []string{
            "字段索引从1开始，0表示整行，-1表示最后一个字段",
            "支持从管道读取输入，此时不需要指定文件参数",
            "正则匹配使用 Go 正则语法",
            "默认分隔符为空格，支持任意字符串作为分隔符",
        },
    }

    if err := AwkCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }

    AwkCmd.SetRun(runAwk)
}

func runAwk(cmd qflag.Command) error {
    config := awk.AwkConfig{
        Pattern:     awkPattern.Get(),
        Fields:      awkField.Get(),
        FieldSep:    awkSeparator.Get(),
        OutputSep:   awkOutputSep.Get(),
        ShowLineNum: awkLineNum.Get(),
        Files:       cmd.Args(),
    }

    // 如果没有指定字段，默认输出整行
    if len(config.Fields) == 0 {
        config.Fields = []int{0}
    }

    return awk.AwkCmdMain(config)
}
```

---

## 六、注册命令

在 `internal/cli/root.go` 的 `SubCmds` 列表中添加：

```go
SubCmds: []qflag.Command{
    // 文件操作类
    CpCmd,
    MvCmd,
    RmCmd,
    MkdirCmd,
    TouchCmd,
    TruncateCmd,
    CatCmd,
    ListCmd,
    
    // 文本处理类
    FindCmd,
    GrepCmd,
    SedCmd,
    AwkCmd,  // <-- 新增
    
    // ... 其他命令
},
```

---

## 七、与传统 awk 对比

| 传统 awk | fck awk |
|----------|---------|
| `awk '{print $1}' file` | `fck awk -f 1 file` |
| `awk '{print $1,$3}' file` | `fck awk -f 1 -f 3 file` |
| `awk -F: '{print $1}' file` | `fck awk -F: -f 1 file` |
| `awk '/pattern/ {print $2}' file` | `fck awk -p pattern -f 2 file` |
| `awk '{print NR,$1}' file` | `fck awk -n -f 1 file` |
| `awk '{print $NF}' file` | `fck awk -f -1 file` |
| `awk -F: -O' -> ' '{print $1,$6}' file` | `fck awk -F: -O" -> " -f 1 -f 6 file` |

---

## 八、开发计划

| 阶段 | 任务 | 预计时间 |
|------|------|----------|
| 1 | 创建目录结构 | 10分钟 |
| 2 | 实现 types.go | 10分钟 |
| 3 | 实现 cmd_awk.go | 40分钟 |
| 4 | 实现 awk.go（CLI） | 20分钟 |
| 5 | 注册命令 | 5分钟 |
| 6 | 编译验证 | 10分钟 |
| 7 | 功能测试 | 15分钟 |

**总计：约 1.5 小时**

---

## 九、注意事项

1. **字段索引**：采用 1-based 索引（与传统 awk 一致），0 表示整行，-1 表示最后一个字段
2. **正则语法**：使用 Go 标准正则语法（RE2）
3. **分隔符**：支持任意字符串，包括特殊字符如 `\t`（制表符）
4. **管道支持**：无文件参数时从标准输入读取
5. **性能**：使用 bufio.Reader，支持任意行长度，适合处理超大文件和管道输出
6. **错误处理**：提供友好的错误信息，包括文件打开失败、正则编译失败等
