# Package hash

```go
import "gitee.com/MM-Q/go-kit/hash"
```

## Constants

### 哈希算法常量

推荐使用这些常量以避免拼写错误，也可以直接传入字符串（如 `"sha256"`）：

```go
const (
    MD5    = "md5"
    SHA1   = "sha1"
    SHA256 = "sha256"
    SHA512 = "sha512"
)
```

## Functions

### func IsAlgorithmSupported

```go
func IsAlgorithmSupported(algorithm string) bool
```

检查给定的哈希算法名称是否受支持。

**参数:**
- `algorithm`: 要检查的哈希算法名称（字符串形式）。

**返回:**
- `bool`: 如果算法受支持则返回 true，否则返回 false。

**使用示例:**

```go
if hash.IsAlgorithmSupported("sha256") {
    // 支持该算法
}
```

### func Checksum

```go
func Checksum(filePath string, algorithm string) (string, error)
```

计算文件哈希值。

**参数:**
- `filePath`: 文件路径
- `algorithm`: 哈希算法名称（如 `"md5"`, `"sha256"`，推荐使用 `hash.MD5`/`hash.SHA256` 等常量）

**返回:**
- `string`: 文件的十六进制哈希值
- `error`: 错误信息，如果计算失败

**注意:**
- 根据文件大小动态分配缓冲区以提高性能
- 使用 io.CopyBuffer 进行高效的文件读取和哈希计算

**使用示例:**

```go
// 使用常量（推荐）
result, err := hash.Checksum("file.txt", hash.SHA256)

// 直接传入字符串
result, err := hash.Checksum("file.txt", "sha256")
```

### func ChecksumProgress

```go
func ChecksumProgress(filePath string, algorithm string) (string, error)
```

计算文件哈希值（带进度条）。

**参数:**
- `filePath`: 文件路径
- `algorithm`: 哈希算法名称

**返回:**
- `string`: 文件的十六进制哈希值
- `error`: 错误信息，如果计算失败

### func HashData

```go
func HashData(data []byte, algorithm string) (string, error)
```

计算内存数据哈希值。

**参数:**
- `data`: 要计算哈希的字节数据
- `algorithm`: 哈希算法名称

**返回:**
- `string`: 数据的十六进制哈希值
- `error`: 错误信息，如果计算失败

### func HashString

```go
func HashString(data string, algorithm string) (string, error)
```

计算字符串哈希值（便利函数）。

**参数:**
- `data`: 要计算哈希的字符串
- `algorithm`: 哈希算法名称

**返回:**
- `string`: 字符串的十六进制哈希值
- `error`: 错误信息，如果计算失败

### func HashReader

```go
func HashReader(reader io.Reader, algorithm string) (string, error)
```

计算 io.Reader 数据哈希值。

**参数:**
- `reader`: 数据源读取器
- `algorithm`: 哈希算法名称

**返回:**
- `string`: 读取数据的十六进制哈希值
- `error`: 错误信息，如果计算失败

### func Verify

```go
func Verify(filePath, expectedHash string, algorithm string) (bool, error)
```

校验指定路径文件的哈希值是否与预期值匹配。

**参数:**
- `filePath`: 文件路径
- `expectedHash`: 预期的十六进制哈希值
- `algorithm`: 哈希算法名称

**返回:**
- `bool`: 哈希值是否匹配
- `error`: 错误信息，如果计算失败

**注意:**
- `expectedHash` 比较时忽略大小写
- 典型场景：下载文件后校验完整性

**使用示例:**

```go
expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
ok, err := hash.Verify("/path/to/downloaded.zip", expected, hash.SHA256)
if err != nil {
    log.Fatal(err)
}
if !ok {
    log.Fatal("文件校验失败！")
}
```

### func Compare

```go
func Compare(path1, path2 string, algorithm string) (bool, error)
```

校验两个路径的文件内容是否一致（通过比较哈希值）。

**参数:**
- `path1`: 第一个文件路径
- `path2`: 第二个文件路径
- `algorithm`: 哈希算法名称

**返回:**
- `bool`: 两个文件内容是否一致
- `error`: 错误信息，如果计算失败

**使用示例:**

```go
ok, err := hash.Compare("/src/file.bin", "/dst/file.bin", hash.SHA256)
if err != nil {
    log.Fatal(err)
}
if !ok {
    log.Fatal("文件不一致！")
}
```

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

    // 直接传入字符串，无需类型转换
    result, err := hash.Checksum(flag.Arg(0), algoStr)
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

// Verify file integrity after download
ok, err := hash.Verify("downloaded.zip", "e3b0c442...", hash.SHA256)

// Compare two files
same, err := hash.Compare("file1.bin", "file2.bin", hash.SHA256)
```
