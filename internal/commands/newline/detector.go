// Package newline 提供文件换行符检测和转换功能
package newline

import (
	"bytes"
	"fmt"

	"gitee.com/MM-Q/fck/internal/types"
)

// DetectionResult 检测结果
type DetectionResult struct {
	Type      string
	LFCount   int
	CRLFCount int
	CRCount   int
}

// Detect 检测数据中的换行符类型
//
// 参数:
//   - data: 要检测的数据
//
// 返回值:
//   - DetectionResult: 检测结果
func Detect(data []byte) DetectionResult {
	lfCount := bytes.Count(data, []byte{'\n'})
	crlfCount := bytes.Count(data, []byte("\r\n"))
	crCount := bytes.Count(data, []byte{'\r'})

	// 纯 CRLF: 所有 \r 都是 \r\n 的一部分，且 \n 数量等于 \r\n 数量
	if crlfCount > 0 && lfCount == crlfCount && crCount == crlfCount {
		return DetectionResult{
			Type:      types.NewlineCRLF,
			CRLFCount: crlfCount,
		}
	}

	// 纯 LF: 有 \n 但没有 \r
	if lfCount > 0 && crlfCount == 0 && crCount == 0 {
		return DetectionResult{
			Type:    types.NewlineLF,
			LFCount: lfCount,
		}
	}

	// 纯 CR: 有 \r 但没有 \n
	if crCount > 0 && lfCount == 0 && crlfCount == 0 {
		return DetectionResult{
			Type:    types.NewlineCR,
			CRCount: crCount,
		}
	}

	// 混合: 存在多种换行符
	if lfCount > 0 || crCount > 0 {
		return DetectionResult{
			Type:      types.NewlineMixed,
			LFCount:   lfCount - crlfCount,
			CRLFCount: crlfCount,
			CRCount:   crCount - crlfCount,
		}
	}

	// 无换行符
	return DetectionResult{Type: types.NewlineNone}
}

// String 返回检测结果的字符串表示
func (r DetectionResult) String() string {
	switch r.Type {
	case types.NewlineLF:
		return fmt.Sprintf("%s (%d lines)", r.Type, r.LFCount)
	case types.NewlineCRLF:
		return fmt.Sprintf("%s (%d lines)", r.Type, r.CRLFCount)
	case types.NewlineCR:
		return fmt.Sprintf("%s (%d lines)", r.Type, r.CRCount)
	case types.NewlineMixed:
		return fmt.Sprintf("%s (LF: %d, CRLF: %d, CR: %d)", r.Type, r.LFCount, r.CRLFCount, r.CRCount)
	default:
		return r.Type
	}
}

// TotalLines 返回总行数
func (r DetectionResult) TotalLines() int {
	return r.LFCount + r.CRLFCount + r.CRCount
}
