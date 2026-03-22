# Utils 包 API 文档

```go
package utils // import "gitee.com/MM-Q/qflag/internal/utils"
```

---

## FUNCTIONS

### CalcOptionMaxWidth()

```go
func CalcOptionMaxWidth(options []types.OptionInfo) int
```

#### CalcOptionMaxWidth 计算选项名称最大宽度

**参数:**
  - options: 选项信息列表

**返回值:**
  - int: 选项名称最大宽度

---

### CalcSubCmdMaxLen()

```go
func CalcSubCmdMaxLen(subCmds []types.SubCmdInfo) int
```

#### CalcSubCmdMaxLen 计算子命令名称最大宽度

**参数:**
  - subCmds: 子命令信息列表

**返回值:**
  - int: 子命令名称最大宽度

---

### FormatDefaultValue()

```go
func FormatDefaultValue(flagType types.FlagType, defValue any) string
```

#### FormatDefaultValue 返回格式化后的默认值

**参数:**
  - flagType: 标志类型
  - defValue: 默认值

**返回值:**
  - string: 格式化后的默认值

---

### FormatFlagName()

```go
func FormatFlagName(longName, shortName string) string
```

#### FormatFlagName 格式化标志名称

**参数:**
  - longName: 长标志名称
  - shortName: 短标志名称

**返回值:**
  - string: 格式化后的标志名称

---

### FormatSize()

```go
func FormatSize(size int64) string
```

#### FormatSize 格式化大小值为带单位的字符串

**参数:**
  - size: 大小值 (字节) 

**返回值:**
  - string: 格式化后的大小字符串, 如 "10KB", "4MB"

**功能说明: **
  - 自动选择合适的单位 (B, KB, MB, GB, TB) 
  - 保留两位小数
  - 对于小于1KB的值, 直接显示字节数

---

### GetCmdName()

```go
func GetCmdName(cmd types.Command) string
```

#### GetCmdName 返回命令名称

**参数:**
  - cmd: 要获取名称的命令

**返回值:**
  - string: 命令名称

---

### SortOptions()

```go
func SortOptions(options []types.OptionInfo)
```

#### SortOptions 对选项列表进行排序

**排序规则:** 长短都有 > 仅长选项 > 仅短选项, 每组内按首字母排序 (忽略大小写)

**参数:**
  - options: 选项信息列表

---

### SortSubCmds()

```go
func SortSubCmds(subCmds []types.SubCmdInfo)
```

#### SortSubCmds 对子命令列表进行排序

**排序规则:** 长短名都有 > 仅单一名字, 每组内按首字母排序 (忽略大小写)

**参数:**
  - subCmds: 子命令信息列表

---

### ToStrSlice()

```go
func ToStrSlice(slice any) ([]string, error)
```

#### ToStrSlice 将任何类型的切片转换为字符串切片

**主要用途:** 主要用于将数字类型切片转换为字符串切片, 以便用于 EnumFlag

**参数:**
  - slice: 任意类型的切片

**返回值:**
  - []string: 转换后的字符串切片
  - error: 如果转换失败返回错误

---

### ValidateFlagName()

```go
func ValidateFlagName(cmd types.Command, longName, shortName string) error
```

#### ValidateFlagName 验证标志名称是否可用

**参数:**
  - cmd: 命令对象, 用于检查标志注册表
  - longName: 长标志名称
  - shortName: 短标志名称

**返回值:**
  - error: 如果验证失败返回错误, 成功返回 nil