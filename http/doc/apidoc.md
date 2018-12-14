生成接口文档
====

kelp/http包提供了根据用户编写的路由和Handler代码自动生成接口文档的功能。

```
// new server
server := http.New(host)

// register routers
initRouter(server)

server.Doc("./api.md")

// ...
```

调用Doc方法，将生成的接口文档输出到指定文件。

接口文档有以下几部分内容：
- 接口名
- 接口说明
- 请求路径
- 请求参数说明
- 返回数据说明
- 异常返回说明

下面分别介绍各部分内容的来源和定义方式。

接口名和请求路径
----

接口名和请求路径都通过注册路由的方法参数传入：

```
todosRouter := server.Handler(
  "获取todo列表", // 接口名
  "/todos", // 请求路径
  TodosHandler,
)
```

接口说明
----

```
todosRouter.Comment(`
该接口根据请求的分页参数返回todo列表。
请求示例：
write your sample code here using markdown schema
`)
```

通过调用router的Comment方法，可以设置接口说明。在接口说明的文本中可以使用markdown语法。

请求参数、返回数据、异常返回说明
----

通过Handler可以指定请求参数、返回数据、异常返回的数据结构。

```
func Handler(in *Param, out *Response) *Error {}
```

对于每一个数据结构，文档都会将其转化为json格式的说明。

例如，下面的Param数据结构：
```
type Param struct {
  Limit int `json:"limit" valid:"(0,),message=必须大于0" comment:"每页多少条"`
  Offset int `json:"offset" valid:"optional,[0,)" comment:"第几页"`
}
```

在生成文档的时候，struct中的可导出属性都会被列出，以json tag的值作为key，value中会列出类型、valid tag和comment tag。

根据这个原则，上面Param所生成的文档为：

```
{
  "limit": "int |(0,),message=必须大于0| // 每页多少条",
  "offset": "int |optional,[0,)| // 第几页"
}
```

相关链接
----

- [阅读上一章：Context](/http/doc/context.md)
- [阅读下一章：日志](/http/doc/logger.md)
- [返回包简介](/http/README.md)
- [返回示例example讲解](/http/example/README.md)
