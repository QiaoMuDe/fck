package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/proc"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var ProcCmd *qflag.Cmd

var (
	procName       *qflag.StringFlag   // 按进程名过滤
	procPID        *qflag.IntFlag      // 单个 PID
	procPIDs       *qflag.IntSliceFlag // 多个 PID
	procSort       *qflag.EnumFlag     // 排序字段
	procAscend     *qflag.BoolFlag     // 升序排列
	procList       *qflag.BoolFlag     // 简洁模式
	procTree       *qflag.BoolFlag     // 树形模式
	procJSON       *qflag.BoolFlag     // JSON 输出
	procTableStyle *qflag.EnumFlag     // 表格样式
)

// SortFields 排序字段选项
var SortFields = []string{"pid", "name", "cpu", "mem", "time"}

func init() {
	ProcCmd = qflag.NewCmd("proc", "", qflag.ExitOnError)

	procName = ProcCmd.String("name", "n", "按进程名过滤，支持部分匹配", "")
	procPID = ProcCmd.Int("pid", "p", "指定单个 PID", 0)
	procPIDs = ProcCmd.IntSlice("pids", "P", "指定多个 PID, 如: 1234,5678", []int{})
	procSort = ProcCmd.Enum("sort", "s", "指定排序字段 (pid/name/cpu/mem/time)", "pid", SortFields)
	procAscend = ProcCmd.Bool("asc", "", "升序排列（默认降序）", false)
	procList = ProcCmd.Bool("list", "l", "简洁模式", false)
	procTree = ProcCmd.Bool("tree", "", "树形显示进程关系", false)
	procJSON = ProcCmd.Bool("json", "", "JSON 格式输出", false)
	procTableStyle = ProcCmd.Enum("table-style", "ts", "指定表格样式，支持以下选项：\n"+
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

	cmdOpts := &qflag.CmdOpts{
		Desc:       "查看系统进程信息",
		UseChinese: true,
		Notes: []string{
			"默认显示所有进程",
			"进程信息获取可能需要管理员权限",
			"复杂过滤可通过管道使用 grep 命令",
		},
		Examples: map[string]string{
			"查看所有进程":     "fck proc",
			"按名称查找":      "fck proc -n chrome",
			"查看指定 PID":   "fck proc -p 1234",
			"查看多个 PID":   "fck proc -P 1234,5678",
			"按 CPU 排序":   "fck proc -s cpu",
			"简洁模式":       "fck proc -l",
			"树形显示":       "fck proc --tree",
			"JSON 输出":    "fck proc --json",
			"配合 grep 过滤": "fck proc -l | grep admin",
		},
		MutexGroups: []qflag.MutexGroup{
			{Name: "display", Flags: []string{"list", "tree", "json"}, AllowNone: true},
		},
	}

	if err := ProcCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	ProcCmd.SetRun(runProc)
}

func runProc(cmd qflag.Command) error {
	config := &proc.ProcConfig{
		Name:       procName.Get(),
		PID:        int32(procPID.Get()),
		PIDs:       procPIDs.Get(),
		SortBy:     procSort.Get(),
		Ascend:     procAscend.Get(),
		ListMode:   procList.Get(),
		TreeMode:   procTree.Get(),
		JSONMode:   procJSON.Get(),
		TableStyle: procTableStyle.Get(),
	}

	return proc.ProcCmdMain(config)
}
