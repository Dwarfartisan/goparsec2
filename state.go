package goP2

import (
	"fmt"
	"math"
)

// State 是基本的状态操作接口
type State interface {
	Pos() int
	SeekTo(int) bool
	Next() (interface{}, error)
	Trap(message string, args ...interface{}) error
}

// TranState 表示支持事务的 State 约定
type TranState interface {
	State
	Begin() int
	Commit(int)
	Rollback(int)
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

// TState 支持事务
type TState struct {
	State
	nextTran int
	begin    int // begin 总是保存最小的事务位置，如果当前没有事务，值为 －1
}

// NewTState 构造一个新的 TState
func NewTState(data []interface{}) *TState {
	state := NewBasicState(data)
	return &TState{
		&state,
		0,
		-1,
	}
}

// TStateFromText 构造一个新的 BasicState
func TStateFromText(str string) *TState {
	state := BasicStateFromText(str)
	return &TState{
		&state,
		0,
		-1,
	}
}

// TStateFromState 将一个无事务的 State 封装为有事务的
func TStateFromState(state State) *TState {
	return &TState{
		state,
		0,
		-1,
	}
}

// Begin 开始一个事务并返回事务号
func (state *TState) Begin() int {
	if state.begin == -1 {
		state.begin = state.Pos()
	} else {
		state.begin = int(math.Min(float64(state.begin), float64(state.Pos())))
	}
	return state.begin
}

// Commit 提交一个事务，将其从注册状态中删除
func (state *TState) Commit(tran int) {
	if state.begin == tran {
		state.begin = -1
	} else {
		state.begin = int(math.Min(float64(state.begin), float64(state.Pos())))
	}
}

// Rollback 取消一个事务，将 pos 移动到 该位置。
func (state *TState) Rollback(tran int) {
	state.SeekTo(tran)
	if state.begin == tran {
		state.begin = -1
	} else {
		state.begin = int(math.Min(float64(state.begin), float64(state.Pos())))
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
