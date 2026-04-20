# head / tail 命令实现方案

## 1. 设计目标

实现独立的 `fck head` 和 `fck tail` 子命令，完全替代 `cat` 命令的 `-u/-d` 功能。

## 2. 使用示例

### head 命令

```bash
# 查看文件前10行（默认）
fck head file.txt

# 查看前20行
fck head -n 20 file.txt

# 查看前100字节
fck head -c 100 file.txt

# 多个文件
fck head -n 5 file1.txt file2.txt

# 从管道读取
cat big.log | fck head -n 100

# 简洁模式
fck head -q file.txt
```

### tail 命令

```bash
# 查看文件后10行（默认）
fck tail file.txt

# 查看后20行
fck tail -n 20 file.txt

# 查看后100字节
fck tail -c 100 file.txt

# 实时追踪文件（核心功能）
fck tail -f /var/log/app.log

# 追踪并显示最后100行
tail -n 100 -f /var/log/app.log

# 多个文件
tail file1.txt file2.txt

# 从管道读取
cat big.log | tail -n 50
```

## 3. 输出格式

### head 标准输出

```
==> file1.txt <==
line 1
line 2
...
line 10

==> file2.txt <==
line 1
line 2
...
line 10
```

### tail 标准输出

```
==> file1.txt <==
line 91
line 92
...
line 100

==> file2.txt <==
line 45
line 46
...
line 50
```

### 简洁模式 (-q)

```
line 1
line 2
...
line 10
```

## 4. 数据结构设计

### head 配置

```go
// HeadConfig head 命令配置
type HeadConfig struct {
    Targets   []string // 目标文件列表
    Lines     int      // -n, 行数（默认10）
    Bytes     int64    // -c, 字节数（与 -n 互斥）
    Quiet     bool     // -q, 不显示文件名标题
    Verbose   bool     // -v, 总是显示文件名标题
    FromStdin bool     // 是否从标准输入读取
}
```

### tail 配置

```go
// TailConfig tail 命令配置
type TailConfig struct {
    Targets   []string // 目标文件列表
    Lines     int      // -n, 行数（默认10）
    Bytes     int64    // -c, 字节数（与 -n 互斥）
    Follow    bool     // -f, 实时追踪
    Quiet     bool     // -q, 不显示文件名标题
    Verbose   bool     // -v, 总是显示文件名标题
    FromStdin bool     // 是否从标准输入读取
    SleepInterval time.Duration // 轮询间隔（默认100ms）
}

// TailFile tail 文件状态追踪
type TailFile struct {
    Path     string
    File     *os.File
    Reader   *bufio.Reader
    Size     int64
}
```

## 5. CLI 参数设计

### head CLI

```go
var (
    headLines   *qflag.IntFlag   // -n, --lines    行数
    headBytes   *qflag.Int64Flag // -c, --bytes    字节数
    headQuiet   *qflag.BoolFlag  // -q, --quiet    不显示文件名
    headVerbose *qflag.BoolFlag  // -v, --verbose  总是显示文件名
)

// 互斥组：-n 和 -c 不能同时使用
MutexGroups: []qflag.MutexGroup{
    {
        Name:      "count-mode",
        Flags:     []string{"lines", "bytes"},
        AllowNone: true, // 默认使用 -n 10
    },
}
```

### tail CLI

```go
var (
    tailLines   *qflag.IntFlag     // -n, --lines    行数
    tailBytes   *qflag.Int64Flag   // -c, --bytes    字节数
    tailFollow  *qflag.BoolFlag    // -f, --follow   实时追踪
    tailQuiet   *qflag.BoolFlag    // -q, --quiet    不显示文件名
    tailVerbose *qflag.BoolFlag    // -v, --verbose  总是显示文件名
)

// 互斥组：-n 和 -c 不能同时使用
MutexGroups: []qflag.MutexGroup{
    {
        Name:      "count-mode",
        Flags:     []string{"lines", "bytes"},
        AllowNone: true,
    },
}
```

## 6. 核心实现逻辑

### head 主函数

```go
func HeadCmdMain(config HeadConfig) error {
    // 默认行数
    if config.Lines == 0 && config.Bytes == 0 {
        config.Lines = 10
    }

    // 从标准输入读取
    if len(config.Targets) == 0 || (len(config.Targets) == 1 && config.Targets[0] == "-") {
        return headStdin(config, os.Stdin)
    }

    // 多个文件处理
    showHeader := !config.Quiet && (config.Verbose || len(config.Targets) > 1)

    for i, path := range config.Targets {
        if i > 0 {
            fmt.Println()
        }

        if showHeader {
            fmt.Printf("==> %s <==\n", path)
        }

        if err := headFile(path, config); err != nil {
            fmt.Fprintf(os.Stderr, "head: %s: %v\n", path, err)
        }
    }

    return nil
}

// headFile 处理单个文件
func headFile(path string, config HeadConfig) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    if config.Bytes > 0 {
        return headByBytes(file, config.Bytes)
    }
    return headByLines(file, config.Lines)
}

// headByLines 按行读取
func headByLines(file *os.File, n int) error {
    scanner := bufio.NewScanner(file)
    count := 0

    for scanner.Scan() {
        fmt.Println(scanner.Text())
        count++
        if count >= n {
            break
        }
    }

    return scanner.Err()
}

// headByBytes 按字节读取
func headByBytes(file *os.File, n int64) error {
    reader := io.LimitReader(file, n)
    _, err := io.Copy(os.Stdout, reader)
    return err
}

// headStdin 从标准输入读取
func headStdin(config HeadConfig, stdin io.Reader) error {
    if config.Bytes > 0 {
        reader := io.LimitReader(stdin, config.Bytes)
        _, err := io.Copy(os.Stdout, reader)
        return err
    }

    scanner := bufio.NewScanner(stdin)
    count := 0

    for scanner.Scan() {
        fmt.Println(scanner.Text())
        count++
        if count >= config.Lines {
            break
        }
    }

    return scanner.Err()
}
```

### tail 主函数

```go
func TailCmdMain(config TailConfig) error {
    // 默认行数
    if config.Lines == 0 && config.Bytes == 0 {
        config.Lines = 10
    }

    // 从标准输入读取
    if len(config.Targets) == 0 || (len(config.Targets) == 1 && config.Targets[0] == "-") {
        return tailStdin(config, os.Stdin)
    }

    // 实时追踪模式
    if config.Follow {
        return tailFollowFiles(config)
    }

    // 普通模式
    return tailFiles(config)
}

// tailFiles 普通模式处理多个文件
func tailFiles(config TailConfig) error {
    showHeader := !config.Quiet && (config.Verbose || len(config.Targets) > 1)

    for i, path := range config.Targets {
        if i > 0 {
            fmt.Println()
        }

        if showHeader {
            fmt.Printf("==> %s <==\n", path)
        }

        if err := tailFile(path, config); err != nil {
            fmt.Fprintf(os.Stderr, "tail: %s: %v\n", path, err)
        }
    }

    return nil
}

// tailFile 处理单个文件
func tailFile(path string, config TailConfig) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    if config.Bytes > 0 {
        return tailByBytes(file, config.Bytes)
    }
    return tailByLines(file, config.Lines)
}

// tailByLines 按行读取（环形缓冲区）
func tailByLines(file *os.File, n int) error {
    ring := make([]string, n)
    index := 0
    count := 0

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        ring[index] = scanner.Text()
        index = (index + 1) % n
        count++
    }

    // 计算起始位置
    start := 0
    if count >= n {
        start = index
    }

    // 输出
    linesToPrint := n
    if count < n {
        linesToPrint = count
    }

    for i := 0; i < linesToPrint; i++ {
        fmt.Println(ring[(start+i)%n])
    }

    return scanner.Err()
}

// tailByBytes 按字节读取
func tailByBytes(file *os.File, n int64) error {
    stat, err := file.Stat()
    if err != nil {
        return err
    }

    size := stat.Size()
    start := int64(0)
    if size > n {
        start = size - n
    }

    _, err = file.Seek(start, io.SeekStart)
    if err != nil {
        return err
    }

    _, err = io.Copy(os.Stdout, file)
    return err
}

// tailFollowFiles 实时追踪多个文件
func tailFollowFiles(config TailConfig) error {
    // 打开所有文件
    files := make([]*TailFile, 0, len(config.Targets))
    for _, path := range config.Targets {
        tf, err := openTailFile(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "tail: %s: %v\n", path, err)
            continue
        }
        files = append(files, tf)
    }

    if len(files) == 0 {
        return fmt.Errorf("no files to follow")
    }

    // 显示文件名标题
    showHeader := !config.Quiet && len(files) > 1

    // 先显示初始内容
    for _, tf := range files {
        if showHeader {
            fmt.Printf("==> %s <==\n", tf.Path)
        }
        if err := tailFile(tf.Path, TailConfig{Lines: config.Lines}); err != nil {
            fmt.Fprintf(os.Stderr, "tail: %s: %v\n", tf.Path, err)
        }
    }

    // 进入追踪模式
    ticker := time.NewTicker(config.SleepInterval)
    defer ticker.Stop()

    for range ticker.C {
        for _, tf := range files {
            if err := followFile(tf, showHeader); err != nil {
                // 文件可能被删除或重命名，尝试重新打开
                if newTf, err := reopenTailFile(tf); err == nil {
                    *tf = *newTf
                }
            }
        }
    }

    return nil
}

// followFile 追踪单个文件的新内容
func followFile(tf *TailFile, showHeader bool) error {
    stat, err := tf.File.Stat()
    if err != nil {
        return err
    }

    newSize := stat.Size()
    if newSize < tf.Size {
        // 文件被截断或替换
        tf.File.Seek(0, io.SeekStart)
        tf.Size = 0
        if showHeader {
            fmt.Printf("\n==> %s <==\n", tf.Path)
        }
    } else if newSize == tf.Size {
        // 没有新内容
        return nil
    }

    // 读取新内容
    for {
        line, err := tf.Reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }

        if showHeader && tf.Size == 0 {
            fmt.Printf("\n==> %s <==\n", tf.Path)
        }

        fmt.Print(line)
        tf.Size += int64(len(line))
    }

    return nil
}

// openTailFile 打开文件用于追踪
func openTailFile(path string) (*TailFile, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    // 定位到文件末尾
    file.Seek(0, io.SeekEnd)

    return &TailFile{
        Path:   path,
        File:   file,
        Reader: bufio.NewReader(file),
        Size:   stat.Size(),
    }, nil
}
```

## 7. 文件结构

```
internal/commands/head/
└── cmd_head.go      # head 业务逻辑

internal/commands/tail/
└── cmd_tail.go      # tail 业务逻辑

internal/cli/
├── head.go          # head CLI 定义
├── tail.go          # tail CLI 定义
└── cat.go           # 修改：移除 -u/-d 参数
```

## 8. 实现步骤

### 步骤 1：创建 head 命令

1. 创建 `internal/commands/head/cmd_head.go`
2. 创建 `internal/cli/head.go`
3. 在 `root.go` 注册 head 命令

### 步骤 2：创建 tail 命令

1. 创建 `internal/commands/tail/cmd_tail.go`
2. 创建 `internal/cli/tail.go`
3. 在 `root.go` 注册 tail 命令

### 步骤 3：修改 cat 命令

1. 移除 `cat.go` 中的 `-u/--head` 和 `-d/--tail` 参数
2. 移除 `cat.go` 中的 `MutexGroups` 相关配置
3. 修改 `cmd_cat.go` 中的 `CatConfig` 结构体
4. 移除 `viewer.go` 中的 head/tail 处理逻辑

### 步骤 4：编译验证

```bash
go build ./...
./fck head -n 10 file.txt
./fck tail -f /var/log/app.log
```

## 9. cat 命令修改清单

### 移除的参数

```go
// 从 cli/cat.go 移除
var (
    catHeadLines *qflag.IntFlag  // -u, --head
    catTailLines *qflag.IntFlag  // -d, --tail
)

// 从 cli/cat.go 移除互斥组
{
    Name:      "head-tail",
    Flags:     []string{"head", "tail"},
    AllowNone: true,
},
```

### 修改 CatConfig

```go
// 从 cmd_cat.go 移除
HeadLines int  // --head 显示前N行
TailLines int  // --tail 显示后N行
```

### 修改 viewer.go

移除 `FileViewer` 中的 head/tail 处理逻辑，只保留分页查看功能。

## 10. 与 Linux 对比

| 功能 | Linux head/tail | fck head/tail |
|------|-----------------|---------------|
| 默认行数 | 10 | 10 |
| -n 行数 | ✅ | ✅ |
| -c 字节 | ✅ | ✅ |
| -q 简洁 | ✅ | ✅ |
| -v 显示文件名 | ✅ | ✅ |
| tail -f 追踪 | ✅ | ✅ |
| 多文件 | ✅ | ✅ |
| 管道输入 | ✅ | ✅ |
| -F 重试追踪 | ❌ | 可考虑 |
| --pid 进程追踪 | ❌ | 可考虑 |

## 11. 注意事项

1. **大文件处理**：tail 使用环形缓冲区，内存固定
2. **文件删除**：tail -f 时文件被删除，尝试重新打开
3. **文件截断**：tail -f 时文件被清空，重新从开头读取
4. **编码问题**：scanner 默认按行分割，处理大行可能需要调整 buffer
5. **跨平台**：使用标准库，完全跨平台

## 12. 依赖

```go
import (
    "bufio"
    "fmt"
    "io"
    "os"
    "time"
)
```

无外部依赖，纯标准库实现。
