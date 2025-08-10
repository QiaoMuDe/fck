# Size 模块测试文档

本目录包含了 size 模块的完整测试套件，包括单元测试、基准测试、集成测试和测试工具。

## 测试文件说明

### 1. `size_test.go` - 单元测试
包含核心函数的单元测试：
- `TestHumanReadableSize` - 测试人类可读大小格式化
- `TestExpandPath` - 测试路径展开功能
- `TestGetPathSize` - 测试路径大小计算
- `TestGetPathSizeWithHiddenFiles` - 测试隐藏文件处理
- `TestHumanReadableSizeEdgeCases` - 测试边界情况

### 2. `benchmark_test.go` - 性能测试
包含各种性能基准测试：
- `BenchmarkHumanReadableSizeVariousSizes` - 不同大小格式化性能
- `BenchmarkExpandPath` - 路径展开性能
- `BenchmarkGetPathSizeFile` - 文件大小计算性能
- `BenchmarkGetPathSizeDirectory` - 目录大小计算性能
- `BenchmarkGetPathSizeDeepDirectory` - 深层目录性能
- `BenchmarkMemoryUsage` - 内存使用测试

### 3. `integration_test.go` - 集成测试
包含完整流程的集成测试：
- `TestIntegrationSizeCmdMain` - 主函数集成测试
- `TestRealFileSystem` - 真实文件系统测试
- `TestErrorHandling` - 错误处理测试
- `TestConcurrentAccess` - 并发访问测试

### 4. `testutil_test.go` - 测试工具
提供测试辅助功能：
- `TestHelper` - 测试辅助结构体
- 文件和目录创建工具
- 断言工具
- 并发测试支持

## 运行测试

### 运行所有测试
```bash
cd commands/size
go test -v
```

### 运行单元测试（跳过耗时的集成测试）
```bash
go test -v -short
```

### 运行特定测试
```bash
# 运行特定测试函数
go test -v -run TestHumanReadableSize

# 运行特定测试文件的所有测试
go test -v -run "Test.*" ./size_test.go
```

### 运行基准测试
```bash
# 运行所有基准测试
go test -v -bench=.

# 运行特定基准测试
go test -v -bench=BenchmarkHumanReadableSize

# 运行基准测试并显示内存分配
go test -v -bench=. -benchmem

# 运行基准测试多次以获得更准确的结果
go test -v -bench=. -count=5
```

### 运行集成测试
```bash
# 运行完整的集成测试（包括耗时测试）
go test -v -run Integration

# 运行真实文件系统测试
go test -v -run TestRealFileSystem
```

### 生成测试覆盖率报告
```bash
# 生成覆盖率报告
go test -v -cover

# 生成详细的覆盖率报告
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率详情
go tool cover -func=coverage.out
```

## 测试数据和环境

### 临时文件
所有测试都使用 `t.TempDir()` 创建临时目录，测试结束后会自动清理。

### 测试文件结构
测试会创建以下类型的测试数据：
- 不同大小的文件（1KB, 10KB, 100KB, 1MB）
- 深层目录结构（10层深度）
- 大量小文件（1000个文件）
- 隐藏文件和目录
- 包含特殊字符的文件名

### 并发测试
并发测试使用以下配置：
- 10个并发goroutine
- 每个goroutine执行100次操作
- 测试文件大小约24KB

## 测试最佳实践

### 1. 测试隔离
- 每个测试使用独立的临时目录
- 测试之间不共享状态
- 使用 `t.Parallel()` 支持并行测试

### 2. 错误处理
- 所有错误都有适当的测试覆盖
- 边界条件和异常情况都有测试
- 使用表驱动测试提高测试覆盖率

### 3. 性能测试
- 基准测试包含不同规模的数据
- 内存分配测试确保没有内存泄漏
- 并发测试验证线程安全性

### 4. 可维护性
- 测试代码有清晰的注释
- 使用测试辅助工具减少重复代码
- 测试名称清晰描述测试目的

## 常见问题

### Q: 为什么有些集成测试被跳过？
A: 集成测试依赖全局变量，需要重构代码以支持依赖注入。使用 `-short` 标志可以跳过这些测试。

### Q: 如何调试失败的测试？
A: 使用 `go test -v` 查看详细输出，或在测试中添加 `t.Logf()` 输出调试信息。

### Q: 基准测试结果如何解读？
A: 基准测试输出格式为：`BenchmarkName-CPU数量 运行次数 每次操作耗时 内存分配`

### Q: 如何添加新的测试？
A: 遵循现有的测试模式，使用 `TestHelper` 创建测试数据，确保测试的独立性和可重复性。

## 持续集成

建议在CI/CD流程中运行以下测试命令：
```bash
# 快速测试（跳过耗时的集成测试）
go test -short -race -coverprofile=coverage.out

# 完整测试（包括基准测试）
go test -v -race -coverprofile=coverage.out -bench=. -benchtime=1s
```

这样可以确保代码质量和性能不会退化。