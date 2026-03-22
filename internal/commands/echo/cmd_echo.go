package echo

import (
	"fmt"
	"strings"
)

// EchoConfig 配置结构体
type EchoConfig struct {
	Text      string
	NoNewline bool
	Raw       bool
	Trim      bool
	Upper     bool
	Lower     bool
	Repeat    int
	Color     string
}

// ANSI颜色代码
var colorCodes = map[string]string{
	"black":          "\033[30m",
	"red":            "\033[31m",
	"green":          "\033[32m",
	"yellow":         "\033[33m",
	"blue":           "\033[34m",
	"magenta":        "\033[35m",
	"cyan":           "\033[36m",
	"white":          "\033[37m",
	"bright-black":   "\033[90m",
	"bright-red":     "\033[91m",
	"bright-green":   "\033[92m",
	"bright-yellow":  "\033[93m",
	"bright-blue":    "\033[94m",
	"bright-magenta": "\033[95m",
	"bright-cyan":    "\033[96m",
	"bright-white":   "\033[97m",
	"reset":          "\033[0m",
}

// EchoCmdMain 主函数
func EchoCmdMain(config EchoConfig) error {
	if config.Repeat <= 0 {
		return nil
	}

	text := config.Text

	text = processText(text, config)

	for i := 0; i < config.Repeat; i++ {
		outputText := text

		if config.Color != "" {
			if colorCode, ok := colorCodes[config.Color]; ok {
				outputText = colorCode + outputText + colorCodes["reset"]
			}
		}

		if config.NoNewline {
			fmt.Print(outputText)
		} else {
			fmt.Println(outputText)
		}
	}

	return nil
}

// processText 处理文本
func processText(text string, config EchoConfig) string {
	result := text

	if !config.Raw {
		result = unescapeString(result)
	}

	if config.Trim {
		result = strings.TrimSpace(result)
	}

	if config.Upper {
		result = strings.ToUpper(result)
	} else if config.Lower {
		result = strings.ToLower(result)
	}

	return result
}

// unescapeString 解析转义字符
func unescapeString(s string) string {
	result := strings.Builder{}
	escaped := false

	for i := 0; i < len(s); i++ {
		c := s[i]

		if escaped {
			switch c {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			case '\'':
				result.WriteByte('\'')
			default:
				result.WriteByte('\\')
				result.WriteByte(c)
			}
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		result.WriteByte(c)
	}

	if escaped {
		result.WriteByte('\\')
	}

	return result.String()
}
