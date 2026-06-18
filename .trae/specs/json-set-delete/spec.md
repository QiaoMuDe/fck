# json-set-delete Spec

## Why
当前 json 命令只支持只读查询和格式化输出，无法直接修改 JSON 内部字段值。用户需要在不手动编辑文件的情况下，通过命令行直接设置或删除 JSON 字段值，配合 `-w` 原地写入实现高效编辑。

## What Changes
- 新增 `--set / -s` 标志，支持通过路径表达式设置 JSON 字段值
- 新增 `--delete / -D` 标志，支持通过路径表达式删除 JSON 字段
- 新增 `--type / -t` 标志，控制设置值的类型推断行为
- 新增依赖 `github.com/tidwall/sjson`（gjson 的配套库）
- 更新 CLI 互斥组和依赖规则
- 更新设计文档

## Impact
- Affected specs: json 命令的格式化/查询/写入流程
- Affected code:
  - `internal/cli/json.go` — 新增 flag 定义、参数解析、互斥组
  - `internal/commands/json/cmd_json.go` — 新增 set/delete 业务逻辑
  - `docs/json命令设计文档.md` — 更新 API 说明和示例
- New dependency: `github.com/tidwall/sjson`

## ADDED Requirements

### Requirement: JSON Set 功能
系统 SHALL 支持通过路径表达式设置 JSON 中的字段值。

#### Scenario: 设置字符串值
- **WHEN** 用户执行 `fck json -s "name=Tom" data.json`
- **THEN** 输出修改后的 JSON，`name` 字段值为 `"Tom"`（字符串）

#### Scenario: 设置数字值（显式类型）
- **WHEN** 用户执行 `fck json -s "age=30" -t number data.json`
- **THEN** 输出修改后的 JSON，`age` 字段值为 `30`（数字，不带引号）

#### Scenario: 设置布尔值
- **WHEN** 用户执行 `fck json -s "active=true" -t bool data.json`
- **THEN** 输出修改后的 JSON，`active` 字段值为 `true`

#### Scenario: 设置 null 值
- **WHEN** 用户执行 `fck json -s "extra=null" data.json`
- **THEN** 输入 `null` 且类型为 auto 时，自动推断为 JSON null

#### Scenario: 多值设置
- **WHEN** 用户执行 `fck json -s "name=Tom,age=30" -t number data.json`
- **THEN** 依次设置 `name` 和 `age` 字段（`-t` 对所有值生效）

#### Scenario: 设置嵌套字段
- **WHEN** 用户执行 `fck json -s "address.city=Beijing" data.json`
- **THEN** 设置嵌套对象中的字段值

#### Scenario: 设置数组元素
- **WHEN** 用户执行 `fck json -s "tags.0=golang" data.json`
- **THEN** 修改数组指定索引的值

#### Scenario: 数组追加元素
- **WHEN** 用户执行 `fck json -s "tags.-1=newtag" data.json`
- **THEN** 在数组末尾追加新值

#### Scenario: 设置 + 原地写入
- **WHEN** 用户执行 `fck json -s "name=Tom" -w data.json`
- **THEN** 修改后的 JSON 写回原文件

#### Scenario: 设置 + 美化输出
- **WHEN** 用户执行 `fck json -s "name=Tom" -p data.json`
- **THEN** 修改后的 JSON 以美化格式输出

#### Scenario: 设置时同时备份
- **WHEN** 用户执行 `fck json -s "name=Tom" -w -b data.json`
- **THEN** 先创建 `.bak` 备份，再写入修改后的 JSON

### Requirement: JSON Delete 功能
系统 SHALL 支持通过路径表达式删除 JSON 中的字段或数组元素。

#### Scenario: 删除字段
- **WHEN** 用户执行 `fck json -D "address" data.json`
- **THEN** 删除 `address` 字段，输出修改后的 JSON

#### Scenario: 删除数组元素
- **WHEN** 用户执行 `fck json -D "tags.1" data.json`
- **THEN** 删除数组索引 1 的元素

#### Scenario: 删除嵌套字段
- **WHEN** 用户执行 `fck json -D "user.phone" data.json`
- **THEN** 删除嵌套对象中的字段

#### Scenario: 删除 + 原地写入
- **WHEN** 用户执行 `fck json -D "temp" -w data.json`
- **THEN** 删除字段后写回原文件

### Requirement: 值类型推断
系统 SHALL 支持通过 `-t` 标志控制值的类型推断。

#### Scenario: auto 模式（默认）
- **WHEN** 用户使用 `-s` 未指定 `-t`（默认为 auto）
- **THEN** sjson 根据 Go 字面量自动推断类型

#### Scenario: string 模式
- **WHEN** 用户执行 `-t string`
- **THEN** 所有值强制作为 JSON 字符串处理（加引号）

#### Scenario: number 模式
- **WHEN** 用户执行 `-t number`
- **THEN** 所有值强制作为 JSON 数字处理

#### Scenario: bool 模式
- **WHEN** 用户执行 `-t bool`
- **THEN** 所有值强制作为 JSON 布尔值处理

### Requirement: 互斥规则
系统 SHALL 保证 set、delete、query、validate 之间互斥。

#### Scenario: set 与 query 互斥
- **WHEN** 用户同时使用 `-s` 和 `-q`
- **THEN** 报错提示互斥

#### Scenario: delete 与 query 互斥
- **WHEN** 用户同时使用 `-D` 和 `-q`
- **THEN** 报错提示互斥

#### Scenario: set 与 delete 互斥
- **WHEN** 用户同时使用 `-s` 和 `-D`
- **THEN** 报错提示互斥

#### Scenario: set/delete 与 validate 互斥
- **WHEN** 用户同时使用 `-s`/`-D` 和 `-v`
- **THEN** 报错提示互斥

### Requirement: 管道输入兼容
系统 SHALL 支持从管道输入 JSON 数据并进行 set/delete 操作。

#### Scenario: 管道输入 + set
- **WHEN** 用户执行 `echo '{"name":"old"}' | fck json -s "name=new"`
- **THEN** 从 stdin 读取 JSON，设置后输出到 stdout

#### Scenario: 管道输入 + delete
- **WHEN** 用户执行 `echo '{"a":1,"b":2}' | fck json -D "a"`
- **THEN** 从 stdin 读取 JSON，删除后输出到 stdout

## ADDED CLI Configuration

### 新增 Flags
| 短标志 | 长标志 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| `-s` | `--set` | string | `""` | 设置字段值，格式: `path1=val1,path2=val2` |
| `-D` | `--delete` | string | `""` | 删除字段，格式: `path1,path2` |
| `-t` | `--type` | enum | `"auto"` | 值类型: auto/string/number/bool |

### 互斥组调整
```
MutexGroups:
  - Name: "operation-mode"
    Flags: ["pretty", "compact", "validate", "query", "set", "delete"]
    AllowNone: true

  - Name: "type-set"
    Flags: ["type", "set"]
    AllowNone: true
    # type 标志必须与 set 搭配使用
```
