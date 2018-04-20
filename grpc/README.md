GRPC
====

依赖
----

`go get google.golang.org/grpc`
`go get golang.org/x/net`

用法
----

```
TOKEN := "input_your_auth_token_here"
HOST := "127.0.0.1:50000" // port:50000-60000
gServer := grpc.New(grpc.Recovery, grpc.Logger, grpc.TokenAuthority(TOKEN))
// serverImpl := YourServerImplement
// RegisterYourServer(gServer, serverImpl)
grpc.Run(gServer, HOST)

```

方法说明
----
- `New(handler...)`：新建grpc server实例
- `Run(gServer, host)`：运行grpc server
- `UnaryInterceptorChain(handler...)`：包装调用链
- `Recovery()`：catch panic，使系统不至于崩溃
- `Logger()`：输出请求日志
- `TokenAuthority(token)`：权限校验，简单校验token
