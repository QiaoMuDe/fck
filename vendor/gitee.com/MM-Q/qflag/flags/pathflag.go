package flags

import (
	"path/filepath"

	"sync"

	"gitee.com/MM-Q/qflag/qerr"
	"gitee.com/MM-Q/qflag/validator"
)

// PathFlag 路径类型标志结构体
// 继承BaseFlag[string]泛型结构体,实现Flag接口
type PathFlag struct {
	BaseFlag[string]
	mu sync.Mutex // 保护validator并发访问
}

// Type 返回标志类型
func (f *PathFlag) Type() FlagType { return FlagTypePath }

// String 实现flag.Value接口,返回当前值的字符串表示
func (f *PathFlag) String() string { return f.Get() }

// Set 实现flag.Value接口,解析并验证路径
//
// 参数:
//   - value: 待解析的路径值
//
// 返回值:
//   - error: 解析错误或验证错误
func (f *PathFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if value == "" {
		return qerr.NewValidationError("path cannot be empty")
	}

	// 规范化路径为绝对路径
	absPath, err := filepath.Abs(value)
	if err != nil {
		return qerr.NewValidationErrorf("failed to get absolute path: %v", err)
	}

	// 调用基类方法设置值(会触发验证器验证)
	return f.BaseFlag.Set(absPath)
}

// Init 初始化路径标志
//
// 参数:
//   - longName: 长名称
//   - shortName: 短名称
//   - defValue: 默认值
//   - usage: 使用说明
//
// 返回值:
//   - error: 初始化错误
func (f *PathFlag) Init(longName, shortName string, defValue string, usage string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 初始化路径标志值指针
	valuePtr := new(string)

	// 规范化默认路径为绝对路径
	absDefValue, err := filepath.Abs(defValue)
	if err != nil {
		return qerr.NewValidationErrorf("failed to normalize default path: %v", err)
	}

	// 设置默认值
	*valuePtr = absDefValue

	// 调用基类方法初始化
	if err := f.BaseFlag.Init(longName, shortName, usage, valuePtr); err != nil {
		return err
	}

	// 设置路径验证器
	f.SetValidator(&validator.PathValidator{
		MustExist: true, // 默认必须存在
	})
	return nil
}

// MustExist 设置路径是否必须存在
//
// 参数值:
//   - mustExist: 是否必须存在
//
// 返回值:
//   - *PathFlag: 当前路径标志对象
//
// 示例:
// cmd.Path("output", "o", "/tmp/output", "输出目录").MustExist(false)
func (f *PathFlag) MustExist(mustExist bool) *PathFlag {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.validator.(*validator.PathValidator); ok {
		v.MustExist = mustExist
	}
	return f
}

// IsDirectory 设置路径是否必须是目录
//
// 参数值:
//   - isDir: 是否必须是目录
//
// 返回值:
//   - *PathFlag: 当前路径标志对象
//
// 示例:
// cmd.Path("log-dir", "l", "/var/log/app", "日志目录").IsDirectory(true)
func (f *PathFlag) IsDirectory(isDir bool) *PathFlag {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.validator.(*validator.PathValidator); ok {
		v.IsDirectory = isDir
	}
	return f
}
