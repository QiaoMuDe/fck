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

// Short 返回短格式
//
// 示例:
//
//	MyApp v1.0.0
func (i *Info) Short() string {
	return fmt.Sprintf("%s %s", i.AppName, i.GitVersion)
}

// Long 返回长格式
//
// 示例:
//
//	MyApp version v1.0.0 linux/amd64 with go1.22.1
func (i *Info) Long() string {
	return fmt.Sprintf("%s version %s %s with %s", i.AppName, i.GitVersion, i.Platform, i.GoVersion)
}

// Standard 返回标准格式
//
// 示例:
//
//	MyApp v1.0.0 (abc1234) [linux/amd64]
func (i *Info) Standard() string {
	return fmt.Sprintf("%s %s (%s) [%s]", i.AppName, i.GitVersion, i.GitCommit, i.Platform)
}

// Full 返回完整格式
//
// 示例:
//
//	MyApp v1.0.0 (abc1234) linux/amd64 built 2024-01-01 with go1.22.1
func (i *Info) Full() string {
	return fmt.Sprintf("%s %s (%s) %s built %s with %s",
		i.AppName, i.GitVersion, i.GitCommit, i.Platform, i.BuildTime, i.GoVersion)
}

// Lines 返回分行格式
//
// 示例:
//
//	MyApp v1.0.0
//	Commit: abc1234 (clean)
//	Platform: linux/amd64
//	Built: 2024-01-01 with go1.22.1
func (i *Info) Lines() string {
	return fmt.Sprintf(`%s %s
Commit: %s (%s)
Platform: %s
Built: %s with %s`,
		i.AppName, i.GitVersion, i.GitCommit, i.GitTreeState, i.Platform, i.BuildTime, i.GoVersion)
}

// Table 返回表格格式
//
// 示例:
//
//	┌─────────────┬──────────────────┐
//	│ Application │ MyApp           │
//	│ Version     │ v1.0.0          │
//	│ Commit      │ abc1234         │
//	│ Platform    │ linux/amd64     │
//	│ Build Time  │ 2024-01-01      │
//	│ Go Version  │ go1.22.1        │
//	└─────────────┴──────────────────┘
func (i *Info) Table() string {
	return fmt.Sprintf(`┌─────────────┬──────────────────┐
│ Application │ %-16s │
│ Version     │ %-16s │
│ Commit      │ %-16s │
│ Platform    │ %-16s │
│ Build Time  │ %-16s │
│ Go Version  │ %-16s │
└─────────────┴──────────────────┘`,
		i.AppName, i.GitVersion, i.GitCommit, i.Platform, i.BuildTime, i.GoVersion)
}

// JSON 返回JSON格式
//
// 示例:
//
//	{"name":"MyApp","version":"v1.0.0","commit":"abc1234","platform":"linux/amd64","buildTime":"2024-01-01T12:00:00Z","goVersion":"go1.22.1"}
func (i *Info) JSON() string {
	return fmt.Sprintf(`{"name":"%s","version":"%s","commit":"%s","platform":"%s","buildTime":"%s","goVersion":"%s"}`,
		i.AppName, i.GitVersion, i.GitCommit, i.Platform, i.GitCommitTime, i.GoVersion)
}

// KV 返回键值对格式
//
// 示例:
//
//	app=MyApp version=v1.0.0 commit=abc1234 platform=linux/amd64 build=2024-01-01 go=go1.22.1
func (i *Info) KV() string {
	return fmt.Sprintf("app=%s version=%s commit=%s platform=%s build=%s go=%s",
		i.AppName, i.GitVersion, i.GitCommit, i.Platform, i.BuildTime, i.GoVersion)
}

// Banner 返回横幅格式
//
// 示例:
//
//	╔═══════════════════════════════════════╗
//	║           MyApp v1.0.0                ║
//	╠═══════════════════════════════════════╣
//	║ Commit:  abc1234 (clean)              ║
//	║ Platform: linux/amd64                 ║
//	║ Built:   2024-01-01 with go1.22.1     ║
//	╚═══════════════════════════════════════╝
func (i *Info) Banner() string {
	return fmt.Sprintf(`╔════════════════════════════════════════════════╗
║                %s %-23s ║
╠════════════════════════════════════════════════╣
║       Commit:  %-31s ║
║       Platform: %-30s ║
║       Built:   %-31s ║
╚════════════════════════════════════════════════╝`,
		i.AppName, i.GitVersion, i.GitCommit+" ("+i.GitTreeState+")", i.Platform, i.BuildTime+" with "+i.GoVersion)
}

// CSV 返回逗号分隔格式
//
// 示例:
//
//	MyApp,v1.0.0,abc1234,linux/amd64,2024-01-01,go1.22.1
func (i *Info) CSV() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s",
		i.AppName, i.GitVersion, i.GitCommit, i.Platform, i.BuildTime, i.GoVersion)
}

// URI 返回URI格式
//
// 示例:
//
//	verman://MyApp/v1.0.0?commit=abc1234&platform=linux/amd64&build=2024-01-01&go=go1.22.1
func (i *Info) URI() string {
	return fmt.Sprintf("verman://%s/%s?commit=%s&platform=%s&build=%s&go=%s",
		i.AppName, i.GitVersion, i.GitCommit, i.Platform, i.BuildTime, i.GoVersion)
}
