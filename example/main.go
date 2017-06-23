package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/kelp/config"
	"github.com/kelp/crontab"
	"github.com/kelp/db"
	"github.com/kelp/log"
	"github.com/kelp/monitor"
	"github.com/kelp/queue"
	"github.com/kelp/web"
)

// a simple implement struct
// which implement Producer, Consumer, Crontab
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

// implement Crontab
func (c SimpleImpl) Triger(task string) {
	// do nothing
}

// implement a handler
func helloHandler(context *web.Context) {
	context.Data = "hello"
}

func initLog() {
	conf := config.Use("config_name")
	log.AddLogger(
		conf.Get("log.NAME"),
		conf.Get("log.PATH"),
		conf.Int("log.MAX_NUMBER"),
		conf.Int64("log.MAX_SIZE"),
		conf.Int("log.MAX_LEVEL"),
		conf.Int("log.MIN_LEVEL"))
	db.SetLogger(log.Log)
	queue.SetLogger(log.Log)
	crontab.SetLogger(log.Log)
}

func initDB() {
	conf := config.Use("config_name")
	db.AddDB(
		"kelp",
		conf.Get("db.DSN"),
		conf.Int("db.MAX_CONNECTION"),
		conf.Int("db.MAX_IDLE"))
}

func initConfig() {
	confFile := flag.String("ini", "./config.ini", "your config file")
	flag.Parse()
	if *confFile == "" {
		panic("run with -h to find usage")
	}
	config.AddConfiger(config.INI, "config_name", *confFile)
}

func main() {
	initConfig()
	initLog()
	initDB()
	conf := config.Use("config_name")

	// regists task
	impl := SimpleImpl{}
	queue.RegistTask("simple1", 1, impl, impl)
	queue.RegistTask("simple10", 10, impl, impl)
	queue.RegistTask("simple100", 10000, impl, impl)
	crontab.Regist("* * * * *", "simple", impl)

	// regist monitor
	monitor.Observe("queue", queue.GetQueueContainer())
	monitor.Observe("crontab", crontab.GetCrontabContainer())
	monitor.Observe("producer", queue.GetProducerContainer())
	monitor.Observe("consumer", queue.GetConsumerContainer())

	// regist web
	server := web.New(conf.Get("web.HOST"))
	server.RegistHandler("/hello", helloHandler)

	// start service
	log.Info("start service ...")
	done := make(chan bool, 1)

	go queue.Run()
	go crontab.Run()
	go monitor.Run(conf.Get("monitor.HOST"))
	go server.Run()

	<-done
}
