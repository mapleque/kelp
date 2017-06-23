package web

import (
	"net/http"
	"testing"
)

func TestServerRun(t *testing.T) {
	server := New("127.0.0.1:9000")
	if server == nil {
		t.Error("server is nil")
	}
	go server.Run()
	server.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {})
	server.RegistHandler("/test", func(context *Context) {
		context.Data = "ok"
	})
}
