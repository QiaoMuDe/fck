package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/port"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var PortCmd *qflag.Cmd

var (
	portPorts       *qflag.IntSliceFlag // 指定端口
	portTCP         *qflag.BoolFlag     // 只显示 TCP
	portUDP         *qflag.BoolFlag     // 只显示 UDP
	portState       *qflag.StringFlag   // 连接状态过滤
	portPID         *qflag.IntFlag      // 指定 PID
	portProcessName *qflag.StringFlag   // 指定进程名
	portList        *qflag.BoolFlag     // 简洁模式
	portTableStyle  *qflag.EnumFlag     // 表格样式
	portListening   *qflag.BoolFlag     // 只显示监听端口
)

func init() {
	PortCmd = qflag.NewCmd("port", "", qflag.ExitOnError)

	portPorts = PortCmd.IntSlice("port", "P", "指定要查看的端口，多个端口用逗号分隔，如: 80,443,8080", []int{})
	portTCP = PortCmd.Bool("tcp", "t", "只显示 TCP 连接", false)
	portUDP = PortCmd.Bool("udp", "u", "只显示 UDP 连接", false)
	portState = PortCmd.String("state", "s", "按连接状态过滤，如: LISTEN, ESTABLISHED, TIME_WAIT", "")
	portPID = PortCmd.Int("pid", "", "指定进程 ID 查看该进程的端口", 0)
	portProcessName = PortCmd.String("name", "n", "按进程名过滤，支持部分匹配", "")
	portList = PortCmd.Bool("list", "l", "简洁模式，只显示地址和进程信息", false)
	portTableStyle = PortCmd.Enum("table-style", "ts", "指定表格样式，支持以下选项：\n"+
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
		"\t\t\t\t\t[none]   - 禁用边框样式", "none", types.TableStyles)
	portListening = PortCmd.Bool("listening", "L", "只显示监听状态的端口", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:       "查看系统端口占用情况",
		UseChinese: true,
		Notes: []string{
			"默认显示所有 TCP 和 UDP 端口",
			"查看所有进程端口通常需要管理员/root权限",
			"Windows 上某些系统进程信息可能无法获取",
		},
		Examples: map[string]string{
			"查看所有端口":       "fck port",
			"查看指定端口":       "fck port -P 8080",
			"查看多个端口":       "fck port -P 80,443,8080",
			"只查看 TCP 监听端口": "fck port -t -L",
			"只查看 UDP 端口":   "fck port -u",
			"查看指定进程的端口":    "fck port -n nginx",
			"查看指定 PID 的端口": "fck port --pid 1234",
			"简洁模式显示":       "fck port -l",
		},
		MutexGroups: []qflag.MutexGroup{
			{Name: "protocol", Flags: []string{"tcp", "udp"}, AllowNone: true},
		},
	}

	if err := PortCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	PortCmd.SetRun(runPort)
}

func runPort(cmd qflag.Command) error {
	// 获取端口列表并验证范围
	ports := portPorts.Get()
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("port number out of range (1-65535): %d", p)
		}
	}

	config := &port.PortConfig{
		Ports:       ports,
		TCPOnly:     portTCP.Get(),
		UDPOnly:     portUDP.Get(),
		State:       portState.Get(),
		PID:         int32(portPID.Get()),
		ProcessName: portProcessName.Get(),
		ListMode:    portList.Get(),
		TableStyle:  portTableStyle.Get(),
		Listening:   portListening.Get(),
	}

	return port.PortCmdMain(config)
}
