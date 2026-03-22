package parser

import (
	"os"

	"gitee.com/MM-Q/qflag/internal/types"
)

// loadEnvVars 加载命令的所有环境变量
//
// 参数:
//   - cmd: 要加载环境变量的命令
//   - envPrefix: 环境变量前缀
//
// 返回值:
//   - error: 如果加载失败返回错误
//
// 注意事项:
//   - 遍历命令的所有标志
//   - 环境变量优先级低于命令行参数
func (p *DefaultParser) loadEnvVars(cmd types.Command, envPrefix string) error {
	flags := cmd.Flags()

	for _, f := range flags {
		if err := p.loadFlagEnv(f, envPrefix); err != nil {
			return err
		}
	}

	return nil
}

// loadFlagEnv 加载单个标志的环境变量
//
// 参数:
//   - f: 要加载环境变量的标志
//   - envPrefix: 环境变量前缀
//
// 返回值:
//   - error: 如果加载失败返回错误
//
// 注意事项:
//   - 只有绑定了环境变量的标志才会加载
//   - 环境变量名由前缀和标志的环境变量名组成
//   - 如果环境变量不存在, 跳过加载
//   - 只有在标志未被命令行参数设置时才加载环境变量
//   - 环境变量值通过标志的Set方法设置
func (p *DefaultParser) loadFlagEnv(f types.Flag, envPrefix string) error {
	envVar := f.GetEnvVar()
	if envVar == "" {
		return nil
	}

	// 如果标志已经被命令行参数设置, 则跳过环境变量加载
	if f.IsSet() {
		return nil
	}

	var fullEnvVar string
	if envPrefix != "" {
		fullEnvVar = envPrefix + envVar
	} else {
		fullEnvVar = envVar
	}

	value, exists := os.LookupEnv(fullEnvVar)
	if !exists {
		return nil
	}

	return f.Set(value)
}
