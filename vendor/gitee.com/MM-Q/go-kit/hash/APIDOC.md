# Package hash

```go
import "gitee.com/MM-Q/go-kit/hash"
```

## Functions

### func Checksum

```go
func Checksum(filePath string, algorithm string) (string, error)
```

Checksum 计算文件哈希值

**参数:**
- `filePath`: 文件路径
- `algorithm`: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）

**返回:**
- `string`: 文件的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 根据文件大小动态分配缓冲区以提高性能
- 支持任何实现hash.Hash接口的哈希算法
- 使用io.CopyBuffer进行高效的文件读取和哈希计算

### func ChecksumProgress

```go
func ChecksumProgress(filePath string, algorithm string) (string, error)
```

ChecksumProgress 计算文件哈希值(带进度条)

**参数:**
- `filePath`: 文件路径
- `algorithm`: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）

**返回:**
- `string`: 文件的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 根据文件大小动态分配缓冲区以提高性能
- 支持任何实现hash.Hash接口的哈希算法
- 使用io.CopyBuffer进行高效的文件读取和哈希计算

### func HashData

```go
func HashData(data []byte, algorithm string) (string, error)
```

HashData 计算内存数据哈希值

**参数:**
- `data`: 要计算哈希的字节数据
- `algorithm`: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）

**返回:**
- `string`: 数据的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 直接在内存中计算，无需文件I/O，性能更高
- 支持任何大小的数据，包括空数据
- 使用标准库优化的hash实现，性能最佳

### func HashReader

```go
func HashReader(reader io.Reader, algorithm string) (string, error)
```

HashReader 计算io.Reader数据哈希值

**参数:**
- `reader`: 数据源读取器
- `algorithm`: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）

**返回:**
- `string`: 读取数据的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 适用于流式数据处理，如网络数据、管道数据等
- 使用缓冲区进行高效读取，避免频繁的小块读取
- 会完全消费Reader中的数据
- 使用对象池优化内存分配

### func HashString

```go
func HashString(data string, algorithm string) (string, error)
```

HashString 计算字符串哈希值（便利函数）

**参数:**
- `data`: 要计算哈希的字符串
- `algorithm`: 哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）

**返回:**
- `string`: 字符串的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 这是HashData的便利包装函数
- 内部将字符串转换为字节切片进行处理
- 适用于文本数据、配置字符串等场景

### func IsAlgorithmSupported

```go
func IsAlgorithmSupported(algorithm string) bool
```

IsAlgorithmSupported 检查给定的哈希算法名称是否受支持。 匹配时会忽略算法名称的大小写。

**参数:**
- `algorithm`: 要检查的哈希算法名称（如 "md5", "sha1", "sha256", "sha512"）。

**返回:**
- `bool`: 如果算法受支持则返回 true，否则返回 false。

