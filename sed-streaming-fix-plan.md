# Sed 命令大文件处理问题修复方案

> **修复目标**：解决 sed 命令处理大文件时的内存溢出问题  
> **方案选择**：方案 A - 流式处理  
> **核心思想**：边读边处理边输出，不缓存整个文件内容  
> **预期效果**：内存占用从 O(n) 降至 O(1)，可处理任意大小文件

---

## 一、问题现状

### 1.1 当前代码问题

```go
// internal/commands/sed/cmd_sed.go
func processFile(config *SedConfig) error {
    // 打开文件...
    
    var resultLines []string  // ⚠️ 问题：缓存所有行
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        processedLine, _ := processLine(line, config, lineNum)
        resultLines = append(resultLines, processedLine)  // ⚠️ 内存持续增长
    }
    
    // 处理完成后才输出或写入
    if config.InPlace {
        return writeFile(config, resultLines)  // 第二次遍历
    }
    // 输出到标准输出...
}
```

### 1.2 内存占用分析

| 文件大小 | 当前内存占用 | 修复后内存占用 |
|----------|-------------|---------------|
| 1 MB | ~2 MB | ~1 KB |
| 100 MB | ~200 MB | ~1 KB |
| 1 GB | ~2 GB (OOM) | ~1 KB |
| 10 GB | OOM | ~1 KB |

---

## 二、修复方案详解

### 2.1 整体架构调整

```
当前架构：
┌─────────────┐    ┌─────────────────┐    ┌─────────────┐
│ 读取文件    │ -> │ 缓存所有行      │ -> │ 输出/写入   │
└─────────────┘    └─────────────────┘    └─────────────┘
                           ↑
                    内存瓶颈在这里

修复后架构：
┌─────────────┐    ┌─────────────────┐    ┌─────────────┐
│ 读取一行    │ -> │ 处理一行        │ -> │ 输出一行    │
└─────────────┘    └─────────────────┘    └─────────────┘
       ↑                                          |
       └────────────── 循环 ──────────────────────┘
                    
                    内存占用恒定
```

### 2.2 核心修改点

#### 修改 1：重构 processFile 函数

**文件**: `internal/commands/sed/cmd_sed.go`

**当前代码**（约第 170-200 行）：
```go
func processFile(config *SedConfig) error {
    // 打开文件
    file, err := os.Open(config.Target)
    if err != nil {
        return fmt.Errorf("cannot open file: %w", err)
    }
    defer func() { _ = file.Close() }()

    // 读取并处理每一行
    var resultLines []string  // 删除这行
    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        processedLine, _ := processLine(line, config, lineNum)
        resultLines = append(resultLines, processedLine)  // 删除这行
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }

    // 输出或写入文件
    if config.InPlace {
        return writeFile(config, resultLines)  // 修改这里
    }

    // 输出到标准输出
    for _, line := range resultLines {  // 修改这里
        fmt.Println(line)
    }

    return nil
}
```

**修复后代码**：
```go
func processFile(config *SedConfig) error {
    // 根据模式选择处理方式
    if config.InPlace {
        return processFileInPlace(config)
    }
    return processFilePreview(config)
}
```

#### 修改 2：新增 processFilePreview 函数（预览模式）

**添加位置**: `processFile` 函数之后

```go
// processFilePreview 预览模式处理文件
// 边读取边处理边输出到标准输出，不缓存任何行
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processFilePreview(config *SedConfig) error {
    file, err := os.Open(config.Target)
    if err != nil {
        return fmt.Errorf("cannot open file: %w", err)
    }
    defer func() { _ = file.Close() }()

    // 设置大缓冲区，避免超长行问题
    scanner := bufio.NewScanner(file)
    const maxCapacity = 1024 * 1024 // 1MB 缓冲区
    buf := make([]byte, maxCapacity)
    scanner.Buffer(buf, maxCapacity)

    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        processedLine, _ := processLine(line, config, lineNum)
        fmt.Println(processedLine)
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }

    return nil
}
```

#### 修改 3：新增 processFileInPlace 函数（原地修改模式）

**添加位置**: `processFilePreview` 函数之后

```go
// processFileInPlace 原地修改模式处理文件
// 边读取边处理边写入临时文件，最后原子替换
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processFileInPlace(config *SedConfig) error {
    // 打开源文件
    sourceFile, err := os.Open(config.Target)
    if err != nil {
        return fmt.Errorf("cannot open file: %w", err)
    }
    defer func() { _ = sourceFile.Close() }()

    // 创建临时文件（在同一目录，确保原子重命名有效）
    dir := filepath.Dir(config.Target)
    tempFile, err := os.CreateTemp(dir, ".sed-*.tmp")
    if err != nil {
        return fmt.Errorf("failed to create temp file: %w", err)
    }
    tempPath := tempFile.Name()
    
    // 确保临时文件在出错时被清理
    cleanup := func() {
        _ = tempFile.Close()
        _ = os.Remove(tempPath)
    }
    defer cleanup()

    // 如果需要备份
    if config.Backup {
        backupPath := config.Target + ".bak"
        if err := copyFile(config.Target, backupPath); err != nil {
            return fmt.Errorf("failed to create backup: %w", err)
        }
    }

    // 设置大缓冲区
    scanner := bufio.NewScanner(sourceFile)
    const maxCapacity = 1024 * 1024 // 1MB 缓冲区
    buf := make([]byte, maxCapacity)
    scanner.Buffer(buf, maxCapacity)

    // 使用缓冲写入器提高性能
    writer := bufio.NewWriter(tempFile)
    
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        processedLine, _ := processLine(line, config, lineNum)
        
        if _, err := writer.WriteString(processedLine); err != nil {
            return fmt.Errorf("failed to write temp file: %w", err)
        }
        if err := writer.WriteByte('\n'); err != nil {
            return fmt.Errorf("failed to write temp file: %w", err)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading file: %w", err)
    }

    // 刷新缓冲区
    if err := writer.Flush(); err != nil {
        return fmt.Errorf("failed to flush temp file: %w", err)
    }

    // 关闭文件
    if err := tempFile.Close(); err != nil {
        return fmt.Errorf("failed to close temp file: %w", err)
    }

    // 取消自动清理（准备提交）
    cleanup = func() {}

    // 原子替换
    if err := os.Rename(tempPath, config.Target); err != nil {
        // 替换失败，尝试恢复
        _ = os.Remove(tempPath)
        return fmt.Errorf("failed to replace file: %w", err)
    }

    return nil
}
```

#### 修改 4：删除旧的 writeFile 函数

**删除范围**: 约第 450-520 行（原 writeFile 和 copyFile 函数）

**注意**: `copyFile` 函数需要保留（用于备份），但 `writeFile` 函数不再需要。

---

## 三、边界情况处理

### 3.1 超长行处理

```go
// 在每个 scanner 创建后添加
const maxCapacity = 1024 * 1024 // 1MB 缓冲区
buf := make([]byte, maxCapacity)
scanner.Buffer(buf, maxCapacity)
```

**说明**:
- 默认 64KB 可能不够
- 1MB 可以处理绝大多数场景
- 如果超过 1MB，scanner 会报错，不会静默截断

### 3.2 错误处理和临时文件清理

```go
// 使用 defer 确保清理
cleanup := func() {
    _ = tempFile.Close()
    _ = os.Remove(tempPath)
}
defer cleanup()

// 成功时取消清理
cleanup = func() {}
```

### 3.3 最后一行无换行符

**当前行为**: `fmt.Println` 会添加换行符

**考虑**: 是否需要保留原文件的换行符行为？

**建议**: 添加标志控制，或默认保留原行为

```go
// 如果需要精确保留原换行符行为
func processFileInPlace(config *SedConfig) error {
    // ...
    firstLine := true
    for scanner.Scan() {
        // ...
        if !firstLine {
            writer.WriteByte('\n')
        }
        firstLine = false
        writer.WriteString(processedLine)
    }
    // ...
}
```

---

## 四、测试方案

### 4.1 单元测试

```go
// cmd_sed_test.go

func TestProcessFilePreview(t *testing.T) {
    // 创建临时测试文件
    content := "hello world\nfoo bar\n"
    tmpFile := createTempFile(content)
    defer os.Remove(tmpFile)
    
    config := SedConfig{
        Target:      tmpFile,
        Pattern:     "hello",
        Replacement: "hi",
    }
    
    // 捕获标准输出
    output := captureStdout(func() {
        err := processFilePreview(&config)
        if err != nil {
            t.Fatal(err)
        }
    })
    
    expected := "hi world\nfoo bar\n"
    if output != expected {
        t.Errorf("expected %q, got %q", expected, output)
    }
}

func TestProcessFileInPlace(t *testing.T) {
    content := "hello world\nfoo bar\n"
    tmpFile := createTempFile(content)
    defer os.Remove(tmpFile)
    
    config := SedConfig{
        Target:      tmpFile,
        Pattern:     "hello",
        Replacement: "hi",
        InPlace:     true,
    }
    
    err := processFileInPlace(&config)
    if err != nil {
        t.Fatal(err)
    }
    
    result := readFile(tmpFile)
    expected := "hi world\nfoo bar\n"
    if result != expected {
        t.Errorf("expected %q, got %q", expected, result)
    }
}

func TestLargeFile(t *testing.T) {
    // 生成 100MB 的测试文件
    // 验证内存占用不超过 10MB
}
```

### 4.2 手动测试命令

```bash
# 1. 基本功能测试
echo -e "hello world\nfoo bar" > /tmp/test.txt
fck sed -p "hello" -r "hi" /tmp/test.txt
fck sed -p "hello" -r "hi" -i /tmp/test.txt

# 2. 大文件测试（100MB）
dd if=/dev/zero of=/tmp/bigfile.txt bs=1M count=100
# 监控内存使用
/usr/bin/time -v fck sed -p "test" -r "TEST" /tmp/bigfile.txt

# 3. 超大文件测试（1GB）
dd if=/dev/zero of=/tmp/hugefile.txt bs=1M count=1024
fck sed -p "test" -r "TEST" -i /tmp/hugefile.txt

# 4. 超长行测试
python3 -c "print('A' * 1000000)" > /tmp/longline.txt
fck sed -p "A" -r "B" /tmp/longline.txt
```

---

## 五、性能对比

### 5.1 预期性能

| 指标 | 修复前 | 修复后 |
|------|--------|--------|
| 内存占用 | O(n) | O(1) |
| 处理速度 | 2-pass | 1-pass |
| 大文件支持 | ❌ | ✅ |
| 磁盘 I/O | 2 次读取 + 1 次写入 | 1 次读取 + 1 次写入 |

### 5.2 实际测试数据（待填充）

| 文件大小 | 修复前内存 | 修复后内存 | 修复前时间 | 修复后时间 |
|----------|-----------|-----------|-----------|-----------|
| 1 MB | | | | |
| 10 MB | | | | |
| 100 MB | | | | |
| 1 GB | OOM | | - | |

---

## 六、风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 引入新 Bug | 中 | 高 | 充分测试，保留原测试用例 |
| 性能下降 | 低 | 中 | 使用缓冲写入器 |
| 换行符行为改变 | 中 | 低 | 文档说明，或添加标志控制 |
| 超长行处理错误 | 低 | 中 | 增大缓冲区，明确报错 |

---

## 七、实施步骤

1. **备份原文件**
   ```bash
   cp internal/commands/sed/cmd_sed.go internal/commands/sed/cmd_sed.go.bak
   ```

2. **按修改点逐步实施**
   - 修改 1：重构 processFile
   - 修改 2：添加 processFilePreview
   - 修改 3：添加 processFileInPlace
   - 修改 4：删除 writeFile

3. **编译验证**
   ```bash
   go build ./...
   ```

4. **运行测试**
   ```bash
   go test ./internal/commands/sed/...
   ```

5. **手动验证**
   ```bash
   go run cmd/main.go sed -p "test" -r "TEST" test.txt
   go run cmd/main.go sed -p "test" -r "TEST" -i test.txt
   ```

---

## 八、总结

**修复收益**:
- ✅ 解决大文件 OOM 问题
- ✅ 内存占用恒定，可处理任意大小文件
- ✅ 代码结构更清晰（分离预览和原地修改）
- ✅ 性能提升（单遍处理）

**实施成本**:
- 中等（需要重构核心函数）
- 需要充分测试

**建议**: 按此方案实施，收益大于风险。

---

**方案完成，等待确认后实施。**
