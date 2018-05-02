package web

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
	"reflect"
)

func (this *Server) GRPC(path string, grpcHandler interface{}) *Router {
	return this.router.GRPC(path, grpcHandler)
}

func (this *Router) GRPC(path string, grpcHandler interface{}) *Router {
	handler := buildHandlerFromGrpcHandler(grpcHandler)
	return this.handle(http.MethodPost, path, handler)
}

func buildHandlerFromGrpcHandler(grpcHandler interface{}) HandlerFunc {
	handlerType := reflect.TypeOf(grpcHandler)
	handler := reflect.ValueOf(grpcHandler)
	paramType := handlerType.In(1)
	param := reflect.New(paramType.Elem()).Interface()
	return func(c *Context) {
		if err := json.Unmarshal(c.Body, &param); err != nil {
			log.Error("grpc param bind error", err)
			c.DieWithHttpStatus(400)
			return
		}
		ret := handler.Call([]reflect.Value{
			reflect.ValueOf(context.Background()),
			reflect.ValueOf(param),
		})
		resp := ret[0].Interface()
		err := ret[1].Interface()
		if err != nil {
			log.Error("grpc deal error", err)
			c.DieWithHttpStatus(500)
			return
		}
		c.Success(resp)
	}
}
