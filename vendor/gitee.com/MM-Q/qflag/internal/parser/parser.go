// Package parser 命令行参数解析器
// 本文件实现了命令行参数的解析逻辑，包括标志解析、参数分离、
// 子命令识别等核心解析功能，为命令行参数处理提供基础支持。
package parser

import (
	"fmt"
	"runtime/debug"

	"gitee.com/MM-Q/qflag/internal/help"
	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/qerr"
)

// ParseCommand 解析单个命令的标志和参数
//
// 参数:
//   - ctx: 命令上下文
//   - args: 命令行参数
//
// 返回值:
//   - error: 如果解析失败，返回错误信息
func ParseCommand(ctx *types.CmdContext, args []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v\nStack: %s", qerr.ErrPanicRecovered, r, debug.Stack())
		}
	}()

	// 添加默认的注意事项
	if ctx.Config.UseChinese {
		ctx.Config.Notes = append(ctx.Config.Notes, help.ChineseTemplate.DefaultNote)
	} else {
		ctx.Config.Notes = append(ctx.Config.Notes, help.EnglishTemplate.DefaultNote)
	}

	// 加载当前命令的环境变量
	if err = LoadEnvVars(ctx); err != nil {
		return fmt.Errorf("%w: %v", qerr.ErrEnvLoadFailed, err)
	}

	// 解析当前命令的标志
	if err = ctx.FlagSet.Parse(args); err != nil {
		return fmt.Errorf("%w: %w", qerr.ErrFlagParseFailed, err)
	}

	// 获取非标志参数
	parsedArgs := ctx.FlagSet.Args()

	// 更新当前命令的非标志参数
	ctx.Args = append(ctx.Args, parsedArgs...) // 设置非标志参数

	return nil
}
