路由
====

在kelp/http包中，可以通过Handle方法注册路由。

```
server := http.New(host)
server.Handle("接口名称", "/api_path", Handler)
```

上面例子中，就注册了一个`http://host/api_path`的路由，开发者可以通过实现Handler处理该请求（Handler的实现方法参考[handler](/http/doc/handler.md)）。

通过Handle方法注册的路由支持全部http的方法（GET、POST等），如果需要针对不同http方法分别处理，可以在Handler中自行实现。

路由的`path`参数，必须以`/`开头，且不能含有url中的非法字符，可以设置多层紫路径。例如：

- `/`
- `/abc`
- `/abc/`
- `/abc/def`

路由组
----
使用路由组可以方便的将路由进行分组管理，有利于代码阅读和维护。

通过Group方法可以创建路由组。

```
user := server.Group("/user")
{
  user.Handle("添加用户", "/add", UserAddHandler)
}
```

上面例子中，添加用户的接口url为：`http://host/user/add`。

任意一个路由都可以创建路由组，例如：

```

wallet := user.Group("/wallet")
money := wallet.Group("/money")
{
  money.Handle("查询用户余额", "/get", UserMoneyGetHandler)
}
```

上面例子中，查询用户余额url为：`http://host/user/wallet/money/get`。


链式调用
----

在注册路由时，可以传入多个Handler。
在处理请求的过程中，第一个传入的Handler将会被首先调用。
此时通过调用Context.Next()方法可以触发进入后面第二个Handler，后面以此类推（更多Context的使用方法方法参考[context](/http/doc/context.md)）。
如果没有调用Next方法，该请求的处理将会到此结束，并按照之前所有Handler处理的最终状态进行返回。

```
mobile := user.Group("/mobile")
{
  mobile.Handle("查询用户手机号", "/get", UserCheckHandler, UserMobileGetHandler)
}

func UserCheckHandler(c *http.Context) *http.Status {
  if err := UserCheck(c.Header.Get("Authorization")); err != nil {
    return http.ErrorStatus(1, err)
  }
  c.Next()
}

func UserMobileGetHandler(in interface{}, out *UserMobileResponse) *http.Status{
  // ...
}
```

上面例子中，当用户权限校验（UserCheck）不通过时，将直接返回错误，UserMobileGetHandler方法不会被调用。
如果校验通过，则会继续调用UserMobileGetHandler方法处理请求。

需要特别说明的是：在使用链式调用模式时，一旦产生返回（正常或异常），该数据都会被立即写入返回流，且不可撤销和重写。
所以注意不要多次写入返回数据。

此外，Context.Next方法可以在Handler的任意位置调用，例如本包[`handler_helper.go`](/http/handler_helper.go)中实现的LogHandler通过在中间调用Next方法，达到记录请求响应时间的目的。

相关链接
----

- [阅读下一章：Handler & 中间件](/http/doc/handler.md)
- [返回包简介](/http/README.md)
- [返回示例example讲解](/http/example/README.md)
