参数校验
====

kelp/http包提供了通过struct tag定义参数校验的方法。

```
type Param struct {
  Data string `json:"data" valid:"(0,),message=data length must not be 0"`
}
param := &Param{}
err := http.BindAndValidJson(param, `{"data":""}`)
fmt.Println(err)
// output:
// data length must not be 0
```

如上面例子中所写，在struct的字段tag中，如果定义了valid标签，那么BindAndValidJson方法将会对该字段进行校验。

valid标签
----

alid标签由三部分组成：

- expression 校验表达式
- optional 是否是可选参数
- message 校验失败的提示信息

每一部分内容都是可选的，因此valid的形式可能会有下面几种情况：

```
valid:"" // 该参数为必选参数，可以传空值，当该参数不存在的时候会返回一个默认的错误提示
valid:"expression" // 该参数为必选参数，且满足表达式expression，当该参数不存在或不满足表达式的时候会返回一个默认的错误提示
valid:"optional" // 该参数为可选参数，不存在返回错误的情况
valid:"message" // 该参数为必选参数，可以传空值，当该参数不存在的时候会返回message
valid:"expression,optional" // 该参数为可选参数，如果传了该参数则必须满足表达式expression，当传了该参数且不满足表达式expression的时候会返回一个默认的错误提示
valid:"expression,message" // 该参数为必选参数，且满足表达式expression，当该参数不存在或不满足表达式的时候会返回message
valid:"optional,message" // 该参数为可选参数，不存在返回错误的情况，所以message无意义
valid:"expression,optional,message" // 该参数为可选参数，如果传了该参数则必须满足表达式expression，当传了该参数且不满足表达式expression的时候会返回message
```

expression表达式
----

目前，kelp/http包支持如下几种表达式模式：
- 正则表达式模式
- 范围表达式模式
- 函数模式


### 正则表达式模式

```
/regexp/
```

正则表达式模式首尾字符必须为`/`，中间写正则表达式。

它的处理逻辑是：在参数校验的时候，不论参数类型都会被当作字符串与当前正则表达式进行正则表达式匹配，如果匹配失败，则校验不通过。

对于参数值的匹配，通常会有如下几种情况：
```
{"data":"str"} // 匹配`str`
{"data": 1} // 匹配`1`
{"data": 1.2} // 匹配`1.2`
{"data": true} // 匹配`true`
```

特别的：部分正则表达式符号是不允许在struct tag中出现的，对于这种情况，建议使用函数模式进行校验。
可以使用ValidRegexpWrapper方法来快速定义校验函数。

### 范围模式
```
[m,n]
(m,n)
[m,n)
(m,n]
(m,)
[m,)
(,n)
(,n]
```

范围模式以`(`或`[`开头，`)`或`]`结尾，其中`(`和`)`表示开区间，`[`和`]`表示闭区间。

中间`m`和`n`是两个非负整数，`m`表示最小值，`n`表示最大值，可以没有其中任何一个表示不限最小(大)值。

范围模式仅对一下数据类型有效：
- int(包括int,int8,int32,int64) 比较值的大小
- float(包括float32,float64) 比较值的大小
- string 比较字符串长度

对于其他类型数据，一律校验不通过。

### 函数模式
```
@funcName
```

自定义的ValidFunc，只要进行注册，即可在valid tag中使用。

在注册的函数名前面使用`@`表示函数模式。

例如：
```
func MyValidFunc(
	fieldType reflect.StructField,
	destSource []byte,
	root reflect.Value,
	rootSource map[string]json.RawMessage,
) bool {
  // TODO valid destSource
  return true
}

// ...
http.RegisterValidFunc("myValidFunc", MyValidFunc)

// ...
type Param struct {
  Data string `json:"data" valid:"@myValidFunc"`
}

```


相关链接
----

- [阅读上一章：Handler & 中间件](/http/doc/handler.md)
- [阅读下一章：Context](/http/doc/context.md)
- [返回包简介](/http/README.md)
- [返回示例example讲解](/http/example/README.md)
