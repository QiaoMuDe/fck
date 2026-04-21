# tr 命令设计方案

## 命令概述

`tr`（translate）命令用于转换或删除标准输入中的字符，支持字符映射、删除、压缩等功能。

## 功能需求

### 核心功能

1. **字符转换** - 将字符集1映射到字符集2
2. **字符删除** - 删除指定字符
3. **字符压缩** - 连续重复字符压缩为单个
4. **大小写转换** - 快速大小写转换
5. **字符类支持** - 预定义字符类（[:alpha:], [:digit:] 等）

### 使用场景

```bash
# 大小写转换
echo "Hello World" | fck tr 'a-z' 'A-Z'    # HELLO WORLD

# 删除字符
echo "hello 123 world" | fck tr -d '0-9'   # hello  world

# 压缩连续空格
echo "a    b" | fck tr -s ' '              # a b

# 使用字符类
echo "abc123" | fck tr '[:digit:]' 'X'     # abcXXX
```

## 命令设计

### 语法

```
fck tr [options] <set1> [set2]
```

### 参数说明

| 参数 | 说明 | 必需 |
|------|------|------|
| set1 | 源字符集 | 是 |
| set2 | 目标字符集（转换时使用） | 否 |

### 选项设计

| 选项 | 短标志 | 说明 | 默认值 |
|------|--------|------|--------|
| --delete | -d | 删除 set1 中的字符 | false |
| --squeeze-repeats | -s | 压缩 set1 中连续重复的字符 | false |
| --complement | -c | 使用 set1 的补集 | false |
| --truncate-set1 | -t | 将 set1 截断为 set2 的长度 | false |

## 字符集语法

### 基本字符
```
'abc'          # 单个字符
'a-z'          # 范围
'0-9'          # 数字范围
```

### 转义字符
```
'\n'           # 换行
'\t'           # 制表符
'\\'           # 反斜杠
```

### 预定义字符类

| 字符类 | 说明 | 等价于 |
|--------|------|--------|
| [:alnum:] | 字母数字 | a-zA-Z0-9 |
| [:alpha:] | 字母 | a-zA-Z |
| [:digit:] | 数字 | 0-9 |
| [:lower:] | 小写字母 | a-z |
| [:upper:] | 大写字母 | A-Z |
| [:space:] | 空白字符 | \t\n\v\f\r |
| [:blank:] | 空格和制表符 | \t |
| [:punct:] | 标点符号 | !"#$%&'()*+,-./:;<=>?@[\]^_`{\|}~ |
| [:xdigit:] | 十六进制数字 | 0-9a-fA-F |

## 实现方案

### 目录结构

```
internal/commands/tr/
└── cmd_tr.go

internal/cli/
└── tr.go
```

### 配置结构体

```go
// TrConfig tr 命令配置
type TrConfig struct {
	Set1             string // 源字符集
	Set2             string // 目标字符集
	Delete           bool   // 删除模式
	SqueezeRepeats   bool   // 压缩重复字符
	Complement       bool   // 使用补集
	TruncateSet1     bool   // 截断 set1
}
```

### 核心算法

#### 1. 字符集解析

```go
// parseCharSet 解析字符集字符串
// 支持：普通字符、范围、转义字符、字符类
func parseCharSet(set string) ([]rune, error)
```

#### 2. 字符映射表构建

```go
// buildMapping 构建字符映射表
func buildMapping(set1, set2 []rune, truncate bool) map[rune]rune
```

#### 3. 处理流程

```go
// TrCmdMain 主函数
func TrCmdMain(config TrConfig) error {
	// 1. 解析字符集
	set1, err := parseCharSet(config.Set1)
	if err != nil {
		return err
	}
	
	// 2. 根据模式处理
	if config.Delete {
		return deleteMode(set1, config.Complement)
	}
	
	if config.SqueezeRepeats {
		return squeezeMode(set1)
	}
	
	// 3. 转换模式（需要 set2）
	if config.Set2 == "" {
		return fmt.Errorf("set2 is required for translate mode")
	}
	
	set2, err := parseCharSet(config.Set2)
	if err != nil {
		return err
	}
	
	mapping := buildMapping(set1, set2, config.TruncateSet1)
	return translateMode(mapping)
}
```

### 处理模式

#### 删除模式 (-d)
```go
func deleteMode(set1 []rune, complement bool) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)
	
	set1Map := make(map[rune]bool)
	for _, r := range set1 {
		set1Map[r] = true
	}
	
	for scanner.Scan() {
		r := []rune(scanner.Text())[0]
		shouldDelete := set1Map[r]
		if complement {
			shouldDelete = !shouldDelete
		}
		
		if !shouldDelete {
			fmt.Print(string(r))
		}
	}
	return scanner.Err()
}
```

#### 压缩模式 (-s)
```go
func squeezeMode(set1 []rune) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)
	
	set1Map := make(map[rune]bool)
	for _, r := range set1 {
		set1Map[r] = true
	}
	
	var lastRune rune
	var hasLast bool
	
	for scanner.Scan() {
		r := []rune(scanner.Text())[0]
		
		if set1Map[r] {
			if hasLast && lastRune == r {
				continue // 跳过重复
			}
		}
		
		fmt.Print(string(r))
		lastRune = r
		hasLast = true
	}
	return scanner.Err()
}
```

#### 转换模式
```go
func translateMode(mapping map[rune]rune) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)
	
	for scanner.Scan() {
		r := []rune(scanner.Text())[0]
		
		if newR, ok := mapping[r]; ok {
			fmt.Print(string(newR))
		} else {
			fmt.Print(string(r))
		}
	}
	return scanner.Err()
}
```

## CLI 定义

```go
package cli

import (
	"fmt"
	"gitee.com/MM-Q/fck/internal/commands/tr"
	"gitee.com/MM-Q/qflag"
)

var TrCmd *qflag.Cmd

var (
	trDelete         *qflag.BoolFlag // -d, --delete
	trSqueezeRepeats *qflag.BoolFlag // -s, --squeeze-repeats
	trComplement     *qflag.BoolFlag // -c, --complement
	trTruncateSet1   *qflag.BoolFlag // -t, --truncate-set1
)

func init() {
	TrCmd = qflag.NewCmd("tr", "", qflag.ExitOnError)
	
	trDelete = TrCmd.Bool("delete", "d", "删除 set1 中的字符", false)
	trSqueezeRepeats = TrCmd.Bool("squeeze-repeats", "s", "压缩 set1 中连续重复的字符", false)
	trComplement = TrCmd.Bool("complement", "c", "使用 set1 的补集", false)
	trTruncateSet1 = TrCmd.Bool("truncate-set1", "t", "将 set1 截断为 set2 的长度", false)
	
	cmdOpts := &qflag.CmdOpts{
		Desc:       "字符转换工具",
		UseChinese: true,
		UsageSyntax: fmt.Sprintf("%s tr [options] <set1> [set2]", qflag.Root.Name()),
		Notes: []string{
			"从标准输入读取数据，转换后输出到标准输出",
			"set1 和 set2 支持字符范围（如 a-z）和字符类（如 [:digit:]）",
			"转换模式需要同时指定 set1 和 set2",
			"删除模式 (-d) 和压缩模式 (-s) 只需要 set1",
		},
		Examples: map[string]string{
			"大小写转换":     fmt.Sprintf("echo 'hello' | %s tr 'a-z' 'A-Z'", qflag.Root.Name()),
			"删除数字":       fmt.Sprintf("echo 'abc123' | %s tr -d '0-9'", qflag.Root.Name()),
			"压缩空格":       fmt.Sprintf("echo 'a    b' | %s tr -s ' '", qflag.Root.Name()),
			"使用字符类":     fmt.Sprintf("echo 'abc123' | %s tr '[:digit:]' 'X'", qflag.Root.Name()),
			"删除非字母":     fmt.Sprintf("echo 'a1b2c3' | %s tr -cd '[:alpha:]'", qflag.Root.Name()),
		},
	}
	
	TrCmd.ApplyOpts(cmdOpts)
	TrCmd.SetRun(runTr)
}

func runTr(cmd qflag.Command) error {
	args := cmd.Args()
	if len(args) < 1 {
		return fmt.Errorf("set1 is required")
	}
	
	config := tr.TrConfig{
		Set1:           args[0],
		Delete:         trDelete.Get(),
		SqueezeRepeats: trSqueezeRepeats.Get(),
		Complement:     trComplement.Get(),
		TruncateSet1:   trTruncateSet1.Get(),
	}
	
	if len(args) >= 2 {
		config.Set2 = args[1]
	}
	
	return tr.TrCmdMain(config)
}
```

## 边缘情况处理

1. **set1 和 set2 长度不一致**
   - 默认：set1 最后一个字符重复映射到 set2 剩余字符
   - `-t` 选项：截断 set1 到 set2 长度

2. **空字符集**
   - 报错：character set cannot be empty

3. **无效字符类**
   - 报错：invalid character class [:xxx:]

4. **模式冲突**
   - `-d` 和 `-s` 同时使用：先删除后压缩
   - `-d` 和 set2 同时指定：忽略 set2，只删除

## 测试用例

```go
// 基本转换
{"大小写转换", "hello", "a-z", "A-Z", "HELLO"},
{"数字转字母", "123", "0-9", "a-j", "bcd"},

// 删除模式
{"删除数字", "abc123", "0-9", "", "abc", true},
{"删除空格", "a b c", " ", "", "abc", true},

// 压缩模式
{"压缩空格", "a    b", " ", "", "a b", false, true},
{"压缩重复字母", "aaabbb", "ab", "", "ab", false, true},

// 字符类
{"数字类", "abc123", "[:digit:]", "X", "abcXXX"},
{"字母类", "abc123", "[:alpha:]", "X", "XXX123"},

// 补集
{"删除非字母", "a1b2c3", "[:alpha:]", "", "abc", true, false, true},
```

## 实现优先级

1. **P1** - 基本字符转换（set1 -> set2）
2. **P1** - 删除模式 (-d)
3. **P2** - 压缩模式 (-s)
4. **P2** - 字符类支持
5. **P3** - 补集模式 (-c)
6. **P3** - 截断模式 (-t)

## 参考文档

- Linux tr 手册：`man tr`
- POSIX tr 规范
