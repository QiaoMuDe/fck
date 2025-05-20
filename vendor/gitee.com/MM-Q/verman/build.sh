#!/bin/bash

#########################################################################
# 设置输出文件名
OUTPUT_FILE="输出文件名"
# 设置项目名
PROJECT_NAME="项目名"
# 入口文件
ENTRY_FILE="main.go" # 入口文件位置
# 是否使用vendor模式构建
USE_VENDOR_IN_BUILD=false
#########################################################################

# 检查go编译器是否安装
echo "正在检查go编译器是否安装..."
if ! command -v go &> /dev/null; then
    echo "错误: go编译器未安装，请先安装go编译器。"
    exit 1
fi

# 检查入口文件是否存在
echo "正在检查入口文件是否存在..."
if [ ! -f "$ENTRY_FILE" ]; then
    echo "错误: 入口文件 $ENTRY_FILE 不存在，请修改 ENTRY_FILE 变量。"
    exit 1
fi

# 检查输出文件名是否为空
echo "正在检查输出文件名是否为空..."
if [ -z "$OUTPUT_FILE" ]; then
    echo "错误: 输出文件名为空，请修改 OUTPUT_FILE 变量"
    exit 1
fi

# 通过 go mod tidy 检查依赖
echo "正在检查依赖..."
go mod tidy > /tmp/tidy.log 2>&1
if [ $? -ne 0 ]; then
    echo "错误: 依赖检查失败，请修复依赖问题。"
    cat /tmp/tidy.log
    exit 1
fi
rm -f /tmp/tidy.log

# 通过 go vet 检查代码
echo "正在检查代码..."
go vet ${ENTRY_FILE} > /tmp/vet.log 2>&1
if [ $? -ne 0 ]; then
    echo "错误: 代码检查失败，请修复代码中的错误。"
    cat /tmp/vet.log
    exit 1
fi
rm -f /tmp/vet.log

# 通过 gofmt 格式化代码
echo "正在格式化代码..."
gofmt -w . > /tmp/gofmt.log 2>&1
if [ $? -ne 0 ]; then
    echo "错误: 代码格式化失败，请修复代码中的错误。"
    cat /tmp/gofmt.log
    exit 1
fi
rm -f /tmp/gofmt.log

# 获取 Git 版本信息
echo "正在获取 Git 版本信息..."
# 获取 Git 版本
GIT_VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")
# 获取 Git 提交哈希
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
# 获取git 提交时间
GIT_COMMIT_TIME=$(git log -1 --format=%cd --date=iso 2>/dev/null || echo "unknown")
# 使用 date 命令转换为所需的格式
FORMAT_TIME=$(date -u -d "$GIT_COMMIT_TIME" +"%Y-%m-%dT%H:%M:%SZ")
# 获取仓库状态
GIT_STATUS=$(git status --porcelain 2>/dev/null | grep -q . && echo "dirty" || echo "clean")

# 检查 Git 版本信息是否获取成功
if [ "$GIT_VERSION" == "unknown" ] || [ "$GIT_COMMIT" == "unknown" ] || [ "$GIT_COMMIT_TIME" == "unknown" ]; then
    echo "警告: 无法获取 Git 版本信息，可能是当前目录不是 Git 仓库。"
    exit 1
fi

# 获取构建时间
echo "正在获取构建时间..."
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# 构建程序
echo "正在构建程序..."
# 编译注入的参数
LD_FLAGS="-X 'gitee.com/MM-Q/verman.appName=${PROJECT_NAME}' -X 'gitee.com/MM-Q/verman.gitVersion=${GIT_VERSION}' -X 'gitee.com/MM-Q/verman.gitCommit=${GIT_COMMIT}' -X 'gitee.com/MM-Q/verman.gitCommitTime=${FORMAT_TIME}' -X 'gitee.com/MM-Q/verman.buildTime=${BUILD_TIME}' -X 'gitee.com/MM-Q/verman.gitTreeState=${GIT_STATUS}' -s -w"

# 编译程序
BUILD_CMD="go build -ldflags '"${LD_FLAGS}"'"
if [ "$USE_VENDOR_IN_BUILD" = true ]; then
    BUILD_CMD="${BUILD_CMD} -mod=vendor"
fi
build_status=$(eval "${BUILD_CMD} -o ${OUTPUT_FILE} ${ENTRY_FILE}" > /tmp/build.log 2>&1; echo $?)
if [ $build_status -eq 1 ]; then
    echo "错误: 程序构建失败，请检查错误信息。"
    cat /tmp/build.log
    exit 1
fi

echo "----------------------------------------------------------"
echo "${PROJECT_NAME} 构建成功，输出文件为 ${OUTPUT_FILE}。"
echo "版本信息:"
echo "Git 仓库状态: $GIT_STATUS"
echo "Git 版本: $GIT_VERSION"
echo "Git 提交: $GIT_COMMIT"
echo "Git 提交时间: $FORMAT_TIME"
echo "构建时间: $BUILD_TIME"
echo "----------------------------------------------------------"

# 清理临时文件
rm -f /tmp/build.log