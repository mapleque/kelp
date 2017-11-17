package web

import (
	"testing"
)

type Model struct {
	Azstr string  `valid:"/^[a-z]+$/"`
	Name  string  `valid:"/.+/"`
	Fun   int     `valid:"required"`
	ir    int     `valid:"[1:10)"`
	fr    float64 `valid:"(-1:1]"`
	sr    string  `valid:"[1:1]"`
}

func TestValid(t *testing.T) {
	passCases := []*Model{
		&Model{"abc", "abc", 1, 1, 1, "a"},
		&Model{"abc", "abc", 1, 9, 0.1, "a"},
	}
	errCases := []*Model{
		&Model{"Abc", "abc", 1, 1, 1, "a"},  // test A a-z
		&Model{"0", "abc", 1, 1, 1, "a"},    // test 0 a-z
		&Model{"abc", "", 1, 1, 1, "a"},     // test "" /.+/
		&Model{"abc", "abc", 1, 10, 1, "a"}, // test 10 [1:10)
		&Model{"abc", "abc", 1, 1, -1, "a"}, // test -1 (-1:1]
		&Model{"abc", "abc", 1, 1, 1, "ab"}, // test ab [1:1]
		&Model{"abc", "abc", 1, 1, -1, ""},  // test "" [1:1]
	}
	for i, c := range passCases {
		if err := Valid(c); err != nil {
			t.Fatal("pass cases", i, err)
		}
	}
	for i, c := range errCases {
		if err := Valid(c); err == nil {
			t.Fatal("err cases", i)
		}
	}
}
