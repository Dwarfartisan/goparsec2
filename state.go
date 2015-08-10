package goparsec2

import "fmt"

// State 是基本的状态操作接口
type State interface {
	Pos() int
	SeekTo(idx int) bool
	Next() (interface{}, error)
	Trap(message string, args ...interface{}) Error
}

// BasicState 实现最基本的 State 操作
type BasicState struct {
	buffer []interface{}
	index  int
}

// NewBasicState 构造一个新的 BasicState
func NewBasicState(data []interface{}) BasicState {
	buffer := make([]interface{}, len(data))
	copy(buffer, data)
	return BasicState{
		buffer,
		0,
	}
}

// BasicStateFromText 构造一个新的 BasicState
func BasicStateFromText(str string) BasicState {
	data := []rune(str)
	buffer := make([]interface{}, 0, len(data))
	for _, r := range data {
		buffer = append(buffer, r)
	}
	return BasicState{
		buffer,
		0,
	}
}

// Pos 返回 state 的当前位置
func (state *BasicState) Pos() int {
	return state.index
}

// SeekTo 实现指针移动，如果越界，会返回 false
func (state *BasicState) SeekTo(idx int) bool {
	if idx < 0 || idx > len(state.buffer) {
		return false
	}
	state.index = idx
	return true
}

// Next 实现迭代逻辑
func (state *BasicState) Next() (interface{}, error) {
	if state.index == len(state.buffer) {
		return nil, state.Trap("eof")
	}
	re := state.buffer[state.index]
	state.index++
	return re, nil
}

// Trap 是构造错误信息的辅助函数，它传递错误的位置，并提供字符串格式化功能
func (state *BasicState) Trap(message string, args ...interface{}) Error {
	return Error{state.index, fmt.Sprintf(message, args...)}
}

// Error 实现基本的错误信息结构
type Error struct {
	Pos     int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("stop at %d : %v", e.Pos, e.Message)
}
