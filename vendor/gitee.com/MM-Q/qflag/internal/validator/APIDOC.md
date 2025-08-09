# Package validator

**Import Path:** `gitee.com/MM-Q/qflag/internal/validator`

**Source File:** `internal/subcmd/validator.go`

Package validator 内部验证器实现。本包实现了内部使用的验证器功能，提供命令和标志的验证逻辑，包括循环引用检测、命名冲突检查等内部验证机制。

## 功能模块

- **内部验证器实现** - 实现了内部使用的验证器功能，提供命令和标志的验证逻辑，包括循环引用检测、命名冲突检查等内部验证机制

## 目录

- [函数](#函数)
  - [GetCmdIdentifier](#getcmdidentifier)
  - [HasCycleFast](#hascyclefast)
  - [ValidateSubCommand](#validatesubcommand)

## 函数

### GetCmdIdentifier

```go
func GetCmdIdentifier(cmd *types.CmdContext) string
```

GetCmdIdentifier 获取命令的标识字符串，用于错误信息。

**参数:**
- `cmd`: 命令对象，类型为 `*types.CmdContext`

**返回值:**
- `string`: 命令标识字符串，如果为空则返回 `<nil>`

**功能特点:**
- 提供统一的命令标识生成逻辑
- 用于错误信息和日志记录
- 处理空命令的边界情况
- 返回可读的命令标识符

**使用场景:**
- 生成错误信息时标识具体命令
- 日志记录中标识命令上下文
- 调试信息中显示命令层级
- 验证失败时提供清晰的错误定位

### HasCycleFast

```go
func HasCycleFast(parent, child *types.CmdContext) bool
```

HasCycleFast 快速检测父命令和子命令之间是否存在循环依赖。

**核心原理:**
1. 只检查child的父链向上遍历，避免复杂的子树遍历
2. 利用CLI工具命令层级浅的特点（通常<10层）
3. 时间复杂度从O(n²)优化到O(d)，其中d是命令深度

**参数:**
- `parent`: 待添加的父命令上下文
- `child`: 待添加的子命令上下文

**返回值:**
- `bool`: true表示存在循环依赖，false表示安全

**算法优势:**
- **高效性**: 时间复杂度为O(d)，其中d为命令层级深度
- **针对性**: 专门针对CLI工具的浅层级特点优化
- **准确性**: 能够准确检测所有可能的循环依赖
- **简洁性**: 算法逻辑简单，易于理解和维护

**使用场景:**
- 在AddSubCmd函数中调用，防止添加会造成循环依赖的子命令
- 动态添加子命令时的安全检查
- 命令树结构验证
- 防止无限递归调用

### ValidateSubCommand

```go
func ValidateSubCommand(parent, child *types.CmdContext) error
```

ValidateSubCommand 验证单个子命令的有效性。

**参数:**
- `parent`: 当前上下文实例（父命令）
- `child`: 待添加的上下文实例（子命令）

**返回值:**
- `error`: 验证失败时返回的错误信息，否则返回nil

**验证内容:**
- **循环依赖检测**: 使用 `HasCycleFast` 检测循环引用
- **命名冲突检查**: 验证子命令名称是否与现有命令冲突
- **有效性验证**: 检查命令上下文的基本有效性
- **层级关系验证**: 确保父子关系的正确性

**验证规则:**
- 子命令不能与父命令形成循环依赖
- 子命令名称在同一父命令下必须唯一
- 子命令的长名称和短名称不能为空
- 子命令不能与系统保留名称冲突

## 使用示例

### 基本验证流程

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/qflag/internal/validator"
    "gitee.com/MM-Q/qflag/internal/types"
)

func main() {
    // 创建父命令
    parent := types.NewCmdContext("app", "a", flag.ExitOnError)
    
    // 创建子命令
    child := types.NewCmdContext("start", "s", flag.ExitOnError)
    
    // 验证子命令
    if err := validator.ValidateSubCommand(parent, child); err != nil {
        fmt.Printf("验证失败: %v\n", err)
        return
    }
    
    fmt.Println("子命令验证通过")
}
```

### 循环依赖检测

```go
func testCycleDetection() {
    // 创建命令链: A -> B -> C
    cmdA := types.NewCmdContext("cmdA", "a", flag.ExitOnError)
    cmdB := types.NewCmdContext("cmdB", "b", flag.ExitOnError)
    cmdC := types.NewCmdContext("cmdC", "c", flag.ExitOnError)
    
    // 建立父子关系
    cmdB.Parent = cmdA
    cmdC.Parent = cmdB
    
    // 检测循环依赖: 尝试让A成为C的子命令（会形成循环）
    hasCycle := validator.HasCycleFast(cmdC, cmdA)
    if hasCycle {
        fmt.Println("检测到循环依赖，操作被阻止")
    } else {
        fmt.Println("无循环依赖，操作安全")
    }
}
```

### 安全添加子命令

```go
func safeAddSubCommand(parent, child *types.CmdContext) error {
    // 1. 首先进行完整验证
    if err := validator.ValidateSubCommand(parent, child); err != nil {
        return fmt.Errorf("子命令验证失败: %w", err)
    }
    
    // 2. 验证通过后安全添加
    child.Parent = parent
    parent.SubCmds = append(parent.SubCmds, child)
    
    // 3. 更新映射表
    if parent.SubCmdMap == nil {
        parent.SubCmdMap = make(map[string]*types.CmdContext)
    }
    
    parent.SubCmdMap[child.LongName] = child
    if child.ShortName != "" {
        parent.SubCmdMap[child.ShortName] = child
    }
    
    return nil
}
```

### 批量验证子命令

```go
func validateMultipleSubCommands(parent *types.CmdContext, children []*types.CmdContext) error {
    // 临时映射表，用于检测批量添加中的名称冲突
    tempMap := make(map[string]*types.CmdContext)
    
    // 复制现有的子命令映射
    if parent.SubCmdMap != nil {
        for k, v := range parent.SubCmdMap {
            tempMap[k] = v
        }
    }
    
    for i, child := range children {
        // 验证单个子命令
        if err := validator.ValidateSubCommand(parent, child); err != nil {
            return fmt.Errorf("第%d个子命令验证失败: %w", i+1, err)
        }
        
        // 检查批量添加中的名称冲突
        if existing, exists := tempMap[child.LongName]; exists {
            return fmt.Errorf("子命令名称冲突: %s 已被 %s 使用", 
                child.LongName, validator.GetCmdIdentifier(existing))
        }
        
        if child.ShortName != "" {
            if existing, exists := tempMap[child.ShortName]; exists {
                return fmt.Errorf("子命令短名称冲突: %s 已被 %s 使用", 
                    child.ShortName, validator.GetCmdIdentifier(existing))
            }
        }
        
        // 添加到临时映射表
        tempMap[child.LongName] = child
        if child.ShortName != "" {
            tempMap[child.ShortName] = child
        }
    }
    
    return nil
}
```

### 命令标识符使用

```go
func demonstrateIdentifier() {
    // 创建命令
    cmd := types.NewCmdContext("deploy", "d", flag.ExitOnError)
    cmd.Config = types.NewCmdConfig()
    cmd.Config.Description = "部署应用程序"
    
    // 获取命令标识符
    identifier := validator.GetCmdIdentifier(cmd)
    fmt.Printf("命令标识符: %s\n", identifier)
    
    // 在错误信息中使用
    if err := someValidation(cmd); err != nil {
        fmt.Printf("命令 %s 验证失败: %v\n", identifier, err)
    }
}

func someValidation(cmd *types.CmdContext) error {
    // 模拟验证逻辑
    if cmd.Config == nil {
        return fmt.Errorf("缺少配置信息")
    }
    return nil
}
```

## 验证规则详解

### 循环依赖检测规则

1. **直接循环**: A -> B -> A
2. **间接循环**: A -> B -> C -> A
3. **自引用**: A -> A
4. **复杂循环**: A -> B -> C -> D -> B

### 命名冲突检测规则

1. **长名称唯一性**: 同一父命令下长名称必须唯一
2. **短名称唯一性**: 同一父命令下短名称必须唯一
3. **跨类型冲突**: 长名称不能与短名称冲突
4. **保留名称**: 不能使用系统保留的名称

### 有效性验证规则

1. **非空验证**: 命令名称不能为空
2. **格式验证**: 名称必须符合命名规范
3. **长度验证**: 名称长度必须在合理范围内
4. **字符验证**: 只能包含允许的字符

## 性能特点

### HasCycleFast 算法优化

- **时间复杂度**: O(d)，其中d为命令层级深度
- **空间复杂度**: O(1)，不需要额外存储空间
- **适用场景**: CLI工具的浅层级结构（通常<10层）
- **优化效果**: 相比传统O(n²)算法有显著性能提升

### 验证性能考虑

- 验证操作在命令注册阶段执行，不影响运行时性能
- 使用快速算法减少验证开销
- 缓存验证结果避免重复计算
- 早期失败策略，遇到错误立即返回

## 错误处理

验证器会返回以下类型的错误：

- **循环依赖错误**: 检测到命令间的循环引用
- **名称冲突错误**: 子命令名称与现有命令冲突
- **无效命令错误**: 命令上下文无效或不完整
- **格式错误**: 命令名称格式不符合规范

## 设计特点

1. **专门优化** - 针对CLI工具特点进行算法优化
2. **早期验证** - 在注册阶段进行验证，避免运行时错误
3. **详细错误** - 提供清晰的错误信息和定位
4. **高性能** - 使用高效算法减少验证开销
5. **可扩展** - 支持添加新的验证规则

## 最佳实践

1. **及时验证**: 在添加子命令时立即进行验证
2. **批量验证**: 批量操作时使用专门的验证函数
3. **错误处理**: 妥善处理验证失败的情况
4. **性能考虑**: 在性能敏感场景中合理使用验证
5. **测试覆盖**: 确保验证逻辑有充分的测试覆盖

## 注意事项

- 验证应在命令注册阶段完成，不要在运行时频繁验证
- 循环依赖检测算法针对浅层级优化，深层级场景需要考虑性能
- 命令标识符主要用于错误信息，不应用于业务逻辑
- 验证失败时应提供清晰的错误信息和修复建议
- 在并发环境中使用时需要注意线程安全