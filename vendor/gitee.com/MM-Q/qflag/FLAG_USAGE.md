# 标志使用语法指南

本文档详细说明了QFlag库中标志的使用语法，包括基础语法和高级标志的用法。

## 基础语法

QFlag基于Go标准库的`flag`包，支持以下标志语法：

### 1. 短标志名（单字符）

```bash
# 使用空格分隔
-f value

# 使用等号分隔
-f=value
```

### 2. 长标志名（多字符）

```bash
# 使用空格分隔
--flag value

# 使用等号分隔
--flag=value
```

### 3. 布尔标志

布尔标志不需要指定值，存在即为true：

```bash
# 短布尔标志
-f

# 长布尔标志
--flag

# 可以使用等号显式设置（可选）
--flag=true
--flag=false
```

### 4. 组合使用

```bash
# 同时使用多个标志
-f value1 --flag2 value2 -g value3
```

## 高级标志语法

除了基础标志外，QFlag还提供了多种高级标志类型。

### 1. 时间标志（Duration）

时间标志接受多种时间格式，支持Go标准库time.ParseDuration所支持的所有格式：

```bash
# 基本时间单位
--timeout 1s           # 1秒
--timeout 100ms         # 100毫秒
--timeout 5m           # 5分钟
--timeout 2h           # 2小时

# 支持小数
--timeout 1.5h         # 1.5小时

# 复合时间
--timeout 1h30m        # 1小时30分钟
--timeout 2d5h30m     # 2天5小时30分钟

# 支持负数
--timeout -30s         # 负30秒
```

支持的时间单位：
- `ns` - 纳秒
- `us` (或 `µs`) - 微秒
- `ms` - 毫秒
- `s` - 秒
- `m` - 分钟
- `h` - 小时

### 2. 时间点标志（Time）

时间点标志支持多种常见时间格式，自动检测并解析：

```bash
# RFC3339 格式
--start-time "2006-01-02T15:04:05Z07:00"
--start-time "2006-01-02T15:04:05.999999999Z07:00"

# ISO8601 格式
--start-time "2006-01-02T15:04:05Z"

# RFC1123 格式
--start-time "Mon, 02 Jan 2006 15:04:05 MST"

# 简单日期时间格式
--start-time "2006-01-02 15:04:05"
--start-time "2006/01/02 15:04:05"

# 仅日期
--start-time "2006-01-02"

# 仅时间
--start-time "15:04:05"

# 时间戳格式
--start-time "Jan _2 15:04:05"
--start-time "Jan _2 15:04:05.000"
```

支持的时间格式（按优先级排序）：
- RFC3339 和 RFC3339Nano
- ISO8601 和 ISO8601Nano
- 日期时间格式（2006-01-02 15:04:05）
- 日期格式（2006-01-02）
- 时间格式（15:04:05）
- RFC1123 和 RFC1123Z
- 时间戳格式（Stamp, StampMilli, StampMicro, StampNano）
- RFC822 和 RFC822Z
- 厨房格式（3:04PM）
- 紧凑格式（20060102150405）

### 3. 大小标志（Size）

大小标志接受多种大小格式，支持二进制和十进制单位：

```bash
# 直接指定数字（默认为字节）
--size 1024           # 1024字节

# 十进制单位（1000进制）
--size 1KB            # 1千字节 (1000字节)
--size 10MB           # 10兆字节 (1000^2字节)
--size 2GB            # 2吉字节 (1000^3字节)

# 二进制单位（1024进制）
--size 1KiB           # 1二进制千字节 (1024字节)
--size 10MiB          # 10二进制兆字节 (1024^2字节)
--size 2GiB           # 2二进制吉字节 (1024^3字节)

# 使用小数
--size 1.5MB         # 1.5兆字节

# 大小写不敏感
--size 1mb            # 与 1MB 相同
```

支持的大小单位：
- `B/b` - 字节
- `KB/kb/K/k` - 十进制千字节 (1000字节)
- `MB/mb/M/m` - 十进制兆字节 (1000^2字节)
- `GB/gb/G/g` - 十进制吉字节 (1000^3字节)
- `TB/tb/T/t` - 十进制太字节 (1000^4字节)
- `PB/pb/P/p` - 十进制拍字节 (1000^5字节)
- `KiB/kib` - 二进制千字节 (1024字节)
- `MiB/mib` - 二进制兆字节 (1024^2字节)
- `GiB/gib` - 二进制吉字节 (1024^3字节)
- `TiB/tib` - 二进制太字节 (1024^4字节)
- `PiB/pib` - 二进制拍字节 (1024^5字节)

### 4. 切片标志（Slice）

切片标志可以接受多个值，支持两种输入方式：

```bash
# 多次使用同一标志
--file file1.txt --file file2.txt --file file3.txt

# 使用逗号分隔
--file file1.txt,file2.txt,file3.txt

# 混合使用
--file file1.txt --file file2.txt,file3.txt --file file4.txt

# 空值处理
--file ""             # 设置为空切片
```

支持的切片类型：
- `StringSlice` - 字符串切片
- `IntSlice` - 整数切片
- `Int64Slice` - 64位整数切片

### 5. 映射标志（Map）

映射标志接受键值对，支持两种输入方式：

```bash
# 使用等号分隔键值
--param key1=value1 --param key2=value2

# 使用逗号分隔多个键值对
--param key1=value1,key2=value2,key3=value3

# 混合使用
--param key1=value1 --param key2=value2,key3=value3

# 空值处理
--param ""             # 设置为空映射
--param ",,,"          # 空对会被跳过
```

### 6. 枚举标志（Enum）

枚举标志限制为预定义的值，使用映射表实现O(1)时间复杂度的值查找：

```bash
# 使用预定义的值
--mode debug          # 正确
--mode production     # 正确
--mode invalid         # 错误，会报错
```

枚举标志特性：
- 使用映射表进行快速值验证
- 不允许空字符串作为枚举值
- 默认值必须在允许值列表中
- 不允许设置空值

## 标志定义示例

```go
// 基础标志
name := cmd.String("name", "n", "用户名", "default")
age := cmd.Int("age", "a", "年龄", 18)
enabled := cmd.Bool("enabled", "e", "是否启用", false)

// 数值类型标志
port := cmd.Uint("port", "p", "端口号", 8080)
timeout := cmd.Duration("timeout", "t", "超时时间", time.Second*30)

// 时间和大小标志
startTime := cmd.Time("start-time", "s", "开始时间", time.Now())
maxSize := cmd.Size("limit", "l", "大小限制", 10*1024*1024) // 10MB

// 集合类型标志
files := cmd.StringSlice("files", "f", "文件列表", []string{})
tags := cmd.IntSlice("tags", "", "标签列表", []int{})
config := cmd.Map("config", "c", "配置参数", map[string]string{})

// 枚举标志
mode := cmd.Enum("mode", "m", "运行模式", []string{"debug", "release"}, "release")
```

## 高级功能

### 1. 时间标志高级功能

```go
// 获取时间标志使用的格式
format := timeFlag.GetFormat()

// 使用相同格式格式化其他时间值
formatted := timeFlag.FormatTime(time.Now())
```

### 2. 切片标志高级功能

```go
// 获取切片长度
length := sliceFlag.Length()

// 检查切片是否为空
empty := sliceFlag.IsEmpty()

// 添加元素
sliceFlag.Append("new-item")

// 清空切片
sliceFlag.Clear()
```

### 3. 映射标志高级功能

```go
// 设置键值对
mapFlag.SetKV("key", "value")

// 获取值
value := mapFlag.Get("key")

// 检查键是否存在
exists := mapFlag.Has("key")

// 删除键值对
mapFlag.Delete("key")

// 清空映射
mapFlag.Clear()
```

## 使用建议

1. **一致性** - 在整个项目中保持一致的命名和风格
2. **默认值** - 为标志提供合理的默认值
3. **帮助文本** - 编写清晰、简洁的帮助文本
4. **验证** - 使用验证器确保输入值的有效性
5. **分组** - 使用互斥组和必需组组织相关标志

## 常见错误

1. **未识别的标志** - 检查标志名称是否正确
2. **缺少参数** - 非布尔标志需要提供值
3. **格式错误** - 检查时间、大小等特殊格式的正确性
4. **类型不匹配** - 确保提供的值可以转换为标志类型
5. **枚举值无效** - 确保枚举标志的值在预定义列表中

## 高级技巧

1. **环境变量** - 可以绑定环境变量作为默认值
2. **子命令** - 使用子命令组织复杂的功能
3. **条件解析** - 根据某些标志的存在决定其他标志的行为
4. **自定义验证** - 添加自定义验证逻辑确保数据完整性
5. **时间格式** - 使用TimeFlag的GetFormat方法获取解析格式，用于格式化其他时间值