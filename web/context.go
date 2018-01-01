package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	Body []byte

	MetaData     map[string]interface{}
	Params       Params
	handlerIndex int
	handlerChain []HandlerFunc

	Response   []byte
	Status     int
	HttpStatus int

	metaInternal map[string]interface{}
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	// 先尝试parse form解码
	req.ParseMultipartForm(1024 * 1000 * 10)
	// 再把body读出来
	body, _ := ioutil.ReadAll(req.Body)
	return &Context{
		Request:        req,
		ResponseWriter: w,
		Body:           body,
		MetaData:       make(map[string]interface{}),
		metaInternal:   make(map[string]interface{}),
	}
}

func (this *Context) Path() string {
	return this.Request.URL.Path
}

func (this *Context) Param(key string) (string, bool) {
	return this.Params.Get(key)
}

func (this *Context) Query(key string) (string, bool) {
	if values := this.QueryArray(key); len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (this *Context) QueryDefault(key string, defaultValue string) string {
	if values := this.QueryArray(key); len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

func (this *Context) QueryArray(key string) []string {
	if values, ok := this.Request.URL.Query()[key]; ok && len(values) > 0 {
		return values
	}
	return []string{}
}

func (this *Context) Set(key string, value interface{}) {
	this.MetaData[key] = value
}

func (this *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = this.MetaData[key]
	return
}

func (this *Context) MustGet(key string) (value interface{}) {
	value, _ = this.MetaData[key]
	return
}

func (this *Context) GetInt(key string) (i int) {
	if val, ok := this.Get(key); ok && val != nil {
		i, _ = val.(int)
		return
	}
	return
}

func (this *Context) GetInt64(key string) (i int64) {
	if val, ok := this.Get(key); ok && val != nil {
		i, _ = val.(int64)
		return
	}
	return
}

func (this *Context) GetString(key string) (s string) {
	if val, ok := this.Get(key); ok && val != nil {
		s, _ = val.(string)
		return
	}
	return
}

func (this *Context) GetBool(key string) (b bool) {
	if val, ok := this.Get(key); ok && val != nil {
		b, _ = val.(bool)
		return
	}
	return
}

func (this *Context) GetFloat64(key string) (f float64) {
	if val, ok := this.Get(key); ok && val != nil {
		f, _ = val.(float64)
		return
	}
	return
}

func (this *Context) Next() {
	this.handlerIndex++
	if this.handlerIndex < len(this.handlerChain) {
		this.handlerChain[this.handlerIndex](this)
	}
}

func (this *Context) Json(data interface{}) {
	this.HttpStatus = 200
	this.ResponseWriter.Header().Add("Content-Type", "text/json;charset=UTF-8")
	out, _ := json.Marshal(data)
	this.Response = out
	this.ResponseWriter.Write(out)
}

func (this *Context) Success(data interface{}) {
	this.Status = 0
	resp := map[string]interface{}{
		"status": 0,
		"data":   data,
	}
	this.Json(resp)
}

func (this *Context) Error(status int, message interface{}) {
	this.Status = status
	resp := map[string]interface{}{
		"status": status,
	}
	switch ret := message.(type) {
	case error:
		resp["message"] = ret.Error()
	default:
		resp["message"] = ret
	}
	this.Json(resp)
	log.Debug(resp)
}

func (this *Context) DieWithHttpStatus(status int) {
	this.HttpStatus = status
	this.ResponseWriter.Header().Add("Content-Type", "text/plain;charset=UTF-8")
	this.ResponseWriter.WriteHeader(status)
	log.Debug("http status", status)
}

func (this *Context) Redirect(status int, location string) {
	if status < 300 || status > 308 {
		status = 302
	}
	this.HttpStatus = status
	http.Redirect(this.ResponseWriter, this.Request, location, status)
	log.Debug("http status", status)
}

func (this *Context) LogDebug(message ...interface{}) {
	log.Debug(message...)
}

func (this *Context) LogInfo(message ...interface{}) {
	log.Info(message...)
}

func (this *Context) LogWarn(message ...interface{}) {
	log.Warn(message...)
}

func (this *Context) LogError(message ...interface{}) {
	log.Error(message...)
}
