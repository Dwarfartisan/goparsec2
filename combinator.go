package goparsec2

// Try 尝试运行给定算子，如果给定算子报错，将state复位再返回错误信息
func Try(psc Parsec) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		idx := state.Pos()
		re, err := psc.Parse(state)
		if err == nil {
			return re, nil
		}
		state.SeekTo(idx)
		return nil, err
	}}
}

// Choice 逐个常识给定的算子，直到某个成功或者 state 无法复位，或者全部失败
func Choice(parsecs ...Parsec) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		var err error
		for _, p := range parsecs {
			idx := state.Pos()
			re, err := p.Parse(state)
			if err == nil {
				return re, nil
			}
			if state.Pos() != idx {
				return nil, err
			}
		}
		return nil, err
	}}
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
func Between(b, psc, e Parsec) Parsec {
	return b.Then(psc).Over(e)
}

// SepBy1 返回匹配 1 到若干次的带分隔符的算子
func SepBy1(p, sep Parsec) Parsec {
	tail := func(value interface{}) Parsec {
		head := Many(Try(sep.Then(p)))
		return head.Bind(func(values interface{}) Parsec {
			return Return(append([]interface{}{value}, values.([]interface{})...))
		})
	}
	return p.Bind(tail)
}

// SepBy 返回匹配 0 到若干次的带分隔符的算子
func SepBy(p, sep Parsec) Parsec {
	return Choice(SepBy1(Try(p), sep), Return([]interface{}{}))
}

// ManyTil 返回以指定算子结尾的  Many
func ManyTil(p, e Parsec) Parsec {
	return Many(p).Over(e)
}

// Many1Til 返回以指定算子结尾的  Many1
func Many1Til(p, e Parsec) Parsec {
	return Many1(p).Over(e)
}

// Skip 忽略指定 0 到若干次算子
func Skip(p Parsec) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		for {
			_, err := Try(p).Parse(state)
			if err != nil {
				return nil, err
			}
		}
	}}
}

// Skip1 忽略指定 1 到若干次算子
func Skip1(p Parsec) Parsec {
	return p.Then(Skip(p))
}
