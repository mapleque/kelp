package web

import (
	"net/http"
	"strings"
)

type Param struct {
	key   string
	value string
}

type Params []Param

func (this Params) Get(key string) (string, bool) {
	for _, param := range this {
		if param.key == key {
			return param.value, true
		}
	}
	return "", false
}

type Router struct {
	path         string
	realPath     string
	method       string
	handlerChain []HandlerFunc
	children     []*Router
}

type HandlerFunc func(*Context)

func (this *Router) Group(path string) *Router {
	router := &Router{
		path:         path,
		realPath:     this.realPath + path,
		method:       "",
		handlerChain: append([]HandlerFunc{}, this.handlerChain...),
		children:     []*Router{},
	}
	this.children = append(this.children, router)
	return router
}

func (this *Router) Use(handler HandlerFunc) *Router {
	this.handlerChain = append(this.handlerChain, handler)
	for _, router := range this.children {
		router.handlerChain = append(router.handlerChain, handler)
	}
	return this
}

func (this *Router) GET(path string, handlers ...HandlerFunc) *Router {
	return this.handle(http.MethodGet, path, handlers...)
}

func (this *Router) POST(path string, handlers ...HandlerFunc) *Router {
	return this.handle(http.MethodPost, path, handlers...)
}

func (this *Router) PUT(path string, handlers ...HandlerFunc) *Router {
	return this.handle(http.MethodPut, path, handlers...)
}

func (this *Router) DELETE(path string, handlers ...HandlerFunc) *Router {
	return this.handle(http.MethodDelete, path, handlers...)
}

func (this *Router) handle(method string, path string, handlers ...HandlerFunc) *Router {
	if len(path) < 1 || path[0] != '/' || strings.Contains(path, "//") {
		log.Error("add router faild, invalid path", path)
		return nil
	}
	if sepIndex := strings.Index(path[1:], "/") + 1; sepIndex > 1 {
		root := path[:sepIndex]
		subpath := path[sepIndex:]
		var group *Router = nil
		for _, router := range this.children {
			if router.method == "" && router.path == root {
				group = router
			}
		}
		if group == nil {
			group = this.Group(root)
		}
		return group.handle(method, subpath, handlers...)
	}
	handlerChain := append([]HandlerFunc{}, this.handlerChain...)
	handlerChain = append(handlerChain, handlers...)
	router := &Router{
		path:         path,
		realPath:     this.realPath + path,
		method:       method,
		handlerChain: handlerChain,
		children:     []*Router{},
	}
	this.children = append(this.children, router)
	log.Info("add router", router.method, router.realPath)
	return router
}

func (this *Router) find(method, path string) (*Router, Params) {
	// path should not like:
	//	1. ""
	//	2. "xxx"
	//	3. "//"
	//	4. "//xxx"
	// path is ok like:
	//	1. "/"
	//	2. "/xxx"
	//	3. "/xxx/"
	//	4. "/xxx/xxx"
	//	5. "/xxx/xxx/"
	//	6. ...
	params := []Param{}
	if len(path) < 1 || path[0] != '/' || strings.HasPrefix(path, "//") {
		log.Error("invalid path", path)
		return nil, params
	}
	// path should not contain chars
	if strings.ContainsAny(path, "\"\"'%&();+[]{}:*<>=") {
		log.Error("illegal path charactor", path)
		return nil, params
	}
	sepIndex := strings.Index(path[1:], "/") + 1
	if sepIndex < 1 {
		// find in this level
		for _, router := range this.children {
			if router.method == method {
				if isParamPath(router.path) {
					params = append(params, Param{router.path[2:], path[1:]})
				}
				if isParamPath(router.path) || router.path == path {
					return router, params
				}
			}
		}
		log.Warn("router not found", this.realPath, method, path)
	} else {
		root := path[:sepIndex]
		subpath := path[sepIndex:]
		// find in next level
		for _, router := range this.children {
			if router.method == "" {
				if isParamPath(router.path) {
					params = append(params, Param{router.path[2:], root[1:]})
				}
				if isParamPath(router.path) || router.path == root {
					subrouter, subparams := router.find(method, subpath)
					return subrouter, append(params, subparams...)
				}
			}
		}
		log.Warn("group not found", this.realPath, method, root, subpath)
	}
	return nil, params
}

func isParamPath(path string) bool {
	return len(path) > 2 && path[1] == ':'
}
