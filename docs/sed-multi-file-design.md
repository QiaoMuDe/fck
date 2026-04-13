# Sed 多文件处理设计方案

## 当前现状

当前 `sed` 命令只支持处理单个文件：
```bash
fck sed -p "old" -r "new" file.txt
```

## 目标

支持处理多个文件：
```bash
# 多个文件
fck sed -p "old" -r "new" file1.txt file2.txt file3.txt

# 通配符（由 shell 展开）
fck sed -p "old" -r "new" *.txt

# 混合模式
fck sed -p "old" -r "new" -i -b file1.txt file2.txt
```

## 设计方案

### 1. CLI 层修改 (internal/cli/sed.go)

**修改点**：
- 接受多个文件参数（`args` 切片）
- 遍历处理每个文件
- 错误处理策略：继续处理还是立即退出

```go
func runSed(cmd qflag.Command) error {
    args := cmd.Args()
    
    // 检查至少有一个文件参数
    if len(args) < 1 {
        return fmt.Errorf("no target file specified")
    }
    
    // 检查缓冲区大小
    maxBuffer := sedMaxBuffer.Get()
    if maxBuffer <= types.InitialBufferSize {
        return fmt.Errorf("buffer size must be greater than %d bytes", types.InitialBufferSize)
    }
    
    // 准备基础配置（所有文件共用）
    baseConfig := sed.SedConfig{
        Pattern:     sedPattern.Get(),
        Replacement: sedReplacement.Get(),
        Regexp:      sedRegexp.Get(),
        LineRange:   sedLineRange.Get(),
        InPlace:     sedInPlace.Get(),
        Backup:      sedBackup.Get(),
        MaxCount:    sedMaxCount.Get(),
        IgnoreCase:  sedIgnoreCase.Get(),
        MaxBuffer:   maxBuffer,
    }
    
    // 处理多个文件
    var hasError bool
    for _, target := range args {
        config := baseConfig
        config.Target = target
        
        // 打印文件名（多文件时）
        if len(args) > 1 {
            fmt.Fprintf(os.Stderr, "Processing: %s\n", target)
        }
        
        if err := sed.SedCmdMain(config); err != nil {
            fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", target, err)
            hasError = true
            // 继续处理下一个文件
        }
    }
    
    if hasError {
        return fmt.Errorf("some files failed to process")
    }
    return nil
}
```

### 2. Core 层修改 (internal/commands/sed/cmd_sed.go)

**修改点**：
- `SedCmdMain` 接收单个文件配置
- 每个文件独立处理，独立计数
- 重置 `replaceCount` 为每个文件单独计数

**关键修改**：

```go
// SedCmdMain 执行 sed 命令（单文件）
// 多文件处理由 CLI 层循环调用
func SedCmdMain(config SedConfig) error {
    // 1. 验证参数
    if err := validateConfig(&config); err != nil {
        return err
    }

    // 2. 解析行号范围
    if config.LineRange != "" {
        if err := parseLineRange(&config); err != nil {
            return err
        }
    }

    // 3. 编译正则（如果需要）
    if config.Regexp {
        if err := compilePattern(&config); err != nil {
            return err
        }
    }

    // 4. 处理文件（每个文件独立，replaceCount 从 0 开始）
    config.replaceCount = 0
    return processFile(&config)
}
```

### 3. 输出格式设计

#### 3.1 预览模式（无 -i）

**单文件**：
```bash
$ fck sed -p "a" -r "X" file.txt
X b c
```

**多文件**：
```bash
$ fck sed -p "a" -r "X" file1.txt file2.txt
==> file1.txt <==
X b c

==> file2.txt <==
X y z
```

#### 3.2 原地修改模式（-i）

**多文件**：
```bash
$ fck sed -p "a" -r "X" -i file1.txt file2.txt
# 无输出（静默处理）
# 或使用 -v 显示处理信息
$ fck sed -p "a" -r "X" -i -v file1.txt file2.txt
processed: file1.txt
processed: file2.txt
```

### 4. 错误处理策略

| 策略 | 行为 | 适用场景 |
|------|------|----------|
| 继续处理 | 出错打印错误，继续下一个文件 | 默认行为 |
| 立即退出 | 第一个错误就退出 | 严格模式（可添加 --strict 标志）|

**默认行为**：
```bash
$ fck sed -p "a" -r "X" file1.txt not_exist.txt file2.txt
Error processing not_exist.txt: file not found: not_exist.txt
# file1.txt 和 file2.txt 已处理
```

### 5. 边界情况处理

#### 5.1 文件不存在
- 打印错误到 stderr
- 继续处理其他文件
- 最后返回错误码

#### 5.2 目录而非文件
- 跳过目录
- 打印警告

#### 5.3 无权限文件
- 打印权限错误
- 继续处理其他文件

#### 5.4 通配符无匹配
- shell 不会展开，直接作为文件名处理
- 显示文件不存在错误

### 6. 实现步骤

1. **修改 CLI 层** (`internal/cli/sed.go`)
   - 循环处理多个文件参数
   - 添加文件分隔输出

2. **修改 Core 层** (`internal/commands/sed/cmd_sed.go`)
   - 确保每个文件独立计数
   - 添加多文件输出格式支持

3. **测试验证**
   - 单文件（向后兼容）
   - 多文件预览模式
   - 多文件原地修改模式
   - 错误处理

### 7. 代码变更范围

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/cli/sed.go` | 修改 | 循环处理多文件 |
| `internal/commands/sed/cmd_sed.go` | 修改 | 支持多文件输出格式 |

### 8. 示例用法

```bash
# 基础多文件
fck sed -p "old" -r "new" a.txt b.txt c.txt

# 原地修改多文件
fck sed -p "old" -r "new" -i a.txt b.txt c.txt

# 带备份的多文件
fck sed -p "old" -r "new" -i -b a.txt b.txt c.txt

# 通配符（shell 展开）
fck sed -p "old" -r "new" -i *.go

# 混合文件和通配符
fck sed -p "old" -r "new" -i a.txt *.go b.txt
```

## 注意事项

1. **向后兼容**：单文件用法完全保持兼容
2. **计数独立**：每个文件的替换计数独立计算
3. **错误聚合**：多文件时收集所有错误，最后统一返回
4. **输出清晰**：多文件时添加文件分隔标识
