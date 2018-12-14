日志
====

kelp/http包默认会在下面几种情况下输出日志：
- 如果使用了LogHandler中间件，在系统接收到请求时输出REQ日志，参考[Handler & 中间件](/http/doc/handler.md)中的介绍。
- 如果使用了RecoverHandler中间件，那么当系统panic的时候，会输出ERROR日志。
- 在系统注册路由时，如果注册失败，会输出DEBUG日志。

主动输出日志
----

kelp/http包提供了输出四个级别日志的方法：

```
http.Info(msg ...interface{})
http.Warn(msg ...interface{})
http.Error(msg ...interface{})
http.Debug(msg ...interface{})
```

用户可以直接使用这些方法输出日志，也可以使用自己的日志模块输出日志。

输出重定向
----

kelp/http包中，定义了一个用于输出日志的接口。

```
type loggerer interface {
	Log(tag string, msg ...interface{})
}
```

用户只需要实现该接口，就可以日志输出重定向到自己实现的组件中。

```
http.SetLogger(logger)
```

相关链接
----

- [阅读上一章：生成接口文档](/http/doc/apidoc.md)
- [阅读下一章：Client](/http/doc/client.md)
- [返回包简介](/http/README.md)
- [返回示例example讲解](/http/example/README.md)
