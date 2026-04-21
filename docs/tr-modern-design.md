# tr 命令现代化改造方案

## 目标
将预定义字符类实现为标志，同时保留位置参数支持自定义字符操作，提升用户体验。

## 设计原则
1. **简洁优先**：常用操作使用标志，减少记忆负担
2. **灵活保留**：复杂场景保留位置参数语法
3. **互斥清晰**：标志之间、标志与位置参数之间互斥
4. **向后兼容**：现有功能不受影响

## 新增标志

| 标志 | 短标志 | 对应字符类 | 说明 |
|------|--------|-----------|------|
| --alnum | | [:alnum:] | 字母数字 (a-zA-Z0-9) |
| --alpha | | [:alpha:] | 字母 (a-zA-Z) |
| --digit | | [:digit:] | 数字 (0-9) |
| --lower | | [:lower:] | 小写字母 (a-z) |
| --upper | | [:upper:] | 大写字母 (A-Z) |
| --space | | [:space:] | 空白字符 (\t\n\v\f\r ) |
| --blank | | [:blank:] | 空格和制表符 (\t ) |
| --punct | | [:punct:] | 标点符号 |
| --xdigit | | [:xdigit:] | 十六进制数字 (0-9a-fA-F) |

## 互斥规则

### 1. 字符类标志互斥
9个字符类标志之间互斥，只能指定其中一个。

### 2. 字符类标志与 set1 位置参数互斥
- 如果指定了字符类标志，则不能使用 set1 位置参数
- 如果指定了 set1 位置参数，则不能使用字符类标志
- 两者都未指定 → 报错

### 3. set2 位置参数规则
- 删除模式 (-d)：不需要 set2
- 压缩模式 (-s)：不需要 set2
- 转换模式：需要 set2（可以是位置参数或默认值）

## 使用场景

### 场景一：删除数字（标志方式）
```bash
# 新方式（简洁）
echo "abc123" | fck tr --digit -d
# 输出：abc

# 旧方式（仍支持）
echo "abc123" | fck tr "[:digit:]" -d
echo "abc123" | fck tr "0-9" -d
```

### 场景二：删除标点和数字
```bash
# 链式处理（两次调用）
echo "a1b2c3!" | fck tr --digit -d | fck tr --punct -d
# 输出：abc

# 或使用位置参数自定义字符集
echo "a1b2c3!" | fck tr "0-9!" -d
```

### 场景三：压缩空格
```bash
# 新方式
echo "a    b" | fck tr --space -s
# 输出：a b

# 旧方式
echo "a    b" | fck tr "[:space:]" -s
echo "a    b" | fck tr " " -s
```

### 场景四：数字转 X
```bash
# 新方式（需要指定 set2）
echo "abc123" | fck tr --digit "X"
# 输出：abcXXX

# 旧方式
echo "abc123" | fck tr "[:digit:]" "X"
```

### 场景五：自定义范围（只能用位置参数）
```bash
# 自定义范围只能用位置参数
echo "abc" | fck tr "a-c" "A-C"
# 输出：ABC

# 混合使用：标志 + 自定义 set2
echo "abc123" | fck tr --digit "X"  # 数字变 X
echo "abc123" | fck tr --digit "0"  # 数字变 0
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
	// 原有选项
	trDelete         *qflag.BoolFlag // -d, --delete
	trSqueezeRepeats *qflag.BoolFlag // -s, --squeeze-repeats
	trComplement     *qflag.BoolFlag // -c, --complement
	trTruncateSet1   *qflag.BoolFlag // -t, --truncate-set1

	// 新增字符类标志
	trAlnum  *qflag.BoolFlag // --alnum
	trAlpha  *qflag.BoolFlag // --alpha
	trDigit  *qflag.BoolFlag // --digit
	trLower  *qflag.BoolFlag // --lower
	trUpper  *qflag.BoolFlag // --upper
	trSpace  *qflag.BoolFlag // --space
	trBlank  *qflag.BoolFlag // --blank
	trPunct  *qflag.BoolFlag // --punct
	trXdigit *qflag.BoolFlag // --xdigit
)

func init() {
	TrCmd = qflag.NewCmd("tr", "", qflag.ExitOnError)

	// 原有选项
	trDelete = TrCmd.Bool("delete", "d", "删除 set1 中的字符", false)
	trSqueezeRepeats = TrCmd.Bool("squeeze-repeats", "s", "压缩 set1 中连续重复的字符", false)
	trComplement = TrCmd.Bool("complement", "c", "使用 set1 的补集", false)
	trTruncateSet1 = TrCmd.Bool("truncate-set1", "t", "将 set1 截断为 set2 的长度", false)

	// 新增字符类标志
	trAlnum = TrCmd.Bool("alnum", "", "使用字母数字字符类 [:alnum:]", false)
	trAlpha = TrCmd.Bool("alpha", "", "使用字母字符类 [:alpha:]", false)
	trDigit = TrCmd.Bool("digit", "", "使用数字字符类 [:digit:]", false)
	trLower = TrCmd.Bool("lower", "", "使用小写字母字符类 [:lower:]", false)
	trUpper = TrCmd.Bool("upper", "", "使用大写字母字符类 [:upper:]", false)
	trSpace = TrCmd.Bool("space", "", "使用空白字符类 [:space:]", false)
	trBlank = TrCmd.Bool("blank", "", "使用空格和制表符字符类 [:blank:]", false)
	trPunct = TrCmd.Bool("punct", "", "使用标点符号字符类 [:punct:]", false)
	trXdigit = TrCmd.Bool("xdigit", "", "使用十六进制数字字符类 [:xdigit:]", false)

	cmdOpts := &qflag.CmdOpts{
		Desc:        "字符转换工具",
		UseChinese:  true,
		UsageSyntax: fmt.Sprintf("%s tr [options] [set1] [set2]", qflag.Root.Name()),
		Notes: []string{
			"从标准输入读取数据，转换后输出到标准输出",
			"支持两种方式指定字符集：",
			"  1. 使用字符类标志（如 --digit, --alpha）",
			"  2. 使用位置参数 set1（如 'a-z', '[:digit:]')",
			"字符类标志与 set1 位置参数互斥，只能使用其一",
			"删除模式 (-d) 和压缩模式 (-s) 不需要 set2",
			"转换模式需要 set2（替换字符）",
		},
		Examples: map[string]string{
			"删除数字(标志方式)": fmt.Sprintf("echo 'abc123' | %s tr --digit -d", qflag.Root.Name()),
			"删除数字(位置参数)": fmt.Sprintf("echo 'abc123' | %s tr '0-9' -d", qflag.Root.Name()),
			"压缩空格":           fmt.Sprintf("echo 'a    b' | %s tr --space -s", qflag.Root.Name()),
			"数字转X":            fmt.Sprintf("echo 'abc123' | %s tr --digit 'X'", qflag.Root.Name()),
			"删除标点":           fmt.Sprintf("echo 'a,b.c' | %s tr --punct -d", qflag.Root.Name()),
			"只保留字母":         fmt.Sprintf("echo 'a1b2c3' | %s tr --alpha -cd", qflag.Root.Name()),
			"大小写转换":         fmt.Sprintf("echo 'hello' | %s tr 'a-z' 'A-Z'", qflag.Root.Name()),
		},
		MutexGroups: []qflag.MutexGroup{
			{
				Name: "char-class",
				Flags: []string{
					"alnum", "alpha", "digit", "lower", "upper",
					"space", "blank", "punct", "xdigit",
				},
				AllowNone: true, // 允许都不选（此时使用位置参数）
			},
		},
	}

	TrCmd.ApplyOpts(cmdOpts)
	TrCmd.SetRun(runTr)
}

func runTr(cmd qflag.Command) error {
	args := cmd.Args()

	// 确定 set1 来源
	var set1 string
	charClassFlag := getCharClassFlag()

	if charClassFlag != "" {
		// 使用字符类标志
		if len(args) >= 1 && args[0] != "" {
			return fmt.Errorf("cannot specify both character class flag (--%s) and set1 positional argument", charClassFlag)
		}
		set1 = fmt.Sprintf("[:%s:]", charClassFlag)
	} else {
		// 使用位置参数
		if len(args) < 1 {
			return fmt.Errorf("set1 is required (use --help to see available character class flags)")
		}
		set1 = args[0]
	}

	// 确定 set2
	var set2 string
	if len(args) >= 2 {
		set2 = args[1]
	}

	config := tr.TrConfig{
		Set1:           set1,
		Set2:           set2,
		Delete:         trDelete.Get(),
		SqueezeRepeats: trSqueezeRepeats.Get(),
		Complement:     trComplement.Get(),
		TruncateSet1:   trTruncateSet1.Get(),
	}

	return tr.TrCmdMain(config)
}

// getCharClassFlag 返回被选中的字符类标志名称
func getCharClassFlag() string {
	if trAlnum.Get() {
		return "alnum"
	}
	if trAlpha.Get() {
		return "alpha"
	}
	if trDigit.Get() {
		return "digit"
	}
	if trLower.Get() {
		return "lower"
	}
	if trUpper.Get() {
		return "upper"
	}
	if trSpace.Get() {
		return "space"
	}
	if trBlank.Get() {
		return "blank"
	}
	if trPunct.Get() {
		return "punct"
	}
	if trXdigit.Get() {
		return "xdigit"
	}
	return ""
}
```

## 错误处理

### 错误场景 1：同时指定标志和位置参数
```bash
$ echo "abc" | fck tr --digit "0-9" -d
Error: cannot specify both character class flag (--digit) and set1 positional argument
```

### 错误场景 2：未指定字符集
```bash
$ echo "abc" | fck tr -d
Error: set1 is required (use --help to see available character class flags)
```

### 错误场景 3：转换模式缺少 set2
```bash
$ echo "abc123" | fck tr --digit
Error: set2 is required for translate mode
```

## 实现步骤

1. **修改 `internal/cli/tr.go`**
   - 添加 9 个字符类标志
   - 添加互斥组配置
   - 修改 `runTr` 函数处理标志和位置参数的优先级
   - 添加 `getCharClassFlag` 辅助函数

2. **无需修改 `internal/commands/tr/cmd_tr.go`**
   - 业务逻辑保持不变
   - 字符类解析逻辑已存在

3. **更新帮助文档**
   - 示例展示新旧两种方式
   - 说明互斥规则

## 向后兼容性

- 所有现有用法完全兼容
- 新增标志是可选的
- 不破坏任何现有功能
