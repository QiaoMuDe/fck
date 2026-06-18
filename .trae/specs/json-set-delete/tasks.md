# Tasks

- [x] Task 1: 添加 `sjson` 依赖并更新设计文档
  - 执行 `go get github.com/tidwall/sjson` 添加依赖
  - 更新 `docs/json命令设计文档.md`，补充 set/delete 功能说明、新标志、使用示例

- [x] Task 2: CLI 层 — 在 `internal/cli/json.go` 中注册新标志和互斥规则
  - 新增 `jsonSet` (`-s/--set`)、`jsonDelete` (`-D/--delete`)、`jsonType` (`-t/--type`) 三个 flag 变量
  - 更新 `MutexGroups`：将 `set`/`delete` 加入 "operation-mode" 互斥组
  - 更新 `cmdOpts.Notes` 和 `Examples` 补充新功能的帮助说明
  - 在 `runJson()` 中将新 flag 值传递到 `JsonConfig`

- [x] Task 3: 业务层 — 在 `internal/commands/json/cmd_json.go` 中实现 set/delete 逻辑
  - 在 `JsonConfig` 结构体中新增 `SetValue`、`DeletePath`、`ValueType` 字段
  - 实现 `setJSON(data []byte, setStr string, valueType string) ([]byte, error)` 函数
    - 解析 `setStr`：按逗号拆分得到多个 `path=value` 对
    - 按 `valueType` 类型调用 `sjson.Set()` 设置值
    - 逐个应用所有修改
  - 实现 `deleteJSON(data []byte, deleteStr string) ([]byte, error)` 函数
    - 解析 `deleteStr`：按逗号拆分得到多个 path
    - 逐个调用 `sjson.Delete()`
  - 在 `JsonCmdMain()` 中增加 set/delete 处理分支（和 query 同级，互斥）

- [x] Task 4: 编译验证和功能测试
  - 运行 `go build ./...` 确保编译通过
  - 手动测试各场景：set/delete/类型推断/管道/原地写入/备份
  - 测试互斥规则是否生效

# 任务依赖
- Task 2 依赖 Task 1
- Task 3 依赖 Task 2
- Task 4 依赖 Task 3
