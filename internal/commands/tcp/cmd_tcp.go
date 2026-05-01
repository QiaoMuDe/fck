// Package tcp 实现 TCP 网络工具功能，包括端口扫描、客户端通信和服务端监听
package tcp

import (
	"time"
)

// 端口状态常量
const (
	PortOpen     = "open"
	PortClosed   = "closed"
	PortFiltered = "filtered"
	PortTimeout  = "timeout"
)

// 默认超时时间
const (
	DefaultTimeout       = 2 * time.Second
	DefaultClientTimeout = 10 * time.Second
	DefaultServerTimeout = 30 * time.Second
	DefaultBufferSize    = 4 * 1024 // 4KB
)

// ScanConfig 端口扫描配置
type ScanConfig struct {
	Target     string
	Ports      string
	Timeout    time.Duration
	Concurrent int
	ShowClosed bool
	Output     string
	Format     string
	Progress   bool
	Summary    bool
}

// PortResult 端口扫描结果
type PortResult struct {
	Port     int
	Status   string
	Service  string
	Response time.Duration
}

// ScanResult 扫描结果集合
type ScanResult struct {
	Target    string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Results   []PortResult
	Stats     ScanStats
}

// ScanStats 扫描统计
type ScanStats struct {
	Total    int
	Open     int
	Closed   int
	Filtered int
	Timeout  int
}

// ClientConfig TCP 客户端配置
type ClientConfig struct {
	Address    string
	Message    string
	Timeout    time.Duration
	BufferSize int
	NoResponse bool
	Delimiter  string
	Hex        bool // 发送十六进制数据
}

// ServerConfig TCP 服务端配置
type ServerConfig struct {
	Address    string
	Port       int
	Timeout    time.Duration
	MaxConn    int
	BufferSize int
	Response   string
	OutputDir  string
	Echo       bool
}

// TransferStats 传输统计
type TransferStats struct {
	BytesSent     int64
	BytesReceived int64
	Duration      time.Duration
	FilesSent     int
	StartTime     time.Time // 连接开始时间
}
