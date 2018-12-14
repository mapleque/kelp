package http

import (
	"strings"
)

// Router Router is a tree indexing by path,
// holding the handler chain for request processing.
// The root is Router with `/` path and children are Routers with subpath.
type Router struct {
	title        string
	comment      string
	path         string
	realPath     string
	method       string
	handlerChain []HandlerFunc
	children     []*Router
}

// Group Group is a Router node, which children are Routers.
// Every Router can create Groups as children.
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

// Use Use register a middleware on the Router, which will work on all children.
func (this *Router) Use(handler HandlerFunc) *Router {
	this.handlerChain = append(this.handlerChain, handler)
	for _, router := range this.children {
		router.handlerChain = append(router.handlerChain, handler)
	}
	return this
}

// Handle Handle register a handler on the Router.
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
	log.Log("DEBUG", "add router", router.realPath)
	return router
}

// Comment Comment add comment on Router, using in doc.
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
		log.Log("DEBUG", "invalid path", path)
		return nil
	}
	// path should not contain chars
	if strings.ContainsAny(path, "\"\"'%&();+[]{}:*<>=") {
		log.Log("DEBUG", "illegal path charactor", path)
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
