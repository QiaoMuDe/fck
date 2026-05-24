package str

import (
	"encoding/base64"
	"strings"
)

// SafeDeref 安全地解引用字符串指针
// 用于安全地获取字符串指针的值，避免空指针异常
//
// 参数:
//   - s: 字符串指针
//
// 返回:
//   - string: 解引用后的字符串，指针为nil时返回空字符串
func SafeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BuildStr 使用字符串构建器执行函数
//
// 参数:
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 创建新的字符串构建器（不使用对象池）
//   - 执行用户提供的函数
//   - 返回构建的字符串结果
//   - 适用于不需要对象池优化的简单字符串构建场景
//
// 使用示例:
//
//	result := str.BuildStr(func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func BuildStr(fn func(*strings.Builder)) string {
	var buf strings.Builder
	fn(&buf)
	return buf.String()
}

// BuildStrCap 使用指定容量的字符串构建器执行函数
//
// 参数:
//   - cap: 字符串构建器初始容量
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 创建指定容量的新字符串构建器（不使用对象池）
//   - 执行用户提供的函数
//   - 返回构建的字符串结果
//   - 适用于已知字符串长度且不需要对象池优化的场景
//
// 使用示例:
//
//	result := str.BuildStrCap(64, func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func BuildStrCap(cap int, fn func(*strings.Builder)) string {
	var buf strings.Builder
	if cap > 0 {
		buf.Grow(cap)
	}
	fn(&buf)
	return buf.String()
}

// IsEmpty 检查字符串是否为空
//
// 参数:
//   - s: 待检查的字符串
//
// 返回:
//   - bool: 字符串为空返回true，否则返回false
func IsEmpty(s string) bool {
	return s == ""
}

// Prefix 获取字符串的前N个字符
//
// 参数:
//   - s: 输入字符串
//   - n: 要获取的字符数量
//
// 返回:
//   - string: 前N个字符，如果字符串长度不足N则返回原字符串
func Prefix(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// Suffix 获取字符串的后N个字符
//
// 参数:
//   - s: 输入字符串
//   - n: 要获取的字符数量
//
// 返回:
//   - string: 后N个字符，如果字符串长度不足N则返回原字符串
func Suffix(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

// StringSuffix8 从给定字符串中获取最后8个字符。
// 如果字符串长度小于等于8，则返回原字符串。
//
// 参数:
//   - s: 输入字符串
//
// 返回:
//   - string: 字符串的最后8个字符，或原字符串（如果长度不足8），或空字符串（如果输入为空）
func StringSuffix8(s string) string {

	// 检查输入字符串是否为空，若为空则直接返回空字符串
	if s == "" {
		return ""
	}

	// 检查输入字符串的长度是否小于等于 8，若是则直接返回该字符串本身
	if len(s) <= 8 {
		return s
	}

	// 若输入字符串长度大于 8，则截取并返回其最后 8 个字符
	return s[len(s)-8:]
}

// Truncate 截断字符串到指定长度
//
// 参数:
//   - s: 输入字符串
//   - maxLen: 最大长度
//
// 返回:
//   - string: 截断后的字符串
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// IsNotEmpty 检查字符串是否不为空
// 用于验证字符串在去除首尾空格后是否包含有效内容
//
// 参数:
//   - s: 待检查的字符串
//
// 返回:
//   - bool: 字符串不为空返回true，否则返回false
func IsNotEmpty(s string) bool {
	// 去除字符串两端的空格
	s = strings.TrimSpace(s)

	// 检查字符串是否为空
	return s != ""
}

// IfBlank 当字符串为空白时返回默认值
//
// 参数:
//   - s: 待检查的字符串
//   - defaultVal: 默认值
//
// 返回:
//   - string: 如果s为空白（空字符串或只包含空白字符）则返回defaultVal，否则返回s
func IfBlank(s, defaultVal string) string {
	if strings.TrimSpace(s) == "" {
		return defaultVal
	}
	return s
}

// IfEmpty 当字符串为空时返回默认值
//
// 参数:
//   - s: 待检查的字符串
//   - defaultVal: 默认值
//
// 返回:
//   - string: 如果s为空则返回defaultVal，否则返回s
func IfEmpty(s, defaultVal string) string {
	if s == "" {
		return defaultVal
	}
	return s
}

// Repeat 重复字符串N次
//
// 参数:
//   - s: 要重复的字符串
//   - count: 重复次数
//
// 返回:
//   - string: 重复后的字符串
func Repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

// PadLeft 在字符串左侧填充字符到指定长度
//
// 参数:
//   - s: 输入字符串
//   - length: 目标长度
//   - pad: 填充字符
//
// 返回:
//   - string: 填充后的字符串
func PadLeft(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	padCount := length - len(s)
	return strings.Repeat(string(pad), padCount) + s
}

// PadRight 在字符串右侧填充字符到指定长度
//
// 参数:
//   - s: 输入字符串
//   - length: 目标长度
//   - pad: 填充字符
//
// 返回:
//   - string: 填充后的字符串
func PadRight(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	padCount := length - len(s)
	return s + strings.Repeat(string(pad), padCount)
}

// SafeIndex 安全地查找子字符串的索引
//
// 参数:
//   - s: 源字符串
//   - substr: 要查找的子字符串
//
// 返回:
//   - int: 子字符串的索引，未找到返回-1
func SafeIndex(s, substr string) int {
	if s == "" || substr == "" {
		return -1
	}
	return strings.Index(s, substr)
}

// Ellipsis 超长字符串显示省略号
//
// 参数:
//   - s: 输入字符串
//   - maxLen: 最大长度（包含省略号）
//
// 返回:
//   - string: 处理后的字符串
func Ellipsis(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return strings.Repeat(".", maxLen)
	}
	return s[:maxLen-3] + "..."
}

// Template 简单模板替换
//
// 参数:
//   - tmpl: 模板字符串，使用 {{key}} 作为占位符
//   - data: 替换数据
//
// 返回:
//   - string: 替换后的字符串
//
// 使用示例:
//
//	result := str.Template("Hello {{name}}, you are {{age}} years old", map[string]string{
//	    "name": "Alice",
//	    "age":  "25",
//	})
func Template(tmpl string, data map[string]string) string {
	// 边界情况1: 模板字符串为空
	if tmpl == "" {
		return ""
	}

	// 边界情况2: 数据为空或nil
	if len(data) == 0 {
		return tmpl
	}

	// 预分配空间, 长度为: 数据对数 * 2
	pairs := make([]string, 0, len(data)*2)

	// 构建替换对: [占位符, 替换值, 占位符, 替换值, ...]
	for k, v := range data {
		// 边界情况3: 键为空字符串，跳过以避免无效占位符 {{}}
		if k == "" {
			continue
		}

		// 边界情况4: 键包含特殊字符，可能导致意外替换
		ph := "{{" + k + "}}"
		pairs = append(pairs, ph, v)
	}

	// 边界情况5: 所有键都为空，没有有效的替换对
	if len(pairs) == 0 {
		return tmpl
	}

	// 使用 strings.Replacer 一次性完成所有替换
	r := strings.NewReplacer(pairs...)
	return r.Replace(tmpl)
}

// ToBase64 将字符串编码为Base64
//
// 参数:
//   - s: 输入字符串
//
// 返回:
//   - string: Base64编码后的字符串
func ToBase64(s string) string {
	if s == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// FromBase64 将Base64字符串解码
//
// 参数:
//   - s: Base64编码的字符串
//
// 返回:
//   - string: 解码后的字符串
//   - error: 解码错误
func FromBase64(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Join 拼接多个字符串
//
// 参数:
//   - parts: 要拼接的字符串切片
//
// 返回:
//   - string: 拼接后的字符串
func Join(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}

	// 预计算总长度以减少内存分配
	totalLen := 0
	for _, part := range parts {
		totalLen += len(part)
	}

	// 复用 BuildStrCap 函数
	return BuildStrCap(totalLen, func(buf *strings.Builder) {
		for _, part := range parts {
			buf.WriteString(part)
		}
	})
}

// JoinNonEmpty 使用分隔符拼接非空字符串
//
// 参数:
//   - sep: 分隔符
//   - parts: 要拼接的字符串切片
//
// 返回:
//   - string: 用分隔符连接的非空字符串
func JoinNonEmpty(sep string, parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	// 单次循环：收集非空字符串并计算总长度
	nonEmpty := make([]string, 0, len(parts))
	totalLen := 0

	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
			totalLen += len(part)
		}
	}

	if len(nonEmpty) == 0 {
		return ""
	}
	if len(nonEmpty) == 1 {
		return nonEmpty[0]
	}

	// 加上分隔符的长度
	totalLen += len(sep) * (len(nonEmpty) - 1)

	// 复用 BuildStrCap 函数
	return BuildStrCap(totalLen, func(buf *strings.Builder) {
		for i, part := range nonEmpty {
			if i > 0 {
				buf.WriteString(sep)
			}
			buf.WriteString(part)
		}
	})
}

// Mask 字符串掩码处理 (如手机号脱敏)
//
// 参数:
//   - s: 输入字符串
//   - start: 开始掩码的位置（包含）
//   - end: 结束掩码的位置（不包含）
//   - maskChar: 掩码字符
//
// 返回:
//   - string: 掩码后的字符串
//
// 使用示例:
//
//	phone := "13812345678"
//	masked := str.Mask(phone, 3, 7, '*') // 138****5678
func Mask(s string, start, end int, maskChar rune) string {
	if s == "" || start < 0 || end < 0 || end <= start || start >= len(s) {
		return s
	}

	// 验证掩码字符，如果为零值或不可打印则使用默认字符
	if maskChar == 0 {
		maskChar = '*'
	}

	// 确保end不超过字符串长度
	if end > len(s) {
		end = len(s)
	}

	runes := []rune(s)
	for i := start; i < end && i < len(runes); i++ {
		runes[i] = maskChar
	}

	return string(runes)
}
