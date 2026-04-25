# Route Test 命令设计方案

## 命令概述

用于测试 HTTP/HTTPS 路由连通性和响应的实用工具，支持批量测试、并发请求、多种输出格式。

## 命令结构

```
fck routetest [选项] <URL...>
```

## 功能特性

### 核心功能
1. **单路由测试** - 测试单个 URL 的连通性
2. **批量路由测试** - 从文件读取 URL 列表批量测试
3. **并发测试** - 支持多并发请求进行压力测试
4. **路由链测试** - 跟随重定向，记录完整路由链
5. **定时测试** - 间隔时间重复测试

### 测试指标
- HTTP 状态码
- 响应时间 (DNS/连接/TLS/首字节/总时间)
- 响应体大小
- 重定向次数
- 错误信息

## CLI 定义

### 文件位置
`internal/cli/routetest.go`

### 标志定义

| 长名称 | 短名称 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| file | f | string | "" | 从文件读取 URL 列表 |
| method | X | string | GET | HTTP 方法 |
| header | H | []string | [] | 请求头 |
| data | d | string | "" | 请求体数据 |
| timeout | t | duration | 30s | 请求超时 |
| connect-timeout | | duration | 10s | 连接超时 |
| concurrent | c | int | 1 | 并发数 |
| interval | i | duration | 0 | 请求间隔 |
| repeat | r | int | 1 | 重复次数 |
| follow | L | bool | false | 跟随重定向 |
| max-redirs | | int | 10 | 最大重定向次数 |
| insecure | k | bool | false | 跳过 SSL 验证 |
| output | o | string | "" | 输出文件 |
| format | | enum | text | 输出格式 (text/json/csv) |
| silent | s | bool | false | 静默模式 |
| verbose | v | bool | false | 详细输出 |
| fail | | bool | false | 错误时非零退出 |

### 代码实现

```go
package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/routetest"
	"gitee.com/MM-Q/qflag"
)

var RouteTestCmd *qflag.Cmd

var (
	rtFile           *qflag.StringFlag
	rtMethod         *qflag.StringFlag
	rtHeader         *qflag.StringSliceFlag
	rtData           *qflag.StringFlag
	rtTimeout        *qflag.DurationFlag
	rtConnectTimeout *qflag.DurationFlag
	rtConcurrent     *qflag.IntFlag
	rtInterval       *qflag.DurationFlag
	rtRepeat         *qflag.IntFlag
	rtFollow         *qflag.BoolFlag
	rtMaxRedirs      *qflag.IntFlag
	rtInsecure       *qflag.BoolFlag
	rtOutput         *qflag.StringFlag
	rtFormat         *qflag.EnumFlag
	rtSilent         *qflag.BoolFlag
	rtVerbose        *qflag.BoolFlag
	rtFail           *qflag.BoolFlag
)

func init() {
	RouteTestCmd = qflag.NewCmd("routetest", "rt", qflag.ExitOnError)

	rtFile = RouteTestCmd.String("file", "f", "从文件读取 URL 列表", "")
	rtMethod = RouteTestCmd.String("method", "X", "HTTP 方法", "GET")
	rtHeader = RouteTestCmd.StringSlice("header", "H", "请求头", []string{})
	rtData = RouteTestCmd.String("data", "d", "请求体数据", "")
	rtTimeout = RouteTestCmd.Duration("timeout", "t", "请求超时", time.Second*30)
	rtConnectTimeout = RouteTestCmd.Duration("connect-timeout", "", "连接超时", time.Second*10)
	rtConcurrent = RouteTestCmd.Int("concurrent", "c", "并发数", 1)
	rtInterval = RouteTestCmd.Duration("interval", "i", "请求间隔", 0)
	rtRepeat = RouteTestCmd.Int("repeat", "r", "重复次数", 1)
	rtFollow = RouteTestCmd.Bool("follow", "L", "跟随重定向", false)
	rtMaxRedirs = RouteTestCmd.Int("max-redirs", "", "最大重定向次数", 10)
	rtInsecure = RouteTestCmd.Bool("insecure", "k", "跳过 SSL 验证", false)
	rtOutput = RouteTestCmd.String("output", "o", "输出文件", "")
	rtFormat = RouteTestCmd.Enum("format", "", "输出格式", "text", []string{"text", "json", "csv"})
	rtSilent = RouteTestCmd.Bool("silent", "s", "静默模式", false)
	rtVerbose = RouteTestCmd.Bool("verbose", "v", "详细输出", false)
	rtFail = RouteTestCmd.Bool("fail", "", "错误时非零退出", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "HTTP/HTTPS 路由测试工具",
		UsageSyntax: fmt.Sprintf("%s routetest [选项] <URL...>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"URL 必须包含协议头 (http:// 或 https://)",
			"使用 -f 从文件读取 URL，每行一个",
			"使用 -c 进行并发压力测试",
			"使用 -L 跟随重定向，查看完整路由链",
			"使用 --format json 输出结构化数据",
		},
		Examples: map[string]string{
			"测试单个路由":         fmt.Sprintf("%s routetest https://api.example.com/health", qflag.Root.Name()),
			"批量测试":           fmt.Sprintf("%s routetest -f urls.txt", qflag.Root.Name()),
			"POST 测试":          fmt.Sprintf("%s routetest -X POST -d '{\"key\":\"value\"}' https://api.example.com/test", qflag.Root.Name()),
			"并发压力测试":        fmt.Sprintf("%s routetest -c 100 -r 10 https://api.example.com/test", qflag.Root.Name()),
			"跟随重定向":         fmt.Sprintf("%s routetest -L https://bit.ly/xxx", qflag.Root.Name()),
			"JSON 输出":          fmt.Sprintf("%s routetest --format json https://api.example.com/test", qflag.Root.Name()),
			"定时重复测试":        fmt.Sprintf("%s routetest -r 5 -i 2s https://api.example.com/test", qflag.Root.Name()),
			"导出 CSV":           fmt.Sprintf("%s routetest -f urls.txt --format csv -o results.csv", qflag.Root.Name()),
		},
	}

	if err := RouteTestCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	RouteTestCmd.SetRun(runRouteTest)
}

func runRouteTest(cmd qflag.Command) error {
	config := routetest.RouteTestConfig{
		URLs:           cmd.Args(),
		File:           rtFile.Get(),
		Method:         rtMethod.Get(),
		Headers:        rtHeader.Get(),
		Data:           rtData.Get(),
		Timeout:        rtTimeout.Get(),
		ConnectTimeout: rtConnectTimeout.Get(),
		Concurrent:     rtConcurrent.Get(),
		Interval:       rtInterval.Get(),
		Repeat:         rtRepeat.Get(),
		Follow:         rtFollow.Get(),
		MaxRedirs:      rtMaxRedirs.Get(),
		Insecure:       rtInsecure.Get(),
		Output:         rtOutput.Get(),
		Format:         rtFormat.Get(),
		Silent:         rtSilent.Get(),
		Verbose:        rtVerbose.Get(),
		Fail:           rtFail.Get(),
	}

	return routetest.RouteTestCmdMain(config)
}
```

## 业务逻辑

### 文件位置
`internal/commands/routetest/cmd_routetest.go`

### 配置结构体

```go
package routetest

import "time"

// RouteTestConfig 路由测试配置
type RouteTestConfig struct {
	URLs           []string
	File           string
	Method         string
	Headers        []string
	Data           string
	Timeout        time.Duration
	ConnectTimeout time.Duration
	Concurrent     int
	Interval       time.Duration
	Repeat         int
	Follow         bool
	MaxRedirs      int
	Insecure       bool
	Output         string
	Format         string
	Silent         bool
	Verbose        bool
	Fail           bool
}

// RouteResult 单个路由测试结果
type RouteResult struct {
	URL           string            `json:"url"`
	Method        string            `json:"method"`
	StatusCode    int               `json:"status_code"`
	Status        string            `json:"status"`
	Headers       map[string]string `json:"headers,omitempty"`
	BodySize      int64             `json:"body_size"`
	TimeDNS       time.Duration     `json:"time_dns_ms"`
	TimeConnect   time.Duration     `json:"time_connect_ms"`
	TimeTLS       time.Duration     `json:"time_tls_ms"`
	TimeFirstByte time.Duration     `json:"time_first_byte_ms"`
	TimeTotal     time.Duration     `json:"time_total_ms"`
	Redirects     int               `json:"redirects"`
	RedirectChain []string          `json:"redirect_chain,omitempty"`
	Error         string            `json:"error,omitempty"`
}

// TestSummary 测试汇总
type TestSummary struct {
	Total      int           `json:"total"`
	Success    int           `json:"success"`
	Failed     int           `json:"failed"`
	TotalTime  time.Duration `json:"total_time_ms"`
	AvgTime    time.Duration `json:"avg_time_ms"`
	MinTime    time.Duration `json:"min_time_ms"`
	MaxTime    time.Duration `json:"max_time_ms"`
	StatusDist map[int]int   `json:"status_distribution"`
}
```

### 主函数

```go
// RouteTestCmdMain 路由测试主函数
func RouteTestCmdMain(config RouteTestConfig) error {
	// 获取所有要测试的 URL
	urls, err := getURLs(config)
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		return fmt.Errorf("未指定测试 URL")
	}

	// 执行测试
	results := executeTests(urls, config)

	// 输出结果
	return outputResults(results, config)
}
```

## 核心功能实现

### 1. URL 获取
- 从命令行参数获取
- 从文件读取（每行一个 URL）
- 去重处理

### 2. 测试执行
- 单 URL 测试：顺序执行
- 批量测试：支持并发
- 重复测试：按间隔重复

### 3. 结果输出
- **text**: 表格格式，适合终端查看
- **json**: 结构化数据，适合程序处理
- **csv**: 表格数据，适合 Excel 分析

## 使用示例

```bash
# 基础测试
fck routetest https://api.example.com/health

# 批量测试
fck routetest -f urls.txt

# 并发压力测试
fck routetest -c 100 -r 10 https://api.example.com/test

# 完整路由链测试
fck routetest -L -v https://bit.ly/xxx

# JSON 输出
fck routetest --format json https://api.example.com/test | jq

# 导出结果到 CSV
fck routetest -f urls.txt --format csv -o results.csv
```

## 注册到根命令

在 `internal/cli/root.go` 的 `SubCmds` 中添加：

```go
SubCmds: []qflag.Command{
	// ... 其他命令
	RouteTestCmd,
},
```

## 实现步骤

1. 创建业务逻辑文件 `internal/commands/routetest/cmd_routetest.go`
2. 创建 CLI 定义文件 `internal/cli/routetest.go`
3. 在 `internal/cli/root.go` 中注册命令
4. 测试编译和功能验证
