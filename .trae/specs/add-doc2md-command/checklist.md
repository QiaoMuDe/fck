# 检查清单

- [x] `internal/commands/doc2md/types.go` 定义 `Doc2mdConfig` 结构与 spec 一致
- [x] `internal/commands/doc2md/cmd_doc2md.go` 实现文件路径转换（含通配符展开）
- [x] `internal/commands/doc2md/cmd_doc2md.go` 实现 stdin 管道输入（`-x` 提示扩展名，缺失时报错）
- [x] 输出行为符合 spec：默认 stdout、`-o` 写文件
- [x] `--keep-data-uris`、`-m/-c` 等选项正确传递到 doc2md 引擎
- [x] 不支持格式与转换失败返回非零退出码并给出可读错误
- [x] `internal/cli/doc2md.go` 声明全部标志并调用 `Doc2mdCmdMain`
- [x] `internal/cli/root.go` 的 `SubCmds` 已注册 `Doc2mdCmd`
- [x] `go.mod` 已添加 `gitee.com/MM-Q/doc2md` 依赖
- [x] `go build ./...`、`go vet ./...`、`go test ./...` 全绿
- [x] 新增代码注释为简体中文，遵循项目规范
- [x] `fck --help` 与 `fck doc2md --help` 能正常显示命令信息