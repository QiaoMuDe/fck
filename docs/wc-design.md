# WC 子命令设计方案

> **设计日期**: 2026-04-19  
> **设计目标**: 实现 Linux 风格的 wc (word count) 命令，支持行数、单词数、字节数、字符数统计

---

## 一、功能概述

WC (Word Count) 是一个文本统计工具，用于统计文件的行数、单词数、字节数和字符数。

### 1.1 核心功能

| 功能 | 说明 | 对应标志 |
|------|------|----------|
| 行数统计 | 统计文件的换行符数量 | `-l, --lines` |
| 单词数统计 | 统计单词数量（以空白分隔） | `-w, --words` |
| 字节数统计 | 统计字节数 | `-c, --bytes` |
| 字符数统计 | 统计 Unicode 字符数 | `-m, --chars` |
| 最大行宽 | 统计最长行的长度 | `-L, --max-line-length` |

### 1.2 与 Linux wc 的兼容性

- 默认行为与 Linux wc 一致（无标志时输出 行数、单词数、字节数）
- 支持多文件统计
- 支持管道输入
- 支持通配符展开（跨平台兼容）
- 多文件时输出总计行

---

## 二、CLI 标志设计

### 2.1 标志定义

| 长标志 | 短标志 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| `--lines` | `-l` | bool | false | 仅统计行数 |
| `--words` | `-w` | bool | false | 仅统计单词数 |
| `--bytes` | `-c` | bool | false | 仅统计字节数 |
| `--chars` | `-m` | bool | false | 仅统计字符数 |
| `--max-line-length` | `-L` | bool | false | 仅统计最大行长度 |

### 2.2 默认行为

- 无任何标志时，默认输出：**行数、单词数、字节数**（与 Linux wc 一致）
- 多个标志可同时使用，按固定顺序输出：行数、单词、字符、字节、最大行宽

### 2.3 输出格式

```
# 单文件（默认）
  10   20  150 file.txt
#  行数 单词 字节 文件名

# 单文件（仅行数）
  10 file.txt

# 多文件
  10   20  150 file1.txt
   5   10   80 file2.txt
  15   30  230 total

# 管道输入
  10   20  150
```

---

## 三、配置结构体

```go
package wc

// WcConfig 命令配置
type WcConfig struct {
	ShowLines    bool     // 显示行数
	ShowWords    bool     // 显示单词数
	ShowBytes    bool     // 显示字节数
	ShowChars    bool     // 显示字符数
	ShowMaxLine  bool     // 显示最大行长度
	Files        []string // 输入文件列表
}

// WcStats 统计结果
type WcStats struct {
	Lines      int64  // 行数
	Words      int64  // 单词数
	Bytes      int64  // 字节数
	Chars      int64  // 字符数
	MaxLineLen int64  // 最大行长度
	Filename   string // 文件名（管道输入时为 "-"）
}
```

---

## 四、核心算法设计

### 4.1 统计逻辑

```go
// 单次遍历完成所有统计
func countStats(reader io.Reader, filename string) (*WcStats, error) {
	stats := &WcStats{Filename: filename}
	scanner := bufio.NewScanner(reader)
	
	// 扩展 scanner 缓冲区以支持长行
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	
	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		
		// 字节数
		lineBytes := len(scanner.Bytes())
		stats.Bytes += int64(lineBytes)
		
		// 字符数（Unicode）
		stats.Chars += int64(utf8.RuneCountInString(line))
		
		// 最大行长度（字符数）
		lineLen := int64(len([]rune(line)))
		if lineLen > stats.MaxLineLen {
			stats.MaxLineLen = lineLen
		}
		
		// 单词数（以空白分隔）
		stats.Words += int64(len(strings.Fields(line)))
	}
	
	return stats, scanner.Err()
}
```

### 4.2 多文件处理流程

```
1. 检查管道输入（优先）
2. 展开通配符
3. 遍历处理每个文件
4. 累加统计总计
5. 按格式输出结果
```

---

## 五、文件结构

```
internal/commands/wc/
├── types.go      # 配置结构体和统计结果定义
└── cmd_wc.go     # 主逻辑实现
internal/cli/wc.go   # CLI 定义
```

---

## 六、使用示例

```bash
# 默认统计（行数、单词数、字节数）
fck wc file.txt

# 仅统计行数
fck wc -l file.txt

# 统计行数和单词数
fck wc -lw file.txt

# 统计字符数（Unicode）
fck wc -m file.txt

# 统计最大行长度
fck wc -L file.txt

# 多文件
fck wc file1.txt file2.txt file3.txt

# 通配符
fck wc *.txt

# 管道输入
cat file.txt | fck wc
ls -la | fck wc -l

# 递归统计目录下所有文件
find . -name "*.go" | fck wc -l
```

---

## 七、技术特点

1. **单次遍历**: 一次读取完成所有统计，性能最优
2. **大文件支持**: 使用 bufio.Scanner 支持大文件
3. **长行支持**: 扩展缓冲区支持超长行（1MB）
4. **Unicode 支持**: 正确统计多字节字符
5. **跨平台**: 通配符展开兼容 Windows/Linux/macOS
6. **管道优先**: 检测管道输入时优先处理

---

## 八、边缘情况处理

| 场景 | 处理方式 |
|------|----------|
| 空文件 | 所有统计为 0 |
| 无换行符结尾的文件 | 最后一行计入行数 |
| 二进制文件 | 正常统计，不报错 |
| 超长行（>1MB） | 扩展缓冲区或报错 |
| 无法读取的文件 | 报错并跳过，继续处理其他文件 |
| 管道输入 | 不显示文件名，仅输出统计 |

---

## 九、与现有命令的一致性

- 遵循 `qflag` 命令开发规范
- 错误信息使用英文
- 注释使用中文
- 支持 `utils.IsStdinPipe()` 检测管道输入
- 使用 `filepath.Glob()` 展开通配符
