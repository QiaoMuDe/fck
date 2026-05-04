/*
Package color 是一个 ANSI 颜色包，用于向标准输出输出彩色或 SGR 定义的文本。
API 可以通过多种方式使用，选择适合你的方式即可。

使用简单且默认的辅助函数，配合预定义的前景色：

	// 直接打印字符串（自动追加换行符）
	color.Cyan("以青色打印文本。")
	color.Red("我们有红色")
	color.Yellow("也有黄色！")

	// 使用格式化版本（不自动追加换行符，需手动添加）
	color.Bluef("以蓝色打印 %s。\n", "文本")

	// 高亮色
	color.HiGreen("亮绿色。")
	color.HiBlack("亮黑色就是灰色..")
	color.HiWhite("闪亮的白色！")

然而，有时需要自定义颜色组合。以下是一些创建自定义颜色对象
并使用每个独立颜色对象的打印函数的示例。

	// 创建一个新的颜色对象
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("打印带下划线的青色文本。")

	// 或者直接添加到 New() 中
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("这也打印粗体青色 %s\n", "！")


	// 混合前景色和背景色，创建新的组合！
	red := color.New(color.FgRed)

	boldRed := red.Add(color.Bold)
	boldRed.Println("这将打印粗体红色文本。")

	whiteBackground := red.Add(color.BgWhite)
	whiteBackground.Println("红色文本配白色背景。")

	// 使用你自己的 io.Writer 输出
	color.New(color.FgBlue).Fprintln(myWriter, "蓝色！")

	blue := color.New(color.FgBlue)
	blue.Fprint(myWriter, "这将打印蓝色文本。")

你可以创建 PrintXxx 函数来进一步简化：

	// 创建自定义打印函数以方便使用
	red := color.New(color.FgRed).PrintfFunc()
	red("警告")
	red("错误：%s", err)

	// 混合多个属性
	notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
	notice("不要忘记这个...")

你也可以使用 FprintXxx 函数传入你自己的 io.Writer：

	blue := color.New(FgBlue).FprintfFunc()
	blue(myWriter, "重要通知：%s", stars)

	// 混合多个属性
	success := color.New(color.Bold, color.FgGreen).FprintlnFunc()
	success(myWriter, "不要忘记这个...")

或者创建 SprintXxx 函数来将字符串与其他非彩色字符串混合：

	yellow := New(FgYellow).SprintFunc()
	red := New(FgRed).SprintFunc()

	fmt.Printf("这是一个 %s，这是一个 %s。\n", yellow("警告"), red("错误"))

	info := New(FgWhite, BgGreen).SprintFunc()
	fmt.Printf("这个 %s 太棒了！\n", info("包"))

Windows 支持默认启用。所有 Print 函数都能按预期工作。
但是，仅对于 color.SprintXXX 函数，用户应该使用 fmt.FprintXXX
并将输出设置为 color.Output：

	fmt.Fprintf(color.Output, "Windows 支持：%s", color.SGreen("通过"))

	info := New(FgWhite, BgGreen).SprintFunc()
	fmt.Fprintf(color.Output, "这个 %s 太棒了！\n", info("包"))

可以与现有代码一起使用。只需使用 Set() 方法将标准输出设置为给定参数。
这样就不需要重写现有代码。

	// 使用便捷的标准颜色。
	color.Set(color.FgYellow)

	fmt.Println("现有文本现在将显示为黄色")
	fmt.Printf("这个也是 %s\n", "黄色")

	color.Unset() // 不要忘记取消设置

	// 你可以混合参数
	color.Set(color.FgMagenta, color.Bold)
	defer color.Unset() // 在你的函数中使用它

	fmt.Println("所有文本现在将显示为粗体洋红色。")

可能会有需要禁用颜色输出的情况（例如将应用程序的标准输出管道传输到其他地方）。
`Color` 支持全局和单个颜色定义禁用颜色。例如，假设你有一个 CLI 应用程序
和一个 `--no-color` 布尔标志。你可以轻松禁用颜色输出：

	var flagNoColor = flag.Bool("no-color", false, "禁用颜色输出")

	if *flagNoColor {
		color.NoColor = true // 禁用彩色输出
	}

你也可以通过将 NO_COLOR 环境变量设置为任何值来禁用颜色。

它还支持单个颜色定义（本地）。你可以随时禁用/启用颜色输出：

	c := color.New(color.FgCyan)
	c.Println("打印青色文本")

	c.DisableColor()
	c.Println("这将不打印任何颜色")

	c.EnableColor()
	c.Println("这又打印青色了...")
*/
package color
