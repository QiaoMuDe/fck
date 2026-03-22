# Date子命令实现方案

## 一、功能概述
date子命令用于获取和格式化时间，支持多种时间格式输出、时区转换、时间戳转换等功能。

## 二、核心功能设计

### 1. 支持的时间格式
- **默认格式**：`2006-01-02 15:04:05`（标准格式）
- **ISO 8601格式**：`2006-01-02T15:04:05Z07:00`
- **Unix时间戳**：秒级和毫秒级
- **RFC 3339格式**：`2006-01-02T15:04:05Z07:00`
- **自定义格式**：支持Go的time格式化模板

### 2. 支持的选项
- `--format, -f`：指定时间格式字符串
- `--timestamp, -t`：将Unix时间戳转换为可读时间
- `--timezone, -z`：指定时区（默认为本地时区）
- `--utc, -u`：使用UTC时间
- `--unix, -U`：输出Unix时间戳

### 3. 命令使用示例
```
fck date                          # 默认格式输出当前时间
fck date -f "2006-01-02"         # 自定义格式
fck date -t 1234567890           # 时间戳转换
fck date -u                       # UTC时间
fck date -z "Asia/Shanghai"        # 指定时区
fck date -U                       # 输出Unix时间戳
```

## 三、文件结构设计

```
internal/
├── cli/
│   └── date.go                  # CLI层：命令定义和参数解析
└── commands/
    └── date/
        ├── cmd_date.go           # 业务逻辑层：核心功能实现
        └── formatter.go         # 格式化输出（可选）
```

## 四、数据结构设计

### CLI层（internal/cli/date.go）
```go
var DateCmd *qflag.Cmd

var (
    dateFormat    *qflag.StringFlag  // 时间格式字符串
    dateTimestamp *qflag.StringFlag  // Unix时间戳
    dateTimezone  *qflag.StringFlag  // 时区
    dateUTC       *qflag.BoolFlag    // 使用UTC
    dateUnix      *qflag.BoolFlag    // 输出Unix时间戳
)

func runDate(cmd qflag.Command) error {
    config := date.DateConfig{
        Format:    dateFormat.Get(),
        Timestamp: dateTimestamp.Get(),
        Timezone:  dateTimezone.Get(),
        UTC:       dateUTC.Get(),
        Unix:      dateUnix.Get(),
    }

    return date.DateCmdMain(config)
}
```

### 业务逻辑层（internal/commands/date/cmd_date.go）
```go
type DateConfig struct {
    Format    string  // 时间格式字符串
    Timestamp string  // Unix时间戳
    Timezone  string  // 时区
    UTC       bool    // 使用UTC
    Unix      bool    // 输出Unix时间戳
}

func DateCmdMain(config DateConfig) error {
    // 1. 解析时间戳（如果提供）
    // 2. 设置时区
    // 3. 格式化输出
}
```

## 五、实现细节

### 1. 时间解析逻辑
- 如果提供了`--timestamp`，解析Unix时间戳（支持秒和毫秒）
- 否则使用当前时间

### 2. 时区处理
- 如果指定了`--timezone`，使用指定时区
- 如果指定了`--utc`，使用UTC时区
- 否则使用本地时区

### 3. 格式化输出
- 如果指定了`--format`，使用自定义格式
- 如果指定了`--unix`，输出Unix时间戳
- 否则使用默认格式

### 4. 预定义格式别名
可以支持一些常用格式的快捷方式：
- `iso`：ISO 8601格式
- `rfc`：RFC 3339格式
- `date`：仅日期 `2006-01-02`
- `time`：仅时间 `15:04:05`

## 六、错误处理
- 无效的时间戳格式
- 无效的时区
- 无效的时间格式字符串

## 七、测试用例设计
1. 默认格式输出测试
2. 自定义格式测试
3. 时间戳转换测试
4. 时区转换测试
5. UTC时间测试
6. Unix时间戳输出测试
7. 错误输入处理测试

## 八、边缘案例考虑
1. 空时间戳字符串
2. 无效的时区名称
3. 格式字符串为空
4. 负数时间戳
5. 超大时间戳

## 九、代码改动范围
- 新增文件：`internal/cli/date.go`
- 新增文件：`internal/commands/date/cmd_date.go`
- 修改文件：`internal/cli/root.go`（注册date子命令）

## 十、预估改动量
- CLI层：约80-100行
- 业务逻辑层：约150-200行
- 总计：约230-300行
