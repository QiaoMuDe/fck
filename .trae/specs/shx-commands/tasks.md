# Tasks

所有任务无依赖关系，可并行实现。

- [x] Task 1: 创建 `internal/commands/shfmt/` — shfmt 业务逻辑实现
  - [x] 1.1 创建 `cmd_shfmt.go`，实现 `ShfmtConfig` 结构体和 `ShfmtCmdMain` 入口函数
  - [x] 1.2 使用 `fs.ExpandFiles()` 展开位置参数中的通配符，获取实际文件列表
  - [x] 1.3 支持从文件读取输入：读取展开后的文件内容
  - [x] 1.4 支持从管道读取输入：使用 `term.IsStdinPipe()` 检测并读取
  - [x] 1.5 调用 `shx.Format()` 进行格式化处理
  - [x] 1.6 支持原地写入（-w）和备份（-b），备份使用 `fs.CopyEx()` 实现

- [x] Task 2: 创建 `internal/commands/shck/` — shck 业务逻辑实现
  - [x] 2.1 创建 `cmd_shck.go`，实现 `ShckConfig` 结构体和 `ShckCmdMain` 入口函数
  - [x] 2.2 使用 `fs.ExpandFiles()` 展开位置参数中的通配符
  - [x] 2.3 支持从文件读取输入
  - [x] 2.4 支持从管道读取输入
  - [x] 2.5 调用 `shx.CheckSyntax()` 进行语法检查，根据结果输出成功或错误信息
  - [x] 2.6 语法错误时从 `*shx.SyntaxError` 提取行号、列号和错误描述

- [x] Task 3: 创建 `internal/commands/shx/` — shx 业务逻辑实现
  - [x] 3.1 创建 `cmd_shx.go`，实现 `ShxConfig` 结构体和 `ShxCmdMain` 入口函数
  - [x] 3.2 检测参数是文件（`.sh` 后缀且存在）还是命令字符串
  - [x] 3.3 支持超时控制（-t）：使用 `WithTimeout()` 或 `RunWith()` / `OutWith()`
  - [x] 3.4 支持工作目录（-d）：使用 `WithDir()`
  - [x] 3.5 支持环境变量（-e）：使用 `WithEnv()`，可多次指定
  - [x] 3.6 支持管道 stdin 传递到命令的 stdin

- [x] Task 4: 创建 CLI 接口文件
  - [x] 4.1 创建 `internal/cli/shfmt.go`：qflag 命令定义、标志注册、`runShfmt` 函数
  - [x] 4.2 创建 `internal/cli/shck.go`：qflag 命令定义、标志注册、`runShck` 函数
  - [x] 4.3 创建 `internal/cli/shx.go`：qflag 命令定义、标志注册、`runShx` 函数

- [x] Task 5: 注册子命令到根命令
  - [x] 5.1 在 `internal/cli/root.go` 的 `SubCmds` 中添加 `ShfmtCmd`、`ShckCmd`、`ShxCmd`
  - [x] 5.2 编译验证所有新增命令可正常注册

# Task Dependencies

所有 Task 相互独立，无依赖关系，可并行实现。
