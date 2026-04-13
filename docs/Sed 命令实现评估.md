## Sed 命令实现评估

### 1. 整体架构

| 方面       | 评价  | 说明                                       |
| ---------- | ----- | ------------------------------------------ |
| 代码结构   | ⭐⭐⭐⭐  | 分层清晰：CLI → Config → 处理函数          |
| 功能完整性 | ⭐⭐⭐⭐⭐ | 支持正则、行号范围、大小写忽略、最大次数等 |
| 错误处理   | ⭐⭐⭐⭐  | 有基本的错误处理和返回                     |
| 安全性     | ⭐⭐⭐⭐  | 使用临时文件 + 原子替换                    |

### 2. 潜在问题

#### 🔴 **问题 1：大文件处理（严重）**

```go
// processFile 函数中
var resultLines []string
scanner := bufio.NewScanner(file)
// ...
resultLines = append(resultLines, processedLine)
```

**问题**：将整个文件内容加载到内存，大文件会 OOM。

**建议**：

```go
// 流式处理，不缓存所有行
if config.InPlace {
    // 直接写入临时文件
    tempFile, _ := os.CreateTemp(dir, ".sed-*.tmp")
    for scanner.Scan() {
        line := scanner.Text()
        processedLine, _ := processLine(line, config, lineNum)
        fmt.Fprintln(tempFile, processedLine)
    }
} else {
    // 直接输出到 stdout
    for scanner.Scan() {
        line := scanner.Text()
        processedLine, _ := processLine(line, config, lineNum)
        fmt.Println(processedLine)
    }
}
```

---

#### 🟡 **问题 2：Scanner 缓冲区限制**

```go
scanner := bufio.NewScanner(file)
```

**问题**：`bufio.Scanner` 默认缓冲区 64KB，超过会报错。

**建议**：
```go
const maxCapacity = 1024 * 1024 // 1MB
buf := make([]byte, maxCapacity)
scanner.Buffer(buf, maxCapacity)
```

---

#### 🟡 **问题 3：替换计数逻辑不一致**

```go
// replaceWithString 中
if config.MaxCount > 0 && config.replaceCount >= config.MaxCount {
    return line, false
}

// replaceWithRegex 中
if config.MaxCount > 0 {
    matches := config.compiledPattern.FindAllStringIndex(line, -1)
    // ...
}
```

**问题**：
- 字符串替换：跨行累计计数
- 正则替换：也是跨行累计，但实现方式不同

**潜在 Bug**：如果一行有多个匹配，且超过 MaxCount，处理逻辑可能不一致。

---

#### 🟡 **问题 4：忽略大小写的字符串替换问题**

```go
func replaceStringIgnoreCase(line, pattern, replacement string, config *SedConfig) {
    lowerLine := strings.ToLower(line)
    lowerPattern := strings.ToLower(pattern)
    // ...
    actualIdx := start + idx
    result.WriteString(line[start:actualIdx])
    result.WriteString(replacement)
    start = actualIdx + len(pattern)  // ⚠️ 问题在这里
}
```

**问题**：`len(pattern)` 是字节长度，如果 pattern 包含多字节字符（如中文），会导致错位。

**示例**：
```
line: "你好世界"
pattern: "你好"  (len=6)
replacement: "Hello"

实际替换后：start = 0 + 6 = 6，但 "你好" 只占 6 字节，正确
但如果 pattern 和 line 的大小写转换后长度不一致...
```

---

#### 🟢 **问题 5：备份文件覆盖**

```go
if config.Backup {
    backupPath := config.Target + ".bak"
    if err := copyFile(config.Target, backupPath); err != nil {
        return fmt.Errorf("failed to create backup: %w", err)
    }
}
```

**问题**：如果 `.bak` 文件已存在，会直接覆盖，无提示。

**建议**：检查备份文件是否存在，或添加时间戳。

---

#### 🟢 **问题 6：临时文件清理**

```go
tempFile, err := os.CreateTemp(dir, ".sed-*.tmp")
// ...
if err := writer.Flush(); err != nil {
    _ = tempFile.Close()
    _ = os.Remove(tempPath)  // ✅ 有清理
    return fmt.Errorf("failed to write temp file: %w", err)
}
```

**评价**：清理逻辑基本完善，但可以使用 `defer` 更优雅。

---

#### 🟢 **问题 7：无替换提示**

```go
// 执行后没有提示用户是否发生了替换
if config.InPlace {
    return writeFile(config, resultLines)
}
```

**建议**：添加 `-v`（verbose）标志，显示替换统计信息。

---

### 3. 功能缺失（非问题，可扩展）

| 功能         | 说明                           |
| ------------ | ------------------------------ |
| 多文件处理   | 目前只支持单个文件             |
| 全局替换计数 | 无法知道总共替换了多少处       |
| 预览差异     | 无法显示哪些行发生了变化       |
| 正则捕获组   | 不支持 `$1`, `$2` 等捕获组引用 |

---

### 4. 测试建议

应该测试以下场景：

```bash
# 1. 基本替换
fck sed -p "old" -r "new" file.txt

# 2. 大文件测试（>100MB）
fck sed -p "test" -r "TEST" -i bigfile.txt

# 3. 超长行测试（>64KB）
fck sed -p "a" -r "b" longline.txt

# 4. 中文替换
fck sed -p "你好" -r "Hello" -I chinese.txt

# 5. 边界：MaxCount
fck sed -p "a" -r "b" -c 5 file.txt

# 6. 边界：行号范围
fck sed -p "old" -r "new" -n 1,10 file.txt
```

---

### 5. 总结

| 类别       | 数量                            |
| ---------- | ------------------------------- |
| 🔴 严重问题 | 1（大文件内存问题）             |
| 🟡 中等问题 | 3（缓冲区、计数逻辑、中文处理） |
| 🟢 轻微问题 | 2（备份覆盖、清理）             |
| 功能建议   | 4                               |

**整体评价**：功能完整，代码结构良好，但需要修复大文件处理问题才能用于生产环境。

