package cmd

// getLast8Chars 函数用于获取输入字符串的最后 8 个字符。
// 如果输入字符串为空，则返回空字符串；
// 如果输入字符串的长度小于等于 8，则返回该字符串本身。
func getLast8Chars(s string) string {
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
