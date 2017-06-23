package web

import (
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	host string
	mux  *http.ServeMux
}

type Handler func(context *Context)

type Context struct {
	req *http.Request
	w   http.ResponseWriter

	startTime time.Time

	Host   string
	Path   string
	Params map[string]interface{}
	Data   interface{}
	Status int
	Errmsg string

	Raw bool
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
		context := &Context{
			req:       req,
			w:         w,
			startTime: time.Now(),
			Host:      req.Host,
			Path:      path,
			Params:    processParams(req),
			Status:    0,
			Raw:       false,
		}
		handler(context)

		if context.Raw { // Raw response
			w.Write([]byte(context.Data.(string)))
		} else { // format response
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
				log.Error("response encode error", err.Error())
			}
			w.Write(out)
		}
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
	})
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
			} else {
				params[k] = dat
			}
		} else {
			params[k] = ""
		}
	}
	return params
}
