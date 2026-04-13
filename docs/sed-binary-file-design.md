# Sed 二进制文件检测设计方案

## 目标

为 sed 命令添加二进制文件检测功能，默认跳过二进制文件，除非用户显式启用处理。

## 参考 grep 的实现

grep 的二进制文件处理策略：
- 检测：读取前 8000 字节，检查是否包含空字符 `\0`
- 默认：输出提示 "bin file matches" 并跳过
- `-a/--text`：强制作为文本处理
- `-I/--ignore-binary`：静默跳过，不输出提示

## 设计方案

### 1. 新增 CLI 标志 (internal/cli/sed.go)

添加两个布尔标志，与 grep 保持一致：

```go
// 二进制文件处理
sedText         *qflag.BoolFlag // -a, --text 强制将二进制文件视为文本处理
sedIgnoreBinary *qflag.BoolFlag // -I 忽略二进制文件（不输出提示）
```

在 `init()` 中定义：

```go
// 二进制文件处理标志
sedText = SedCmd.Bool("text", "a", "强制将二进制文件视为文本处理", false)
sedIgnoreBinary = SedCmd.Bool("ignore-binary", "I", "完全忽略二进制文件，不输出提示", false)
```

### 2. 修改 SedConfig (internal/commands/sed/cmd_sed.go)

在配置结构体中添加二进制文件处理字段：

```go
type SedConfig struct {
    // ... 原有字段 ...
    
    // 二进制文件处理
    Text         bool // -a, --text 强制将二进制文件视为文本处理
    IgnoreBinary bool // -I 忽略二进制文件（不输出提示，静默跳过）
}
```

### 3. 二进制文件检测与处理

在 `processFile` 函数中，打开文件后添加检测逻辑：

```go
func processFile(config *SedConfig) error {
    // 1. 验证参数
    // 2. 解析行号范围
    // 3. 编译正则
    // 4. 重置计数
    
    // 5. 打开文件并检测二进制
    file, err := os.Open(config.Target)
    if err != nil {
        return fmt.Errorf("cannot open file: %w", err)
    }
    defer func() { _ = file.Close() }()
    
    // 检测二进制文件
    isBinary, err := utils.IsBinaryFile(file)
    if err != nil {
        return fmt.Errorf("failed to detect file type: %w", err)
    }
    
    if isBinary {
        if config.Text {
            // -a/--text 模式：强制作为文本处理，继续执行
        } else if config.IgnoreBinary {
            // -I 模式：静默跳过
            return nil
        } else {
            // 默认行为：输出提示并跳过
            fmt.Printf("bin file %s matches\n", config.Target)
            return nil
        }
    }
    
    // 6. 处理文件（预览或原地修改）
    // ...
}
```

### 4. 原地修改模式的特殊处理

原地修改模式 (`processFileInPlace`) 中，二进制文件检测需要在创建备份之前进行：

```go
func processFileInPlace(config *SedConfig) (err error) {
    // 打开源文件
    sourceFile, err := os.Open(config.Target)
    if err != nil {
        return fmt.Errorf("cannot open file: %w", err)
    }
    
    // 检测二进制文件
    isBinary, err := utils.IsBinaryFile(sourceFile)
    if err != nil {
        _ = sourceFile.Close()
        return fmt.Errorf("failed to detect file type: %w", err)
    }
    
    if isBinary {
        _ = sourceFile.Close()
        if config.Text {
            // 重新打开文件继续处理
            sourceFile, err = os.Open(config.Target)
            if err != nil {
                return fmt.Errorf("cannot open file: %w", err)
            }
        } else if config.IgnoreBinary {
            // 静默跳过
            return nil
        } else {
            // 默认行为：输出提示并跳过
            fmt.Printf("bin file %s matches\n", config.Target)
            return nil
        }
    }
    
    defer func() { _ = sourceFile.Close() }()
    
    // 创建备份（如果需要）
    if config.Backup {
        // ...
    }
    
    // ... 后续处理
}
```

### 5. 多文件处理

在 CLI 层的多文件循环中，每个文件独立进行二进制检测：

```go
for _, target := range targets {
    config := baseConfig
    config.Target = target
    
    // 多文件时打印文件名（预览模式）
    if len(targets) > 1 && !config.InPlace {
        fmt.Printf("==> %s <==\n", target)
    }
    
    if err := sed.SedCmdMain(config); err != nil {
        fmt.Fprintf(os.Stderr, "error processing %s: %v\n", target, err)
        hasError = true
    }
    
    // 多文件预览时添加空行分隔
    if len(targets) > 1 && !config.InPlace {
        fmt.Println()
    }
}
```

## 实现步骤

1. **修改 CLI 层** (`internal/cli/sed.go`)
   - 添加 `sedText` 和 `sedIgnoreBinary` 标志
   - 在 `baseConfig` 中传递这两个配置

2. **修改 Core 层配置** (`internal/commands/sed/cmd_sed.go`)
   - 在 `SedConfig` 中添加 `Text` 和 `IgnoreBinary` 字段

3. **添加二进制文件检测** (`internal/commands/sed/cmd_sed.go`)
   - 在 `processFile` 中添加检测逻辑
   - 在 `processFileInPlace` 中添加检测逻辑
   - 使用 `utils.IsBinaryFile()` 进行检测

4. **测试验证**
   - 测试二进制文件默认跳过
   - 测试 `-a` 强制处理
   - 测试 `-I` 静默跳过
   - 测试多文件处理

## 代码变更范围

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/cli/sed.go` | 修改 | 添加二进制文件处理标志 |
| `internal/commands/sed/cmd_sed.go` | 修改 | 添加配置字段和检测逻辑 |

## 示例用法

```bash
# 默认行为：检测到二进制文件输出提示并跳过
fck sed -p "old" -r "new" binary_file.exe
# 输出: Binary file binary_file.exe matches

# -a/--text 强制作为文本处理
fck sed -p "old" -r "new" -a binary_file.exe

# -I 静默跳过二进制文件
fck sed -p "old" -r "new" -I binary_file.exe
# 无输出，直接跳过

# 多文件处理
fck sed -p "old" -r "new" file1.txt binary_file.exe file2.txt
# file1.txt: 正常处理
# binary_file.exe: 输出提示并跳过
# file2.txt: 正常处理
```

## 注意事项

1. **与 grep 保持一致**：使用相同的标志 `-a` 和 `-I`，相同的行为
2. **原地修改模式**：检测在创建备份之前进行，避免不必要的备份
3. **多文件处理**：每个文件独立检测，互不影响
4. **错误处理**：检测失败时返回错误，由调用方决定如何处理
