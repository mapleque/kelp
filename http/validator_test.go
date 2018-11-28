package http

import (
	"testing"
)

type Case struct {
	assert  bool   // 预期结果
	message string // 预期提示
	json    string // 数据
}

func TestValid(t *testing.T) {
	for _, c := range []Case{
		Case{true, "", `{"name":"abc"}`},
		Case{true, "", `{"name":""}`},
		Case{false, "", `{}`},
		Case{false, "", `{"other":"abc"}`},
	} {
		assertValid(t, c, Valid(&struct {
			Name string `json:"name" valid:""` // name必填，可以为空字符串
		}{}, []byte(c.json)))
	}

	// message
	for _, c := range []Case{
		Case{false, "invalid name", `{}`},
	} {
		assertValid(t, c, Valid(&struct {
			Name string `json:"name" valid:"message=invalid name"`
		}{}, []byte(c.json)))
	}

	// optional
	for _, c := range []Case{
		Case{true, "", `{}`},
		Case{false, "", `{"name":""}`},
	} {
		assertValid(t, c, Valid(&struct {
			Name string `json:"name" valid:"(0,),optional"`
		}{}, []byte(c.json)))
	}

	// reg
	for _, c := range []Case{
		Case{true, "", `{"name":"123"}`},
		Case{false, "", `{"name":"abc"}`},
	} {
		assertValid(t, c, Valid(&struct {
			Name string `json:"name" valid:"/\\d+/"`
		}{}, []byte(c.json)))
	}
}

func assertValid(t *testing.T, c Case, err error) {
	if err != nil {
		if c.assert {
			t.Error("valid should return nil but", err)
		} else {
			if c.message != "" && c.message != err.Error() {
				t.Error(c, "valid message should be", c.message, "but", err)
			}
		}
	} else {
		if !c.assert {
			t.Error(c, "valid should return", c.message, "but nil")
		}
	}
}

type validParseAssert struct {
	isOptional bool
	message    string
	rulesValue []string
	tagsStr    string
}

func TestValidParse(t *testing.T) {
	for _, a := range []validParseAssert{
		validParseAssert{false, "", []string{}, ""},
		validParseAssert{false, "", []string{"func"}, "@func"},
		validParseAssert{false, "", []string{"func1", "func2"}, "@func1,@func2"},
		validParseAssert{false, "", []string{"1<=x<=2"}, "[1,2]"},
		validParseAssert{false, "", []string{"-2<=x<=-1"}, "[-2,-1]"},
		validParseAssert{false, "", []string{"1<x<=2"}, "(1,2]"},
		validParseAssert{false, "", []string{"1<=x<2"}, "[1,2)"},
		validParseAssert{false, "", []string{"1<x<2"}, "(1,2)"},
		validParseAssert{false, "", []string{"1<x<"}, "(1,)"},
		validParseAssert{false, "", []string{"<x<2"}, "(,2)"},
		validParseAssert{false, "", []string{".*"}, "/.*/"},
		validParseAssert{false, "", []string{`\d+`}, `/\d+/`},
		validParseAssert{true, "", []string{}, "optional"},
		validParseAssert{false, "aaa", []string{}, "message=aaa"},
		validParseAssert{true, "", []string{"func"}, "@func,optional"},
		validParseAssert{true, "", []string{"1<=x<=2"}, "[1,2],optional"},
		validParseAssert{false, "", []string{"func", "1<=x<=2"}, "@func,[1,2]"},
		validParseAssert{true, "", []string{"func", "1<=x<=2"}, "@func,[1,2],optional"},
		validParseAssert{true, "hello 123", []string{"func", "1<=x<=2"}, "@func,[1,2],optional,message=hello 123"},
	} {
		if g, err := validParse(a.tagsStr); err != nil {
			t.Error(err)
		} else {
			if a.isOptional != g.optional ||
				a.message != g.message ||
				len(a.rulesValue) != len(g.rules) {
				t.Error(a, g, a.rulesValue, len(g.rules))
			}
			for _, r := range a.rulesValue {
				if !hasRule(g.rules, r) {
					for _, tr := range g.rules {
						t.Error(r, "vs", tr.value())
					}
				}
			}
		}
	}
}

func hasRule(rules []rule, value string) bool {
	for _, rule := range rules {
		if rule.value() == value {
			return true
		}
	}
	return false
}

type structForTestBind struct {
	A int                   `json:"a" valid:"message=invalid a"`
	B structSubForTestBind  `json:"b"`
	C *structSubForTestBind `json:"c"`
}

type structSubForTestBind struct {
	d int `json:"d" valid:"message=invalid d in sub"`
}

func TestBindAndValidJson(t *testing.T) {
	tm := &structForTestBind{}
	if err := BindAndValidJson(tm, []byte(`{"a":1.1}`)); err != nil {
		if err.Error() != "invalid a" {
			t.Error("failed in bind error when unmarshal type error", err)
		}
	}
	if err := BindAndValidJson(tm, []byte(`{"a":1,"b":{"d":1.3}}`)); err != nil {
		if err.Error() != "invalid d in sub" {
			t.Error("failed in bind error when unmarshal type error", err)
		}
	}
	if err := BindAndValidJson(tm, []byte(`{"a":1,"c":{"d":1.3}}`)); err != nil {
		if err.Error() != "invalid d in sub" {
			t.Error("failed in bind error when unmarshal type error", err)
		}
	}
}
