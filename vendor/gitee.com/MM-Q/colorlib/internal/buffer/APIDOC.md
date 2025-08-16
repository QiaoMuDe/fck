# Package buffer

Package buffer 提供了高性能的字符串构建器对象池实现。该文件实现了 `BuilderPool` 结构体，用于管理 `strings.Builder` 对象的复用，支持智能缩容、统计信息收集等功能，有效减少内存分配和垃圾回收压力。

## CONSTANTS

```go
const (
    DefaultMaxSize = 64 * 1024 // 64KB 默认最大大小
)
```

常量定义默认大小。

## TYPES

### BuilderPool

```go
type BuilderPool struct {
    // Has unexported fields.
}
```

`BuilderPool` 字符串构建器对象池实现。

### Functions

#### NewBuilderPool

```go
func NewBuilderPool(maxSize int) *BuilderPool
```

`NewBuilderPool` 创建新的字符串构建器池。

- 参数:
  - `maxSize`: 最大构建器大小限制，超过此大小的构建器不会被放回池中
- 返回值:
  - `*BuilderPool`: 字符串构建器池实例

#### NewDefaultBuilderPool

```go
func NewDefaultBuilderPool() *BuilderPool
```

`NewDefaultBuilderPool` 创建默认配置的字符串构建器池。

- 返回值:
  - `*BuilderPool`: 使用默认配置的字符串构建器池实例

### Methods

#### BuildString

```go
func (bp *BuilderPool) BuildString(fn func(*strings.Builder)) string
```

`BuildString` 构建字符串。

- 参数:
  - `fn`: 构建函数，接收字符串构建器作为参数
- 返回值:
  - `string`: 构建的字符串

#### Get

```go
func (bp *BuilderPool) Get() *strings.Builder
```

`Get` 从池中获取字符串构建器。

- 返回值:
  - `*strings.Builder`: 干净的字符串构建器实例

#### Put

```go
func (bp *BuilderPool) Put(builder *strings.Builder)
```

`Put` 将字符串构建器归还到池中（智能缩容优化版）。

- 参数:
  - `builder`: 要归还的字符串构建器

#### WithBuilder

```go
func (bp *BuilderPool) WithBuilder(fn func(*strings.Builder))
```

`WithBuilder` 使用字符串构建器执行函数。

- 参数:
  - `fn`: 要执行的函数，接收字符串构建器作为参数

### Pool

```go
type Pool interface {
    Get() *strings.Builder                     // 获取字符串构建器
    Put(*strings.Builder)                      // 归还字符串构建器
    WithBuilder(func(*strings.Builder))        // 使用字符串构建器
    BuildString(func(*strings.Builder)) string // 构建字符串
}
```

`Pool` 字符串构建器池接口。