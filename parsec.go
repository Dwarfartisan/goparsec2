package goparsec2

// Parsec 是算子的公共抽象类型，实现 Monad 和解析逻辑
type Parsec struct {
	Psc func(state State) (interface{}, error)
}

//Parse 简单的调用被封装的算子逻辑
func (parsec Parsec) Parse(state State) (interface{}, error) {
	return parsec.Psc(state)
}

//Exec 调用被封装的算子，如果返回错误，用panic抛出
func (parsec Parsec) Exec(state State) interface{} {
	re, err := parsec.Psc(state)
	if err != nil {
		panic(err)
	}
	return re
}

// Bind 方法实现 Monad >>= 运算
func (parsec Parsec) Bind(binder func(interface{}) Parsec) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		x, err := parsec.Psc(state)
		if err != nil {
			return nil, err
		}
		psc := binder(x)
		return psc.Psc(state)
	}}
}

// Then 方法实现 Monad 的 >> 运算
func (parsec Parsec) Then(psc Parsec) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		_, err := parsec.Psc(state)
		if err != nil {
			return nil, err
		}
		return psc.Psc(state)
	}}
}

// Over 方法实现一个简化的 bind 逻辑，如果两个算子都成功，返回前一个算子的结果，否则返回第一个发生的错误
func (parsec Parsec) Over(psc Parsec) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		re, err := parsec.Psc(state)
		if err != nil {
			return nil, err
		}
		_, err = psc.Psc(state)
		if err != nil {
			return nil, err
		}
		return re, nil
	}}
}

// Env 函数与 Parsec 算子的 Exec 方法对应，将其抛出的 error 还原为无 panic 的流程，用于模拟
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
func Do(fn func(State) interface{}) Parsec {
	return Parsec{func(state State) (re interface{}, err error) {
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
	}}
}
