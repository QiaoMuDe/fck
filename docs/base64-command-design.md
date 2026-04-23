# base64 命令设计方案

## 功能概述

实现 base64 编码和解码功能，支持从文件、标准输入或直接输入字符串进行编解码操作。

## 参考实现

参考 Linux `base64` 命令的行为和选项设计。

## 命令结构

```
fck base64 [options] [string...]     # 默认：编码位置参数（字符串）
fck base64 -d [options] [string...]  # 解码位置参数（字符串）
fck base64 -f file.txt               # 编码文件内容
fck base64 -d -f encoded.txt         # 解码文件内容
```

## 配置结构体

```go
// Base64Config 配置结构体
type Base64Config struct {
	Decode    bool     // 解码模式
	Strings   []string // 位置参数字符串（默认编码/解码这些字符串）
	FilePath  string   // 输入文件路径（通过 -f 指定）
	Output    string   // 输出文件路径（默认标准输出）
	Wrap      int      // 每行字符数（0表示不换行）
	URLSafe   bool     // 使用 URL 安全的 base64 变体
	NoPadding bool     // 禁用填充
}
```

## 命令选项

| 长选项 | 短选项 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| decode | d | bool | false | 解码模式 |
| file | f | string | "" | 从文件读取输入 |
| output | o | string | "" | 输出文件路径 |
| wrap | w | int | 76 | 每行最大字符数（0=不换行） |
| url-safe | u | bool | false | 使用 URL 安全变体（- 替换 +，_ 替换 /） |
| no-padding | p | bool | false | 禁用填充（=） |

## 功能设计

### 1. 编码模式（默认）

```bash
# 编码位置参数（字符串）
fck base64 "hello world"

# 编码多个字符串（空格连接）
fck base64 hello world

# 从文件编码
fck base64 -f file.txt

# 从标准输入编码
echo "hello" | fck base64
fck base64 < file.txt

# 输出到文件
fck base64 -o output.txt "hello"

# 不换行输出
fck base64 -w 0 "hello"

# URL 安全编码
fck base64 -u "hello"
```

### 2. 解码模式

```bash
# 解码位置参数（字符串）
fck base64 -d "aGVsbG8="

# 从文件解码
fck base64 -d -f encoded.txt

# 从标准输入解码
echo "aGVsbG8=" | fck base64 -d
```

### 3. 输入优先级

1. 标准输入（管道/重定向）- 使用 `utils.IsStdinPipe()` 检测
2. 文件输入（`-f`）
3. 位置参数字符串

## 错误处理

| 错误场景 | 错误信息 |
|----------|----------|
| 解码非法字符 | "invalid base64 data: illegal character at position X" |
| 文件不存在 | "file not found: [path]" |
| 输入为空 | "no input provided" |
| 解码填充错误 | "invalid base64 padding" |
| 同时指定文件和字符串 | "cannot use both file and string arguments" |

## 核心实现逻辑

### 读取输入

```go
func readInput(config Base64Config) ([]byte, error) {
	// 1. 优先检测管道/重定向输入
	if utils.IsStdinPipe() {
		return io.ReadAll(os.Stdin)
	}
	
	// 2. 从文件读取
	if config.FilePath != "" {
		return os.ReadFile(config.FilePath)
	}
	
	// 3. 从位置参数读取（空格连接）
	if len(config.Strings) > 0 {
		return []byte(strings.Join(config.Strings, " ")), nil
	}
	
	return nil, fmt.Errorf("no input provided")
}
```

### 编码流程

```go
func encode(config Base64Config) error {
	// 1. 读取输入
	data, err := readInput(config)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no input provided")
	}
	
	// 2. 创建 base64 编码器
	enc := base64.StdEncoding
	if config.URLSafe {
		enc = base64.URLEncoding
	}
	if config.NoPadding {
		enc = enc.WithPadding(base64.NoPadding)
	}
	
	// 3. 编码
	encoded := enc.EncodeToString(data)
	
	// 4. 换行处理
	if config.Wrap > 0 {
		encoded = wrapLines(encoded, config.Wrap)
	}
	
	// 5. 输出
	return writeOutput(encoded+"\n", config.Output)
}
```

### 解码流程

```go
func decode(config Base64Config) error {
	// 1. 读取输入
	data, err := readInput(config)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no input provided")
	}
	
	// 2. 清理换行符（解码时忽略）
	data = bytes.ReplaceAll(data, []byte("\n"), nil)
	data = bytes.ReplaceAll(data, []byte("\r"), nil)
	
	// 3. 创建解码器
	enc := base64.StdEncoding
	if config.URLSafe {
		enc = base64.URLEncoding
	}
	
	// 4. 解码
	decoded, err := enc.DecodeString(string(data))
	if err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}
	
	// 5. 输出（二进制数据不加换行）
	return writeOutputBytes(decoded, config.Output)
}
```

## 文件结构

```
internal/
├── commands/
│   └── base64/
│       └── cmd_base64.go    # 业务逻辑
└── cli/
    └── base64.go            # CLI 定义
```

## 使用示例

```go
// 示例配置
cmdOpts := &qflag.CmdOpts{
    Desc: "Base64 编解码工具",
    Notes: []string{
        "优先从管道/重定向读取输入",
        "其次使用 -f 选项从文件读取",
        "最后编码位置参数字符串",
        "使用 -d 选项切换到解码模式",
        "URL 安全变体将 + 替换为 -，/ 替换为 _",
    },
    Examples: map[string]string{
        "编码字符串":     "fck base64 'hello world'",
        "编码多个词":     "fck base64 hello world",
        "编码文件":       "fck base64 -f file.txt",
        "解码字符串":     "fck base64 -d 'aGVsbG8='",
        "解码文件":       "fck base64 -d -f encoded.txt",
        "不换行输出":     "fck base64 -w 0 'hello'",
        "URL 安全编码":   "fck base64 -u 'hello'",
        "输出到文件":     "fck base64 -o out.txt 'hello'",
        "管道使用":       "echo 'hello' | fck base64",
    },
}
```

## 边缘情况处理

1. **空输入**：返回错误 "no input provided"
2. **二进制数据**：正确处理二进制数据的编解码
3. **大文件**：使用流式处理避免内存溢出
4. **解码错误**：精确定位非法字符位置
5. **换行处理**：编码时按 wrap 换行，解码时自动忽略换行符
6. **解码输出**：二进制数据直接输出，不添加换行符

## 测试用例建议

1. 基本字符串编码/解码
2. 多字符串参数（空格连接）
3. 文件输入编码/解码
4. 标准输入编码/解码
5. URL 安全变体
6. 无填充模式
7. 换行处理（编码 wrap，解码忽略）
8. 二进制文件处理
9. 错误输入处理（非法字符、错误填充）
