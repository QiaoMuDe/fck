package curl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Execute 执行 curl 命令
//
// 参数:
//   - config: 配置
//
// 返回值:
//   - error: 错误
func Execute(config Config) error {
	// 设置超时
	ctx := context.Background()
	if config.MaxTime > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.MaxTime)
		defer cancel()
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !config.Location {
				return http.ErrUseLastResponse
			}
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// 构建请求
	req, err := buildRequest(ctx, config)
	if err != nil {
		return err
	}

	// 执行请求（带重试）
	var resp *http.Response
	startTime := time.Now()

	for attempt := 0; attempt <= config.Retry; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		resp, err = client.Do(req)
		if err == nil {
			break
		}

		if ctx.Err() != nil {
			return fmt.Errorf("request timeout")
		}
	}

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close response body: %v\n", err)
		}
	}()

	elapsed := time.Since(startTime)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// 构建响应对象
	response := &Response{
		StatusCode:    resp.StatusCode,
		Status:        resp.Status,
		Headers:       resp.Header,
		Body:          body,
		ContentLength: resp.ContentLength,
		Time:          elapsed,
	}

	// 输出响应
	return outputResponse(response, config)
}

// buildRequest 构建 HTTP 请求
//
// 参数:
//   - ctx: 上下文
//   - config: 配置
//
// 返回值:
//   - *http.Request: 请求对象
//   - error: 错误
func buildRequest(ctx context.Context, config Config) (*http.Request, error) {
	method := strings.ToUpper(config.Method)
	if method == "" {
		method = "GET"
	}

	// 如果指定了 -I/--head，强制使用 HEAD 方法
	if config.Head {
		method = "HEAD"
	}

	// 自动添加 http:// 前缀
	url := config.URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	var body io.Reader
	var contentType string

	// 处理表单数据
	if len(config.Form) > 0 {
		body, contentType = buildFormBody(config.Form)
	} else if config.Data != "" {
		body = strings.NewReader(config.Data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for _, header := range config.Headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			req.Header.Set(key, value)
		}
	}

	// 设置认证
	if config.User != "" {
		parts := strings.SplitN(config.User, ":", 2)
		if len(parts) == 2 {
			req.SetBasicAuth(parts[0], parts[1])
		} else {
			req.SetBasicAuth(parts[0], "")
		}
	}

	// 设置默认 User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "fck-curl/1.0")
	}

	return req, nil
}

// buildFormBody 构建表单请求体
//
// 参数:
//   - forms: 表单数据列表
//
// 返回值:
//   - io.Reader: 请求体
//   - string: Content-Type
func buildFormBody(forms []string) (io.Reader, string) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	for _, form := range forms {
		parts := strings.SplitN(form, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// 检查是否是文件上传
		if strings.HasPrefix(value, "@") {
			filePath := value[1:]
			file, err := os.Open(filePath)
			if err != nil {
				continue
			}
			func() {
				defer func() {
					if err := file.Close(); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", err)
					}
				}()

				part, err := writer.CreateFormFile(key, filepath.Base(filePath))
				if err != nil {
					return
				}
				if _, err := io.Copy(part, file); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to copy file: %v\n", err)
				}
			}()
		} else {
			if err := writer.WriteField(key, value); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to write field: %v\n", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to close writer: %v\n", err)
	}
	return &b, writer.FormDataContentType()
}

// outputResponse 输出响应
//
// 参数:
//   - resp: 响应对象
//   - config: 配置
//
// 返回值:
//   - error: 错误
func outputResponse(resp *Response, config Config) error {
	// 格式化输出
	formatter := NewFormatter(config.Color, config.Pretty)

	// Head 模式：仅显示响应头
	if config.Head {
		formatter.PrintHeaders(resp)
		return nil
	}

	// 静默模式：只输出响应体
	if config.Silent {
		if config.Output != "" {
			return os.WriteFile(config.Output, resp.Body, 0644)
		}
		_, err := os.Stdout.Write(resp.Body)
		return err
	}

	// 保存到文件
	if config.Output != "" {
		return os.WriteFile(config.Output, resp.Body, 0644)
	}

	// Verbose 模式：显示请求信息
	if config.Verbose {
		formatter.PrintVerbose(resp)
		return nil
	}

	// 显示响应头
	if config.Include {
		formatter.PrintHeaders(resp)
	}

	// 显示响应体
	formatter.PrintBody(resp)

	return nil
}
