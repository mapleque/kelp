package web

import (
	"testing"
)

func TestSession(t *testing.T) {
	server.RegistHandlerChain("/session", func(context *Context) bool {
		if context.SessionGet("key") != nil {
			t.Fatal("session not nil")
		}
		context.SessionSet("key", "value")
		if context.SessionGet("key").(string) != "value" {
			t.Fatal("session set faild")
		}
		return true
	}, func(context *Context) bool {
		if context.SessionGet("key").(string) != "value" {
			t.Fatal("session set faild")
		}
		context.SessionDestroy()
		if context.SessionGet("key") != nil {
			t.Fatal("session destroy faild")
		}
		context.SessionSet("key", "value")
		if context.SessionGet("key") != nil {
			t.Fatal("session destroy faild")
		}
		return true
	})
	action("/session", nil)
}
