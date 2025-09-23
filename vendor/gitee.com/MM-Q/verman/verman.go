package verman

// 功能说明:
// 1. 提供编译时注入版本信息的私有变量
// 2. 通过 Info 结构体提供获取不同格式版本信息的方法
// 3. 简洁直接的API设计，通过全局实例 V 使用

/*
示例编译时注入版本信息:
go build -ldflags "-X 'gitee.com/MM-Q/verman.appName=myapp' -X 'gitee.com/MM-Q/verman.gitVersion=v1.0.0' -X 'gitee.com/MM-Q/verman.gitCommit=abc1234' -X 'gitee.com/MM-Q/verman.gitTreeState=clean' -X 'gitee.com/MM-Q/verman.gitCommitTime=2024-01-01T12:00:00Z' -X 'gitee.com/MM-Q/verman.buildTime=2024-01-01T12:00:00Z'" main.go
*/

import (
	"fmt"
	"runtime"
)

// 私有版本信息变量，在编译时注入
var (
	appName       string // 应用程序名称
	gitVersion    string // Git 语义化版本号(如 v1.0.0)
	gitCommit     string // Git 提交哈希值(如 abc1234)
	gitTreeState  string // Git 仓库状态(如 clean, dirty)
	gitCommitTime string // Git 提交时间(如 2024-01-01T12:00:00Z)
	buildTime     string // 构建时间(如 2024-01-01T12:00:00Z)
)

// Info 版本信息结构体
type Info struct {
	AppName       string // 应用程序名称
	GitVersion    string // Git 语义化版本号
	GitCommit     string // Git 提交哈希值
	GitTreeState  string // Git 仓库状态
	GitCommitTime string // Git 提交时间
	BuildTime     string // 构建时间
	GoVersion     string // Go 运行时版本
	Platform      string // 平台信息
}

// V 全局版本信息实例，供外部使用
var V *Info

// 初始化函数，设置默认值和运行时信息
func init() {
	// 设置默认值
	if appName == "" {
		appName = "unknown"
	}
	if gitVersion == "" {
		gitVersion = "unknown"
	}
	if gitCommit == "" {
		gitCommit = "unknown"
	}
	if gitTreeState == "" {
		gitTreeState = "unknown"
	}
	if gitCommitTime == "" {
		gitCommitTime = "unknown"
	}
	if buildTime == "" {
		buildTime = "unknown"
	}

	// 创建全局实例
	V = &Info{
		AppName:       appName,                                            // 应用程序名称
		GitVersion:    gitVersion,                                         // Git 语义化版本号
		GitCommit:     gitCommit,                                          // Git 提交哈希值
		GitTreeState:  gitTreeState,                                       // Git 仓库状态
		GitCommitTime: gitCommitTime,                                      // Git 提交时间
		BuildTime:     buildTime,                                          // 构建时间
		GoVersion:     runtime.Version(),                                  // Go 运行时版本
		Platform:      fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH), // 平台信息
	}
}

// Version 返回格式为"程序名 version 版本号 平台/架构"的字符串
//
// 示例:
//
//	MyApp version v1.0.0 linux/amd64
func (i *Info) Version() string {
	return fmt.Sprintf("%s version %s %s", i.AppName, i.GitVersion, i.Platform)
}

// Simple 返回格式为"程序名 v1.0.0"的字符串
//
// 示例:
//
//	MyApp v1.0.0
func (i *Info) Simple() string {
	return fmt.Sprintf("%s %s", i.AppName, i.GitVersion)
}

// Full 返回格式为"程序名 version 版本号 平台/架构 (commit: abc1234)"的字符串
//
// 示例:
//
//	MyApp version v1.0.0 linux/amd64 (commit: abc1234)
func (i *Info) Full() string {
	return fmt.Sprintf("%s version %s %s (commit: %s)", i.AppName, i.GitVersion, i.Platform, i.GitCommit)
}

// Detail 返回格式为"程序名 v1.0.0 linux/amd64 built at 2024-01-01"的字符串
//
// 示例:
//
//	MyApp v1.0.0 linux/amd64 built at 2024-01-01
func (i *Info) Detail() string {
	return fmt.Sprintf("%s %s %s built at %s", i.AppName, i.GitVersion, i.Platform, i.BuildTime)
}

// Complete 返回包含所有信息的完整字符串
//
// 示例:
//
//	MyApp v1.0.0 linux/amd64 (commit: abc1234, tree: clean, built: 2024-01-01T12:00:00Z, go: go1.19)"
func (i *Info) Complete() string {
	return fmt.Sprintf("%s %s %s (commit: %s, tree: %s, built: %s, go: %s)",
		i.AppName, i.GitVersion, i.Platform, i.GitCommit, i.GitTreeState, i.BuildTime, i.GoVersion)
}

// Banner 返回横幅格式(多行)
//
// 示例:
//
//	MyApp v2.1.0
//	Platform: linux/amd64 | Go: go1.22.1
func (i *Info) Banner() string {
	return fmt.Sprintf(`%s %s
Platform: %s | Go: %s`, i.AppName, i.GitVersion, i.Platform, i.GoVersion)
}

// Table 返回表格格式(多行)
//
// 示例:
//
//	Application : MyApp
//	Version     : v2.1.0
//	Platform    : linux/amd64
//	Commit      : a1b2c3d4e5f6
//	Tree State  : clean
//	Build Time  : 2024-03-15T15:00:00Z
//	Go Version  : go1.22.1
func (i *Info) Table() string {
	return fmt.Sprintf(`Application : %s
Version     : %s
Platform    : %s
Commit      : %s
Tree State  : %s
Build Time  : %s
Go Version  : %s`,
		i.AppName, i.GitVersion, i.Platform, i.GitCommit, i.GitTreeState, i.BuildTime, i.GoVersion)
}

// Build 返回构建信息格式
//
// 示例:
//
//	MyApp v2.1.0
//	Built at 2024-03-15T15:00:00Z with go1.22.1
func (i *Info) Build() string {
	return fmt.Sprintf("%s %s\nBuilt at %s with %s", i.AppName, i.GitVersion, i.BuildTime, i.GoVersion)
}

// Git 返回Git信息格式
//
// 示例:
//
//	Version: v2.1.0
//	Commit: a1b2c3d4e5f6 (clean)
//	Commit Time: 2024-03-15T14:30:00Z
func (i *Info) Git() string {
	return fmt.Sprintf(`Version: %s
Commit: %s (%s)
Commit Time: %s`,
		i.GitVersion, i.GitCommit, i.GitTreeState, i.GitCommitTime)
}

// JSON 返回JSON格式
//
// 示例:
//
//	 {
//		  "appName": "MyApp",
//		  "gitVersion": "v2.1.0",
//		  "gitCommit": "a1b2c3d4e5f6",
//		  "gitTreeState": "clean",
//		  "gitCommitTime": "2024-03-15T14:30:00Z",
//		  "buildTime": "2024-03-15T15:00:00Z",
//		  "goVersion": "go1.22.1",
//		  "platform": "linux/amd64"
//	 }
func (i *Info) JSON() string {
	return fmt.Sprintf(`{
  "appName": "%s",
  "gitVersion": "%s",
  "gitCommit": "%s",
  "gitTreeState": "%s",
  "gitCommitTime": "%s",
  "buildTime": "%s",
  "goVersion": "%s",
  "platform": "%s"
}`, i.AppName, i.GitVersion, i.GitCommit, i.GitTreeState, i.GitCommitTime, i.BuildTime, i.GoVersion, i.Platform)
}
