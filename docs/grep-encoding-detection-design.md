# Grep 命令非 UTF-8 编码文件检测方案

## 1. 功能概述

为 `fck grep` 命令增加非 UTF-8 编码文件检测功能，当遇到 GBK 或其他非 UTF-8 编码文件时，打印警告提示，避免中文搜索时出现乱码问题。

## 2. 检测策略

### 2.1 检测时机

在读取文件内容后、执行搜索前进行检测。

### 2.2 检测方法

使用 Go 标准库 `unicode/utf8` 的 `Valid` 函数检测：

```go
import "unicode/utf8"

// 检测数据是否为有效 UTF-8
if !utf8.Valid(data) {
    // 非 UTF-8 编码
}
```

### 2.3 处理逻辑

```
读取文件内容
    ↓
检测是否为有效 UTF-8
    ├─ 是 → 正常执行搜索
    └─ 否 → 打印警告提示，跳过该文件
```

## 3. 命令行参数复用

### 复用现有标志

| 短标志 | 长标志 | 作用 |
|--------|--------|------|
| -a | --text | 强制将非 UTF-8 文件视为文本处理 |
| -I | --ignore-binary | 静默跳过非 UTF-8 文件（不输出提示） |

### 标志组合行为

| 场景 | 命令 | 行为 |
|------|------|------|
| 默认 | `grep pattern file` | 检测到非 UTF-8 时输出警告提示 |
| 强制处理 | `grep -a pattern file` | 强制处理非 UTF-8 文件（可能显示乱码） |
| 静默跳过 | `grep -I pattern file` | 静默跳过非 UTF-8 文件 |
| 强制+静默 | `grep -a -I pattern file` | 强制处理且静默（无提示） |

## 4. 代码结构修改

### 4.1 cli/grep.go

无需修改，复用现有的 `--text/-a` 标志。

### 4.2 commands/grep/cmd_grep.go

#### 新增导入

```go
import (
    // ... 现有导入
    "unicode/utf8"
)
```

#### 修改 GrepConfig 结构体（如有需要）

```go
// 复用现有 Text 和 IgnoreBinary 字段
// Text: 强制将非 UTF-8/二进制文件视为文本处理
// IgnoreBinary: 静默跳过非 UTF-8/二进制文件
```

#### 新增检测函数

```go
// isValidUTF8 检测数据是否为有效 UTF-8
// 为提高性能，只检测前 GrepEncodingCheckSize 字节（默认 8KB）
//
// 参数:
//   - data: 文件内容
//
// 返回值:
//   - bool: 是否为有效 UTF-8
func isValidUTF8(data []byte) bool {
    sampleSize := types.GrepEncodingCheckSize
    if len(data) < sampleSize {
        return utf8.Valid(data)
    }
    return utf8.Valid(data[:sampleSize])
}

// printEncodingWarning 打印编码警告
//
// 参数:
//   - filename: 文件名
//   - noMessages: 是否静默模式（-I 标志）
func printEncodingWarning(filename string, noMessages bool) {
    if noMessages {
        return
    }
    fmt.Fprintf(os.Stderr, "grep: %s: file appears to be non-UTF-8 encoded\n", filename)
    fmt.Fprintf(os.Stderr, "      (consider converting to UTF-8 with: iconv -f GBK -t UTF-8 %s)\n", filename)
}
```

#### 修改文件处理逻辑

在 `processFile` 或 `searchInFile` 函数中：

```go
func processFile(filename string, config GrepConfig) error {
    // 读取文件内容
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }

    // 检测是否为有效 UTF-8
    if !isValidUTF8(data) {
        // 如果启用了 -I 标志，静默跳过
        if config.IgnoreBinary {
            return nil
        }
        // 如果启用了 -a 标志，强制处理
        if config.Text {
            // 尝试将非 UTF-8 字符替换为 � 后继续搜索
            data = []byte(string(data)) // 这会强制转换，乱码字符变为 �
        } else {
            // 默认：打印警告并跳过
            printEncodingWarning(filename, config.NoMessages)
            return nil // 返回 nil 表示正常跳过，不中断其他文件处理
        }
    }

    // 继续执行搜索...
    return searchInContent(data, config)
}
```

### 4.3 处理非 UTF-8 内容的方式

当 `--text/-a` 启用时，将非 UTF-8 字符替换为替换字符（�）：

```go
// toValidUTF8 将非 UTF-8 数据转换为有效 UTF-8（替换非法字节）
func toValidUTF8(data []byte) []byte {
    // 先尝试直接转换为字符串（Go 会自动处理）
    str := string(data)
    // 再转回字节切片
    return []byte(str)
}
```

Go 的 `string` 转换会自动将非法 UTF-8 序列替换为 Unicode 替换字符（U+FFFD，即 �）。

## 5. 输出示例

### 正常情况（UTF-8 文件）

```bash
$ fck grep "中文" utf8-file.txt
1:这是一行中文文本
```

### 检测到非 UTF-8 文件（默认行为）

```bash
$ fck grep "中文" gbk-file.txt
grep: gbk-file.txt: file appears to be non-UTF-8 encoded
      (consider converting to UTF-8 with: iconv -f GBK -t UTF-8 gbk-file.txt)
```

### 强制处理非 UTF-8 文件

```bash
$ fck grep -a "中文" gbk-file.txt
1:������������  # 乱码，但会显示
```

### 静默跳过非 UTF-8 文件

```bash
$ fck grep -I "中文" gbk-file.txt
$  # 无任何输出（静默跳过）
```

## 6. 与现有功能的整合

### 与二进制文件检测的关系

当前代码可能已经检测二进制文件，需要确保编码检测在二进制检测之后：

```go
func processFile(filename string, config GrepConfig) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }

    // 1. 检测是否为二进制文件
    if isBinary(data) {
        if config.IgnoreBinary {
            return nil // 静默跳过
        }
        if !config.Text {
            fmt.Fprintf(os.Stderr, "Binary file %s matches\n", filename)
            return nil
        }
    }

    // 2. 检测是否为有效 UTF-8（如果不是二进制）
    if !isBinary(data) && !isValidUTF8(data) {
        if !config.Text {
            printEncodingWarning(filename, config.NoMessages)
            return nil
        }
        // 强制处理：转换为有效 UTF-8
        data = toValidUTF8(data)
    }

    // 执行搜索...
}
```

## 7. 边界情况处理

### 7.1 部分非 UTF-8 字节

某些文件可能大部分是 UTF-8，但包含少量非 UTF-8 字节（如损坏的文件）。`utf8.Valid` 会返回 `false`，这是预期行为。

### 7.2 UTF-8 BOM

UTF-8 BOM（字节顺序标记）是合法的 UTF-8 序列，`utf8.Valid` 会返回 `true`，无需特殊处理。

### 7.3 大文件性能优化

对于大文件，只检测前 `GrepEncodingCheckSize`（默认 8KB）数据，平衡准确性和性能：

```go
func isValidUTF8(data []byte) bool {
    sampleSize := types.GrepEncodingCheckSize
    if len(data) < sampleSize {
        return utf8.Valid(data)
    }
    return utf8.Valid(data[:sampleSize])
}
```

**为什么选择 8KB：**
- 能覆盖文件头信息和前几页内容
- 大多数编码特征（如 GBK 中文）在前 8KB 就能体现
- 性能影响极小（8KB 检测非常快）
- 误判概率较低

**常量定义位置：** `internal/types/types.go` 中的 `GrepEncodingCheckSize`，方便统一修改

## 8. 实现步骤

1. **添加 UTF-8 检测函数** - 在 `cmd_grep.go` 中添加 `isValidUTF8`
2. **修改文件处理逻辑** - 在读取文件后添加编码检测
3. **添加警告输出函数** - 实现 `printEncodingWarning`
4. **整合 `--text` 标志** - 支持强制处理非 UTF-8 文件
5. **测试** - 使用 GBK、UTF-8、二进制文件进行测试

## 9. 注意事项

1. **性能影响** - `utf8.Valid` 需要遍历整个文件，对大文件可能有性能影响
2. **误判可能** - 某些 GBK 文件可能恰好是有效的 UTF-8（罕见），会正常搜索但结果可能不对
3. **用户教育** - 提示信息应指导用户如何转换编码
