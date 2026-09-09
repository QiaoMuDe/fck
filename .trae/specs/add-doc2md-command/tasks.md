# Tasks

- [x] Task 1: 添加 doc2md 依赖并构建验证
  - 在 `go.mod` 添加 `gitee.com/MM-Q/doc2md` 依赖并执行 `go get`，确认 `go build ./...` 通过
- [x] Task 2: 创建业务实现层 `internal/commands/doc2md/`
  - [x] SubTask 2.1: 创建 `types.go`，定义 `Doc2mdConfig`（Files / Extension / MIMEType / Charset / Output / KeepDataURIs）
  - [x] SubTask 2.2: 创建 `cmd_doc2md.go`，实现 `Doc2mdCmdMain(config)`：
    - 使用 `gitee.com/MM-Q/doc2md` 创建引擎（`markitdown.New()`）
    - 若 `KeepDataURIs` 为 true，应用 `WithKeepDataURIs(true)` 选项
    - 优先判断 stdin 管道：`term.IsStdinPipe()` 且 `-x` 已提供 → 读取 stdin 到内存，构造 `markitdown.StreamInfo{Extension: ...}` 并 `ConvertReader`
    - 否则展开通配符（`fs.Expand`）逐个 `ConvertFile`
    - 输出：指定 `-o` 写入文件；否则打印到 stdout
    - 使用 `markitdown.IsUnsupportedFormat` 区分错误并给出中文可读提示
- [x] Task 3: 创建命令定义层 `internal/cli/doc2md.go`
  - 定义 `Doc2mdCmd`，声明 `-o/-output`、`-x/-extension`、`-m/-mime-type`、`-c/-charset`、`--keep-data-uris` 标志
  - 配置 `CmdOpts`（Desc/UseChinese/Examples/Notes），`SetRun(runDoc2md)` 将参数映射为 `Doc2mdConfig` 后调用 `Doc2mdCmdMain`
- [x] Task 4: 在 `internal/cli/root.go` 注册 `Doc2mdCmd` 到 `SubCmds`
- [x] Task 5: 验证与文档同步
  - [x] SubTask 5.1: 运行 `go build ./...`、`go vet ./...`、`go test ./...` 全绿
  - [x] SubTask 5.2: 手动验证 `fck doc2md input.docx` 与管道输入两种路径
  - [x] SubTask 5.3: 在 fck-skill 以及对应命令帮助中补充 doc2md 命令用法（若该技能列表涉及）

# Task Dependencies

- [Task 2] depends on [Task 1]
- [Task 3] depends on [Task 2]
- [Task 4] depends on [Task 3]
- [Task 5] depends on [Task 1]、[Task 2]、[Task 3]、[Task 4]