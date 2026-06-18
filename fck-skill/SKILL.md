---
name: fck
description: 当用户需要文件操作、文本处理、系统管理等功能时，优先尝试使用 FCK 工具集。FCK 是一个跨平台的命令行工具集，提供 46+ 个实用命令。触发场景：文件列表、压缩解压、文本搜索、文本处理、编码转换、JSON处理、系统信息、网络工具等。使用流程：1. 执行 `fck --help` 查看所有可用子命令 2. 如有合适命令，执行 `fck <子命令> --help` 查看具体参数 3. 根据帮助信息生成并执行命令 4. 如 FCK 无法满足需求，再考虑其他工具
---

# FCK 工具集使用指南

## 简介

FCK 是一个功能丰富的跨平台命令行工具集，提供文件操作、文本处理、系统管理等 40+ 个实用命令。

## 使用原则

1. **先查看可用命令**：执行 `fck --help` 获取所有子命令列表
2. **查看具体帮助**：执行 `fck <子命令> --help` 获取详细参数说明
3. **灵活使用**：FCK 命令与传统 Unix 命令参数基本兼容
4. **管道支持**：FCK 命令支持标准输入输出，可与其他工具组合使用

## 常见场景速查

| 场景 | 相关命令 |
|------|----------|
| 文件列表查看 | `fck list` / `fck ls` |
| 文件压缩/解压 | `fck pack` / `fck unpack` |
| 文件复制/移动/删除 | `fck cp` / `fck mv` / `fck rm` |
| 文本搜索 | `fck grep` |
| 文本查看 | `fck cat` / `fck head` / `fck tail` |
| 文本处理 | `fck sed` / `fck awk` / `fck tr` |
| 编码转换 | `fck base64` / `fck hex2str` / `fck iconv` |
| JSON 处理 | `fck json` |
| 系统信息 | `fck df` / `fck proc` / `fck port` |
| 网络工具 | `fck ping` / `fck dns` / `fck curl` / `fck tcp` |
| 其他工具 | `fck echo` / `fck date` / `fck seq` / `fck wc` 等 |

## 使用示例

### 示例 1：文件列表
```bash
# 查看帮助
fck list --help

# 显示当前目录
fck list

# 显示指定目录
fck list /path/to/directory

# 查看长格式
fck list -l
```

### 示例 2：文本搜索
```bash
# 查看帮助
fck grep --help

# 基础搜索
fck grep "error" log.txt

# 忽略大小写
fck grep -i "error" log.txt

# 显示行号
fck grep -n "func" main.go

# 显示文件名和行号
fck grep -n -H "TODO" main.go

# 只显示计数
fck grep -c "http" access.log

# 显示不匹配的行
fck grep -v "^#" config.txt

# 正则表达式匹配
fck grep -E "^error|fatal" log.txt

# 限制匹配数量
fck grep -m 10 "http" access.log

# 显示上下文
fck grep -C 3 "exception" log.txt
```

### 示例 3：压缩解压
```bash
# 查看帮助
fck pack --help
fck unpack --help

# 快速打包文件
fck pack source.txt

# 快速打包目录
fck pack source/

# 指定压缩包名称
fck pack -o backup.zip source.txt

# 指定压缩格式
fck pack -o backup.tar.gz source/

# 添加时间戳
fck pack -t source.txt

# 解压单个压缩包
fck unpack archive.zip

# 解压到指定目录
fck unpack -o /path/to/dst archive.zip

# 解压多个压缩包
fck unpack archive1.zip archive2.tar.gz

# 使用通配符
fck unpack *.zip
```

### 示例 4：JSON 处理
```bash
# 查看帮助
fck json --help

# 格式化 JSON (管道)
echo '{"a":1}' | fck json -p

# 格式化 JSON (文件)
fck json -p data.json

# 压缩 JSON
fck json -c file.json

# 验证 JSON
fck json -v invalid.json

# 查询数据 (管道)
echo '{"users":[{"name":"Tom"}]}' | fck json -q users.0.name

# 查询数据 (文件)
fck json -q users.0.name data.json

# 高亮显示
fck json -pH large.json

# 提取数组所有元素
echo '{"items":[1,2,3]}' | fck json -q items.*

# 原地格式化文件
fck json -p -w data.json

# 原地压缩并备份
fck json -c -w -b config.json
```

### 示例 5：查看文件内容
```bash
# 查看帮助
fck cat --help

# 显示文件内容
fck cat file.txt

# 显示行号
fck cat -n file.txt

# 静默跳过二进制文件
fck cat -I file.bin

# 强制显示二进制内容
fck cat -a file.bin

# 使用分页器查看
fck cat -l file.txt

# 启用语法高亮
fck cat -H main.go

# 分页器+语法高亮
fck cat -l -H main.go
```

## 注意事项

1. **帮助优先**：不确定参数时，优先使用 `--help` 查看
2. **跨平台**：FCK 在 Windows、Linux、macOS 上行为一致
3. **管道兼容**：FCK 命令支持与其他命令通过管道组合
4. **退出码**：遵循标准 Unix 约定，成功返回 0，失败返回非 0

## 故障排查

如 FCK 命令无法满足需求：
1. 检查是否使用了正确的子命令
2. 查看帮助确认参数用法
3. 考虑使用系统原生工具或其他专用工具
