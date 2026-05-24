# Package hash

```go
import "gitee.com/MM-Q/go-kit/hash"
```

## Types

### type Algorithm

```go
type Algorithm string
```

Algorithm 定义支持的哈希算法类型。使用自定义类型可以在编译时进行类型检查，避免运行时错误。

**常量:**
- `MD5` - MD5 算法
- `SHA1` - SHA1 算法
- `SHA256` - SHA256 算法
- `SHA512` - SHA512 算法

## Functions

### func ParseAlgorithm

```go
func ParseAlgorithm(algorithm string) (Algorithm, error)
```

ParseAlgorithm 将字符串解析为 Algorithm 类型。

**参数:**
- `algorithm`: 哈希算法名称字符串（如 "md5", "sha256"）。

**返回:**
- `Algorithm`: 对应的算法类型常量。
- `error`: 如果不支持该算法，则返回错误。

**使用示例:**

```go
algo, err := hash.ParseAlgorithm("sha256")
if err != nil {
    return err
}
result, err := hash.Checksum(filename, algo)
```

### func IsAlgorithmSupported

```go
func IsAlgorithmSupported(algorithm string) bool
```

IsAlgorithmSupported 检查给定的哈希算法名称是否受支持。

**参数:**
- `algorithm`: 要检查的哈希算法名称（字符串形式）。

**返回:**
- `bool`: 如果算法受支持则返回 true，否则返回 false。

### func Checksum

```go
func Checksum(filePath string, algorithm Algorithm) (string, error)
```

Checksum 计算文件哈希值。

**参数:**
- `filePath`: 文件路径
- `algorithm`: 哈希算法类型（如 `hash.MD5`, `hash.SHA256`）

**返回:**
- `string`: 文件的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 根据文件大小动态分配缓冲区以提高性能
- 支持任何实现hash.Hash接口的哈希算法
- 使用io.CopyBuffer进行高效的文件读取和哈希计算

**使用示例:**

```go
// 使用常量
result, err := hash.Checksum("file.txt", hash.SHA256)

// 从字符串解析
algo, _ := hash.ParseAlgorithm("md5")
result, err := hash.Checksum("file.txt", algo)
```

### func ChecksumProgress

```go
func ChecksumProgress(filePath string, algorithm Algorithm) (string, error)
```

ChecksumProgress 计算文件哈希值(带进度条)。

**参数:**
- `filePath`: 文件路径
- `algorithm`: 哈希算法类型

**返回:**
- `string`: 文件的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 根据文件大小动态分配缓冲区以提高性能
- 支持任何实现hash.Hash接口的哈希算法
- 使用io.CopyBuffer进行高效的文件读取和哈希计算

### func HashData

```go
func HashData(data []byte, algorithm Algorithm) (string, error)
```

HashData 计算内存数据哈希值。

**参数:**
- `data`: 要计算哈希的字节数据
- `algorithm`: 哈希算法类型

**返回:**
- `string`: 数据的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 直接在内存中计算，无需文件I/O，性能更高
- 支持任何大小的数据，包括空数据
- 使用标准库优化的hash实现，性能最佳

### func HashString

```go
func HashString(data string, algorithm Algorithm) (string, error)
```

HashString 计算字符串哈希值（便利函数）。

**参数:**
- `data`: 要计算哈希的字符串
- `algorithm`: 哈希算法类型

**返回:**
- `string`: 字符串的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 这是HashData的便利包装函数
- 内部将字符串转换为字节切片进行处理
- 适用于文本数据、配置字符串等场景

### func HashReader

```go
func HashReader(reader io.Reader, algorithm Algorithm) (string, error)
```

HashReader 计算io.Reader数据哈希值。

**参数:**
- `reader`: 数据源读取器
- `algorithm`: 哈希算法类型

**返回:**
- `string`: 读取数据的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 适用于流式数据处理，如网络数据、管道数据等
- 使用缓冲区进行高效读取，避免频繁的小块读取
- 会完全消费Reader中的数据
- 使用对象池优化内存分配

## Usage Examples

### CLI Tool Integration

```go
package main

import (
    "flag"
    "fmt"
    "gitee.com/MM-Q/go-kit/hash"
)

func main() {
    var algoStr string
    flag.StringVar(&algoStr, "algo", "sha256", "Hash algorithm (md5, sha1, sha256, sha512)")
    flag.Parse()
    
    // Parse and validate algorithm
    algo, err := hash.ParseAlgorithm(algoStr)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Calculate hash
    result, err := hash.Checksum("file.txt", algo)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Hash: %s\n", result)
}
```

### Direct Usage with Constants

```go
// Calculate MD5 hash of a file
md5Hash, err := hash.Checksum("data.txt", hash.MD5)

// Calculate SHA256 hash of string data
sha256Hash, err := hash.HashString("Hello, World!", hash.SHA256)

// Calculate hash with progress bar
hashWithProgress, err := hash.ChecksumProgress("largefile.zip", hash.SHA512)
```
