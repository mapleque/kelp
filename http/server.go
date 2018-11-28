package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"
)

var (
	SERVICE_VERSION string
	KELP_VERSION    string
	BUILD_TIME      string
	GO_VERSION      string
)

type Server struct {
	host   string
	router *Router

	start time.Time
}

func New(host string) *Server {
	return &Server{
		host: host,
		router: &Router{
			path:         "",
			realPath:     "",
			method:       "",
			handlerChain: []HandlerFunc{},
			children:     []*Router{},
		},
	}
}

func (this *Server) RunTest() *httptest.Server {
	return httptest.NewServer(this)
}

func (this *Server) Run() {
	this.start = time.Now()
	err := http.ListenAndServe(this.host, this)
	if err != nil {
		panic(err)
	}
}

func (this *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	this.handle(c)
}

func (this *Server) handle(c *Context) {
	defer c.response()
	path := c.Request.URL.Path
	if path == "/_kelp/metric" {
		this.metric(c)
		return
	}
	router := this.router.find(path)
	if router == nil || len(router.handlerChain) <= 0 {
		c.Json(STATUS_NOT_FOUND)
		return
	}
	c.handlerChain = router.handlerChain
	c.handlerIndex = 0
	processHandlerFunc(c.handlerChain[c.handlerIndex], c)
}

func (this *Server) Group(path string) *Router {
	return this.router.Group(path)
}

func (this *Server) Use(handler HandlerFunc) *Router {
	return this.router.Use(handler)
}

func (this *Server) Handle(comment, path string, handler ...HandlerFunc) *Router {
	return this.router.Handle(comment, path, handler...)
}

func (this *Server) metric(c *Context) {
	ret := map[string]interface{}{}
	ret["hostname"] = os.Getenv("HOSTNAME")
	ret["listening"] = this.host
	if pwd, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		ret["pwd"] = pwd
	} else {
		ret["pwd"] = err
	}
	ret["args"] = os.Args
	ret["last_start_at"] = this.start.Format("2006-01-02 15:04:05")
	ret["running_seconds"] = time.Now().Sub(this.start).Seconds()
	ret["service_version"] = SERVICE_VERSION
	ret["kelp_version"] = KELP_VERSION
	ret["go_version"] = GO_VERSION
	ret["build_time"] = BUILD_TIME
	c.Json(ret)
}
