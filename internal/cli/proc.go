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
	ProcCmd = qflag.NewCmd("proc", "ps", qflag.ExitOnError)

	procName = ProcCmd.String("name", "n", "按进程名过滤，支持部分匹配", "")
	procPID = ProcCmd.Int("pid", "p", "指定单个 PID", 0)
	procPIDs = ProcCmd.IntSlice("pids", "P", "指定多个 PID, 如: 1234,5678", []int{})
	procSort = ProcCmd.Enum("sort", "s", "指定排序字段 (pid/name/cpu/mem/time)", "pid", SortFields)
	procAscend = ProcCmd.Bool("asc", "", "升序排列（默认降序）", false)
	procList = ProcCmd.Bool("list", "l", "简洁模式", false)
	procTree = ProcCmd.Bool("tree", "", "树形显示进程关系", false)
	procJSON = ProcCmd.Bool("json", "", "JSON 格式输出", false)
	procTableStyle = ProcCmd.Enum("table-style", "ts", qflag.EnumHelp("指定表格样式，支持以下选项:", types.TableStyleOptions, "\t\t\t\t\t"), "none", types.TableStyles)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "查看系统进程信息",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s proc [options]", qflag.Root.Name()),
		Notes: []string{
			"默认显示所有进程",
			"进程信息获取可能需要管理员权限",
			"复杂过滤可通过管道使用 grep 命令",
		},
		Examples: map[string]string{
			"查看所有进程":     fmt.Sprintf("%s proc", qflag.Root.Name()),
			"按名称查找":      fmt.Sprintf("%s proc -n chrome", qflag.Root.Name()),
			"查看指定 PID":   fmt.Sprintf("%s proc -p 1234", qflag.Root.Name()),
			"查看多个 PID":   fmt.Sprintf("%s proc -P 1234,5678", qflag.Root.Name()),
			"按 CPU 排序":   fmt.Sprintf("%s proc -s cpu", qflag.Root.Name()),
			"简洁模式":       fmt.Sprintf("%s proc -l", qflag.Root.Name()),
			"树形显示":       fmt.Sprintf("%s proc --tree", qflag.Root.Name()),
			"JSON 输出":    fmt.Sprintf("%s proc --json", qflag.Root.Name()),
			"配合 grep 过滤": fmt.Sprintf("%s proc -l | grep admin", qflag.Root.Name()),
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
