# Package verman

```go
import "gitee.com/MM-Q/verman"
```

## TYPES

### type Info

```go
type Info struct {
	AppName       string // 应用程序名称
	GitVersion    string // Git 语义化版本号
	GitCommit     string // Git 提交哈希值
	GitTreeState  string // Git 仓库状态
	GitCommitTime string // Git 提交时间
	BuildTime     string // 构建时间
	GoVersion     string // Go 运行时版本
	Platform      string // 平台信息
}
```

Info 版本信息结构体

### var V

```go
var V *Info
```

V 全局版本信息实例，供外部使用

#### func (*Info) Banner

```go
func (i *Info) Banner() string
```

Banner 返回横幅格式(多行)

#### func (*Info) Build

```go
func (i *Info) Build() string
```

Build 返回构建信息格式

#### func (*Info) Complete

```go
func (i *Info) Complete() string
```

Complete 返回包含所有信息的完整字符串"程序名 v1.0.0 linux/amd64 (commit: abc1234, tree: clean, built: 2024-01-01T12:00:00Z, go: go1.19)"

#### func (*Info) Detail

```go
func (i *Info) Detail() string
```

Detail 返回格式为"程序名 v1.0.0 linux/amd64 built at 2024-01-01"的字符串

#### func (*Info) Full

```go
func (i *Info) Full() string
```

Full 返回格式为"程序名 version 版本号 平台/架构 (commit: abc1234)"的字符串

#### func (*Info) Git

```go
func (i *Info) Git() string
```

Git 返回Git信息格式

#### func (*Info) JSON

```go
func (i *Info) JSON() string
```

JSON 返回JSON格式

#### func (*Info) Simple

```go
func (i *Info) Simple() string
```

Simple 返回格式为"程序名 v1.0.0"的字符串

#### func (*Info) Table

```go
func (i *Info) Table() string
```

Table 返回表格格式(多行)

#### func (*Info) Version

```go
func (i *Info) Version() string
```

Version 返回格式为"程序名 version 版本号 平台/架构"的字符串

