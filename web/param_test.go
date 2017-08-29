package web

import (
	"net/url"
	"testing"
)

func TestInArray(t *testing.T) {
	assertEqual(1, t, InArray([]interface{}{1, 2, 3})(1), true)
	assertEqual(2, t, InArray([]interface{}{1, 2, 3})(2), true)
	assertEqual(3, t, InArray([]interface{}{1, 2, 3})(4), false)
	assertEqual(4, t, InArray([]interface{}{1, 2, 3})("aaa"), false)
}

func TestExist(t *testing.T) {
	param := "{\"key\":\"aaa\"}"
	handlerTest(
		t, "/testExist",
		func(context *Context) {
			assertEqual(1, t, Exist(context, "key"), true)
		},
		url.Values{"data": {param}},
		"{\"status\":0}",
	)
}

func TestCheck(t *testing.T) {
	param := `{
		"str":"aaa",
		"num":1,
		"enum":1,
		"enumstr":"a",
		"subset":[1,"a"]
	}`
	handlerTest(
		t, "/testCheck",
		func(context *Context) {
			assertEqual(1, t, Check(context, "str", IsString(-1, -1), nil, 3), true)
			assertEqual(2, t, Check(context, "num", IsNum(-1, -1), nil, 3), true)
			assertEqual(3, t, Check(context, "enum", InArray([]interface{}{"a", 1}), nil, 3), true)
			assertEqual(4, t, Check(context, "enumstr", InArray([]interface{}{"a", 1}), nil, 3), true)
			assertEqual(5, t, Check(context, "subset", IsSubSet([]interface{}{1, "a", 2, "b"}), nil, 3), true)

			assertEqual(11, t, Check(context, "num", IsString(-1, -1), nil, 3), false)
			assertEqual(12, t, Check(context, "str", IsNum(-1, -1), nil, 3), false)
			assertEqual(13, t, Check(context, "str", InArray([]interface{}{"a", 1}), nil, 3), false)
			assertEqual(14, t, Check(context, "str", IsSubSet([]interface{}{1, "a", 2, "b"}), nil, 3), false)

			context.Status = 0
		},
		url.Values{"data": {param}},
		"{\"status\":0}",
	)
}

func transSimple(tar string) ParamTransFunc {
	return func(context *Context, field string) {
		param := context.GetParam(field)
		if _, ok := context.TransParams[tar]; !ok {
			context.TransParams[tar] = make(map[string]interface{})
		}
		context.TransParams[tar].(map[string]interface{})[field] = param
	}
}

func TestTrans(t *testing.T) {
	param := `{
		"str":"aaa",
		"num":1,
		"enum":1,
		"enumstr":"a",
		"subset":[1,"a"]
	}`
	handlerTest(
		t, "/testTrans",
		func(context *Context) {
			assertEqual(1, t, Check(context, "str", IsString(-1, -1), transSimple("str"), 3), true)
			assertEqual(2, t, Check(context, "num", IsNum(-1, -1), transSimple("num"), 3), true)
			assertEqual(3, t, Check(context, "enum", InArray([]interface{}{"a", 1}), transSimple("num"), 3), true)
			assertEqual(4, t, Check(context, "enumstr", InArray([]interface{}{"a", 1}), transSimple("str"), 3), true)
			assertEqual(5, t, Check(context, "subset", IsSubSet([]interface{}{1, "a", 2, "b"}), transSimple("arr"), 3), true)

			assertEqual(21, t, context.GetTrans("str").(map[string]interface{})["str"], "aaa")
			assertEqual(22, t, context.GetTrans("num").(map[string]interface{})["num"], 1)
			assertEqual(23, t, context.GetTrans("num").(map[string]interface{})["enum"], 1)
			assertEqual(24, t, context.GetTrans("str").(map[string]interface{})["enumstr"], "a")
			assertEqual(25, t, context.GetTrans("arr").(map[string]interface{})["subset"].([]interface{})[1], "a")

			assertEqual(11, t, Check(context, "num", IsString(-1, -1), nil, 3), false)
			assertEqual(12, t, Check(context, "str", IsNum(-1, -1), nil, 3), false)
			assertEqual(13, t, Check(context, "str", InArray([]interface{}{"a", 1}), nil, 3), false)
			assertEqual(14, t, Check(context, "str", IsSubSet([]interface{}{1, "a", 2, "b"}), nil, 3), false)

			context.Status = 0
		},
		url.Values{"data": {param}},
		"{\"status\":0}",
	)
}
