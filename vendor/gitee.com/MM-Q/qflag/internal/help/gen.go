package help

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/internal/utils"
)

// GenHelp 生成帮助信息
//
// 参数:
//   - cmd: 要生成帮助信息的命令
//
// 返回值:
//   - string: 生成的帮助信息字符串
func GenHelp(cmd types.Command) string {
	var buf strings.Builder
	cfg := cmd.Config()
	if cfg == nil {
		return "cmd config is nil"
	}

	// 写入命令logo
	writeLogo(cfg, &buf)

	// 写入命令名称
	writeName(cmd, cfg, &buf)

	// 写入命令描述
	writeDesc(cmd, cfg, &buf)

	// 写入命令使用方法
	writeUsage(cmd, cfg, &buf)

	// 写入命令选项
	writeOptions(cmd, cfg, &buf)

	// 写入命令子命令
	writeSubCmds(cmd, cfg, &buf)

	// 写入命令示例
	writeExample(cfg, &buf)

	// 写入命令注意事项
	writeNotes(cfg, &buf)

	return buf.String()
}

// writeName 写入命令名称
//
// 参数:
//   - cmd: 要生成帮助信息的命令
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeName(cmd types.Command, cfg *types.CmdConfig, buf *strings.Builder) {
	if cfg.UseChinese {
		buf.WriteString(types.HelpNameCN)
	} else {
		buf.WriteString(types.HelpNameEN)
	}

	// 写入命令名称
	buf.WriteString(types.HelpPrefix)
	buf.WriteString(utils.GetCmdName(cmd))
}

// writeLogo 写入命令logo
//
// 参数:
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeLogo(cfg *types.CmdConfig, buf *strings.Builder) {
	if cfg.LogoText != "" {
		fmt.Fprintf(buf, "\n\t\t%s\n", cfg.LogoText)
	}
}

// writeDesc 写入命令描述
//
// 参数:
//   - cmd: 要生成帮助信息的命令
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeDesc(cmd types.Command, cfg *types.CmdConfig, buf *strings.Builder) {
	if cmd.Desc() == "" {
		return
	}

	if cfg.UseChinese {
		buf.WriteString(types.HelpDescCN)
	} else {
		buf.WriteString(types.HelpDescEN)
	}

	// 写入命令描述
	buf.WriteString(types.HelpPrefix + cmd.Desc() + "\n")
}

// writeUsage 写入命令使用方法
//
// 参数:
//   - cmd: 要生成帮助信息的命令
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeUsage(cmd types.Command, cfg *types.CmdConfig, buf *strings.Builder) {
	if cfg.UseChinese {
		buf.WriteString(types.HelpUsageCN)
	} else {
		buf.WriteString(types.HelpUsageEN)
	}

	// 检查命令是否有使用方法
	if cfg.UsageSyntax != "" {
		buf.WriteString(types.HelpPrefix + cfg.UsageSyntax + "\n")
		return
	}

	// 没有指定时默认生成
	fmt.Fprintf(buf, "%s%s [options] [args...]\n", types.HelpPrefix, cmd.Path())
}

// writeOptions 写入命令选项
//
// 参数:
//   - cmd: 要生成帮助信息的命令
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeOptions(cmd types.Command, cfg *types.CmdConfig, buf *strings.Builder) {
	flags := cmd.Flags()
	if len(flags) == 0 {
		return
	}

	if cfg.UseChinese {
		buf.WriteString(types.HelpOptionsCN)
	} else {
		buf.WriteString(types.HelpOptionsEN)
	}

	// 收集命令选项
	options := make([]types.OptionInfo, 0, len(flags))
	for _, f := range flags {
		opt := types.OptionInfo{
			Desc:     f.Desc(),
			DefValue: utils.FormatDefaultValue(f.Type(), f.GetDef()),
		}

		if f.LongName() != "" && f.ShortName() != "" {
			opt.NamePart = fmt.Sprintf("-%s, --%s <%s>", f.ShortName(), f.LongName(), f.Type().String())
		} else if f.LongName() != "" {
			opt.NamePart = fmt.Sprintf("--%s <%s>", f.LongName(), f.Type().String())
		} else if f.ShortName() != "" {
			opt.NamePart = fmt.Sprintf("-%s <%s>", f.ShortName(), f.Type().String())
		}

		options = append(options, opt)
	}

	// 排序选项
	utils.SortOptions(options)

	// 计算选项名称最大宽度
	maxWidth := utils.CalcOptionMaxWidth(options)

	// 写入选项
	for _, opt := range options {
		fmt.Fprintf(buf, "  %-*s%s", maxWidth, opt.NamePart, types.HelpOptionSubCmdSpace)
		if opt.Desc != "" {
			buf.WriteString(opt.Desc)
		}

		if opt.DefValue != "" {
			fmt.Fprintf(buf, " (default: %s)", opt.DefValue)
		}
		buf.WriteByte('\n')
	}
}

// writeSubCmds 写入命令子命令
//
// 参数:
//   - cmd: 要生成帮助信息的命令
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeSubCmds(cmd types.Command, cfg *types.CmdConfig, buf *strings.Builder) {
	SubCmds := cmd.SubCmds()
	if len(SubCmds) == 0 {
		return
	}

	if cfg.UseChinese {
		buf.WriteString(types.HelpSubCmdsCN)
	} else {
		buf.WriteString(types.HelpSubCmdsEN)
	}

	// 收集子命令信息
	subCmds := make([]types.SubCmdInfo, 0, len(SubCmds))
	for _, subCmd := range SubCmds {
		info := types.SubCmdInfo{Desc: subCmd.Desc()}

		if subCmd.LongName() != "" && subCmd.ShortName() != "" {
			info.Name = fmt.Sprintf("%s, %s", subCmd.LongName(), subCmd.ShortName())
		} else if subCmd.LongName() != "" {
			info.Name = subCmd.LongName()
		} else if subCmd.ShortName() != "" {
			info.Name = subCmd.ShortName()
		} else {
			info.Name = subCmd.Name()
		}
		subCmds = append(subCmds, info)
	}

	// 排序子命令
	utils.SortSubCmds(subCmds)

	// 计算子命令名称最大宽度
	maxLen := utils.CalcSubCmdMaxLen(subCmds)

	// 写入子命令
	for _, info := range subCmds {
		fmt.Fprintf(buf, "  %-*s%s%s\n", maxLen, info.Name, types.HelpOptionSubCmdSpace, info.Desc)
	}
}

// writeExample 写入命令示例
//
// 参数:
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeExample(cfg *types.CmdConfig, buf *strings.Builder) {
	if len(cfg.Example) == 0 {
		return
	}

	if cfg.UseChinese {
		buf.WriteString(types.HelpExamplesCN)
	} else {
		buf.WriteString(types.HelpExamplesEN)
	}

	total := len(cfg.Example)
	jd := 0
	for k, v := range cfg.Example {
		jd++
		if jd == total {
			fmt.Fprintf(buf, "%s%d. %s\n     %s\n", types.HelpPrefix, jd, k, v)
		} else {
			fmt.Fprintf(buf, "%s%d. %s\n     %s\n\n", types.HelpPrefix, jd, k, v)
		}
	}
}

// writeNotes 写入命令注意事项
//
// 参数:
//   - cfg: 命令配置
//   - buf: 用于写入帮助信息的字符串构建器
func writeNotes(cfg *types.CmdConfig, buf *strings.Builder) {
	if len(cfg.Notes) == 0 {
		return
	}

	if cfg.UseChinese {
		buf.WriteString(types.HelpNotesCN)
	} else {
		buf.WriteString(types.HelpNotesEN)
	}

	for i, note := range cfg.Notes {
		fmt.Fprintf(buf, "%s%d. %s\n", types.HelpPrefix, i+1, note)
	}
}
