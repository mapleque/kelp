Handler & 中间件
====

Handler
----
每一个路由，都需要有Handler来处理请求。kelp/http包支持多种Handler的定义形式：

```
func ()
func () (err *<struct>)
func (c *Context)
func (c *Context) (err *<struct>)
func (in *<struct>)
func (in *<struct>) (err *<struct>)
func (in *<struct>, out *<struct>)
func (in <non-ptr>, out *<struct>)
func (in *<struct>, out *<struct>) (err *<struct>)
func (in <non-ptr>, out *<struct>) (err *<struct>)
func (in *<struct>, out *<struct>, c *Context)
func (in <non-ptr>, out *<struct>, c *Context)
func (in <non-ptr>, out <non-ptr>, c *Context)
func (in *<struct>, out <non-ptr>, c *Context)
func (in *<struct>, out *<struct>, c *Context) (err *<struct>)
func (in <non-ptr>, out *<struct>, c *Context) (err *<struct>)
func (in <non-ptr>, out <non-ptr>, c *Context) (err *<struct>)
func (in *<struct>, out <non-ptr>, c *Context) (err *<struct>)
```

其中：

- `in` 表示输入参数
  - 如果是一个结构体指针，则会在handler调用之前将请求的body以json的形式绑定到in上，如果有valid tag，还会做参数校验，参考[参数校验](/http/doc/validator.md)
  - 如果不是指针，则在handler调用之前忽略这个参数，通常用于占位

- `out` 表示输出参数
  - 如果是一个结构体指针，则会在handler调用之前将其初始化，并且在handler调用之后将它作为返回数据返回
  - 如果不是指针或不存在，则返回默认的`{"status":0, "message":"成功"}`

- `c` 表示Context实例
  - 可以在第一个或者第三个参数位置获取到
  - Context的使用方法参考[Context](/http/doc/context.md)

- `err` 表示错误信息，必须是个结构体指针
  - 如果有err，则会将此error直接json格式化后返回
  - 在handler chain中后面的error会覆盖前面的error
  - 特别的，这里的err推荐使用`*http.Status`


中间件
----
中间件也是Handler。

路由和路由组都可以通过注册中间件来增加对请求处理的逻辑。

中间件常用的注册形式是：

```
server := http.New(host)

// All request will be process with ServerHandler
server.Use(ServerHandler)

// Create a resource group
resource := server.Group("/resource")

// Request with path prefix `/resource` will be process with ResourceHandler
resource.Use(ResourceHandler)

// Request `/resource/add` will be process with beforeAddHandler, addHandler
resource.Handle("add", "/add", beforeAddHandler, addHandler)
// Your can add multiple middleware before addHandler.

// Attention:
// Do not write like follow:
// ```
// resource.Handle("add", "/add", beforeAddHandler, addHandler, afterAddHandler).Use(finalAddHandler)
// ```
// Cause of afterAddHandler and finalAddHandler is the behind the addHandler, as while,
// we usually write the main logic in addHandler without calling Next.
// Therefore, afterAddHandler and finalAddHanler will be never processing.
```

所有已注册的中间件都会按照其注册顺序加载到Context中，使用Context的Next方法调用下一个中间件。

注意：请避免在主逻辑Handler之后注册中间件，如果需要，在不同的位置调用Next方法可以达到修改执行顺序的目的。
例如：
```
// Register a router
server.Handle("/example", middlewareHandler, mainHandler)

// middlewareHandler
func middlewareHandler(c *http.Context) {
  // do something before mainHandler
  c.Next() // call mainHandler
  // do something after mainHandler
}
```

常用Handler（中间件）
----

所有中间件的实现代码在[handler_helper.go](handler_helper.go)中。

### LogHandler

```
server.Use(http.LogHandler)
```

当使用了LogHandler中间件后，所有经过中间件处理的请求都会输出请求日志，按照下面的格式输出：

`$request_start [INFO] $remote_ip $request_end $latency $method $uri_path($uri_param) $traceid $uuid """$request_body""" """$repsonse_body"""`

其中：

- request_start表示请求开始时间，格式：yyyy/MM/dd HHmmss
- remote_ip表示请求方ip，注意多层代理情况下使用X-Forwarded-For的第一个值
- request_end表示请求返回时间，格式：yyyy/MM/dd HHmmss
- latency表示请求响应时长，整数，单位是毫秒(ms)
- method表示http method
- uri_path表示请求path
- uri_param表示请求参数，如果有则以?开头，如果没有则留空
- traceid表示请求链的唯一标识，如果没有则输出`-`
- uuid表示用户的唯一标识，如果没有则输出`-`
- request_body表示请求数据，如果没有则输出`-`
- response_body表示请求返回数据，如果没有则输出`-`

下面是使用logstash解析该日志的代码：

```
dissect {
  mapping => {
    "msg" => "%{request_start} [%{log_flag}] %{remote_ip} %{request_end} %{+request_end} %{latency} %{method} %{url} %{traceid} %{uuid} '''%{request_body}''' '''%{repsonse_body}'''"
  }
  remove_field => ["msg"]
}
if [remote_ip]{
  grok{
    match => { "url" => "%{URIPATH:uri_path}(?:(%{URIPARAM:uri_param})?)"  }
    remove_field => [ "url" ]
  }
}
date {
  match => ["request_end", "yyyy/MM/dd HH:mm:ss"]
}
mutate {
  convert => {
    "latency" => "integer"
  }
}

```

### RecoveryHandler

```
server.Use(http.RecoveryHandler)
```

用于捕获handler中的panic，避免由于panic造成整个服务崩溃。

所有在该中间件之后注册的Handler所产生的panic都能被捕获。

### TraceHandler

```
server.Use(http.TraceHandler)
```

在请求头自动添加Traceid。

当所开发的服务是一组微服务的入口时，可以使用这个中间件进行请求跟踪。

相关链接
----

- [阅读上一章：路由](/http/doc/validator.md)
- [阅读下一章：参数校验](/http/doc/validator.md)
- [返回包简介](/http/README.md)
- [返回示例example讲解](/http/example/README.md)
