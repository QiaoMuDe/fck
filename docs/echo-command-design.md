# Echo 命令设计方案

## 功能概述
实现一个用于输出内容的控制台命令 `echo`，支持文本输出、格式化、特殊字符处理等功能。

## 核心功能

### 1. 基本文本输出
- 支持输出普通文本到控制台
- 支持多行文本输出
- 支持空字符串输出（仅换行）

### 2. 特殊字符处理
- 支持转义字符：`\n`（换行）、`\t`（制表符）、`\r`（回车）、`\\`（反斜杠）、`\"`（双引号）
- 支持原始输出模式（不解析转义字符）

### 3. 格式化选项
- 支持不换行输出（默认换行）
- 支持去除首尾空白字符
- 支持大写/小写转换
- 支持重复输出（指定重复次数）

### 4. 颜色输出（可选）
- 支持ANSI颜色代码输出
- 支持预设颜色别名（red, green, yellow, blue等）

## 命令选项设计

| 选项 | 短选项 | 说明 | 默认值 |
|------|--------|------|--------|
| `--no-newline` | `-n` | 不在末尾添加换行符 | false |
| `--raw` | `-r` | 原始输出，不解析转义字符 | false |
| `--trim` | `-t` | 去除首尾空白字符 | false |
| `--upper` | `-u` | 转换为大写 | false |
| `--lower` | `-l` | 转换为小写 | false |
| `--repeat` | `-R` | 重复输出次数 | 1 |
| `--color` | `-c` | 颜色输出（red, green, yellow, blue等） | "" |

**注意**：输出内容通过位置参数传递，多个参数会用空格连接，与 Linux echo 命令一致。

## 使用示例

### 基本输出
```bash
fck echo "Hello, World!"
# 输出: Hello, World!
```

### 多个参数输出（自动用空格连接）
```bash
fck echo "Hello," "World!"
# 输出: Hello, World!
```

### 不换行输出
```bash
fck echo -n "Hello, "
fck echo "World!"
# 输出: Hello, World!
```

### 原始输出
```bash
fck echo -r "Line1\nLine2"
# 输出: Line1\nLine2
```

### 解析转义字符
```bash
fck echo "Line1\nLine2"
# 输出:
# Line1
# Line2
```

### 格式化输出
```bash
fck echo -u "hello"
# 输出: HELLO

fck echo -t "  hello  "
# 输出: hello
```

### 重复输出
```bash
fck echo -R 3 "Hello"
# 输出:
# Hello
# Hello
# Hello
```

### 颜色输出
```bash
fck echo -c red "Error message"
# 输出红色文本: Error message
```

## 文件结构

```
internal/
├── commands/
│   └── echo/
│       └── cmd_echo.go          # 业务逻辑实现
└── cli/
    └── echo.go                   # CLI定义和配置
```

## 实现细节

### 1. 配置结构体

```go
type EchoConfig struct {
    Text       string
    NoNewline  bool
    Raw        bool
    Trim       bool
    Upper      bool
    Lower      bool
    Repeat     int
    Color      string
}
```

### 2. 核心处理流程

```
1. 解析命令行参数，构建 EchoConfig
2. 文本预处理：
   - 如果 Raw=false，解析转义字符
   - 如果 Trim=true，去除首尾空白
   - 如果 Upper=true，转换为大写
   - 如果 Lower=true，转换为小写
3. 重复输出：
   - 根据 Repeat 次数循环输出
   - 如果 Color 不为空，应用颜色
   - 如果 NoNewline=false，添加换行符
4. 输出到控制台
```

### 3. 转义字符处理

支持的转义字符：
- `\n` → 换行符
- `\t` → 制表符
- `\r` → 回车符
- `\\` → 反斜杠
- `\"` → 双引号
- `\'` → 单引号

### 4. 颜色支持

支持的ANSI颜色：
- `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- `bright-black`, `bright-red`, `bright-green`, `bright-yellow`, `bright-blue`, `bright-magenta`, `bright-cyan`, `bright-white`

颜色实现：
```go
var colorCodes = map[string]string{
    "red":   "\033[31m",
    "green": "\033[32m",
    // ... 其他颜色
    "reset": "\033[0m",
}
```

## 边界情况处理

1. **空文本**：输出空行或什么都不输出（取决于 NoNewline）
2. **重复次数为0或负数**：不输出任何内容
3. **无效颜色**：忽略颜色选项，正常输出
4. **同时指定 Upper 和 Lower**：优先级冲突，默认 Lower 优先
5. **转义字符解析失败**：保持原样输出

## 测试用例

### 单元测试
- 测试基本输出
- 测试转义字符解析
- 测试格式化选项（trim, upper, lower）
- 测试重复输出
- 测试不换行输出
- 测试原始输出模式

### 集成测试
- 测试命令行参数解析
- 测试多个选项组合使用
- 测试边界情况（空文本、重复次数为0等）

## 注意事项

1. 与系统 `echo` 命令的兼容性：保持基本用法一致
2. 跨平台支持：Windows和Linux/macOS的换行符处理
3. 性能考虑：大量重复输出时的性能优化
4. 错误处理：无效参数的友好提示

## 后续扩展（可选）

1. 支持输出到文件（重定向）
2. 支持变量替换（如 ${ENV_VAR}）
3. 支持进度条输出
4. 支持富文本格式（粗体、下划线等）
