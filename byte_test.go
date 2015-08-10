package goparsec2

import (
	"encoding/json"
	"testing"
)

func testState(data string) State {
	buffer := []byte(data)
	toBuf := make([]interface{}, len(buffer), len(buffer))
	for idx := range buffer {
		toBuf[idx] = buffer[idx]
	}
	state := NewBasicState(toBuf)
	return &state
}

func content() Parsec {
	return Skip(Choice(Try(ByteNone("\\\"")), Try(Byte('\\').Then(One()))))
}

func TestContent(t *testing.T) {
	state := testState("It is string content \\\"")
	_, err := content().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

var str = Choice(Between(Byte('"'), content(), Byte('"')),
	Between(Byte('\''), content(), Byte('\'')))

func TestStr0(t *testing.T) {
	state := testState("\"content\"")
	_, err := str.Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStr1(t *testing.T) {
	state := testState("\"It is \\\" string.\"")
	_, err := str.Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

var atom = Skip1(ByteNone(" ,[]{}\"'"))

func TestAtom0(t *testing.T) {
	state := testState("atom.")
	_, err := atom.Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAtom1(t *testing.T) {
	state := testState("123243")
	_, err := atom.Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

var spaces = Skip(Byte(' '))
var skipColon = spaces.Then(Byte(':')).Then(spaces)
var skipComma = spaces.Then(Byte(',')).Then(spaces)

func node() Parsec {
	return Parsec{func(state State) (interface{}, error) {
		re, err := Try(str).Parse(state)
		if err == nil {
			return re, nil
		}
		re, err = Try(atom).Parse(state)
		if err == nil {
			return re, nil
		}
		return Try(j()).Parse(state)
	}}
}
func TestNode0(t *testing.T) {
	state := testState("\"It is a node .\"")
	_, err := node().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNode1(t *testing.T) {
	state := testState("12343")
	_, err := node().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNode2(t *testing.T) {
	state := testState("1234.3")
	_, err := node().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func arrayBody() Parsec {
	return SepBy(node(), skipComma)
}

func TestArrayBody0(t *testing.T) {
	state := testState("\"content\"")
	_, err := str.Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayBody1(t *testing.T) {
	state := testState("\"It\", \"is\", \"a\", \"content\"")
	_, err := arrayBody().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayBody2(t *testing.T) {
	state := testState("1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89")
	_, err := arrayBody().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func TestArrayBody3(t *testing.T) {
	state := testState("1,1,2,3,5,8,13,21,34,55,89")
	_, err := arrayBody().Then(EOF()).Parse(state)
	//_, err := SepBy1(Byte(','), Skip1(ByteNone(","))).Then(EOF()).Parse(state)

	if err != nil {
		t.Fatal(err)
	}
}

func pair() Parsec {
	return str.Then(skipColon).Then(node())
}
func TestPair(t *testing.T) {
	state := testState("\"content\" : [\"quit\"]")
	_, err := pair().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}

func dictBody() Parsec {
	return SepBy(pair(), skipComma)
}
func TestDictBody(t *testing.T) {
	state := testState("\"content\" : [\"quit\"]")
	_, err := dictBody().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}
func array() Parsec {
	return Between(Byte('[').Then(spaces), arrayBody(), spaces.Then(Byte(']')))
}
func dict() Parsec {
	return Between(Byte('{').Then(spaces), dictBody(), spaces.Then(Byte('}')))
}
func TestDict0(t *testing.T) {
	data := map[string]interface{}{
		// "meta": map[string]interface{}{
		// 	"category": "command",
		// },
		"content": []interface{}{"quit"},
	}
	buffer, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	state := testState(string(buffer))
	_, err = dict().Parse(state)
	if err != nil {
		t.Fatal(err)
	}

}

func j() Parsec {
	return Choice(Try(array()), Try(dict()))
}

func TestJsonTrue0(t *testing.T) {
	data := map[string]interface{}{
		"meta": map[string]interface{}{
			"category": "command",
		},
		"content": []interface{}{"quite"},
	}
	buffer, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	state := testState(string(buffer))
	_, err = j().Then(EOF()).Parse(state)
	if err != nil {
		t.Fatal(err)
	}
}
