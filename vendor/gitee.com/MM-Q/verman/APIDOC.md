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

Banner 返回横幅格式

#### func (*Info) CSV

```go
func (i *Info) CSV() string
```

CSV 返回逗号分隔格式

#### func (*Info) Full

```go
func (i *Info) Full() string
```

Full 返回完整格式

#### func (*Info) JSON

```go
func (i *Info) JSON() string
```

JSON 返回JSON格式

#### func (*Info) KV

```go
func (i *Info) KV() string
```

KV 返回键值对格式

#### func (*Info) Lines

```go
func (i *Info) Lines() string
```

Lines 返回分行格式

#### func (*Info) Long

```go
func (i *Info) Long() string
```

Long 返回长格式

#### func (*Info) Short

```go
func (i *Info) Short() string
```

Short 返回短格式

#### func (*Info) Standard

```go
func (i *Info) Standard() string
```

Standard 返回标准格式

#### func (*Info) Table

```go
func (i *Info) Table() string
```

Table 返回表格格式

#### func (*Info) URI

```go
func (i *Info) URI() string
```

URI 返回URI格式

#### func (*Info) Version

```go
func (i *Info) Version() string
```

Version 返回格式为"程序名 version 版本号 平台/架构"的字符串

