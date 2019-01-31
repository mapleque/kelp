package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	MetaData     map[string]interface{}
	handlerIndex int
	handlerChain []HandlerFunc

	ManuResponse     bool
	HasResponse      bool
	Response         []byte
	HttpStatus       int
	ContentType      string
	RedirectLocation string

	metaInternal *sync.Map

	body        []byte
	hasReadBody bool
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	c := &Context{
		Request:        req,
		ResponseWriter: w,
		MetaData:       make(map[string]interface{}),
		metaInternal:   new(sync.Map),
	}
	return c
}

func (this *Context) Body() []byte {
	if !this.hasReadBody {
		body, _ := ioutil.ReadAll(this.Request.Body)
		this.hasReadBody = true
		this.body = body
	}
	return this.body
}

func (this *Context) Path() string {
	return this.Request.URL.Path
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

func (this *Context) Next() {
	this.handlerIndex++
	if this.handlerIndex < len(this.handlerChain) {
		processHandlerFunc(this.handlerChain[this.handlerIndex], this)
	}
}

func (this *Context) Text(msg ...interface{}) {
	this.HasResponse = true
	this.HttpStatus = 200
	this.ContentType = "text/plain;charset=UTF-8"
	out := []byte(fmt.Sprint(msg...))
	this.Response = out
}

func (this *Context) Json(data interface{}) {
	this.HasResponse = true
	this.HttpStatus = 200
	this.ContentType = "text/json;charset=UTF-8"
	out, _ := json.Marshal(data)
	this.Response = out
}

func (this *Context) DieWithHttpStatus(status int) {
	this.HasResponse = true
	this.HttpStatus = status
	this.ContentType = "text/plain;charset=UTF-8"
}

func (this *Context) Redirect(status int, location string) {
	this.HasResponse = true
	if status < 300 || status > 308 {
		status = 302
	}
	this.RedirectLocation = location
	this.HttpStatus = status
	this.ContentType = "text/plain;charset=UTF-8"
}

func (this *Context) response() {
	if this.ManuResponse {
		return
	}
	if this.ContentType != "" {
		this.ResponseWriter.Header().Add("Content-Type", this.ContentType)
	}
	if this.HttpStatus >= 300 && this.HttpStatus <= 308 {
		http.Redirect(this.ResponseWriter, this.Request, this.RedirectLocation, this.HttpStatus)
	} else if this.HttpStatus == 200 {
		this.ResponseWriter.Write(this.Response)
	} else {
		this.ResponseWriter.WriteHeader(this.HttpStatus)
	}
}
