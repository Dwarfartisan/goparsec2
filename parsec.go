package goparsec2
import "fmt"
// Parsec 是算子的公共抽象类型，实现 Monad 和解析逻辑
type Parsec struct {
	Psc func(state State) (interface{}, error)
}

//Parse 简单的调用被封装的算子逻辑
func (parsec Parsec) Parse(state State) (interface{}, error) {
	return parsec.Psc(state)
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
		fmt.Println("over run")
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
