# Tasks

三个任务无依赖关系，可并行实现。

- [x] Task 1: shck 新增静默模式（-q/--quiet）
  - [x] 1.1 CLI 层：在 `internal/cli/shck.go` 注册 `-q`/`--quiet` bool 标志
  - [x] 1.2 CLI 层：将 `quiet` 字段传入 `ShckConfig`
  - [x] 1.3 业务层：在 `ShckConfig` 结构体添加 `Quiet bool` 字段
  - [x] 1.4 业务层：在 `ShckCmdMain` 中，当 `Quiet=true` 时跳过成功信息的 `Printf`，仅输出错误
  - [x] 1.5 编译验证

- [x] Task 2: shfmt 新增列表模式（-l/--list）
  - [x] 2.1 CLI 层：在 `internal/cli/shfmt.go` 注册 `-l`/`--list` bool 标志
  - [x] 2.2 CLI 层：将 `list` 字段传入 `ShfmtConfig`
  - [x] 2.3 业务层：在 `ShfmtConfig` 结构体添加 `List bool` 字段
  - [x] 2.4 业务层：在 `processFile` 或 `ShfmtCmdMain` 中，当 `List=true` 时：读取→格式化→与原内容比较→有差异则输出文件名，不写入文件
  - [x] 2.5 `-l` 优先级高于 `-w`，设置 `-l` 时忽略 `-w`
  - [x] 2.6 编译验证

- [x] Task 3: shx 退出码透传
  - [x] 3.1 在 `internal/cli/shx.go` 的 `runShx` 函数中，捕获 `ShxCmdMain` 返回的错误
  - [x] 3.2 使用 `shx.IsExitStatus()` 检测是否是退出状态错误
  - [x] 3.3 如果是退出状态错误，提取退出码并调用 `os.Exit(int(code))` 直接退出，不返回错误给 qflag
  - [x] 3.4 其他错误（超时、取消等）正常返回给 qflag 处理
  - [x] 3.5 编译验证

# Task Dependencies

所有 Task 相互独立，无依赖关系，可并行实现。
