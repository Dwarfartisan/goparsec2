package goP2

import "fmt"

// State 是基本的状态操作接口
type State interface {
	Pos() int
	SeekTo(int) bool
	Next() (interface{}, error)
	Trap(message string, args ...interface{}) error
	Begin() int
	Commit(int)
	Rollback(int)
}

// BasicState 实现最基本的 State 操作
type BasicState struct {
	buffer []interface{}
	index  int
	begin  int
}

// NewBasicState 构造一个新的 BasicState
func NewBasicState(data []interface{}) BasicState {
	buffer := make([]interface{}, len(data))
	copy(buffer, data)
	return BasicState{
		buffer,
		0,
		-1,
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
		-1,
	}
}

// Pos 返回 state 的当前位置
func (state *BasicState) Pos() int {
	return state.index
}

//SeekTo 将指针移动到指定位置
func (state *BasicState) SeekTo(pos int) bool {
	if 0 <= pos && pos < len(state.buffer) {
		state.index = pos
		return true
	}
	return false
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
func (state *BasicState) Trap(message string, args ...interface{}) error {
	return Error{state.index, fmt.Sprintf(message, args...)}
}

// Begin 开始一个事务并返回事务号，State 的 Begin 总是记录比较靠后的位置。
func (state *BasicState) Begin() int {
	if state.begin == -1 {
		state.begin = state.Pos()
	} else {
		state.begin = min(state.begin, state.Pos())
	}
	return state.Pos()
}

// Commit 提交一个事务，将其从注册状态中删除，将事务位置保存为比较靠前的位置
func (state *BasicState) Commit(tran int) {
	if state.begin == tran {
		state.begin = -1
	}
}

// Rollback 取消一个事务，将 pos 移动到 该位置，将事务位置保存为比较靠前的位置
func (state *BasicState) Rollback(tran int) {
	state.SeekTo(tran)
	if state.begin == tran {
		state.begin = -1
	}
}

// Error 实现基本的错误信息结构
type Error struct {
	Pos     int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("stop at %d : %v", e.Pos, e.Message)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
