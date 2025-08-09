// Package help 测试辅助工具
// 本文件提供了帮助信息模块的测试辅助函数和工具，
// 用于支持帮助信息生成和格式化功能的单元测试。
package help

import (
	"flag"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// createTestContext 创建测试用的命令上下文
func createTestContext(longName, shortName string) *types.CmdContext {
	ctx := types.NewCmdContext(longName, shortName, flag.ContinueOnError)
	return ctx
}

// addTestFlag 添加测试标志到上下文
func addTestFlag(ctx *types.CmdContext, longName, shortName, usage, flagType string, defValue interface{}) {
	var testFlag flags.Flag

	switch flagType {
	case "bool":
		boolFlag := &flags.BoolFlag{}
		currentBool := new(bool)
		*currentBool = defValue.(bool)
		if err := boolFlag.Init(longName, shortName, usage, currentBool); err != nil {
			return
		}
		testFlag = boolFlag
	case "string":
		stringFlag := &flags.StringFlag{}
		currentStr := new(string)
		*currentStr = defValue.(string)
		if err := stringFlag.Init(longName, shortName, usage, currentStr); err != nil {
			return
		}
		testFlag = stringFlag
	case "int":
		intFlag := &flags.IntFlag{}
		currentInt := new(int)
		*currentInt = defValue.(int)
		if err := intFlag.Init(longName, shortName, usage, currentInt); err != nil {
			return
		}
		testFlag = intFlag
	}

	if testFlag != nil {
		if err := ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: testFlag}); err != nil {
			return
		}
	}
}
