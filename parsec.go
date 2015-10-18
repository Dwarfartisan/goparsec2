package goP2

// P 是算子的公共抽象类型，实现 Monad 和解析逻辑
type P func(state State) (interface{}, error)

//Parse 简单的调用被封装的算子逻辑
func (p P) Parse(state State) (interface{}, error) {
	return p(state)
}

//Exec 调用被封装的算子，如果返回错误，用panic抛出
func (p P) Exec(state State) interface{} {
	re, err := p(state)
	if err != nil {
		panic(err)
	}
	return re
}

// Bind 方法实现 Monad >>= 运算
func (p P) Bind(binder func(interface{}) P) P {
	return func(state State) (interface{}, error) {
		x, err := p(state)
		if err != nil {
			return nil, err
		}
		psc := binder(x)
		return psc(state)
	}
}

// Then 方法实现 Monad 的 >> 运算
func (p P) Then(psc P) P {
	return func(state State) (interface{}, error) {
		_, err := p(state)
		if err != nil {
			return nil, err
		}
		return psc(state)
	}
}

// Over 方法实现一个简化的 bind 逻辑，如果两个算子都成功，返回前一个算子的结果，否则返回第一个发生的错误
func (p P) Over(psc P) P {
	return func(state State) (interface{}, error) {
		re, err := p(state)
		if err != nil {
			return nil, err
		}
		_, err = psc(state)
		if err != nil {
			return nil, err
		}
		return re, nil
	}
}

// Env 函数与 P 算子的 Exec 方法对应，将其抛出的 error 还原为无 panic 的流程，用于模拟
// Haskell Monad Do 。需要注意的是，捕获的非 error 类型的panic会重新抛出。
func Env(fn func() interface{}) (re interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				re = nil
			} else {
				panic(r)
			}
		}
	}()
	re = fn()
	err = nil
	return
}

// Do 构造一个算子，其内部类似 Monad Do Environment ，将 Exec 形式恢复成 Parse 形式。
// 需要注意的是，捕获的非 error 类型的panic会重新抛出。
func Do(fn func(State) interface{}) P {
	return func(state State) (re interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					err = e
					re = nil
				} else {
					panic(r)
				}
			}
		}()
		re = fn(state)
		err = nil
		return
	}
}

// M 工具函数实现将函数明确转型为 P 算子的逻辑
// func P(fn func(State) (interface{}, error)) P {
// 	return fn
// }
