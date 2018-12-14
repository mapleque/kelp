Http Package
====
[![godoc reference](https://godoc.org/github.com/mapleque/kelp/http?status.svg)](http://godoc.org/pkg/github.com/mapleque/kelp/http)

本组件主要用于快速实现一个http服务，全部基于go基础包实现，不需要额外引用任何包。

如何开始
----

使用kelp自带的工具创建项目：

```
# add $GO_PATH/bin to your $PATH
go get github.com/mapleque/kelp
kelp create github.com/your_account/hello -http
```

或者参考[kelp/http/example](./example/)

文档说明
----

- [路由](/http/doc/router.md) 路由、路由组、链式调用
- [Handler & 中间件](/http/doc/handler.md) Handler和中间件的使用方法以及常用Handler说明
- [参数校验](/http/doc/validator.md) 使用json tag进行参数校验
- [Context](/http/doc/context.md) 请求上下文，提供了更多扩展的方法
- [生成接口文档](/http/doc/apidoc.md) kelp支持通过加载用户定义的路由和Handler自动生成接口文档
- [日志](/http/doc/logger.md) 如何在程序中输出日志以及如何重定向日志输出目标
- [Client](/http/doc/client.md) 提供便捷快速的请求http服务的方法，还提供了请求基于kelp/http构建的Server的接口封装。
