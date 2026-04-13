# Sed 命令替换计数逻辑重构方案

## 问题描述

当前 `sed` 命令的替换计数逻辑分散在多个函数中，存在以下问题：
1. **计数逻辑不一致**：`replaceWithString` 和 `replaceStringIgnoreCase` 处理超限制的方式不同
2. **重复检查**：`processLine` 和子函数都检查 `replaceCount >= MaxCount`
3. **职责不清**：子函数既负责替换又负责全局计数管理

## 重构目标

1. **单一职责**：`processLine` 负责流程控制，子函数只负责替换逻辑
2. **统一接口**：所有替换函数返回 `(string, int)` 表示（结果, 本次替换次数）
3. **集中计数**：`processLine` 统一更新 `config.replaceCount`

## 重构范围

涉及文件：`internal/commands/sed/cmd_sed.go`

### 1. 修改函数签名

#### processLine
```go
// 修改前
func processLine(line string, config *SedConfig, lineNum int) (string, bool)

// 修改后
func processLine(line string, config *SedConfig, lineNum int) (string, bool)
// 返回值含义不变，但内部逻辑调整
```

#### replaceWithString
```go
// 修改前
func replaceWithString(line string, config *SedConfig) (string, bool)

// 修改后
func replaceWithString(line string, config *SedConfig) (string, int)
// 返回 (处理后的行, 本次替换次数)
```

#### replaceStringIgnoreCase
```go
// 修改前
func replaceStringIgnoreCase(line, pattern, replacement string, config *SedConfig) (string, bool)

// 修改后
func replaceStringIgnoreCase(line string, config *SedConfig) (string, int)
// 简化参数，从 config 获取 pattern/replacement
```

#### replaceWithRegex
```go
// 修改前
func replaceWithRegex(line string, config *SedConfig) (string, bool)

// 修改后
func replaceWithRegex(line string, config *SedConfig) (string, int)
```

### 2. 新增辅助函数

```go
// inLineRange 检查行号是否在指定范围内
func inLineRange(config *SedConfig, lineNum int) bool {
    if config.lineStart == 0 {
        return true // 未指定范围
    }
    if lineNum < config.lineStart {
        return false
    }
    if config.lineEnd > 0 && lineNum > config.lineEnd {
        return false
    }
    return true
}

// calcActualReplaceCount 计算实际可替换次数
func calcActualReplaceCount(desired, remaining int) int {
    if desired > remaining {
        return remaining
    }
    return desired
}
```

### 3. 核心逻辑修改

#### processLine 新实现
```go
func processLine(line string, config *SedConfig, lineNum int) (string, bool) {
    // 1. 行范围检查
    if !inLineRange(config, lineNum) {
        return line, false
    }
    
    // 2. 全局替换次数检查（快速路径）
    if config.MaxCount > 0 && config.replaceCount >= config.MaxCount {
        return line, false
    }
    
    // 3. 执行替换
    var result string
    var count int
    
    switch {
    case config.Regexp:
        result, count = replaceWithRegex(line, config)
    case config.IgnoreCase:
        result, count = replaceStringIgnoreCase(line, config)
    default:
        result, count = replaceWithString(line, config)
    }
    
    // 4. 更新全局计数
    if count > 0 {
        config.replaceCount += count
        return result, true
    }
    
    return line, false
}
```

#### replaceWithString 新实现
```go
func replaceWithString(line string, config *SedConfig) (string, int) {
    pattern := config.Pattern
    replacement := config.Replacement
    
    // 统计匹配次数
    count := strings.Count(line, pattern)
    if count == 0 {
        return line, 0
    }
    
    // 计算实际可替换次数（考虑全局限制）
    actual := count
    if config.MaxCount > 0 {
        remaining := config.MaxCount - config.replaceCount
        actual = calcActualReplaceCount(count, remaining)
    }
    
    if actual <= 0 {
        return line, 0
    }
    
    result := replaceN(line, pattern, replacement, actual)
    return result, actual
}
```

#### replaceStringIgnoreCase 新实现
```go
func replaceStringIgnoreCase(line string, config *SedConfig) (string, int) {
    pattern := config.Pattern
    replacement := config.Replacement
    
    lowerLine := strings.ToLower(line)
    lowerPattern := strings.ToLower(pattern)
    
    // 统计匹配次数
    count := 0
    start := 0
    for {
        idx := strings.Index(lowerLine[start:], lowerPattern)
        if idx == -1 {
            break
        }
        count++
        start += idx + len(pattern)
    }
    
    if count == 0 {
        return line, 0
    }
    
    // 计算实际可替换次数
    actual := count
    if config.MaxCount > 0 {
        remaining := config.MaxCount - config.replaceCount
        actual = calcActualReplaceCount(count, remaining)
    }
    
    if actual <= 0 {
        return line, 0
    }
    
    // 执行替换
    var result strings.Builder
    start = 0
    replaced := 0
    for replaced < actual {
        idx := strings.Index(lowerLine[start:], lowerPattern)
        if idx == -1 {
            break
        }
        actualIdx := start + idx
        result.WriteString(line[start:actualIdx])
        result.WriteString(replacement)
        start = actualIdx + len(pattern)
        replaced++
    }
    result.WriteString(line[start:])
    
    return result.String(), actual
}
```

#### replaceWithRegex 新实现
```go
func replaceWithRegex(line string, config *SedConfig) (string, int) {
    if config.compiledPattern == nil {
        return line, 0
    }
    
    // 检查是否有匹配
    matches := config.compiledPattern.FindAllStringIndex(line, -1)
    if len(matches) == 0 {
        return line, 0
    }
    
    // 计算实际可替换次数
    actual := len(matches)
    if config.MaxCount > 0 {
        remaining := config.MaxCount - config.replaceCount
        actual = calcActualReplaceCount(len(matches), remaining)
    }
    
    if actual <= 0 {
        return line, 0
    }
    
    // 执行替换
    if actual == len(matches) {
        // 全部替换
        result := config.compiledPattern.ReplaceAllString(line, config.Replacement)
        return result, actual
    }
    
    // 部分替换
    result := replaceRegexN(line, config.compiledPattern, config.Replacement, actual)
    return result, actual
}
```

### 4. 删除的函数

- `replaceN` - 被内联到 `replaceWithString` 中
- `replaceRegexN` - 保留，但只在 `replaceWithRegex` 内部使用

### 5. 测试用例

#### 测试1：基础替换计数
```go
config := SedConfig{
    Pattern:     "a",
    Replacement: "X",
    MaxCount:    3,
}
// 输入: "aaa aaa"
// 第1行: 替换3个 -> "Xaa aaa", count=3
// 第2行: 已达上限 -> "aaa aaa", count=3 (不变)
```

#### 测试2：大小写不敏感模式
```go
config := SedConfig{
    Pattern:     "A",
    Replacement: "X",
    IgnoreCase:  true,
    MaxCount:    2,
}
// 输入: "aAaA"
// 期望: 替换2个 -> "XAXa", count=2
```

#### 测试3：正则模式
```go
config := SedConfig{
    Pattern:     `\d+`,
    Replacement: "N",
    Regexp:      true,
    MaxCount:    2,
}
// 输入: "1 2 3 4"
// 期望: 替换2个 -> "N N 3 4", count=2
```

#### 测试4：行范围 + 计数限制
```go
config := SedConfig{
    Pattern:     "a",
    Replacement: "X",
    LineRange:   "1,3",
    MaxCount:    5,
}
// 第1行: "aaa" -> "Xaa", count=1
// 第2行: "aaa" -> "Xaa", count=2
// 第3行: "aaa" -> "Xaa", count=3
// 第4行: 超出范围，不处理
```

## 实施步骤

1. **添加辅助函数**：`inLineRange`, `calcActualReplaceCount`
2. **修改子函数**：按新签名修改 `replaceWithString`, `replaceStringIgnoreCase`, `replaceWithRegex`
3. **修改 processLine**：使用新逻辑
4. **删除冗余代码**：`replaceN`（如不再需要）
5. **运行测试**：验证所有测试用例通过

## 回滚方案

如果出现问题，可以通过 git 回滚到重构前的版本：
```bash
git checkout HEAD -- internal/commands/sed/cmd_sed.go
```

## 预期收益

1. **代码清晰**：每个函数职责单一
2. **易于测试**：子函数纯逻辑，不依赖全局状态
3. **一致行为**：所有替换模式处理计数的方式统一
4. **便于扩展**：新增替换模式只需遵循约定
