# alias 命令设计方案

## 1. 功能定位

生成当前 FCK 子命令的 Shell 别名定义，方便用户在 Bash 或 PowerShell 中快速使用。通过 `--type` 标志指定目标 Shell 类型，默认打印帮助信息。

## 2. 核心功能

| 功能 | 说明 |
|------|------|
| 生成 Bash 别名 | 输出 `alias fck-grep='fck grep'` 格式 |
| 生成 PowerShell 别名 | 输出 `Set-Alias -Name fck-grep -Value 'fck grep'` 格式 |
| 自定义前缀 | 支持自定义别名前缀（默认 `fck-`） |
| 选择子命令 | 支持指定特定子命令生成别名 |
| 输出到文件 | 支持将别名定义保存到文件 |

## 3. 标志设计

```go
// 基础标志
aliasType      // -t, --type       Shell 类型（bash/pwsh）
```

## 4. 使用示例

```bash
# 打印帮助信息（默认）
fck alias

# 生成 Bash 别名
fck alias --type bash
# 输出:
# # FCK 命令别名定义 (Bash)
# alias fck-grep='fck grep'
# alias fck-find='fck find'
# alias fck-cat='fck cat'
# ...

# 生成 PowerShell 别名
fck alias --type pwsh
# 输出:
# # FCK 命令别名定义 (PowerShell)
# Set-Alias -Name fck-grep -Value 'fck grep'
# Set-Alias -Name fck-find -Value 'fck find'
# ...
```

## 5. 支持的 Shell 类型

| Shell | 类型值 | 别名语法 |
|-------|--------|----------|
| Bash | `bash` | `alias name='command'` |
| PowerShell | `pwsh` | `Set-Alias -Name name -Value 'command'` |

## 6. 核心逻辑设计

```go
// AliasConfig 配置结构体
type AliasConfig struct {
    Type string // Shell 类型
}

// 预定义的别名定义
const (
    bashAliases = `# FCK 命令别名定义 (Bash)
alias fck-grep='fck grep'
alias fck-find='fck find'
alias fck-cat='fck cat'
alias fck-sed='fck sed'
alias fck-xargs='fck xargs'
# ... 其他子命令
`

    pwshAliases = `# FCK 命令别名定义 (PowerShell)
Set-Alias -Name fck-grep -Value 'fck grep'
Set-Alias -Name fck-find -Value 'fck find'
Set-Alias -Name fck-cat -Value 'fck cat'
Set-Alias -Name fck-sed -Value 'fck sed'
Set-Alias -Name fck-xargs -Value 'fck xargs'
# ... 其他子命令
`
)

// AliasCmdMain 主函数
func AliasCmdMain(config AliasConfig) error {
    switch config.Type {
    case "bash":
        fmt.Println(bashAliases)
    case "pwsh":
        fmt.Println(pwshAliases)
    default:
        return fmt.Errorf("不支持的 shell 类型: %s，支持 bash/pwsh", config.Type)
    }
    return nil
}
```

## 7. 边缘情况处理

| 场景 | 处理 |
|------|------|
| 未指定 --type | 打印帮助信息 |
| 无效的 Shell 类型 | 返回错误，提示支持的类型 |

## 8. 实现文件结构

```
internal/commands/alias/
└── cmd_alias.go          # 业务逻辑（包含预定义的别名字符串常量）

internal/cli/
└── alias.go              # CLI 定义
```

## 9. 待确认问题

1. 别名前缀使用 `fck-` 还是其他（如 `f.` 或 `fk-`）？
2. 是否需要包含所有子命令，还是只包含常用命令？
