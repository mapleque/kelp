package http

import (
	"strings"
)

type Router struct {
	title        string
	comment      string
	path         string
	realPath     string
	method       string
	handlerChain []HandlerFunc
	children     []*Router
}

func (this *Router) Group(path string) *Router {
	router := &Router{
		path:         path,
		realPath:     this.realPath + path,
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

func (this *Router) Handle(title, path string, handlers ...HandlerFunc) *Router {
	if len(path) < 1 || path[0] != '/' || strings.Contains(path, "//") {
		panic("add router faild, invalid path " + path)
	}
	if sepIndex := strings.Index(path[1:], "/") + 1; sepIndex > 1 {
		root := path[:sepIndex]
		subpath := path[sepIndex:]
		var group *Router = nil
		for _, router := range this.children {
			if router.path == root {
				group = router
			}
		}
		if group == nil {
			group = this.Group(root)
		}
		return group.Handle(title, subpath, handlers...)
	}
	handlerChain := append([]HandlerFunc{}, this.handlerChain...)
	handlerChain = append(handlerChain, handlers...)
	router := &Router{
		title:        title,
		path:         path,
		realPath:     this.realPath + path,
		handlerChain: handlerChain,
		children:     []*Router{},
	}
	this.children = append(this.children, router)
	log.Debug("add router", router.realPath)
	return router
}

func (this *Router) Comment(comment string) *Router {
	this.comment = comment
	return this
}

func (this *Router) find(path string) *Router {
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
	if len(path) < 1 || path[0] != '/' || strings.HasPrefix(path, "//") {
		log.Debug("invalid path", path)
		return nil
	}
	// path should not contain chars
	if strings.ContainsAny(path, "\"\"'%&();+[]{}:*<>=") {
		log.Debug("illegal path charactor", path)
		return nil
	}
	sepIndex := strings.Index(path[1:], "/") + 1
	if sepIndex < 1 {
		// find in this level
		for _, router := range this.children {
			if router.path == path {
				return router
			}
		}
	} else {
		root := path[:sepIndex]
		subpath := path[sepIndex:]
		// find in next level
		for _, router := range this.children {
			if router.path == root {
				subrouter := router.find(subpath)
				return subrouter
			}
		}
	}
	return nil
}
