package color

// 代码由 gen/helper_gen.go 自动生成，请勿手动修改
// 生成命令: go generate . || go run gen/helper_gen.go

// Black 以黑色前景打印文本
func Black(message string) { colorPrint(FgBlack, "%s\n", message) }

// Blackf 以黑色前景打印文本
func Blackf(format string, a ...interface{}) { colorPrint(FgBlack, format, a...) }

// SBlack 返回黑色前景字符串
func SBlack(message string) string { return colorString(FgBlack, "%s", message) }

// SBlackf 返回黑色前景字符串
func SBlackf(format string, a ...interface{}) string { return colorString(FgBlack, format, a...) }

// Red 以红色前景打印文本
func Red(message string) { colorPrint(FgRed, "%s\n", message) }

// Redf 以红色前景打印文本
func Redf(format string, a ...interface{}) { colorPrint(FgRed, format, a...) }

// SRed 返回红色前景字符串
func SRed(message string) string { return colorString(FgRed, "%s", message) }

// SRedf 返回红色前景字符串
func SRedf(format string, a ...interface{}) string { return colorString(FgRed, format, a...) }

// Green 以绿色前景打印文本
func Green(message string) { colorPrint(FgGreen, "%s\n", message) }

// Greenf 以绿色前景打印文本
func Greenf(format string, a ...interface{}) { colorPrint(FgGreen, format, a...) }

// SGreen 返回绿色前景字符串
func SGreen(message string) string { return colorString(FgGreen, "%s", message) }

// SGreenf 返回绿色前景字符串
func SGreenf(format string, a ...interface{}) string { return colorString(FgGreen, format, a...) }

// Yellow 以黄色前景打印文本
func Yellow(message string) { colorPrint(FgYellow, "%s\n", message) }

// Yellowf 以黄色前景打印文本
func Yellowf(format string, a ...interface{}) { colorPrint(FgYellow, format, a...) }

// SYellow 返回黄色前景字符串
func SYellow(message string) string { return colorString(FgYellow, "%s", message) }

// SYellowf 返回黄色前景字符串
func SYellowf(format string, a ...interface{}) string { return colorString(FgYellow, format, a...) }

// Blue 以蓝色前景打印文本
func Blue(message string) { colorPrint(FgBlue, "%s\n", message) }

// Bluef 以蓝色前景打印文本
func Bluef(format string, a ...interface{}) { colorPrint(FgBlue, format, a...) }

// SBlue 返回蓝色前景字符串
func SBlue(message string) string { return colorString(FgBlue, "%s", message) }

// SBluef 返回蓝色前景字符串
func SBluef(format string, a ...interface{}) string { return colorString(FgBlue, format, a...) }

// Magenta 以洋红色前景打印文本
func Magenta(message string) { colorPrint(FgMagenta, "%s\n", message) }

// Magentaf 以洋红色前景打印文本
func Magentaf(format string, a ...interface{}) { colorPrint(FgMagenta, format, a...) }

// SMagenta 返回洋红色前景字符串
func SMagenta(message string) string { return colorString(FgMagenta, "%s", message) }

// SMagentaf 返回洋红色前景字符串
func SMagentaf(format string, a ...interface{}) string { return colorString(FgMagenta, format, a...) }

// Cyan 以青色前景打印文本
func Cyan(message string) { colorPrint(FgCyan, "%s\n", message) }

// Cyanf 以青色前景打印文本
func Cyanf(format string, a ...interface{}) { colorPrint(FgCyan, format, a...) }

// SCyan 返回青色前景字符串
func SCyan(message string) string { return colorString(FgCyan, "%s", message) }

// SCyanf 返回青色前景字符串
func SCyanf(format string, a ...interface{}) string { return colorString(FgCyan, format, a...) }

// White 以白色前景打印文本
func White(message string) { colorPrint(FgWhite, "%s\n", message) }

// Whitef 以白色前景打印文本
func Whitef(format string, a ...interface{}) { colorPrint(FgWhite, format, a...) }

// SWhite 返回白色前景字符串
func SWhite(message string) string { return colorString(FgWhite, "%s", message) }

// SWhitef 返回白色前景字符串
func SWhitef(format string, a ...interface{}) string { return colorString(FgWhite, format, a...) }

// Gray 以灰色前景打印文本
func Gray(message string) { colorPrint(FgHiBlack, "%s\n", message) }

// Grayf 以灰色前景打印文本
func Grayf(format string, a ...interface{}) { colorPrint(FgHiBlack, format, a...) }

// SGray 返回灰色前景字符串
func SGray(message string) string { return colorString(FgHiBlack, "%s", message) }

// SGrayf 返回灰色前景字符串
func SGrayf(format string, a ...interface{}) string { return colorString(FgHiBlack, format, a...) }

// HiBlack 以高亮黑色前景打印文本
func HiBlack(message string) { colorPrint(FgHiBlack, "%s\n", message) }

// HiBlackf 以高亮黑色前景打印文本
func HiBlackf(format string, a ...interface{}) { colorPrint(FgHiBlack, format, a...) }

// SHiBlack 返回高亮黑色前景字符串
func SHiBlack(message string) string { return colorString(FgHiBlack, "%s", message) }

// SHiBlackf 返回高亮黑色前景字符串
func SHiBlackf(format string, a ...interface{}) string { return colorString(FgHiBlack, format, a...) }

// HiRed 以高亮红色前景打印文本
func HiRed(message string) { colorPrint(FgHiRed, "%s\n", message) }

// HiRedf 以高亮红色前景打印文本
func HiRedf(format string, a ...interface{}) { colorPrint(FgHiRed, format, a...) }

// SHiRed 返回高亮红色前景字符串
func SHiRed(message string) string { return colorString(FgHiRed, "%s", message) }

// SHiRedf 返回高亮红色前景字符串
func SHiRedf(format string, a ...interface{}) string { return colorString(FgHiRed, format, a...) }

// HiGreen 以高亮绿色前景打印文本
func HiGreen(message string) { colorPrint(FgHiGreen, "%s\n", message) }

// HiGreenf 以高亮绿色前景打印文本
func HiGreenf(format string, a ...interface{}) { colorPrint(FgHiGreen, format, a...) }

// SHiGreen 返回高亮绿色前景字符串
func SHiGreen(message string) string { return colorString(FgHiGreen, "%s", message) }

// SHiGreenf 返回高亮绿色前景字符串
func SHiGreenf(format string, a ...interface{}) string { return colorString(FgHiGreen, format, a...) }

// HiYellow 以高亮黄色前景打印文本
func HiYellow(message string) { colorPrint(FgHiYellow, "%s\n", message) }

// HiYellowf 以高亮黄色前景打印文本
func HiYellowf(format string, a ...interface{}) { colorPrint(FgHiYellow, format, a...) }

// SHiYellow 返回高亮黄色前景字符串
func SHiYellow(message string) string { return colorString(FgHiYellow, "%s", message) }

// SHiYellowf 返回高亮黄色前景字符串
func SHiYellowf(format string, a ...interface{}) string { return colorString(FgHiYellow, format, a...) }

// HiBlue 以高亮蓝色前景打印文本
func HiBlue(message string) { colorPrint(FgHiBlue, "%s\n", message) }

// HiBluef 以高亮蓝色前景打印文本
func HiBluef(format string, a ...interface{}) { colorPrint(FgHiBlue, format, a...) }

// SHiBlue 返回高亮蓝色前景字符串
func SHiBlue(message string) string { return colorString(FgHiBlue, "%s", message) }

// SHiBluef 返回高亮蓝色前景字符串
func SHiBluef(format string, a ...interface{}) string { return colorString(FgHiBlue, format, a...) }

// HiMagenta 以高亮洋红色前景打印文本
func HiMagenta(message string) { colorPrint(FgHiMagenta, "%s\n", message) }

// HiMagentaf 以高亮洋红色前景打印文本
func HiMagentaf(format string, a ...interface{}) { colorPrint(FgHiMagenta, format, a...) }

// SHiMagenta 返回高亮洋红色前景字符串
func SHiMagenta(message string) string { return colorString(FgHiMagenta, "%s", message) }

// SHiMagentaf 返回高亮洋红色前景字符串
func SHiMagentaf(format string, a ...interface{}) string {
	return colorString(FgHiMagenta, format, a...)
}

// HiCyan 以高亮青色前景打印文本
func HiCyan(message string) { colorPrint(FgHiCyan, "%s\n", message) }

// HiCyanf 以高亮青色前景打印文本
func HiCyanf(format string, a ...interface{}) { colorPrint(FgHiCyan, format, a...) }

// SHiCyan 返回高亮青色前景字符串
func SHiCyan(message string) string { return colorString(FgHiCyan, "%s", message) }

// SHiCyanf 返回高亮青色前景字符串
func SHiCyanf(format string, a ...interface{}) string { return colorString(FgHiCyan, format, a...) }

// HiWhite 以高亮白色前景打印文本
func HiWhite(message string) { colorPrint(FgHiWhite, "%s\n", message) }

// HiWhitef 以高亮白色前景打印文本
func HiWhitef(format string, a ...interface{}) { colorPrint(FgHiWhite, format, a...) }

// SHiWhite 返回高亮白色前景字符串
func SHiWhite(message string) string { return colorString(FgHiWhite, "%s", message) }

// SHiWhitef 返回高亮白色前景字符串
func SHiWhitef(format string, a ...interface{}) string { return colorString(FgHiWhite, format, a...) }
