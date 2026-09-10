package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/df"
	"gitee.com/MM-Q/fck/internal/types"
	"gitee.com/MM-Q/qflag"
)

var DFCmd *qflag.Cmd

var (
	dfLocal      *qflag.BoolFlag   // -l, --local      只显示本地文件系统
	dfType       *qflag.StringFlag // -t, --type       按文件系统类型过滤
	dfTotal      *qflag.BoolFlag   // --total          显示总计
	dfList       *qflag.BoolFlag   // -s, --simple     简洁模式
	dfTableStyle *qflag.EnumFlag   // --table-style    表格样式
	dfBytes      *qflag.BoolFlag   // -b, --bytes      以字节为单位显示大小
)

func init() {
	DFCmd = qflag.NewCmd("df", "", qflag.ExitOnError)

	dfLocal = DFCmd.Bool("local", "l", "只显示本地文件系统", false)
	dfType = DFCmd.String("type", "t", "按文件系统类型过滤", "")
	dfTotal = DFCmd.Bool("total", "T", "显示总计行", false)
	dfList = DFCmd.Bool("simple", "s", "简洁模式", false)
	dfBytes = DFCmd.Bool("bytes", "b", "以字节为单位显示大小", false)
	dfTableStyle = DFCmd.Enum("table-style", "ts", qflag.EnumHelp("指定表格样式，支持以下选项:", types.TableStyleOptions, "\t\t\t\t\t"), "none", types.TableStyles)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "查看磁盘空间使用情况",
		UsageSyntax: fmt.Sprintf("%s df [options]", qflag.Root.Name()),
		UseChinese:  true,
		Examples: map[string]string{
			"查看所有分区":  fmt.Sprintf("%s df", qflag.Root.Name()),
			"只显示本地分区": fmt.Sprintf("%s df -l", qflag.Root.Name()),
			"按类型过滤":   fmt.Sprintf("%s df -t ext4", qflag.Root.Name()),
			"显示总计":    fmt.Sprintf("%s df --total", qflag.Root.Name()),
			"简洁模式":    fmt.Sprintf("%s df -s", qflag.Root.Name()),
		},
		Notes: []string{
			"默认显示所有分区 (包括网络文件系统)",
			"大小自动转换为人类可读格式 (GB/MB/KB)",
			"Windows 上显示为 C:, D: 等盘符",
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name:      "local-type",
				Flags:     []string{"local", "type"},
				AllowNone: true,
			},
		},
	}

	if err := DFCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	DFCmd.SetRun(runDF)
}

func runDF(cmd qflag.Command) error {
	config := df.DFConfig{
		LocalOnly:  dfLocal.Get(),
		FSFilter:   dfType.Get(),
		ShowTotal:  dfTotal.Get(),
		ListMode:   dfList.Get(),
		TableStyle: dfTableStyle.Get(),
		ShowBytes:  dfBytes.Get(),
	}

	return df.DFCmdMain(config)
}
