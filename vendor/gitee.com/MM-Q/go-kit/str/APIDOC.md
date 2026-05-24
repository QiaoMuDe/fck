# Package str

**Import:** `gitee.com/MM-Q/go-kit/str`

## Functions

### func BuildStr

```go
func BuildStr(fn func(*strings.Builder)) string
```

BuildStr 使用字符串构建器执行函数

**参数:**
- `fn`: 使用字符串构建器的函数

**返回值:**
- `string`: 函数执行后构建的字符串结果

**说明:**
- 创建新的字符串构建器（不使用对象池）
- 执行用户提供的函数
- 返回构建的字符串结果
- 适用于不需要对象池优化的简单字符串构建场景

**使用示例:**

```go
result := str.BuildStr(func(buf *strings.Builder) {
    buf.WriteString("Hello")
    buf.WriteByte(' ')
    buf.WriteString("World")
})
```

### func BuildStrCap

```go
func BuildStrCap(cap int, fn func(*strings.Builder)) string
```

BuildStrCap 使用指定容量的字符串构建器执行函数

**参数:**
- `cap`: 字符串构建器初始容量
- `fn`: 使用字符串构建器的函数

**返回值:**
- `string`: 函数执行后构建的字符串结果

**说明:**
- 创建指定容量的新字符串构建器（不使用对象池）
- 执行用户提供的函数
- 返回构建的字符串结果
- 适用于已知字符串长度且不需要对象池优化的场景

**使用示例:**

```go
result := str.BuildStrCap(64, func(buf *strings.Builder) {
    buf.WriteString("Hello")
    buf.WriteByte(' ')
    buf.WriteString("World")
})
```

### func Ellipsis

```go
func Ellipsis(s string, maxLen int) string
```

Ellipsis 超长字符串显示省略号

**参数:**
- `s`: 输入字符串
- `maxLen`: 最大长度（包含省略号）

**返回:**
- `string`: 处理后的字符串

### func FromBase64

```go
func FromBase64(s string) (string, error)
```

FromBase64 将Base64字符串解码

**参数:**
- `s`: Base64编码的字符串

**返回:**
- `string`: 解码后的字符串
- `error`: 解码错误

### func IfBlank

```go
func IfBlank(s, defaultVal string) string
```

IfBlank 当字符串为空白时返回默认值

**参数:**
- `s`: 待检查的字符串
- `defaultVal`: 默认值

**返回:**
- `string`: 如果s为空白（空字符串或只包含空白字符）则返回defaultVal，否则返回s

### func IfEmpty

```go
func IfEmpty(s, defaultVal string) string
```

IfEmpty 当字符串为空时返回默认值

**参数:**
- `s`: 待检查的字符串
- `defaultVal`: 默认值

**返回:**
- `string`: 如果s为空则返回defaultVal，否则返回s

### func IsEmpty

```go
func IsEmpty(s string) bool
```

IsEmpty 检查字符串是否为空

**参数:**
- `s`: 待检查的字符串

**返回:**
- `bool`: 字符串为空返回true，否则返回false

### func IsNotEmpty

```go
func IsNotEmpty(s string) bool
```

IsNotEmpty 检查字符串是否不为空 用于验证字符串在去除首尾空格后是否包含有效内容

**参数:**
- `s`: 待检查的字符串

**返回:**
- `bool`: 字符串不为空返回true，否则返回false

### func Join

```go
func Join(parts ...string) string
```

Join 拼接多个字符串

**参数:**
- `parts`: 要拼接的字符串切片

**返回:**
- `string`: 拼接后的字符串

### func JoinNonEmpty

```go
func JoinNonEmpty(sep string, parts ...string) string
```

JoinNonEmpty 使用分隔符拼接非空字符串

**参数:**
- `sep`: 分隔符
- `parts`: 要拼接的字符串切片

**返回:**
- `string`: 用分隔符连接的非空字符串

### func Mask

```go
func Mask(s string, start, end int, maskChar rune) string
```

Mask 字符串掩码处理 (如手机号脱敏)

**参数:**
- `s`: 输入字符串
- `start`: 开始掩码的位置（包含）
- `end`: 结束掩码的位置（不包含）
- `maskChar`: 掩码字符

**返回:**
- `string`: 掩码后的字符串

**使用示例:**

```go
phone := "13812345678"
masked := str.Mask(phone, 3, 7, '*') // 138****5678
```

### func PadLeft

```go
func PadLeft(s string, length int, pad rune) string
```

PadLeft 在字符串左侧填充字符到指定长度

**参数:**
- `s`: 输入字符串
- `length`: 目标长度
- `pad`: 填充字符

**返回:**
- `string`: 填充后的字符串

### func PadRight

```go
func PadRight(s string, length int, pad rune) string
```

PadRight 在字符串右侧填充字符到指定长度

**参数:**
- `s`: 输入字符串
- `length`: 目标长度
- `pad`: 填充字符

**返回:**
- `string`: 填充后的字符串

### func Prefix

```go
func Prefix(s string, n int) string
```

Prefix 获取字符串的前N个字符

**参数:**
- `s`: 输入字符串
- `n`: 要获取的字符数量

**返回:**
- `string`: 前N个字符，如果字符串长度不足N则返回原字符串

### func Repeat

```go
func Repeat(s string, count int) string
```

Repeat 重复字符串N次

**参数:**
- `s`: 要重复的字符串
- `count`: 重复次数

**返回:**
- `string`: 重复后的字符串

### func SafeDeref

```go
func SafeDeref(s *string) string
```

SafeDeref 安全地解引用字符串指针 用于安全地获取字符串指针的值，避免空指针异常

**参数:**
- `s`: 字符串指针

**返回:**
- `string`: 解引用后的字符串，指针为nil时返回空字符串

### func SafeIndex

```go
func SafeIndex(s, substr string) int
```

SafeIndex 安全地查找子字符串的索引

**参数:**
- `s`: 源字符串
- `substr`: 要查找的子字符串

**返回:**
- `int`: 子字符串的索引，未找到返回-1

### func StringSuffix8

```go
func StringSuffix8(s string) string
```

StringSuffix8 从给定字符串中获取最后8个字符。 如果字符串长度小于等于8，则返回原字符串。

**参数:**
- `s`: 输入字符串

**返回:**
- `string`: 字符串的最后8个字符，或原字符串（如果长度不足8），或空字符串（如果输入为空）

### func Suffix

```go
func Suffix(s string, n int) string
```

Suffix 获取字符串的后N个字符

**参数:**
- `s`: 输入字符串
- `n`: 要获取的字符数量

**返回:**
- `string`: 后N个字符，如果字符串长度不足N则返回原字符串

### func Template

```go
func Template(tmpl string, data map[string]string) string
```

Template 简单模板替换

**参数:**
- `tmpl`: 模板字符串，使用 {{key}} 作为占位符
- `data`: 替换数据

**返回:**
- `string`: 替换后的字符串

**使用示例:**

```go
result := str.Template("Hello {{name}}, you are {{age}} years old", map[string]string{
    "name": "Alice",
    "age":  "25",
})
```

### func ToBase64

```go
func ToBase64(s string) string
```

ToBase64 将字符串编码为Base64

**参数:**
- `s`: 输入字符串

**返回:**
- `string`: Base64编码后的字符串

### func Truncate

```go
func Truncate(s string, maxLen int) string
```

Truncate 截断字符串到指定长度

**参数:**
- `s`: 输入字符串
- `maxLen`: 最大长度

**返回:**
- `string`: 截断后的字符串