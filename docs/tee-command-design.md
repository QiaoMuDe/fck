# tee 命令设计方案

## 功能概述

实现 Unix/Linux `tee` 命令的功能：从标准输入读取数据，同时输出到标准输出和指定的文件。

## 参考实现

参考 Linux `tee` 命令的行为和选项设计。

## 命令结构

```
fck tee [options] [file...]
```

## 配置结构体

```go
// TeeConfig 配置结构体
type TeeConfig struct {
	Append           bool     // 追加模式
	IgnoreInterrupts bool     // 忽略中断信号
	Files            []string // 输出文件列表
}
```

## 命令选项

| 长选项 | 短选项 | 类型 | 默认值 | 说明 |
|--------|--------|------|--------|------|
| append | a | bool | false | 追加到文件（默认覆盖） |
| ignore-interrupts | i | bool | false | 忽略中断信号（SIGINT） |

## 功能设计

### 1. 基本用法

```bash
# 输出到屏幕并写入文件
echo "hello" | fck tee file.txt

# 追加到文件
echo "line 2" | fck tee -a file.txt

# 多文件输出
echo "data" | fck tee file1.txt file2.txt file3.txt

# 仅输出到屏幕（无文件参数）
echo "test" | fck tee
```

### 2. 管道链中使用

```bash
# 查看输出同时保存日志
make 2>&1 | fck tee build.log

# 处理数据并保存中间结果
cat data.txt | grep "error" | fck tee errors.log | wc -l

# base64 编码并保存
cat file.bin | fck base64 | fck tee encoded.b64
```

### 3. 信号处理

```bash
# 忽略 Ctrl+C，确保数据完整写入
long_running_command | fck tee -a output.log
```

## 核心实现逻辑

### 主流程

```go
// TeeCmdMain 执行 tee 命令
func TeeCmdMain(config TeeConfig) error {
	// 打开所有输出文件
	files, err := openFiles(config.Files, config.Append)
	if err != nil {
		return err
	}
	defer closeFiles(files)

	// 设置信号处理（如需要）
	if config.IgnoreInterrupts {
		setupSignalHandler()
	}

	// 创建多路写入器
	writers := createWriters(files)

	// 从 stdin 读取并同时写入所有目标
	return copyToMultiple(os.Stdin, writers)
}
```

### 文件操作

```go
// openFiles 打开所有输出文件
func openFiles(paths []string, append bool) ([]*os.File, error) {
	flag := os.O_CREATE | os.O_WRONLY
	if append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	var files []*os.File
	for _, path := range paths {
		file, err := os.OpenFile(path, flag, 0644)
		if err != nil {
			closeFiles(files)
			return nil, fmt.Errorf("open file %s: %w", path, err)
		}
		files = append(files, file)
	}
	return files, nil
}

// closeFiles 关闭所有文件
func closeFiles(files []*os.File) {
	for _, f := range files {
		f.Close()
	}
}
```

### 多路写入

```go
// multiWriter 同时写入多个目标
type multiWriter struct {
	writers []io.Writer
}

func (mw *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}

// copyToMultiple 从 reader 复制到多个 writer
func copyToMultiple(r io.Reader, writers []io.Writer) error {
	// 添加 stdout 到写入列表
	allWriters := append([]io.Writer{os.Stdout}, writers...)
	
	mw := &multiWriter{writers: allWriters}
	_, err := io.Copy(mw, r)
	return err
}
```

### 信号处理

```go
// setupSignalHandler 设置忽略 SIGINT
func setupSignalHandler() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	
	go func() {
		for range sigCh {
			// 忽略信号，继续执行
		}
	}()
}
```

## 文件结构

```
internal/
├── commands/
│   └── tee/
│       └── cmd_tee.go       # 业务逻辑
└── cli/
    └── tee.go               # CLI 定义
```

## 使用示例

```go
// 示例配置
cmdOpts := &qflag.CmdOpts{
    Desc: "从标准输入读取并输出到多个目标",
    Notes: []string{
        "同时输出到标准输出和指定的文件",
        "支持追加模式 (-a)",
        "支持多文件同时写入",
        "使用 -i 忽略中断信号，确保数据完整",
    },
    UseChinese:  true,
    UsageSyntax: fmt.Sprintf("%s tee [options] [file...]", qflag.Root.Name()),
    Examples: map[string]string{
        "写入文件":         fmt.Sprintf("echo 'hello' | %s tee file.txt", qflag.Root.Name()),
        "追加到文件":       fmt.Sprintf("echo 'line 2' | %s tee -a file.txt", qflag.Root.Name()),
        "多文件输出":       fmt.Sprintf("echo 'data' | %s tee file1.txt file2.txt", qflag.Root.Name()),
        "保存构建日志":     fmt.Sprintf("make 2>&1 | %s tee build.log", qflag.Root.Name()),
        "忽略中断信号":     fmt.Sprintf("long_cmd | %s tee -i output.log", qflag.Root.Name()),
        "管道链中使用":     fmt.Sprintf("cat data | grep 'err' | %s tee errors.log | wc -l", qflag.Root.Name()),
    },
}
```

## 边缘情况处理

1. **无文件参数**：仅输出到 stdout（相当于 `cat`）
2. **文件无法创建**：返回错误，不写入任何文件
3. **部分文件写入失败**：已写入的数据不撤回，返回错误
4. **管道断开**：优雅处理，已写入数据保留
5. **权限不足**：返回友好的错误信息
6. **磁盘满**：返回 IO 错误

## 错误处理

| 错误场景 | 错误信息 |
|----------|----------|
| 文件无法创建 | "open file [path]: [reason]" |
| 写入失败 | "write to [path]: [reason]" |
| 权限不足 | "permission denied: [path]" |

## 测试用例建议

1. 基本写入功能
2. 追加模式
3. 多文件写入
4. 无文件参数（仅 stdout）
5. 大文件处理
6. 信号忽略功能
7. 错误处理（权限、磁盘满等）
8. 管道链中使用

## 与现有命令的配合

```bash
# 与 base64 配合
fck base64 -f file.bin | fck tee encoded.b64

# 与 grep 配合
cat log.txt | fck grep "ERROR" | fck tee errors.txt | fck wc -l

# 与 seq 配合
fck seq 100 | fck tee numbers.txt | fck head -5
```
