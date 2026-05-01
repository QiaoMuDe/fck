package curl

import (
	"time"
)

// Config curl 命令配置
type Config struct {
	URL      string        // 请求 URL
	Method   string        // HTTP 方法
	Data     string        // 请求体数据
	Headers  []string      // 请求头列表
	Output   string        // 输出文件路径
	Include  bool          // 显示响应头
	Head     bool          // 仅显示响应头（使用 HEAD 方法）
	Silent   bool          // 静默模式
	Verbose  bool          // 详细模式
	Form     []string      // 表单数据
	User     string        // 认证信息
	Location bool          // 跟随重定向
	MaxTime  time.Duration // 最大执行时间
	Retry    int           // 重试次数
	Pretty   bool          // 格式化 JSON
	Color    bool          // 彩色输出
}

// Response HTTP 响应信息
type Response struct {
	StatusCode    int
	Status        string
	Headers       map[string][]string
	Body          []byte
	ContentLength int64
	Time          time.Duration
}
