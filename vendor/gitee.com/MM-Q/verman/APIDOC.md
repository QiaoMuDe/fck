# Package verman

`import "gitee.com/MM-Q/verman"`

Package verman 提供应用程序版本信息管理功能。

---

## VARIABLES

```go
var Info = VerMan{}
```

定义一个全局变量，用于存储应用程序版本信息。

---

## FUNCTIONS

### GetAppName

```go
func GetAppName() string
```

GetAppName 返回应用程序的名称。

---

### Get

```go
func Get() *VerMan
```

Get 返回包含当前应用程序版本信息的 `VerMan` 结构体。

**返回值**  
- `*VerMan`：包含当前应用程序版本信息的 `VerMan` 结构体。

---

## TYPES

### VerMan

```go
type VerMan struct {
    AppName       string `json:"appName"`       // 应用程序名称
    GitVersion    string `json:"gitVersion"`    // Git 语义化版本号(如 v1.0.0)
    GitCommit     string `json:"gitCommit"`     // Git 提交哈希值(如 abc1234)
    GitTreeState  string `json:"gitTreeState"`  // Git 仓库状态(如 clean, dirty)
    GitCommitTime string `json:"gitCommitTime"` // Git 提交时间(如 2024-01-01T12:00:00Z)
    BuildTime     string `json:"buildTime"`     // 构建时间(如 2024-01-01T12:00:00Z)
    GoVersion     string `json:"goVersion"`     // Go 运行时版本(如 go1.19)
    Platform      string `json:"platform"`      // 平台信息(如 linux/amd64)
}
```

VerMan 包含版本信息的结构体。

---

### Methods (on VerMan)

#### JSON

```go
func (v VerMan) JSON() (string, error)
```

返回 `VerMan` 的 JSON 字符串表示。

**返回值**  
- `string`：JSON 格式的字符串表示  
- `error`：转换错误信息

---

#### PrintVersion

```go
func (v VerMan) PrintVersion(format string)
```

根据指定的格式打印版本信息。

**参数**  
- `format`：输出格式，支持 `"json"`、`"text"` 和 `"simple"`。

| 格式值   | 说明                 |
| -------- | -------------------- |
| `json`   | 以 JSON 格式输出     |
| `text`   | 以完整详细格式输出   |
| `simple` | 以简洁格式输出       |

> 如果传入的格式不是支持的格式，将打印错误信息。

---

#### SprintVersion

```go
func (v VerMan) SprintVersion(format string) (string, error)
```

根据指定的格式返回版本信息字符串。

**参数**  
- `format`：返回格式，支持 `"json"`、`"text"` 和 `"simple"`。

| 格式值   | 说明                 |
| -------- | -------------------- |
| `json`   | 以 JSON 格式返回     |
| `text`   | 以完整详细格式返回   |
| `simple` | 以简洁格式返回       |

**返回值**  
- `string`：按指定格式生成的版本信息  
- `error`：如果 `format` 不在支持范围内，则返回错误

---

#### String

```go
func (v VerMan) String() string
```

返回版本信息的字符串表示。