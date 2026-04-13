# Sed 管道输入支持设计方案

## 目标

为 sed 命令添加管道输入支持，允许从标准输入读取内容进行替换处理，类似于 Linux sed 的行为。

## Linux sed 管道行为参考

```bash
# 基本用法：从管道读取
echo "hello world" | sed 's/hello/hi/'

# 与文件处理的区别：
# - 管道输入不支持原地修改 (-i)
# - 输出直接到 stdout
# - 没有文件名信息
```

## 当前代码结构分析

### CLI 层 (`internal/cli/sed.go`)
- 当前要求至少一个文件参数 `len(args) < 1`
- 通过 `filepath.Glob` 展开通配符
- 多文件循环处理

### Core 层 (`internal/commands/sed/cmd_sed.go`)
- `SedCmdMain`: 主入口，验证配置后调用 `processFile`
- `validateConfig`: 验证文件存在性
- `processFile`: 根据 `InPlace` 选择预览或原地修改模式
- `processFilePreview`: 从文件读取，处理后输出到 stdout
- `processFileInPlace`: 从文件读取，处理后写回文件

## 设计方案

### 1. 检测管道输入

使用 `utils.IsStdinPipe()` 检测 stdin 是否为管道/重定向：

```go
// 在 CLI 层检测
isPipe := utils.IsStdinPipe()
```

### 2. 修改 CLI 层 (`internal/cli/sed.go`)

#### 2.1 文件参数改为可选

```go
func runSed(cmd qflag.Command) error {
    args := cmd.Args()
    
    // 检测是否为管道输入
    isPipe := utils.IsStdinPipe()
    
    // 检查参数：非管道模式下需要文件参数
    if len(args) < 1 && !isPipe {
        return fmt.Errorf("no target file specified")
    }
    
    // ... 后续处理
}
```

#### 2.2 管道输入的处理流程

```go
// 管道输入模式
if isPipe {
    // 检查冲突标志
    if sedInPlace.Get() {
        return fmt.Errorf("cannot use -i flag when reading from stdin")
    }
    
    // 准备配置
    config := sed.SedConfig{
        Target:       "-",  // 使用 "-" 表示 stdin
        Pattern:      sedPattern.Get(),
        Replacement:  sedReplacement.Get(),
        Regexp:       sedRegexp.Get(),
        LineRange:    sedLineRange.Get(),
        InPlace:      false,  // 管道模式强制预览
        Backup:       false,
        MaxCount:     sedMaxCount.Get(),
        IgnoreCase:   sedIgnoreCase.Get(),
        MaxBuffer:    maxBuffer,
        Text:         sedText.Get(),
        IgnoreBinary: sedIgnoreBinary.Get(),
    }
    
    return sed.SedCmdMain(config)
}

// 文件模式（原有逻辑）
// ... 展开通配符、多文件循环等
```

### 3. 修改 Core 层 (`internal/commands/sed/cmd_sed.go`)

#### 3.1 修改 `validateConfig`

支持 `"-"` 作为 stdin 的标识：

```go
func validateConfig(config *SedConfig) error {
    // 检查模式参数
    if config.Pattern == "" {
        return fmt.Errorf("please use -p to specify the pattern")
    }
    if config.Replacement == "" {
        return fmt.Errorf("please use -r to specify the replacement")
    }
    
    // stdin 模式不需要检查文件存在性
    if config.Target == "-" {
        // 管道模式不支持原地修改
        if config.InPlace {
            return fmt.Errorf("cannot use -i flag when reading from stdin")
        }
        return nil
    }
    
    // 原有文件存在性检查
    if config.Target == "" {
        return fmt.Errorf("no target file specified")
    }
    if _, err := os.Stat(config.Target); err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("file not found: %s", config.Target)
        }
        return fmt.Errorf("cannot access file: %w", err)
    }
    
    return nil
}
```

#### 3.2 修改 `processFile`

根据 Target 是否为 `"-"` 选择处理方式：

```go
func processFile(config *SedConfig) error {
    // stdin 模式
    if config.Target == "-" {
        return processStdin(config)
    }
    
    // 原地修改模式
    if config.InPlace {
        return processFileInPlace(config)
    }
    
    // 预览模式
    return processFilePreview(config)
}
```

#### 3.3 新增 `processStdin` 函数

```go
// processStdin 处理标准输入
// 从 stdin 读取，处理后输出到 stdout
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processStdin(config *SedConfig) error {
    // stdin 不检测二进制（管道输入视为文本）
    
    // 设置动态缓冲区
    scanner := bufio.NewScanner(os.Stdin)
    buf := make([]byte, types.InitialBufferSize)
    scanner.Buffer(buf, int(config.MaxBuffer))
    
    // 遍历输入行，边读边处理边输出
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        processedLine, _ := processLine(line, config, lineNum)
        fmt.Println(processedLine)
    }
    
    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading stdin: %w", err)
    }
    
    return nil
}
```

### 4. 使用示例

```bash
# 基本管道替换
echo "hello world" | fck sed -p "hello" -r "hi"
# 输出: hi world

# 多行管道输入
cat file.txt | fck sed -p "old" -r "new"

# 结合其他命令
ls -la | fck sed -p "total" -r "TOTAL"

# 管道输入不支持 -i（会报错）
echo "test" | fck sed -p "test" -r "ok" -i
# 错误: cannot use -i flag when reading from stdin
```

## 实现步骤

1. **修改 CLI 层** (`internal/cli/sed.go`)
   - 导入 `gitee.com/MM-Q/fck/internal/utils`
   - 修改 `runSed` 函数，添加管道检测和分支处理
   - 文件参数改为可选（仅在非管道模式下需要）

2. **修改 Core 层** (`internal/commands/sed/cmd_sed.go`)
   - 修改 `validateConfig`，支持 `"-"` 作为 stdin 标识
   - 修改 `processFile`，添加 stdin 分支
   - 新增 `processStdin` 函数

3. **测试验证**
   - 测试基本管道输入
   - 测试多行管道输入
   - 测试 `-i` 冲突检测
   - 测试与文件模式的共存

## 代码变更范围

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `internal/cli/sed.go` | 修改 | 添加管道检测和 stdin 处理分支 |
| `internal/commands/sed/cmd_sed.go` | 修改 | 支持 `"-"` 标识，添加 `processStdin` |

## 注意事项

1. **与 `-i` 互斥**：管道输入不支持原地修改，需要检测并报错
2. **不检测二进制**：管道输入默认视为文本，不进行二进制检测
3. **错误处理**：stdin 读取错误需要单独处理
4. **与文件模式共存**：CLI 层优先处理管道，然后是文件
