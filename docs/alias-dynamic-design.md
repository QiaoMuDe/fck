# alias 动态生成设计方案

## 1. 核心思路

创建一个别名注册中心，所有子命令（除 alias 自身外）在初始化时注册自己的别名信息（包括额外参数），alias 命令遍历注册表动态生成别名定义。

## 2. 数据结构

```go
// AliasRegistry 别名注册中心
type AliasRegistry struct {
    commands map[string]*AliasInfo // 命令名 -> 别名信息
}

// AliasInfo 单个命令的别名信息
type AliasInfo struct {
    Name      string // 别名（如 grep、ll）
    Command   string // 实际命令（如 grep、list），为空时与 Name 相同
    ExtraArgs string // 额外参数（如 -n -H -r）
}

// 全局注册中心实例
var GlobalAliasRegistry = &AliasRegistry{
    commands: make(map[string]*AliasInfo),
}
```

## 3. 注册机制

```go
// Register 注册别名信息
// 由各个子命令在 init() 中调用
func (r *AliasRegistry) Register(info *AliasInfo) {
    r.commands[info.Name] = info
}

// GetAll 获取所有注册的命令（排除指定命令）
func (r *AliasRegistry) GetAll(exclude string) []*AliasInfo {
    var result []*AliasInfo
    for name, info := range r.commands {
        if name != exclude {
            result = append(result, info)
        }
    }
    // 按名称排序
    sort.Slice(result, func(i, j int) bool {
        return result[i].Name < result[j].Name
    })
    return result
}
```

## 4. 子命令注册示例

```go
// internal/cli/grep.go
func init() {
    // ... 其他初始化代码 ...
    
    // 注册别名信息
    alias.GlobalAliasRegistry.Register(&alias.AliasInfo{
        Name:      "grep",
        ExtraArgs: "-n -H -r",
    })
}

// internal/cli/list.go
func init() {
    // ...
    alias.GlobalAliasRegistry.Register(&alias.AliasInfo{
        Name:      "list",
        ExtraArgs: "-c -i",
    })
    // 注册 ll 别名，实际执行 list 命令
    alias.GlobalAliasRegistry.Register(&alias.AliasInfo{
        Name:      "ll",
        Command:   "list",
        ExtraArgs: "-l -c -i",
    })
}
```

## 5. 动态生成逻辑

```go
// GenerateBashAliases 生成 Bash 别名定义
func GenerateBashAliases(commands []*AliasInfo, rootName string) string {
    var builder strings.Builder
    
    for _, cmd := range commands {
        command := cmd.Command
        if command == "" {
            command = cmd.Name
        }
        if cmd.ExtraArgs != "" {
            builder.WriteString(fmt.Sprintf("alias %s='%s %s %s'\n",
                cmd.Name, rootName, command, cmd.ExtraArgs))
        } else {
            builder.WriteString(fmt.Sprintf("alias %s='%s %s'\n",
                cmd.Name, rootName, command))
        }
    }
    
    return builder.String()
}

// GeneratePwshAliases 生成 PowerShell 别名定义
func GeneratePwshAliases(commands []*AliasInfo, rootName string) string {
    var builder strings.Builder
    
    // 生成函数定义
    for _, cmd := range commands {
        command := cmd.Command
        if command == "" {
            command = cmd.Name
        }
        funcName := rootName + "-" + cmd.Name
        if cmd.ExtraArgs != "" {
            builder.WriteString(fmt.Sprintf("function %s {\n", funcName))
            builder.WriteString(fmt.Sprintf("    %s %s %s @args\n",
                rootName, command, cmd.ExtraArgs))
            builder.WriteString("}\n\n")
        } else {
            builder.WriteString(fmt.Sprintf("function %s {\n", funcName))
            builder.WriteString(fmt.Sprintf("    %s %s @args\n", rootName, command))
            builder.WriteString("}\n\n")
        }
    }
    
    // 生成别名映射
    for _, cmd := range commands {
        funcName := rootName + "-" + cmd.Name
        builder.WriteString(fmt.Sprintf("Set-Alias -Name %s -Value %s\n",
            cmd.Name, funcName))
    }
    
    return builder.String()
}
```

## 6. alias 命令主逻辑

```go
// AliasCmdMain 执行 alias 命令
func AliasCmdMain(config AliasConfig, rootName string) error {
    // 获取所有注册的命令（排除 alias 自身）
    commands := GlobalAliasRegistry.GetAll("alias")
    
    if len(commands) == 0 {
        return fmt.Errorf("没有注册的命令")
    }
    
    switch config.Type {
    case "bash":
        fmt.Println(GenerateBashAliases(commands, rootName))
    case "pwsh":
        fmt.Println(GeneratePwshAliases(commands, rootName))
    default:
        return fmt.Errorf("不支持的 shell 类型: %s, 支持 bash/pwsh", config.Type)
    }
    
    return nil
}
```

## 7. 文件结构

```
internal/commands/alias/
├── cmd_alias.go          # 主逻辑和注册中心
├── generator.go          # 生成器（Bash/Pwsh）
└── types.go              # 类型定义

internal/cli/
├── alias.go              # CLI 定义
├── grep.go               # 注册别名信息
├── list.go               # 注册别名信息
└── ...                   # 其他命令注册别名
```

## 8. 使用示例

```go
// types.go
package alias

type AliasInfo struct {
    Name      string
    Command   string
    ExtraArgs string
}

// cmd_alias.go
package alias

type AliasRegistry struct {
    commands map[string]*AliasInfo
}

var GlobalAliasRegistry = &AliasRegistry{
    commands: make(map[string]*AliasInfo),
}

func (r *AliasRegistry) Register(info *AliasInfo) {
    r.commands[info.Name] = info
}

// generator.go
package alias

func GenerateBashAliases(commands []*AliasInfo, rootName string) string { ... }
func GeneratePwshAliases(commands []*AliasInfo, rootName string) string { ... }
```

## 9. 生成的别名示例

### Bash
```bash
alias cat='fck cat'
alias grep='fck grep -n -H -r'
alias ls='fck list -c -i'
alias ll='fck list -l -c -i'
```

### PowerShell
```powershell
function fck-cat {
    fck cat @args
}

function fck-grep {
    fck grep -n -H -r @args
}

function fck-ls {
    fck list -c -i @args
}

function fck-ll {
    fck list -l -c -i @args
}

Set-Alias -Name cat -Value fck-cat
Set-Alias -Name grep -Value fck-grep
Set-Alias -Name ls -Value fck-ls
Set-Alias -Name ll -Value fck-ll
```

> 注：如果 rootName 为 "myapp"，则生成 `myapp-cat`、`myapp grep` 等形式

## 10. 优势

1. **解耦**：alias 命令不硬编码命令列表
2. **自注册**：每个命令自己决定别名参数
3. **可扩展**：新增命令自动包含在别名生成中
4. **灵活**：每个命令可独立配置 Bash/Pwsh 参数
