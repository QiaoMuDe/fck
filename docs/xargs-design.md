# xargs 命令设计方案

## 1. 功能定位

从标准输入读取参数，批量执行指定命令。类似 Linux xargs，支持多参数组合、并行执行、参数替换等。

## 2. 核心功能

| 功能 | 说明 |
|------|------|
| 基础执行 | `echo a b c \| xargs echo` → `echo a b c` |
| 参数替换 | `-i {}` 占位符替换 |
| 批量大小 | `-n N` 每批 N 个参数 |
| 并行执行 | `-P N` N 个并行进程 |
| 参数限制 | `-s N` 命令长度限制 |
| 空输入处理 | `-r` 空输入不执行 |
| 文件输入 | `-a file` 从文件读取 |
| 分隔符 | `-d delim` 自定义分隔符 |
| 确认模式 | `-p` 执行前确认 |
| 最大执行 | `-L N` 最多执行 N 行 |

## 3. 标志设计

```go
// 基础标志
xargsDelimiter    // -d, --delimiter  输入分隔符（默认空格/换行）
xargsNull         // -0, --null       使用 \0 作为分隔符
xargsArgFile      // -a, --arg-file   从文件读取参数

// 批量控制
xargsMaxArgs      // -n, --max-args   每批最大参数个数
xargsMaxLines     // -L, --max-lines  每批最大行数
xargsMaxChars     // -s, --max-chars  命令最大长度

// 执行控制
xargsReplaceStr   // -i, --replace    启用占位符替换模式（默认占位符 {}）
xargsReplaceDelim // -I, --replace-delim  自定义占位符字符串
xargsMaxProcs     // -P, --max-procs  并行进程数
xargsNoRunIfEmpty // -r, --no-run-if-empty  空输入不执行
xargsInteractive  // -p, --interactive 执行前确认
xargsVerbose      // -t, --verbose    打印执行的命令
xargsExitOnError  // -e, --exit-on-error  出错立即停止

// 特殊处理
xargsOpenTty      // -o, --open-tty   在子进程打开 tty
xargsShowLimits   // --show-limits    显示系统限制
```

## 4. 使用示例

```bash
# 基础用法（所有参数追加到命令后，执行一次）
echo a b c | fck xargs echo
# 实际执行: echo a b c

# 参数替换（-i 模式，每个参数单独替换占位符，执行多次）
echo a b c | fck xargs -i {} echo "file: {}"
# 实际执行: echo "file: a"
# 实际执行: echo "file: b"
# 实际执行: echo "file: c"

# 参数替换（文件重命名示例）
find . -name "*.log" | fck xargs -i {} mv {} {}.bak

# 自定义占位符
echo a | fck xargs -i -I %% echo %%

# 批量处理（每批2个）
echo a b c d | fck xargs -n 2 echo

# 并行执行（4进程）
cat urls.txt | fck xargs -P 4 curl -O

# 从文件读取
fck xargs -a files.txt rm

# 空分隔符（处理含空格的文件名）
find . -print0 | fck xargs -0 rm

# 执行前确认
echo file1 | fck xargs -p rm

# 限制命令长度
fck xargs -s 1024 echo

# 组合使用
find . -name "*.tmp" | fck xargs -0 -i {} -P 4 mv {} /backup/
```

## 5. 核心逻辑设计

```go
type XargsConfig struct {
    // 输入配置
    Delimiter    string   // 分隔符
    NullMode     bool     // \0 分隔模式
    ArgFile      string   // 参数文件
    
    // 批量控制
    MaxArgs      int      // 每批最大参数
    MaxLines     int      // 每批最大行数
    MaxChars     int      // 命令最大长度
    
    // 执行配置
    ReplaceStr   string   // 占位符（如 {}）
    ReplaceDelim string   // 自定义占位符字符串
    MaxProcs     int      // 并行数
    NoRunIfEmpty bool     // 空输入不执行
    Interactive  bool     // 确认模式
    Verbose      bool     // 显示命令
    ExitOnError  bool     // 出错停止
    
    // 目标命令
    Command      string   // 要执行的命令
    CommandArgs  []string // 命令固定参数
}

// 主流程
func XargsCmdMain(config XargsConfig) error {
    // 1. 读取所有参数（标准输入或文件）
    args, err := readArgs(config)
    
    // 2. 按规则分批
    batches := splitBatches(args, config)
    
    // 3. 执行批次
    if config.MaxProcs > 1 {
        return runParallel(batches, config)
    }
    return runSequential(batches, config)
}
```

## 6. 关键实现点

### 6.1 参数分批算法

- 按 `-n` 参数个数分批
- 按 `-s` 命令长度分批
- 按 `-L` 行数分批
- 组合限制条件

### 6.2 并行执行

- 使用 `sync.WaitGroup` + 信号量控制并发
- 收集所有进程退出码
- 支持 `-e` 出错立即停止

### 6.3 参数替换

- 扫描命令参数中的占位符
- 替换为当前批次参数
- 支持多次出现

## 7. 边缘情况处理

| 场景 | 处理 |
|------|------|
| 空输入 | 默认不执行，`-r` 时直接退出 |
| 命令不存在 | 返回错误 |
| 参数过长 | 自动拆分到下一批 |
| 子进程失败 | 收集退出码，`-e` 时立即停止 |
| 信号中断 | 优雅关闭，等待进行中的任务 |

## 8. 与 Linux xargs 的差异

| 特性 | Linux xargs | 本实现 |
|------|-------------|--------|
| 默认分隔符 | 空格/换行/制表 | 相同 |
| `-i` 替换 | 支持 | 支持 |
| `-P` 并行 | 支持 | 支持 |
| `-p` 确认 | 支持 | 支持 |
| `--show-limits` | 支持 | 支持 |
| `-E EOF` | 支持 | 暂不支持（可用 `-L` 替代） |

## 9. 实现文件结构

```
internal/commands/xargs/
└── cmd_xargs.go          # 业务逻辑

internal/cli/
└── xargs.go              # CLI 定义
```

## 10. 依赖库

命令执行使用 `gitee.com/MM-Q/shellx/shx` 库：

```go
import "gitee.com/MM-Q/shellx/shx"

// 执行命令并输出到终端
if err := shx.RunToTerminal(cmdStr); err != nil {
    return fmt.Errorf("failed to execute command: %v", err)
}
```

## 11. 待确认问题

1. 是否需要支持 `-E EOF`（EOF 标记）功能？
2. 并行执行时，输出是否需要加锁保证顺序？
3. 是否需要支持 `--process-slot-var`（进程槽位变量）？
4. Windows 平台是否需要特殊处理（如命令查找）？
