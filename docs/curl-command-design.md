# Curl 子命令设计方案

## 1. 概述

设计一个 `curl` 子命令，用于执行 HTTP/HTTPS 请求测试、API 调试和 Web 服务探测。通过 URL 协议自动区分 HTTP 和 HTTPS。

## 2. 命令定位

| 命令 | 用途 | 说明 |
|------|------|------|
| `curl` | HTTP/HTTPS 请求工具 | 支持 http:// 和 https:// 协议 |

## 3. 目录结构

```
internal/
├── cli/
│   ├── curl.go          # Curl 命令 CLI 定义
│   └── root.go          # 注册 CurlCmd
├── commands/
│   └── curl/
│       ├── cmd_curl.go      # 主业务逻辑
│       ├── client.go        # HTTP 客户端封装
│       ├── request.go       # 请求构建
│       ├── response.go      # 响应处理
│       ├── output.go        # 输出格式化
│       └── tls.go           # TLS/SSL 相关
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
- 发送文件作为请求体 (--data-binary)
- 从文件读取请求体 (-T, --upload-file)
- URL 查询参数 (-G, --get + --data-urlencode)
- 表单数据提交 (-F, --form)
- JSON 自动 Content-Type (-j, --json)
- Cookie 设置 (-b, --cookie / -c, --cookie-jar)
- 用户代理设置 (-A, --user-agent)
- 基础认证 (-u, --user)
- Bearer Token 认证 (--oauth2-bearer / --token)

### 4.3 连接与性能
- 连接超时控制 (--connect-timeout)
- 总超时控制 (-m, --max-time)
- 请求重试次数 (--retry)
- 重试间隔 (--retry-delay)
- 并发请求 (-c, --concurrent)
- 请求间隔 (-i, --interval)
- 跟随重定向 (-L, --location)
- 最大重定向次数 (--max-redirs)
- 禁用 SSL 证书验证 (-k, --insecure)

### 4.4 输出控制
- 静默模式 (-s, --silent)
- 进度条显示 (-#, --progress-bar)
- 输出到文件 (-o, --output)
- 远程文件名 (-O, --remote-name)
- 仅显示响应头 (-I, --head)
- 显示请求头 (-v, --verbose)
- 显示详细过程 (--trace)
- 格式化 JSON 输出 (--pretty)
- 显示响应时间 (--write-out)

### 4.5 TLS/SSL 配置
- 客户端证书 (--cert, --key)
- CA 证书 (--cacert)
- 证书格式 (--cert-type)
- TLS 版本 (--tlsv1.0, --tlsv1.1, --tlsv1.2, --tlsv1.3)
- 禁用证书验证 (-k, --insecure)
- 指定密码套件 (--ciphers)

### 4.6 代理支持
- HTTP 代理 (-x, --proxy)
- 代理认证 (-U, --proxy-user)
- 代理隧道 (-p, --proxytunnel)
- SOCKS 代理 (--socks5)
- 不使用代理 (--noproxy)

### 4.7 高级功能
- 自定义 DNS 解析 (--resolve)
- HTTP 版本指定 (--http1.1, --http2, --http3)
- 断点续传 (-C, --continue-at)
- 限速 (--limit-rate)
- 最大下载大小 (--max-filesize)
- 自动跳转Referer (--referer)
- 自定义请求方法 (-X, --request)

## 5. 命令行接口设计

```go
// internal/cli/curl.go
package cli

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/fck/internal/commands/curl"
	"gitee.com/MM-Q/qflag"
)

var CurlCmd *qflag.Cmd

var (
	// 请求方法
	curlRequest     *qflag.StringFlag   // -X, --request    HTTP 方法
	curlHeader      *qflag.StringSliceFlag // -H, --header     请求头
	curlData        *qflag.StringFlag   // -d, --data       请求体数据
	curlDataBinary  *qflag.StringFlag   // --data-binary    二进制数据
	curlUploadFile  *qflag.StringFlag   // -T, --upload-file 上传文件
	curlGet         *qflag.BoolFlag     // -G, --get        将 data 作为查询参数
	curlForm        *qflag.StringSliceFlag // -F, --form       表单数据
	curlJSON        *qflag.BoolFlag     // -j, --json       发送 JSON
	curlCookie      *qflag.StringFlag   // -b, --cookie     发送 Cookie
	curlCookieJar   *qflag.StringFlag   // -c, --cookie-jar 保存 Cookie
	curlUserAgent   *qflag.StringFlag   // -A, --user-agent User-Agent
	curlReferer     *qflag.StringFlag   // -e, --referer    Referer

	// 认证
	curlUser        *qflag.StringFlag   // -u, --user       基础认证
	curlToken       *qflag.StringFlag   // --token          Bearer Token

	// 超时与重试
	curlConnectTimeout *qflag.DurationFlag // --connect-timeout 连接超时
	curlMaxTime        *qflag.DurationFlag // -m, --max-time   总超时
	curlRetry          *qflag.IntFlag      // --retry          重试次数
	curlRetryDelay     *qflag.DurationFlag // --retry-delay    重试间隔
	curlRetryMaxTime   *qflag.DurationFlag // --retry-max-time 最大重试时间

	// 并发控制
	curlConcurrent  *qflag.IntFlag      // -c, --concurrent 并发数
	curlInterval    *qflag.DurationFlag // -i, --interval   请求间隔

	// 重定向
	curlLocation    *qflag.BoolFlag     // -L, --location   跟随重定向
	curlMaxRedirs   *qflag.IntFlag      // --max-redirs     最大重定向数

	// SSL/TLS
	curlInsecure    *qflag.BoolFlag     // -k, --insecure   跳过 SSL 验证
	curlCert        *qflag.StringFlag   // --cert           客户端证书
	curlCertType    *qflag.StringFlag   // --cert-type      证书类型
	curlKey         *qflag.StringFlag   // --key            私钥
	curlKeyType     *qflag.StringFlag   // --key-type       私钥类型
	curlCACert      *qflag.StringFlag   // --cacert         CA 证书
	curlCiphers     *qflag.StringFlag   // --ciphers        密码套件
	curlTLSv10      *qflag.BoolFlag     // --tlsv1.0        强制 TLS 1.0
	curlTLSv11      *qflag.BoolFlag     // --tlsv1.1        强制 TLS 1.1
	curlTLSv12      *qflag.BoolFlag     // --tlsv1.2        强制 TLS 1.2
	curlTLSv13      *qflag.BoolFlag     // --tlsv1.3        强制 TLS 1.3

	// 代理
	curlProxy       *qflag.StringFlag   // -x, --proxy      代理地址
	curlProxyUser   *qflag.StringFlag   // -U, --proxy-user 代理认证
	curlProxyTunnel *qflag.BoolFlag     // -p, --proxytunnel 代理隧道
	curlSocks5      *qflag.StringFlag   // --socks5         SOCKS5 代理
	curlNoProxy     *qflag.StringFlag   // --noproxy        不使用代理的主机

	// 输出控制
	curlSilent      *qflag.BoolFlag     // -s, --silent     静默模式
	curlVerbose     *qflag.BoolFlag     // -v, --verbose    详细输出
	curlTrace       *qflag.StringFlag   // --trace          跟踪输出到文件
	curlHead        *qflag.BoolFlag     // -I, --head       仅响应头
	curlOutput      *qflag.StringFlag   // -o, --output     输出到文件
	curlRemoteName  *qflag.BoolFlag     // -O, --remote-name 使用远程文件名
	curlCreateDirs  *qflag.BoolFlag     // --create-dirs    自动创建目录
	curlProgressBar *qflag.BoolFlag     // -#, --progress-bar 显示进度条
	curlWriteOut    *qflag.StringFlag   // -w, --write-out  自定义输出格式
	curlPretty      *qflag.BoolFlag     // --pretty         格式化 JSON

	// 下载控制
	curlContinueAt  *qflag.StringFlag   // -C, --continue-at 断点续传
	curlLimitRate   *qflag.SizeFlag     // --limit-rate     限速
	curlMaxFileSize *qflag.SizeFlag     // --max-filesize   最大文件大小

	// 其他
	curlResolve     *qflag.StringSliceFlag // --resolve        自定义 DNS
	curlHTTP11      *qflag.BoolFlag     // --http1.1        强制 HTTP/1.1
	curlHTTP2       *qflag.BoolFlag     // --http2          强制 HTTP/2
	curlHTTP3       *qflag.BoolFlag     // --http3          尝试 HTTP/3
	curlInclude     *qflag.BoolFlag     // -i, --include    包含响应头
	curlFail        *qflag.BoolFlag     // -f, --fail       HTTP 错误时返回非零
	curlShowError   *qflag.BoolFlag     // -S, --show-error 显示错误信息
)

func init() {
	CurlCmd = qflag.NewCmd("curl", "", qflag.ExitOnError)

	// 请求方法
	curlRequest = CurlCmd.String("request", "X", "HTTP 方法 (GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS)", "GET")
	curlHeader = CurlCmd.StringSlice("header", "H", "请求头, 格式: 'Key: Value'", []string{})
	curlData = CurlCmd.String("data", "d", "请求体数据", "")
	curlDataBinary = CurlCmd.String("data-binary", "", "二进制请求体数据", "")
	curlUploadFile = CurlCmd.String("upload-file", "T", "上传文件路径", "")
	curlGet = CurlCmd.Bool("get", "G", "将 --data 数据作为 URL 查询参数", false)
	curlForm = CurlCmd.StringSlice("form", "F", "表单数据, 格式: 'name=value' 或 'name=@file'", []string{})
	curlJSON = CurlCmd.Bool("json", "j", "发送 JSON (自动设置 Content-Type)", false)
	curlCookie = CurlCmd.String("cookie", "b", "发送 Cookie 字符串或从文件读取", "")
	curlCookieJar = CurlCmd.String("cookie-jar", "c", "保存 Cookie 到文件", "")
	curlUserAgent = CurlCmd.String("user-agent", "A", "User-Agent", "fck-curl/1.0")
	curlReferer = CurlCmd.String("referer", "e", "Referer URL", "")

	// 认证
	curlUser = CurlCmd.String("user", "u", "基础认证, 格式: username:password", "")
	curlToken = CurlCmd.String("token", "", "Bearer Token 认证", "")

	// 超时与重试
	curlConnectTimeout = CurlCmd.Duration("connect-timeout", "", "连接超时时间", time.Second*30)
	curlMaxTime = CurlCmd.Duration("max-time", "m", "最大传输时间", 0)
	curlRetry = CurlCmd.Int("retry", "", "请求失败重试次数", 0)
	curlRetryDelay = CurlCmd.Duration("retry-delay", "", "重试间隔", time.Second)
	curlRetryMaxTime = CurlCmd.Duration("retry-max-time", "", "最大重试时间", 0)

	// 并发控制
	curlConcurrent = CurlCmd.Int("concurrent", "c", "并发请求数", 1)
	curlInterval = CurlCmd.Duration("interval", "i", "并发请求间隔", 0)

	// 重定向
	curlLocation = CurlCmd.Bool("location", "L", "跟随重定向", false)
	curlMaxRedirs = CurlCmd.Int("max-redirs", "", "最大重定向次数", 50)

	// SSL/TLS
	curlInsecure = CurlCmd.Bool("insecure", "k", "跳过 SSL 证书验证", false)
	curlCert = CurlCmd.String("cert", "", "客户端证书文件", "")
	curlCertType = CurlCmd.String("cert-type", "", "证书类型 (PEM/DER/P12)", "PEM")
	curlKey = CurlCmd.String("key", "", "私钥文件", "")
	curlKeyType = CurlCmd.String("key-type", "", "私钥类型 (PEM/DER)", "PEM")
	curlCACert = CurlCmd.String("cacert", "", "CA 证书文件", "")
	curlCiphers = CurlCmd.String("ciphers", "", "SSL 密码套件列表", "")
	curlTLSv10 = CurlCmd.Bool("tlsv1.0", "", "强制使用 TLS 1.0", false)
	curlTLSv11 = CurlCmd.Bool("tlsv1.1", "", "强制使用 TLS 1.1", false)
	curlTLSv12 = CurlCmd.Bool("tlsv1.2", "", "强制使用 TLS 1.2", false)
	curlTLSv13 = CurlCmd.Bool("tlsv1.3", "", "强制使用 TLS 1.3", false)

	// 代理
	curlProxy = CurlCmd.String("proxy", "x", "代理地址, 如: http://proxy:8080", "")
	curlProxyUser = CurlCmd.String("proxy-user", "U", "代理认证, 格式: username:password", "")
	curlProxyTunnel = CurlCmd.Bool("proxytunnel", "p", "通过代理建立隧道", false)
	curlSocks5 = CurlCmd.String("socks5", "", "SOCKS5 代理地址", "")
	curlNoProxy = CurlCmd.String("noproxy", "", "不使用代理的主机列表(逗号分隔)", "")

	// 输出控制
	curlSilent = CurlCmd.Bool("silent", "s", "静默模式", false)
	curlVerbose = CurlCmd.Bool("verbose", "v", "显示详细通信过程", false)
	curlTrace = CurlCmd.String("trace", "", "跟踪信息输出到文件", "")
	curlHead = CurlCmd.Bool("head", "I", "仅获取响应头", false)
	curlOutput = CurlCmd.String("output", "o", "保存响应到指定文件", "")
	curlRemoteName = CurlCmd.Bool("remote-name", "O", "使用远程文件名保存", false)
	curlCreateDirs = CurlCmd.Bool("create-dirs", "", "自动创建输出目录", false)
	curlProgressBar = CurlCmd.Bool("progress-bar", "#", "显示进度条", false)
	curlWriteOut = CurlCmd.String("write-out", "w", "请求完成后自定义输出格式", "")
	curlPretty = CurlCmd.Bool("pretty", "", "格式化 JSON 响应", false)

	// 下载控制
	curlContinueAt = CurlCmd.String("continue-at", "C", "断点续传偏移量或 '-' 自动", "")
	curlLimitRate = CurlCmd.Size("limit-rate", "", "下载限速, 如: 100K, 1M", 0)
	curlMaxFileSize = CurlCmd.Size("max-filesize", "", "最大下载文件大小", 0)

	// 其他
	curlResolve = CurlCmd.StringSlice("resolve", "", "自定义 DNS, 格式: 'host:port:ip'", []string{})
	curlHTTP11 = CurlCmd.Bool("http1.1", "", "强制使用 HTTP/1.1", false)
	curlHTTP2 = CurlCmd.Bool("http2", "", "强制使用 HTTP/2", false)
	curlHTTP3 = CurlCmd.Bool("http3", "", "尝试使用 HTTP/3", false)
	curlInclude = CurlCmd.Bool("include", "i", "输出包含响应头", false)
	curlFail = CurlCmd.Bool("fail", "f", "HTTP 错误(>=400)时返回非零退出码", false)
	curlShowError = CurlCmd.Bool("show-error", "S", "与 -s 一起使用时仍显示错误", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "HTTP/HTTPS 请求工具 (类 curl 风格)",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s curl [options] <URL>", qflag.Root.Name()),
		Notes: []string{
			"URL 必须包含协议头 (http:// 或 https://)",
			"默认使用 GET 方法, 使用 -X 指定其他方法",
			"使用 -d 发送数据时自动使用 POST 方法",
			"使用 -L 跟随重定向, --max-redirs 限制次数",
			"使用 -k 跳过 SSL 验证(仅测试环境使用)",
			"使用 -o 指定输出文件, -O 使用远程文件名",
			"使用 -c 并发请求进行压力测试",
		},
		Examples: map[string]string{
			"简单 GET 请求":           fmt.Sprintf("%s curl http://api.example.com/users", qflag.Root.Name()),
			"HTTPS 请求":              fmt.Sprintf("%s curl https://api.example.com/secure", qflag.Root.Name()),
			"POST JSON 数据":          fmt.Sprintf("%s curl -X POST -j -d '{\"name\":\"test\"}' http://api.example.com/users", qflag.Root.Name()),
			"添加请求头":              fmt.Sprintf("%s curl -H 'Authorization: Bearer token' http://api.example.com/data", qflag.Root.Name()),
			"发送表单":                fmt.Sprintf("%s curl -X POST -F 'name=value' -F 'file=@upload.txt' http://api.example.com/upload", qflag.Root.Name()),
			"下载文件":                fmt.Sprintf("%s curl -O http://example.com/file.zip", qflag.Root.Name()),
			"断点续传":                fmt.Sprintf("%s curl -C - -O http://example.com/large.zip", qflag.Root.Name()),
			"保存到指定文件":          fmt.Sprintf("%s curl -o output.json http://api.example.com/data", qflag.Root.Name()),
			"跟随重定向":              fmt.Sprintf("%s curl -L http://bit.ly/xxx", qflag.Root.Name()),
			"仅显示响应头":            fmt.Sprintf("%s curl -I http://example.com", qflag.Root.Name()),
			"显示详细过程":            fmt.Sprintf("%s curl -v http://example.com", qflag.Root.Name()),
			"基础认证":                fmt.Sprintf("%s curl -u admin:password http://api.example.com/secure", qflag.Root.Name()),
			"Bearer Token":            fmt.Sprintf("%s curl --token 'xyz123' http://api.example.com/data", qflag.Root.Name()),
			"跳过 SSL 验证":           fmt.Sprintf("%s curl -k https://self-signed.example.com", qflag.Root.Name()),
			"使用客户端证书":          fmt.Sprintf("%s curl --cert client.crt --key client.key https://api.example.com/mtls", qflag.Root.Name()),
			"使用代理":                fmt.Sprintf("%s curl -x http://proxy:8080 http://example.com", qflag.Root.Name()),
			"SOCKS5 代理":             fmt.Sprintf("%s curl --socks5 127.0.0.1:1080 http://example.com", qflag.Root.Name()),
			"并发压力测试":            fmt.Sprintf("%s curl -c 100 -i 10ms http://api.example.com/test", qflag.Root.Name()),
			"格式化 JSON 输出":        fmt.Sprintf("%s curl --pretty http://api.example.com/data", qflag.Root.Name()),
			"自定义 DNS":              fmt.Sprintf("%s curl --resolve 'example.com:443:192.168.1.1' https://example.com", qflag.Root.Name()),
			"限速下载":                fmt.Sprintf("%s curl --limit-rate 1M -O http://example.com/large.zip", qflag.Root.Name()),
			"保存 Cookie":             fmt.Sprintf("%s curl -c cookies.txt http://example.com/login", qflag.Root.Name()),
			"发送 Cookie":             fmt.Sprintf("%s curl -b cookies.txt http://example.com/profile", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "data-source",
				Flags:     []string{"data", "data-binary", "upload-file", "form"},
				AllowNone: true,
			},
			{
				Name:      "output-file",
				Flags:     []string{"output", "remote-name"},
				AllowNone: true,
			},
			{
				Name:      "tls-version",
				Flags:     []string{"tlsv1.0", "tlsv1.1", "tlsv1.2", "tlsv1.3"},
				AllowNone: true,
			},
			{
				Name:      "http-version",
				Flags:     []string{"http1.1", "http2", "http3"},
				AllowNone: true,
			},
		},
	}

	if err := CurlCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	CurlCmd.SetRun(runCurl)
}

func runCurl(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("缺少 URL 参数")
	}

	url := args[0]

	// 确定 TLS 版本
	tlsVersion := ""
	switch {
	case curlTLSv10.Get():
		tlsVersion = "1.0"
	case curlTLSv11.Get():
		tlsVersion = "1.1"
	case curlTLSv12.Get():
		tlsVersion = "1.2"
	case curlTLSv13.Get():
		tlsVersion = "1.3"
	}

	// 确定 HTTP 版本
	httpVersion := ""
	switch {
	case curlHTTP11.Get():
		httpVersion = "1.1"
	case curlHTTP2.Get():
		httpVersion = "2"
	case curlHTTP3.Get():
		httpVersion = "3"
	}

	config := curl.CurlConfig{
		URL:            url,
		Method:         curlRequest.Get(),
		Headers:        curlHeader.Get(),
		Data:           curlData.Get(),
		DataBinary:     curlDataBinary.Get(),
		UploadFile:     curlUploadFile.Get(),
		Get:            curlGet.Get(),
		Form:           curlForm.Get(),
		JSON:           curlJSON.Get(),
		Cookie:         curlCookie.Get(),
		CookieJar:      curlCookieJar.Get(),
		UserAgent:      curlUserAgent.Get(),
		Referer:        curlReferer.Get(),
		User:           curlUser.Get(),
		Token:          curlToken.Get(),
		ConnectTimeout: curlConnectTimeout.Get(),
		MaxTime:        curlMaxTime.Get(),
		Retry:          curlRetry.Get(),
		RetryDelay:     curlRetryDelay.Get(),
		RetryMaxTime:   curlRetryMaxTime.Get(),
		Concurrent:     curlConcurrent.Get(),
		Interval:       curlInterval.Get(),
		Location:       curlLocation.Get(),
		MaxRedirs:      curlMaxRedirs.Get(),
		Insecure:       curlInsecure.Get(),
		Cert:           curlCert.Get(),
		CertType:       curlCertType.Get(),
		Key:            curlKey.Get(),
		KeyType:        curlKeyType.Get(),
		CACert:         curlCACert.Get(),
		Ciphers:        curlCiphers.Get(),
		TLSVersion:     tlsVersion,
		Proxy:          curlProxy.Get(),
		ProxyUser:      curlProxyUser.Get(),
		ProxyTunnel:    curlProxyTunnel.Get(),
		Socks5:         curlSocks5.Get(),
		NoProxy:        curlNoProxy.Get(),
		Silent:         curlSilent.Get(),
		Verbose:        curlVerbose.Get(),
		Trace:          curlTrace.Get(),
		Head:           curlHead.Get(),
		Output:         curlOutput.Get(),
		RemoteName:     curlRemoteName.Get(),
		CreateDirs:     curlCreateDirs.Get(),
		ProgressBar:    curlProgressBar.Get(),
		WriteOut:       curlWriteOut.Get(),
		Pretty:         curlPretty.Get(),
		ContinueAt:     curlContinueAt.Get(),
		LimitRate:      curlLimitRate.Get(),
		MaxFileSize:    curlMaxFileSize.Get(),
		Resolve:        curlResolve.Get(),
		HTTPVersion:    httpVersion,
		Include:        curlInclude.Get(),
		Fail:           curlFail.Get(),
		ShowError:      curlShowError.Get(),
	}

	return curl.CurlCmdMain(config)
}
```

## 6. 配置结构体设计

```go
// internal/commands/curl/cmd_curl.go
package curl

import "time"

// CurlConfig curl 命令配置
type CurlConfig struct {
	// 请求配置
	URL        string   // 目标 URL
	Method     string   // HTTP 方法
	Headers    []string // 请求头列表
	Data       string   // 请求体数据
	DataBinary string   // 二进制数据
	UploadFile string   // 上传文件路径
	Get        bool     // 将 data 作为查询参数
	Form       []string // 表单数据
	JSON       bool     // 发送 JSON
	Cookie     string   // Cookie 字符串或文件
	CookieJar  string   // Cookie 保存文件
	UserAgent  string   // User-Agent
	Referer    string   // Referer

	// 认证
	User  string // 基础认证
	Token string // Bearer Token

	// 超时与重试
	ConnectTimeout time.Duration // 连接超时
	MaxTime        time.Duration // 总超时
	Retry          int           // 重试次数
	RetryDelay     time.Duration // 重试间隔
	RetryMaxTime   time.Duration // 最大重试时间

	// 并发控制
	Concurrent int           // 并发数
	Interval   time.Duration // 请求间隔

	// 重定向
	Location  bool // 跟随重定向
	MaxRedirs int  // 最大重定向数

	// SSL/TLS
	Insecure   bool   // 跳过验证
	Cert       string // 客户端证书
	CertType   string // 证书类型
	Key        string // 私钥
	KeyType    string // 私钥类型
	CACert     string // CA 证书
	Ciphers    string // 密码套件
	TLSVersion string // TLS 版本

	// 代理
	Proxy       string // 代理地址
	ProxyUser   string // 代理认证
	ProxyTunnel bool   // 代理隧道
	Socks5      string // SOCKS5 代理
	NoProxy     string // 不使用代理的主机

	// 输出控制
	Silent      bool   // 静默模式
	Verbose     bool   // 详细输出
	Trace       string // 跟踪文件
	Head        bool   // 仅响应头
	Output      string // 输出文件
	RemoteName  bool   // 使用远程文件名
	CreateDirs  bool   // 自动创建目录
	ProgressBar bool   // 进度条
	WriteOut    string // 自定义输出格式
	Pretty      bool   // 格式化 JSON

	// 下载控制
	ContinueAt  string // 断点续传
	LimitRate   int64  // 限速
	MaxFileSize int64  // 最大文件大小

	// 其他
	Resolve     []string // 自定义 DNS
	HTTPVersion string   // HTTP 版本
	Include     bool     // 包含响应头
	Fail        bool     // 错误时非零退出
	ShowError   bool     // 显示错误
}

// CurlResult 请求结果
type CurlResult struct {
	URL        string            `json:"url"`
	Method     string            `json:"method"`
	StatusCode int               `json:"status_code"`
	Status     string            `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body,omitempty"`
	BodySize   int64             `json:"body_size"`
	TimeTotal  time.Duration     `json:"time_total_ms"`
	TimeDNS    time.Duration     `json:"time_dns_ms"`
	TimeConnect time.Duration    `json:"time_connect_ms"`
	TimeTLS    time.Duration     `json:"time_tls_ms"`
	Error      string            `json:"error,omitempty"`
}

// CurlStats 并发统计
type CurlStats struct {
	Total       int           // 总请求数
	Success     int           // 成功数
	Failed      int           // 失败数
	TotalTime   time.Duration // 总耗时
	AvgTime     time.Duration // 平均耗时
	MinTime     time.Duration // 最小耗时
	MaxTime     time.Duration // 最大耗时
	StatusDist  map[int]int   // 状态码分布
	SizeTotal   int64         // 总下载大小
	SizeAvg     int64         // 平均下载大小
}

// CurlCmdMain curl 命令主函数
func CurlCmdMain(config CurlConfig) error {
	// 实现逻辑
	return nil
}
```

## 7. 核心功能实现思路

### 7.1 客户端封装

```go
// internal/commands/curl/client.go
package curl

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Client HTTP 客户端
type Client struct {
	httpClient *http.Client
	config     CurlConfig
}

// NewClient 创建客户端
func NewClient(config CurlConfig) (*Client, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.ConnectTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: config.MaxTime,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// 配置 TLS
	if err := configureTLS(transport, config); err != nil {
		return nil, err
	}

	// 配置代理
	if err := configureProxy(transport, config); err != nil {
		return nil, err
	}

	// 配置 HTTP 版本
	configureHTTPVersion(transport, config)

	// 配置自定义 DNS
	if len(config.Resolve) > 0 {
		configureDNS(transport, config)
	}

	return &Client{
		httpClient: &http.Client{
			Transport:     transport,
			Timeout:       config.MaxTime,
			CheckRedirect: makeRedirectPolicy(config),
		},
		config: config,
	}, nil
}
```

### 7.2 请求构建

```go
// internal/commands/curl/request.go
package curl

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// BuildRequest 构建 HTTP 请求
func BuildRequest(config CurlConfig) (*http.Request, error) {
	u, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	// 确定请求方法
	method := config.Method
	if config.Head {
		method = "HEAD"
	} else if config.Data != "" && method == "GET" {
		method = "POST"
	}

	// 构建请求体
	body, contentType, err := buildBody(config)
	if err != nil {
		return nil, err
	}

	// 处理 -G 参数 (将 data 作为查询参数)
	if config.Get && config.Data != "" {
		q := u.Query()
		// 解析 data 为查询参数
		dataParams, _ := url.ParseQuery(config.Data)
		for k, v := range dataParams {
			q[k] = v
		}
		u.RawQuery = q.Encode()
		body = nil
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	setHeaders(req, config, contentType)

	// 设置认证
	setAuth(req, config)

	return req, nil
}

// buildBody 构建请求体
func buildBody(config CurlConfig) (io.Reader, string, error) {
	// 表单数据
	if len(config.Form) > 0 {
		return buildFormBody(config.Form)
	}

	// 上传文件
	if config.UploadFile != "" {
		return buildFileBody(config.UploadFile)
	}

	// 二进制数据
	if config.DataBinary != "" {
		return strings.NewReader(config.DataBinary), "application/octet-stream", nil
	}

	// 普通数据
	if config.Data != "" {
		return strings.NewReader(config.Data), "", nil
	}

	return nil, "", nil
}

// buildFormBody 构建表单请求体
func buildFormBody(form []string) (io.Reader, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, f := range form {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		value := parts[1]

		// 文件上传 (@file)
		if strings.HasPrefix(value, "@") {
			filePath := value[1:]
			file, err := os.Open(filePath)
			if err != nil {
				return nil, "", err
			}
			defer file.Close()

			part, err := writer.CreateFormFile(name, filepath.Base(filePath))
			if err != nil {
				return nil, "", err
			}
			io.Copy(part, file)
		} else {
			writer.WriteField(name, value)
		}
	}

	writer.Close()
	return &buf, writer.FormDataContentType(), nil
}
```

### 7.3 响应处理与输出

```go
// internal/commands/curl/output.go
package curl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Output 处理输出
func Output(result *CurlResult, config CurlConfig) error {
	// 静默模式
	if config.Silent && !config.ShowError {
		return nil
	}

	// 仅显示状态码
	if config.Fail && result.StatusCode >= 400 {
		if config.ShowError {
			fmt.Fprintf(os.Stderr, "curl: (22) HTTP error %d\n", result.StatusCode)
		}
		return fmt.Errorf("HTTP error %d", result.StatusCode)
	}

	// 保存到文件
	if config.Output != "" || config.RemoteName {
		return saveToFile(result, config)
	}

	// 自定义输出格式
	if config.WriteOut != "" {
		return writeOut(result, config.WriteOut)
	}

	// 标准输出
	return standardOutput(result, config)
}

// standardOutput 标准输出
func standardOutput(result *CurlResult, config CurlConfig) error {
	// 包含响应头
	if config.Include || config.Verbose {
		fmt.Printf("HTTP/1.1 %d %s\r\n", result.StatusCode, result.Status)
		for k, v := range result.Headers {
			fmt.Printf("%s: %s\r\n", k, v)
		}
		fmt.Println()
	}

	// 仅响应头模式
	if config.Head {
		return nil
	}

	// 格式化 JSON
	if config.Pretty && json.Valid(result.Body) {
		var prettyJSON bytes.Buffer
		json.Indent(&prettyJSON, result.Body, "", "  ")
		fmt.Println(prettyJSON.String())
	} else {
		fmt.Print(string(result.Body))
	}

	return nil
}

// saveToFile 保存到文件
func saveToFile(result *CurlResult, config CurlConfig) error {
	filename := config.Output
	if config.RemoteName {
		// 从 URL 提取文件名
		u, _ := url.Parse(result.URL)
		filename = path.Base(u.Path)
		if filename == "" || filename == "/" {
			filename = "index.html"
		}
	}

	// 自动创建目录
	if config.CreateDirs {
		dir := filepath.Dir(filename)
		if dir != "" {
			os.MkdirAll(dir, 0755)
		}
	}

	return os.WriteFile(filename, result.Body, 0644)
}

// writeOut 自定义输出格式
func writeOut(result *CurlResult, format string) error {
	// 支持的变量:
	// %{url}, %{method}, %{http_code}, %{size_download}, %{time_total}, etc.
	output := format
	output = strings.ReplaceAll(output, "%{url}", result.URL)
	output = strings.ReplaceAll(output, "%{method}", result.Method)
	output = strings.ReplaceAll(output, "%{http_code}", fmt.Sprintf("%d", result.StatusCode))
	output = strings.ReplaceAll(output, "%{size_download}", fmt.Sprintf("%d", result.BodySize))
	output = strings.ReplaceAll(output, "%{time_total}", fmt.Sprintf("%.3f", result.TimeTotal.Seconds()))

	fmt.Println(output)
	return nil
}
```

### 7.4 并发请求

```go
// internal/commands/curl/concurrent.go
package curl

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// ExecuteConcurrent 执行并发请求
func ExecuteConcurrent(config CurlConfig) (*CurlStats, error) {
	if config.Concurrent <= 1 {
		return nil, fmt.Errorf("concurrent must be > 1")
	}

	stats := &CurlStats{
		Total:      config.Concurrent,
		StatusDist: make(map[int]int),
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.Concurrent)
	results := make(chan *CurlResult, config.Concurrent)

	var minTime, maxTime int64 = -1, 0
	var totalSize int64

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

			// 更新统计
			if result.Error == "" {
				atomic.AddInt64(&totalSize, result.BodySize)
				timeMs := result.TimeTotal.Milliseconds()
				if minTime == -1 || timeMs < minTime {
					atomic.StoreInt64(&minTime, timeMs)
				}
				if timeMs > maxTime {
					atomic.StoreInt64(&maxTime, timeMs)
				}
			}
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
		stats.MinTime = time.Duration(minTime) * time.Millisecond
		stats.MaxTime = time.Duration(maxTime) * time.Millisecond
		stats.SizeTotal = totalSize
		stats.SizeAvg = totalSize / int64(stats.Success)
	}

	return stats, nil
}
```

## 8. 使用示例

### 8.1 基础用法

```bash
# 简单 GET 请求
fck curl http://api.example.com/users

# HTTPS 请求
fck curl https://api.example.com/secure

# POST JSON 数据
fck curl -X POST -j -d '{"name":"test"}' http://api.example.com/users

# 添加请求头
fck curl -H "Authorization: Bearer token" http://api.example.com/data

# 发送表单
fck curl -X POST -F "name=value" -F "file=@upload.txt" http://api.example.com/upload
```

### 8.2 下载功能

```bash
# 下载文件
fck curl -O http://example.com/file.zip

# 保存到指定文件
fck curl -o output.json http://api.example.com/data

# 断点续传
fck curl -C - -O http://example.com/large.zip

# 限速下载
fck curl --limit-rate 1M -O http://example.com/large.zip

# 显示进度条
fck curl -# -O http://example.com/file.zip
```

### 8.3 SSL/TLS

```bash
# 跳过 SSL 验证
fck curl -k https://self-signed.example.com

# 使用客户端证书
fck curl --cert client.crt --key client.key https://api.example.com/mtls

# 指定 CA 证书
fck curl --cacert ca.crt https://internal.example.com

# 强制 TLS 1.3
fck curl --tlsv1.3 https://example.com
```

### 8.4 代理

```bash
# HTTP 代理
fck curl -x http://proxy:8080 http://example.com

# 代理认证
fck curl -x http://proxy:8080 -U user:pass http://example.com

# SOCKS5 代理
fck curl --socks5 127.0.0.1:1080 http://example.com

# 不使用代理
fck curl --noproxy "*.local,localhost" http://example.com
```

### 8.5 高级用法

```bash
# 并发压力测试
fck curl -c 100 -i 10ms http://api.example.com/test

# 带重试的请求
fck curl --retry 3 --retry-delay 2s http://api.example.com/data

# 自定义 DNS
fck curl --resolve "example.com:443:192.168.1.1" https://example.com

# 格式化 JSON 输出
fck curl --pretty http://api.example.com/data

# 自定义输出格式
fck curl -w "Status: %{http_code}, Time: %{time_total}s" http://example.com

# 保存和发送 Cookie
fck curl -c cookies.txt http://example.com/login
fck curl -b cookies.txt http://example.com/profile

# 详细调试
fck curl -v http://example.com
fck curl --trace trace.txt http://example.com
```

## 9. 实现计划

### 第一阶段：基础功能
1. 创建目录结构
2. 实现基本 HTTP/HTTPS 请求
3. 支持 GET/POST/PUT/DELETE 方法
4. 支持自定义请求头和请求体
5. 支持 -o/-O 输出到文件

### 第二阶段：增强功能
1. 实现表单提交 (-F)
2. 实现 Cookie 支持 (-b/-c)
3. 实现认证功能 (-u/--token)
4. 实现重定向跟随 (-L)
5. 实现代理支持 (-x/--socks5)

### 第三阶段：SSL/TLS 功能
1. 实现证书验证控制 (-k)
2. 实现客户端证书 (--cert/--key)
3. 实现 CA 证书 (--cacert)
4. 实现 TLS 版本控制
5. 实现密码套件选择

### 第四阶段：高级功能
1. 实现并发请求 (-c)
2. 实现重试机制 (--retry)
3. 实现断点续传 (-C)
4. 实现限速 (--limit-rate)
5. 实现 HTTP/2 和 HTTP/3

### 第五阶段：优化完善
1. 实现进度条显示
2. 实现详细跟踪 (--trace)
3. 完善错误处理
4. 添加更多示例
5. 编写单元测试

## 10. 与 curl 命令的兼容性

### 兼容的常用选项
| 选项 | 说明 |
|------|------|
| -X, --request | HTTP 方法 |
| -H, --header | 请求头 |
| -d, --data | 请求体数据 |
| -F, --form | 表单数据 |
| -o, --output | 输出文件 |
| -O, --remote-name | 使用远程文件名 |
| -L, --location | 跟随重定向 |
| -I, --head | 仅响应头 |
| -v, --verbose | 详细输出 |
| -s, --silent | 静默模式 |
| -u, --user | 基础认证 |
| -k, --insecure | 跳过 SSL 验证 |
| -x, --proxy | 代理 |
| -A, --user-agent | User-Agent |
| -b, --cookie | Cookie |
| -c, --cookie-jar | Cookie 保存 |
| -e, --referer | Referer |
| -C, --continue-at | 断点续传 |
| --max-time | 最大时间 |
| --connect-timeout | 连接超时 |

### 扩展选项
| 选项 | 说明 |
|------|------|
| -j, --json | 自动设置 JSON Content-Type |
| --token | Bearer Token 认证 |
| -c, --concurrent | 并发请求 |
| -i, --interval | 请求间隔 |
| --pretty | 格式化 JSON |
| --trace | 跟踪输出到文件 |

## 11. 注意事项

1. **SSL 安全**: `-k/--insecure` 仅用于测试环境，生产环境应避免使用
2. **并发控制**: 大量并发可能对目标服务造成压力，请谨慎使用
3. **超时设置**: 合理设置超时时间，避免长时间挂起
4. **文件权限**: 下载文件时注意检查目录写入权限
5. **内存使用**: 大文件下载应使用流式处理，避免内存溢出
6. **Cookie 安全**: Cookie 文件可能包含敏感信息，注意权限控制
