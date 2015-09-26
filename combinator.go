package goparsec2

import "fmt"

// Try 尝试运行给定算子，如果给定算子报错，将state复位再返回错误信息
func Try(psc Parsec) Parsec {
	return func(state State) (interface{}, error) {
		idx := state.Pos()
		re, err := psc.Parse(state)
		if err == nil {
			return re, nil
		}
		state.SeekTo(idx)
		return nil, err
	}
}

// Choice 逐个尝试给定的算子，直到某个成功或者 state 无法复位，或者全部失败
func Choice(parsecs ...Parsec) Parsec {
	return func(state State) (interface{}, error) {
		var err error
		for _, p := range parsecs {
			var re interface{}
			idx := state.Pos()
			re, err = p.Parse(state)
			if err == nil {
				return re, nil
			}
			if state.Pos() != idx {
				return nil, err
			}
		}
		//下面这个分支确保最后一个算子是 Fail 之类的零步进算子时，也能把错误信息传递出来。
		return nil, err
	}
}

// Many 匹配 0 到若干次 psc 并返回结果序列
func Many(psc Parsec) Parsec {
	return Choice(Try(Many1(psc)), Return([]interface{}{}))
}

// Many1 匹配 1 到若干次 psc 并返回结果序列
func Many1(psc Parsec) Parsec {
	tail := func(value interface{}) Parsec {
		head := Many(Try(psc))
		return head.Bind(func(values interface{}) Parsec {
			return Return(append([]interface{}{value}, values.([]interface{})...))
		})
	}
	return psc.Bind(tail)
}

//Between 构造一个有边界算子的 Parsec
func Between(open, close, psc Parsec) Parsec {
	return open.Then(psc).Over(close)
}

// SepBy1 返回匹配 1 到若干次的带分隔符的算子
func SepBy1(p, sep Parsec) Parsec {
	binder := func(value interface{}) Parsec {
		head := Many(sep.Then(p))
		return head.Bind(func(values interface{}) Parsec {
			return Return(append([]interface{}{value}, values.([]interface{})...))
		})
	}
	return p.Bind(binder)
}

// SepBy 返回匹配 0 到若干次的带分隔符的算子
func SepBy(p, sep Parsec) Parsec {
	return Choice(Try(SepBy1(p, sep)), Return([]interface{}{}))
}

// ManyTil 返回以指定算子结尾的  Many
func ManyTil(p, e Parsec) Parsec {
	return Many(p).Over(e)
}

// Many1Til 返回以指定算子结尾的  Many1
func Many1Til(p, e Parsec) Parsec {
	return Many1(p).Over(e)
}

// Skip 忽略 0 到若干次指定算子
func Skip(p Parsec) Parsec {
	return func(state State) (interface{}, error) {
		for {
			_, err := Try(p).Parse(state)
			if err != nil {
				return nil, nil
			}
		}
	}
}

// Skip1 忽略 1 到若干次指定算子
func Skip1(p Parsec) Parsec {
	return p.Then(Skip(p))
}

// FailIf 是算子的否定检查，如果给定算子匹配成功，返回错误信息。否则退换复位并且返回 nil，
// 可以用于边界检查。
func FailIf(psc Parsec) Parsec {
	message := fmt.Sprintf("Expect the parsec %v failed but it success.", psc)
	return Choice(Try(psc).Then(Fail(message)), Return(nil))
}

// Repeat 函数生成一个 parsec 算子，它匹配指定算子x到y次。
func Repeat(x, y int, psc Parsec) Parsec {
	if x >= y {
		message := fmt.Sprintf("x must greater than y but x=%d and y=%d", x, y)
		panic(message)
	}
	return Times(x, psc).Bind(func(val interface{}) Parsec {
		return UpTo(y-x, psc).Bind(func(y interface{}) Parsec {
			buffer := val.([]interface{})
			buffer = append(buffer, y.([]interface{})...)
			return Return(buffer)
		})
	})
}

// InRange 函数生成一个 parsec 算子，它匹配指定算子x到y次。如果第 y+1 次仍然成功，返回错误信息
func InRange(x, y int, psc Parsec) Parsec {
	if x >= y {
		message := fmt.Sprintf("x must greater than y but x=%d and y=%d", x, y)
		panic(message)
	}
	return Times(x, psc).Bind(func(val interface{}) Parsec {
		return AtMost(y-x, psc).Bind(func(y interface{}) Parsec {
			buffer := val.([]interface{})
			buffer = append(buffer, y.([]interface{})...)
			return Return(buffer)
		})
	})
}

// UpTo 函数匹配 0 到 x 次 psc
func UpTo(x int, psc Parsec) Parsec {
	return func(state State) (interface{}, error) {
		var re = make([]interface{}, 0, x)
		for i := 0; i < x; i++ {
			item, err := Try(psc).Parse(state)
			if err != nil {
				return re, nil
			}
			re = append(re, item)
		}
		return re, nil
	}
}

// AtMost 函数匹配至多 x 次 psc ，如果后续的数据仍然匹配成功，返回错误信息
func AtMost(x int, psc Parsec) Parsec {
	return UpTo(x, psc).Bind(func(val interface{}) Parsec {
		re := val.([]interface{})
		if len(re) < x {
			return Return(val)
		}
		return FailIf(psc)
	})
}

// AtLeast 函数匹配至少 x 次 psc
func AtLeast(x int, psc Parsec) Parsec {
	return Times(x, psc).Bind(func(valx interface{}) Parsec {
		return Many(psc).Bind(func(valy interface{}) Parsec {
			var re = valx.([]interface{})
			re = append(re, valy.([]interface{})...)
			return Return(re)
		})
	})
}

// Times 函数生成一个 parsec 算子，它匹配指定算子x次。我们在这里用它构造一个不严谨的ip判定
func Times(x int, psc Parsec) Parsec {
	return func(state State) (interface{}, error) {
		var re = make([]interface{}, 0, x)
		for i := 0; i < x; i++ {
			item, err := psc.Parse(state)
			if err != nil {
				return nil, err
			}
			re = append(re, item)
		}
		return re, nil
	}
}

// Union 逐个执行每个给定的算子，直到出错或者完整，将结果中非空值合成一个 []interface{} 返回
func Union(parsers ...Parsec) Parsec {
	return func(st State) (interface{}, error) {
		var ret = make([]interface{}, 0, len(parsers))
		for _, parser := range parsers {
			val, err := parser(st)
			if err == nil {
				if val != nil {
					ret = append(ret, val)
				}
			} else {
				return nil, err
			}
		}
		return ret, nil
	}
}

// UnionAll 逐个执行每个给定的算子，直到出错或者完整，将结果合成一个 []interface{} 返回
func UnionAll(parsers ...Parsec) Parsec {
	return func(st State) (interface{}, error) {
		var ret = make([]interface{}, 0, len(parsers))
		for _, parser := range parsers {
			val, err := parser(st)
			if err == nil {
				ret = append(ret, val)
			} else {
				return nil, err
			}
		}
		return ret, nil
	}
}

// Option 实现一个默认值封装，如果给定算子失败，返回默认值 x
func Option(v interface{}, x Parsec) Parsec {
	return func(st State) (interface{}, error) {
		return Choice(x, Return(v))(st)
	}
}

// Maybe 在算子失败时返回 nil
func Maybe(p Parsec) Parsec {
	return Option(Return(nil), p)
}
