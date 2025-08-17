package find

import (
	"testing"
)

// initTestFlags 初始化测试用的标志变量
// 这个函数在每个测试开始前调用，确保标志变量不为nil
func initTestFlags() {
	// 如果标志变量为nil，则初始化find命令
	if findCmd == nil {
		findCmd = InitFindCmd()
	}
}

// TestMain 在所有测试运行前执行初始化
func TestMain(m *testing.M) {
	// 初始化find命令和标志
	findCmd = InitFindCmd()

	// 运行测试
	m.Run()
}
