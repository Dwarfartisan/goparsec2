package goP2

import "fmt"

// State 是基本的状态操作接口
type State interface {
	Pos() int
	Begin() int
	Commit(int)
	Rollback(int)
	Next() (interface{}, error)
	Trap(message string, args ...interface{}) error
}

// BasicState 实现最基本的 State 操作
type BasicState struct {
	buffer   []interface{}
	index    int
	nextTran int
	trans    map[int]int
}

// NewBasicState 构造一个新的 BasicState
func NewBasicState(data []interface{}) BasicState {
	buffer := make([]interface{}, len(data))
	copy(buffer, data)
	return BasicState{
		buffer,
		0,
		0,
		map[int]int{},
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
		0,
		map[int]int{},
	}
}

// Pos 返回 state 的当前位置
func (state *BasicState) Pos() int {
	return state.index
}

//Begin 注册并返回一个事务号
func (state *BasicState) Begin() int {
	state.trans[state.nextTran] = state.Pos()
	var re = state.nextTran
	state.nextTran++
	return re
}

// Commit 表示事务成功，删除该事务号
func (state *BasicState) Commit(num int) {
	delete(state.trans, num)
}

// Rollback 表示事务失败，删除事务号，并将 state 的 pos 还原到该事务开始时的位置
func (state *BasicState) Rollback(num int) {
	state.index = state.trans[num]
	delete(state.trans, num)
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

// Error 实现基本的错误信息结构
type Error struct {
	Pos     int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("stop at %d : %v", e.Pos, e.Message)
}
