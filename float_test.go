package goparsec2

import "testing"

func TestFloat0(t *testing.T) {
	data := "3.14"
	state := BasicStateFromText(data)
	re, err := M(UFloat).Parse(&state)
	if err != nil {
		t.Fatal(err)
	}
	if output, ok := re.(string); ok {
		if output != "3.14" {
			t.Fatalf("Expect 3.14 but %v", output)
		}
	} else {
		t.Fatalf("Expect string 3.14 but %v is %t", output, output)
	}
}
func TestFloat1(t *testing.T) {
	data := "3.14f"
	state := BasicStateFromText(data)
	re, err := M(Float).Parse(&state)
	if err != nil {
		t.Fatal(err)
	}
	if output, ok := re.(string); ok {
		if output != "3.14" {
			t.Fatalf("Expect 3.14 but %v", output)
		}
	} else {
		t.Fatalf("Expect string 3.14 but %v is %t", output, output)
	}
}

func TestFloat2(t *testing.T) {
	data := "e.14"
	state := BasicStateFromText(data)
	re, err := M(Float).Parse(&state)
	if err == nil {
		t.Fatalf("expect a error when data e.14 but got %v.", re)
	}
}
