http客户端
====

直接发送请求
----
使用Request方法，直接发送请求：

```
response, body, err := http.Request("http://localhost/hello", "POST", []byte(`{"name":"jack"}`))
```

其中，response是标准包http的Response对象，它的Body已经被读取到返回值body中。


构建服务对象
----

通过构建服务对象来发送请求，可以使服务的概念更清晰：

```
myService := http.NewClient("http://host")
response, body, err := myService.Request("/hello", "POST", []byte(`{"name":"jack"}`))
```

扩展请求体
----

如有需要，可以对Request进行扩展：

```
myService := http.NewClient("http://host")
req, err := myService.BuildRequest("/hello", "POST", []byte(`{"name":"jack"}`)

// TODO check err
req.Header.Set("My-Header", "value")
response. body, err := myService.Do(req)

```

这里通过BuildRequest方法创建的Request是标准包http的Request对象，用户可以对其进行随意的修改定制。

当然，用户也可以直接使用标准包http创建Request对象。

请求kelp
----

在kelp/http包中，额外提供了直接请求通过kelp/http创建的http server(server创建参考[kelp/httpexample](/http/example/))的方法：

```
// use package method
status, err := http.RequestKelp("http://host/hello", "my token if exist", in, out, lastContext)

// use client service mode
myService := http.NewKelpClient("http://host", "my token if exist")
status, err := myService.RequestKelp("/hello", in, out, lastContext)
```

其中in和out按照Handler中定义的数据结构定义即可。

当请求失败时，status和err都可能非空，注意分情况判断。
如果status非空，那么它一定是server中定义的Status对象（这里要求Handler的返回值必须使用kelp/http包提供的Status的形式定义）。
