# GM命令实施文档

## 一、命令概述

### 1.1 命令名称
**gm** (Git Metadata) - Git仓库元数据获取命令

### 1.2 命令描述
获取Git仓库的各种元数据信息，包括版本号、提交哈希、提交时间、仓库状态等。

### 1.3 命令用途
- 快速获取Git仓库的版本信息
- 检查当前目录是否为Git仓库
- 获取提交相关的元数据
- 查看仓库状态（clean/dirty）

---

## 二、功能需求分析

### 2.1 功能列表

#### 功能1：检查Git仓库
- **标志**: `--check` / `-c`
- **描述**: 检查当前目录是否为Git仓库
- **输出**:
  - 是Git仓库：不输出任何信息，直接正常退出（返回码0）
  - 不是Git仓库：输出错误信息并返回错误码
- **实现方式**: 使用 `git rev-parse --is-inside-work-tree` 命令

#### 功能2：获取版本号
- **标志**: `--version` / `-v`
- **描述**: 获取Git仓库的版本号（标签或提交描述）
- **输出**: 版本号字符串（支持JSON格式输出）
- **实现方式**: 使用 `git describe --tags --always --dirty` 命令
- **默认行为**: 如果不指定任何功能标志，默认执行此功能

#### 功能3：获取提交哈希
- **标志**: `--hash` / `-H`
- **描述**: 获取当前提交的哈希值
- **输出**: 提交哈希字符串（支持JSON格式输出）
- **可选参数**: `--abbrev` / `-a` 控制哈希缩写长度（默认7，0表示完整哈希）
- **实现方式**: 使用 `git rev-parse HEAD` 或 `git rev-parse --short=N HEAD` 命令

#### 功能4：获取提交时间
- **标志**: `--time` / `-t`
- **描述**: 获取最新提交的时间
- **输出**: 格式化的时间字符串（支持JSON格式输出）
- **可选参数**: `--format` / `-f` 自定义时间格式（默认："2006-01-02 15:04:05"）
- **实现方式**: 使用 `git log -1 --format=%cd --date=iso` 命令

#### 功能5：获取仓库状态
- **标志**: `--status` / `-s`
- **描述**: 获取Git仓库的状态
- **输出**: "clean"（干净）或 "dirty"（有未提交的更改）（支持JSON格式输出）
- **实现方式**: 使用 `git status --porcelain` 命令

### 2.2 标志设计

| 标志 | 短标志 | 类型 | 默认值 | 描述 |
|------|--------|------|--------|------|
| --check | -c | Bool | false | 检查是否为Git仓库 |
| --version | -v | Bool | false | 获取版本号 |
| --hash | -H | Bool | false | 获取提交哈希 |
| --time | -t | Bool | false | 获取提交时间 |
| --status | -s | Bool | false | 获取仓库状态 |
| --abbrev | -a | Int | 7 | 哈希缩写长度（配合--hash使用） |
| --format | -f | String | "2006-01-02 15:04:05" | 时间格式（配合--time使用） |
| --json | -j | Bool | false | 以JSON格式输出结果 |

### 2.3 互斥性设计
- **互斥组**: 暂不强制互斥，后期根据需要添加到互斥组
- **默认行为**: 如果不指定任何功能标志，默认执行 `--version`
- **JSON输出**: `--json` 标志可以与其他功能标志配合使用

### 2.4 JSON输出格式

当使用 `--json` / `-j` 标志时，所有功能将以JSON格式输出结果：

#### 版本号JSON格式
```json
{
  "version": "v1.2.3-5-g7a8b9c0"
}
```

#### 提交哈希JSON格式
```json
{
  "hash": "7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6"
}
```

#### 提交时间JSON格式
```json
{
  "time": "2026-03-24 15:30:45",
  "timestamp": "2026-03-24T15:30:45+08:00"
}
```

#### 仓库状态JSON格式
```json
{
  "status": "clean"
}
```

或

```json
{
  "status": "dirty"
}
```

---

## 三、技术方案设计

### 3.1 技术选型

#### 依赖库
- **qflag**: 命令行参数解析
- **shellx/shx**: 执行Git命令
- **time**: 时间格式化处理
- **strings**: 字符串处理

#### Git命令映射

| 功能 | Git命令 | 说明 |
|------|---------|------|
| 检查仓库 | `git rev-parse --is-inside-work-tree` | 返回"true"或"false" |
| 获取版本 | `git describe --tags --always --dirty` | 返回标签或提交描述 |
| 获取哈希 | `git rev-parse HEAD` | 返回完整哈希 |
| 获取哈希 | `git rev-parse --short=N HEAD` | 返回缩写哈希 |
| 获取时间 | `git log -1 --format=%cd --date=iso` | 返回ISO格式时间 |
| 获取状态 | `git status --porcelain` | 返回状态列表，空表示clean |

### 3.2 架构设计

#### 模块划分
```
gm命令
├── CLI层 (internal/cli/gm.go)
│   ├── 命令定义
│   ├── 标志定义
│   └── 运行函数
└── 业务逻辑层 (internal/commands/gm/cmd_gm.go)
    ├── 配置结构体
    ├── Git命令执行器
    ├── 功能实现函数
    └── 主函数
```

#### 数据流
```
用户输入 → CLI参数解析 → 功能选择 → Git命令执行 → 结果处理 → 输出
```

### 3.3 错误处理策略

#### 错误类型
1. **Git未安装**: 检查 `git --version` 是否可用
2. **不是Git仓库**: 检查 `git rev-parse --is-inside-work-tree` 返回值
3. **Git命令执行失败**: 捕获命令执行错误
4. **参数错误**: 验证标志参数的有效性

#### 错误处理流程
```
执行前检查 → Git可用性检查 → Git仓库检查 → 命令执行 → 结果验证 → 输出
    ↓            ↓                 ↓            ↓          ↓         ↓
  提示错误    提示安装Git      提示非仓库   捕获错误   验证结果   正常输出
```

---

## 四、目录结构

### 4.1 文件结构
```
internal/
├── cli/
│   └── gm.go                    # CLI定义文件
└── commands/
    └── gm/                      # gm命令目录
        └── cmd_gm.go           # 业务逻辑文件
```

### 4.2 文件职责

#### internal/cli/gm.go
- 定义gm命令对象
- 定义所有标志变量
- 实现init()函数
- 实现runGm()函数

#### internal/commands/gm/cmd_gm.go
- 定义GmConfig配置结构体
- 实现Git命令执行相关函数
- 实现5个功能函数
- 实现GmCmdMain()主函数

---

## 五、实现步骤

### 5.1 第一步：创建目录和文件
```bash
# 创建命令目录
mkdir internal/commands/gm

# 创建业务逻辑文件
touch internal/commands/gm/cmd_gm.go

# 创建CLI定义文件
touch internal/cli/gm.go
```

### 5.2 第二步：编写业务逻辑文件
1. 定义GmConfig结构体
2. 实现Git命令执行器
3. 实现5个功能函数
4. 实现主函数GmCmdMain()

### 5.3 第三步：编写CLI定义文件
1. 定义命令变量
2. 定义标志变量（7个标志）
3. 实现init()函数
4. 实现runGm()函数

### 5.4 第四步：注册命令
在 `internal/cli/root.go` 的SubCmds列表中添加GmCmd

### 5.5 第五步：编译验证
```bash
go build ./internal/cli/...
```

### 5.6 第六步：功能测试
测试每个功能标志和组合使用场景

---

## 六、接口设计

### 6.1 配置结构体

```go
type GmConfig struct {
    Check   bool     // 检查Git仓库
    Version bool     // 获取版本号
    Hash    bool     // 获取提交哈希
    Time    bool     // 获取提交时间
    Status  bool     // 获取仓库状态
    Abbrev  int      // 哈希缩写长度
    Format  string   // 时间格式
    JSON    bool     // JSON格式输出
}
```

### 6.2 核心函数接口

#### Git命令执行器
```go
// 执行Git命令
func runGitCommand(workDir string, args ...string) (string, error)

// 检查Git是否可用
func checkGitAvailable() (bool, error)

// 检查是否为Git仓库
func checkGitRepo(workDir string) (bool, error)
```

#### 功能函数
```go
// 检查Git仓库
func checkGitRepo(config GmConfig) (string, error)

// 获取版本号
func getVersion(config GmConfig) (string, error)

// 获取提交哈希
func getHash(config GmConfig) (string, error)

// 获取提交时间
func getTime(config GmConfig) (string, error)

// 获取仓库状态
func getStatus(config GmConfig) (string, error)
```

#### 主函数
```go
// 主函数
func GmCmdMain(config GmConfig) error
```

### 6.3 CLI函数接口

```go
// 运行函数
func runGm(cmd qflag.Command) error
```

---

## 七、错误处理

### 7.1 错误类型定义

```go
// 错误类型
const (
    ErrGitNotInstalled = "git命令未找到，请确保git已安装"
    ErrNotGitRepo     = "不是git仓库"
    ErrCommandFailed   = "git命令执行失败"
    ErrInvalidFormat  = "无效的时间格式"
)
```

### 7.2 错误处理流程

#### 检查Git可用性
```go
func checkGitAvailable() (bool, error) {
    cmd := shx.New("git --version")
    _, err := cmd.ExecOutput()
    if err != nil {
        return false, fmt.Errorf(ErrGitNotInstalled)
    }
    return true, nil
}
```

#### 检查Git仓库
```go
func checkGitRepo(workDir string) (bool, error) {
    cmd := shx.New("git rev-parse --is-inside-work-tree").WithDir(workDir)
    output, err := cmd.ExecOutput()
    if err != nil {
        return false, fmt.Errorf("%s: %s", ErrNotGitRepo, workDir)
    }
    
    if strings.TrimSpace(string(output)) != "true" {
        return false, fmt.Errorf("%s: %s", ErrNotGitRepo, workDir)
    }
    
    return true, nil
}
```

#### 执行Git命令
```go
func runGitCommand(workDir string, args ...string) (string, error) {
    // 检查Git可用性
    if _, err := checkGitAvailable(); err != nil {
        return "", err
    }
    
    // 检查Git仓库
    if _, err := checkGitRepo(workDir); err != nil {
        return "", err
    }
    
    // 执行命令
    cmd := shx.NewArgs("git", args...).WithDir(workDir)
    output, err := cmd.ExecOutput()
    if err != nil {
        return "", fmt.Errorf("%s: %v", ErrCommandFailed, err)
    }
    
    return strings.TrimSpace(string(output)), nil
}
```

### 7.3 参数验证

#### 验证abbrev参数
```go
if config.Abbrev < 0 {
    return fmt.Errorf("abbrev参数必须大于等于0")
}
```

#### 验证format参数
```go
if config.Time && config.Format == "" {
    config.Format = "2006-01-02 15:04:05"
}
```

---

## 八、使用示例

### 8.1 基本使用

```bash
# 检查是否为Git仓库
fck gm --check

# 获取版本号（默认行为）
fck gm
fck gm --version

# 获取提交哈希（默认缩写7位）
fck gm --hash

# 获取完整哈希
fck gm --hash --abbrev=0

# 获取提交时间（默认格式）
fck gm --time

# 自定义时间格式
fck gm --time --format="2006/01/02 15:04"

# 获取仓库状态
fck gm --status

# JSON格式输出
fck gm --version --json
fck gm --hash --abbrev=8 --json
fck gm --time --format="2006-01-02" --json
fck gm --status --json
```

### 8.2 短选项使用

```bash
# 使用短选项
fck gm -c
fck gm -v
fck gm -h
fck gm -h -a=8
fck gm -t
fck gm -t -f="2006-01-02"
fck gm -s
```

### 8.3 组合使用

```bash
# 获取8位缩写哈希
fck gm --hash --abbrev=8

# 自定义时间格式
fck gm --time --format="2006年01月02日 15:04:05"
```

---

## 九、测试计划

### 9.1 单元测试

#### 测试用例1：检查Git仓库
- 测试在Git仓库中执行
- 测试在非Git仓库中执行
- 测试Git未安装的情况

#### 测试用例2：获取版本号
- 测试有标签的仓库
- 测试无标签的仓库
- 测试有未提交更改的仓库

#### 测试用例3：获取提交哈希
- 测试默认缩写长度
- 测试自定义缩写长度
- 测试完整哈希

#### 测试用例4：获取提交时间
- 测试默认时间格式
- 测试自定义时间格式
- 测试无效时间格式

#### 测试用例5：获取仓库状态
- 测试clean状态
- 测试dirty状态

### 9.2 集成测试

#### 测试场景1：完整流程
```bash
cd /path/to/git/repo
fck gm --check
fck gm --version
fck gm --hash
fck gm --time
fck gm --status
```

#### 测试场景2：错误处理
```bash
cd /path/to/non-git/dir
fck gm --check
fck gm --version
```

#### 测试场景3：参数组合
```bash
fck gm --hash --abbrev=8
fck gm --time --format="2006/01/02"
```

### 9.3 边界测试

#### 测试边界值
- `--abbrev=0`（完整哈希）
- `--abbrev=1`（最小缩写）
- `--abbrev=40`（最大缩写，SHA-1长度）
- 无效的`--abbrev`值（负数、非数字）
- 无效的时间格式

---

## 十、注意事项

### 10.1 跨平台兼容性
- Git命令在不同平台上的行为一致
- 路径处理使用`filepath`包
- 字符串处理使用`strings`包

### 10.2 性能考虑
- Git命令执行速度较快，无需特殊优化
- 避免重复执行相同的Git命令
- 缓存Git仓库检查结果

### 10.3 安全性考虑
- 验证用户输入的参数
- 防止命令注入攻击
- 使用shx库安全执行命令

### 10.4 用户体验
- 提供清晰的错误提示
- 支持中文帮助信息
- 提供使用示例

---

## 十一、待讨论问题

### 11.1 功能相关问题

1. **默认行为**: 如果不指定任何功能标志，是否默认执行`--version`？
   - **已确定**: 是，这样用户直接输入`fck gm`就能获取版本号

2. **互斥性**: 是否强制功能标志互斥？
   - **已确定**: 暂不强制互斥，后期根据需要添加到互斥组

3. **组合使用**: `--abbrev`和`--format`是否只能配合对应的功能标志使用？
   - **已确定**: 是，`--abbrev`只能配合`--hash`，`--format`只能配合`--time`

4. **JSON输出**: 是否支持JSON格式输出？
   - **已确定**: 是，添加`--json` / `-j`标志支持JSON格式输出

5. **检查Git仓库行为**: 检查Git仓库时，是仓库是否输出信息？
   - **已确定**: 是仓库时不输出任何信息，直接正常退出（返回码0）；不是仓库时输出错误信息并返回错误码

### 11.2 实现相关问题

1. **工作目录**: 是否需要支持指定工作目录？
   - **已确定**: 暂不支持，默认使用当前目录

2. **输出格式**: 输出是否需要支持JSON格式？
   - **已确定**: 是，添加`--json` / `-j`标志支持JSON格式输出

3. **颜色输出**: 是否需要支持彩色输出？
   - **已确定**: 否，不支持彩色输出，保持简单文本输出

### 11.3 扩展性问题

1. **未来扩展**: 是否需要支持更多Git元数据？
   - **建议**: 先实现基础功能，后续根据需求扩展

2. **批量操作**: 是否需要支持批量获取多个Git仓库的信息？
   - **建议**: 暂不支持，保持简单

### 11.4 依赖性问题

1. **shellx库**: 是否确定使用shellx库执行Git命令？
   - **建议**: 是，与示例代码保持一致

2. **其他依赖**: 是否需要引入其他依赖库？
   - **建议**: 不需要，使用标准库即可

---

## 十二、实施时间估算

### 12.1 开发任务分解

| 任务 | 预计时间 | 优先级 |
|------|----------|--------|
| 创建目录和文件 | 5分钟 | 高 |
| 编写业务逻辑文件 | 30分钟 | 高 |
| 编写CLI定义文件 | 20分钟 | 高 |
| 注册命令 | 5分钟 | 高 |
| 编译验证 | 5分钟 | 高 |
| 功能测试 | 20分钟 | 高 |
| 文档更新 | 10分钟 | 中 |
| **总计** | **95分钟** | - |

### 12.2 风险评估

| 风险 | 可能性 | 影响 | 应对措施 |
|------|--------|------|----------|
| Git命令执行失败 | 中 | 高 | 完善错误处理 |
| 跨平台兼容性问题 | 低 | 中 | 充分测试 |
| 性能问题 | 低 | 低 | 无需特殊优化 |
| 用户需求变更 | 中 | 中 | 保持架构灵活 |

---

## 十三、总结

### 13.1 方案优势
1. **功能完整**: 覆盖了5个核心Git元数据获取功能
2. **设计清晰**: 模块化设计，职责分明
3. **易于扩展**: 架构灵活，便于后续扩展
4. **用户友好**: 提供清晰的错误提示和使用示例

### 13.2 方案特点
1. **遵循规范**: 严格遵循项目的新命令开发规范
2. **技术成熟**: 使用成熟的技术栈和依赖库
3. **错误处理**: 完善的错误处理机制
4. **测试充分**: 详细的测试计划和测试用例

### 13.3 下一步行动
1. 与用户讨论本实施文档
2. 根据讨论结果调整方案
3. 开始实施代码编写
4. 进行功能测试和验证
5. 完善文档和示例

---

**文档版本**: v1.1  
**创建日期**: 2026-03-24  
**作者**: 技术架构师  
**状态**: 待实施
