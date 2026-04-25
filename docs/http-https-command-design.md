# HTTP/HTTPS 子命令设计方案

## 1. 概述

设计两个独立的子命令 `http` 和 `https`，用于执行 HTTP/HTTPS 请求测试、API 调试和 Web 服务探测。

## 2. 命令定位

| 命令 | 用途 | 默认端口 |
|------|------|----------|
| `http` | 执行 HTTP 请求 | 80 |
| `https` | 执行 HTTPS 请求 | 443 |

## 3. 目录结构

```
internal/
├── cli/
│   ├── http.go          # HTTP 命令 CLI 定义
│   ├── https.go         # HTTPS 命令 CLI 定义
│   └── root.go          # 注册 HttpCmd 和 HttpsCmd
├── commands/
│   ├── http/
│   │   ├── cmd_http.go      # HTTP 业务逻辑
│   │   ├── client.go        # HTTP 客户端封装
│   │   ├── request.go       # 请求构建
│   │   ├── response.go      # 响应处理
│   │   └── output.go        # 输出格式化
│   └── https/
│       ├── cmd_https.go     # HTTPS 业务逻辑
│       └── ...              # 可复用 http 包的代码
```

## 4. 功能特性

### 4.1 支持的 HTTP 方法
- GET (默认)
- POST
- PUT
- DELETE
- PATCH
- HEAD
- OPTIONS

### 4.2 请求功能
- 自定义请求头 (-H, --header)
- 请求体数据 (-d, --data)
- 发送文件作为请求体 (-f, --file)
- URL 查询参数 (-q, --query)
- 表单数据提交 (-F, --form)
- JSON 自动 Content-Type (-j, --json)
- Cookie 设置 (-b, --cookie)
- 用户代理设置 (-A, --user-agent)
- 基础认证 (-u, --user)
- Bearer Token 认证 (--token)

### 4.3 连接与性能
- 连接超时控制 (-t, --timeout)
- 请求重试次数 (--retry)
- 重试间隔 (--retry-delay)
- 并发请求 (-c, --concurrent)
- 请求间隔 (-i, --interval)
- 跟随重定向 (-L, --location)
- 最大重定向次数 (--max-redirs)
- 禁用 SSL 证书验证 (-k, --insecure) [仅 https]

### 4.4 输出控制
- 静默模式 (-s, --silent)
- JSON 格式输出 (-j, --json) [注意与请求 JSON 冲突]
- 仅显示响应头 (-I, --head)
- 显示详细过程 (-v, --verbose)
- 保存响应到文件 (-o, --output)
- 显示响应时间 (--show-time)
- 显示状态码 (--show-code)
- 十六进制显示响应体 (--hex)
- 格式化 JSON 响应 (--pretty)

### 4.5 高级功能
- 代理支持 (-x, --proxy)
- 自定义 DNS 解析 (--resolve)
- HTTP 版本指定 (--http1.1, --http2)
- WebSocket 升级 (--websocket)
- 下载文件并显示进度 (--download)
- 断点续传 (--continue)

## 5. 命令行接口设计

### 5.1 HTTP 命令

```go
// internal/cli/http.go
package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/http"
	"gitee.com/MM-Q/qflag"
)

var HttpCmd *qflag.Cmd

var (
	httpMethod      *qflag.StringFlag   // -X, --request    HTTP 方法
	httpHeader      *qflag.StringSliceFlag // -H, --header     请求头
	httpData        *qflag.StringFlag   // -d, --data       请求体数据
	httpFile        *qflag.StringFlag   // -f, --file       从文件读取请求体
	httpQuery       *qflag.StringSliceFlag // -q, --query      URL 查询参数
	httpForm        *qflag.StringSliceFlag // -F, --form       表单数据
	httpJSON        *qflag.BoolFlag     // -j, --json       发送 JSON 数据
	httpCookie      *qflag.StringFlag   // -b, --cookie     Cookie
	httpUserAgent   *qflag.StringFlag   // -A, --user-agent 用户代理
	httpAuth        *qflag.StringFlag   // -u, --user       基础认证 user:password
	httpToken       *qflag.StringFlag   // --token          Bearer Token
	httpTimeout     *qflag.DurationFlag // -t, --timeout    超时时间
	httpRetry       *qflag.IntFlag      // --retry          重试次数
	httpRetryDelay  *qflag.DurationFlag // --retry-delay    重试间隔
	httpConcurrent  *qflag.IntFlag      // -c, --concurrent 并发数
	httpInterval    *qflag.DurationFlag // -i, --interval   请求间隔
	httpLocation    *qflag.BoolFlag     // -L, --location   跟随重定向
	httpMaxRedirs   *qflag.IntFlag      // --max-redirs     最大重定向次数
	httpSilent      *qflag.BoolFlag     // -s, --silent     静默模式
	httpOutputJSON  *qflag.BoolFlag     // --output-json    JSON 格式输出
	httpHeadOnly    *qflag.BoolFlag     // -I, --head       仅显示响应头
	httpVerbose     *qflag.BoolFlag     // -v, --verbose    详细输出
	httpOutput      *qflag.StringFlag   // -o, --output     保存响应到文件
	httpShowTime    *qflag.BoolFlag     // --show-time      显示响应时间
	httpShowCode    *qflag.BoolFlag     // --show-code      仅显示状态码
	httpHex         *qflag.BoolFlag     // --hex            十六进制显示
	httpPretty      *qflag.BoolFlag     // --pretty         格式化 JSON
	httpProxy       *qflag.StringFlag   // -x, --proxy      代理地址
	httpResolve     *qflag.StringSliceFlag // --resolve        自定义 DNS
	httpHTTP2       *qflag.BoolFlag     // --http2          强制使用 HTTP/2
	httpDownload    *qflag.BoolFlag     // --download       下载模式
	httpContinue    *qflag.BoolFlag     // --continue       断点续传
)

func init() {
	HttpCmd = qflag.NewCmd("http", "", qflag.ExitOnError)

	httpMethod = HttpCmd.String("request", "X", "HTTP 方法 (GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS)", "GET")
	httpHeader = HttpCmd.StringSlice("header", "H", "请求头, 可多次使用, 如: -H 'Content-Type: application/json'", []string{})
	httpData = HttpCmd.String("data", "d", "请求体数据", "")
	httpFile = HttpCmd.String("file", "f", "从文件读取请求体", "")
	httpQuery = HttpCmd.StringSlice("query", "q", "URL 查询参数, 如: -q 'key=value'", []string{})
	httpForm = HttpCmd.StringSlice("form", "F", "表单数据, 如: -F 'name=value'", []string{})
	httpJSON = HttpCmd.Bool("json", "j", "发送 JSON 数据(自动设置 Content-Type)", false)
	httpCookie = HttpCmd.String("cookie", "b", "Cookie 字符串", "")
	httpUserAgent = HttpCmd.String("user-agent", "A", "User-Agent", "fck-http/1.0")
	httpAuth = HttpCmd.String("user", "u", "基础认证, 格式: username:password", "")
	httpToken = HttpCmd.String("token", "", "Bearer Token 认证", "")
	httpTimeout = HttpCmd.Duration("timeout", "t", "请求超时时间", time.Second*30)
	httpRetry = HttpCmd.Int("retry", "", "请求失败重试次数", 0)
	httpRetryDelay = HttpCmd.Duration("retry-delay", "", "重试间隔", time.Second)
	httpConcurrent = HttpCmd.Int("concurrent", "c", "并发请求数", 1)
	httpInterval = HttpCmd.Duration("interval", "i", "并发请求间隔", 0)
	httpLocation = HttpCmd.Bool("location", "L", "跟随重定向", false)
	httpMaxRedirs = HttpCmd.Int("max-redirs", "", "最大重定向次数", 10)
	httpSilent = HttpCmd.Bool("silent", "s", "静默模式, 不显示进度", false)
	httpOutputJSON = HttpCmd.Bool("output-json", "", "JSON 格式输出结果", false)
	httpHeadOnly = HttpCmd.Bool("head", "I", "仅发送 HEAD 请求获取响应头", false)
	httpVerbose = HttpCmd.Bool("verbose", "v", "显示详细请求/响应信息", false)
	httpOutput = HttpCmd.String("output", "o", "保存响应体到文件", "")
	httpShowTime = HttpCmd.Bool("show-time", "", "显示请求耗时", false)
	httpShowCode = HttpCmd.Bool("show-code", "", "仅显示 HTTP 状态码", false)
	httpHex = HttpCmd.Bool("hex", "", "以十六进制格式显示响应体", false)
	httpPretty = HttpCmd.Bool("pretty", "", "格式化 JSON 响应", false)
	httpProxy = HttpCmd.String("proxy", "x", "代理地址, 如: http://proxy:8080", "")
	httpResolve = HttpCmd.StringSlice("resolve", "", "自定义 DNS 解析, 如: --resolve 'host:port:ip'", []string{})
	httpHTTP2 = HttpCmd.Bool("http2", "", "强制使用 HTTP/2", false)
	httpDownload = HttpCmd.Bool("download", "", "下载模式, 显示进度条", false)
	httpContinue = HttpCmd.Bool("continue", "", "断点续传(需配合 -o 使用)", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "HTTP 请求测试工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s http [options] <URL>", qflag.Root.Name()),
		Notes: []string{
			"URL 必须包含协议头 (http://)",
			"使用 -X 指定 HTTP 方法, 默认为 GET",
			"使用 -H 可多次添加请求头",
			"使用 -j 自动设置 Content-Type: application/json",
			"使用 -I 等价于 -X HEAD",
		},
		Examples: map[string]string{
			"简单 GET 请求":           fmt.Sprintf("%s http http://api.example.com/users", qflag.Root.Name()),
			"POST JSON 数据":          fmt.Sprintf("%s http -X POST -j -d '{\"name\":\"test\"}' http://api.example.com/users", qflag.Root.Name()),
			"添加自定义请求头":         fmt.Sprintf("%s http -H 'Authorization: Bearer token' http://api.example.com/data", qflag.Root.Name()),
			"发送表单数据":            fmt.Sprintf("%s http -X POST -F 'name=value' -F 'file=@upload.txt' http://api.example.com/upload", qflag.Root.Name()),
			"下载文件":               fmt.Sprintf("%s http --download -o file.zip http://example.com/file.zip", qflag.Root.Name()),
			"仅显示状态码":           fmt.Sprintf("%s http --show-code http://example.com", qflag.Root.Name()),
			"带重试的请求":           fmt.Sprintf("%s http --retry 3 --retry-delay 2s http://api.example.com/data", qflag.Root.Name()),
			"并发压力测试":           fmt.Sprintf("%s http -c 100 -i 10ms http://api.example.com/test", qflag.Root.Name()),
			"使用代理":               fmt.Sprintf("%s http -x http://proxy:8080 http://example.com", qflag.Root.Name()),
			"基础认证":               fmt.Sprintf("%s http -u admin:password http://api.example.com/secure", qflag.Root.Name()),
			"Bearer Token":           fmt.Sprintf("%s http --token 'xyz123' http://api.example.com/data", qflag.Root.Name()),
			"JSON 格式输出":          fmt.Sprintf("%s http --output-json http://api.example.com/users", qflag.Root.Name()),
			"格式化 JSON 响应":        fmt.Sprintf("%s http --pretty http://api.example.com/data", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "data-source",
				Flags:     []string{"data", "file"},
				AllowNone: true,
			},
			{
				Name:      "output-format",
				Flags:     []string{"show-code", "head", "hex", "silent"},
				AllowNone: true,
			},
		},
	}

	if err := HttpCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	HttpCmd.SetRun(runHttp)
}

func runHttp(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("缺少 URL 参数")
	}

	url := args[0]

	config := http.HttpConfig{
		URL:           url,
		Method:        httpMethod.Get(),
		Headers:       httpHeader.Get(),
		Data:          httpData.Get(),
		File:          httpFile.Get(),
		Query:         httpQuery.Get(),
		Form:          httpForm.Get(),
		JSON:          httpJSON.Get(),
		Cookie:        httpCookie.Get(),
		UserAgent:     httpUserAgent.Get(),
		Auth:          httpAuth.Get(),
		Token:         httpToken.Get(),
		Timeout:       httpTimeout.Get(),
		Retry:         httpRetry.Get(),
		RetryDelay:    httpRetryDelay.Get(),
		Concurrent:    httpConcurrent.Get(),
		Interval:      httpInterval.Get(),
		FollowRedirect:httpLocation.Get(),
		MaxRedirects:  httpMaxRedirs.Get(),
		Silent:        httpSilent.Get(),
		OutputJSON:    httpOutputJSON.Get(),
		HeadOnly:      httpHeadOnly.Get(),
		Verbose:       httpVerbose.Get(),
		Output:        httpOutput.Get(),
		ShowTime:      httpShowTime.Get(),
		ShowCode:      httpShowCode.Get(),
		Hex:           httpHex.Get(),
		Pretty:        httpPretty.Get(),
		Proxy:         httpProxy.Get(),
		Resolve:       httpResolve.Get(),
		HTTP2:         httpHTTP2.Get(),
		Download:      httpDownload.Get(),
		Continue:      httpContinue.Get(),
	}

	return http.HttpCmdMain(config)
}
```

### 5.2 HTTPS 命令

```go
// internal/cli/https.go
package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/https"
	"gitee.com/MM-Q/qflag"
)

var HttpsCmd *qflag.Cmd

var (
	httpsMethod      *qflag.StringFlag   // -X, --request    HTTP 方法
	httpsHeader      *qflag.StringSliceFlag // -H, --header     请求头
	httpsData        *qflag.StringFlag   // -d, --data       请求体数据
	httpsFile        *qflag.StringFlag   // -f, --file       从文件读取请求体
	httpsQuery       *qflag.StringSliceFlag // -q, --query      URL 查询参数
	httpsForm        *qflag.StringSliceFlag // -F, --form       表单数据
	httpsJSON        *qflag.BoolFlag     // -j, --json       发送 JSON 数据
	httpsCookie      *qflag.StringFlag   // -b, --cookie     Cookie
	httpsUserAgent   *qflag.StringFlag   // -A, --user-agent 用户代理
	httpsAuth        *qflag.StringFlag   // -u, --user       基础认证
	httpsToken       *qflag.StringFlag   // --token          Bearer Token
	httpsTimeout     *qflag.DurationFlag // -t, --timeout    超时时间
	httpsRetry       *qflag.IntFlag      // --retry          重试次数
	httpsRetryDelay  *qflag.DurationFlag // --retry-delay    重试间隔
	httpsConcurrent  *qflag.IntFlag      // -c, --concurrent 并发数
	httpsInterval    *qflag.DurationFlag // -i, --interval   请求间隔
	httpsLocation    *qflag.BoolFlag     // -L, --location   跟随重定向
	httpsMaxRedirs   *qflag.IntFlag      // --max-redirs     最大重定向次数
	httpsInsecure    *qflag.BoolFlag     // -k, --insecure   跳过 SSL 验证
	httpsCert        *qflag.StringFlag   // --cert           客户端证书
	httpsKey         *qflag.StringFlag   // --key            客户端私钥
	httpsCA          *qflag.StringFlag   // --cacert         CA 证书
	httpsTLSVersion  *qflag.StringFlag   // --tls-version    TLS 版本
	httpsSilent      *qflag.BoolFlag     // -s, --silent     静默模式
	httpsOutputJSON  *qflag.BoolFlag     // --output-json    JSON 格式输出
	httpsHeadOnly    *qflag.BoolFlag     // -I, --head       仅显示响应头
	httpsVerbose     *qflag.BoolFlag     // -v, --verbose    详细输出
	httpsOutput      *qflag.StringFlag   // -o, --output     保存响应到文件
	httpsShowTime    *qflag.BoolFlag     // --show-time      显示响应时间
	httpsShowCode    *qflag.BoolFlag     // --show-code      仅显示状态码
	httpsHex         *qflag.BoolFlag     // --hex            十六进制显示
	httpsPretty      *qflag.BoolFlag     // --pretty         格式化 JSON
	httpsProxy       *qflag.StringFlag   // -x, --proxy      代理地址
	httpsResolve     *qflag.StringSliceFlag // --resolve        自定义 DNS
	httpsHTTP2       *qflag.BoolFlag     // --http2          强制使用 HTTP/2
	httpsHTTP3       *qflag.BoolFlag     // --http3          尝试使用 HTTP/3
	httpsDownload    *qflag.BoolFlag     // --download       下载模式
	httpsContinue    *qflag.BoolFlag     // --continue       断点续传
)

func init() {
	HttpsCmd = qflag.NewCmd("https", "", qflag.ExitOnError)

	httpsMethod = HttpsCmd.String("request", "X", "HTTP 方法 (GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS)", "GET")
	httpsHeader = HttpsCmd.StringSlice("header", "H", "请求头, 可多次使用", []string{})
	httpsData = HttpsCmd.String("data", "d", "请求体数据", "")
	httpsFile = HttpsCmd.String("file", "f", "从文件读取请求体", "")
	httpsQuery = HttpsCmd.StringSlice("query", "q", "URL 查询参数", []string{})
	httpsForm = HttpsCmd.StringSlice("form", "F", "表单数据", []string{})
	httpsJSON = HttpsCmd.Bool("json", "j", "发送 JSON 数据", false)
	httpsCookie = HttpsCmd.String("cookie", "b", "Cookie 字符串", "")
	httpsUserAgent = HttpsCmd.String("user-agent", "A", "User-Agent", "fck-https/1.0")
	httpsAuth = HttpsCmd.String("user", "u", "基础认证, 格式: username:password", "")
	httpsToken = HttpsCmd.String("token", "", "Bearer Token 认证", "")
	httpsTimeout = HttpsCmd.Duration("timeout", "t", "请求超时时间", time.Second*30)
	httpsRetry = HttpsCmd.Int("retry", "", "请求失败重试次数", 0)
	httpsRetryDelay = HttpsCmd.Duration("retry-delay", "", "重试间隔", time.Second)
	httpsConcurrent = HttpsCmd.Int("concurrent", "c", "并发请求数", 1)
	httpsInterval = HttpsCmd.Duration("interval", "i", "并发请求间隔", 0)
	httpsLocation = HttpsCmd.Bool("location", "L", "跟随重定向", false)
	httpsMaxRedirs = HttpsCmd.Int("max-redirs", "", "最大重定向次数", 10)
	httpsInsecure = HttpsCmd.Bool("insecure", "k", "跳过 SSL 证书验证", false)
	httpsCert = HttpsCmd.String("cert", "", "客户端证书文件路径", "")
	httpsKey = HttpsCmd.String("key", "", "客户端私钥文件路径", "")
	httpsCA = HttpsCmd.String("cacert", "", "CA 证书文件路径", "")
	httpsTLSVersion = HttpsCmd.String("tls-version", "", "TLS 版本 (1.0/1.1/1.2/1.3)", "")
	httpsSilent = HttpsCmd.Bool("silent", "s", "静默模式", false)
	httpsOutputJSON = HttpsCmd.Bool("output-json", "", "JSON 格式输出结果", false)
	httpsHeadOnly = HttpsCmd.Bool("head", "I", "仅发送 HEAD 请求", false)
	httpsVerbose = HttpsCmd.Bool("verbose", "v", "显示详细信息", false)
	httpsOutput = HttpsCmd.String("output", "o", "保存响应体到文件", "")
	httpsShowTime = HttpsCmd.Bool("show-time", "", "显示请求耗时", false)
	httpsShowCode = HttpsCmd.Bool("show-code", "", "仅显示 HTTP 状态码", false)
	httpsHex = HttpsCmd.Bool("hex", "", "以十六进制格式显示响应体", false)
	httpsPretty = HttpsCmd.Bool("pretty", "", "格式化 JSON 响应", false)
	httpsProxy = HttpsCmd.String("proxy", "x", "代理地址", "")
	httpsResolve = HttpsCmd.StringSlice("resolve", "", "自定义 DNS 解析", []string{})
	httpsHTTP2 = HttpsCmd.Bool("http2", "", "强制使用 HTTP/2", false)
	httpsHTTP3 = HttpsCmd.Bool("http3", "", "尝试使用 HTTP/3", false)
	httpsDownload = HttpsCmd.Bool("download", "", "下载模式", false)
	httpsContinue = HttpsCmd.Bool("continue", "", "断点续传", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "HTTPS 请求测试工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s https [options] <URL>", qflag.Root.Name()),
		Notes: []string{
			"URL 必须包含协议头 (https://)",
			"使用 -k 跳过 SSL 证书验证(不推荐用于生产环境)",
			"使用 --cert/--key 指定客户端证书",
			"使用 --cacert 指定自定义 CA 证书",
		},
		Examples: map[string]string{
			"简单 GET 请求":           fmt.Sprintf("%s https https://api.example.com/users", qflag.Root.Name()),
			"跳过 SSL 验证":           fmt.Sprintf("%s https -k https://self-signed.example.com", qflag.Root.Name()),
			"指定客户端证书":          fmt.Sprintf("%s https --cert client.crt --key client.key https://api.example.com", qflag.Root.Name()),
			"POST JSON 数据":          fmt.Sprintf("%s https -X POST -j -d '{\"name\":\"test\"}' https://api.example.com/users", qflag.Root.Name()),
			"下载文件":               fmt.Sprintf("%s https --download -o file.zip https://example.com/file.zip", qflag.Root.Name()),
			"仅显示状态码":           fmt.Sprintf("%s https --show-code https://example.com", qflag.Root.Name()),
			"使用 HTTP/2":            fmt.Sprintf("%s https --http2 https://api.example.com", qflag.Root.Name()),
			"测试 TLS 版本":          fmt.Sprintf("%s https --tls-version 1.3 https://example.com", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "data-source",
				Flags:     []string{"data", "file"},
				AllowNone: true,
			},
			{
				Name:      "output-format",
				Flags:     []string{"show-code", "head", "hex", "silent"},
				AllowNone: true,
			},
		},
	}

	if err := HttpsCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	HttpsCmd.SetRun(runHttps)
}

func runHttps(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("缺少 URL 参数")
	}

	url := args[0]

	config := https.HttpsConfig{
		URL:           url,
		Method:        httpsMethod.Get(),
		Headers:       httpsHeader.Get(),
		Data:          httpsData.Get(),
		File:          httpsFile.Get(),
		Query:         httpsQuery.Get(),
		Form:          httpsForm.Get(),
		JSON:          httpsJSON.Get(),
		Cookie:        httpsCookie.Get(),
		UserAgent:     httpsUserAgent.Get(),
		Auth:          httpsAuth.Get(),
		Token:         httpsToken.Get(),
		Timeout:       httpsTimeout.Get(),
		Retry:         httpsRetry.Get(),
		RetryDelay:    httpsRetryDelay.Get(),
		Concurrent:    httpsConcurrent.Get(),
		Interval:      httpsInterval.Get(),
		FollowRedirect:httpsLocation.Get(),
		MaxRedirects:  httpsMaxRedirs.Get(),
		Insecure:      httpsInsecure.Get(),
		Cert:          httpsCert.Get(),
		Key:           httpsKey.Get(),
		CA:            httpsCA.Get(),
		TLSVersion:    httpsTLSVersion.Get(),
		Silent:        httpsSilent.Get(),
		OutputJSON:    httpsOutputJSON.Get(),
		HeadOnly:      httpsHeadOnly.Get(),
		Verbose:       httpsVerbose.Get(),
		Output:        httpsOutput.Get(),
		ShowTime:      httpsShowTime.Get(),
		ShowCode:      httpsShowCode.Get(),
		Hex:           httpsHex.Get(),
		Pretty:        httpsPretty.Get(),
		Proxy:         httpsProxy.Get(),
		Resolve:       httpsResolve.Get(),
		HTTP2:         httpsHTTP2.Get(),
		HTTP3:         httpsHTTP3.Get(),
		Download:      httpsDownload.Get(),
		Continue:      httpsContinue.Get(),
	}

	return https.HttpsCmdMain(config)
}
```

## 6. 配置结构体设计

### 6.1 HTTP 配置

```go
// internal/commands/http/cmd_http.go
package http

import "time"

// HttpConfig HTTP 请求配置
type HttpConfig struct {
	// 请求配置
	URL       string   // 目标 URL
	Method    string   // HTTP 方法
	Headers   []string // 请求头列表 (格式: "Key: Value")
	Data      string   // 请求体数据
	File      string   // 请求体文件路径
	Query     []string // URL 查询参数 (格式: "key=value")
	Form      []string // 表单数据 (格式: "key=value" 或 "key=@file")
	JSON      bool     // 发送 JSON 数据
	Cookie    string   // Cookie 字符串
	UserAgent string   // User-Agent
	Auth      string   // 基础认证 (格式: "username:password")
	Token     string   // Bearer Token

	// 连接配置
	Timeout        time.Duration // 请求超时
	Retry          int           // 重试次数
	RetryDelay     time.Duration // 重试间隔
	Concurrent     int           // 并发数
	Interval       time.Duration // 并发间隔
	FollowRedirect bool          // 跟随重定向
	MaxRedirects   int           // 最大重定向次数
	Proxy          string        // 代理地址
	Resolve        []string      // 自定义 DNS (格式: "host:port:ip")
	HTTP2          bool          // 强制 HTTP/2

	// 输出配置
	Silent     bool   // 静默模式
	OutputJSON bool   // JSON 格式输出
	HeadOnly   bool   // 仅 HEAD 请求
	Verbose    bool   // 详细输出
	Output     string // 输出文件路径
	ShowTime   bool   // 显示耗时
	ShowCode   bool   // 仅显示状态码
	Hex        bool   // 十六进制显示
	Pretty     bool   // 格式化 JSON
	Download   bool   // 下载模式
	Continue   bool   // 断点续传
}

// HttpResult HTTP 请求结果
type HttpResult struct {
	URL        string        `json:"url"`
	Method     string        `json:"method"`
	StatusCode int           `json:"status_code"`
	Status     string        `json:"status"`
	Headers    http.Header   `json:"headers"`
	Body       []byte        `json:"body,omitempty"`
	BodySize   int64         `json:"body_size"`
	TimeCost   time.Duration `json:"time_cost_ms"`
	Error      string        `json:"error,omitempty"`
}

// HttpStats 并发请求统计
type HttpStats struct {
	Total      int           // 总请求数
	Success    int           // 成功数
	Failed     int           // 失败数
	TotalTime  time.Duration // 总耗时
	AvgTime    time.Duration // 平均耗时
	MinTime    time.Duration // 最小耗时
	MaxTime    time.Duration // 最大耗时
	StatusDist map[int]int   // 状态码分布
}

// HttpCmdMain HTTP 命令主函数
func HttpCmdMain(config HttpConfig) error {
	// 实现逻辑
	return nil
}
```

### 6.2 HTTPS 配置

```go
// internal/commands/https/cmd_https.go
package https

import (
	"time"
	"gitee.com/MM-Q/fck/internal/commands/http"
)

// HttpsConfig HTTPS 请求配置
type HttpsConfig struct {
	// 嵌入 HTTP 通用配置 (使用组合而非继承)
	URL       string
	Method    string
	Headers   []string
	Data      string
	File      string
	Query     []string
	Form      []string
	JSON      bool
	Cookie    string
	UserAgent string
	Auth      string
	Token     string

	// 连接配置
	Timeout        time.Duration
	Retry          int
	RetryDelay     time.Duration
	Concurrent     int
	Interval       time.Duration
	FollowRedirect bool
	MaxRedirects   int
	Proxy          string
	Resolve        []string
	HTTP2          bool
	HTTP3          bool // HTTPS 特有

	// SSL/TLS 配置 (HTTPS 特有)
	Insecure   bool   // 跳过证书验证
	Cert       string // 客户端证书
	Key        string // 客户端私钥
	CA         string // CA 证书
	TLSVersion string // TLS 版本

	// 输出配置
	Silent     bool
	OutputJSON bool
	HeadOnly   bool
	Verbose    bool
	Output     string
	ShowTime   bool
	ShowCode   bool
	Hex        bool
	Pretty     bool
	Download   bool
	Continue   bool
}

// HttpsResult HTTPS 请求结果
type HttpsResult struct {
	http.HttpResult        // 嵌入 HTTP 结果
	TLSVersion      string `json:"tls_version"`
	CipherSuite     string `json:"cipher_suite"`
	ServerCert      *CertInfo `json:"server_cert,omitempty"`
}

// CertInfo 证书信息
type CertInfo struct {
	Subject    string    `json:"subject"`
	Issuer     string    `json:"issuer"`
	NotBefore  time.Time `json:"not_before"`
	NotAfter   time.Time `json:"not_after"`
	DNSNames   []string  `json:"dns_names"`
	IPAddresses []string `json:"ip_addresses"`
}

// HttpsCmdMain HTTPS 命令主函数
func HttpsCmdMain(config HttpsConfig) error {
	// 实现逻辑
	return nil
}
```

## 7. 核心功能实现思路

### 7.1 HTTP 客户端封装

```go
// internal/commands/http/client.go
package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Client HTTP 客户端封装
type Client struct {
	httpClient *http.Client
	config     HttpConfig
}

// NewClient 创建 HTTP 客户端
func NewClient(config HttpConfig) (*Client, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.Timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2: config.HTTP2,
	}

	// 配置代理
	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// 配置自定义 DNS
	if len(config.Resolve) > 0 {
		// 实现自定义解析器
	}

	// 配置重定向策略
	checkRedirect := func(req *http.Request, via []*http.Request) error {
		if !config.FollowRedirect {
			return http.ErrUseLastResponse
		}
		if len(via) >= config.MaxRedirects {
			return fmt.Errorf("stopped after %d redirects", config.MaxRedirects)
		}
		return nil
	}

	return &Client{
		httpClient: &http.Client{
			Transport:     transport,
			Timeout:       config.Timeout,
			CheckRedirect: checkRedirect,
		},
		config: config,
	}, nil
}

// Do 执行 HTTP 请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
```

### 7.2 请求构建

```go
// internal/commands/http/request.go
package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// BuildRequest 构建 HTTP 请求
func BuildRequest(config HttpConfig) (*http.Request, error) {
	// 解析 URL
	u, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	// 添加查询参数
	if len(config.Query) > 0 {
		q := u.Query()
		for _, param := range config.Query {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) == 2 {
				q.Add(parts[0], parts[1])
			}
		}
		u.RawQuery = q.Encode()
	}

	// 确定请求体
	var body io.Reader
	if config.File != "" {
		// 从文件读取
		data, err := os.ReadFile(config.File)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	} else if config.Data != "" {
		body = strings.NewReader(config.Data)
	}

	// 创建请求
	method := config.Method
	if config.HeadOnly {
		method = "HEAD"
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for _, header := range config.Headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// 自动设置 JSON Content-Type
	if config.JSON {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置 User-Agent
	if config.UserAgent != "" {
		req.Header.Set("User-Agent", config.UserAgent)
	}

	// 设置 Cookie
	if config.Cookie != "" {
		req.Header.Set("Cookie", config.Cookie)
	}

	// 设置基础认证
	if config.Auth != "" {
		parts := strings.SplitN(config.Auth, ":", 2)
		if len(parts) == 2 {
			req.SetBasicAuth(parts[0], parts[1])
		}
	}

	// 设置 Bearer Token
	if config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+config.Token)
	}

	return req, nil
}
```

### 7.3 响应处理

```go
// internal/commands/http/response.go
package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ProcessResponse 处理 HTTP 响应
func ProcessResponse(resp *http.Response, config HttpConfig, startTime time.Time) (*HttpResult, error) {
	result := &HttpResult{
		URL:        resp.Request.URL.String(),
		Method:     resp.Request.Method,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		TimeCost:   time.Since(startTime),
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result.Body = body
	result.BodySize = int64(len(body))

	// 保存到文件
	if config.Output != "" {
		err = os.WriteFile(config.Output, body, 0644)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// PrintResult 输出请求结果
func PrintResult(result *HttpResult, config HttpConfig) {
	if config.ShowCode {
		fmt.Println(result.StatusCode)
		return
	}

	if config.OutputJSON {
		// JSON 格式输出
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
		return
	}

	if config.Silent {
		return
	}

	// 标准格式输出
	if config.Verbose {
		fmt.Printf("HTTP/%d.%d %d %s\n", 
			result.ProtoMajor, result.ProtoMinor,
			result.StatusCode, result.Status)
		fmt.Println("Headers:")
		for k, v := range result.Headers {
			fmt.Printf("  %s: %s\n", k, strings.Join(v, ", "))
		}
	}

	if config.ShowTime {
		fmt.Printf("Time: %v\n", result.TimeCost)
	}

	if !config.HeadOnly {
		if config.Hex {
			// 十六进制输出
			fmt.Println(hex.Dump(result.Body))
		} else if config.Pretty && json.Valid(result.Body) {
			// 格式化 JSON
			var prettyJSON bytes.Buffer
			json.Indent(&prettyJSON, result.Body, "", "  ")
			fmt.Println(prettyJSON.String())
		} else {
			fmt.Println(string(result.Body))
		}
	}
}
```

### 7.4 并发请求处理

```go
// internal/commands/http/concurrent.go
package http

import (
	"context"
	"sync"
	"time"
)

// ExecuteConcurrent 执行并发请求
func ExecuteConcurrent(config HttpConfig) (*HttpStats, error) {
	if config.Concurrent <= 1 {
		return nil, fmt.Errorf("concurrent must be > 1")
	}

	stats := &HttpStats{
		Total:      config.Concurrent,
		StatusDist: make(map[int]int),
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Concurrent)
	results := make(chan *HttpResult, config.Concurrent)

	start := time.Now()

	for i := 0; i < config.Concurrent; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 请求间隔
			if config.Interval > 0 && index > 0 {
				time.Sleep(config.Interval)
			}

			result := executeSingle(config)
			results <- result
		}(i)
	}

	// 关闭结果通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	for result := range results {
		if result.Error != "" {
			stats.Failed++
		} else {
			stats.Success++
			stats.StatusDist[result.StatusCode]++
		}
	}

	stats.TotalTime = time.Since(start)
	if stats.Success > 0 {
		stats.AvgTime = stats.TotalTime / time.Duration(stats.Success)
	}

	return stats, nil
}
```

## 8. 使用示例

### 8.1 基础用法

```bash
# 简单 GET 请求
fck http http://api.example.com/users

# POST JSON 数据
fck http -X POST -j -d '{"name":"test"}' http://api.example.com/users

# 添加请求头
fck http -H "Authorization: Bearer token" -H "X-Custom: value" http://api.example.com/data

# 发送表单数据
fck http -X POST -F "name=value" -F "file=@upload.txt" http://api.example.com/upload

# 下载文件
fck http --download -o file.zip http://example.com/file.zip

# 仅显示状态码
fck http --show-code http://example.com
```

### 8.2 HTTPS 用法

```bash
# 简单 HTTPS 请求
fck https https://api.example.com/secure

# 跳过 SSL 验证(测试用)
fck https -k https://self-signed.example.com

# 指定客户端证书
fck https --cert client.crt --key client.key https://api.example.com/mtls

# 指定 CA 证书
fck https --cacert ca.crt https://internal.example.com

# 强制 TLS 1.3
fck https --tls-version 1.3 https://example.com
```

### 8.3 高级用法

```bash
# 并发压力测试
fck http -c 100 -i 10ms http://api.example.com/test

# 带重试的请求
fck http --retry 3 --retry-delay 2s http://api.example.com/data

# 使用代理
fck http -x http://proxy:8080 http://example.com

# 自定义 DNS 解析
fck http --resolve "example.com:80:192.168.1.1" http://example.com

# 断点续传下载
fck http --download --continue -o large.zip http://example.com/large.zip

# 格式化 JSON 输出
fck http --pretty http://api.example.com/data

# 详细调试信息
fck http -v http://api.example.com/debug
```

## 9. 实现计划

### 第一阶段：基础功能
1. 创建目录结构
2. 实现 HTTP 基础请求功能
3. 实现 HTTPS 基础请求功能
4. 支持 GET/POST/PUT/DELETE 方法
5. 支持自定义请求头和请求体

### 第二阶段：增强功能
1. 实现并发请求
2. 实现重试机制
3. 实现下载功能(带进度条)
4. 支持表单数据和文件上传
5. 支持代理和自定义 DNS

### 第三阶段：高级功能
1. 实现 WebSocket 支持
2. 实现 HTTP/2 和 HTTP/3 支持
3. 实现断点续传
4. 完善 SSL/TLS 配置选项
5. 添加性能统计和报告

### 第四阶段：优化完善
1. 完善错误处理
2. 优化输出格式
3. 添加更多示例
4. 编写单元测试
5. 性能优化

## 10. 注意事项

1. **SSL 安全**: `-k/--insecure` 选项仅用于测试环境，生产环境应避免使用
2. **并发控制**: 大量并发请求可能对目标服务造成压力，请谨慎使用
3. **超时设置**: 合理设置超时时间，避免长时间挂起
4. **文件权限**: 下载文件时注意检查目录写入权限
5. **内存使用**: 大文件下载应使用流式处理，避免内存溢出
