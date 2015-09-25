package goparsec2

// Nil 判断当前元素是否为 nil
func Nil(state State) (interface{}, error) {
	data, err := state.Next()
	if err == nil {
		return nil, state.Trap("Expect nil but error %v", err)
	}
	if data == nil {
		return nil, nil
	}
	return nil, state.Trap("Expect nil but %v", data)
}

// AsInt 判断当前元素是否为一个 int
func AsInt(state State) (interface{}, error) {
	data, err := state.Next()
	if err == nil {
		return nil, state.Trap("Expect a int value but error %v", err)
	}
	if _, ok := data.(int); ok {
		return data, nil
	}
	return nil, state.Trap("Expect a int value but %v", data)
}

// AsFloat64 判断当前元素是否为 float64
func AsFloat64(state State) (interface{}, error) {
	data, err := state.Next()
	if err == nil {
		return nil, state.Trap("Expect a float64 value but error %v", err)
	}
	if _, ok := data.(float64); ok {
		return data, nil
	}
	return nil, state.Trap("Expect a float64 value but %v", data)
}

// AsFloat32 判断当前元素是否为 float32
func AsFloat32(state State) (interface{}, error) {
	data, err := state.Next()
	if err == nil {
		return nil, state.Trap("Expect a float32 value but error %v", err)
	}
	if _, ok := data.(float32); ok {
		return data, nil
	}
	return nil, state.Trap("Expect a float32 value but %v", data)
}

// AsString 判断当前元素是否为 string
func AsString(state State) (interface{}, error) {
	data, err := state.Next()
	if err == nil {
		return nil, state.Trap("Expect a string value but error %v", err)
	}
	if _, ok := data.(string); ok {
		return data, nil
	}
	return nil, state.Trap("Expect a string value but %v", data)
}
