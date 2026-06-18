# Checklist

- [x] `sjson` 依赖已添加，`go.mod` 和 `vendor` 已更新
- [x] `docs/json命令设计文档.md` 已更新，包含 set/delete 功能说明和示例
- [x] `internal/cli/json.go` 中新增了 `jsonSet`、`jsonDelete`、`jsonType` 三个 flag
- [x] `internal/cli/json.go` 的 `MutexGroups` 已更新，set/delete/query/validate 互斥
- [x] `internal/cli/json.go` 的 `Notes` 和 `Examples` 已更新
- [x] `internal/commands/json/cmd_json.go` 的 `JsonConfig` 新增了 `SetValue`、`DeletePath`、`ValueType` 字段
- [x] `setJSON()` 函数已实现：解析逗号分隔的 `path=value` 对，按类型调用 sjson
- [x] `deleteJSON()` 函数已实现：解析逗号分隔的 paths，调用 sjson.Delete
- [x] `JsonCmdMain()` 中已增加 set/delete 处理分支
- [x] set/delete 模式下支持管道输入
- [x] set/delete 模式下支持 `-w` 原地写入和 `-b` 备份
- [x] set/delete 模式下支持 `-p` 美化输出
- [x] 编译通过 (`go build ./...`)
- [x] 手动测试通过：set 字符串/数字/布尔/多值/嵌套/数组/追加
- [x] 手动测试通过：delete 字段/数组元素/嵌套字段
- [x] 互斥规则验证：`-s`+`-q`、`-D`+`-q`、`-s`+`-D`、`-s`+`-v` 均报错
