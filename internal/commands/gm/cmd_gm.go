package gm

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitee.com/MM-Q/shellx/shx"
)

// GmConfig gm命令配置结构体
type GmConfig struct {
	WorkDir string // 工作目录
	Check   bool   // 检查Git仓库
	Version bool   // 获取版本号
	Hash    bool   // 获取提交哈希
	Time    bool   // 获取提交时间
	Status  bool   // 获取仓库状态
	Abbrev  int    // 哈希缩写长度
	Format  string // 时间格式
	JSON    bool   // JSON格式输出
}

// GmCmdMain 执行gm命令
// 参数:
//   - config: gm命令配置
//
// 返回值:
//   - error: 执行错误
func GmCmdMain(config GmConfig) error {
	// 如果工作目录为空，使用当前目录
	if config.WorkDir == "" {
		config.WorkDir = "."
	}

	// 执行对应功能标志
	var err error

	switch {
	case config.Check:
		// 检查Git仓库
		_, err = checkGitRepo(config.WorkDir)
		if err != nil {
			return err
		}
		// 是Git仓库，不输出任何信息，直接正常退出
		return nil

	case config.Version:
		// 获取版本号
		result, err := getVersion(config.WorkDir, config)
		if err != nil {
			return err
		}
		fmt.Println(result)

	case config.Hash:
		// 获取提交哈希
		result, err := getHash(config.WorkDir, config)
		if err != nil {
			return err
		}
		fmt.Println(result)

	case config.Time:
		// 获取提交时间
		result, err := getTime(config.WorkDir, config)
		if err != nil {
			return err
		}
		fmt.Println(result)

	case config.Status:
		// 获取仓库状态
		result, err := getStatus(config.WorkDir, config)
		if err != nil {
			return err
		}
		fmt.Println(result)

	default:
		// 默认行为：获取所有元数据
		return getAllMetadata(config)
	}

	return nil
}

// getAllMetadata 获取所有Git元数据
// 参数:
//   - config: gm命令配置
//
// 返回值:
//   - error: 执行错误
func getAllMetadata(config GmConfig) error {
	// 获取所有元数据
	version, err := getVersion(config.WorkDir, config)
	if err != nil {
		return err
	}

	hash, err := getHash(config.WorkDir, config)
	if err != nil {
		return err
	}

	time, err := getTime(config.WorkDir, config)
	if err != nil {
		return err
	}

	status, err := getStatus(config.WorkDir, config)
	if err != nil {
		return err
	}

	// 输出结果
	if config.JSON {
		// JSON格式输出：合并所有信息到一个JSON对象
		var versionData, hashData, timeData, statusData map[string]interface{}
		if err := json.Unmarshal([]byte(version), &versionData); err != nil {
			return fmt.Errorf("failed to parse version JSON: %v", err)
		}
		if err := json.Unmarshal([]byte(hash), &hashData); err != nil {
			return fmt.Errorf("failed to parse hash JSON: %v", err)
		}
		if err := json.Unmarshal([]byte(time), &timeData); err != nil {
			return fmt.Errorf("failed to parse time JSON: %v", err)
		}
		if err := json.Unmarshal([]byte(status), &statusData); err != nil {
			return fmt.Errorf("failed to parse status JSON: %v", err)
		}

		// 合并所有数据
		allData := make(map[string]interface{})
		for k, v := range versionData {
			allData[k] = v
		}
		for k, v := range hashData {
			allData[k] = v
		}
		for k, v := range timeData {
			allData[k] = v
		}
		for k, v := range statusData {
			allData[k] = v
		}

		jsonData, err := json.MarshalIndent(allData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	} else {
		// 普通文本输出
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Hash: %s\n", hash)
		fmt.Printf("Time: %s\n", time)
		fmt.Printf("Status: %s\n", status)
	}

	return nil
}

// checkGitAvailable 检查git命令是否可用
// 返回值:
//   - bool: git命令是否可用
//   - error: 错误信息
func checkGitAvailable() (bool, error) {
	cmd := shx.New("git --version")
	_, err := cmd.ExecOutput()
	if err != nil {
		return false, fmt.Errorf("failed to execute git --version: %v", err)
	}
	return true, nil
}

// checkGitRepo 检查是否是git仓库
// 参数:
//   - workDir: 工作目录
//
// 返回值:
//   - bool: 是否是git仓库
//   - error: 错误信息
func checkGitRepo(workDir string) (bool, error) {
	cmd := shx.New("git rev-parse --is-inside-work-tree").WithDir(workDir)
	output, err := cmd.ExecOutput()
	if err != nil {
		return false, fmt.Errorf("not a git repository: %s", workDir)
	}

	if strings.TrimSpace(string(output)) != "true" {
		return false, fmt.Errorf("not a git repository: %s", workDir)
	}

	return true, nil
}

// runGitCommand 执行git命令
// 参数:
//   - workDir: 工作目录
//   - args: git命令参数
//
// 返回值:
//   - string: 命令输出
//   - error: 执行错误
func runGitCommand(workDir string, args ...string) (string, error) {
	// 检查git是否可用
	if _, err := checkGitAvailable(); err != nil {
		return "", err
	}

	// 检查是否是git仓库
	if _, err := checkGitRepo(workDir); err != nil {
		return "", err
	}

	// 执行命令
	cmd := shx.NewArgs("git", args...).WithDir(workDir)
	output, err := cmd.ExecOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute git command: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// getVersion 获取版本号
// 参数:
//   - workDir: 工作目录
//   - config: gm命令配置
//
// 返回值:
//   - string: 版本号（JSON格式或普通格式）
//   - error: 执行错误
func getVersion(workDir string, config GmConfig) (string, error) {
	version, err := runGitCommand(workDir, "describe", "--tags", "--always", "--dirty")
	if err != nil {
		return "", fmt.Errorf("failed to get version info: %v", err)
	}

	if config.JSON {
		// JSON格式输出
		result := map[string]string{"version": version}
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON version: %v", err)
		}
		return string(jsonData), nil
	}

	return version, nil
}

// getHash 获取提交哈希
// 参数:
//   - workDir: 工作目录
//   - config: gm命令配置
//
// 返回值:
//   - string: 提交哈希（JSON格式或普通格式）
//   - error: 执行错误
func getHash(workDir string, config GmConfig) (string, error) {
	var hash string
	var err error

	if config.Abbrev == 0 {
		// 获取完整哈希
		hash, err = runGitCommand(workDir, "rev-parse", "HEAD")
	} else {
		// 获取缩写哈希
		hash, err = runGitCommand(workDir, "rev-parse", fmt.Sprintf("--short=%d", config.Abbrev), "HEAD")
	}

	if err != nil {
		return "", fmt.Errorf("failed to get hash: %v", err)
	}

	if config.JSON {
		// JSON格式输出
		result := map[string]string{"hash": hash}
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON hash: %v", err)
		}
		return string(jsonData), nil
	}

	return hash, nil
}

// getTime 获取提交时间
// 参数:
//   - workDir: 工作目录
//   - config: gm命令配置
//
// 返回值:
//   - string: 提交时间（JSON格式或普通格式）
//   - error: 执行错误
func getTime(workDir string, config GmConfig) (string, error) {
	// 获取时间格式
	timeFormat := config.Format
	if timeFormat == "" {
		timeFormat = "2006-01-02 15:04:05"
	}

	// 执行git log命令获取提交时间
	gitTime, err := runGitCommand(workDir, "log", "-1", "--format=%cd", "--date=iso")
	if err != nil {
		return "", fmt.Errorf("failed to get time info: %v", err)
	}

	// 解析git时间格式
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700", gitTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse time info: %v", err)
	}

	// 格式化为用户指定的格式
	formattedTime := parsedTime.Format(timeFormat)

	if config.JSON {
		// JSON格式输出
		result := map[string]string{"time": formattedTime}
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON time: %v", err)
		}
		return string(jsonData), nil
	}

	return formattedTime, nil
}

// getStatus 获取仓库状态
// 参数:
//   - workDir: 工作目录
//   - config: gm命令配置
//
// 返回值:
//   - string: 仓库状态（JSON格式或普通格式）
//   - error: 执行错误
func getStatus(workDir string, config GmConfig) (string, error) {
	status, err := runGitCommand(workDir, "status", "--porcelain")
	if err != nil {
		return "", fmt.Errorf("failed to get status info: %v", err)
	}

	// 判断仓库状态
	var repoStatus string
	if status == "" {
		repoStatus = "clean"
	} else {
		repoStatus = "dirty"
	}

	if config.JSON {
		// JSON格式输出
		result := map[string]string{"status": repoStatus}
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON status: %v", err)
		}
		return string(jsonData), nil
	}

	return repoStatus, nil
}
