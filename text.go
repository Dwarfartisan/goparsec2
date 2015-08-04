package goparsec2

import (
	"fmt"
	"unicode"
)

// Chr 判断下一个字符是否与给定值相等
func Chr(val rune) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(int32); ok {
			if c == val {
				return c, nil
			}
			return nil, state.Trap("Expact '%v' but '%v'", string([]rune{val}), string([]rune{c}))
		}
		return nil, state.Trap("Expact a rune '%s' but x is %v", string([]rune{val}), x)
	}}
}

// NChr 判断下一个字符是否与给定值不相等
func NChr(val rune) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(int32); ok {
			if c == val {
				return nil, state.Trap("Expact not '%v' but '%v'", string([]rune{val}), string([]rune{c}))
			}
			return c, nil
		}
		return nil, state.Trap("Expact a rune '%s' but x is %v", string([]rune{val}), x)
	}}
}

// RuneOf 检查后续的字符是否是给定值中的某一个
func RuneOf(str string) Parsec {
	data := []rune(str)
	return Parsec{func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(int32); ok {
			for _, r := range data {
				if c == r {
					return c, nil
				}
			}
			return nil, state.Trap("Expact rune in '%s' but '%s'", str, string([]rune{c}))
		}
		return nil, state.Trap("Expact rune in '%s' but x=%v is %t", str, x, x)
	}}
}

// RuneNone 检查后续的字符是否不是给定值中的任一个
func RuneNone(str string) Parsec {
	data := []rune(str)
	return Parsec{func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(int32); ok {
			for _, r := range data {
				if c == r {
					return nil, state.Trap("Expact rune none of '%s' but '%s'", str, string([]rune{c}))
				}
			}
			return c, nil
		}
		return nil, state.Trap("Expact rune none of '%s' but x=%v is %t", str, x, x)
	}}
}

// Str 判断后续的字符串是否匹配给定的串
func Str(str string) Parsec {
	data := []rune(str)
	return Parsec{func(state State) (interface{}, error) {
		for _, r := range data {
			_, err := Chr(r).Parse(state)
			if err != nil {
				return nil, err
			}
		}
		return str, nil
	}}
}

// RuneParsec 通过一个谓词参数，提供通用的 rune 算子生成判断
func RuneParsec(name string, pred func(r rune) bool) Parsec {
	return Parsec{func(state State) (interface{}, error) {
		x, err := state.Next()
		if err != nil {
			return nil, err
		}
		if c, ok := x.(int32); ok {
			r := rune(c)
			if pred(r) {
				return c, nil
			}
			return nil, state.Trap("Expact %s but '%v'", name, string([]rune{r}))
		}
		return nil, state.Trap("Expact %s but x=%v is %t", name, x, x)
	}}
}

// Space 构造一个空格校验算子
func Space() Parsec {
	return RuneParsec("space", unicode.IsSpace)
}

// Letter 构造一个字母校验算子
func Letter() Parsec {
	return RuneParsec("letter", unicode.IsLetter)
}

// Number 构造一个 Number 校验算子
func Number() Parsec {
	return RuneParsec("number", unicode.IsNumber)
}

// Digit 构造一个数字字符校验算子
func Digit() Parsec {
	return RuneParsec("digit", unicode.IsDigit)
}

// UInt 返回一个无符号整型的解析算子
func UInt() Parsec {
	return Many1(Digit()).Bind(func(values interface{}) Parsec {
		buffer := values.([]interface{})
		data := make([]rune, 0, len(buffer))
		for _, value := range buffer {
			data = append(data, value.(rune))
		}
		return Return(string(data))
	})
}

// Int 返回一个有符号整型的解析算子
func Int() Parsec {
	binder := func(value interface{}) Parsec {
		return Return(fmt.Sprintf("-%v" + value.(string)))
	}
	return Choice(Try(Chr('-').Then(UInt()).Bind(binder)), Return("").Then(UInt()))
}

// UFloat 返回一个无符号实数的解析算子
func UFloat() Parsec {
	return UInt().Over(Chr('.')).Bind(func(left interface{}) Parsec {
		return UInt().Bind(func(right interface{}) Parsec {
			return Return(fmt.Sprintf("%s.%s", left, right))
		})
	})
}

// Float 返回一个有符号实数的解析算子
func Float() Parsec {
	binder := func(value interface{}) Parsec {
		return Return("-" + value.(string))
	}
	return Choice(Try(Chr('-').Then(UFloat()).Bind(binder)), UFloat())
}
