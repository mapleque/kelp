package web

import (
	"net/http"
)

type Server struct {
	host   string
	router *Router

	session *_SessionServer
}

func New(host string) *Server {
	return &Server{
		host: host,
		router: &Router{
			path:         "",
			realPath:     host,
			method:       "",
			handlerChain: []HandlerFunc{},
			children:     []*Router{},
		},
	}
}

func (this *Server) Run() {
	err := http.ListenAndServe(this.host, this)
	if err != nil {
		panic("web server start faild " + err.Error())
	}
}

func (this *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	this.handle(c)
}

func (this *Server) handle(c *Context) {
	httpMethod := c.Request.Method
	path := c.Request.URL.Path
	router, params := this.router.find(httpMethod, path)
	if router == nil || len(router.handlerChain) <= 0 {
		c.DieWithHttpStatus(404)
		return
	}
	c.Params = params
	c.handlerChain = router.handlerChain
	c.handlerIndex = 0
	c.handlerChain[c.handlerIndex](c)
}

func (this *Server) Group(path string) *Router {
	return this.router.Group(path)
}

func (this *Server) Use(handler HandlerFunc) *Router {
	return this.router.Use(handler)
}

func (this *Server) GET(path string, handler ...HandlerFunc) *Router {
	return this.router.GET(path, handler...)
}

func (this *Server) POST(path string, handler ...HandlerFunc) *Router {
	return this.router.POST(path, handler...)
}

func (this *Server) PUT(path string, handler ...HandlerFunc) *Router {
	return this.router.PUT(path, handler...)
}

func (this *Server) DELETE(path string, handler ...HandlerFunc) *Router {
	return this.router.DELETE(path, handler...)
}
