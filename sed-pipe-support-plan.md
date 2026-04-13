# Sed 命令管道输入支持方案

> **目标**：支持从管道读取输入，实现 `echo "hello" | fck sed -p "hello" -r "hi"`  
> **方案**：`-` 特殊文件参数 + 标准输入流式处理  
> **预期效果**：与 Linux sed 行为一致，支持管道和重定向

---

## 一、需求分析

### 1.1 目标使用场景

```bash
# 场景 1：管道输入
echo "hello world" | fck sed -p "hello" -r "hi"

# 场景 2：重定向输入
fck sed -p "hello" -r "hi" < input.txt

# 场景 3：显式指定标准输入
cat file.txt | fck sed -p "hello" -r "hi" -

# 场景 4：与其他命令组合
cat log.txt | grep "ERROR" | fck sed -p "ERROR" -r "WARN" -
```

### 1.2 约束条件

| 约束 | 说明 |
|------|------|
| `-i` 原地修改 | 管道输入时不支持（无文件可修改） |
| `-b` 备份 | 管道输入时不支持 |
| 目标文件参数 | 变为可选，支持 `-` 表示 stdin |

---

## 二、方案设计

### 2.1 整体架构

```
用户输入
    │
    ├─ 有文件参数 ──> 原流程（文件处理）
    │
    └─ 无文件参数 / `-` ──> 新流程（stdin 处理）
                              │
                              ├─ 有 `-i` 标志 ──> 报错（不支持）
                              │
                              └─ 无 `-i` 标志 ──> 流式处理 stdin
```

### 2.2 核心修改点

#### 修改 1：CLI 层 - 修改参数校验

**文件**: `internal/cli/sed.go`

**当前代码**（约第 70-95 行）：
```go
func runSed(cmd qflag.Command) error {
    args := cmd.Args()

    // 检查文件参数
    if len(args) < 1 {
        return fmt.Errorf("no target file specified")
    }

    config := sed.SedConfig{
        Target:      args[0],
        // ...
    }
    // ...
}
```

**修复后代码**：
```go
func runSed(cmd qflag.Command) error {
    args := cmd.Args()

    var target string
    var fromStdin bool

    if len(args) == 0 {
        // 无参数：从标准输入读取
        fromStdin = true
    } else if args[0] == "-" {
        // 显式指定 -：从标准输入读取
        fromStdin = true
    } else {
        // 文件路径
        target = args[0]
    }

    config := sed.SedConfig{
        Target:      target,
        FromStdin:   fromStdin,  // 新增字段
        // ...
    }

    // 管道输入时不支持原地修改
    if fromStdin && sedInPlace.Get() {
        return fmt.Errorf("cannot use -i/--in-place with stdin input")
    }

    // ...
}
```

#### 修改 2：Config 结构体 - 添加 stdin 支持

**文件**: `internal/commands/sed/cmd_sed.go`

**当前代码**（约第 20-40 行）：
```go
type SedConfig struct {
    Target      string // 目标文件路径
    Pattern     string // 搜索模式
    Replacement string // 替换内容
    Regexp      bool   // 是否使用正则表达式
    LineRange   string // 行号范围（如 "1,10" 或 "5,$"）
    InPlace     bool   // 是否原地修改
    Backup      bool   // 是否创建备份
    MaxCount    int    // 最大替换次数（0表示无限制）
    IgnoreCase  bool   // 是否忽略大小写

    // 内部使用
    compiledPattern *regexp.Regexp
    lineStart       int
    lineEnd         int
}
```

**修复后代码**：
```go
type SedConfig struct {
    Target      string // 目标文件路径（FromStdin 为 true 时为空）
    FromStdin   bool   // 是否从标准输入读取
    Pattern     string // 搜索模式
    Replacement string // 替换内容
    Regexp      bool   // 是否使用正则表达式
    LineRange   string // 行号范围（如 "1,10" 或 "5,$"）
    InPlace     bool   // 是否原地修改
    Backup      bool   // 是否创建备份
    MaxCount    int    // 最大替换次数（0表示无限制）
    IgnoreCase  bool   // 是否忽略大小写

    // 内部使用
    compiledPattern *regexp.Regexp
    lineStart       int
    lineEnd         int
}
```

#### 修改 3：核心逻辑 - 添加 stdin 处理

**文件**: `internal/commands/sed/cmd_sed.go`

**修改位置**: `processFile` 函数

**当前代码**：
```go
func processFile(config *SedConfig) error {
    if config.InPlace {
        return processFileInPlace(config)
    }
    return processFilePreview(config)
}
```

**修复后代码**：
```go
func processFile(config *SedConfig) error {
    // 从标准输入读取
    if config.FromStdin {
        return processStdin(config)
    }

    if config.InPlace {
        return processFileInPlace(config)
    }
    return processFilePreview(config)
}
```

#### 修改 4：新增 processStdin 函数

**文件**: `internal/commands/sed/cmd_sed.go`

**添加位置**: `processFilePreview` 函数之后

```go
// processStdin 从标准输入读取并处理
// 边读取边处理边输出到标准输出，支持管道和重定向
//
// 参数:
//   - config: 命令配置指针
//
// 返回:
//   - error: 处理错误
func processStdin(config *SedConfig) error {
    // 设置大缓冲区，避免超长行问题
    scanner := bufio.NewScanner(os.Stdin)
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
        return fmt.Errorf("error reading stdin: %w", err)
    }

    return nil
}
```

---

## 三、边界情况处理

### 3.1 互斥标志检查

```go
// CLI 层检查
if fromStdin {
    if sedInPlace.Get() {
        return fmt.Errorf("cannot use -i/--in-place with stdin input")
    }
    if sedBackup.Get() {
        return fmt.Errorf("cannot use -b/--backup with stdin input")
    }
}
```

### 3.2 空输入处理

```go
func processStdin(config *SedConfig) error {
    scanner := bufio.NewScanner(os.Stdin)
    // ...

    lineNum := 0
    for scanner.Scan() {
        lineNum++
        // 处理每一行...
    }

    // 空输入是合法的，直接返回
    if lineNum == 0 {
        return nil
    }

    // ...
}
```

### 3.3 信号处理（可选增强）

```go
// 处理 Ctrl+C 等信号，确保输出不中断
func processStdin(config *SedConfig) error {
    // 设置行缓冲，确保及时输出
    // Go 的 fmt.Println 默认会缓冲，但通常够用了
    // 如果需要更及时的输出，可以使用无缓冲输出

    scanner := bufio.NewScanner(os.Stdin)
    // ...
}
```

---

## 四、使用示例

### 4.1 基本用法

```bash
# 管道输入
echo "hello world" | fck sed -p "hello" -r "hi"
# 输出: hi world

# 重定向输入
echo "foo bar" > /tmp/test.txt
fck sed -p "foo" -r "baz" < /tmp/test.txt
# 输出: baz bar

# 显式指定 -
echo "abc def" | fck sed -p "abc" -r "xyz" -
# 输出: xyz def
```

### 4.2 高级用法

```bash
# 多行处理
cat <<EOF | fck sed -p "^" -r "> "
line 1
line 2
line 3
EOF
# 输出:
# > line 1
# > line 2
# > line 3

# 与其他命令组合
cat /var/log/syslog | grep "error" | fck sed -p "error" -r "ERROR" - | head -10

# 正则替换
echo "user123@example.com" | fck sed -p "user[0-9]+" -r "guest" -r
# 输出: guest@example.com
```

### 4.3 错误示例

```bash
# 错误：管道输入不支持原地修改
echo "test" | fck sed -p "test" -r "TEST" -i
# 错误: cannot use -i/--in-place with stdin input

# 错误：管道输入不支持备份
echo "test" | fck sed -p "test" -r "TEST" -b
# 错误: cannot use -b/--backup with stdin input
```

---

## 五、测试方案

### 5.1 单元测试

```go
// cmd_sed_test.go

func TestProcessStdin(t *testing.T) {
    // 重定向标准输入
    oldStdin := os.Stdin
    defer func() { os.Stdin = oldStdin }()

    r, w, _ := os.Pipe()
    os.Stdin = r

    // 写入测试数据
    go func() {
        w.WriteString("hello world\nfoo bar\n")
        w.Close()
    }()

    config := SedConfig{
        FromStdin:   true,
        Pattern:     "hello",
        Replacement: "hi",
    }

    // 捕获标准输出
    output := captureStdout(func() {
        err := processStdin(&config)
        if err != nil {
            t.Fatal(err)
        }
    })

    expected := "hi world\nfoo bar\n"
    if output != expected {
        t.Errorf("expected %q, got %q", expected, output)
    }
}

func TestStdinWithInPlaceError(t *testing.T) {
    config := SedConfig{
        FromStdin:   true,
        InPlace:     true,
        Pattern:     "test",
        Replacement: "TEST",
    }

    err := processFile(&config)
    if err == nil {
        t.Error("expected error for stdin with in-place")
    }
}
```

### 5.2 手动测试命令

```bash
# 1. 基本管道测试
echo "hello world" | go run cmd/main.go sed -p "hello" -r "hi"

# 2. 重定向测试
echo "foo bar" > /tmp/test.txt
go run cmd/main.go sed -p "foo" -r "baz" < /tmp/test.txt

# 3. 显式 - 参数测试
echo "abc def" | go run cmd/main.go sed -p "abc" -r "xyz" -

# 4. 错误测试
echo "test" | go run cmd/main.go sed -p "test" -r "TEST" -i

# 5. 多行管道测试
seq 1 10 | go run cmd/main.go sed -p "5" -r "FIVE"

# 6. 正则测试
echo "user123@example.com" | go run cmd/main.go sed -p "user[0-9]+" -r "guest" -r
```

---

## 六、与现有代码的兼容性

### 6.1 向后兼容

| 现有用法 | 兼容性 |
|----------|--------|
| `fck sed -p "x" -r "y" file.txt` | ✅ 完全兼容 |
| `fck sed -p "x" -r "y" -i file.txt` | ✅ 完全兼容 |
| `fck sed -p "x" -r "y" -b file.txt` | ✅ 完全兼容 |

### 6.2 新增功能

| 新用法 | 状态 |
|--------|------|
| `echo "x" \| fck sed -p "x" -r "y"` | ✅ 新增支持 |
| `fck sed -p "x" -r "y" < file.txt` | ✅ 新增支持 |
| `cat file \| fck sed -p "x" -r "y" -` | ✅ 新增支持 |

---

## 七、实施步骤

1. **修改 Config 结构体**
   - 添加 `FromStdin bool` 字段

2. **修改 CLI 层**
   - 更新参数解析逻辑
   - 添加互斥标志检查

3. **修改核心逻辑**
   - 更新 `processFile` 函数
   - 添加 `processStdin` 函数

4. **编译验证**
   ```bash
   go build ./...
   ```

5. **运行测试**
   ```bash
   go test ./internal/commands/sed/...
   ```

6. **手动验证**
   ```bash
   # 管道测试
   echo "hello" | go run cmd/main.go sed -p "hello" -r "hi"

   # 重定向测试
   go run cmd/main.go sed -p "hello" -r "hi" < test.txt

   # 错误测试
   echo "test" | go run cmd/main.go sed -p "test" -r "TEST" -i
   ```

---

## 八、总结

**方案优势**:
- ✅ 符合 Unix 工具惯例（`-` 表示 stdin）
- ✅ 向后完全兼容
- ✅ 实现简单，改动范围小
- ✅ 流式处理，无内存问题

**实施成本**:
- 低（约 3 个文件，50 行代码）

**建议**: 按此方案实施，增强 sed 命令的实用性。

---

**方案完成，等待确认后实施。**
