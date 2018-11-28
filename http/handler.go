package http

import (
	"reflect"
)

// HandlerFunc 是一个重载函数，可以有如下形式：
// func ()
// func () (err *<struct>)
// func (c *Context)
// func (c *Context) (err *<struct>)
// func (in *<struct>)
// func (in *<struct>) (err *<struct>)
// func (in *<struct>, out *<struct>)
// func (in <non-ptr>, out *<struct>)
// func (in *<struct>, out *<struct>) (err *<struct>)
// func (in <non-ptr>, out *<struct>) (err *<struct>)
// func (in *<struct>, out *<struct>, c *Context)
// func (in <non-ptr>, out *<struct>, c *Context)
// func (in <non-ptr>, out <non-ptr>, c *Context)
// func (in *<struct>, out <non-ptr>, c *Context)
// func (in *<struct>, out *<struct>, c *Context) (err *<struct>)
// func (in <non-ptr>, out *<struct>, c *Context) (err *<struct>)
// func (in <non-ptr>, out <non-ptr>, c *Context) (err *<struct>)
// func (in *<struct>, out <non-ptr>, c *Context) (err *<struct>)
//
// 其中：
//
// in 表示输入参数
//   - 如果是一个结构体指针，则会在handler调用之前将请求的body以json的形式绑定到in上，如果有valid tag，还会做参数校验
//   - 如果不是指针，则在handler调用之前忽略这个参数，通常用于占位
//
// out 表示输出参数
//   - 如果是一个结构体指针，则会在handler调用之前将其初始化，并且在handler调用之后将它作为返回数据返回
//
// c 表示Context实例
//   - 可以在第一个或者第三个参数位置获取到
//
// err 表示错误信息，必须是个结构体指针
//   - 如果有err，则会将此error直接json格式化后返回
//   - 在handler chain中后面的error会覆盖前面的error
type HandlerFunc interface{}

// processHandlerFunc 用反射实现了HandlerFunc的重载
func processHandlerFunc(handlerFunc interface{}, c *Context) {
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		panic("handler type must be func but " + handlerType.Name())
	}
	handler := reflect.ValueOf(handlerFunc)
	var args []reflect.Value
	var response reflect.Value
	hasResponse := false
	switch handlerType.NumIn() {
	case 0:
		args = []reflect.Value{}
	case 1:
		paramType := handlerType.In(0)
		param := reflect.New(paramType.Elem()).Interface()
		// 先判断是不是context
		if paramType.Elem().Name() == "Context" {
			// 如果是
			args = []reflect.Value{
				reflect.ValueOf(c),
			}
		} else {
			// 如果不是，就是in
			if err := c.BindAndValidJson(param); err != nil {
				c.Json(StatusInvalidParam(err))
				return
			}
			args = []reflect.Value{
				reflect.ValueOf(param),
			}
		}
	case 2:
		paramType := handlerType.In(0)
		var param reflect.Value
		if paramType.Kind() == reflect.Ptr {
			param = reflect.New(paramType.Elem())
			if err := c.BindAndValidJson(param.Interface()); err != nil {
				c.Json(StatusInvalidParam(err))
				return
			}
		} else {
			param = reflect.New(paramType)
		}

		responseType := handlerType.In(1)
		if responseType.Kind() == reflect.Ptr {
			hasResponse = true
			response = reflect.New(responseType.Elem())
		} else {
			response = reflect.New(responseType)
		}

		args = []reflect.Value{
			reflect.ValueOf(param.Interface()),
			reflect.ValueOf(response.Interface()),
		}
	case 3:
		paramType := handlerType.In(0)
		var param reflect.Value
		if paramType.Kind() == reflect.Ptr {
			param = reflect.New(paramType.Elem())
			if err := c.BindAndValidJson(param.Interface()); err != nil {
				c.Json(StatusInvalidParam(err))
				return
			}
		} else {
			param = reflect.New(paramType)
		}

		responseType := handlerType.In(1)
		if responseType.Kind() == reflect.Ptr {
			hasResponse = true
			response = reflect.New(responseType.Elem())
		} else {
			response = reflect.New(responseType)
		}

		args = []reflect.Value{
			reflect.ValueOf(param.Interface()),
			reflect.ValueOf(response.Interface()),
			reflect.ValueOf(c),
		}
	default:
		panic("illegal handler define: " + handler.String())
	}
	ret := handler.Call(args)
	if len(ret) < 1 {
		if hasResponse {
			c.Json(map[string]interface{}{
				"status": 0,
				"data":   response.Interface(),
			})
		} else if !c.HasResponse {
			c.Json(STATUS_SUCCESS)
		}
	} else {
		if !ret[0].IsNil() {
			c.Json(ret[0].Interface())
		} else {
			if hasResponse {
				c.Json(map[string]interface{}{
					"status": 0,
					"data":   response.Interface(),
				})
			} else if !c.HasResponse {
				c.Json(STATUS_SUCCESS)
			}
		}
	}
}
