# Cat 命令二进制文件检测功能方案

> **目标**：为 cat 命令添加二进制文件检测功能，避免直接输出二进制文件导致终端乱码  
> **参考**：grep 命令的 `-a` 和 `-I` 标志实现  
> **预期效果**：与 grep 保持一致的用户体验，增强工具一致性

---

## 一、需求分析

### 1.1 当前问题

```bash
# 当前行为：直接输出二进制文件内容
$ fck cat /bin/ls
[一堆乱码，可能包含控制字符，破坏终端显示]
```

### 1.2 目标使用场景

```bash
# 场景 1：默认行为（输出警告但继续显示）
$ fck cat /bin/ls
Binary file /bin/ls matches
[乱码内容...]

# 场景 2：静默跳过二进制文件
$ fck cat -I /bin/ls
# 无输出

# 场景 3：强制作为文本处理
$ fck cat -a /bin/ls
[直接输出，无警告]

# 场景 4：批量处理时跳过二进制
$ fck cat -I *.log *.bin
# 只显示 .log 文件内容，.bin 文件被静默跳过
```

---

## 二、方案设计

### 2.1 标志定义（与 grep 保持一致）

| 标志 | 长选项 | 说明 | 行为 |
|------|--------|------|------|
| `-a` | `--text` | 强制文本模式 | 跳过二进制检测，直接输出内容 |
| `-I` | `--ignore-binary` | 忽略二进制文件 | 检测到二进制时跳过，不输出内容和警告 |

### 2.2 默认行为

- **未指定标志**：检测二进制文件，输出警告但继续显示内容
- **与 grep 保持一致**：`Binary file <filename> matches`

---

## 三、代码修改方案

### 3.1 CLI 层修改

**文件**: `internal/cli/cat.go`

#### 修改 1：添加 flag 定义

**位置**: `init()` 函数中，其他 flag 定义之后

```go
var (
    // 原有 flags...
    catShowLineNum  *qflag.BoolFlag
    catShowNonBlank *qflag.BoolFlag
    // ...
    catQuiet        *qflag.BoolFlag

    // 新增：二进制文件处理
    catText         *qflag.BoolFlag // -a, --text 强制文本模式
    catIgnoreBinary *qflag.BoolFlag // -I, --ignore-binary 忽略二进制文件
)
```

#### 修改 2：初始化 flag

**位置**: `init()` 函数中

```go
func init() {
    CatCmd = qflag.NewCmd("cat", "", qflag.ExitOnError)

    // 原有 flags...
    catQuiet = CatCmd.Bool("quiet", "q", "静默模式 (不显示错误)", false)

    // 新增：二进制文件处理
    catText = CatCmd.Bool("text", "a", "强制将二进制文件视为文本处理", false)
    catIgnoreBinary = CatCmd.Bool("ignore-binary", "I", "完全忽略二进制文件，不输出提示", false)

    // ...
}
```

#### 修改 3：更新 CmdOpts

**位置**: `cmdOpts` 结构体

```go
cmdOpts := &qflag.CmdOpts{
    Desc:        "显示文件内容",
    UseChinese:  true,
    UsageSyntax: fmt.Sprintf("%s cat [options] <file...>", qflag.Root.Name()),
    Examples: map[string]string{
        "显示文件内容":       fmt.Sprintf("%s cat file.txt", qflag.Root.Name()),
        "显示行号":         fmt.Sprintf("%s cat -n file.txt", qflag.Root.Name()),
        "显示前10行":       fmt.Sprintf("%s cat -u 10 file.txt", qflag.Root.Name()),
        "显示后5行":        fmt.Sprintf("%s cat -d 5 file.txt", qflag.Root.Name()),
        "显示换行符类型":     fmt.Sprintf("%s cat -N file.txt", qflag.Root.Name()),
        "静默跳过二进制文件":  fmt.Sprintf("%s cat -I file.bin", qflag.Root.Name()),
        "强制显示二进制内容":  fmt.Sprintf("%s cat -a file.bin", qflag.Root.Name()),
    },
    Notes: []string{
        "换行符检测仅支持 Windows(CRLF) 和 Unix(LF) 格式，不支持旧版 Mac(CR) 格式",
        "默认跳过二进制文件并输出提示，使用 -a 强制处理，使用 -I 静默跳过",
    },
}
```

#### 修改 4：更新 runCat 函数

**位置**: `runCat` 函数

```go
func runCat(cmd qflag.Command) error {
    config := cat.CatConfig{
        Targets:      cmd.Args(),
        ShowLineNum:  catShowLineNum.Get(),
        ShowNonBlank: catShowNonBlank.Get(),
        ShowEnd:      catShowEnd.Get(),
        ShowTabs:     catShowTabs.Get(),
        ShowAll:      catShowAll.Get(),
        ShowNewline:  catShowNewline.Get(),
        HeadLines:    catHeadLines.Get(),
        TailLines:    catTailLines.Get(),
        Quiet:        catQuiet.Get(),
        Text:         catText.Get(),           // 新增
        IgnoreBinary: catIgnoreBinary.Get(),   // 新增
    }

    return cat.CatCmdMain(config)
}
```

---

### 3.2 Commands 层修改

**文件**: `internal/commands/cat/cmd_cat.go`

#### 修改 1：更新导入

```go
import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"

    "gitee.com/MM-Q/fck/internal/utils"  // 新增：导入 utils 包
)
```

#### 修改 2：更新 CatConfig 结构体

```go
// CatConfig cat 命令配置
type CatConfig struct {
    // CLI 参数
    Targets      []string // 目标文件列表
    ShowLineNum  bool     // -n 显示所有行号
    ShowNonBlank bool     // -b 显示非空行行号
    ShowEnd      bool     // -E 显示行尾$
    ShowTabs     bool     // -T 显示制表符为^I
    ShowAll      bool     // -A 等价于 -ET
    ShowNewline  bool     // -N 显示换行符类型
    HeadLines    int      // --head 显示前N行 (0表示全部)
    TailLines    int      // --tail 显示后N行 (0表示全部)
    Quiet        bool     // -q 静默模式 (不显示错误信息)
    Text         bool     // -a, --text 强制将二进制文件视为文本处理
    IgnoreBinary bool     // -I, --ignore-binary 完全忽略二进制文件

    // 运行时
    LineCounter int // 行号计数器
}
```

#### 修改 3：更新 processFile 函数

```go
// processFile 处理单个文件
//
// 参数:
//   - path: 文件路径
//   - config: 命令配置
//
// 返回:
//   - error: 处理错误 (如果有)
func processFile(path string, config *CatConfig) error {
    // 打开文件
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open file %s: %w", path, err)
    }
    defer func() {
        _ = file.Close()
    }()

    // 获取文件信息 (用于判断是否是目录)
    info, err := file.Stat()
    if err != nil {
        return fmt.Errorf("failed to get file info %s: %w", path, err)
    }
    if info.IsDir() {
        return fmt.Errorf("%s is a directory", path)
    }

    // 二进制文件检测（除非强制文本模式）
    if !config.Text {
        isBinary, err := utils.IsBinaryFile(file)
        if err != nil {
            // 检测失败，非静默模式下输出警告
            if !config.Quiet {
                fmt.Fprintf(os.Stderr, "Warning: cannot detect file type for %s: %v\n", path, err)
            }
        } else if isBinary {
            // 是二进制文件
            if config.IgnoreBinary {
                // 静默跳过
                return nil
            }
            // 输出提示（与 grep 保持一致）
            if !config.Quiet {
                fmt.Fprintf(os.Stderr, "Binary file %s matches\n", path)
            }
        }
    }

    // 根据 head/tail 选项处理
    if config.HeadLines > 0 {
        return processHead(file, config)
    }

    if config.TailLines > 0 {
        return processTail(file, config)
    }

    // 普通处理：使用 bufio.Reader 逐行读取
    reader := bufio.NewReader(file)
    for {
        line, newline, err := readLine(reader)
        if err != nil && err != io.EOF {
            return err
        }

        processLine(line, newline, config)

        if err == io.EOF {
            break
        }
    }

    return nil
}
```

---

## 四、边界情况处理

### 4.1 特殊文件类型

| 文件类型 | 处理方式 | 说明 |
|----------|----------|------|
| 普通二进制文件 | 按标志处理 | 检测空字符 |
| 管道/设备文件 | 视为文本 | IsBinaryFile 返回 false |
| 空文件 | 视为文本 | 无内容可检测 |
| 符号链接 | 跟随链接 | 打开后检测目标文件 |

### 4.2 错误处理

```go
// 检测失败时的处理
isBinary, err := utils.IsBinaryFile(file)
if err != nil {
    if !config.Quiet {
        // 输出警告但继续处理
        fmt.Fprintf(os.Stderr, "Warning: cannot detect file type for %s: %v\n", path, err)
    }
    // 继续处理文件（保守策略）
}
```

### 4.3 多文件处理

```bash
# 每个文件独立检测
$ fck cat -I file1.txt file2.bin file3.txt
# file1.txt: 正常显示
# file2.bin: 跳过（无输出）
# file3.txt: 正常显示
```

---

## 五、测试方案

### 5.1 单元测试

```go
// cmd_cat_test.go

func TestProcessFileBinaryDetection(t *testing.T) {
    // 创建临时二进制文件
    binaryFile := createTempFile([]byte{0x00, 0x01, 0x02, 0x03})
    defer os.Remove(binaryFile)

    tests := []struct {
        name         string
        config       CatConfig
        expectOutput bool
        expectWarn   bool
    }{
        {
            name: "默认行为-输出警告",
            config: CatConfig{
                Targets:      []string{binaryFile},
                Text:         false,
                IgnoreBinary: false,
                Quiet:        false,
            },
            expectOutput: true,
            expectWarn:   true,
        },
        {
            name: "-I 标志-跳过二进制",
            config: CatConfig{
                Targets:      []string{binaryFile},
                Text:         false,
                IgnoreBinary: true,
                Quiet:        false,
            },
            expectOutput: false,
            expectWarn:   false,
        },
        {
            name: "-a 标志-强制文本",
            config: CatConfig{
                Targets:      []string{binaryFile},
                Text:         true,
                IgnoreBinary: false,
                Quiet:        false,
            },
            expectOutput: true,
            expectWarn:   false,
        },
        {
            name: "静默模式-不输出警告",
            config: CatConfig{
                Targets:      []string{binaryFile},
                Text:         false,
                IgnoreBinary: false,
                Quiet:        true,
            },
            expectOutput: true,
            expectWarn:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 捕获 stdout 和 stderr
            stdout := captureStdout(func() {
                stderr := captureStderr(func() {
                    err := CatCmdMain(tt.config)
                    if err != nil {
                        t.Logf("Error: %v", err)
                    }
                })
                
                if tt.expectWarn && stderr == "" {
                    t.Error("expected warning in stderr, got none")
                }
                if !tt.expectWarn && stderr != "" {
                    t.Errorf("expected no warning, got: %s", stderr)
                }
            })

            if tt.expectOutput && stdout == "" {
                t.Error("expected output in stdout, got none")
            }
            if !tt.expectOutput && stdout != "" {
                t.Errorf("expected no output, got: %s", stdout)
            }
        })
    }
}
```

### 5.2 手动测试命令

```bash
# 1. 准备测试文件
echo "hello world" > /tmp/text.txt
printf '\x00\x01\x02\x03' > /tmp/binary.bin

# 2. 测试默认行为（输出警告）
fck cat /tmp/binary.bin
# 预期: Binary file /tmp/binary.bin matches
#       [乱码]

# 3. 测试 -I 标志（跳过）
fck cat -I /tmp/binary.bin
# 预期: 无输出

# 4. 测试 -a 标志（强制文本）
fck cat -a /tmp/binary.bin
# 预期: [乱码]，无警告

# 5. 测试混合文件
fck cat /tmp/text.txt /tmp/binary.bin
# 预期: hello world
#       Binary file /tmp/binary.bin matches
#       [乱码]

# 6. 测试 -I 混合文件
fck cat -I /tmp/text.txt /tmp/binary.bin
# 预期: hello world

# 7. 测试静默模式
fck cat -q /tmp/binary.bin
# 预期: [乱码]，无警告

# 8. 测试管道/设备文件（应正常处理）
fck cat -I /dev/null
# 预期: 无输出（空文件）
```

---

## 六、与 grep 的对比

### 6.1 行为一致性

| 场景 | grep 行为 | cat 行为（修改后） |
|------|-----------|-------------------|
| 默认 | `Binary file matches` + 跳过内容 | `Binary file matches` + 继续显示 |
| `-I` | 静默跳过 | 静默跳过 |
| `-a` | 强制处理 | 强制处理 |

**差异说明**：cat 和 grep 的默认行为略有不同
- grep：检测到二进制时**跳过内容**（因为是搜索工具）
- cat：检测到二进制时**继续显示**（因为是查看工具，用户可能想看）

### 6.2 提示信息格式

```
grep: Binary file <filename> matches
cat:  Binary file <filename> matches
```

保持一致，便于用户理解。

---

## 七、实施步骤

1. **修改 CLI 层** (`internal/cli/cat.go`)
   - 添加 flag 定义
   - 更新 CmdOpts
   - 更新 runCat 函数

2. **修改 Commands 层** (`internal/commands/cat/cmd_cat.go`)
   - 添加 utils 导入
   - 更新 CatConfig 结构体
   - 更新 processFile 函数

3. **编译验证**
   ```bash
   go build ./...
   ```

4. **运行测试**
   ```bash
   go test ./internal/commands/cat/...
   ```

5. **手动验证**
   ```bash
   # 创建测试文件
   echo "text" > /tmp/test.txt
   printf '\x00\x01' > /tmp/test.bin

   # 测试各种场景
   go run cmd/main.go cat /tmp/test.bin
   go run cmd/main.go cat -I /tmp/test.bin
   go run cmd/main.go cat -a /tmp/test.bin
   ```

---

## 八、总结

**方案优势**:
- ✅ 与 grep 命令保持一致的用户体验
- ✅ 向后兼容（默认行为只是增加警告）
- ✅ 实现简单，复用现有 `IsBinaryFile` 函数
- ✅ 提供灵活的控制选项

**实施成本**:
- 低（约 2 个文件，50 行代码）

**建议**: 按此方案实施，增强 cat 命令的健壮性。

---

**方案完成，等待确认后实施。**
