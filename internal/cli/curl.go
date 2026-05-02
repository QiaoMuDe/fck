package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/curl"
	"gitee.com/MM-Q/qflag"
)

// CurlCmd curl 命令
var CurlCmd *qflag.Cmd

var (
	curlRequest  *qflag.EnumFlag
	curlData     *qflag.StringFlag
	curlHeader   *qflag.StringSliceFlag
	curlOutput   *qflag.StringFlag
	curlInclude  *qflag.BoolFlag
	curlHead     *qflag.BoolFlag
	curlSilent   *qflag.BoolFlag
	curlVerbose  *qflag.BoolFlag
	curlForm     *qflag.StringSliceFlag
	curlUser     *qflag.StringFlag
	curlLocation *qflag.BoolFlag
	curlMaxTime  *qflag.DurationFlag
	curlRetry    *qflag.IntFlag
	curlColor    *qflag.BoolFlag
	curlInsecure *qflag.BoolFlag
)

// requestMethods 支持的 HTTP 方法
var requestMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

func init() {
	CurlCmd = qflag.NewCmd("curl", "", qflag.ExitOnError)

	curlRequest = CurlCmd.Enum("request", "X", fmt.Sprintf("HTTP 方法, 支持: %v", requestMethods), "GET", requestMethods)
	curlData = CurlCmd.String("data", "d", "请求体数据", "")
	curlHeader = CurlCmd.StringSlice("header", "H", "请求头，多个用逗号分隔", []string{})
	curlOutput = CurlCmd.String("output", "o", "输出到文件", "")
	curlInclude = CurlCmd.Bool("include", "i", "显示响应头", false)
	curlHead = CurlCmd.Bool("head", "I", "仅显示响应头（使用 HEAD 方法）", false)
	curlSilent = CurlCmd.Bool("silent", "s", "静默模式，只输出响应体", false)
	curlVerbose = CurlCmd.Bool("verbose", "v", "显示完整请求/响应详情", false)
	curlForm = CurlCmd.StringSlice("form", "F", "multipart/form-data 数据，多个用逗号分隔", []string{})
	curlUser = CurlCmd.String("user", "u", "用户名密码认证 (user:password)", "")
	curlLocation = CurlCmd.Bool("location", "L", "跟随重定向", false)
	curlMaxTime = CurlCmd.Duration("max-time", "m", "最大执行时间", 0)
	curlRetry = CurlCmd.Int("retry", "R", "失败重试次数", 0)
	curlColor = CurlCmd.Bool("color", "c", "启用彩色输出", false)
	curlInsecure = CurlCmd.Bool("insecure", "k", "跳过 HTTPS 证书验证", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "HTTP 客户端工具",
		UsageSyntax: fmt.Sprintf("%s curl [options] <url>", qflag.Root.Name()),
		UseChinese:  true,
		Notes: []string{
			"类似 curl 的 HTTP 客户端，支持 GET/POST/PUT/DELETE 等方法",
			"使用 -v 查看完整请求/响应详情",
			"使用 -s 静默模式，便于管道处理",
		},
		Examples: map[string]string{
			"简单 GET":    fmt.Sprintf("%s curl https://api.example.com/users", qflag.Root.Name()),
			"POST JSON": fmt.Sprintf("%s curl -X POST -H \"Content-Type: application/json\" -d '{\"name\":\"test\"}' https://api.example.com/users", qflag.Root.Name()),
			"POST 表单":   fmt.Sprintf("%s curl -X POST -d 'username=admin&password=123' https://api.example.com/login", qflag.Root.Name()),
			"文件上传":      fmt.Sprintf("%s curl -F \"file=@/path/to/image.jpg\" https://api.example.com/upload", qflag.Root.Name()),
			"下载文件":      fmt.Sprintf("%s curl -o result.json https://api.example.com/users", qflag.Root.Name()),
			"显示响应头":     fmt.Sprintf("%s curl -i https://api.example.com/users", qflag.Root.Name()),
			"完整详情":      fmt.Sprintf("%s curl -v https://api.example.com/users", qflag.Root.Name()),
			"Basic 认证":  fmt.Sprintf("%s curl -u admin:password https://api.example.com/admin", qflag.Root.Name()),
			"跟随重定向":     fmt.Sprintf("%s curl -L https://bit.ly/xxx", qflag.Root.Name()),
		},
	}

	if err := CurlCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	CurlCmd.SetRun(runCurl)
}

// runCurl 运行 curl 命令
//
// 参数:
//   - cmd: 命令接口
//
// 返回值:
//   - error: 错误
func runCurl(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("URL not specified")
	}

	config := curl.Config{
		URL:      args[0],
		Method:   curlRequest.Get(),
		Data:     curlData.Get(),
		Headers:  curlHeader.Get(),
		Output:   curlOutput.Get(),
		Include:  curlInclude.Get(),
		Head:     curlHead.Get(),
		Silent:   curlSilent.Get(),
		Verbose:  curlVerbose.Get(),
		Form:     curlForm.Get(),
		User:     curlUser.Get(),
		Location: curlLocation.Get(),
		MaxTime:  curlMaxTime.Get(),
		Retry:    curlRetry.Get(),
		Color:    curlColor.Get(),
		Insecure: curlInsecure.Get(),
	}

	return curl.Execute(config)
}
