# 接口文档

## 添加todo

新增一条todo数据
```
curl -H'Authorization: Basic a2VscDprZWxw' -d'{"title":"kelp example","content":"write kelp example for every package","alert_time":"2018-10-26 14:00:00"}' http://host:port/todo/create
```


请求路径：` /todo/create `

请求参数：
```
{
  "alert_time": "string |optional,@datetime,message=alert_time不合法|",
  "content": "string |optional,(0,1024),message=content不合法|",
  "title": "string |(0,128],message=titil不合法|"
}
```

返回数据：
```
{
  "data": "interface // 默认为空字符串",
  "status": "int // 默认为0"
}
```

异常返回：
```
{
  "message": "interface // 用于联调测试时参考的错误信息",
  "status": "int // 请参考开发者定义的Status列表"
}
```


## 修改todo

修改一条todo数据
```
curl -H'Authorization: Basic a2VscDprZWxw' -d'{"id":"1","title":"kelp example","content":"write kelp example for every package","alert_time":"2018-10-26 14:00:00"}' http://host:port/todo/update
```


请求路径：` /todo/update `

请求参数：
```
{
  "alert_time": "string |optional,@datetime,message=alert_time不合法|",
  "content": "string |optional,(0,1024),message=content不合法|",
  "id": "int |[0,),message=id不合法|",
  "title": "string |(0,128],message=titil不合法|"
}
```

返回数据：
```
{
  "data": "interface // 默认为空字符串",
  "status": "int // 默认为0"
}
```

异常返回：
```
{
  "message": "interface // 用于联调测试时参考的错误信息",
  "status": "int // 请参考开发者定义的Status列表"
}
```


## 删除todo

删除一条todo数据
```
curl -H'Authorization: Basic a2VscDprZWxw' -d'{"id":"1"}' http://host:port/todo/delete
```


请求路径：` /todo/delete `

请求参数：
```
{
  "id": "int |[0,),message=id不合法|"
}
```

返回数据：
```
{
  "data": "interface // 默认为空字符串",
  "status": "int // 默认为0"
}
```

异常返回：
```
{
  "message": "interface // 用于联调测试时参考的错误信息",
  "status": "int // 请参考开发者定义的Status列表"
}
```


## 列表todo

查看所有todo数据
```
curl -H'Authorization: Basic a2VscDprZWxw' http://host:port/todo/retrieve
```


请求路径：` /todo/retrieve `


返回数据：
```
{
  "list": [
    {
      "alert_time": "string",
      "content": "string",
      "id": "int",
      "title": "string"
    }
  ]
}
```

异常返回：
```
{
  "message": "interface // 用于联调测试时参考的错误信息",
  "status": "int // 请参考开发者定义的Status列表"
}
```

