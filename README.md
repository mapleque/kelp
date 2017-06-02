# Kelp & 后台任务处理调度框架

## 功能
将业务数据，经过一定的清洗和重组，同步到数据分析数据库。

## 部署

1. 依赖第三方mysql驱动
```
go get github.com/go-sql-driver/mysql
```
1. clone本项目到gopath/github.com/下
1. 拷贝config.ini.example为config.ini并修改
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

## Queue
队列模块实现了一个非阻塞异步队列生产消费机制，对于一个队列，需要分别指定他的Producer和Consumer。
Producer和Consumer只需要实现对应的接口方法即可。参考```example/main.go```。

## Crontab

后台任务模块提供了定时执行任务机制，任务周期表达式参考crontab的标准，任务执行者需要实现Crontab对应的接口方法。参考```example/main.go```。

## Database

数据库模块支持mysql数据库操作。参考```example/main.go```。

## Log

日志模块支持日志定制输出。参考```example/main.go```。

## Monitor

监控模块可以通过开启的指定端口监控系统状态。

### 查看队列

```
request url : /queue
response : {}
```

### 查看生产者

```
request url : /producer
response : {}
```

### 查看消费者

```
request url : /consumer
response : {}
```

### 查看后台任务

```
request url : /crontab
response : {}
```

## Util

工具模块提供一些常用工具。

### File

读文件方法：读取指定文件返回包含了每行数据的数组。

### Type

提供了一些强制类型转换的方法。
