package flags

import (
	"fmt"
	"net"
	"sync"
)

// IP6Flag IPv6地址类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type IP6Flag struct {
	BaseFlag[string]
	mu sync.Mutex // 互斥锁
}

// Type 返回标志类型
func (f *IP6Flag) Type() FlagType { return FlagTypeIP6 }

// String 实现flag.Value接口,返回当前值的字符串表示
func (f *IP6Flag) String() string { return f.Get() }

// Set 实现flag.Value接口,解析并验证IPv6地址
func (f *IP6Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return fmt.Errorf("ipv6 address cannot be empty")
	}

	// 解析IPv6地址
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() != nil {
		return fmt.Errorf("invalid ipv6 address: %s", value)
	}

	return f.BaseFlag.Set(ip.String())
}
