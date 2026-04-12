# grep 二进制文件跳过功能设计方案

## 1. 核心思路

参考 Linux grep 的行为：
- **默认**：检测到二进制文件时跳过，输出提示 "Binary file xxx matches"
- **-a/--text**：强制将二进制文件视为文本处理
- **-I**：完全忽略二进制文件（不输出任何提示）

## 2. 数据结构修改

```go
// GrepConfig 新增字段
type GrepConfig struct {
    // ... 现有字段 ...

    // 二进制文件处理
    Text         bool // -a, --text 强制将二进制文件视为文本处理
    IgnoreBinary bool // -I 忽略二进制文件（不输出提示，静默跳过）
}
```

## 3. 二进制检测实现

```go
// isBinaryFile 检测文件是否为二进制文件
//
// 原理：读取文件前 8000 字节，检查是否包含空字符(\0)
// 这是 Linux grep 使用的检测方法
//
// 注意：
//   - 只支持普通文件，stdin/pipe 默认返回 false（视为文本）
//   - 检测后会重置文件指针到开头
//
// 参数:
//   - file: 已打开的文件句柄
//
// 返回:
//   - bool: true 表示二进制文件，false 表示文本文件或无法检测
//   - error: 读取或重置指针错误
func isBinaryFile(file *os.File) (bool, error) {
    // 获取文件信息，检查是否为普通文件
    info, err := file.Stat()
    if err != nil {
        return false, err
    }

    // 非普通文件（stdin、pipe、设备文件等）默认视为文本
    // 这些文件不支持 Seek 操作
    if !info.Mode().IsRegular() {
        return false, nil
    }

    // 读取文件前 8000 字节
    buf := make([]byte, 8000)
    n, err := file.Read(buf)
    if err != nil && err != io.EOF {
        return false, err
    }

    // 检查是否包含空字符（二进制文件特征）
    isBinary := bytes.Contains(buf[:n], []byte{0})

    // 重置文件指针到开头，供后续读取使用
    if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
        return false, seekErr
    }

    return isBinary, nil
}
```

## 4. 处理流程修改

### 4.1 单文件处理 (processFile)

```go
func processFile(file *os.File, config *GrepConfig) error {
    // 1. 检测是否为二进制文件
    isBinary, err := isBinaryFile(file)
    if err != nil {
        if config.NoMessages {
            return nil
        }
        return fmt.Errorf("检测文件类型失败: %w", err)
    }
    
    // 2. 根据配置处理二进制文件
    if isBinary {
        if config.Text {
            // -a/--text 模式：强制作为文本处理，继续执行
        } else if config.IgnoreBinary {
            // -I 模式：完全忽略，不输出任何提示
            return nil
        } else {
            // 默认行为：输出提示并跳过
            if config.WithFilename || config.filename != "" {
                fmt.Printf("Binary file %s matches\n", config.filename)
            } else {
                fmt.Println("Binary file matches")
            }
            return nil
        }
    }
    
    // 3. 继续原有处理逻辑...
    return processLines(file, config)
}
```

### 4.2 递归模式处理 (processSingleFile)

```go
func processSingleFile(path string, config *GrepConfig) error {
    // ... 原有代码 ...
    
    file, err := os.Open(path)
    if err != nil {
        if config.NoMessages {
            return nil
        }
        return nil
    }
    defer func() { _ = file.Close() }()
    
    config.filename = path
    return processFile(file, config) // 复用上面的逻辑
}
```

## 5. CLI 标志定义

```go
// internal/cli/grep.go

var (
    // ... 现有标志 ...

    grepText         *qflag.BoolFlag // -a, --text 强制将二进制文件视为文本处理
    grepIgnoreBinary *qflag.BoolFlag // -I 忽略二进制文件（不输出提示）
)

func init() {
    // ... 原有代码 ...

    // 二进制文件处理标志
    grepText = GrepCmd.Bool("text", "a",
        "强制将二进制文件视为文本处理", false)
    grepIgnoreBinary = GrepCmd.Bool("ignore-binary", "I",
        "完全忽略二进制文件，不输出提示", false)

    // 更新帮助文档
    cmdOpts.Notes = append(cmdOpts.Notes,
        "默认跳过二进制文件并输出提示，使用 -a 强制处理，使用 -I 静默跳过",
    )
}

func runGrep(cmd qflag.Command) error {
    // ...
    config := grep.GrepConfig{
        // ... 现有字段 ...
        Text:         grepText.Get(),
        IgnoreBinary: grepIgnoreBinary.Get(),
    }
    // ...
}
```

## 6. 文件修改清单

| 文件 | 修改内容 |
|------|----------|
| `internal/commands/grep/cmd_grep.go` | 添加 Text/IgnoreBinary 字段，添加 isBinaryFile 函数，修改 processFile |
| `internal/commands/grep/recursive.go` | 修改 processSingleFile 复用新逻辑 |
| `internal/cli/grep.go` | 添加 -a 和 -I 两个标志，更新帮助文档 |

## 7. 使用示例

```bash
# 默认行为：检测到二进制文件输出提示并跳过
fck grep "pattern" file.bin
# 输出: Binary file file.bin matches

# -a 强制作为文本处理（可能输出乱码）
fck grep -a "pattern" file.bin

# -I 完全静默跳过二进制文件
fck grep -I "pattern" .
```

## 8. 边缘情况处理

1. **空文件**：不会包含 \0，视为文本文件处理
2. **UTF-16 编码**：包含大量 \0，会被误判为二进制（与 Linux grep 行为一致）
3. **大文件**：只读取前 8000 字节，性能开销小
4. **管道输入**：从 stdin 读取时无法检测（不支持 Seek），默认视为文本处理
5. **设备文件**：/dev/null、/dev/random 等默认视为文本处理
6. **符号链接**：跟随链接后检测目标文件类型

## 9. 优势

1. **兼容性**：与 Linux grep 行为一致
2. **灵活性**：提供多种处理方式满足不同需求
3. **性能**：只读取文件头部，开销极小
4. **用户体验**：避免终端显示二进制乱码
