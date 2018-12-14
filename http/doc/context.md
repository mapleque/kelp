Context
====

Context是kelp/http包中所定义的用来承载请求上下文的对象。

在Handler的参数中，可以声明使用Context：
```
func MyHandler(c *http.Context) {}
func MyHandler(in interface{}, out interface{}, c *http.Context) {}
```

使用Context，可以在以下三个方面得到扩展：
- 请求相关 获得更多的请求信息
- 调用链相关 控制调用链和在调用链之间传递数据
- 返回相关 更灵活的定制返回信息

请求相关
----

获取标准包的Request对象：

```
req := c.Request
```

> 获取的Request对象可以调用标准包的任何方法。

获取请求体：

```
body := c.Body()
```

> body是一个`[]byte`，用户可以根据自己的需要进行处理。    
> ** 注意，一旦获取过请求体，那么就无法再通过标准包的Request读取Body了。 **


获取Query参数：

```
param, exist := c.Query("args")
param := c.QueryDefault("args", "deafultValue")
params := c.QueryArray("args")
```

> Context提供了三个获取Query参数的方法，除此之外用户还可以选择通过标准包Request提供的方法获取。

获取请求路径：

```
path := c.Path()
```

> 除此之外用户还可以通过标准包Reqeust提供的方法获取。

调用链相关
----

执行调用链中下一个Handler：
```
c.Next()
```
> 这个方法常用于中间件中。

给调用链中后面的Handler传递数据：
```
func HandlerA(c *http.Context) {
  c.Set("key", value)
}

func HandlerB(c *http.Context) {
  value, exist := c.Get("key")
}

func HandlerC(c *http.Context) {
  value := c.MustGet("key")
}

func HandlerD(c *http.Context) {
  meta := c.MetaData
  value, exist := meta[key]
}

```

> 用户通过Set方法存储数据，然后可以通过三种形式获取数据。

返回相关
----

使用预定义格式返回数据：

```
c.Text(message...)
c.Json(obj)
```

> Context提供了两种返回数据的格式，Text和Json方法分别对应`text/plain`和`text/json`。

返回错误Http Status：

```
c.DieWithHttpStatus(400)
```

> 可以使用DieWithHttpStatus方法返回指定的Http Status。

返回重定向：

```
c.Redirect(302, "redirect_location")
```

> 使用Redirect方法返回重定向标记，用户可以选择300到308之间的任何status返回。

自定义返回：

```
c.ManuResponse = true
writer := c.ResponseWriter
// Do anything you want with writer
// may like follow:
// writer.Header().Add("Context-Type", "text/plain;charset=UTF-8")
// writer.Writer("some thing your want to response")
```

> 其中ResponseWriter是标准包中的对象，用户可以使用该对象的任何方法。    
> ** 注意，这里一定要将ManuResponse标记设置为true，否则框架将会对返回数据进行格式化。 **

相关链接
----

- [阅读上一章：参数校验](/http/doc/validator.md)
- [阅读下一章：生成接口文档](/http/doc/apidoc.md)
- [返回包简介](/http/README.md)
- [返回示例example讲解](/http/example/README.md)
