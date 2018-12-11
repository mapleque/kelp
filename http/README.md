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

- [router](/http/doc/router.md) 路由、路由组、链式调用
- [middleware](/http/doc/middleware.md)
- [handler](/http/doc/handler.md) Handler
- [validator](/http/doc/validator.md) 参数校验
- [apidoc](/http/doc/apidoc.md) 生成接口文档
- [logger](/http/doc/logger.md) 日志输出相关的方法封装
- [client](/http/doc/client.md) 一个http客户端的实现和方法封装
- [crypto](/http/doc/crypto.md) 一些随机和加解密方法的实现
