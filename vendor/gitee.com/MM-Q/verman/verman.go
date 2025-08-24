package verman

// 功能说明:
// 1. 定义版本信息结构体
// 2. 提供获取版本信息的方法
// 3. 提供获取应用程序名称的方法
// 4. 提供获取应用程序版本信息的方法
// 5. 提供获取应用程序版本信息的 JSON 字符串的方法

/*
示例编译时注入版本信息:
go build -ldflags "-X 'github.com/your-username/your-repo/verman.gitVersion=v1.0.0' -X 'github.com/your-username/your-repo/verman.gitCommit=abc1234' -X 'github.com/your-username/your-repo/verman.gitTreeState=clean' -X 'github.com/your-username/your-repo/verman.gitCommitTime=2024-01-01T12:00:00Z' -X 'github.com/your-username/your-repo/verman.buildTime=2024-01-01T12:00:00Z'" main.go
*/

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// 应用程序名称, 无需修改，会在编译时自动注入
var appName = ""

// 定义用于在编译时注入的版本信息
var (
	gitVersion    string // 语义化版本号(如 v1.0.0)
	gitCommit     string // Git 提交哈希值(如 abc1234)
	gitTreeState  string // Git 仓库状态(如 clean, dirty)
	gitCommitTime string // Git 提交时间(如 2024-01-01T12:00:00Z)
	buildTime     string // 构建时间(如 2024-01-01T12:00:00Z)
)

// VerMan 包含版本信息的结构体
type VerMan struct {
	AppName       string `json:"appName"`       // 应用程序名称
	GitVersion    string `json:"gitVersion"`    // Git 语义化版本号(如 v1.0.0)
	GitCommit     string `json:"gitCommit"`     // Git 提交哈希值(如 abc1234)
	GitTreeState  string `json:"gitTreeState"`  // Git 仓库状态(如 clean, dirty)
	GitCommitTime string `json:"gitCommitTime"` // Git 提交时间(如 2024-01-01T12:00:00Z)
	BuildTime     string `json:"buildTime"`     // 构建时间(如 2024-01-01T12:00:00Z)
	GoVersion     string `json:"goVersion"`     // Go 运行时版本(如 go1.19)
	Platform      string `json:"platform"`      // 平台信息(如 linux/amd64)
}

// 定义一个全局变量，用于存储应用程序版本信息
var Info = VerMan{}

// 初始化函数，在导入包时自动赋值版本信息
func init() {
	// 如果 Git 语义化版本号为空，则返回默认值
	if gitVersion == "" {
		gitVersion = "unknown"
	}

	// 如果 Git 提交哈希值为空，则返回默认值
	if gitCommit == "" {
		gitCommit = "unknown"
	}

	// 如果 Git 仓库状态为空，则返回默认值
	if gitTreeState == "" {
		gitTreeState = "unknown"
	}

	// 如果 Git 提交时间为空，则返回默认值
	if gitCommitTime == "" {
		gitCommitTime = "unknown"
	}

	// 如果构建时间为空，则返回默认值
	if buildTime == "" {
		buildTime = "unknown"
	}

	// 将版本信息存储到全局变量中
	Info = VerMan{
		AppName:       GetAppName(),                                  // 应用程序名称
		GitVersion:    gitVersion,                                    // Git 语义化版本号
		GitCommit:     gitCommit,                                     // Git 提交哈希值
		GitTreeState:  gitTreeState,                                  // Git 仓库状态
		GitCommitTime: gitCommitTime,                                 // Git 提交时间
		BuildTime:     buildTime,                                     // 构建时间
		GoVersion:     runtime.Version(),                             // Go 运行时版本
		Platform:      fmt.Sprint(runtime.GOOS, "/", runtime.GOARCH), // 平台信息(平台/架构)
	}
}

// GetAppName 返回应用程序的名称
func GetAppName() string {
	// 如果应用程序名称为空，则返回默认值
	if appName == "" {
		appName = "unknown"
	}
	return appName
}

// Get 返回包含当前应用程序版本信息的 VersionInfo 结构体。
//
// 返回值：
//   - *VersionInfo：包含当前应用程序版本信息的 VersionInfo 结构体。
func Get() *VerMan {
	return &Info // 返回全局变量中的版本信息
}

// String 返回版本信息的字符串表示
func (v VerMan) String() string {
	return v.GitVersion
}

// 定义一个名为JSON的方法，该方法属于VersionInfo类型
//
// 返回值：
//   - string：JSON格式的字符串表示
//   - error：转换错误信息
func (v VerMan) JSON() (string, error) {
	// 使用json.MarshalIndent将VersionInfo类型的v转换为JSON格式的字节切片
	// 第二个参数""表示在JSON输出的每个键之前不添加任何前缀
	// 第三个参数"  "表示在JSON输出的每个级别之间添加两个空格的缩进
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// 如果转换过程中发生错误，返回空字符串和错误信息
		return "", err
	}

	// 如果转换成功，将字节切片转换为字符串并返回，同时返回nil表示没有错误
	return string(data), nil
}

// PrintVersion 根据指定的格式打印版本信息。
//
// 参数:
//   - format: 输出格式，支持 "json", "text" 和 "simple"。
//
// 支持的格式包括:
//   - "json": 以 JSON 格式输出版本信息
//   - "text": 以完整详细格式输出版本信息
//   - "simple": 以简洁格式输出版本信息
//
// 注意:
//   - 如果传入的格式不是支持的格式，将打印错误信息。
func (v VerMan) PrintVersion(format string) {
	switch format {
	// JSON 格式
	case "json":
		jsonStr, err := v.JSON()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(jsonStr)
	// 完整详细格式
	case "text":
		fmt.Printf("%-15s %s\n", "Application:", v.AppName)       // 打印应用程序名称和版本号
		fmt.Printf("%-15s %s\n", "Version:", v.GitVersion)        // 打印 Git 语义化版本号
		fmt.Printf("%-15s %s\n", "Git Commit:", v.GitCommit)      // 打印 Git 提交哈希值
		fmt.Printf("%-15s %s\n", "Commit Time:", v.GitCommitTime) // 打印 Git 提交时间
		fmt.Printf("%-15s %s\n", "Build Time:", v.BuildTime)      // 打印构建时间
		fmt.Printf("%-15s %s\n", "Go Version:", v.GoVersion)      // 打印 Go 运行时版本
		fmt.Printf("%-15s %s\n", "Platform:", v.Platform)         // 打印平台信息
	// 简洁格式
	case "simple":
		fmt.Printf("%s Version: %s, Git Commit: %s, Build Time: %s\n", v.AppName, v.GitVersion, v.GitCommit, v.BuildTime) // 打印应用程序名称和版本号
	default:
		// 如果传入的格式不是支持的格式，打印错误信息
		fmt.Printf("Unknown format: %s\n", format)
	}
}

// SprintVersion 根据指定的格式返回版本信息。
// 该方法接受一个字符串参数 format，用于指定返回格式。
//
// 参数:
//   - format: 返回格式，支持 "json", "text" 和 "simple"。
//
// 支持的格式包括:
//   - "json": 以 JSON 格式返回版本信息
//   - "text": 以完整详细格式返回版本信息
//   - "simple": 以简洁格式返回版本信息
//
// 注意:
//   - 如果传入的格式不是支持的格式，将返回错误信息。
func (v VerMan) SprintVersion(format string) (string, error) {
	// 根据传入的格式参数进行不同的处理
	switch format {
	// JSON 格式
	case "json":
		// 调用 JSON 方法将 VersionInfo 结构体转换为 JSON 字符串
		jsonStr, err := v.JSON()
		// 检查转换过程中是否发生错误
		if err != nil {
			return "", err
		}
		// 如果转换成功，返回 JSON 字符串
		return jsonStr, nil
	// 完整详细格式
	case "text":
		return fmt.Sprintf("%-15s %s\n"+
			"%-15s %s\n"+
			"%-15s %s\n"+
			"%-15s %s\n"+
			"%-15s %s\n"+
			"%-15s %s\n"+
			"%-15s %s",
			"Application:", v.AppName,
			"Version:", v.GitVersion,
			"Git Commit:", v.GitCommit,
			"Commit Time:", v.GitCommitTime,
			"Build Time:", v.BuildTime,
			"Go Version:", v.GoVersion,
			"Platform:", v.Platform), nil
	// 简洁格式
	case "simple":
		simpleInfo := fmt.Sprintf("%s Version: %s, Git Commit: %s, Build Time: %s", v.AppName, v.GitVersion, v.GitCommit, v.BuildTime) // 将版本信息格式化为字符串
		return simpleInfo, nil
	default:
		// 如果传入的格式不是支持的格式，返回错误信息
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
