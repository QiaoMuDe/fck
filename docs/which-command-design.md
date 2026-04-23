# which 命令设计方案

## 一、需求概述

在环境中查找指定命令是否存在。存在则正常退出并输出命令的完整路径，不存在则报错退出（错误信息使用英文）。

## 二、核心思路

1. **查找逻辑**：遍历系统的 PATH 环境变量中的每个目录，查找指定的可执行文件
2. **跨平台支持**：
   - Windows：查找 `.exe`, `.bat`, `.cmd`, `.com` 等可执行扩展名
   - Linux/macOS：直接查找文件，检查可执行权限
3. **输出**：找到则输出完整路径，未找到返回英文错误信息

## 三、涉及文件/模块范围

| 文件路径 | 说明 |
|---------|------|
| `internal/commands/which/cmd_which.go` | 业务逻辑文件（新增） |
| `internal/cli/which.go` | CLI 定义文件（新增） |
| `internal/cli/root.go` | 注册命令（修改） |

## 四、关键逻辑设计

### 4.1 配置结构体

```go
// WhichConfig 配置结构体
type WhichConfig struct {
    Commands []string  // 要查找的命令列表（支持多个）
    All      bool      // 显示所有匹配的路径（而不仅是第一个）
    Silent   bool      // 静默模式，不输出，只返回退出码
}
```

### 4.2 查找算法

```
1. 获取 PATH 环境变量，按分隔符分割（Windows: ;, Unix: :）
2. 对于每个要查找的命令：
   a. 如果命令包含路径分隔符，直接检查该路径
   b. 否则遍历 PATH 中的每个目录：
      - Windows: 尝试添加可执行扩展名 (.exe, .bat, .cmd, .com)
      - Unix: 直接检查文件是否存在且有可执行权限
   c. 找到匹配则记录路径
3. 根据配置输出结果或返回错误
```

### 4.3 错误处理

| 场景 | 退出码 | 错误信息（英文） |
|-----|-------|----------------|
| 未提供命令参数 | 1 | `no command specified` |
| 命令未找到 | 1 | `<command>: command not found` |
| 系统错误（如无法读取 PATH） | 2 | `failed to search command: <details>` |

## 五、选项设计

| 长选项 | 短选项 | 类型 | 默认值 | 说明 |
|-------|-------|-----|-------|-----|
| `--all` | `-a` | bool | false | 显示所有匹配的路径 |
| `--silent` | `-s` | bool | false | 静默模式，不输出任何内容 |

## 六、使用示例

```bash
# 查找单个命令
fck which go
# 输出: C:\Program Files\Go\bin\go.exe

# 查找多个命令
fck which go git node
# 输出:
# C:\Program Files\Go\bin\go.exe
# C:\Program Files\Git\bin\git.exe
# C:\Program Files\nodejs\node.exe

# 显示所有匹配（Windows 可能同时匹配 go.exe 和 go.bat）
fck which -a python
# 输出:
# C:\Python39\python.exe
# C:\Users\xxx\AppData\Local\Microsoft\WindowsApps\python.exe

# 静默模式（用于脚本检查）
fck which -s go && echo "go is installed"
```

## 七、技术选型依据

1. **标准库优先**：使用 `os`, `path/filepath`, `strings` 等标准库，无额外依赖
2. **跨平台**：使用 `runtime.GOOS` 判断操作系统，分别处理 Windows 和 Unix-like 系统
3. **错误信息**：按照 Unix `which` 命令惯例，使用英文错误信息

## 八、预估改动量

- 新增文件：2 个
- 修改文件：1 个（root.go 添加命令注册）
- 预估代码行数：约 150-200 行

## 九、边缘案例

1. **空参数**：未提供命令时返回错误
2. **路径分隔符**：命令名包含 `/` 或 `\` 时，视为相对/绝对路径直接检查
3. **大小写敏感**：Windows 不区分大小写，Unix 区分
4. **特殊字符**：命令名包含空格或特殊字符的处理
5. **权限问题**：Unix 系统需要检查可执行权限
6. **空 PATH**：PATH 环境变量为空或不存在时的处理

## 十、测试用例建议

| 测试场景 | 输入 | 预期结果 |
|---------|-----|---------|
| 查找存在的命令 | `which go` | 输出完整路径，退出码 0 |
| 查找不存在的命令 | `which notexistcmd` | 错误信息，退出码 1 |
| 无参数 | `which` | 错误信息 `no command specified`，退出码 1 |
| 多个命令混合 | `which go notexist git` | 输出存在的路径，最后报错，退出码 1 |
| 静默模式存在 | `which -s go` | 无输出，退出码 0 |
| 静默模式不存在 | `which -s notexist` | 无输出，退出码 1 |
| 绝对路径 | `which C:\Windows\System32\cmd.exe` | 直接检查该路径 |

## 十一、实现文件结构

### 11.1 internal/commands/which/cmd_which.go

```go
package which

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

// WhichConfig 配置结构体
type WhichConfig struct {
    Commands []string
    All      bool
    Silent   bool
}

// WhichCmdMain 主函数
func WhichCmdMain(config WhichConfig) error {
    // 实现查找逻辑
}

// findCommand 查找单个命令
func findCommand(name string, all bool) ([]string, error) {
    // 实现查找算法
}

// isExecutable 检查文件是否可执行
func isExecutable(path string) bool {
    // Windows/Unix 分别处理
}
```

### 11.2 internal/cli/which.go

```go
package cli

import (
    "fmt"
    "gitee.com/MM-Q/fck/internal/commands/which"
    "gitee.com/MM-Q/qflag"
)

var WhichCmd *qflag.Cmd

var (
    whichAll    *qflag.BoolFlag
    whichSilent *qflag.BoolFlag
)

func init() {
    // 命令定义和注册
}

func runWhich(cmd qflag.Command) error {
    // 运行函数
}
```

### 11.3 internal/cli/root.go 修改

在 `SubCmds` 列表中添加 `WhichCmd`。
