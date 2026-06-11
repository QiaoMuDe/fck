# json 命令设计文档

## 1. 功能概述

`json` 命令是一个 JSON 数据处理工具，提供格式化、验证、查询、设置字段值、删除字段等功能，类似于 `jq` 的简化版，但更加易用。

## 2. 功能特性

| 功能 | 描述 | 标志 |
|------|------|------|
| 格式化 | 美化/压缩 JSON 输出 | `-p/--pretty`, `-c/--compact` |
| 验证 | 验证 JSON 语法有效性 | `-v/--validate` |
| 查询 | 使用路径表达式提取数据 | `-q/--query` |
| 设置 | 使用路径表达式设置字段值 | `-s/--set` |
| 删除 | 使用路径表达式删除字段 | `-D/--delete` |
| 类型 | 设置值类型（auto/string/number/bool） | `-t/--type` |
| 颜色 | 语法高亮显示 | `-H/--highlight` |

## 3. 命令行接口设计

### 3.1 命令定义

```go
var JsonCmd *qflag.Cmd

var (
    jsonPretty         *qflag.BoolFlag   // -p, --pretty      美化输出
    jsonCompact        *qflag.BoolFlag   // -c, --compact     压缩输出
    jsonValidate       *qflag.BoolFlag   // -v, --validate    验证模式
    jsonSet            *qflag.StringFlag // -s, --set         设置字段值
    jsonDelete         *qflag.StringFlag // -D, --delete      删除字段
    jsonType           *qflag.EnumFlag   // -t, --type        值类型
    jsonQuery          *qflag.StringFlag // -q, --query       查询路径
    jsonHighlight      *qflag.BoolFlag   // -H, --highlight   语法高亮
    jsonRaw            *qflag.BoolFlag   // -r, --raw         原始字符串输出
    jsonWrite          *qflag.BoolFlag   // -w, --write       原地写入
    jsonBackup         *qflag.BoolFlag   // -b, --backup      写入前备份
)
```

### 3.2 命令配置

```go
cmdOpts := &qflag.CmdOpts{
    Desc:        "JSON 数据处理工具 - 格式化、验证、查询、设置、删除",
    UseChinese:  true,
    UsageSyntax: fmt.Sprintf("%s json [options] [file]", qflag.Root.Name()),
    Notes: []string{
        "输入方式: 管道传递JSON字符串 或 位置参数指定文件路径",
        "查询路径使用点号分隔，如: users.0.name",
        "数组索引支持负数，-1 表示最后一个元素",
        "支持通配符 * 匹配数组所有元素，如: users.*.name",
        "使用 -s 设置值，格式: path1=val1,path2=val2 (逗号分隔多个操作)",
        "使用 -D 删除字段，格式: path1,path2 (逗号分隔多个路径)",
        "使用 -t 指定值类型: auto(默认), string, number, bool",
        "使用 -w 原地写入时，必须指定文件路径参数（不支持管道）",
        "使用 -b 备份时，会创建 .bak 后缀的备份文件",
    },
    Examples: map[string]string{
        "格式化 JSON (管道)":         fmt.Sprintf("echo '{\"a\":1}' | %s json -p", qflag.Root.Name()),
        "格式化 JSON (文件)":         fmt.Sprintf("%s json -p data.json", qflag.Root.Name()),
        "压缩 JSON":                fmt.Sprintf("%s json -c file.json", qflag.Root.Name()),
        "验证 JSON":                fmt.Sprintf("%s json -v invalid.json", qflag.Root.Name()),
        "查询数据 (管道)":            fmt.Sprintf("echo '{\"users\":[{\"name\":\"Tom\"}]}' | %s json -q users.0.name", qflag.Root.Name()),
        "查询数据 (文件)":            fmt.Sprintf("%s json -q users.0.name data.json", qflag.Root.Name()),
        "设置字段值 (字符串)":         fmt.Sprintf("%s json -s \"name=Tom\" data.json", qflag.Root.Name()),
        "设置字段值 (数字)":          fmt.Sprintf("%s json -s \"age=30\" -t number data.json", qflag.Root.Name()),
        "多值设置":                  fmt.Sprintf("%s json -s \"name=Tom,age=30\" -t number data.json", qflag.Root.Name()),
        "设置嵌套字段":               fmt.Sprintf("%s json -s \"address.city=Beijing\" data.json", qflag.Root.Name()),
        "数组追加元素":               fmt.Sprintf("%s json -s \"tags.-1=newtag\" data.json", qflag.Root.Name()),
        "删除字段":                  fmt.Sprintf("%s json -D \"address\" data.json", qflag.Root.Name()),
        "删除数组元素":               fmt.Sprintf("%s json -D \"tags.1\" data.json", qflag.Root.Name()),
        "高亮显示":                  fmt.Sprintf("%s json -pH large.json", qflag.Root.Name()),
        "提取数组所有元素":            fmt.Sprintf("echo '{\"items\":[1,2,3]}' | %s json -q items.*", qflag.Root.Name()),
        "原地格式化文件":              fmt.Sprintf("%s json -p -w data.json", qflag.Root.Name()),
        "原地压缩并备份":              fmt.Sprintf("%s json -c -w -b config.json", qflag.Root.Name()),
        "原地设置并备份":              fmt.Sprintf("%s json -s \"name=Tom\" -w -b data.json", qflag.Root.Name()),
    },
    MutexGroups: []qflag.MutexGroup{
        {
            // 操作模式互斥: set/delete/query/validate/pretty/compact 之间互斥
            Name:      "operation-mode",
            Flags:     []string{"pretty", "compact", "validate", "query", "set", "delete"},
            AllowNone: true,
        },
        {
            // type 只能与 set 搭配使用（依赖于 set）
            Name:      "type-set",
            Flags:     []string{"type", "set"},
            AllowNone: true,
        },
    },
}
```

## 4. 业务逻辑设计

### 4.1 配置结构体

```go
package json

// JsonConfig 配置结构体
type JsonConfig struct {
    Pretty     bool     // 美化输出
    Compact    bool     // 压缩输出
    Validate   bool     // 验证模式
    SetValue   string   // 设置值表达式: path1=val1,path2=val2
    DeletePath string   // 删除路径: path1,path2
    ValueType  string   // 值类型: auto/string/number/bool
    Query      string   // 查询路径
    Highlight  bool     // 语法高亮
    Raw        bool     // 原始字符串输出
    Write      bool     // 原地写入
    Backup     bool     // 写入前备份
    Files      []string // 位置参数（文件路径）
}

// JsonStats 操作统计
type JsonStats struct {
    IsValid    bool   // 是否有效
    ParseTime  string // 解析耗时
    OutputSize int    // 输出大小
}
```

### 4.2 主函数

```go
// JsonCmdMain 执行 json 命令
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - error: 执行错误
func JsonCmdMain(config JsonConfig) error {
    // 1. 读取输入
    data, err := readInput(config)
    if err != nil {
        return err
    }

    // 2. 验证模式
    if config.Validate {
        return validateJSON(data)
    }

    // 3. 设置模式 — 使用 sjson 设置字段值
    if config.SetValue != "" {
        result, err := setJSON(data, config.SetValue, config.ValueType)
        if err != nil {
            return err
        }
        return writeOutput(result, config)
    }

    // 4. 删除模式 — 使用 sjson 删除字段
    if config.DeletePath != "" {
        result, err := deleteJSON(data, config.DeletePath)
        if err != nil {
            return err
        }
        return writeOutput(result, config)
    }

    // 5. 查询处理 (优先于解析，直接操作原始数据)
    if config.Query != "" {
        result, err := queryJSON(data, config.Query)
        if err != nil {
            return err
        }
        // 查询结果直接输出
        return writeOutput([]byte(result.Raw), config)
    }

    // 6. 解析 JSON
    var jsonData interface{}
    if err := json.Unmarshal(data, &jsonData); err != nil {
        return fmt.Errorf("JSON 解析失败: %w", err)
    }

    // 7. 格式化输出
    output, err := formatOutput(jsonData, config)
    if err != nil {
        return err
    }

    // 8. 输出结果
    return writeOutput(output, config)
}
```

### 4.3 核心功能函数

```go
// readInput 读取输入数据
//
// 输入方式:
//   1. 管道/重定向输入 (使用 term.IsStdinPipe() 检测) - 传递JSON字符串
//   2. 位置参数 - 指定文件路径，读取文件内容
//
// 参数:
//   - config: 命令配置
//
// 返回值:
//   - []byte: 输入数据
//   - error: 读取错误
func readInput(config JsonConfig) ([]byte, error)

// validateJSON 验证 JSON 有效性
//
// 参数:
//   - data: JSON 数据
//
// 返回值:
//   - error: 验证错误 (nil 表示有效)
func validateJSON(data []byte) error

// setJSON 使用路径表达式设置 JSON 字段值
//
// 使用 tidwall/sjson 库实现高性能设置
// 支持多个操作逗号分隔: path1=val1,path2=val2
// 路径语法与 gjson 兼容，支持:
//   - 对象属性: users.name
//   - 数组索引: users.0
//   - 负数索引: users.-1 (追加)
//   - 嵌套路径: address.city
//
// 参数:
//   - data: JSON 原始数据
//   - setStr: 设置表达式
//   - valueType: 值类型 (auto/string/number/bool)
//
// 返回值:
//   - []byte: 修改后的 JSON
//   - error: 设置错误
func setJSON(data []byte, setStr string, valueType string) ([]byte, error)

// deleteJSON 使用路径表达式删除 JSON 字段
//
// 使用 tidwall/sjson 库实现
// 支持多个路径逗号分隔: path1,path2
//
// 参数:
//   - data: JSON 原始数据
//   - deleteStr: 删除路径
//
// 返回值:
//   - []byte: 修改后的 JSON
//   - error: 删除错误
func deleteJSON(data []byte, deleteStr string) ([]byte, error)

// queryJSON 使用路径表达式查询 JSON
//
// 使用 tidwall/gjson 库实现高性能查询
// 路径语法:
//   - 对象属性: users.name
//   - 数组索引: users.0
//   - 负数索引: users.-1 (最后一个)
//   - 通配符:   users.*.name (所有元素, 内部转换为 users.#.name)
//
// 参数:
//   - data: JSON 原始数据 ([]byte)
//   - path: 查询路径
//
// 返回值:
//   - gjson.Result: 查询结果
//   - error: 查询错误
func queryJSON(data []byte, path string) (gjson.Result, error)

// formatOutput 格式化输出
//
// 参数:
//   - data: JSON 数据
//   - config: 命令配置
//
// 返回值:
//   - []byte: 格式化后的数据
//   - error: 格式化错误
func formatOutput(data interface{}, config JsonConfig) ([]byte, error)

// writeOutput 输出结果
//
// 参数:
//   - data: 输出数据
//   - config: 命令配置
//
// 返回值:
//   - error: 输出错误
func writeOutput(data []byte, config JsonConfig) error
```

## 5. 查询路径语法

### 5.1 基本语法

| 语法 | 说明 | 示例 |
|------|------|------|
| `.key` | 对象属性访问 | `users.name` |
| `[n]` | 数组索引 | `items[0]` |
| `[-n]` | 负数索引 | `items[-1]` |
| `.*` | 通配符 | `users.*.email` |

### 5.2 示例数据

```json
{
  "users": [
    {"name": "Alice", "age": 30, "email": "alice@example.com"},
    {"name": "Bob", "age": 25, "email": "bob@example.com"}
  ],
  "meta": {
    "total": 2,
    "page": 1
  }
}
```

### 5.3 查询示例

| 查询路径 | 结果 |
|----------|------|
| `users` | 整个 users 数组 |
| `users.0` | 第一个用户对象 |
| `users.0.name` | "Alice" |
| `users.-1.name` | "Bob" |
| `users.*.name` | ["Alice", "Bob"] |
| `meta.total` | 2 |

## 6. 输出格式

### 6.1 JSON 美化输出 (-p)

```json
{
  "users": [
    {
      "name": "Alice",
      "age": 30
    }
  ]
}
```

### 6.2 JSON 压缩输出 (-c)

```json
{"users":[{"name":"Alice","age":30}]}
```

## 7. 文件结构

```
internal/
├── cli/
│   └── json.go              # CLI 定义
└── commands/
    └── json/
        └── cmd_json.go      # 业务逻辑
```

## 8. 依赖库

| 库 | 用途 | 说明 |
|----|------|------|
| encoding/json | 标准库 JSON 解析 | Go 标准库 |
| github.com/tidwall/gjson | JSON 路径查询 | 高性能查询库 |
| github.com/tidwall/sjson | JSON 字段设置/删除 | 高性能设置库 (gjson 配套库) |
| github.com/alecthomas/chroma | 语法高亮 | 代码高亮显示 |

### 8.1 gjson 路径语法对照

| 设计语法 | gjson 语法 | 说明 |
|----------|-----------|------|
| `users.0.name` | `users.0.name` | 对象属性 + 数组索引 |
| `users.-1.name` | `users.-1.name` | 负数索引（最后一个） |
| `users.*.name` | `users.#.name` | 通配符匹配所有元素 |

> **注**: gjson 使用 `#` 表示数组通配符，设计文档中保持 `*` 语法，在代码中进行转换

### 8.2 sjson 路径语法

sjson 使用与 gjson 兼容的路径语法，额外支持：

| 语法 | 说明 | 示例 |
|------|------|------|
| `path=value` | 设置字段值 | `name=Tom` |
| `path` (delete) | 删除字段 | `address` |
| `array.-1` | 数组末尾追加 | `tags.-1=newtag` |

## 9. 错误处理

| 错误场景 | 错误信息 |
|----------|----------|
| 无输入 | "未提供输入数据，请通过管道、文件或参数指定" |
| JSON 解析失败 | "JSON 解析失败: <详细错误>" |
| 查询路径无效 | "查询路径无效: <路径>" |
| 数组索引越界 | "数组索引越界: <索引>" |
| 设置表达式格式错误 | "设置表达式格式无效: <表达式>" |
| 删除路径错误 | "删除路径无效: <路径>" |

## 10. 测试用例

### 10.1 基础功能测试

```bash
# 格式化测试
echo '{"a":1,"b":2}' | fck json -p

# 压缩测试
fck json -c pretty.json

# 验证测试
fck json -v invalid.json

# 查询测试
echo '{"users":[{"name":"Tom"}]}' | fck json -q users.0.name

# 设置字段值
echo '{"name":"old"}' | fck json -s "name=new"

# 设置数字值
echo '{"age":10}' | fck json -s "age=30" -t number

# 多值设置
echo '{"a":1,"b":2}' | fck json -s "a=10,b=20" -t number

# 删除字段
echo '{"a":1,"b":2}' | fck json -D "a"

# 原地设置并备份
fck json -s "name=Tom" -w -b data.json
```

### 10.2 边缘案例

```bash
# 空对象
echo '{}' | fck json -p

# 空数组
echo '[]' | fck json -p

# 嵌套深度
echo '{"a":{"b":{"c":{"d":1}}}}' | fck json -p

# 特殊字符
echo '{"text":"hello\nworld\t!"}' | fck json -p

# Unicode
echo '{"name":"你好世界"}' | fck json -p

# 大数字
echo '{"big":9007199254740993}' | fck json -p

# 设置 null 值
echo '{"a":1}' | fck json -s "a=null"

# 数组追加
echo '{"items":[1,2]}' | fck json -s "items.-1=3" -t number
```

## 11. 集成到 root.go

在 `internal/cli/root.go` 的 `SubCmds` 列表中添加：

```go
SubCmds: []qflag.Command{
    // ... 现有命令
    Base64Cmd,
    JsonCmd,    // 新增
    TeeCmd,
},
```

## 12. 实现步骤

1. 创建 `internal/commands/json/cmd_json.go` 业务逻辑文件
2. 创建 `internal/cli/json.go` CLI 定义文件
3. 在 `internal/cli/root.go` 中注册命令
4. 运行 `go mod tidy` 安装依赖
5. 编译测试 `go build ./...`
6. 功能测试验证

---

**设计日期**: 2026-04-24  
**版本**: v1.1  
**状态**: 已实现 set/delete 功能
