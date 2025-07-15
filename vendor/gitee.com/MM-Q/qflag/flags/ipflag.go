package flags

import (
	"net"
	"sync"

	"gitee.com/MM-Q/qflag/qerr"
)

// IP4Flag IPv4地址类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type IP4Flag struct {
	BaseFlag[string]
	mu sync.Mutex // 互斥锁
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *IP4Flag) Type() FlagType { return FlagTypeIP4 }

// String 实现flag.Value接口,返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
func (f *IP4Flag) String() string { return f.Get() }

// Set 实现flag.Value接口,解析并验证IPv4地址
//
// 参数:
//   - value: 待解析的IPv4地址值
//
// 返回值:
//   - error: 解析或验证错误
func (f *IP4Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return qerr.NewValidationError("ipv4 address cannot be empty")
	}

	// 解析IPv4地址
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() == nil {
		return qerr.NewValidationErrorf("invalid ipv4 address: %s", value)
	}

	return f.BaseFlag.Set(ip.String())
}

// IP6Flag IPv6地址类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type IP6Flag struct {
	BaseFlag[string]
	mu sync.Mutex // 互斥锁
}

// Type 返回标志类型
//
// 返回值:
//   - FlagType: 标志类型枚举值
func (f *IP6Flag) Type() FlagType { return FlagTypeIP6 }

// String 实现flag.Value接口,返回当前值的字符串表示
//
// 返回值:
//   - string: 当前值的字符串表示
func (f *IP6Flag) String() string { return f.Get() }

// Set 实现flag.Value接口,解析并验证IPv6地址
//
// 参数:
//   - value: 待解析的IPv6地址值
//
// 返回值:
//   - error: 解析或验证错误
func (f *IP6Flag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return qerr.NewValidationError("ipv6 address cannot be empty")
	}

	// 解析IPv6地址
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() != nil {
		return qerr.NewValidationErrorf("invalid ipv6 address: %s", value)
	}

	return f.BaseFlag.Set(ip.String())
}
