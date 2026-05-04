package color

// 代码由 gen/global_gen.go 自动生成，请勿手动修改
// 生成命令: go generate . || go run gen/global_gen.go

// Black 使用黑色样式打印
func (g *GlobalColor) Black(message string) { g.printColor(FgBlack, "%s\n", message) }

// Blackf 使用黑色样式打印
func (g *GlobalColor) Blackf(format string, a ...interface{}) { g.printColor(FgBlack, format, a...) }

// SBlack 返回黑色样式的字符串
func (g *GlobalColor) SBlack(message string) string { return g.sprintColor(FgBlack, "%s", message) }

// SBlackf 返回黑色样式的字符串
func (g *GlobalColor) SBlackf(format string, a ...interface{}) string {
	return g.sprintColor(FgBlack, format, a...)
}

// Red 使用红色样式打印
func (g *GlobalColor) Red(message string) { g.printColor(FgRed, "%s\n", message) }

// Redf 使用红色样式打印
func (g *GlobalColor) Redf(format string, a ...interface{}) { g.printColor(FgRed, format, a...) }

// SRed 返回红色样式的字符串
func (g *GlobalColor) SRed(message string) string { return g.sprintColor(FgRed, "%s", message) }

// SRedf 返回红色样式的字符串
func (g *GlobalColor) SRedf(format string, a ...interface{}) string {
	return g.sprintColor(FgRed, format, a...)
}

// Green 使用绿色样式打印
func (g *GlobalColor) Green(message string) { g.printColor(FgGreen, "%s\n", message) }

// Greenf 使用绿色样式打印
func (g *GlobalColor) Greenf(format string, a ...interface{}) { g.printColor(FgGreen, format, a...) }

// SGreen 返回绿色样式的字符串
func (g *GlobalColor) SGreen(message string) string { return g.sprintColor(FgGreen, "%s", message) }

// SGreenf 返回绿色样式的字符串
func (g *GlobalColor) SGreenf(format string, a ...interface{}) string {
	return g.sprintColor(FgGreen, format, a...)
}

// Yellow 使用黄色样式打印
func (g *GlobalColor) Yellow(message string) { g.printColor(FgYellow, "%s\n", message) }

// Yellowf 使用黄色样式打印
func (g *GlobalColor) Yellowf(format string, a ...interface{}) { g.printColor(FgYellow, format, a...) }

// SYellow 返回黄色样式的字符串
func (g *GlobalColor) SYellow(message string) string { return g.sprintColor(FgYellow, "%s", message) }

// SYellowf 返回黄色样式的字符串
func (g *GlobalColor) SYellowf(format string, a ...interface{}) string {
	return g.sprintColor(FgYellow, format, a...)
}

// Blue 使用蓝色样式打印
func (g *GlobalColor) Blue(message string) { g.printColor(FgBlue, "%s\n", message) }

// Bluef 使用蓝色样式打印
func (g *GlobalColor) Bluef(format string, a ...interface{}) { g.printColor(FgBlue, format, a...) }

// SBlue 返回蓝色样式的字符串
func (g *GlobalColor) SBlue(message string) string { return g.sprintColor(FgBlue, "%s", message) }

// SBluef 返回蓝色样式的字符串
func (g *GlobalColor) SBluef(format string, a ...interface{}) string {
	return g.sprintColor(FgBlue, format, a...)
}

// Magenta 使用洋红色样式打印
func (g *GlobalColor) Magenta(message string) { g.printColor(FgMagenta, "%s\n", message) }

// Magentaf 使用洋红色样式打印
func (g *GlobalColor) Magentaf(format string, a ...interface{}) {
	g.printColor(FgMagenta, format, a...)
}

// SMagenta 返回洋红色样式的字符串
func (g *GlobalColor) SMagenta(message string) string { return g.sprintColor(FgMagenta, "%s", message) }

// SMagentaf 返回洋红色样式的字符串
func (g *GlobalColor) SMagentaf(format string, a ...interface{}) string {
	return g.sprintColor(FgMagenta, format, a...)
}

// Cyan 使用青色样式打印
func (g *GlobalColor) Cyan(message string) { g.printColor(FgCyan, "%s\n", message) }

// Cyanf 使用青色样式打印
func (g *GlobalColor) Cyanf(format string, a ...interface{}) { g.printColor(FgCyan, format, a...) }

// SCyan 返回青色样式的字符串
func (g *GlobalColor) SCyan(message string) string { return g.sprintColor(FgCyan, "%s", message) }

// SCyanf 返回青色样式的字符串
func (g *GlobalColor) SCyanf(format string, a ...interface{}) string {
	return g.sprintColor(FgCyan, format, a...)
}

// White 使用白色样式打印
func (g *GlobalColor) White(message string) { g.printColor(FgWhite, "%s\n", message) }

// Whitef 使用白色样式打印
func (g *GlobalColor) Whitef(format string, a ...interface{}) { g.printColor(FgWhite, format, a...) }

// SWhite 返回白色样式的字符串
func (g *GlobalColor) SWhite(message string) string { return g.sprintColor(FgWhite, "%s", message) }

// SWhitef 返回白色样式的字符串
func (g *GlobalColor) SWhitef(format string, a ...interface{}) string {
	return g.sprintColor(FgWhite, format, a...)
}

// Gray 使用灰色样式打印
func (g *GlobalColor) Gray(message string) { g.printColor(FgHiBlack, "%s\n", message) }

// Grayf 使用灰色样式打印
func (g *GlobalColor) Grayf(format string, a ...interface{}) { g.printColor(FgHiBlack, format, a...) }

// SGray 返回灰色样式的字符串
func (g *GlobalColor) SGray(message string) string { return g.sprintColor(FgHiBlack, "%s", message) }

// SGrayf 返回灰色样式的字符串
func (g *GlobalColor) SGrayf(format string, a ...interface{}) string {
	return g.sprintColor(FgHiBlack, format, a...)
}

// HiBlack 使用高亮黑色样式打印
func (g *GlobalColor) HiBlack(message string) { g.printColor(FgHiBlack, "%s\n", message) }

// HiBlackf 使用高亮黑色样式打印
func (g *GlobalColor) HiBlackf(format string, a ...interface{}) {
	g.printColor(FgHiBlack, format, a...)
}

// SHiBlack 返回高亮黑色样式的字符串
func (g *GlobalColor) SHiBlack(message string) string { return g.sprintColor(FgHiBlack, "%s", message) }

// SHiBlackf 返回高亮黑色样式的字符串
func (g *GlobalColor) SHiBlackf(format string, a ...interface{}) string {
	return g.sprintColor(FgHiBlack, format, a...)
}

// HiRed 使用高亮红色样式打印
func (g *GlobalColor) HiRed(message string) { g.printColor(FgHiRed, "%s\n", message) }

// HiRedf 使用高亮红色样式打印
func (g *GlobalColor) HiRedf(format string, a ...interface{}) { g.printColor(FgHiRed, format, a...) }

// SHiRed 返回高亮红色样式的字符串
func (g *GlobalColor) SHiRed(message string) string { return g.sprintColor(FgHiRed, "%s", message) }

// SHiRedf 返回高亮红色样式的字符串
func (g *GlobalColor) SHiRedf(format string, a ...interface{}) string {
	return g.sprintColor(FgHiRed, format, a...)
}

// HiGreen 使用高亮绿色样式打印
func (g *GlobalColor) HiGreen(message string) { g.printColor(FgHiGreen, "%s\n", message) }

// HiGreenf 使用高亮绿色样式打印
func (g *GlobalColor) HiGreenf(format string, a ...interface{}) {
	g.printColor(FgHiGreen, format, a...)
}

// SHiGreen 返回高亮绿色样式的字符串
func (g *GlobalColor) SHiGreen(message string) string { return g.sprintColor(FgHiGreen, "%s", message) }

// SHiGreenf 返回高亮绿色样式的字符串
func (g *GlobalColor) SHiGreenf(format string, a ...interface{}) string {
	return g.sprintColor(FgHiGreen, format, a...)
}

// HiYellow 使用高亮黄色样式打印
func (g *GlobalColor) HiYellow(message string) { g.printColor(FgHiYellow, "%s\n", message) }

// HiYellowf 使用高亮黄色样式打印
func (g *GlobalColor) HiYellowf(format string, a ...interface{}) {
	g.printColor(FgHiYellow, format, a...)
}

// SHiYellow 返回高亮黄色样式的字符串
func (g *GlobalColor) SHiYellow(message string) string {
	return g.sprintColor(FgHiYellow, "%s", message)
}

// SHiYellowf 返回高亮黄色样式的字符串
func (g *GlobalColor) SHiYellowf(format string, a ...interface{}) string {
	return g.sprintColor(FgHiYellow, format, a...)
}

// HiBlue 使用高亮蓝色样式打印
func (g *GlobalColor) HiBlue(message string) { g.printColor(FgHiBlue, "%s\n", message) }

// HiBluef 使用高亮蓝色样式打印
func (g *GlobalColor) HiBluef(format string, a ...interface{}) { g.printColor(FgHiBlue, format, a...) }

// SHiBlue 返回高亮蓝色样式的字符串
func (g *GlobalColor) SHiBlue(message string) string { return g.sprintColor(FgHiBlue, "%s", message) }

// SHiBluef 返回高亮蓝色样式的字符串
func (g *GlobalColor) SHiBluef(format string, a ...interface{}) string {
	return g.sprintColor(FgHiBlue, format, a...)
}

// HiMagenta 使用高亮洋红色样式打印
func (g *GlobalColor) HiMagenta(message string) { g.printColor(FgHiMagenta, "%s\n", message) }

// HiMagentaf 使用高亮洋红色样式打印
func (g *GlobalColor) HiMagentaf(format string, a ...interface{}) {
	g.printColor(FgHiMagenta, format, a...)
}

// SHiMagenta 返回高亮洋红色样式的字符串
func (g *GlobalColor) SHiMagenta(message string) string {
	return g.sprintColor(FgHiMagenta, "%s", message)
}

// SHiMagentaf 返回高亮洋红色样式的字符串
func (g *GlobalColor) SHiMagentaf(format string, a ...interface{}) string {
	return g.sprintColor(FgHiMagenta, format, a...)
}

// HiCyan 使用高亮青色样式打印
func (g *GlobalColor) HiCyan(message string) { g.printColor(FgHiCyan, "%s\n", message) }

// HiCyanf 使用高亮青色样式打印
func (g *GlobalColor) HiCyanf(format string, a ...interface{}) { g.printColor(FgHiCyan, format, a...) }

// SHiCyan 返回高亮青色样式的字符串
func (g *GlobalColor) SHiCyan(message string) string { return g.sprintColor(FgHiCyan, "%s", message) }

// SHiCyanf 返回高亮青色样式的字符串
func (g *GlobalColor) SHiCyanf(format string, a ...interface{}) string {
	return g.sprintColor(FgHiCyan, format, a...)
}

// HiWhite 使用高亮白色样式打印
func (g *GlobalColor) HiWhite(message string) { g.printColor(FgHiWhite, "%s\n", message) }

// HiWhitef 使用高亮白色样式打印
func (g *GlobalColor) HiWhitef(format string, a ...interface{}) {
	g.printColor(FgHiWhite, format, a...)
}

// SHiWhite 返回高亮白色样式的字符串
func (g *GlobalColor) SHiWhite(message string) string { return g.sprintColor(FgHiWhite, "%s", message) }

// SHiWhitef 返回高亮白色样式的字符串
func (g *GlobalColor) SHiWhitef(format string, a ...interface{}) string {
	return g.sprintColor(FgHiWhite, format, a...)
}
