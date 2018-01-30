A golang web framework
====

We provide a http server, a http(s) client, and some other helpful component, such as:

- encode and crypto methods
- common handlers

How about a Hello world:
```
// new a go file main.go

// import this web package in your go env

// implement a main method
func main() {
    server := web.New(":9999")
    server.GET("/hello", func(c *web.Context) {
        c.Success("Hello World!")
    })
    server.Run()
}

// save file and run with command: go run main.go

// curl http://127.0.0.1:9999/hello
// you will recieve:
// {"status":0,"data":"Hello World!"}
```

Usage above, performance data followed.

Start
----

Run as a service.
```
server := web.New(":9999")
server.Run()

// server will listen on 127.0.0.1:9999
```

Router
----

Here are router and group.
```
// GET http://127.0.0.1:9999/get
server.GET("/get", func(c *web.Context){
    // ...
})

// GET http://127.0.0.1:9999/path/get
server.GET("/path/get", func(c *web.Context){
    // ...
})

// GET http://127.0.0.1:9999/group/get
group := server.Group("/group")
group.GET("/get", func(c *web.Context){
    // ...
})

```

Method & Param
----

Method and Path param can be used for RESTful.
```
// POST http://127.0.0.1:9999/post
server.POST("/post", func(c *web.Context){
    // ...
})

// PUT http://127.0.0.1:9999/put
server.PUT("/put", func(c *web.Context){
    // ...
})

// DELETE http://127.0.0.1:9999/delete
server.PUT("/delete", func(c *web.Context){
    // ...
})

// GET http://127.0.0.1:9999/param/123/ok
server.GET("/param/:id/:status", func(c *web.Context){
    // all param is string
    //  id == "123"
    //  status == "ok"
    //  other == "", exist == false
    id, _ := c.Param("id")
    status, _ := c.Param("status")
    other, exist := c.Param("other")

    // ...
})

```

Query
----

Get a query param if there exist.
```
// GET http://127.0.0.1:9999/query?id=123
server.GET("/query", func(c *web.Context){
    // all query is string
    //  id == "123"
    //  status == "ok"
    //  other == "", exist == false
    id, _ := c.Query("id")
    status, _ := c.QueryDefault("status", "ok")
    other, exist := c.Query("other")

    // ...
})
```

Bind
----

Bind body to a struct.
```
type ReqParam struct {
    Name string `json:"name"`
    Age int `json:"age"`
}

// request with json body:
//      {"name":"cookie","age":12}
func handler(c *web.Context) {
    req := &ReqParam{}
    c.Bind(req)
    name := req.Name // cookie
    age := req.Age // 12
    // ...
}
```

Response
----

Response with json body or http status.
```
c.Json(obj) // response obj json object
c.Success(data) // response {"status":0, "data":<data json object>}
c.Error(1,message) // reponse {"status":1, "message":<message json object>}
c.DieWithHttpStatus(404) // response a 404 http status
```

Handler Chain
----

Use handler chain as middleware.
```
auth := server.Group("/auth")
server.Use(sessionHandler)
auth.Get("/info", authHandler, infoHandler)

// GET /auth/info will call sessionHandler()
func sessionHandler(c *web.Context) {
    // do something
    c.Next() // call authHandler() here
    // do other thing
}

func authHandler(c *web.Context) {
    // do something
    return
    // because of not call infoHandler(), here will response directly
}

// ...
```

We provide some common handlers for your convenience. See [handlers.go](handlers.go)

Meta Data
----

Save meta data with context.
```
// in a handler
c.Set("username", "cookie")
c.Set("age", 12)

// in other handler
name := c.GetString("username")
age := c.GetInt("age")

undef, exist := c.Get("undefined")
// exist is false
iname, _ := c.Get("username")
// iname is interface{} which can be assert as string
```

Cookie
----

Use cookie as you wish.
```
token := c.GetCookie("token")
c.SetCookie("token", token, 24*time.Hour)
```

Session
----

We advise using session in handler chain.
```
// before start use session middleware
server.UseMemSession(2*time.Hour, 5*time.Minute)
// server.UseSession(YourSessionImplementiSessionPool)

// in handler chain
c.SetSession("session_key", value)
value := c.GetSession("session_key")
```

Client
----

```
body, err := web.PostJson("http://127.0.0.1/jsonapi", jsonstr)

body, err := web.PostForm("http://127.0.0.1/formapi", formValues)

body, err := web.Get("http://127.0.0.1/getapi?query=xxx")

err := Mail(
    "mapleque@163.com",
    "password",
    "smtp.163.com:465",
    "mapleque@163.com",
    "Mail Title",
    "Mail Content",
    )
```

Crypt
----

Log
----

Redirect log to a component witch implements logInterface. Otherwise default golang log package will be use.
```
web.SetLogger(log.Log)
```


