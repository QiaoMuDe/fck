package cmd

// getLast8Chars 函数用于获取输入字符串的最后8个字符
func getLast8Chars(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return s
	}
	return s[len(s)-8:]
}
