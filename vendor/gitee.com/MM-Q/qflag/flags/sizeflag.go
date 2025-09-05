package flags

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// 单位常量定义
const (
	Byte = 1
	KiB  = 1024
	MiB  = KiB * 1024
	GiB  = MiB * 1024
	TiB  = GiB * 1024
	PiB  = TiB * 1024

	KB = 1000
	MB = KB * 1000
	GB = MB * 1000
	TB = GB * 1000
	PB = TB * 1000
)

// int64 边界值常量
const (
	maxInt64     = int64(1<<63 - 1)          // 9223372036854775807
	minInt64     = int64(-1 << 63)           // -9223372036854775808
	safeMaxFloat = float64(maxInt64) * 0.999 // 安全的最大浮点值
	safeMinFloat = float64(minInt64) * 0.999 // 安全的最小浮点值
)

// sizeUnits 单位映射表
var sizeUnits = map[string]int64{
	// 字节单位
	"bytes": Byte,
	"byte":  Byte,
	"b":     Byte,

	// 二进制单位 (1024) - IEC标准
	"pib": PiB,
	"tib": TiB,
	"gib": GiB,
	"mib": MiB,
	"kib": KiB,

	// 十进制单位 (1000) - SI标准
	"pb": PB,
	"tb": TB,
	"gb": GB,
	"mb": MB,
	"kb": KB,

	// 简写单位 (默认二进制)
	"p": PiB,
	"t": TiB,
	"g": GiB,
	"m": MiB,
	"k": KiB,
}

// formatUnits 格式化单位定义（使用标准十进制单位）
var formatUnits = []struct {
	name string
	size int64
}{
	{"PB", PB},
	{"TB", TB},
	{"GB", GB},
	{"MB", MB},
	{"KB", KB},
	{"B", Byte},
}

// sizeRegex 预编译的正则表达式，用于解析大小值（支持负数）
var sizeRegex = regexp.MustCompile(`^(-?\d+(?:\.\d+)?)\s*([a-zA-Z]+)$`)

// numberRegex 预编译的正则表达式，用于检查纯数字格式
var numberRegex = regexp.MustCompile(`^-?\d+(?:\.\d+)?$`)

// SizeFlag 大小标志结构体
type SizeFlag struct {
	BaseFlag[int64]
	allowDecimal  bool         // 是否允许小数
	allowNegative bool         // 是否允许负数
	mu            sync.RWMutex // 读写锁
	initOnce      sync.Once    // 确保只初始化一次
}

// Init 初始化大小标志（使用 sync.Once 确保只初始化一次）
//
// 参数:
//   - longName: 长标志名称
//   - shortName: 短标志字符
//   - defValue: 默认值（字节数）
//   - usage: 帮助说明
//
// 返回值:
//   - error: 初始化错误信息
//
// 注意: 重复调用此方法是安全的，后续调用将被忽略
func (f *SizeFlag) Init(longName, shortName string, defValue int64, usage string) error {
	var initErr error
	f.initOnce.Do(func() {
		f.mu.Lock()
		defer f.mu.Unlock()

		// 创建大小值指针
		sizePtr := new(int64)
		*sizePtr = defValue

		// 调用基类的 Init 方法
		err := f.BaseFlag.Init(longName, shortName, usage, sizePtr)
		if err != nil {
			initErr = err
			return
		}

		// 设置大小标志特有的属性
		f.allowDecimal = true   // 默认允许小数
		f.allowNegative = false // 默认不允许负数
	})
	return initErr
}

// Set 设置标志值
func (f *SizeFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 处理空值 - 设为 0
	if strings.TrimSpace(value) == "" {
		return f.BaseFlag.Set(int64(0))
	}

	// 解析大小值
	size, err := f.parseSize(value)
	if err != nil {
		return err // 直接返回原始错误，避免重复包装
	}

	// 调用基类方法设置值
	return f.BaseFlag.Set(size)
}

// String 返回标志的字符串表示
func (f *SizeFlag) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.formatSize(f.Get())
}

// Type 返回标志类型
func (f *SizeFlag) Type() FlagType {
	return FlagTypeSize
}

// SetAllowDecimal 设置是否允许小数
func (f *SizeFlag) SetAllowDecimal(allow bool) *SizeFlag {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowDecimal = allow
	return f
}

// SetAllowNegative 设置是否允许负数
func (f *SizeFlag) SetAllowNegative(allow bool) *SizeFlag {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowNegative = allow
	return f
}

// GetBytes 获取字节数
func (f *SizeFlag) GetBytes() int64 {
	return f.Get()
}

// GetKiB 获取KiB数
func (f *SizeFlag) GetKiB() float64 {
	return float64(f.Get()) / KiB
}

// GetMiB 获取MiB数
func (f *SizeFlag) GetMiB() float64 {
	return float64(f.Get()) / MiB
}

// GetGiB 获取GiB数
func (f *SizeFlag) GetGiB() float64 {
	return float64(f.Get()) / GiB
}

// GetTiB 获取TiB数
func (f *SizeFlag) GetTiB() float64 {
	return float64(f.Get()) / TiB
}

// GetPiB 获取PiB数
func (f *SizeFlag) GetPiB() float64 {
	return float64(f.Get()) / PiB
}

// IsZero 检查是否为零
func (f *SizeFlag) IsZero() bool {
	return f.Get() == 0
}

// IsPositive 检查是否为正数
func (f *SizeFlag) IsPositive() bool {
	return f.Get() > 0
}

// IsNegative 检查是否为负数
func (f *SizeFlag) IsNegative() bool {
	return f.Get() < 0
}

// GetAllowDecimal 获取是否允许小数设置
func (f *SizeFlag) GetAllowDecimal() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.allowDecimal
}

// GetAllowNegative 获取是否允许负数设置
func (f *SizeFlag) GetAllowNegative() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.allowNegative
}

// formatSize 格式化大小为人类可读的字符串（高效实现，避免递归）
//
// 参数：
//   - bytes：输入的字节数
//
// 返回：
//   - string：格式化后的字符串
func (f *SizeFlag) formatSize(bytes int64) string {
	// 处理零值
	if bytes == 0 {
		return "0B"
	}

	// 处理负数，避免递归调用
	negative := bytes < 0
	if negative {
		bytes = -bytes
	}

	var result string

	// 查找合适的单位
	for _, unit := range formatUnits {
		if bytes >= unit.size {
			// 检查是否为整数倍
			if bytes%unit.size == 0 {
				// 整数倍，使用整数除法
				quotient := bytes / unit.size
				result = strconv.FormatInt(quotient, 10) + unit.name
			} else {
				// 非整数倍，计算一位小数
				quotient := bytes / unit.size
				remainder := bytes % unit.size
				decimal := (remainder * 10) / unit.size
				result = strconv.FormatInt(quotient, 10) + "." + strconv.FormatInt(decimal, 10) + unit.name
			}
			break
		}
	}

	// 小于 1KB 的情况
	if result == "" {
		result = strconv.FormatInt(bytes, 10) + "B"
	}

	// 添加负号
	if negative {
		result = "-" + result
	}

	return result
}

// parseSize 解析大小字符串（使用正则表达式）
func (f *SizeFlag) parseSize(input string) (int64, error) {
	// 清理输入并转换为小写
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return 0, fmt.Errorf("empty input")
	}

	// 提取数字和单位
	numStr, multiplier, err := f.extractNumberAndUnit(input)
	if err != nil {
		return 0, err
	}

	// 解析数字部分
	num, err := f.parseNumber(numStr)
	if err != nil {
		return 0, err
	}

	// 检查负数权限
	if num < 0 && !f.allowNegative {
		return 0, fmt.Errorf("negative values not allowed")
	}

	// 计算最终大小并检查溢出
	return f.calculateSize(num, multiplier)
}

// extractNumberAndUnit 提取数字和单位 - 使用正则表达式实现
//
// 参数：
//   - input：输入字符串
//
// 返回：
//   - numStr：数字部分字符串
//   - multiplier：单位对应的乘数
//   - err：错误信息
func (f *SizeFlag) extractNumberAndUnit(input string) (string, int64, error) {
	// 特殊处理零值
	if input == "0" {
		return "0", 1, nil
	}

	// 使用正则表达式匹配
	matches := sizeRegex.FindStringSubmatch(input)
	if matches == nil {
		// 检查是否为纯数字（非零值）
		if f.isPureNumber(input) {
			return "", 0, fmt.Errorf("size value must include a unit (e.g., 1GB, 512MB, 1024KB) or use '0' for zero")
		}
		return "", 0, fmt.Errorf("invalid size format: expected format like '1GB', '512MB', '1.5GiB'")
	}

	numStr := matches[1]
	unitStr := strings.ToLower(matches[2])

	// 查找单位对应的乘数
	multiplier, exists := sizeUnits[unitStr]
	if !exists {
		return "", 0, fmt.Errorf("unrecognized unit '%s', supported units: B/Byte/Bytes, KB/MB/GB/TB/PB, KiB/MiB/GiB/TiB/PiB, K/M/G/T/P", matches[2])
	}

	return numStr, multiplier, nil
}

// isPureNumber 检查字符串是否为纯数字（使用预编译正则）
func (f *SizeFlag) isPureNumber(s string) bool {
	return numberRegex.MatchString(s)
}

// parseNumber 解析数字部分（优化版本）
func (f *SizeFlag) parseNumber(numStr string) (float64, error) {
	// 检查是否允许小数
	if strings.Contains(numStr, ".") && !f.allowDecimal {
		return 0, fmt.Errorf("decimal values not allowed")
	}

	// 统一使用 ParseFloat 处理所有数字格式
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format '%s'", numStr)
	}

	return num, nil
}

// calculateSize 计算最终大小，精确检查溢出
//
// 参数：
//   - num：输入的数值
//   - multiplier：单位对应的乘数
//
// 返回：
//   - int64：计算后的大小
//   - error：错误信息
func (f *SizeFlag) calculateSize(num float64, multiplier int64) (int64, error) {
	// 处理零值和单位为1的情况
	if multiplier == 0 {
		return 0, nil
	}
	if multiplier == 1 {
		if num > float64(maxInt64) || num < float64(minInt64) {
			return 0, fmt.Errorf("size out of range")
		}
		return int64(num), nil
	}

	// 检查是否为整数
	if num == float64(int64(num)) {
		// 整数情况，使用精确的整数溢出检查
		intNum := int64(num)

		// 正数溢出检查
		if intNum > 0 && multiplier > 0 {
			if intNum > maxInt64/multiplier {
				return 0, fmt.Errorf("size too large: overflow")
			}
		}
		// 负数溢出检查
		if intNum < 0 && multiplier > 0 {
			if intNum < minInt64/multiplier {
				return 0, fmt.Errorf("size too small: underflow")
			}
		}

		return intNum * multiplier, nil
	}

	// 小数情况，使用浮点运算但加强检查
	result := num * float64(multiplier)

	if result > safeMaxFloat {
		return 0, fmt.Errorf("size too large: overflow")
	}
	if result < safeMinFloat {
		return 0, fmt.Errorf("size too small: underflow")
	}

	return int64(result), nil
}
