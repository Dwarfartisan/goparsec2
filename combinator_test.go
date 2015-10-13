package goP2

import (
	"fmt"
	"testing"
	"unicode"
)

var iprPsc = Repeat(1, 3, Digit).Bind(func(x interface{}) P {
	buffer := x.([]interface{})
	data := make([]rune, 0, len(buffer))
	for _, r := range buffer {
		data = append(data, r.(rune))
	}
	return Return(string(data))
})
var ipPsc = iprPsc.Bind(func(first interface{}) P {
	return Times(3, Chr('.').Then(iprPsc)).Bind(func(postfix interface{}) P {
		buffer := make([]interface{}, 1, 4)
		buffer[0] = first
		buffer = append(buffer, postfix.([]interface{})...)
		data := fmt.Sprintf("%s.%s.%s.%s", buffer...)
		return Return(data)
	})
})
var dnPsc = Many1(RuneP("word", func(x rune) bool { return !unicode.IsSpace(x) }))
var hPsc = Choice(ipPsc, dnPsc)
var ptPsc = P(UInt)

var listen = Try(ipPsc).Over(Chr(':')).Bind(func(ip interface{}) P {
	return ptPsc.Bind(func(port interface{}) P {
		return Return([]string{ip.(string), port.(string)})
	})
})

func TestListen(t *testing.T) {
	data := "127.0.0.1:8080"
	state := BasicStateFromText(data)
	re, err := listen.Parse(&state)
	if err != nil {
		t.Fatal(err)
	}
	var output []string
	var ok bool
	if output, ok = re.([]string); ok {
		if len(output) == 2 {
			if output[0] != "127.0.0.1" {
				t.Fatalf("Expect 127.0.0.1 but %v", output[0])
			}
			if output[1] != "8080" {
				t.Fatalf("Expect 8080 but %v", output[1])
			}
		}
	} else {
		t.Fatalf("Expect [\"127.0.0.1\", \"8080\"] but %v is %t", output, output)
	}
}
