package goparsec2

import "reflect"

// One 仅仅简单的返回下一个迭代结果，或者得到 eof 错误
func One() Parsec {
	return func(state State) (interface{}, error) {
		return state.Next()
	}
}

// Eq 判断下一个数据是否与给定值相等，这里简单的使用了反射
func Eq(val interface{}) Parsec {
	return func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if reflect.DeepEqual(x, val) {
			return x, nil
		}
		return nil, state.Trap("Expact %v but %v", val, x)
	}
}

// Ne 判断下一个数据是否与给定值不相等，这里简单的使用了反射
func Ne(val interface{}) Parsec {
	return func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if reflect.DeepEqual(x, val) {
			return nil, state.Trap("Expact not %v but %v", val, x)
		}
		return x, nil
	}
}

// Return 生成的算子总是返回给定值
func Return(val interface{}) Parsec {
	return func(state State) (interface{}, error) {
		return val, nil
	}
}

// Fail 生成的算子总是返回给定错误
func Fail(message string, args ...interface{}) Parsec {
	return func(state State) (interface{}, error) {
		return nil, state.Trap(message, args...)
	}
}

// EOF 仅仅到达结尾时匹配成功
func EOF() Parsec {
	return func(state State) (interface{}, error) {
		data, err := state.Next()
		if err == nil {
			return nil, state.Trap("Expect eof but %v", data)
		}
		return nil, nil
	}
}

// OneOf 期待下一个元素属于给定的参数中的一个
func OneOf(args ...interface{}) Parsec {
	return func(state State) (interface{}, error) {
		data, err := state.Next()
		if err != nil {
			return nil, err
		}
		for element := range args {
			if reflect.DeepEqual(data, element) {
				return data, nil
			}
		}
		return nil, state.Trap("Expect one of [%v] but %v", args, data)
	}
}

// NoneOf 期待下一个元素不属于给定的参数中的任一个
func NoneOf(args ...interface{}) Parsec {
	return func(state State) (interface{}, error) {
		data, err := state.Next()
		if err != nil {
			return nil, err
		}
		for element := range args {
			if reflect.DeepEqual(data, element) {
				return nil, state.Trap("Expect none of [%v] but %v", args, data)
			}
		}
		return data, nil
	}
}
