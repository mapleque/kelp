package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	host string
	mux  *http.ServeMux
}

type Handler func(context *Context)
type HandlerChain func(context *Context) bool

type Context struct {
	req *http.Request
	w   http.ResponseWriter

	startTime time.Time

	Host        string
	Path        string
	Params      map[string]interface{}
	TransParams map[string]interface{}
	Data        interface{}
	Status      int
	Errmsg      string

	Raw    bool
	Header map[string]string
}

func New(host string) *Server {
	return &Server{
		host: host,
		mux:  http.NewServeMux()}
}

func (server *Server) Run() {
	log.Info("web listening on", server.host)
	err := http.ListenAndServe(server.host, server.mux)
	if err != nil {
		panic("web server start faild " + err.Error())
	}
}

func (server *Server) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	server.mux.HandleFunc(path, handler)
}

func (server *Server) RegistHandler(path string, handler Handler) {
	server.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		context := newContext(path, w, req)
		handler(context)
		processResponse(context)
		logHandler(context)
	})
}

func (server *Server) RegistHandlerChain(path string, handlers ...HandlerChain) {
	server.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		context := newContext(path, w, req)
		for _, handler := range handlers {
			if !handler(context) {
				break
			}
		}
		processResponse(context)
		logHandler(context)
	})
}

func newContext(path string, w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		req:         req,
		w:           w,
		startTime:   time.Now(),
		Host:        req.Host,
		Path:        path,
		Params:      processParams(req),
		TransParams: make(map[string]interface{}),
		Status:      0,
		Raw:         false,
		Header:      make(map[string]string),
	}
}

func processParams(req *http.Request) map[string]interface{} {
	err := req.ParseForm()
	if err != nil {
		log.Error("request decode error", err.Error())
	}
	params := make(map[string]interface{})
	for k, vs := range req.Form {
		if len(vs) > 0 {
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(vs[0]), &dat); err != nil {
				params[k] = vs[0]
				log.Debug("param parse warn:", err)
			} else {
				params[k] = dat
			}
		} else {
			params[k] = ""
		}
	}
	return params
}

func processResponse(context *Context) {
	for key, value := range context.Header {
		context.w.Header().Add(key, value)
	}
	if context.Raw { // Raw response
		context.w.Write([]byte(context.Data.(string)))
	} else { // format response
		// default json
		context.w.Header().Add("Content-Type", "text/json;charset=UTF-8")

		resp := make(map[string]interface{})
		resp["status"] = context.Status
		if context.Status == 0 { // success response
			if context.Data != nil {
				resp["data"] = context.Data
			}
		} else { // error response
			resp["errmsg"] = context.Errmsg
		}
		out, err := json.Marshal(resp)
		if err != nil {
			// acturally, there will never error here
			log.Error("response encode error", err.Error())
		}
		context.w.Write(out)
	}
}

func logHandler(context *Context) {
	if context.Raw {
		log.Info(
			context.Status,
			time.Now().Sub(context.startTime),
			context.Host,
			context.Path,
			context.Params,
			context.Status,
			"RawData",
			context.Errmsg,
		)
	} else {
		log.Info(
			context.Status,
			time.Now().Sub(context.startTime),
			context.Host,
			context.Path,
			context.Params,
			context.Status,
			context.Data,
			context.Errmsg,
		)
	}
}

func (context *Context) GetParam(field string) interface{} {
	if data, ok := context.Params["data"]; ok {
		return getField(data.(map[string]interface{}), field)
	}
	return getField(context.Params, field)
}

func (context *Context) GetTrans(field string) interface{} {
	return getField(context.TransParams, field)
}

func (context *Context) GetCookie(name string) string {
	cookie, err := context.req.Cookie(name)
	log.Debug("get cookie", context.req.Cookies())
	if err != nil {
		log.Warn("read cookie error", err)
		return ""
	}
	return cookie.Value
}

func (context *Context) SetCookie(name, value string, expires time.Duration) {
	now := time.Now()
	cookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Expires: now.Add(expires),
	}
	context.req.AddCookie(cookie)
	http.SetCookie(context.w, cookie)
}

func getField(params map[string]interface{}, field string) interface{} {
	if params == nil {
		return nil
	}
	tmpArr := strings.SplitN(field, ".", 2)
	value, ok := params[tmpArr[0]]
	if !ok {
		return nil
	}
	if len(tmpArr) == 2 {
		return getField(value.(map[string]interface{}), tmpArr[1])
	}
	return value
}
