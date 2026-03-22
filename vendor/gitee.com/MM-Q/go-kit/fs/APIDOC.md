# Package fs

```go
import "gitee.com/MM-Q/go-kit/fs"
```

## Functions

### func Collect

```go
func Collect(targetPath string, recursive bool) ([]string, error)
```

Collect 收集指定路径下的所有文件 用于收集文件或目录中的文件，支持通配符匹配和递归遍历

**参数:**
- `targetPath`: 目标路径，支持通配符(*?[]{})
- `recursive`: 是否递归遍历目录

**返回:**
- `[]string`: 收集到的文件路径切片
- `error`: 收集失败时返回错误

### func Copy

```go
func Copy(src, dst string) error
```

Copy 通用复制函数，自动判断源路径类型并调用相应的复制函数。支持复制普通文件、目录、符号链接和特殊文件（设备文件、命名管道等）。默认不覆盖已存在的目标文件/目录。

**特性:**
- 自动识别文件类型（普通文件、目录、符号链接、特殊文件）
- 保持文件权限和目录结构
- 安全的备份恢复机制，失败时自动回滚
- 对空文件进行性能优化
- 使用临时文件+原子重命名确保数据安全
- 智能路径处理：如果目标是已存在的目录，自动追加源文件名/目录名

**参数:**
- `src`: 源路径（支持文件、目录、符号链接、特殊文件）
- `dst`: 目标路径

**返回:**
- `error`: 复制失败时返回错误，如果目标已存在则返回错误

**示例:**

```go
// 精确路径模式
Copy("a.txt", "b.txt")           // 创建 b.txt
Copy("dirA", "dirB")             // 创建 dirB/

// 自动追加文件名/目录名
Copy("a.txt", "existingDir")     // 创建 existingDir/a.txt
Copy("dirA", "existingDir")      // 创建 existingDir/dirA/

// 自动创建父目录
Copy("a.txt", "newDir/b.txt")    // 创建 newDir/b.txt
Copy("dirA", "newDir/subDir")    // 创建 newDir/subDir/
```

### func CopyEx

```go
func CopyEx(src, dst string, overwrite bool) error
```

CopyEx 通用复制函数（可控制是否覆盖），自动判断源路径类型并调用相应的复制函数。支持复制普通文件、目录、符号链接和特殊文件（设备文件、命名管道等）。

**特性:**
- 自动识别文件类型（普通文件、目录、符号链接、特殊文件）
- 保持文件权限和目录结构
- 安全的备份恢复机制，失败时自动回滚
- 对空文件进行性能优化
- 使用临时文件+原子重命名确保数据安全
- 可控制覆盖行为
- 智能路径处理：如果目标是已存在的目录，自动追加源文件名/目录名

**参数:**
- `src`: 源路径（支持文件、目录、符号链接、特殊文件）
- `dst`: 目标路径
- `overwrite`: 是否允许覆盖已存在的目标文件/目录

**返回:**
- `error`: 复制失败时返回错误

**示例:**

```go
// 不覆盖已存在的目标
CopyEx("a.txt", "b.txt", false)

// 覆盖已存在的目标
CopyEx("a.txt", "b.txt", true)

// 智能路径处理
CopyEx("a.txt", "existingDir", true)  // 创建/覆盖 existingDir/a.txt
CopyEx("dirA", "existingDir", true)   // 创建/覆盖 existingDir/dirA/
```

### func Move

```go
func Move(src, dst string) error
```

Move 通用移动函数，将文件或目录移动到目标位置。支持移动普通文件、目录、符号链接和特殊文件。默认不覆盖已存在的目标文件/目录。优先使用 os.Rename（同文件系统内），失败时降级使用复制+删除（支持跨文件系统）。

**特性:**
- 自动识别文件类型（普通文件、目录、符号链接、特殊文件）
- 保持文件权限和目录结构
- 优先使用 os.Rename（原子操作，同文件系统内高效）
- 失败时降级使用 CopyEx + os.RemoveAll（支持跨文件系统）
- 智能路径处理：如果目标是已存在的目录，自动追加源文件名/目录名

**参数:**
- `src`: 源路径（支持文件、目录、符号链接、特殊文件）
- `dst`: 目标路径

**返回:**
- `error`: 移动失败时返回错误，如果目标已存在则返回错误

**示例:**

```go
// 精确路径模式
Move("a.txt", "b.txt")           // 移动到 b.txt
Move("dirA", "dirB")             // 移动到 dirB/

// 自动追加文件名/目录名
Move("a.txt", "existingDir")     // 移动到 existingDir/a.txt
Move("dirA", "existingDir")      // 移动到 existingDir/dirA/

// 自动创建父目录
Move("a.txt", "newDir/b.txt")    // 移动到 newDir/b.txt
Move("dirA", "newDir/subDir")    // 移动到 newDir/subDir/
```

### func MoveEx

```go
func MoveEx(src, dst string, overwrite bool) error
```

MoveEx 通用移动函数（可控制是否覆盖），将文件或目录移动到目标位置。支持移动普通文件、目录、符号链接和特殊文件。优先使用 os.Rename（同文件系统内），失败时降级使用复制+删除（支持跨文件系统）。

**特性:**
- 自动识别文件类型（普通文件、目录、符号链接、特殊文件）
- 保持文件权限和目录结构
- 优先使用 os.Rename（原子操作，同文件系统内高效）
- 失败时降级使用 CopyEx + os.RemoveAll（支持跨文件系统）
- 可控制覆盖行为
- 智能路径处理：如果目标是已存在的目录，自动追加源文件名/目录名

**参数:**
- `src`: 源路径（支持文件、目录、符号链接、特殊文件）
- `dst`: 目标路径
- `overwrite`: 是否允许覆盖已存在的目标文件/目录

**返回:**
- `error`: 移动失败时返回错误

**示例:**

```go
// 不覆盖已存在的目标
MoveEx("a.txt", "b.txt", false)

// 覆盖已存在的目标
MoveEx("a.txt", "b.txt", true)

// 智能路径处理
MoveEx("a.txt", "existingDir", true)  // 移动/覆盖到 existingDir/a.txt
MoveEx("dirA", "existingDir", true)   // 移动/覆盖到 existingDir/dirA/

// 跨文件系统移动（自动降级到复制+删除）
MoveEx("/mnt/disk1/file.txt", "/mnt/disk2/file.txt", true)
```

### func Exists

```go
func Exists(path string) bool
```

Exists 检查指定路径的文件或目录是否存在 用于验证文件系统中指定路径是否存在，权限错误等异常情况视为不存在

**参数:**
- `path`: 要检查的路径

**返回:**
- `bool`: 文件或目录存在返回true，否则返回false

### func GetDefaultBinPath

```go
func GetDefaultBinPath() string
```

GetDefaultBinPath 返回默认bin路径 用于获取Go程序的默认bin路径，采用多级回退策略确保总能返回有效路径

**返回:**
- `string`: 默认bin路径，优先级为GOPATH/bin > 用户主目录/go/bin > 当前工作目录/bin

### func GetExecutablePath

```go
func GetExecutablePath() string
```

GetExecutablePath 获取程序的绝对安装路径 用于获取当前可执行文件的绝对路径，提供多级降级策略确保总能返回路径

**返回:**
- `string`: 程序的绝对路径，失败时降级为相对路径

### func GetSize

```go
func GetSize(path string) (int64, error)
```

GetSize 获取文件或目录的大小 用于计算文件或目录的总字节数，目录会递归计算所有普通文件的大小

**参数:**
- `path`: 文件或目录路径

**返回:**
- `int64`: 文件或目录的总大小(字节)
- `error`: 路径不存在或访问失败时返回错误

### func GetUserHomeDir

```go
func GetUserHomeDir() string
```

GetUserHomeDir 获取用户家目录 用于获取用户家目录路径，提供多级降级策略确保总能返回有效路径

**返回:**
- `string`: 用户家目录路径，失败时依次降级为工作目录或当前目录

### func IsDir

```go
func IsDir(path string) bool
```

IsDir 检查指定路径是否为目录 用于验证指定路径是否为目录

**参数:**
- `path`: 要检查的路径

**返回:**
- `bool`: 是目录返回true，否则返回false

### func IsFile

```go
func IsFile(path string) bool
```

IsFile 检查指定路径是否为文件 用于验证指定路径是否为普通文件

**参数:**
- `path`: 要检查的路径

**返回:**
- `bool`: 是文件返回true，否则返回false

### func IsHidden

```go
func IsHidden(path string) bool
```

IsHidden 判断文件或目录是否为隐藏 用于跨平台检查文件或目录的隐藏属性

**参数:**
- `path`: 文件或目录路径

**返回:**
- `bool`: 文件为隐藏返回true，否则返回false

### func IsOctPerm

```go
func IsOctPerm(permission string) bool
```

IsOctPerm 检查输入的权限是否是合法的4位八进制数 用于验证权限字符串格式是否符合八进制权限规范

**参数:**
- `permission`: 输入的权限字符串，例如 "0755" 或 "0644"

**返回:**
- `bool`: 权限格式合法返回true，否则返回false

### func IsReadOnly

```go
func IsReadOnly(path string) bool
```

IsReadOnly 判断文件或目录是否为只读 用于跨平台检查文件或目录的只读属性

**参数:**
- `path`: 文件或目录路径

**返回:**
- `bool`: 文件为只读返回true，否则返回false

### func OctStrToMode

```go
func OctStrToMode(octalStr string) (os.FileMode, error)
```

OctStrToMode 将4位八进制字符串权限转换为 os.FileMode 类型 用于将八进制权限字符串转换为Go标准库的文件权限类型

**参数:**
- `octalStr`: 4位八进制字符串，例如 "0755" 或 "0644"

**返回:**
- `os.FileMode`: 转换后的文件权限
- `error`: 输入不合法时返回错误

**示例:**

```go
mode, err := OctStrToMode("0755")
if err != nil {
    log.Fatal(err)
}
fmt.Println(mode) // 输出: -rwxr-xr-x
```

