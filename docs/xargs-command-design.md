# xargs 子命令实施方案

## 1. 功能需求分析

### 核心功能
`xargs` 命令用于从标准输入读取数据，将其作为参数传递给其他命令执行。类似于 Unix/Linux 的 xargs 工具。

### 主要功能点
1. **从 stdin 读取输入** - 支持多行输入
2. **参数传递** - 将输入作为参数传递给指定命令
3. **批量执行** - 支持一次传递多个参数（-n 选项）
4. **并行执行** - 支持并发执行（-P 选项）
5. **占位符替换** - 支持 `{}` 占位符自定义参数位置
6. **空值处理** - 控制是否处理空行（-r 选项）
7. **限制执行次数** - 限制最大执行次数（-L 选项）
8. **显示命令** - 打印要执行的命令而不执行（-t 选项）

---

## 2. 命令设计

### 命令名称
- **长名称**: `xargs`
- **短名称**: `x`

### 使用场景示例

```bash
# 基本用法：将输入行作为参数传递给命令
echo "file1.txt" | fck xargs cat

# 批量处理：每次传递2个参数
ls *.txt | fck xargs -n 2 cat

# 并行执行：同时运行3个进程
cat urls.txt | fck xargs -P 4 curl -O

# 使用占位符：自定义参数位置
echo "file.txt" | fck xargs -I {} cp {} {}.bak

# 限制执行次数：最多执行10次
seq 1 100 | fck xargs -L 10 echo

# 显示命令但不执行
echo "test" | fck xargs -t echo

# 忽略空行
echo -e "a\n\nb" | fck xargs -r echo
```

---

## 3. 标志设计

| 标志 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-n, --max-args` | int | 1 | 每次执行传递的最大参数数量 |
| `-P, --max-procs` | int | 1 | 最大并行进程数 |
| `-I, --replace` | string | `{}` | 占位符（如 `{}`），用于替换命令中的参数位置 |
| `-r, --no-run-if-empty` | bool | false | 如果输入为空，不执行命令 |
| `-L, --max-lines` | int | 0 | 最大执行次数（0表示无限制） |
| `-t, --verbose` | bool | false | 打印要执行的命令 |
| `-d, --delimiter` | string | "\n" | 输入分隔符 |
| `-0, --null` | bool | false | 使用空字符（\0）作为分隔符 |
| `-e, --eof` | string | "" | 指定结束标志，遇到此行停止读取 |

---

## 4. 目录结构设计

根据 FCK 项目规范，采用**分层架构**：

```
internal/
├── commands/
│   └── xargs/                    # 业务逻辑层（新增）
│       └── cmd_xargs.go          # 核心逻辑实现
└── cli/
    ├── root.go                   # 根命令
    ├── gm.go                     # gm 命令
    ├── xargs.go                  # xargs CLI 定义（新增）
    └── ...
```

### 文件分工

| 文件 | 职责 |
|------|------|
| `internal/commands/xargs/cmd_xargs.go` | 核心逻辑（读取输入、分批、执行命令） |
| `internal/cli/xargs.go` | CLI 定义（标志、配置、调用核心逻辑） |

---

## 5. 业务逻辑文件设计

### 文件位置
`internal/commands/xargs/cmd_xargs.go`

### 代码结构

```go
package xargs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"gitee.com/MM-Q/shellx/shx"
)

// XargsConfig 配置结构体
type XargsConfig struct {
	MaxArgs      int
	MaxProcs     int
	Replace      string
	NoRunIfEmpty bool
	MaxLines     int
	Verbose      bool
	Delimiter    string
	UseNull      bool
	Eof          string
	TargetCmd    string
	InitialArgs  []string
}

// XargsStats 执行统计（可选）
type XargsStats struct {
	Executed int
	Errors   int
}

// XargsCmdMain 主函数
func XargsCmdMain(config XargsConfig, stdin io.Reader) error {
	// 根据 -0 选项设置分隔符
	delimiter := config.Delimiter
	if config.UseNull {
		delimiter = "\x00"
	}

	// 读取标准输入
	inputs, err := readInputs(stdin, delimiter, config.Eof)
	if err != nil {
		return fmt.Errorf("读取输入失败: %w", err)
	}

	// 检查空输入
	if len(inputs) == 0 && config.NoRunIfEmpty {
		return nil
	}

	// 限制执行次数
	if config.MaxLines > 0 && len(inputs) > config.MaxLines {
		inputs = inputs[:config.MaxLines]
	}

	// 执行命令
	return executeXargs(inputs, config)
}

// readInputs 从 reader 读取输入，按分隔符分割
func readInputs(r io.Reader, delimiter string, eof string) ([]string, error) {
	var inputs []string
	scanner := bufio.NewScanner(r)

	// 自定义分隔函数
	if delimiter != "\n" {
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := strings.Index(string(data), delimiter); i >= 0 {
				return i + len(delimiter), data[0:i], nil
			}
			if atEOF {
				return len(data), data, nil
			}
			return 0, nil, nil
		})
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行（除非使用 -0 选项）
		if line == "" && delimiter != "\x00" {
			continue
		}

		// 检查结束标志
		if eof != "" && line == eof {
			break
		}

		inputs = append(inputs, line)
	}

	return inputs, scanner.Err()
}

// executeXargs 执行 xargs 核心逻辑
func executeXargs(inputs []string, config XargsConfig) error {
	// 批量分组
	batches := batchInputs(inputs, config.MaxArgs)

	// 串行执行
	if config.MaxProcs == 1 {
		for _, batch := range batches {
			if err := executeBatch(batch, config); err != nil {
				return err
			}
		}
		return nil
	}

	// 并行执行
	return executeParallel(batches, config)
}

// batchInputs 将输入分批
func batchInputs(inputs []string, batchSize int) [][]string {
	if batchSize <= 0 {
		batchSize = 1
	}

	var batches [][]string
	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		batches = append(batches, inputs[i:end])
	}

	return batches
}

// executeBatch 执行一批命令
func executeBatch(batch []string, config XargsConfig) error {
	// 构建参数
	args := make([]string, len(config.InitialArgs))
	copy(args, config.InitialArgs)

	if config.Replace != "" {
		// 使用占位符替换
		replacement := strings.Join(batch, " ")
		for i, arg := range args {
			args[i] = strings.ReplaceAll(arg, config.Replace, replacement)
		}
		// 如果参数列表为空，将替换后的占位符作为参数
		if len(args) == 0 {
			args = append(args, replacement)
		}
	} else {
		// 追加到参数列表
		args = append(args, batch...)
	}

	// 打印命令（如果启用 verbose）
	if config.Verbose {
		fmt.Printf("+ %s %s\n", config.TargetCmd, strings.Join(args, " "))
	}

	// 使用 shx 执行命令
	cmd := shx.NewArgs(config.TargetCmd, args...).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithStdin(os.Stdin)

	return cmd.Exec()
}

// executeParallel 并行执行
func executeParallel(batches [][]string, config XargsConfig) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxProcs)
	errChan := make(chan error, len(batches))

	for _, batch := range batches {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(b []string) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			if err := executeBatch(b, config); err != nil {
				errChan <- err
			}
		}(batch)
	}

	wg.Wait()
	close(errChan)

	// 检查错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
```

---

## 6. CLI 定义文件设计

### 文件位置
`internal/cli/xargs.go`

### 代码结构

```go
package cli

import (
	"fmt"
	"os"

	"gitee.com/MM-Q/fck/internal/commands/xargs"
	"gitee.com/MM-Q/qflag"
)

// ============================================
// 1. 全局命令变量
// ============================================
var (
	XargsCmd *qflag.Cmd
)

// ============================================
// 2. 全局标志变量
// ============================================
var (
	xargsMaxArgs      *qflag.IntFlag
	xargsMaxProcs     *qflag.IntFlag
	xargsReplace      *qflag.StringFlag
	xargsNoRunEmpty   *qflag.BoolFlag
	xargsMaxLines     *qflag.IntFlag
	xargsVerbose      *qflag.BoolFlag
	xargsDelimiter    *qflag.StringFlag
	xargsNull         *qflag.BoolFlag
	xargsEof          *qflag.StringFlag
)

// ============================================
// 3. init() 初始化
// ============================================
func init() {
	XargsCmd = qflag.NewCmd("xargs", "x", qflag.ExitOnError)

	// 定义标志
	xargsMaxArgs = XargsCmd.Int("max-args", "n", "每次执行传递的最大参数数量", 1)
	xargsMaxProcs = XargsCmd.Int("max-procs", "P", "最大并行进程数", 1)
	xargsReplace = XargsCmd.String("replace", "I", "占位符（用于替换命令中的参数位置）", "{}")
	xargsNoRunEmpty = XargsCmd.Bool("no-run-if-empty", "r", "如果输入为空，不执行命令", false)
	xargsMaxLines = XargsCmd.Int("max-lines", "L", "最大执行次数（0表示无限制）", 0)
	xargsVerbose = XargsCmd.Bool("verbose", "t", "打印要执行的命令", false)
	xargsDelimiter = XargsCmd.String("delimiter", "d", "输入分隔符", "\n")
	xargsNull = XargsCmd.Bool("null", "0", "使用空字符（\\0）作为分隔符", false)
	xargsEof = XargsCmd.String("eof", "e", "指定结束标志，遇到此行停止读取", "")

	// 应用命令配置
	cmdOpts := &qflag.CmdOpts{
		Desc:        "从标准输入读取数据并执行命令",
		UsageSyntax: fmt.Sprintf("%s xargs [选项] 命令 [初始参数]", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"从标准输入读取数据，将每行作为参数传递给指定命令",
			"如果不指定命令，默认使用 echo",
			"默认使用 {} 作为占位符，可通过 -I 选项自定义",
		},
		Examples: map[string]string{
			"基本用法":         `echo "file.txt" | fck xargs cat`,
			"批量处理":         `ls *.txt | fck xargs -n 2 cat`,
			"并行执行":         `cat urls.txt | fck xargs -P 4 curl -O`,
			"使用占位符":       `echo "file.txt" | fck xargs -I {} cp {} {}.bak`,
		},
	}

	if err := XargsCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	XargsCmd.SetRun(runXargs)
}

// ============================================
// 4. run() 函数
// ============================================
func runXargs(cmd qflag.Command) error {
	// 获取标志值
	config := xargs.XargsConfig{
		MaxArgs:      xargsMaxArgs.Get(),
		MaxProcs:     xargsMaxProcs.Get(),
		Replace:      xargsReplace.Get(),
		NoRunIfEmpty: xargsNoRunEmpty.Get(),
		MaxLines:     xargsMaxLines.Get(),
		Verbose:      xargsVerbose.Get(),
		Delimiter:    xargsDelimiter.Get(),
		UseNull:      xargsNull.Get(),
		Eof:          xargsEof.Get(),
	}

	// 获取要执行的命令和初始参数
	args := cmd.Args()
	if len(args) == 0 {
		config.TargetCmd = "echo" // 默认命令
	} else {
		config.TargetCmd = args[0]
		config.InitialArgs = args[1:]
	}

	// 调用业务逻辑
	return xargs.XargsCmdMain(config, os.Stdin)
}
```

---

## 7. 注册到根命令

在 `internal/cli/root.go` 的 `SubCmds` 中添加：

```go
SubCmds: []qflag.Command{
	// 现有命令...
	GmCmd,
	XargsCmd,  // 新增 xargs 命令
},
```

---

## 8. 测试用例设计

| 测试场景 | 输入 | 命令 | 预期结果 |
|---------|------|------|---------|
| 基本用法 | `echo "hello"` | `xargs echo` | 输出 "hello" |
| 批量处理 | `seq 1 5` | `xargs -n 2 echo` | 分3次输出: "1 2", "3 4", "5" |
| 并行执行 | `seq 1 4` | `xargs -P 2 -n 1 sleep` | 并行执行4个 sleep |
| 占位符 | `echo "file.txt"` | `xargs -I {} echo {}.bak` | 输出 "file.txt.bak" |
| 空输入处理 | `echo ""` | `xargs -r echo` | 不执行命令 |
| 最大执行次数 | `seq 1 100` | `xargs -L 5 echo` | 只执行5次 |
| 显示命令 | `echo "test"` | `xargs -t echo` | 先打印 "+ echo test" |
| 自定义分隔符 | `echo "a,b,c"` | `xargs -d ',' echo` | 输出 "a b c" |
| 空字符分隔 | `printf "a\0b\0c"` | `xargs -0 echo` | 输出 "a b c" |

---

## 9. 边缘案例考虑

1. **命令不存在** - 返回友好错误信息
2. **命令执行失败** - 返回退出码和错误信息
3. **输入包含特殊字符** - 正确处理引号和空格
4. **并行执行时的资源竞争** - 使用信号量控制并发
5. **信号处理** - 支持 Ctrl+C 中断
6. **大输入处理** - 流式读取，避免内存溢出

---

## 10. 实施步骤

### 步骤1：创建业务逻辑文件
```bash
mkdir -p internal/commands/xargs
touch internal/commands/xargs/cmd_xargs.go
```

### 步骤2：创建 CLI 定义文件
```bash
touch internal/cli/xargs.go
```

### 步骤3：注册命令
编辑 `internal/cli/root.go`，在 `SubCmds` 中添加 `XargsCmd`

### 步骤4：编译测试
```bash
go build ./...
```

### 步骤5：功能测试
使用文档中的测试用例进行验证

---

## 11. 实施时间估算

| 任务 | 预计时间 |
|------|---------|
| 创建目录和文件 | 5分钟 |
| 实现业务逻辑（cmd_xargs.go） | 40分钟 |
| 实现 CLI 定义（xargs.go） | 15分钟 |
| 注册命令 | 5分钟 |
| 编译测试 | 10分钟 |
| 功能测试 | 20分钟 |
| **总计** | **约1.5-2小时** |

---

## 总结

这个 `xargs` 命令设计方案：
- ✅ 符合 FCK 项目规范（分层架构：commands + cli）
- ✅ 业务逻辑与 CLI 定义分离
- ✅ 功能完整（支持所有常用选项）
- ✅ 易于使用（清晰的帮助文档和示例）
- ✅ 性能优化（支持并行执行）
- ✅ 错误处理完善
