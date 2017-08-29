package web

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

var server *Server

func init() {
	if server != nil {
		return
	}
	server = New("127.0.0.1:9999")
	go server.Run()
}

func action(path string, data url.Values) (string, error) {
	resp, err := http.PostForm("http://127.0.0.1:9999"+path, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func assertEqual(seq int, t *testing.T, a interface{}, b interface{}) {
	if !equal(a, b) {
		t.Fatal("euqual assert", seq, a, b)
	}
}

func handlerTest(t *testing.T, path string, handler func(*Context), values url.Values, resp string) {
	server.RegistHandler(path, handler)
	body, err := action(path, values)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != resp {
		t.Fatal(string(body), resp)
	}
}
