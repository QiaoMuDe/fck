// Package parser 环境变量解析和处理
// 本文件实现了环境变量的解析和处理逻辑，支持从环境变量中读取标志值，
// 为命令行参数提供环境变量绑定和默认值设置功能。
package parser

import (
	"flag"
	"os"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/qerr"
)

// LoadEnvVars 从环境变量加载参数值
// 纯函数设计，不依赖结构体状态
//
// 参数:
//   - ctx: 命令上下文
//
// 返回值:
//   - error: 错误信息
func LoadEnvVars(ctx *types.CmdContext) error {
	// 存储读取错误
	var errors []error

	// 预分配map容量以提高性能,初始容量为已注册标志数量
	// 使用所有标志总数作为容量最大基准, 确保独立长/短标志场景下容量充足
	processedEnvs := make(map[string]bool, ctx.FlagRegistry.GetAllFlagsCount()) // 跟踪已处理的环境变量，避免重复处理

	// 遍历所有已注册的标志
	ctx.FlagSet.VisitAll(func(f *flag.Flag) {
		// 获取标志实例
		flagInstance, ok := f.Value.(flags.Flag)
		if !ok {
			return
		}

		// 获取环境变量名称
		envVar := flagInstance.GetEnvVar()
		if envVar == "" {
			// 环境变量未设置，提前返回
			return
		}

		// 检查是否已处理过该环境变量（避免长短标志重复处理）
		if processedEnvs[envVar] {
			return
		}

		// 读取环境变量值
		envValue := os.Getenv(envVar)
		if envValue == "" {
			return // 环境变量未设置，提前返回
		}

		// 标记该环境变量为已处理
		processedEnvs[envVar] = true

		// 设置标志值(使用现有Set方法进行类型转换)
		if err := f.Value.Set(envValue); err != nil {
			errors = append(errors, qerr.NewValidationErrorf("Failed to parse environment variable %s for flag %s: %v", envVar, f.Name, err))
		}
	})

	// 函数末尾返回聚合错误
	if len(errors) > 0 {
		return qerr.NewValidationErrorf("Failed to load environment variables: %v", qerr.JoinErrors(errors))
	}

	return nil
}
