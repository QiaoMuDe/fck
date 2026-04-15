# xargs Shell 模式切换设计方案

## 概述

为 xargs 命令新增 `--shell` 标志，支持两种执行模式：
- **默认模式**：直接执行（安全，无注入风险）
- **Shell 模式**：通过 shell 执行（支持管道、重定向等高级特性）

## 新增标志

| 标志 | 长格式 | 说明 | 默认值 |
|------|--------|------|--------|
| 无 | `--shell` | 通过 shell 执行命令 | `false` |

## 行为对比

| 场景 | 默认模式（直接执行） | Shell 模式（`--shell`） |
|------|----------------------|------------------------|
| 执行方式 | `exec.Command()` 直接调用 | `sh -c "command"` 通过 shell |
| 特殊字符 | 安全处理，无注入风险 | 按 shell 规则解析 |
| 管道支持 | ❌ 不支持 | ✅ `cat {} \| grep pattern` |
| 重定向支持 | ❌ 不支持 | ✅ `cat {} > output.txt` |
| 多命令组合 | ❌ 不支持 | ✅ `cp {} {}.bak && echo done` |
| 通配符展开 | ❌ 不支持 | ✅ `cat *.txt` |

## 修改文件清单

### 1. internal/cli/xargs.go

**新增 flag：**
```go
var (
    // 现有 flags...
    xargsShell *qflag.BoolFlag  // --shell 通过 shell 执行
)

func init() {
    // ...
    xargsShell = XargsCmd.Bool("shell", "", "通过 shell 执行命令（支持管道、重定向等）", false)
    
    cmdOpts := &qflag.CmdOpts{
        // ...
        Examples: map[string]string{
            // 现有示例...
            "使用管道":      `echo "file.txt" | fck xargs --shell "cat {} | grep pattern"`,
            "使用重定向":    `echo "file.txt" | fck xargs --shell "cat {} > output.txt"`,
        },
        Notes: []string{
            // 现有 notes...
            "默认使用直接执行模式，避免 shell 注入风险",
            "如需使用管道、重定向等 shell 特性，请添加 --shell 标志",
            "--shell 模式下请注意输入安全性",
        },
    }
}

func runXargs(cmd qflag.Command) error {
    config := xargs.XargsConfig{
        // 现有字段...
        Shell: xargsShell.Get(),
    }
    return xargs.XargsCmdMain(config)
}
```

### 2. internal/commands/xargs/cmd_xargs.go

#### 2.1 新增配置字段
```go
type XargsConfig struct {
    // 现有字段...
    Shell bool // --shell 通过 shell 执行
}
```

#### 2.2 重构 executeBatch 函数
```go
// executeBatch 执行单个批次
func executeBatch(batch []string, config XargsConfig, stats *XargsStats) error {
    stats.Executed++

    // 构建命令
    var cmdStr string
    if config.ReplaceStr != "" || config.ReplaceDelim != "" {
        cmdStr = buildReplaceCommand(batch, config)
    } else {
        cmdStr = buildDefaultCommand(batch, config)
    }

    // 显示命令
    if config.Verbose {
        fmt.Fprintln(os.Stderr, cmdStr)
    }

    // 确认模式
    if config.Interactive {
        fmt.Fprintf(os.Stderr, "执行? (y/n): ")
        var response string
        if _, err := fmt.Scanln(&response); err != nil {
            return fmt.Errorf("读取确认失败: %w", err)
        }
        if response != "y" && response != "Y" {
            return nil
        }
    }

    // 根据模式选择执行方式
    if config.Shell {
        return executeWithShell(cmdStr, stats)
    }
    return executeDirectly(batch, config, stats)
}
```

#### 2.3 直接执行模式（新增）
```go
// executeDirectly 直接执行命令（安全模式）
func executeDirectly(batch []string, config XargsConfig, stats *XargsStats) error {
    // 替换模式：每个参数单独执行
    if config.ReplaceStr != "" || config.ReplaceDelim != "" {
        placeholder := config.ReplaceStr
        if placeholder == "" {
            placeholder = config.ReplaceDelim
        }
        if placeholder == "" {
            placeholder = "{}"
        }

        for _, arg := range batch {
            // 构建参数列表
            var args []string
            
            // 处理 CommandArgs 中的占位符
            for _, fixedArg := range config.CommandArgs {
                replacedArg := strings.ReplaceAll(fixedArg, placeholder, arg)
                args = append(args, replacedArg)
            }

            cmd := exec.Command(config.Command, args...)
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr

            if err := cmd.Run(); err != nil {
                stats.Failed++
                return fmt.Errorf("执行失败: %w", err)
            }
        }
        return nil
    }

    // 默认模式：所有参数追加到命令后
    args := append(config.CommandArgs, batch...)
    cmd := exec.Command(config.Command, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        stats.Failed++
        return fmt.Errorf("执行失败: %w", err)
    }
    return nil
}
```

#### 2.4 Shell 执行模式（现有逻辑提取）
```go
// executeWithShell 通过 shell 执行命令（兼容模式）
func executeWithShell(cmdStr string, stats *XargsStats) error {
    if err := shx.RunToTerminal(cmdStr); err != nil {
        stats.Failed++
        return fmt.Errorf("执行失败: %w", err)
    }
    return nil
}

// buildReplaceCommand 构建替换模式的命令字符串（供 shell 模式使用）
func buildReplaceCommand(batch []string, config XargsConfig) string {
    placeholder := config.ReplaceStr
    if placeholder == "" {
        placeholder = config.ReplaceDelim
    }
    if placeholder == "" {
        placeholder = "{}"
    }

    var cmds []string
    for _, arg := range batch {
        cmd := config.Command
        cmd = strings.ReplaceAll(cmd, placeholder, arg)
        
        for _, fixedArg := range config.CommandArgs {
            replacedArg := strings.ReplaceAll(fixedArg, placeholder, arg)
            cmd += " " + replacedArg
        }
        cmds = append(cmds, cmd)
    }
    return strings.Join(cmds, " && ")
}

// buildDefaultCommand 构建默认模式的命令字符串（供 shell 模式使用）
func buildDefaultCommand(batch []string, config XargsConfig) string {
    cmd := config.Command
    for _, arg := range config.CommandArgs {
        cmd += " " + arg
    }
    for _, arg := range batch {
        cmd += " " + arg
    }
    return cmd
}
```

#### 2.5 新增导入
```go
import (
    "bufio"
    "fmt"
    "os"
    "os/exec"  // 新增
    "strings"
    "sync"

    "gitee.com/MM-Q/shellx/shx"
)
```

## 使用示例

### 默认模式（安全，推荐）
```bash
# 基本用法
echo "file.txt" | fck xargs cat

# 批量处理
echo "file1.txt file2.txt" | fck xargs cat

# 占位符替换
echo "file.txt" | fck xargs -i cat {}
ls *.txt | fck xargs -i mv {} {}.bak
```

### Shell 模式（高级特性）
```bash
# 使用管道
echo "file.txt" | fck xargs --shell "cat {} | grep pattern"

# 使用重定向
echo "file.txt" | fck xargs --shell "cat {} > output.txt"

# 多命令组合
echo "file.txt" | fck xargs --shell "cp {} {}.bak && echo done"

# 通配符展开
echo "/path" | fck xargs --shell "ls {}/*.txt"
```

## 安全考虑

### 默认模式的优势
```bash
# 恶意输入在默认模式下是安全的
echo "file; rm -rf /" | fck xargs cat
# 执行: cat "file; rm -rf /"
# 结果: 尝试打开名为 "file; rm -rf /" 的文件（失败，但无危害）
```

### Shell 模式的注意事项
```bash
# 恶意输入在 shell 模式下有危险
echo "file; rm -rf /" | fck xargs --shell cat
# 执行: sh -c "cat file; rm -rf /"
# 结果: 执行了两个命令！

# 建议：使用 --shell 时确保输入可信
# 或使用 -0 模式处理文件名
echo -e "file\x00file2" | fck xargs -0 --shell "cat {}"
```

## 测试用例

### 默认模式测试
```bash
# 1. 基本执行
echo "file.txt" | fck xargs cat

# 2. 多个参数
echo "a b c" | fck xargs echo

# 3. 占位符替换
echo "file.txt" | fck xargs -i echo "{}"

# 4. 安全测试（特殊字符）
echo "file; echo hacked" | fck xargs echo
# 预期输出: file; echo hacked（作为一个整体）
```

### Shell 模式测试
```bash
# 1. 管道
echo "file.txt" | fck xargs --shell "cat {} | head -1"

# 2. 重定向
echo "file.txt" | fck xargs --shell "cat {} > /tmp/output"

# 3. 多命令
echo "file.txt" | fck xargs --shell "echo start && cat {} && echo end"

# 4. 特殊字符（按 shell 解析）
echo "file" | fck xargs --shell "echo {} && echo second"
# 预期输出: file
#           second
```

## 兼容性说明

### 与 Linux xargs 的兼容性

| 特性 | Linux xargs | fck xargs（修改后） |
|------|-------------|---------------------|
| 默认执行方式 | 直接执行 | ✅ 直接执行 |
| 通过 shell 执行 | 需显式 `sh -c` | ✅ `--shell` 标志 |
| 占位符替换 | `-I {}` | ✅ `-i` / `-I` |
| 管道支持 | 需 `sh -c` | ✅ `--shell` 直接支持 |

### 破坏性变更

**这是一个破坏性变更**：
- 修改前：默认通过 shell 执行
- 修改后：默认直接执行

**影响**：
- 依赖 shell 特性的命令需要添加 `--shell`
- 例如：`echo "file" | fck xargs 