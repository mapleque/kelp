Kelp
----

A go libs

1. 依赖第三方mysql驱动
```
go get github.com/go-sql-driver/mysql
```
1. clone本项目到gopath/github.com/下
1. 参考下文实现自己的业务逻辑

## Queue
队列模块实现了一个非阻塞异步队列生产消费机制，对于一个队列，需要分别指定他的Producer和Consumer。
Producer和Consumer要实现对应的接口方法```Push```和```Pop```。
可以通过下面的方式初始化一个队列：
```
// import package queue

// a simple implement struct
type SimpleImpl struct{}

// implement Producer
func (p SimpleImpl) Push(q *queue.Queue, task string) {
	for i := 0; i < 5+rand.Intn(5); i++ {
		qItem := q.Push(task, i, "item data")
		log.Info("push", qItem)
	}
	time.Sleep(time.Duration(rand.Intn(2000))*time.Millisecond + 7*time.Second)
}

// implement Consumer
func (c SimpleImpl) Pop(q *queue.Queue, task string) {
	qItem := q.Pop()
	log.Info("pop", qItem)
	time.Sleep(time.Duration(rand.Intn(2000))*time.Millisecond + time.Second)
}

impl := SimpleImpl{}

// regist impl as producer because of it implement Producer interface
// regist impl as consumer because of it implement Consumer interface
// regist both producer and consumer at the same time
aQueue := queue.RegistTask("task_name", 10, impl, impl)
// here we use the task name as queue name
// nil producer or consumer will not be regist on queue in this method

// get a queue by name
targetQueue, ok := queue.GetQueue("task_name")
// aQueue and targetQueue is the same point

otherQueue := queue.CreateQueue("queue_name", 10)

// only regist producer
otherQueue.RegistProducer("producer_name", impl)

// only regist consumer
otherQueue.RegistConsumer("consumer_name", impl)

// here run all producer and consumer
// while running
// all producer.Push and consumer.Pop method will be call in loop
// if queue is empty, Pop will be block
go queue.Run()
```

## Crontab

后台任务模块提供了定时执行任务机制，任务周期表达式参考crontab的标准，任务执行者需要实现Crontab对应的接口方法```Triger```。
可以通过下面的方式初始化后台任务：
```
// a simple implement struct
type SimpleImpl struct{}

// implement Crontab
func (c SimpleImpl) Triger(task string) {
	// do nothing
}

impl := SimpleImpl{}

crontab.Regist("* * * * *", "crontab_name", impl)
go crontab.Run()
```

## Monitor

监控模块可以通过开启的指定端口监控系统状态。

可以被监控的模块都实现了```monitor.Observable```接口。
注册这些模块就可以通过监控接口获得运行时状态数据。
```
monitor.Observe("queue", queue.GetQueueContainer())
monitor.Observe("crontab", crontab.GetCrontabContainer())
monitor.Observe("producer", queue.GetProducerContainer())
monitor.Observe("consumer", queue.GetConsumerContainer())
go monitor.Run("127.0.0.1:9998")
```

## Web

web模块提供了一个http的服务框架，通过注册路由和对应的handler方法，可以实现http服务。
```
// implement a handler
func helloHandler(context *web.Context) {
	context.Data = "hello"
}

server := web.New("127.0.0.1:9000")
server.RegistHandler("/hello", helloHandler)
go server.Run()
```

## Config

配置文件加载模块读取配置。
```
;config_file.ini
[section]
KEY=xxx
```

```
// init a configer
config.Add(config.INI, "config_name", "config_file_path")

// read a configer
configer := config.Use("config_name")
some_config_value := configer.Get("section.KEY")
```

此外，如果在系统中只需要读取一个配置文件，还可以使用简化的方法
```
// init default configer
configer := config.InitDefault(config.INI, "ini", "./config_file.ini", "the config file path param --ini")

// get configer in other scope
configer := config.Default()
```

## Mysql

数据库模块支持mysql数据库操作。
使用数据库模块需要在系统启动时初始化数据库连接：
```
// init db connection
// ping db server when add db
// if ping failed, it will fatal
db.AddMysql("db_name",
    "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8",
    10,// 最大连接数
    10,// 最大闲置连接数
)

// db operation
// Select返回查询结果数组[]map[string]interface{}
// Insert返回插入数据id
// Update返回受影响行数
// Excute返回受影响行数

// 先选择database，再调用方法
query := db.UseMysql("db_name")
ret := query.Select("your sql", your_param...)

// 使用事物：
trans := db.Begin("db_name")
// 也可以先选择database，在开启事物
// query := db.UseMysql("db_name")
// trans := query.Begin()
trans.Update("your sql", your_param...)
ret := trans.Insert("your sql", your_param...)

// commit or rollback
if ret < 0 {
    trans.Rollback()
} else {
    trans.Commit()
}

// 获取链接自行操作
```
conn := db.UseMysql("db_name").GetConn() // return *sql.DB
transConn := db.Begin("db_name").GetConn() // return *sql.Tx
```

## Log

日志模块支持日志定制输出。
本项目中所有模块日志都支持重定向到logger输出，在应用中也可以单独使用logger。
```
log.AddLogger(
    "log_file_name.log",
    "log_file_path",
    10, // 日志文件数量，到达这个数量时，新切分的日志会覆盖最早的
    10000000, // 单个日志文件最大值，当文件超过这个值时，会被切分为log_file_name.log.n
    5, // 日志输出最高级别
    2, // 日志输出最低级别
)
// 只有当日志级别满足对应logger时，对应的日志才会被输出到文件中

// 日志输出方法调用
log.Error("your error massage", "other additional object")
log.Warn("some massage")
log.Info("some massage")
log.Debug("some massage")
// 其中 Error, Warn, Debug会输出callstack
// callstack也可以被直接调用
log.Callstack(log.INFO)

// 也可以选择指定的日志输出
log.Get("log_file_name.log").Info("some massage")

```

## Util

工具模块提供一些常用工具。

### File

读文件方法：读取指定文件返回包含了每行数据的数组。

### Type

提供了一些强制类型转换的方法。
=======
- [web](./web) Http(s) 服务组件，包含服务端和客户端方法实现
- [db](./db) 数据库封装组件
- [log](./log) 日志组件，支持分级输出和文件轮转
- [config](./config) 配置加载组件封装
