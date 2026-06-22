package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/ifconfig"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var IfconfigCmd *qflag.Cmd

var (
	ifconfigAll        *qflag.BoolFlag // -a, --all 显示所有接口, 包含虚拟网卡
	ifconfigShort      *qflag.BoolFlag // -s, --short 简略输出
	ifconfigJSON       *qflag.BoolFlag // -j, --json JSON 格式输出
	ifconfigStats      *qflag.BoolFlag // --stats 显示接口流量统计信息
	ifconfigTableStyle *qflag.EnumFlag // -ts, --table-style 指定表格样式
)

func init() {
	IfconfigCmd = qflag.NewCmd("ifconfig", "", qflag.ExitOnError)

	// 注册命令行标志
	ifconfigAll = IfconfigCmd.Bool("all", "a", "显示所有接口, 包含虚拟网卡", false)
	ifconfigShort = IfconfigCmd.Bool("short", "s", "简略输出 (名称、状态、IP)", false)
	ifconfigJSON = IfconfigCmd.Bool("json", "j", "JSON 格式输出", false)
	ifconfigStats = IfconfigCmd.Bool("stats", "", "显示接口流量统计信息", false)
	ifconfigTableStyle = IfconfigCmd.Enum("table-style", "ts", "指定表格样式, 支持以下选项：\n"+
		"\t\t\t\t\t[def ]   - 默认样式\n"+
		"\t\t\t\t\t[l   ]   - 浅色样式\n"+
		"\t\t\t\t\t[r   ]   - 圆角样式\n"+
		"\t\t\t\t\t[bd  ]   - 粗体样式\n"+
		"\t\t\t\t\t[cb  ]   - 亮色彩色样式\n"+
		"\t\t\t\t\t[cd  ]   - 暗色彩色样式\n"+
		"\t\t\t\t\t[db  ]   - 双线样式\n"+
		"\t\t\t\t\t[cbb ]   - 黑色背景蓝色字体\n"+
		"\t\t\t\t\t[cbc ]   - 青色背景蓝色字体\n"+
		"\t\t\t\t\t[cbg ]   - 绿色背景蓝色字体\n"+
		"\t\t\t\t\t[cbm ]   - 紫色背景蓝色字体\n"+
		"\t\t\t\t\t[cby ]   - 黄色背景蓝色字体\n"+
		"\t\t\t\t\t[cbr ]   - 红色背景蓝色字体\n"+
		"\t\t\t\t\t[cwb ]   - 蓝色背景白色字体\n"+
		"\t\t\t\t\t[ccw ]   - 青色背景白色字体\n"+
		"\t\t\t\t\t[cgw ]   - 绿色背景白色字体\n"+
		"\t\t\t\t\t[cmw ]   - 紫色背景白色字体\n"+
		"\t\t\t\t\t[crw ]   - 红色背景白色字体\n"+
		"\t\t\t\t\t[cyw ]   - 黄色背景白色字体\n"+
		"\t\t\t\t\t[none]   - 禁用表格样式", "none", types.TableStyles)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "网络接口信息查看工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s ifconfig [options] [interface...]", qflag.Root.Name()),
		Notes: []string{
			"支持通过标准输入或参数指定接口名称；不指定时显示所有有效接口",
			"默认过滤虚拟网卡, 使用 -a/--all 显示全部接口",
			"JSON 输出与表格样式互斥, JSON 模式优先于表格样式",
			"--stats 和 -s/--short 互斥, 不能同时使用",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "output-mode",
				Flags:     []string{"stats", "short"},
				AllowNone: true,
			},
		},
		Examples: map[string]string{
			"查看本机接口":       fmt.Sprintf("%s ifconfig", qflag.Root.Name()),
			"显示全部接口 (含虚拟)": fmt.Sprintf("%s ifconfig -a", qflag.Root.Name()),
			"简略输出":         fmt.Sprintf("%s ifconfig -s", qflag.Root.Name()),
			"JSON 输出":      fmt.Sprintf("%s ifconfig -j", qflag.Root.Name()),
			"查看指定接口":       fmt.Sprintf("%s ifconfig eth0", qflag.Root.Name()),
			"显示流量统计":       fmt.Sprintf("%s ifconfig --stats", qflag.Root.Name()),
		},
	}

	if err := IfconfigCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	IfconfigCmd.SetRun(runIfconfig)
}

// runIfconfig 是 ifconfig 命令的运行函数, 将 CLI 参数转换为配置并执行业务逻辑
func runIfconfig(cmd qflag.Command) error {
	config := ifconfig.IfconfigConfig{
		All:        ifconfigAll.Get(),
		Short:      ifconfigShort.Get(),
		JSON:       ifconfigJSON.Get(),
		Stats:      ifconfigStats.Get(),
		TableStyle: ifconfigTableStyle.Get(),
		Interfaces: cmd.Args(),
	}

	return ifconfig.IfconfigCmdMain(config)
}
