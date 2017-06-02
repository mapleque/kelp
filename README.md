# Kelp

## 部署

1. 依赖第三方mysql驱动
```
go get github.com/go-sql-driver/mysql
```
1. clone本项目到gopath/github.com/下
1. 参考example实现自己的业务逻辑

## 设计
```
                                    queue
                                     | |
                prducer              | |
 create task --->| push data         | |
                 +------------------>| |
                                     | |        consumer
                                     | | pop data   |
                                     | |----------->|
                                     | |            +---> do something
                                     | |


                 regist crontab task
                          | |
ticker event second ----->| | triger
                          | |--------> run task
```

### Queue
队列模块实现了一个非阻塞异步队列生产消费机制，对于一个队列，需要分别指定他的Producer和Consumer。
Producer和Consumer要实现对应的接口方法```Push```和```Pop```。参考```example/main.go```。
可以通过下面的方式初始化一个队列：
```
q := queue.RegistTask("task_name", 10, producerImpl, consumerImpl)
q.RegistProducer("producer_name", otherProducerImpl)
q.RegistConsumer("consumer_name", otherConsumerImpl)

otherQueue := queue.CreateQueue("queue_name", 10)

targetQueue, ok := queue.GetQueue("target_queue_name")

go q.Run()
```
对于RegistTask方法，如果producer或者consumer为nil的时候，对应的生产者或者消费者不会被注册。
该方法会返回创建的queue对象指针，这样就可以通过其提供的方法单独注册生产者和消费者。

当然，也可以通过CreateQueue方法直接创建queue对象。

对于已经创建的queue对象，可以通过GetQueue方法获取对应指针。

### Crontab

后台任务模块提供了定时执行任务机制，任务周期表达式参考crontab的标准，任务执行者需要实现Crontab对应的接口方法```Triger```。参考```example/main.go```。
可以通过下面的方式初始化后台任务：
```
crontab.Regist("* * * * *", "crontab_name", crontabImpl)
go crontab.Run()
```

### Monitor

监控模块可以通过开启的指定端口监控系统状态。

在系统启动时加入下面方法调用，即可启动monitor：
```
monitor.Run("127.0.0.1:9998")
```
monitor目前提供了下面功能，通过http协议访问即可：
#### 查看队列

```
request url : /queue
response : {}
```

#### 查看生产者

```
request url : /producer
response : {}
```

#### 查看消费者

```
request url : /consumer
response : {}
```

#### 查看后台任务

```
request url : /crontab
response : {}
```

至此，你已经可以开始使用本框架工作了，下面提供了一些扩展模块的封装，请按需使用。

## Config

配置文件加载模块支持读取配置文件（目前只支持ini格式）。参考```example/main.go```。

```
config.AddConfiger(config.INI, "config_name", "config_file_path")
```

使用时通过config_name获取对应实例：
```
conf := config.Use("config_name")

some_config_value := conf.Get("section.KEY")
```

## Database

数据库模块支持mysql数据库操作。参考```example/main.go```。
使用数据库模块需要在系统启动时初始化数据库连接：
```
db.AddDB("db_name",
    "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8",
    10,
    10)

```
业务逻辑中可以有以下几种使用方式：
```
// 直接调用方法Select,Update,Insert
ret := db.Select("db_name", "your sql", your_param...)

// 先选择database，再调用方法
query := db.UseDB("db_name")
ret := query.Select("your sql", your_param...)

```
事物使用方法：
```
trans := db.Begin("db_name")
// 也支持先选择database
// query := db.UseDB("db_name")
// trans := query.Begin()
trans.Update("your sql", your_param...)
ret := trans.Update("your sql", your_param...)

if ret < 0 {
    trans.Rollback()
} else {
    trans.Commit()
}
```
数据库模块支持添加多个数据库连接，在操作数据库的时候选择不同的db_name即可。

## Log

日志模块支持日志定制输出。参考```example/main.go```。

## Util

工具模块提供一些常用工具。

### File

读文件方法：读取指定文件返回包含了每行数据的数组。

### Type

提供了一些强制类型转换的方法。
