# iconv 编码转换子命令设计方案

## 1. 功能概述

实现文件编码检测和转换功能，支持多种常见编码格式（UTF-8、GBK、GB2312、GB18030、Big5、Latin1 等）之间的相互转换。

## 2. 命令设计

### 2.1 命令结构

```
fck iconv [options] <file...>
```

### 2.2 命令行参数

| 短标志 | 长标志 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| -f | --from | enum | auto | 源文件编码（auto 表示自动检测） |
| -t | --to | enum | none | 目标编码（none 表示只检测不转换） |
| -w | --write | bool | false | 原地写入（覆盖原文件） |
| -b | --backup | bool | false | 原地写入时创建备份（.bak） |
| -l | --list | bool | false | 列出支持的编码格式 |
| -o | --output | string | "." | 输出路径（目录或文件，默认当前目录） |
| -v | --verbose | bool | false | 显示详细转换信息 |
| -q | --quiet | bool | false | 静默模式，不输出任何信息 |

### 2.3 支持的编码格式

```
none        只检测不转换 - types.EncodingNone（仅用于 -t）
auto        自动检测 - types.EncodingAuto（仅用于 -f）
UTF-8       UTF-8 编码（默认）- types.EncodingUTF8
UTF-8-BOM   带 BOM 的 UTF-8 - types.EncodingUTF8BOM
UTF-16      UTF-16 编码 - types.EncodingUTF16
UTF-16LE    UTF-16 小端序 - types.EncodingUTF16LE
UTF-16BE    UTF-16 大端序 - types.EncodingUTF16BE
GBK         简体中文（Windows 默认）- types.EncodingGBK
GB2312      简体中文（国标）- types.EncodingGB2312
GB18030     简体中文（最新国标）- types.EncodingGB18030
Big5        繁体中文 - types.EncodingBig5
Latin1      ISO-8859-1 - types.EncodingLatin1
ShiftJIS    日文 - types.EncodingShiftJIS
EUC-KR      韩文 - types.EncodingEUCKR
```

**编码常量定义位置：** `internal/types/types.go`

## 3. 模块划分

### 3.1 目录结构

```
internal/cli/
├── iconv.go              # CLI 层：命令定义和参数解析
└── iconv/
    └── (无需子命令，单一功能)

internal/commands/iconv/
├── cmd_iconv.go          # 核心逻辑：编码转换主入口
├── detector.go           # 编码检测功能
└── converter.go          # 编码转换功能
```

### 3.2 文件职责

| 文件 | 职责 |
|------|------|
| cli/iconv.go | 定义命令行参数、验证输入、调用核心逻辑 |
| commands/iconv/cmd_iconv.go | 配置结构体、主入口函数、流程控制 |
| commands/iconv/detector.go | 编码自动检测实现 |
| commands/iconv/converter.go | 编码转换实现、文件读写 |

## 4. 核心实现逻辑

### 4.1 cli/iconv.go

**编码列表常量定义：** `internal/types/types.go` 中的 `SupportedEncodings`，供 CLI 和核心逻辑共享使用。

### 4.2 配置结构体

```go
// IconvConfig iconv 命令配置
type IconvConfig struct {
    // 输入参数
    Files        []string // 目标文件列表（支持通配符和多个文件）
    FromEncoding string   // 源编码（auto 表示自动检测）
    ToEncoding   string   // 目标编码

    // 输出选项
    Write   bool   // 原地写入（覆盖原文件）
    Backup  bool   // 原地写入时创建备份
    Output  string // 输出路径（目录或文件，空表示当前目录）

    // 其他选项
    List    bool // 列出支持的编码
    Verbose bool // 详细输出
    Quiet   bool // 静默模式
}

// ConversionResult 转换结果
type ConversionResult struct {
    FilePath     string // 文件路径
    Success      bool   // 是否成功
    FromEncoding string // 检测/指定的源编码
    ToEncoding   string // 目标编码
    BytesRead    int64  // 读取字节数
    BytesWritten int64  // 写入字节数
    Error        error  // 错误信息（如果有）
}
```

### 4.2 主流程

```
解析参数
    ↓
检查 -l（列出编码）→ 直接输出支持的编码列表
    ↓
验证参数（-f 和 -t 的有效性）
    ↓
展开文件列表（处理通配符、递归目录）
    ↓
遍历处理每个文件
    ├─ 检测源编码（如果是 auto）
    ├─ 执行编码转换
    ├─ 输出结果（非静默模式）
    └─ 统计结果
    ↓
返回整体执行状态
```

### 4.3 编码检测逻辑

```go
// DetectEncoding 自动检测文件编码
//
// 参数:
//   - data: 文件内容样本（前 8KB）
//
// 返回值:
//   - string: 检测到的编码名称
//   - float64: 置信度（0-1）
//   - error: 检测错误（如果有）
func DetectEncoding(data []byte) (string, float64, error) {
    // 1. 检查 UTF-8 BOM
    if hasBOM(data, types.EncodingUTF8) {
        return types.EncodingUTF8BOM, 1.0, nil
    }

    // 2. 检查 UTF-16 BOM
    if hasBOM(data, types.EncodingUTF16LE) {
        return types.EncodingUTF16LE, 1.0, nil
    }
    if hasBOM(data, types.EncodingUTF16BE) {
        return types.EncodingUTF16BE, 1.0, nil
    }

    // 3. 检查是否为纯 ASCII（也是有效的 UTF-8）
    if isASCII(data) {
        return types.EncodingUTF8, 1.0, nil
    }

    // 4. 检查是否为有效的 UTF-8
    if utf8.Valid(data) {
        return types.EncodingUTF8, 0.95, nil
    }

    // 5. 使用启发式算法检测中文编码（GBK/Big5）
    // 统计字节分布特征
    encoding, confidence := detectChineseEncoding(data)
    if confidence > 0.8 {
        return encoding, confidence, nil
    }

    // 6. 默认返回 GBK（Windows 中文环境最常见）
    return types.EncodingGBK, 0.5, nil
}
```

### 4.4 编码转换逻辑

```go
// ConvertFile 转换单个文件编码
//
// 参数:
//   - srcPath: 源文件路径
//   - config: 转换配置
//
// 返回值:
//   - ConversionResult: 转换结果
//   - error: 执行错误（如果有）
func ConvertFile(srcPath string, config IconvConfig) (ConversionResult, error) {
    result := ConversionResult{FilePath: srcPath}

    // 1. 打开源文件
    srcFile, err := os.Open(srcPath)
    if err != nil {
        result.Error = err
        return result, err
    }
    defer srcFile.Close()

    // 2. 读取样本用于编码检测
    sample := make([]byte, types.GrepEncodingCheckSize)
    n, _ := srcFile.Read(sample)
    sample = sample[:n]
    srcFile.Seek(0, 0)

    // 3. 确定源编码
    fromEncoding := config.FromEncoding
    if fromEncoding == types.EncodingAuto {
        detected, confidence, _ := DetectEncoding(sample)
        fromEncoding = detected
        result.FromEncoding = fromEncoding

        // 目标编码为 none：只检测并打印
        if config.ToEncoding == types.EncodingNone {
            result.Success = true
            if !config.Quiet {
                printDetectionResult(srcPath, fromEncoding, confidence)
            }
            return result, nil
        }
    }
    result.FromEncoding = fromEncoding

    // 4. 如果源编码和目标编码相同，跳过
    if strings.EqualFold(fromEncoding, config.ToEncoding) {
        result.Success = true
        result.ToEncoding = config.ToEncoding
        if !config.Quiet {
            fmt.Printf("%s: %s (no conversion needed)\n", srcPath, fromEncoding)
        }
        return result, nil
    }

    // 5. 确定目标路径
    var dstPath string
    if config.Write {
        // 原地写入：使用临时文件，后续重命名
        dstPath = srcPath + ".tmp"
    } else {
        // 非原地写入
        if config.Output == "." {
            // 默认：当前目录下生成 原文件名.converted
            dstPath = srcPath + ".converted"
        } else {
            // 检查 -o 指定的是目录还是文件
            outputInfo, err := os.Stat(config.Output)
            if err == nil && outputInfo.IsDir() {
                // 是目录：目录下生成 原文件名.converted
                baseName := filepath.Base(srcPath)
                dstPath = filepath.Join(config.Output, baseName+".converted")
            } else {
                // 是文件路径（或不存在）：直接写入指定路径
                dstPath = config.Output
            }
        }
    }

    dstFile, err := os.Create(dstPath)
    if err != nil {
        result.Error = err
        return result, err
    }

    // 6. 执行转换
    bytesRead, bytesWritten, err := doConvert(srcFile, dstFile, fromEncoding, config.ToEncoding)
    dstFile.Close()

    result.BytesRead = bytesRead
    result.BytesWritten = bytesWritten

    if err != nil {
        os.Remove(dstPath)
        result.Error = err
        return result, err
    }

    // 7. 处理输出
    if config.Write {
        // 原地写入：备份原文件，重命名临时文件
        if config.Backup {
            backupPath := srcPath + ".bak"
            fs.MoveEx(srcPath, backupPath, true)
        } else {
            os.Remove(srcPath)
        }
        err = fs.MoveEx(dstPath, srcPath, true)
    }

    result.Success = err == nil
    result.ToEncoding = config.ToEncoding
    return result, err
}

// doConvert 执行实际的编码转换
func doConvert(src io.Reader, dst io.Writer, from, to string) (int64, int64, error) {
    // 获取解码器和编码器
    decoder := getEncoding(from).NewDecoder()
    encoder := getEncoding(to).NewEncoder()

    // 创建转换链：src → decoder → encoder → dst
    reader := transform.NewReader(src, decoder)
    writer := transform.NewWriter(dst, encoder)

    // 执行转换
    bytesRead, err := io.Copy(writer, reader)
    if err != nil {
        return bytesRead, 0, err
    }

    // 刷新编码器缓冲区
    err = writer.Close()

    // 获取实际写入字节数（需要额外统计）
    // 简化处理：返回读取字节数作为参考
    return bytesRead, bytesRead, err
}
```

### 4.5 编码映射

```go
import (
    "golang.org/x/text/encoding"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/encoding/traditionalchinese"
    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/encoding/korean"
    "golang.org/x/text/encoding/unicode"
)

// getEncoding 根据名称获取编码
func getEncoding(name string) encoding.Encoding {
    switch strings.ToUpper(name) {
    case types.EncodingUTF8:
        return encoding.Nop // UTF-8 无需转换
    case types.EncodingUTF8BOM:
        return unicode.UTF8BOM
    case types.EncodingUTF16:
        return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
    case types.EncodingUTF16LE:
        return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
    case types.EncodingUTF16BE:
        return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
    case types.EncodingGBK:
        return simplifiedchinese.GBK
    case types.EncodingGB2312:
        return simplifiedchinese.GB18030 // GB18030 兼容 GB2312
    case types.EncodingGB18030:
        return simplifiedchinese.GB18030
    case types.EncodingBig5:
        return traditionalchinese.Big5
    case types.EncodingLatin1:
        return charmap.ISO8859_1
    case types.EncodingShiftJIS:
        return japanese.ShiftJIS
    case types.EncodingEUCKR:
        return korean.EUCKR
    default:
        return encoding.Nop
    }
}

// printDetectionResult 打印编码检测结果
//
// 参数:
//   - path: 文件路径
//   - encoding: 检测到的编码
//   - confidence: 置信度
func printDetectionResult(path, encoding string, confidence float64) {
    if confidence >= 0.9 {
        fmt.Printf("%s: %s (high confidence)\n", path, encoding)
    } else if confidence >= 0.7 {
        fmt.Printf("%s: %s (medium confidence)\n", path, encoding)
    } else {
        fmt.Printf("%s: %s (low confidence)\n", path, encoding)
    }
}

## 5. CLI 层实现

### 5.1 命令定义

```go
// internal/cli/iconv.go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/qflag"
    "gitee.com/MM-Q/fck/internal/commands/iconv"
)

var IconvCmd *qflag.Cmd

var (
    iconvFrom         *qflag.StringFlag  // -f 源编码
    iconvTo           *qflag.StringFlag  // -t 目标编码
    iconvWrite        *qflag.BoolFlag    // -w 原地写入
    iconvBackup       *qflag.BoolFlag    // -b 创建备份
    iconvRecursive    *qflag.BoolFlag    // -r 递归处理
    iconvFollowSymlink *qflag.BoolFlag   // -L 跟随符号链接
    iconvList         *qflag.BoolFlag    // -l 列出支持的编码
    iconvVerbose      *qflag.BoolFlag    // -v 详细输出
    iconvQuiet        *qflag.BoolFlag    // -q 静默模式
)

func init() {
    IconvCmd = qflag.NewCmd("iconv", "", qflag.ExitOnError)

    // 使用 types 包中定义的编码列表
    iconvFrom = IconvCmd.Enum("from", "f", "源文件编码", types.EncodingAuto, types.SupportedEncodings)
    iconvTo = IconvCmd.Enum("to", "t", "目标编码", types.EncodingNone, types.TargetEncodings)
    iconvWrite = IconvCmd.Bool("write", "w", "原地写入（覆盖原文件）", false)
    iconvBackup = IconvCmd.Bool("backup", "b", "原地写入时创建备份（.bak）", false)
    iconvOutput = IconvCmd.String("output", "o", "输出目录（默认覆盖原文件或添加 .converted 后缀）", "")
    iconvList = IconvCmd.Bool("list", "l", "列出支持的编码格式", false)
    iconvVerbose = IconvCmd.Bool("verbose", "v", "显示详细转换信息", false)
    iconvQuiet = IconvCmd.Bool("quiet", "q", "静默模式", false)

    cmdOpts := &qflag.CmdOpts{
        Desc:        "文件编码转换工具",
        UseChinese:  true,
        UsageSyntax: fmt.Sprintf("%s iconv [options] <file...>", qflag.Root.Name()),
        Examples: map[string]string{
            "自动检测并转为 UTF-8": fmt.Sprintf("%s iconv -t UTF-8 gbk-file.txt", qflag.Root.Name()),
            "指定编码转换":        fmt.Sprintf("%s iconv -f GBK -t UTF-8 file.txt", qflag.Root.Name()),
            "检测文件编码":      fmt.Sprintf("%s iconv file.txt", qflag.Root.Name()),
            "检测多个文件编码":  fmt.Sprintf("%s iconv file1.txt file2.txt", qflag.Root.Name()),
            "指定编码转换":      fmt.Sprintf("%s iconv -f GBK -t UTF-8 file.txt", qflag.Root.Name()),
            "自动检测并转换":    fmt.Sprintf("%s iconv -t UTF-8 file.txt", qflag.Root.Name()),
            "原地转换":          fmt.Sprintf("%s iconv -t UTF-8 -w file.txt", qflag.Root.Name()),
            "原地转换并备份":    fmt.Sprintf("%s iconv -t UTF-8 -wb file.txt", qflag.Root.Name()),
            "批量转换":          fmt.Sprintf("%s iconv -t UTF-8 *.txt", qflag.Root.Name()),
            "输出到指定目录":    fmt.Sprintf("%s iconv -t UTF-8 -o output/ file.txt", qflag.Root.Name()),
            "输出到指定文件":    fmt.Sprintf("%s iconv -t UTF-8 -o out.txt file.txt", qflag.Root.Name()),
            "列出支持的编码":    fmt.Sprintf("%s iconv -l", qflag.Root.Name()),
        },
        Notes: []string{
            "-t 默认为 none（只检测不转换），指定目标编码时执行转换",
            "-f 默认为 auto（自动检测源编码）",
            "使用 -w 原地写入会直接覆盖原文件",
            "使用 -o 可指定输出目录或文件",
            "建议重要文件使用 -b 创建备份",
            "不支持递归目录处理，请使用通配符或管道",
        },
    }

    if err := IconvCmd.ApplyOpts(cmdOpts); err != nil {
        panic(fmt.Errorf("apply opts err: %w", err))
    }

    IconvCmd.SetRun(runIconv)
}

// runIconv 执行 iconv 命令
func runIconv(cmd qflag.Command) error {
    // 列出支持的编码
    if iconvList.Get() {
        return iconv.ListEncodings()
    }

    // 检查文件参数
    files := cmd.Args()
    if len(files) == 0 {
        return fmt.Errorf("no input files specified")
    }

    config := iconv.IconvConfig{
        Files:        files,
        FromEncoding: iconvFrom.Get(),
        ToEncoding:   iconvTo.Get(),
        Write:        iconvWrite.Get(),
        Backup:       iconvBackup.Get(),
        Output:       iconvOutput.Get(),
        Verbose:      iconvVerbose.Get(),
        Quiet:        iconvQuiet.Get(),
    }

    return iconv.IconvCmdMain(config)
}
```

## 6. 错误处理

### 6.1 可能的错误

| 错误类型 | 处理方式 |
|----------|----------|
| 文件不存在 | 输出错误信息，继续处理其他文件 |
| 权限不足 | 输出错误信息，继续处理其他文件 |
| 编码不支持 | 输出错误信息，跳过该文件 |
| 转换失败 | 输出错误信息，删除临时文件 |
| 磁盘空间不足 | 输出错误信息，清理临时文件 |

### 6.2 错误输出格式

```
iconv: file.txt: file not found
iconv: file.txt: unsupported encoding: XXX
iconv: file.txt: conversion failed: invalid byte sequence
```

## 7. 输出格式

### 7.1 标准输出（非静默模式）

```bash
# 单文件转换
$ fck iconv -f GBK -t UTF-8 file.txt
file.txt: GBK → UTF-8 (1024 bytes)

# 多文件转换
$ fck iconv -f GBK -t UTF-8 *.txt
file1.txt: GBK → UTF-8 (1024 bytes)
file2.txt: GBK → UTF-8 (2048 bytes)
file3.txt: skipped (already UTF-8)

# 详细模式
$ fck iconv -f GBK -t UTF-8 -v file.txt
file.txt:
  Source encoding: GBK (auto-detected, confidence: 0.95)
  Target encoding: UTF-8
  Bytes read: 1024
  Bytes written: 1024
  Status: success
```

### 7.2 统计信息

```bash
$ fck iconv -f GBK -t UTF-8 -wr dir/
...

Total: 10 files
Success: 8 files
Skipped: 1 file (already UTF-8)
Failed: 1 file
```

## 8. 依赖

```go
// go.mod
go get golang.org/x/text
```

## 9. 实现步骤

1. 创建目录结构
2. 实现 `commands/iconv/detector.go` - 编码检测
3. 实现 `commands/iconv/converter.go` - 编码转换
4. 实现 `commands/iconv/cmd_iconv.go` - 主逻辑
5. 实现 `cli/iconv.go` - CLI 层
6. 注册命令到 root
7. 测试各种编码转换场景
