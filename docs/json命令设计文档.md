# json 命令设计文档

## 1. 功能概述

`json` 命令是一个 JSON 数据处理工具，提供格式化、验证、查询、转换等功能，类似于 `jq` 的简化版，但更加易用。

## 2. 功能特性

| 功能 | 描述 | 标志 |
|------|------|------|
| 格式化 | 美化/压缩 JSON 输出 | `-p/--pretty`, `-c/--compact` |
| 验证 | 验证 JSON 语法有效性 | `-v/--validate` |
| 查询 | 使用路径表达式提取数据 | `-q/--query` |
| 颜色 | 语法高亮显示 | `-H/--highlight` |

## 3. 命令行接口设计

### 3.1 命令定义

```go
var JsonCmd *qflag.Cmd

var (
    jsonPretty         *qflag.BoolFlag   // -p, --pretty      美化输出
    jsonCompact        *qflag.BoolFlag   // -c, --compact     压缩输出
    jsonValidate       *qflag.BoolFlag   // -v, --validate    验证模式
    jsonQuery          *qflag.StringFlag // -q, --query       查询路径
    jsonHighlight      *qflag.BoolFlag   // -H, --highlight   语法高亮
    jsonRaw            *qflag.BoolFlag   // -r, --raw         原始字符串输出
)
```

### 3.2 命令配置

```go
cmdOpts := &qflag.CmdOpts{
    Desc:        "JSON 数据处理工具 - 格式化、验证、查询",
    UseChinese:  true,
    UsageSyntax: fmt.Sprintf("%s json [options] [json-string]", qflag.Root.Name()),
    Notes: []string{
        "输入方式: 管道传递JSON字符串 或 位置参数指定文件路径",
        "查询路径使用点号分隔，如: users.0.name",
        "数组索引支持负数，-1 表示最后一个元素",
        "支持通配符 * 匹配数组所有元素，如: users.*.name",
    },
    Examples: map[string]string{
        "格式化 JSON":              fmt.Sprintf("echo '{\"a\":1}' | %s json -p", qflag.Root.Name()),
        "压缩 JSON":               fmt.Sprintf("%s json -c file.json", qflag.Root.Name()),
        "验证 JSON":               fmt.Sprintf("%s json -v invalid.json", qflag.Root.Name()),
        "查询数据":                fmt.Sprintf("echo '{\"users\":[{\"name\":\"Tom\"}]}' | %s json -q users.0.name", qflag.Root.Name()),
        "高亮显示":                fmt.Sprintf("%s json -pH large.json", qflag.Root.Name()),
        "提取数组所有元素":         fmt.Sprintf("echo '{\"items\":[1,2,3]}' | %s json -q items.*", qflag.Root.Name()),
    },
    MutexGroups: []qflag.MutexGroup{
        {
            Name:      "output-format",
            Flags:     []string{"pretty", "compact"},
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
    Pretty    bool     // 美化输出
    Compact   bool     // 压缩输出
    Validate  bool     // 验证模式
    Query     string   // 查询路径
    Highlight bool     // 语法高亮
    Raw       bool     // 原始字符串输出
    Files     []string // 位置参数（文件路径）
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

    // 3. 查询处理 (优先于解析，直接操作原始数据)
    if config.Query != "" {
        result, err := queryJSON(data, config.Query)
        if err != nil {
            return err
        }
        // 查询结果直接输出
        return writeOutput([]byte(result.Raw), config)
    }

    // 4. 解析 JSON
    var jsonData interface{}
    if err := json.Unmarshal(data, &jsonData); err != nil {
        return fmt.Errorf("JSON 解析失败: %w", err)
    }

    // 5. 格式化输出
    output, err := formatOutput(jsonData, config)
    if err != nil {
        return err
    }

    // 6. 输出结果
    return writeOutput(output, config)
}
```

### 4.3 核心功能函数

```go
// readInput 读取输入数据
//
// 输入方式:
//   1. 管道/重定向输入 (使用 utils.IsStdinPipe() 检测) - 传递JSON字符串
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
| github.com/alecthomas/chroma | 语法高亮 | 代码高亮显示 |

### 8.1 gjson 路径语法对照

| 设计语法 | gjson 语法 | 说明 |
|----------|-----------|------|
| `users.0.name` | `users.0.name` | 对象属性 + 数组索引 |
| `users.-1.name` | `users.-1.name` | 负数索引（最后一个） |
| `users.*.name` | `users.#.name` | 通配符匹配所有元素 |

> **注**: gjson 使用 `#` 表示数组通配符，设计文档中保持 `*` 语法，在代码中进行转换

## 9. 错误处理

| 错误场景 | 错误信息 |
|----------|----------|
| 无输入 | "未提供输入数据，请通过管道、文件或参数指定" |
| JSON 解析失败 | "JSON 解析失败: <详细错误>" |
| 查询路径无效 | "查询路径无效: <路径>" |
| 数组索引越界 | "数组索引越界: <索引>" |

| 文件读取失败 | "读取文件失败: <路径>" |

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
**版本**: v1.0  
**状态**: 待评审
