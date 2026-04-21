package cli

import (
	"fmt"

	"gitee.com/MM-Q/fck/internal/commands/date"
	"gitee.com/MM-Q/qflag"
)

var DateCmd *qflag.Cmd

var (
	dateFormat    *qflag.StringFlag // 时间格式字符串
	dateTimestamp *qflag.StringFlag // 时间戳字符串
	dateTimezone  *qflag.StringFlag // 时区字符串
	dateUTC       *qflag.BoolFlag   // 是否使用UTC时间
	dateUnix      *qflag.BoolFlag   // 输出Unix时间戳
)

func init() {
	DateCmd = qflag.NewCmd("date", "", qflag.ExitOnError)

	dateFormat = DateCmd.String("format", "f", "指定时间格式字符串", "")
	dateTimestamp = DateCmd.String("timestamp", "t", "将Unix时间戳转换为可读时间", "")
	dateTimezone = DateCmd.String("timezone", "z", "指定时区（默认为本地时区）", "")
	dateUTC = DateCmd.Bool("utc", "u", "使用UTC时间", false)
	dateUnix = DateCmd.Bool("unix", "U", "输出Unix时间戳", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "时间获取和格式化工具",
		UsageSyntax: fmt.Sprintf("%s date [options]", qflag.Root.Name()),
		Notes: []string{
			"预定义格式别名:\n" +
				"\tiso - 2006-01-02T15:04:05Z07:00\n" +
				"\tdate - 2006-01-02\n" +
				"\ttime - 15:04:05\n" +
				"\tdatetime - 2006-01-02 15:04:05\n" +
				"\tcompact - 2006-01-02T15:04:05\n" +
				"\tfull - 2006年01月02日 15时04分05秒\n" +
				"\tcn-date - 2006年01月02日\n" +
				"\tcn-time - 15时04分05秒\n" +
				"\tcn - 2006年01月02日 15时04分05秒\n" +
				"\tus-date - 01/02/2006\n" +
				"\tus-time - 03:04:05 PM\n" +
				"\tus - 01/02/2006 03:04:05 PM\n" +
				"\teu-date - 02.01.2006\n" +
				"\teu-time - 15:04:05\n" +
				"\teu - 02.01.2006 15:04:05\n" +
				"\tlog - 2006-01-02 15:04:05.000\n" +
				"\tfilename - 20060102_150405\n" +
				"\tnum - 20060102150405\n" +
				"\tdnum - 20060102\n" +
				"\ttnum - 150405",
			"时间戳支持秒级和毫秒级自动识别",
		},
		UseChinese: true,
		Examples: map[string]string{
			"获取当前时间":      fmt.Sprintf("%s date", qflag.Root.Name()),
			"将时间戳转换为可读时间": fmt.Sprintf("%s date -t 1694502400", qflag.Root.Name()),
			"获取指定时间格式":    fmt.Sprintf("%s date -f compact", qflag.Root.Name()),
		},
	}

	if err := DateCmd.ApplyOpts(cmdOpts); err != nil {
		panic(fmt.Errorf("apply opts err: %w", err))
	}

	DateCmd.SetRun(runDate)
}

func runDate(cmd qflag.Command) error {
	config := date.DateConfig{
		Format:    dateFormat.Get(),
		Timestamp: dateTimestamp.Get(),
		Timezone:  dateTimezone.Get(),
		UTC:       dateUTC.Get(),
		Unix:      dateUnix.Get(),
	}

	return date.DateCmdMain(config)
}
