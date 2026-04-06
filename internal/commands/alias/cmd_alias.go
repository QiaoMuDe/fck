package alias

import (
	_ "embed"
	"fmt"
)

//go:embed templates/bash_aliases.tmpl
var bashAliases string

//go:embed templates/pwsh_aliases.tmpl
var pwshAliases string

// AliasConfig alias命令配置
type AliasConfig struct {
	Type string // Shell 类型
}

// AliasCmdMain 执行alias命令
//
// 参数:
//   - config: 命令配置
//
// 返回:
//   - error: 执行错误
func AliasCmdMain(config AliasConfig) error {
	switch config.Type {
	case "bash":
		fmt.Println(bashAliases)
	case "pwsh":
		fmt.Println(pwshAliases)
	default:
		return fmt.Errorf("不支持的 shell 类型: %s, 支持 bash/pwsh", config.Type)
	}
	return nil
}
