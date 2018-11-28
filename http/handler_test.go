package http

import (
	"testing"
)

type in struct {
	Name string `json:"name"`
}
type out struct {
	Result string `json:"result"`
}

func p0() {
}
func p1() *Status {
	return STATUS_ERROR_DB
}
func p2(c *Context) *Status {
	c.Text("hello")
	return nil
}
func p3(in *in) *Status {
	return &Status{100, in.Name}
}
func p4(in *in, out *out) {
	out.Result = in.Name
}
func p5(in interface{}, out *out) {
	out.Result = "hello"
}
func p6(in *in, out *out, c *Context) {
	out.Result = in.Name + string(c.Body())
}
func p7(in *in, out interface{}, c *Context) {
	c.Text(in.Name)
}

func TestProcessHandlerFunc(t *testing.T) {
	var c *Context
	c = nc()
	processHandlerFunc(p0, c)
	assert(t, c, `{"status":0,"message":"成功"}`)
	c = nc()
	processHandlerFunc(p1, c)
	assert(t, c, `{"status":2,"message":"数据库错误"}`)
	c = nc()
	processHandlerFunc(p2, c)
	assert(t, c, "hello")
	c = nc()
	processHandlerFunc(p3, c)
	assert(t, c, `{"status":100,"message":"kelp"}`)
	c = nc()
	processHandlerFunc(p4, c)
	assert(t, c, `{"data":{"result":"kelp"},"status":0}`)
	c = nc()
	processHandlerFunc(p5, c)
	assert(t, c, `{"data":{"result":"hello"},"status":0}`)
	c = nc()
	processHandlerFunc(p6, c)
	assert(t, c, `{"data":{"result":"kelp{\"name\":\"kelp\"}"},"status":0}`)
	c = nc()
	processHandlerFunc(p7, c)
	assert(t, c, "kelp")
}

func nc() *Context {
	return &Context{
		body:        []byte(`{"name":"kelp"}`),
		hasReadBody: true,
	}

}

func assert(t *testing.T, c *Context, res string) {
	if string(c.Response) != res {
		t.Error("wrong response", string(c.Response), res)
	}
}
