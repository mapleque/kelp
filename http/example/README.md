kelp/http/example
====

这是一个kelp/http包使用的示例，实现了一个todo list的http服务，提供了增删改查四个接口，旨在帮助您更深入的了解全部kelp/http包所实现的功能。

同时本示例也提供了一种组织代码的形式，供您借鉴参考。

代码概述
----

在[example/server.go](example.go)中，定义并实现了一个http server。
这个server对外提供NewServer、Init和Run接口，
其中，New接口用于设置当前http服务运行时所依赖的其他服务，比如用于数据持久化的mysql（kelp包提供了mysql操作的封装，参考[kelp/mysql](/kelp/mysql)）;
Init接口对http服务进行设置，包括监听端口、路由、中间件等;
Run接口启动http服务。

在本示例中还实现了一个todo list的数据管理服务。
由于该服务比较简单，它被直接放在了example包中[example/todolist.go](todolist.go)，对于复杂的服务也可以考虑在example创建子包并实现。

在example包中，子包main用于启动服务。下文将按照[example/main/main.go](main/main.go)中执行的顺序进行讲解。

服务创建
----

在创建服务时，可以将一些服务所需的配置作为参数传入。
当服务配置较为复杂时，也可以考虑单独实现一些用于配置的接口。
在实际应用中，配置项通常通过配置文件或者环境变量读取
（kelp的config包提供了读取配置文件和环境变量的方法，参考[kelp/config](/kelp/config)）。

在本示例中：

```
// example/main/main.go:8
ts := service.NewTodoList()
ss := service.NewServer(ts)
```

首先通过`service.NewTodoList`方法创建了一个todo list的数据管理服务实例，该服务实现了todo数据的存储功能，并提供了增删改查四个接口。

然后通过`service.NewServer`方法创建http服务实例，并将todo list数据管理服务实例绑定到http服务上。

服务初始化
----

服务初始化主要是实现一些http服务所需要的基本设置。

在本示例中，首先使用kelp/http包创建一个http server的实例，在创建实例时指定服务所监听的地址和端口：

```
// example/main/main.go:10
ss.Init("0.0.0.0:9999")

// example/server.go:25
server := http.New(host)

```

然后给这个http server指定一系列基础中间件（kelp/http包提供了一些常用中间件的封装，参考[中间件说明文档](/kelp/doc/middleware.md)）：

```
// exmaple/server.go:27
server.Use(http.LogHandler)
server.Use(http.RecoveryHandler)
server.Use(http.TraceHandler)
```

其中：
- `LogHandler`用于记录请求日志
- `RecoveryHandler`的作用是当业务代码运行过程中出现panic的时候捕获异常并返回500错误，从而避免整个http服务因此而停止运行。
- `TraceHandler`会在请求头部增加一个`traceid`，当多个http服务互相调用时，可以通过这个`traceid`进行请求的跟踪。

接着注册一些列路由（kelp/http包实现了路由、路由组、链式调用等功能，参考[路由说明文档](/kelp/doc/router.md)）：

```
// example/server.go:36
this.initRouter()

// example/server.go:51
func (this *Server) initRouter() {
  ...
}
```

其中，四个路由使用了四个不同的handler负责处理请求，handler的实现在`example/handler.go`中（kelp/http包支持多种形式的handler，参考[Handler说明文档](/kelp/doc/handler.md)）。

在本示例中，对于路由组`todo`，单独注册了一个用于认证的中间件`this.Auth`，其实现在`example/middlewares.go`中。这种注册中间件的方式，可以做到让其只对形如`/todo/*`的请求起作用。

最后自定义一些校验函数（kelp/http包提供了基于struct tag的参数校验机制，参考[参数校验说明文档](/kelp/doc/validator.md)）：
```
// example/server.go:38
this.initValidator()

// example/server.go:84
http.RegisterValidFunc("datetime", http.ValidRegexpWrapper(`^\d{4}-\d{2}-\d{2} \d{2}\:\d{2}\:\d{2}$`))
```

在本示例中，定义了一个`datatime`的校验函数。
该函数简单校验了日期时间的数据格式。
在接口参数定义中使用了该函数：

```
// example/handler.go:10
AlertTime string `json:"alert_time" valid:"optional,@datetime,message=alert_time不合法"`
```

以上，就是一个http服务常用的初始化过程。

在代码组织方面，有以下一些建议：
- 如果handler比较少且实现简单，可以考虑将代码放在`server.go`中
- 如果handler很多，也可以考虑按照路由组拆分为一系列形如`handler_todo.go`的文件
- 如果需要定义很多路由，可以考虑将路由初始化部分单独拆分到`router.go`文件或者形如`router_todo.go`的文件系列中
- 如果自定义的校验函数很多，可以考虑将其查分到`validator.go`文件中
- 不要按照传统的MVC模式组织代码，那会增加代码复杂度和维护成本
- 使用形如`TodoCreateParam`和`TodoRetrieveResponse`的模式命名参数和返回的数据结构类型
- 避免使用组合的模式定义参数，大多情况下组合依赖不利于代码的扩展和维护
- 将一些常量单独定义到一个文件中，方便查找和维护，如`status.go`

接口文档
----

kelp/http包支持通过代码生成文档，参考[生成接口文档说明](/kelp/doc/apidoc.md)。

```
// example/main/main.go:11
ss.Doc("./api.md")
```

本示例中，每次启动服务的时候，都将会在运行路径下生成一个`api.md`的接口文档。

接口文档的内容来自：
- 注册路由方法调用时所传入的`title`和`path`参数
- 路由实例调用Comment方法所传入的`comment`参数
- handler的`in`和`out`参数数据结构中所定义的`json`、`valid`、`comment`等struct tag

运行服务
----

```
// example/main/main.go:12
ss.Run()
```

这是一个同步阻塞方法，除非发送信号量或者服务崩溃，否则服务会一直执行永不停歇。

其他功能
----

- [logger](/kelp/doc/logger.md) 日志输出相关的方法封装
- [client](/kelp/doc/client.md) 一个http客户端的实现和方法封装
- [crypto](/kelp/doc/crypto.md) 一些随机和加解密方法的实现
