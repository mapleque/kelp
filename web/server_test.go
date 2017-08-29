package web

import (
	"net/http"
	"net/url"
	"testing"
)

func TestServe(t *testing.T) {
	server.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {})
	server.RegistHandler("/test", func(context *Context) {
		context.Data = "ok"
	})
	body, err := action("/test", nil)
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"data":"ok","status":0}` {
		t.Error(string(body))
	}
	server.RegistHandler("/echo", func(context *Context) {
		context.Data = context.Params["data"]
	})
	resp := `{"filter":"filters","range":[0,10],"sort":"sorts"}`
	body, err = action("/echo", url.Values{"data": {resp}})
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"data":`+resp+`,"status":0}` {
		t.Error(string(body))
	}
	server.RegistHandlerChain("/chain", func(context *Context) bool {
		context.Data = "chain1"
		return true
	}, func(context *Context) bool {
		context.Data = context.Data.(string) + "chain2"
		return true
	})
	body, err = action("/chain", nil)
	resp = `"chain1chain2"`
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"data":`+resp+`,"status":0}` {
		t.Error(string(body))
	}
}

func TestGetField(t *testing.T) {
	params := map[string]interface{}{
		"a": "abc",
		"b": map[string]interface{}{
			"num": 1,
			"str": "sss",
		},
		"c": []interface{}{1, 2, 3},
	}
	assertEqual(1, t, getField(params, "a"), "abc")
	assertEqual(2, t, getField(params, "b.num"), 1)
	assertEqual(3, t, getField(params, "b.str"), "sss")
	assertEqual(4, t, getField(params, "c").([]interface{})[0], 1)
}
