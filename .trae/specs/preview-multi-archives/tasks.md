# Tasks

- [x] Task 1: 修改命令层 `PreviewConfig` 和 `PreviewCmdMain` 支持多压缩包
  - 将 `PackPath string` 改为 `PackPaths []string`
  - 新增 `PreviewArchives` 函数循环处理多个压缩包
  - 多个压缩包时在输出间打印 `==> filename <==` 分隔标题
  - 单个压缩包时保持原有行为（不打印标题）

- [x] Task 2: 修改 CLI 层 `preview.go` 支持通配符展开
  - 导入 `gitee.com/MM-Q/go-kit/fs`
  - 调用 `fs.ExpandFiles(args)` 展开通配符
  - 传递展开后的路径列表到命令层

# Task Dependencies

- [Task 2] depends on [Task 1]
