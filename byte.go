package goP2

// Byte 判断下一个字节是否与给定值相等
func Byte(val byte) P {
	return func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(byte); ok {
			if c == val {
				return c, nil
			}
			return nil, state.Trap("Expect '%v' but '%v'", string([]byte{val}), string([]byte{c}))
		}
		return nil, state.Trap("Expect a byte '%s' but x is %v", string([]byte{val}), x)
	}
}

// NByte 判断下一个字符是否与给定值不相等
func NByte(val byte) P {
	return func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(byte); ok {
			if c == val {
				return nil, state.Trap("Expect not '%v' but '%v'", string([]byte{val}), string([]byte{c}))
			}
			return c, nil
		}
		return nil, state.Trap("Expect a rune '%s' but x is %v", string([]byte{val}), x)
	}
}

// ByteOf 检查后续的字符是否是给定值中的某一个
func ByteOf(str string) P {
	data := []byte(str)
	return func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(byte); ok {
			for _, r := range data {
				if c == r {
					return c, nil
				}
			}
			return nil, state.Trap("Expect rune in '%s' but '%s'", str, string([]byte{c}))
		}
		return nil, state.Trap("Expect rune in '%s' but x=%v is %t", str, x, x)
	}
}

// ByteNone 检查后续的字符是否不是给定值中的任一个
func ByteNone(str string) P {
	data := []byte(str)
	return func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(byte); ok {
			for _, r := range data {
				if c == r {
					return nil, state.Trap("Expect rune none of '%s' but '%s'", str, string([]byte{c}))
				}
			}
			return c, nil
		}
		return nil, state.Trap("Expect rune none of '%s' but x=%v is %t", str, x, x)
	}
}

// Bytes 判断后续的字节串是否匹配给定的串
func Bytes(str string) P {
	data := []byte(str)
	return func(state State) (interface{}, error) {
		for _, r := range data {
			_, err := Byte(r).Parse(state)
			if err != nil {
				return nil, err
			}
		}
		return str, nil
	}
}

// ByteP 通过一个谓词参数，提供通用的 rune 算子生成判断
func ByteP(name string, pred func(r byte) bool) P {
	return func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(byte); ok {
			r := byte(c)
			if pred(r) {
				return c, nil
			}
			return nil, state.Trap("Expect %s but '%v'", name, string([]byte{r}))
		}
		return nil, state.Trap("Expect %s but x=%v is %t", name, x, x)
	}
}
