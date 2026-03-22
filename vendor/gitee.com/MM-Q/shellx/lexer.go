package shellx

import (
	"strings"
	"unicode"
)

// splitState 命令拆分过程中的状态信息
//
// 封装了命令拆分过程中需要的所有状态变量，
// 便于状态管理和参数传递
type splitState struct {
	result         []string        // 拆分结果
	builder        strings.Builder // 当前命令片段构建器
	inQuotes       bool            // 是否在引号中
	quote          rune            // 当前引号类型
	hasQuoteInWord bool            // 当前片段是否包含过引号(用于处理空引号情况)
	emptyQuote     bool            // 当前引号是否为空(用于区分空引号和非空引号）
}

// newSplitState 创建新的拆分状态
//
// 返回值:
//   - *splitState: 初始化后的拆分状态指针
func newSplitState() *splitState {
	return &splitState{
		result:         make([]string, 0, 8),
		builder:        strings.Builder{},
		inQuotes:       false,
		quote:          0,
		hasQuoteInWord: false,
		emptyQuote:     false,
	}
}

// isQuote 判断字符是否为引号（单引号、双引号、反引号）
//
// 支持的引号类型:
//   - 双引号: "
//   - 单引号: '
//   - 反引号: `
//
// 参数:
//   - ch: 要判断的字符
//
// 返回值:
//   - bool: 如果是引号返回 true，否则返回 false
func isQuote(ch rune) bool {
	switch ch {
	case '"', '\'', '`':
		return true
	default:
		return false
	}
}

// isSpecialChar 判断字符是否为shell特殊字符
//
// 支持的特殊字符:
//   - ; : 命令分隔符
//   - | : 管道符
//   - > : 重定向符
//   - < : 重定向符
//   - & : 后台运行/逻辑运算
//   - ! : 逻辑非
//   - ` : 命令替换
//   - * : 通配符
//   - ? : 通配符
//   - # : 注释符
//   - ~ : 家目录
//
// 参数:
//   - ch: 要判断的字符
//
// 返回值:
//   - bool: 如果是特殊字符返回 true，否则返回 false
func isSpecialChar(ch rune) bool {
	switch ch {
	case ';', '|', '>', '<', '&', '!', '`', '*', '?', '#', '~':
		return true

	default:
		return false
	}
}

// handleSpecialChar 处理特殊字符（分号、管道符等）的逻辑
//
// 在非引号状态下，特殊字符作为独立token处理
// 在引号内，特殊字符作为普通字符处理
//
// 参数:
//   - state: 拆分状态
//   - ch: 当前字符
func (s *splitState) handleSpecialChar(ch rune) {
	// 根据是否在引号内采用不同的处理策略
	if s.inQuotes {
		// 引号内：特殊字符作为普通字符处理
		s.builder.WriteRune(ch)
	} else {
		// 引号外：特殊字符作为独立token处理

		// 1. 先处理当前累积的token
		s.flushBuilder()

		// 2. 将特殊字符作为独立token添加到结果中
		s.result = append(s.result, string(ch))
	}

	// 3. 重置空引号状态标记
	s.emptyQuote = false
}

// flushBuilder 将 builder 中的内容添加到结果中并重置
//
// 如果 builder 中有内容，则将其作为独立的 token 添加到结果中，
// 然后重置 builder，准备下一个 token 的构建
//
// 参数:
//   - 无
//
// 返回值:
//   - 无
func (s *splitState) flushBuilder() {
	if s.builder.Len() > 0 {
		s.result = append(s.result, s.builder.String())
		s.builder.Reset()
	}
}

// handleEscapeChar 处理转义字符的逻辑
//
// 将转义字符和下一个字符作为一个整体处理
// 保持转义符和字符，不改变原始内容
//
// 参数:
//   - state: 拆分状态
//   - runes: 输入的 rune 切片
//   - i: 当前位置
//
// 返回值:
//   - int: 新的位置索引
func (s *splitState) handleEscapeChar(runes []rune, i int) int {
	// 转义符是最后一个字符的情况
	if i+1 >= len(runes) {
		s.builder.WriteString("\\")
		return i + 1
	}

	// 正常情况：保持转义符和下一个字符
	nextChar := runes[i+1]
	s.builder.WriteString("\\" + string(nextChar))
	return i + 2 // 跳过转义符和被转义的字符
}

// checkMultiCharOperator 检查并处理多字符操作符
//
// 支持的多字符操作符:
//   - && : 逻辑与
//   - || : 逻辑或
//   - >> : 追加重定向
//   - << : here document
//
// 参数:
//   - state: 拆分状态
//   - runes: 输入的 rune 切片
//   - i: 当前位置
//
// 返回值:
//   - bool: 如果是多字符操作符返回 true，否则返回 false
//   - int: 新的位置索引（如果是多字符操作符）
func checkMultiCharOperator(state *splitState, runes []rune, i int) (bool, int) {
	if i+1 >= len(runes) {
		return false, i
	}

	switch {
	case runes[i] == '&' && runes[i+1] == '&':
		// 处理&&操作符
		state.flushBuilder()
		state.result = append(state.result, "&&")
		state.emptyQuote = false
		return true, i + 1

	case runes[i] == '|' && runes[i+1] == '|':
		// 处理||操作符
		state.flushBuilder()
		state.result = append(state.result, "||")
		state.emptyQuote = false
		return true, i + 1

	case runes[i] == '>' && runes[i+1] == '>':
		// 处理>>操作符
		state.flushBuilder()
		state.result = append(state.result, ">>")
		state.emptyQuote = false
		return true, i + 1

	case runes[i] == '<' && runes[i+1] == '<':
		// 处理<<操作符
		state.flushBuilder()
		state.result = append(state.result, "<<")
		state.emptyQuote = false
		return true, i + 1

	default:
		// 其他情况：不是多字符操作符
		return false, i
	}
}

// splitInternal 将命令字符串拆分为命令切片（内部函数）
//
// 实现原理：
//  1. 去除首尾空白
//  2. 遍历每个字符
//  3. 处理引号状态切换
//  4. 统一处理转义字符：
//     - 无论在引号内外，都保持转义符和被转义字符原样
//     - 将转义符和被转义字符作为一个整体处理
//  5. 在非引号状态下遇到空格时分割
//  6. 检测未闭合的引号
//
// 参数:
//   - cmdStr: 要拆分的命令字符串
//
// 返回值:
//   - []string: 拆分后的命令切片
//   - error: 拆分错误，成功时为 nil
func splitInternal(cmdStr string) ([]string, error) {
	// 去除首尾空白
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return []string{}, nil
	}

	// 初始化拆分状态
	state := newSplitState()
	runes := []rune(cmdStr)

	// 遍历每个字符
	for i := 0; i < len(runes); i++ {
		currentRune := runes[i]

		// 处理转义字符（统一处理，保持原样）
		if currentRune == '\\' && i+1 < len(runes) {
			i = state.handleEscapeChar(runes, i) - 1
			continue
		}

		// 检查多字符操作符
		if isMultiOp, newPos := checkMultiCharOperator(state, runes, i); isMultiOp {
			i = newPos
			continue
		}

		switch {
		case isQuote(currentRune):
			// 引号字符处理
			state.handleQuoteChar(currentRune)

		case isSpecialChar(currentRune):
			// 特殊字符处理
			state.handleSpecialChar(currentRune)

		case unicode.IsSpace(currentRune) && !state.inQuotes:
			// 处理空格分隔符换行符等
			state.handleSeparator()

		default:
			// 处理普通字符
			state.builder.WriteRune(currentRune)
			state.emptyQuote = false
		}
	}

	// 添加最后一个命令片段
	if state.builder.Len() > 0 || state.hasQuoteInWord {
		state.result = append(state.result, state.builder.String())
	}

	// 检查引号是否闭合，如果未闭合返回带有未闭合引号类型的错误
	if state.inQuotes {
		return state.result, &UnclosedQuoteError{QuoteType: state.quote}
	}

	return state.result, nil
}

// handleQuoteChar 处理引号字符的逻辑
//
// 参数:
//   - ch: 当前字符
func (s *splitState) handleQuoteChar(ch rune) {
	switch {
	case !s.inQuotes: // 进入引号状态
		s.inQuotes = true   // 标记进入引号状态
		s.quote = ch        // 记录当前引号类型
		s.emptyQuote = true // 初始化空引号状态为 true

	case ch == s.quote: // 退出引号状态
		s.inQuotes = false // 标记退出引号状态
		// 检查引号内是否有内容（非空）
		if !s.emptyQuote || s.builder.Len() == 0 {
			s.hasQuoteInWord = true
		}

	default: // 引号内：普通字符处理
		s.builder.WriteRune(ch)
		s.emptyQuote = false
	}
}

// handleSeparator 处理分隔符（空格或制表符）的逻辑
//
// 判断逻辑:
//   - builder.Len() > 0：片段中有内容（非空）
//   - 只在有内容时添加token，避免空字符串token
func (s *splitState) handleSeparator() {
	// 判断是否需要添加当前命令片段
	s.flushBuilder()
	s.hasQuoteInWord = false // 重置标记，为下一个片段做准备
}
