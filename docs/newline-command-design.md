# newline 换行符检测与转换命令设计方案

## 一、需求概述

实现一个文件换行符检测和转换的子命令，支持：
- 检测文件换行符类型（LF、CRLF、CR、Mixed）
- 转换换行符格式（Unix ↔ Windows ↔ Mac）
- 批量处理多个文件
- 原地写入支持备份

## 二、命令命名

**主命令**: `newline`（完整）或 `nl`（简写）

考虑到 FCK 已有 `nl`（可能冲突），建议使用 `newline` 或 `eol`（End Of Line）。

参考 Unix 传统工具，使用 `dos2unix`/`unix2dos` 风格或统一的 `newline` 命令。

**建议**: 使用 `newline` 作为主命令，语义清晰。

## 三、命令行参数设计

### 3.1 标志定义

```go
var (
    // 输入参数
    newlineTo      *qflag.EnumFlag   // -t, --to      目标换行符格式
    newlineWrite   *qflag.BoolFlag   // -w, --write   原地写入
    newlineBackup  *qflag.BoolFlag   // -b, --backup  创建备份
    newlineOutput  *qflag.StringFlag // -o, --output  输出路径
    
    // 其他选项
    newlineList    *qflag.BoolFlag   // -l, --list    列出支持的换行符类型
    newlineQuiet   *qflag.BoolFlag   // -q, --quiet   静默模式
    newlineForce   *qflag.BoolFlag   // -f, --force   强制转换（即使已经是目标格式）
)
```

### 3.2 支持的换行符类型（Enum）

定义在 `types` 包：

```go
const (
    NewlineAuto   = "auto"   // 自动检测（默认）
    NewlineLF     = "lf"     // Unix/Linux/macOS (\n)
    NewlineCRLF   = "crlf"   // Windows (\r\n)
    NewlineCR     = "cr"     // Classic Mac (\r)
    NewlineNone   = "none"   // 不转换，仅检测
)

var SupportedNewlines = []string{
    NewlineAuto, NewlineLF, NewlineCRLF, NewlineCR, NewlineNone,
}

var TargetNewlines = []string{
    NewlineLF, NewlineCRLF, NewlineCR, NewlineNone,
}
```

### 3.3 命令配置

```go
cmdOpts := &qflag.CmdOpts{
    Desc:        "文件换行符检测与转换工具",
    UseChinese:  true,
    UsageSyntax: fmt.Sprintf("%s newline [options] <file...>", qflag.Root.Name()),
    Examples: map[string]string{
        "检测文件换行符":        fmt.Sprintf("%s newline file.txt", qflag.Root.Name()),
        "检测多个文件":         fmt.Sprintf("%s newline file1.txt file2.txt", qflag.Root.Name()),
        "转换为 Unix 格式":    fmt.Sprintf("%s newline -t lf file.txt", qflag.Root.Name()),
        "转换为 Windows 格式": fmt.Sprintf("%s newline -t crlf file.txt", qflag.Root.Name()),
        "原地转换":           fmt.Sprintf("%s newline -t lf -w file.txt", qflag.Root.Name()),
        "原地转换并备份":      fmt.Sprintf("%s newline -t lf -w -b file.txt", qflag.Root.Name()),
        "批量转换":           fmt.Sprintf("%s newline -t lf *.txt", qflag.Root.Name()),
        "列出支持的类型":      fmt.Sprintf("%s newline -l", qflag.Root.Name()),
    },
    Notes: []string{
        "-t 默认为 none (只检测不转换), 指定目标格式时执行转换",
        "使用 -w 原地写入会直接覆盖原文件",
        "建议重要文件使用 -b 创建备份",
        "Mixed 表示文件中存在多种换行符",
    },
    MutexGroups: []qflag.MutexGroup{
        {
            Name:      "to-list",
            Flags:     []string{"list", "to"},
            AllowNone: true,
        },
    },
    FlagDependencies: []qflag.FlagDependency{
        {
            Name:    "backup-write",
            Trigger: "backup",
            Targets: []string{"write"},
            Type:    qflag.DepRequired,
        },
        {
            Name:    "write-to",
            Trigger: "write",
            Targets: []string{"to"},
            Type:    qflag.DepRequired,
        },
    },
}
```

## 四、模块划分与代码结构

```
internal/
├── cli/
│   └── newline.go          # CLI 定义和参数解析
├── commands/
│   └── newline/
│       ├── cmd_newline.go  # 主逻辑和文件处理
│       ├── detector.go     # 换行符检测逻辑
│       └── converter.go    # 换行符转换逻辑
```

## 五、核心功能实现逻辑

### 5.1 换行符检测

```go
// NewlineType 换行符类型
type NewlineType string

const (
    TypeLF    NewlineType = "LF"      // \n
    TypeCRLF  NewlineType = "CRLF"    // \r\n
    TypeCR    NewlineType = "CR"      // \r
    TypeMixed NewlineType = "Mixed"   // 混合
    TypeNone  NewlineType = "None"    // 无换行符
)

// DetectionResult 检测结果
type DetectionResult struct {
    FilePath   string
    Type       NewlineType
    LFCount    int
    CRLFCount  int
    CRCount    int
    TotalLines int
}

// DetectNewline 检测文件换行符类型
func DetectNewline(data []byte) DetectionResult {
    lfCount := bytes.Count(data, []byte{'\n'})
    crlfCount := bytes.Count(data, []byte{"\r\n"})
    crCount := bytes.Count(data, []byte{'\r'})
    
    // 纯 CRLF
    if crlfCount > 0 && lfCount == crlfCount && crCount == crlfCount {
        return DetectionResult{Type: TypeCRLF, CRLFCount: crlfCount}
    }
    
    // 纯 LF
    if lfCount > 0 && crlfCount == 0 && crCount == 0 {
        return DetectionResult{Type: TypeLF, LFCount: lfCount}
    }
    
    // 纯 CR
    if crCount > 0 && lfCount == 0 && crlfCount == 0 {
        return DetectionResult{Type: TypeCR, CRCount: crCount}
    }
    
    // 混合
    if lfCount > 0 || crCount > 0 {
        return DetectionResult{
            Type:      TypeMixed,
            LFCount:   lfCount - crlfCount,
            CRLFCount: crlfCount,
            CRCount:   crCount - crlfCount,
        }
    }
    
    return DetectionResult{Type: TypeNone}
}
```

### 5.2 换行符转换

```go
// ConvertNewline 转换换行符
func ConvertNewline(data []byte, targetType string) []byte {
    // 先统一转换为 LF
    normalized := normalizeToLF(data)
    
    switch targetType {
    case NewlineLF:
        return normalized
    case NewlineCRLF:
        return bytes.ReplaceAll(normalized, []byte{'\n'}, []byte{"\r\n"})
    case NewlineCR:
        return bytes.ReplaceAll(normalized, []byte{'\n'}, []byte{'\r'})
    default:
        return data
    }
}

// normalizeToLF 将所有换行符统一为 LF
func normalizeToLF(data []byte) []byte {
    // 先将 CRLF 转为 LF
    result := bytes.ReplaceAll(data, []byte{"\r\n"}, []byte{'\n'})
    // 再将 CR 转为 LF
    result = bytes.ReplaceAll(result, []byte{'\r'}, []byte{'\n'})
    return result
}
```

### 5.3 文件处理流程

参考 iconv 命令的实现模式：

```go
// NewlineConfig 配置
type NewlineConfig struct {
    Files      []string
    ToNewline  string
    Write      bool
    Backup     bool
    Output     string
    Quiet      bool
    Force      bool
}

// ConvertFile 转换单个文件
func ConvertFile(srcPath string, config NewlineConfig) (ConversionResult, error) {
    // 1. 读取文件
    data, err := os.ReadFile(srcPath)
    if err != nil {
        return result, err
    }
    
    // 2. 检测当前换行符
    detection := DetectNewline(data)
    result.FromType = detection.Type
    
    // 3. 目标为 none：只检测
    if config.ToNewline == NewlineNone {
        result.Success = true
        if !config.Quiet {
            printDetectionResult(srcPath, detection)
        }
        return result, nil
    }
    
    // 4. 已经是目标格式且非强制模式
    if string(detection.Type) == strings.ToUpper(config.ToNewline) && !config.Force {
        result.Success = true
        result.ToType = detection.Type
        if !config.Quiet {
            fmt.Printf("%s: %s (no conversion needed)\n", srcPath, detection.Type)
        }
        return result, nil
    }
    
    // 5. 执行转换
    converted := ConvertNewline(data, config.ToNewline)
    result.ToType = NewlineType(strings.ToUpper(config.ToNewline))
    
    // 6. 写入输出
    // ... 类似 iconv 的实现
}
```

## 六、输出格式

### 6.1 检测模式输出

```
$ fck newline file.txt
file.txt: CRLF (100 lines)

$ fck newline mixed.txt
mixed.txt: Mixed (LF: 50, CRLF: 30, CR: 20)
```

### 6.2 转换模式输出

```
$ fck newline -t lf file.txt
file.txt: CRLF → LF (100 lines)

$ fck newline -t lf -w file.txt
file.txt: CRLF → LF (in-place)
```

### 6.3 统计信息（多个文件）

```
Total: 10 files
Success: 8 files
Skipped: 1 files
Failed: 1 files
```

## 七、错误处理

1. **文件读取错误**: 返回错误，非静默模式下打印到 stderr
2. **目录作为输入**: 报错 "is a directory, not a file"
3. **输出路径不存在**: 报错 "output path does not exist"
4. **权限错误**: 返回系统错误信息

## 八、与现有命令的对比

| 特性 | iconv | newline |
|------|-------|---------|
| 检测内容 | 文件编码 | 换行符类型 |
| 转换内容 | 编码格式 | 换行符格式 |
| 原地写入 | -w | -w |
| 备份 | -b | -b |
| 静默模式 | -q | -q |
| 输出路径 | -o | -o |
| 批量处理 | 支持通配符 | 支持通配符 |

## 九、实现建议

1. **代码复用**: 参考 `iconv` 命令的文件处理逻辑，保持一致的代码风格
2. **常量定义**: 在 `types` 包定义换行符相关常量
3. **后缀常量**: 
   - `NewlineTempExt = ".newline.tmp"`
   - `NewlineConvertedExt = ".unix"` 或 `.win`（根据转换类型）
   - `NewlineBackupExt = ".bak"`
4. **大文件处理**: 使用流式读取，避免一次性加载大文件到内存

## 十、测试用例建议

1. **检测测试**: LF、CRLF、CR、Mixed、None 各种情况
2. **转换测试**: 各种格式之间的相互转换
3. **边界测试**: 空文件、单行文件、无换行符文件
4. **错误测试**: 目录输入、权限不足、路径不存在
5. **批量测试**: 通配符匹配、多个文件统计

---

**设计日期**: 2026-05-11  
**参考命令**: iconv  
**预计工作量**: 中等（参考 iconv 实现，约 400-500 行代码）
