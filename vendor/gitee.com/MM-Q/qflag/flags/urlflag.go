package flags

import (
	"net/url"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
)

// URLFlag URL类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type URLFlag struct {
	BaseFlag[string]
	mu sync.Mutex // 互斥锁
}

// Type 返回标志类型
func (f *URLFlag) Type() FlagType { return FlagTypeURL }

// String 实现flag.Value接口,返回当前值的字符串表示
func (f *URLFlag) String() string { return f.Get() }

// Set 实现flag.Value接口,解析并验证URL格式
func (f *URLFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return qerr.NewValidationError("url cannot be empty")
	}

	// 解析URL
	parsedURL, err := url.ParseRequestURI(value)
	if err != nil {
		return qerr.NewValidationErrorf("invalid url format: %v", err)
	}

	// 检查URL是否包含有效的方案
	if parsedURL.Scheme == "" {
		return qerr.NewValidationError("url must include scheme (http/https)")
	}

	return f.BaseFlag.Set(parsedURL.String())
}
